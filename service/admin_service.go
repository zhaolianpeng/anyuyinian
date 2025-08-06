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

	// 根据管理员级别获取不同的数据
	if admin.AdminLevel == 2 { // 超级管理员
		dbCli.Model(&model.UserModel{}).Count(&totalUsers)
		dbCli.Model(&model.OrderModel{}).Count(&totalOrders)
		dbCli.Model(&model.OrderModel{}).Where("DATE(createdAt) = CURDATE()").Count(&todayOrders)
		dbCli.Model(&model.OrderModel{}).Select("IFNULL(SUM(totalAmount),0)").Row().Scan(&totalAmount)
		// 已支付总金额
		dbCli.Model(&model.OrderModel{}).Where("status = 'paid'").Select("IFNULL(SUM(totalAmount),0)").Row().Scan(&paidAmount)
		// 待支付总金额
		dbCli.Model(&model.OrderModel{}).Where("status = 'pending'").Select("IFNULL(SUM(totalAmount),0)").Row().Scan(&unpaidAmount)
		// 退款总金额
		dbCli.Model(&model.OrderModel{}).Where("status = 'refunded'").Select("IFNULL(SUM(totalAmount),0)").Row().Scan(&refundAmount)
	} else { // 一级管理员
		// 获取该管理员推广的用户ID列表
		var promotedUserIds []string
		dbCli.Model(&model.UserModel{}).Where("referrerId = ?", adminUserId).Pluck("userId", &promotedUserIds)
		promotedUserIds = append(promotedUserIds, adminUserId) // 包含管理员自己的用户ID

		dbCli.Model(&model.UserModel{}).Where("userId IN (?)", promotedUserIds).Count(&totalUsers)
		dbCli.Model(&model.OrderModel{}).Where("userId IN (?)", promotedUserIds).Count(&totalOrders)
		dbCli.Model(&model.OrderModel{}).Where("userId IN (?) AND DATE(createdAt) = CURDATE()", promotedUserIds).Count(&todayOrders)
		dbCli.Model(&model.OrderModel{}).Where("userId IN (?)", promotedUserIds).Select("IFNULL(SUM(totalAmount),0)").Row().Scan(&totalAmount)
		// 已支付总金额
		dbCli.Model(&model.OrderModel{}).Where("userId IN (?) AND status = 'paid'", promotedUserIds).Select("IFNULL(SUM(totalAmount),0)").Row().Scan(&paidAmount)
		// 待支付总金额
		dbCli.Model(&model.OrderModel{}).Where("userId IN (?) AND status = 'pending'", promotedUserIds).Select("IFNULL(SUM(totalAmount),0)").Row().Scan(&unpaidAmount)
		// 退款总金额
		dbCli.Model(&model.OrderModel{}).Where("userId IN (?) AND status = 'refunded'", promotedUserIds).Select("IFNULL(SUM(totalAmount),0)").Row().Scan(&refundAmount)
	}

	stats := map[string]interface{}{
		"totalUsers":   totalUsers,
		"totalOrders":  totalOrders,
		"todayOrders":  todayOrders,
		"totalAmount":  totalAmount,
		"paidAmount":   paidAmount,
		"unpaidAmount": unpaidAmount,
		"refundAmount": refundAmount,
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
