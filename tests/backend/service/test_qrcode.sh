#!/bin/bash

# 二维码功能测试脚本

BASE_URL="http://localhost:80"

echo "=== 二维码功能测试 ==="

# 1. 测试生成二维码URL
echo "1. 测试生成二维码URL..."
PROMOTER_CODE="ABC123"
curl -X GET "${BASE_URL}/api/qrcode/generate?promoterCode=${PROMOTER_CODE}" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n\n"

# 2. 测试生成Base64编码的二维码
echo "2. 测试生成Base64编码的二维码..."
curl -X GET "${BASE_URL}/api/qrcode/generate_base64?promoterCode=${PROMOTER_CODE}" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n\n"

# 3. 测试无效推广码
echo "3. 测试无效推广码..."
curl -X GET "${BASE_URL}/api/qrcode/generate?promoterCode=INVALID" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n\n"

# 4. 测试空推广码
echo "4. 测试空推广码..."
curl -X GET "${BASE_URL}/api/qrcode/generate?promoterCode=" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n\n"

# 5. 测试推广员信息（包含二维码URL）
echo "5. 测试推广员信息（包含二维码URL）..."
USER_ID="507f1f77bcf86cd799439011"
curl -X GET "${BASE_URL}/api/promoter/info?userId=${USER_ID}" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n\n"

echo "=== 测试完成 ===" 