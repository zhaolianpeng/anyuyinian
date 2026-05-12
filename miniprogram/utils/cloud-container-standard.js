/**
 * 统一服务端调用工具
 * 所有业务接口统一走用户自己的服务端
 */

const config = require('../config')

const SERVER_BASE_URL = config.cos.bucketDomain

// 云托管配置
const CONTAINER_CONFIG = {
  env: 'prod-5g94mx7a3d07e78c',
  service: 'golang-lfwy'
}

function appendQueryParams(url, query = {}) {
  const entries = Object.entries(query).filter(([, value]) => value !== undefined && value !== null)
  if (entries.length === 0) {
    return url
  }

  const queryString = entries
    .map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(value)}`)
    .join('&')

  return `${url}${url.includes('?') ? '&' : '?'}${queryString}`
}

/**
 * 标准服务端HTTP调用
 * @param {string} path - API路径
 * @param {string} method - HTTP方法
 * @param {object} data - 请求数据
 * @param {object} options - 额外选项
 */
const callContainer = (path, method = 'GET', data = {}, options = {}) => {
  return new Promise((resolve, reject) => {
    if (typeof wx === 'undefined') {
      reject(new Error('不在微信小程序环境中'))
      return
    }

    if (!wx.request) {
      reject(new Error('当前环境不支持网络请求'))
      return
    }

    const { header = {}, query, ...restOptions } = options
    const requestMethod = method === 'GET' ? 'POST' : method
    const requestUrl = requestMethod === 'POST'
      ? `${SERVER_BASE_URL}${path}`
      : appendQueryParams(`${SERVER_BASE_URL}${path}`, query)
    const requestData = requestMethod === 'POST' && method === 'GET'
      ? {
          ...(data || {}),
          ...(query || {})
        }
      : data

    console.log('服务端调用参数:', {
      url: requestUrl,
      method: requestMethod,
      originalMethod: method,
      data: requestData,
      options: restOptions
    })

    const requestConfig = {
      url: requestUrl,
      method: requestMethod,
      data: requestData,
      header: {
        'content-type': 'application/json',
        ...(method === 'GET' ? { 'X-HTTP-Method-Override': 'GET' } : {}),
        ...header
      },
      ...restOptions
    }

    wx.request({
      ...requestConfig,
      success: (result) => {
        console.log('服务端调用成功:', result)

        if (result.statusCode === 200) {
          resolve(result.data)
          return
        }

        const responseMessage = (result.data && result.data.message) || result.data || '请求失败'
        const errorMsg = `HTTP ${result.statusCode}: ${responseMessage}`
        console.error('服务端调用失败:', errorMsg)
        console.error('失败详情:', {
          statusCode: result.statusCode,
          data: result.data,
          header: result.header
        })
        reject(new Error(errorMsg))
      },
      fail: (error) => {
        console.error('服务端调用异常:', error)
        console.error('异常详情:', {
          message: error.message,
          errMsg: error.errMsg,
          errCode: error.errCode
        })

        let errorMessage = error.message || '网络请求失败'
        if (error.errMsg) {
          if (error.errMsg.includes('timeout')) {
            errorMessage = '请求超时，请检查网络连接'
          } else if (error.errMsg.includes('fail')) {
            errorMessage = '网络连接失败，请检查网络设置'
          } else if (error.errMsg.includes('ssl') || error.errMsg.includes('certificate')) {
            errorMessage = '服务端证书异常，请检查域名与证书配置'
          }
        }

        reject(new Error(errorMessage))
      }
    })
  })
}

/**
 * API接口封装
 */
const api = {
  // 计数器相关
  count: {
    get: () => callContainer('/api/count', 'GET'),
    increment: () => callContainer('/api/count', 'POST', { action: 'inc' }),
    decrement: () => callContainer('/api/count', 'POST', { action: 'dec' })
  },

  // 用户认证
  wxLogin: (data) => callContainer('/api/wx/login', 'POST', data),
  
  // 首页相关
  homeInit: (params) => callContainer('/api/home/init', 'GET', params),
  mediaImages: (data) => callContainer('/api/media/images', 'POST', data),
  
  // 用户管理
  userInfo: (params) => callContainer('/api/user/info', 'GET', params),
  bindPhone: (data) => callContainer('/api/user/bind_phone', 'POST', data),
  updateUserInfo: (data) => callContainer('/api/user/update_info', 'POST', data),
  decryptPhoneNumber: (data) => callContainer('/api/user/decrypt_phone', 'POST', data),
  userAddress: (params) => callContainer('/api/user/address', 'GET', params),
  userAddressAdd: (data) => callContainer('/api/user/address', 'POST', data),
  userAddressUpdate: (data) => callContainer('/api/user/address', 'PUT', data),
  userAddressDelete: (id) => callContainer(`/api/user/address/${id}`, 'DELETE'),
  userPatient: (params) => callContainer('/api/user/patient', 'GET', params),
  userPatientAdd: (data) => callContainer('/api/user/patient', 'POST', data),
  userPatientUpdate: (data) => callContainer('/api/user/patient', 'PUT', data),
  userPatientDelete: (id) => callContainer(`/api/user/patient/${id}`, 'DELETE'),
  
  // 服务管理
  serviceList: (params) => callContainer('/api/service/list', 'GET', params),
  serviceCategories: () => callContainer('/api/service/categories', 'GET'),
  serviceDetail: (data) => callContainer('/api/service/detail', 'POST', data),
  serviceFormConfig: (id) => callContainer(`/api/service/form_config/${id}`, 'GET'),
  
  // 订单管理
  orderSubmit: (data) => callContainer('/api/order/submit', 'POST', data),
  smartElderlyOrderSubmit: (data) => callContainer('/api/order/smart-elderly', 'POST', data),
  orderPay: (orderId, data) => callContainer(`/api/order/pay/${orderId}`, 'POST', data),
  orderPayConfirm: (orderId, data) => callContainer(`/api/order/pay_confirm/${orderId}`, 'POST', data),
  orderCancel: (orderId, data) => callContainer(`/api/order/cancel/${orderId}`, 'POST', data),
  orderRefund: (orderId, data) => callContainer(`/api/order/refund/${orderId}`, 'POST', data),
  orderList: (params) => callContainer('/api/order/list', 'GET', params),
  orderDetail: (data) => callContainer('/api/order/detail', 'POST', data),
  orderDetailById: (orderId) => callContainer(`/api/order/detail/${orderId}`, 'GET'),
  orderTimeSlots: (params) => callContainer('/api/order/time_slots', 'POST', params),
  
  // 推荐系统
  referralQrcode: (params) => callContainer('/api/referral/qrcode', 'GET', params),
  referralReport: (params) => callContainer('/api/referral/report', 'GET', params),
  referralConfig: () => callContainer('/api/referral/config', 'GET'),
  applyCashout: (data) => callContainer('/api/referral/apply_cashout', 'POST', data),
  
  // 推广中心
  promoterInfo: (params) => callContainer('/api/promoter/info', 'GET', params),
  commissionList: (params) => callContainer('/api/promoter/commission_list', 'GET', params),
  cashoutList: (params) => callContainer('/api/promoter/cashout_list', 'GET', params),
  
  // 客服医院
  kefuSendMsg: (data) => callContainer('/api/kefu/send_msg', 'POST', data),
  kefuFaq: () => callContainer('/api/kefu/faq', 'GET'),
  hospitalList: (params) => callContainer('/api/hospital/list', 'GET', params),
  hospitalDetail: (id) => callContainer(`/api/hospital/detail/${id}`, 'GET'),
  
  // 文件管理
  upload: (data) => callContainer('/api/upload', 'POST', data),
  fileList: (params) => callContainer('/api/files', 'GET', params),
  deleteFile: (fileId) => callContainer('/api/file/delete', 'DELETE', {}, { fileId }),
  updateFilePermission: (data) => callContainer('/api/file/permission', 'PUT', data),
  getFilePermission: (params) => callContainer('/api/file/permission/get', 'GET', params),
  
  // 二维码生成
  generateQRCode: (params) => callContainer('/api/qrcode/generate', 'GET', {}, { query: params }),
  generateQRCodeBase64: (params) => callContainer('/api/qrcode/generate_base64', 'GET', {}, { query: params }),
  
  // 系统配置
  config: () => callContainer('/api/config', 'GET'),
  
  // 管理员相关接口
  adminLogin: (data) => callContainer('/api/admin/login', 'POST', data),
  adminCheckStatus: (data) => callContainer('/api/admin/check-status', 'POST', data),
  adminUsers: (params) => callContainer('/api/admin/users', 'GET', params),
  adminOrders: (params) => callContainer('/api/admin/orders', 'GET', params),
  adminStats: (params) => callContainer('/api/admin/stats', 'GET', params),
  adminAdmins: (params) => callContainer('/api/admin/admins', 'GET', params),
  adminSetAdmin: (data) => callContainer('/api/admin/set-admin', 'POST', data),
  adminRemoveAdmin: (data) => callContainer('/api/admin/remove-admin', 'POST', data),
  adminUpdateOrderAmount: (data) => callContainer('/api/admin/order/update-amount', 'POST', data),
  adminRefundOrder: (data) => callContainer('/api/admin/order/refund', 'POST', data),
  adminServices: (params) => callContainer('/api/admin/services', 'GET', params),
  adminUpdateServicePrice: (data) => callContainer('/api/admin/service/updateprice', 'POST', data),
  
  // 咨询相关接口
  consultationStats: () => callContainer('/api/consultation/stats', 'GET')
}

module.exports = {
  callContainer,
  api,
  CONTAINER_CONFIG,
  SERVER_BASE_URL
} 