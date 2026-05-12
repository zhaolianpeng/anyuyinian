package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const mediaServerBaseURL = "https://api.succ.online/anyuyinian"
const maxImageFileSize = 5 * 1024 * 1024

type MediaImagesRequest struct {
	Paths []string `json:"paths"`
}

type MediaImageItem struct {
	Path     string `json:"path"`
	MimeType string `json:"mimeType,omitempty"`
	Base64   string `json:"base64,omitempty"`
	Error    string `json:"error,omitempty"`
	Status   int    `json:"status"`
}

type MediaImagesResponse struct {
	Items []MediaImageItem `json:"items"`
}

func MediaImagesHandler(w http.ResponseWriter, r *http.Request) {
	LogInfo("开始处理媒体图片批量读取请求", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodPost {
		LogError("请求方法不支持", fmt.Errorf("期望POST方法，实际为%s", r.Method))
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	var req MediaImagesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogError("解析媒体图片请求失败", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&ServiceResponse{
			Code:     -1,
			ErrorMsg: "请求参数解析失败: " + err.Error(),
		})
		return
	}

	staticRoot := os.Getenv("STATIC_ROOT")
	if staticRoot == "" {
		staticRoot = "/home/ubuntu/data"
	}

	items := make([]MediaImageItem, 0, len(req.Paths))
	for _, mediaPath := range req.Paths {
		item := MediaImageItem{Path: mediaPath, Status: 1}

		fullPath, err := resolveMediaImagePath(staticRoot, mediaPath)
		if err != nil {
			item.Error = err.Error()
			items = append(items, item)
			continue
		}

		fileInfo, err := os.Stat(fullPath)
		if err != nil {
			item.Error = "读取图片信息失败: " + err.Error()
			items = append(items, item)
			continue
		}

		if fileInfo.Size() > maxImageFileSize {
			item.Error = fmt.Sprintf("图片过大，超过 %d 字节限制", maxImageFileSize)
			items = append(items, item)
			continue
		}

		content, err := os.ReadFile(fullPath)
		if err != nil {
			item.Error = "读取图片失败: " + err.Error()
			items = append(items, item)
			continue
		}

		mimeType := mime.TypeByExtension(strings.ToLower(filepath.Ext(fullPath)))
		if mimeType == "" {
			mimeType = http.DetectContentType(content)
		}

		item.MimeType = mimeType
		item.Base64 = base64.StdEncoding.EncodeToString(content)
		item.Status = 0
		items = append(items, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&ServiceResponse{
		Code: 0,
		Data: &MediaImagesResponse{Items: items},
	})

	LogInfo("媒体图片批量读取完成", map[string]interface{}{
		"count": len(items),
	})
}

func resolveMediaImagePath(staticRoot, mediaPath string) (string, error) {
	trimmed := strings.TrimSpace(mediaPath)
	if trimmed == "" {
		return "", fmt.Errorf("图片路径为空")
	}

	trimmed = strings.TrimPrefix(trimmed, mediaServerBaseURL)

	trimmed = strings.ReplaceAll(trimmed, "\\", "/")

	var baseDir string
	var relativePath string

	switch {
	case strings.HasPrefix(trimmed, "/static/"):
		baseDir = filepath.Join(staticRoot, "static")
		relativePath = strings.TrimPrefix(trimmed, "/static/")
	case strings.HasPrefix(trimmed, "/images/"):
		baseDir = filepath.Join(staticRoot, "images")
		relativePath = strings.TrimPrefix(trimmed, "/images/")
	default:
		return "", fmt.Errorf("不支持的图片路径: %s", mediaPath)
	}

	cleanBaseDir := filepath.Clean(baseDir)
	fullPath := filepath.Clean(filepath.Join(cleanBaseDir, filepath.FromSlash(relativePath)))
	relToBase, err := filepath.Rel(cleanBaseDir, fullPath)
	if err != nil {
		return "", fmt.Errorf("解析图片路径失败: %w", err)
	}

	if strings.HasPrefix(relToBase, "..") {
		return "", fmt.Errorf("图片路径非法: %s", mediaPath)
	}

	return fullPath, nil
}
