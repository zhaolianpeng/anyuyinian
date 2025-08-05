package dao

import (
	"wxcloudrun-golang/db/model"
)

// ReferralInterface 推荐数据接口
type ReferralInterface interface {
	CreateReferral(referral *model.ReferralModel) error
	GetReferralByUserId(userId string) (*model.ReferralModel, error)
	GetReferralByPromoterCode(promoterCode string) (*model.ReferralModel, error) // 通过推广码查找用户
	GetReferralsByReferrerId(referrerId string, page, pageSize int) ([]*model.ReferralModel, int64, error)
	UpdateReferral(referral *model.ReferralModel) error

	// 佣金相关
	CreateCommission(commission *model.CommissionModel) error
	GetCommissionsByUserId(userId string, page, pageSize int) ([]*model.CommissionModel, int64, error)
	UpdateCommissionStatus(id int32, status int) error

	// 提现相关
	CreateCashout(cashout *model.CashoutModel) error
	GetCashoutsByUserId(userId string, page, pageSize int) ([]*model.CashoutModel, int64, error)
	UpdateCashoutStatus(id int32, status int) error
}

// ReferralInterfaceImp 推荐数据实现
type ReferralInterfaceImp struct{}

// ReferralImp 推荐实现实例
var ReferralImp ReferralInterface = &ReferralInterfaceImp{}
