#!/bin/bash

# 平台配置接口测试脚本

# 设置测试服务器地址
SERVER_URL="http://localhost:80"

echo "测试平台配置接口..."
echo "服务器地址: $SERVER_URL"
echo ""

# 测试获取平台配置
echo "=== 测试获取平台配置 ==="
curl -X GET "$SERVER_URL/api/config"
echo ""
echo ""

echo "测试完成！"
echo "注意：确保数据库已创建Configs表并插入了配置数据" 