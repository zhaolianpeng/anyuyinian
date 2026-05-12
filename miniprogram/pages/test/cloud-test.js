// 服务端连通性测试页面
const app = getApp()
const { api } = require('../../utils/cloud-container-standard')

Page({
  data: {
    testResults: [],
    loading: false
  },

  onLoad() {
    console.log('服务端连通性测试页面加载')
  },

  // 测试基本连接
  async testBasicConnection() {
    this.addTestResult('开始测试基本连接...')
    
    try {
      const result = await app.callContainer('/', 'GET')
      this.addTestResult('✅ 基本连接测试成功', result)
    } catch (error) {
      this.addTestResult('❌ 基本连接测试失败', error.message)
    }
  },

  // 测试计数器API
  async testCountAPI() {
    this.addTestResult('开始测试计数器API...')
    
    try {
      const result = await api.count.get()
      this.addTestResult('✅ 计数器API测试成功', result)
    } catch (error) {
      this.addTestResult('❌ 计数器API测试失败', error.message)
    }
  },

  // 测试首页API
  async testHomeAPI() {
    this.addTestResult('开始测试首页API...')
    
    try {
      const result = await api.homeInit({
        longitude: 121.4737,
        latitude: 31.2304,
        limit: 5
      })
      this.addTestResult('✅ 首页API测试成功', result)
    } catch (error) {
      this.addTestResult('❌ 首页API测试失败', error.message)
    }
  },

  // 测试订单列表API
  async testOrderListAPI() {
    this.addTestResult('开始测试订单列表API...')
    
    try {
      const result = await api.orderList({
        page: 1,
        pageSize: 5
      })
      this.addTestResult('✅ 订单列表API测试成功', result)
    } catch (error) {
      this.addTestResult('❌ 订单列表API测试失败', error.message)
    }
  },

  // 测试服务列表API
  async testServiceListAPI() {
    this.addTestResult('开始测试服务列表API...')
    
    try {
      const result = await api.serviceList({
        page: 1,
        pageSize: 5
      })
      this.addTestResult('✅ 服务列表API测试成功', result)
    } catch (error) {
      this.addTestResult('❌ 服务列表API测试失败', error.message)
    }
  },

  // 测试用户信息API
  async testUserInfoAPI() {
    this.addTestResult('开始测试用户信息API...')
    
    try {
      const result = await api.userInfo()
      this.addTestResult('✅ 用户信息API测试成功', result)
    } catch (error) {
      this.addTestResult('❌ 用户信息API测试失败', error.message)
    }
  },

  // 测试统一服务端调用环境
  testCloudContainerStatus() {
    this.addTestResult('开始测试统一服务端调用环境...')
    
    try {
      if (typeof wx.request === 'function') {
        this.addTestResult('✅ wx.request 可用')
      } else {
        this.addTestResult('❌ wx.request 不可用')
      }
    } catch (error) {
      this.addTestResult('❌ 统一服务端调用环境检查失败', error.message)
    }
  },

  // 测试错误处理
  async testErrorHandling() {
    this.addTestResult('开始测试错误处理...')
    
    try {
      // 故意调用一个不存在的API
      await app.callContainer('/api/nonexistent', 'GET')
      this.addTestResult('❌ 错误处理测试失败 - 应该返回错误')
    } catch (error) {
      this.addTestResult('✅ 错误处理正常', error.message)
    }
  },

  // 运行所有测试
  async runAllTests() {
    this.setData({ loading: true, testResults: [] })
    
    console.log('开始运行所有服务端连通性测试...')
    
    // 测试1: 统一服务端调用环境
    this.testCloudContainerStatus()
    await this.delay(500)
    
    // 测试2: 基本连接
    await this.testBasicConnection()
    await this.delay(500)
    
    // 测试3: 计数器API
    await this.testCountAPI()
    await this.delay(500)
    
    // 测试4: 首页API
    await this.testHomeAPI()
    await this.delay(500)
    
    // 测试5: 订单列表API
    await this.testOrderListAPI()
    await this.delay(500)
    
    // 测试6: 服务列表API
    await this.testServiceListAPI()
    await this.delay(500)
    
    // 测试7: 用户信息API
    await this.testUserInfoAPI()
    await this.delay(500)
    
    // 测试8: 错误处理
    await this.testErrorHandling()
    
    this.setData({ loading: false })
    
    console.log('所有测试完成')
    
    // 显示测试结果
    const successCount = this.data.testResults.filter(r => r.message.includes('✅')).length
    const totalCount = this.data.testResults.length
    
    wx.showModal({
      title: '测试完成',
      content: `成功: ${successCount}/${totalCount}`,
      showCancel: false
    })
  },

  // 添加测试结果
  addTestResult(message, data = null) {
    const testResults = [...this.data.testResults]
    testResults.push({
      message,
      data,
      timestamp: new Date().toLocaleTimeString()
    })
    
    this.setData({ testResults })
    
    console.log('测试结果:', message, data)
  },

  // 延迟函数
  delay(ms) {
    return new Promise(resolve => setTimeout(resolve, ms))
  },

  // 清除测试结果
  clearTestResults() {
    this.setData({ testResults: [] })
  },

  // 复制测试结果
  copyTestResults() {
    const results = this.data.testResults.map(r => 
      `${r.timestamp}: ${r.message}${r.data ? ` - ${JSON.stringify(r.data)}` : ''}`
    ).join('\n')
    
    wx.setClipboardData({
      data: results,
      success: () => {
        wx.showToast({
          title: '测试结果已复制',
          icon: 'success'
        })
      }
    })
  }
}) 