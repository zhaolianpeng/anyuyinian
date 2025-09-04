#!/bin/bash

# 微信支付配置设置脚本
# 商户号: 1726638701

echo "=== 微信支付配置设置 ==="
echo "商户号: 1726638701"
echo ""

# 检查是否已设置环境变量
if [ -z "$WECHAT_PAY_MCH_KEY" ]; then
    echo "⚠️  请先设置微信支付商户密钥环境变量:"
    echo "export WECHAT_PAY_MCH_KEY=\"您的商户密钥\""
    echo ""
    echo "获取商户密钥的步骤:"
    echo "1. 登录微信商户平台: https://pay.weixin.qq.com/"
    echo "2. 进入 账户中心 -> API安全 -> API密钥"
    echo "3. 设置API密钥（32位字符串）"
    echo "4. 复制密钥并设置环境变量"
    echo ""
fi

# 设置其他必要的环境变量
echo "设置微信支付环境变量..."

# 微信小程序AppID
export WECHAT_PAY_APP_ID="wx101090677bd5219e"
echo "✅ AppID: $WECHAT_PAY_APP_ID"

# 商户号
export WECHAT_PAY_MCH_ID="1726638701"
echo "✅ 商户号: $WECHAT_PAY_MCH_ID"

# 支付通知URL（需要替换为您的实际域名）
export WECHAT_PAY_NOTIFY_URL="https://your-domain.com/api/payment/notify"
echo "⚠️  通知URL: $WECHAT_PAY_NOTIFY_URL (请替换为您的实际域名)"

# 环境设置
export WECHAT_PAY_ENVIRONMENT="production"
echo "✅ 环境: $WECHAT_PAY_ENVIRONMENT"

echo ""
echo "=== 配置完成 ==="
echo "请确保以下配置正确:"
echo "1. 商户密钥已设置: WECHAT_PAY_MCH_KEY"
echo "2. 通知URL已更新为您的实际域名"
echo "3. 微信小程序已关联商户号1726638701"
echo "4. 商户平台已配置支付回调地址"
echo ""

# 生成环境变量配置文件
cat > .env.wechat_pay << EOF
# 微信支付配置
WECHAT_PAY_APP_ID=wx101090677bd5219e
WECHAT_PAY_MCH_ID=1726638701
WECHAT_PAY_MCH_KEY=请设置您的商户密钥
WECHAT_PAY_NOTIFY_URL=https://your-domain.com/api/payment/notify
WECHAT_PAY_ENVIRONMENT=production
EOF

echo "已生成环境变量配置文件: .env.wechat_pay"
echo "请编辑此文件并设置正确的商户密钥和通知URL"
