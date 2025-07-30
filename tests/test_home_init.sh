#!/bin/bash

# 首页初始化接口测试脚本

# 设置测试服务器地址
SERVER_URL="http://localhost:80"

echo "测试首页初始化接口..."
echo "服务器地址: $SERVER_URL"
echo ""

# 测试1: 获取首页数据（无位置信息）
echo "=== 测试1: 获取首页数据（无位置信息） ==="
curl -X GET "$SERVER_URL/api/home/init"
echo ""
echo ""

# 测试2: 获取首页数据（带位置信息）
echo "=== 测试2: 获取首页数据（带位置信息） ==="
curl -X GET "$SERVER_URL/api/home/init?longitude=114.0579&latitude=22.5431&limit=5"
echo ""
echo ""

# 测试3: 获取首页数据（限制数量）
echo "=== 测试3: 获取首页数据（限制数量） ==="
curl -X GET "$SERVER_URL/api/home/init?limit=3"
echo ""
echo ""

echo "测试完成！"
echo "注意：确保数据库已创建相关表并插入了示例数据" 