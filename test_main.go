package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	fmt.Println("=== 启动测试服务 ===")

	// 健康检查端点
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","message":"测试服务运行正常","timestamp":"%s"}`,
			time.Now().Format("2006-01-02 15:04:05"))
	})

	// 基础页面
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
    <title>安语颐年测试服务</title>
</head>
<body>
    <h1>安语颐年护理陪诊服务</h1>
    <p>服务运行正常</p>
    <p>时间: %s</p>
    <p>路径: %s</p>
</body>
</html>`, time.Now().Format("2006-01-02 15:04:05"), r.URL.Path)
	})

	// 测试API
	http.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"code":0,"message":"测试API正常","timestamp":"%s"}`,
			time.Now().Format("2006-01-02 15:04:05"))
	})

	// 管理员API测试
	http.HandleFunc("/api/admin/service/update-price", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"code":0,"message":"管理员API测试正常","timestamp":"%s"}`,
			time.Now().Format("2006-01-02 15:04:05"))
	})

	fmt.Println("服务启动在端口 80")
	fmt.Println("健康检查: http://localhost/health")
	fmt.Println("测试API: http://localhost/api/test")
	fmt.Println("管理员API: http://localhost/api/admin/service/update-price")

	if err := http.ListenAndServe(":80", nil); err != nil {
		fmt.Printf("服务启动失败: %v\n", err)
	}
}
