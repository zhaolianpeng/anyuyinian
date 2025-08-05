#!/bin/bash

# 测试API返回的serviceitemid字段
echo "=== 测试API返回的serviceitemid字段 ==="

# 测试首页init接口
echo "1. 调用首页init接口..."
RESPONSE=$(curl -s -X POST "http://localhost:8080/api/home/init" \
  -H "Content-Type: application/json" \
  -d '{
    "longitude": 121.4737,
    "latitude": 31.2304,
    "limit": 10
  }')

echo "API响应:"
echo "$RESPONSE" | jq '.'

echo ""
echo "2. 检查services数组中的serviceitemid字段..."
echo "$RESPONSE" | jq '.data.services[] | {id, serviceitemid, name, description}'

echo ""
echo "3. 验证serviceitemid字段值..."
echo "$RESPONSE" | jq '.data.services[] | "ID: \(.id), ServiceItemId: \(.serviceitemid), Name: \(.name)"'

echo ""
echo "=== 测试完成 ===" 