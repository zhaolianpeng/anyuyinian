#!/bin/bash

# WebSocket连接测试脚本

echo "=== WebSocket连接测试 ==="

# 测试后端服务是否运行
echo "1. 测试后端服务状态..."

# 检查服务是否在运行
if curl -s http://localhost:80/ > /dev/null; then
    echo "✅ 后端服务正在运行"
else
    echo "❌ 后端服务未运行或无法访问"
    echo "请确保后端服务已启动并监听在端口80"
    exit 1
fi

# 测试WebSocket端点
echo ""
echo "2. 测试WebSocket端点..."

# 使用wscat测试WebSocket连接（如果安装了wscat）
if command -v wscat &> /dev/null; then
    echo "使用wscat测试WebSocket连接..."
    timeout 10s wscat -c ws://localhost:80/ws || echo "WebSocket连接失败"
else
    echo "wscat未安装，使用curl测试WebSocket升级..."
    curl -i -N -H "Connection: Upgrade" -H "Upgrade: websocket" -H "Sec-WebSocket-Version: 13" -H "Sec-WebSocket-Key: x3JJHMbDL1EzLkh9GBhXDw==" http://localhost:80/ws
fi

echo ""
echo "3. 检查服务日志..."
echo "请检查后端服务日志以获取更多信息"

echo ""
echo "=== 测试完成 ===" 