#!/bin/bash

# 测试微信手机号解密API
# 用于验证完善资料页面的手机号获取功能

echo "=== 测试微信手机号解密API ==="

# 配置
API_BASE_URL="http://localhost:80"
USER_ID="507f1f77bcf86cd799439011"  # 使用实际的用户ID
ENCRYPTED_DATA="0CaCNQSekx3+3I82M0cbiyY6vr6PizU9UAl9YtB6wYVKjQIsq4"
IV="FFBz015C3VHz4ku/4QEHWQ=="

echo "测试参数:"
echo "  API地址: $API_BASE_URL"
echo "  用户ID: $USER_ID"
echo "  加密数据: ${ENCRYPTED_DATA:0:20}..."
echo "  IV: $IV"
echo ""

# 测试解密手机号API
echo "1. 测试解密手机号API..."
RESPONSE=$(curl -s -X POST "$API_BASE_URL/api/user/decrypt_phone" \
  -H "Content-Type: application/json" \
  -d "{
    \"userId\": \"$USER_ID\",
    \"encryptedData\": \"$ENCRYPTED_DATA\",
    \"iv\": \"$IV\"
  }")

echo "响应: $RESPONSE"
echo ""

# 解析响应
PHONE_NUMBER=$(echo "$RESPONSE" | jq -r '.data.phoneNumber // "null"')
CODE=$(echo "$RESPONSE" | jq -r '.code // "null"')

echo "解析结果:"
echo "  状态码: $CODE"
echo "  手机号: $PHONE_NUMBER"
echo ""

# 验证结果
if [ "$CODE" = "0" ] && [ "$PHONE_NUMBER" != "null" ] && [ "$PHONE_NUMBER" != "" ]; then
    echo "✅ 手机号解密成功!"
    echo "   获取到的手机号: $PHONE_NUMBER"
else
    echo "❌ 手机号解密失败!"
    echo "   错误信息: $RESPONSE"
fi

echo ""
echo "=== 测试完成 ==="
