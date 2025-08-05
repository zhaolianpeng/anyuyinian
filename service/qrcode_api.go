package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"wxcloudrun-golang/utils"
)

// QRCodeResponse 二维码响应
type QRCodeResponse struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

// GenerateQRCodeHandler 生成二维码接口
func GenerateQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理生成二维码请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodGet {
		LogError("请求方法不支持", fmt.Errorf("期望GET方法，实际为%s", r.Method))
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	// 获取推广码参数
	promoterCode := r.URL.Query().Get("promoterCode")
	if promoterCode == "" {
		http.Error(w, "缺少promoterCode参数", http.StatusBadRequest)
		return
	}

	// 验证推广码格式
	if !utils.ValidatePromoterCode(promoterCode) {
		response := &QRCodeResponse{
			Code:     -1,
			ErrorMsg: "推广码格式无效",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 生成二维码URL
	qrCodeUrl := generatePromoterQrCodeUrl(promoterCode)

	// 构建小程序页面URL
	pageURL := fmt.Sprintf("https://your-domain.com/pages/index/index?promoterCode=%s", promoterCode)

	response := &QRCodeResponse{
		Code: 0,
		Data: map[string]interface{}{
			"promoterCode": promoterCode,
			"qrCodeUrl":    qrCodeUrl,
			"pageURL":      pageURL,
		},
	}

	LogStep("二维码生成成功", map[string]interface{}{
		"promoterCode": promoterCode,
		"qrCodeUrl":    qrCodeUrl,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GenerateQRCodeBase64Handler 生成Base64编码的二维码接口
func GenerateQRCodeBase64Handler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理生成Base64二维码请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodGet {
		LogError("请求方法不支持", fmt.Errorf("期望GET方法，实际为%s", r.Method))
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	// 获取推广码参数
	promoterCode := r.URL.Query().Get("promoterCode")
	if promoterCode == "" {
		http.Error(w, "缺少promoterCode参数", http.StatusBadRequest)
		return
	}

	// 验证推广码格式
	if !utils.ValidatePromoterCode(promoterCode) {
		response := &QRCodeResponse{
			Code:     -1,
			ErrorMsg: "推广码格式无效",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 生成Base64编码的二维码
	base64QRCode := GeneratePromoterQRCodeBase64(promoterCode)

	// 构建小程序页面URL
	pageURL := fmt.Sprintf("https://your-domain.com/pages/index/index?promoterCode=%s", promoterCode)

	response := &QRCodeResponse{
		Code: 0,
		Data: map[string]interface{}{
			"promoterCode": promoterCode,
			"base64QRCode": base64QRCode,
			"pageURL":      pageURL,
		},
	}

	LogStep("Base64二维码生成成功", map[string]interface{}{
		"promoterCode": promoterCode,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
} 