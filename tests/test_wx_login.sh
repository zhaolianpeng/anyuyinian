#!/bin/bash

# 微信登录接口测试脚本

# 设置测试服务器地址
SERVER_URL="http://localhost:80"

echo "测试微信登录接口..."
echo "服务器地址: $SERVER_URL"
echo ""

# 测试微信登录
echo "=== 测试微信登录 ==="
curl -X POST \
  -H "Content-Type: application/json" \
  -d "{
    \"code\": \"test_code_123456\",
    \"userInfo\": {
      \"nickName\": \"测试用户\",
      \"avatarUrl\": \"https://example.com/avatar.jpg\",
      \"gender\": 1,
      \"country\": \"中国\",
      \"province\": \"广东省\",
      \"city\": \"深圳市\",
      \"language\": \"zh_CN\"
    }
  }" \
  "$SERVER_URL/api/wx/login"
echo ""
echo ""

echo "测试完成！"
echo "注意：测试脚本中的code是模拟的，实际使用时需要从小程序获取真实的code" 