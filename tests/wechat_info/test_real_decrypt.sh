#!/bin/bash

# 测试真正的微信手机号解密功能
# 使用用户提供的真实加密数据

echo "=== 测试真正的微信手机号解密 ==="

# 配置
API_BASE_URL="http://localhost:80"
USER_ID="507f1f77bcf86cd799439011"
# 用户提供的真实加密数据
ENCRYPTED_DATA="wtiE31rLGUTJUyJxC3Gx+ML0CMxqv3qjc7OyYh3TWy7bH0om/0BBTGAXTGYAxJhULQ6dQjoSerHvv++DTgbJ9fpgLw1MTLsj5g3geV3EuARPqRfiTHkJZrf4rHtqXVW4agj8HaxZ9KRXIX8Os2/QXG55LXsgf0hk18EgbLWQVWysx3xsvnLNTKRKOO+YdV4XE8g09XWrivDGr0vqRGUjog=="
IV="a11z+4Y6k4jYq9SpzGNWHA=="

echo "测试参数:"
echo "  API地址: $API_BASE_URL"
echo "  用户ID: $USER_ID"
echo "  加密数据: ${ENCRYPTED_DATA:0:30}..."
echo "  IV: $IV"
echo "  期望手机号: 13691028481"
echo ""

# 测试解密手机号API
echo "1. 测试真正的微信手机号解密..."
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
ERROR_MSG=$(echo "$RESPONSE" | jq -r '.errorMsg // "null"')

echo "解析结果:"
echo "  状态码: $CODE"
echo "  手机号: $PHONE_NUMBER"
echo "  错误信息: $ERROR_MSG"
echo ""

# 验证结果
if [ "$CODE" = "0" ] && [ "$PHONE_NUMBER" != "null" ] && [ "$PHONE_NUMBER" != "" ]; then
    echo "✅ 手机号解密成功!"
    echo "   获取到的手机号: $PHONE_NUMBER"
    
    if [ "$PHONE_NUMBER" = "13691028481" ]; then
        echo "🎉 解密结果正确! 与用户真实手机号匹配!"
    else
        echo "⚠️  解密结果与期望不符:"
        echo "   期望: 13691028481"
        echo "   实际: $PHONE_NUMBER"
    fi
else
    echo "❌ 手机号解密失败!"
    echo "   错误信息: $ERROR_MSG"
    echo "   完整响应: $RESPONSE"
fi

echo ""
echo "=== 测试完成 ==="
