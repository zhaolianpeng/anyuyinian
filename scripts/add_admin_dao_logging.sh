#!/bin/bash

# 为admin_dao.go添加SQL日志的脚本

echo "开始为admin_dao.go添加SQL日志..."

# 备份原文件
cp anyuyinian/db/dao/admin_dao.go anyuyinian/db/dao/admin_dao.go.backup

# 使用sed命令批量添加SQL日志
# 注意：这里只是示例，实际需要根据具体的方法来调整

echo "已备份原文件到 admin_dao.go.backup"
echo "请手动完成剩余方法的SQL日志添加"

# 显示需要添加日志的方法
echo "需要添加SQL日志的方法："
echo "1. RemoveAdmin"
echo "2. GetSubAdmins" 
echo "3. GetVisibleUsers"
echo "4. GetVisibleOrders"
echo "5. LogAdminLogin"
echo "6. GetAdminLoginLogs"

echo "脚本执行完成"
