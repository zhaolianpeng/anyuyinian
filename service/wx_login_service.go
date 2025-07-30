package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体
	var req WxLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "请求参数解析失败", http.StatusBadRequest)
		return
	}

	// 验证必要参数
	if req.Code == "" {
		http.Error(w, "缺少code参数", http.StatusBadRequest)
		return
	}

	// 调用微信API获取用户信息
	wxResp, err := getWxSession(req.Code)
	if err != nil {
		http.Error(w, "微信API调用失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 检查微信API返回错误
	if wxResp.ErrCode != 0 {
		http.Error(w, "微信API错误: "+wxResp.ErrMsg, http.StatusBadRequest)
		return
	}

	// 处理用户登录
	result, err := processUserLogin(wxResp, &req)
	if err != nil {
		http.Error(w, "用户登录处理失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回结果
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// getWxSession 调用微信API获取session_key和openid
func getWxSession(code string) (*WxLoginResponse, error) {
	// 获取微信配置
	wxConfig := config.GetWxConfig()
	appID := wxConfig.AppID
	appSecret := wxConfig.AppSecret

	// 构建请求URL
	baseURL := "https://api.weixin.qq.com/sns/jscode2session"
	params := url.Values{}
	params.Add("appid", appID)
	params.Add("secret", appSecret)
	params.Add("js_code", code)
	params.Add("grant_type", "authorization_code")

	// 发送HTTP请求
	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("微信API请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析响应
	var wxResp WxLoginResponse
	if err := json.Unmarshal(body, &wxResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &wxResp, nil
}

// processUserLogin 处理用户登录逻辑
func processUserLogin(wxResp *WxLoginResponse, req *WxLoginRequest) (*WxLoginResult, error) {
	// 查询用户是否已存在
	existingUser, err := dao.UserImp.GetUserByOpenId(wxResp.OpenId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	var user *model.UserModel
	isNewUser := false

	if err == gorm.ErrRecordNotFound {
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

		if err := dao.UserImp.CreateUser(user); err != nil {
			return nil, fmt.Errorf("创建用户失败: %v", err)
		}
		isNewUser = true
	} else {
		// 老用户，更新登录时间和session_key
		existingUser.SessionKey = wxResp.SessionKey
		existingUser.LastLoginAt = time.Now()
		existingUser.UpdatedAt = time.Now()

		// 如果请求中包含用户信息，则更新
		if req.NickName != "" {
			existingUser.NickName = req.NickName
		}
		if req.AvatarUrl != "" {
			existingUser.AvatarUrl = req.AvatarUrl
		}
		if req.Gender != 0 {
			existingUser.Gender = req.Gender
		}
		if req.Country != "" {
			existingUser.Country = req.Country
		}
		if req.Province != "" {
			existingUser.Province = req.Province
		}
		if req.City != "" {
			existingUser.City = req.City
		}
		if req.Language != "" {
			existingUser.Language = req.Language
		}

		if err := dao.UserImp.UpdateUser(existingUser); err != nil {
			return nil, fmt.Errorf("更新用户失败: %v", err)
		}
		user = existingUser
	}

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
	}

	return &WxLoginResult{
		Code: 0,
		Data: userData,
	}, nil
}
