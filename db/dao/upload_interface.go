package dao

import (
	"wxcloudrun-golang/db/model"
)

// UploadInterface 文件上传数据接口
type UploadInterface interface {
	// 文件相关
	CreateFile(file *model.FileModel) error
	GetFileById(id int32) (*model.FileModel, error)
	GetFilesByUserId(userId int32, limit int) ([]*model.FileModel, error)
	GetFilesByCategory(category string, limit int) ([]*model.FileModel, error)
	UpdateFile(file *model.FileModel) error
	DeleteFile(id int32) error
}

// UploadInterfaceImp 文件上传数据实现
type UploadInterfaceImp struct{}

// UploadImp 文件上传实现实例
var UploadImp UploadInterface = &UploadInterfaceImp{}
