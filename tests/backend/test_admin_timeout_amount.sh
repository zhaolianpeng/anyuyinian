#!/bin/bash

# 测试管理员首页超时未支付金额功能

echo "=== 测试管理员首页超时未支付金额功能 ==="

# 设置测试环境
BASE_URL="http://localhost:80"
ADMIN_USER_ID="anyuyinian"  # 超级管理员用户ID

echo "1. 测试管理员数据概览接口（包含超时未支付金额）"
echo "请求URL: ${BASE_URL}/api/admin/stats?adminUserId=${ADMIN_USER_ID}"

# 发送请求
response=$(curl -s -X GET "${BASE_URL}/api/admin/stats?adminUserId=${ADMIN_USER_ID}")

echo "响应内容:"
echo "$response" | jq '.'

# 检查响应格式
echo ""
echo "2. 验证响应字段"
echo "$response" | jq -r '.data | keys[]' | while read key; do
    echo "字段: $key"
done

# 检查是否包含超时未支付金额字段
if echo "$response" | jq -e '.data.timeoutUnpaidAmount' > /dev/null; then
    echo "✅ 超时未支付金额字段存在"
    timeout_amount=$(echo "$response" | jq -r '.data.timeoutUnpaidAmount')
    echo "超时未支付金额: ¥$timeout_amount"
else
    echo "❌ 超时未支付金额字段不存在"
fi

# 检查其他必要字段
required_fields=("totalUsers" "totalOrders" "todayOrders" "totalAmount" "paidAmount" "unpaidAmount" "refundAmount")
for field in "${required_fields[@]}"; do
    if echo "$response" | jq -e ".data.$field" > /dev/null; then
        echo "✅ $field 字段存在"
    else
        echo "❌ $field 字段缺失"
    fi
done

echo ""
echo "=== 测试完成 ===" 