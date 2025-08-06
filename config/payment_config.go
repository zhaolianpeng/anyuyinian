package config

import (
	"os"
)

// PaymentConfig 支付配置
type PaymentConfig struct {
	WechatPay WechatPayConfig `json:"wechatPay"`
}

// WechatPayConfig 微信支付配置
type WechatPayConfig struct {
	AppID       string `json:"appId"`       // 小程序AppID
	MchID       string `json:"mchId"`       // 商户号
	MchKey      string `json:"mchKey"`      // 商户密钥
	NotifyURL   string `json:"notifyUrl"`   // 支付结果通知地址
	CertPath    string `json:"certPath"`    // 证书路径
	KeyPath     string `json:"keyPath"`     // 私钥路径
	Environment string `json:"environment"` // 环境：sandbox或production
}

// GetPaymentConfig 获取支付配置
func GetPaymentConfig() *PaymentConfig {
	return &PaymentConfig{
		WechatPay: WechatPayConfig{
			AppID:       getPaymentEnv("WECHAT_PAY_APP_ID", "wx101090677bd5219e"),
			MchID:       getPaymentEnv("WECHAT_PAY_MCH_ID", ""),
			MchKey:      getPaymentEnv("WECHAT_PAY_MCH_KEY", ""),
			NotifyURL:   getPaymentEnv("WECHAT_PAY_NOTIFY_URL", "https://your-domain.com/api/payment/notify"),
			CertPath:    getPaymentEnv("WECHAT_PAY_CERT_PATH", ""),
			KeyPath:     getPaymentEnv("WECHAT_PAY_KEY_PATH", ""),
			Environment: getPaymentEnv("WECHAT_PAY_ENVIRONMENT", "sandbox"),
		},
	}
}

// getPaymentEnv 获取支付环境变量，如果不存在则返回默认值
func getPaymentEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
