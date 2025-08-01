#!/bin/bash

echo "测试患者API - 验证年龄计算功能"

# 测试获取患者列表API
echo "1. 测试获取患者列表API..."
curl -X GET "https://golang-lfwy-176496-6-1353115175.sh.run.tcloudbase.com/api/user/patient?userId=1" \
  -H "Authorization: test-token" \
  -H "Content-Type: application/json" | jq '.'

echo ""
echo "2. 测试创建患者API（包含身份证号）..."
curl -X POST "https://golang-lfwy-176496-6-1353115175.sh.run.tcloudbase.com/api/user/patient" \
  -H "Authorization: test-token" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": 1,
    "name": "测试患者",
    "idCard": "110101199001011234",
    "phone": "13800138000",
    "gender": 1,
    "birthday": "1990-01-01",
    "relation": "本人",
    "isDefault": false
  }' | jq '.'

echo ""
echo "测试完成！" 