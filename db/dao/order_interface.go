package dao

import (
	"time"
	"wxcloudrun-golang/db/model"
)

// OrderInterface 订单数据接口
type OrderInterface interface {
	CreateOrder(order *model.OrderModel) error
	GetOrderById(id int32) (*model.OrderModel, error)
	GetOrderByOrderNo(orderNo string) (*model.OrderModel, error)
	GetOrdersByUserId(userId int32, page, pageSize int) ([]*model.OrderModel, int64, error)
	UpdateOrder(order *model.OrderModel) error
	UpdateOrderStatus(id int32, status int) error
	UpdatePayStatus(id int32, payStatus int, payTime *time.Time, transactionId string) error
	UpdateRefundStatus(id int32, refundStatus int, refundAmount float64, refundReason string) error
	GetExpiredOrders() ([]*model.OrderModel, error)
	BatchCancelExpiredOrders() (int64, error)
	GetOrdersByStatus(status int, page, pageSize int) ([]*model.OrderModel, int64, error)
	GetOrdersByStatusAndUserId(status int, userId int32, page, pageSize int) ([]*model.OrderModel, int64, error)
}

// OrderInterfaceImp 订单数据实现
type OrderInterfaceImp struct{}

// OrderImp 订单实现实例
var OrderImp OrderInterface = &OrderInterfaceImp{}
