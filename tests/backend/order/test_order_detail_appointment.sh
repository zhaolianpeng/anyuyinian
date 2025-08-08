#!/bin/bash

echo "=== 订单详情预约信息测试 ==="

# 设置基础URL
BASE_URL="http://localhost:8080"

echo "1. 测试订单详情接口"
echo "请求参数: orderNo=ORDER20241201000001"

# 发送POST请求获取订单详情
response=$(curl -s -X POST "${BASE_URL}/api/order/detail" \
  -H "Content-Type: application/json" \
  -d '{"orderNo":"ORDER20241201000001"}')

echo "响应状态码: $?"
echo "响应数据:"
echo "$response" | jq '.'

echo -e "\n2. 检查预约信息字段"
order_data=$(echo "$response" | jq '.data // empty')

if [ "$order_data" != "null" ] && [ -n "$order_data" ]; then
  echo "✅ 订单数据存在"
  
  echo -e "\n3. 检查基本信息字段"
  appointment_date=$(echo "$order_data" | jq -r '.appointmentDate // empty')
  appointment_time=$(echo "$order_data" | jq -r '.appointmentTime // empty')
  quantity=$(echo "$order_data" | jq -r '.quantity // empty')
  
  echo "预约日期: '$appointment_date'"
  echo "预约时间: '$appointment_time'"
  echo "预约数量: '$quantity'"
  
  echo -e "\n4. 检查患者信息字段"
  patient_name=$(echo "$order_data" | jq -r '.patientName // empty')
  patient_phone=$(echo "$order_data" | jq -r '.patientPhone // empty')
  
  echo "患者姓名: '$patient_name'"
  echo "患者电话: '$patient_phone'"
  
  echo -e "\n5. 检查地址信息字段"
  address_info=$(echo "$order_data" | jq -r '.addressInfo // empty')
  
  echo "地址信息: '$address_info'"
  
  echo -e "\n6. 检查病史信息字段"
  disease_info=$(echo "$order_data" | jq -r '.diseaseInfo // empty')
  need_toilet_assist=$(echo "$order_data" | jq -r '.needToiletAssist // empty')
  
  echo "既往病史: '$disease_info'"
  echo "助排二便: '$need_toilet_assist'"
  
  echo -e "\n7. 检查备注信息字段"
  remark=$(echo "$order_data" | jq -r '.remark // empty')
  
  echo "备注信息: '$remark'"
  
  echo -e "\n8. 检查表单数据字段"
  form_data=$(echo "$order_data" | jq -r '.formData // empty')
  
  echo "表单数据: '$form_data'"
  
  if [ -n "$form_data" ] && [ "$form_data" != "null" ]; then
    echo "解析表单数据:"
    echo "$form_data" | jq '.'
  fi
  
  echo -e "\n9. 检查价格信息字段"
  price=$(echo "$order_data" | jq -r '.price // empty')
  total_amount=$(echo "$order_data" | jq -r '.totalAmount // empty')
  formatted_price=$(echo "$order_data" | jq -r '.formattedPrice // empty')
  
  echo "价格: '$price'"
  echo "总金额: '$total_amount'"
  echo "格式化价格: '$formatted_price'"
  
else
  echo "❌ 订单数据不存在或为空"
fi

echo -e "\n=== 测试完成 ===" 