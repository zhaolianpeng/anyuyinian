#!/bin/bash

# 测试UserId迁移功能
echo "开始测试UserId迁移功能..."

# 1. 测试ID生成工具
echo "1. 测试ID生成工具..."
curl -X POST http://localhost:8080/api/test/generate-userid \
  -H "Content-Type: application/json" \
  -d '{}'

echo -e "\n"

# 2. 测试用户登录（新用户）
echo "2. 测试用户登录（新用户）..."
curl -X POST http://localhost:8080/api/wx/login \
  -H "Content-Type: application/json" \
  -d '{
    "code": "test_code_new_user",
    "nickName": "测试用户",
    "avatarUrl": "https://example.com/avatar.jpg",
    "gender": 1
  }'

echo -e "\n"

# 3. 测试获取用户信息
echo "3. 测试获取用户信息..."
# 这里需要使用上一步返回的userId
USER_ID="test_user_id_here"
curl -X GET "http://localhost:8080/api/user/info?userId=$USER_ID"

echo -e "\n"

# 4. 测试地址管理
echo "4. 测试地址管理..."
curl -X POST http://localhost:8080/api/user/address \
  -H "Content-Type: application/json" \
  -d "{
    \"userId\": \"$USER_ID\",
    \"name\": \"张三\",
    \"phone\": \"13800138000\",
    \"province\": \"广东省\",
    \"city\": \"深圳市\",
    \"district\": \"南山区\",
    \"address\": \"科技园路1号\",
    \"isDefault\": true
  }"

echo -e "\n"

# 5. 测试就诊人管理
echo "5. 测试就诊人管理..."
curl -X POST http://localhost:8080/api/user/patient \
  -H "Content-Type: application/json" \
  -d "{
    \"userId\": \"$USER_ID\",
    \"name\": \"张三\",
    \"idCard\": \"440301199001011234\",
    \"phone\": \"13800138000\",
    \"gender\": 1,
    \"birthday\": \"1990-01-01\",
    \"relation\": \"本人\",
    \"isDefault\": true
  }"

echo -e "\n"

echo "UserId迁移功能测试完成！" 