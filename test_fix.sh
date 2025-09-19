#!/bin/bash

# 测试修复后的API

echo "=== 测试修复后的API ==="

SERVICE_URL="https://golang-lfwy-prod-5g94mx7a3d07e78c.service.tcloudbaseapp.com"

# 等待部署完成
echo "等待部署完成..."
sleep 10

# 1. 测试管理员服务列表
echo "1. 测试管理员服务列表..."
ADMIN_SERVICES_RESPONSE=$(curl -k -s -w "HTTP_CODE:%{http_code}" -o /tmp/admin_services_response.txt -X GET "$SERVICE_URL/api/admin/services?page=1&pageSize=10&adminUserId=admin_super" 2>/dev/null)
ADMIN_SERVICES_HTTP_CODE=$(echo "$ADMIN_SERVICES_RESPONSE" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)

echo "管理员服务列表响应: HTTP $ADMIN_SERVICES_HTTP_CODE"

# 2. 测试服务价格更新
echo "2. 测试服务价格更新..."
UPDATE_PRICE_RESPONSE=$(curl -k -s -w "HTTP_CODE:%{http_code}" -o /tmp/update_price_response.txt -X POST "$SERVICE_URL/api/admin/service/updateprice" \
  -H "Content-Type: application/json" \
  -H "X-WX-SERVICE: golang-lfwy" \
  -d '{"serviceId": 21, "newPrice": 280, "newOriginalPrice": 200, "reason": "测试修复"}' 2>/dev/null)
UPDATE_PRICE_HTTP_CODE=$(echo "$UPDATE_PRICE_RESPONSE" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
UPDATE_PRICE_CONTENT=$(cat /tmp/update_price_response.txt 2>/dev/null)

echo "服务价格更新响应: HTTP $UPDATE_PRICE_HTTP_CODE"
if [ -n "$UPDATE_PRICE_CONTENT" ]; then
    echo "响应内容: $UPDATE_PRICE_CONTENT"
fi

# 3. 分析结果
echo ""
echo "=== 分析结果 ==="

if [ "$ADMIN_SERVICES_HTTP_CODE" = "200" ] && [ "$UPDATE_PRICE_HTTP_CODE" = "200" ]; then
    echo "✅ 修复成功！两个API都正常工作"
    echo "问题根因: 路由冲突已解决"
    echo "解决方案: 调整路由注册顺序，将具体路由放在通用路由之前"
elif [ "$ADMIN_SERVICES_HTTP_CODE" = "200" ] && [ "$UPDATE_PRICE_HTTP_CODE" = "404" ]; then
    echo "❌ 部分修复 - 管理员服务列表正常，但价格更新仍返回404"
    echo "问题根因: 路由冲突未完全解决"
    echo "解决方案: 需要进一步检查路由配置"
elif [ "$ADMIN_SERVICES_HTTP_CODE" = "200" ] && [ "$UPDATE_PRICE_HTTP_CODE" = "418" ]; then
    echo "❌ 服务问题 - 管理员服务列表正常，但价格更新返回418"
    echo "问题根因: 服务启动问题"
    echo "解决方案: 检查服务日志"
else
    echo "❓ 未知状态 - 需要进一步诊断"
    echo "管理员服务列表: HTTP $ADMIN_SERVICES_HTTP_CODE"
    echo "价格更新: HTTP $UPDATE_PRICE_HTTP_CODE"
fi

# 清理临时文件
rm -f /tmp/admin_services_response.txt /tmp/update_price_response.txt

echo ""
echo "=== 测试完成 ==="
