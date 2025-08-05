package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

// KefuResponse 客服响应
type KefuResponse struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	UserId     string   `json:"userId"`
	UserName   string   `json:"userName"`
	UserAvatar string   `json:"userAvatar"`
	Content    string   `json:"content"`
	Images     []string `json:"images"`
}

// FaqListResponse FAQ列表响应
type FaqListResponse struct {
	List     []*model.FaqModel `json:"list"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"pageSize"`
	HasMore  bool              `json:"hasMore"`
}

// SendMessageHandler 提交用户咨询问题接口
func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理发送消息请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodPost {
		LogError("请求方法不支持", fmt.Errorf("期望POST方法，实际为%s", r.Method))
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogError("请求参数解析失败", err)
		http.Error(w, "请求参数解析失败", http.StatusBadRequest)
		return
	}

	LogStep("解析发送消息请求参数", map[string]interface{}{
		"userId":     req.UserId,
		"userName":   req.UserName,
		"content":    req.Content,
		"imageCount": len(req.Images),
	})

	// 验证参数
	if req.UserId == "" || req.Content == "" {
		LogError("缺少必要参数", fmt.Errorf("userId=%s, content=%s", req.UserId, req.Content))
		http.Error(w, "缺少必要参数", http.StatusBadRequest)
		return
	}

	// 转换图片数组为JSON
	LogStep("处理图片数据", map[string]interface{}{
		"imageCount": len(req.Images),
	})
	imagesJson, _ := json.Marshal(req.Images)

	// 创建消息
	LogStep("开始创建消息对象", map[string]interface{}{
		"userId":     req.UserId,
		"userName":   req.UserName,
		"content":    req.Content,
		"imageCount": len(req.Images),
	})

	message := &model.KefuMessageModel{
		UserId:     req.UserId,
		UserName:   req.UserName,
		UserAvatar: req.UserAvatar,
		Type:       1, // 用户消息
		Content:    req.Content,
		Images:     string(imagesJson),
		Status:     0, // 未读
	}

	LogStep("开始保存消息到数据库", map[string]interface{}{
		"userId":  message.UserId,
		"content": message.Content,
	})

	if err := dao.KefuImp.CreateMessage(message); err != nil {
		LogError("数据库保存消息失败", err)
		response := &KefuResponse{
			Code:     -1,
			ErrorMsg: "发送消息失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	LogStep("消息保存成功", map[string]interface{}{
		"messageId": message.Id,
		"userId":    message.UserId,
	})

	response := &KefuResponse{
		Code: 0,
		Data: map[string]interface{}{
			"messageId": message.Id,
			"message":   "消息发送成功，客服将尽快回复您",
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	LogInfo("消息发送成功", map[string]interface{}{
		"messageId": message.Id,
		"userId":    message.UserId,
		"content":   message.Content,
	})
}

// FaqHandler 常见问题列表接口
func FaqHandler(w http.ResponseWriter, r *http.Request) {
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

	var faqs []*model.FaqModel
	var total int64
	var err error

	// 根据分类获取FAQ列表
	if category != "" {
		faqs, total, err = dao.KefuImp.GetFaqsByCategory(category, page, pageSize)
	} else {
		faqs, total, err = dao.KefuImp.GetAllFaqs(page, pageSize)
	}

	if err != nil {
		response := &KefuResponse{
			Code:     -1,
			ErrorMsg: "获取FAQ列表失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 计算是否有更多数据
	hasMore := int64(page*pageSize) < total

	response := &KefuResponse{
		Code: 0,
		Data: &FaqListResponse{
			List:     faqs,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
			HasMore:  hasMore,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
