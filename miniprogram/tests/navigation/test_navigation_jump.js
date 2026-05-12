// 导航跳转测试脚本
// 用于测试首页导航项的跳转功能

const testNavigationJump = () => {
  console.log('=== 开始测试导航跳转功能 ===')
  
  // 模拟导航数据
  const testNavigations = [
    {
      id: 1,
      name: '服务预约',
      icon: 'https://i.postimg.cc/p5W10Vw7/fuwuyuyue-logo.png',
      linkUrl: '/pages/service/list',
      sort: 1
    },
    {
      id: 2,
      name: '我的订单',
      icon: 'https://i.postimg.cc/phY6Y56z/wodedingdan-logo.png',
      linkUrl: '/pages/order/list',
      sort: 2
    },
    {
      id: 3,
      name: '医院信息',
      icon: 'https://i.postimg.cc/BLGzKBMJ/yiyuanxinxi-logo.png',
      linkUrl: '/pages/hospital/list',
      sort: 3
    },
    {
      id: 4,
      name: '个人中心',
      icon: 'https://i.postimg.cc/XZ5hqzK2/gerenzhongxin-logo.png',
      linkUrl: '/pages/user/profile',
      sort: 4
    }
  ]
  
  console.log('测试导航数据:', testNavigations)
  
  // 测试每个导航项的跳转
  testNavigations.forEach((nav, index) => {
    console.log(`\n--- 测试导航项 ${index + 1}: ${nav.name} ---`)
    console.log('导航信息:', {
      name: nav.name,
      linkUrl: nav.linkUrl,
      icon: nav.icon
    })
    
    // 检查linkUrl格式
    if (!nav.linkUrl) {
      console.error('❌ linkUrl为空')
      return
    }
    
    if (!nav.linkUrl.startsWith('/pages/')) {
      console.error('❌ linkUrl格式不正确，应该以/pages/开头')
      return
    }
    
    console.log('✅ linkUrl格式正确')
    
    // 检查是否需要登录
    const needLoginPages = ['/pages/user/profile', '/pages/order/list']
    const isNeedLogin = needLoginPages.some(page => nav.linkUrl.includes(page))
    
    if (isNeedLogin) {
      console.log('⚠️ 该页面需要登录')
      
      // 检查用户登录状态
      const userInfo = wx.getStorageSync('userInfo')
      if (!userInfo) {
        console.log('❌ 用户未登录，应该跳转到登录页面')
      } else {
        console.log('✅ 用户已登录，可以直接跳转')
      }
    } else {
      console.log('✅ 该页面不需要登录')
    }
    
    // 检查目标页面是否存在
    const pagePath = nav.linkUrl.replace('/pages/', 'pages/')
    console.log('目标页面路径:', pagePath)
    
    // 这里可以添加页面存在性检查
    console.log('✅ 页面跳转测试完成')
  })
  
  console.log('\n=== 导航跳转测试完成 ===')
}

// 测试页面跳转逻辑
const testNavigateToPage = (url) => {
  console.log('测试页面跳转:', url)
  
  if (url.startsWith('/pages/')) {
    // 检查是否需要登录
    const needLoginPages = ['/pages/user/profile', '/pages/order/list']
    const isNeedLogin = needLoginPages.some(page => url.includes(page))
    
    if (isNeedLogin) {
      const userInfo = wx.getStorageSync('userInfo')
      if (!userInfo) {
        console.log('需要登录，应该跳转到登录页面')
        return false
      }
    }
    
    console.log('可以直接跳转到:', url)
    return true
  } else if (url.startsWith('http')) {
    console.log('外部链接，应该通过webview打开')
    return true
  } else {
    console.log('无效的页面URL')
    return false
  }
}

// 导出测试函数
module.exports = {
  testNavigationJump,
  testNavigateToPage
}

// 如果直接运行此文件，执行测试
if (typeof wx !== 'undefined') {
  testNavigationJump()
} 