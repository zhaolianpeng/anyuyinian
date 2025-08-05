package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

// GeneratePromoterCode 生成六位随机推广码
// 格式：字母数字组合，如 ABC123
func GeneratePromoterCode() string {
	// 定义字符集：大写字母和数字
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 6

	// 生成随机码
	code := make([]byte, codeLength)
	for i := 0; i < codeLength; i++ {
		// 生成随机索引
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// 如果随机数生成失败，使用时间戳作为备选方案
			randomIndex = big.NewInt(int64(i*7) % int64(len(charset)))
		}
		code[i] = charset[randomIndex.Int64()]
	}

	return string(code)
}

// ValidatePromoterCode 验证推广码格式
func ValidatePromoterCode(code string) bool {
	if len(code) != 6 {
		return false
	}

	// 检查是否只包含大写字母和数字
	for _, char := range code {
		if !((char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
			return false
		}
	}

	return true
}

// FormatPromoterCode 格式化推广码显示
// 在适当位置添加分隔符，如 ABC-123
func FormatPromoterCode(code string) string {
	if len(code) != 6 {
		return code
	}

	return fmt.Sprintf("%s-%s", code[:3], code[3:])
}

// GenerateUniquePromoterCode 生成唯一的推广码
// 需要检查数据库中是否已存在
func GenerateUniquePromoterCode(checkExists func(string) bool) string {
	maxAttempts := 100 // 最大尝试次数，避免无限循环

	for i := 0; i < maxAttempts; i++ {
		code := GeneratePromoterCode()
		if !checkExists(code) {
			return code
		}
	}

	// 如果生成失败，返回带时间戳的码
	return fmt.Sprintf("P%d", time.Now().Unix()%100000)
}
