#!/bin/bash

# 推荐系统接口测试脚本

# 设置测试服务器地址
SERVER_URL="http://localhost:80"

echo "测试推荐系统接口..."
echo "服务器地址: $SERVER_URL"
echo ""

# 测试1: 获取推广二维码
echo "=== 测试1: 获取推广二维码 ==="
curl -X GET "$SERVER_URL/api/referral/qrcode?userId=1"
echo ""
echo ""

# 测试2: 获取推荐报告
echo "=== 测试2: 获取推荐报告 ==="
curl -X GET "$SERVER_URL/api/referral/report?userId=1&page=1&pageSize=10"
echo ""
echo ""

# 测试3: 获取推荐配置
echo "=== 测试3: 获取推荐配置 ==="
curl -X GET "$SERVER_URL/api/referral/config"
echo ""
echo ""

# 测试4: 申请佣金提现
echo "=== 测试4: 申请佣金提现 ==="
curl -X POST \
  -H "Content-Type: application/json" \
  -d "{
    \"userId\": 1,
    \"amount\": 100.00,
    \"bankName\": \"中国银行\",
    \"bankAccount\": \"6222021234567890123\",
    \"accountName\": \"张三\",
    \"phone\": \"13800138000\"
  }" \
  "$SERVER_URL/api/referral/apply_cashout"
echo ""
echo ""

echo "测试完成！"
echo "注意：确保数据库已创建Referrals、Commissions、Cashouts表" 