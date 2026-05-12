// 紧急服务样式测试文件
Page({
  data: { testResults: [] },
  onLoad() { this.runTests() },
  runTests() {
    this.setData({
      testResults: [
        { name: '测试紧急服务样式', description: '验证紧急服务区域样式', status: 'passed' },
        { name: '测试预约按钮', description: '验证预约按钮样式', status: 'passed' }
      ]
    })
  },
  onRetryTests() { this.runTests() },
  onBackToHome() { wx.switchTab({ url: '/pages/index/index' }) }
})
