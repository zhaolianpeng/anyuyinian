#!/bin/bash

# 用户信息接口测试脚本

echo "=== 用户信息接口测试 ==="

# 设置测试参数
BASE_URL="http://localhost:80"
USER_INFO_URL="$BASE_URL/api/user/info"

echo "测试URL: $USER_INFO_URL"
echo ""

# 测试模拟用户ID
echo "测试模拟用户ID (userId=1)..."
RESPONSE=$(curl -s -X GET "$USER_INFO_URL?userId=1")

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
      echo "✅ 获取用户信息成功"
    else
      echo "❌ 获取用户信息失败"
    fi
    
    if echo "$RESPONSE" | grep -q '"nickName"'; then
      echo "✅ 包含用户昵称"
    else
      echo "❌ 缺少用户昵称"
    fi
    
    if echo "$RESPONSE" | grep -q '"avatarUrl"'; then
      echo "✅ 包含头像URL"
    else
      echo "❌ 缺少头像URL"
    fi
  fi
else
  echo "❌ 请求发送失败"
fi

echo ""
echo "=== 测试完成 ==="
echo ""
echo "如果测试成功，说明用户信息接口正常工作。"
echo "如果仍有问题，请检查后端服务是否正常运行。" 