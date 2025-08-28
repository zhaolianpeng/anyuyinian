package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"wxcloudrun-golang/service"
)

type ConsultationHandler struct {
	consultationService *service.ConsultationService
}

func NewConsultationHandler() *ConsultationHandler {
	return &ConsultationHandler{
		consultationService: service.NewConsultationService(),
	}
}

// CreateConsultationRequest 创建咨询请求
type CreateConsultationRequest struct {
	UserID    string `json:"userId"`
	UserName  string `json:"userName"`
	UserPhone string `json:"userPhone"`
}

// CreateConsultationResponse 创建咨询响应
type CreateConsultationResponse struct {
	ConsultationID uint   `json:"consultationId"`
	Status         string `json:"status"`
}

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	ConsultationID uint   `json:"consultationId"`
	Content        string `json:"content"`
	SenderType     string `json:"senderType"`
}

// GetMessagesRequest 获取消息请求
type GetMessagesRequest struct {
	ConsultationID uint `json:"consultationId"`
}

// GetStatusRequest 获取状态请求
type GetStatusRequest struct {
	ConsultationID uint `json:"consultationId"`
}

// CloseConsultationRequest 关闭咨询请求
type CloseConsultationRequest struct {
	ConsultationID uint `json:"consultationId"`
}

// CreateConsultation 创建咨询会话
func (h *ConsultationHandler) CreateConsultation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateConsultationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 验证请求参数
	if req.UserName == "" {
		http.Error(w, "UserName is required", http.StatusBadRequest)
		return
	}

	// 创建咨询会话
	consultation, err := h.consultationService.CreateConsultation(req.UserID, req.UserName, req.UserPhone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回响应
	response := CreateConsultationResponse{
		ConsultationID: consultation.ID,
		Status:         consultation.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"data": response,
	})
}

// SendMessage 发送消息
func (h *ConsultationHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

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
	message, err := h.consultationService.SendMessage(req.ConsultationID, req.Content, req.SenderType)
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

// GetMessages 获取咨询消息
func (h *ConsultationHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
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
	messages, err := h.consultationService.GetConsultationMessages(uint(consultationID))
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

// GetStatus 获取咨询状态
func (h *ConsultationHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
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
	status, err := h.consultationService.GetConsultationStatus(uint(consultationID))
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

// CloseConsultation 关闭咨询会话
func (h *ConsultationHandler) CloseConsultation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CloseConsultationRequest
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
	err := h.consultationService.CloseConsultation(req.ConsultationID)
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

// GetActiveConsultations 获取活跃咨询列表
func (h *ConsultationHandler) GetActiveConsultations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 获取活跃咨询
	consultations, err := h.consultationService.GetActiveConsultations()
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

// GetUnreadNotifications 获取未读通知
func (h *ConsultationHandler) GetUnreadNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 获取未读通知
	notifications, err := h.consultationService.GetUnreadNotifications()
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

// MarkNotificationAsRead 标记通知为已读
func (h *ConsultationHandler) MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
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
	err := h.consultationService.MarkNotificationAsRead(req.NotificationID)
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

// GetConsultationStats 获取咨询统计
func (h *ConsultationHandler) GetConsultationStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 获取统计信息
	stats, err := h.consultationService.GetConsultationStats()
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
