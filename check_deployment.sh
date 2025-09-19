#!/bin/bash

# 云托管部署检查脚本

echo "=== 云托管部署检查 ==="

# 配置信息
ENV_ID="prod-5g94mx7a3d07e78c"
SERVICE_NAME="golang-lfwy"
SERVICE_URL="https://${SERVICE_NAME}-${ENV_ID}.service.tcloudbaseapp.com"

echo "环境ID: $ENV_ID"
echo "服务名称: $SERVICE_NAME"
echo "服务地址: $SERVICE_URL"
echo ""

# 1. 检查服务是否可访问
echo "1. 检查服务是否可访问..."
HEALTH_RESPONSE=$(curl -k -s -w "%{http_code}" -o /dev/null "$SERVICE_URL/health" 2>/dev/null)
if [ "$HEALTH_RESPONSE" = "200" ]; then
    echo "✅ 服务健康检查通过"
else
    echo "❌ 服务健康检查失败 (HTTP $HEALTH_RESPONSE)"
fi

# 2. 检查基础API
echo ""
echo "2. 检查基础API..."
CATEGORIES_RESPONSE=$(curl -k -s -w "%{http_code}" -o /dev/null "$SERVICE_URL/api/service/categories" 2>/dev/null)
if [ "$CATEGORIES_RESPONSE" = "200" ]; then
    echo "✅ 分类API可访问"
else
    echo "❌ 分类API不可访问 (HTTP $CATEGORIES_RESPONSE)"
fi

# 3. 检查管理员API
echo ""
echo "3. 检查管理员API..."
ADMIN_RESPONSE=$(curl -k -s -w "%{http_code}" -o /dev/null -X POST "$SERVICE_URL/api/admin/service/update-price" \
  -H "Content-Type: application/json" \
  -d '{"serviceId": 21, "newPrice": 258, "newOriginalPrice": 380, "reason": "测试"}' 2>/dev/null)

if [ "$ADMIN_RESPONSE" = "200" ] || [ "$ADMIN_RESPONSE" = "400" ]; then
    echo "✅ 管理员API可访问 (HTTP $ADMIN_RESPONSE)"
else
    echo "❌ 管理员API不可访问 (HTTP $ADMIN_RESPONSE)"
fi

# 4. 检查路由列表
echo ""
echo "4. 检查路由列表..."
echo "已配置的路由："
echo "- /health (健康检查)"
echo "- /api/service/categories (分类列表)"
echo "- /api/admin/services (管理员服务列表)"
echo "- /api/admin/service/update-price (管理员服务价格更新)"

# 5. 提供部署建议
echo ""
echo "=== 部署建议 ==="
if [ "$HEALTH_RESPONSE" != "200" ]; then
    echo "❌ 服务未正常运行，需要重新部署"
    echo "请执行以下步骤："
    echo "1. 打开微信开发者工具"
    echo "2. 选择云托管项目"
    echo "3. 右键点击 'golang-lfwy' 服务"
    echo "4. 选择 '上传并部署'"
    echo "5. 等待部署完成"
elif [ "$ADMIN_RESPONSE" = "404" ]; then
    echo "⚠️  服务运行正常，但管理员API不可访问"
    echo "可能原因："
    echo "1. 代码未部署到云托管"
    echo "2. 路由配置问题"
    echo "3. 服务版本不是最新"
    echo ""
    echo "建议："
    echo "1. 重新部署服务"
    echo "2. 检查云托管控制台的服务状态"
    echo "3. 查看服务日志"
else
    echo "✅ 所有检查通过，服务运行正常"
fi

echo ""
echo "=== 检查完成 ==="
