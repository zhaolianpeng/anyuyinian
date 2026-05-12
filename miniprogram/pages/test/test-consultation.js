Page({
  data: {
    testResults: []
  },

  onLoad() {
    this.runTests()
  },

  // 运行测试
  runTests() {
    this.addTestResult('开始测试咨询功能...')
    
    // 测试1：检查页面路径
    this.testPagePath()
    
    // 测试2：测试API导入
    this.testAPIImport()
    
    // 测试3：测试页面跳转
    this.testPageNavigation()
  },

  // 测试页面路径
  testPagePath() {
    this.addTestResult('测试1：检查咨询页面路径')
    try {
      const pages = getCurrentPages()
      this.addTestResult(`当前页面路径: ${pages[pages.length - 1].route}`)
      this.addTestResult('✅ 页面路径检查通过')
    } catch (error) {
      this.addTestResult(`❌ 页面路径检查失败: ${error.message}`)
    }
  },

  // 测试API导入
  testAPIImport() {
    this.addTestResult('测试2：检查API导入')
    try {
      const { consultationAPI, api } = require('../../utils/request.js')
      if (consultationAPI && api) {
        this.addTestResult('✅ API导入成功')
        this.addTestResult(`咨询API方法数量: ${Object.keys(consultationAPI).length}`)
      } else {
        this.addTestResult('❌ API导入失败')
      }
    } catch (error) {
      this.addTestResult(`❌ API导入异常: ${error.message}`)
    }
  },

  // 测试页面跳转
  testPageNavigation() {
    this.addTestResult('测试3：测试页面跳转')
    try {
      wx.navigateTo({
        url: '/pages/consultation/consultation',
        success: () => {
          this.addTestResult('✅ 页面跳转成功')
        },
        fail: (err) => {
          this.addTestResult(`❌ 页面跳转失败: ${err.errMsg}`)
        }
      })
    } catch (error) {
      this.addTestResult(`❌ 页面跳转异常: ${error.message}`)
    }
  },

  // 添加测试结果
  addTestResult(message) {
    const timestamp = new Date().toLocaleTimeString()
    const result = `[${timestamp}] ${message}`
    this.setData({
      testResults: [...this.data.testResults, result]
    })
    console.log(result)
  },

  // 重新运行测试
  onRetryTests() {
    this.setData({ testResults: [] })
    this.runTests()
  },

  // 返回首页
  onBackToHome() {
    wx.navigateBack()
  }
})
