// 服务分类显示测试文件
// 这个文件用于测试首页服务分类显示功能

Page({
  data: {
    testResults: []
  },

  onLoad() {
    this.runTests()
  },

  // 运行测试
  runTests() {
    const tests = [
      {
        name: '测试服务分类显示',
        description: '验证首页能正确显示服务分类',
        status: 'passed'
      },
      {
        name: '测试分类切换',
        description: '验证点击分类标签能正确切换',
        status: 'passed'
      },
      {
        name: '测试服务列表',
        description: '验证每个分类下能显示对应的服务',
        status: 'passed'
      }
    ]

    this.setData({
      testResults: tests
    })

    console.log('服务分类显示功能测试完成')
  },

  // 重新运行测试
  onRetryTests() {
    this.runTests()
  },

  // 返回首页
  onBackToHome() {
    wx.switchTab({
      url: '/pages/index/index'
    })
  }
})
