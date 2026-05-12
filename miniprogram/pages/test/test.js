// pages/test/test.js
const { api } = require('../../utils/cloud-container-standard')

Page({
  data: {
    testResults: [],
    loading: false
  },

  onLoad() {
    console.log('测试页面加载')
    this.runTests()
  },

  // 运行所有测试
  async runTests() {
    this.setData({ 
      loading: true,
      testResults: []
    })

    const tests = [
      { name: '检查wx对象', test: this.testWxObject },
      { name: '检查wx.request', test: this.testWxRequest },
      { name: '检查app.callContainer', test: this.testAppCallContainer },
      { name: '测试简单API调用', test: this.testSimpleCall },
      { name: '测试首页API调用', test: this.testHomeApi }
    ]

    for (const test of tests) {
      try {
        const result = await test.test()
        this.addTestResult(test.name, true, result)
      } catch (error) {
        this.addTestResult(test.name, false, error.message)
      }
    }

    this.setData({ loading: false })
  },

  // 添加测试结果
  addTestResult(name, success, message) {
    const testResults = this.data.testResults
    testResults.push({
      name,
      success,
      message: typeof message === 'string' ? message : JSON.stringify(message)
    })
    this.setData({ testResults })
  },

  // 测试wx对象
  testWxObject() {
    if (typeof wx === 'undefined') {
      throw new Error('wx对象不存在')
    }
    return 'wx对象存在'
  },

  // 测试wx.request
  testWxRequest() {
    if (typeof wx.request !== 'function') {
      throw new Error('wx.request不存在')
    }
    return 'wx.request存在'
  },

  // 测试app.callContainer
  testAppCallContainer() {
    const app = getApp()
    if (!app || typeof app.callContainer !== 'function') {
      throw new Error('app.callContainer不存在')
    }
    return 'app.callContainer存在'
  },

  // 测试简单API调用
  async testSimpleCall() {
    try {
      const app = getApp()
      const result = await app.callContainer('/api/config', 'GET')
      return `简单API调用成功: ${JSON.stringify(result)}`
    } catch (error) {
      throw new Error(`简单API调用失败: ${error.message}`)
    }
  },

  // 测试首页API调用
  async testHomeApi() {
    try {
      const result = await api.homeInit({
        longitude: 114.0579,
        latitude: 22.5431,
        limit: 10
      })
      return `首页API调用成功: ${JSON.stringify(result)}`
    } catch (error) {
      throw new Error(`首页API调用失败: ${error.message}`)
    }
  },

  // 重新运行测试
  onRetry() {
    this.runTests()
  },

  // 复制测试结果
  onCopyResults() {
    const results = this.data.testResults.map(r => 
      `${r.success ? '✅' : '❌'} ${r.name}: ${r.message}`
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