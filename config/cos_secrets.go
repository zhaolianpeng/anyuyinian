package config

import (
	"os"
)

// GetCOSSecrets 从环境变量获取COS密钥
func GetCOSSecrets() (string, string) {
	secretID := os.Getenv("COS_SECRET_ID")
	secretKey := os.Getenv("COS_SECRET_KEY")

	// 如果环境变量未设置，使用默认值（仅用于开发环境）
	if secretID == "" {
		secretID = "your_secret_id_here"
	}
	if secretKey == "" {
		secretKey = "your_secret_key_here"
	}

	return secretID, secretKey
}
