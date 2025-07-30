#!/bin/bash

# 用户地址管理接口测试脚本

# 设置测试服务器地址
SERVER_URL="http://localhost:80"

echo "测试用户地址管理接口..."
echo "服务器地址: $SERVER_URL"
echo ""

# 测试1: 获取地址列表
echo "=== 测试1: 获取地址列表 ==="
curl -X GET "$SERVER_URL/api/user/address?userId=1"
echo ""
echo ""

# 测试2: 添加地址
echo "=== 测试2: 添加地址 ==="
curl -X POST \
  -H "Content-Type: application/json" \
  -d "{
    \"userId\": 1,
    \"name\": \"张三\",
    \"phone\": \"13800138000\",
    \"province\": \"广东省\",
    \"city\": \"深圳市\",
    \"district\": \"罗湖区\",
    \"address\": \"东门北路1017号\",
    \"isDefault\": true
  }" \
  "$SERVER_URL/api/user/address"
echo ""
echo ""

# 测试3: 更新地址
echo "=== 测试3: 更新地址 ==="
curl -X PUT \
  -H "Content-Type: application/json" \
  -d "{
    \"id\": 1,
    \"userId\": 1,
    \"name\": \"张三\",
    \"phone\": \"13800138000\",
    \"province\": \"广东省\",
    \"city\": \"深圳市\",
    \"district\": \"罗湖区\",
    \"address\": \"东门北路1017号更新\",
    \"isDefault\": true
  }" \
  "$SERVER_URL/api/user/address"
echo ""
echo ""

# 测试4: 删除地址
echo "=== 测试4: 删除地址 ==="
curl -X DELETE "$SERVER_URL/api/user/address?id=1&userId=1"
echo ""
echo ""

echo "测试完成！"
echo "注意：确保数据库已创建UserAddresses表" 