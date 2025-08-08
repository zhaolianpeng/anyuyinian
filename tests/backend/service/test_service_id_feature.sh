#!/bin/bash

# 测试serviceitemid功能
echo "=== 测试serviceitemid功能 ==="

# 测试首页init接口，检查返回的services是否包含serviceitemid字段
echo "1. 测试首页init接口..."
curl -X POST "http://localhost:8080/api/home/init" \
  -H "Content-Type: application/json" \
  -d '{
    "longitude": 121.4737,
    "latitude": 31.2304,
    "limit": 10
  }' | jq '.data.services[0] | {id, serviceitemid, name, description}'

echo ""
echo "2. 检查返回的服务数据结构..."
curl -X POST "http://localhost:8080/api/home/init" \
  -H "Content-Type: application/json" \
  -d '{
    "longitude": 121.4737,
    "latitude": 31.2304,
    "limit": 10
  }' | jq '.data.services | length'

echo ""
echo "3. 验证serviceitemid字段值..."
curl -X POST "http://localhost:8080/api/home/init" \
  -H "Content-Type: application/json" \
  -d '{
    "longitude": 121.4737,
    "latitude": 31.2304,
    "limit": 10
  }' | jq '.data.services[] | {id, serviceitemid, name}'

echo ""
echo "=== 测试完成 ===" 