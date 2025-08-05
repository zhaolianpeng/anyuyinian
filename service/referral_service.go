package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

// ReferralResponse 推荐响应
type ReferralResponse struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

// ReferralReportResponse 推荐报告响应
type ReferralReportResponse struct {
	Referrer        *model.UserModel         `json:"referrer"`
	Referrals       []*model.ReferralModel   `json:"referrals"`
	Commissions     []*model.CommissionModel `json:"commissions"`
	TotalCommission float64                  `json:"totalCommission"`
}

// ReferralConfigResponse 推荐配置响应
type ReferralConfigResponse struct {
	CommissionRate float64 `json:"commissionRate"` // 佣金比例
	MinCashout     float64 `json:"minCashout"`     // 最低提现金额
	Rules          string  `json:"rules"`          // 规则说明
}

// ApplyCashoutRequest 申请提现请求
type ApplyCashoutRequest struct {
	UserId  string  `json:"userId"`
	Amount  float64 `json:"amount"`
	Method  string  `json:"method"` // wechat, alipay, bank
	Account string  `json:"account"`
}

// ReferralQrCodeHandler 获取用户专属推广二维码接口
func ReferralQrCodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	// 获取用户ID参数
	userIdStr := r.URL.Query().Get("userId")
	if userIdStr == "" {
		http.Error(w, "缺少userId参数", http.StatusBadRequest)
		return
	}

	// 直接使用字符串类型的userId
	userId := userIdStr

	// 获取或创建推荐关系
	referral, err := dao.ReferralImp.GetReferralByUserId(userId)
	if err != nil {
		// 如果不存在，创建一个新的推荐关系
		referral = &model.ReferralModel{
			UserId:     userId,
			ReferrerId: nil, // 设为nil，表示没有推荐人
			QrCodeUrl:  generateQrCodeUrl(userId),
			Status:     1,
		}
		if err := dao.ReferralImp.CreateReferral(referral); err != nil {
			response := &ReferralResponse{
				Code:     -1,
				ErrorMsg: "创建推荐关系失败: " + err.Error(),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	response := &ReferralResponse{
		Code: 0,
		Data: map[string]interface{}{
			"qrCodeUrl": referral.QrCodeUrl,
			"userId":    referral.UserId,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ReferralReportHandler 获取我的推荐人及下单记录接口
func ReferralReportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	// 获取用户ID参数
	userIdStr := r.URL.Query().Get("userId")
	if userIdStr == "" {
		http.Error(w, "缺少userId参数", http.StatusBadRequest)
		return
	}

	// 直接使用字符串类型的userId
	userId := userIdStr

	// 获取推荐关系
	referral, err := dao.ReferralImp.GetReferralByUserId(userId)
	if err != nil {
		response := &ReferralResponse{
			Code:     -1,
			ErrorMsg: "获取推荐关系失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 获取推荐人信息（简化处理）
	var referrer *model.UserModel
	if referral.ReferrerId != nil && *referral.ReferrerId != "" {
		referrer, _ = dao.UserImp.GetUserByUserId(*referral.ReferrerId)
	}

	// 获取我推荐的用户列表
	referrals, _, err := dao.ReferralImp.GetReferralsByReferrerId(userId, 1, 100)
	if err != nil {
		response := &ReferralResponse{
			Code:     -1,
			ErrorMsg: "获取推荐列表失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 获取佣金记录
	commissions, _, err := dao.ReferralImp.GetCommissionsByUserId(userId, 1, 100)
	if err != nil {
		response := &ReferralResponse{
			Code:     -1,
			ErrorMsg: "获取佣金记录失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 计算总佣金
	var totalCommission float64
	for _, commission := range commissions {
		if commission.Status == 1 { // 已结算
			totalCommission += commission.Amount
		}
	}

	response := &ReferralResponse{
		Code: 0,
		Data: &ReferralReportResponse{
			Referrer:        referrer,
			Referrals:       referrals,
			Commissions:     commissions,
			TotalCommission: totalCommission,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ReferralConfigHandler 获取推荐返佣规则说明接口
func ReferralConfigHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	config := &ReferralConfigResponse{
		CommissionRate: 0.05, // 5%佣金比例
		MinCashout:     10.0, // 最低提现10元
		Rules: `推荐返佣规则：
1. 成功推荐好友注册并下单，可获得订单金额5%的佣金
2. 佣金在订单完成后自动结算
3. 累计佣金达到10元后可申请提现
4. 提现支持微信、支付宝、银行卡等方式
5. 提现申请将在1-3个工作日内处理完成`,
	}

	response := &ReferralResponse{
		Code: 0,
		Data: config,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ApplyCashoutHandler 申请佣金提现接口
func ApplyCashoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	var req ApplyCashoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "请求参数解析失败", http.StatusBadRequest)
		return
	}

	// 验证参数
	if req.UserId == "" || req.Amount <= 0 || req.Method == "" || req.Account == "" {
		http.Error(w, "缺少必要参数", http.StatusBadRequest)
		return
	}

	// 检查最低提现金额
	if req.Amount < 10.0 {
		response := &ReferralResponse{
			Code:     -1,
			ErrorMsg: "提现金额不能少于10元",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 获取用户佣金记录，检查可提现金额
	commissions, _, err := dao.ReferralImp.GetCommissionsByUserId(req.UserId, 1, 1000)
	if err != nil {
		response := &ReferralResponse{
			Code:     -1,
			ErrorMsg: "获取佣金记录失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 计算可提现金额
	var availableAmount float64
	for _, commission := range commissions {
		if commission.Status == 1 { // 已结算且未提现
			availableAmount += commission.Amount
		}
	}

	if req.Amount > availableAmount {
		response := &ReferralResponse{
			Code:     -1,
			ErrorMsg: "可提现金额不足",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 创建提现记录
	cashout := &model.CashoutModel{
		UserId:  req.UserId,
		Amount:  req.Amount,
		Method:  req.Method,
		Account: req.Account,
		Status:  0, // 待审核
	}

	if err := dao.ReferralImp.CreateCashout(cashout); err != nil {
		response := &ReferralResponse{
			Code:     -1,
			ErrorMsg: "创建提现记录失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := &ReferralResponse{
		Code: 0,
		Data: map[string]interface{}{
			"cashoutId": cashout.Id,
			"amount":    cashout.Amount,
			"method":    cashout.Method,
			"status":    cashout.Status,
			"message":   "提现申请提交成功，将在1-3个工作日内处理",
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 生成二维码URL
func generateQrCodeUrl(userId string) string {
	return fmt.Sprintf("https://example.com/qrcode/user_%s.png", userId)
}
