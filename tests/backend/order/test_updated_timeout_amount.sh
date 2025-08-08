#!/bin/bash

# 测试修改后的超时未支付金额功能

echo "=== 测试修改后的超时未支付金额功能 ==="

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
    echo "❌ 超时未支付金额仍为0或null"
    echo ""
    echo "3. 可能的原因:"
    echo "   a) 后端服务未重启，仍使用旧代码"
    echo "   b) 数据库中没有超时未支付的订单"
    echo "   c) 订单的payDeadline字段为NULL"
    echo "   d) 所有订单都已支付或已退款"
    echo ""
    echo "4. 建议检查:"
    echo "   - 重启后端服务"
    echo "   - 运行 tests/backend/verify_timeout_query.sql 验证数据"
    echo "   - 检查数据库中是否有status=3且payStatus=0的订单"
else
    echo "✅ 超时未支付金额正常: ¥$timeout_amount"
    echo ""
    echo "3. 功能验证成功！"
    echo "   - 修改后的查询逻辑已生效"
    echo "   - 包含已取消但未支付的订单"
fi

# 检查其他金额字段
echo ""
echo "4. 检查其他金额字段"
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
echo "5. 检查订单数量"
total_orders=$(echo "$response" | jq -r '.data.totalOrders // 0')
today_orders=$(echo "$response" | jq -r '.data.todayOrders // 0')
total_users=$(echo "$response" | jq -r '.data.totalUsers // 0')

echo "总订单数: $total_orders"
echo "今日订单数: $today_orders"
echo "总用户数: $total_users"

echo ""
echo "=== 测试完成 ==="
echo ""
echo "如果超时未支付金额仍为0，请:"
echo "1. 重启后端服务"
echo "2. 运行 tests/backend/verify_timeout_query.sql 检查数据"
echo "3. 确认数据库中有超时未支付的订单" 