package model

import "time"

// UserModel 用户模型
type UserModel struct {
	Id          int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OpenId      string    `gorm:"column:openId;uniqueIndex;not null" json:"openId"`
	UnionId     string    `gorm:"column:unionId" json:"unionId"`
	NickName    string    `gorm:"column:nickName" json:"nickName"`
	AvatarUrl   string    `gorm:"column:avatarUrl" json:"avatarUrl"`
	Gender      int       `gorm:"column:gender" json:"gender"`
	Phone       string    `gorm:"column:phone" json:"phone"`
	Country     string    `gorm:"column:country" json:"country"`
	Province    string    `gorm:"column:province" json:"province"`
	City        string    `gorm:"column:city" json:"city"`
	Language    string    `gorm:"column:language" json:"language"`
	SessionKey  string    `gorm:"column:sessionKey" json:"-"`
	LastLoginAt time.Time `gorm:"column:lastLoginAt" json:"lastLoginAt"`
	CreatedAt   time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

// TableName 指定表名
func (UserModel) TableName() string {
	return "Users"
}
