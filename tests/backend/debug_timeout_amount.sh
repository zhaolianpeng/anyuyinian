#!/bin/bash

# 调试超时未支付金额显示为0的问题

echo "=== 调试超时未支付金额显示为0的问题 ==="

# 设置测试环境
BASE_URL="http://localhost:80"
ADMIN_USER_ID="anyuyinian"  # 超级管理员用户ID

echo "1. 测试管理员数据概览接口"
echo "请求URL: ${BASE_URL}/api/admin/stats?adminUserId=${ADMIN_USER_ID}"

# 发送请求
response=$(curl -s -X GET "${BASE_URL}/api/admin/stats?adminUserId=${ADMIN_USER_ID}")

echo "响应内容:"
echo "$response" | jq '.'

# 检查超时未支付金额
echo ""
echo "2. 检查超时未支付金额"
timeout_amount=$(echo "$response" | jq -r '.data.timeoutUnpaidAmount // 0')
echo "超时未支付金额: ¥$timeout_amount"

if [ "$timeout_amount" = "0" ] || [ "$timeout_amount" = "null" ]; then
    echo "❌ 超时未支付金额为0或null"
    echo ""
    echo "3. 可能的原因分析:"
    echo "   a) 数据库中没有超时未支付的订单"
    echo "   b) 订单的payDeadline字段为NULL"
    echo "   c) 订单状态不是待支付(status != 0 或 payStatus != 0)"
    echo "   d) 支付截止时间未过期(payDeadline >= NOW())"
    echo ""
    echo "4. 建议检查数据库中的订单数据:"
    echo "   - 运行 tests/backend/debug_timeout_orders.sql 脚本"
    echo "   - 检查是否有status=0且payStatus=0的订单"
    echo "   - 检查这些订单的payDeadline字段是否有值"
    echo "   - 检查payDeadline是否已经过期"
else
    echo "✅ 超时未支付金额正常: ¥$timeout_amount"
fi

# 检查其他金额字段
echo ""
echo "5. 检查其他金额字段"
paid_amount=$(echo "$response" | jq -r '.data.paidAmount // 0')
unpaid_amount=$(echo "$response" | jq -r '.data.unpaidAmount // 0')
refund_amount=$(echo "$response" | jq -r '.data.refundAmount // 0')
total_amount=$(echo "$response" | jq -r '.data.totalAmount // 0')

echo "已支付金额: ¥$paid_amount"
echo "待支付金额: ¥$unpaid_amount"
echo "退款金额: ¥$refund_amount"
echo "总金额: ¥$total_amount"

# 检查订单数量
echo ""
echo "6. 检查订单数量"
total_orders=$(echo "$response" | jq -r '.data.totalOrders // 0')
today_orders=$(echo "$response" | jq -r '.data.todayOrders // 0')
total_users=$(echo "$response" | jq -r '.data.totalUsers // 0')

echo "总订单数: $total_orders"
echo "今日订单数: $today_orders"
echo "总用户数: $total_users"

echo ""
echo "=== 调试完成 ==="
echo ""
echo "如果超时未支付金额为0，请检查:"
echo "1. 数据库中是否有订单数据"
echo "2. 订单的payDeadline字段是否正确设置"
echo "3. 是否有status=0且payStatus=0的订单"
echo "4. 这些订单的支付截止时间是否已过期" 