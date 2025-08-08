#!/bin/bash

# 测试服务分类筛选功能
echo "=== 测试服务分类筛选功能 ==="

# 设置基础URL
BASE_URL="http://localhost:80"

# 测试数据
ADMIN_USER_ID="test_admin_001"

echo "1. 测试获取所有服务列表..."
curl -X GET "${BASE_URL}/api/admin/services?adminUserId=${ADMIN_USER_ID}&page=1&pageSize=10" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n"

echo ""
echo "2. 测试获取服务分类列表..."
curl -X GET "${BASE_URL}/api/service/list?page=1&pageSize=1000" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n"

echo ""
echo "3. 测试按分类筛选服务（假设分类为'护理'）..."
curl -X GET "${BASE_URL}/api/admin/services?adminUserId=${ADMIN_USER_ID}&page=1&pageSize=10&category=护理" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n"

echo ""
echo "4. 测试按分类筛选服务（假设分类为'陪诊'）..."
curl -X GET "${BASE_URL}/api/admin/services?adminUserId=${ADMIN_USER_ID}&page=1&pageSize=10&category=陪诊" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n"

echo ""
echo "=== 测试完成 ===" 