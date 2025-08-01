#!/bin/bash

echo "=== 订单列表显示测试 ==="

# 设置基础URL
BASE_URL="http://localhost:8080"

echo "1. 测试订单列表接口"
echo "请求参数: userId=1&page=1&pageSize=5"

response=$(curl -s -X GET "${BASE_URL}/api/order/list?userId=1&page=1&pageSize=5")

echo "响应状态码: $?"
echo "响应数据:"
echo "$response" | jq '.'

echo -e "\n2. 检查订单金额字段"
orders=$(echo "$response" | jq '.data.list // []')

if [ "$(echo "$orders" | jq 'length')" -gt 0 ]; then
  echo "找到订单数量: $(echo "$orders" | jq 'length')"
  
  echo -e "\n3. 详细检查每个订单的金额字段"
  echo "$orders" | jq '.[] | {id, orderNo, serviceName, price, totalAmount, formattedAmount}'
  
  echo -e "\n4. 检查金额字段的值"
  echo "$orders" | jq '.[] | "订单ID: \(.id), 订单号: \(.orderNo), 价格: \(.price), 总金额: \(.totalAmount), 格式化金额: \(.formattedAmount)"'
  
  echo -e "\n5. 验证金额字段是否正确"
  price_check=$(echo "$orders" | jq '.[] | select(.price > 0) | .id')
  if [ -n "$price_check" ]; then
    echo "✅ 找到有价格的订单: $price_check"
  else
    echo "❌ 没有找到有价格的订单"
  fi
  
  total_amount_check=$(echo "$orders" | jq '.[] | select(.totalAmount > 0) | .id')
  if [ -n "$total_amount_check" ]; then
    echo "✅ 找到有总金额的订单: $total_amount_check"
  else
    echo "❌ 没有找到有总金额的订单"
  fi
  
  echo -e "\n6. 检查金额字段的数据类型"
  echo "$orders" | jq '.[] | "订单ID: \(.id), 价格类型: \(.price | type), 价格值: \(.price), 总金额类型: \(.totalAmount | type), 总金额值: \(.totalAmount)"'
  
  echo -e "\n7. 模拟前端formatAmount处理"
  echo "$orders" | jq -r '.[] | "订单ID: \(.id), 总金额: \(.totalAmount), 格式化后: \(.totalAmount | tostring | . + ".00" | .[0:length-3] + "." + .[length-2:])"'
  
else
  echo "❌ 没有找到订单数据"
fi

echo -e "\n=== 测试完成 ===" 