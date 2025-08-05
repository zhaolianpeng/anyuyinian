package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"
)

// PromoterResponse 推广响应
type PromoterResponse struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

// PromoterInfo 推广员信息
type PromoterInfo struct {
	UserId       string  `json:"userId"`
	PromoterCode string  `json:"promoterCode"` // 六位推广码
	NickName     string  `json:"nickName"`
	AvatarUrl    string  `json:"avatarUrl"`
	QrCodeUrl    string  `json:"qrCodeUrl"`
	TotalIncome  float64 `json:"totalIncome"`
	TodayIncome  float64 `json:"todayIncome"`
	MonthIncome  float64 `json:"monthIncome"`
	TotalOrders  int     `json:"totalOrders"`
	TodayOrders  int     `json:"todayOrders"`
	MonthOrders  int     `json:"monthOrders"`
}

// CommissionInfo 佣金信息
type CommissionInfo struct {
	Id          int32      `json:"id"`
	OrderId     int32      `json:"orderId"`
	OrderNo     string     `json:"orderNo"`
	Amount      float64    `json:"amount"`
	Rate        float64    `json:"rate"`
	Status      int        `json:"status"`
	StatusText  string     `json:"statusText"`
	CashoutTime *time.Time `json:"cashoutTime"`
	CreatedAt   time.Time  `json:"createdAt"`
}

// CashoutInfo 提现信息
type CashoutInfo struct {
	Id          int32      `json:"id"`
	Amount      float64    `json:"amount"`
	Method      string     `json:"method"`
	MethodText  string     `json:"methodText"`
	Account     string     `json:"account"`
	Status      int        `json:"status"`
	StatusText  string     `json:"statusText"`
	Remark      string     `json:"remark"`
	ProcessTime *time.Time `json:"processTime"`
	CreatedAt   time.Time  `json:"createdAt"`
}

// PromoterStats 推广统计
type PromoterStats struct {
	TotalIncome     float64 `json:"totalIncome"`
	TodayIncome     float64 `json:"todayIncome"`
	MonthIncome     float64 `json:"monthIncome"`
	TotalOrders     int     `json:"totalOrders"`
	TodayOrders     int     `json:"todayOrders"`
	MonthOrders     int     `json:"monthOrders"`
	PendingAmount   float64 `json:"pendingAmount"`
	SettledAmount   float64 `json:"settledAmount"`
	WithdrawnAmount float64 `json:"withdrawnAmount"`
}

