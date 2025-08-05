package service

import (
	"encoding/json"
	"net/http"
	"wxcloudrun-golang/utils"
)

// TestResponse 测试响应
type TestResponse struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

// GenerateUserIdHandler 生成UserId的测试接口
func GenerateUserIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	// 生成UserId
	userId := utils.GenerateUserID()

	response := &TestResponse{
		Code: 0,
		Data: map[string]interface{}{
			"userId": userId,
			"length": len(userId),
			"valid":  utils.IsValidMongoID(userId),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// TestUserIdMigrationHandler 测试UserId迁移的接口
func TestUserIdMigrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	// 创建迁移服务
	migrationService := NewMigrationService()

	// 执行迁移
	err := migrationService.MigrateExistingUsers()
	if err != nil {
		response := &TestResponse{
			Code:     -1,
			ErrorMsg: "迁移失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 验证迁移结果
	err = migrationService.ValidateUserIds()
	if err != nil {
		response := &TestResponse{
			Code:     -1,
			ErrorMsg: "验证失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := &TestResponse{
		Code: 0,
		Data: map[string]string{
			"message": "UserId迁移成功",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
