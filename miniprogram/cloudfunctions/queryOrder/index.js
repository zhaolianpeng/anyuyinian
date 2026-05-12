// 云开发微信支付查询订单云函数
const cloud = require('wx-server-sdk')

// 初始化云开发
cloud.init({
  env: cloud.DYNAMIC_CURRENT_ENV
})

// 云开发支付
const pay = cloud.pay

/**
 * 查询订单
 * @param {Object} event 事件对象
 * @param {String} event.out_trade_no 商户订单号
 * @param {String} event.transaction_id 微信支付订单号（可选）
 */
exports.main = async (event, context) => {
  const { out_trade_no, transaction_id } = event
  
  try {
    console.log('查询订单:', { out_trade_no, transaction_id })
    
    // 调用云开发查询订单接口
    const result = await pay.orderQuery({
      // 商户订单号
      out_trade_no: out_trade_no,
      // 微信支付订单号（可选）
      transaction_id: transaction_id
    })
    
    console.log('查询订单成功:', result)
    
    return {
      success: true,
      data: result
    }
    
  } catch (error) {
    console.error('查询订单失败:', error)
    return {
      success: false,
      error: error.message || '查询订单失败'
    }
  }
}