// GetPromoterInfoHandler 获取推广员信息接口
func GetPromoterInfoHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理获取推广员信息请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodGet {
		LogError("请求方法不支持", fmt.Errorf("期望GET方法，实际为%s", r.Method))
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	// 获取用户ID参数
	userIdStr := r.URL.Query().Get("userId")
	LogStep("解析请求参数", map[string]interface{}{
		"userId": userIdStr,
	})

	if userIdStr == "" {
		LogError("缺少必要参数", fmt.Errorf("userId参数为空"))
		http.Error(w, "缺少userId参数", http.StatusBadRequest)
		return
	}

	// 获取用户信息
	user, err := dao.UserImp.GetUserByUserId(userIdStr)
	if err != nil {
		LogError("获取用户信息失败", err)
		response := &PromoterResponse{
			Code:     -1,
			ErrorMsg: "获取用户信息失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 获取或创建推荐关系
	referral, err := dao.ReferralImp.GetReferralByUserId(userIdStr)
	if err != nil {
		// 如果不存在，创建一个新的推荐关系
		referral = &model.ReferralModel{
			UserId:       userIdStr,
			ReferrerId:   nil,                          // 设为nil，表示没有推荐人
			PromoterCode: generateUniquePromoterCode(), // 生成唯一推广码
			QrCodeUrl:    generateQrCodeUrl(userIdStr),
			Status:       1,
		}
		if err := dao.ReferralImp.CreateReferral(referral); err != nil {
			LogError("创建推荐关系失败", err)
			response := &PromoterResponse{
				Code:     -1,
				ErrorMsg: "创建推荐关系失败: " + err.Error(),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
	} else {
		// 如果存在但没有推广码，生成一个
		if referral.PromoterCode == "" {
			referral.PromoterCode = generateUniquePromoterCode()
			if err := dao.ReferralImp.UpdateReferral(referral); err != nil {
				LogError("更新推广码失败", err)
				response := &PromoterResponse{
					Code:     -1,
					ErrorMsg: "更新推广码失败: " + err.Error(),
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				return
			}
		}
	}

	// 获取推广统计
	stats, err := getPromoterStats(userIdStr)
	if err != nil {
		LogError("获取推广统计失败", err)
		response := &PromoterResponse{
			Code:     -1,
			ErrorMsg: "获取推广统计失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 构建推广员信息
	promoterInfo := &PromoterInfo{
		UserId:       user.UserId,
		PromoterCode: referral.PromoterCode, // 使用数据库中的推广码
		NickName:     user.NickName,
		AvatarUrl:    user.AvatarUrl,
		QrCodeUrl:    referral.QrCodeUrl,
		TotalIncome:  stats.TotalIncome,
		TodayIncome:  stats.TodayIncome,
		MonthIncome:  stats.MonthIncome,
		TotalOrders:  stats.TotalOrders,
		TodayOrders:  stats.TodayOrders,
		MonthOrders:  stats.MonthOrders,
	}

	LogStep("推广员信息获取成功", map[string]interface{}{
		"userId":       promoterInfo.UserId,
		"promoterCode": promoterInfo.PromoterCode,
		"nickName":     promoterInfo.NickName,
		"totalIncome":  promoterInfo.TotalIncome,
		"totalOrders":  promoterInfo.TotalOrders,
	})

	response := &PromoterResponse{
		Code: 0,
		Data: promoterInfo,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetCommissionListHandler 获取佣金记录列表接口
func GetCommissionListHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理获取佣金记录列表请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodGet {
		LogError("请求方法不支持", fmt.Errorf("期望GET方法，实际为%s", r.Method))
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	// 获取用户ID参数
	userIdStr := r.URL.Query().Get("userId")
	if userIdStr == "" {
		http.Error(w, "缺少userId参数", http.StatusBadRequest)
		return
	}

	// 获取分页参数
	page := 1
	pageSize := 20
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if _, err := fmt.Sscanf(pageStr, "%d", &page); err != nil {
			page = 1
		}
	}
	if pageSizeStr := r.URL.Query().Get("pageSize"); pageSizeStr != "" {
		if _, err := fmt.Sscanf(pageSizeStr, "%d", &pageSize); err != nil {
			pageSize = 20
		}
	}

	// 获取佣金记录列表
	commissions, total, err := dao.CommissionImp.GetCommissionsByUserId(userIdStr, page, pageSize)
	if err != nil {
		LogError("获取佣金记录失败", err)
		response := &PromoterResponse{
			Code:     -1,
			ErrorMsg: "获取佣金记录失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 转换为前端格式
	var commissionList []*CommissionInfo
	for _, commission := range commissions {
		commissionInfo := &CommissionInfo{
			Id:          commission.Id,
			OrderId:     commission.OrderId,
			OrderNo:     commission.OrderNo,
			Amount:      commission.Amount,
			Rate:        commission.Rate,
			Status:      commission.Status,
			StatusText:  getCommissionStatusText(commission.Status),
			CashoutTime: commission.CashoutTime,
			CreatedAt:   commission.CreatedAt,
		}
		commissionList = append(commissionList, commissionInfo)
	}

	response := &PromoterResponse{
		Code: 0,
		Data: map[string]interface{}{
			"list":     commissionList,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
			"hasMore":  int64(page*pageSize) < total,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetCashoutListHandler 获取提现记录列表接口
func GetCashoutListHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理获取提现记录列表请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodGet {
		LogError("请求方法不支持", fmt.Errorf("期望GET方法，实际为%s", r.Method))
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	// 获取用户ID参数
	userIdStr := r.URL.Query().Get("userId")
	if userIdStr == "" {
		http.Error(w, "缺少userId参数", http.StatusBadRequest)
		return
	}

	// 获取分页参数
	page := 1
	pageSize := 20
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if _, err := fmt.Sscanf(pageStr, "%d", &page); err != nil {
			page = 1
		}
	}
	if pageSizeStr := r.URL.Query().Get("pageSize"); pageSizeStr != "" {
		if _, err := fmt.Sscanf(pageSizeStr, "%d", &pageSize); err != nil {
			pageSize = 20
		}
	}

	// 获取提现记录列表
	cashouts, total, err := dao.CashoutImp.GetCashoutsByUserId(userIdStr, page, pageSize)
	if err != nil {
		LogError("获取提现记录失败", err)
		response := &PromoterResponse{
			Code:     -1,
			ErrorMsg: "获取提现记录失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 转换为前端格式
	var cashoutList []*CashoutInfo
	for _, cashout := range cashouts {
		cashoutInfo := &CashoutInfo{
			Id:          cashout.Id,
			Amount:      cashout.Amount,
			Method:      cashout.Method,
			MethodText:  getCashoutMethodText(cashout.Method),
			Account:     cashout.Account,
			Status:      cashout.Status,
			StatusText:  getCashoutStatusText(cashout.Status),
			Remark:      cashout.Remark,
			ProcessTime: cashout.ProcessTime,
			CreatedAt:   cashout.CreatedAt,
		}
		cashoutList = append(cashoutList, cashoutInfo)
	}

	response := &PromoterResponse{
		Code: 0,
		Data: map[string]interface{}{
			"list":     cashoutList,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
			"hasMore":  int64(page*pageSize) < total,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 获取推广统计
func getPromoterStats(userId string) (*PromoterStats, error) {
	// 获取佣金记录
	commissions, _, err := dao.CommissionImp.GetCommissionsByUserId(userId, 1, 1000)
	if err != nil {
		return nil, err
	}

	stats := &PromoterStats{}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	for _, commission := range commissions {
		// 总收入和订单数
		stats.TotalIncome += commission.Amount
		stats.TotalOrders++

		// 今日收入和订单数
		if commission.CreatedAt.After(today) {
			stats.TodayIncome += commission.Amount
			stats.TodayOrders++
		}

		// 本月收入和订单数
		if commission.CreatedAt.After(monthStart) {
			stats.MonthIncome += commission.Amount
			stats.MonthOrders++
		}

		// 按状态分类
		switch commission.Status {
		case 0: // 待结算
			stats.PendingAmount += commission.Amount
		case 1: // 已结算
			stats.SettledAmount += commission.Amount
		case 2: // 已提现
			stats.WithdrawnAmount += commission.Amount
		}
	}

	return stats, nil
}

// 获取可提现金额
func getAvailableCashoutAmount(userId string) (float64, error) {
	// 获取已结算但未提现的佣金总额
	commissions, _, err := dao.CommissionImp.GetCommissionsByUserId(userId, 1, 1000)
	if err != nil {
		return 0, err
	}

	var availableAmount float64
	for _, commission := range commissions {
		if commission.Status == 1 { // 已结算
			availableAmount += commission.Amount
		}
	}

	return availableAmount, nil
}

// 获取佣金状态文本
func getCommissionStatusText(status int) string {
	switch status {
	case 0:
		return "待结算"
	case 1:
		return "已结算"
	case 2:
		return "已提现"
	default:
		return "未知"
	}
}

// 获取提现方式文本
func getCashoutMethodText(method string) string {
	switch method {
	case "wechat":
		return "微信"
	case "alipay":
		return "支付宝"
	case "bank":
		return "银行卡"
	default:
		return method
	}
}

// 获取提现状态文本
func getCashoutStatusText(status int) string {
	switch status {
	case 0:
		return "待审核"
	case 1:
		return "已通过"
	case 2:
		return "已拒绝"
	case 3:
		return "已到账"
	default:
		return "未知"
	}
}

// generateUniquePromoterCode 生成唯一的推广码
func generateUniquePromoterCode() string {
	// 检查推广码是否已存在的函数
	checkExists := func(code string) bool {
		var count int64
		cli := db.Get()
		cli.Table("Referrals").Where("promoterCode = ?", code).Count(&count)
		return count > 0
	}

	return utils.GenerateUniquePromoterCode(checkExists)
}

// GetUserByPromoterCodeHandler 通过推广码查找用户接口
func GetUserByPromoterCodeHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理通过推广码查找用户请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodGet {
		LogError("请求方法不支持", fmt.Errorf("期望GET方法，实际为%s", r.Method))
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	// 获取推广码参数
	promoterCode := r.URL.Query().Get("promoterCode")
	LogStep("解析请求参数", map[string]interface{}{
		"promoterCode": promoterCode,
	})

	if promoterCode == "" {
		LogError("缺少必要参数", fmt.Errorf("promoterCode参数为空"))
		http.Error(w, "缺少promoterCode参数", http.StatusBadRequest)
		return
	}

	// 验证推广码格式
	if !utils.ValidatePromoterCode(promoterCode) {
		LogError("推广码格式无效", fmt.Errorf("推广码格式错误: %s", promoterCode))
		response := &PromoterResponse{
			Code:     -1,
			ErrorMsg: "推广码格式无效",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 通过推广码查找用户
	referral, err := dao.ReferralImp.GetReferralByPromoterCode(promoterCode)
	if err != nil {
		LogError("通过推广码查找用户失败", err)
		response := &PromoterResponse{
			Code:     -1,
			ErrorMsg: "推广码不存在或已失效",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 获取用户信息
	user, err := dao.UserImp.GetUserByUserId(referral.UserId)
	if err != nil {
		LogError("获取用户信息失败", err)
		response := &PromoterResponse{
			Code:     -1,
			ErrorMsg: "获取用户信息失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 构建返回数据
	userInfo := map[string]interface{}{
		"userId":       user.UserId,
		"nickName":     user.NickName,
		"avatarUrl":    user.AvatarUrl,
		"promoterCode": referral.PromoterCode,
	}

	LogStep("通过推广码查找用户成功", map[string]interface{}{
		"promoterCode": promoterCode,
		"userId":       user.UserId,
		"nickName":     user.NickName,
	})

	response := &PromoterResponse{
		Code: 0,
		Data: userInfo,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GeneratePromoterCodesHandler 批量生成推广码接口
func GeneratePromoterCodesHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理批量生成推广码请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodPost {
		LogError("请求方法不支持", fmt.Errorf("期望POST方法，实际为%s", r.Method))
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	// 获取所有没有推广码的用户
	var referrals []*model.ReferralModel
	cli := db.Get()
	err := cli.Table("Referrals").Where("promoterCode IS NULL OR promoterCode = ''").Find(&referrals).Error
	if err != nil {
		LogError("查询用户失败", err)
		response := &PromoterResponse{
			Code:     -1,
			ErrorMsg: "查询用户失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 为每个用户生成推广码
	successCount := 0
	failedCount := 0
	var results []map[string]interface{}

	for _, referral := range referrals {
		// 生成唯一推广码
		promoterCode := generateUniquePromoterCode()

		// 更新数据库
		err := cli.Table("Referrals").Where("id = ?", referral.Id).Update("promoterCode", promoterCode).Error
		if err != nil {
			LogError("更新推广码失败", fmt.Errorf("用户ID: %s, 错误: %v", referral.UserId, err))
			failedCount++
			results = append(results, map[string]interface{}{
				"userId":       referral.UserId,
				"promoterCode": "",
				"status":       "失败",
				"error":        err.Error(),
			})
		} else {
			successCount++
			results = append(results, map[string]interface{}{
				"userId":       referral.UserId,
				"promoterCode": promoterCode,
				"status":       "成功",
			})
		}
	}

	LogStep("批量生成推广码完成", map[string]interface{}{
		"totalCount":   len(referrals),
		"successCount": successCount,
		"failedCount":  failedCount,
	})

	response := &PromoterResponse{
		Code: 0,
		Data: map[string]interface{}{
			"totalCount":   len(referrals),
			"successCount": successCount,
			"failedCount":  failedCount,
			"results":      results,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
