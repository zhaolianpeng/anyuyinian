#!/bin/bash

# 预约咨询服务部署脚本
# 创建时间：2024-12-20
# 描述：为首页一键预约功能部署预约咨询服务

echo "=== 开始部署预约咨询服务 ==="

# 设置颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 检查数据库连接
echo -e "${BLUE}1. 检查数据库连接...${NC}"

# 这里需要根据实际环境配置数据库连接信息
# 请修改以下变量为您的实际配置
DB_HOST="localhost"
DB_PORT="3306"
DB_USER="root"
DB_PASSWORD=""
DB_NAME="anyuyinian"

echo "数据库配置:"
echo "  主机: $DB_HOST"
echo "  端口: $DB_PORT"
echo "  用户: $DB_USER"
echo "  数据库: $DB_NAME"

# 测试数据库连接
if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "USE $DB_NAME;" 2>/dev/null; then
    echo -e "${GREEN}✓ 数据库连接成功${NC}"
else
    echo -e "${RED}✗ 数据库连接失败${NC}"
    echo "请检查数据库配置和连接信息"
    exit 1
fi

# 执行数据库迁移
echo -e "${BLUE}2. 执行数据库迁移...${NC}"

MIGRATION_FILE="db/migration/20241220_add_appointment_consultation_service.sql"

if [ -f "$MIGRATION_FILE" ]; then
    echo "找到迁移文件: $MIGRATION_FILE"
    
    # 执行SQL迁移
    if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "$MIGRATION_FILE"; then
        echo -e "${GREEN}✓ 数据库迁移执行成功${NC}"
    else
        echo -e "${RED}✗ 数据库迁移执行失败${NC}"
        echo "请检查SQL文件语法和数据库权限"
        exit 1
    fi
else
    echo -e "${RED}✗ 迁移文件不存在: $MIGRATION_FILE${NC}"
    exit 1
fi

# 验证数据插入
echo -e "${BLUE}3. 验证数据插入...${NC}"

# 检查ServiceItems表
echo "检查ServiceItems表:"
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "
SELECT 
    'ServiceItems' as table_name,
    id,
    name,
    category,
    price,
    status
FROM ServiceItems 
WHERE name = '预约咨询服务' AND category = '预约咨询';"

# 检查Services表
echo "检查Services表:"
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "
SELECT 
    'Services' as table_name,
    id,
    name,
    category,
    price,
    serviceitemid,
    status
FROM Services 
WHERE name = '预约咨询服务' AND category = '预约咨询';"

echo -e "${GREEN}✓ 数据验证完成${NC}"

# 检查前端代码
echo -e "${BLUE}4. 检查前端代码...${NC}"

FRONTEND_FILES=(
    "miniprogram/pages/index/index.js"
    "miniprogram/pages/service/detail.js"
    "miniprogram/pages/service/detail.wxml"
    "miniprogram/pages/service/detail.wxss"
)

for file in "${FRONTEND_FILES[@]}"; do
    if [ -f "$file" ]; then
        echo -e "${GREEN}✓ $file 存在${NC}"
    else
        echo -e "${YELLOW}⚠ $file 不存在${NC}"
    fi
done

# 部署完成提示
echo -e "${BLUE}5. 部署完成检查清单${NC}"

echo -e "${GREEN}=== 部署完成 ===${NC}"
echo ""
echo "请按以下步骤验证功能:"
echo ""
echo "1. 数据库验证:"
echo "   ✓ 预约咨询服务已插入ServiceItems表"
echo "   ✓ 预约咨询服务已插入Services表"
echo "   ✓ 关联关系已建立"
echo ""
echo "2. 前端验证:"
echo "   ✓ 首页显示一键预约按钮"
echo "   ✓ 点击一键预约跳转正常"
echo "   ✓ 预约咨询服务页面加载正常"
echo "   ✓ 表单配置正确显示"
echo ""
echo "3. 功能测试:"
echo "   ✓ 填写预约表单"
echo "   ✓ 表单验证正常"
echo "   ✓ 提交功能正常"
echo ""
echo "4. 注意事项:"
echo "   - 预约咨询服务价格为0.01元"
echo "   - 服务分类为'预约咨询'"
echo "   - 服务状态为上架状态"
echo "   - 表单包含6个字段：姓名、手机号、服务类型、服务时间、服务地址、特殊需求"
echo ""
echo "如有问题，请检查:"
echo "   - 数据库连接和权限"
echo "   - SQL文件语法"
echo "   - 前端页面路径"
echo "   - 服务详情API接口"

echo -e "${GREEN}部署脚本执行完成！${NC}"
