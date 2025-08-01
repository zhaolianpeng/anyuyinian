package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

// ServiceResponse 服务响应
type ServiceResponse struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

// ServiceListRequest 服务列表请求
type ServiceListRequest struct {
	Category string `json:"category"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}

// ServiceListResponse 服务列表响应
type ServiceListResponse struct {
	List     []*model.ServiceItemModel `json:"list"`
	Total    int64                     `json:"total"`
	Page     int                       `json:"page"`
	PageSize int                       `json:"pageSize"`
	HasMore  bool                      `json:"hasMore"`
}

// FormConfig 表单配置
type FormConfig struct {
	Fields []FormField `json:"fields"`
}

// FormField 表单字段
type FormField struct {
	Name        string       `json:"name"`
	Label       string       `json:"label"`
	Type        string       `json:"type"` // text, textarea, select, radio, checkbox, date, file
	Required    bool         `json:"required"`
	Placeholder string       `json:"placeholder"`
	Options     []FormOption `json:"options,omitempty"`
	Validation  string       `json:"validation,omitempty"`
}

// FormOption 表单选项
type FormOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// ServiceListHandler 获取服务列表接口
func ServiceListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	// 解析查询参数
	category := r.URL.Query().Get("category")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

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

	var services []*model.ServiceItemModel
	var total int64
	var err error

	// 根据分类获取服务列表
	if category != "" {
		services, total, err = dao.ServiceImp.GetServicesByCategory(category, page, pageSize)
	} else {
		services, total, err = dao.ServiceImp.GetAllServices(page, pageSize)
	}

	if err != nil {
		response := &ServiceResponse{
			Code:     -1,
			ErrorMsg: "获取服务列表失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 计算是否有更多数据
	hasMore := int64(page*pageSize) < total

	response := &ServiceResponse{
		Code: 0,
		Data: &ServiceListResponse{
			List:     services,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
			HasMore:  hasMore,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ServiceDetailRequest 服务详情请求
type ServiceDetailRequest struct {
	ServiceId int32 `json:"serviceId"`
}

// ServiceDetailHandler 获取服务详情接口
func ServiceDetailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	var req ServiceDetailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := &ServiceResponse{
			Code:     -1,
			ErrorMsg: "请求参数解析失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 验证服务ID
	if req.ServiceId <= 0 {
		response := &ServiceResponse{
			Code:     -1,
			ErrorMsg: "无效的服务ID",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 获取服务详情
	service, err := dao.ServiceImp.GetServiceById(req.ServiceId)
	if err != nil {
		// 检查是否是记录不存在错误
		if strings.Contains(err.Error(), "record not found") || strings.Contains(err.Error(), "no rows") {
			response := &ServiceResponse{
				Code:     -1,
				ErrorMsg: fmt.Sprintf("服务ID %d 不存在，请检查服务ID是否正确", req.ServiceId),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		
		response := &ServiceResponse{
			Code:     -1,
			ErrorMsg: "获取服务详情失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 检查服务状态
	if service.Status == 0 {
		response := &ServiceResponse{
			Code:     -1,
			ErrorMsg: fmt.Sprintf("服务ID %d 已下架，暂不可用", req.ServiceId),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := &ServiceResponse{
		Code: 0,
		Data: service,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ServiceFormConfigHandler 获取服务表单配置接口
func ServiceFormConfigHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	// 从URL路径中获取服务ID
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "缺少服务ID参数", http.StatusBadRequest)
		return
	}

	serviceIdStr := pathParts[3]
	serviceId, err := strconv.Atoi(serviceIdStr)
	if err != nil {
		http.Error(w, "无效的服务ID", http.StatusBadRequest)
		return
	}

	// 获取服务详情
	service, err := dao.ServiceImp.GetServiceById(int32(serviceId))
	if err != nil {
		response := &ServiceResponse{
			Code:     -1,
			ErrorMsg: "获取服务详情失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 解析表单配置
	var formConfig FormConfig
	if service.FormConfig != "" {
		if err := json.Unmarshal([]byte(service.FormConfig), &formConfig); err != nil {
			response := &ServiceResponse{
				Code:     -1,
				ErrorMsg: "解析表单配置失败: " + err.Error(),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	response := &ServiceResponse{
		Code: 0,
		Data: formConfig,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
