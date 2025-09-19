#!/bin/bash

# API调用方式测试脚本

echo "=== API调用方式测试 ==="

SERVICE_URL="https://golang-lfwy-prod-5g94mx7a3d07e78c.service.tcloudbaseapp.com"

# 1. 测试根路径
echo "1. 测试根路径..."
ROOT_RESPONSE=$(curl -k -s -w "HTTP_CODE:%{http_code}" -o /tmp/root_response.txt "$SERVICE_URL/" 2>/dev/null)
ROOT_HTTP_CODE=$(echo "$ROOT_RESPONSE" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
ROOT_CONTENT=$(cat /tmp/root_response.txt 2>/dev/null)

echo "根路径响应: HTTP $ROOT_HTTP_CODE"
if [ -n "$ROOT_CONTENT" ]; then
    echo "响应内容: $ROOT_CONTENT"
fi

# 2. 测试健康检查
echo ""
echo "2. 测试健康检查..."
HEALTH_RESPONSE=$(curl -k -s -w "HTTP_CODE:%{http_code}" -o /tmp/health_response.txt "$SERVICE_URL/health" 2>/dev/null)
HEALTH_HTTP_CODE=$(echo "$HEALTH_RESPONSE" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
HEALTH_CONTENT=$(cat /tmp/health_response.txt 2>/dev/null)

echo "健康检查响应: HTTP $HEALTH_HTTP_CODE"
if [ -n "$HEALTH_CONTENT" ]; then
    echo "响应内容: $HEALTH_CONTENT"
fi

# 3. 测试管理员API（直接HTTP调用）
echo ""
echo "3. 测试管理员API（直接HTTP调用）..."
ADMIN_RESPONSE=$(curl -k -s -w "HTTP_CODE:%{http_code}" -o /tmp/admin_response.txt -X POST "$SERVICE_URL/api/admin/service/updateprice" \
  -H "Content-Type: application/json" \
  -H "X-WX-SERVICE: golang-lfwy" \
  -d '{"serviceId": 21, "newPrice": 280, "newOriginalPrice": 200, "reason": "测试"}' 2>/dev/null)
ADMIN_HTTP_CODE=$(echo "$ADMIN_RESPONSE" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
ADMIN_CONTENT=$(cat /tmp/admin_response.txt 2>/dev/null)

echo "管理员API响应: HTTP $ADMIN_HTTP_CODE"
if [ -n "$ADMIN_CONTENT" ]; then
    echo "响应内容: $ADMIN_CONTENT"
fi

# 4. 测试其他API
echo ""
echo "4. 测试其他API..."
OTHER_RESPONSE=$(curl -k -s -w "HTTP_CODE:%{http_code}" -o /tmp/other_response.txt -X GET "$SERVICE_URL/api/service/categories" 2>/dev/null)
OTHER_HTTP_CODE=$(echo "$OTHER_RESPONSE" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
OTHER_CONTENT=$(cat /tmp/other_response.txt 2>/dev/null)

echo "其他API响应: HTTP $OTHER_HTTP_CODE"
if [ -n "$OTHER_CONTENT" ]; then
    echo "响应内容: $OTHER_CONTENT"
fi

# 5. 分析结果
echo ""
echo "=== 分析结果 ==="

if [ "$ROOT_HTTP_CODE" = "418" ]; then
    echo "❌ 服务启动失败 - 返回418错误"
    echo "问题根因: 服务未正确启动"
    echo "解决方案: 检查云托管日志，修复启动问题"
elif [ "$ROOT_HTTP_CODE" = "200" ]; then
    echo "✅ 服务启动正常 - 返回200成功"
    if [ "$ADMIN_HTTP_CODE" = "200" ]; then
        echo "✅ 管理员API正常 - 返回200成功"
        echo "问题根因: 小程序调用方式问题"
        echo "解决方案: 检查小程序调用配置"
    else
        echo "❌ 管理员API异常 - 返回HTTP $ADMIN_HTTP_CODE"
        echo "问题根因: 管理员API路由问题"
        echo "解决方案: 检查路由配置"
    fi
else
    echo "❓ 服务返回未知状态 - HTTP $ROOT_HTTP_CODE"
    echo "问题根因: 未知错误"
    echo "解决方案: 检查云托管配置"
fi

# 6. 提供建议
echo ""
echo "=== 建议 ==="

if [ "$ROOT_HTTP_CODE" = "418" ]; then
    echo "1. 立即检查云托管控制台的服务日志"
    echo "2. 查看服务启动错误信息"
    echo "3. 检查数据库连接配置"
    echo "4. 修复启动问题后重新部署"
elif [ "$ROOT_HTTP_CODE" = "200" ] && [ "$ADMIN_HTTP_CODE" = "200" ]; then
    echo "1. 服务正常，问题在小程序调用"
    echo "2. 检查小程序云托管配置"
    echo "3. 验证小程序调用方式"
    echo "4. 测试小程序功能"
else
    echo "1. 需要进一步诊断"
    echo "2. 检查云托管控制台"
    echo "3. 查看服务日志"
    echo "4. 联系技术支持"
fi

# 清理临时文件
rm -f /tmp/root_response.txt /tmp/health_response.txt /tmp/admin_response.txt /tmp/other_response.txt

echo ""
echo "=== 测试完成 ==="
