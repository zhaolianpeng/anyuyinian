#!/bin/bash

# 设置咨询表脚本
echo "开始设置咨询相关数据表..."

# 数据库连接信息（请根据实际情况修改）
DB_HOST="localhost"
DB_PORT="3306"
DB_USER="root"
DB_PASSWORD="your_password"
DB_NAME="your_database_name"

# 执行SQL文件
echo "执行数据库迁移脚本..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < db/migration/20241220_create_consultation_tables.sql

if [ $? -eq 0 ]; then
    echo "✅ 咨询表创建成功！"
    echo ""
    echo "已创建的表："
    echo "- consultations (咨询会话表)"
    echo "- consultation_messages (咨询消息表)"
    echo "- consultation_notifications (咨询通知表)"
    echo ""
    echo "测试数据已插入，可以开始测试咨询功能。"
else
    echo "❌ 咨询表创建失败，请检查数据库连接和SQL脚本。"
    exit 1
fi

echo ""
echo "下一步："
echo "1. 重新编译并部署后端服务"
echo "2. 测试咨询API接口"
echo "3. 在小程序中测试咨询功能"
