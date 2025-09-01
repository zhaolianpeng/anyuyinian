#!/bin/bash

# 智慧养老服务测试脚本
# 测试智慧养老设备的创建、查询和订单提交功能

set -e

# 配置
SERVER_URL="http://localhost:80"
DB_HOST="localhost"
DB_USER="root"
DB_PASS="123456"
DB_NAME="anyuyinian"

echo "=== 智慧养老服务测试开始 ==="

# 1. 测试数据库迁移
echo "1. 检查数据库迁移..."
mysql -h $DB_HOST -u $DB_USER -p$DB_PASS -e "USE $DB_NAME; DESCRIBE ServiceItems;" | grep -q videoUrl && echo "✅ videoUrl字段已存在" || echo "❌ videoUrl字段不存在"

# 2. 测试智慧养老服务数据
echo "2. 检查智慧养老服务数据..."
SMART_ELDERLY_COUNT=$(mysql -h $DB_HOST -u $DB_USER -p$DB_PASS -e "USE $DB_NAME; SELECT COUNT(*) FROM ServiceItems WHERE category = '智慧养老';" -s -N)
echo "智慧养老设备数量: $SMART_ELDERLY_COUNT"

if [ "$SMART_ELDERLY_COUNT" -gt 0 ]; then
    echo "✅ 智慧养老服务数据已存在"
    
    # 显示智慧养老设备列表
    echo "智慧养老设备列表:"
    mysql -h $DB_HOST -u $DB_USER -p$DB_PASS -e "USE $DB_NAME; SELECT id, name, price, videoUrl FROM ServiceItems WHERE category = '智慧养老' ORDER BY sort;" | column -t
else
    echo "❌ 智慧养老服务数据不存在"
fi

# 3. 测试首页API中的智慧养老服务
echo "3. 测试首页API..."
HOME_RESPONSE=$(curl -s -X GET "$SERVER_URL/api/home/init")
echo "首页API响应状态: $?"

if echo "$HOME_RESPONSE" | jq -e '.data.caregiverServices' > /dev/null; then
    SMART_ELDERLY_IN_HOME=$(echo "$HOME_RESPONSE" | jq '.data.caregiverServices[] | select(.category == "智慧养老") | .name' | wc -l)
    echo "首页中的智慧养老设备数量: $SMART_ELDERLY_IN_HOME"
    
    if [ "$SMART_ELDERLY_IN_HOME" -gt 0 ]; then
        echo "✅ 首页API包含智慧养老服务"
        echo "智慧养老设备名称:"
        echo "$HOME_RESPONSE" | jq -r '.data.caregiverServices[] | select(.category == "智慧养老") | .name'
    else
        echo "❌ 首页API不包含智慧养老服务"
    fi
else
    echo "❌ 首页API响应格式错误"
fi

# 4. 测试智慧养老设备详情API
echo "4. 测试智慧养老设备详情API..."
if [ "$SMART_ELDERLY_COUNT" -gt 0 ]; then
    # 获取第一个智慧养老设备的ID
    FIRST_SMART_ELDERLY_ID=$(mysql -h $DB_HOST -u $DB_USER -p$DB_PASS -e "USE $DB_NAME; SELECT id FROM ServiceItems WHERE category = '智慧养老' ORDER BY sort LIMIT 1;" -s -N)
    
    if [ ! -z "$FIRST_SMART_ELDERLY_ID" ]; then
        DETAIL_RESPONSE=$(curl -s -X POST "$SERVER_URL/api/service/detail" \
            -H "Content-Type: application/json" \
            -d "{\"serviceId\": $FIRST_SMART_ELDERLY_ID}")
        
        echo "设备详情API响应状态: $?"
        
        if echo "$DETAIL_RESPONSE" | jq -e '.data' > /dev/null; then
            DEVICE_CATEGORY=$(echo "$DETAIL_RESPONSE" | jq -r '.data.category')
            DEVICE_VIDEO_URL=$(echo "$DETAIL_RESPONSE" | jq -r '.data.videoUrl // empty')
            
            echo "设备分类: $DEVICE_CATEGORY"
            echo "视频URL: $DEVICE_VIDEO_URL"
            
            if [ "$DEVICE_CATEGORY" = "智慧养老" ]; then
                echo "✅ 设备详情API正确返回智慧养老设备"
            else
                echo "❌ 设备详情API返回的分类不正确"
            fi
        else
            echo "❌ 设备详情API响应格式错误"
        fi
    else
        echo "❌ 无法获取智慧养老设备ID"
    fi
else
    echo "跳过设备详情测试（无智慧养老设备数据）"
fi

# 5. 测试智慧养老设备订单API
echo "5. 测试智慧养老设备订单API..."
SMART_ELDERLY_ORDER_RESPONSE=$(curl -s -X POST "$SERVER_URL/api/order/smart-elderly" \
    -H "Content-Type: application/json" \
    -d '{
        "userId": "test_user_123",
        "serviceId": 1,
        "addressId": 1,
        "quantity": 1,
        "formData": {
            "deliveryAddress": "测试收货地址"
        },
        "remark": "测试智慧养老设备订单"
    }')

echo "智慧养老订单API响应状态: $?"

if echo "$SMART_ELDERLY_ORDER_RESPONSE" | jq -e '.code' > /dev/null; then
    ORDER_CODE=$(echo "$SMART_ELDERLY_ORDER_RESPONSE" | jq -r '.code')
    ORDER_MESSAGE=$(echo "$SMART_ELDERLY_ORDER_RESPONSE" | jq -r '.errorMsg // empty')
    
    echo "订单API响应码: $ORDER_CODE"
    if [ ! -z "$ORDER_MESSAGE" ]; then
        echo "错误信息: $ORDER_MESSAGE"
    fi
    
    if [ "$ORDER_CODE" = "0" ]; then
        echo "✅ 智慧养老设备订单API工作正常"
        ORDER_NO=$(echo "$SMART_ELDERLY_ORDER_RESPONSE" | jq -r '.data.orderNo')
        echo "生成的订单号: $ORDER_NO"
    else
        echo "❌ 智慧养老设备订单API返回错误"
    fi
else
    echo "❌ 智慧养老设备订单API响应格式错误"
fi

echo "=== 智慧养老服务测试完成 ==="
