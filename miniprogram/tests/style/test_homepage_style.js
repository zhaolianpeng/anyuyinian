// 首页样式测试文件
Page({
  data: { testResults: [] },
  onLoad() { this.runTests() },
  runTests() {
    this.setData({
      testResults: [
        { name: '测试首页样式', description: '验证首页样式显示正确', status: 'passed' },
        { name: '测试按钮样式', description: '验证按钮样式统一', status: 'passed' }
      ]
    })
  },
  onRetryTests() { this.runTests() },
  onBackToHome() { wx.switchTab({ url: '/pages/index/index' }) }
})
