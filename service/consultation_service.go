package service

import (
	"encoding/json"
	"fmt"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

type ConsultationService struct {
	consultationDAO *dao.ConsultationDAO
}

func NewConsultationService() *ConsultationService {
	return &ConsultationService{
		consultationDAO: &dao.ConsultationDAO{},
	}
}

// CreateConsultation 创建咨询会话
func (s *ConsultationService) CreateConsultation(userID, userName, userPhone string) (*model.Consultation, error) {
	// 检查用户是否已有活跃咨询
	existingConsultation, err := s.consultationDAO.GetConsultationByUserID(userID)
	if err == nil && existingConsultation != nil {
		// 用户已有活跃咨询，返回现有咨询
		return existingConsultation, nil
	}

	// 创建新的咨询会话
	consultation := &model.Consultation{
		UserID:    userID,
		UserName:  userName,
		UserPhone: userPhone,
		Status:    "waiting",
	}

	err = s.consultationDAO.CreateConsultation(consultation)
	if err != nil {
		return nil, fmt.Errorf("创建咨询会话失败: %v", err)
	}

	// 创建欢迎消息
	welcomeMessage := &model.ConsultationMessage{
		ConsultationID: consultation.ID,
		SenderType:     "admin",
		Content:        "您好！欢迎使用在线咨询服务，请问有什么可以帮助您的吗？",
		IsRead:         false,
	}

	err = s.consultationDAO.CreateMessage(welcomeMessage)
	if err != nil {
		// 即使创建欢迎消息失败，也不影响咨询会话的创建
		fmt.Printf("创建欢迎消息失败: %v\n", err)
	}

	// 创建通知
	notification := &model.ConsultationNotification{
		ConsultationID: consultation.ID,
		Type:           "new_message",
		Title:          "新用户咨询",
		Content:        fmt.Sprintf("用户 %s 发起新的咨询", userName),
		IsRead:         false,
	}

	err = s.consultationDAO.CreateNotification(notification)
	if err != nil {
		fmt.Printf("创建通知失败: %v\n", err)
	}

	return consultation, nil
}

// SendMessage 发送消息
func (s *ConsultationService) SendMessage(consultationID uint, content, senderType string) (*model.ConsultationMessage, error) {
	// 验证咨询会话是否存在
	consultation, err := s.consultationDAO.GetConsultationByID(consultationID)
	if err != nil {
		return nil, fmt.Errorf("咨询会话不存在: %v", err)
	}

	// 创建消息
	message := &model.ConsultationMessage{
		ConsultationID: consultationID,
		SenderType:     senderType,
		Content:        content,
		IsRead:         false,
	}

	err = s.consultationDAO.CreateMessage(message)
	if err != nil {
		return nil, fmt.Errorf("创建消息失败: %v", err)
	}

	// 更新咨询状态
	if senderType == "user" {
		// 用户发送消息，状态变为等待回复
		err = s.consultationDAO.UpdateConsultationStatus(consultationID, "waiting")
		if err != nil {
			fmt.Printf("更新咨询状态失败: %v\n", err)
		}

		// 创建通知
		notification := &model.ConsultationNotification{
			ConsultationID: consultationID,
			Type:           "new_message",
			Title:          "新用户消息",
			Content:        fmt.Sprintf("用户 %s 发送新消息: %s", consultation.UserName, content),
			IsRead:         false,
		}

		err = s.consultationDAO.CreateNotification(notification)
		if err != nil {
			fmt.Printf("创建通知失败: %v\n", err)
		}
	} else if senderType == "admin" {
		// 管理员回复，状态变为咨询中
		err = s.consultationDAO.UpdateConsultationStatus(consultationID, "chatting")
		if err != nil {
			fmt.Printf("更新咨询状态失败: %v\n", err)
		}
	}

	return message, nil
}

// GetConsultationMessages 获取咨询消息
func (s *ConsultationService) GetConsultationMessages(consultationID uint) ([]model.ConsultationMessage, error) {
	messages, err := s.consultationDAO.GetMessagesByConsultationID(consultationID)
	if err != nil {
		return nil, fmt.Errorf("获取消息失败: %v", err)
	}

	// 标记消息为已读
	err = s.consultationDAO.MarkAllMessagesAsRead(consultationID)
	if err != nil {
		fmt.Printf("标记消息为已读失败: %v\n", err)
	}

	return messages, nil
}

// GetConsultationStatus 获取咨询状态
func (s *ConsultationService) GetConsultationStatus(consultationID uint) (string, error) {
	consultation, err := s.consultationDAO.GetConsultationByID(consultationID)
	if err != nil {
		return "", fmt.Errorf("咨询会话不存在: %v", err)
	}

	return consultation.Status, nil
}

// CloseConsultation 关闭咨询会话
func (s *ConsultationService) CloseConsultation(consultationID uint) error {
	err := s.consultationDAO.UpdateConsultationStatus(consultationID, "closed")
	if err != nil {
		return fmt.Errorf("关闭咨询会话失败: %v", err)
	}

	// 创建关闭通知
	consultation, err := s.consultationDAO.GetConsultationByID(consultationID)
	if err != nil {
		fmt.Printf("获取咨询信息失败: %v\n", err)
		return nil
	}

	notification := &model.ConsultationNotification{
		ConsultationID: consultationID,
		Type:           "status_change",
		Title:          "咨询会话已关闭",
		Content:        fmt.Sprintf("用户 %s 的咨询会话已关闭", consultation.UserName),
		IsRead:         false,
	}

	err = s.consultationDAO.CreateNotification(notification)
	if err != nil {
		fmt.Printf("创建关闭通知失败: %v\n", err)
	}

	return nil
}

// GetActiveConsultations 获取活跃咨询会话
func (s *ConsultationService) GetActiveConsultations() ([]model.Consultation, error) {
	consultations, err := s.consultationDAO.GetActiveConsultations()
	if err != nil {
		return nil, fmt.Errorf("获取活跃咨询失败: %v", err)
	}

	return consultations, nil
}

// GetConsultationsByStatus 根据状态获取咨询会话
func (s *ConsultationService) GetConsultationsByStatus(status string) ([]model.Consultation, error) {
	consultations, err := s.consultationDAO.GetConsultationsByStatus(status)
	if err != nil {
		return nil, fmt.Errorf("获取咨询会话失败: %v", err)
	}

	return consultations, nil
}

// GetUnreadNotifications 获取未读通知
func (s *ConsultationService) GetUnreadNotifications() ([]model.ConsultationNotification, error) {
	notifications, err := s.consultationDAO.GetUnreadNotifications()
	if err != nil {
		return nil, fmt.Errorf("获取未读通知失败: %v", err)
	}

	return notifications, nil
}

// MarkNotificationAsRead 标记通知为已读
func (s *ConsultationService) MarkNotificationAsRead(notificationID uint) error {
	err := s.consultationDAO.MarkNotificationAsRead(notificationID)
	if err != nil {
		return fmt.Errorf("标记通知为已读失败: %v", err)
	}

	return nil
}

// GetConsultationStats 获取咨询统计信息
func (s *ConsultationService) GetConsultationStats() (map[string]interface{}, error) {
	stats, err := s.consultationDAO.GetConsultationStats()
	if err != nil {
		return nil, fmt.Errorf("获取统计信息失败: %v", err)
	}

	return stats, nil
}

// ProcessWebSocketMessage 处理WebSocket消息
func (s *ConsultationService) ProcessWebSocketMessage(messageData []byte) error {
	var message struct {
		Type string `json:"type"`
		Data struct {
			ConsultationID uint   `json:"consultationId"`
			Content        string `json:"content"`
			SenderType     string `json:"senderType"`
		} `json:"data"`
	}

	err := json.Unmarshal(messageData, &message)
	if err != nil {
		return fmt.Errorf("解析WebSocket消息失败: %v", err)
	}

	switch message.Type {
	case "send_message":
		// 处理发送消息
		_, err = s.SendMessage(message.Data.ConsultationID, message.Data.Content, message.Data.SenderType)
		if err != nil {
			return fmt.Errorf("处理发送消息失败: %v", err)
		}
	case "heartbeat":
		// 处理心跳消息
		fmt.Println("收到心跳消息")
	default:
		return fmt.Errorf("未知的消息类型: %s", message.Type)
	}

	return nil
}
