package dao

import (
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

// CashoutDao 提现数据访问实现
type CashoutDao struct{}

// CreateCashout 创建提现记录
func (c *CashoutDao) CreateCashout(cashout *model.CashoutModel) error {
	cli := db.Get()
	return cli.Table("Cashouts").Create(cashout).Error
}

// GetCashoutById 根据ID获取提现记录
func (c *CashoutDao) GetCashoutById(id int32) (*model.CashoutModel, error) {
	var cashout model.CashoutModel
	cli := db.Get()
	err := cli.Table("Cashouts").Where("id = ?", id).First(&cashout).Error
	if err != nil {
		return nil, err
	}
	return &cashout, nil
}

// GetCashoutsByUserId 根据用户ID获取提现记录列表
func (c *CashoutDao) GetCashoutsByUserId(userId string, page, pageSize int) ([]*model.CashoutModel, int64, error) {
	var cashouts []*model.CashoutModel
	var total int64

	cli := db.Get()

	// 获取总数
	err := cli.Table("Cashouts").Where("userId = ?", userId).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = cli.Table("Cashouts").Where("userId = ?", userId).
		Order("createdAt DESC").
		Offset(offset).Limit(pageSize).
		Find(&cashouts).Error

	if err != nil {
		return nil, 0, err
	}

	return cashouts, total, nil
}

// UpdateCashout 更新提现记录
func (c *CashoutDao) UpdateCashout(cashout *model.CashoutModel) error {
	cli := db.Get()
	return cli.Table("Cashouts").Where("id = ?", cashout.Id).Updates(cashout).Error
}

// DeleteCashout 删除提现记录
func (c *CashoutDao) DeleteCashout(id int32) error {
	cli := db.Get()
	return cli.Table("Cashouts").Where("id = ?", id).Delete(&model.CashoutModel{}).Error
}

// GetCashoutsByStatus 根据状态获取提现记录
func (c *CashoutDao) GetCashoutsByStatus(status int, page, pageSize int) ([]*model.CashoutModel, int64, error) {
	var cashouts []*model.CashoutModel
	var total int64

	cli := db.Get()

	// 获取总数
	err := cli.Table("Cashouts").Where("status = ?", status).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = cli.Table("Cashouts").Where("status = ?", status).
		Order("createdAt DESC").
		Offset(offset).Limit(pageSize).
		Find(&cashouts).Error

	if err != nil {
		return nil, 0, err
	}

	return cashouts, total, nil
}
