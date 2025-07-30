package config

import (
	"net/http"
	"net/url"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// COSConfig 腾讯云COS配置
type COSConfig struct {
	SecretID  string
	SecretKey string
	Region    string
	Bucket    string
	Domain    string
	ACL       string // 访问控制列表
}

// GetCOSConfig 获取COS配置
func GetCOSConfig() *COSConfig {
	secretID, secretKey := GetCOSSecrets()
	return &COSConfig{
		SecretID:  secretID,
		SecretKey: secretKey,
		Region:    "ap-shanghai", // 上海地域
		Bucket:    "7072-prod-5g94mx7a3d07e78c-1353115175",
		Domain:    "https://7072-prod-5g94mx7a3d07e78c-1353115175.cos.ap-shanghai.myqcloud.com",
		ACL:       "public-read", // 所有用户可读，仅创建者可读写
	}
}

// GetCOSClient 获取COS客户端
func GetCOSClient() *cos.Client {
	config := GetCOSConfig()

	// 将 https://xxx.cos.ap-shanghai.myqcloud.com 格式的域名转换为 https://cos.ap-shanghai.myqcloud.com
	u, _ := url.Parse(config.Domain)

	b := &cos.BaseURL{BucketURL: u}

	// 创建客户端
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.SecretID,
			SecretKey: config.SecretKey,
		},
		Timeout: 30 * time.Second,
	})

	return client
}
