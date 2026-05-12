// 调试工具

// 检查统一服务端调用环境
const checkServerRequestEnv = () => {
  console.log('=== 统一服务端调用环境检查 ===')

  if (typeof wx === 'undefined') {
    console.error('❌ wx 不存在')
    return false
  }

  if (typeof wx.request !== 'function') {
    console.error('❌ wx.request 不存在')
    return false
  }

  console.log('✅ wx.request 可用')
  return true
}

// 测试统一服务端连接
const testServerConnection = () => {
  return new Promise((resolve, reject) => {
    console.log('=== 测试统一服务端连接 ===')
    
    if (!checkServerRequestEnv()) {
      reject(new Error('统一服务端环境检查失败'))
      return
    }
    
    // 获取app实例
    const app = getApp()
    if (!app) {
      reject(new Error('无法获取app实例'))
      return
    }
    
    // 测试简单接口
    app.callContainer('/api/count', 'GET')
      .then(res => {
        console.log('✅ 统一服务端连接测试成功:', res)
        resolve(res)
      })
      .catch(err => {
        console.error('❌ 统一服务端连接测试失败:', err)
        reject(err)
      })
  })
}

// 检查网络状态
const checkNetwork = () => {
  return new Promise((resolve) => {
    console.log('=== 网络状态检查 ===')
    
    wx.getNetworkType({
      success: (res) => {
        console.log('网络类型:', res.networkType)
        console.log('是否连接:', res.networkType !== 'none')
        resolve({
          isConnected: res.networkType !== 'none',
          networkType: res.networkType
        })
      },
      fail: () => {
        console.error('❌ 无法获取网络状态')
        resolve({
          isConnected: false,
          networkType: 'unknown'
        })
      }
    })
  })
}

// 检查微信环境
const checkWechatEnv = () => {
  console.log('=== 微信环境检查 ===')
  
  const checks = {
    'wx对象存在': typeof wx !== 'undefined',
    'wx.request存在': typeof wx !== 'undefined' && typeof wx.request === 'function'
  }
  
  Object.entries(checks).forEach(([name, result]) => {
    console.log(`${result ? '✅' : '❌'} ${name}`)
  })
  
  return Object.values(checks).every(Boolean)
}

// 检查app实例
const checkAppInstance = () => {
  console.log('=== App实例检查 ===')
  
  const app = getApp()
  if (!app) {
    console.error('❌ 无法获取app实例')
    return false
  }
  
  const checks = {
    'app.callContainer方法存在': typeof app.callContainer === 'function'
  }
  
  Object.entries(checks).forEach(([name, result]) => {
    console.log(`${result ? '✅' : '❌'} ${name}`)
  })
  
  return Object.values(checks).every(Boolean)
}

// 完整诊断
const runDiagnostics = async () => {
  console.log('=== 开始诊断 ===')
  
  try {
    // 1. 检查微信环境
    const wechatOk = checkWechatEnv()
    if (!wechatOk) {
      throw new Error('微信环境检查失败')
    }
    
    // 2. 检查网络状态
    const networkStatus = await checkNetwork()
    if (!networkStatus.isConnected) {
      throw new Error('网络未连接')
    }
    
    // 3. 检查app实例
    const appOk = checkAppInstance()
    if (!appOk) {
      throw new Error('App实例检查失败')
    }
    
    // 4. 检查统一服务端环境
    const serverOk = checkServerRequestEnv()
    if (!serverOk) {
      throw new Error('统一服务端环境检查失败')
    }
    
    // 5. 测试统一服务端连接
    await testServerConnection()
    
    console.log('✅ 所有检查通过')
    return true
    
  } catch (error) {
    console.error('❌ 诊断失败:', error.message)
    return false
  }
}

// 获取详细错误信息
const getDetailedError = (error) => {
  console.log('=== 错误详情 ===')
  console.log('错误类型:', error.constructor.name)
  console.log('错误消息:', error.message)
  console.log('错误堆栈:', error.stack)
  
  if (error.errMsg) {
    console.log('微信错误消息:', error.errMsg)
  }
  
  if (error.errCode) {
    console.log('微信错误码:', error.errCode)
  }
  
  return {
    type: error.constructor.name,
    message: error.message,
    wxErrMsg: error.errMsg,
    wxErrCode: error.errCode,
    stack: error.stack
  }
}

module.exports = {
  checkServerRequestEnv,
  testServerConnection,
  checkNetwork,
  checkWechatEnv,
  checkAppInstance,
  runDiagnostics,
  getDetailedError
} 