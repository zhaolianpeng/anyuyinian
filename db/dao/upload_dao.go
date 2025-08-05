package dao

import (
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

const fileTableName = "Files"

// CreateFile 创建文件记录
func (imp *UploadInterfaceImp) CreateFile(file *model.FileModel) error {
	cli := db.Get()
	file.CreatedAt = time.Now()
	file.UpdatedAt = time.Now()
	return cli.Table(fileTableName).Create(file).Error
}

// GetFileById 根据ID获取文件
func (imp *UploadInterfaceImp) GetFileById(id int32) (*model.FileModel, error) {
	var file = new(model.FileModel)
	cli := db.Get()
	err := cli.Table(fileTableName).Where("id = ? AND status = ?", id, 1).First(file).Error
	return file, err
}

// GetFilesByUserId 根据用户ID获取文件列表
func (imp *UploadInterfaceImp) GetFilesByUserId(userId string, limit int) ([]*model.FileModel, error) {
	var files []*model.FileModel
	cli := db.Get()
	err := cli.Table(fileTableName).
		Where("userId = ? AND status = ?", userId, 1).
		Order("createdAt DESC").
		Limit(limit).
		Find(&files).Error
	return files, err
}

// GetFilesByCategory 根据分类获取文件列表
func (imp *UploadInterfaceImp) GetFilesByCategory(category string, limit int) ([]*model.FileModel, error) {
	var files []*model.FileModel
	cli := db.Get()
	err := cli.Table(fileTableName).
		Where("category = ? AND status = ?", category, 1).
		Order("createdAt DESC").
		Limit(limit).
		Find(&files).Error
	return files, err
}

// UpdateFile 更新文件信息
func (imp *UploadInterfaceImp) UpdateFile(file *model.FileModel) error {
	cli := db.Get()
	file.UpdatedAt = time.Now()
	return cli.Table(fileTableName).Where("id = ?", file.Id).Updates(file).Error
}

// DeleteFile 删除文件（软删除）
func (imp *UploadInterfaceImp) DeleteFile(id int32) error {
	cli := db.Get()
	return cli.Table(fileTableName).Where("id = ?", id).Update("status", 0).Error
}
