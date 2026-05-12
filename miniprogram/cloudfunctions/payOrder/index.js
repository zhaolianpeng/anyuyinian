// 云开发微信支付云函数
const cloud = require('wx-server-sdk')

// 初始化云开发
cloud.init({
  env: cloud.DYNAMIC_CURRENT_ENV
})

// 云开发数据库
const db = cloud.database()

// 云开发支付
const pay = cloud.pay

/**
 * 统一下单
 * @param {Object} event 事件对象
 * @param {Object} event.orderInfo 订单信息
 * @param {String} event.orderInfo.out_trade_no 商户订单号
 * @param {String} event.orderInfo.body 商品描述
 * @param {Number} event.orderInfo.total_fee 订单总金额，单位为分
 * @param {String} event.orderInfo.spbill_create_ip 用户端IP
 * @param {String} event.orderInfo.openid 用户openid
 */
exports.main = async (event, context) => {
  const { orderInfo } = event
  
  try {
    console.log('开始云开发统一下单:', orderInfo)
    
    // 调用云开发统一下单接口
    const result = await pay.unifiedOrder({
      // 商户订单号
      out_trade_no: orderInfo.out_trade_no,
      // 商品描述
      body: orderInfo.body,
      // 订单总金额，单位为分
      total_fee: orderInfo.total_fee,
      // 用户端IP
      spbill_create_ip: orderInfo.spbill_create_ip,
      // 支付结果通知地址（云函数名）
      notify_url: 'payNotify',
      // 交易类型
      trade_type: 'JSAPI',
      // 用户openid
      openid: orderInfo.openid,
      // 附加数据
      attach: orderInfo.attach || '',
      // 商品详情
      detail: orderInfo.detail || '',
      // 订单优惠标记
      goods_tag: orderInfo.goods_tag || '',
      // 订单生成时间
      time_start: orderInfo.time_start || '',
      // 订单失效时间
      time_expire: orderInfo.time_expire || '',
      // 商品标记
      product_id: orderInfo.product_id || ''
    })
    
    console.log('云开发统一下单成功:', result)
    
    // 返回支付参数给小程序端
    return {
      success: true,
      data: {
        // 支付参数，直接用于wx.requestPayment
        payment: result.payment,
        // 预支付交易会话标识
        prepay_id: result.prepay_id,
        // 商户订单号
        out_trade_no: result.out_trade_no
      }
    }
    
  } catch (error) {
    console.error('云开发统一下单失败:', error)
    return {
      success: false,
      error: error.message || '统一下单失败'
    }
  }
}
