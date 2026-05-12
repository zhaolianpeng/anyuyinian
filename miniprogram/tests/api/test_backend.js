// 后端服务测试脚本
const testBackend = {
  // 测试后端服务是否可用
  testBackendService() {
    console.log('开始测试后端服务...')
    
    const { baseURL } = require('./config')
    console.log('当前使用的后端地址:', baseURL)
    
    return new Promise((resolve, reject) => {
      wx.request({
        url: baseURL + '/config',
        method: 'GET',
        success: (res) => {
          console.log('✅ 后端服务测试成功:', res)
          resolve(res)
        },
        fail: (err) => {
          console.error('❌ 后端服务测试失败:', err)
          reject(err)
        }
      })
    })
  },

  // 测试登录接口
  testLoginAPI(code, userInfo) {
    console.log('开始测试登录接口...')
    
    const { baseURL } = require('./config')
    console.log('登录接口地址:', baseURL + '/wx/login')
    console.log('请求参数:', { code, userInfo })
    
    return new Promise((resolve, reject) => {
      wx.request({
        url: baseURL + '/wx/login',
        method: 'POST',
        data: {
          code: code,
          userInfo: userInfo
        },
        success: (res) => {
          console.log('✅ 登录接口测试成功:', res)
          resolve(res)
        },
        fail: (err) => {
          console.error('❌ 登录接口测试失败:', err)
          reject(err)
        }
      })
    })
  },

  // 测试网络连接
  testNetwork() {
    console.log('开始测试网络连接...')
    
    return new Promise((resolve, reject) => {
      wx.request({
        url: 'https://www.baidu.com',
        method: 'GET',
        success: (res) => {
          console.log('✅ 网络连接正常:', res.statusCode)
          resolve(res)
        },
        fail: (err) => {
          console.error('❌ 网络连接失败:', err)
          reject(err)
        }
      })
    })
  },

  // 完整的后端测试
  async runFullTest() {
    console.log('开始完整的后端测试...')
    
    try {
      // 1. 测试网络连接
      console.log('1. 测试网络连接...')
      await this.testNetwork()
      
      // 2. 测试后端服务
      console.log('2. 测试后端服务...')
      await this.testBackendService()
      
      console.log('✅ 所有测试通过')
      return true
    } catch (error) {
      console.error('❌ 测试失败:', error)
      return false
    }
  }
}

// 导出测试对象
module.exports = testBackend

// 如果直接运行此文件，执行测试
if (typeof wx !== 'undefined') {
  console.log('在微信小程序环境中，可以调用 testBackend.runFullTest() 进行测试')
} 