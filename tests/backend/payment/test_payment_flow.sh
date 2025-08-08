#!/bin/bash

# 测试支付流程

echo "=== 测试支付流程 ==="

# 设置测试环境
BASE_URL="http://localhost:80"
TEST_USER_ID="507f1f77bcf86cd799439011"
TEST_ORDER_ID="1"

echo "1. 测试订单支付接口"
echo "请求URL: ${BASE_URL}/api/order/pay/${TEST_ORDER_ID}"

# 构建支付请求数据
payment_data=$(cat <<EOF
{
  "orderId": ${TEST_ORDER_ID},
  "payMethod": "wechat_pay",
  "openId": "test_openid_${TEST_USER_ID}"
}
EOF
)

echo "请求数据:"
echo "$payment_data" | jq '.'

# 发送支付请求
response=$(curl -s -X POST "${BASE_URL}/api/order/pay/${TEST_ORDER_ID}" \
  -H "Content-Type: application/json" \
  -d "$payment_data")

echo ""
echo "响应内容:"
echo "$response" | jq '.'

# 检查支付响应
echo ""
echo "2. 检查支付响应"
code=$(echo "$response" | jq -r '.code // -1')

if [ "$code" = "0" ]; then
    echo "✅ 支付请求成功"
    
    # 检查支付参数
    payment_params=$(echo "$response" | jq -r '.data.paymentParams // {}')
    
    if [ "$payment_params" != "{}" ] && [ "$payment_params" != "null" ]; then
        echo "✅ 支付参数生成成功"
        echo "支付参数详情:"
        echo "$payment_params" | jq '.'
        
        # 检查必要的支付参数
        timeStamp=$(echo "$payment_params" | jq -r '.timeStamp // ""')
        nonceStr=$(echo "$payment_params" | jq -r '.nonceStr // ""')
        package=$(echo "$payment_params" | jq -r '.package // ""')
        signType=$(echo "$payment_params" | jq -r '.signType // ""')
        paySign=$(echo "$payment_params" | jq -r '.paySign // ""')
        
        echo ""
        echo "3. 验证支付参数完整性:"
        
        if [ "$timeStamp" != "" ]; then
            echo "✅ timeStamp: $timeStamp"
        else
            echo "❌ timeStamp 缺失"
        fi
        
        if [ "$nonceStr" != "" ]; then
            echo "✅ nonceStr: $nonceStr"
        else
            echo "❌ nonceStr 缺失"
        fi
        
        if [ "$package" != "" ]; then
            echo "✅ package: $package"
        else
            echo "❌ package 缺失"
        fi
        
        if [ "$signType" != "" ]; then
            echo "✅ signType: $signType"
        else
            echo "❌ signType 缺失"
        fi
        
        if [ "$paySign" != "" ]; then
            echo "✅ paySign: $paySign"
        else
            echo "❌ paySign 缺失"
        fi
        
        echo ""
        echo "4. 支付参数验证完成"
        echo "这些参数可以直接用于微信小程序支付"
        
    else
        echo "❌ 支付参数生成失败"
        echo "可能的原因:"
        echo "   - 微信支付配置不完整"
        echo "   - 商户号或密钥未配置"
        echo "   - 网络连接问题"
    fi
    
else
    echo "❌ 支付请求失败"
    error_msg=$(echo "$response" | jq -r '.errorMsg // "未知错误"')
    echo "错误信息: $error_msg"
    
    echo ""
    echo "可能的原因:"
    echo "   - 订单不存在或状态不正确"
    echo "   - 微信支付配置问题"
    echo "   - 网络连接问题"
fi

echo ""
echo "=== 测试完成 ==="
echo ""
echo "如果支付失败，请检查:"
echo "1. 订单是否存在且状态为待支付"
echo "2. 微信支付配置是否完整"
echo "3. 网络连接是否正常"
echo "4. 后端服务是否正常运行" 