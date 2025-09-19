package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	fmt.Println("=== 启动简化服务 ===")

	// 健康检查
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","message":"简化服务运行正常","timestamp":"%s"}`, 
			time.Now().Format("2006-01-02 15:04:05"))
	})

	// 根路径
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
    <title>安语颐年服务</title>
</head>
<body>
    <h1>安语颐年护理陪诊服务</h1>
    <p>服务运行正常</p>
    <p>时间: %s</p>
    <p>路径: %s</p>
    <p>方法: %s</p>
</body>
</html>`, time.Now().Format("2006-01-02 15:04:05"), r.URL.Path, r.Method)
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
		fmt.Printf("=== 收到管理员API请求 ===\n")
		fmt.Printf("方法: %s\n", r.Method)
		fmt.Printf("路径: %s\n", r.URL.Path)
		fmt.Printf("时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"code":0,"message":"管理员API测试正常","timestamp":"%s","path":"%s","method":"%s"}`, 
			time.Now().Format("2006-01-02 15:04:05"), r.URL.Path, r.Method)
	})

	// 分类API测试
	http.HandleFunc("/api/service/categories", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"code":0,"data":[{"id":1,"name":"智慧养老","icon":"elderly"}]}`)
	})

	fmt.Println("简化服务启动在端口 80")
	fmt.Println("健康检查: http://localhost/health")
	fmt.Println("测试API: http://localhost/api/test")
	fmt.Println("管理员API: http://localhost/api/admin/service/update-price")
	fmt.Println("分类API: http://localhost/api/service/categories")

	if err := http.ListenAndServe(":80", nil); err != nil {
		fmt.Printf("服务启动失败: %v\n", err)
	}
}
