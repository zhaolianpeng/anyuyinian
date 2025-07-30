package dao

import (
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

const configTableName = "Configs"

// GetConfigByKey 根据key获取配置
func (imp *ConfigInterfaceImp) GetConfigByKey(key string) (*model.ConfigModel, error) {
	var config = new(model.ConfigModel)
	cli := db.Get()
	err := cli.Table(configTableName).Where("key = ? AND status = ?", key, 1).First(config).Error
	return config, err
}

// GetAllConfigs 获取所有配置
func (imp *ConfigInterfaceImp) GetAllConfigs() ([]*model.ConfigModel, error) {
	var configs []*model.ConfigModel
	cli := db.Get()
	err := cli.Table(configTableName).Where("status = ?", 1).Find(&configs).Error
	return configs, err
}

// CreateConfig 创建配置
func (imp *ConfigInterfaceImp) CreateConfig(config *model.ConfigModel) error {
	cli := db.Get()
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()
	return cli.Table(configTableName).Create(config).Error
}

// UpdateConfig 更新配置
func (imp *ConfigInterfaceImp) UpdateConfig(config *model.ConfigModel) error {
	cli := db.Get()
	config.UpdatedAt = time.Now()
	return cli.Table(configTableName).Where("id = ?", config.Id).Updates(config).Error
}

// DeleteConfig 删除配置（软删除）
func (imp *ConfigInterfaceImp) DeleteConfig(id int32) error {
	cli := db.Get()
	return cli.Table(configTableName).Where("id = ?", id).Update("status", 0).Error
}
