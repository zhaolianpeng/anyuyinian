#!/bin/bash

# 用户就诊人管理接口测试脚本

# 设置测试服务器地址
SERVER_URL="http://localhost:80"

echo "测试用户就诊人管理接口..."
echo "服务器地址: $SERVER_URL"
echo ""

# 测试1: 获取就诊人列表
echo "=== 测试1: 获取就诊人列表 ==="
curl -X GET "$SERVER_URL/api/user/patient?userId=1"
echo ""
echo ""

# 测试2: 添加就诊人
echo "=== 测试2: 添加就诊人 ==="
curl -X POST \
  -H "Content-Type: application/json" \
  -d "{
    \"userId\": 1,
    \"name\": \"张三\",
    \"idCard\": \"440301199001011234\",
    \"phone\": \"13800138000\",
    \"gender\": 1,
    \"birthday\": \"1990-01-01\",
    \"relationship\": \"本人\",
    \"isDefault\": true
  }" \
  "$SERVER_URL/api/user/patient"
echo ""
echo ""

# 测试3: 更新就诊人
echo "=== 测试3: 更新就诊人 ==="
curl -X PUT \
  -H "Content-Type: application/json" \
  -d "{
    \"id\": 1,
    \"userId\": 1,
    \"name\": \"张三\",
    \"idCard\": \"440301199001011234\",
    \"phone\": \"13800138000\",
    \"gender\": 1,
    \"birthday\": \"1990-01-01\",
    \"relationship\": \"本人\",
    \"isDefault\": true
  }" \
  "$SERVER_URL/api/user/patient"
echo ""
echo ""

# 测试4: 删除就诊人
echo "=== 测试4: 删除就诊人 ==="
curl -X DELETE "$SERVER_URL/api/user/patient?id=1&userId=1"
echo ""
echo ""

echo "测试完成！"
echo "注意：确保数据库已创建Patients表" 