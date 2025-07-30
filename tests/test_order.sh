#!/bin/bash

# 订单相关接口测试脚本

# 设置测试服务器地址
SERVER_URL="http://localhost:80"

echo "测试订单相关接口..."
echo "服务器地址: $SERVER_URL"
echo ""

# 测试1: 提交订单
echo "=== 测试1: 提交订单 ==="
curl -X POST \
  -H "Content-Type: application/json" \
  -d "{
    \"userId\": 1,
    \"serviceId\": 1,
    \"patientId\": 1,
    \"addressId\": 1,
    \"appointmentDate\": \"2024-01-15\",
    \"appointmentTime\": \"morning\",
    \"specialRequirements\": \"无特殊要求\",
    \"formData\": {
      \"patientName\": \"张三\",
      \"patientPhone\": \"13800138000\",
      \"appointmentDate\": \"2024-01-15\",
      \"appointmentTime\": \"morning\",
      \"specialRequirements\": \"无特殊要求\"
    }
  }" \
  "$SERVER_URL/api/order/submit"
echo ""
echo ""

# 测试2: 发起支付
echo "=== 测试2: 发起支付 ==="
curl -X POST \
  -H "Content-Type: application/json" \
  -d "{
    \"paymentMethod\": \"wechat_pay\",
    \"openId\": \"test_openid_123456\"
  }" \
  "$SERVER_URL/api/order/pay/1"
echo ""
echo ""

# 测试3: 获取订单列表
echo "=== 测试3: 获取订单列表 ==="
curl -X GET "$SERVER_URL/api/order/list?userId=1&page=1&pageSize=10"
echo ""
echo ""

# 测试4: 获取订单详情
echo "=== 测试4: 获取订单详情 ==="
curl -X GET "$SERVER_URL/api/order/detail/1"
echo ""
echo ""

# 测试5: 取消订单
echo "=== 测试5: 取消订单 ==="
curl -X POST \
  -H "Content-Type: application/json" \
  -d "{
    \"reason\": \"个人原因取消\"
  }" \
  "$SERVER_URL/api/order/cancel/1"
echo ""
echo ""

# 测试6: 申请退款
echo "=== 测试6: 申请退款 ==="
curl -X POST \
  -H "Content-Type: application/json" \
  -d "{
    \"reason\": \"服务不满意\",
    \"refundAmount\": 299.00
  }" \
  "$SERVER_URL/api/order/refund/1"
echo ""
echo ""

echo "测试完成！"
echo "注意：确保数据库已创建Orders表" 