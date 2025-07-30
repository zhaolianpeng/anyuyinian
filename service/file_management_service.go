package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"wxcloudrun-golang/db/dao"
)

// FileManagementResponse 文件管理响应
type FileManagementResponse struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

// DeleteFileHandler 删除文件接口
func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	LogStep("开始处理删除文件请求", map[string]string{"method": r.Method, "path": r.URL.Path})

	if r.Method != http.MethodDelete {
		LogError("请求方法不支持", fmt.Errorf("期望DELETE方法，实际为%s", r.Method))
		http.Error(w, "只支持DELETE请求", http.StatusMethodNotAllowed)
		return
	}

	// 获取文件ID
	fileIdStr := r.URL.Query().Get("fileId")
	if fileIdStr == "" {
		LogError("缺少文件ID参数", fmt.Errorf("fileId参数为空"))
		http.Error(w, "缺少文件ID参数", http.StatusBadRequest)
		return
	}

	fileId, err := strconv.Atoi(fileIdStr)
	if err != nil {
		LogError("无效的文件ID", err)
		http.Error(w, "无效的文件ID", http.StatusBadRequest)
		return
	}

	LogStep("开始删除文件", map[string]interface{}{"fileId": fileId})

	// 查询文件信息
	LogDBOperation("查询", "files", map[string]interface{}{"id": fileId})
	file, err := dao.UploadImp.GetFileById(int32(fileId))
	LogDBResult("查询", "files", file, err)

	if err != nil {
		LogError("查询文件信息失败", err)
		response := &FileManagementResponse{
			Code:     -1,
			ErrorMsg: "文件不存在或查询失败",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		LogResponse(response, err)
		return
	}

	// 删除COS中的文件
	permissionService := NewCOSPermissionService()
	fileName := file.FileName
	err = permissionService.DeleteObject(fileName)
	if err != nil {
		LogError("删除COS文件失败", err)
		// 注意：这里不返回错误，继续删除数据库记录
	} else {
		LogStep("COS文件删除成功", map[string]string{"fileName": fileName})
	}

	// 删除数据库记录
	LogDBOperation("删除", "files", map[string]interface{}{"id": fileId})
	err = dao.UploadImp.DeleteFile(int32(fileId))
	LogDBResult("删除", "files", nil, err)

	if err != nil {
		LogError("删除数据库记录失败", err)
		response := &FileManagementResponse{
			Code:     -1,
			ErrorMsg: "删除文件记录失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		LogResponse(response, err)
		return
	}

	LogStep("文件删除完成", map[string]interface{}{
		"fileId":   fileId,
		"fileName": fileName,
	})

	response := &FileManagementResponse{
		Code: 0,
		Data: map[string]interface{}{
			"fileId":   fileId,
			"fileName": fileName,
			"message":  "文件删除成功",
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	LogResponse(response, nil)
}

// UpdateFilePermissionHandler 更新文件权限接口
func UpdateFilePermissionHandler(w http.ResponseWriter, r *http.Request) {
	LogStep("开始处理更新文件权限请求", map[string]string{"method": r.Method, "path": r.URL.Path})

	if r.Method != http.MethodPut {
		LogError("请求方法不支持", fmt.Errorf("期望PUT方法，实际为%s", r.Method))
		http.Error(w, "只支持PUT请求", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体
	var req struct {
		FileId int32  `json:"fileId"`
		ACL    string `json:"acl"` // "public-read" 或 "private"
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogError("解析请求体失败", err)
		http.Error(w, "解析请求体失败", http.StatusBadRequest)
		return
	}

	LogStep("解析请求参数", map[string]interface{}{
		"fileId": req.FileId,
		"acl":    req.ACL,
	})

	// 验证ACL参数
	if req.ACL != "public-read" && req.ACL != "private" {
		LogError("无效的ACL参数", fmt.Errorf("ACL必须是public-read或private，实际为%s", req.ACL))
		http.Error(w, "无效的ACL参数", http.StatusBadRequest)
		return
	}

	// 查询文件信息
	LogDBOperation("查询", "files", map[string]interface{}{"id": req.FileId})
	file, err := dao.UploadImp.GetFileById(req.FileId)
	LogDBResult("查询", "files", file, err)

	if err != nil {
		LogError("查询文件信息失败", err)
		response := &FileManagementResponse{
			Code:     -1,
			ErrorMsg: "文件不存在或查询失败",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		LogResponse(response, err)
		return
	}

	// 更新COS对象权限
	permissionService := NewCOSPermissionService()
	var updateErr error

	if req.ACL == "public-read" {
		updateErr = permissionService.SetObjectPublicRead(file.FileName)
	} else {
		updateErr = permissionService.SetObjectPrivate(file.FileName)
	}

	if updateErr != nil {
		LogError("更新COS对象权限失败", updateErr)
		response := &FileManagementResponse{
			Code:     -1,
			ErrorMsg: "更新文件权限失败: " + updateErr.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		LogResponse(response, updateErr)
		return
	}

	LogStep("文件权限更新成功", map[string]interface{}{
		"fileId":   req.FileId,
		"fileName": file.FileName,
		"acl":      req.ACL,
	})

	response := &FileManagementResponse{
		Code: 0,
		Data: map[string]interface{}{
			"fileId":   req.FileId,
			"fileName": file.FileName,
			"acl":      req.ACL,
			"message":  "文件权限更新成功",
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	LogResponse(response, nil)
}

// GetFilePermissionHandler 获取文件权限接口
func GetFilePermissionHandler(w http.ResponseWriter, r *http.Request) {
	LogStep("开始处理获取文件权限请求", map[string]string{"method": r.Method, "path": r.URL.Path})

	if r.Method != http.MethodGet {
		LogError("请求方法不支持", fmt.Errorf("期望GET方法，实际为%s", r.Method))
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	// 获取文件ID
	fileIdStr := r.URL.Query().Get("fileId")
	if fileIdStr == "" {
		LogError("缺少文件ID参数", fmt.Errorf("fileId参数为空"))
		http.Error(w, "缺少文件ID参数", http.StatusBadRequest)
		return
	}

	fileId, err := strconv.Atoi(fileIdStr)
	if err != nil {
		LogError("无效的文件ID", err)
		http.Error(w, "无效的文件ID", http.StatusBadRequest)
		return
	}

	LogStep("开始获取文件权限", map[string]interface{}{"fileId": fileId})

	// 查询文件信息
	LogDBOperation("查询", "files", map[string]interface{}{"id": fileId})
	file, err := dao.UploadImp.GetFileById(int32(fileId))
	LogDBResult("查询", "files", file, err)

	if err != nil {
		LogError("查询文件信息失败", err)
		response := &FileManagementResponse{
			Code:     -1,
			ErrorMsg: "文件不存在或查询失败",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		LogResponse(response, err)
		return
	}

	// 获取COS对象权限
	permissionService := NewCOSPermissionService()
	acl, err := permissionService.GetObjectACL(file.FileName)
	if err != nil {
		LogError("获取COS对象权限失败", err)
		response := &FileManagementResponse{
			Code:     -1,
			ErrorMsg: "获取文件权限失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		LogResponse(response, err)
		return
	}

	LogStep("获取文件权限成功", map[string]interface{}{
		"fileId":   fileId,
		"fileName": file.FileName,
		"acl":      acl,
	})

	response := &FileManagementResponse{
		Code: 0,
		Data: map[string]interface{}{
			"fileId":   fileId,
			"fileName": file.FileName,
			"acl":      acl,
			"fileUrl":  file.FileUrl,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	LogResponse(response, nil)
}
