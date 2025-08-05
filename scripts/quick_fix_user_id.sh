#!/bin/bash

# 快速修复脚本 - 解决"record not found"问题
# 使用方法: ./quick_fix_user_id.sh [base_url]

BASE_URL=${1:-"http://localhost:80"}
echo "🔧 快速修复UserId问题"
echo "使用基础URL: $BASE_URL"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_message() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# 检查服务状态
check_service() {
    print_message $BLUE "🔍 检查服务状态..."
    
    response=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/api/count")
    if [ "$response" = "200" ]; then
        print_message $GREEN "✅ 服务运行正常"
        return 0
    else
        print_message $RED "❌ 服务连接失败 (HTTP $response)"
        return 1
    fi
}

# 执行用户迁移
migrate_users() {
    print_message $BLUE "🔄 执行用户UserId迁移..."
    
    response=$(curl -s -X POST "$BASE_URL/api/migration/migrate_users")
    echo "迁移响应: $response"
    
    if echo "$response" | grep -q '"code":0'; then
        print_message $GREEN "✅ 用户迁移成功"
        return 0
    else
        print_message $RED "❌ 用户迁移失败"
        return 1
    fi
}

# 测试修复效果
test_fix() {
    print_message $BLUE "🧪 测试修复效果..."
    
    # 测试用户信息API
    response=$(curl -s -X GET "$BASE_URL/api/user/info?userId=1")
    echo "用户信息API响应: $response"
    
    if echo "$response" | grep -q "record not found"; then
        print_message $YELLOW "⚠️  仍然存在'record not found'错误"
        print_message $YELLOW "💡 建议: 用户需要重新登录以获取新的UserId"
        return 1
    else
        print_message $GREEN "✅ 修复成功！用户信息API正常"
        return 0
    fi
}

# 显示解决方案
show_solution() {
    echo ""
    print_message $BLUE "📋 问题解决方案:"
    echo "1. ✅ 已执行数据库迁移，为现有用户生成UserId"
    echo "2. ⚠️  前端仍在使用旧的数字userId (如 '1')"
    echo "3. 💡 解决方案:"
    echo "   - 用户需要重新登录"
    echo "   - 或者清除本地存储的userId"
    echo "   - 或者等待前端自动处理userId兼容性"
    echo ""
    print_message $YELLOW "🔧 前端兼容性处理:"
    echo "- 前端已实现userId兼容性处理"
    echo "- 会自动检测旧格式userId并提示重新登录"
    echo "- 或者自动清除本地存储"
}

# 主函数
main() {
    print_message $BLUE "🚀 开始快速修复..."
    echo ""
    
    # 检查服务
    if ! check_service; then
        print_message $RED "❌ 服务不可用，请先启动服务"
        exit 1
    fi
    
    # 执行迁移
    if ! migrate_users; then
        print_message $RED "❌ 迁移失败"
        exit 1
    fi
    
    # 测试修复效果
    test_fix
    
    # 显示解决方案
    show_solution
    
    echo ""
    print_message $GREEN "🎉 快速修复完成！"
}

# 运行主函数
main "$@" 