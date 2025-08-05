package model

import "time"

// BannerModel 轮播图模型
type BannerModel struct {
	Id        int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title     string    `gorm:"column:title" json:"title"`
	ImageUrl  string    `gorm:"column:imageUrl;not null" json:"imageUrl"`
	LinkUrl   string    `gorm:"column:linkUrl" json:"linkUrl"`
	Sort      int       `gorm:"column:sort;default:0" json:"sort"`
	Status    int       `gorm:"column:status;default:1" json:"status"` // 1-启用，0-禁用
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

// NavigationModel 导航模型
type NavigationModel struct {
	Id        int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;not null" json:"name"`
	Icon      string    `gorm:"column:icon;not null" json:"icon"`
	LinkUrl   string    `gorm:"column:linkUrl" json:"linkUrl"`
	Sort      int       `gorm:"column:sort;default:0" json:"sort"`
	Status    int       `gorm:"column:status;default:1" json:"status"` // 1-启用，0-禁用
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

// ServiceModel 服务项模型
type ServiceModel struct {
	Id            int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ServiceItemId int32     `gorm:"column:serviceitemid" json:"serviceId"` // 服务项目ID，用于前端跳转
	Name          string    `gorm:"column:name;not null" json:"name"`
	Description   string    `gorm:"column:description" json:"description"`
	Icon          string    `gorm:"column:icon;not null" json:"icon"`
	ImageUrl      string    `gorm:"column:imageUrl" json:"imageUrl"`
	LinkUrl       string    `gorm:"column:linkUrl" json:"linkUrl"`
	Sort          int       `gorm:"column:sort;default:0" json:"sort"`
	Status        int       `gorm:"column:status;default:1" json:"status"` // 1-启用，0-禁用
	CreatedAt     time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

// HospitalModel 医院模型
type HospitalModel struct {
	Id          int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"column:name;not null" json:"name"`
	Logo        string    `gorm:"column:logo" json:"logo"`
	Address     string    `gorm:"column:address" json:"address"`
	Phone       string    `gorm:"column:phone" json:"phone"`
	Description string    `gorm:"column:description" json:"description"`
	Level       string    `gorm:"column:level" json:"level"` // 医院等级
	Type        string    `gorm:"column:type" json:"type"`   // 医院类型
	Longitude   float64   `gorm:"column:longitude" json:"longitude"`
	Latitude    float64   `gorm:"column:latitude" json:"latitude"`
	Sort        int       `gorm:"column:sort;default:0" json:"sort"`
	Status      int       `gorm:"column:status;default:1" json:"status"` // 1-启用，0-禁用
	CreatedAt   time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

// TableName 指定表名
func (BannerModel) TableName() string {
	return "Banners"
}

func (NavigationModel) TableName() string {
	return "Navigations"
}

func (ServiceModel) TableName() string {
	return "Services"
}

func (HospitalModel) TableName() string {
	return "Hospitals"
}
