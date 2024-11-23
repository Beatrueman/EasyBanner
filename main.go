package main

import (
	"EasyBanner/pkgs/message"
	"EasyBanner/utils/base"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	base.InitConfig()

	r := gin.Default()

	r.POST("/webhook", message.HandleWebhook)
	r.POST("/event", message.HandleCardCallback)
	r.POST("/event/alert", message.HandleAlert)
	// 启动服务
	r.Run(":8080")
}
