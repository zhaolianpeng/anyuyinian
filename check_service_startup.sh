#!/bin/bash

# 服务启动检查脚本

echo "=== 服务启动检查 ==="

SERVICE_URL="https://golang-lfwy-prod-5g94mx7a3d07e78c.service.tcloudbaseapp.com"

# 1. 检查服务根路径
echo "1. 检查服务根路径..."
ROOT_RESPONSE=$(curl -k -s -w "HTTP_CODE:%{http_code}" -o /tmp/root_response.txt "$SERVICE_URL/" 2>/dev/null)
ROOT_HTTP_CODE=$(echo "$ROOT_RESPONSE" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
ROOT_CONTENT=$(cat /tmp/root_response.txt 2>/dev/null)

echo "根路径响应: HTTP $ROOT_HTTP_CODE"
if [ -n "$ROOT_CONTENT" ]; then
    echo "响应内容: $ROOT_CONTENT"
fi

# 2. 检查健康检查端点
echo ""
echo "2. 检查健康检查端点..."
HEALTH_RESPONSE=$(curl -k -s -w "HTTP_CODE:%{http_code}" -o /tmp/health_response.txt "$SERVICE_URL/health" 2>/dev/null)
HEALTH_HTTP_CODE=$(echo "$HEALTH_RESPONSE" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
HEALTH_CONTENT=$(cat /tmp/health_response.txt 2>/dev/null)

echo "健康检查响应: HTTP $HEALTH_HTTP_CODE"
if [ -n "$HEALTH_CONTENT" ]; then
    echo "响应内容: $HEALTH_CONTENT"
fi

# 3. 检查服务头信息
echo ""
echo "3. 检查服务头信息..."
HEADERS_RESPONSE=$(curl -k -s -I "$SERVICE_URL/" 2>/dev/null)
echo "服务头信息:"
echo "$HEADERS_RESPONSE"

# 4. 分析问题
echo ""
echo "=== 问题分析 ==="

if [ "$ROOT_HTTP_CODE" = "418" ]; then
    echo "❌ 服务返回418错误"
    echo "可能原因："
    echo "1. 服务启动失败"
    echo "2. 数据库连接失败"
    echo "3. 端口配置错误"
    echo "4. 内存不足"
    echo "5. 服务内部错误"
elif [ "$ROOT_HTTP_CODE" = "200" ]; then
    echo "✅ 服务根路径正常"
elif [ "$ROOT_HTTP_CODE" = "404" ]; then
    echo "⚠️  服务返回404错误"
    echo "可能原因："
    echo "1. 路由配置问题"
    echo "2. 服务未完全启动"
    echo "3. 缓存问题"
else
    echo "❓ 服务返回未知错误: HTTP $ROOT_HTTP_CODE"
fi

# 5. 提供解决方案
echo ""
echo "=== 解决方案 ==="

if [ "$ROOT_HTTP_CODE" = "418" ]; then
    echo "🚨 服务启动失败，建议："
    echo "1. 检查云托管控制台的服务日志"
    echo "2. 查看服务启动错误信息"
    echo "3. 检查数据库连接配置"
    echo "4. 检查服务资源限制"
    echo "5. 尝试使用简化版本部署"
elif [ "$ROOT_HTTP_CODE" = "404" ]; then
    echo "⚠️  路由问题，建议："
    echo "1. 等待服务完全启动（通常需要2-5分钟）"
    echo "2. 检查路由配置"
    echo "3. 清除浏览器缓存"
    echo "4. 重新测试"
else
    echo "✅ 服务运行正常"
fi

# 清理临时文件
rm -f /tmp/root_response.txt /tmp/health_response.txt

echo ""
echo "=== 检查完成 ==="
