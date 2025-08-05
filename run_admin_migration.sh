#!/bin/bash

# 管理员数据库迁移脚本

echo "开始执行管理员数据库迁移..."

# 设置数据库连接信息
DB_HOST="localhost"
DB_PORT="3306"
DB_NAME="your_database_name"
DB_USER="your_username"
DB_PASS="your_password"

# 执行迁移SQL
mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASS $DB_NAME < db/migration/add_admin_fields.sql

if [ $? -eq 0 ]; then
    echo "管理员数据库迁移成功！"
    echo "默认超级管理员账号：anyuyinian"
    echo "默认密码：000000"
else
    echo "管理员数据库迁移失败！"
    exit 1
fi 