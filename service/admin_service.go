package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

// AdminResponse 管理员响应
type AdminResponse struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

// AdminLoginRequest 管理员登录请求
type AdminLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AdminLoginResponse 管理员登录响应
type AdminLoginResponse struct {
	UserId        string `json:"userId"`
	NickName      string `json:"nickName"`
	AvatarUrl     string `json:"avatarUrl"`
	AdminLevel    int    `json:"adminLevel"`
	AdminUsername string `json:"adminUsername"`
}

// AdminUserInfo 管理员用户信息
type AdminUserInfo struct {
	UserId         string     `json:"userId"`
	NickName       string     `json:"nickName"`
	AvatarUrl      string     `json:"avatarUrl"`
	Phone          string     `json:"phone"`
	IsAdmin        int        `json:"isAdmin"`
	AdminLevel     int        `json:"adminLevel"`
	AdminUsername  string     `json:"adminUsername"`
	ParentAdminId  string     `json:"parentAdminId"`
	AdminCreatedAt *time.Time `json:"adminCreatedAt"`
	CreatedAt      time.Time  `json:"createdAt"`
}

// AdminOrderInfo 管理员订单信息
type AdminOrderInfo struct {
	Id           int32     `json:"id"`
	OrderNo      string    `json:"orderNo"`
	UserId       string    `json:"userId"`
	UserNickName string    `json:"userNickName"`
	ServiceId    int32     `json:"serviceId"`
	ServiceName  string    `json:"serviceName"`
	Amount       float64   `json:"amount"`
	Status       int       `json:"status"`
	StatusText   string    `json:"statusText"`
	CreatedAt    time.Time `json:"createdAt"`
}

// UpdateOrderAmountRequest 修改订单金额请求
type UpdateOrderAmountRequest struct {
	OrderId   int32   `json:"orderId"`
	NewAmount float64 `json:"newAmount"`
	Reason    string  `json:"reason"`
}

// AdminRefundOrderRequest 管理员退款请求
type AdminRefundOrderRequest struct {
	OrderId      int32   `json:"orderId"`
	RefundAmount float64 `json:"refundAmount"`
	Reason       string  `json:"reason"`
	RefundStatus int     `json:"refundStatus"` // 1-退款中，2-已退款
}

