// pages/test/setup-profile-debug.js
const { api } = require('../../utils/cloud-container-standard')

Page({
  data: {
    debugInfo: '',
    testResults: []
  },

  onLoad() {
    this.addDebugInfo('页面加载完成')
    this.testBasicFunctions()
  },

  // 添加调试信息
  addDebugInfo(message) {
    const timestamp = new Date().toLocaleTimeString()
    const debugInfo = `[${timestamp}] ${message}\n`
    this.setData({
      debugInfo: this.data.debugInfo + debugInfo
    })
    console.log(message)
  },

  // 测试基本功能
  async testBasicFunctions() {
    this.addDebugInfo('开始测试基本功能...')
    
    try {
      // 1. 检查用户登录状态
      const userId = wx.getStorageSync('userId')
      this.addDebugInfo(`用户ID: ${userId || '未登录'}`)
      
      if (!userId) {
        this.addDebugInfo('❌ 用户未登录')
        return
      }

      // 2. 测试API调用
      this.addDebugInfo('测试用户信息API...')
      const userInfoResult = await api.userInfo({ userId })
      this.addDebugInfo(`用户信息API结果: ${JSON.stringify(userInfoResult)}`)
      
      // 3. 测试页面跳转
      this.addDebugInfo('测试页面跳转...')
      wx.navigateTo({
        url: '/pages/user/setup-profile',
        success: () => {
          this.addDebugInfo('✅ 页面跳转成功')
        },
        fail: (error) => {
          this.addDebugInfo(`❌ 页面跳转失败: ${JSON.stringify(error)}`)
        }
      })

    } catch (error) {
      this.addDebugInfo(`❌ 测试失败: ${error.message}`)
    }
  },

  // 测试微信API
  testWxAPI() {
    this.addDebugInfo('测试微信API...')
    
    // 测试getUserProfile
    if (wx.getUserProfile) {
      this.addDebugInfo('✅ wx.getUserProfile 可用')
    } else {
      this.addDebugInfo('❌ wx.getUserProfile 不可用')
    }
    
    // 测试getPhoneNumber
    if (wx.getPhoneNumber) {
      this.addDebugInfo('✅ wx.getPhoneNumber 可用')
    } else {
      this.addDebugInfo('❌ wx.getPhoneNumber 不可用')
    }
  },

  // 测试后端API
  async testBackendAPI() {
    this.addDebugInfo('测试后端API...')
    
    try {
      const userId = wx.getStorageSync('userId')
      if (!userId) {
        this.addDebugInfo('❌ 用户未登录，无法测试后端API')
        return
      }

      // 测试更新用户信息API
      this.addDebugInfo('测试更新用户信息API...')
      const updateResult = await api.updateUserInfo({
        userId: userId,
        nickName: '测试昵称_' + Date.now()
      })
      this.addDebugInfo(`更新用户信息结果: ${JSON.stringify(updateResult)}`)

    } catch (error) {
      this.addDebugInfo(`❌ 后端API测试失败: ${error.message}`)
    }
  },

  // 清除调试信息
  clearDebugInfo() {
    this.setData({
      debugInfo: '',
      testResults: []
    })
  },

  // 复制调试信息
  copyDebugInfo() {
    wx.setClipboardData({
      data: this.data.debugInfo,
      success: () => {
        wx.showToast({
          title: '调试信息已复制',
          icon: 'success'
        })
      }
    })
  }
}) 