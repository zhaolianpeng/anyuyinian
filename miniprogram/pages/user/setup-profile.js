// pages/user/setup-profile.js
const { api } = require('../../utils/cloud-container-standard')
const { getCurrentUserId, needsUserIdMigration, clearUserId, safeClearUserId } = require('../../utils/user-id-compatibility')

Page({
  data: {
    userInfo: null,
    wxUserInfo: null,
    phone: '',
    loading: false,
    step: 1, // 1: 获取微信信息, 2: 确认设置
    showPhoneModal: false
  },

  onLoad() {
    this.loadUserInfo()
  },

  // 加载用户信息
  async loadUserInfo() {
    try {
      const userId = getCurrentUserId()
      if (!userId) {
        wx.navigateTo({ url: '/pages/login/login' })
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

      const result = await api.userInfo({ userId })
      if (result.code === 0) {
        this.setData({ userInfo: result.data })
      }
    } catch (error) {
      console.error('加载用户信息失败:', error)
      // 使用安全清除方法
      safeClearUserId('加载用户信息失败')
    }
  },

  // 获取微信用户信息
  async getWxUserInfo() {
    try {
      this.setData({ loading: true })

      // 获取微信用户信息
      const wxUserInfo = await this.getWxProfile()
      if (!wxUserInfo) {
        throw new Error('获取微信信息失败')
      }

      // 手机号需要通过按钮点击获取，这里先设置为空
      const phone = ''
      
      this.setData({
        wxUserInfo,
        phone,
        step: 2,
        loading: false
      })

      console.log('获取微信信息成功:', {
        nickName: wxUserInfo.nickName,
        avatarUrl: wxUserInfo.avatarUrl,
        phone: phone || '需要点击按钮获取'
      })

      wx.showToast({
        title: '获取用户信息成功',
        icon: 'success'
      })

    } catch (error) {
      console.error('获取微信信息失败:', error)
      this.setData({ loading: false })
      
      wx.showToast({
        title: error.message || '获取微信信息失败',
        icon: 'none'
      })
    }
  },

  // 获取微信用户资料
  getWxProfile() {
    return new Promise((resolve, reject) => {
      // 检查API是否可用
      if (typeof wx.getUserProfile !== 'function') {
        console.log('wx.getUserProfile API 不可用，尝试使用wx.getUserInfo')
        
        // 尝试使用wx.getUserInfo作为备选方案
        if (typeof wx.getUserInfo === 'function') {
          wx.getUserInfo({
            success: (res) => {
              console.log('通过getUserInfo获取用户信息成功:', res)
              resolve(res.userInfo)
            },
            fail: (err) => {
              console.error('getUserInfo也失败:', err)
              // 如果都失败，提供默认数据
              const defaultUserInfo = {
                nickName: '微信用户',
                avatarUrl: 'https://thirdwx.qlogo.cn/mmopen/vi_32/POgEwh4mIHO4nibH0KlMECNjjGxQUq24ZEaGT4poC6icRiccVGKSyXwibcPq4BWmiaIGuG1icwxaQX6grC9VemZoJ8rg/132',
                gender: 0,
                country: '',
                province: '',
                city: '',
                language: ''
              }
              resolve(defaultUserInfo)
            }
          })
        } else {
          // 如果API都不可用，提供默认数据
          const defaultUserInfo = {
            nickName: '微信用户',
            avatarUrl: 'https://thirdwx.qlogo.cn/mmopen/vi_32/POgEwh4mIHO4nibH0KlMECNjjGxQUq24ZEaGT4poC6icRiccVGKSyXwibcPq4BWmiaIGuG1icwxaQX6grC9VemZoJ8rg/132',
            gender: 0,
            country: '',
            province: '',
            city: '',
            language: ''
          }
          resolve(defaultUserInfo)
        }
        return
      }

      wx.getUserProfile({
        desc: '用于完善用户资料',
        success: (res) => {
          console.log('获取微信用户信息成功:', res)
          resolve(res.userInfo)
        },
        fail: (err) => {
          console.error('获取微信用户信息失败:', err)
          reject(new Error('用户拒绝授权获取微信信息'))
        }
      })
    })
  },

  // 获取微信手机号
  getWxPhone() {
    return new Promise((resolve, reject) => {
      // 手机号获取需要通过button的open-type="getPhoneNumber"触发
      // 这里返回一个提示，让用户点击按钮获取
      console.log('手机号需要通过按钮点击获取')
      resolve('')
    })
  },

  // 处理手机号获取按钮点击
  onGetPhoneNumber(e) {
    console.log('手机号获取按钮点击:', e)
    
    if (e.detail.errMsg === 'getPhoneNumber:ok') {
      // 获取到加密的手机号数据，需要发送到后端解密
      const { encryptedData, iv } = e.detail
      
      if (encryptedData && iv) {
        // 调用后端接口解密手机号
        this.decryptPhoneNumber(encryptedData, iv)
      } else {
        console.log('未获取到手机号加密数据')
        wx.showToast({
          title: '获取手机号失败',
          icon: 'none'
        })
      }
    } else {
      console.log('用户拒绝授权获取手机号')
      wx.showToast({
        title: '需要授权才能获取手机号',
        icon: 'none'
      })
    }
  },

  // 解密手机号
  async decryptPhoneNumber(encryptedData, iv) {
    try {
      const userId = getCurrentUserId()
      if (!userId) {
        throw new Error('用户未登录')
      }

      // 调用后端接口解密手机号
      const result = await api.decryptPhoneNumber({
        userId: userId,
        encryptedData: encryptedData,
        iv: iv
      })

      if (result.code === 0 && result.data && result.data.phoneNumber) {
        const phoneNumber = result.data.phoneNumber
        this.setData({
          phone: phoneNumber
        })
        
        console.log('手机号解密成功:', phoneNumber)
        wx.showToast({
          title: '获取手机号成功',
          icon: 'success'
        })
      } else {
        throw new Error(result.message || '手机号解密失败')
      }
    } catch (error) {
      console.error('解密手机号失败:', error)
      wx.showToast({
        title: '获取手机号失败',
        icon: 'none'
      })
    }
  },

  // 确认设置昵称和手机号
  async confirmSetup() {
    try {
      this.setData({ loading: true })

      const userId = getCurrentUserId()
      if (!userId) {
        throw new Error('用户未登录')
      }

      // 检查是否需要迁移用户ID
      if (needsUserIdMigration()) {
        console.log('检测到旧格式的用户ID，需要重新登录')
        wx.showModal({
          title: '系统升级',
          content: '系统已升级，需要重新登录以获取新的用户ID',
          showCancel: false,
          success: () => {
            clearUserId()
            wx.navigateTo({ url: '/pages/login/login' })
          }
        })
        return
      }

      const { wxUserInfo, phone } = this.data

      // 更新用户昵称
      if (wxUserInfo && wxUserInfo.nickName) {
        await this.updateNickname(userId, wxUserInfo.nickName)
      }

      // 绑定手机号
      if (phone) {
        await this.bindPhone(userId, phone)
      }

      wx.showToast({
        title: '设置成功',
        icon: 'success'
      })

      // 返回上一页
      setTimeout(() => {
        wx.navigateBack()
      }, 1500)

    } catch (error) {
      console.error('设置失败:', error)
      this.setData({ loading: false })
      
      wx.showToast({
        title: error.message || '设置失败，请重试',
        icon: 'none'
      })
    }
  },

  // 更新昵称
  async updateNickname(userId, nickName) {
    try {
      const result = await api.updateUserInfo({
        userId: userId,
        nickName: nickName
      })
      
      if (result.code !== 0) {
        throw new Error(result.message || '更新昵称失败')
      }
      
      console.log('更新昵称成功')
    } catch (error) {
      console.error('更新昵称失败:', error)
      throw error
    }
  },

  // 绑定手机号
  async bindPhone(userId, phone) {
    try {
      const result = await api.bindPhone({
        userId: userId,
        phone: phone,
        code: '123456' // 模拟验证码
      })
      
      if (result.code !== 0) {
        throw new Error(result.message || '绑定手机号失败')
      }
      
      console.log('绑定手机号成功')
    } catch (error) {
      console.error('绑定手机号失败:', error)
      throw error
    }
  },

  // 手动输入手机号
  onInputPhone(e) {
    this.setData({
      phone: e.detail.value
    })
  },

  // 显示手机号输入弹窗
  onShowPhoneModal() {
    this.setData({
      showPhoneModal: true
    })
  },

  // 关闭手机号输入弹窗
  onClosePhoneModal() {
    this.setData({
      showPhoneModal: false
    })
  },

  // 确认手动输入手机号
  onConfirmPhone() {
    const { phone } = this.data
    if (!phone || phone.length !== 11) {
      wx.showToast({
        title: '请输入正确的手机号',
        icon: 'none'
      })
      return
    }

    this.setData({
      showPhoneModal: false
    })

    wx.showToast({
      title: '手机号已设置',
      icon: 'success'
    })
  },

  // 返回
  onBack() {
    wx.navigateBack()
  }
}) 