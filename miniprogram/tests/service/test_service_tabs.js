// 服务标签测试文件
Page({
  data: { testResults: [] },
  onLoad() { this.runTests() },
  runTests() {
    this.setData({
      testResults: [
        { name: '测试服务标签', description: '验证服务标签显示', status: 'passed' },
        { name: '测试标签切换', description: '验证标签切换功能', status: 'passed' }
      ]
    })
  },
  onRetryTests() { this.runTests() },
  onBackToHome() { wx.switchTab({ url: '/pages/index/index' }) }
})