// AdminLoginHandler 管理员登录接口
func AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理管理员登录请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodPost {
		LogError("请求方法不支持", fmt.Errorf("期望POST方法，实际为%s", r.Method))
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体
	var req AdminLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogError("解析请求体失败", err)
		http.Error(w, "请求体格式错误", http.StatusBadRequest)
		return
	}

	LogStep("解析登录请求", map[string]interface{}{
		"username": req.Username,
	})

	// 验证参数
	if req.Username == "" || req.Password == "" {
		LogError("缺少必要参数", fmt.Errorf("用户名或密码为空"))
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "用户名和密码不能为空",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 管理员登录
	adminImp := &dao.AdminImp{}
	admin, err := adminImp.AdminLogin(req.Username, req.Password)
	if err != nil {
		LogError("管理员登录失败", err)
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "用户名或密码错误",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 记录登录日志
	log := &model.AdminLoginLogModel{
		AdminUserId: admin.UserId,
		LoginTime:   time.Now(),
		LoginIp:     r.RemoteAddr,
		UserAgent:   r.UserAgent(),
		Status:      1,
		Remark:      "管理员登录成功",
	}
	adminImp.LogAdminLogin(log)

	// 构建响应数据
	loginResponse := &AdminLoginResponse{
		UserId:        admin.UserId,
		NickName:      admin.NickName,
		AvatarUrl:     admin.AvatarUrl,
		AdminLevel:    admin.AdminLevel,
		AdminUsername: admin.AdminUsername,
	}

	LogStep("管理员登录成功", map[string]interface{}{
		"userId":        admin.UserId,
		"adminLevel":    admin.AdminLevel,
		"adminUsername": admin.AdminUsername,
	})

	response := &AdminResponse{
		Code: 0,
		Data: loginResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAdminUsersHandler 获取管理员可见的用户列表
func GetAdminUsersHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理获取管理员用户列表请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodGet {
		LogError("请求方法不支持", fmt.Errorf("期望GET方法，实际为%s", r.Method))
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	// 获取管理员用户ID
	adminUserId := r.URL.Query().Get("adminUserId")
	if adminUserId == "" {
		LogError("缺少必要参数", fmt.Errorf("adminUserId参数为空"))
		http.Error(w, "缺少adminUserId参数", http.StatusBadRequest)
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

	// 获取用户列表
	adminImp := &dao.AdminImp{}
	users, total, err := adminImp.GetVisibleUsers(adminUserId, page, pageSize)
	if err != nil {
		LogError("获取用户列表失败", err)
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "获取用户列表失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 转换为前端格式
	var userList []*AdminUserInfo
	if users != nil {
		for _, user := range users {
			userInfo := &AdminUserInfo{
				UserId:         user.UserId,
				NickName:       user.NickName,
				AvatarUrl:      user.AvatarUrl,
				Phone:          user.Phone,
				IsAdmin:        user.IsAdmin,
				AdminLevel:     user.AdminLevel,
				AdminUsername:  user.AdminUsername,
				ParentAdminId:  user.ParentAdminId,
				AdminCreatedAt: user.AdminCreatedAt,
				CreatedAt:      user.CreatedAt,
			}
			userList = append(userList, userInfo)
		}
	}

	// 确保返回空数组而不是null
	if userList == nil {
		userList = []*AdminUserInfo{}
	}

	response := &AdminResponse{
		Code: 0,
		Data: map[string]interface{}{
			"list":     userList,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
			"hasMore":  int64(page*pageSize) < total,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAdminOrdersHandler 获取管理员可见的订单列表
func GetAdminOrdersHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理获取管理员订单列表请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodGet {
		LogError("请求方法不支持", fmt.Errorf("期望GET方法，实际为%s", r.Method))
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	// 获取管理员用户ID
	adminUserId := r.URL.Query().Get("adminUserId")
	if adminUserId == "" {
		LogError("缺少必要参数", fmt.Errorf("adminUserId参数为空"))
		http.Error(w, "缺少adminUserId参数", http.StatusBadRequest)
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

	// 获取订单列表
	adminImp := &dao.AdminImp{}
	orders, total, err := adminImp.GetVisibleOrders(adminUserId, page, pageSize)
	if err != nil {
		LogError("获取订单列表失败", err)
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "获取订单列表失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 转换为前端格式
	var orderList []*AdminOrderInfo
	if orders != nil {
		for _, order := range orders {
			// 获取用户信息
			user, _ := dao.UserImp.GetUserByUserId(order.UserId)
			userNickName := ""
			if user != nil {
				userNickName = user.NickName
			}

			// 获取服务信息
			service, _ := dao.ServiceImp.GetServiceById(order.ServiceId)
			serviceName := ""
			if service != nil {
				serviceName = service.Name
			}

			orderInfo := &AdminOrderInfo{
				Id:           order.Id,
				OrderNo:      order.OrderNo,
				UserId:       order.UserId,
				UserNickName: userNickName,
				ServiceId:    order.ServiceId,
				ServiceName:  serviceName,
				Amount:       order.TotalAmount,
				Status:       order.Status,
				StatusText:   getOrderStatusText(order.Status),
				CreatedAt:    order.CreatedAt,
			}
			orderList = append(orderList, orderInfo)
		}
	}

	// 确保返回空数组而不是null
	if orderList == nil {
		orderList = []*AdminOrderInfo{}
	}

	response := &AdminResponse{
		Code: 0,
		Data: map[string]interface{}{
			"list":     orderList,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
			"hasMore":  int64(page*pageSize) < total,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SetUserAsAdminHandler 设置用户为管理员
func SetUserAsAdminHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理设置用户为管理员请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodPost {
		LogError("请求方法不支持", fmt.Errorf("期望POST方法，实际为%s", r.Method))
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体
	var req struct {
		UserId        string `json:"userId"`
		AdminLevel    int    `json:"adminLevel"`
		ParentAdminId string `json:"parentAdminId"`
		AdminUsername string `json:"adminUsername"`
		AdminPassword string `json:"adminPassword"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogError("解析请求体失败", err)
		http.Error(w, "请求体格式错误", http.StatusBadRequest)
		return
	}

	// 验证参数
	if req.UserId == "" || req.AdminUsername == "" || req.AdminPassword == "" {
		LogError("缺少必要参数", fmt.Errorf("用户ID、管理员用户名或密码为空"))
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "用户ID、管理员用户名和密码不能为空",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 检查用户是否存在
	_, err := dao.UserImp.GetUserByUserId(req.UserId)
	if err != nil {
		LogError("用户不存在", err)
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "用户不存在",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 设置用户为管理员
	adminImp := &dao.AdminImp{}
	err = adminImp.SetUserAsAdmin(req.UserId, req.AdminLevel, req.ParentAdminId)
	if err != nil {
		LogError("设置用户为管理员失败", err)
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "设置用户为管理员失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 更新管理员用户名和密码
	cli := db.Get()
	updates := map[string]interface{}{
		"adminUsername": req.AdminUsername,
		"adminPassword": req.AdminPassword,
		"updatedAt":     time.Now(),
	}
	err = cli.Table("Users").Where("userId = ?", req.UserId).Updates(updates).Error
	if err != nil {
		LogError("更新管理员信息失败", err)
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "更新管理员信息失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	LogStep("设置用户为管理员成功", map[string]interface{}{
		"userId":        req.UserId,
		"adminLevel":    req.AdminLevel,
		"adminUsername": req.AdminUsername,
	})

	response := &AdminResponse{
		Code: 0,
		Data: map[string]interface{}{
			"message": "设置用户为管理员成功",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RemoveAdminHandler 取消用户管理员权限
func RemoveAdminHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理取消管理员权限请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodPost {
		LogError("请求方法不支持", fmt.Errorf("期望POST方法，实际为%s", r.Method))
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体
	var req struct {
		UserId string `json:"userId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogError("解析请求体失败", err)
		http.Error(w, "请求体格式错误", http.StatusBadRequest)
		return
	}

	// 验证参数
	if req.UserId == "" {
		LogError("缺少必要参数", fmt.Errorf("userId参数为空"))
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "用户ID不能为空",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 取消管理员权限
	adminImp := &dao.AdminImp{}
	err := adminImp.RemoveAdmin(req.UserId)
	if err != nil {
		LogError("取消管理员权限失败", err)
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "取消管理员权限失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	LogStep("取消管理员权限成功", map[string]interface{}{
		"userId": req.UserId,
	})

	response := &AdminResponse{
		Code: 0,
		Data: map[string]interface{}{
			"message": "取消管理员权限成功",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CheckAdminStatusHandler 检查用户管理员状态接口
func CheckAdminStatusHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理检查管理员状态请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodGet {
		LogError("请求方法不支持", fmt.Errorf("期望GET方法，实际为%s", r.Method))
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	// 获取用户ID参数
	userId := r.URL.Query().Get("userId")
	if userId == "" {
		LogError("缺少必要参数", fmt.Errorf("userId参数为空"))
		http.Error(w, "缺少userId参数", http.StatusBadRequest)
		return
	}

	// 检查用户是否为管理员
	adminImp := &dao.AdminImp{}
	admin, err := adminImp.GetAdminByUserId(userId)
	if err != nil {
		// 用户不是管理员
		response := &AdminResponse{
			Code: 0,
			Data: map[string]interface{}{
				"isAdmin": false,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 用户是管理员，返回管理员信息
	adminInfo := &AdminLoginResponse{
		UserId:        admin.UserId,
		NickName:      admin.NickName,
		AvatarUrl:     admin.AvatarUrl,
		AdminLevel:    admin.AdminLevel,
		AdminUsername: admin.AdminUsername,
	}

	LogStep("检查管理员状态成功", map[string]interface{}{
		"userId":        admin.UserId,
		"adminLevel":    admin.AdminLevel,
		"adminUsername": admin.AdminUsername,
	})

	response := &AdminResponse{
		Code: 0,
		Data: map[string]interface{}{
			"isAdmin":   true,
			"adminInfo": adminInfo,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 管理员数据概览统计接口
func AdminStatsHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理管理员数据概览统计请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodGet {
		LogError("请求方法不支持", fmt.Errorf("期望GET方法，实际为%s", r.Method))
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	adminUserId := r.URL.Query().Get("adminUserId")
	if adminUserId == "" {
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "缺少adminUserId参数",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	adminImp := &dao.AdminImp{}
	admin, err := adminImp.GetAdminByUserId(adminUserId)
	if err != nil || admin == nil || admin.IsAdmin == 0 {
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "无效的管理员账号",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	dbCli := db.Get()
	var totalUsers int64
	var totalOrders int64
	var todayOrders int64
	var totalAmount float64
	var paidAmount float64
	var unpaidAmount float64
	var refundAmount float64
	var timeoutUnpaidAmount float64

	// 根据管理员级别获取不同的数据
	if admin.AdminLevel == 2 { // 超级管理员
		dbCli.Model(&model.UserModel{}).Count(&totalUsers)
		dbCli.Model(&model.OrderModel{}).Count(&totalOrders)
		dbCli.Model(&model.OrderModel{}).Where("DATE(createdAt) = CURDATE()").Count(&todayOrders)
		dbCli.Model(&model.OrderModel{}).Select("IFNULL(SUM(totalAmount),0)").Row().Scan(&totalAmount)
		// 已支付总金额（status = 1 或 payStatus = 1）
		dbCli.Model(&model.OrderModel{}).Where("status = 1 OR payStatus = 1").Select("IFNULL(SUM(totalAmount),0)").Row().Scan(&paidAmount)
		// 待支付总金额（status = 0 且 payStatus = 0）
		dbCli.Model(&model.OrderModel{}).Where("status = 0 AND payStatus = 0").Select("IFNULL(SUM(totalAmount),0)").Row().Scan(&unpaidAmount)
		// 退款总金额（status = 4 或 refundStatus = 2）
		dbCli.Model(&model.OrderModel{}).Where("status = 4 OR refundStatus = 2").Select("IFNULL(SUM(totalAmount),0)").Row().Scan(&refundAmount)
		// 超时未支付总金额（status = 0 或 status = 3 且 payStatus = 0 且 payDeadline < NOW()）
		dbCli.Model(&model.OrderModel{}).Where("(status = 0 OR status = 3) AND payStatus = 0 AND payDeadline IS NOT NULL AND payDeadline < NOW()").Select("IFNULL(SUM(totalAmount),0)").Row().Scan(&timeoutUnpaidAmount)
	} else { // 一级管理员
		// 获取该管理员推广的用户ID列表
		var promotedUserIds []string
		dbCli.Model(&model.UserModel{}).Where("referrerId = ?", adminUserId).Pluck("userId", &promotedUserIds)
		promotedUserIds = append(promotedUserIds, adminUserId) // 包含管理员自己的用户ID

		dbCli.Model(&model.UserModel{}).Where("userId IN (?)", promotedUserIds).Count(&totalUsers)
		dbCli.Model(&model.OrderModel{}).Where("userId IN (?)", promotedUserIds).Count(&totalOrders)
		dbCli.Model(&model.OrderModel{}).Where("userId IN (?) AND DATE(createdAt) = CURDATE()", promotedUserIds).Count(&todayOrders)
		dbCli.Model(&model.OrderModel{}).Where("userId IN (?)", promotedUserIds).Select("IFNULL(SUM(totalAmount),0)").Row().Scan(&totalAmount)
		// 已支付总金额（status = 1 或 payStatus = 1）
		dbCli.Model(&model.OrderModel{}).Where("userId IN (?) AND (status = 1 OR payStatus = 1)", promotedUserIds).Select("IFNULL(SUM(totalAmount),0)").Row().Scan(&paidAmount)
		// 待支付总金额（status = 0 且 payStatus = 0）
		dbCli.Model(&model.OrderModel{}).Where("userId IN (?) AND status = 0 AND payStatus = 0", promotedUserIds).Select("IFNULL(SUM(totalAmount),0)").Row().Scan(&unpaidAmount)
		// 退款总金额（status = 4 或 refundStatus = 2）
		dbCli.Model(&model.OrderModel{}).Where("userId IN (?) AND (status = 4 OR refundStatus = 2)", promotedUserIds).Select("IFNULL(SUM(totalAmount),0)").Row().Scan(&refundAmount)
		// 超时未支付总金额（status = 0 或 status = 3 且 payStatus = 0 且 payDeadline < NOW()）
		dbCli.Model(&model.OrderModel{}).Where("userId IN (?) AND (status = 0 OR status = 3) AND payStatus = 0 AND payDeadline IS NOT NULL AND payDeadline < NOW()", promotedUserIds).Select("IFNULL(SUM(totalAmount),0)").Row().Scan(&timeoutUnpaidAmount)
	}

	stats := map[string]interface{}{
		"totalUsers":          totalUsers,
		"totalOrders":         totalOrders,
		"todayOrders":         todayOrders,
		"totalAmount":         totalAmount,
		"paidAmount":          paidAmount,
		"unpaidAmount":        unpaidAmount,
		"refundAmount":        refundAmount,
		"timeoutUnpaidAmount": timeoutUnpaidAmount,
	}
	response := &AdminResponse{
		Code: 0,
		Data: stats,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 管理员列表接口
func AdminAdminsHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理获取管理员列表请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodGet {
		LogError("请求方法不支持", fmt.Errorf("期望GET方法，实际为%s", r.Method))
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	adminUserId := r.URL.Query().Get("adminUserId")
	if adminUserId == "" {
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "缺少adminUserId参数",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	adminImp := &dao.AdminImp{}
	admin, err := adminImp.GetAdminByUserId(adminUserId)
	if err != nil || admin == nil || admin.IsAdmin == 0 {
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "无效的管理员账号",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	// 根据管理员级别获取不同的数据
	var admins []*model.UserModel
	var total int64

	if admin.AdminLevel == 2 { // 超级管理员
		// 获取所有管理员
		admins, total, err = adminImp.GetAllAdmins(page, pageSize)
	} else { // 一级管理员
		// 只获取自己的下级管理员
		admins, total, err = adminImp.GetSubAdmins(adminUserId, page, pageSize)
	}

	if err != nil {
		LogError("获取管理员列表失败", err)
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "获取管理员列表失败",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 转换为前端需要的格式
	var adminList []map[string]interface{}
	for _, admin := range admins {
		// 统计该管理员推荐码下单的总金额
		var totalAmount float64
		if admin.AdminLevel > 0 { // 只统计管理员
			dbCli := db.Get()
			dbCli.Table("Orders").Where("referrerId = ?", admin.UserId).Select("IFNULL(SUM(totalAmount),0)").Scan(&totalAmount)
		}

		adminInfo := map[string]interface{}{
			"userId":         admin.UserId,
			"nickName":       admin.NickName,
			"avatarUrl":      admin.AvatarUrl,
			"phone":          admin.Phone,
			"isAdmin":        admin.IsAdmin,
			"adminLevel":     admin.AdminLevel,
			"adminUsername":  admin.AdminUsername,
			"parentAdminId":  admin.ParentAdminId,
			"adminCreatedAt": admin.AdminCreatedAt,
			"createdAt":      admin.CreatedAt,
			"totalAmount":    totalAmount, // 添加推荐码下单总金额
		}
		adminList = append(adminList, adminInfo)
	}

	hasMore := (page * pageSize) < int(total)

	response := &AdminResponse{
		Code: 0,
		Data: map[string]interface{}{
			"list":    adminList,
			"total":   total,
			"page":    page,
			"hasMore": hasMore,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 设置管理员接口
func SetAdminHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理设置管理员请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodPost {
		LogError("请求方法不支持", fmt.Errorf("期望POST方法，实际为%s", r.Method))
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体
	var req struct {
		UserId        string `json:"userId"`
		AdminLevel    int    `json:"adminLevel"`
		ParentAdminId string `json:"parentAdminId"`
		AdminUsername string `json:"adminUsername"`
		AdminPassword string `json:"adminPassword"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogError("解析请求体失败", err)
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "请求体格式错误",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 验证参数
	if req.UserId == "" || req.AdminUsername == "" || req.AdminPassword == "" {
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "缺少必要参数",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	adminImp := &dao.AdminImp{}

	// 检查要设置的用户是否存在
	dbCli := db.Get()
	var user model.UserModel
	if err := dbCli.Where("userId = ?", req.UserId).First(&user).Error; err != nil {
		LogError("用户不存在", err)
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "用户不存在",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 检查用户是否已经是管理员
	if user.IsAdmin == 1 {
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "用户已经是管理员",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 设置用户为管理员
	err := adminImp.SetUserAsAdmin(req.UserId, req.AdminLevel, req.ParentAdminId)
	if err != nil {
		LogError("设置管理员失败", err)
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "设置管理员失败",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 更新用户的用户名和密码
	updateData := map[string]interface{}{
		"adminUsername":  req.AdminUsername,
		"adminPassword":  req.AdminPassword,
		"adminCreatedAt": time.Now(),
	}

	if err := dbCli.Model(&user).Updates(updateData).Error; err != nil {
		LogError("更新管理员信息失败", err)
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "更新管理员信息失败",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	LogStep("设置管理员成功", map[string]interface{}{
		"userId":        req.UserId,
		"adminLevel":    req.AdminLevel,
		"adminUsername": req.AdminUsername,
	})

	response := &AdminResponse{
		Code: 0,
		Data: map[string]interface{}{
			"userId":        req.UserId,
			"adminLevel":    req.AdminLevel,
			"adminUsername": req.AdminUsername,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 获取订单状态文本
func getOrderStatusText(status int) string {
	switch status {
	case 0:
		return "待支付"
	case 1:
		return "已支付"
	case 2:
		return "已取消"
	case 3:
		return "已完成"
	default:
		return "未知"
	}
}

// UpdateOrderAmountHandler 修改订单金额接口
func UpdateOrderAmountHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理修改订单金额请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodPost {
		LogError("请求方法不支持", fmt.Errorf("期望POST方法，实际为%s", r.Method))
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	// 获取管理员用户ID
	adminUserId := r.URL.Query().Get("adminUserId")
	if adminUserId == "" {
		LogError("缺少必要参数", fmt.Errorf("adminUserId参数为空"))
		http.Error(w, "缺少adminUserId参数", http.StatusBadRequest)
		return
	}

	// 解析请求体
	var req UpdateOrderAmountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogError("请求参数解析失败", err)
		http.Error(w, "请求参数解析失败", http.StatusBadRequest)
		return
	}

	LogStep("解析修改订单金额请求参数", map[string]interface{}{
		"orderId":   req.OrderId,
		"newAmount": req.NewAmount,
		"reason":    req.Reason,
	})

	// 验证参数
	if req.OrderId <= 0 {
		LogError("订单ID无效", fmt.Errorf("orderId=%d", req.OrderId))
		http.Error(w, "订单ID无效", http.StatusBadRequest)
		return
	}

	if req.NewAmount <= 0 {
		LogError("新金额无效", fmt.Errorf("newAmount=%f", req.NewAmount))
		http.Error(w, "新金额必须大于0", http.StatusBadRequest)
		return
	}

	// 检查管理员权限
	adminImp := &dao.AdminImp{}
	admin, err := adminImp.GetAdminByUserId(adminUserId)
	if err != nil || admin == nil || admin.IsAdmin == 0 {
		LogError("管理员权限验证失败", err)
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "无效的管理员账号",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 只有超级管理员可以修改订单金额
	if admin.AdminLevel != 2 {
		LogError("权限不足", fmt.Errorf("adminLevel=%d, 需要adminLevel=2", admin.AdminLevel))
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "只有超级管理员可以修改订单金额",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 获取订单信息
	order, err := dao.OrderImp.GetOrderById(req.OrderId)
	if err != nil {
		LogError("获取订单信息失败", err)
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "获取订单信息失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if order == nil {
		LogError("订单不存在", fmt.Errorf("orderId=%d", req.OrderId))
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "订单不存在",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 检查订单状态，只有未支付的订单可以修改金额
	if order.Status != 0 || order.PayStatus != 0 {
		LogError("订单状态不允许修改金额", fmt.Errorf("status=%d, payStatus=%d", order.Status, order.PayStatus))
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "只有未支付的订单可以修改金额",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 记录修改前的金额
	oldAmount := order.TotalAmount

	// 更新订单金额
	err = dao.OrderImp.UpdateOrderAmount(req.OrderId, req.NewAmount)
	if err != nil {
		LogError("更新订单金额失败", err)
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "更新订单金额失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	LogStep("订单金额修改成功", map[string]interface{}{
		"orderId":   req.OrderId,
		"oldAmount": oldAmount,
		"newAmount": req.NewAmount,
		"adminId":   adminUserId,
		"reason":    req.Reason,
	})

	response := &AdminResponse{
		Code: 0,
		Data: map[string]interface{}{
			"orderId":   req.OrderId,
			"orderNo":   order.OrderNo,
			"oldAmount": oldAmount,
			"newAmount": req.NewAmount,
			"reason":    req.Reason,
			"adminId":   adminUserId,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AdminRefundOrderHandler 管理员退款订单接口
func AdminRefundOrderHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理管理员退款订单请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodPost {
		LogError("请求方法不支持", fmt.Errorf("期望POST方法，实际为%s", r.Method))
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	// 获取管理员用户ID
	adminUserId := r.URL.Query().Get("adminUserId")
	if adminUserId == "" {
		LogError("缺少必要参数", fmt.Errorf("adminUserId参数为空"))
		http.Error(w, "缺少adminUserId参数", http.StatusBadRequest)
		return
	}

	// 解析请求体
	var req AdminRefundOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogError("请求参数解析失败", err)
		http.Error(w, "请求参数解析失败", http.StatusBadRequest)
		return
	}

	LogStep("解析管理员退款请求参数", map[string]interface{}{
		"orderId":      req.OrderId,
		"refundAmount": req.RefundAmount,
		"reason":       req.Reason,
		"refundStatus": req.RefundStatus,
	})

	// 验证参数
	if req.OrderId <= 0 {
		LogError("订单ID无效", fmt.Errorf("orderId=%d", req.OrderId))
		http.Error(w, "订单ID无效", http.StatusBadRequest)
		return
	}

	if req.RefundAmount <= 0 {
		LogError("退款金额无效", fmt.Errorf("refundAmount=%f", req.RefundAmount))
		http.Error(w, "退款金额必须大于0", http.StatusBadRequest)
		return
	}

	if req.RefundStatus != 1 && req.RefundStatus != 2 {
		LogError("退款状态无效", fmt.Errorf("refundStatus=%d", req.RefundStatus))
		http.Error(w, "退款状态无效", http.StatusBadRequest)
		return
	}

	// 检查管理员权限
	adminImp := &dao.AdminImp{}
	admin, err := adminImp.GetAdminByUserId(adminUserId)
	if err != nil || admin == nil || admin.IsAdmin == 0 {
		LogError("管理员权限验证失败", err)
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "无效的管理员账号",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 只有超级管理员可以处理退款
	if admin.AdminLevel != 2 {
		LogError("权限不足", fmt.Errorf("adminLevel=%d, 需要adminLevel=2", admin.AdminLevel))
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "只有超级管理员可以处理退款",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 获取订单信息
	order, err := dao.OrderImp.GetOrderById(req.OrderId)
	if err != nil {
		LogError("获取订单信息失败", err)
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "获取订单信息失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if order == nil {
		LogError("订单不存在", fmt.Errorf("orderId=%d", req.OrderId))
		response := &AdminResponse{
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
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "只有已支付的订单可以退款",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 检查退款金额不能超过订单金额
	if req.RefundAmount > order.TotalAmount {
		LogError("退款金额超过订单金额", fmt.Errorf("refundAmount=%f, totalAmount=%f", req.RefundAmount, order.TotalAmount))
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "退款金额不能超过订单金额",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 更新退款状态
	if err := dao.OrderImp.UpdateRefundStatus(req.OrderId, req.RefundStatus, req.RefundAmount, req.Reason); err != nil {
		LogError("处理退款失败", err)
		response := &AdminResponse{
			Code:     -1,
			ErrorMsg: "处理退款失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 如果设置为已退款，同时更新订单状态
	if req.RefundStatus == 2 {
		if err := dao.OrderImp.UpdateOrderStatus(req.OrderId, 4); err != nil {
			LogError("更新订单状态失败", err)
			response := &AdminResponse{
				Code:     -1,
				ErrorMsg: "更新订单状态失败: " + err.Error(),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	LogStep("管理员退款处理成功", map[string]interface{}{
		"orderId":      req.OrderId,
		"orderNo":      order.OrderNo,
		"refundAmount": req.RefundAmount,
		"reason":       req.Reason,
		"refundStatus": req.RefundStatus,
		"adminId":      adminUserId,
	})

	response := &AdminResponse{
		Code: 0,
		Data: map[string]interface{}{
			"orderId":      req.OrderId,
			"orderNo":      order.OrderNo,
			"refundAmount": req.RefundAmount,
			"reason":       req.Reason,
			"refundStatus": req.RefundStatus,
			"adminId":      adminUserId,
			"message": func() string {
				if req.RefundStatus == 2 {
					return "退款处理成功"
				} else {
					return "退款状态更新成功"
				}
			}(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// =============================
// 服务管理：列表与修改价格
// =============================

// AdminServiceInfo 管理员服务信息
type AdminServiceInfo struct {
	Id            int32     `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Category      string    `json:"category"`
	Price         float64   `json:"price"`
	OriginalPrice float64   `json:"originalPrice"`
	ImageUrl      string    `json:"imageUrl"`
	Status        int       `json:"status"`
	StatusText    string    `json:"statusText"`
	Sort          int       `json:"sort"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// UpdateServicePriceRequest 修改服务价格请求
type UpdateServicePriceRequest struct {
	ServiceId        int32   `json:"serviceId"`
	NewPrice         float64 `json:"newPrice"`
	NewOriginalPrice float64 `json:"newOriginalPrice"`
	Reason           string  `json:"reason"`
}

// GetAdminServicesHandler 获取管理员服务列表接口
func GetAdminServicesHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理获取管理员服务列表请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodGet {
		LogError("请求方法不支持", fmt.Errorf("期望GET方法，实际为%s", r.Method))
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	// 读取管理员身份（支持query或header）
	adminUserId := r.URL.Query().Get("adminUserId")
	if adminUserId == "" {
		adminUserId = r.Header.Get("adminUserId")
	}
	if adminUserId == "" {
		LogError("缺少必要参数", fmt.Errorf("adminUserId参数为空"))
		http.Error(w, "缺少adminUserId参数", http.StatusBadRequest)
		return
	}

	// 解析查询参数
	page := 1
	pageSize := 20
	if v := r.URL.Query().Get("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			page = p
		}
	}
	if v := r.URL.Query().Get("pageSize"); v != "" {
		if ps, err := strconv.Atoi(v); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}
	category := r.URL.Query().Get("category")

	// 获取服务列表
	var (
		services []*model.ServiceItemModel
		total    int64
		err      error
	)
	if category != "" {
		services, total, err = dao.ServiceImp.GetServicesByCategory(category, page, pageSize)
	} else {
		services, total, err = dao.ServiceImp.GetAllServices(page, pageSize)
	}
	if err != nil {
		LogError("获取服务列表失败", err)
		response := &AdminResponse{Code: -1, ErrorMsg: "获取服务列表失败: " + err.Error()}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 转换返回
	list := make([]*AdminServiceInfo, 0, len(services))
	for _, s := range services {
		list = append(list, &AdminServiceInfo{
			Id:            s.Id,
			Name:          s.Name,
			Description:   s.Description,
			Category:      s.Category,
			Price:         s.Price,
			OriginalPrice: s.OriginalPrice,
			ImageUrl:      s.ImageUrl,
			Status:        s.Status,
			StatusText:    getServiceStatusText(s.Status),
			Sort:          s.Sort,
			CreatedAt:     s.CreatedAt,
			UpdatedAt:     s.UpdatedAt,
		})
	}

	response := &AdminResponse{Code: 0, Data: map[string]interface{}{
		"list":     list,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
		"hasMore":  int64(page*pageSize) < total,
	}}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateServicePriceHandler 修改服务价格接口
func UpdateServicePriceHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理修改服务价格请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodPost {
		LogError("请求方法不支持", fmt.Errorf("期望POST方法，实际为%s", r.Method))
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	// 读取管理员身份（支持query或header）
	adminUserId := r.URL.Query().Get("adminUserId")
	if adminUserId == "" {
		adminUserId = r.Header.Get("adminUserId")
	}
	if adminUserId == "" {
		LogError("缺少必要参数", fmt.Errorf("adminUserId参数为空"))
		http.Error(w, "缺少adminUserId参数", http.StatusBadRequest)
		return
	}

	// 解析请求体
	var req UpdateServicePriceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogError("解析请求体失败", err)
		response := &AdminResponse{Code: -1, ErrorMsg: "请求体格式错误"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if req.ServiceId <= 0 {
		response := &AdminResponse{Code: -1, ErrorMsg: "服务ID无效"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	if req.NewPrice < 0 || req.NewOriginalPrice < 0 {
		response := &AdminResponse{Code: -1, ErrorMsg: "价格不能为负数"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 获取服务并更新
	s, err := dao.ServiceImp.GetServiceById(req.ServiceId)
	if err != nil || s == nil {
		if err == nil {
			err = fmt.Errorf("service not found")
		}
		LogError("获取服务信息失败", err)
		response := &AdminResponse{Code: -1, ErrorMsg: "获取服务信息失败: " + err.Error()}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	oldPrice := s.Price
	oldOriginalPrice := s.OriginalPrice
	s.Price = req.NewPrice
	s.OriginalPrice = req.NewOriginalPrice
	s.UpdatedAt = time.Now()
	if err := dao.ServiceImp.UpdateService(s); err != nil {
		LogError("更新服务价格失败", err)
		response := &AdminResponse{Code: -1, ErrorMsg: "更新服务价格失败: " + err.Error()}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := &AdminResponse{Code: 0, Data: map[string]interface{}{
		"serviceId":        s.Id,
		"serviceName":      s.Name,
		"oldPrice":         oldPrice,
		"newPrice":         s.Price,
		"oldOriginalPrice": oldOriginalPrice,
		"newOriginalPrice": s.OriginalPrice,
		"reason":           req.Reason,
		"adminId":          adminUserId,
		"message":          "服务价格更新成功",
	}}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getServiceStatusText 获取服务状态文本
func getServiceStatusText(status int) string {
	switch status {
	case 1:
		return "上架"
	case 0:
		return "下架"
	default:
		return "未知"
	}
}
