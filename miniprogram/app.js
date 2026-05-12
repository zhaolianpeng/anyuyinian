// app.js
const config = require('./config.js')
const { preloadCriticalImages } = require('./utils/imagePreloader')
const { callContainer: requestServer } = require('./utils/cloud-container-standard')

App({
  globalData: {
    userInfo: null,
    hasLogin: false,
    baseUrl: config.baseURL
  },

  onLaunch() {
    // 小程序启动时执行
    console.log('小程序启动')
    
    // 合规：启动时检查隐私授权（不再默认同意）
    this.ensurePrivacyAuthorization().catch(err => {
      console.warn('隐私授权未同意或检查失败：', err)
    })

    // 检查登录状态
    this.checkLoginStatus()
    
    // 预加载关键图片
    this.preloadCriticalImages()
  },

  onShow() {
    // 小程序显示时执行
    console.log('小程序显示')
  },

  onHide() {
    // 小程序隐藏时执行
    console.log('小程序隐藏')
  },

  onError(msg) {
    // 小程序发生错误时执行
    console.error('小程序错误:', msg)
  },

  // 全局分享配置
  onShareAppMessage(res) {
    console.log('全局分享被触发:', res)
    
    // 根据来源页面返回不同的分享内容
    if (res.from === 'button') {
      // 来自按钮点击
      return {
        title: res.target.dataset.title || '安语颐年护理陪诊',
        desc: res.target.dataset.desc || '专业护理陪诊服务，让您就医更安心',
        path: res.target.dataset.path || 'pages/index/index',
        imageUrl: res.target.dataset.imageUrl || ''
      }
    } else {
      // 来自右上角菜单
      return {
        title: '安语颐年护理陪诊',
        desc: '专业护理陪诊服务，让您就医更安心',
        path: 'pages/index/index'
      }
    }
  },

  // 分享到朋友圈
  onShareTimeline(res) {
    console.log('分享到朋友圈被触发:', res)
    
    return {
      title: '安语颐年护理陪诊 - 专业护理陪诊服务',
      query: '',
      imageUrl: ''
    }
  },

  // 检查登录状态
  checkLoginStatus() {
    const token = wx.getStorageSync('token')
    const userInfo = wx.getStorageSync('userInfo')
    
    if (token && userInfo) {
      this.globalData.hasLogin = true
      this.globalData.userInfo = userInfo
    }
  },

  // 设置用户信息
  setUserInfo(userInfo) {
    this.globalData.userInfo = userInfo
    this.globalData.hasLogin = true
    wx.setStorageSync('userInfo', userInfo)
  },

  // 清除用户信息
  clearUserInfo() {
    this.globalData.userInfo = null
    this.globalData.hasLogin = false
    wx.removeStorageSync('token')
    wx.removeStorageSync('userInfo')
    wx.removeStorageSync('userId')
  },

  // 检查是否已登录
  isLoggedIn() {
    return this.globalData.hasLogin
  },

  // 获取用户信息
  getUserInfo() {
    return this.globalData.userInfo
  },

  /**
   * 当需要获取用户隐私授权时触发（由基础库调用）
   * 应展示说明并由用户自主选择同意或不同意
   */
  onNeedPrivacyAuthorization(resolve, reject) {
    const showConsentDialog = () => {
      wx.showModal({
        title: '隐私保护提示',
        content: '为更好地提供服务，请您先阅读并同意《隐私政策》和《用户服务协议》。如不同意，将无法继续使用相关功能。',
        confirmText: '同意',
        cancelText: '不同意',
        success: (res) => {
          if (res.confirm) {
            try { resolve && resolve() } catch (e) {}
          } else {
            try { reject && reject() } catch (e) {}
          }
        },
        fail: () => {
          try { reject && reject() } catch (e) {}
        }
      })
    }

    // 先提供查看隐私协议入口，再给出同意/不同意选择
    if (wx.openPrivacyContract) {
      wx.openPrivacyContract({
        success: () => { showConsentDialog() },
        fail: () => { showConsentDialog() },
        complete: () => {}
      })
    } else {
      showConsentDialog()
    }
  },

  /**
   * 主动检查并触发隐私授权流程
   */
  async ensurePrivacyAuthorization() {
    try {
      if (!wx.getPrivacySetting || !wx.requirePrivacyAuthorize) {
        console.log('当前基础库不支持隐私授权流程，跳过检查')
        return
      }

      const setting = await new Promise((resolve, reject) => {
        wx.getPrivacySetting({
          success: resolve,
          fail: reject
        })
      })

      if (setting && setting.needAuthorization) {
        // 触发隐私授权流程，回调中由 onNeedPrivacyAuthorization 处理
        await new Promise((resolve, reject) => {
          wx.requirePrivacyAuthorize({
            success: resolve,
            fail: reject
          })
        })
      }
    } catch (error) {
      console.warn('ensurePrivacyAuthorization 执行异常：', error)
      throw error
    }
  },

  /**
   * 对外提供查看隐私协议入口
   */
  openPrivacyContract() {
    if (wx.openPrivacyContract) {
      wx.openPrivacyContract({})
    } else {
      wx.showToast({ title: '当前版本不支持查看隐私协议', icon: 'none' })
    }
  },

  /**
   * 统一的服务端调用方法
   */
  async callContainer(path, method = 'GET', data = {}, options = {}) {
    return requestServer(path, method, data, options)
  },

  // 预加载关键图片
  async preloadCriticalImages() {
    try {
      console.log('开始预加载关键图片...')
      
      await preloadCriticalImages(
        (progress) => {
          console.log(`关键图片预加载进度: ${progress.progress}% (${progress.loaded}/${progress.total})`)
        },
        (result) => {
          console.log('关键图片预加载完成:', result)
        }
      )
    } catch (error) {
      console.error('关键图片预加载失败:', error)
    }
  }
})