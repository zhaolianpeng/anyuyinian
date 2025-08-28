package service

import (
	"encoding/json"
	"fmt"
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
	Banners           []interface{} `json:"banners"`           // 轮播图
	Navigations       []interface{} `json:"navigations"`       // 导航
	Services          []interface{} `json:"services"`          // 服务项
	Hospitals         []interface{} `json:"hospitals"`         // 医院列表
	CaregiverServices []interface{} `json:"caregiverServices"` // 护工服务列表
}

// HomeInitHandler 首页初始化接口
func HomeInitHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理首页初始化请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		LogError("请求方法不支持", fmt.Errorf("期望GET或POST方法，实际为%s", r.Method))
		http.Error(w, "只支持GET和POST请求", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求参数
	var req HomeInitRequest

	LogStep("开始解析请求参数", map[string]interface{}{
		"method": r.Method,
	})

	if r.Method == http.MethodPost {
		// POST请求，从请求体解析
		LogStep("解析POST请求体参数", nil)
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			LogError("POST请求参数解析失败", err)
			http.Error(w, "请求参数解析失败", http.StatusBadRequest)
			return
		}
		LogStep("POST请求参数解析成功", map[string]interface{}{
			"longitude": req.Longitude,
			"latitude":  req.Latitude,
			"limit":     req.Limit,
		})
	} else {
		// GET请求，从URL参数解析
		LogStep("解析GET请求URL参数", nil)
		longitudeStr := r.URL.Query().Get("longitude")
		latitudeStr := r.URL.Query().Get("latitude")
		limitStr := r.URL.Query().Get("limit")

		LogStep("获取URL参数", map[string]interface{}{
			"longitude": longitudeStr,
			"latitude":  latitudeStr,
			"limit":     limitStr,
		})

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

		LogStep("GET请求参数解析成功", map[string]interface{}{
			"longitude": req.Longitude,
			"latitude":  req.Latitude,
			"limit":     req.Limit,
		})
	}

	// 设置默认值
	if req.Limit <= 0 {
		req.Limit = 10 // 默认返回10家医院
		LogStep("设置默认医院数量限制", map[string]interface{}{
			"limit": req.Limit,
		})
	}

	LogStep("开始获取首页数据", map[string]interface{}{
		"longitude": req.Longitude,
		"latitude":  req.Latitude,
		"limit":     req.Limit,
	})

	// 获取首页数据
	data, err := getHomeInitData(&req)
	if err != nil {
		LogError("获取首页数据失败", err)
		response := &HomeInitResponse{
			Code:     -1,
			ErrorMsg: "获取首页数据失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	LogStep("首页数据获取成功", map[string]interface{}{
		"bannerCount":     len(data.Banners),
		"navigationCount": len(data.Navigations),
		"serviceCount":    len(data.Services),
		"hospitalCount":   len(data.Hospitals),
	})

	// 返回成功响应
	response := &HomeInitResponse{
		Code: 0,
		Data: data,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	LogInfo("首页初始化成功", map[string]interface{}{
		"bannerCount":     len(data.Banners),
		"navigationCount": len(data.Navigations),
		"serviceCount":    len(data.Services),
		"hospitalCount":   len(data.Hospitals),
	})
}

