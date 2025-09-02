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

	// 记录SQL操作日志
	logger := NewSQLLogger("插入", orderTableName, map[string]interface{}{
		"orderNo": order.OrderNo,
		"userId":  order.UserId,
		"status":  order.Status,
	})

	err := cli.Table(orderTableName).Create(order).Error
	logger.LogInsert(order, err)

	return err
}

// GetOrderById 根据ID获取订单
func (imp *OrderInterfaceImp) GetOrderById(id int32) (*model.OrderModel, error) {
	var order = new(model.OrderModel)
	cli := db.Get()

	// 记录SQL操作日志
	logger := NewSQLLogger("查询", orderTableName, map[string]interface{}{
		"id": id,
	})

	err := cli.Table(orderTableName).Where("id = ?", id).First(order).Error
	logger.LogQuery(order, err)

	return order, err
}

// GetOrderByOrderNo 根据订单号获取订单
func (imp *OrderInterfaceImp) GetOrderByOrderNo(orderNo string) (*model.OrderModel, error) {
	var order = new(model.OrderModel)
	cli := db.Get()

	// 记录SQL操作日志
	logger := NewSQLLogger("查询", orderTableName, map[string]interface{}{
		"orderNo": orderNo,
	})

	err := cli.Table(orderTableName).Where("orderNo = ?", orderNo).First(order).Error
	logger.LogQuery(order, err)

	return order, err
}

