#!/bin/bash

# 护工服务平台重新组织后的服务分类测试脚本

# 设置测试服务器地址
SERVER_URL="http://localhost:80"

echo "=== 护工服务平台重新组织后的服务分类测试 ==="
echo "服务器地址: $SERVER_URL"
echo ""

# 测试1: 获取首页数据
echo "=== 测试1: 获取首页数据 ==="
response=$(curl -s -X GET "$SERVER_URL/api/home/init")
echo "$response" | jq '.'
echo ""

# 测试2: 验证护工服务数据结构
echo "=== 测试2: 验证护工服务数据结构 ==="
if echo "$response" | jq -e '.data.caregiverServices' > /dev/null; then
    echo "✅ 护工服务数据存在"
    caregiver_count=$(echo "$response" | jq '.data.caregiverServices | length')
    echo "护工服务总数: $caregiver_count"
else
    echo "❌ 护工服务数据缺失"
    exit 1
fi
echo ""

# 测试3: 验证服务分类
echo "=== 测试3: 验证服务分类 ==="
categories=$(echo "$response" | jq -r '.data.caregiverServices[].category' | sort | uniq)
echo "服务分类:"
echo "$categories"
echo ""

# 测试4: 验证居家照护类服务
echo "=== 测试4: 验证居家照护类服务 ==="
home_care_services=$(echo "$response" | jq '.data.caregiverServices[] | select(.category == "居家照护")')
if [ ! -z "$home_care_services" ]; then
    echo "✅ 居家照护类服务存在"
    home_care_count=$(echo "$home_care_services" | jq -s 'length')
    echo "居家照护服务数量: $home_care_count"
    echo "服务列表:"
    echo "$home_care_services" | jq '{name: .name, price: .price, description: .description}'
else
    echo "❌ 居家照护类服务缺失"
fi
echo ""

# 测试5: 验证医院陪诊类服务
echo "=== 测试5: 验证医院陪诊类服务 ==="
hospital_escort_services=$(echo "$response" | jq '.data.caregiverServices[] | select(.category == "医院陪诊")')
if [ ! -z "$hospital_escort_services" ]; then
    echo "✅ 医院陪诊类服务存在"
    hospital_escort_count=$(echo "$hospital_escort_services" | jq -s 'length')
    echo "医院陪诊服务数量: $hospital_escort_count"
    echo "服务列表:"
    echo "$hospital_escort_services" | jq '{name: .name, price: .price, description: .description}'
else
    echo "❌ 医院陪诊类服务缺失"
fi
echo ""

# 测试6: 验证周期护理类服务
echo "=== 测试6: 验证周期护理类服务 ==="
periodic_care_services=$(echo "$response" | jq '.data.caregiverServices[] | select(.category == "周期护理")')
if [ ! -z "$periodic_care_services" ]; then
    echo "✅ 周期护理类服务存在"
    periodic_care_count=$(echo "$periodic_care_services" | jq -s 'length')
    echo "周期护理服务数量: $periodic_care_count"
    echo "服务列表:"
    echo "$periodic_care_services" | jq '{name: .name, price: .price, description: .description}'
else
    echo "❌ 周期护理类服务缺失"
fi
echo ""

# 测试7: 验证家政服务类服务
echo "=== 测试7: 验证家政服务类服务 ==="
housekeeping_services=$(echo "$response" | jq '.data.caregiverServices[] | select(.category == "家政服务")')
if [ ! -z "$housekeeping_services" ]; then
    echo "✅ 家政服务类服务存在"
    housekeeping_count=$(echo "$housekeeping_services" | jq -s 'length')
    echo "家政服务数量: $housekeeping_count"
    echo "服务列表:"
    echo "$housekeeping_services" | jq '{name: .name, price: .price, description: .description}'
else
    echo "❌ 家政服务类服务缺失"
fi
echo ""

# 测试8: 验证服务价格信息
echo "=== 测试8: 验证服务价格信息 ==="
services_with_price=$(echo "$response" | jq '.data.caregiverServices[] | select(.price != null)')
if [ ! -z "$services_with_price" ]; then
    echo "✅ 服务价格信息完整"
    echo "价格统计:"
    echo "$services_with_price" | jq '{name: .name, price: .price, category: .category}' | head -10
else
    echo "❌ 服务价格信息缺失"
fi
echo ""

# 测试9: 验证服务描述信息
echo "=== 测试9: 验证服务描述信息 ==="
services_with_desc=$(echo "$response" | jq '.data.caregiverServices[] | select(.description != null and .description != "")')
if [ ! -z "$services_with_desc" ]; then
    echo "✅ 服务描述信息完整"
    echo "描述示例:"
    echo "$services_with_desc" | jq '{name: .name, description: .description}' | head -5
else
    echo "❌ 服务描述信息缺失"
fi
echo ""

# 测试10: 验证数据库中的分类数据
echo "=== 测试10: 验证数据库中的分类数据 ==="
echo "请手动执行以下SQL查询来验证数据库中的分类数据："
echo ""
echo "1. 查询居家照护类服务："
echo "SELECT id, name, price, category FROM Services WHERE category = '居家照护' ORDER BY sort;"
echo ""
echo "2. 查询医院陪诊类服务："
echo "SELECT id, name, price, category FROM Services WHERE category = '医院陪诊' ORDER BY sort;"
echo ""
echo "3. 查询周期护理类服务："
echo "SELECT id, name, price, category FROM Services WHERE category = '周期护理' ORDER BY sort;"
echo ""
echo "4. 查询家政服务类服务："
echo "SELECT id, name, price, category FROM Services WHERE category = '家政服务' ORDER BY sort;"
echo ""
echo "5. 统计各分类服务数量："
echo "SELECT category, COUNT(*) as count FROM Services WHERE category IN ('居家照护', '医院陪诊', '周期护理', '家政服务') GROUP BY category;"
echo ""

echo "=== 测试完成 ==="
echo ""
echo "📊 测试总结："
echo "- 护工服务总数: $caregiver_count"
echo "- 服务分类: $(echo "$categories" | tr '\n' ', ')"
echo ""
echo "✅ 如果所有测试都通过，说明护工服务数据重新组织成功！"
echo ""
echo "⚠️  注意事项："
echo "1. 确保已执行新的迁移文件: 20241220_reorganize_caregiver_services.sql"
echo "2. 验证Services表中的分类字段是否正确设置"
echo "3. 检查各分类下的服务数量是否符合预期"
echo "4. 确认价格和描述信息是否完整"
