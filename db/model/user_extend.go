package model

import "time"

// UserAddressModel 用户地址模型
type UserAddressModel struct {
	Id        int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId    string    `gorm:"column:userId;not null;type:varchar(24)" json:"userId"`
	Name      string    `gorm:"column:name;not null" json:"name"`
	Phone     string    `gorm:"column:phone;not null" json:"phone"`
	Province  string    `gorm:"column:province" json:"province"`
	City      string    `gorm:"column:city" json:"city"`
	District  string    `gorm:"column:district" json:"district"`
	Address   string    `gorm:"column:address;not null" json:"address"`
	IsDefault int       `gorm:"column:isDefault;default:0" json:"isDefault"` // 1-默认地址，0-非默认
	Status    int       `gorm:"column:status;default:1" json:"status"`       // 1-正常，0-删除
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

// PatientModel 就诊人信息模型
type PatientModel struct {
	Id        int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId    string    `gorm:"column:userId;not null;type:varchar(24)" json:"userId"`
	Name      string    `gorm:"column:name;not null" json:"name"`
	IdCard    string    `gorm:"column:idCard" json:"idCard"`
	Phone     string    `gorm:"column:phone" json:"phone"`
	Gender    int       `gorm:"column:gender;default:0" json:"gender"` // 0-未知，1-男，2-女
	Birthday  string    `gorm:"column:birthday" json:"birthday"`
	Relation  string    `gorm:"column:relation" json:"relation"`             // 与用户关系：本人、父亲、母亲等
	IsDefault int       `gorm:"column:isDefault;default:0" json:"isDefault"` // 1-默认就诊人，0-非默认
	Status    int       `gorm:"column:status;default:1" json:"status"`       // 1-正常，0-删除
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

// TableName 指定表名
func (UserAddressModel) TableName() string {
	return "UserAddresses"
}

func (PatientModel) TableName() string {
	return "Patients"
}
