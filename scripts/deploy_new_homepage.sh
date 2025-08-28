#!/bin/bash

# 护工服务平台首页部署脚本（使用现有Services表）

echo "=== 护工服务平台首页部署脚本 ==="
echo ""
echo "注意：此部署使用现有的Services和ServiceItems表，不创建新表"
echo ""

# 检查是否在正确的目录
if [ ! -f "main.go" ]; then
    echo "❌ 错误：请在项目根目录运行此脚本"
    exit 1
fi

echo "✅ 检测到项目根目录"
echo ""

# 1. 数据库迁移
echo "=== 步骤1: 执行数据库迁移 ==="
if [ -f "db/migration/20241220_create_caregiver_services.sql" ]; then
    echo "发现数据库迁移文件，此文件将："
    echo "- 为Services表添加price和category字段"
    echo "- 为ServiceItems表添加duration和features字段"
    echo "- 插入护工服务数据到现有表中"
    echo ""
    echo "请手动执行以下SQL："
    echo "mysql -u your_username -p your_database < db/migration/20241220_create_caregiver_services.sql"
    echo ""
    read -p "数据库迁移是否已完成？(y/n): " db_migrated
    
    if [ "$db_migrated" != "y" ] && [ "$db_migrated" != "Y" ]; then
        echo "请先完成数据库迁移，然后重新运行此脚本"
        exit 1
    fi
else
    echo "❌ 未找到数据库迁移文件"
    exit 1
fi

# 2. 编译后端
echo "=== 步骤2: 编译后端服务 ==="
echo "正在编译Go程序..."
if go build -o main .; then
    echo "✅ 后端编译成功"
else
    echo "❌ 后端编译失败"
    exit 1
fi

# 3. 检查小程序代码
echo "=== 步骤3: 检查小程序代码 ==="
if [ -d "miniprogram" ]; then
    echo "✅ 检测到小程序代码目录"
    
    # 检查关键文件
    required_files=(
        "miniprogram/pages/index/index.wxml"
        "miniprogram/pages/index/index.wxss"
        "miniprogram/pages/index/index.js"
    )
    
    for file in "${required_files[@]}"; do
        if [ -f "$file" ]; then
            echo "✅ $file 存在"
        else
            echo "❌ $file 缺失"
            exit 1
        fi
    done
else
    echo "❌ 未找到小程序代码目录"
    exit 1
fi

# 4. 启动后端服务
echo "=== 步骤4: 启动后端服务 ==="
echo "正在启动后端服务..."
./main &
BACKEND_PID=$!

# 等待服务启动
sleep 3

# 检查服务是否启动成功
if curl -s http://localhost:80/api/home/init > /dev/null; then
    echo "✅ 后端服务启动成功"
else
    echo "❌ 后端服务启动失败"
    kill $BACKEND_PID 2>/dev/null
    exit 1
fi

# 5. 测试API接口
echo "=== 步骤5: 测试API接口 ==="
echo "测试首页初始化接口..."
response=$(curl -s http://localhost:80/api/home/init)

if echo "$response" | jq -e '.data.caregiverServices' > /dev/null 2>&1; then
    echo "✅ 护工服务数据接口正常"
else
    echo "❌ 护工服务数据接口异常"
fi

if echo "$response" | jq -e '.data.platformGuarantees' > /dev/null 2>&1; then
    echo "✅ 平台保障数据接口正常"
else
    echo "❌ 平台保障数据接口异常"
fi

if echo "$response" | jq -e '.data.companyInfo' > /dev/null 2>&1; then
    echo "✅ 公司信息数据接口正常"
else
    echo "❌ 公司信息数据接口异常"
fi

# 6. 部署说明
echo ""
echo "=== 部署完成 ==="
echo ""
echo "🎉 护工服务平台首页部署成功！"
echo ""
echo "📱 小程序端："
echo "1. 在微信开发者工具中打开 miniprogram 目录"
echo "2. 上传代码到微信小程序后台"
echo "3. 提交审核并发布"
echo ""
echo "🔧 后端服务："
echo "- 服务已启动，PID: $BACKEND_PID"
echo "- 服务地址: http://localhost:80"
echo "- 首页API: http://localhost:80/api/home/init"
echo ""
echo "🗄️  数据库变更："
echo "- Services表：添加了price和category字段"
echo "- ServiceItems表：添加了duration和features字段"
echo "- 插入了5条护工服务数据，category为'护工服务'"
echo ""
echo "📊 测试验证："
echo "- 运行测试脚本: tests/homepage/test_new_homepage.sh"
echo "- 检查Services表的护工服务数据"
echo "- 验证category字段和价格信息"
echo ""
echo "⚠️  注意事项："
echo "1. 确保数据库连接配置正确"
echo "2. 检查图片资源路径是否存在"
echo "3. 验证所有API接口返回正常"
echo "4. 测试小程序页面显示效果"
echo "5. 确认Services表中的护工服务数据完整"
echo ""
echo "🛑 停止服务："
echo "kill $BACKEND_PID"
echo ""
echo "部署完成！如有问题，请查看日志或联系技术支持。"
