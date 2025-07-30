package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

// HomeInitRequest 首页初始化请求
type HomeInitRequest struct {
	Longitude float64 `json:"longitude,omitempty"` // 经度
	Latitude  float64 `json:"latitude,omitempty"`  // 纬度
	Limit     int     `json:"limit,omitempty"`     // 医院列表限制数量
}

// HomeInitResponse 首页初始化响应
type HomeInitResponse struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

// HomeInitData 首页数据
type HomeInitData struct {
	Banners     []interface{} `json:"banners"`     // 轮播图
	Navigations []interface{} `json:"navigations"` // 导航
	Services    []interface{} `json:"services"`    // 服务项
	Hospitals   []interface{} `json:"hospitals"`   // 医院列表
}

// HomeInitHandler 首页初始化接口
func HomeInitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		http.Error(w, "只支持GET和POST请求", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求参数
	var req HomeInitRequest

	if r.Method == http.MethodPost {
		// POST请求，从请求体解析
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "请求参数解析失败", http.StatusBadRequest)
			return
		}
	} else {
		// GET请求，从URL参数解析
		longitudeStr := r.URL.Query().Get("longitude")
		latitudeStr := r.URL.Query().Get("latitude")
		limitStr := r.URL.Query().Get("limit")

		if longitudeStr != "" {
			if longitude, err := strconv.ParseFloat(longitudeStr, 64); err == nil {
				req.Longitude = longitude
			}
		}

		if latitudeStr != "" {
			if latitude, err := strconv.ParseFloat(latitudeStr, 64); err == nil {
				req.Latitude = latitude
			}
		}

		if limitStr != "" {
			if limit, err := strconv.Atoi(limitStr); err == nil {
				req.Limit = limit
			}
		}
	}

	// 设置默认值
	if req.Limit <= 0 {
		req.Limit = 10 // 默认返回10家医院
	}

	// 获取首页数据
	data, err := getHomeInitData(&req)
	if err != nil {
		response := &HomeInitResponse{
			Code:     -1,
			ErrorMsg: "获取首页数据失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 返回成功响应
	response := &HomeInitResponse{
		Code: 0,
		Data: data,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getHomeInitData 获取首页初始化数据
func getHomeInitData(req *HomeInitRequest) (*HomeInitData, error) {
	data := &HomeInitData{}

	// 获取轮播图
	banners, err := dao.HomeImp.GetBanners()
	if err != nil {
		return nil, err
	}
	data.Banners = convertBannersToInterface(banners)

	// 获取导航
	navigations, err := dao.HomeImp.GetNavigations()
	if err != nil {
		return nil, err
	}
	data.Navigations = convertNavigationsToInterface(navigations)

	// 获取服务项
	services, err := dao.HomeImp.GetServices()
	if err != nil {
		return nil, err
	}
	data.Services = convertServicesToInterface(services)

	// 获取医院列表
	var hospitals []*model.HospitalModel
	if req.Longitude != 0 && req.Latitude != 0 {
		// 如果有位置信息，按距离排序
		hospitals, err = dao.HomeImp.GetHospitalsByLocation(req.Longitude, req.Latitude, req.Limit)
	} else {
		// 否则按默认排序
		hospitals, err = dao.HomeImp.GetHospitals(req.Limit)
	}
	if err != nil {
		return nil, err
	}
	data.Hospitals = convertHospitalsToInterface(hospitals)

	return data, nil
}

// convertBannersToInterface 转换轮播图数据
func convertBannersToInterface(banners []*model.BannerModel) []interface{} {
	result := make([]interface{}, len(banners))
	for i, banner := range banners {
		result[i] = map[string]interface{}{
			"id":       banner.Id,
			"title":    banner.Title,
			"imageUrl": banner.ImageUrl,
			"linkUrl":  banner.LinkUrl,
			"sort":     banner.Sort,
		}
	}
	return result
}

// convertNavigationsToInterface 转换导航数据
func convertNavigationsToInterface(navigations []*model.NavigationModel) []interface{} {
	result := make([]interface{}, len(navigations))
	for i, nav := range navigations {
		result[i] = map[string]interface{}{
			"id":      nav.Id,
			"name":    nav.Name,
			"icon":    nav.Icon,
			"linkUrl": nav.LinkUrl,
			"sort":    nav.Sort,
		}
	}
	return result
}

// convertServicesToInterface 转换服务项数据
func convertServicesToInterface(services []*model.ServiceModel) []interface{} {
	result := make([]interface{}, len(services))
	for i, service := range services {
		result[i] = map[string]interface{}{
			"id":          service.Id,
			"name":        service.Name,
			"description": service.Description,
			"icon":        service.Icon,
			"imageUrl":    service.ImageUrl,
			"linkUrl":     service.LinkUrl,
			"sort":        service.Sort,
		}
	}
	return result
}

// convertHospitalsToInterface 转换医院数据
func convertHospitalsToInterface(hospitals []*model.HospitalModel) []interface{} {
	result := make([]interface{}, len(hospitals))
	for i, hospital := range hospitals {
		result[i] = map[string]interface{}{
			"id":          hospital.Id,
			"name":        hospital.Name,
			"logo":        hospital.Logo,
			"address":     hospital.Address,
			"phone":       hospital.Phone,
			"description": hospital.Description,
			"level":       hospital.Level,
			"type":        hospital.Type,
			"longitude":   hospital.Longitude,
			"latitude":    hospital.Latitude,
			"sort":        hospital.Sort,
		}
	}
	return result
}
