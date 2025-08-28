package dao

import (
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

type ConsultationDAO struct{}

const consultationTableName = "consultations"
const consultationMessageTableName = "consultation_messages"
const consultationNotificationTableName = "consultation_notifications"

// CreateConsultation 创建咨询会话
func (dao *ConsultationDAO) CreateConsultation(consultation *model.Consultation) error {
	cli := db.Get()
	consultation.CreatedAt = time.Now()
	consultation.UpdatedAt = time.Now()
	consultation.LastMessage = time.Now()
	return cli.Table(consultationTableName).Create(consultation).Error
}

// GetConsultationByID 根据ID获取咨询会话
func (dao *ConsultationDAO) GetConsultationByID(id uint) (*model.Consultation, error) {
	var consultation model.Consultation
	cli := db.Get()
	err := cli.Table(consultationTableName).Where("id = ?", id).First(&consultation).Error
	if err != nil {
		return nil, err
	}

	// 获取消息
	messages, err := dao.GetMessagesByConsultationID(id)
	if err != nil {
		return nil, err
	}
	consultation.Messages = messages

	return &consultation, nil
}

// GetConsultationByUserID 根据用户ID获取咨询会话
func (dao *ConsultationDAO) GetConsultationByUserID(userID string) (*model.Consultation, error) {
	var consultation model.Consultation
	cli := db.Get()
	err := cli.Table(consultationTableName).
		Where("user_id = ? AND status != ?", userID, "closed").
		Order("created_at DESC").
		First(&consultation).Error
	if err != nil {
		return nil, err
	}

	// 获取消息
	messages, err := dao.GetMessagesByConsultationID(consultation.ID)
	if err != nil {
		return nil, err
	}
	consultation.Messages = messages

	return &consultation, nil
}

// UpdateConsultationStatus 更新咨询状态
func (dao *ConsultationDAO) UpdateConsultationStatus(id uint, status string) error {
	cli := db.Get()
	return cli.Table(consultationTableName).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       status,
			"updated_at":   time.Now(),
			"last_message": time.Now(),
		}).Error
}

// GetActiveConsultations 获取活跃的咨询会话
func (dao *ConsultationDAO) GetActiveConsultations() ([]model.Consultation, error) {
	var consultations []model.Consultation
	cli := db.Get()
	err := cli.Table(consultationTableName).
		Where("status IN ?", []string{"waiting", "chatting"}).
		Order("last_message DESC").
		Find(&consultations).Error
	if err != nil {
		return nil, err
	}

	// 为每个咨询获取消息
	for i := range consultations {
		messages, err := dao.GetMessagesByConsultationID(consultations[i].ID)
		if err != nil {
			continue
		}
		consultations[i].Messages = messages
	}

	return consultations, nil
}

// GetConsultationsByStatus 根据状态获取咨询会话
func (dao *ConsultationDAO) GetConsultationsByStatus(status string) ([]model.Consultation, error) {
	var consultations []model.Consultation
	cli := db.Get()
	err := cli.Table(consultationTableName).
		Where("status = ?", status).
		Order("created_at DESC").
		Find(&consultations).Error
	if err != nil {
		return nil, err
	}

	// 为每个咨询获取消息
	for i := range consultations {
		messages, err := dao.GetMessagesByConsultationID(consultations[i].ID)
		if err != nil {
			continue
		}
		consultations[i].Messages = messages
	}

	return consultations, nil
}

// CreateMessage 创建消息
func (dao *ConsultationDAO) CreateMessage(message *model.ConsultationMessage) error {
	cli := db.Get()
	message.CreatedAt = time.Now()
	return cli.Table(consultationMessageTableName).Create(message).Error
}

// GetMessagesByConsultationID 根据咨询ID获取消息
func (dao *ConsultationDAO) GetMessagesByConsultationID(consultationID uint) ([]model.ConsultationMessage, error) {
	var messages []model.ConsultationMessage
	cli := db.Get()
	err := cli.Table(consultationMessageTableName).
		Where("consultation_id = ?", consultationID).
		Order("created_at ASC").
		Find(&messages).Error
	return messages, err
}

// MarkMessageAsRead 标记消息为已读
func (dao *ConsultationDAO) MarkMessageAsRead(messageID uint) error {
	cli := db.Get()
	return cli.Table(consultationMessageTableName).
		Where("id = ?", messageID).
		Update("is_read", true).Error
}

// MarkAllMessagesAsRead 标记咨询会话的所有消息为已读
func (dao *ConsultationDAO) MarkAllMessagesAsRead(consultationID uint) error {
	cli := db.Get()
	return cli.Table(consultationMessageTableName).
		Where("consultation_id = ?", consultationID).
		Update("is_read", true).Error
}

// GetUnreadMessageCount 获取未读消息数量
func (dao *ConsultationDAO) GetUnreadMessageCount(consultationID uint) (int64, error) {
	var count int64
	cli := db.Get()
	err := cli.Table(consultationMessageTableName).
		Where("consultation_id = ? AND is_read = ?", consultationID, false).
		Count(&count).Error
	return count, err
}

// CreateNotification 创建通知
func (dao *ConsultationDAO) CreateNotification(notification *model.ConsultationNotification) error {
	cli := db.Get()
	notification.CreatedAt = time.Now()
	return cli.Table(consultationNotificationTableName).Create(notification).Error
}

// GetUnreadNotifications 获取未读通知
func (dao *ConsultationDAO) GetUnreadNotifications() ([]model.ConsultationNotification, error) {
	var notifications []model.ConsultationNotification
	cli := db.Get()
	err := cli.Table(consultationNotificationTableName).
		Where("is_read = ?", false).
		Order("created_at DESC").
		Find(&notifications).Error
	if err != nil {
		return nil, err
	}

	// 为每个通知获取咨询信息
	for i := range notifications {
		consultation, err := dao.GetConsultationByID(notifications[i].ConsultationID)
		if err != nil {
			continue
		}
		notifications[i].Consultation = *consultation
	}

	return notifications, nil
}

// MarkNotificationAsRead 标记通知为已读
func (dao *ConsultationDAO) MarkNotificationAsRead(notificationID uint) error {
	cli := db.Get()
	return cli.Table(consultationNotificationTableName).
		Where("id = ?", notificationID).
		Update("is_read", true).Error
}

// GetConsultationStats 获取咨询统计信息
func (dao *ConsultationDAO) GetConsultationStats() (map[string]interface{}, error) {
	var stats = make(map[string]interface{})
	cli := db.Get()

	// 总咨询数
	var totalCount int64
	err := cli.Table(consultationTableName).Count(&totalCount).Error
	if err != nil {
		return nil, err
	}
	stats["totalCount"] = totalCount

	// 等待回复数
	var waitingCount int64
	err = cli.Table(consultationTableName).Where("status = ?", "waiting").Count(&waitingCount).Error
	if err != nil {
		return nil, err
	}
	stats["waitingCount"] = waitingCount

	// 咨询中数量
	var chattingCount int64
	err = cli.Table(consultationTableName).Where("status = ?", "chatting").Count(&chattingCount).Error
	if err != nil {
		return nil, err
	}
	stats["chattingCount"] = chattingCount

	// 今日新增咨询数
	var todayCount int64
	today := time.Now().Format("2006-01-02")
	err = cli.Table(consultationTableName).
		Where("DATE(created_at) = ?", today).
		Count(&todayCount).Error
	if err != nil {
		return nil, err
	}
	stats["todayCount"] = todayCount

	return stats, nil
}
