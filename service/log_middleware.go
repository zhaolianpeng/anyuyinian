package service

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

// LogMiddleware 日志中间件
type LogMiddleware struct {
	handler http.HandlerFunc
}

// NewLogMiddleware 创建日志中间件
func NewLogMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// 记录请求开始
		log.Printf("[API] 开始处理请求: %s %s", r.Method, r.URL.Path)

		// 记录请求头信息
		log.Printf("[API] 请求头: User-Agent=%s, Content-Type=%s",
			r.Header.Get("User-Agent"), r.Header.Get("Content-Type"))

		// 记录请求参数
		if r.Method == "GET" {
			log.Printf("[API] 查询参数: %s", r.URL.RawQuery)
		} else if r.Method == "POST" {
			// 读取请求体
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("[API] 读取请求体失败: %v", err)
			} else {
				log.Printf("[API] 请求体: %s", string(body))
				// 重新设置请求体，因为已经被读取了
				// 使用更可靠的方式重新设置请求体
				r.Body = io.NopCloser(bytes.NewReader(body))
			}
		}

		// 创建响应记录器
		responseRecorder := &ResponseRecorder{
			ResponseWriter: w,
			body:           &bytes.Buffer{},
		}

		// 调用原始处理器
		handler(responseRecorder, r)

		// 记录响应
		duration := time.Since(startTime)
		log.Printf("[API] 响应状态: %d", responseRecorder.statusCode)
		log.Printf("[API] 响应体: %s", responseRecorder.body.String())
		log.Printf("[API] 请求处理完成，耗时: %v", duration)
	}
}

// ResponseRecorder 响应记录器
type ResponseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

// WriteHeader 重写WriteHeader方法
func (r *ResponseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// Write 重写Write方法
func (r *ResponseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// LogRequest 记录请求日志
func LogRequest(method, path string, params interface{}) {
	if params != nil {
		paramsJSON, _ := json.Marshal(params)
		log.Printf("[API] 请求参数: %s", string(paramsJSON))
	}
}

// LogResponse 记录响应日志
func LogResponse(data interface{}, err error) {
	if err != nil {
		log.Printf("[API] 响应错误: %v", err)
	} else if data != nil {
		dataJSON, _ := json.Marshal(data)
		log.Printf("[API] 响应数据: %s", string(dataJSON))
	}
}

// LogDBOperation 记录数据库操作日志
func LogDBOperation(operation string, table string, params interface{}) {
	if params != nil {
		paramsJSON, _ := json.Marshal(params)
		log.Printf("[DB] %s 表: %s, 参数: %s", operation, table, string(paramsJSON))
	} else {
		log.Printf("[DB] %s 表: %s", operation, table)
	}
}

// LogDBResult 记录数据库操作结果
func LogDBResult(operation string, table string, result interface{}, err error) {
	if err != nil {
		log.Printf("[DB] %s 表: %s 失败: %v", operation, table, err)
	} else {
		if result != nil {
			resultJSON, _ := json.Marshal(result)
			log.Printf("[DB] %s 表: %s 成功: %s", operation, table, string(resultJSON))
		} else {
			log.Printf("[DB] %s 表: %s 成功", operation, table)
		}
	}
}

// LogStep 记录中间步骤日志
func LogStep(step string, data interface{}) {
	if data != nil {
		dataJSON, _ := json.Marshal(data)
		log.Printf("[STEP] %s: %s", step, string(dataJSON))
	} else {
		log.Printf("[STEP] %s", step)
	}
}

// LogError 记录错误日志
func LogError(message string, err error) {
	if err != nil {
		log.Printf("[ERROR] %s: %v", message, err)
	} else {
		log.Printf("[ERROR] %s", message)
	}
}

// LogInfo 记录信息日志
func LogInfo(message string, data interface{}) {
	if data != nil {
		dataJSON, _ := json.Marshal(data)
		log.Printf("[INFO] %s: %s", message, string(dataJSON))
	} else {
		log.Printf("[INFO] %s", message)
	}
}
