package service

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

// SmartElderlyOrderRequest 智慧养老设备订单请求
type SmartElderlyOrderRequest struct {
	UserId    string                 `json:"userId"`
	ServiceId int32                  `json:"serviceId"`
	AddressId int32                  `json:"addressId"`
	Quantity  int                    `json:"quantity"`
	FormData  map[string]interface{} `json:"formData"`
	Remark    string                 `json:"remark"`
}

// SmartElderlyOrderHandler 智慧养老设备订单提交接口
func SmartElderlyOrderHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理智慧养老设备订单请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodPost {
		LogError("请求方法不支持", fmt.Errorf("期望POST方法，实际为%s", r.Method))
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	var req SmartElderlyOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogError("请求参数解析失败", err)
		http.Error(w, "请求参数解析失败", http.StatusBadRequest)
		return
	}

	LogStep("解析智慧养老设备订单请求参数", map[string]interface{}{
		"userId":        req.UserId,
		"serviceId":     req.ServiceId,
		"addressId":     req.AddressId,
		"quantity":      req.Quantity,
		"formDataCount": len(req.FormData),
	})

	// 验证参数
	if req.UserId == "" || req.ServiceId == 0 || req.AddressId == 0 {
		LogError("缺少必要参数", fmt.Errorf("userId=%s, serviceId=%d, addressId=%d", req.UserId, req.ServiceId, req.AddressId))
		http.Error(w, "缺少必要参数", http.StatusBadRequest)
		return
	}

	// 获取服务信息
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

	// 验证是否为智慧养老设备
	if service.Category != "智慧养老" {
		LogError("服务类型不是智慧养老设备", fmt.Errorf("serviceCategory=%s", service.Category))
		http.Error(w, "该服务不是智慧养老设备", http.StatusBadRequest)
		return
	}

	LogStep("服务信息验证通过", map[string]interface{}{
		"serviceId":   service.Id,
		"serviceName": service.Name,
		"price":       service.Price,
		"category":    service.Category,
	})

	// 生成订单号
	orderNo := generateSmartElderlyOrderNo()

	// 计算总金额
	totalAmount := service.Price * float64(req.Quantity)

	// 创建订单
	order := &model.OrderModel{
		OrderNo:         orderNo,
		UserId:          req.UserId,
		ServiceId:       req.ServiceId,
		PatientId:       nil, // 智慧养老设备不需要患者信息
		AddressId:       req.AddressId,
		AppointmentDate: nil, // 智慧养老设备不需要预约时间
		AppointmentTime: nil, // 智慧养老设备不需要预约时间
		ServiceName:     service.Name,
		Price:           service.Price,
		Quantity:        req.Quantity,
		TotalAmount:     totalAmount,
		Status:          0, // 待支付
		PayStatus:       0, // 未支付
		Remark:          req.Remark,
	}

	// 处理表单数据
	if req.FormData != nil {
		formDataJSON, err := json.Marshal(req.FormData)
		if err != nil {
			LogError("表单数据序列化失败", err)
			http.Error(w, "表单数据处理失败", http.StatusInternalServerError)
			return
		}
		order.FormData = string(formDataJSON)
	}

	// 保存订单到数据库
	err = dao.OrderImp.CreateOrder(order)
	if err != nil {
		LogError("保存订单失败", err)
		response := &OrderResponse{
			Code:     -1,
			ErrorMsg: "订单创建失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	LogStep("智慧养老设备订单创建成功", map[string]interface{}{
		"orderNo":     orderNo,
		"totalAmount": totalAmount,
	})

	// 返回成功响应
	response := &OrderResponse{
		Code: 0,
		Data: map[string]interface{}{
			"orderNo":     orderNo,
			"totalAmount": totalAmount,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 生成智慧养老设备订单号
func generateSmartElderlyOrderNo() string {
	now := time.Now()
	dateStr := now.Format("20060102")
	timeStr := now.Format("150405")
	randomNum := rand.Intn(1000)
	return fmt.Sprintf("SE%s%s%03d", dateStr, timeStr, randomNum)
}
