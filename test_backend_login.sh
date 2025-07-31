#!/bin/bash

# 后端登录测试脚本

echo "=== 后端登录测试 ==="

# 设置测试参数
BASE_URL="http://localhost:80"
LOGIN_URL="$BASE_URL/api/wx/login"

# 模拟登录请求数据
LOGIN_DATA='{
  "code": "test_code_123",
  "userInfo": {
    "nickName": "测试用户",
    "avatarUrl": "https://example.com/avatar.jpg",
    "gender": 0,
    "country": "China",
    "province": "Guangdong",
    "city": "Shenzhen",
    "language": "zh_CN"
  }
}'

echo "测试URL: $LOGIN_URL"
echo "请求数据: $LOGIN_DATA"
echo ""

# 发送登录请求
echo "发送登录请求..."
RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d "$LOGIN_DATA" \
  "$LOGIN_URL")

echo "响应状态码: $?"
echo "响应内容: $RESPONSE"
echo ""

# 解析响应
if [ $? -eq 0 ]; then
  echo "✅ 请求发送成功"
  
  # 检查响应是否包含错误信息
  if echo "$RESPONSE" | grep -q "error"; then
    echo "❌ 响应包含错误信息"
    echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
  else
    echo "✅ 响应格式正确"
    echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
  fi
else
  echo "❌ 请求发送失败"
fi

echo ""
echo "=== 测试完成 ==="
echo ""
echo "如果看到配置错误，请设置正确的微信AppID和AppSecret："
echo "export WX_APP_ID=\"你的真实AppID\""
echo "export WX_APP_SECRET=\"你的真实AppSecret\""
echo "然后重启后端服务：go run main.go" 