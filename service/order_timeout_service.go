package service

import (
	"log"
	"time"
	"wxcloudrun-golang/db/dao"
)

// OrderTimeoutService 订单超时处理服务
type OrderTimeoutService struct {
	ticker *time.Ticker
	done   chan bool
}

// NewOrderTimeoutService 创建订单超时处理服务
func NewOrderTimeoutService() *OrderTimeoutService {
	return &OrderTimeoutService{
		done: make(chan bool),
	}
}

// Start 启动订单超时处理服务
func (s *OrderTimeoutService) Start() {
	// 每分钟检查一次超时订单
	s.ticker = time.NewTicker(1 * time.Minute)

	log.Println("订单超时处理服务已启动")

	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.checkAndCancelExpiredOrders()
			case <-s.done:
				return
			}
		}
	}()
}

// Stop 停止订单超时处理服务
func (s *OrderTimeoutService) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	close(s.done)
	log.Println("订单超时处理服务已停止")
}

// checkAndCancelExpiredOrders 检查并取消超时订单
func (s *OrderTimeoutService) checkAndCancelExpiredOrders() {
	log.Println("开始检查超时订单...")

	// 获取超时订单
	expiredOrders, err := dao.OrderImp.GetExpiredOrders()
	if err != nil {
		log.Printf("获取超时订单失败: %v", err)
		return
	}

	if len(expiredOrders) == 0 {
		log.Println("没有发现超时订单")
		return
	}

	log.Printf("发现 %d 个超时订单", len(expiredOrders))

	// 批量取消超时订单
	affectedRows, err := dao.OrderImp.BatchCancelExpiredOrders()
	if err != nil {
		log.Printf("批量取消超时订单失败: %v", err)
		return
	}

	log.Printf("成功取消 %d 个超时订单", affectedRows)

	// 记录取消的订单详情
	for _, order := range expiredOrders {
		log.Printf("订单 %s 因超时未支付已自动取消", order.OrderNo)
	}
}

// ManualCheckExpiredOrders 手动检查超时订单（用于测试）
func (s *OrderTimeoutService) ManualCheckExpiredOrders() {
	log.Println("手动检查超时订单...")
	s.checkAndCancelExpiredOrders()
}

// GetExpiredOrdersCount 获取超时订单数量
func (s *OrderTimeoutService) GetExpiredOrdersCount() (int, error) {
	expiredOrders, err := dao.OrderImp.GetExpiredOrders()
	if err != nil {
		return 0, err
	}
	return len(expiredOrders), nil
}

// 全局订单超时服务实例
var orderTimeoutService *OrderTimeoutService

// InitOrderTimeoutService 初始化订单超时服务
func InitOrderTimeoutService() {
	orderTimeoutService = NewOrderTimeoutService()
	orderTimeoutService.Start()
}

// StopOrderTimeoutService 停止订单超时服务
func StopOrderTimeoutService() {
	if orderTimeoutService != nil {
		orderTimeoutService.Stop()
	}
}

// GetOrderTimeoutService 获取订单超时服务实例
func GetOrderTimeoutService() *OrderTimeoutService {
	return orderTimeoutService
}
