#!/bin/bash

echo "=== 订单价格字段调试测试 ==="

# 设置基础URL
BASE_URL="http://localhost:8080"

echo "1. 测试订单列表接口"
echo "请求参数: userId=1&page=1&pageSize=5"

response=$(curl -s -X GET "${BASE_URL}/api/order/list?userId=1&page=1&pageSize=5")

echo "响应状态码: $?"
echo "响应数据:"
echo "$response" | jq '.'

echo -e "\n2. 检查订单数据结构"
orders=$(echo "$response" | jq '.data.list // []')

if [ "$(echo "$orders" | jq 'length')" -gt 0 ]; then
  echo "找到订单数量: $(echo "$orders" | jq 'length')"
  
  echo -e "\n3. 详细检查第一个订单的价格字段"
  first_order=$(echo "$orders" | jq '.[0]')
  echo "第一个订单完整数据:"
  echo "$first_order" | jq '.'
  
  echo -e "\n4. 检查价格相关字段"
  price=$(echo "$first_order" | jq -r '.price // empty')
  total_amount=$(echo "$first_order" | jq -r '.totalAmount // empty')
  formatted_amount=$(echo "$first_order" | jq -r '.formattedAmount // empty')
  
  echo "price字段: '$price'"
  echo "totalAmount字段: '$total_amount'"
  echo "formattedAmount字段: '$formatted_amount'"
  
  if [ -n "$price" ] && [ "$price" != "null" ]; then
    echo "✅ price字段存在且不为空"
  else
    echo "❌ price字段为空或不存在"
  fi
  
  if [ -n "$total_amount" ] && [ "$total_amount" != "null" ]; then
    echo "✅ totalAmount字段存在且不为空"
  else
    echo "❌ totalAmount字段为空或不存在"
  fi
  
else
  echo "❌ 没有找到订单数据"
fi

echo -e "\n=== 测试完成 ===" 