package dao

import (
	"fmt"
	"math"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

// GetBanners 获取轮播图列表
func (imp *HomeInterfaceImp) GetBanners() ([]*model.BannerModel, error) {
	var banners []*model.BannerModel
	cli := db.Get()
	err := cli.Table("Banners").
		Where("status = ?", 1).
		Order("sort ASC, id DESC").
		Find(&banners).Error
	return banners, err
}

// GetNavigations 获取导航列表
func (imp *HomeInterfaceImp) GetNavigations() ([]*model.NavigationModel, error) {
	var navigations []*model.NavigationModel
	cli := db.Get()
	err := cli.Table("Navigations").
		Where("status = ?", 1).
		Order("sort ASC, id DESC").
		Find(&navigations).Error
	return navigations, err
}

// GetServices 获取服务项列表
func (imp *HomeInterfaceImp) GetServices() ([]*model.ServiceModel, error) {
	var services []*model.ServiceModel
	cli := db.Get()
	err := cli.Table("Services").
		Where("status = ?", 1).
		Order("sort ASC, id DESC").
		Find(&services).Error
	return services, err
}

// GetHospitals 获取医院列表
func (imp *HomeInterfaceImp) GetHospitals(limit int) ([]*model.HospitalModel, error) {
	var hospitals []*model.HospitalModel
	cli := db.Get()
	err := cli.Table("Hospitals").
		Where("status = ?", 1).
		Order("sort ASC, id DESC").
		Limit(limit).
		Find(&hospitals).Error
	return hospitals, err
}

// GetHospitalsByLocation 根据位置获取医院列表（按距离排序）
func (imp *HomeInterfaceImp) GetHospitalsByLocation(longitude, latitude float64, limit int) ([]*model.HospitalModel, error) {
	var hospitals []*model.HospitalModel
	cli := db.Get()

	// 使用距离计算排序（简化版本，实际项目中可能需要更复杂的距离计算）
	// 这里使用经纬度的简单差值作为排序依据
	err := cli.Table("Hospitals").
		Where("status = ?", 1).
		Where("longitude IS NOT NULL AND latitude IS NOT NULL").
		Order("ABS(longitude - " + fmt.Sprintf("%f", longitude) + ") + ABS(latitude - " + fmt.Sprintf("%f", latitude) + ") ASC").
		Limit(limit).
		Find(&hospitals).Error

	return hospitals, err
}

// calculateDistance 计算两点之间的距离（简化版本）
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // 地球半径（公里）

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
