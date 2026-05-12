// 简单的登录测试脚本
const testLogin = {
  // 测试getUserProfile调用
  testGetUserProfile() {
    console.log('开始测试getUserProfile...')
    
    return new Promise((resolve, reject) => {
      wx.getUserProfile({
        desc: '用于完善用户资料',
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

  // 测试完整的登录流程
  async testFullLogin() {
    console.log('开始测试完整登录流程...')
    
    try {
      // 1. 获取登录code
      console.log('1. 获取微信登录code...')
      const loginRes = await this.getWxLoginCode()
      console.log('✅ 登录code获取成功:', loginRes.code)
      
      // 2. 获取用户信息
      console.log('2. 获取用户信息...')
      const userInfoRes = await this.testGetUserProfile()
      console.log('✅ 用户信息获取成功:', userInfoRes.userInfo.nickName)
      
      // 3. 调用后端登录接口
      console.log('3. 调用后端登录接口...')
      const res = await this.callLoginAPI(loginRes.code, userInfoRes.userInfo)
      console.log('✅ 后端登录成功:', res)
      
      return res
    } catch (error) {
      console.error('❌ 登录流程测试失败:', error)
      throw error
    }
  },

  // 获取微信登录code
  getWxLoginCode() {
    return new Promise((resolve, reject) => {
      wx.login({
        success: (res) => {
          console.log('wx.login成功:', res)
          resolve(res)
        },
        fail: (err) => {
          console.error('wx.login失败:', err)
          reject(err)
        }
      })
    })
  },

  // 调用后端登录接口
  callLoginAPI(code, userInfo) {
    return new Promise((resolve, reject) => {
      wx.request({
        url: 'https://golang-lfwy-176496-6-1353115175.sh.run.tcloudbase.com/api/wx/login',
        method: 'POST',
        data: {
          code: code,
          userInfo: userInfo
        },
        success: (res) => {
          console.log('后端登录接口响应:', res.data)
          resolve(res.data)
        },
        fail: (err) => {
          console.error('后端登录接口失败:', err)
          reject(err)
        }
      })
    })
  }
}

// 导出测试对象
module.exports = testLogin

// 如果直接运行此文件，执行测试
if (typeof wx !== 'undefined') {
  console.log('在微信小程序环境中，可以调用 testLogin.testFullLogin() 进行测试')
} 