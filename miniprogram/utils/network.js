// 网络状态检测工具
const checkNetworkStatus = () => {
  return new Promise((resolve) => {
    wx.getNetworkType({
      success: (res) => {
        console.log('当前网络类型:', res.networkType)
        resolve({
          isConnected: res.networkType !== 'none',
          networkType: res.networkType
        })
      },
      fail: () => {
        resolve({
          isConnected: false,
          networkType: 'unknown'
        })
      }
    })
  })
}

// 监听网络状态变化
const onNetworkStatusChange = (callback) => {
  wx.onNetworkStatusChange((res) => {
    console.log('网络状态变化:', res)
    callback && callback(res)
  })
}

// 检查是否在微信环境中
const isWechatEnv = () => {
  return typeof wx !== 'undefined'
}

// 获取系统信息
const getSystemInfo = () => {
  return new Promise((resolve) => {
    wx.getSystemInfo({
      success: (res) => {
        resolve(res)
      },
      fail: () => {
        resolve(null)
      }
    })
  })
}

// 显示网络错误提示
const showNetworkError = (message = '网络连接异常，请检查网络设置') => {
  wx.showModal({
    title: '网络错误',
    content: message,
    showCancel: false,
    confirmText: '确定'
  })
}

// 检查网络并提示
const checkNetworkAndAlert = async () => {
  const networkStatus = await checkNetworkStatus()
  
  if (!networkStatus.isConnected) {
    showNetworkError('当前无网络连接，请检查网络设置')
    return false
  }
  
  if (networkStatus.networkType === '2g') {
    wx.showToast({
      title: '当前为2G网络，可能影响使用体验',
      icon: 'none',
      duration: 2000
    })
  }
  
  return true
}

module.exports = {
  checkNetworkStatus,
  onNetworkStatusChange,
  isWechatEnv,
  getSystemInfo,
  showNetworkError,
  checkNetworkAndAlert
} 