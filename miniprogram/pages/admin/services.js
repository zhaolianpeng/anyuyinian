// pages/admin/services.js
const { api } = require('../../utils/cloud-container-standard')
const app = getApp()

Page({
  data: {
    services: [],
    loading: false,
    page: 1,
    pageSize: 20,
    hasMore: true,
    categories: [],
    selectedCategoryIndex: -1,
    selectedCategory: '',
    showPriceModal: false,
    currentService: null,
    newPrice: '',
    newOriginalPrice: '',
    reason: '',
    submitting: false
  },

  onLoad(options) {
    this.loadCategories()
    this.loadServices()
  },

  onShow() {
    // 页面显示时刷新数据
    this.refreshData()
  },

  onPullDownRefresh() {
    this.refreshData()
  },

  onReachBottom() {
    if (this.data.hasMore && !this.data.loading) {
      this.loadMoreServices()
    }
  },

  // 加载服务分类
  loadCategories() {
    api.serviceList({
      page: 1,
      pageSize: 1000
    }).then(res => {
      if (res.code === 0) {
        const categories = [...new Set(res.data.list.map(item => item.category))]
        this.setData({
          categories: categories
        })
      }
    }).catch(err => {
      console.error('获取分类失败:', err)
      wx.showToast({
        title: '获取分类失败',
        icon: 'none'
      })
    })
  },

  // 加载服务列表
  loadServices(refresh = false) {
    if (this.data.loading) return

    this.setData({ loading: true })

    const page = refresh ? 1 : this.data.page
    // 从本地缓存获取管理员ID，避免全局变量为空
    const adminInfo = wx.getStorageSync('adminInfo') || {}
    const adminUserId = adminInfo.userId || app.globalData.adminUserId
    
    const queryParams = {
      page: page,
      pageSize: this.data.pageSize,
      adminUserId: adminUserId
    }

    if (this.data.selectedCategory) {
      queryParams.category = this.data.selectedCategory
    }

    console.log('请求参数:', queryParams)
    api.adminServices(queryParams).then(res => {
      console.log('API响应:', res)
      if (res.code === 0) {
        const newServicesRaw = res.data.list || []
        console.log('原始服务数据:', newServicesRaw)
        const newServices = newServicesRaw.map(s => ({
          ...s,
          displayPrice: (Number(s.price) || 0).toFixed(2),
          displayOriginalPrice: (Number(s.originalPrice) || 0).toFixed(2),
          displayUpdatedAt: formatDateTime(s.updatedAt)
        }))
        const services = refresh ? newServices : [...this.data.services, ...newServices]
        
        this.setData({
          services: services,
          page: page + 1,
          hasMore: res.data.hasMore,
          loading: false
        })
      } else {
        wx.showToast({
          title: res.errorMsg || '获取服务列表失败',
          icon: 'none'
        })
        this.setData({ loading: false })
      }
    }).catch(err => {
      console.error('获取服务列表失败:', err)
      wx.showToast({
        title: '网络错误',
        icon: 'none'
      })
      this.setData({ loading: false })
    })
  },

  // 加载更多服务
  loadMoreServices() {
    this.loadServices(false)
  },

  // 刷新数据
  refreshData() {
    this.setData({
      page: 1,
      hasMore: true,
      services: [],
      selectedCategoryIndex: -1,
      selectedCategory: ''
    })
    this.loadServices(true)
    wx.stopPullDownRefresh()
  },

  // 选择分类
  onCategoryChange(e) {
    const index = e.detail.value
    const category = this.data.categories[index] || ''
    this.setData({
      selectedCategoryIndex: index,
      selectedCategory: category,
      page: 1,
      hasMore: true,
      services: []
    })
    this.loadServices(true)
  },

  // 显示价格修改弹窗
  showPriceModal(e) {
    const service = e.currentTarget.dataset.service
    this.setData({
      showPriceModal: true,
      currentService: service,
      newPrice: service.price.toString(),
      newOriginalPrice: service.originalPrice.toString(),
      reason: ''
    })
  },

  // 隐藏价格修改弹窗
  hidePriceModal() {
    this.setData({
      showPriceModal: false,
      currentService: null,
      newPrice: '',
      newOriginalPrice: '',
      reason: '',
      submitting: false
    })
  },

  // 价格输入
  onPriceInput(e) {
    this.setData({
      newPrice: e.detail.value
    })
  },

  // 原价输入
  onOriginalPriceInput(e) {
    this.setData({
      newOriginalPrice: e.detail.value
    })
  },

  // 原因输入
  onReasonInput(e) {
    this.setData({
      reason: e.detail.value
    })
  },

  // 提交价格修改
  submitPriceUpdate() {
    if (this.data.submitting) return

    const newPrice = parseFloat(this.data.newPrice)
    const newOriginalPrice = parseFloat(this.data.newOriginalPrice)
    const reason = this.data.reason.trim()

    // 验证输入
    if (isNaN(newPrice) || newPrice < 0) {
      wx.showToast({
        title: '请输入有效的价格',
        icon: 'none'
      })
      return
    }

    if (isNaN(newOriginalPrice) || newOriginalPrice < 0) {
      wx.showToast({
        title: '请输入有效的原价',
        icon: 'none'
      })
      return
    }

    if (!reason) {
      wx.showToast({
        title: '请输入修改原因',
        icon: 'none'
      })
      return
    }

    this.setData({ submitting: true })

    // 再次获取管理员ID，防止作用域变量不可用
    const adminInfo = wx.getStorageSync('adminInfo') || {}
    const adminUserId = adminInfo.userId || app.globalData.adminUserId

    api.adminUpdateServicePrice({
      serviceId: this.data.currentService.id,
      newPrice: newPrice,
      newOriginalPrice: newOriginalPrice,
      reason: reason,
      adminUserId: adminUserId
    }).then(res => {
      if (res.code === 0) {
        wx.showToast({
          title: '价格修改成功',
          icon: 'success'
        })
        
        // 更新本地数据
        const services = this.data.services.map(service => {
          if (service.id === this.data.currentService.id) {
            return {
              ...service,
              price: newPrice,
              originalPrice: newOriginalPrice,
              displayPrice: newPrice.toFixed(2),
              displayOriginalPrice: newOriginalPrice.toFixed(2),
              displayUpdatedAt: formatDateTime(new Date())
            }
          }
          return service
        })
        
        this.setData({
          services: services
        })
        
        this.hidePriceModal()
      } else {
        wx.showToast({
          title: res.errorMsg || '价格修改失败',
          icon: 'none'
        })
      }
    }).catch(err => {
      console.error('价格修改失败:', err)
      wx.showToast({
        title: '网络错误',
        icon: 'none'
      })
    }).finally(() => {
      this.setData({ submitting: false })
    })
  },

  // 阻止事件冒泡
  preventClose() {
    // 空函数，用于阻止事件冒泡
  },
  
})

// 辅助函数（文件作用域）
function formatDateTime(ts) {
  const d = new Date(ts)
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const hh = String(d.getHours()).padStart(2, '0')
  const mm = String(d.getMinutes()).padStart(2, '0')
  return `${y}-${m}-${day} ${hh}:${mm}`
}