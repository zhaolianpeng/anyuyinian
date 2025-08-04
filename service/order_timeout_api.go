package service

import (
	"encoding/json"
	"net/http"
)

// OrderTimeoutResponse 订单超时响应
type OrderTimeoutResponse struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

// CheckExpiredOrdersHandler 手动检查超时订单接口
func CheckExpiredOrdersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	// 获取订单超时服务实例
	timeoutService := GetOrderTimeoutService()
	if timeoutService == nil {
		response := &OrderTimeoutResponse{
			Code:     -1,
			ErrorMsg: "订单超时服务未初始化",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 手动检查超时订单
	timeoutService.ManualCheckExpiredOrders()

	// 获取超时订单数量
	expiredCount, err := timeoutService.GetExpiredOrdersCount()
	if err != nil {
		response := &OrderTimeoutResponse{
			Code:     -1,
			ErrorMsg: "获取超时订单数量失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := &OrderTimeoutResponse{
		Code: 0,
		Data: map[string]interface{}{
			"message":      "超时订单检查完成",
			"expiredCount": expiredCount,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetExpiredOrdersCountHandler 获取超时订单数量接口
func GetExpiredOrdersCountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	// 获取订单超时服务实例
	timeoutService := GetOrderTimeoutService()
	if timeoutService == nil {
		response := &OrderTimeoutResponse{
			Code:     -1,
			ErrorMsg: "订单超时服务未初始化",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 获取超时订单数量
	expiredCount, err := timeoutService.GetExpiredOrdersCount()
	if err != nil {
		response := &OrderTimeoutResponse{
			Code:     -1,
			ErrorMsg: "获取超时订单数量失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := &OrderTimeoutResponse{
		Code: 0,
		Data: map[string]interface{}{
			"expiredCount": expiredCount,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
