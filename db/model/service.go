package model

import "time"

// ServiceItemModel 服务项目模型
type ServiceItemModel struct {
	Id            int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name          string    `gorm:"column:name;not null" json:"name"`
	Description   string    `gorm:"column:description" json:"description"`
	Category      string    `gorm:"column:category;not null" json:"category"`
	Price         float64   `gorm:"column:price;not null" json:"price"`
	OriginalPrice float64   `gorm:"column:originalPrice" json:"originalPrice"`
	ImageUrl      string    `gorm:"column:imageUrl" json:"imageUrl"`
	DetailImages  string    `gorm:"column:detailImages" json:"detailImages"` // JSON数组
	FormConfig    string    `gorm:"column:formConfig" json:"formConfig"`     // JSON配置
	Status        int       `gorm:"column:status;default:1" json:"status"`   // 1-上架，0-下架
	Sort          int       `gorm:"column:sort;default:0" json:"sort"`
	CreatedAt     time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

// TableName 指定表名
func (ServiceItemModel) TableName() string {
	return "ServiceItems"
}
