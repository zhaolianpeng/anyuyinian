// 测试订单导航功能
function testOrderNavigation() {
  console.log('=== 测试订单导航功能 ===')
  
  try {
    // 1. 检查用户登录状态
    const userId = wx.getStorageSync('userId')
    console.log('当前用户ID:', userId || '未登录')
    
    if (!userId) {
      console.log('❌ 用户未登录')
      wx.showToast({
        title: '请先登录',
        icon: 'none'
      })
      return
    }
    
    // 2. 测试跳转到订单列表页面
    console.log('准备跳转到订单列表页面...')
    wx.navigateTo({
      url: '/pages/order/list',
      success: () => {
        console.log('✅ 跳转到订单列表页面成功')
      },
      fail: (error) => {
        console.error('❌ 跳转到订单列表页面失败:', error)
        wx.showToast({
          title: '页面跳转失败',
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

// 测试订单列表API
async function testOrderListAPI() {
  console.log('=== 测试订单列表API ===')
  
  try {
    const { api } = require('../utils/cloud-container-standard')
    const userId = wx.getStorageSync('userId')
    
    if (!userId) {
      console.log('❌ 用户未登录，无法测试API')
      return
    }
    
    console.log('测试订单列表API...')
    const result = await api.orderList({
      userId: userId,
      page: 1,
      pageSize: 10,
      status: ''
    })
    
    console.log('订单列表API结果:', result)
    
    if (result.code === 0) {
      console.log('✅ 订单列表API调用成功')
      const orders = result.data?.list || []
      console.log('订单数量:', orders.length)
      
      if (orders.length > 0) {
        console.log('第一个订单信息:', {
          id: orders[0].id,
          orderNo: orders[0].orderNo,
          status: orders[0].status,
          totalAmount: orders[0].totalAmount
        })
      }
    } else {
      console.log('❌ 订单列表API调用失败:', result.message)
    }
    
  } catch (error) {
    console.error('订单列表API测试失败:', error)
  }
}

// 测试tabBar订单页面
function testTabBarOrder() {
  console.log('=== 测试TabBar订单页面 ===')
  
  try {
    // 检查当前页面
    const pages = getCurrentPages()
    const currentPage = pages[pages.length - 1]
    console.log('当前页面:', currentPage.route)
    
    // 如果当前在用户资料页面，测试跳转到订单页面
    if (currentPage.route === 'pages/user/profile') {
      console.log('当前在用户资料页面，测试跳转到订单列表')
      testOrderNavigation()
    } else {
      console.log('当前不在用户资料页面，请先进入用户资料页面')
    }
    
  } catch (error) {
    console.error('TabBar订单页面测试失败:', error)
  }
}

// 检查页面配置
function checkPageConfig() {
  console.log('=== 检查页面配置 ===')
  
  try {
    // 检查订单列表页面是否存在
    const pages = getCurrentPages()
    console.log('当前页面栈:', pages.map(p => p.route))
    
    // 检查是否可以访问订单列表页面
    wx.navigateTo({
      url: '/pages/order/list',
      success: () => {
        console.log('✅ 订单列表页面可以正常访问')
        // 立即返回
        setTimeout(() => {
          wx.navigateBack()
        }, 100)
      },
      fail: (error) => {
        console.error('❌ 订单列表页面访问失败:', error)
      }
    })
    
  } catch (error) {
    console.error('页面配置检查失败:', error)
  }
}

// 运行所有测试
function runOrderTests() {
  console.log('🚀 开始运行订单功能测试...')
  
  // 检查页面配置
  checkPageConfig()
  
  // 测试API
  setTimeout(() => {
    testOrderListAPI()
  }, 500)
  
  // 测试导航
  setTimeout(() => {
    testOrderNavigation()
  }, 1000)
  
  // 测试TabBar
  setTimeout(() => {
    testTabBarOrder()
  }, 1500)
}

// 导出测试函数
module.exports = {
  testOrderNavigation,
  testOrderListAPI,
  testTabBarOrder,
  checkPageConfig,
  runOrderTests
}

// 如果直接运行此文件，执行测试
if (typeof wx !== 'undefined') {
  console.log('订单导航测试脚本已加载')
  console.log('运行 runOrderTests() 执行所有测试')
  console.log('运行 testOrderNavigation() 测试页面跳转')
  console.log('运行 testOrderListAPI() 测试API')
  console.log('运行 checkPageConfig() 检查页面配置')
} else {
  console.log('请在微信小程序环境中运行此测试')
} 