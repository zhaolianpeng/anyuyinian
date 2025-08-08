#!/bin/bash

echo "=== 测试订单API修改 ==="

# 设置基础URL
BASE_URL="http://localhost:8080"

echo "1. 测试订单详情接口 (POST /api/order/detail)"
echo "请求参数: {\"orderNo\": \"ORDER20241201000001\"}"
curl -X POST "${BASE_URL}/api/order/detail" \
  -H "Content-Type: application/json" \
  -d '{"orderNo": "ORDER20241201000001"}' \
  | jq '.'

echo -e "\n2. 测试订单列表接口 (GET /api/order/list)"
echo "请求参数: userId=1&page=1&pageSize=5"
curl -X GET "${BASE_URL}/api/order/list?userId=1&page=1&pageSize=5" \
  | jq '.'

echo -e "\n3. 验证订单列表返回字段"
echo "检查是否包含以下字段："
echo "- serviceName (服务名称)"
echo "- appointmentDate (预约日期)"
echo "- appointmentTime (预约时间)"
echo "- consultTime (沟通时间)"
echo "- totalAmount (订单金额)"
echo "- statusText (状态文本)"
echo "- formattedAmount (格式化金额)"

echo -e "\n=== 测试完成 ===" 