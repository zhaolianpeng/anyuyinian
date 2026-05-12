#!/bin/bash

# 云函数部署脚本
# 使用方法: ./deploy_cloud_functions.sh

echo "☁️ 开始部署云开发微信支付云函数..."

# 检查是否在微信开发者工具中
if [ ! -d "cloudfunctions" ]; then
    echo "❌ 请在微信开发者工具中运行此脚本"
    exit 1
fi

# 云函数列表
FUNCTIONS=("payOrder" "payNotify" "queryOrder" "refundOrder")

echo "📦 准备部署以下云函数:"
for func in "${FUNCTIONS[@]}"; do
    echo "  - $func"
done

echo ""
echo "🚀 开始部署..."

# 部署每个云函数
for func in "${FUNCTIONS[@]}"; do
    echo ""
    echo "📤 部署 $func 云函数..."
    
    if [ -d "cloudfunctions/$func" ]; then
        echo "✅ $func 云函数目录存在"
        
        # 检查必要文件
        if [ -f "cloudfunctions/$func/index.js" ]; then
            echo "✅ $func/index.js 存在"
        else
            echo "❌ $func/index.js 不存在"
            continue
        fi
        
        if [ -f "cloudfunctions/$func/package.json" ]; then
            echo "✅ $func/package.json 存在"
        else
            echo "❌ $func/package.json 不存在"
            continue
        fi
        
        echo "📋 $func 云函数文件检查完成"
        echo "💡 请在微信开发者工具中右键 cloudfunctions/$func → 上传并部署"
        
    else
        echo "❌ $func 云函数目录不存在"
    fi
done

echo ""
echo "🎉 云函数部署准备完成！"
echo ""
echo "📝 下一步操作："
echo "1. 在微信开发者工具中右键每个云函数目录"
echo "2. 选择'上传并部署'"
echo "3. 等待部署完成"
echo "4. 测试云函数调用"
echo ""
echo "🔍 验证部署："
echo "1. 在云开发控制台查看云函数列表"
echo "2. 检查云函数日志"
echo "3. 测试云函数调用"
