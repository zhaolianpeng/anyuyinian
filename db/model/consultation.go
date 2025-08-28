package model

import (
	"time"
)

// Consultation 咨询会话模型
type Consultation struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      string    `json:"userId" gorm:"type:varchar(100);not null;index"`
	UserName    string    `json:"userName" gorm:"type:varchar(100);not null"`
	UserPhone   string    `json:"userPhone" gorm:"type:varchar(20)"`
	Status      string    `json:"status" gorm:"type:varchar(20);default:'waiting';comment:waiting-等待回复,chatting-咨询中,closed-已结束"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	LastMessage time.Time `json:"lastMessage" gorm:"comment:最后一条消息时间"`

	// 关联关系
	Messages []ConsultationMessage `json:"messages,omitempty" gorm:"foreignKey:ConsultationID"`
}

// ConsultationMessage 咨询消息模型
type ConsultationMessage struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ConsultationID uint      `json:"consultationId" gorm:"not null;index"`
	SenderType     string    `json:"senderType" gorm:"type:varchar(20);not null;comment:user-用户,admin-管理员"`
	Content        string    `json:"content" gorm:"type:text;not null"`
	IsRead         bool      `json:"isRead" gorm:"default:false;comment:是否已读"`
	CreatedAt      time.Time `json:"createdAt" gorm:"autoCreateTime"`

	// 关联关系
	Consultation Consultation `json:"consultation,omitempty" gorm:"foreignKey:ConsultationID"`
}

// ConsultationNotification 咨询通知模型
type ConsultationNotification struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ConsultationID uint      `json:"consultationId" gorm:"not null;index"`
	Type           string    `json:"type" gorm:"type:varchar(20);not null;comment:new_message-新消息,status_change-状态变更"`
	Title          string    `json:"title" gorm:"type:varchar(200);not null"`
	Content        string    `json:"content" gorm:"type:text"`
	IsRead         bool      `json:"isRead" gorm:"default:false"`
	CreatedAt      time.Time `json:"createdAt" gorm:"autoCreateTime"`

	// 关联关系
	Consultation Consultation `json:"consultation,omitempty" gorm:"foreignKey:ConsultationID"`
}

// TableName 指定表名
func (Consultation) TableName() string {
	return "consultations"
}

func (ConsultationMessage) TableName() string {
	return "consultation_messages"
}

func (ConsultationNotification) TableName() string {
	return "consultation_notifications"
}
