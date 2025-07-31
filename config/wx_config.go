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
		AppID:     getEnv("WX_APP_ID", "wx101090677bd5219e"),
		AppSecret: getEnv("WX_APP_SECRET", "042ff9921818ada9336df6e91fc2287e"),
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
