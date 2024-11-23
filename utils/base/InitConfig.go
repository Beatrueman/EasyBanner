package base

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

func InitConfig() {
	// 读取配置
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	// 配置文件所在目录
	viper.AddConfigPath("EasyBanner")

	// 从环境变量中读取配置
	viper.SetDefault("AppID", os.Getenv("APP_ID"))
	viper.SetDefault("AppSecret", os.Getenv("APP_SECRET"))
	viper.SetDefault("URL", os.Getenv("URL"))

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("error reading config: %v", err)
		return
	}
}
