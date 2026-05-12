// 云开发微信支付申请退款云函数
const cloud = require('wx-server-sdk')

// 初始化云开发
cloud.init({
  env: cloud.DYNAMIC_CURRENT_ENV
})

// 云开发支付
const pay = cloud.pay

// 云开发数据库
const db = cloud.database()

/**
 * 申请退款
 * @param {Object} event 事件对象
 * @param {String} event.out_trade_no 商户订单号
 * @param {String} event.out_refund_no 商户退款单号
 * @param {Number} event.total_fee 订单总金额，单位为分
 * @param {Number} event.refund_fee 退款金额，单位为分
 * @param {String} event.refund_desc 退款原因（可选）
 */
exports.main = async (event, context) => {
  const { out_trade_no, out_refund_no, total_fee, refund_fee, refund_desc } = event
  
  try {
    console.log('申请退款:', { out_trade_no, out_refund_no, total_fee, refund_fee, refund_desc })
    
    // 调用云开发申请退款接口
    const result = await pay.refund({
      // 商户订单号
      out_trade_no: out_trade_no,
      // 商户退款单号
      out_refund_no: out_refund_no,
      // 订单总金额，单位为分
      total_fee: total_fee,
      // 退款金额，单位为分
      refund_fee: refund_fee,
      // 退款原因
      refund_desc: refund_desc || '用户申请退款'
    })
    
    console.log('申请退款成功:', result)
    
    // 更新订单状态为退款中
    await updateOrderStatus(out_trade_no, {
      status: 'refunding',
      refund_no: out_refund_no,
      refund_fee: refund_fee,
      refund_desc: refund_desc
    })
    
    return {
      success: true,
      data: result
    }
    
  } catch (error) {
    console.error('申请退款失败:', error)
    return {
      success: false,
      error: error.message || '申请退款失败'
    }
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
