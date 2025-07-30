package service

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

// UploadResponse 上传响应
type UploadResponse struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

// FileInfo 文件信息
type FileInfo struct {
	Id           int32  `json:"id"`
	FileName     string `json:"fileName"`
	OriginalName string `json:"originalName"`
	FileUrl      string `json:"fileUrl"`
	FileSize     int64  `json:"fileSize"`
	FileType     string `json:"fileType"`
	MimeType     string `json:"mimeType"`
	Category     string `json:"category"`
	Description  string `json:"description"`
}

// UploadConfig 上传配置
type UploadConfig struct {
	MaxFileSize  int64    // 最大文件大小（字节）
	AllowedTypes []string // 允许的文件类型
	UploadPath   string   // 上传路径
	BaseURL      string   // 基础URL
}

var uploadConfig = &UploadConfig{
	MaxFileSize: 10 * 1024 * 1024, // 10MB
	AllowedTypes: []string{
		"image/jpeg", "image/png", "image/gif", "image/webp",
		"application/pdf", "application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"text/plain", "application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	},
	UploadPath: "./uploads",
	BaseURL:    "http://localhost:80/uploads",
}

// UploadHandler 文件上传接口
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	// 解析multipart表单
	err := r.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		http.Error(w, "解析表单失败", http.StatusBadRequest)
		return
	}

	// 获取上传的文件
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "获取上传文件失败", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 获取其他参数
	userIdStr := r.FormValue("userId")
	category := r.FormValue("category")
	description := r.FormValue("description")

	// 验证用户ID
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "无效的用户ID", http.StatusBadRequest)
		return
	}

	// 验证文件
	if err := validateFile(header); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 保存文件
	fileInfo, err := saveFile(file, header, int32(userId), category, description)
	if err != nil {
		http.Error(w, "保存文件失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回成功响应
	response := &UploadResponse{
		Code: 0,
		Data: fileInfo,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// validateFile 验证文件
func validateFile(header *multipart.FileHeader) error {
	// 检查文件大小
	if header.Size > uploadConfig.MaxFileSize {
		return fmt.Errorf("文件大小超过限制，最大允许 %d 字节", uploadConfig.MaxFileSize)
	}

	// 检查文件类型
	contentType := header.Header.Get("Content-Type")
	isAllowed := false
	for _, allowedType := range uploadConfig.AllowedTypes {
		if contentType == allowedType {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		return fmt.Errorf("不支持的文件类型: %s", contentType)
	}

	return nil
}

// saveFile 保存文件
func saveFile(file multipart.File, header *multipart.FileHeader, userId int32, category, description string) (*FileInfo, error) {
	// 创建上传目录
	if err := os.MkdirAll(uploadConfig.UploadPath, 0755); err != nil {
		return nil, fmt.Errorf("创建上传目录失败: %v", err)
	}

	// 生成文件名
	originalName := header.Filename
	fileName := generateFileName(originalName)
	fullPath := filepath.Join(uploadConfig.UploadPath, fileName)

	// 创建目标文件
	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %v", err)
	}
	defer dst.Close()

	// 复制文件内容
	_, err = io.Copy(dst, file)
	if err != nil {
		return nil, fmt.Errorf("保存文件失败: %v", err)
	}

	// 获取文件信息
	fileInfo, err := dst.Stat()
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %v", err)
	}

	// 确定文件类型和分类
	contentType := header.Header.Get("Content-Type")
	fileType := getFileType(contentType)
	if category == "" {
		category = getCategoryByType(fileType)
	}

	// 创建数据库记录
	dbFile := &model.FileModel{
		FileName:     fileName,
		OriginalName: originalName,
		FilePath:     fullPath,
		FileUrl:      uploadConfig.BaseURL + "/" + fileName,
		FileSize:     fileInfo.Size(),
		FileType:     fileType,
		MimeType:     contentType,
		Category:     category,
		Description:  description,
		UserId:       userId,
	}

	if err := dao.UploadImp.CreateFile(dbFile); err != nil {
		// 如果数据库保存失败，删除已上传的文件
		os.Remove(fullPath)
		return nil, fmt.Errorf("保存文件记录失败: %v", err)
	}

	// 返回文件信息
	return &FileInfo{
		Id:           dbFile.Id,
		FileName:     dbFile.FileName,
		OriginalName: dbFile.OriginalName,
		FileUrl:      dbFile.FileUrl,
		FileSize:     dbFile.FileSize,
		FileType:     dbFile.FileType,
		MimeType:     dbFile.MimeType,
		Category:     dbFile.Category,
		Description:  dbFile.Description,
	}, nil
}

// generateFileName 生成文件名
func generateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	name := strings.TrimSuffix(originalName, ext)

	// 使用时间戳和MD5生成唯一文件名
	timestamp := time.Now().UnixNano()
	hash := md5.Sum([]byte(fmt.Sprintf("%s_%d", name, timestamp)))
	hashStr := hex.EncodeToString(hash[:])

	return fmt.Sprintf("%d_%s%s", timestamp, hashStr[:8], ext)
}

// getFileType 根据MIME类型获取文件类型
func getFileType(mimeType string) string {
	switch {
	case strings.HasPrefix(mimeType, "image/"):
		return "image"
	case strings.HasPrefix(mimeType, "video/"):
		return "video"
	case strings.HasPrefix(mimeType, "audio/"):
		return "audio"
	case mimeType == "application/pdf":
		return "pdf"
	case strings.Contains(mimeType, "word"):
		return "document"
	case strings.Contains(mimeType, "excel"):
		return "spreadsheet"
	case mimeType == "text/plain":
		return "text"
	default:
		return "other"
	}
}

// getCategoryByType 根据文件类型获取分类
func getCategoryByType(fileType string) string {
	switch fileType {
	case "image":
		return "image"
	case "video":
		return "video"
	case "audio":
		return "audio"
	case "pdf":
		return "document"
	case "document":
		return "document"
	case "spreadsheet":
		return "document"
	case "text":
		return "document"
	default:
		return "other"
	}
}

// GetFileListHandler 获取文件列表接口
func GetFileListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	// 获取查询参数
	userIdStr := r.URL.Query().Get("userId")
	category := r.URL.Query().Get("category")
	limitStr := r.URL.Query().Get("limit")

	limit := 20 // 默认限制
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	var files []*model.FileModel
	var err error

	if userIdStr != "" {
		// 根据用户ID查询
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			http.Error(w, "无效的用户ID", http.StatusBadRequest)
			return
		}
		files, err = dao.UploadImp.GetFilesByUserId(int32(userId), limit)
	} else if category != "" {
		// 根据分类查询
		files, err = dao.UploadImp.GetFilesByCategory(category, limit)
	} else {
		http.Error(w, "缺少查询参数", http.StatusBadRequest)
		return
	}

	if err != nil {
		response := &UploadResponse{
			Code:     -1,
			ErrorMsg: "获取文件列表失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 转换文件信息
	fileInfos := make([]*FileInfo, len(files))
	for i, file := range files {
		fileInfos[i] = &FileInfo{
			Id:           file.Id,
			FileName:     file.FileName,
			OriginalName: file.OriginalName,
			FileUrl:      file.FileUrl,
			FileSize:     file.FileSize,
			FileType:     file.FileType,
			MimeType:     file.MimeType,
			Category:     file.Category,
			Description:  file.Description,
		}
	}

	response := &UploadResponse{
		Code: 0,
		Data: fileInfos,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
