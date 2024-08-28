package data

import (
	"bytes"
	"github.com/goccy/go-json"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

// 将请求体发送至 ban 的 API
func ExecuteBanIP(ipDataJSON string) (bool, error) {
	URL := viper.GetString("URL")
	banAPI := URL + "/ban"

	// 创建请求对象
	req, err := http.NewRequest("POST", banAPI, bytes.NewBuffer([]byte(ipDataJSON)))
	if err != nil {
		log.Println("error creating request:", err)
		return false, err
	}

	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("error sending request:", err)
		return false, err
	}
	defer resp.Body.Close()

	var results []struct {
		IP     string `json:"ip"`
		Status string `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return false, err
	}

	// 检查是否所有 IP 都成功处理
	for _, result := range results {
		log.Printf("IP: %s, Status: %s", result.IP, result.Status)
		if result.Status != "success" {
			return false, nil
		}
	}

	log.Println("Successfully sent ban request to API")

	return true, nil
}
