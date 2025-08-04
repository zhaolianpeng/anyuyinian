#!/bin/bash

# 构建脚本 - 用于本地测试和验证

echo "=== 开始构建项目 ==="

# 1. 清理之前的构建产物
echo "1. 清理构建产物..."
rm -f main
rm -f anyuyinian

# 2. 更新依赖
echo "2. 更新Go模块依赖..."
go mod tidy

# 3. 验证依赖
echo "3. 验证依赖..."
go mod verify

# 4. 构建项目
echo "4. 构建项目..."
GOOS=linux go build -o main .

# 5. 检查构建结果
if [ $? -eq 0 ]; then
    echo "✅ 构建成功！"
    echo "构建产物:"
    ls -la main
    echo ""
    echo "文件大小: $(du -h main | cut -f1)"
else
    echo "❌ 构建失败！"
    exit 1
fi

# 6. 检查文件类型
echo "6. 检查文件类型..."
file main

echo "=== 构建完成 ===" 