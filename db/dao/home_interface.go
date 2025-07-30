package dao

import (
	"wxcloudrun-golang/db/model"
)

// HomeInterface 首页数据接口
type HomeInterface interface {
	// 轮播图相关
	GetBanners() ([]*model.BannerModel, error)

	// 导航相关
	GetNavigations() ([]*model.NavigationModel, error)

	// 服务项相关
	GetServices() ([]*model.ServiceModel, error)

	// 医院相关
	GetHospitals(limit int) ([]*model.HospitalModel, error)
	GetHospitalsByLocation(longitude, latitude float64, limit int) ([]*model.HospitalModel, error)
}

// HomeInterfaceImp 首页数据实现
type HomeInterfaceImp struct{}

// HomeImp 首页实现实例
var HomeImp HomeInterface = &HomeInterfaceImp{}
