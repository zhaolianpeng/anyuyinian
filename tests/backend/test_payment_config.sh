#!/bin/bash

# 测试支付配置

echo "=== 测试支付配置 ==="

# 设置测试环境
BASE_URL="http://localhost:80"

echo "1. 测试支付配置接口"
echo "请求URL: ${BASE_URL}/api/config"

# 发送请求
response=$(curl -s -X GET "${BASE_URL}/api/config")

echo "响应内容:"
echo "$response" | jq '.'

# 检查支付配置
echo ""
echo "2. 检查支付配置"
payment_config=$(echo "$response" | jq -r '.data.payment_config // {}')

if [ "$payment_config" = "{}" ] || [ "$payment_config" = "null" ]; then
    echo "❌ 支付配置为空"
    echo ""
    echo "3. 需要配置支付参数:"
    echo "   - 设置环境变量 WECHAT_PAY_MCH_ID"
    echo "   - 设置环境变量 WECHAT_PAY_MCH_KEY"
    echo "   - 设置环境变量 WECHAT_PAY_NOTIFY_URL"
    echo ""
    echo "4. 配置示例:"
    echo "   export WECHAT_PAY_MCH_ID='你的商户号'"
    echo "   export WECHAT_PAY_MCH_KEY='你的商户密钥'"
    echo "   export WECHAT_PAY_NOTIFY_URL='https://your-domain.com/api/payment/notify'"
else
    echo "✅ 支付配置存在"
    echo "支付配置详情:"
    echo "$payment_config" | jq '.'
fi

# 检查微信配置
echo ""
echo "5. 检查微信配置"
wechat_appid=$(echo "$response" | jq -r '.data.wechat_appid // ""')
wechat_secret=$(echo "$response" | jq -r '.data.wechat_secret // ""')

if [ "$wechat_appid" = "" ] || [ "$wechat_appid" = "null" ]; then
    echo "❌ 微信AppID未配置"
else
    echo "✅ 微信AppID已配置: $wechat_appid"
fi

if [ "$wechat_secret" = "" ] || [ "$wechat_secret" = "null" ]; then
    echo "❌ 微信AppSecret未配置"
else
    echo "✅ 微信AppSecret已配置: $wechat_secret"
fi

echo ""
echo "=== 测试完成 ==="
echo ""
echo "如果支付配置不完整，请:"
echo "1. 设置必要的环境变量"
echo "2. 重启后端服务"
echo "3. 重新运行此测试脚本" 