// GetOrdersByUserId 根据用户ID获取订单列表（分页）
func (imp *OrderInterfaceImp) GetOrdersByUserId(userId string, page, pageSize int) ([]*model.OrderModel, int64, error) {
	var orders []*model.OrderModel
	var total int64
	cli := db.Get()

	// 记录SQL操作日志
	logger := NewSQLLogger("查询", orderTableName, map[string]interface{}{
		"userId":   userId,
		"page":     page,
		"pageSize": pageSize,
	})

	// 获取总数
	err := cli.Table(orderTableName).Where("userId = ?", userId).Count(&total).Error
	if err != nil {
		logger.LogCount(total, err)
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

	logger.LogQuery(orders, err)
	return orders, total, err
}

// UpdateOrder 更新订单
func (imp *OrderInterfaceImp) UpdateOrder(order *model.OrderModel) error {
	cli := db.Get()
	order.UpdatedAt = time.Now()

	// 记录SQL操作日志
	logger := NewSQLLogger("更新", orderTableName, map[string]interface{}{
		"id":      order.Id,
		"orderNo": order.OrderNo,
	})

	err := cli.Table(orderTableName).Where("id = ?", order.Id).Updates(order).Error
	logger.LogUpdate(order, err)

	return err
}

// UpdateOrderStatus 更新订单状态
func (imp *OrderInterfaceImp) UpdateOrderStatus(id int32, status int) error {
	cli := db.Get()

	// 记录SQL操作日志
	logger := NewSQLLogger("更新", orderTableName, map[string]interface{}{
		"id":     id,
		"status": status,
	})

	updates := map[string]interface{}{
		"status":    status,
		"updatedAt": time.Now(),
	}

	err := cli.Table(orderTableName).Where("id = ?", id).Updates(updates).Error
	logger.LogUpdate(updates, err)

	return err
}

// UpdatePayStatus 更新支付状态
func (imp *OrderInterfaceImp) UpdatePayStatus(id int32, payStatus int, payTime *time.Time, transactionId string) error {
	cli := db.Get()

	// 记录SQL操作日志
	logger := NewSQLLogger("更新", orderTableName, map[string]interface{}{
		"id":            id,
		"payStatus":     payStatus,
		"transactionId": transactionId,
	})

	updates := map[string]interface{}{
		"payStatus":     payStatus,
		"updatedAt":     time.Now(),
		"transactionId": transactionId,
	}
	if payTime != nil {
		updates["payTime"] = payTime
	}

	err := cli.Table(orderTableName).Where("id = ?", id).Updates(updates).Error
	logger.LogUpdate(updates, err)

	return err
}

// UpdateRefundStatus 更新退款状态
func (imp *OrderInterfaceImp) UpdateRefundStatus(id int32, refundStatus int, refundAmount float64, refundReason string) error {
	cli := db.Get()

	// 记录SQL操作日志
	logger := NewSQLLogger("更新", orderTableName, map[string]interface{}{
		"id":            id,
		"refundStatus":  refundStatus,
		"refundAmount":  refundAmount,
		"refundReason":  refundReason,
	})

	updates := map[string]interface{}{
		"refundStatus": refundStatus,
		"refundAmount": refundAmount,
		"refundReason": refundReason,
		"updatedAt":    time.Now(),
	}
	if refundStatus == 2 { // 已退款
		updates["refundTime"] = time.Now()
	}

	err := cli.Table(orderTableName).Where("id = ?", id).Updates(updates).Error
	logger.LogUpdate(updates, err)

	return err
}

// UpdateOrderAmount 更新订单金额
func (imp *OrderInterfaceImp) UpdateOrderAmount(id int32, newAmount float64) error {
	cli := db.Get()

	// 记录SQL操作日志
	logger := NewSQLLogger("更新", orderTableName, map[string]interface{}{
		"id":        id,
		"newAmount": newAmount,
	})

	updates := map[string]interface{}{
		"totalAmount": newAmount,
		"updatedAt":   time.Now(),
	}

	err := cli.Table(orderTableName).Where("id = ?", id).Updates(updates).Error
	logger.LogUpdate(updates, err)

	return err
}

// GetExpiredOrders 获取已超时的待支付订单
func (imp *OrderInterfaceImp) GetExpiredOrders() ([]*model.OrderModel, error) {
	var orders []*model.OrderModel
	cli := db.Get()
	now := time.Now()

	// 记录SQL操作日志
	logger := NewSQLLogger("查询", orderTableName, map[string]interface{}{
		"operation": "GetExpiredOrders",
		"now":       now,
	})

	err := cli.Table(orderTableName).
		Where("status = ? AND payStatus = ? AND payDeadline < ?", 0, 0, now).
		Find(&orders).Error

	logger.LogQuery(orders, err)
	return orders, err
}

// BatchCancelExpiredOrders 批量取消超时订单
func (imp *OrderInterfaceImp) BatchCancelExpiredOrders() (int64, error) {
	cli := db.Get()
	now := time.Now()

	// 记录SQL操作日志
	logger := NewSQLLogger("更新", orderTableName, map[string]interface{}{
		"operation": "BatchCancelExpiredOrders",
		"now":       now,
	})

	updates := map[string]interface{}{
		"status":    3, // 已取消
		"updatedAt": now,
	}

	result := cli.Table(orderTableName).
		Where("status = ? AND payStatus = ? AND payDeadline < ?", 0, 0, now).
		Updates(updates)

	logger.LogUpdate(updates, result.Error)
	return result.RowsAffected, result.Error
}

// GetOrdersByStatus 根据状态获取订单列表
func (imp *OrderInterfaceImp) GetOrdersByStatus(status int, page, pageSize int) ([]*model.OrderModel, int64, error) {
	var orders []*model.OrderModel
	var total int64
	cli := db.Get()

	// 记录SQL操作日志
	logger := NewSQLLogger("查询", orderTableName, map[string]interface{}{
		"status":   status,
		"page":     page,
		"pageSize": pageSize,
	})

	// 获取总数
	err := cli.Table(orderTableName).Where("status = ?", status).Count(&total).Error
	if err != nil {
		logger.LogCount(total, err)
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

	logger.LogQuery(orders, err)
	return orders, total, err
}

// GetOrdersByStatusAndUserId 根据状态和用户ID获取订单列表（分页）
func (imp *OrderInterfaceImp) GetOrdersByStatusAndUserId(status int, userId string, page, pageSize int) ([]*model.OrderModel, int64, error) {
	var orders []*model.OrderModel
	var total int64
	cli := db.Get()

	// 记录SQL操作日志
	logger := NewSQLLogger("查询", orderTableName, map[string]interface{}{
		"status":   status,
		"userId":   userId,
		"page":     page,
		"pageSize": pageSize,
	})

	// 获取总数
	err := cli.Table(orderTableName).Where("status = ? AND userId = ?", status, userId).Count(&total).Error
	if err != nil {
		logger.LogCount(total, err)
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = cli.Table(orderTableName).
		Where("status = ? AND userId = ?", status, userId).
		Order("createdAt DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&orders).Error

	logger.LogQuery(orders, err)
	return orders, total, err
}
