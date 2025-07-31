#!/bin/bash

# 模拟登录测试脚本

echo "=== 模拟登录测试 ==="

# 设置测试参数
BASE_URL="http://localhost:80"
LOGIN_URL="$BASE_URL/api/wx/login"

# 模拟登录请求数据
LOGIN_DATA='{
  "code": "mock_code_123",
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
echo "发送模拟登录请求..."
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
    
    # 检查是否包含必要的字段
    if echo "$RESPONSE" | grep -q '"code":0'; then
      echo "✅ 登录成功"
    else
      echo "❌ 登录失败"
    fi
    
    if echo "$RESPONSE" | grep -q '"token"'; then
      echo "✅ 包含token"
    else
      echo "❌ 缺少token"
    fi
    
    if echo "$RESPONSE" | grep -q '"userId"'; then
      echo "✅ 包含userId"
    else
      echo "❌ 缺少userId"
    fi
  fi
else
  echo "❌ 请求发送失败"
fi

echo ""
echo "=== 测试完成 ==="
echo ""
echo "如果测试成功，说明模拟登录功能正常工作。"
echo "如果仍有问题，请检查后端服务是否正常运行。" 