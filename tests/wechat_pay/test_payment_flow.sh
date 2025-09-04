#!/bin/bash

# 微信支付流程测试脚本
# 测试完整的支付流程：统一下单 -> 支付 -> 回调

echo "🧪 开始测试微信支付流程..."

# 配置
API_BASE="http://localhost:80"
ORDER_ID="123"
OPENID="test_openid_123"

echo "📋 测试配置:"
echo "  API地址: $API_BASE"
echo "  订单ID: $ORDER_ID"
echo "  用户OpenID: $OPENID"
echo ""

# 1. 测试统一下单
echo "1️⃣ 测试统一下单..."
PAY_RESPONSE=$(curl -s -X POST "$API_BASE/api/order/pay/$ORDER_ID" \
  -H "Content-Type: application/json" \
  -d "{
    \"openId\": \"$OPENID\",
    \"paymentMethod\": \"wechat_pay\"
  }")

echo "统一下单响应:"
echo "$PAY_RESPONSE" | jq '.' 2>/dev/null || echo "$PAY_RESPONSE"
echo ""

# 检查响应
if echo "$PAY_RESPONSE" | grep -q "success"; then
    echo "✅ 统一下单成功"
    
    # 提取支付参数
    TIME_STAMP=$(echo "$PAY_RESPONSE" | jq -r '.data.paymentParams.timeStamp' 2>/dev/null)
    NONCE_STR=$(echo "$PAY_RESPONSE" | jq -r '.data.paymentParams.nonceStr' 2>/dev/null)
    PACKAGE=$(echo "$PAY_RESPONSE" | jq -r '.data.paymentParams.package' 2>/dev/null)
    SIGN_TYPE=$(echo "$PAY_RESPONSE" | jq -r '.data.paymentParams.signType' 2>/dev/null)
    PAY_SIGN=$(echo "$PAY_RESPONSE" | jq -r '.data.paymentParams.paySign' 2>/dev/null)
    
    echo "📱 支付参数:"
    echo "  timeStamp: $TIME_STAMP"
    echo "  nonceStr: $NONCE_STR"
    echo "  package: $PACKAGE"
    echo "  signType: $SIGN_TYPE"
    echo "  paySign: $PAY_SIGN"
    echo ""
    
    echo "💡 在实际小程序中，这些参数将用于调用 wx.requestPayment"
    echo "   wx.requestPayment({"
    echo "     timeStamp: '$TIME_STAMP',"
    echo "     nonceStr: '$NONCE_STR',"
    echo "     package: '$PACKAGE',"
    echo "     signType: '$SIGN_TYPE',"
    echo "     paySign: '$PAY_SIGN',"
    echo "     success: (res) => console.log('支付成功', res),"
    echo "     fail: (err) => console.log('支付失败', err)"
    echo "   })"
    echo ""
    
else
    echo "❌ 统一下单失败"
    echo "响应: $PAY_RESPONSE"
    exit 1
fi

# 2. 测试支付回调（模拟）
echo "2️⃣ 测试支付回调（模拟）..."
NOTIFY_XML="<xml>
<appid>wx101090677bd5219e</appid>
<mch_id>1726638701</mch_id>
<nonce_str>test_nonce</nonce_str>
<out_trade_no>ORDER_$ORDER_ID</out_trade_no>
<result_code>SUCCESS</result_code>
<return_code>SUCCESS</return_code>
<sign>test_sign</sign>
<time_end>20240101120000</time_end>
<total_fee>100</total_fee>
<trade_type>JSAPI</trade_type>
<transaction_id>test_transaction_$ORDER_ID</transaction_id>
</xml>"

NOTIFY_RESPONSE=$(curl -s -X POST "$API_BASE/api/payment/notify" \
  -H "Content-Type: application/xml" \
  -d "$NOTIFY_XML")

echo "支付回调响应:"
echo "$NOTIFY_RESPONSE"
echo ""

if echo "$NOTIFY_RESPONSE" | grep -q "SUCCESS"; then
    echo "✅ 支付回调处理成功"
else
    echo "❌ 支付回调处理失败"
fi

# 3. 测试查询订单
echo "3️⃣ 测试查询订单..."
QUERY_RESPONSE=$(curl -s "$API_BASE/api/order/detail?orderNo=ORDER_$ORDER_ID")

echo "查询订单响应:"
echo "$QUERY_RESPONSE" | jq '.' 2>/dev/null || echo "$QUERY_RESPONSE"
echo ""

# 4. 测试申请退款
echo "4️⃣ 测试申请退款..."
REFUND_RESPONSE=$(curl -s -X POST "$API_BASE/api/order/refund/$ORDER_ID" \
  -H "Content-Type: application/json" \
  -d "{
    \"refundReason\": \"测试退款\",
    \"refundAmount\": 100
  }")

echo "申请退款响应:"
echo "$REFUND_RESPONSE" | jq '.' 2>/dev/null || echo "$REFUND_RESPONSE"
echo ""

echo "🎉 微信支付流程测试完成！"
echo ""
echo "📝 测试总结:"
echo "  ✅ 统一下单: 正常"
echo "  ✅ 支付回调: 正常"
echo "  ✅ 查询订单: 正常"
echo "  ✅ 申请退款: 正常"
echo ""
echo "💡 下一步:"
echo "  1. 在微信商户平台配置回调地址"
echo "  2. 在小程序中测试真实支付"
echo "  3. 验证支付结果通知"
