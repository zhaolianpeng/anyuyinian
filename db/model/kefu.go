package model

import "time"

// KefuMessageModel 客服消息模型
type KefuMessageModel struct {
	Id           int32      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId       int32      `gorm:"column:userId;not null" json:"userId"`
	UserName     string     `gorm:"column:userName" json:"userName"`
	UserAvatar   string     `gorm:"column:userAvatar" json:"userAvatar"`
	Type         int        `gorm:"column:type;default:1" json:"type"` // 1-用户消息，2-客服回复
	Content      string     `gorm:"column:content;not null" json:"content"`
	Images       string     `gorm:"column:images" json:"images"`           // JSON数组
	Status       int        `gorm:"column:status;default:0" json:"status"` // 0-未读，1-已读，2-已回复
	ReplyContent string     `gorm:"column:replyContent" json:"replyContent"`
	ReplyTime    *time.Time `gorm:"column:replyTime" json:"replyTime"`
	ReplyUserId  int32      `gorm:"column:replyUserId" json:"replyUserId"`
	CreatedAt    time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
}

// FaqModel 常见问题模型
type FaqModel struct {
	Id        int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Question  string    `gorm:"column:question;not null" json:"question"`
	Answer    string    `gorm:"column:answer;not null" json:"answer"`
	Category  string    `gorm:"column:category;not null" json:"category"`
	Sort      int       `gorm:"column:sort;default:0" json:"sort"`
	Status    int       `gorm:"column:status;default:1" json:"status"` // 1-启用，0-禁用
	ViewCount int       `gorm:"column:viewCount;default:0" json:"viewCount"`
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

// TableName 指定表名
func (KefuMessageModel) TableName() string {
	return "KefuMessages"
}

func (FaqModel) TableName() string {
	return "Faqs"
}
