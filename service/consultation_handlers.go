package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

// CreateConsultationHandler 创建咨询会话处理器
func CreateConsultationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID    interface{} `json:"userId"` // 使用interface{}接受任何类型
		UserName  string      `json:"userName"`
		UserPhone string      `json:"userPhone"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 验证请求参数
	if req.UserName == "" {
		http.Error(w, "UserName is required", http.StatusBadRequest)
		return
	}

	// 转换UserID为字符串
	var userID string
	switch v := req.UserID.(type) {
	case string:
		userID = v
	case int, int32, int64, float32, float64:
		userID = fmt.Sprintf("%v", v)
	default:
		userID = fmt.Sprintf("%v", v)
	}

	// 创建咨询会话
	consultationService := NewConsultationService()
	consultation, err := consultationService.CreateConsultation(userID, req.UserName, req.UserPhone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回响应
	response := map[string]interface{}{
		"consultationId": consultation.ID,
		"status":         consultation.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"data": response,
	})
}

// GetConsultationMessagesHandler 获取咨询消息处理器
func GetConsultationMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 从查询参数获取consultationId
	consultationIDStr := r.URL.Query().Get("consultationId")
	if consultationIDStr == "" {
		http.Error(w, "ConsultationID is required", http.StatusBadRequest)
		return
	}

	consultationID, err := strconv.ParseUint(consultationIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ConsultationID", http.StatusBadRequest)
		return
	}

	// 获取消息
	consultationService := NewConsultationService()
	messages, err := consultationService.GetConsultationMessages(uint(consultationID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"messages": messages,
		},
	})
}

// SendConsultationMessageHandler 发送咨询消息处理器
func SendConsultationMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 添加调试日志
	log.Printf("[DEBUG] 开始处理发送消息请求")

	// 重新读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] 读取请求体失败: %v", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	log.Printf("[DEBUG] 请求体内容: %s", string(body))

	var req struct {
		ConsultationID uint   `json:"consultationId"`
		Content        string `json:"content"`
		SenderType     string `json:"senderType"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		log.Printf("[ERROR] JSON解析失败: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("[DEBUG] 解析后的请求: ConsultationID=%d, Content=%s, SenderType=%s",
		req.ConsultationID, req.Content, req.SenderType)

	// 验证请求参数
	if req.ConsultationID == 0 {
		http.Error(w, "ConsultationID is required", http.StatusBadRequest)
		return
	}
	if req.Content == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}
	if req.SenderType == "" {
		http.Error(w, "SenderType is required", http.StatusBadRequest)
		return
	}

	// 发送消息
	consultationService := NewConsultationService()
	message, err := consultationService.SendMessage(req.ConsultationID, req.Content, req.SenderType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"data": message,
	})
}

// GetConsultationStatusHandler 获取咨询状态处理器
func GetConsultationStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 从查询参数获取consultationId
	consultationIDStr := r.URL.Query().Get("consultationId")
	if consultationIDStr == "" {
		http.Error(w, "ConsultationID is required", http.StatusBadRequest)
		return
	}

	consultationID, err := strconv.ParseUint(consultationIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ConsultationID", http.StatusBadRequest)
		return
	}

	// 获取状态
	consultationService := NewConsultationService()
	status, err := consultationService.GetConsultationStatus(uint(consultationID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"status": status,
		},
	})
}

// CloseConsultationHandler 关闭咨询会话处理器
func CloseConsultationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ConsultationID uint `json:"consultationId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 验证请求参数
	if req.ConsultationID == 0 {
		http.Error(w, "ConsultationID is required", http.StatusBadRequest)
		return
	}

	// 关闭咨询会话
	consultationService := NewConsultationService()
	err := consultationService.CloseConsultation(req.ConsultationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"message": "咨询会话已关闭",
		},
	})
}

// GetActiveConsultationsHandler 获取活跃咨询列表处理器
func GetActiveConsultationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 获取活跃咨询
	consultationService := NewConsultationService()
	consultations, err := consultationService.GetActiveConsultations()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"consultations": consultations,
		},
	})
}

// GetConsultationStatsHandler 获取咨询统计处理器
func GetConsultationStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 获取统计信息
	consultationService := NewConsultationService()
	stats, err := consultationService.GetConsultationStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"data": stats,
	})
}

// GetUnreadNotificationsHandler 获取未读通知处理器
func GetUnreadNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 获取未读通知
	consultationService := NewConsultationService()
	notifications, err := consultationService.GetUnreadNotifications()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"notifications": notifications,
		},
	})
}

// MarkNotificationAsReadHandler 标记通知为已读处理器
func MarkNotificationAsReadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		NotificationID uint `json:"notificationId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 验证请求参数
	if req.NotificationID == 0 {
		http.Error(w, "NotificationID is required", http.StatusBadRequest)
		return
	}

	// 标记通知为已读
	consultationService := NewConsultationService()
	err := consultationService.MarkNotificationAsRead(req.NotificationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"message": "通知已标记为已读",
		},
	})
}
