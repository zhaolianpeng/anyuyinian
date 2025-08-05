package service

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/skip2/go-qrcode"
)

// QRCodeService 二维码服务
type QRCodeService struct {
	BaseURL    string // 小程序页面基础URL
	OutputDir  string // 二维码图片输出目录
	PublicURL  string // 公共访问URL
}

// NewQRCodeService 创建二维码服务实例
func NewQRCodeService() *QRCodeService {
	return &QRCodeService{
		BaseURL:   "https://your-domain.com/pages/index/index?promoterCode=",
		OutputDir: "./static/qrcode",
		PublicURL: "https://your-domain.com/static/qrcode",
	}
}

// GeneratePromoterQRCode 生成推广二维码
func (s *QRCodeService) GeneratePromoterQRCode(promoterCode string) (string, error) {
	// 构建小程序页面URL
	pageURL := fmt.Sprintf("%s%s", s.BaseURL, promoterCode)
	
	// 生成二维码
	qrCode, err := qrcode.Encode(pageURL, qrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("生成二维码失败: %v", err)
	}

	// 确保输出目录存在
	if err := os.MkdirAll(s.OutputDir, 0755); err != nil {
		return "", fmt.Errorf("创建输出目录失败: %v", err)
	}

	// 生成文件名
	filename := fmt.Sprintf("promoter_%s.png", promoterCode)
	filepath := filepath.Join(s.OutputDir, filename)

	// 保存二维码图片
	if err := os.WriteFile(filepath, qrCode, 0644); err != nil {
		return "", fmt.Errorf("保存二维码图片失败: %v", err)
	}

	// 返回公共访问URL
	publicURL := fmt.Sprintf("%s/%s", s.PublicURL, filename)
	
	log.Printf("二维码生成成功: %s -> %s", pageURL, publicURL)
	return publicURL, nil
}

// GenerateQRCodeBase64 生成Base64编码的二维码（用于内联显示）
func (s *QRCodeService) GenerateQRCodeBase64(promoterCode string) (string, error) {
	// 构建小程序页面URL
	pageURL := fmt.Sprintf("%s%s", s.BaseURL, promoterCode)
	
	// 创建二维码
	qr, err := qrcode.New(pageURL, qrcode.Medium)
	if err != nil {
		return "", fmt.Errorf("创建二维码失败: %v", err)
	}

	// 编码为PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, qr.Image(256)); err != nil {
		return "", fmt.Errorf("编码PNG失败: %v", err)
	}

	// 转换为Base64
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
	dataURL := fmt.Sprintf("data:image/png;base64,%s", base64Str)
	
	return dataURL, nil
}

// 全局二维码服务实例
var qrCodeService = NewQRCodeService()

// GeneratePromoterQRCodeURL 生成推广二维码URL（便捷函数）
func GeneratePromoterQRCodeURL(promoterCode string) string {
	url, err := qrCodeService.GeneratePromoterQRCode(promoterCode)
	if err != nil {
		log.Printf("生成二维码失败: %v", err)
		// 返回一个占位符URL
		return fmt.Sprintf("https://via.placeholder.com/256x256/CCCCCC/666666?text=QR+Code+Error")
	}
	return url
}

// GeneratePromoterQRCodeBase64 生成Base64编码的二维码（便捷函数）
func GeneratePromoterQRCodeBase64(promoterCode string) string {
	base64Str, err := qrCodeService.GenerateQRCodeBase64(promoterCode)
	if err != nil {
		log.Printf("生成Base64二维码失败: %v", err)
		// 返回一个占位符图片
		return "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg=="
	}
	return base64Str
} 