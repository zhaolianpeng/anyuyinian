#!/bin/bash

echo "=== 测试订单列表接口数据问题 ==="

# 设置基础URL
BASE_URL="http://localhost:8080"

echo "1. 测试订单列表接口"
echo "请求参数: userId=1&page=1&pageSize=5"

response=$(curl -s -X GET "${BASE_URL}/api/order/list?userId=1&page=1&pageSize=5")

echo "响应数据:"
echo "$response" | jq '.'

echo -e "\n2. 检查订单金额字段"
orders=$(echo "$response" | jq '.data.list // []')

if [ "$orders" != "[]" ]; then
  echo "订单数量: $(echo "$orders" | jq 'length')"
  
  echo -e "\n3. 详细检查每个订单的金额字段"
  echo "$orders" | jq '.[] | {id, orderNo, totalAmount, price, serviceName}'
  
  echo -e "\n4. 检查时间字段"
  echo "$orders" | jq '.[] | {id, orderNo, createdAt, formattedDate}'
  
  echo -e "\n5. 检查是否有金额为0的订单"
  zero_amount_orders=$(echo "$orders" | jq '[.[] | select(.totalAmount == 0 or .totalAmount == null)]')
  
  if [ "$(echo "$zero_amount_orders" | jq 'length')" -gt 0 ]; then
    echo "❌ 发现金额为0的订单:"
    echo "$zero_amount_orders" | jq '.[] | {id, orderNo, totalAmount, serviceName}'
  else
    echo "✅ 所有订单金额都正常"
  fi
else
  echo "❌ 没有找到订单数据"
fi

echo -e "\n6. 测试订单详情接口对比"
first_order_no=$(echo "$response" | jq -r '.data.list[0].orderNo // empty')

if [ -n "$first_order_no" ]; then
  echo "使用第一个订单号测试详情接口: $first_order_no"
  
  detail_response=$(curl -s -X POST "${BASE_URL}/api/order/detail" \
    -H "Content-Type: application/json" \
    -d "{\"orderNo\": \"$first_order_no\"}")
  
  echo "订单详情响应:"
  echo "$detail_response" | jq '.data | {id, orderNo, totalAmount, price, serviceName}'
  
  list_total=$(echo "$response" | jq -r '.data.list[0].totalAmount // empty')
  list_price=$(echo "$response" | jq -r '.data.list[0].price // empty')
  detail_total=$(echo "$detail_response" | jq -r '.data.totalAmount // empty')
  detail_price=$(echo "$detail_response" | jq -r '.data.price // empty')
  
  echo "列表接口金额: $list_total"
  echo "列表接口单价: $list_price"
  echo "详情接口金额: $detail_total"
  echo "详情接口单价: $detail_price"
  
  if [ "$list_total" != "$detail_total" ]; then
    echo "❌ 列表和详情接口的金额不一致"
  else
    echo "✅ 列表和详情接口的金额一致"
  fi
  
  if [ "$list_price" != "$detail_price" ]; then
    echo "❌ 列表和详情接口的单价不一致"
  else
    echo "✅ 列表和详情接口的单价一致"
  fi
else
  echo "❌ 没有找到订单号"
fi

echo -e "\n=== 测试完成 ===" 