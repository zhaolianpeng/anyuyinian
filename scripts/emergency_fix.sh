#!/bin/bash

# 紧急修复脚本 - 解决"record not found"问题
# 使用方法: ./emergency_fix.sh [base_url]

BASE_URL=${1:-"https://prod-5g94mx7a3d07e78c.service.tcloudbase.com"}
echo "🚨 紧急修复UserId问题"
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

# 获取用户状态
get_user_status() {
    print_message $BLUE "📊 获取用户状态..."
    
    response=$(curl -s -X GET "$BASE_URL/api/emergency/user_status")
    echo "用户状态响应: $response"
    
    if echo "$response" | grep -q '"code":0'; then
        print_message $GREEN "✅ 用户状态获取成功"
        return 0
    else
        print_message $RED "❌ 用户状态获取失败"
        return 1
    fi
}

# 执行紧急修复
fix_user_ids() {
    print_message $BLUE "🔧 执行紧急修复..."
    
    response=$(curl -s -X POST "$BASE_URL/api/emergency/fix_user_ids")
    echo "修复响应: $response"
    
    if echo "$response" | grep -q '"code":0'; then
        print_message $GREEN "✅ 紧急修复成功"
        return 0
    else
        print_message $RED "❌ 紧急修复失败"
        return 1
    fi
}

# 测试用户信息API
test_user_info() {
    print_message $BLUE "🧪 测试用户信息API..."
    
    response=$(curl -s -X GET "$BASE_URL/api/emergency/test_user_info?userId=1")
    echo "用户信息API响应: $response"
    
    if echo "$response" | grep -q '"code":0'; then
        print_message $GREEN "✅ 用户信息API测试成功"
        return 0
    else
        print_message $YELLOW "⚠️  用户信息API仍然有问题"
        return 1
    fi
}

# 显示修复结果
show_results() {
    echo ""
    print_message $BLUE "📋 修复结果:"
    echo "1. ✅ 已执行数据库修复，为现有用户生成UserId"
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
    print_message $BLUE "🚀 开始紧急修复..."
    echo ""
    
    # 检查服务
    if ! check_service; then
        print_message $RED "❌ 服务不可用，请先启动服务"
        exit 1
    fi
    
    # 获取用户状态
    if ! get_user_status; then
        print_message $RED "❌ 无法获取用户状态"
        exit 1
    fi
    
    # 执行修复
    if ! fix_user_ids; then
        print_message $RED "❌ 修复失败"
        exit 1
    fi
    
    # 测试修复效果
    test_user_info
    
    # 显示结果
    show_results
    
    echo ""
    print_message $GREEN "🎉 紧急修复完成！"
}

# 运行主函数
main "$@" 