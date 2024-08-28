package message

import (
	"EasyBanner/model"
	"EasyBanner/pkgs/data"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
)

// 解析飞书 challenge 测试
func HandleLarkCallBack(data []byte) (string, error) {
	// 解析请求体
	var challengeRequest model.ChallengeRequest
	if err := json.Unmarshal(data, &challengeRequest); err == nil && challengeRequest.Challenge != "" {
		// 如果请求体包含 challenge 字段，则认为是验证请求
		return challengeRequest.Challenge, nil
	}
	return "", nil
}

// 处理请求，确保body可以存值
// 由于 HTTP 请求的 body 是一个 io.ReadCloser，在第一次读取后，它就会被耗尽。如果 HandleLarkEvent 在 GetMessageBody 之前被调用，它可能会导致 c.Request.Body 在 GetMessageBody 调用时为空。
func HandleWebhook(c *gin.Context) {
	// 先读取请求体内容
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("error reading message body:", err)
		return
	}

	// 检查 Challenge
	challenge, err := HandleLarkCallBack(body)
	if err == nil && challenge != "" {
		c.JSON(http.StatusOK, gin.H{"challenge": challenge})
		return
	}

	// 再次设置请求体，以便 GetMessageBody 读取
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// 立即返回 HTTP 200 响应
	c.JSON(http.StatusOK, gin.H{"status": "ok"})

	go func() {

		// 调用 GetMessageBody
		messageBody, err := GetMessageBody(c)
		if err != nil {
			log.Println("error get message body:", err)
			return
		}

		// 只处理消息接收事件
		if messageBody.Header.EventType == "im.message.receive_v1" {
			// 判断是否为 @机器人事件
			flag := CheckAtBot(messageBody)

			if flag {

				client := lark.NewClient(viper.GetString("AppID"), viper.GetString("AppSecret"))

				err = SendInteractiveMsg(client, messageBody.Event.Message.ChatID, messageBody.Event.Message.MessageID)
				if err != nil {
					return
				}

				return
			} else {
				log.Println("Not an @bot event.")
				return
			}
		} else {
			// 忽略其他事件类型
			log.Println("Event type not handled:", messageBody.Header.EventType)
			return
		}
	}()
}

// 处理卡片回调
func HandleCardCallback(c *gin.Context) {

	// 先读取请求体内容
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("error reading message body:", err)
		return
	}

	// 检查 Challenge
	challenge, err := HandleLarkCallBack(body)
	if err == nil && challenge != "" {
		c.JSON(http.StatusOK, gin.H{"challenge": challenge})
		return
	}

	// 再次设置请求体，以便 GetMessageBody 读取
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// 立即返回 HTTP 200 响应
	c.JSON(http.StatusOK, gin.H{"status": "ok"})

	go func() {
		// 调用 GetCardCallbackBody
		cardCallback, err := GetCardCallbackBody(c)
		if err != nil {
			log.Println("error getting card callback body:", err)
			return
		}

		// 处理回调
		if cardCallback.Header.EventType == "card.action.trigger" {
			action := cardCallback.Event.Action

			// 处理特定操作类型
			if action.Tag == "button" && action.Value["action"] == "ban_ip" {
				// 获取需要 ban 的 IP 列表
				ipDataJSON, err := data.GetNeedBanIPList()
				if err != nil {
					log.Println("Error getting IPs!")
					return
				}

				// 将生成的 JSON 发送到 ban API
				success, err := data.ExecuteBanIP(ipDataJSON)
				if err != nil {
					log.Println("Failed to ban IPs!")
				}

				// 更新卡片内容
				client := lark.NewClient(viper.GetString("AppID"), viper.GetString("AppSecret"))
				messageID := cardCallback.Event.Context.OpenMessageID // 获取消息 ID
				ipDataList := data.ParseIPDataJSON(ipDataJSON)        // 从 JSON 获取 IP 数据

				// 根据封禁操作的结果更新卡片
				resultText := ""
				if success {
					resultText = "✅ 封禁 IP 操作成功!"
				} else {
					resultText = "🔴 封禁 IP 操作失败，请稍后重试。"
				}

				updateErr := UpdateInteractiveMsg(client, messageID, ipDataList, resultText)
				if updateErr != nil {
					log.Println("Error updating card:", updateErr)
				} else {
					log.Println("Card updated successfully!")
				}
			} else {
				log.Println("Unknown action!")
			}
		} else {
			log.Println("Event type not handled!")
		}
	}()
}
