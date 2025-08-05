package dao

import (
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

const referralTableName = "Referrals"
const commissionTableName = "Commissions"
const cashoutTableName = "Cashouts"

// 推荐关系相关方法

// CreateReferral 创建推荐关系
func (imp *ReferralInterfaceImp) CreateReferral(referral *model.ReferralModel) error {
	cli := db.Get()
	referral.CreatedAt = time.Now()
	referral.UpdatedAt = time.Now()
	return cli.Table(referralTableName).Create(referral).Error
}

// GetReferralByUserId 根据用户ID获取推荐关系
func (imp *ReferralInterfaceImp) GetReferralByUserId(userId string) (*model.ReferralModel, error) {
	var referral = new(model.ReferralModel)
	cli := db.Get()
	err := cli.Table(referralTableName).Where("userId = ? AND status = ?", userId, 1).First(referral).Error
	return referral, err
}

// GetReferralsByReferrerId 根据推荐人ID获取推荐列表（分页）
func (imp *ReferralInterfaceImp) GetReferralsByReferrerId(referrerId string, page, pageSize int) ([]*model.ReferralModel, int64, error) {
	var referrals []*model.ReferralModel
	var total int64
	cli := db.Get()

	// 获取总数
	err := cli.Table(referralTableName).Where("referrerId = ? AND status = ?", referrerId, 1).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = cli.Table(referralTableName).
		Where("referrerId = ? AND status = ?", referrerId, 1).
		Order("createdAt DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&referrals).Error

	return referrals, total, err
}

// UpdateReferral 更新推荐关系
func (imp *ReferralInterfaceImp) UpdateReferral(referral *model.ReferralModel) error {
	cli := db.Get()
	referral.UpdatedAt = time.Now()
	return cli.Table(referralTableName).Where("id = ?", referral.Id).Updates(referral).Error
}

// 佣金相关方法

// CreateCommission 创建佣金记录
func (imp *ReferralInterfaceImp) CreateCommission(commission *model.CommissionModel) error {
	cli := db.Get()
	commission.CreatedAt = time.Now()
	commission.UpdatedAt = time.Now()
	return cli.Table(commissionTableName).Create(commission).Error
}

// GetCommissionsByUserId 根据用户ID获取佣金记录（分页）
func (imp *ReferralInterfaceImp) GetCommissionsByUserId(userId string, page, pageSize int) ([]*model.CommissionModel, int64, error) {
	var commissions []*model.CommissionModel
	var total int64
	cli := db.Get()

	// 获取总数
	err := cli.Table(commissionTableName).Where("userId = ?", userId).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = cli.Table(commissionTableName).
		Where("userId = ?", userId).
		Order("createdAt DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&commissions).Error

	return commissions, total, err
}

// UpdateCommissionStatus 更新佣金状态
func (imp *ReferralInterfaceImp) UpdateCommissionStatus(id int32, status int) error {
	cli := db.Get()
	updates := map[string]interface{}{
		"status":    status,
		"updatedAt": time.Now(),
	}
	if status == 2 { // 已提现
		updates["cashoutTime"] = time.Now()
	}
	return cli.Table(commissionTableName).Where("id = ?", id).Updates(updates).Error
}

// 提现相关方法

// CreateCashout 创建提现记录
func (imp *ReferralInterfaceImp) CreateCashout(cashout *model.CashoutModel) error {
	cli := db.Get()
	cashout.CreatedAt = time.Now()
	cashout.UpdatedAt = time.Now()
	return cli.Table(cashoutTableName).Create(cashout).Error
}

// GetCashoutsByUserId 根据用户ID获取提现记录（分页）
func (imp *ReferralInterfaceImp) GetCashoutsByUserId(userId string, page, pageSize int) ([]*model.CashoutModel, int64, error) {
	var cashouts []*model.CashoutModel
	var total int64
	cli := db.Get()

	// 获取总数
	err := cli.Table(cashoutTableName).Where("userId = ?", userId).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = cli.Table(cashoutTableName).
		Where("userId = ?", userId).
		Order("createdAt DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&cashouts).Error

	return cashouts, total, err
}

// UpdateCashoutStatus 更新提现状态
func (imp *ReferralInterfaceImp) UpdateCashoutStatus(id int32, status int) error {
	cli := db.Get()
	updates := map[string]interface{}{
		"status":    status,
		"updatedAt": time.Now(),
	}
	if status == 1 || status == 2 { // 已通过或已拒绝
		updates["processTime"] = time.Now()
	}
	return cli.Table(cashoutTableName).Where("id = ?", id).Updates(updates).Error
}
