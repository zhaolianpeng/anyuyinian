#!/bin/bash

echo "=== 订单详情总计价格测试 ==="

# 设置基础URL
BASE_URL="http://localhost:8080"

echo "1. 测试订单详情接口"
echo "请求参数: orderNo=ORDER202581373396"

# 发送POST请求获取订单详情
response=$(curl -s -X POST "${BASE_URL}/api/order/detail" \
  -H "Content-Type: application/json" \
  -d '{"orderNo":"ORDER202581373396"}')

echo "响应状态码: $?"
echo "响应数据:"
echo "$response" | jq '.'

echo -e "\n2. 检查价格相关字段"
order_data=$(echo "$response" | jq '.data // empty')

if [ "$order_data" != "null" ] && [ -n "$order_data" ]; then
  echo "✅ 订单数据存在"
  
  echo -e "\n3. 检查价格计算"
  price=$(echo "$order_data" | jq -r '.price // empty')
  quantity=$(echo "$order_data" | jq -r '.quantity // empty')
  total_amount=$(echo "$order_data" | jq -r '.totalAmount // empty')
  
  echo "服务单价: '$price'"
  echo "数量: '$quantity'"
  echo "总金额: '$total_amount'"
  
  # 计算正确的总金额
  if [ -n "$price" ] && [ -n "$quantity" ] && [ "$price" != "null" ] && [ "$quantity" != "null" ]; then
    calculated_total=$(echo "$price * $quantity" | bc -l)
    echo "计算的总金额: $calculated_total"
    
    if [ "$(echo "$total_amount == $calculated_total" | bc -l)" -eq 1 ]; then
      echo "✅ 总金额计算正确"
    else
      echo "❌ 总金额计算错误"
      echo "   - 数据库中的总金额: $total_amount"
      echo "   - 计算出的总金额: $calculated_total"
    fi
  else
    echo "❌ 价格或数量字段为空"
  fi
  
  echo -e "\n4. 检查其他价格字段"
  formatted_price=$(echo "$order_data" | jq -r '.formattedPrice // empty')
  
  echo "格式化价格: '$formatted_price'"
  
  echo -e "\n5. 验证价格字段的数据类型"
  echo "价格类型: $(echo "$order_data" | jq -r '.price | type')"
  echo "数量类型: $(echo "$order_data" | jq -r '.quantity | type')"
  echo "总金额类型: $(echo "$order_data" | jq -r '.totalAmount | type')"
  
else
  echo "❌ 订单数据不存在或为空"
fi

echo -e "\n=== 测试完成 ===" 