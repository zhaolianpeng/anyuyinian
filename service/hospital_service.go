package service

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"

	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

// HospitalResponse 医院响应
type HospitalResponse struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

// HospitalListResponse 医院列表响应
type HospitalListResponse struct {
	List     []*model.HospitalModel `json:"list"`
	Total    int64                  `json:"total"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"pageSize"`
	HasMore  bool                   `json:"hasMore"`
}

// HospitalDetailResponse 医院详情响应
type HospitalDetailResponse struct {
	Hospital   *model.HospitalModel `json:"hospital"`
	Navigation NavigationInfo       `json:"navigation"`
}

// NavigationInfo 导航信息
type NavigationInfo struct {
	Distance    float64 `json:"distance"`    // 距离（公里）
	Duration    int     `json:"duration"`    // 预计时间（分钟）
	RouteType   string  `json:"routeType"`   // 路线类型：driving, walking, transit
	RoutePoints []Point `json:"routePoints"` // 路线点
}

// Point 坐标点
type Point struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// HospitalListHandler 获取可选医院列表接口
func HospitalListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	// 解析查询参数
	longitudeStr := r.URL.Query().Get("longitude")
	latitudeStr := r.URL.Query().Get("latitude")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	// 设置默认值
	page := 1
	pageSize := 10
	var longitude, latitude float64

	if longitudeStr != "" {
		if lng, err := strconv.ParseFloat(longitudeStr, 64); err == nil {
			longitude = lng
		}
	}

	if latitudeStr != "" {
		if lat, err := strconv.ParseFloat(latitudeStr, 64); err == nil {
			latitude = lat
		}
	}

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 50 {
			pageSize = ps
		}
	}

	// 获取医院列表
	var hospitals []*model.HospitalModel
	var err error

	if longitude != 0 && latitude != 0 {
		hospitals, err = dao.HomeImp.GetHospitalsByLocation(longitude, latitude, pageSize)
	} else {
		hospitals, err = dao.HomeImp.GetHospitals(pageSize)
	}

	if err != nil {
		response := &HospitalResponse{
			Code:     -1,
			ErrorMsg: "获取医院列表失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 计算总数（简化处理）
	total := int64(len(hospitals))
	hasMore := int64(page*pageSize) < total

	response := &HospitalResponse{
		Code: 0,
		Data: &HospitalListResponse{
			List:     hospitals,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
			HasMore:  hasMore,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HospitalDetailHandler 医院详情接口
func HospitalDetailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	// 从URL路径中获取医院ID
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "缺少医院ID参数", http.StatusBadRequest)
		return
	}

	hospitalIdStr := pathParts[3]
	hospitalId, err := strconv.Atoi(hospitalIdStr)
	if err != nil {
		http.Error(w, "无效的医院ID", http.StatusBadRequest)
		return
	}

	// 获取医院详情（简化处理，从医院列表中查找）
	hospitals, err := dao.HomeImp.GetHospitals(100) // 获取所有医院
	if err != nil {
		response := &HospitalResponse{
			Code:     -1,
			ErrorMsg: "获取医院详情失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 查找指定医院
	var hospital *model.HospitalModel
	for _, h := range hospitals {
		if h.Id == int32(hospitalId) {
			hospital = h
			break
		}
	}

	if hospital == nil {
		response := &HospitalResponse{
			Code:     -1,
			ErrorMsg: "医院不存在",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 获取用户当前位置（从查询参数）
	userLongitudeStr := r.URL.Query().Get("userLongitude")
	userLatitudeStr := r.URL.Query().Get("userLatitude")

	var navigation NavigationInfo
	if userLongitudeStr != "" && userLatitudeStr != "" {
		if userLng, err := strconv.ParseFloat(userLongitudeStr, 64); err == nil {
			if userLat, err := strconv.ParseFloat(userLatitudeStr, 64); err == nil {
				// 计算距离（简化计算，实际项目中可以使用更精确的算法）
				distance := calculateDistance(userLat, userLng, hospital.Latitude, hospital.Longitude)
				duration := int(distance * 2) // 假设平均速度30km/h

				navigation = NavigationInfo{
					Distance:  distance,
					Duration:  duration,
					RouteType: "driving",
					RoutePoints: []Point{
						{Latitude: userLat, Longitude: userLng},
						{Latitude: hospital.Latitude, Longitude: hospital.Longitude},
					},
				}
			}
		}
	}

	response := &HospitalResponse{
		Code: 0,
		Data: &HospitalDetailResponse{
			Hospital:   hospital,
			Navigation: navigation,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 计算两点间距离（公里）
func calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371 // 地球半径（公里）

	// 转换为弧度
	lat1Rad := lat1 * 3.141592653589793 / 180
	lat2Rad := lat2 * 3.141592653589793 / 180
	deltaLat := (lat2 - lat1) * 3.141592653589793 / 180
	deltaLng := (lng2 - lng1) * 3.141592653589793 / 180

	// Haversine公式
	a := deltaLat/2*deltaLat/2 + lat1Rad*lat2Rad*deltaLng/2*deltaLng/2
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
