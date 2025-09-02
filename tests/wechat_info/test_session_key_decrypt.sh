#!/bin/bash

# 测试使用session_key的微信手机号解密功能
# 使用用户提供的真实加密数据

echo "=== 测试使用session_key的微信手机号解密 ==="

# 配置
API_BASE_URL="http://localhost:80"
USER_ID="507f1f77bcf86cd799439011"
# 用户提供的真实加密数据
ENCRYPTED_DATA="pfB6t6suxcI+d8LLv57u2TnRmkkccNTYwUIpoq9QOAdw2DTlB832otL9j8oZ4OecpKcT+c4+j44/LwNhp+zF/aPLXYoEuDhZ6QAc/cr8bIgPZjNyLm2dbYPDv0cmQUnN2S881pjQwKsbpklvOaNXDnzWwA9qAK84aPKHzhuBltGVHfUu0rH3+iwRYDgxQI9zbLogOX3+C2vslIc5a+lwwg=="
IV="iG13YFHeizSoY3y2kFPo9Q=="

echo "测试参数:"
echo "  API地址: $API_BASE_URL"
echo "  用户ID: $USER_ID"
echo "  加密数据: ${ENCRYPTED_DATA:0:30}..."
echo "  IV: $IV"
echo "  期望手机号: 13691028481"
echo "  使用密钥: session_key (从用户记录中获取)"
echo ""

# 测试解密手机号API
echo "1. 测试使用session_key的微信手机号解密..."
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
        echo "✅ 使用session_key解密成功!"
    else
        echo "⚠️  解密结果与期望不符:"
        echo "   期望: 13691028481"
        echo "   实际: $PHONE_NUMBER"
    fi
else
    echo "❌ 手机号解密失败!"
    echo "   错误信息: $ERROR_MSG"
    
    # 分析错误类型
    if [[ "$ERROR_MSG" == *"session_key"* ]]; then
        echo "💡 提示: 可能是session_key问题，请检查用户是否已登录"
    elif [[ "$ERROR_MSG" == *"解析解密数据失败"* ]]; then
        echo "💡 提示: 解密算法可能仍有问题，需要进一步调试"
    elif [[ "$ERROR_MSG" == *"用户不存在"* ]]; then
        echo "💡 提示: 用户ID不存在，请检查用户是否已注册"
    fi
    
    echo "   完整响应: $RESPONSE"
fi

echo ""
echo "=== 测试完成 ==="
