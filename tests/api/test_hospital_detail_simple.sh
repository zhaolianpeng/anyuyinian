#!/bin/bash

# 简单的医院详情API测试

echo "=== 简单测试医院详情API ==="

# 测试医院列表API
echo "1. 测试医院列表API"
curl -X GET "http://localhost:80/api/hospital/list" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n"

echo -e "\n2. 测试医院详情API (ID=1)"
curl -X GET "http://localhost:80/api/hospital/detail/1" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n"

echo -e "\n3. 测试医院详情API (ID=5)"
curl -X GET "http://localhost:80/api/hospital/detail/5" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n"

echo -e "\n=== 测试完成 ===" 