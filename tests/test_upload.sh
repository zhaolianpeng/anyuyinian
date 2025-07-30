#!/bin/bash

# 文件上传接口测试脚本

# 设置测试服务器地址
SERVER_URL="http://localhost:80"

echo "测试文件上传接口..."
echo "服务器地址: $SERVER_URL"
echo ""

# 测试1: 上传文件
echo "=== 测试1: 上传文件 ==="
# 创建一个测试文件
echo "这是一个测试文件内容" > test_file.txt

curl -X POST \
  -F "file=@test_file.txt" \
  -F "userId=1" \
  -F "category=test" \
  "$SERVER_URL/api/upload"
echo ""
echo ""

# 测试2: 获取文件列表
echo "=== 测试2: 获取文件列表 ==="
curl -X GET "$SERVER_URL/api/files?userId=1&page=1&pageSize=10"
echo ""
echo ""

# 测试3: 获取指定分类的文件列表
echo "=== 测试3: 获取指定分类的文件列表 ==="
curl -X GET "$SERVER_URL/api/files?userId=1&category=test&page=1&pageSize=5"
echo ""
echo ""

# 清理测试文件
rm -f test_file.txt

echo "测试完成！"
echo "注意：确保uploads目录存在且有写入权限" 