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
	UserId           int32                  `json:"userId"`
	ServiceId        int32                  `json:"serviceId"`
	PatientId        int32                  `json:"patientId"`       // 就诊人ID
	AddressId        int32                  `json:"addressId"`       // 地址ID
	AppointmentDate  string                 `json:"appointmentDate"` // 预约日期
	AppointmentTime  string                 `json:"appointmentTime"` // 预约时间
	Quantity         int                    `json:"quantity"`
	FormData         map[string]interface{} `json:"formData"`
	ReferrerId       int32                  `json:"referrerId,omitempty"`
	Remark           string                 `json:"remark"`
	DiseaseInfo      string                 `json:"diseaseInfo"`      // 既往病史
	NeedToiletAssist string                 `json:"needToiletAssist"` // 是否需要助排二便
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
	LogInfo("开始处理提交订单请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodPost {
		LogError("请求方法不支持", fmt.Errorf("期望POST方法，实际为%s", r.Method))
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	var req SubmitOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogError("请求参数解析失败", err)
		http.Error(w, "请求参数解析失败", http.StatusBadRequest)
		return
	}

	LogStep("解析提交订单请求参数", map[string]interface{}{
		"userId":          req.UserId,
		"serviceId":       req.ServiceId,
		"patientId":       req.PatientId,
		"addressId":       req.AddressId,
		"appointmentDate": req.AppointmentDate,
		"appointmentTime": req.AppointmentTime,
		"quantity":        req.Quantity,
		"referrerId":      req.ReferrerId,
		"formDataCount":   len(req.FormData),
	})

	// 验证参数
	if req.UserId == 0 || req.ServiceId == 0 || req.PatientId == 0 || req.AddressId == 0 {
		LogError("缺少必要参数", fmt.Errorf("userId=%d, serviceId=%d, patientId=%d, addressId=%d", req.UserId, req.ServiceId, req.PatientId, req.AddressId))
		http.Error(w, "缺少必要参数", http.StatusBadRequest)
		return
	}

	// 验证预约时间
	if req.AppointmentDate == "" || req.AppointmentTime == "" {
		LogError("缺少预约时间", fmt.Errorf("appointmentDate=%s, appointmentTime=%s", req.AppointmentDate, req.AppointmentTime))
		http.Error(w, "请选择预约时间", http.StatusBadRequest)
		return
	}

	// 验证预约时间是否在允许范围内（明天开始，未来7天）
	appointmentDateTime, err := time.Parse("2006-01-02 15:04", req.AppointmentDate+" "+req.AppointmentTime)
	if err != nil {
		LogError("预约时间格式错误", err)
		http.Error(w, "预约时间格式错误", http.StatusBadRequest)
		return
	}

	tomorrow := time.Now().AddDate(0, 0, 1)
	tomorrow = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())

	maxDate := time.Now().AddDate(0, 0, 7)
	maxDate = time.Date(maxDate.Year(), maxDate.Month(), maxDate.Day(), 23, 59, 59, 0, maxDate.Location())

	if appointmentDateTime.Before(tomorrow) {
		LogError("预约时间过早", fmt.Errorf("appointmentDateTime=%v, tomorrow=%v", appointmentDateTime, tomorrow))
		http.Error(w, "预约时间不能早于明天", http.StatusBadRequest)
		return
	}

	if appointmentDateTime.After(maxDate) {
		LogError("预约时间过晚", fmt.Errorf("appointmentDateTime=%v, maxDate=%v", appointmentDateTime, maxDate))
		http.Error(w, "预约时间不能超过7天后", http.StatusBadRequest)
		return
	}

	// 验证时间槽是否在允许范围内
	allowedTimeSlots := []string{
		"08:00", "09:00", "10:00", "11:00",
		"14:00", "15:00", "16:00", "17:00", "18:00", "19:00",
	}

	isValidTimeSlot := false
	for _, slot := range allowedTimeSlots {
		if req.AppointmentTime == slot {
			isValidTimeSlot = true
			break
		}
	}

	if !isValidTimeSlot {
		LogError("预约时间不在允许的时间段内", fmt.Errorf("appointmentTime=%s, allowedSlots=%v", req.AppointmentTime, allowedTimeSlots))
		http.Error(w, "预约时间不在允许的时间段内", http.StatusBadRequest)
		return
	}

	LogStep("预约时间验证通过", map[string]interface{}{
		"appointmentDateTime": appointmentDateTime,
		"tomorrow":            tomorrow,
		"maxDate":             maxDate,
		"appointmentTime":     req.AppointmentTime,
		"allowedTimeSlots":    allowedTimeSlots,
	})

	// 获取服务信息
	LogStep("开始查询服务信息", map[string]interface{}{
		"serviceId": req.ServiceId,
	})

	service, err := dao.ServiceImp.GetServiceById(req.ServiceId)
	if err != nil {
		LogError("数据库查询服务信息失败", err)
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "获取服务信息失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	LogStep("服务信息查询成功", map[string]interface{}{
		"serviceId":   service.Id,
		"serviceName": service.Name,
		"price":       service.Price,
	})

	// 生成订单号
	orderNo := generateOrderNo()
	LogStep("生成订单号", map[string]interface{}{
		"orderNo": orderNo,
	})

	// 计算总金额
	totalAmount := service.Price * float64(req.Quantity)
	LogStep("计算订单金额", map[string]interface{}{
		"unitPrice":   service.Price,
		"quantity":    req.Quantity,
		"totalAmount": totalAmount,
	})

	// 计算佣金（示例：5%）
	commission := totalAmount * 0.05
	LogStep("计算佣金", map[string]interface{}{
		"commission": commission,
		"rate":       0.05,
	})

	// 转换表单数据为JSON
	formDataJson, _ := json.Marshal(req.FormData)
	LogStep("处理表单数据", map[string]interface{}{
		"formDataLength": len(string(formDataJson)),
	})

	// 创建订单
	LogStep("开始创建订单对象", map[string]interface{}{
		"orderNo":         orderNo,
		"userId":          req.UserId,
		"serviceId":       req.ServiceId,
		"patientId":       req.PatientId,
		"addressId":       req.AddressId,
		"appointmentDate": req.AppointmentDate,
		"appointmentTime": req.AppointmentTime,
		"serviceName":     service.Name,
		"totalAmount":     totalAmount,
	})

	// 转换助排二便字段
	needToiletAssist := 0
	if req.NeedToiletAssist == "1" {
		needToiletAssist = 1
	}

	order := &model.OrderModel{
		OrderNo:          orderNo,
		UserId:           req.UserId,
		ServiceId:        req.ServiceId,
		PatientId:        req.PatientId,
		AddressId:        req.AddressId,
		AppointmentDate:  req.AppointmentDate,
		AppointmentTime:  req.AppointmentTime,
		DiseaseInfo:      req.DiseaseInfo,
		NeedToiletAssist: needToiletAssist,
		ServiceName:      service.Name,
		Price:            service.Price,
		Quantity:         req.Quantity,
		TotalAmount:      totalAmount,
		FormData:         string(formDataJson),
		Status:           0, // 待支付
		PayStatus:        0, // 未支付
		ReferrerId:       req.ReferrerId,
		Commission:       commission,
		Remark:           req.Remark,
	}

	LogStep("开始保存订单到数据库", map[string]interface{}{
		"orderNo":     order.OrderNo,
		"userId":      order.UserId,
		"totalAmount": order.TotalAmount,
	})

	if err := dao.OrderImp.CreateOrder(order); err != nil {
		LogError("数据库创建订单失败", err)
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "创建订单失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	LogStep("订单创建成功", map[string]interface{}{
		"orderId":     order.Id,
		"orderNo":     order.OrderNo,
		"totalAmount": order.TotalAmount,
	})

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

	LogInfo("订单提交成功", map[string]interface{}{
		"orderId":     order.Id,
		"orderNo":     order.OrderNo,
		"userId":      order.UserId,
		"totalAmount": order.TotalAmount,
	})
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

	// 生成微信支付参数
	paymentParams, err := generateWechatPayParams(order, req.PayMethod)
	if err != nil {
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "生成支付参数失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := &OrderResponse{
		Code: 0,
		Data: map[string]interface{}{
			"orderId":       order.Id,
			"orderNo":       order.OrderNo,
			"totalAmount":   order.TotalAmount,
			"paymentParams": paymentParams,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// PayConfirmHandler 支付确认接口
func PayConfirmHandler(w http.ResponseWriter, r *http.Request) {
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

	var req struct {
		TransactionId string `json:"transactionId"`
		PayMethod     string `json:"payMethod"`
	}
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
	transactionId := req.TransactionId
	if transactionId == "" {
		transactionId = generateTransactionId()
	}

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

// 生成微信支付参数
func generateWechatPayParams(order *model.OrderModel, payMethod string) (map[string]interface{}, error) {
	// 这里应该调用微信支付API生成支付参数
	// 目前返回模拟数据
	return map[string]interface{}{
		"timeStamp": strconv.FormatInt(time.Now().Unix(), 10),
		"nonceStr":  generateTransactionId(),
		"package":   "prepay_id=wx" + generateTransactionId(),
		"signType":  "MD5",
		"paySign":   "mock_pay_sign_" + generateTransactionId(),
	}, nil
}

// GetAvailableTimeSlotsRequest 获取可用时间槽请求
type GetAvailableTimeSlotsRequest struct {
	Date string `json:"date"` // 日期格式：YYYY-MM-DD
}

// GetAvailableTimeSlotsResponse 获取可用时间槽响应
type GetAvailableTimeSlotsResponse struct {
	Date      string   `json:"date"`
	TimeSlots []string `json:"timeSlots"`
}

// GetAvailableTimeSlotsHandler 获取可用时间槽接口
func GetAvailableTimeSlotsHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理获取可用时间槽请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodPost {
		LogError("请求方法不支持", fmt.Errorf("期望POST方法，实际为%s", r.Method))
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	var req GetAvailableTimeSlotsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogError("请求参数解析失败", err)
		http.Error(w, "请求参数解析失败", http.StatusBadRequest)
		return
	}

	LogStep("解析获取可用时间槽请求参数", map[string]interface{}{
		"date": req.Date,
	})

	// 验证日期格式
	if req.Date == "" {
		LogError("缺少日期参数", fmt.Errorf("date=%s", req.Date))
		http.Error(w, "请提供日期参数", http.StatusBadRequest)
		return
	}

	// 解析日期
	requestDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		LogError("日期格式错误", err)
		http.Error(w, "日期格式错误，请使用YYYY-MM-DD格式", http.StatusBadRequest)
		return
	}

	// 验证日期范围（明天开始，未来7天）
	tomorrow := time.Now().AddDate(0, 0, 1)
	tomorrow = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())

	maxDate := time.Now().AddDate(0, 0, 7)
	maxDate = time.Date(maxDate.Year(), maxDate.Month(), maxDate.Day(), 23, 59, 59, 0, maxDate.Location())

	requestDateTime := time.Date(requestDate.Year(), requestDate.Month(), requestDate.Day(), 0, 0, 0, 0, requestDate.Location())

	if requestDateTime.Before(tomorrow) {
		LogError("请求日期过早", fmt.Errorf("requestDate=%v, tomorrow=%v", requestDateTime, tomorrow))
		http.Error(w, "只能查询明天开始的日期", http.StatusBadRequest)
		return
	}

	if requestDateTime.After(maxDate) {
		LogError("请求日期过晚", fmt.Errorf("requestDate=%v, maxDate=%v", requestDateTime, maxDate))
		http.Error(w, "只能查询7天内的日期", http.StatusBadRequest)
		return
	}

	// 定义允许的时间槽
	allowedTimeSlots := []string{
		"08:00", "09:00", "10:00", "11:00",
		"14:00", "15:00", "16:00", "17:00", "18:00", "19:00",
	}

	// TODO: 这里可以添加业务逻辑来检查哪些时间段已被预约
	// 例如：查询数据库中该日期的已预约时间段，然后从允许的时间槽中排除
	// 目前返回所有允许的时间槽

	LogStep("获取可用时间槽成功", map[string]interface{}{
		"date":      req.Date,
		"timeSlots": allowedTimeSlots,
	})

	response := &OrderResponse{
		Code: 0,
		Data: &GetAvailableTimeSlotsResponse{
			Date:      req.Date,
			TimeSlots: allowedTimeSlots,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	LogInfo("获取可用时间槽成功", map[string]interface{}{
		"date":      req.Date,
		"timeSlots": allowedTimeSlots,
	})
}
