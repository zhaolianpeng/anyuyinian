#!/bin/bash

# 添加 imageCosId 字段到 ServiceItems 表的脚本

echo "=== 添加 imageCosId 字段到 ServiceItems 表 ==="

# 数据库配置
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-3306}"
DB_NAME="${DB_NAME:-anyuyinian}"
DB_USER="${DB_USER:-root}"
DB_PASSWORD="${DB_PASSWORD:-}"

# 检查 MySQL 是否可用
if ! command -v mysql &> /dev/null; then
    echo "❌ MySQL 客户端未安装"
    exit 1
fi

# 执行迁移
echo "1. 执行数据库迁移..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < db/migration/add_image_cos_id_to_service_items.sql

if [ $? -eq 0 ]; then
    echo "✅ 数据库迁移成功"
else
    echo "❌ 数据库迁移失败"
    exit 1
fi

# 验证字段是否添加成功
echo "2. 验证字段添加..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "DESCRIBE ServiceItems;" | grep imageCosId

if [ $? -eq 0 ]; then
    echo "✅ imageCosId 字段添加成功"
else
    echo "❌ imageCosId 字段添加失败"
    exit 1
fi

echo ""
echo "=== 迁移完成 ==="
echo "ServiceItems 表已成功添加 imageCosId 字段"
echo "字段类型: VARCHAR(255)"
echo "注释: 腾讯云对象存储图片ID"
echo "索引: idx_service_items_image_cos_id"
