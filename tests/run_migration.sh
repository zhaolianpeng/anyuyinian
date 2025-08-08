#!/bin/bash

# 数据库迁移修复脚本
echo "开始运行数据库迁移修复..."

# 检查数据库连接
echo "检查数据库连接..."

# 运行迁移文件
echo "执行迁移文件: fix_referrer_id_nullable.sql"
mysql -u root -p < db/migration/fix_referrer_id_nullable.sql

echo "迁移完成！" 