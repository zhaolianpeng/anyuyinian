package dao

import "wxcloudrun-golang/db/model"

// CashoutInterface 提现数据访问接口
type CashoutInterface interface {
	// CreateCashout 创建提现记录
	CreateCashout(cashout *model.CashoutModel) error

	// GetCashoutById 根据ID获取提现记录
	GetCashoutById(id int32) (*model.CashoutModel, error)

	// GetCashoutsByUserId 根据用户ID获取提现记录列表
	GetCashoutsByUserId(userId string, page, pageSize int) ([]*model.CashoutModel, int64, error)

	// UpdateCashout 更新提现记录
	UpdateCashout(cashout *model.CashoutModel) error

	// DeleteCashout 删除提现记录
	DeleteCashout(id int32) error

	// GetCashoutsByStatus 根据状态获取提现记录
	GetCashoutsByStatus(status int, page, pageSize int) ([]*model.CashoutModel, int64, error)
}

// CashoutInterfaceImp 提现数据访问实现
type CashoutInterfaceImp struct{}

// CreateCashout 创建提现记录
func (c *CashoutInterfaceImp) CreateCashout(cashout *model.CashoutModel) error {
	return (&CashoutDao{}).CreateCashout(cashout)
}

// GetCashoutById 根据ID获取提现记录
func (c *CashoutInterfaceImp) GetCashoutById(id int32) (*model.CashoutModel, error) {
	return (&CashoutDao{}).GetCashoutById(id)
}

// GetCashoutsByUserId 根据用户ID获取提现记录列表
func (c *CashoutInterfaceImp) GetCashoutsByUserId(userId string, page, pageSize int) ([]*model.CashoutModel, int64, error) {
	return (&CashoutDao{}).GetCashoutsByUserId(userId, page, pageSize)
}

// UpdateCashout 更新提现记录
func (c *CashoutInterfaceImp) UpdateCashout(cashout *model.CashoutModel) error {
	return (&CashoutDao{}).UpdateCashout(cashout)
}

// DeleteCashout 删除提现记录
func (c *CashoutInterfaceImp) DeleteCashout(id int32) error {
	return (&CashoutDao{}).DeleteCashout(id)
}

// GetCashoutsByStatus 根据状态获取提现记录
func (c *CashoutInterfaceImp) GetCashoutsByStatus(status int, page, pageSize int) ([]*model.CashoutModel, int64, error) {
	return (&CashoutDao{}).GetCashoutsByStatus(status, page, pageSize)
}

// Imp 实现实例
var CashoutImp CashoutInterface = &CashoutInterfaceImp{}
