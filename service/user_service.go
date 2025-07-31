package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

// UserResponse 用户响应
type UserResponse struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

// UserInfo 用户信息
type UserInfo struct {
	Id        int32  `json:"id"`
	OpenId    string `json:"openId"`
	NickName  string `json:"nickName"`
	AvatarUrl string `json:"avatarUrl"`
	Gender    int    `json:"gender"`
	Phone     string `json:"phone"`
	Country   string `json:"country"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Language  string `json:"language"`
}

// BindPhoneRequest 绑定手机号请求
type BindPhoneRequest struct {
	UserId int32  `json:"userId"`
	Phone  string `json:"phone"`
	Code   string `json:"code"` // 验证码
}

// AddressRequest 地址请求
type AddressRequest struct {
	Id        int32  `json:"id,omitempty"`
	UserId    int32  `json:"userId"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Province  string `json:"province"`
	City      string `json:"city"`
	District  string `json:"district"`
	Address   string `json:"address"`
	IsDefault bool   `json:"isDefault"`
}

// PatientRequest 就诊人请求
type PatientRequest struct {
	Id        int32  `json:"id,omitempty"`
	UserId    int32  `json:"userId"`
	Name      string `json:"name"`
	IdCard    string `json:"idCard"`
	Phone     string `json:"phone"`
	Gender    int    `json:"gender"`
	Birthday  string `json:"birthday"`
	Relation  string `json:"relation"`
	IsDefault bool   `json:"isDefault"`
}

