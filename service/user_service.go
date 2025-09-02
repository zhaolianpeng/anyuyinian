package service

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

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
	UserId    string `json:"userId"`
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
	UserId string `json:"userId"`
	Phone  string `json:"phone"`
	Code   string `json:"code"` // 验证码
}

// UpdateUserInfoRequest 更新用户信息请求
type UpdateUserInfoRequest struct {
	UserId    string `json:"userId"`
	NickName  string `json:"nickName,omitempty"`
	AvatarUrl string `json:"avatarUrl,omitempty"`
	Gender    int    `json:"gender,omitempty"`
}

// AddressRequest 地址请求
type AddressRequest struct {
	Id        int32  `json:"id,omitempty"`
	UserId    string `json:"userId"`
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
	UserId    string `json:"userId"`
	Name      string `json:"name"`
	IdCard    string `json:"idCard"`
	Phone     string `json:"phone"`
	Gender    int    `json:"gender"`
	Birthday  string `json:"birthday"`
	Relation  string `json:"relation"`
	IsDefault bool   `json:"isDefault"`
}

// PatientResponse 就诊人响应（包含计算出的年龄）
type PatientResponse struct {
	Id        int32     `json:"id"`
	UserId    string    `json:"userId"`
	Name      string    `json:"name"`
	IdCard    string    `json:"idCard"`
	Phone     string    `json:"phone"`
	Gender    int       `json:"gender"`
	Birthday  string    `json:"birthday"`
	Age       int       `json:"age"` // 根据身份证计算的年龄
	Relation  string    `json:"relation"`
	IsDefault int       `json:"isDefault"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// calculateAgeFromIdCard 根据身份证号计算年龄
func calculateAgeFromIdCard(idCard string) int {
	if len(idCard) != 18 {
		return 0
	}

	// 提取出生日期（身份证第7-14位）
	birthDateStr := idCard[6:14]
	birthYear, _ := strconv.Atoi(birthDateStr[:4])
	birthMonth, _ := strconv.Atoi(birthDateStr[4:6])
	birthDay, _ := strconv.Atoi(birthDateStr[6:8])

	// 获取当前日期
	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())
	currentDay := now.Day()

	// 计算年龄
	age := currentYear - birthYear

	// 如果今年的生日还没到，年龄减1
	if currentMonth < birthMonth || (currentMonth == birthMonth && currentDay < birthDay) {
		age--
	}

	return age
}

// convertPatientToResponse 将PatientModel转换为PatientResponse
func convertPatientToResponse(patient *model.PatientModel) *PatientResponse {
	age := 0
	if patient.IdCard != "" {
		age = calculateAgeFromIdCard(patient.IdCard)
		fmt.Printf("DEBUG: 患者 %s 身份证 %s 计算年龄 %d\n", patient.Name, patient.IdCard, age)
	}

	response := &PatientResponse{
		Id:        patient.Id,
		UserId:    patient.UserId,
		Name:      patient.Name,
		IdCard:    patient.IdCard,
		Phone:     patient.Phone,
		Gender:    patient.Gender,
		Birthday:  patient.Birthday,
		Age:       age,
		Relation:  patient.Relation,
		IsDefault: patient.IsDefault,
		Status:    patient.Status,
		CreatedAt: patient.CreatedAt,
		UpdatedAt: patient.UpdatedAt,
	}

	fmt.Printf("DEBUG: 转换后的响应包含年龄字段: %+v\n", response)
	return response
}

// GetUserInfoHandler 获取用户个人信息接口
func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理获取用户信息请求", map[string]interface{}{
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

	LogStep("开始查询用户信息", map[string]interface{}{
		"userId": userIdStr,
	})

	// 获取用户信息 - 使用新的UserId字段查询
	user, err := dao.UserImp.GetUserByUserId(userIdStr)
	if err != nil {
		LogError("数据库查询用户信息失败", err)
		response := &UserResponse{
			Code:     -1,
			ErrorMsg: "获取用户信息失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	LogStep("用户信息查询成功", map[string]interface{}{
		"userId":   user.UserId,
		"openId":   user.OpenId,
		"nickName": user.NickName,
	})

	userInfo := &UserInfo{
		Id:        user.Id,
		UserId:    user.UserId,
		OpenId:    user.OpenId,
		NickName:  user.NickName,
		AvatarUrl: user.AvatarUrl,
		Gender:    user.Gender,
		Phone:     user.Phone, // 从用户表获取手机号
		Country:   user.Country,
		Province:  user.Province,
		City:      user.City,
		Language:  user.Language,
	}

	response := &UserResponse{
		Code: 0,
		Data: userInfo,
	}

	LogStep("准备返回用户信息", map[string]interface{}{
		"userId":   userInfo.UserId,
		"nickName": userInfo.NickName,
		"phone":    userInfo.Phone,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	LogInfo("用户信息获取成功", map[string]interface{}{
		"userId": userInfo.UserId,
	})
}

// BindPhoneHandler 用户绑定手机号接口
func BindPhoneHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理绑定手机号请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodPost {
		LogError("请求方法不支持", fmt.Errorf("期望POST方法，实际为%s", r.Method))
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	var req BindPhoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogError("请求参数解析失败", err)
		http.Error(w, "请求参数解析失败", http.StatusBadRequest)
		return
	}

	LogStep("解析绑定手机号请求参数", map[string]interface{}{
		"userId": req.UserId,
		"phone":  req.Phone,
		"code":   req.Code,
	})

	// 验证参数
	if req.UserId == "" || req.Phone == "" || req.Code == "" {
		LogError("缺少必要参数", fmt.Errorf("userId=%s, phone=%s, code=%s", req.UserId, req.Phone, req.Code))
		http.Error(w, "缺少必要参数", http.StatusBadRequest)
		return
	}

	// 这里应该验证验证码，简化处理
	LogStep("开始验证验证码", map[string]interface{}{
		"inputCode":    req.Code,
		"expectedCode": "123456",
	})

	if req.Code != "123456" {
		LogError("验证码错误", fmt.Errorf("输入验证码=%s，期望验证码=123456", req.Code))
		response := &UserResponse{
			Code:     -1,
			ErrorMsg: "验证码错误",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	LogStep("验证码验证成功", nil)

	// 获取用户信息
	LogStep("开始查询用户信息", map[string]interface{}{
		"userId": req.UserId,
	})

	user, err := dao.UserImp.GetUserByUserId(req.UserId)
	if err != nil {
		LogError("数据库查询用户信息失败", err)
		response := &UserResponse{
			Code:     -1,
			ErrorMsg: "用户不存在: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	LogStep("用户信息查询成功", map[string]interface{}{
		"userId":   user.UserId,
		"nickName": user.NickName,
		"oldPhone": user.Phone,
	})

	// 更新用户手机号
	LogStep("开始更新用户手机号", map[string]interface{}{
		"userId":   user.UserId,
		"oldPhone": user.Phone,
		"newPhone": req.Phone,
	})

	user.Phone = req.Phone
	user.UpdatedAt = time.Now()

	if err := dao.UserImp.UpdateUser(user); err != nil {
		LogError("数据库更新用户手机号失败", err)
		response := &UserResponse{
			Code:     -1,
			ErrorMsg: "更新手机号失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	LogStep("用户手机号更新成功", map[string]interface{}{
		"userId":   user.UserId,
		"newPhone": req.Phone,
	})

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

	LogInfo("手机号绑定成功", map[string]interface{}{
		"userId": req.UserId,
		"phone":  req.Phone,
	})
}

// AddressHandler 地址管理接口
func AddressHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理地址管理请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	switch r.Method {
	case http.MethodGet:
		LogStep("处理获取地址列表请求", nil)
		handleGetAddresses(w, r)
	case http.MethodPost:
		LogStep("处理创建地址请求", nil)
		handleCreateAddress(w, r)
	case http.MethodPut:
		LogStep("处理更新地址请求", nil)
		handleUpdateAddress(w, r)
	case http.MethodDelete:
		LogStep("处理删除地址请求", nil)
		handleDeleteAddress(w, r)
	default:
		LogError("不支持的请求方法", fmt.Errorf("方法=%s", r.Method))
		http.Error(w, "不支持的请求方法", http.StatusMethodNotAllowed)
	}
}

// handleGetAddresses 获取地址列表
func handleGetAddresses(w http.ResponseWriter, r *http.Request) {
	LogStep("开始获取地址列表", nil)

	userIdStr := r.URL.Query().Get("userId")
	LogStep("解析请求参数", map[string]interface{}{
		"userId": userIdStr,
	})

	if userIdStr == "" {
		LogError("缺少必要参数", fmt.Errorf("userId参数为空"))
		http.Error(w, "缺少userId参数", http.StatusBadRequest)
		return
	}

	LogStep("开始查询地址列表", map[string]interface{}{
		"userId": userIdStr,
	})

	addresses, err := dao.UserExtendImp.GetAddressesByUserId(userIdStr)
	if err != nil {
		LogError("数据库查询地址列表失败", err)
		response := &UserResponse{
			Code:     -1,
			ErrorMsg: "获取地址列表失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	LogStep("地址列表查询成功", map[string]interface{}{
		"userId":       userIdStr,
		"addressCount": len(addresses),
	})

	response := &UserResponse{
		Code: 0,
		Data: addresses,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	LogInfo("地址列表获取成功", map[string]interface{}{
		"userId":       userIdStr,
		"addressCount": len(addresses),
	})
}

// handleCreateAddress 创建地址
func handleCreateAddress(w http.ResponseWriter, r *http.Request) {
	LogStep("开始创建地址", nil)

	var req AddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogError("请求参数解析失败", err)
		http.Error(w, "请求参数解析失败", http.StatusBadRequest)
		return
	}

	LogStep("解析创建地址请求参数", map[string]interface{}{
		"userId":    req.UserId,
		"name":      req.Name,
		"phone":     req.Phone,
		"province":  req.Province,
		"city":      req.City,
		"district":  req.District,
		"isDefault": req.IsDefault,
	})

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

	LogStep("开始创建地址到数据库", map[string]interface{}{
		"userId":    address.UserId,
		"name":      address.Name,
		"isDefault": address.IsDefault,
	})

	if err := dao.UserExtendImp.CreateAddress(address); err != nil {
		LogError("数据库创建地址失败", err)
		response := &UserResponse{
			Code:     -1,
			ErrorMsg: "创建地址失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	LogStep("地址创建成功", map[string]interface{}{
		"addressId": address.Id,
		"userId":    address.UserId,
	})

	// 如果设置为默认地址
	if req.IsDefault {
		LogStep("设置默认地址", map[string]interface{}{
			"userId":    req.UserId,
			"addressId": address.Id,
		})
		dao.UserExtendImp.SetDefaultAddress(req.UserId, address.Id)
		LogStep("默认地址设置完成", nil)
	}

	response := &UserResponse{
		Code: 0,
		Data: address,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	LogInfo("地址创建成功", map[string]interface{}{
		"addressId": address.Id,
		"userId":    address.UserId,
		"name":      address.Name,
	})
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

	patients, err := dao.UserExtendImp.GetPatientsByUserId(userIdStr)
	if err != nil {
		response := &UserResponse{
			Code:     -1,
			ErrorMsg: "获取就诊人列表失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	var patientResponses []*PatientResponse
	for _, patient := range patients {
		patientResponses = append(patientResponses, convertPatientToResponse(patient))
	}

	response := &UserResponse{
		Code: 0,
		Data: patientResponses,
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
		Data: convertPatientToResponse(patient),
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
		Data: convertPatientToResponse(patient),
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

// UpdateUserInfoHandler 更新用户信息
func UpdateUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理更新用户信息请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodPost {
		LogError("请求方法不支持", fmt.Errorf("期望POST方法，实际为%s", r.Method))
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	var req UpdateUserInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogError("请求参数解析失败", err)
		http.Error(w, "请求参数解析失败", http.StatusBadRequest)
		return
	}

	LogStep("解析更新用户信息请求参数", map[string]interface{}{
		"userId":    req.UserId,
		"nickName":  req.NickName,
		"avatarUrl": req.AvatarUrl,
		"gender":    req.Gender,
	})

	// 验证用户是否存在
	user, err := dao.UserImp.GetUserByUserId(req.UserId)
	if err != nil {
		LogError("查询用户失败", err)
		response := &UserResponse{
			Code:     -1,
			ErrorMsg: "用户不存在: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 更新用户信息
	if req.NickName != "" {
		user.NickName = req.NickName
	}
	if req.AvatarUrl != "" {
		user.AvatarUrl = req.AvatarUrl
	}
	if req.Gender > 0 {
		user.Gender = req.Gender
	}

	if err := dao.UserImp.UpdateUser(user); err != nil {
		LogError("更新用户信息失败", err)
		response := &UserResponse{
			Code:     -1,
			ErrorMsg: "更新用户信息失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	LogStep("用户信息更新成功", map[string]interface{}{
		"userId":    user.UserId,
		"nickName":  user.NickName,
		"avatarUrl": user.AvatarUrl,
		"gender":    user.Gender,
	})

	response := &UserResponse{
		Code: 0,
		Data: map[string]string{"message": "用户信息更新成功"},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	LogInfo("更新用户信息成功", map[string]interface{}{
		"userId": user.UserId,
	})
}

// DecryptPhoneNumberRequest 解密手机号请求
type DecryptPhoneNumberRequest struct {
	UserId        string `json:"userId"`
	EncryptedData string `json:"encryptedData"`
	IV            string `json:"iv"`
}

// DecryptPhoneNumberResponse 解密手机号响应
type DecryptPhoneNumberResponse struct {
	PhoneNumber string `json:"phoneNumber"`
}

// DecryptPhoneNumberHandler 解密微信手机号接口
func DecryptPhoneNumberHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理解密手机号请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodPost {
		LogError("请求方法不支持", fmt.Errorf("期望POST方法，实际为%s", r.Method))
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	var req DecryptPhoneNumberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogError("请求参数解析失败", err)
		http.Error(w, "请求参数解析失败", http.StatusBadRequest)
		return
	}

	LogStep("解析解密手机号请求参数", map[string]interface{}{
		"userId":        req.UserId,
		"encryptedData": req.EncryptedData[:20] + "...", // 只显示前20个字符
		"iv":            req.IV,
	})

	// 验证参数
	if req.UserId == "" || req.EncryptedData == "" || req.IV == "" {
		LogError("缺少必要参数", fmt.Errorf("userId=%s, encryptedData长度=%d, iv=%s",
			req.UserId, len(req.EncryptedData), req.IV))
		http.Error(w, "缺少必要参数", http.StatusBadRequest)
		return
	}

	// 验证用户是否存在
	user, err := dao.UserImp.GetUserByUserId(req.UserId)
	if err != nil {
		LogError("查询用户失败", err)
		response := &UserResponse{
			Code:     -1,
			ErrorMsg: "用户不存在: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	LogStep("用户验证成功", map[string]interface{}{
		"userId":     user.UserId,
		"nickName":   user.NickName,
		"sessionKey": user.SessionKey[:10] + "...", // 只显示前10个字符
	})

	// 检查session_key是否存在
	if user.SessionKey == "" {
		LogError("用户session_key为空", fmt.Errorf("userId=%s, sessionKey为空", req.UserId))
		response := &UserResponse{
			Code:     -1,
			ErrorMsg: "用户session_key不存在，请重新登录",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 使用微信解密算法解密手机号
	LogStep("开始解密手机号", map[string]interface{}{
		"userId":        req.UserId,
		"encryptedData": req.EncryptedData[:20] + "...",
		"iv":            req.IV,
		"sessionKey":    user.SessionKey[:10] + "...",
	})

	// 调用微信解密算法，使用session_key
	phoneNumber, err := decryptWechatPhoneNumber(req.EncryptedData, req.IV, user.SessionKey)
	if err != nil {
		LogError("微信手机号解密失败", err)
		response := &UserResponse{
			Code:     -1,
			ErrorMsg: "手机号解密失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	LogStep("手机号解密成功", map[string]interface{}{
		"userId":      req.UserId,
		"phoneNumber": phoneNumber,
	})

	response := &UserResponse{
		Code: 0,
		Data: &DecryptPhoneNumberResponse{
			PhoneNumber: phoneNumber,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	LogInfo("手机号解密成功", map[string]interface{}{
		"userId":      req.UserId,
		"phoneNumber": phoneNumber,
	})
}

// decryptWechatPhoneNumber 解密微信手机号
func decryptWechatPhoneNumber(encryptedData, iv, sessionKey string) (string, error) {
	// 解码Base64数据
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", fmt.Errorf("解码加密数据失败: %v", err)
	}

	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return "", fmt.Errorf("解码IV失败: %v", err)
	}

	// 解码session_key
	key, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return "", fmt.Errorf("解码session_key失败: %v", err)
	}

	// 创建AES解密器
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("创建AES解密器失败: %v", err)
	}

	// 检查数据长度
	if len(encryptedBytes)%aes.BlockSize != 0 {
		return "", fmt.Errorf("加密数据长度不是AES块大小的倍数")
	}

	// 创建CBC模式解密器
	mode := cipher.NewCBCDecrypter(block, ivBytes)

	// 解密数据
	decrypted := make([]byte, len(encryptedBytes))
	mode.CryptBlocks(decrypted, encryptedBytes)

	// 去除PKCS7填充
	decrypted = removePKCS7Padding(decrypted)

	// 解析JSON数据
	var phoneData struct {
		PhoneNumber     string `json:"phoneNumber"`
		PurePhoneNumber string `json:"purePhoneNumber"`
		CountryCode     string `json:"countryCode"`
	}

	if err := json.Unmarshal(decrypted, &phoneData); err != nil {
		return "", fmt.Errorf("解析解密数据失败: %v", err)
	}

	// 返回手机号（优先使用purePhoneNumber）
	if phoneData.PurePhoneNumber != "" {
		return phoneData.PurePhoneNumber, nil
	}
	return phoneData.PhoneNumber, nil
}

// removePKCS7Padding 移除PKCS7填充
func removePKCS7Padding(data []byte) []byte {
	if len(data) == 0 {
		return data
	}

	padding := int(data[len(data)-1])
	if padding > len(data) || padding == 0 {
		return data
	}

	// 验证填充是否正确
	for i := len(data) - padding; i < len(data); i++ {
		if data[i] != byte(padding) {
			return data
		}
	}

	return data[:len(data)-padding]
}
