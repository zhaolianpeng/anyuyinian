package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// GenerateMongoID 生成MongoDB风格的24位十六进制ID
func GenerateMongoID() string {
	// 生成12字节的随机数
	bytes := make([]byte, 12)
	_, err := rand.Read(bytes)
	if err != nil {
		// 如果随机数生成失败，使用时间戳作为备选方案
		now := time.Now()
		bytes = []byte{
			byte(now.Unix() >> 24),
			byte(now.Unix() >> 16),
			byte(now.Unix() >> 8),
			byte(now.Unix()),
			byte(now.Nanosecond() >> 24),
			byte(now.Nanosecond() >> 16),
			byte(now.Nanosecond() >> 8),
			byte(now.Nanosecond()),
			byte(now.UnixNano() >> 56),
			byte(now.UnixNano() >> 48),
			byte(now.UnixNano() >> 40),
			byte(now.UnixNano() >> 32),
		}
	}
	return hex.EncodeToString(bytes)
}

// GenerateUserID 生成用户ID（MongoDB风格）
func GenerateUserID() string {
	return GenerateMongoID()
}

// IsValidMongoID 验证是否为有效的MongoDB ID格式
func IsValidMongoID(id string) bool {
	if len(id) != 24 {
		return false
	}
	_, err := hex.DecodeString(id)
	return err == nil
}

// FormatUserID 格式化用户ID，确保是字符串格式
func FormatUserID(id interface{}) string {
	switch v := id.(type) {
	case string:
		return v
	case int, int32, int64:
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
