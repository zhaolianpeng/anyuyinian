// 登录调试脚本
const debugLogin = {
  // 测试getUserProfile调用
  testGetUserProfile() {
    console.log('开始测试getUserProfile...')
    
    return new Promise((resolve, reject) => {
      wx.getUserProfile({
        desc: '用于完善用户资料和提供个性化服务',
        lang: 'zh_CN',
        success: (res) => {
          console.log('✅ getUserProfile调用成功:', res)
          resolve(res)
        },
        fail: (err) => {
          console.error('❌ getUserProfile调用失败:', err)
          reject(err)
        }
      })
    })
  },

  // 测试错误处理
  testErrorHandling(errorMsg) {
    console.log('测试错误处理:', errorMsg)
    
    if (errorMsg.includes('getUserProfile:fail can only be invoked by user TAP gesture')) {
      console.log('✅ 检测到TAP gesture错误')
      return '请点击"微信登录"按钮进行登录，不能通过其他方式触发'
    } else if (errorMsg.includes('getUserProfile:fail user deny')) {
      console.log('✅ 检测到用户拒绝授权')
      return '需要授权才能登录，请重试'
    } else if (errorMsg.includes('getUserProfile:fail')) {
      console.log('✅ 检测到其他getUserProfile错误')
      return '获取用户信息失败，请重试'
    } else {
      console.log('✅ 检测到其他错误')
      return '登录失败，请重试'
    }
  },

  // 模拟各种错误情况
  simulateErrors() {
    const testErrors = [
      'getUserProfile:fail can only be invoked by user TAP gesture.',
      'getUserProfile:fail user deny',
      'getUserProfile:fail network error',
      'wx.login:fail network error'
    ]
    
    testErrors.forEach(error => {
      const result = this.testErrorHandling(error)
      console.log(`错误: ${error} -> 处理结果: ${result}`)
    })
  }
}

// 导出调试对象
module.exports = debugLogin

// 如果直接运行此文件，执行测试
if (typeof wx !== 'undefined') {
  console.log('在微信小程序环境中，可以调用 debugLogin.simulateErrors() 进行测试')
} 