// GetUserInfoHandler 获取用户个人信息接口
func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
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

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "无效的用户ID", http.StatusBadRequest)
		return
	}

	// 检查是否为模拟用户ID
	if userId == 1 {
		// 返回模拟用户信息
		userInfo := &UserInfo{
			Id:        1,
			OpenId:    "user_1",
			NickName:  "微信用户",
			AvatarUrl: "https://thirdwx.qlogo.cn/mmopen/vi_32/POgEwh4mIHO4nibH0KlMECNjjGxQUq24ZEaGT4poC6icRiccVGKSyXwibcPq4BWmiaIGuG1icwxaQX6grC9VemZoJ8rg/132",
			Gender:    0,
			Phone:     "",
			Country:   "China",
			Province:  "Guangdong",
			City:      "Shenzhen",
			Language:  "zh_CN",
		}

		response := &UserResponse{
			Code: 0,
			Data: userInfo,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 获取用户信息
	user, err := dao.UserImp.GetUserByOpenId(fmt.Sprintf("user_%d", userId)) // 这里简化处理
	if err != nil {
		response := &UserResponse{
			Code:     -1,
			ErrorMsg: "获取用户信息失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	userInfo := &UserInfo{
		Id:        user.Id,
		OpenId:    user.OpenId,
		NickName:  user.NickName,
		AvatarUrl: user.AvatarUrl,
		Gender:    user.Gender,
		Phone:     "", // 这里需要从其他表获取手机号
		Country:   user.Country,
		Province:  user.Province,
		City:      user.City,
		Language:  user.Language,
	}

	response := &UserResponse{
		Code: 0,
		Data: userInfo,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// BindPhoneHandler 用户绑定手机号接口
func BindPhoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	var req BindPhoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "请求参数解析失败", http.StatusBadRequest)
		return
	}

	// 验证参数
	if req.UserId == 0 || req.Phone == "" || req.Code == "" {
		http.Error(w, "缺少必要参数", http.StatusBadRequest)
		return
	}

	// 这里应该验证验证码，简化处理
	if req.Code != "123456" {
		response := &UserResponse{
			Code:     -1,
			ErrorMsg: "验证码错误",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 更新用户手机号（这里简化处理，实际应该更新用户表或创建手机号关联表）
	response := &UserResponse{
		Code: 0,
		Data: map[string]interface{}{
			"userId":  req.UserId,
			"phone":   req.Phone,
			"message": "手机号绑定成功",
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AddressHandler 地址管理接口
func AddressHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetAddresses(w, r)
	case http.MethodPost:
		handleCreateAddress(w, r)
	case http.MethodPut:
		handleUpdateAddress(w, r)
	case http.MethodDelete:
		handleDeleteAddress(w, r)
	default:
		http.Error(w, "不支持的请求方法", http.StatusMethodNotAllowed)
	}
}

// handleGetAddresses 获取地址列表
func handleGetAddresses(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.URL.Query().Get("userId")
	if userIdStr == "" {
		http.Error(w, "缺少userId参数", http.StatusBadRequest)
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "无效的用户ID", http.StatusBadRequest)
		return
	}

	addresses, err := dao.UserExtendImp.GetAddressesByUserId(int32(userId))
	if err != nil {
		response := &UserResponse{
			Code:     -1,
			ErrorMsg: "获取地址列表失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := &UserResponse{
		Code: 0,
		Data: addresses,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleCreateAddress 创建地址
func handleCreateAddress(w http.ResponseWriter, r *http.Request) {
	var req AddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "请求参数解析失败", http.StatusBadRequest)
		return
	}

	address := &model.UserAddressModel{
		UserId:    req.UserId,
		Name:      req.Name,
		Phone:     req.Phone,
		Province:  req.Province,
		City:      req.City,
		District:  req.District,
		Address:   req.Address,
		IsDefault: 0,
	}

	if req.IsDefault {
		address.IsDefault = 1
	}

	if err := dao.UserExtendImp.CreateAddress(address); err != nil {
		response := &UserResponse{
			Code:     -1,
			ErrorMsg: "创建地址失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 如果设置为默认地址
	if req.IsDefault {
		dao.UserExtendImp.SetDefaultAddress(req.UserId, address.Id)
	}

	response := &UserResponse{
		Code: 0,
		Data: address,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleUpdateAddress 更新地址
func handleUpdateAddress(w http.ResponseWriter, r *http.Request) {
	var req AddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "请求参数解析失败", http.StatusBadRequest)
		return
	}

	address := &model.UserAddressModel{
		Id:        req.Id,
		UserId:    req.UserId,
		Name:      req.Name,
		Phone:     req.Phone,
		Province:  req.Province,
		City:      req.City,
		District:  req.District,
		Address:   req.Address,
		IsDefault: 0,
	}

	if req.IsDefault {
		address.IsDefault = 1
	}

	if err := dao.UserExtendImp.UpdateAddress(address); err != nil {
		response := &UserResponse{
			Code:     -1,
			ErrorMsg: "更新地址失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 如果设置为默认地址
	if req.IsDefault {
		dao.UserExtendImp.SetDefaultAddress(req.UserId, address.Id)
	}

	response := &UserResponse{
		Code: 0,
		Data: address,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleDeleteAddress 删除地址
func handleDeleteAddress(w http.ResponseWriter, r *http.Request) {
	addressIdStr := r.URL.Query().Get("id")
	if addressIdStr == "" {
		http.Error(w, "缺少地址ID参数", http.StatusBadRequest)
		return
	}

	addressId, err := strconv.Atoi(addressIdStr)
	if err != nil {
		http.Error(w, "无效的地址ID", http.StatusBadRequest)
		return
	}

	if err := dao.UserExtendImp.DeleteAddress(int32(addressId)); err != nil {
		response := &UserResponse{
			Code:     -1,
			ErrorMsg: "删除地址失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := &UserResponse{
		Code: 0,
		Data: map[string]string{"message": "地址删除成功"},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// PatientHandler 就诊人管理接口
func PatientHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetPatients(w, r)
	case http.MethodPost:
		handleCreatePatient(w, r)
	case http.MethodPut:
		handleUpdatePatient(w, r)
	case http.MethodDelete:
		handleDeletePatient(w, r)
	default:
		http.Error(w, "不支持的请求方法", http.StatusMethodNotAllowed)
	}
}

// handleGetPatients 获取就诊人列表
func handleGetPatients(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.URL.Query().Get("userId")
	if userIdStr == "" {
		http.Error(w, "缺少userId参数", http.StatusBadRequest)
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "无效的用户ID", http.StatusBadRequest)
		return
	}

	patients, err := dao.UserExtendImp.GetPatientsByUserId(int32(userId))
	if err != nil {
		response := &UserResponse{
			Code:     -1,
			ErrorMsg: "获取就诊人列表失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := &UserResponse{
		Code: 0,
		Data: patients,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleCreatePatient 创建就诊人
func handleCreatePatient(w http.ResponseWriter, r *http.Request) {
	var req PatientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "请求参数解析失败", http.StatusBadRequest)
		return
	}

	patient := &model.PatientModel{
		UserId:    req.UserId,
		Name:      req.Name,
		IdCard:    req.IdCard,
		Phone:     req.Phone,
		Gender:    req.Gender,
		Birthday:  req.Birthday,
		Relation:  req.Relation,
		IsDefault: 0,
	}

	if req.IsDefault {
		patient.IsDefault = 1
	}

	if err := dao.UserExtendImp.CreatePatient(patient); err != nil {
		response := &UserResponse{
			Code:     -1,
			ErrorMsg: "创建就诊人失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 如果设置为默认就诊人
	if req.IsDefault {
		dao.UserExtendImp.SetDefaultPatient(req.UserId, patient.Id)
	}

	response := &UserResponse{
		Code: 0,
		Data: patient,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleUpdatePatient 更新就诊人
func handleUpdatePatient(w http.ResponseWriter, r *http.Request) {
	var req PatientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "请求参数解析失败", http.StatusBadRequest)
		return
	}

	patient := &model.PatientModel{
		Id:        req.Id,
		UserId:    req.UserId,
		Name:      req.Name,
		IdCard:    req.IdCard,
		Phone:     req.Phone,
		Gender:    req.Gender,
		Birthday:  req.Birthday,
		Relation:  req.Relation,
		IsDefault: 0,
	}

	if req.IsDefault {
		patient.IsDefault = 1
	}

	if err := dao.UserExtendImp.UpdatePatient(patient); err != nil {
		response := &UserResponse{
			Code:     -1,
			ErrorMsg: "更新就诊人失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 如果设置为默认就诊人
	if req.IsDefault {
		dao.UserExtendImp.SetDefaultPatient(req.UserId, patient.Id)
	}

	response := &UserResponse{
		Code: 0,
		Data: patient,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleDeletePatient 删除就诊人
func handleDeletePatient(w http.ResponseWriter, r *http.Request) {
	patientIdStr := r.URL.Query().Get("id")
	if patientIdStr == "" {
		http.Error(w, "缺少就诊人ID参数", http.StatusBadRequest)
		return
	}

	patientId, err := strconv.Atoi(patientIdStr)
	if err != nil {
		http.Error(w, "无效的就诊人ID", http.StatusBadRequest)
		return
	}

	if err := dao.UserExtendImp.DeletePatient(int32(patientId)); err != nil {
		response := &UserResponse{
			Code:     -1,
			ErrorMsg: "删除就诊人失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := &UserResponse{
		Code: 0,
		Data: map[string]string{"message": "就诊人删除成功"},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
