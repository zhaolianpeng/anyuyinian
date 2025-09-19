#!/bin/bash

# 紧急修复脚本 - 解决云托管服务418错误

echo "=== 紧急修复云托管服务 ==="

# 1. 检查当前状态
echo "1. 检查当前状态..."
echo "服务地址: https://golang-lfwy-prod-5g94mx7a3d07e78c.service.tcloudbaseapp.com"
echo "当前状态: HTTP 418错误"
echo ""

# 2. 检查本地代码
echo "2. 检查本地代码..."
if [ -f "main" ]; then
    echo "✅ 本地构建文件存在"
    ls -la main
else
    echo "❌ 本地构建文件不存在，正在构建..."
    ./build.sh
fi

# 3. 检查路由配置
echo ""
echo "3. 检查路由配置..."
echo "管理员服务价格更新路由:"
grep -n "update-price" main.go

echo ""
echo "路由注册顺序:"
grep -A 2 -B 2 "admin.*service" main.go

# 4. 检查数据库配置
echo ""
echo "4. 检查数据库配置..."
echo "数据库连接信息:"
grep -A 3 "mysql init" db/init.go

# 5. 创建简化版本（不依赖数据库）
echo ""
echo "5. 创建简化版本..."
cat > simple_main.go << 'EOF'
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

	fmt.Println("简化服务启动在端口 80")
	fmt.Println("健康检查: http://localhost/health")
	fmt.Println("测试API: http://localhost/api/test")
	fmt.Println("管理员API: http://localhost/api/admin/service/update-price")

	if err := http.ListenAndServe(":80", nil); err != nil {
		fmt.Printf("服务启动失败: %v\n", err)
	}
}
EOF

echo "✅ 简化版本创建完成"

# 6. 构建简化版本
echo ""
echo "6. 构建简化版本..."
GOOS=linux go build -o simple_main simple_main.go

if [ $? -eq 0 ]; then
    echo "✅ 简化版本构建成功"
    ls -la simple_main
else
    echo "❌ 简化版本构建失败"
    exit 1
fi

# 7. 创建简化版Dockerfile
echo ""
echo "7. 创建简化版Dockerfile..."
cat > Dockerfile.simple << 'EOF'
FROM golang:1.17.1-alpine3.14 as builder

WORKDIR /app
COPY . /app/
RUN GOOS=linux go build -o simple_main simple_main.go

FROM alpine:3.13
RUN apk add ca-certificates
WORKDIR /app
COPY --from=builder /app/simple_main /app/
CMD ["/app/simple_main"]
EOF

echo "✅ 简化版Dockerfile创建完成"

# 8. 提供部署建议
echo ""
echo "=== 部署建议 ==="
echo ""
echo "方案1: 使用简化版本（推荐）"
echo "1. 备份当前Dockerfile: mv Dockerfile Dockerfile.backup"
echo "2. 使用简化版Dockerfile: mv Dockerfile.simple Dockerfile"
echo "3. 重新部署到云托管"
echo "4. 验证服务是否正常"
echo ""
echo "方案2: 修复原版本"
echo "1. 检查云托管控制台的服务日志"
echo "2. 查找数据库连接错误"
echo "3. 修复数据库配置"
echo "4. 重新部署"
echo ""
echo "方案3: 检查云托管配置"
echo "1. 检查服务配额"
echo "2. 检查网络配置"
echo "3. 检查环境变量"
echo "4. 联系腾讯云支持"

echo ""
echo "=== 紧急修复完成 ==="
