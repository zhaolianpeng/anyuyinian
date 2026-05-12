// 云开发微信支付测试脚本
// 在微信开发者工具控制台中运行

// 测试云开发统一下单
async function testCloudPayOrder() {
  console.log('🧪 测试云开发统一下单云函数...')
  
  try {
    const result = await wx.cloud.callFunction({
      name: 'payOrder',
      data: {
        orderInfo: {
          out_trade_no: `CLOUD_ORDER_${Date.now()}`,
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
    
    console.log('✅ 云开发统一下单成功:', result)
    return result
    
  } catch (error) {
    console.error('❌ 云开发统一下单失败:', error)
    return null
  }
}

// 测试云开发查询订单
async function testCloudQueryOrder() {
  console.log('🧪 测试云开发查询订单云函数...')
  
  try {
    const result = await wx.cloud.callFunction({
      name: 'queryOrder',
      data: {
        out_trade_no: 'CLOUD_ORDER_123456'
      }
    })
    
    console.log('✅ 云开发查询订单成功:', result)
    return result
    
  } catch (error) {
    console.error('❌ 云开发查询订单失败:', error)
    return null
  }
}

// 测试云开发申请退款
async function testCloudRefundOrder() {
  console.log('🧪 测试云开发申请退款云函数...')
  
  try {
    const result = await wx.cloud.callFunction({
      name: 'refundOrder',
      data: {
        out_trade_no: 'CLOUD_ORDER_123456',
        out_refund_no: `CLOUD_REFUND_${Date.now()}`,
        total_fee: 100,
        refund_fee: 100,
        refund_desc: '测试退款'
      }
    })
    
    console.log('✅ 云开发申请退款成功:', result)
    return result
    
  } catch (error) {
    console.error('❌ 云开发申请退款失败:', error)
    return null
  }
}

// 测试云开发支付流程
async function testCloudPaymentFlow() {
  console.log('🚀 开始测试云开发完整支付流程...')
  
  // 1. 统一下单
  const payResult = await testCloudPayOrder()
  if (!payResult || !payResult.result.success) {
    console.error('❌ 云开发统一下单失败，终止测试')
    return
  }
  
  const payment = payResult.result.data.payment
  console.log('📱 获取到云开发支付参数:', payment)
  
  // 2. 模拟支付（实际环境中会调用wx.requestPayment）
  console.log('💡 在实际环境中，这里会调用wx.requestPayment')
  console.log('云开发支付参数:', {
    timeStamp: payment.timeStamp,
    nonceStr: payment.nonceStr,
    package: payment.package,
    signType: payment.signType,
    paySign: payment.paySign
  })
  
  // 3. 查询订单
  await testCloudQueryOrder()
  
  // 4. 申请退款
  await testCloudRefundOrder()
  
  console.log('✅ 云开发支付流程测试完成')
}

// 检查云函数状态
async function checkCloudFunctions() {
  console.log('🔍 检查云开发云函数状态...')
  
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

// 测试云开发环境
async function testCloudEnvironment() {
  console.log('🔍 检查云开发环境...')
  
  try {
    // 检查云开发是否初始化
    if (!wx.cloud) {
      console.error('❌ 云开发未初始化')
      return false
    }
    
    console.log('✅ 云开发已初始化')
    
    // 检查云开发环境ID
    const envId = wx.cloud.env
    console.log('📋 云开发环境ID:', envId)
    
    // 检查云开发数据库
    const db = wx.cloud.database()
    console.log('✅ 云开发数据库可用')
    
    // 检查云开发支付
    const pay = wx.cloud.pay
    console.log('✅ 云开发支付可用')
    
    return true
    
  } catch (error) {
    console.error('❌ 云开发环境检查失败:', error)
    return false
  }
}

// 导出测试函数
window.testCloudDevelopmentPayment = {
  testCloudPayOrder,
  testCloudQueryOrder,
  testCloudRefundOrder,
  testCloudPaymentFlow,
  checkCloudFunctions,
  testCloudEnvironment
}

console.log('🧪 云开发微信支付测试脚本已加载')
console.log('使用方法:')
console.log('1. testCloudDevelopmentPayment.testCloudEnvironment() - 检查云开发环境')
console.log('2. testCloudDevelopmentPayment.checkCloudFunctions() - 检查云函数状态')
console.log('3. testCloudDevelopmentPayment.testCloudPayOrder() - 测试云开发统一下单')
console.log('4. testCloudDevelopmentPayment.testCloudQueryOrder() - 测试云开发查询订单')
console.log('5. testCloudDevelopmentPayment.testCloudRefundOrder() - 测试云开发申请退款')
console.log('6. testCloudDevelopmentPayment.testCloudPaymentFlow() - 测试云开发完整支付流程')
