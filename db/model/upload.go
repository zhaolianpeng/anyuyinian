package model

import "time"

// FileModel 文件模型
type FileModel struct {
	Id           int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	FileName     string    `gorm:"column:fileName;not null" json:"fileName"`
	OriginalName string    `gorm:"column:originalName;not null" json:"originalName"`
	FilePath     string    `gorm:"column:filePath;not null" json:"filePath"`
	FileUrl      string    `gorm:"column:fileUrl;not null" json:"fileUrl"`
	FileSize     int64     `gorm:"column:fileSize" json:"fileSize"`
	FileType     string    `gorm:"column:fileType" json:"fileType"`
	MimeType     string    `gorm:"column:mimeType" json:"mimeType"`
	Category     string    `gorm:"column:category" json:"category"` // 文件分类：image, document, video等
	Description  string    `gorm:"column:description" json:"description"`
	UserId       string    `gorm:"column:userId;type:varchar(24)" json:"userId"` // 上传用户ID
	Status       int       `gorm:"column:status;default:1" json:"status"`        // 1-正常，0-删除
	CreatedAt    time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

// TableName 指定表名
func (FileModel) TableName() string {
	return "Files"
}
