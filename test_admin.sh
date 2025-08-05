#!/bin/bash

# 管理员功能测试脚本

echo "开始测试管理员功能..."

# 测试超级管理员登录
echo "1. 测试超级管理员登录..."
curl -X POST http://localhost:8080/admin/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "anyuyinian",
    "password": "000000"
  }'

echo -e "\n\n2. 测试错误密码..."
curl -X POST http://localhost:8080/admin/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "anyuyinian",
    "password": "wrong_password"
  }'

echo -e "\n\n3. 测试检查管理员状态..."
curl -X GET "http://localhost:8080/admin/check-status?userId=test_user_id"

echo -e "\n\n4. 测试获取管理员用户列表..."
curl -X GET "http://localhost:8080/admin/users?adminUserId=admin_super&page=1&pageSize=10"

echo -e "\n\n5. 测试获取管理员订单列表..."
curl -X GET "http://localhost:8080/admin/orders?adminUserId=admin_super&page=1&pageSize=10"

echo -e "\n\n管理员功能测试完成！" 