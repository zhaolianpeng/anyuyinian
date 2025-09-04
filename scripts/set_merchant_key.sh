#!/bin/bash

# 设置微信支付商户密钥脚本
# 商户号: 1726638701

echo "=== 微信支付商户密钥设置 ==="
echo "商户号: 1726638701"
echo ""

# 检查是否已设置环境变量
if [ -n "$WECHAT_PAY_MCH_KEY" ]; then
    echo "当前商户密钥: ${WECHAT_PAY_MCH_KEY:0:8}****"
    echo ""
    read -p "是否要重新设置商户密钥? (y/n): " confirm
    if [ "$confirm" != "y" ]; then
        echo "保持当前设置"
        exit 0
    fi
fi

echo "请按照以下步骤获取商户密钥:"
echo ""
echo "1. 访问微信商户平台: https://pay.weixin.qq.com/"
echo "2. 登录并选择商户号: 1726638701"
echo "3. 进入 账户中心 -> API安全 -> API密钥"
echo "4. 点击 '设置API密钥' 或 '修改API密钥'"
echo "5. 设置32位字符串作为商户密钥"
echo "6. 复制生成的密钥"
echo ""

# 提示用户输入商户密钥
echo "请输入您的32位商户密钥:"
read -s MERCHANT_KEY

# 验证密钥长度
if [ ${#MERCHANT_KEY} -ne 32 ]; then
    echo "❌ 错误: 商户密钥必须是32位字符串，当前为${#MERCHANT_KEY}位"
    exit 1
fi

# 设置环境变量
export WECHAT_PAY_MCH_KEY="$MERCHANT_KEY"
echo "✅ 商户密钥设置成功"

# 更新环境变量文件
if [ -f ".env.wechat_pay" ]; then
    # 更新现有文件
    sed -i.bak "s/WECHAT_PAY_MCH_KEY=.*/WECHAT_PAY_MCH_KEY=$MERCHANT_KEY/" .env.wechat_pay
    echo "✅ 已更新 .env.wechat_pay 文件"
else
    # 创建新文件
    cat > .env.wechat_pay << EOF
# 微信支付配置
WECHAT_PAY_APP_ID=wx101090677bd5219e
WECHAT_PAY_MCH_ID=1726638701
WECHAT_PAY_MCH_KEY=$MERCHANT_KEY
WECHAT_PAY_NOTIFY_URL=https://your-domain.com/api/payment/notify
WECHAT_PAY_ENVIRONMENT=production
EOF
    echo "✅ 已创建 .env.wechat_pay 文件"
fi

echo ""
echo "=== 配置完成 ==="
echo "商户号: 1726638701"
echo "商户密钥: ${MERCHANT_KEY:0:8}****"
echo ""

# 提供加载环境变量的命令
echo "请运行以下命令加载环境变量:"
echo "source .env.wechat_pay"
echo ""
echo "或者重启应用以加载新的环境变量"
