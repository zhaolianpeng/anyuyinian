/**
 * 历史兼容配置
 * 该文件仅保留部分旧字段，当前小程序主链已统一走服务端直连。
 */

// 历史配置
const CLOUD_CONFIG = {
  // 历史环境ID
  env: 'prod-5g94mx7a3d07e78c',
  // 历史服务名称
  service: 'golang-lfwy',
  // 小程序AppID
  appId: 'wx101090677bd5219e',
  // 默认超时时间（毫秒）
  timeout: 10000,
  // WebSocket路径
  websocketPath: '/ws'
}

// 标准请求头
const DEFAULT_HEADERS = {
  'X-WX-SERVICE': CLOUD_CONFIG.service,
  'content-type': 'application/json'
}

// 标准调用配置
const getCallConfig = (path, method = 'GET', data = {}, options = {}) => {
  return {
    config: {
      env: CLOUD_CONFIG.env
    },
    path: path,
    method: method,
    data: data,
    header: {
      ...DEFAULT_HEADERS,
      ...options.header
    },
    timeout: options.timeout || CLOUD_CONFIG.timeout,
    ...options
  }
}

// 标准WebSocket连接配置
const getWebSocketConfig = (path = CLOUD_CONFIG.websocketPath, options = {}) => {
  return {
    config: {
      env: CLOUD_CONFIG.env
    },
    service: CLOUD_CONFIG.service,
    path: path,
    timeout: options.timeout || CLOUD_CONFIG.timeout,
    ...options
  }
}

// 错误处理函数
const handleCloudError = (error) => {
  console.error('服务端调用错误:', error)
  
  let message = '请求失败'
  if (error.message) {
    if (error.message.includes('timeout')) {
      message = '请求超时，请检查网络'
    } else if (error.message.includes('fail')) {
      message = '服务暂时不可用，请稍后重试'
    } else if (error.message.includes('service')) {
      message = '服务端异常，请稍后重试'
    } else if (error.message.includes('network')) {
      message = '网络连接异常，请检查网络'
    } else {
      message = error.message
    }
  }
  
  // 显示错误提示
  if (typeof wx !== 'undefined') {
    wx.showToast({
      title: message,
      icon: 'none',
      duration: 2000
    })
  }
  
  return error
}

// 网络状态检查
const checkNetworkStatus = () => {
  return new Promise((resolve) => {
    if (typeof wx === 'undefined') {
      resolve(false)
      return
    }
    
    wx.getNetworkType({
      success: (res) => {
        console.log('网络类型:', res.networkType)
        resolve(res.networkType !== 'none')
      },
      fail: () => {
        resolve(false)
      }
    })
  })
}

// 重试机制
const callWithRetry = async (apiCall, maxRetries = 3, delay = 1000) => {
  for (let i = 0; i < maxRetries; i++) {
    try {
      return await apiCall()
    } catch (error) {
      if (i === maxRetries - 1) {
        throw error
      }
      console.log(`第${i + 1}次重试...`)
      await new Promise(resolve => setTimeout(resolve, delay * (i + 1)))
    }
  }
}

module.exports = {
  CLOUD_CONFIG,
  DEFAULT_HEADERS,
  getCallConfig,
  getWebSocketConfig,
  handleCloudError,
  checkNetworkStatus,
  callWithRetry
} 