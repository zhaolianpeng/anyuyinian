#!/bin/bash

# 测试预约页面相关的后端API

BASE_URL="https://golang-lfwy-176496-6-1353115175.sh.run.tcloudbase.com"

echo "=== 测试预约页面后端API ==="

# 使用模拟的用户ID进行测试
USER_ID=1

echo "使用模拟用户ID: $USER_ID"

# 测试就诊人API
echo ""
echo "1. 测试就诊人API..."
PATIENT_RESPONSE=$(curl -s -X GET "$BASE_URL/api/user/patient?userId=$USER_ID")
echo "就诊人API响应: $PATIENT_RESPONSE"

# 测试地址API
echo ""
echo "2. 测试地址API..."
ADDRESS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/user/address?userId=$USER_ID")
echo "地址API响应: $ADDRESS_RESPONSE"

# 测试服务详情API
echo ""
echo "3. 测试服务详情API..."
SERVICE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/service/detail" \
  -H "Content-Type: application/json" \
  -d '{"serviceId": 1}')
echo "服务详情API响应: $SERVICE_RESPONSE"

# 测试订单提交API
echo ""
echo "4. 测试订单提交API..."
ORDER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/order/submit" \
  -H "Content-Type: application/json" \
  -d "{
    \"userId\": $USER_ID,
    \"serviceId\": 1,
    \"patientId\": 1,
    \"addressId\": 1,
    \"appointmentDate\": \"2024-12-20\",
    \"appointmentTime\": \"10:00\",
    \"remark\": \"测试订单\"
  }")
echo "订单提交API响应: $ORDER_RESPONSE"

echo ""
echo "=== 测试完成 ===" 