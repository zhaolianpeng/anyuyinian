package dao

import (
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

const orderTableName = "Orders"

// CreateOrder 创建订单
func (imp *OrderInterfaceImp) CreateOrder(order *model.OrderModel) error {
	cli := db.Get()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	return cli.Table(orderTableName).Create(order).Error
}

// GetOrderById 根据ID获取订单
func (imp *OrderInterfaceImp) GetOrderById(id int32) (*model.OrderModel, error) {
	var order = new(model.OrderModel)
	cli := db.Get()
	err := cli.Table(orderTableName).Where("id = ?", id).First(order).Error
	return order, err
}

// GetOrderByOrderNo 根据订单号获取订单
func (imp *OrderInterfaceImp) GetOrderByOrderNo(orderNo string) (*model.OrderModel, error) {
	var order = new(model.OrderModel)
	cli := db.Get()
	err := cli.Table(orderTableName).Where("orderNo = ?", orderNo).First(order).Error
	return order, err
}

// GetOrdersByUserId 根据用户ID获取订单列表（分页）
func (imp *OrderInterfaceImp) GetOrdersByUserId(userId int32, page, pageSize int) ([]*model.OrderModel, int64, error) {
	var orders []*model.OrderModel
	var total int64
	cli := db.Get()

	// 获取总数
	err := cli.Table(orderTableName).Where("userId = ?", userId).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = cli.Table(orderTableName).
		Where("userId = ?", userId).
		Order("createdAt DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&orders).Error

	return orders, total, err
}

// UpdateOrder 更新订单
func (imp *OrderInterfaceImp) UpdateOrder(order *model.OrderModel) error {
	cli := db.Get()
	order.UpdatedAt = time.Now()
	return cli.Table(orderTableName).Where("id = ?", order.Id).Updates(order).Error
}

// UpdateOrderStatus 更新订单状态
func (imp *OrderInterfaceImp) UpdateOrderStatus(id int32, status int) error {
	cli := db.Get()
	return cli.Table(orderTableName).Where("id = ?", id).Updates(map[string]interface{}{
		"status":    status,
		"updatedAt": time.Now(),
	}).Error
}

// UpdatePayStatus 更新支付状态
func (imp *OrderInterfaceImp) UpdatePayStatus(id int32, payStatus int, payTime *time.Time, transactionId string) error {
	cli := db.Get()
	updates := map[string]interface{}{
		"payStatus":     payStatus,
		"updatedAt":     time.Now(),
		"transactionId": transactionId,
	}
	if payTime != nil {
		updates["payTime"] = payTime
	}
	return cli.Table(orderTableName).Where("id = ?", id).Updates(updates).Error
}

// UpdateRefundStatus 更新退款状态
func (imp *OrderInterfaceImp) UpdateRefundStatus(id int32, refundStatus int, refundAmount float64, refundReason string) error {
	cli := db.Get()
	updates := map[string]interface{}{
		"refundStatus": refundStatus,
		"refundAmount": refundAmount,
		"refundReason": refundReason,
		"updatedAt":    time.Now(),
	}
	if refundStatus == 2 { // 已退款
		updates["refundTime"] = time.Now()
	}
	return cli.Table(orderTableName).Where("id = ?", id).Updates(updates).Error
}

// GetExpiredOrders 获取已超时的待支付订单
func (imp *OrderInterfaceImp) GetExpiredOrders() ([]*model.OrderModel, error) {
	var orders []*model.OrderModel
	cli := db.Get()
	now := time.Now()

	err := cli.Table(orderTableName).
		Where("status = ? AND payStatus = ? AND payDeadline < ?", 0, 0, now).
		Find(&orders).Error

	return orders, err
}

// BatchCancelExpiredOrders 批量取消超时订单
func (imp *OrderInterfaceImp) BatchCancelExpiredOrders() (int64, error) {
	cli := db.Get()
	now := time.Now()

	result := cli.Table(orderTableName).
		Where("status = ? AND payStatus = ? AND payDeadline < ?", 0, 0, now).
		Updates(map[string]interface{}{
			"status":    3, // 已取消
			"updatedAt": now,
		})

	return result.RowsAffected, result.Error
}

// GetOrdersByStatus 根据状态获取订单列表
func (imp *OrderInterfaceImp) GetOrdersByStatus(status int, page, pageSize int) ([]*model.OrderModel, int64, error) {
	var orders []*model.OrderModel
	var total int64
	cli := db.Get()

	// 获取总数
	err := cli.Table(orderTableName).Where("status = ?", status).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = cli.Table(orderTableName).
		Where("status = ?", status).
		Order("createdAt DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&orders).Error

	return orders, total, err
}
