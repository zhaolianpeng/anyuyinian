#!/bin/bash

# 推广码功能测试脚本

BASE_URL="http://localhost:80"

echo "=== 推广码功能测试 ==="

# 1. 测试生成推广码
echo "1. 测试批量生成推广码..."
curl -X POST "${BASE_URL}/api/promoter/generate_codes" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n\n"

# 2. 测试获取推广员信息（需要有效的userId）
echo "2. 测试获取推广员信息..."
# 这里需要替换为实际的userId
USER_ID="test_user_123"
curl -X GET "${BASE_URL}/api/promoter/info?userId=${USER_ID}" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n\n"

# 3. 测试通过推广码查找用户
echo "3. 测试通过推广码查找用户..."
# 这里需要替换为实际的推广码
PROMOTER_CODE="ABC123"
curl -X GET "${BASE_URL}/api/promoter/find_user?promoterCode=${PROMOTER_CODE}" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n\n"

# 4. 测试无效推广码
echo "4. 测试无效推广码..."
curl -X GET "${BASE_URL}/api/promoter/find_user?promoterCode=INVALID" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n\n"

# 5. 测试空推广码
echo "5. 测试空推广码..."
curl -X GET "${BASE_URL}/api/promoter/find_user?promoterCode=" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n\n"

echo "=== 测试完成 ===" 