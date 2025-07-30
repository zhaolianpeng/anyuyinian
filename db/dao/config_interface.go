package dao

import (
	"wxcloudrun-golang/db/model"
)

// ConfigInterface 配置数据接口
type ConfigInterface interface {
	GetConfigByKey(key string) (*model.ConfigModel, error)
	GetAllConfigs() ([]*model.ConfigModel, error)
	CreateConfig(config *model.ConfigModel) error
	UpdateConfig(config *model.ConfigModel) error
	DeleteConfig(id int32) error
}

// ConfigInterfaceImp 配置数据实现
type ConfigInterfaceImp struct{}

// ConfigImp 配置实现实例
var ConfigImp ConfigInterface = &ConfigInterfaceImp{}
