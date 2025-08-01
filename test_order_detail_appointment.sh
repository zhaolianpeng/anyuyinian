#!/bin/bash

echo "=== 测试订单详情页预约信息显示 ==="

# 设置基础URL
BASE_URL="http://localhost:8080"

echo "1. 测试订单详情接口，检查预约信息字段"
echo "请求参数: {\"orderNo\": \"ORDER20241201000001\"}"

response=$(curl -s -X POST "${BASE_URL}/api/order/detail" \
  -H "Content-Type: application/json" \
  -d '{"orderNo": "ORDER20241201000001"}')

echo "响应数据:"
echo "$response" | jq '.'

echo -e "\n2. 检查预约信息字段是否存在"
echo "$response" | jq '.data | {appointmentDate, appointmentTime, formData}'

echo -e "\n3. 验证预约信息字段"
appointment_date=$(echo "$response" | jq -r '.data.appointmentDate // empty')
appointment_time=$(echo "$response" | jq -r '.data.appointmentTime // empty')
form_data=$(echo "$response" | jq -r '.data.formData // empty')

echo "预约日期: $appointment_date"
echo "预约时间: $appointment_time"
echo "表单数据: $form_data"

if [ -n "$appointment_date" ]; then
  echo "✅ 预约日期字段存在"
else
  echo "❌ 预约日期字段为空"
fi

if [ -n "$appointment_time" ]; then
  echo "✅ 预约时间字段存在"
else
  echo "❌ 预约时间字段为空"
fi

if [ -n "$form_data" ]; then
  echo "✅ 表单数据字段存在"
  echo "表单数据内容:"
  echo "$form_data" | jq '.' 2>/dev/null || echo "表单数据不是有效的JSON格式"
else
  echo "❌ 表单数据字段为空"
fi

echo -e "\n=== 测试完成 ===" 