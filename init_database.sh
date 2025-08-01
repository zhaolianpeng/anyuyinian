#!/bin/bash

# 数据库初始化脚本
# 用于初始化数据库表结构和示例数据

echo "开始初始化数据库..."

# 数据库配置（请根据实际情况修改）
DB_HOST="localhost"
DB_PORT="3306"
DB_NAME="anyuyinian"
DB_USER="root"
DB_PASSWORD="your_password"

# 检查MySQL连接
echo "检查数据库连接..."
mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASSWORD -e "SELECT 1;" > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "错误：无法连接到数据库，请检查数据库配置"
    exit 1
fi

# 创建数据库（如果不存在）
echo "创建数据库..."
mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASSWORD -e "CREATE DATABASE IF NOT EXISTS $DB_NAME CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 执行数据库迁移脚本
echo "执行数据库迁移脚本..."
mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASSWORD $DB_NAME < db/migration/create_service_order_referral_tables.sql

# 初始化服务数据
echo "初始化服务数据..."
mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASSWORD $DB_NAME < init_service_data.sql

# 验证数据
echo "验证数据..."
mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASSWORD $DB_NAME -e "SELECT COUNT(*) as service_count FROM ServiceItems;"
mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASSWORD $DB_NAME -e "SELECT id, name, category, price FROM ServiceItems WHERE status = 1 ORDER BY sort;"

echo "数据库初始化完成！"
echo ""
echo "可用的服务ID："
mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASSWORD $DB_NAME -e "SELECT id, name FROM ServiceItems WHERE status = 1 ORDER BY id;" 