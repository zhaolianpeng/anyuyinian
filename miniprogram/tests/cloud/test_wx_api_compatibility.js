// 测试微信API兼容性
function testWxAPICompatibility() {
  console.log('=== 测试微信API兼容性 ===')
  
  const results = {
    getUserProfile: false,
    getPhoneNumber: false,
    environment: 'unknown'
  }
  
  // 检测环境
  if (typeof wx !== 'undefined') {
    if (wx.getSystemInfoSync) {
      const systemInfo = wx.getSystemInfoSync()
      results.environment = systemInfo.platform
      console.log('当前环境:', systemInfo.platform)
    }
  }
  
  // 测试 wx.getUserProfile
  if (typeof wx.getUserProfile === 'function') {
    results.getUserProfile = true
    console.log('✅ wx.getUserProfile 可用')
  } else {
    console.log('❌ wx.getUserProfile 不可用')
  }
  
  // 测试 wx.getPhoneNumber
  if (typeof wx.getPhoneNumber === 'function') {
    results.getPhoneNumber = true
    console.log('✅ wx.getPhoneNumber 可用')
  } else {
    console.log('❌ wx.getPhoneNumber 不可用')
  }
  
  // 输出总结
  console.log('=== API兼容性总结 ===')
  console.log('环境:', results.environment)
  console.log('getUserProfile:', results.getUserProfile ? '✅ 可用' : '❌ 不可用')
  console.log('getPhoneNumber:', results.getPhoneNumber ? '✅ 可用' : '❌ 不可用')
  
  // 提供建议
  if (!results.getUserProfile || !results.getPhoneNumber) {
    console.log('⚠️ 建议：')
    if (!results.getUserProfile) {
      console.log('- 在真机上测试 getUserProfile 功能')
      console.log('- 检查小程序基础库版本是否支持')
    }
    if (!results.getPhoneNumber) {
      console.log('- 在真机上测试 getPhoneNumber 功能')
      console.log('- 检查是否配置了手机号获取权限')
    }
  }
  
  return results
}

// 测试设置页面功能
function testSetupProfileFunction() {
  console.log('=== 测试设置页面功能 ===')
  
  try {
    // 模拟获取微信用户信息
    const mockGetWxProfile = () => {
      return new Promise((resolve, reject) => {
        if (typeof wx.getUserProfile !== 'function') {
          console.log('使用模拟数据')
          resolve({
            nickName: '微信用户',
            avatarUrl: 'https://thirdwx.qlogo.cn/mmopen/vi_32/POgEwh4mIHO4nibH0KlMECNjjGxQUq24ZEaGT4poC6icRiccVGKSyXwibcPq4BWmiaIGuG1icwxaQX6grC9VemZoJ8rg/132',
            gender: 0
          })
        } else {
          wx.getUserProfile({
            desc: '用于完善用户资料',
            success: (res) => resolve(res.userInfo),
            fail: (err) => reject(err)
          })
        }
      })
    }
    
    // 模拟获取手机号
    const mockGetWxPhone = () => {
      return new Promise((resolve) => {
        if (typeof wx.getPhoneNumber !== 'function') {
          console.log('使用模拟手机号')
          resolve('13800138000')
        } else {
          wx.getPhoneNumber({
            success: (res) => resolve('13800138000'),
            fail: (err) => resolve('')
          })
        }
      })
    }
    
    // 测试完整流程
    async function testFullFlow() {
      console.log('1. 获取微信用户信息...')
      const userInfo = await mockGetWxProfile()
      console.log('用户信息:', userInfo)
      
      console.log('2. 获取手机号...')
      const phone = await mockGetWxPhone()
      console.log('手机号:', phone)
      
      console.log('3. 模拟更新用户信息...')
      const userId = wx.getStorageSync('userId')
      if (userId) {
        console.log('用户ID:', userId)
        console.log('准备更新昵称:', userInfo.nickName)
        console.log('准备绑定手机号:', phone)
      }
      
      console.log('✅ 测试流程完成')
    }
    
    testFullFlow()
    
  } catch (error) {
    console.error('测试失败:', error)
  }
}

// 导出测试函数
module.exports = {
  testWxAPICompatibility,
  testSetupProfileFunction
}

// 如果直接运行此文件，执行测试
if (typeof wx !== 'undefined') {
  // 在小程序环境中
  testWxAPICompatibility()
  setTimeout(() => {
    testSetupProfileFunction()
  }, 1000)
} else {
  // 在Node.js环境中
  console.log('请在微信小程序环境中运行此测试')
} 