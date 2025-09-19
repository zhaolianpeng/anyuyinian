#!/bin/bash

# 云托管服务详细诊断脚本

echo "=== 云托管服务详细诊断 ==="

# 配置信息
ENV_ID="prod-5g94mx7a3d07e78c"
SERVICE_NAME="golang-lfwy"
SERVICE_URL="https://${SERVICE_NAME}-${ENV_ID}.service.tcloudbaseapp.com"

echo "环境ID: $ENV_ID"
echo "服务名称: $SERVICE_NAME"
echo "服务地址: $SERVICE_URL"
echo ""

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

# 3. 检查分类API
echo ""
echo "3. 检查分类API..."
CATEGORIES_RESPONSE=$(curl -k -s -w "HTTP_CODE:%{http_code}" -o /tmp/categories_response.txt "$SERVICE_URL/api/service/categories" 2>/dev/null)
CATEGORIES_HTTP_CODE=$(echo "$CATEGORIES_RESPONSE" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
CATEGORIES_CONTENT=$(cat /tmp/categories_response.txt 2>/dev/null)

echo "分类API响应: HTTP $CATEGORIES_HTTP_CODE"
if [ -n "$CATEGORIES_CONTENT" ]; then
    echo "响应内容: $CATEGORIES_CONTENT"
fi

# 4. 检查管理员API
echo ""
echo "4. 检查管理员API..."
ADMIN_RESPONSE=$(curl -k -s -w "HTTP_CODE:%{http_code}" -o /tmp/admin_response.txt -X POST "$SERVICE_URL/api/admin/service/update-price" \
  -H "Content-Type: application/json" \
  -d '{"serviceId": 21, "newPrice": 258, "newOriginalPrice": 380, "reason": "测试"}' 2>/dev/null)
ADMIN_HTTP_CODE=$(echo "$ADMIN_RESPONSE" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
ADMIN_CONTENT=$(cat /tmp/admin_response.txt 2>/dev/null)

echo "管理员API响应: HTTP $ADMIN_HTTP_CODE"
if [ -n "$ADMIN_CONTENT" ]; then
    echo "响应内容: $ADMIN_CONTENT"
fi

# 5. 检查服务头信息
echo ""
echo "5. 检查服务头信息..."
HEADERS_RESPONSE=$(curl -k -s -I "$SERVICE_URL/" 2>/dev/null)
echo "服务头信息:"
echo "$HEADERS_RESPONSE"

# 6. 分析问题
echo ""
echo "=== 问题分析 ==="

if [ "$ROOT_HTTP_CODE" = "418" ]; then
    echo "❌ 服务返回418错误 - 这通常表示："
    echo "   1. 服务未正确启动"
    echo "   2. 端口配置错误"
    echo "   3. 服务内部错误"
    echo "   4. 云托管配置问题"
elif [ "$ROOT_HTTP_CODE" = "404" ]; then
    echo "⚠️  服务返回404错误 - 这通常表示："
    echo "   1. 路由配置问题"
    echo "   2. 服务未部署最新代码"
    echo "   3. 服务启动但路由不匹配"
elif [ "$ROOT_HTTP_CODE" = "200" ]; then
    echo "✅ 服务根路径正常"
else
    echo "❓ 服务返回未知错误: HTTP $ROOT_HTTP_CODE"
fi

# 7. 提供解决方案
echo ""
echo "=== 解决方案 ==="

if [ "$ROOT_HTTP_CODE" = "418" ] || [ "$ROOT_HTTP_CODE" = "500" ]; then
    echo "🚨 服务内部错误，建议："
    echo "   1. 检查云托管控制台的服务日志"
    echo "   2. 检查服务配置是否正确"
    echo "   3. 重新部署服务"
    echo "   4. 检查数据库连接"
elif [ "$ROOT_HTTP_CODE" = "404" ]; then
    echo "⚠️  路由问题，建议："
    echo "   1. 确认代码已正确部署"
    echo "   2. 检查路由配置"
    echo "   3. 重新部署服务"
else
    echo "✅ 服务运行正常，问题可能在其他地方"
fi

# 清理临时文件
rm -f /tmp/root_response.txt /tmp/health_response.txt /tmp/categories_response.txt /tmp/admin_response.txt

echo ""
echo "=== 诊断完成 ==="
