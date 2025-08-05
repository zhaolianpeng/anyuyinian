#!/bin/bash

# 用户信息相关接口测试脚本

# 设置测试服务器地址
SERVER_URL="http://localhost:80"

echo "测试用户信息相关接口..."
echo "服务器地址: $SERVER_URL"
echo ""

# 测试1: 获取用户信息
echo "=== 测试1: 获取用户信息 ==="
curl -X GET "$SERVER_URL/api/user/info?userId=1"
echo ""
echo ""

# 测试2: 绑定手机号
echo "=== 测试2: 绑定手机号 ==="
curl -X POST \
  -H "Content-Type: application/json" \
  -d "{
    \"userId\": 1,
    \"phone\": \"13800138000\",
    \"code\": \"123456\"
  }" \
  "$SERVER_URL/api/user/bind_phone"
echo ""
echo ""

echo "测试完成！"
echo "注意：确保数据库已创建Users表并插入了用户数据" 