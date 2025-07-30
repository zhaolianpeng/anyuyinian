package service

import (
	"context"
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

	"wxcloudrun-golang/config"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"

	"github.com/tencentyun/cos-go-sdk-v5"
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
	UploadPath   string   // 本地上传路径（临时存储）
	BaseURL      string   // 基础URL
	UseCOS       bool     // 是否使用腾讯云COS
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
	UseCOS:     true, // 启用腾讯云COS
}

// UploadHandler 文件上传接口
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	LogStep("开始处理文件上传请求", map[string]string{"method": r.Method, "path": r.URL.Path})

	if r.Method != http.MethodPost {
		LogError("请求方法不支持", fmt.Errorf("期望POST方法，实际为%s", r.Method))
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	// 解析multipart表单
	LogStep("开始解析multipart表单", nil)
	err := r.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		LogError("解析表单失败", err)
		http.Error(w, "解析表单失败", http.StatusBadRequest)
		return
	}
	LogStep("表单解析成功", nil)

	// 获取上传的文件
	file, header, err := r.FormFile("file")
	if err != nil {
		LogError("获取上传文件失败", err)
		http.Error(w, "获取上传文件失败", http.StatusBadRequest)
		return
	}
	defer file.Close()

	LogStep("获取上传文件成功", map[string]interface{}{
		"fileName":    header.Filename,
		"fileSize":    header.Size,
		"contentType": header.Header.Get("Content-Type"),
	})

	// 获取其他参数
	userIdStr := r.FormValue("userId")
	category := r.FormValue("category")
	description := r.FormValue("description")

	LogStep("获取表单参数", map[string]interface{}{
		"userId":      userIdStr,
		"category":    category,
		"description": description,
	})

	// 验证用户ID
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		LogError("无效的用户ID", err)
		http.Error(w, "无效的用户ID", http.StatusBadRequest)
		return
	}
	LogStep("用户ID验证成功", map[string]interface{}{"userId": userId})

	// 验证文件
	LogStep("开始验证文件", map[string]interface{}{
		"fileName":    header.Filename,
		"fileSize":    header.Size,
		"contentType": header.Header.Get("Content-Type"),
	})
	if err := validateFile(header); err != nil {
		LogError("文件验证失败", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	LogStep("文件验证通过", nil)

	// 保存文件
	fileInfo, err := saveFile(file, header, int32(userId), category, description)
	if err != nil {
		LogError("保存文件失败", err)
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
	LogResponse(response, nil)
}

// validateFile 验证文件
func validateFile(header *multipart.FileHeader) error {
	LogStep("开始验证文件", map[string]interface{}{
		"fileName":    header.Filename,
		"fileSize":    header.Size,
		"maxFileSize": uploadConfig.MaxFileSize,
		"contentType": header.Header.Get("Content-Type"),
	})

	// 检查文件大小
	if header.Size > uploadConfig.MaxFileSize {
		LogError("文件大小超过限制", fmt.Errorf("文件大小: %d, 最大允许: %d", header.Size, uploadConfig.MaxFileSize))
		return fmt.Errorf("文件大小超过限制，最大允许 %d 字节", uploadConfig.MaxFileSize)
	}
	LogStep("文件大小验证通过", map[string]interface{}{"fileSize": header.Size})

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
		LogError("不支持的文件类型", fmt.Errorf("文件类型: %s", contentType))
		return fmt.Errorf("不支持的文件类型: %s", contentType)
	}
	LogStep("文件类型验证通过", map[string]string{"contentType": contentType})

	return nil
}

// saveFile 保存文件
func saveFile(file multipart.File, header *multipart.FileHeader, userId int32, category, description string) (*FileInfo, error) {
	LogStep("开始保存文件", map[string]string{"originalName": header.Filename, "size": fmt.Sprintf("%d", header.Size)})

	// 生成文件名
	originalName := header.Filename
	fileName := generateFileName(originalName)

	var fileUrl string
	var fileSize int64

	if uploadConfig.UseCOS {
		// 使用腾讯云COS上传
		LogStep("使用腾讯云COS上传文件", map[string]string{"fileName": fileName})
		cosUrl, err := uploadToCOS(file, fileName, header.Header.Get("Content-Type"))
		if err != nil {
			LogError("COS上传失败", err)
			return nil, fmt.Errorf("上传到COS失败: %v", err)
		}
		fileUrl = cosUrl
		fileSize = header.Size
		LogStep("COS上传成功", map[string]string{"fileUrl": fileUrl})
	} else {
		// 本地保存
		LogStep("使用本地存储保存文件", map[string]string{"fileName": fileName})

		// 创建上传目录
		if err := os.MkdirAll(uploadConfig.UploadPath, 0755); err != nil {
			LogError("创建上传目录失败", err)
			return nil, fmt.Errorf("创建上传目录失败: %v", err)
		}

		fullPath := filepath.Join(uploadConfig.UploadPath, fileName)

		// 创建目标文件
		dst, err := os.Create(fullPath)
		if err != nil {
			LogError("创建文件失败", err)
			return nil, fmt.Errorf("创建文件失败: %v", err)
		}
		defer dst.Close()

		// 复制文件内容
		_, err = io.Copy(dst, file)
		if err != nil {
			LogError("保存文件失败", err)
			return nil, fmt.Errorf("保存文件失败: %v", err)
		}

		// 获取文件信息
		fileInfo, err := dst.Stat()
		if err != nil {
			LogError("获取文件信息失败", err)
			return nil, fmt.Errorf("获取文件信息失败: %v", err)
		}

		fileUrl = uploadConfig.BaseURL + "/" + fileName
		fileSize = fileInfo.Size()
		LogStep("本地文件保存成功", map[string]string{"filePath": fullPath, "fileUrl": fileUrl})
	}

	// 确定文件类型和分类
	contentType := header.Header.Get("Content-Type")
	fileType := getFileType(contentType)
	if category == "" {
		category = getCategoryByType(fileType)
	}

	LogStep("准备保存文件记录到数据库", map[string]interface{}{
		"fileName": fileName,
		"fileUrl":  fileUrl,
		"fileSize": fileSize,
		"fileType": fileType,
		"category": category,
	})

	// 创建数据库记录
	dbFile := &model.FileModel{
		FileName:     fileName,
		OriginalName: originalName,
		FilePath:     fileUrl, // 对于COS，FilePath存储URL
		FileUrl:      fileUrl,
		FileSize:     fileSize,
		FileType:     fileType,
		MimeType:     contentType,
		Category:     category,
		Description:  description,
		UserId:       userId,
	}

	LogDBOperation("创建", "files", dbFile)
	if err := dao.UploadImp.CreateFile(dbFile); err != nil {
		LogError("保存文件记录失败", err)
		// 如果数据库保存失败，删除已上传的文件（仅本地文件）
		if !uploadConfig.UseCOS {
			os.Remove(filepath.Join(uploadConfig.UploadPath, fileName))
		}
		return nil, fmt.Errorf("保存文件记录失败: %v", err)
	}
	LogDBResult("创建", "files", dbFile, nil)

	LogStep("文件保存完成", map[string]interface{}{
		"fileId":  dbFile.Id,
		"fileUrl": dbFile.FileUrl,
	})

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

// uploadToCOS 上传文件到腾讯云COS
func uploadToCOS(file multipart.File, fileName, contentType string) (string, error) {
	LogStep("开始上传文件到COS", map[string]string{"fileName": fileName, "contentType": contentType})

	// 获取COS客户端
	client := config.GetCOSClient()
	cosConfig := config.GetCOSConfig()

	// 重置文件指针到开始位置
	file.Seek(0, 0)

	// 上传文件到COS，设置ACL权限为public-read（所有用户可读，仅创建者可读写）
	_, err := client.Object.Put(context.Background(), fileName, file, &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: contentType,
		},
	})
	if err != nil {
		LogError("COS上传失败", err)
		return "", fmt.Errorf("上传到COS失败: %v", err)
	}

	// 构建文件URL
	fileUrl := cosConfig.Domain + "/" + fileName
	LogStep("COS上传成功", map[string]string{"fileUrl": fileUrl, "acl": cosConfig.ACL})

	// 设置对象ACL权限为public-read（所有用户可读，仅创建者可读写）
	permissionService := NewCOSPermissionService()
	err = permissionService.SetObjectPublicRead(fileName)
	if err != nil {
		LogError("设置对象ACL失败", err)
		// 注意：这里不返回错误，因为文件已经上传成功，只是权限设置失败
	} else {
		LogStep("对象ACL设置成功", map[string]string{"fileName": fileName, "acl": "public-read"})
	}

	return fileUrl, nil
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
	LogStep("开始处理获取文件列表请求", map[string]string{"method": r.Method, "path": r.URL.Path})

	if r.Method != http.MethodGet {
		LogError("请求方法不支持", fmt.Errorf("期望GET方法，实际为%s", r.Method))
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

	LogStep("获取查询参数", map[string]interface{}{
		"userId":   userIdStr,
		"category": category,
		"limit":    limit,
	})

	var files []*model.FileModel
	var err error

	if userIdStr != "" {
		// 根据用户ID查询
		LogStep("根据用户ID查询文件", map[string]interface{}{"userId": userIdStr, "limit": limit})
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			LogError("无效的用户ID", err)
			http.Error(w, "无效的用户ID", http.StatusBadRequest)
			return
		}
		LogDBOperation("查询", "files", map[string]interface{}{"userId": userId, "limit": limit})
		files, err = dao.UploadImp.GetFilesByUserId(int32(userId), limit)
		LogDBResult("查询", "files", files, err)
	} else if category != "" {
		// 根据分类查询
		LogStep("根据分类查询文件", map[string]interface{}{"category": category, "limit": limit})
		LogDBOperation("查询", "files", map[string]interface{}{"category": category, "limit": limit})
		files, err = dao.UploadImp.GetFilesByCategory(category, limit)
		LogDBResult("查询", "files", files, err)
	} else {
		LogError("缺少查询参数", fmt.Errorf("userId和category参数都为空"))
		http.Error(w, "缺少查询参数", http.StatusBadRequest)
		return
	}

	if err != nil {
		LogError("获取文件列表失败", err)
		response := &UploadResponse{
			Code:     -1,
			ErrorMsg: "获取文件列表失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		LogResponse(response, err)
		return
	}

	LogStep("文件列表查询成功", map[string]interface{}{"fileCount": len(files)})

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
	LogResponse(response, nil)
}
