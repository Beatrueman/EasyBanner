package base

import (
	"EasyBanner/model"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

func GetTenantAccessToken() string {
	url := "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"

	InitConfig()
	AppID := viper.GetString("AppID")
	AppSecret := viper.GetString("AppSecret")

	// 校验是否已配置环境变量
	if AppID == "" || AppSecret == "" {
		log.Println("Error: AppID or AppSecret is not set")
		return ""
	}

	// 构造请求参数
	params := map[string]string{
		"app_id":     AppID,
		"app_secret": AppSecret,
	}
	jsonData, err := json.Marshal(params)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return ""
	}

	// 发起 POST 请求
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error making POST request:", err)
		return ""
	}
	defer resp.Body.Close()

	// 处理响应
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to get token, status code: %d\n", resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("Response:", string(body))
		return ""
	}

	// 解析响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return ""
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Println("Error unmarshalling response JSON:", err)
		return ""
	}

	// 输出 token 信息
	if token, ok := result["tenant_access_token"].(string); ok {
		//log.Println("Tenant Access Token:", token)
		return token
	} else {
		log.Println("Failed to retrieve tenant_access_token:", result)
		return ""
	}
}

// 获取机器人所在群聊 chat_id
func GetChatID() string {
	url := "https://open.feishu.cn/open-apis/im/v1/chats"
	accessToken := GetTenantAccessToken()
	if accessToken == "" {
		log.Println("Failed to get tenant access token")
		return ""
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Failed to create GET request:", err)
		return ""
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error making GET request:", err)
		return ""
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Failed to get chat_id, status code: %d, response: %s\n", resp.StatusCode, string(body))
		return ""
	}

	// 解析响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return ""
	}

	// 处理 interface{}动态结构化数据
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Println("Error unmarshalling response JSON:", err)
		return ""
	}

	// 校验返回 code 和数据结构
	code, ok := result["code"].(float64)
	if !ok || code != 0 {
		log.Printf("Unexpected response code: %v, result: %v\n", code, result)
		return ""
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		log.Println("No 'data' field in response")
		return ""
	}

	items, ok := data["items"].([]interface{})
	if !ok || len(items) == 0 {
		log.Println("No chat_id found in 'items'")
		return ""
	}

	// 遍历 items 提取 chat_id
	for _, item := range items {
		chat, ok := item.(map[string]interface{})
		if ok {
			if chatID, exists := chat["chat_id"].(string); exists {
				//log.Println("Retrieved chat_id:", chatID)
				return chatID
			}
		}
	}

	// 未找到 chat_id 的情况
	log.Println("Failed to retrieve chat_id from items")
	return ""
}

// 获取来自EasyBanner的最新一个卡片的message_id
func GetLatestMessageID(AppID string) (string, error) {
	// 设置请求 URL，并添加 container_id 和 container_id_type 参数
	url := fmt.Sprintf("https://open.feishu.cn/open-apis/im/v1/messages?container_id=%s&container_id_type=chat", GetChatID())

	token := GetTenantAccessToken()

	// 设置请求头
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("failed to create request: %v", err)
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	// 发起 HTTP 请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to send request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read response body: %v", err)
		return "", err
	}

	// 解析 JSON 数据
	var response model.MessageListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Println("failed to parse response: %v", err)
		return "", err
	}

	// 筛选指定 senderID 的消息
	var filteredMessages []model.Message
	for _, message := range response.Data.Items {
		if message.Sender.ID == AppID {
			filteredMessages = append(filteredMessages, message)
		}
	}

	// 按 create_time 排序，获取最新消息
	if len(filteredMessages) == 0 {
		log.Printf("no messages found for sender ID: %s", AppID)
		return "", nil
	}
	sort.Slice(filteredMessages, func(i, j int) bool {
		return filteredMessages[i].CreateTime > filteredMessages[j].CreateTime
	})

	// 返回最新消息的 message_id
	log.Println("Latest message ID:", filteredMessages[0].MessageID)
	return filteredMessages[0].MessageID, nil
}
