// 首页导航功能测试脚本

/**
 * 测试导航数据
 */
function testNavigationData() {
  console.log('=== 测试导航数据 ===')
  
  const expectedNavigations = [
    { name: '服务预约', linkUrl: '/pages/service/list' },
    { name: '我的订单', linkUrl: '/pages/order/list' },
    { name: '个人中心', linkUrl: '/pages/user/profile' },
    { name: '客服咨询', linkUrl: '/pages/kefu/chat' },
    { name: '推广中心', linkUrl: '/pages/promoter/home' },
    { name: '健康资讯', linkUrl: '/pages/webview/webview?url=https://example.com/news' },
    { name: '附近医院', linkUrl: '/pages/hospital/list' }
  ]
  
  console.log('期望的导航数据:')
  expectedNavigations.forEach((nav, index) => {
    console.log(`${index + 1}. ${nav.name} -> ${nav.linkUrl}`)
  })
  
  return expectedNavigations
}

/**
 * 测试页面跳转逻辑
 */
function testPageNavigation() {
  console.log('\n=== 测试页面跳转逻辑 ===')
  
  const testUrls = [
    '/pages/service/list',
    '/pages/order/list',
    '/pages/user/profile',
    '/pages/kefu/chat',
    '/pages/promoter/home',
    '/pages/webview/webview?url=https://example.com/news',
    '/pages/hospital/list',
    'https://example.com/external-link',
    'invalid-url'
  ]
  
  console.log('测试URL列表:')
  testUrls.forEach((url, index) => {
    console.log(`${index + 1}. ${url}`)
  })
  
  return testUrls
}

/**
 * 测试登录检查逻辑
 */
function testLoginCheck() {
  console.log('\n=== 测试登录检查逻辑 ===')
  
  const needLoginPages = ['/pages/user/profile', '/pages/order/list']
  const testPages = [
    '/pages/service/list',
    '/pages/order/list',
    '/pages/user/profile',
    '/pages/kefu/chat',
    '/pages/promoter/home'
  ]
  
  console.log('需要登录的页面:')
  needLoginPages.forEach(page => {
    console.log(`- ${page}`)
  })
  
  console.log('\n测试页面登录检查:')
  testPages.forEach(page => {
    const needLogin = needLoginPages.some(loginPage => page.includes(loginPage))
    console.log(`${page} -> ${needLogin ? '需要登录' : '无需登录'}`)
  })
}

/**
 * 模拟页面跳转函数
 */
function simulateNavigateToPage(url) {
  console.log(`\n模拟跳转到: ${url}`)
  
  if (url.startsWith('/pages/')) {
    // 检查是否需要登录
    const needLoginPages = ['/pages/user/profile', '/pages/order/list']
    const isNeedLogin = needLoginPages.some(page => url.includes(page))
    
    if (isNeedLogin) {
      const userInfo = wx.getStorageSync('userInfo')
      if (!userInfo) {
        console.log('需要登录，跳转到登录页面')
        return '/pages/login/login'
      }
    }
    
    console.log('直接跳转到页面')
    return url
  } else if (url.startsWith('http')) {
    console.log('跳转到WebView页面')
    return `/pages/webview/webview?url=${encodeURIComponent(url)}`
  } else {
    console.log('无效的页面URL')
    return null
  }
}

/**
 * 运行所有测试
 */
function runAllTests() {
  console.log('🧭 开始首页导航功能测试\n')
  
  const expectedNavigations = testNavigationData()
  const testUrls = testPageNavigation()
  testLoginCheck()
  
  console.log('\n=== 模拟页面跳转测试 ===')
  testUrls.forEach((url, index) => {
    const result = simulateNavigateToPage(url)
    console.log(`测试 ${index + 1}: ${url} -> ${result || '跳转失败'}`)
  })
  
  console.log('\n🧭 首页导航功能测试完成')
}

// 如果直接运行此脚本
if (typeof module !== 'undefined' && module.exports) {
  module.exports = {
    testNavigationData,
    testPageNavigation,
    testLoginCheck,
    simulateNavigateToPage,
    runAllTests
  }
}

// 在微信小程序环境中运行测试
if (typeof wx !== 'undefined') {
  // 延迟执行，确保页面加载完成
  setTimeout(() => {
    runAllTests()
  }, 1000)
} 