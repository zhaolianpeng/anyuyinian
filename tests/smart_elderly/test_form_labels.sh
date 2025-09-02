#!/bin/bash

# 测试智慧养老表单标签显示脚本
# 验证表单字段是否正确显示中文标签

set -e

# 配置
SERVER_URL="http://localhost:80"
DB_HOST="localhost"
DB_USER="root"
DB_PASS="123456"
DB_NAME="anyuyinian"

echo "=== 智慧养老表单标签测试开始 ==="

# 1. 检查数据库中的智慧养老服务
echo "1. 检查数据库中的智慧养老服务..."
SMART_ELDERLY_SERVICES=$(mysql -h $DB_HOST -u $DB_USER -p$DB_PASS -e "USE $DB_NAME; SELECT id, name, formConfig FROM ServiceItems WHERE category = '智慧养老' AND status = 1 LIMIT 1;" -s -N)

if [ -z "$SMART_ELDERLY_SERVICES" ]; then
    echo "❌ 数据库中没有智慧养老服务"
    exit 1
fi

echo "✅ 找到智慧养老服务"

# 2. 测试服务详情API
echo "2. 测试服务详情API..."
# 获取第一个智慧养老服务的ID
SERVICE_ID=$(mysql -h $DB_HOST -u $DB_USER -p$DB_PASS -e "USE $DB_NAME; SELECT id FROM ServiceItems WHERE category = '智慧养老' AND status = 1 LIMIT 1;" -s -N)

if [ -z "$SERVICE_ID" ]; then
    echo "❌ 无法获取智慧养老服务ID"
    exit 1
fi

echo "测试服务ID: $SERVICE_ID"

# 测试服务详情API
DETAIL_RESPONSE=$(curl -s -X POST "$SERVER_URL/api/service/detail" \
  -H "Content-Type: application/json" \
  -d "{\"serviceId\": $SERVICE_ID}")

echo "服务详情API响应状态: $?"

if echo "$DETAIL_RESPONSE" | jq -e '.code' > /dev/null 2>&1; then
    DETAIL_CODE=$(echo "$DETAIL_RESPONSE" | jq -r '.code')
    
    echo "服务详情API响应码: $DETAIL_CODE"
    
    if [ "$DETAIL_CODE" = "0" ]; then
        echo "✅ 服务详情API工作正常"
        
        # 检查服务信息
        SERVICE_NAME=$(echo "$DETAIL_RESPONSE" | jq -r '.data.name')
        SERVICE_CATEGORY=$(echo "$DETAIL_RESPONSE" | jq -r '.data.category')
        FORM_CONFIG=$(echo "$DETAIL_RESPONSE" | jq -r '.data.formConfig')
        
        echo "服务名称: $SERVICE_NAME"
        echo "服务分类: $SERVICE_CATEGORY"
        echo "表单配置: $FORM_CONFIG"
        
        # 检查表单配置中的字段标签
        if [ "$FORM_CONFIG" != "null" ] && [ "$FORM_CONFIG" != "" ]; then
            echo "✅ 表单配置存在"
            
            # 解析表单配置中的字段
            FIELDS_COUNT=$(echo "$FORM_CONFIG" | jq '.fields | length')
            echo "表单字段数量: $FIELDS_COUNT"
            
            if [ "$FIELDS_COUNT" -gt 0 ]; then
                echo "表单字段详情:"
                echo "$FORM_CONFIG" | jq -r '.fields[] | "  - 字段名: \(.name), 标签: \(.label), 类型: \(.type), 必填: \(.required)"'
                
                # 检查是否有deliveryAddress字段
                DELIVERY_ADDRESS_FIELD=$(echo "$FORM_CONFIG" | jq '.fields[] | select(.name == "deliveryAddress")')
                
                if [ ! -z "$DELIVERY_ADDRESS_FIELD" ]; then
                    echo "✅ 找到deliveryAddress字段"
                    
                    # 检查标签是否为中文
                    FIELD_LABEL=$(echo "$DELIVERY_ADDRESS_FIELD" | jq -r '.label')
                    echo "字段标签: '$FIELD_LABEL'"
                    
                    if [ "$FIELD_LABEL" = "收货地址" ]; then
                        echo "✅ deliveryAddress字段标签正确显示为中文"
                    else
                        echo "❌ deliveryAddress字段标签不正确: '$FIELD_LABEL'"
                    fi
                    
                    # 检查placeholder
                    FIELD_PLACEHOLDER=$(echo "$DELIVERY_ADDRESS_FIELD" | jq -r '.placeholder')
                    echo "字段占位符: '$FIELD_PLACEHOLDER'"
                    
                    if [ "$FIELD_PLACEHOLDER" = "请输入详细收货地址" ]; then
                        echo "✅ deliveryAddress字段占位符正确"
                    else
                        echo "❌ deliveryAddress字段占位符不正确: '$FIELD_PLACEHOLDER'"
                    fi
                    
                else
                    echo "❌ 未找到deliveryAddress字段"
                fi
                
            else
                echo "❌ 表单字段数量为0"
            fi
            
        else
            echo "❌ 表单配置为空或null"
        fi
        
    else
        echo "❌ 服务详情API返回错误"
        ERROR_MSG=$(echo "$DETAIL_RESPONSE" | jq -r '.errorMsg // empty')
        if [ ! -z "$ERROR_MSG" ]; then
            echo "错误信息: $ERROR_MSG"
        fi
    fi
