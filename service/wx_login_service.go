package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"wxcloudrun-golang/config"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"

	"gorm.io/gorm"
)

// WxLoginRequest 微信登录请求
type WxLoginRequest struct {
	Code      string `json:"code"`
	NickName  string `json:"nickName,omitempty"`
	AvatarUrl string `json:"avatarUrl,omitempty"`
	Gender    int    `json:"gender,omitempty"`
	Country   string `json:"country,omitempty"`
	Province  string `json:"province,omitempty"`
	City      string `json:"city,omitempty"`
	Language  string `json:"language,omitempty"`
}

// WxLoginResponse 微信登录响应
type WxLoginResponse struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// WxLoginResult 登录结果
type WxLoginResult struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

// WxLoginHandler 微信小程序登录接口
func WxLoginHandler(w http.ResponseWriter, r *http.Request) {
	LogStep("开始处理微信登录请求", map[string]string{"method": r.Method, "path": r.URL.Path})

	if r.Method != http.MethodPost {
		LogError("请求方法不支持", fmt.Errorf("期望POST方法，实际为%s", r.Method))
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体
	var req WxLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogError("请求参数解析失败", err)
		http.Error(w, "请求参数解析失败", http.StatusBadRequest)
		return
	}
	LogRequest("POST", "/api/wx/login", req)

	// 验证必要参数
	if req.Code == "" {
		LogError("缺少code参数", fmt.Errorf("code参数为空"))
		http.Error(w, "缺少code参数", http.StatusBadRequest)
		return
	}

	LogStep("开始调用微信API", map[string]string{"code": req.Code})
	// 调用微信API获取用户信息
	wxResp, err := getWxSession(req.Code)
	if err != nil {
		LogError("微信API调用失败", err)
		http.Error(w, "微信API调用失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	LogStep("微信API调用成功", wxResp)

	// 检查微信API返回错误
	if wxResp.ErrCode != 0 {
		LogError("微信API返回错误", fmt.Errorf("错误码: %d, 错误信息: %s", wxResp.ErrCode, wxResp.ErrMsg))
		http.Error(w, "微信API错误: "+wxResp.ErrMsg, http.StatusBadRequest)
		return
	}

	// 检查openId是否有效
	if wxResp.OpenId == "" {
		LogError("微信API返回的openId为空", fmt.Errorf("openId为空"))
		http.Error(w, "微信API返回的openId为空", http.StatusBadRequest)
		return
	}

	// 检查是否为模拟数据
	if strings.HasPrefix(wxResp.OpenId, "user_") {
		LogStep("检测到模拟openId，使用模拟数据", map[string]string{"openId": wxResp.OpenId})
		// 对于模拟数据，直接返回成功响应
		result := &WxLoginResult{
			Code: 0,
			Data: map[string]interface{}{
				"id":          1,
				"openId":      wxResp.OpenId,
				"nickName":    req.NickName,
				"avatarUrl":   req.AvatarUrl,
				"gender":      req.Gender,
				"country":     req.Country,
				"province":    req.Province,
				"city":        req.City,
				"language":    req.Language,
				"lastLoginAt": time.Now(),
				"isNewUser":   false,
				"token":       generateToken(1),
				"userId":      1,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
		LogResponse(result, nil)
		return
	}

	LogStep("开始处理用户登录", map[string]string{"openId": wxResp.OpenId})
	// 处理用户登录
	result, err := processUserLogin(wxResp, &req)
	if err != nil {
		LogError("用户登录处理失败", err)
		http.Error(w, "用户登录处理失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回结果
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
	LogResponse(result, nil)
}

// getWxSession 调用微信API获取session_key和openid
func getWxSession(code string) (*WxLoginResponse, error) {
	LogStep("开始获取微信配置", nil)
	// 获取微信配置
	wxConfig := config.GetWxConfig()
	appID := wxConfig.AppID
	appSecret := wxConfig.AppSecret

	// 添加配置验证
	if appID == "" || appSecret == "" {
		LogError("微信配置为空", fmt.Errorf("AppID或AppSecret未配置"))
		return nil, fmt.Errorf("微信配置未正确设置")
	}

	// 检查是否为默认配置
	if appID == "your_app_id" || appSecret == "your_app_secret" {
		LogError("微信配置为默认值", fmt.Errorf("请配置真实的微信AppID和AppSecret"))
		return nil, fmt.Errorf("微信配置为默认值，请设置真实配置")
	}

	LogStep("微信配置获取成功", map[string]string{"appID": appID})

	// 构建请求URL
	baseURL := "https://api.weixin.qq.com/sns/jscode2session"
	params := url.Values{}
	params.Add("appid", appID)
	params.Add("secret", appSecret)
	params.Add("js_code", code)
	params.Add("grant_type", "authorization_code")

	requestURL := baseURL + "?" + params.Encode()
	LogStep("构建微信API请求", map[string]string{"url": requestURL})

	// 发送HTTP请求
	resp, err := http.Get(requestURL)
	if err != nil {
		LogError("微信API请求失败", err)
		return nil, fmt.Errorf("微信API请求失败: %v", err)
	}
	defer resp.Body.Close()
	LogStep("微信API请求发送成功", map[string]int{"status_code": resp.StatusCode})

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		LogError("读取微信API响应失败", err)
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}
	LogStep("微信API响应读取成功", map[string]string{"response_body": string(body)})

	// 解析响应
	var wxResp WxLoginResponse
	if err := json.Unmarshal(body, &wxResp); err != nil {
		LogError("解析微信API响应失败", err)
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}
	LogStep("微信API响应解析成功", wxResp)

	return &wxResp, nil
}

// processUserLogin 处理用户登录逻辑
func processUserLogin(wxResp *WxLoginResponse, req *WxLoginRequest) (*WxLoginResult, error) {
	LogStep("开始查询用户是否存在", map[string]string{"openId": wxResp.OpenId})
	// 查询用户是否已存在
	existingUser, err := dao.UserImp.GetUserByOpenId(wxResp.OpenId)
	LogDBResult("查询", "users", existingUser, err)

	if err != nil && err != gorm.ErrRecordNotFound {
		LogError("查询用户失败", err)
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	var user *model.UserModel
	isNewUser := false

	if err == gorm.ErrRecordNotFound {
		LogStep("用户不存在，开始创建新用户", map[string]string{"openId": wxResp.OpenId})
		// 新用户，创建用户记录
		user = &model.UserModel{
			OpenId:     wxResp.OpenId,
			UnionId:    wxResp.UnionId,
			SessionKey: wxResp.SessionKey,
			NickName:   req.NickName,
			AvatarUrl:  req.AvatarUrl,
			Gender:     req.Gender,
			Country:    req.Country,
			Province:   req.Province,
			City:       req.City,
			Language:   req.Language,
		}

		LogDBOperation("创建", "users", user)
		if err := dao.UserImp.CreateUser(user); err != nil {
			LogError("创建用户失败", err)
			return nil, fmt.Errorf("创建用户失败: %v", err)
		}
		LogDBResult("创建", "users", user, nil)
		isNewUser = true
		LogStep("新用户创建成功", map[string]interface{}{"userId": user.Id, "isNewUser": isNewUser})
	} else {
		LogStep("用户已存在，开始更新用户信息", map[string]interface{}{"userId": existingUser.Id, "openId": existingUser.OpenId})
		// 老用户，更新登录时间和session_key
		existingUser.SessionKey = wxResp.SessionKey
		existingUser.LastLoginAt = time.Now()
		existingUser.UpdatedAt = time.Now()

		// 如果请求中包含用户信息，则更新
		updateFields := make(map[string]interface{})
		if req.NickName != "" {
			existingUser.NickName = req.NickName
			updateFields["nickName"] = req.NickName
		}
		if req.AvatarUrl != "" {
			existingUser.AvatarUrl = req.AvatarUrl
			updateFields["avatarUrl"] = req.AvatarUrl
		}
		if req.Gender != 0 {
			existingUser.Gender = req.Gender
			updateFields["gender"] = req.Gender
		}
		if req.Country != "" {
			existingUser.Country = req.Country
			updateFields["country"] = req.Country
		}
		if req.Province != "" {
			existingUser.Province = req.Province
			updateFields["province"] = req.Province
		}
		if req.City != "" {
			existingUser.City = req.City
			updateFields["city"] = req.City
		}
		if req.Language != "" {
			existingUser.Language = req.Language
			updateFields["language"] = req.Language
		}

		LogStep("准备更新用户字段", updateFields)
		LogDBOperation("更新", "users", existingUser)
		if err := dao.UserImp.UpdateUser(existingUser); err != nil {
			LogError("更新用户失败", err)
			return nil, fmt.Errorf("更新用户失败: %v", err)
		}
		LogDBResult("更新", "users", existingUser, nil)
		user = existingUser
		LogStep("用户信息更新成功", map[string]interface{}{"userId": user.Id, "isNewUser": isNewUser})
	}

	LogStep("开始构建返回数据", nil)
	// 构建返回数据（不包含敏感信息如session_key）
	userData := map[string]interface{}{
		"id":          user.Id,
		"openId":      user.OpenId,
		"nickName":    user.NickName,
		"avatarUrl":   user.AvatarUrl,
		"gender":      user.Gender,
		"country":     user.Country,
		"province":    user.Province,
		"city":        user.City,
		"language":    user.Language,
		"lastLoginAt": user.LastLoginAt,
		"isNewUser":   isNewUser,
		"token":       generateToken(user.Id), // 添加token生成
		"userId":      user.Id,                // 添加userId
	}

	result := &WxLoginResult{
		Code: 0,
		Data: userData,
	}
	LogStep("登录处理完成", result)
	return result, nil
}

// generateToken 生成简单的token
func generateToken(userId int32) string {
	return fmt.Sprintf("token_%d_%d", userId, time.Now().Unix())
}
