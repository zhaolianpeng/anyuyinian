// 云开发微信支付回调云函数
const cloud = require('wx-server-sdk')

// 初始化云开发
cloud.init({
  env: cloud.DYNAMIC_CURRENT_ENV
})

// 云开发数据库
const db = cloud.database()

/**
 * 支付结果通知
 * @param {Object} event 事件对象
 * @param {Object} event.payResult 支付结果
 */
exports.main = async (event, context) => {
  const { payResult } = event
  
  try {
    console.log('收到云开发支付回调:', payResult)
    
    // 验证支付结果
    if (!payResult || !payResult.out_trade_no) {
      console.error('支付回调数据无效')
      return { errcode: -1, errmsg: '支付回调数据无效' }
    }
    
    const { out_trade_no, transaction_id, total_fee, result_code, return_code } = payResult
    
    // 检查支付状态
    if (return_code === 'SUCCESS' && result_code === 'SUCCESS') {
      console.log('支付成功:', {
        out_trade_no,
        transaction_id,
        total_fee
      })
      
      // 更新订单状态为已支付
      await updateOrderStatus(out_trade_no, {
        status: 'paid',
        transaction_id,
        paid_at: new Date(),
        total_fee
      })
      
      // 这里可以添加其他业务逻辑
      // 比如：发送支付成功通知、更新库存等
      
    } else {
      console.log('支付失败:', payResult)
      
      // 更新订单状态为支付失败
      await updateOrderStatus(out_trade_no, {
        status: 'pay_failed',
        fail_reason: payResult.err_code_des || '支付失败'
      })
    }
    
    // 返回成功响应（必须返回errcode: 0）
    return { errcode: 0, errmsg: 'success' }
    
  } catch (error) {
    console.error('处理云开发支付回调失败:', error)
    return { errcode: -1, errmsg: '处理支付回调失败' }
  }
}

/**
 * 更新订单状态
 * @param {String} outTradeNo 商户订单号
 * @param {Object} updateData 更新数据
 */
async function updateOrderStatus(outTradeNo, updateData) {
  try {
    const result = await db.collection('orders').where({
      out_trade_no: outTradeNo
    }).update({
      data: {
        ...updateData,
        updated_at: new Date()
      }
    })
    
    console.log('订单状态更新成功:', result)
    return result
    
  } catch (error) {
    console.error('更新订单状态失败:', error)
    throw error
  }
}
