package message

import (
	"EasyBanner/model"
	"EasyBanner/pkgs/data"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"io"
	"log"
	"os"
)

// 读取模板文件
func readTemplateFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	byteValue, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(byteValue), nil
}

// 发送交互消息
func SendInteractiveMsg(client *lark.Client, chatID string, messageID string) error {
	flag := data.CheckCount()

	var templatePath string
	var cardContext string
	var err error

	// 判断 ip 次数，发送不同的模板
	switch flag {
	case 0:
		templatePath = "/app/EasyBanner/templates/failure.json"
		// 读取 JSON 模板文件
		cardContext, err = readTemplateFile(templatePath)
		if err != nil {
			log.Println("Failed to read JSON file:", err)
			return err
		}
	case 1:
		// 生成包含 IP 数据的 JSON 模板
		cardContext, err = data.GetNeedBanIP()
		if err != nil {
			return err
		}
		if cardContext == "" { // 如果没有需要 ban 的 IP，返回错误
			log.Println("No IPs to ban")
			return nil
		}
	case 2:
		templatePath = "/app/EasyBanner/templates/common.json"
		// 读取 JSON 模板文件
		cardContext, err = readTemplateFile(templatePath)
		if err != nil {
			log.Println("Failed to read JSON file:", err)
			return err
		}
	}

	// 构建消息请求
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeChatId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeInteractive).
			ReceiveId(chatID).
			Content(cardContext).
			Build()).
		Build()

	// 发送消息
	resp, err := client.Im.Message.Create(context.Background(), req)
	if err != nil {
		log.Println("Failed to send message:", err)
		return err
	}

	if !resp.Success() {
		log.Println("Failed to send message!", resp)
		return err
	}

	return nil
}

// 接收整个消息体
func GetMessageBody(c *gin.Context) (*model.MessageBody, error) {
	// 读取消息体内容
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("error reading message body:", err)
		return nil, err
	}

	if len(body) == 0 {
		log.Println("Received empty request body")
		return nil, errors.New("empty request body")
	}

	log.Println("Request body:", string(body))

	var messageBody model.MessageBody
	if err := json.Unmarshal(body, &messageBody); err != nil {
		log.Println("error parsing json:", err)
		return nil, err
	}

	return &messageBody, nil
}

// 接收整个卡片回调消息体
func GetCardCallbackBody(c *gin.Context) (*model.CardCallback, error) {
	// 读取消息体内容
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("error reading message body:", err)
		return nil, err
	}

	if len(body) == 0 {
		log.Println("Received empty request body")
		return nil, errors.New("empty request body")
	}

	log.Println("Request body:", string(body))

	var cardCallback model.CardCallback
	if err := json.Unmarshal(body, &cardCallback); err != nil {
		log.Println("error parsing json:", err)
		return nil, err
	}

	return &cardCallback, nil
}

// 更新消息卡片内容
func UpdateInteractiveMsg(client *lark.Client, messageID string, ipDataList []model.IPData, resultText string) error {
	// 生成新的卡片内容
	cardContent, err := data.GenerateTemplate(ipDataList, false, resultText)
	if err != nil {
		return err
	}

	// 创建请求对象
	req := larkim.NewPatchMessageReqBuilder().
		MessageId(messageID).
		Body(larkim.NewPatchMessageReqBodyBuilder().
			Content(cardContent).
			Build()).
		Build()

	// 发起请求
	resp, err := client.Im.Message.Patch(context.Background(), req)
	if err != nil {
		return err
	}

	if !resp.Success() {
		return fmt.Errorf("failed to update message: %v", resp)
	}

	log.Println("卡片更新成功！")
	return nil
}

// 判断@事件
func CheckAtBot(messageBody *model.MessageBody) bool {
	// 检测是否有@事件
	if messageBody == nil {
		log.Println("MessageBody or Event or Message is nil!")
		return false
	}

	if len(messageBody.Event.Message.Mentions) > 0 {
		log.Println("detected @ event!")

		// 检测是否为 @机器人
		for _, mention := range messageBody.Event.Message.Mentions {
			if mention.ID.UserID == "" {
				// 检测到是 @机器人 的情况
				log.Println("Bot is mentioned!")
				return true
			}
		}

		// 如果循环结束都没有检测到 @机器人
		log.Println("This message is not for bot!")
		return false
	}

	log.Println("no detected @ event!")
	return false
}
