#!/bin/bash

echo "=== 测试订单倒计时功能 ==="

# 设置基础URL
BASE_URL="http://localhost:8080"

echo "1. 测试订单列表接口，检查待支付订单"
echo "请求参数: userId=1&page=1&pageSize=10"

response=$(curl -s -X GET "${BASE_URL}/api/order/list?userId=1&page=1&pageSize=10")

echo "响应数据:"
echo "$response" | jq '.'

echo -e "\n2. 检查待支付订单"
pending_orders=$(echo "$response" | jq '.data.list[] | select(.status == 0)')

if [ -n "$pending_orders" ]; then
  echo "找到待支付订单:"
  echo "$pending_orders" | jq '.'
  
  echo -e "\n3. 检查订单金额字段"
  total_amount=$(echo "$pending_orders" | jq -r '.totalAmount // empty')
  formatted_amount=$(echo "$pending_orders" | jq -r '.formattedAmount // empty')
  
  echo "订单金额: $total_amount"
  echo "格式化金额: $formatted_amount"
  
  if [ -n "$total_amount" ] && [ "$total_amount" != "0" ]; then
    echo "✅ 订单金额字段正常"
  else
    echo "❌ 订单金额字段异常"
  fi
else
  echo "❌ 没有找到待支付订单"
fi

echo -e "\n4. 测试订单详情接口"
echo "使用第一个订单的订单号测试详情接口"

first_order_no=$(echo "$response" | jq -r '.data.list[0].orderNo // empty')

if [ -n "$first_order_no" ]; then
  echo "订单号: $first_order_no"
  
  detail_response=$(curl -s -X POST "${BASE_URL}/api/order/detail" \
    -H "Content-Type: application/json" \
    -d "{\"orderNo\": \"$first_order_no\"}")
  
  echo "订单详情响应:"
  echo "$detail_response" | jq '.'
  
  order_status=$(echo "$detail_response" | jq -r '.data.status // empty')
  created_at=$(echo "$detail_response" | jq -r '.data.createdAt // empty')
  
  echo "订单状态: $order_status"
  echo "创建时间: $created_at"
  
  if [ "$order_status" = "0" ]; then
    echo "✅ 这是待支付订单，应该显示倒计时"
  else
    echo "ℹ️ 这不是待支付订单"
  fi
else
  echo "❌ 没有找到订单号"
fi

echo -e "\n=== 测试完成 ===" 