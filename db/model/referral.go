package model

import "time"

// ReferralModel 推荐关系模型
type ReferralModel struct {
	Id         int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId     string    `gorm:"column:userId;uniqueIndex;not null;type:varchar(24)" json:"userId"`
	ReferrerId int32     `gorm:"column:referrerId;not null" json:"referrerId"` // 推荐人ID
	QrCodeUrl  string    `gorm:"column:qrCodeUrl" json:"qrCodeUrl"`            // 专属二维码URL
	Status     int       `gorm:"column:status;default:1" json:"status"`        // 1-正常，0-禁用
	CreatedAt  time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

// CommissionModel 佣金记录模型
type CommissionModel struct {
	Id          int32      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId      string     `gorm:"column:userId;not null;type:varchar(24)" json:"userId"`
	OrderId     int32      `gorm:"column:orderId;not null" json:"orderId"`
	OrderNo     string     `gorm:"column:orderNo;not null" json:"orderNo"`
	Amount      float64    `gorm:"column:amount;not null" json:"amount"`  // 佣金金额
	Rate        float64    `gorm:"column:rate;not null" json:"rate"`      // 佣金比例
	Status      int        `gorm:"column:status;default:0" json:"status"` // 0-待结算，1-已结算，2-已提现
	CashoutTime *time.Time `gorm:"column:cashoutTime" json:"cashoutTime"`
	CreatedAt   time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
}

// CashoutModel 提现记录模型
type CashoutModel struct {
	Id          int32      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId      string     `gorm:"column:userId;not null;type:varchar(24)" json:"userId"`
	Amount      float64    `gorm:"column:amount;not null" json:"amount"`
	Method      string     `gorm:"column:method;not null" json:"method"`   // 提现方式：wechat, alipay, bank
	Account     string     `gorm:"column:account;not null" json:"account"` // 提现账户
	Status      int        `gorm:"column:status;default:0" json:"status"`  // 0-待审核，1-已通过，2-已拒绝，3-已到账
	Remark      string     `gorm:"column:remark" json:"remark"`
	ProcessTime *time.Time `gorm:"column:processTime" json:"processTime"`
	CreatedAt   time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
}

// TableName 指定表名
func (ReferralModel) TableName() string {
	return "Referrals"
}

func (CommissionModel) TableName() string {
	return "Commissions"
}

func (CashoutModel) TableName() string {
	return "Cashouts"
}