// getHomeInitData 获取首页初始化数据
func getHomeInitData(req *HomeInitRequest) (*HomeInitData, error) {
	LogStep("开始获取首页初始化数据", nil)

	data := &HomeInitData{}

	// 获取轮播图
	LogStep("开始查询轮播图数据", nil)
	banners, err := dao.HomeImp.GetBanners()
	if err != nil {
		LogError("数据库查询轮播图失败", err)
		return nil, err
	}
	LogStep("轮播图数据查询成功", map[string]interface{}{
		"bannerCount": len(banners),
	})
	data.Banners = convertBannersToInterface(banners)

	// 获取导航
	LogStep("开始查询导航数据", nil)
	navigations, err := dao.HomeImp.GetNavigations()
	if err != nil {
		LogError("数据库查询导航失败", err)
		return nil, err
	}
	LogStep("导航数据查询成功", map[string]interface{}{
		"navigationCount": len(navigations),
	})
	data.Navigations = convertNavigationsToInterface(navigations)

	// 获取服务项
	LogStep("开始查询服务项数据", nil)
	services, err := dao.HomeImp.GetServices()
	if err != nil {
		LogError("数据库查询服务项失败", err)
		return nil, err
	}
	LogStep("服务项数据查询成功", map[string]interface{}{
		"serviceCount": len(services),
	})
	data.Services = convertServicesToInterface(services)

	// 获取医院列表
	LogStep("开始查询医院列表数据", map[string]interface{}{
		"longitude": req.Longitude,
		"latitude":  req.Latitude,
		"limit":     req.Limit,
	})

	var hospitals []*model.HospitalModel
	if req.Longitude != 0 && req.Latitude != 0 {
		// 如果有位置信息，按距离排序
		LogStep("按位置信息查询医院", map[string]interface{}{
			"longitude": req.Longitude,
			"latitude":  req.Latitude,
		})
		hospitals, err = dao.HomeImp.GetHospitalsByLocation(req.Longitude, req.Latitude, req.Limit)
	} else {
		// 否则按默认排序
		LogStep("按默认排序查询医院", nil)
		hospitals, err = dao.HomeImp.GetHospitals(req.Limit)
	}
	if err != nil {
		LogError("数据库查询医院列表失败", err)
		return nil, err
	}
	LogStep("医院列表数据查询成功", map[string]interface{}{
		"hospitalCount": len(hospitals),
	})
	data.Hospitals = convertHospitalsToInterface(hospitals)

	// 获取护工服务列表（从Services表获取，按分类组织）
	LogStep("开始查询护工服务数据", nil)

	// 获取所有分类的护工服务
	caregiverServices := make([]interface{}, 0)

	// 居家照护类服务
	homeCareServices, err := dao.HomeImp.GetServicesByCategory("居家照护")
	if err == nil && len(homeCareServices) > 0 {
		caregiverServices = append(caregiverServices, convertServicesToInterface(homeCareServices)...)
	}

	// 医院陪诊类服务
	hospitalEscortServices, err := dao.HomeImp.GetServicesByCategory("医院陪诊")
	if err == nil && len(hospitalEscortServices) > 0 {
		caregiverServices = append(caregiverServices, convertServicesToInterface(hospitalEscortServices)...)
	}

	// 周期护理类服务
	periodicCareServices, err := dao.HomeImp.GetServicesByCategory("周期护理")
	if err == nil && len(periodicCareServices) > 0 {
		caregiverServices = append(caregiverServices, convertServicesToInterface(periodicCareServices)...)
	}

	// 家政服务类服务
	housekeepingServices, err := dao.HomeImp.GetServicesByCategory("家政服务")
	if err == nil && len(housekeepingServices) > 0 {
		caregiverServices = append(caregiverServices, convertServicesToInterface(housekeepingServices)...)
	}

	// 如果没有找到任何护工服务，使用默认数据
	if len(caregiverServices) == 0 {
		LogError("未找到护工服务数据，使用默认数据", nil)
		data.CaregiverServices = getDefaultCaregiverServices()
	} else {
		LogStep("护工服务数据查询成功", map[string]interface{}{
			"caregiverServiceCount": len(caregiverServices),
		})
		data.CaregiverServices = caregiverServices
	}

	LogStep("首页数据获取完成", map[string]interface{}{
		"bannerCount":           len(data.Banners),
		"navigationCount":       len(data.Navigations),
		"serviceCount":          len(data.Services),
		"hospitalCount":         len(data.Hospitals),
		"caregiverServiceCount": len(data.CaregiverServices),
	})

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
			"id":            service.Id,
			"serviceitemid": service.ServiceItemId, // 使用数据库中的serviceitemid字段
			"name":          service.Name,
			"description":   service.Description,
			"icon":          service.Icon,
			"imageUrl":      service.ImageUrl,
			"linkUrl":       service.LinkUrl,
			"sort":          service.Sort,
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

// getDefaultCaregiverServices 获取默认护工服务数据
func getDefaultCaregiverServices() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"id":            1,
			"serviceitemid": 1,
			"name":          "慢病照护",
			"description":   "生活支援,守护健康",
			"imageUrl":      "/images/service/chronic-care.jpg",
			"price":         4880.0,
			"category":      "居家照护",
			"sort":          1,
		},
		map[string]interface{}{
			"id":            2,
			"serviceitemid": 2,
			"name":          "居家术后照护",
			"description":   "省心省力、全天照护",
			"imageUrl":      "/images/service/post-surgery.jpg",
			"price":         5580.0,
			"category":      "居家照护",
			"sort":          2,
		},
		map[string]interface{}{
			"id":            3,
			"serviceitemid": 3,
			"name":          "康复照护",
			"description":   "偏瘫、肢体康复训练",
			"imageUrl":      "/images/service/rehabilitation.jpg",
			"price":         6280.0,
			"category":      "周期护理",
			"sort":          3,
		},
		map[string]interface{}{
			"id":            4,
			"serviceitemid": 4,
			"name":          "认知症照护",
			"description":   "认知症(阿尔兹海默病)习惯培养、守护健康",
			"imageUrl":      "/images/service/dementia-care.jpg",
			"price":         4980.0,
			"category":      "居家照护",
			"sort":          4,
		},
		map[string]interface{}{
			"id":            5,
			"serviceitemid": 5,
			"name":          "肌无力照护",
			"description":   "行动支持、安全防护",
			"imageUrl":      "/images/service/muscle-weakness.jpg",
			"price":         5380.0,
			"category":      "居家照护",
			"sort":          5,
		},
	}
}
