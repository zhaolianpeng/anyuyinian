package dao

import (
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

// CommissionDao 佣金数据访问实现
type CommissionDao struct{}

// CreateCommission 创建佣金记录
func (c *CommissionDao) CreateCommission(commission *model.CommissionModel) error {
	cli := db.Get()
	return cli.Table("Commissions").Create(commission).Error
}

// GetCommissionById 根据ID获取佣金记录
func (c *CommissionDao) GetCommissionById(id int32) (*model.CommissionModel, error) {
	var commission model.CommissionModel
	cli := db.Get()
	err := cli.Table("Commissions").Where("id = ?", id).First(&commission).Error
	if err != nil {
		return nil, err
	}
	return &commission, nil
}

// GetCommissionsByUserId 根据用户ID获取佣金记录列表
func (c *CommissionDao) GetCommissionsByUserId(userId string, page, pageSize int) ([]*model.CommissionModel, int64, error) {
	var commissions []*model.CommissionModel
	var total int64

	cli := db.Get()

	// 获取总数
	err := cli.Table("Commissions").Where("userId = ?", userId).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = cli.Table("Commissions").Where("userId = ?", userId).
		Order("createdAt DESC").
		Offset(offset).Limit(pageSize).
		Find(&commissions).Error

	if err != nil {
		return nil, 0, err
	}

	return commissions, total, nil
}

// GetCommissionsByOrderId 根据订单ID获取佣金记录
func (c *CommissionDao) GetCommissionsByOrderId(orderId int32) (*model.CommissionModel, error) {
	var commission model.CommissionModel
	cli := db.Get()
	err := cli.Table("Commissions").Where("orderId = ?", orderId).First(&commission).Error
	if err != nil {
		return nil, err
	}
	return &commission, nil
}

// UpdateCommission 更新佣金记录
func (c *CommissionDao) UpdateCommission(commission *model.CommissionModel) error {
	cli := db.Get()
	return cli.Table("Commissions").Where("id = ?", commission.Id).Updates(commission).Error
}

// DeleteCommission 删除佣金记录
func (c *CommissionDao) DeleteCommission(id int32) error {
	cli := db.Get()
	return cli.Table("Commissions").Where("id = ?", id).Delete(&model.CommissionModel{}).Error
}

// GetCommissionsByStatus 根据状态获取佣金记录
func (c *CommissionDao) GetCommissionsByStatus(status int, page, pageSize int) ([]*model.CommissionModel, int64, error) {
	var commissions []*model.CommissionModel
	var total int64

	cli := db.Get()

	// 获取总数
	err := cli.Table("Commissions").Where("status = ?", status).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = cli.Table("Commissions").Where("status = ?", status).
		Order("createdAt DESC").
		Offset(offset).Limit(pageSize).
		Find(&commissions).Error

	if err != nil {
		return nil, 0, err
	}

	return commissions, total, nil
}
