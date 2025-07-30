package dao

import (
	"wxcloudrun-golang/db/model"
)

// ServiceInterface 服务数据接口
type ServiceInterface interface {
	GetServiceById(id int32) (*model.ServiceItemModel, error)
	GetServicesByCategory(category string, page, pageSize int) ([]*model.ServiceItemModel, int64, error)
	GetAllServices(page, pageSize int) ([]*model.ServiceItemModel, int64, error)
	GetServiceCategories() ([]string, error)
	CreateService(service *model.ServiceItemModel) error
	UpdateService(service *model.ServiceItemModel) error
	DeleteService(id int32) error
}

// ServiceInterfaceImp 服务数据实现
type ServiceInterfaceImp struct{}

// ServiceImp 服务实现实例
var ServiceImp ServiceInterface = &ServiceInterfaceImp{}
