// pages/promoter/home.js
const app = getApp()
const { api } = require('../../utils/cloud-container-standard')
const { getCurrentUserId } = require('../../utils/user-id-compatibility')

Page({
  data: {
    loading: false,
    promoterInfo: null,
    activeTab: 0, // 0-推广信息，1-佣金记录，2-提现记录
    commissionList: [],
    cashoutList: [],
    commissionPage: 1,
    cashoutPage: 1,
    hasMoreCommission: true,
    hasMoreCashout: true,
    showCashoutModal: false,
    cashoutForm: {
      amount: '',
      method: 'wechat',
      account: ''
    }
  },

  onLoad() {
    console.log('推广中心页面加载')
    this.loadPromoterInfo()
  },

  onShow() {
    console.log('推广中心页面显示')
    if (this.data.activeTab === 1) {
      this.loadCommissionList(true)
    } else if (this.data.activeTab === 2) {
      this.loadCashoutList(true)
    }
  },

  // 加载推广员信息
  async loadPromoterInfo() {
    try {
      this.setData({ loading: true })
      
      const userId = getCurrentUserId()
      if (!userId) {
        console.log('用户未登录，跳转到登录页面')
        wx.navigateTo({ url: '/pages/login/login' })
        return
      }

      console.log('开始获取推广员信息，userId:', userId)
      const result = await api.promoterInfo({ userId })
      
      if (result.code === 0) {
        this.setData({ 
          promoterInfo: result.data,
          loading: false
        })
        console.log('推广员信息获取成功:', result.data)
      } else {
        throw new Error(result.errorMsg || '获取推广员信息失败')
      }
    } catch (error) {
      console.error('加载推广员信息失败:', error)
      this.setData({ loading: false })
      
      wx.showToast({
        title: '加载失败，请重试',
        icon: 'none'
      })
    }
  },

  // 切换标签页
  onTabChange(e) {
    const activeTab = parseInt(e.currentTarget.dataset.index)
    console.log('切换标签页:', activeTab)
    this.setData({ activeTab })
    
    // 添加安全检查，确保数组存在且不为null
    if (activeTab === 1 && (!this.data.commissionList || this.data.commissionList.length === 0)) {
      this.loadCommissionList(true)
    } else if (activeTab === 2 && (!this.data.cashoutList || this.data.cashoutList.length === 0)) {
      this.loadCashoutList(true)
    }
  },

  // 加载佣金记录列表
  async loadCommissionList(refresh = false) {
    try {
      const userId = getCurrentUserId()
      if (!userId) return

      const page = refresh ? 1 : this.data.commissionPage
      
      console.log('开始获取佣金记录，page:', page)
      const result = await api.commissionList({ 
        userId, 
        page, 
        pageSize: 20 
      })
      
      if (result.code === 0) {
        const { list, hasMore } = result.data
        // 添加安全检查，确保list不为null
        const safeList = list || []
        // 添加安全检查，确保this.data.commissionList存在
        const currentList = this.data.commissionList || []
        const commissionList = refresh ? safeList : [...currentList, ...safeList]
        
        this.setData({
          commissionList,
          commissionPage: page + 1,
          hasMoreCommission: hasMore
        })
        
        console.log('佣金记录获取成功，数量:', safeList.length)
      } else {
        throw new Error(result.errorMsg || '获取佣金记录失败')
      }
    } catch (error) {
      console.error('加载佣金记录失败:', error)
      wx.showToast({
        title: '加载失败，请重试',
        icon: 'none'
      })
    }
  },

  // 加载提现记录列表
  async loadCashoutList(refresh = false) {
    try {
      const userId = getCurrentUserId()
      if (!userId) return

      const page = refresh ? 1 : this.data.cashoutPage
      
      console.log('开始获取提现记录，page:', page)
      const result = await api.cashoutList({ 
        userId, 
        page, 
        pageSize: 20 
      })
      
      if (result.code === 0) {
        const { list, hasMore } = result.data
        // 添加安全检查，确保list不为null
        const safeList = list || []
        // 添加安全检查，确保this.data.cashoutList存在
        const currentList = this.data.cashoutList || []
        const cashoutList = refresh ? safeList : [...currentList, ...safeList]
        
        this.setData({
          cashoutList,
          cashoutPage: page + 1,
          hasMoreCashout: hasMore
        })
        
        console.log('提现记录获取成功，数量:', safeList.length)
      } else {
        throw new Error(result.errorMsg || '获取提现记录失败')
      }
    } catch (error) {
      console.error('加载提现记录失败:', error)
      wx.showToast({
        title: '加载失败，请重试',
        icon: 'none'
      })
    }
  },

  // 显示提现申请弹窗
  onShowCashoutModal() {
    this.setData({ 
      showCashoutModal: true,
      cashoutForm: {
        amount: '',
        method: 'wechat',
        account: ''
      }
    })
  },

  // 隐藏提现申请弹窗
  onHideCashoutModal() {
    this.setData({ showCashoutModal: false })
  },

  // 提现表单输入
  onCashoutFormInput(e) {
    const { field } = e.currentTarget.dataset
    const { value } = e.detail
    this.setData({
      [`cashoutForm.${field}`]: value
    })
  },

  // 提交提现申请
  async onSubmitCashout() {
    try {
      const { amount, method, account } = this.data.cashoutForm
      
      if (!amount || !method || !account) {
        wx.showToast({
          title: '请填写完整信息',
          icon: 'none'
        })
        return
      }

      const amountNum = parseFloat(amount)
      if (isNaN(amountNum) || amountNum <= 0) {
        wx.showToast({
          title: '请输入有效金额',
          icon: 'none'
        })
        return
      }

      const userId = getCurrentUserId()
      if (!userId) return

      console.log('开始提交提现申请:', { userId, amount: amountNum, method, account })
      
      const result = await api.applyCashout({
        userId,
        amount: amountNum,
        method,
        account
      })
      
      if (result.code === 0) {
        wx.showToast({
          title: '申请提交成功',
          icon: 'success'
        })
        
        this.onHideCashoutModal()
        this.loadCashoutList(true) // 刷新提现记录
      } else {
        throw new Error(result.errorMsg || '申请提现失败')
      }
    } catch (error) {
      console.error('申请提现失败:', error)
      wx.showToast({
        title: error.message || '申请失败，请重试',
        icon: 'none'
      })
    }
  },

  // 分享小程序
  onShareAppMessage() {
    const promoterInfo = this.data.promoterInfo
    if (!promoterInfo) {
      return {
        title: '安语颐年护理陪诊 - 专业护理服务',
        path: '/pages/index/index'
      }
    }

    return {
      title: '安语颐年护理陪诊 - 专业护理服务',
      path: `/pages/index/index?promoterCode=${promoterInfo.promoterCode}`,
      imageUrl: promoterInfo.qrCodeUrl || '/images/service-default.jpg'
    }
  },

  // 复制推广码
  onCopyPromoterId() {
    const promoterInfo = this.data.promoterInfo
    if (!promoterInfo) return

    wx.setClipboardData({
      data: promoterInfo.promoterCode,
      success: () => {
        wx.showToast({
          title: '推广码已复制',
          icon: 'success'
        })
      }
    })
  },

  // 生成并预览二维码
  async onPreviewQrCode() {
    const promoterInfo = this.data.promoterInfo
    if (!promoterInfo || !promoterInfo.promoterCode) {
      wx.showToast({
        title: '推广码不存在',
        icon: 'none'
      })
      return
    }

    try {
      wx.showLoading({ title: '生成二维码中...' })
      
      // 构建小程序页面路径
      const pagePath = `pages/index/index?promoterCode=${promoterInfo.promoterCode}`
      
      // 调用微信二维码生成API
      const qrCodeResult = await this.generateWechatQRCode(pagePath)
      
      wx.hideLoading()
      
      if (qrCodeResult.success) {
        // 预览二维码
        wx.previewImage({
          urls: [qrCodeResult.qrCodeUrl],
          current: qrCodeResult.qrCodeUrl
        })
      } else {
        throw new Error(qrCodeResult.error || '生成二维码失败')
      }
    } catch (error) {
      wx.hideLoading()
      console.error('生成二维码失败:', error)
      wx.showToast({
        title: '生成二维码失败，请重试',
        icon: 'none'
      })
    }
  },

  // 生成微信二维码
  generateWechatQRCode(pagePath) {
    return new Promise((resolve) => {
      // 使用微信小程序的二维码生成API
      wx.createQRCode({
        path: pagePath,
        width: 300,
        success: (res) => {
          console.log('微信二维码生成成功:', res)
          resolve({
            success: true,
            qrCodeUrl: res.path
          })
        },
        fail: (error) => {
          console.error('微信二维码生成失败:', error)
          // 如果微信API失败，尝试使用后端API生成
          this.generateBackendQRCode(pagePath).then(resolve).catch(() => {
            resolve({
              success: false,
              error: '二维码生成失败'
            })
          })
        }
      })
    })
  },

  // 使用后端API生成二维码（备选方案）
  async generateBackendQRCode(pagePath) {
    try {
      const result = await api.generateQRCode({ 
        promoterCode: this.data.promoterInfo.promoterCode 
      })
      
      if (result.code === 0) {
        return {
          success: true,
          qrCodeUrl: result.data.qrCodeUrl
        }
      } else {
        throw new Error(result.errorMsg || '后端二维码生成失败')
      }
    } catch (error) {
      console.error('后端二维码生成失败:', error)
      throw error
    }
  },

  // 下拉刷新
  onPullDownRefresh() {
    if (this.data.activeTab === 0) {
      this.loadPromoterInfo()
    } else if (this.data.activeTab === 1) {
      this.loadCommissionList(true)
    } else if (this.data.activeTab === 2) {
      this.loadCashoutList(true)
    }
    
    setTimeout(() => {
      wx.stopPullDownRefresh()
    }, 1000)
  },

  // 上拉加载更多
  onReachBottom() {
    if (this.data.activeTab === 1 && this.data.hasMoreCommission) {
      this.loadCommissionList()
    } else if (this.data.activeTab === 2 && this.data.hasMoreCashout) {
      this.loadCashoutList()
    }
  }
})