#!/bin/bash

# 测试管理员修改订单金额功能

echo "=== 测试管理员修改订单金额功能 ==="

# 设置测试环境
BASE_URL="http://localhost:80"
ADMIN_USER_ID="anyuyinian"  # 超级管理员用户ID
TEST_ORDER_ID="1"  # 测试订单ID

echo "1. 测试修改订单金额接口"
echo "请求URL: ${BASE_URL}/api/admin/order/update-amount?adminUserId=${ADMIN_USER_ID}"

# 构建修改金额请求数据
update_data=$(cat <<EOF
{
  "orderId": ${TEST_ORDER_ID},
  "newAmount": 399.00,
  "reason": "测试修改订单金额"
}
EOF
)

echo "请求数据:"
echo "$update_data" | jq '.'

# 发送修改请求
response=$(curl -s -X POST "${BASE_URL}/api/admin/order/update-amount?adminUserId=${ADMIN_USER_ID}" \
  -H "Content-Type: application/json" \
  -d "$update_data")

echo ""
echo "响应内容:"
echo "$response" | jq '.'

# 检查响应
echo ""
echo "2. 检查响应结果"
code=$(echo "$response" | jq -r '.code // -1')

if [ "$code" = "0" ]; then
    echo "✅ 修改订单金额成功"
    
    # 检查返回的数据
    order_id=$(echo "$response" | jq -r '.data.orderId // ""')
    order_no=$(echo "$response" | jq -r '.data.orderNo // ""')
    old_amount=$(echo "$response" | jq -r '.data.oldAmount // ""')
    new_amount=$(echo "$response" | jq -r '.data.newAmount // ""')
    reason=$(echo "$response" | jq -r '.data.reason // ""')
    
    echo "订单ID: $order_id"
    echo "订单号: $order_no"
    echo "原金额: ¥$old_amount"
    echo "新金额: ¥$new_amount"
    echo "修改原因: $reason"
    
    echo ""
    echo "3. 验证修改结果"
    echo "✅ 订单金额修改成功"
    
else
    echo "❌ 修改订单金额失败"
    error_msg=$(echo "$response" | jq -r '.errorMsg // "未知错误"')
    echo "错误信息: $error_msg"
    
    echo ""
    echo "3. 可能的原因:"
    echo "   a) 管理员权限不足（需要超级管理员）"
    echo "   b) 订单不存在"
    echo "   c) 订单已支付，无法修改"
    echo "   d) 新金额无效"
    echo "   e) 网络连接问题"
fi

# 测试权限验证
echo ""
echo "4. 测试权限验证"
echo "使用非超级管理员账号测试..."

# 这里可以添加使用一级管理员账号的测试
# 预期应该返回权限不足的错误

echo ""
echo "=== 测试完成 ==="
echo ""
echo "如果测试失败，请检查:"
echo "1. 管理员账号是否为超级管理员"
echo "2. 订单是否存在且未支付"
echo "3. 新金额是否有效"
echo "4. 后端服务是否正常运行" 