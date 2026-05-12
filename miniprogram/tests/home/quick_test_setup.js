// 快速测试设置页面功能
// 在用户资料页面控制台运行此脚本

function quickTestSetup() {
  console.log('=== 快速测试设置页面功能 ===')
  
  try {
    // 1. 检查用户登录状态
    const userId = wx.getStorageSync('userId')
    console.log('用户ID:', userId || '未登录')
    
    if (!userId) {
      console.log('❌ 请先登录')
      return
    }
    
    // 2. 测试页面跳转
    console.log('测试页面跳转...')
    wx.navigateTo({
      url: '/pages/user/setup-profile',
      success: () => {
        console.log('✅ 页面跳转成功')
        console.log('请在设置页面中测试获取微信信息功能')
      },
      fail: (error) => {
        console.log('❌ 页面跳转失败:', error)
      }
    })
    
  } catch (error) {
    console.error('测试失败:', error)
  }
}

// 测试微信API兼容性
function testWxAPI() {
  console.log('=== 测试微信API兼容性 ===')
  
  const apis = {
    getUserProfile: typeof wx.getUserProfile === 'function',
    getPhoneNumber: typeof wx.getPhoneNumber === 'function'
  }
  
  console.log('API状态:')
  console.log('- getUserProfile:', apis.getUserProfile ? '✅ 可用' : '❌ 不可用')
  console.log('- getPhoneNumber:', apis.getPhoneNumber ? '✅ 可用' : '❌ 不可用')
  
  if (!apis.getUserProfile || !apis.getPhoneNumber) {
    console.log('⚠️ 在开发环境中将使用模拟数据')
  }
  
  return apis
}

// 测试后端API
async function testBackendAPI() {
  console.log('=== 测试后端API ===')
  
  try {
    const { api } = require('../utils/cloud-container-standard')
    const userId = wx.getStorageSync('userId')
    
    if (!userId) {
      console.log('❌ 用户未登录')
      return
    }
    
    // 测试用户信息API
    console.log('测试用户信息API...')
    const userInfo = await api.userInfo({ userId })
    console.log('用户信息API结果:', userInfo)
    
    // 测试更新用户信息API
    console.log('测试更新用户信息API...')
    const updateResult = await api.updateUserInfo({
      userId: userId,
      nickName: '测试昵称_' + Date.now()
    })
    console.log('更新用户信息API结果:', updateResult)
    
  } catch (error) {
    console.error('后端API测试失败:', error)
  }
}

// 运行所有测试
function runAllTests() {
  console.log('🚀 开始运行所有测试...')
  
  // 测试微信API
  testWxAPI()
  
  // 测试后端API
  setTimeout(() => {
    testBackendAPI()
  }, 500)
  
  // 测试页面跳转
  setTimeout(() => {
    quickTestSetup()
  }, 1000)
}

// 导出函数
module.exports = {
  quickTestSetup,
  testWxAPI,
  testBackendAPI,
  runAllTests
}

// 如果直接运行，执行所有测试
if (typeof wx !== 'undefined') {
  console.log('快速测试脚本已加载')
  console.log('运行 runAllTests() 执行所有测试')
  console.log('运行 quickTestSetup() 测试页面跳转')
  console.log('运行 testWxAPI() 测试微信API')
  console.log('运行 testBackendAPI() 测试后端API')
} else {
  console.log('请在微信小程序环境中运行此测试')
} 