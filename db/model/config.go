package model

import "time"

// ConfigModel 平台配置模型
type ConfigModel struct {
	Id          int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Key         string    `gorm:"column:key;uniqueIndex;not null" json:"key"`
	Value       string    `gorm:"column:value;not null" json:"value"`
	Description string    `gorm:"column:description" json:"description"`
	Type        string    `gorm:"column:type" json:"type"`               // string, number, boolean, json
	Status      int       `gorm:"column:status;default:1" json:"status"` // 1-启用，0-禁用
	CreatedAt   time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

// TableName 指定表名
func (ConfigModel) TableName() string {
	return "Configs"
}
