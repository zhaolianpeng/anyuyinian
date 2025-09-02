#!/bin/bash

# SQL日志测试脚本
# 用于验证SQL日志系统是否正常工作

echo "=== SQL日志系统测试 ==="

# 检查必要的文件是否存在
echo "1. 检查必要文件..."
if [ -f "db/dao/sql_logger.go" ]; then
    echo "✅ sql_logger.go 存在"
else
    echo "❌ sql_logger.go 不存在"
    exit 1
fi

if [ -f "service/log_middleware.go" ]; then
    echo "✅ log_middleware.go 存在"
else
    echo "❌ log_middleware.go 不存在"
    exit 1
fi

# 检查已添加日志的DAO文件
echo ""
echo "2. 检查已添加SQL日志的DAO文件..."

dao_files=(
    "db/dao/service_dao.go"
    "db/dao/order_dao.go"
    "db/dao/user_dao.go"
)

for file in "${dao_files[@]}"; do
    if [ -f "$file" ]; then
        if grep -q "NewSQLLogger" "$file"; then
            echo "✅ $file 已包含SQL日志"
        else
            echo "⚠️  $file 未包含SQL日志"
        fi
    else
        echo "❌ $file 不存在"
    fi
done

# 编译检查
echo ""
echo "3. 编译检查..."
if go build -o /tmp/test_build main.go 2>/dev/null; then
    echo "✅ 编译成功"
    rm -f /tmp/test_build
else
    echo "❌ 编译失败"
    echo "请检查代码语法错误"
    exit 1
fi

# 创建测试日志文件
echo ""
echo "4. 创建测试日志文件..."
test_log_file="/tmp/sql_logging_test.log"
echo "测试日志将写入: $test_log_file"

# 启动应用进行测试（后台运行）
echo ""
echo "5. 启动应用进行测试..."
echo "注意: 应用将在后台运行，测试完成后会自动停止"

# 启动应用
go run main.go > "$test_log_file" 2>&1 &
APP_PID=$!

# 等待应用启动
echo "等待应用启动..."
sleep 3

# 测试API调用
echo ""
echo "6. 测试API调用..."

# 测试服务列表API
echo "测试服务列表API..."
curl -s -X GET "http://localhost:8080/api/service/list?page=1&pageSize=5" > /dev/null
sleep 1

# 测试服务分类API
echo "测试服务分类API..."
curl -s -X GET "http://localhost:8080/api/service/categories" > /dev/null
sleep 1

# 测试首页初始化API
echo "测试首页初始化API..."
curl -s -X GET "http://localhost:8080/api/home/init" > /dev/null
sleep 1

# 停止应用
echo ""
echo "7. 停止应用..."
kill $APP_PID 2>/dev/null
wait $APP_PID 2>/dev/null

# 分析日志
echo ""
echo "8. 分析日志输出..."

if [ -f "$test_log_file" ]; then
    echo "日志文件大小: $(wc -l < "$test_log_file") 行"
    
    # 统计不同类型的日志
    echo ""
    echo "日志统计:"
    echo "SQL日志: $(grep -c '\[SQL\]' "$test_log_file" 2>/dev/null || echo 0)"
    echo "API日志: $(grep -c '\[API\]' "$test_log_file" 2>/dev/null || echo 0)"
    echo "错误日志: $(grep -c '\[ERROR\]' "$test_log_file" 2>/dev/null || echo 0)"
    echo "信息日志: $(grep -c '\[INFO\]' "$test_log_file" 2>/dev/null || echo 0)"
    
    # 显示最近的SQL日志
    echo ""
    echo "最近的SQL日志:"
    grep '\[SQL\]' "$test_log_file" | tail -5 || echo "未找到SQL日志"
    
    # 显示错误日志
    echo ""
    echo "错误日志:"
    grep '\[ERROR\]' "$test_log_file" | tail -3 || echo "未找到错误日志"
    
    # 清理测试文件
    echo ""
    echo "9. 清理测试文件..."
    rm -f "$test_log_file"
    
else
    echo "❌ 未找到测试日志文件"
fi

echo ""
echo "=== 测试完成 ==="
echo ""
echo "如果看到SQL日志输出，说明SQL日志系统工作正常"
echo "如果没有看到SQL日志，请检查："
echo "1. DAO文件是否正确添加了SQL日志"
echo "2. 应用是否正确启动"
echo "3. API调用是否成功"
