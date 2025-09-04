#!/bin/bash

# 微信支付测试脚本
# 商户号: 1726638701

echo "=== 微信支付功能测试 ==="
echo "商户号: 1726638701"
echo ""

# 配置
API_BASE_URL="http://localhost:80"
USER_ID="507f1f77bcf86cd799439011"
SERVICE_ID="1"
OPEN_ID="test_openid_123456"

echo "测试参数:"
echo "  API地址: $API_BASE_URL"
echo "  用户ID: $USER_ID"
echo "  服务ID: $SERVICE_ID"
echo "  OpenID: $OPEN_ID"
echo ""

# 1. 测试创建订单
echo "1. 测试创建订单..."
ORDER_RESPONSE=$(curl -s -X POST "$API_BASE_URL/api/order/submit" \
  -H "Content-Type: application/json" \
  -d "{
    \"serviceId\": \"$SERVICE_ID\",
    \"userId\": \"$USER_ID\",
    \"patientId\": \"1\",
    \"addressId\": \"1\",
    \"appointmentDate\": \"2025-01-03\",
    \"appointmentTime\": \"09:00\",
    \"remark\": \"测试订单\"
  }")

echo "订单创建响应: $ORDER_RESPONSE"
echo ""

# 解析订单ID
ORDER_ID=$(echo "$ORDER_RESPONSE" | jq -r '.data.orderId // "null"')
if [ "$ORDER_ID" = "null" ] || [ -z "$ORDER_ID" ]; then
    echo "❌ 订单创建失败，无法继续测试支付"
    exit 1
fi

echo "✅ 订单创建成功，订单ID: $ORDER_ID"
echo ""

# 2. 测试支付订单
echo "2. 测试支付订单..."
PAYMENT_RESPONSE=$(curl -s -X POST "$API_BASE_URL/api/order/pay" \
  -H "Content-Type: application/json" \
  -d "{
    \"orderId\": \"$ORDER_ID\",
    \"paymentMethod\": \"wechat_pay\",
    \"openId\": \"$OPEN_ID\"
  }")

echo "支付响应: $PAYMENT_RESPONSE"
echo ""

# 解析支付参数
PAYMENT_CODE=$(echo "$PAYMENT_RESPONSE" | jq -r '.code // "null"')
if [ "$PAYMENT_CODE" = "0" ]; then
    echo "✅ 支付参数获取成功"
    
    # 显示支付参数
    echo "支付参数详情:"
    echo "$PAYMENT_RESPONSE" | jq '.data.paymentParams'
    
    # 验证支付参数完整性
    TIMESTAMP=$(echo "$PAYMENT_RESPONSE" | jq -r '.data.paymentParams.timeStamp // "null"')
    NONCE_STR=$(echo "$PAYMENT_RESPONSE" | jq -r '.data.paymentParams.nonceStr // "null"')
    PACKAGE=$(echo "$PAYMENT_RESPONSE" | jq -r '.data.paymentParams.package // "null"')
    PAY_SIGN=$(echo "$PAYMENT_RESPONSE" | jq -r '.data.paymentParams.paySign // "null"')
    
    if [ "$TIMESTAMP" != "null" ] && [ "$NONCE_STR" != "null" ] && [ "$PACKAGE" != "null" ] && [ "$PAY_SIGN" != "null" ]; then
        echo "✅ 支付参数完整，可以调起微信支付"
    else
        echo "❌ 支付参数不完整"
    fi
else
    echo "❌ 支付参数获取失败"
    ERROR_MSG=$(echo "$PAYMENT_RESPONSE" | jq -r '.message // "未知错误"')
    echo "错误信息: $ERROR_MSG"
fi

echo ""
echo "=== 测试完成 ==="
echo ""
echo "注意事项:"
echo "1. 确保已设置正确的商户密钥环境变量"
echo "2. 确保商户号1726638701已正确配置"
echo "3. 确保通知URL已配置为您的实际域名"
echo "4. 确保微信小程序已关联该商户号"
