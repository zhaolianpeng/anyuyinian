#!/bin/bash

# 为所有DAO文件添加SQL日志的脚本
# 这个脚本会为所有DAO操作添加详细的SQL日志记录

echo "=== 开始为DAO文件添加SQL日志 ==="

# 定义DAO文件目录
DAO_DIR="db/dao"

# 需要处理的DAO文件列表（排除接口文件和已处理的文件）
DAO_FILES=(
    "order_dao.go"
    "user_dao.go"
    "user_extend_dao.go"
    "home_dao.go"
    "admin_dao.go"
    "upload_dao.go"
    "referral_dao.go"
    "commission_dao.go"
    "cashout_dao.go"
    "kefu_dao.go"
    "config_dao.go"
    "consultation_dao.go"
)

echo "将处理以下DAO文件:"
for file in "${DAO_FILES[@]}"; do
    echo "  - $file"
done

echo ""
echo "注意：此脚本需要手动执行，因为每个DAO文件的结构不同"
echo "建议按以下步骤手动添加SQL日志："
echo ""
echo "1. 在每个DAO方法开始处添加："
echo "   logger := NewSQLLogger(\"操作类型\", \"表名\", 参数map)"
echo ""
echo "2. 在SQL操作后添加："
echo "   logger.LogQuery/LogInsert/LogUpdate/LogDelete(结果, 错误)"
echo ""
echo "3. 确保导入sql_logger包"
echo ""

# 检查sql_logger.go是否存在
if [ -f "$DAO_DIR/sql_logger.go" ]; then
    echo "✅ sql_logger.go 已存在"
else
    echo "❌ sql_logger.go 不存在，请先创建"
    exit 1
fi

# 检查每个DAO文件
echo ""
echo "检查DAO文件状态："
for file in "${DAO_FILES[@]}"; do
    if [ -f "$DAO_DIR/$file" ]; then
        echo "✅ $file 存在"
        
        # 检查是否已经包含SQL日志
        if grep -q "NewSQLLogger" "$DAO_DIR/$file"; then
            echo "  - 已包含SQL日志"
        else
            echo "  - 需要添加SQL日志"
        fi
    else
        echo "❌ $file 不存在"
    fi
done

echo ""
echo "=== 脚本执行完成 ==="
echo ""
echo "下一步操作："
echo "1. 手动为每个DAO文件添加SQL日志"
echo "2. 运行测试确保功能正常"
echo "3. 检查日志输出是否符合预期"
