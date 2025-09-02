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

	// 记录SQL操作日志
	logger := NewSQLLogger("查询", serviceTableName, map[string]interface{}{
		"id":     id,
		"status": 1,
	})

	err := cli.Table(serviceTableName).Where("id = ? AND status = ?", id, 1).First(service).Error
	logger.LogQuery(service, err)

	return service, err
}

// GetServicesByCategory 根据分类获取服务列表（分页）
func (imp *ServiceInterfaceImp) GetServicesByCategory(category string, page, pageSize int) ([]*model.ServiceItemModel, int64, error) {
	var services []*model.ServiceItemModel
	var total int64
	cli := db.Get()

	// 记录SQL操作日志
	logger := NewSQLLogger("查询", serviceTableName, map[string]interface{}{
		"category": category,
		"page":     page,
		"pageSize": pageSize,
		"status":   1,
	})

	// 获取总数
	err := cli.Table(serviceTableName).Where("category = ? AND status = ?", category, 1).Count(&total).Error
	if err != nil {
		logger.LogCount(0, err)
		return nil, 0, err
	}
	logger.LogCount(total, nil)

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = cli.Table(serviceTableName).
		Where("category = ? AND status = ?", category, 1).
		Order("sort ASC, createdAt DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&services).Error

	logger.LogQuery(services, err)
	return services, total, err
}

// GetAllServices 获取所有服务列表（分页）
func (imp *ServiceInterfaceImp) GetAllServices(page, pageSize int) ([]*model.ServiceItemModel, int64, error) {
	var services []*model.ServiceItemModel
	var total int64
	cli := db.Get()

	// 记录SQL操作日志
	logger := NewSQLLogger("查询", serviceTableName, map[string]interface{}{
		"page":     page,
		"pageSize": pageSize,
		"status":   1,
	})

	// 获取总数
	err := cli.Table(serviceTableName).Where("status = ?", 1).Count(&total).Error
	if err != nil {
		logger.LogCount(0, err)
		return nil, 0, err
	}
	logger.LogCount(total, nil)

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = cli.Table(serviceTableName).
		Where("status = ?", 1).
		Order("sort ASC, createdAt DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&services).Error

	logger.LogQuery(services, err)
	return services, total, err
}

// CategoryCount 分类数量结构
type CategoryCount struct {
	Category string `json:"category"`
	Name     string `json:"name"`
	Count    int64  `json:"count"`
}

// GetServiceCategories 获取服务分类列表
func (imp *ServiceInterfaceImp) GetServiceCategories() ([]string, error) {
	var categories []string
	cli := db.Get()

	// 记录SQL操作日志
	logger := NewSQLLogger("查询", serviceTableName, map[string]interface{}{
		"status": 1,
		"field":  "category",
	})

	err := cli.Table(serviceTableName).
		Where("status = ?", 1).
		Distinct("category").
		Pluck("category", &categories).Error

	logger.LogQuery(categories, err)
	return categories, err
}

// GetServiceCategoriesWithCount 获取服务分类列表及其数量
func (imp *ServiceInterfaceImp) GetServiceCategoriesWithCount() ([]CategoryCount, error) {
	var categories []CategoryCount
	cli := db.Get()

	// 记录SQL操作日志
	logger := NewSQLLogger("查询", serviceTableName, map[string]interface{}{
		"status": 1,
		"group":  "category",
	})

	// 查询每个分类的服务数量
	err := cli.Table(serviceTableName).
		Where("status = ?", 1).
		Select("category, COUNT(*) as count").
		Group("category").
		Order("count DESC").
		Scan(&categories).Error

	if err != nil {
		logger.LogQuery(nil, err)
		return nil, err
	}

	logger.LogQuery(categories, nil)

	// 为每个分类设置显示名称
	for i := range categories {
		switch categories[i].Category {
		case "居家照护":
			categories[i].Name = "居家照护"
		case "医院陪诊":
			categories[i].Name = "医院陪诊"
		case "周期护理":
			categories[i].Name = "周期护理"
		case "家政服务":
			categories[i].Name = "家政服务"
		case "预约咨询":
			categories[i].Name = "预约咨询"
		case "智慧养老":
			categories[i].Name = "智慧养老"
		default:
			categories[i].Name = categories[i].Category
		}
	}

	return categories, nil
}

// CreateService 创建服务
func (imp *ServiceInterfaceImp) CreateService(service *model.ServiceItemModel) error {
	cli := db.Get()
	service.CreatedAt = time.Now()
	service.UpdatedAt = time.Now()

	// 记录SQL操作日志
	logger := NewSQLLogger("插入", serviceTableName, map[string]interface{}{
		"name":     service.Name,
		"category": service.Category,
		"price":    service.Price,
	})

	err := cli.Table(serviceTableName).Create(service).Error
	logger.LogInsert(service, err)

	return err
}

// UpdateService 更新服务
func (imp *ServiceInterfaceImp) UpdateService(service *model.ServiceItemModel) error {
	cli := db.Get()
	service.UpdatedAt = time.Now()

	// 记录SQL操作日志
	logger := NewSQLLogger("更新", serviceTableName, map[string]interface{}{
		"id":       service.Id,
		"name":     service.Name,
		"category": service.Category,
	})

	result := cli.Table(serviceTableName).Where("id = ?", service.Id).Updates(service)
	logger.LogUpdate(result.RowsAffected, result.Error)

	return result.Error
}

// DeleteService 删除服务（软删除）
func (imp *ServiceInterfaceImp) DeleteService(id int32) error {
	cli := db.Get()

	// 记录SQL操作日志
	logger := NewSQLLogger("删除", serviceTableName, map[string]interface{}{
		"id":     id,
		"status": 0,
	})

	result := cli.Table(serviceTableName).Where("id = ?", id).Update("status", 0)
	logger.LogDelete(result.RowsAffected, result.Error)

	return result.Error
}
