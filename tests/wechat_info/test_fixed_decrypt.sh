#!/bin/bash

# 测试修复后的微信手机号解密功能
# 使用用户提供的最新加密数据

echo "=== 测试修复后的微信手机号解密 ==="

# 配置
API_BASE_URL="http://localhost:80"
USER_ID="507f1f77bcf86cd799439011"
# 用户提供的最新加密数据
ENCRYPTED_DATA="9sMW1YO/iA5SflKcNUl7RONDXBJdxGdbpKwZiNMJai8pQtDMsM37v0VmVbx5/LtPvz8izDZ0DYhYf2r4IdA/dV0Cd/Nj7Ctc/WoAkT7kWdRw8z9ydV0Efpf5LJnNiz2USXpIsUQRecmERvbyeiVllbie+6MeuIhKws24F8KpyldGFyNtS6Q/1IgMuCsN7hYRxAmj4//XQHjHYUhwceXuvw=="
IV="w9hfnau2fnxnwCFuk8Cl3g=="

echo "测试参数:"
echo "  API地址: $API_BASE_URL"
echo "  用户ID: $USER_ID"
echo "  加密数据: ${ENCRYPTED_DATA:0:30}..."
echo "  IV: $IV"
echo "  期望手机号: 13691028481"
echo "  使用密钥: session_key (Base64解码后)"
echo ""

# 测试解密手机号API
echo "1. 测试修复后的微信手机号解密..."
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
    elif [[ "$ERROR_MSG" == *"解码session_key失败"* ]]; then
        echo "💡 提示: session_key格式不正确，需要Base64解码"
    elif [[ "$ERROR_MSG" == *"加密数据长度"* ]]; then
        echo "💡 提示: 加密数据长度不是AES块大小的倍数"
    fi
    
    echo "   完整响应: $RESPONSE"
fi

echo ""
echo "=== 测试完成 ==="
