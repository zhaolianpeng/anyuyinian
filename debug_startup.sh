#!/bin/bash

# 服务启动诊断脚本

echo "=== 服务启动诊断 ==="

# 1. 检查本地构建
echo "1. 检查本地构建..."
if [ -f "main" ]; then
    echo "✅ 本地构建文件存在"
    ls -la main
    file main
else
    echo "❌ 本地构建文件不存在"
    echo "正在构建..."
    ./build.sh
fi

# 2. 检查Go模块
echo ""
echo "2. 检查Go模块..."
go mod tidy
go mod verify

# 3. 检查代码语法
echo ""
echo "3. 检查代码语法..."
if go build -o test_build .; then
    echo "✅ 代码语法检查通过"
    rm -f test_build
else
    echo "❌ 代码语法检查失败"
    exit 1
fi

# 4. 检查数据库配置
echo ""
echo "4. 检查数据库配置..."
echo "数据库配置信息："
grep -A 5 "mysql init" db/init.go

# 5. 检查环境变量
echo ""
echo "5. 检查环境变量..."
echo "当前环境变量："
env | grep -E "(MYSQL|DB|DATABASE)" || echo "未找到数据库相关环境变量"

# 6. 检查端口占用
echo ""
echo "6. 检查端口占用..."
if command -v lsof >/dev/null 2>&1; then
    PORT_80=$(lsof -i :80 2>/dev/null | wc -l)
    if [ "$PORT_80" -gt 0 ]; then
        echo "⚠️  端口80被占用："
        lsof -i :80
    else
        echo "✅ 端口80可用"
    fi
else
    echo "⚠️  无法检查端口占用（lsof命令不可用）"
fi

# 7. 检查Dockerfile
echo ""
echo "7. 检查Dockerfile..."
if [ -f "Dockerfile" ]; then
    echo "✅ Dockerfile存在"
    echo "Dockerfile内容："
    head -20 Dockerfile
else
    echo "❌ Dockerfile不存在"
fi

# 8. 检查容器配置
echo ""
echo "8. 检查容器配置..."
if [ -f "container.config.json" ]; then
    echo "✅ 容器配置文件存在"
    echo "容器端口配置："
    grep -E "(containerPort|port)" container.config.json
else
    echo "❌ 容器配置文件不存在"
fi

# 9. 提供解决方案
echo ""
echo "=== 解决方案建议 ==="

echo "1. 检查云托管服务日志："
echo "   - 登录腾讯云控制台"
echo "   - 进入云托管服务"
echo "   - 查看 golang-lfwy 服务日志"
echo "   - 查找启动错误信息"

echo ""
echo "2. 如果数据库连接失败："
echo "   - 检查数据库服务是否可用"
echo "   - 检查网络连接"
echo "   - 检查数据库凭据"

echo ""
echo "3. 如果服务启动失败："
echo "   - 检查内存和CPU限制"
echo "   - 检查端口配置"
echo "   - 检查依赖服务"

echo ""
echo "4. 使用测试版本："
echo "   - 部署不依赖数据库的测试版本"
echo "   - 验证服务基本功能"

echo ""
echo "=== 诊断完成 ==="
