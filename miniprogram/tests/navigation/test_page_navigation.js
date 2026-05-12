// 测试页面跳转功能
function testPageNavigation() {
  console.log('=== 测试页面跳转功能 ===')
  
  try {
    // 1. 检查用户登录状态
    const userId = wx.getStorageSync('userId')
    console.log('当前用户ID:', userId)
    
    if (!userId) {
      console.log('❌ 用户未登录')
      wx.showToast({
        title: '请先登录',
        icon: 'none'
      })
      return
    }
    
    // 2. 测试跳转到设置页面
    console.log('准备跳转到设置页面...')
    wx.navigateTo({
      url: '/pages/user/setup-profile',
      success: () => {
        console.log('✅ 跳转成功')
        wx.showToast({
          title: '跳转成功',
          icon: 'success'
        })
      },
      fail: (error) => {
        console.error('❌ 跳转失败:', error)
        wx.showToast({
          title: '跳转失败: ' + (error.errMsg || '未知错误'),
          icon: 'none'
        })
      }
    })
    
  } catch (error) {
    console.error('测试失败:', error)
    wx.showToast({
      title: '测试失败: ' + error.message,
      icon: 'none'
    })
  }
}

// 测试用户资料页面功能
function testProfilePage() {
  console.log('=== 测试用户资料页面功能 ===')
  
  try {
    // 检查页面数据
    const pages = getCurrentPages()
    const currentPage = pages[pages.length - 1]
    
    console.log('当前页面:', currentPage.route)
    console.log('页面数据:', currentPage.data)
    
    // 检查用户信息
    const userInfo = currentPage.data.userInfo
    console.log('用户信息:', userInfo)
    
    if (userInfo) {
      console.log('昵称:', userInfo.nickName || '未设置')
      console.log('手机号:', userInfo.phone || '未绑定')
      console.log('头像:', userInfo.avatarUrl || '无头像')
    }
    
  } catch (error) {
    console.error('测试失败:', error)
  }
}

// 导出测试函数
module.exports = {
  testPageNavigation,
  testProfilePage
}

// 如果直接运行此文件，执行测试
if (typeof wx !== 'undefined') {
  // 在小程序环境中
  console.log('请在用户资料页面运行测试')
} else {
  // 在Node.js环境中
  console.log('请在微信小程序环境中运行此测试')
} 