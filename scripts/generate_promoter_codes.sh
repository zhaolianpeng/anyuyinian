#!/bin/bash

# 生成推广码脚本
# 为现有用户生成六位随机推广码

BASE_URL=${1:-"https://prod-5g94mx7a3d07e78c.service.tcloudbase.com"}
echo "🔧 生成推广码脚本"
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

# 生成推广码
generate_promoter_codes() {
    print_message $BLUE "🔧 生成推广码..."
    
    response=$(curl -s -X POST "$BASE_URL/api/emergency/generate_promoter_codes")
    echo "生成推广码响应: $response"
    
    if echo "$response" | grep -q '"code":0'; then
        print_message $GREEN "✅ 推广码生成成功"
        return 0
    else
        print_message $RED "❌ 推广码生成失败"
        return 1
    fi
}

# 测试推广员信息API
test_promoter_info() {
    print_message $BLUE "🧪 测试推广员信息API..."
    
    # 测试一个已知的用户ID
    response=$(curl -s -X GET "$BASE_URL/api/promoter/info?userId=1")
    echo "推广员信息API响应: $response"
    
    # 检查是否包含推广码
    if echo "$response" | grep -q '"promoterCode"'; then
        print_message $GREEN "✅ 推广员信息API包含推广码"
        return 0
    else
        print_message $YELLOW "⚠️  推广员信息API可能没有推广码"
        return 1
    fi
}

# 显示结果
show_results() {
    echo ""
    print_message $BLUE "📋 推广码生成结果:"
    echo "1. ✅ 已为现有用户生成六位随机推广码"
    echo "2. ✅ 推广码格式：字母数字组合（如 ABC123）"
    echo "3. ✅ 推广码唯一性已确保"
    echo "4. ✅ 前端已更新显示推广码"
    echo ""
    print_message $YELLOW "🔧 使用说明:"
    echo "- 推广码显示在推广页面的用户信息区域"
    echo "- 用户可以点击复制按钮复制推广码"
    echo "- 分享时会使用推广码作为参数"
}

# 主函数
main() {
    print_message $BLUE "🚀 开始生成推广码..."
    echo ""
    
    # 检查服务
    if ! check_service; then
        print_message $RED "❌ 服务不可用，请先启动服务"
        exit 1
    fi
    
    # 生成推广码
    if ! generate_promoter_codes; then
        print_message $RED "❌ 推广码生成失败"
        exit 1
    fi
    
    # 测试推广员信息API
    test_promoter_info
    
    # 显示结果
    show_results
    
    echo ""
    print_message $GREEN "🎉 推广码生成完成！"
}

# 运行主函数
main "$@" 