#!/bin/bash

# 云托管部署脚本
# 用于重新部署更新后的服务到微信云托管

set -e

echo "=== 开始部署云托管服务 ==="

# 配置信息
ENV_ID="prod-5g94mx7a3d07e78c"
SERVICE_NAME="golang-lfwy"

echo "环境ID: $ENV_ID"
echo "服务名称: $SERVICE_NAME"

# 1. 构建项目
echo "1. 构建项目..."
./build.sh

if [ $? -ne 0 ]; then
    echo "❌ 构建失败，部署终止"
    exit 1
fi

# 2. 检查构建产物
echo "2. 检查构建产物..."
if [ ! -f "main" ]; then
    echo "❌ 构建产物不存在，部署终止"
    exit 1
fi

echo "✅ 构建产物检查通过"

# 3. 部署到云托管
echo "3. 部署到云托管..."
echo "正在上传并部署服务，请稍候..."

# 使用微信开发者工具的命令行工具部署
# 注意：这里需要根据您的实际部署方式调整
echo "请手动执行以下步骤："
echo ""
echo "方法1: 使用微信开发者工具"
echo "1. 打开微信开发者工具"
echo "2. 选择云托管项目"
echo "3. 右键点击 'golang-lfwy' 服务"
echo "4. 选择 '上传并部署'"
echo "5. 等待部署完成"
echo ""
echo "方法2: 使用命令行工具"
echo "1. 确保已安装微信开发者工具命令行"
echo "2. 执行: cli cloud deploy --env $ENV_ID --service $SERVICE_NAME"
echo ""

# 4. 等待用户确认部署完成
echo "4. 等待部署确认..."
read -p "部署完成后，按回车键继续测试..."

# 5. 测试API
echo "5. 测试API..."
echo "测试分类API..."

# 获取云托管服务地址
SERVICE_URL="https://${SERVICE_NAME}-${ENV_ID}.service.tcloudbaseapp.com"

echo "服务地址: $SERVICE_URL"

# 测试分类API
echo "测试分类API: $SERVICE_URL/api/service/categories"
RESPONSE=$(curl -s -X GET "$SERVICE_URL/api/service/categories")

echo "响应状态: $?"
echo "响应内容:"
echo "$RESPONSE" | head -20

# 检查响应
if echo "$RESPONSE" | grep -q "智慧养老"; then
    echo "✅ 分类API测试成功，包含智慧养老分类"
elif echo "$RESPONSE" | grep -q "code.*0"; then
    echo "✅ 分类API测试成功"
else
    echo "❌ 分类API测试失败"
    echo "可能的原因："
    echo "1. 服务未正确部署"
    echo "2. 服务未完全启动"
    echo "3. 网络连接问题"
    echo ""
    echo "建议："
    echo "1. 检查云托管控制台的服务状态"
    echo "2. 查看服务日志"
    echo "3. 等待几分钟后重试"
fi

echo ""
echo "=== 部署完成 ==="
echo "如果API测试失败，请："
echo "1. 检查云托管控制台的服务状态"
echo "2. 查看服务日志确认是否有错误"
echo "3. 等待服务完全启动后重试"
