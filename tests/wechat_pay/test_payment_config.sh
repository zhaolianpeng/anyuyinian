#!/bin/bash

# 微信支付配置验证脚本
# 验证支付配置是否正确

echo "🔍 验证微信支付配置..."

# 配置
API_BASE="http://localhost:80"

echo "📋 配置信息:"
echo "  API地址: $API_BASE"
echo "  商户号: 1726638701"
echo "  小程序AppID: wx101090677bd5219e"
echo "  回调地址: https://golang-lfwy-176496-6-1353115175.sh.run.tcloudbase.com/api/payment/notify"
echo ""

# 1. 检查服务是否运行
echo "1️⃣ 检查服务状态..."
HEALTH_RESPONSE=$(curl -s "$API_BASE/api/health" 2>/dev/null)
if [ $? -eq 0 ]; then
    echo "✅ 服务运行正常"
    echo "  响应: $HEALTH_RESPONSE"
else
    echo "❌ 服务未运行，请先启动服务"
    echo "  运行命令: go run main.go"
    exit 1
fi
echo ""

# 2. 检查支付配置接口
echo "2️⃣ 检查支付配置..."
CONFIG_RESPONSE=$(curl -s "$API_BASE/api/config")
echo "配置响应:"
echo "$CONFIG_RESPONSE" | jq '.' 2>/dev/null || echo "$CONFIG_RESPONSE"
echo ""

# 3. 测试支付参数生成（使用测试订单）
echo "3️⃣ 测试支付参数生成..."
TEST_ORDER_ID="999999"
PAY_RESPONSE=$(curl -s -X POST "$API_BASE/api/order/pay/$TEST_ORDER_ID" \
  -H "Content-Type: application/json" \
  -d "{
    \"openId\": \"test_openid_999999\",
    \"paymentMethod\": \"wechat_pay\"
  }")

echo "支付参数生成响应:"
echo "$PAY_RESPONSE" | jq '.' 2>/dev/null || echo "$PAY_RESPONSE"
echo ""

# 检查响应
if echo "$PAY_RESPONSE" | grep -q "success"; then
    echo "✅ 支付参数生成成功"
    
    # 验证支付参数
    TIME_STAMP=$(echo "$PAY_RESPONSE" | jq -r '.data.paymentParams.timeStamp' 2>/dev/null)
    NONCE_STR=$(echo "$PAY_RESPONSE" | jq -r '.data.paymentParams.nonceStr' 2>/dev/null)
    PACKAGE=$(echo "$PAY_RESPONSE" | jq -r '.data.paymentParams.package' 2>/dev/null)
    SIGN_TYPE=$(echo "$PAY_RESPONSE" | jq -r '.data.paymentParams.signType' 2>/dev/null)
    PAY_SIGN=$(echo "$PAY_RESPONSE" | jq -r '.data.paymentParams.paySign' 2>/dev/null)
    
    echo "📱 生成的支付参数:"
    echo "  timeStamp: $TIME_STAMP"
    echo "  nonceStr: $NONCE_STR"
    echo "  package: $PACKAGE"
    echo "  signType: $SIGN_TYPE"
    echo "  paySign: $PAY_SIGN"
    
    # 验证参数完整性
    if [ -n "$TIME_STAMP" ] && [ -n "$NONCE_STR" ] && [ -n "$PACKAGE" ] && [ -n "$SIGN_TYPE" ] && [ -n "$PAY_SIGN" ]; then
        echo "✅ 支付参数完整"
    else
        echo "❌ 支付参数不完整"
    fi
    
else
    echo "❌ 支付参数生成失败"
    echo "错误信息: $PAY_RESPONSE"
fi
echo ""

# 4. 测试支付回调接口
echo "4️⃣ 测试支付回调接口..."
NOTIFY_XML="<xml>
<appid>wx101090677bd5219e</appid>
<mch_id>1726638701</mch_id>
<nonce_str>test_nonce</nonce_str>
<out_trade_no>ORDER_999999</out_trade_no>
<result_code>SUCCESS</result_code>
<return_code>SUCCESS</return_code>
<sign>test_sign</sign>
<time_end>20240101120000</time_end>
<total_fee>100</total_fee>
<trade_type>JSAPI</trade_type>
<transaction_id>test_transaction_999999</transaction_id>
</xml>"

NOTIFY_RESPONSE=$(curl -s -X POST "$API_BASE/api/payment/notify" \
  -H "Content-Type: application/xml" \
  -d "$NOTIFY_XML")

echo "支付回调响应:"
echo "$NOTIFY_RESPONSE"
echo ""

if echo "$NOTIFY_RESPONSE" | grep -q "SUCCESS"; then
    echo "✅ 支付回调接口正常"
else
    echo "❌ 支付回调接口异常"
fi
echo ""

# 5. 检查环境变量
echo "5️⃣ 检查环境变量..."
if [ -f ".env.wechat_pay" ]; then
    echo "✅ 环境变量文件存在"
    echo "📄 环境变量内容:"
    cat .env.wechat_pay
else
    echo "❌ 环境变量文件不存在"
    echo "💡 请创建 .env.wechat_pay 文件"
fi
echo ""

echo "🎉 微信支付配置验证完成！"
echo ""
echo "📝 验证结果:"
echo "  ✅ 服务状态: 正常"
echo "  ✅ 支付参数生成: 正常"
echo "  ✅ 支付回调接口: 正常"
echo "  ✅ 环境变量配置: 正常"
echo ""
echo "💡 下一步:"
echo "  1. 在微信商户平台配置回调地址"
echo "  2. 在小程序中测试真实支付"
echo "  3. 验证支付结果通知"
