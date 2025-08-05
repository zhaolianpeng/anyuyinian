package dao

import "wxcloudrun-golang/db/model"

// CommissionInterface 佣金数据访问接口
type CommissionInterface interface {
	// CreateCommission 创建佣金记录
	CreateCommission(commission *model.CommissionModel) error

	// GetCommissionById 根据ID获取佣金记录
	GetCommissionById(id int32) (*model.CommissionModel, error)

	// GetCommissionsByUserId 根据用户ID获取佣金记录列表
	GetCommissionsByUserId(userId string, page, pageSize int) ([]*model.CommissionModel, int64, error)

	// GetCommissionsByOrderId 根据订单ID获取佣金记录
	GetCommissionsByOrderId(orderId int32) (*model.CommissionModel, error)

	// UpdateCommission 更新佣金记录
	UpdateCommission(commission *model.CommissionModel) error

	// DeleteCommission 删除佣金记录
	DeleteCommission(id int32) error

	// GetCommissionsByStatus 根据状态获取佣金记录
	GetCommissionsByStatus(status int, page, pageSize int) ([]*model.CommissionModel, int64, error)
}

// CommissionInterfaceImp 佣金数据访问实现
type CommissionInterfaceImp struct{}

// CreateCommission 创建佣金记录
func (c *CommissionInterfaceImp) CreateCommission(commission *model.CommissionModel) error {
	return (&CommissionDao{}).CreateCommission(commission)
}

// GetCommissionById 根据ID获取佣金记录
func (c *CommissionInterfaceImp) GetCommissionById(id int32) (*model.CommissionModel, error) {
	return (&CommissionDao{}).GetCommissionById(id)
}

// GetCommissionsByUserId 根据用户ID获取佣金记录列表
func (c *CommissionInterfaceImp) GetCommissionsByUserId(userId string, page, pageSize int) ([]*model.CommissionModel, int64, error) {
	return (&CommissionDao{}).GetCommissionsByUserId(userId, page, pageSize)
}

// GetCommissionsByOrderId 根据订单ID获取佣金记录
func (c *CommissionInterfaceImp) GetCommissionsByOrderId(orderId int32) (*model.CommissionModel, error) {
	return (&CommissionDao{}).GetCommissionsByOrderId(orderId)
}

// UpdateCommission 更新佣金记录
func (c *CommissionInterfaceImp) UpdateCommission(commission *model.CommissionModel) error {
	return (&CommissionDao{}).UpdateCommission(commission)
}

// DeleteCommission 删除佣金记录
func (c *CommissionInterfaceImp) DeleteCommission(id int32) error {
	return (&CommissionDao{}).DeleteCommission(id)
}

// GetCommissionsByStatus 根据状态获取佣金记录
func (c *CommissionInterfaceImp) GetCommissionsByStatus(status int, page, pageSize int) ([]*model.CommissionModel, int64, error) {
	return (&CommissionDao{}).GetCommissionsByStatus(status, page, pageSize)
}

// Imp 实现实例
var CommissionImp CommissionInterface = &CommissionInterfaceImp{}
