#!/bin/bash

# 推广中心功能测试脚本

echo "🧪 开始测试推广中心功能..."

# 设置测试环境
BASE_URL="http://localhost:80"
TEST_USER_ID="507f1f77bcf86cd799439011"

echo "📋 测试1: 获取推广员信息"
curl -X GET "${BASE_URL}/api/promoter/info?userId=${TEST_USER_ID}" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n" \
  -s

echo -e "\n📋 测试2: 获取佣金记录列表"
curl -X GET "${BASE_URL}/api/promoter/commission_list?userId=${TEST_USER_ID}&page=1&pageSize=10" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n" \
  -s

echo -e "\n📋 测试3: 获取提现记录列表"
curl -X GET "${BASE_URL}/api/promoter/cashout_list?userId=${TEST_USER_ID}&page=1&pageSize=10" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n" \
  -s

echo -e "\n📋 测试4: 申请提现"
curl -X POST "${BASE_URL}/api/referral/apply_cashout" \
  -H "Content-Type: application/json" \
  -d "{
    \"userId\": \"${TEST_USER_ID}\",
    \"amount\": 100.00,
    \"method\": \"wechat\",
    \"account\": \"test_account\"
  }" \
  -w "\nHTTP状态码: %{http_code}\n" \
  -s

echo -e "\n📋 测试5: 获取推荐二维码"
curl -X GET "${BASE_URL}/api/referral/qrcode?userId=${TEST_USER_ID}" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n" \
  -s

echo -e "\n📋 测试6: 获取推荐报告"
curl -X GET "${BASE_URL}/api/referral/report?userId=${TEST_USER_ID}" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n" \
  -s

echo -e "\n📋 测试7: 获取推荐配置"
curl -X GET "${BASE_URL}/api/referral/config" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n" \
  -s

echo -e "\n🎉 推广中心功能测试完成！"
echo "📊 测试结果说明："
echo "  - HTTP状态码200: 接口正常"
echo "  - HTTP状态码404: 接口不存在"
echo "  - HTTP状态码500: 服务器错误"
echo "  - 响应体包含code字段: 业务逻辑正常" 