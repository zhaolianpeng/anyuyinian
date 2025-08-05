#!/bin/bash

# 医院详情API测试脚本

echo "=== 测试医院详情API ==="

# 测试参数
BASE_URL="http://localhost:80"
# 如果是在云托管环境中，需要替换为实际的域名
# BASE_URL="https://your-domain.com"

echo "1. 测试有效的医院ID (ID=1)"
curl -X GET "${BASE_URL}/api/hospital/detail/1" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n"

echo -e "\n2. 测试有效的医院ID (ID=5)"
curl -X GET "${BASE_URL}/api/hospital/detail/5" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n"

echo -e "\n3. 测试无效的医院ID (ID=999)"
curl -X GET "${BASE_URL}/api/hospital/detail/999" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n"

echo -e "\n4. 测试带位置参数的医院详情 (ID=1)"
curl -X GET "${BASE_URL}/api/hospital/detail/1?userLongitude=121.4737&userLatitude=31.2304" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n"

echo -e "\n5. 测试缺少医院ID的请求"
curl -X GET "${BASE_URL}/api/hospital/detail/" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n"

echo -e "\n=== 测试完成 ===" 