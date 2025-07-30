#!/bin/bash

# 客服、医院相关接口测试脚本

# 设置测试服务器地址
SERVER_URL="http://localhost:80"

echo "测试客服、医院相关接口..."
echo "服务器地址: $SERVER_URL"
echo ""

# 测试1: 提交用户咨询问题
echo "=== 测试1: 提交用户咨询问题 ==="
curl -X POST \
  -H "Content-Type: application/json" \
  -d "{
    \"userId\": 1,
    \"userName\": \"张三\",
    \"userAvatar\": \"https://example.com/avatar.jpg\",
    \"content\": \"我想咨询一下体检套餐的具体内容\",
    \"images\": [\"https://example.com/image1.jpg\"]
  }" \
  "$SERVER_URL/api/kefu/send_msg"
echo ""
echo ""

# 测试2: 获取常见问题列表
echo "=== 测试2: 获取常见问题列表 ==="
curl -X GET "$SERVER_URL/api/kefu/faq"
echo ""
echo ""

# 测试3: 获取指定分类的常见问题
echo "=== 测试3: 获取预约服务分类的常见问题 ==="
curl -X GET "$SERVER_URL/api/kefu/faq?category=预约服务&page=1&pageSize=5"
echo ""
echo ""

# 测试4: 获取医院列表
echo "=== 测试4: 获取医院列表 ==="
curl -X GET "$SERVER_URL/api/hospital/list"
echo ""
echo ""

# 测试5: 根据位置获取医院列表
echo "=== 测试5: 根据位置获取医院列表 ==="
curl -X GET "$SERVER_URL/api/hospital/list?longitude=114.0579&latitude=22.5431&page=1&pageSize=10"
echo ""
echo ""

# 测试6: 获取医院详情
echo "=== 测试6: 获取医院详情 ==="
curl -X GET "$SERVER_URL/api/hospital/detail/1"
echo ""
echo ""

# 测试7: 获取医院详情（含导航信息）
echo "=== 测试7: 获取医院详情（含导航信息） ==="
curl -X GET "$SERVER_URL/api/hospital/detail/1?userLongitude=114.0579&userLatitude=22.5431"
echo ""
echo ""

echo "测试完成！"
echo "注意：确保数据库已创建相关表并插入了示例数据" 