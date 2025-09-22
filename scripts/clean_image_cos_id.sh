#!/bin/bash

# 清理 imageCosId 字段中的无效字符

echo "=== 清理 imageCosId 字段中的无效字符 ==="

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

# 执行清理
echo "1. 执行数据清理..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < scripts/clean_image_cos_id.sql

if [ $? -eq 0 ]; then
    echo "✅ 数据清理成功"
else
    echo "❌ 数据清理失败"
    exit 1
fi

# 验证清理结果
echo "2. 验证清理结果..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "
SELECT 
  id, 
  name, 
  imageCosId,
  LENGTH(imageCosId) as cos_id_length,
  CHAR_LENGTH(imageCosId) as char_length
FROM ServiceItems 
WHERE imageCosId IS NOT NULL 
  AND imageCosId != ''
ORDER BY id;"

echo ""
echo "=== 清理完成 ==="
echo "已移除 imageCosId 字段中的制表符、换行符等无效字符"
