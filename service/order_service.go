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
	UserId           string                 `json:"userId"`
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
	OpenID    string `json:"openId"`    // 用户openID
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
	UserId   string `json:"userId"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}

// OrderListItem 订单列表项（增强版）
type OrderListItem struct {
	Id              int32     `json:"id"`
	OrderNo         string    `json:"orderNo"`
	ServiceName     string    `json:"serviceName"`     // 服务名称
	ServiceTitle    string    `json:"serviceTitle"`    // 服务标题
	AppointmentDate string    `json:"appointmentDate"` // 预约日期
	AppointmentTime string    `json:"appointmentTime"` // 预约时间
	ConsultTime     string    `json:"consultTime"`     // 服务沟通时间（从formData中提取）
	Price           float64   `json:"price"`           // 服务单价
	TotalAmount     float64   `json:"totalAmount"`     // 订单金额
	Status          int       `json:"status"`          // 订单状态
	PayStatus       int       `json:"payStatus"`       // 支付状态
	CreatedAt       time.Time `json:"createdAt"`       // 创建时间
	StatusText      string    `json:"statusText"`      // 状态文本
	PayStatusText   string    `json:"payStatusText"`   // 支付状态文本
	FormattedAmount string    `json:"formattedAmount"` // 格式化金额
	FormattedDate   string    `json:"formattedDate"`   // 格式化日期
	Amount          float64   `json:"amount"`          // 兼容字段
}

// OrderListResponse 订单列表响应
type OrderListResponse struct {
	List     []*OrderListItem `json:"list"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"pageSize"`
	HasMore  bool             `json:"hasMore"`
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
	if req.UserId == "" || req.ServiceId == 0 || req.PatientId == 0 || req.AddressId == 0 {
		LogError("缺少必要参数", fmt.Errorf("userId=%s, serviceId=%d, patientId=%d, addressId=%d", req.UserId, req.ServiceId, req.PatientId, req.AddressId))
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

	// 设置支付截止时间（30分钟后）
	payDeadline := time.Now().Add(30 * time.Minute)

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
		Status:           0,            // 待支付
		PayStatus:        0,            // 未支付
		PayDeadline:      &payDeadline, // 支付截止时间
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
			"totalAmount": totalAmount, // 使用计算出的金额，而不是order.TotalAmount
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
	if len(pathParts) < 5 {
		http.Error(w, "缺少订单ID参数", http.StatusBadRequest)
		return
	}

	orderIdStr := pathParts[4] // 修复：使用索引4而不是3
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
	paymentParams, err := GenerateWechatPayParams(order, req.OpenID)
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
	if len(pathParts) < 5 {
		http.Error(w, "缺少订单ID参数", http.StatusBadRequest)
		return
	}

	orderIdStr := pathParts[4] // 修复：使用索引4而不是3
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
			UserId:  fmt.Sprintf("%d", order.ReferrerId), // 将int32转换为string
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
	LogInfo("开始处理取消订单请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
		"url":    r.URL.String(),
	})

	if r.Method != http.MethodPost {
		LogError("请求方法不支持", fmt.Errorf("期望POST方法，实际为%s", r.Method))
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	// 从URL路径中获取订单ID
	pathParts := strings.Split(r.URL.Path, "/")
	LogStep("路径解析", map[string]interface{}{
		"path":      r.URL.Path,
		"pathParts": pathParts,
		"length":    len(pathParts),
	})

	if len(pathParts) < 5 {
		LogError("URL路径格式错误", fmt.Errorf("路径段数不足: %d", len(pathParts)))
		http.Error(w, "缺少订单ID参数", http.StatusBadRequest)
		return
	}

	orderIdStr := pathParts[4] // 修复：使用索引4而不是3
	LogStep("解析订单ID", map[string]interface{}{
		"orderIdStr": orderIdStr,
		"pathParts":  pathParts,
	})

	// 检查orderIdStr是否为空
	if orderIdStr == "" {
		LogError("订单ID为空", fmt.Errorf("orderIdStr为空"))
		http.Error(w, "无效的订单ID", http.StatusBadRequest)
		return
	}

	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		LogError("订单ID解析失败", err)
		http.Error(w, "无效的订单ID", http.StatusBadRequest)
		return
	}

	LogStep("订单ID解析成功", map[string]interface{}{
		"orderId": orderId,
	})

	var req CancelOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogError("请求参数解析失败", err)
		http.Error(w, "请求参数解析失败", http.StatusBadRequest)
		return
	}

	LogStep("请求参数解析成功", map[string]interface{}{
		"orderId": req.OrderId,
		"reason":  req.Reason,
	})

	// 获取订单信息
	order, err := dao.OrderImp.GetOrderById(int32(orderId))
	if err != nil {
		LogError("获取订单信息失败", err)
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "获取订单信息失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	LogStep("获取订单信息成功", map[string]interface{}{
		"orderId": order.Id,
		"orderNo": order.OrderNo,
		"status":  order.Status,
		"userId":  order.UserId,
	})

	// 检查订单状态
	if order.Status != 0 {
		LogError("订单状态不正确", fmt.Errorf("期望状态0，实际状态%d", order.Status))
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "订单状态不正确，只有待支付的订单可以取消",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 更新订单状态为已取消
	if err := dao.OrderImp.UpdateOrderStatus(int32(orderId), 3); err != nil {
		LogError("更新订单状态失败", err)
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "取消订单失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	LogStep("订单取消成功", map[string]interface{}{
		"orderId": orderId,
		"orderNo": order.OrderNo,
	})

	response := &OrderResponse{
		Code: 0,
		Data: map[string]string{"message": "订单取消成功"},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RefundOrderHandler 退款订单接口
func RefundOrderHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理退款订单请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodPost {
		LogError("请求方法不支持", fmt.Errorf("期望POST方法，实际为%s", r.Method))
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	// 从URL路径中获取订单ID
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		LogError("缺少订单ID参数", fmt.Errorf("路径段数不足: %d", len(pathParts)))
		http.Error(w, "缺少订单ID参数", http.StatusBadRequest)
		return
	}

	orderIdStr := pathParts[4]
	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		LogError("无效的订单ID", fmt.Errorf("orderId=%s", orderIdStr))
		http.Error(w, "无效的订单ID", http.StatusBadRequest)
		return
	}

	var req RefundOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogError("请求参数解析失败", err)
		http.Error(w, "请求参数解析失败", http.StatusBadRequest)
		return
	}

	LogStep("解析退款请求参数", map[string]interface{}{
		"orderId":      orderId,
		"refundAmount": req.RefundAmount,
		"reason":       req.Reason,
	})

	// 验证退款金额
	if req.RefundAmount <= 0 {
		LogError("退款金额无效", fmt.Errorf("refundAmount=%f", req.RefundAmount))
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "退款金额必须大于0",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 获取订单信息
	order, err := dao.OrderImp.GetOrderById(int32(orderId))
	if err != nil {
		LogError("获取订单信息失败", err)
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "获取订单信息失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if order == nil {
		LogError("订单不存在", fmt.Errorf("orderId=%d", orderId))
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "订单不存在",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 检查订单状态
	if order.Status != 1 || order.PayStatus != 1 {
		LogError("订单状态不正确", fmt.Errorf("status=%d, payStatus=%d", order.Status, order.PayStatus))
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "只有已支付的订单可以申请退款",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 检查退款金额不能超过订单金额
	if req.RefundAmount > order.TotalAmount {
		LogError("退款金额超过订单金额", fmt.Errorf("refundAmount=%f, totalAmount=%f", req.RefundAmount, order.TotalAmount))
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "退款金额不能超过订单金额",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 检查是否已经申请过退款
	if order.RefundStatus > 0 {
		LogError("订单已申请退款", fmt.Errorf("refundStatus=%d", order.RefundStatus))
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "订单已申请退款，请勿重复申请",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 更新退款状态
	if err := dao.OrderImp.UpdateRefundStatus(int32(orderId), 1, req.RefundAmount, req.Reason); err != nil {
		LogError("申请退款失败", err)
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "申请退款失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	LogStep("退款申请成功", map[string]interface{}{
		"orderId":      orderId,
		"orderNo":      order.OrderNo,
		"refundAmount": req.RefundAmount,
		"reason":       req.Reason,
	})

	response := &OrderResponse{
		Code: 0,
		Data: map[string]interface{}{
			"orderId":      orderId,
			"orderNo":      order.OrderNo,
			"refundAmount": req.RefundAmount,
			"reason":       req.Reason,
			"message":      "退款申请提交成功",
		},
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
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	statusStr := r.URL.Query().Get("status") // 新增状态筛选参数

	// 获取用户ID参数
	userIdStr := r.URL.Query().Get("userId")
	if userIdStr == "" {
		http.Error(w, "缺少userId参数", http.StatusBadRequest)
		return
	}

	// 直接使用字符串类型的userId
	userId := userIdStr

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

	LogStep("开始状态筛选逻辑", map[string]interface{}{
		"statusStr": statusStr,
		"userId":    userId,
		"page":      page,
		"pageSize":  pageSize,
	})

	// 先检查是否有该用户的任何订单（不按状态筛选）
	allOrders, allTotal, allErr := dao.OrderImp.GetOrdersByUserId(userId, 1, 1)
	LogStep("检查用户订单总数", map[string]interface{}{
		"userId":    userId,
		"allOrders": len(allOrders),
		"allTotal":  allTotal,
		"allErr":    allErr,
	})

	// 状态筛选逻辑
	var orders []*model.OrderModel
	var total int64
	var err error

	if statusStr != "" {
		// 状态映射：支持数字和字符串状态值
		var status int
		switch statusStr {
		case "pending_pay", "0":
			status = 0
		case "paid", "1":
			status = 1
		case "cancelled", "3":
			status = 3
		case "refunded", "4":
			status = 4
		default:
			LogStep("状态值不匹配，返回空列表", map[string]interface{}{
				"statusStr": statusStr,
			})
			// 如果状态不匹配，返回空列表
			response := &OrderResponse{
				Code: 0,
				Data: &OrderListResponse{
					List:     []*OrderListItem{},
					Total:    0,
					Page:     page,
					PageSize: pageSize,
					HasMore:  false,
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		LogStep("按状态筛选订单", map[string]interface{}{
			"status":   status,
			"userId":   userId,
			"page":     page,
			"pageSize": pageSize,
		})

		// 按状态筛选订单
		orders, total, err = dao.OrderImp.GetOrdersByStatusAndUserId(status, userId, page, pageSize)
	} else {
		LogStep("获取所有订单", map[string]interface{}{
			"userId":   userId,
			"page":     page,
			"pageSize": pageSize,
		})

		// 获取所有订单
		orders, total, err = dao.OrderImp.GetOrdersByUserId(userId, page, pageSize)
	}

	LogStep("数据库查询结果", map[string]interface{}{
		"orderCount": len(orders),
		"total":      total,
		"error":      err,
	})

	if err != nil {
		LogError("数据库查询失败", err)
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "获取订单列表失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	LogStep("数据库查询成功", map[string]interface{}{
		"orderCount": len(orders),
		"total":      total,
	})

	// 如果没有订单数据，直接返回空列表
	if len(orders) == 0 {
		LogStep("没有找到订单数据", map[string]interface{}{
			"userId":   userId,
			"status":   statusStr,
			"page":     page,
			"pageSize": pageSize,
		})

		// 临时调试：尝试获取所有订单来检查数据库状态
		allOrdersDebug, allTotalDebug, allErrDebug := dao.OrderImp.GetOrdersByUserId(userId, 1, 10)
		LogStep("调试：检查用户所有订单", map[string]interface{}{
			"userId":    userId,
			"allOrders": len(allOrdersDebug),
			"allTotal":  allTotalDebug,
			"allErr":    allErrDebug,
		})

		// 如果用户有订单但状态筛选后为空，记录详细信息
		if allTotalDebug > 0 && len(allOrdersDebug) > 0 {
			LogStep("调试：用户有订单但状态筛选后为空", map[string]interface{}{
				"userId":           userId,
				"statusStr":        statusStr,
				"allOrders":        len(allOrdersDebug),
				"firstOrderStatus": allOrdersDebug[0].Status,
			})
		}

		response := &OrderResponse{
			Code: 0,
			Data: &OrderListResponse{
				List:     []*OrderListItem{},
				Total:    0,
				Page:     page,
				PageSize: pageSize,
				HasMore:  false,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 转换为增强的订单列表项
	var orderList []*OrderListItem
	for _, order := range orders {
		// 从formData中提取咨询时间
		consultTime := ""
		if order.FormData != "" {
			var formData map[string]interface{}
			if err := json.Unmarshal([]byte(order.FormData), &formData); err == nil {
				if ct, ok := formData["consultTime"].(string); ok {
					consultTime = ct
				}
			}
		}

		// 状态文本映射
		statusText := "待支付"
		switch order.Status {
		case 1:
			statusText = "已支付"
		case 2:
			statusText = "已完成"
		case 3:
			statusText = "已取消"
		case 4:
			statusText = "已退款"
		}

		// 支付状态文本映射
		payStatusText := "未支付"
		if order.PayStatus == 1 {
			payStatusText = "已支付"
		}

		// 格式化金额
		formattedAmount := fmt.Sprintf("¥%.2f", order.TotalAmount)

		// 格式化日期 - 使用 UTC 时间，让前端处理时区转换
		formattedDate := order.CreatedAt.UTC().Format("2006-01-02T15:04:05Z")

		// 添加调试日志
		LogStep("处理订单列表项", map[string]interface{}{
			"orderId":       order.Id,
			"orderNo":       order.OrderNo,
			"totalAmount":   order.TotalAmount,
			"price":         order.Price,
			"quantity":      order.Quantity,
			"createdAt":     order.CreatedAt,
			"formattedDate": formattedDate,
		})

		orderItem := &OrderListItem{
			Id:              order.Id,
			OrderNo:         order.OrderNo,
			ServiceName:     order.ServiceName,
			ServiceTitle:    order.ServiceName, // 使用服务名称作为标题
			AppointmentDate: order.AppointmentDate,
			AppointmentTime: order.AppointmentTime,
			ConsultTime:     consultTime,
			Price:           order.Price, // 添加价格字段
			TotalAmount:     order.TotalAmount,
			Amount:          order.TotalAmount, // 兼容字段
			Status:          order.Status,
			PayStatus:       order.PayStatus,
			CreatedAt:       order.CreatedAt,
			StatusText:      statusText,
			PayStatusText:   payStatusText,
			FormattedAmount: formattedAmount,
			FormattedDate:   formattedDate,
		}
		orderList = append(orderList, orderItem)
	}

	// 计算是否有更多数据
	hasMore := int64(page*pageSize) < total

	response := &OrderResponse{
		Code: 0,
		Data: &OrderListResponse{
			List:     orderList,
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
// OrderDetailRequest 订单详情请求
type OrderDetailRequest struct {
	OrderNo string `json:"orderNo"` // 订单号
}

// OrderDetailResponse 订单详情响应
type OrderDetailResponse struct {
	*model.OrderModel
	PatientName    string `json:"patientName,omitempty"`    // 患者姓名
	PatientPhone   string `json:"patientPhone,omitempty"`   // 患者电话
	AddressInfo    string `json:"addressInfo,omitempty"`    // 地址信息
	ServiceTitle   string `json:"serviceTitle,omitempty"`   // 服务标题
	FormattedPrice string `json:"formattedPrice,omitempty"` // 格式化价格
}

func OrderDetailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	var req OrderDetailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogError("请求参数解析失败", err)
		http.Error(w, "请求参数解析失败", http.StatusBadRequest)
		return
	}

	LogStep("解析订单详情请求参数", map[string]interface{}{
		"orderNo": req.OrderNo,
	})

	// 验证订单号
	if req.OrderNo == "" {
		LogError("缺少订单号", fmt.Errorf("orderNo is empty"))
		http.Error(w, "缺少订单号", http.StatusBadRequest)
		return
	}

	// 按订单号查询订单详情
	order, err := dao.OrderImp.GetOrderByOrderNo(req.OrderNo)
	if err != nil {
		LogError("获取订单详情失败", err)
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "获取订单详情失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	LogStep("订单详情查询成功", map[string]interface{}{
		"orderId":     order.Id,
		"orderNo":     order.OrderNo,
		"serviceName": order.ServiceName,
		"totalAmount": order.TotalAmount,
	})

	// 构建增强的订单详情响应
	detailResponse := &OrderDetailResponse{
		OrderModel: order,
	}

	// 获取患者信息
	if order.PatientId > 0 {
		patient, err := dao.UserExtendImp.GetPatientById(order.PatientId)
		if err == nil && patient != nil {
			detailResponse.PatientName = patient.Name
			detailResponse.PatientPhone = patient.Phone
		} else {
			LogError("获取患者信息失败", err)
		}
	}

	// 获取地址信息
	if order.AddressId > 0 {
		address, err := dao.UserExtendImp.GetAddressById(order.AddressId)
		if err == nil && address != nil {
			detailResponse.AddressInfo = address.Province + address.City + address.District + address.Address
		} else {
			LogError("获取地址信息失败", err)
		}
	}

	// 格式化价格
	detailResponse.FormattedPrice = fmt.Sprintf("%.2f", order.Price)
	detailResponse.ServiceTitle = order.ServiceName

	LogStep("订单详情增强信息", map[string]interface{}{
		"patientName":    detailResponse.PatientName,
		"patientPhone":   detailResponse.PatientPhone,
		"addressInfo":    detailResponse.AddressInfo,
		"formattedPrice": detailResponse.FormattedPrice,
	})

	response := &OrderResponse{
		Code: 0,
		Data: detailResponse,
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
	// 从请求中获取openID，这里需要修改PayOrderHandler来传递openID
	// 暂时使用模拟的openID
	openID := "mock_openid_" + order.UserId

	// 调用微信支付服务生成支付参数
	return GenerateWechatPayParams(order, openID)
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
