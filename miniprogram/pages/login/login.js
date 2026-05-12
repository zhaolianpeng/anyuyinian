// pages/login/login.js
const app = getApp()
const { api } = require('../../utils/cloud-container-standard')
const { setUserId, getCurrentUserId } = require('../../utils/user-id-compatibility')

Page({
  data: {
    loading: false,
    redirectUrl: '',
    protocolAgreed: false
  },

  onLoad(options) {
    console.log('登录页面加载，参数:', options)
    
    // 检查是否有重定向URL
    if (options.redirect) {
      this.setData({
        redirectUrl: decodeURIComponent(options.redirect)
      })
      console.log('设置重定向URL:', this.data.redirectUrl)
    }
    
    // 检查是否已经登录
    const userId = getCurrentUserId()
    if (userId) {
      console.log('用户已登录，userId:', userId)
      // 如果已登录，直接跳转
      this.redirectAfterLogin()
    }

    // 初始化协议勾选状态（根据本地存储）
    const ua = wx.getStorageSync('agree_user_agreement') || false
    const pp = wx.getStorageSync('agree_privacy_policy') || false
    this.setData({ protocolAgreed: !!(ua && pp) })
  },

  onShow() {
    // 返回本页时，检测两个协议是否都已同意，若都同意则自动勾选
    // 但如果用户已经手动勾选了，则保持用户的选择
    const ua = wx.getStorageSync('agree_user_agreement') || false
    const pp = wx.getStorageSync('agree_privacy_policy') || false
    const bothAgreed = !!(ua && pp)
    
    // 如果两个协议都已同意，或者用户已经勾选，则保持勾选状态
    if (bothAgreed || this.data.protocolAgreed) {
      this.setData({ protocolAgreed: true })
    }
  },

  onLoginTap() {
    console.log('点击登录按钮')
    
    if (this.data.loading) {
      console.log('正在登录中，忽略重复点击')
      return
    }

    if (!this.data.protocolAgreed) {
      wx.showToast({ title: '请先勾选并同意协议', icon: 'none' })
      return
    }

    this.setData({ loading: true })
    console.log('开始登录流程')

    // 获取用户信息
    wx.getUserProfile({
      desc: '用于完善用户资料',
      success: (userInfoRes) => {
        console.log('获取用户信息成功:', userInfoRes)
        
        // 调用微信登录
        wx.login({
          success: (loginRes) => {
            console.log('微信登录成功:', loginRes)
            
            if (!loginRes.code) {
              throw new Error('微信登录失败，未获取到code')
            }

            // 调用后端登录接口
            api.wxLogin({
              code: loginRes.code,
              nickName: userInfoRes.userInfo.nickName,
              avatarUrl: userInfoRes.userInfo.avatarUrl,
              gender: userInfoRes.userInfo.gender,
              country: userInfoRes.userInfo.country,
              province: userInfoRes.userInfo.province,
              city: userInfoRes.userInfo.city,
              language: userInfoRes.userInfo.language
            })
            .then((res) => {
              console.log('后端登录响应:', res)
              
              const responseCode = res.code
              const responseData = res.data || {}
              
              if (responseCode === 0 || !res.errorMsg) {
                // 登录成功，保存用户信息
                const userInfo = {
                  ...userInfoRes.userInfo,
                  ...responseData
                }
                
                // 保存token
                wx.setStorageSync('token', responseData.token || 'mock_token_' + Date.now())
                wx.setStorageSync('userInfo', userInfo)
                
                // 保存userId（使用新的格式）
                if (responseData.userId) {
                  setUserId(responseData.userId)
                } else {
                  console.warn('后端未返回userId，使用模拟ID')
                  setUserId('mock_user_' + Date.now())
                }
                
                console.log('用户信息已保存到本地存储')
                
                wx.showToast({
                  title: '登录成功',
                  icon: 'success'
                })
                
                // 登录成功后跳转
                setTimeout(() => {
                  this.redirectAfterLogin()
                }, 1500)
              } else {
                throw new Error(res.errorMsg || '登录失败')
              }
            })
            .catch((error) => {
              console.error('登录失败:', error)
              console.error('错误详情:', {
                message: error.message,
                errMsg: error.errMsg,
                stack: error.stack
              })
              
              // 针对网络错误的特殊处理
              if (error.message && error.message.includes('网络错误')) {
                wx.showModal({
                  title: '网络错误',
                  content: '请检查网络连接，或稍后重试',
                  showCancel: false
                })
              } else if (error.message && error.message.includes('url not in domain list')) {
                wx.showModal({
                  title: '域名限制',
                  content: '当前环境无法访问后端服务，请使用模拟数据',
                  showCancel: false
                })
              } else {
                wx.showToast({
                  title: error.errMsg || error.message || '登录失败，请重试',
                  icon: 'none'
                })
              }
            })
            .finally(() => {
              this.setData({ loading: false })
            })
          },
          fail: (err) => {
            console.error('wx.login失败:', err)
            wx.showToast({
              title: '登录失败，请重试',
              icon: 'none'
            })
            this.setData({ loading: false })
          }
        })
      },
      fail: (err) => {
        console.error('wx.getUserProfile失败:', err)
        console.error('登录失败:', err)
        
        // 针对getUserProfile的特殊错误处理
        if (err.errMsg && err.errMsg.includes('getUserProfile:fail can only be invoked by user TAP gesture')) {
          wx.showModal({
            title: '获取用户信息失败',
            content: '请点击"获取用户信息"按钮重新授权',
            showCancel: false
          })
        } else if (err.errMsg && err.errMsg.includes('getUserProfile:fail auth deny')) {
          wx.showModal({
            title: '用户拒绝授权',
            content: '需要获取您的用户信息才能正常使用小程序',
            showCancel: false
          })
        } else {
          wx.showToast({
            title: '获取用户信息失败，请重试',
            icon: 'none'
          })
        }
        
        this.setData({ loading: false })
      }
    })
  },

  // 勾选框变化
  onProtocolChange(e) {
    const values = e.detail.value || []
    const checked = values.includes('agree')
    // 用户可以直接勾选，无需强制阅读协议
    this.setData({ protocolAgreed: checked })
  },

  // 登录成功后的跳转逻辑
  redirectAfterLogin() {
    // 如果有重定向URL，跳转到目标页面
    if (this.data.redirectUrl) {
      console.log('登录成功，跳转到目标页面:', this.data.redirectUrl)
      wx.navigateTo({
        url: this.data.redirectUrl,
        fail: () => {
          // 如果跳转失败，返回上一页
          wx.navigateBack({
            fail: () => {
              wx.switchTab({
                url: '/pages/index/index'
              })
            }
          })
        }
      })
    } else {
      // 没有重定向URL，返回上一页
      wx.navigateBack({
        fail: () => {
          wx.switchTab({
            url: '/pages/index/index'
          })
        }
      })
    }
  },

  

  onUserAgreement() {
    wx.navigateTo({
      url: '/pages/agreement/user-agreement'
    })
  },

  onPrivacyPolicy() {
    wx.navigateTo({
      url: '/pages/agreement/privacy-policy'
    })
  }
})