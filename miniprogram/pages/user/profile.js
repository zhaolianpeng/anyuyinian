// pages/user/profile.js
const app = getApp()
const { api } = require('../../utils/cloud-container-standard')
const { getCurrentUserId, clearUserId, needsUserIdMigration, safeClearUserId } = require('../../utils/user-id-compatibility')

Page({
  data: {
    userInfo: null,
    loading: false,
    isAdmin: false,
    needLogin: false
  },

  onLoad() {
    console.log('用户资料页面加载')
  },

  onShow() {
    console.log('用户资料页面显示')
    this.loadUserInfo()
  },

  // 加载用户信息
  async loadUserInfo() {
    try {
      this.setData({ loading: true })
      
      const userId = getCurrentUserId()
      console.log('开始加载用户信息，当前用户ID:', userId)
      
      if (!userId) {
        console.log('用户未登录，显示登录门')
        this.setData({ needLogin: true, loading: false })
        return
      }

      // 检查是否需要迁移用户ID
      if (needsUserIdMigration()) {
        console.log('检测到旧格式的用户ID，需要重新登录')
        wx.showModal({
          title: '系统升级',
          content: '系统已升级，需要重新登录以获取新的用户ID',
          showCancel: false,
          success: () => {
            clearUserId('系统升级，需要重新登录')
            wx.navigateTo({ url: '/pages/login/login' })
          }
        })
        return
      }

      console.log('开始获取用户信息，userId:', userId)
      const result = await api.userInfo({ userId })
      
      if (result.code === 0) {
        this.setData({ 
          userInfo: result.data,
          loading: false,
          needLogin: false
        })
        console.log('用户信息获取成功:', result.data)
      } else {
        throw new Error(result.errorMsg || '获取用户信息失败')
      }
    } catch (error) {
      console.error('加载用户信息失败:', error)
      this.setData({ loading: false })
      
      // 如果是用户不存在错误，安全清除本地存储并跳转登录
      if (error.message && error.message.includes('record not found')) {
        console.log('用户不存在，安全清除本地存储并跳转登录')
        safeClearUserId('用户不存在，需要重新登录')
        this.setData({ needLogin: true })
        return
      }
      
      wx.showToast({
        title: '加载失败，请重试',
        icon: 'none'
      })
    }
  },

  // 去登录
  onGoLogin() {
    const current = '/pages/user/profile'
    wx.navigateTo({
      url: `/pages/login/login?redirect=${encodeURIComponent(current)}`
    })
  },

  // 完善资料
  onSetupProfile() {
    console.log('点击完善资料按钮')
    
    // 检查用户登录状态
    const userId = getCurrentUserId()
    if (!userId) {
      wx.showToast({
        title: '请先登录',
        icon: 'none'
      })
      return
    }
    
    console.log('准备跳转到设置页面')
    wx.navigateTo({
      url: '/pages/user/setup-profile',
      success: () => {
        console.log('跳转成功')
      },
      fail: (error) => {
        console.error('跳转失败:', error)
        wx.showToast({
          title: '页面跳转失败',
          icon: 'none'
        })
      }
    })
  },

  // 绑定手机号（保留原有功能）
  onBindPhone() {
    // 检查是否支持获取手机号
    if (wx.getPhoneNumber) {
      wx.showModal({
        title: '绑定手机号',
        content: '是否获取微信绑定的手机号？',
        success: (res) => {
          if (res.confirm) {
            this.getPhoneNumber()
          }
        }
      })
    } else {
      wx.showToast({
        title: '当前版本不支持获取手机号',
        icon: 'none'
      })
    }
  },

  // 获取手机号
  getPhoneNumber(e) {
    console.log('获取手机号:', e)
    
    if (e.detail.errMsg !== 'getPhoneNumber:ok') {
      wx.showToast({
        title: '获取手机号失败',
        icon: 'none'
      })
      return
    }

    // 这里应该调用后端接口绑定手机号
    // 由于getPhoneNumber返回的是加密数据，需要后端解密
    console.log('手机号加密数据:', e.detail)
    
    wx.showToast({
      title: '手机号获取成功',
      icon: 'success'
    })
  },

  // 地址管理
  onAddressManage() {
    wx.navigateTo({
      url: '/pages/user/address/list'
    })
  },

  // 就诊人管理
  onPatientManage() {
    wx.navigateTo({
      url: '/pages/user/patient/list'
    })
  },

  // 订单列表
  onOrderList() {
    wx.switchTab({
      url: '/pages/order/list'
    })
  },

  // 推广中心
  onReferralCenter() {
    wx.navigateTo({
      url: '/pages/promoter/home'
    })
  },

  // 客服
  onCustomerService() {
    wx.navigateTo({
      url: '/pages/kefu/chat'
    })
  },

  // 管理员登录
  onAdminLogin() {
    console.log('点击管理员入口按钮')
    wx.navigateTo({
      url: '/pages/admin/login'
    })
  },

  // 管理员中心跳转
  onGoAdminHome() {
    wx.redirectTo({
      url: '/pages/admin/home'
    })
  },

  // 检查管理员状态
  checkAdminStatus() {
    // 仅在已登录时检查管理员状态
    const userId = getCurrentUserId()
    if (!userId) {
      this.setData({ isAdmin: false })
      return
    }
    const adminInfo = wx.getStorageSync('adminInfo')
    const isAdmin = wx.getStorageSync('isAdmin')
    this.setData({ isAdmin: !!(adminInfo && isAdmin) })
  },

  // 退出登录
  onLogout() {
    wx.showModal({
      title: '确认退出',
      content: '确定要退出登录吗？',
      success: (res) => {
        if (res.confirm) {
          // 清除本地存储
          clearUserId()
          wx.removeStorageSync('token')
          wx.removeStorageSync('userInfo')
          
          console.log('用户已退出登录')
          
          wx.showToast({
            title: '已退出登录',
            icon: 'success'
          })
          
          // 跳转到首页
          setTimeout(() => {
            wx.switchTab({
              url: '/pages/index/index'
            })
          }, 1500)
        }
      }
    })
  },

  // 分享给好友
  onShareAppMessage(res) {
    console.log('个人中心分享被触发:', res)
    
    return {
      title: '安语颐年护理陪诊 - 个人中心',
      desc: '专业护理陪诊服务，让您就医更安心',
      path: 'pages/index/index'
    }
  },

  // 分享到朋友圈
  onShareTimeline(res) {
    console.log('个人中心分享到朋友圈被触发:', res)
    
    return {
      title: '安语颐年护理陪诊 - 专业护理陪诊服务',
      query: '',
      imageUrl: ''
    }
  }
}) 