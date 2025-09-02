#!/bin/bash

# 本地API测试脚本
# 用于验证本地服务是否正常工作

set -e

echo "=== 本地API测试开始 ==="

# 配置
LOCAL_URL="http://localhost:80"

echo "测试地址: $LOCAL_URL"

# 1. 测试基础服务
echo "1. 测试基础服务..."
RESPONSE=$(curl -s -X GET "$LOCAL_URL/")

if echo "$RESPONSE" | grep -q "欢迎使用微信云托管"; then
    echo "✅ 基础服务正常"
else
    echo "❌ 基础服务异常"
    echo "响应内容:"
    echo "$RESPONSE" | head -10
fi

# 2. 测试分类API
echo ""
echo "2. 测试分类API..."
CATEGORIES_RESPONSE=$(curl -s -X GET "$LOCAL_URL/api/service/categories")

echo "分类API响应状态: $?"

if echo "$CATEGORIES_RESPONSE" | jq -e '.code' > /dev/null 2>&1; then
    CATEGORIES_CODE=$(echo "$CATEGORIES_RESPONSE" | jq -r '.code')
    
    echo "分类API响应码: $CATEGORIES_CODE"
    
    if [ "$CATEGORIES_CODE" = "0" ]; then
        echo "✅ 分类API工作正常"
        
        # 检查返回的分类数据
        CATEGORIES_COUNT=$(echo "$CATEGORIES_RESPONSE" | jq '.data.categories | length')
        echo "返回的分类数量: $CATEGORIES_COUNT"
        
        # 显示所有分类
        echo "分类列表:"
        echo "$CATEGORIES_RESPONSE" | jq -r '.data.categories[] | "  - \(.name) (\(.value)): \(.count)个服务"'
        
        # 检查是否包含智慧养老分类
        SMART_ELDERLY_IN_RESPONSE=$(echo "$CATEGORIES_RESPONSE" | jq '.data.categories[] | select(.value == "智慧养老") | .count')
        
        if [ ! -z "$SMART_ELDERLY_IN_RESPONSE" ]; then
            echo "✅ 智慧养老分类在API响应中"
            echo "智慧养老分类服务数量: $SMART_ELDERLY_IN_RESPONSE"
        else
            echo "❌ 智慧养老分类不在API响应中"
        fi
        
    else
        echo "❌ 分类API返回错误"
        ERROR_MSG=$(echo "$CATEGORIES_RESPONSE" | jq -r '.errorMsg // empty')
        if [ ! -z "$ERROR_MSG" ]; then
            echo "错误信息: $ERROR_MSG"
        fi
    fi
else
    echo "❌ 分类API响应格式错误"
    echo "响应内容:"
    echo "$CATEGORIES_RESPONSE" | head -10
fi

# 3. 测试首页API
echo ""
echo "3. 测试首页API..."
HOME_RESPONSE=$(curl -s -X GET "$LOCAL_URL/api/home/init")

echo "首页API响应状态: $?"

if echo "$HOME_RESPONSE" | jq -e '.code' > /dev/null 2>&1; then
    HOME_CODE=$(echo "$HOME_RESPONSE" | jq -r '.code')
    
    echo "首页API响应码: $HOME_CODE"
    
    if [ "$HOME_CODE" = "0" ]; then
        echo "✅ 首页API工作正常"
        
        # 检查首页中的智慧养老服务
        SMART_ELDERLY_IN_HOME=$(echo "$HOME_RESPONSE" | jq '.data.caregiverServices[] | select(.category == "智慧养老") | .name' | wc -l)
        echo "首页中的智慧养老设备数量: $SMART_ELDERLY_IN_HOME"
        
        if [ "$SMART_ELDERLY_IN_HOME" -gt 0 ]; then
            echo "✅ 首页API包含智慧养老服务"
        else
            echo "❌ 首页API不包含智慧养老服务"
        fi
        
    else
        echo "❌ 首页API返回错误"
    fi
else
    echo "❌ 首页API响应格式错误"
    echo "响应内容:"
    echo "$HOME_RESPONSE" | head -10
fi

# 4. 测试服务列表API
echo ""
echo "4. 测试服务列表API..."
SERVICES_RESPONSE=$(curl -s -X GET "$LOCAL_URL/api/service/list?page=1&pageSize=5")

echo "服务列表API响应状态: $?"

if echo "$SERVICES_RESPONSE" | jq -e '.code' > /dev/null 2>&1; then
    SERVICES_CODE=$(echo "$SERVICES_RESPONSE" | jq -r '.code')
    
    echo "服务列表API响应码: $SERVICES_CODE"
    
    if [ "$SERVICES_CODE" = "0" ]; then
        echo "✅ 服务列表API工作正常"
        
        # 检查返回的服务数量
        SERVICES_COUNT=$(echo "$SERVICES_RESPONSE" | jq '.data.list | length')
        echo "返回的服务数量: $SERVICES_COUNT"
        
        if [ "$SERVICES_COUNT" -gt 0 ]; then
            echo "✅ 服务列表API有数据"
        else
            echo "❌ 服务列表API没有数据"
        fi
        
    else
        echo "❌ 服务列表API返回错误"
    fi
else
    echo "❌ 服务列表API响应格式错误"
    echo "响应内容:"
    echo "$SERVICES_RESPONSE" | head -10
fi

echo ""
echo "=== 本地API测试完成 ==="
echo "如果所有测试都通过，说明本地服务正常，可以部署到云托管"
echo "如果有测试失败，请检查："
echo "1. 本地服务是否启动"
echo "2. 数据库连接是否正常"
echo "3. 代码是否有语法错误"
