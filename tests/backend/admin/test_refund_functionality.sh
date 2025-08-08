#!/bin/bash

# 测试退款功能

echo "=== 测试退款功能 ==="

# 设置测试环境
BASE_URL="http://localhost:80"
ADMIN_USER_ID="anyuyinian"  # 超级管理员用户ID
TEST_ORDER_ID="1"  # 测试订单ID

echo "1. 测试用户申请退款"
echo "请求URL: ${BASE_URL}/api/order/refund/${TEST_ORDER_ID}"

# 构建退款请求数据
refund_data=$(cat <<EOF
{
  "orderId": ${TEST_ORDER_ID},
  "refundAmount": 299.00,
  "reason": "用户申请退款测试"
}
EOF
)

echo "请求数据:"
echo "$refund_data" | jq '.'

# 发送退款申请请求
response=$(curl -s -X POST "${BASE_URL}/api/order/refund/${TEST_ORDER_ID}" \
  -H "Content-Type: application/json" \
  -d "$refund_data")

echo ""
echo "响应内容:"
echo "$response" | jq '.'

# 检查响应
echo ""
echo "2. 检查退款申请响应结果"
code=$(echo "$response" | jq -r '.code // -1')

if [ "$code" = "0" ]; then
    echo "✅ 退款申请成功"
    
    # 检查返回的数据
    order_id=$(echo "$response" | jq -r '.data.orderId // ""')
    order_no=$(echo "$response" | jq -r '.data.orderNo // ""')
    refund_amount=$(echo "$response" | jq -r '.data.refundAmount // ""')
    reason=$(echo "$response" | jq -r '.data.reason // ""')
    
    echo "订单ID: $order_id"
    echo "订单号: $order_no"
    echo "退款金额: ¥$refund_amount"
    echo "退款原因: $reason"
    
else
    echo "❌ 退款申请失败"
    error_msg=$(echo "$response" | jq -r '.errorMsg // "未知错误"')
    echo "错误信息: $error_msg"
fi

echo ""
echo "3. 测试管理员处理退款"
echo "请求URL: ${BASE_URL}/api/admin/order/refund?adminUserId=${ADMIN_USER_ID}"

# 构建管理员退款请求数据
admin_refund_data=$(cat <<EOF
{
  "orderId": ${TEST_ORDER_ID},
  "refundAmount": 299.00,
  "reason": "管理员处理退款测试",
  "refundStatus": 2
}
EOF
)

echo "请求数据:"
echo "$admin_refund_data" | jq '.'

# 发送管理员退款请求
admin_response=$(curl -s -X POST "${BASE_URL}/api/admin/order/refund?adminUserId=${ADMIN_USER_ID}" \
  -H "Content-Type: application/json" \
  -d "$admin_refund_data")

echo ""
echo "响应内容:"
echo "$admin_response" | jq '.'

# 检查管理员退款响应
echo ""
echo "4. 检查管理员退款响应结果"
admin_code=$(echo "$admin_response" | jq -r '.code // -1')

if [ "$admin_code" = "0" ]; then
    echo "✅ 管理员退款处理成功"
    
    # 检查返回的数据
    order_id=$(echo "$admin_response" | jq -r '.data.orderId // ""')
    order_no=$(echo "$admin_response" | jq -r '.data.orderNo // ""')
    refund_amount=$(echo "$admin_response" | jq -r '.data.refundAmount // ""')
    refund_status=$(echo "$admin_response" | jq -r '.data.refundStatus // ""')
    reason=$(echo "$admin_response" | jq -r '.data.reason // ""')
    message=$(echo "$admin_response" | jq -r '.data.message // ""')
    
    echo "订单ID: $order_id"
    echo "订单号: $order_no"
    echo "退款金额: ¥$refund_amount"
    echo "退款状态: $refund_status"
    echo "退款原因: $reason"
    echo "处理结果: $message"
    
else
    echo "❌ 管理员退款处理失败"
    error_msg=$(echo "$admin_response" | jq -r '.errorMsg // "未知错误"')
    echo "错误信息: $error_msg"
fi

echo ""
echo "5. 测试错误情况"

# 测试1: 无效的订单ID
echo "测试1: 无效的订单ID"
invalid_response=$(curl -s -X POST "${BASE_URL}/api/order/refund/999" \
  -H "Content-Type: application/json" \
  -d '{"orderId": 999, "refundAmount": 100.00, "reason": "测试"}')

echo "响应:"
echo "$invalid_response" | jq '.'

# 测试2: 无效的退款金额
echo ""
echo "测试2: 无效的退款金额"
invalid_amount_response=$(curl -s -X POST "${BASE_URL}/api/order/refund/${TEST_ORDER_ID}" \
  -H "Content-Type: application/json" \
  -d '{"orderId": 1, "refundAmount": -100.00, "reason": "测试"}')

echo "响应:"
echo "$invalid_amount_response" | jq '.'

# 测试3: 权限不足
echo ""
echo "测试3: 权限不足"
unauthorized_response=$(curl -s -X POST "${BASE_URL}/api/admin/order/refund?adminUserId=invalid_user" \
  -H "Content-Type: application/json" \
  -d '{"orderId": 1, "refundAmount": 100.00, "reason": "测试", "refundStatus": 2}')

echo "响应:"
echo "$unauthorized_response" | jq '.'

echo ""
echo "=== 测试完成 ==="
echo ""
echo "如果所有测试都通过，说明退款功能正常工作！"
echo ""
echo "退款功能特性:"
echo "1. ✅ 用户申请退款"
echo "2. ✅ 管理员处理退款"
echo "3. ✅ 权限控制"
echo "4. ✅ 参数验证"
echo "5. ✅ 状态管理"
echo "6. ✅ 错误处理" 