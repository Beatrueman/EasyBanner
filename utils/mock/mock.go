package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	url := "http://localhost:8080/event/alert"
	// 读取 JSON 文件
	jsonFile, err := ioutil.ReadFile("EasyBanner/templates/alert.json")
	if err != nil {
		log.Fatalf("Failed to read JSON file: %v", err)
	}

	// 发送 POST 请求
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonFile))
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 打印响应
	log.Printf("Response status: %s\n", resp.Status)
}
