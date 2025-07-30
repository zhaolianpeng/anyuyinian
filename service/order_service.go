package service

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

// OrderResponse 订单响应
type OrderResponse struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

// SubmitOrderRequest 提交订单请求
type SubmitOrderRequest struct {
	UserId     int32                  `json:"userId"`
	ServiceId  int32                  `json:"serviceId"`
	Quantity   int                    `json:"quantity"`
	FormData   map[string]interface{} `json:"formData"`
	ReferrerId int32                  `json:"referrerId,omitempty"`
	Remark     string                 `json:"remark"`
}

// PayOrderRequest 支付订单请求
type PayOrderRequest struct {
	OrderId   int32  `json:"orderId"`
	PayMethod string `json:"payMethod"` // wechat, alipay
}

// CancelOrderRequest 取消订单请求
type CancelOrderRequest struct {
	OrderId int32  `json:"orderId"`
	Reason  string `json:"reason"`
}

// RefundOrderRequest 退款订单请求
type RefundOrderRequest struct {
	OrderId      int32   `json:"orderId"`
	RefundAmount float64 `json:"refundAmount"`
	Reason       string  `json:"reason"`
}

// OrderListRequest 订单列表请求
type OrderListRequest struct {
	UserId   int32 `json:"userId"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}

// OrderListResponse 订单列表响应
type OrderListResponse struct {
	List     []*model.OrderModel `json:"list"`
	Total    int64               `json:"total"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"pageSize"`
	HasMore  bool                `json:"hasMore"`
}

// SubmitOrderHandler 提交订单接口
func SubmitOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	var req SubmitOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "请求参数解析失败", http.StatusBadRequest)
		return
	}

	// 验证参数
	if req.UserId == 0 || req.ServiceId == 0 || req.Quantity <= 0 {
		http.Error(w, "缺少必要参数", http.StatusBadRequest)
		return
	}

	// 获取服务信息
	service, err := dao.ServiceImp.GetServiceById(req.ServiceId)
	if err != nil {
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "获取服务信息失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 生成订单号
	orderNo := generateOrderNo()

	// 计算总金额
	totalAmount := service.Price * float64(req.Quantity)

	// 计算佣金（示例：5%）
	commission := totalAmount * 0.05

	// 转换表单数据为JSON
	formDataJson, _ := json.Marshal(req.FormData)

	// 创建订单
	order := &model.OrderModel{
		OrderNo:     orderNo,
		UserId:      req.UserId,
		ServiceId:   req.ServiceId,
		ServiceName: service.Name,
		Price:       service.Price,
		Quantity:    req.Quantity,
		TotalAmount: totalAmount,
		FormData:    string(formDataJson),
		Status:      0, // 待支付
		PayStatus:   0, // 未支付
		ReferrerId:  req.ReferrerId,
		Commission:  commission,
		Remark:      req.Remark,
	}

	if err := dao.OrderImp.CreateOrder(order); err != nil {
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "创建订单失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := &OrderResponse{
		Code: 0,
		Data: map[string]interface{}{
			"orderId":     order.Id,
			"orderNo":     order.OrderNo,
			"totalAmount": order.TotalAmount,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// PayOrderHandler 支付订单接口
func PayOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	// 从URL路径中获取订单ID
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "缺少订单ID参数", http.StatusBadRequest)
		return
	}

	orderIdStr := pathParts[3]
	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		http.Error(w, "无效的订单ID", http.StatusBadRequest)
		return
	}

	var req PayOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "请求参数解析失败", http.StatusBadRequest)
		return
	}

	// 获取订单信息
	order, err := dao.OrderImp.GetOrderById(int32(orderId))
	if err != nil {
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "获取订单信息失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 检查订单状态
	if order.Status != 0 {
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "订单状态不正确",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 生成交易号
	transactionId := generateTransactionId()

	// 更新支付状态
	payTime := time.Now()
	if err := dao.OrderImp.UpdatePayStatus(int32(orderId), 1, &payTime, transactionId); err != nil {
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "更新支付状态失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 更新订单状态为已支付
	if err := dao.OrderImp.UpdateOrderStatus(int32(orderId), 1); err != nil {
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "更新订单状态失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 如果有推荐人，创建佣金记录
	if order.ReferrerId > 0 && order.Commission > 0 {
		commission := &model.CommissionModel{
			UserId:  order.ReferrerId,
			OrderId: order.Id,
			OrderNo: order.OrderNo,
			Amount:  order.Commission,
			Rate:    0.05, // 5%
			Status:  0,    // 待结算
		}
		dao.ReferralImp.CreateCommission(commission)
	}

	response := &OrderResponse{
		Code: 0,
		Data: map[string]interface{}{
			"orderId":       order.Id,
			"orderNo":       order.OrderNo,
			"transactionId": transactionId,
			"payTime":       payTime,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CancelOrderHandler 取消订单接口
func CancelOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	// 从URL路径中获取订单ID
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "缺少订单ID参数", http.StatusBadRequest)
		return
	}

	orderIdStr := pathParts[3]
	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		http.Error(w, "无效的订单ID", http.StatusBadRequest)
		return
	}

	var req CancelOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "请求参数解析失败", http.StatusBadRequest)
		return
	}

	// 获取订单信息
	order, err := dao.OrderImp.GetOrderById(int32(orderId))
	if err != nil {
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "获取订单信息失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 检查订单状态
	if order.Status != 0 {
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "订单状态不正确",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 更新订单状态为已取消
	if err := dao.OrderImp.UpdateOrderStatus(int32(orderId), 3); err != nil {
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "取消订单失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := &OrderResponse{
		Code: 0,
		Data: map[string]string{"message": "订单取消成功"},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RefundOrderHandler 退款订单接口
func RefundOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	// 从URL路径中获取订单ID
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "缺少订单ID参数", http.StatusBadRequest)
		return
	}

	orderIdStr := pathParts[3]
	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		http.Error(w, "无效的订单ID", http.StatusBadRequest)
		return
	}

	var req RefundOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "请求参数解析失败", http.StatusBadRequest)
		return
	}

	// 获取订单信息
	order, err := dao.OrderImp.GetOrderById(int32(orderId))
	if err != nil {
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "获取订单信息失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 检查订单状态
	if order.Status != 1 {
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "订单状态不正确",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 更新退款状态
	if err := dao.OrderImp.UpdateRefundStatus(int32(orderId), 1, req.RefundAmount, req.Reason); err != nil {
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "申请退款失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := &OrderResponse{
		Code: 0,
		Data: map[string]string{"message": "退款申请提交成功"},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// OrderListHandler 获取订单列表接口
func OrderListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	// 解析查询参数
	userIdStr := r.URL.Query().Get("userId")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	if userIdStr == "" {
		http.Error(w, "缺少userId参数", http.StatusBadRequest)
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "无效的用户ID", http.StatusBadRequest)
		return
	}

	// 设置默认值
	page := 1
	pageSize := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 50 {
			pageSize = ps
		}
	}

	// 获取订单列表
	orders, total, err := dao.OrderImp.GetOrdersByUserId(int32(userId), page, pageSize)
	if err != nil {
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "获取订单列表失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 计算是否有更多数据
	hasMore := int64(page*pageSize) < total

	response := &OrderResponse{
		Code: 0,
		Data: &OrderListResponse{
			List:     orders,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
			HasMore:  hasMore,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// OrderDetailHandler 获取订单详情接口
func OrderDetailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	// 从URL路径中获取订单ID
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "缺少订单ID参数", http.StatusBadRequest)
		return
	}

	orderIdStr := pathParts[3]
	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		http.Error(w, "无效的订单ID", http.StatusBadRequest)
		return
	}

	// 获取订单详情
	order, err := dao.OrderImp.GetOrderById(int32(orderId))
	if err != nil {
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "获取订单详情失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := &OrderResponse{
		Code: 0,
		Data: order,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 生成订单号
func generateOrderNo() string {
	now := time.Now()
	return fmt.Sprintf("ORDER%d%d%d", now.Year(), now.Month(), now.Day()) + fmt.Sprintf("%06d", rand.Intn(999999))
}

// 生成交易号
func generateTransactionId() string {
	now := time.Now()
	return fmt.Sprintf("TXN%d%d%d", now.Year(), now.Month(), now.Day()) + fmt.Sprintf("%06d", rand.Intn(999999))
}
