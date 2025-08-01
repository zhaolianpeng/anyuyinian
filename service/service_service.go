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

// ServiceDetailRequest 服务详情请求
type ServiceDetailRequest struct {
	ServiceId json.Number `json:"serviceId"`
}

// ServiceListHandler 获取服务列表接口
func ServiceListHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理获取服务列表请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodGet {
		LogError("请求方法不支持", fmt.Errorf("期望GET方法，实际为%s", r.Method))
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	// 解析查询参数
	category := r.URL.Query().Get("category")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	LogStep("解析查询参数", map[string]interface{}{
		"category": category,
		"page":     pageStr,
		"pageSize": pageSizeStr,
	})

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

	LogStep("设置分页参数", map[string]interface{}{
		"page":     page,
		"pageSize": pageSize,
	})

	var services []*model.ServiceItemModel
	var total int64
	var err error

	// 根据分类获取服务列表
	if category != "" {
		LogStep("开始按分类查询服务列表", map[string]interface{}{
			"category": category,
			"page":     page,
			"pageSize": pageSize,
		})
		services, total, err = dao.ServiceImp.GetServicesByCategory(category, page, pageSize)
	} else {
		LogStep("开始查询所有服务列表", map[string]interface{}{
			"page":     page,
			"pageSize": pageSize,
		})
		services, total, err = dao.ServiceImp.GetAllServices(page, pageSize)
	}

	if err != nil {
		LogError("数据库查询服务列表失败", err)
		response := &ServiceResponse{
			Code:     -1,
			ErrorMsg: "获取服务列表失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	LogStep("服务列表查询成功", map[string]interface{}{
		"serviceCount": len(services),
		"total":        total,
		"category":     category,
	})

	// 计算是否有更多数据
	hasMore := int64(page*pageSize) < total
	LogStep("计算分页信息", map[string]interface{}{
		"hasMore": hasMore,
		"total":   total,
		"current": int64(page * pageSize),
	})

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

	LogInfo("服务列表获取成功", map[string]interface{}{
		"serviceCount": len(services),
		"total":        total,
		"page":         page,
		"pageSize":     pageSize,
		"hasMore":      hasMore,
	})
}

// ServiceDetailHandler 获取服务详情接口
func ServiceDetailHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理获取服务详情请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodPost {
		LogError("请求方法不支持", fmt.Errorf("期望POST方法，实际为%s", r.Method))
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	var req ServiceDetailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogError("请求参数解析失败", err)
		response := &ServiceResponse{
			Code:     -1,
			ErrorMsg: "请求参数解析失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	LogStep("解析服务详情请求参数", map[string]interface{}{
		"serviceId": req.ServiceId.String(),
	})

	// 将json.Number转换为int32
	LogStep("转换服务ID格式", map[string]interface{}{
		"serviceIdString": req.ServiceId.String(),
	})

	serviceId, err := req.ServiceId.Int64()
	if err != nil {
		LogError("服务ID格式转换失败", err)
		response := &ServiceResponse{
			Code:     -1,
			ErrorMsg: "无效的服务ID格式: " + req.ServiceId.String(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	LogStep("服务ID转换成功", map[string]interface{}{
		"serviceId": serviceId,
	})

	// 验证服务ID
	if serviceId <= 0 {
		LogError("服务ID无效", fmt.Errorf("serviceId=%d", serviceId))
		response := &ServiceResponse{
			Code:     -1,
			ErrorMsg: "无效的服务ID",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 获取服务详情
	LogStep("开始查询服务详情", map[string]interface{}{
		"serviceId": serviceId,
	})

	service, err := dao.ServiceImp.GetServiceById(int32(serviceId))
	if err != nil {
		LogError("数据库查询服务详情失败", err)
		// 检查是否是记录不存在错误
		if strings.Contains(err.Error(), "record not found") || strings.Contains(err.Error(), "no rows") {
			LogError("服务不存在", fmt.Errorf("serviceId=%d", serviceId))
			response := &ServiceResponse{
				Code:     -1,
				ErrorMsg: fmt.Sprintf("服务ID %d 不存在，请检查服务ID是否正确", serviceId),
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

	LogStep("服务详情查询成功", map[string]interface{}{
		"serviceId":   service.Id,
		"serviceName": service.Name,
		"price":       service.Price,
	})

	// 检查服务状态
	if service.Status == 0 {
		response := &ServiceResponse{
			Code:     -1,
			ErrorMsg: fmt.Sprintf("服务ID %d 已下架，暂不可用", serviceId),
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
