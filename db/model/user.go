package model

import "time"

// UserModel 用户模型
type UserModel struct {
	Id          int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId      string    `gorm:"column:userId;uniqueIndex;type:varchar(24);not null" json:"userId"`
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
	// 管理员相关字段
	IsAdmin        int        `gorm:"column:isAdmin;default:0" json:"isAdmin"`
	AdminLevel     int        `gorm:"column:adminLevel;default:0" json:"adminLevel"`
	AdminPassword  string     `gorm:"column:adminPassword" json:"-"`
	AdminUsername  string     `gorm:"column:adminUsername" json:"adminUsername"`
	ParentAdminId  string     `gorm:"column:parentAdminId" json:"parentAdminId"`
	AdminCreatedAt *time.Time `gorm:"column:adminCreatedAt" json:"adminCreatedAt"`
}

// AdminLoginLogModel 管理员登录记录模型
type AdminLoginLogModel struct {
	Id          int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	AdminUserId string    `gorm:"column:adminUserId;not null" json:"adminUserId"`
	LoginTime   time.Time `gorm:"column:loginTime;default:CURRENT_TIMESTAMP" json:"loginTime"`
	LoginIp     string    `gorm:"column:loginIp" json:"loginIp"`
	UserAgent   string    `gorm:"column:userAgent" json:"userAgent"`
	Status      int       `gorm:"column:status;default:1" json:"status"`
	Remark      string    `gorm:"column:remark" json:"remark"`
	CreatedAt   time.Time `gorm:"column:createdAt;default:CURRENT_TIMESTAMP" json:"createdAt"`
}

// TableName 指定表名
func (UserModel) TableName() string {
	return "Users"
}

// TableName 指定表名
func (AdminLoginLogModel) TableName() string {
	return "AdminLoginLogs"
}
