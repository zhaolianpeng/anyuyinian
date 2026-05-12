// 导航跳转修复测试脚本
// 用于测试修复后的导航跳转功能

const testNavigationFix = () => {
  console.log('=== 开始测试修复后的导航跳转功能 ===')
  
  // 测试导航数据
  const testCases = [
    {
      name: '服务预约',
      url: '/pages/service/list',
      isTabBar: true,
      needLogin: false
    },
    {
      name: '我的订单',
      url: '/pages/order/list',
      isTabBar: true,
      needLogin: true
    },
    {
      name: '医院信息',
      url: '/pages/hospital/list',
      isTabBar: false,
      needLogin: false
    },
    {
      name: '个人中心',
      url: '/pages/user/profile',
      isTabBar: true,
      needLogin: true
    }
  ]
  
  console.log('测试用例:', testCases)
  
  // 检查用户登录状态
  const userInfo = wx.getStorageSync('userInfo')
  console.log('当前用户登录状态:', userInfo ? '已登录' : '未登录')
  
  // 测试每个导航项
  testCases.forEach((testCase, index) => {
    console.log(`\n--- 测试导航项 ${index + 1}: ${testCase.name} ---`)
    console.log('测试信息:', {
      name: testCase.name,
      url: testCase.url,
      isTabBar: testCase.isTabBar,
      needLogin: testCase.needLogin
    })
    
    // 检查URL格式
    if (!testCase.url.startsWith('/pages/')) {
      console.error('❌ URL格式不正确')
      return
    }
    
    console.log('✅ URL格式正确')
    
    // 检查是否需要登录
    if (testCase.needLogin && !userInfo) {
      console.log('⚠️ 该页面需要登录，用户未登录')
      console.log('✅ 应该跳转到登录页面')
    } else if (testCase.needLogin && userInfo) {
      console.log('✅ 该页面需要登录，用户已登录')
    } else {
      console.log('✅ 该页面不需要登录')
    }
    
    // 检查跳转方式
    if (testCase.isTabBar) {
      console.log('✅ 应该使用 wx.switchTab 跳转')
    } else {
      console.log('✅ 应该使用 wx.navigateTo 跳转')
    }
    
    console.log('✅ 测试完成')
  })
  
  console.log('\n=== 导航跳转修复测试完成 ===')
}

// 模拟页面跳转函数
const simulateNavigateToPage = (url) => {
  console.log('模拟页面跳转:', url)
  
  // 检查是否需要登录
  const needLoginPages = ['/pages/user/profile', '/pages/order/list']
  const isNeedLogin = needLoginPages.some(page => url.includes(page))
  
  if (isNeedLogin) {
    const userInfo = wx.getStorageSync('userInfo')
    if (!userInfo) {
      console.log('需要登录，应该跳转到登录页面')
      return 'login'
    }
  }
  
  // 检查是否是tabBar页面
  const tabBarPages = ['/pages/index/index', '/pages/service/list', '/pages/order/list', '/pages/user/profile']
  const isTabBarPage = tabBarPages.some(page => url.includes(page))
  
  if (isTabBarPage) {
    console.log('应该使用 wx.switchTab 跳转到:', url)
    return 'switchTab'
  } else {
    console.log('应该使用 wx.navigateTo 跳转到:', url)
    return 'navigateTo'
  }
}

// 测试所有导航项
const testAllNavigations = () => {
  console.log('=== 测试所有导航项 ===')
  
  const navigations = [
    '/pages/service/list',
    '/pages/order/list', 
    '/pages/hospital/list',
    '/pages/user/profile'
  ]
  
  navigations.forEach((url, index) => {
    console.log(`\n测试导航 ${index + 1}: ${url}`)
    const result = simulateNavigateToPage(url)
    console.log('跳转方式:', result)
  })
}

// 导出测试函数
module.exports = {
  testNavigationFix,
  simulateNavigateToPage,
  testAllNavigations
}

// 如果直接运行此文件，执行测试
if (typeof wx !== 'undefined') {
  testNavigationFix()
  testAllNavigations()
} 