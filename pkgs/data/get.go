package data

import (
	"EasyBanner/model"
	"encoding/json"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
)

type IPResult []model.IPData

func GetIP() IPResult {
	URL := viper.GetString("URL")

	url := URL + "/execute"
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Failed to send request:", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response body:", err)
		return nil
	}

	var ipResult IPResult
	if err := json.Unmarshal(body, &ipResult); err != nil {
		log.Println("Failed to parse json response:", err)
		return nil
	}

	return ipResult
}

// 检查第一个 IP 的 count 是否大于 250
func CheckCount() int {
	// 获取IP和Count
	ipResult := GetIP()

	// 如果结果为空也返回false
	if len(ipResult) == 0 {
		log.Println("result is null!")
		return 0
	}

	// 检查第一个IP的count是否大于250
	if ipResult[0].Count >= 250 {
		return 1
	}

	return 2
}

// 获取所有次数大于250次的IP以及对应的次数填充至 JSON 模板
func GetNeedBanIP() (string, error) {
	result := GetIP()

	// 列表保存需要禁止的IP
	var needBanIPs []model.IPData

	for _, ipData := range result {
		if ipData.Count >= 250 {
			needBanIPs = append(needBanIPs, ipData)
		}
	}

	// 检查是否有需要禁止的IP
	if len(needBanIPs) == 0 {
		log.Println("没有大于250次的IP！")
		return "", nil
	}

	// 输出需要 ban 的 IP 和 次数
	log.Println("需要禁止的IP以及对应次数：")
	for _, ipData := range needBanIPs {
		log.Printf("IP: %s, count: %d", ipData.IP, ipData.Count)
	}

	// 生成包含 IP 数据的 JSON 模板
	finalJSON, err := GenerateTemplate(needBanIPs, true, "")
	if err != nil {
		log.Println("生成模板失败:", err)
		return "", nil
	}

	return finalJSON, nil
}

// 获取所有次数大于250次的IP以及对应的次数填充至 JSON 模板，无按键版
func GetNeedBanIPNoButton() (string, error) {
	result := GetIP()

	// 列表保存需要禁止的IP
	var needBanIPs []model.IPData

	for _, ipData := range result {
		if ipData.Count >= 250 {
			needBanIPs = append(needBanIPs, ipData)
		}
	}

	// 检查是否有需要禁止的IP
	if len(needBanIPs) == 0 {
		log.Println("没有大于250次的IP！")
		return "", nil
	}

	// 输出需要 ban 的 IP 和 次数
	log.Println("需要禁止的IP以及对应次数：")
	for _, ipData := range needBanIPs {
		log.Printf("IP: %s, count: %d", ipData.IP, ipData.Count)
	}

	// 生成包含 IP 数据的 JSON 模板
	finalJSON, err := GenerateNoButtonTemplate(needBanIPs, true, "")
	if err != nil {
		log.Println("生成模板失败:", err)
		return "", nil
	}

	return finalJSON, nil
}

// 获取所有次数大于250次的IP以及对应的次数
func GetNeedBanIPList() (string, error) {
	result := GetIP()

	// 列表保存需要禁止的IP
	var needBanIPs []model.IPData

	for _, ipData := range result {
		if ipData.Count >= 250 {
			needBanIPs = append(needBanIPs, ipData)
		}
	}

	// 检查是否有需要禁止的IP
	if len(needBanIPs) == 0 {
		log.Println("没有大于250次的IP！")
		return "[]", nil
	}

	// 将 IP 列表转为 JSON
	ipDataJSON, err := json.Marshal(needBanIPs)
	if err != nil {
		return "", nil
	}

	return string(ipDataJSON), nil
}

// 从 JSON 解析 IP 数据
func ParseIPDataJSON(ipDataJSON string) []model.IPData {
	var ipDataList []model.IPData
	json.Unmarshal([]byte(ipDataJSON), &ipDataList)
	return ipDataList
}
