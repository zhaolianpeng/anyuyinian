package config

import (
	"os"
)

// WxConfig 微信小程序配置
type WxConfig struct {
	AppID     string
	AppSecret string
}

// GetWxConfig 获取微信配置
func GetWxConfig() *WxConfig {
	return &WxConfig{
		AppID:     getEnv("WX_APP_ID", "your_app_id"),
		AppSecret: getEnv("WX_APP_SECRET", "your_app_secret"),
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