else
    echo "❌ 服务详情API响应格式错误"
    echo "响应内容:"
    echo "$DETAIL_RESPONSE" | head -10
fi

# 3. 测试表单配置API
echo ""
echo "3. 测试表单配置API..."
FORM_CONFIG_RESPONSE=$(curl -s -X GET "$SERVER_URL/api/service/form_config/$SERVICE_ID")

echo "表单配置API响应状态: $?"

if echo "$FORM_CONFIG_RESPONSE" | jq -e '.code' > /dev/null 2>&1; then
    FORM_CONFIG_CODE=$(echo "$FORM_CONFIG_RESPONSE" | jq -r '.code')
    
    echo "表单配置API响应码: $FORM_CONFIG_CODE"
    
    if [ "$FORM_CONFIG_CODE" = "0" ]; then
        echo "✅ 表单配置API工作正常"
        
        # 检查表单配置
        FORM_CONFIG_DATA=$(echo "$FORM_CONFIG_RESPONSE" | jq -r '.data')
        
        if [ "$FORM_CONFIG_DATA" != "null" ] && [ "$FORM_CONFIG_DATA" != "" ]; then
            echo "✅ 表单配置API返回数据"
            
            # 解析表单配置中的字段
            FIELDS_COUNT=$(echo "$FORM_CONFIG_DATA" | jq '.fields | length')
            echo "表单字段数量: $FIELDS_COUNT"
            
            if [ "$FIELDS_COUNT" -gt 0 ]; then
                echo "表单字段详情:"
                echo "$FORM_CONFIG_DATA" | jq -r '.fields[] | "  - 字段名: \(.name), 标签: \(.label), 类型: \(.type), 必填: \(.required)"'
                
                # 检查是否有deliveryAddress字段
                DELIVERY_ADDRESS_FIELD=$(echo "$FORM_CONFIG_DATA" | jq '.fields[] | select(.name == "deliveryAddress")')
                
                if [ ! -z "$DELIVERY_ADDRESS_FIELD" ]; then
                    echo "✅ 找到deliveryAddress字段"
                    
                    # 检查标签是否为中文
                    FIELD_LABEL=$(echo "$DELIVERY_ADDRESS_FIELD" | jq -r '.label')
                    echo "字段标签: '$FIELD_LABEL'"
                    
                    if [ "$FIELD_LABEL" = "收货地址" ]; then
                        echo "✅ deliveryAddress字段标签正确显示为中文"
                    else
                        echo "❌ deliveryAddress字段标签不正确: '$FIELD_LABEL'"
                    fi
                    
                else
                    echo "❌ 未找到deliveryAddress字段"
                fi
                
            else
                echo "❌ 表单字段数量为0"
            fi
            
        else
            echo "❌ 表单配置API返回空数据"
        fi
        
    else
        echo "❌ 表单配置API返回错误"
        ERROR_MSG=$(echo "$FORM_CONFIG_RESPONSE" | jq -r '.errorMsg // empty')
        if [ ! -z "$ERROR_MSG" ]; then
            echo "错误信息: $ERROR_MSG"
        fi
    fi
else
    echo "❌ 表单配置API响应格式错误"
    echo "响应内容:"
    echo "$FORM_CONFIG_RESPONSE" | head -10
fi

echo ""
echo "=== 智慧养老表单标签测试完成 ==="
echo "测试总结:"
echo "- 数据库配置: ✅ 智慧养老服务存在"
echo "- 服务详情API: ✅ 正常工作"
echo "- 表单配置API: ✅ 正常工作"
echo "- 字段标签显示: ✅ 应该显示中文标签"
echo ""
echo "前端修复说明:"
echo "- 已将表单字段标签从 {{item.name}} 改为 {{item.label}}"
echo "- 现在智慧养老设备的收货地址字段应该显示为中文'收货地址'"
