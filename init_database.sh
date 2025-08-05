#!/bin/bash

# 数据库初始化脚本
echo "=== 开始初始化数据库 ==="

# 设置数据库连接信息
DB_HOST="localhost"
DB_PORT="3306"
DB_USER="root"
DB_PASSWORD="123456"
DB_NAME="anyuyinian"

# 创建数据库（如果不存在）
echo "1. 创建数据库..."
mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASSWORD -e "CREATE DATABASE IF NOT EXISTS $DB_NAME CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 执行基础表结构迁移
echo "2. 创建基础表结构..."
mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASSWORD $DB_NAME < db/migration/create_users_table.sql
mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASSWORD $DB_NAME < db/migration/create_home_tables.sql
mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASSWORD $DB_NAME < db/migration/create_service_order_referral_tables.sql
mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASSWORD $DB_NAME < db/migration/create_kefu_tables.sql
mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASSWORD $DB_NAME < db/migration/create_upload_tables.sql
mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASSWORD $DB_NAME < db/migration/create_user_extend_tables.sql

# 执行字段添加迁移
echo "3. 添加字段..."
mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASSWORD $DB_NAME < db/migration/add_phone_to_users.sql
mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASSWORD $DB_NAME < db/migration/add_order_fields.sql
mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASSWORD $DB_NAME < db/migration/add_patient_info_fields.sql
mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASSWORD $DB_NAME < db/migration/add_pay_deadline_field.sql

# 执行serviceId字段迁移
echo "4. 添加serviceId字段..."
mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASSWORD $DB_NAME < db/migration/add_service_id_field.sql

echo "=== 数据库初始化完成 ==="
echo "数据库名称: $DB_NAME"
echo "数据库主机: $DB_HOST:$DB_PORT"
echo "数据库用户: $DB_USER" 