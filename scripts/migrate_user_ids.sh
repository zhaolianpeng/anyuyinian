#!/bin/bash

# 数据迁移脚本 - 为现有用户生成UserId
# 使用方法: ./migrate_user_ids.sh [base_url]

BASE_URL=${1:-"http://localhost:80"}
echo "使用基础URL: $BASE_URL"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_message() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# 测试API连接
test_connection() {
    print_message $BLUE "🔍 测试API连接..."
    
    response=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/api/count")
    if [ "$response" = "200" ]; then
        print_message $GREEN "✅ API连接成功"
        return 0
    else
        print_message $RED "❌ API连接失败 (HTTP $response)"
        return 1
    fi
}

# 生成UserId测试
test_generate_user_id() {
    print_message $BLUE "🧪 测试UserId生成..."
    
    response=$(curl -s -X POST "$BASE_URL/api/migration/generate_user_ids")
    echo "响应: $response"
    
    # 检查响应
    if echo "$response" | grep -q '"code":0'; then
        print_message $GREEN "✅ UserId生成测试成功"
        return 0
    else
        print_message $RED "❌ UserId生成测试失败"
        return 1
    fi
}

# 迁移用户UserId
migrate_users() {
    print_message $BLUE "🔄 开始迁移用户UserId..."
    
    response=$(curl -s -X POST "$BASE_URL/api/migration/migrate_users")
    echo "响应: $response"
    
    # 检查响应
    if echo "$response" | grep -q '"code":0'; then
        print_message $GREEN "✅ 用户UserId迁移成功"
        return 0
    else
        print_message $RED "❌ 用户UserId迁移失败"
        return 1
    fi
}

# 迁移所有表的UserId
migrate_all_tables() {
    print_message $BLUE "🔄 开始迁移所有表的UserId..."
    
    response=$(curl -s -X POST "$BASE_URL/api/migration/migrate_all_tables")
    echo "响应: $response"
    
    # 检查响应
    if echo "$response" | grep -q '"code":0'; then
        print_message $GREEN "✅ 所有表UserId迁移成功"
        return 0
    else
        print_message $RED "❌ 所有表UserId迁移失败"
        return 1
    fi
}

# 验证UserId
validate_user_ids() {
    print_message $BLUE "🔍 验证UserId..."
    
    response=$(curl -s -X GET "$BASE_URL/api/migration/validate")
    echo "响应: $response"
    
    # 检查响应
    if echo "$response" | grep -q '"code":0'; then
        print_message $GREEN "✅ UserId验证成功"
        return 0
    else
        print_message $RED "❌ UserId验证失败"
        return 1
    fi
}

# 测试用户信息API
test_user_info_api() {
    print_message $BLUE "🧪 测试用户信息API..."
    
    # 测试一个已知的用户ID
    response=$(curl -s -X GET "$BASE_URL/api/user/info?userId=1")
    echo "用户信息API响应: $response"
    
    # 检查是否还有"record not found"错误
    if echo "$response" | grep -q "record not found"; then
        print_message $YELLOW "⚠️  用户信息API仍然返回'record not found'，可能需要重新登录"
        return 1
    else
        print_message $GREEN "✅ 用户信息API测试成功"
        return 0
    fi
}

# 主函数
main() {
    print_message $BLUE "🚀 开始数据迁移流程..."
    echo ""
    
    # 测试连接
    if ! test_connection; then
        print_message $RED "❌ 无法连接到API，请检查服务是否运行"
        exit 1
    fi
    
    # 测试UserId生成
    if ! test_generate_user_id; then
        print_message $RED "❌ UserId生成测试失败"
        exit 1
    fi
    
    # 迁移用户UserId
    if ! migrate_users; then
        print_message $RED "❌ 用户UserId迁移失败"
        exit 1
    fi
    
    # 迁移所有表的UserId
    if ! migrate_all_tables; then
        print_message $RED "❌ 所有表UserId迁移失败"
        exit 1
    fi
    
    # 验证UserId
    if ! validate_user_ids; then
        print_message $RED "❌ UserId验证失败"
        exit 1
    fi
    
    # 测试用户信息API
    test_user_info_api
    
    echo ""
    print_message $GREEN "🎉 数据迁移完成！"
    print_message $YELLOW "💡 提示: 如果用户信息API仍然返回错误，用户可能需要重新登录以获取新的UserId"
}

# 运行主函数
main "$@" 