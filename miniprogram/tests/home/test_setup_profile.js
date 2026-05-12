// 测试设置昵称和绑定手机号功能
const { api } = require('../utils/cloud-container-standard')

async function testSetupProfile() {
  console.log('=== 测试设置昵称和绑定手机号功能 ===')
  
  try {
    // 1. 测试获取用户信息
    console.log('1. 测试获取用户信息...')
    const userId = wx.getStorageSync('userId')
    if (!userId) {
      console.log('❌ 用户未登录')
      return
    }

    const userInfoResult = await api.userInfo({ userId })
    console.log('用户信息:', userInfoResult)
    
    if (userInfoResult.code === 0) {
      const userInfo = userInfoResult.data
      console.log('当前用户信息:', {
        nickName: userInfo.nickName || '未设置',
        phone: userInfo.phone || '未绑定',
        avatarUrl: userInfo.avatarUrl || '无头像'
      })
    }

    // 2. 测试更新用户昵称
    console.log('2. 测试更新用户昵称...')
    const updateResult = await api.updateUserInfo({
      userId: userId,
      nickName: '测试昵称_' + Date.now()
    })
    
    if (updateResult.code === 0) {
      console.log('✅ 更新昵称成功')
    } else {
      console.log('❌ 更新昵称失败:', updateResult.message)
    }

    // 3. 测试绑定手机号
    console.log('3. 测试绑定手机号...')
    const bindPhoneResult = await api.bindPhone({
      userId: userId,
      phone: '13800138000',
      code: '123456'
    })
    
    if (bindPhoneResult.code === 0) {
      console.log('✅ 绑定手机号成功')
    } else {
      console.log('❌ 绑定手机号失败:', bindPhoneResult.message)
    }

    // 4. 验证更新后的用户信息
    console.log('4. 验证更新后的用户信息...')
    const updatedUserInfo = await api.userInfo({ userId })
    if (updatedUserInfo.code === 0) {
      console.log('更新后的用户信息:', updatedUserInfo.data)
    }

  } catch (error) {
    console.error('测试失败:', error)
  }
  
  console.log('=== 测试完成 ===')
}

// 导出测试函数
module.exports = {
  testSetupProfile
}

// 如果直接运行此文件，执行测试
if (typeof wx !== 'undefined') {
  // 在小程序环境中
  testSetupProfile()
} else {
  // 在Node.js环境中
  console.log('请在微信小程序环境中运行此测试')
} 