#!/bin/bash

# 护工服务平台首页功能测试脚本（使用现有Services表）

# 设置测试服务器地址
SERVER_URL="http://localhost:80"

echo "=== 护工服务平台首页功能测试 ==="
echo "服务器地址: $SERVER_URL"
echo ""

# 测试1: 获取首页数据（包含护工服务）
echo "=== 测试1: 获取首页数据（包含护工服务） ==="
curl -X GET "$SERVER_URL/api/home/init" | jq '.'
echo ""
echo ""

# 测试2: 验证护工服务数据结构
echo "=== 测试2: 验证护工服务数据结构 ==="
response=$(curl -s -X GET "$SERVER_URL/api/home/init")
if echo "$response" | jq -e '.data.caregiverServices' > /dev/null; then
    echo "✅ 护工服务数据存在"
    caregiver_count=$(echo "$response" | jq '.data.caregiverServices | length')
    echo "护工服务数量: $caregiver_count"
else
    echo "❌ 护工服务数据缺失"
fi
echo ""

# 测试3: 验证平台保障数据结构
echo "=== 测试3: 验证平台保障数据结构 ==="
if echo "$response" | jq -e '.data.platformGuarantees' > /dev/null; then
    echo "✅ 平台保障数据存在"
    guarantee_count=$(echo "$response" | jq '.data.platformGuarantees | length')
    echo "平台保障数量: $guarantee_count"
else
    echo "❌ 平台保障数据缺失"
fi
echo ""

# 测试4: 验证公司信息数据
echo "=== 测试4: 验证公司信息数据 ==="
if echo "$response" | jq -e '.data.companyInfo' > /dev/null; then
    echo "✅ 公司信息数据存在"
    company_name=$(echo "$response" | jq -r '.data.companyInfo.name')
    echo "公司名称: $company_name"
else
    echo "❌ 公司信息数据缺失"
fi
echo ""

# 测试5: 验证护工服务详细内容
echo "=== 测试5: 验证护工服务详细内容 ==="
if echo "$response" | jq -e '.data.caregiverServices[0]' > /dev/null; then
    first_service=$(echo "$response" | jq '.data.caregiverServices[0]')
    echo "第一个服务详情:"
    echo "$first_service" | jq '.'
else
    echo "❌ 护工服务详细数据缺失"
fi
echo ""

# 测试6: 验证服务价格信息
echo "=== 测试6: 验证服务价格信息 ==="
services_with_price=$(echo "$response" | jq '.data.caregiverServices[] | select(.price != null)')
if [ ! -z "$services_with_price" ]; then
    echo "✅ 服务价格信息完整"
    echo "$services_with_price" | jq '{name: .name, price: .price, category: .category}'
else
    echo "❌ 服务价格信息缺失"
fi
echo ""

# 测试7: 验证服务分类
echo "=== 测试7: 验证服务分类 ==="
categories=$(echo "$response" | jq -r '.data.caregiverServices[].category' | sort | uniq)
echo "服务分类:"
echo "$categories"
echo ""

# 测试8: 验证平台保障图标
echo "=== 测试8: 验证平台保障图标 ==="
guarantees_with_icons=$(echo "$response" | jq '.data.platformGuarantees[] | select(.icon != null)')
if [ ! -z "$guarantees_with_icons" ]; then
    echo "✅ 平台保障图标完整"
    echo "$guarantees_with_icons" | jq '{title: .title, icon: .icon}'
else
    echo "❌ 平台保障图标缺失"
fi
echo ""

# 测试9: 验证公司发展历程
echo "=== 测试9: 验证公司发展历程 ==="
if echo "$response" | jq -e '.data.companyInfo.timeline' > /dev/null; then
    echo "✅ 公司发展历程数据存在"
    timeline=$(echo "$response" | jq '.data.companyInfo.timeline')
    echo "发展历程:"
    echo "$timeline" | jq '.'
else
    echo "❌ 公司发展历程数据缺失"
fi
echo ""

# 测试10: 验证集团背书信息
echo "=== 测试10: 验证集团背书信息 ==="
if echo "$response" | jq -e '.data.companyInfo.endorsements' > /dev/null; then
    echo "✅ 集团背书数据存在"
    endorsements=$(echo "$response" | jq '.data.companyInfo.endorsements')
    echo "集团背书:"
    echo "$endorsements" | jq '.'
else
    echo "❌ 集团背书数据缺失"
fi
echo ""

# 测试11: 验证Services表数据
echo "=== 测试11: 验证Services表护工服务数据 ==="
echo "检查数据库中Services表的护工服务数据..."
echo "请手动执行以下SQL查询："
echo "SELECT id, name, description, price, category, status FROM Services WHERE category = '护工服务';"
echo ""

# 测试12: 验证ServiceItems表数据
echo "=== 测试12: 验证ServiceItems表护工服务数据 ==="
echo "检查数据库中ServiceItems表的护工服务数据..."
echo "请手动执行以下SQL查询："
echo "SELECT id, name, description, price, category, duration, features, status FROM ServiceItems WHERE category = '护工服务';"
echo ""

echo "=== 测试完成 ==="
echo ""
echo "注意事项："
echo "1. 确保已执行数据库迁移文件: 20241220_create_caregiver_services.sql"
echo "2. 确保Services表已添加price和category字段"
echo "3. 确保ServiceItems表已添加duration和features字段"
echo "4. 验证护工服务数据是否正确插入到现有表中"
echo "5. 检查category字段是否设置为'护工服务'"
echo "6. 验证价格信息是否完整"
