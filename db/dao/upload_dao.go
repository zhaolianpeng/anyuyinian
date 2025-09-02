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

	// 记录SQL操作日志
	logger := NewSQLLogger("插入", fileTableName, map[string]interface{}{
		"fileName": file.FileName,
		"userId":   file.UserId,
		"category": file.Category,
	})

	err := cli.Table(fileTableName).Create(file).Error
	logger.LogInsert(file, err)

	return err
}

// GetFileById 根据ID获取文件
func (imp *UploadInterfaceImp) GetFileById(id int32) (*model.FileModel, error) {
	var file = new(model.FileModel)
	cli := db.Get()

	// 记录SQL操作日志
	logger := NewSQLLogger("查询", fileTableName, map[string]interface{}{
		"id": id,
	})

	err := cli.Table(fileTableName).Where("id = ? AND status = ?", id, 1).First(file).Error
	logger.LogQuery(file, err)

	return file, err
}

// GetFilesByUserId 根据用户ID获取文件列表
func (imp *UploadInterfaceImp) GetFilesByUserId(userId string, limit int) ([]*model.FileModel, error) {
	var files []*model.FileModel
	cli := db.Get()

	// 记录SQL操作日志
	logger := NewSQLLogger("查询", fileTableName, map[string]interface{}{
		"userId": userId,
		"limit":  limit,
	})

	err := cli.Table(fileTableName).
		Where("userId = ? AND status = ?", userId, 1).
		Order("createdAt DESC").
		Limit(limit).
		Find(&files).Error
	logger.LogQuery(files, err)

	return files, err
}

// GetFilesByCategory 根据分类获取文件列表
func (imp *UploadInterfaceImp) GetFilesByCategory(category string, limit int) ([]*model.FileModel, error) {
	var files []*model.FileModel
	cli := db.Get()

	// 记录SQL操作日志
	logger := NewSQLLogger("查询", fileTableName, map[string]interface{}{
		"category": category,
		"limit":    limit,
	})

	err := cli.Table(fileTableName).
		Where("category = ? AND status = ?", category, 1).
		Order("createdAt DESC").
		Limit(limit).
		Find(&files).Error
	logger.LogQuery(files, err)

	return files, err
}

// UpdateFile 更新文件信息
func (imp *UploadInterfaceImp) UpdateFile(file *model.FileModel) error {
	cli := db.Get()
	file.UpdatedAt = time.Now()

	// 记录SQL操作日志
	logger := NewSQLLogger("更新", fileTableName, map[string]interface{}{
		"id":       file.Id,
		"fileName": file.FileName,
	})

	err := cli.Table(fileTableName).Where("id = ?", file.Id).Updates(file).Error
	logger.LogUpdate(file, err)

	return err
}

// DeleteFile 删除文件（软删除）
func (imp *UploadInterfaceImp) DeleteFile(id int32) error {
	cli := db.Get()

	// 记录SQL操作日志
	logger := NewSQLLogger("删除", fileTableName, map[string]interface{}{
		"id": id,
	})

	err := cli.Table(fileTableName).Where("id = ?", id).Update("status", 0).Error
	logger.LogDelete(map[string]interface{}{"id": id, "status": 0}, err)

	return err
}
