// 云开发微信支付测试脚本
// 在微信开发者工具控制台中运行

// 测试统一下单
async function testPayOrder() {
  console.log('🧪 测试统一下单云函数...')
  
  try {
    const result = await wx.cloud.callFunction({
      name: 'payOrder',
      data: {
        orderInfo: {
          out_trade_no: `TEST_ORDER_${Date.now()}`,
          body: '安语一年服务测试',
          total_fee: 1, // 1分钱测试
          spbill_create_ip: '127.0.0.1',
          openid: 'test_openid',
          attach: 'test_order',
          detail: '安语一年服务测试',
          goods_tag: '测试',
          time_start: '20240101120000',
          time_expire: '20240101123000'
        }
      }
    })
    
    console.log('✅ 统一下单成功:', result)
    return result
    
  } catch (error) {
    console.error('❌ 统一下单失败:', error)
    return null
  }
}

// 测试查询订单
async function testQueryOrder() {
  console.log('🧪 测试查询订单云函数...')
  
  try {
    const result = await wx.cloud.callFunction({
      name: 'queryOrder',
      data: {
        out_trade_no: 'TEST_ORDER_123456'
      }
    })
    
    console.log('✅ 查询订单成功:', result)
    return result
    
  } catch (error) {
    console.error('❌ 查询订单失败:', error)
    return null
  }
}

// 测试申请退款
async function testRefundOrder() {
  console.log('🧪 测试申请退款云函数...')
  
  try {
    const result = await wx.cloud.callFunction({
      name: 'refundOrder',
      data: {
        out_trade_no: 'TEST_ORDER_123456',
        out_refund_no: `TEST_REFUND_${Date.now()}`,
        total_fee: 100,
        refund_fee: 100,
        refund_desc: '测试退款'
      }
    })
    
    console.log('✅ 申请退款成功:', result)
    return result
    
  } catch (error) {
    console.error('❌ 申请退款失败:', error)
    return null
  }
}

// 测试支付流程
async function testPaymentFlow() {
  console.log('🚀 开始测试完整支付流程...')
  
  // 1. 统一下单
  const payResult = await testPayOrder()
  if (!payResult || !payResult.result.success) {
    console.error('❌ 统一下单失败，终止测试')
    return
  }
  
  const payment = payResult.result.data.payment
  console.log('📱 获取到支付参数:', payment)
  
  // 2. 模拟支付（实际环境中会调用wx.requestPayment）
  console.log('💡 在实际环境中，这里会调用wx.requestPayment')
  console.log('支付参数:', {
    timeStamp: payment.timeStamp,
    nonceStr: payment.nonceStr,
    package: payment.package,
    signType: payment.signType,
    paySign: payment.paySign
  })
  
  // 3. 查询订单
  await testQueryOrder()
  
  // 4. 申请退款
  await testRefundOrder()
  
  console.log('✅ 支付流程测试完成')
}

// 检查云函数状态
async function checkCloudFunctions() {
  console.log('🔍 检查云函数状态...')
  
  const functions = ['payOrder', 'payNotify', 'queryOrder', 'refundOrder']
  
  for (const funcName of functions) {
    try {
      // 尝试调用云函数（使用空数据）
      await wx.cloud.callFunction({
        name: funcName,
        data: {}
      })
      console.log(`✅ ${funcName} 云函数可访问`)
    } catch (error) {
      console.log(`❌ ${funcName} 云函数不可访问:`, error.message)
    }
  }
}

// 导出测试函数
window.testCloudPayment = {
  testPayOrder,
  testQueryOrder,
  testRefundOrder,
  testPaymentFlow,
  checkCloudFunctions
}

console.log('🧪 云开发微信支付测试脚本已加载')
console.log('使用方法:')
console.log('1. testCloudPayment.checkCloudFunctions() - 检查云函数状态')
console.log('2. testCloudPayment.testPayOrder() - 测试统一下单')
console.log('3. testCloudPayment.testQueryOrder() - 测试查询订单')
console.log('4. testCloudPayment.testRefundOrder() - 测试申请退款')
console.log('5. testCloudPayment.testPaymentFlow() - 测试完整支付流程')
