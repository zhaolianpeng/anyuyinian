#!/bin/bash

# 测试订单列表状态筛选功能
echo "=== 测试订单列表状态筛选功能 ==="

# 基础URL
BASE_URL="http://localhost:80"

# 测试用户ID
USER_ID="1"

echo ""
echo "1. 测试获取全部订单列表"
curl -X GET "${BASE_URL}/api/order/list?userId=${USER_ID}&page=1&pageSize=10" \
  -H "Content-Type: application/json" | jq '.'

echo ""
echo "2. 测试获取待支付订单列表"
curl -X GET "${BASE_URL}/api/order/list?userId=${USER_ID}&status=pending_pay&page=1&pageSize=10" \
  -H "Content-Type: application/json" | jq '.'

echo ""
echo "3. 测试获取已支付订单列表"
curl -X GET "${BASE_URL}/api/order/list?userId=${USER_ID}&status=paid&page=1&pageSize=10" \
  -H "Content-Type: application/json" | jq '.'

echo ""
echo "4. 测试获取已取消订单列表"
curl -X GET "${BASE_URL}/api/order/list?userId=${USER_ID}&status=cancelled&page=1&pageSize=10" \
  -H "Content-Type: application/json" | jq '.'

echo ""
echo "5. 测试获取已退款订单列表"
curl -X GET "${BASE_URL}/api/order/list?userId=${USER_ID}&status=refunded&page=1&pageSize=10" \
  -H "Content-Type: application/json" | jq '.'

echo ""
echo "6. 测试无效状态参数"
curl -X GET "${BASE_URL}/api/order/list?userId=${USER_ID}&status=invalid_status&page=1&pageSize=10" \
  -H "Content-Type: application/json" | jq '.'

echo ""
echo "=== 测试完成 ===" 