#!/bin/bash

# 服务相关接口测试脚本

# 设置测试服务器地址
SERVER_URL="http://localhost:80"

echo "测试服务相关接口..."
echo "服务器地址: $SERVER_URL"
echo ""

# 测试1: 获取服务列表
echo "=== 测试1: 获取服务列表 ==="
curl -X GET "$SERVER_URL/api/service/list"
echo ""
echo ""

# 测试2: 获取指定分类的服务列表
echo "=== 测试2: 获取指定分类的服务列表 ==="
curl -X GET "$SERVER_URL/api/service/list?category=体检套餐&page=1&pageSize=5"
echo ""
echo ""

# 测试3: 获取服务详情
echo "=== 测试3: 获取服务详情 ==="
curl -X GET "$SERVER_URL/api/service/detail/1"
echo ""
echo ""

# 测试4: 获取服务表单配置
echo "=== 测试4: 获取服务表单配置 ==="
curl -X GET "$SERVER_URL/api/service/form_config/1"
echo ""
echo ""

echo "测试完成！"
echo "注意：确保数据库已创建ServiceItems表并插入了服务数据" 