package service

import (
	"encoding/json"
	"log"
	"net/http"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"
)

// EmergencyFixResponse 紧急修复响应
type EmergencyFixResponse struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

// EmergencyFixUserIdsHandler 紧急修复用户UserId
func EmergencyFixUserIdsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	log.Println("开始紧急修复用户UserId...")

	cli := db.Get()

	// 1. 查询所有没有userId的用户
	var users []*model.UserModel
	err := cli.Table("Users").Where("userId IS NULL OR userId = ''").Find(&users).Error
	if err != nil {
		log.Printf("查询用户失败: %v", err)
		response := &EmergencyFixResponse{
			Code:     -1,
			ErrorMsg: "查询用户失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	log.Printf("找到 %d 个需要修复的用户", len(users))

	// 2. 为每个用户生成userId
	fixedCount := 0
	for _, user := range users {
		// 生成新的userId
		newUserId := utils.GenerateUserID()

		// 更新用户记录
		err := cli.Table("Users").Where("id = ?", user.Id).Update("userId", newUserId).Error
		if err != nil {
			log.Printf("更新用户 %d 的userId失败: %v", user.Id, err)
			continue
		}

		log.Printf("用户 %d 的userId已更新为: %s", user.Id, newUserId)
		fixedCount++
	}

	// 3. 验证修复结果
	var totalUsers, usersWithUserId int64
	cli.Table("Users").Count(&totalUsers)
	cli.Table("Users").Where("userId IS NOT NULL AND userId != ''").Count(&usersWithUserId)

	response := &EmergencyFixResponse{
		Code: 0,
		Data: map[string]interface{}{
			"message":        "紧急修复完成",
			"fixedCount":     fixedCount,
			"totalUsers":     totalUsers,
			"usersWithUserId": usersWithUserId,
			"success":        totalUsers == usersWithUserId,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	log.Printf("紧急修复完成，修复了 %d 个用户", fixedCount)
}

// TestUserInfoHandler 测试用户信息API
func TestUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	// 获取测试userId参数
	userId := r.URL.Query().Get("userId")
	if userId == "" {
		userId = "1" // 默认测试userId
	}

	log.Printf("测试用户信息API，userId: %s", userId)

	// 尝试获取用户信息
	user, err := dao.UserImp.GetUserByUserId(userId)
	if err != nil {
		log.Printf("获取用户信息失败: %v", err)
		response := &EmergencyFixResponse{
			Code:     -1,
			ErrorMsg: "获取用户信息失败: " + err.Error(),
			Data: map[string]interface{}{
				"userId": userId,
				"error":  err.Error(),
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := &EmergencyFixResponse{
		Code: 0,
		Data: map[string]interface{}{
			"message": "用户信息获取成功",
			"user": map[string]interface{}{
				"id":     user.Id,
				"userId": user.UserId,
				"openId": user.OpenId,
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	log.Printf("用户信息API测试成功，用户ID: %d, UserId: %s", user.Id, user.UserId)
}

// GetUserStatusHandler 获取用户状态信息
func GetUserStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	log.Println("获取用户状态信息...")

	cli := db.Get()

	var totalUsers, usersWithUserId, usersWithoutUserId int64

	// 查询总用户数
	cli.Table("Users").Count(&totalUsers)

	// 查询有userId的用户数
	cli.Table("Users").Where("userId IS NOT NULL AND userId != ''").Count(&usersWithUserId)

	// 查询没有userId的用户数
	cli.Table("Users").Where("userId IS NULL OR userId = ''").Count(&usersWithoutUserId)

	// 查询前几个用户的详细信息
	var sampleUsers []*model.UserModel
	cli.Table("Users").Limit(5).Find(&sampleUsers)

	userDetails := make([]map[string]interface{}, 0)
	for _, user := range sampleUsers {
		userDetails = append(userDetails, map[string]interface{}{
			"id":     user.Id,
			"userId": user.UserId,
			"openId": user.OpenId,
			"status": func() string {
				if user.UserId == "" {
					return "需要修复"
				}
				return "正常"
			}(),
		})
	}

	response := &EmergencyFixResponse{
		Code: 0,
		Data: map[string]interface{}{
			"totalUsers":        totalUsers,
			"usersWithUserId":   usersWithUserId,
			"usersWithoutUserId": usersWithoutUserId,
			"fixNeeded":         usersWithoutUserId > 0,
			"sampleUsers":       userDetails,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	log.Printf("用户状态: 总数=%d, 有userId=%d, 无userId=%d", 
		totalUsers, usersWithUserId, usersWithoutUserId)
} 