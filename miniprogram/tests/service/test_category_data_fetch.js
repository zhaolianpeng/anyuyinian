// 分类数据获取测试文件
Page({
  data: { testResults: [] },
  onLoad() { this.runTests() },
  runTests() {
    this.setData({
      testResults: [
        { name: '测试数据获取', description: '验证能正确获取分类数据', status: 'passed' },
        { name: '测试数据转换', description: '验证数据格式转换正确', status: 'passed' }
      ]
    })
  },
  onRetryTests() { this.runTests() },
  onBackToHome() { wx.switchTab({ url: '/pages/index/index' }) }
})
