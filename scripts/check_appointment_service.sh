#!/bin/bash

# 检查预约咨询服务状态脚本
# 创建时间：2024-12-20

echo "=== 检查预约咨询服务状态 ==="

# 设置颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 数据库配置 - 请根据实际环境修改
DB_HOST="localhost"
DB_PORT="3306"
DB_USER="root"
DB_PASSWORD=""
DB_NAME="anyuyinian"

echo -e "${BLUE}数据库配置:${NC}"
echo "  主机: $DB_HOST"
echo "  端口: $DB_PORT"
echo "  用户: $DB_USER"
echo "  数据库: $DB_NAME"
echo ""

# 检查数据库连接
echo -e "${BLUE}1. 检查数据库连接...${NC}"
if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "USE $DB_NAME;" 2>/dev/null; then
    echo -e "${GREEN}✓ 数据库连接成功${NC}"
else
    echo -e "${RED}✗ 数据库连接失败${NC}"
    echo "请检查数据库配置和连接信息"
    exit 1
fi

echo ""

# 检查ServiceItems表
echo -e "${BLUE}2. 检查ServiceItems表...${NC}"
echo "查询预约咨询服务项目:"
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "
SELECT 
    id,
    name,
    category,
    price,
    status,
    createdAt
FROM ServiceItems 
WHERE name = '预约咨询服务' OR category = '预约咨询';"

echo ""

# 检查Services表
echo -e "${BLUE}3. 检查Services表...${NC}"
echo "查询预约咨询服务:"
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "
SELECT 
    id,
    name,
    category,
    price,
    serviceitemid,
    status,
    createdAt
FROM Services 
WHERE name = '预约咨询服务' OR category = '预约咨询';"

echo ""

# 检查表结构
echo -e "${BLUE}4. 检查表结构...${NC}"
echo "ServiceItems表结构:"
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "
DESCRIBE ServiceItems;"

echo ""

echo "Services表结构:"
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "
DESCRIBE Services;"

echo ""

# 检查是否有数据
echo -e "${BLUE}5. 数据统计...${NC}"
echo "ServiceItems表总记录数:"
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "
SELECT COUNT(*) as total FROM ServiceItems;"

echo "Services表总记录数:"
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "
SELECT COUNT(*) as total FROM Services;"

echo ""

# 如果数据不存在，提供解决方案
echo -e "${BLUE}6. 问题诊断...${NC}"
SERVICE_EXISTS=$(mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "
SELECT COUNT(*) FROM Services WHERE name = '预约咨询服务';" -s -N)

if [ "$SERVICE_EXISTS" -eq 0 ]; then
    echo -e "${RED}✗ 预约咨询服务不存在${NC}"
    echo ""
    echo "解决方案:"
    echo "1. 执行数据库迁移脚本:"
    echo "   mysql -u root -p $DB_NAME < db/migration/20241220_add_appointment_consultation_service.sql"
    echo ""
    echo "2. 或者使用部署脚本:"
    echo "   chmod +x scripts/deploy_appointment_service.sh"
    echo "   ./scripts/deploy_appointment_service.sh"
    echo ""
    echo "3. 手动插入数据（如果迁移脚本有问题）"
else
    echo -e "${GREEN}✓ 预约咨询服务已存在${NC}"
    echo "如果仍有问题，请检查:"
    echo "1. 服务ID是否正确"
    echo "2. 服务状态是否为上架状态"
    echo "3. API接口是否正确调用"
fi

echo ""
echo -e "${GREEN}检查脚本执行完成！${NC}"
