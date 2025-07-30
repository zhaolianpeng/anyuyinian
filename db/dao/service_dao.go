package dao

import (
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

const serviceTableName = "ServiceItems"

// GetServiceById 根据ID获取服务
func (imp *ServiceInterfaceImp) GetServiceById(id int32) (*model.ServiceItemModel, error) {
	var service = new(model.ServiceItemModel)
	cli := db.Get()
	err := cli.Table(serviceTableName).Where("id = ? AND status = ?", id, 1).First(service).Error
	return service, err
}

// GetServicesByCategory 根据分类获取服务列表（分页）
func (imp *ServiceInterfaceImp) GetServicesByCategory(category string, page, pageSize int) ([]*model.ServiceItemModel, int64, error) {
	var services []*model.ServiceItemModel
	var total int64
	cli := db.Get()

	// 获取总数
	err := cli.Table(serviceTableName).Where("category = ? AND status = ?", category, 1).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = cli.Table(serviceTableName).
		Where("category = ? AND status = ?", category, 1).
		Order("sort ASC, createdAt DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&services).Error

	return services, total, err
}

// GetAllServices 获取所有服务列表（分页）
func (imp *ServiceInterfaceImp) GetAllServices(page, pageSize int) ([]*model.ServiceItemModel, int64, error) {
	var services []*model.ServiceItemModel
	var total int64
	cli := db.Get()

	// 获取总数
	err := cli.Table(serviceTableName).Where("status = ?", 1).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = cli.Table(serviceTableName).
		Where("status = ?", 1).
		Order("sort ASC, createdAt DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&services).Error

	return services, total, err
}

// GetServiceCategories 获取服务分类列表
func (imp *ServiceInterfaceImp) GetServiceCategories() ([]string, error) {
	var categories []string
	cli := db.Get()
	err := cli.Table(serviceTableName).
		Where("status = ?", 1).
		Distinct("category").
		Pluck("category", &categories).Error
	return categories, err
}

// CreateService 创建服务
func (imp *ServiceInterfaceImp) CreateService(service *model.ServiceItemModel) error {
	cli := db.Get()
	service.CreatedAt = time.Now()
	service.UpdatedAt = time.Now()
	return cli.Table(serviceTableName).Create(service).Error
}

// UpdateService 更新服务
func (imp *ServiceInterfaceImp) UpdateService(service *model.ServiceItemModel) error {
	cli := db.Get()
	service.UpdatedAt = time.Now()
	return cli.Table(serviceTableName).Where("id = ?", service.Id).Updates(service).Error
}

// DeleteService 删除服务（软删除）
func (imp *ServiceInterfaceImp) DeleteService(id int32) error {
	cli := db.Get()
	return cli.Table(serviceTableName).Where("id = ?", id).Update("status", 0).Error
}
