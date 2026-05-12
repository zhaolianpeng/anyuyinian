// pages/user/address/list.js
const app = getApp()
const { api } = require('../../../utils/cloud-container-standard')

Page({
  data: {
    addressList: [],
    loading: false,
    userId: null,
    selectMode: false
  },

  onLoad(options) {
    const userId = wx.getStorageSync('userId')
    this.setData({ 
      userId,
      selectMode: options.selectMode === 'true'
    })
  },

  onShow() {
    this.loadAddressList()
  },

  async loadAddressList() {
    if (!this.data.userId) {
      wx.showToast({
        title: '请先登录',
        icon: 'none'
      })
      return
    }

    this.setData({ loading: true })

    try {
      console.log('开始加载地址列表，userId:', this.data.userId)
      
      const result = await api.userAddress({ userId: this.data.userId })
      
      console.log('地址列表API返回:', result)
      
      if (result.code === 0) {
        this.setData({
          addressList: result.data || []
        })
        console.log('地址列表加载成功，数量:', result.data ? result.data.length : 0)
      } else {
        throw new Error(result.errorMsg || result.message || '获取地址列表失败')
      }
    } catch (error) {
      console.error('获取地址列表失败:', error)
      wx.showToast({
        title: error.message || '网络错误，请重试',
        icon: 'none',
        duration: 3000
      })
    } finally {
      this.setData({ loading: false })
    }
  },

  // 添加地址
  addAddress() {
    wx.navigateTo({
      url: '/pages/user/address/add'
    })
  },

  // 编辑地址
  editAddress(e) {
    const address = e.currentTarget.dataset.address
    wx.navigateTo({
      url: `/pages/user/address/add?edit=1&address=${JSON.stringify(address)}`
    })
  },

  // 删除地址
  deleteAddress(e) {
    const id = e.currentTarget.dataset.id
    
    wx.showModal({
      title: '确认删除',
      content: '确定要删除这个地址吗？',
      success: (res) => {
        if (res.confirm) {
          this.deleteAddressById(id)
        }
      }
    })
  },

  // 删除地址API调用
  async deleteAddressById(id) {
    try {
      console.log('删除地址，id:', id)
      
      const result = await api.userAddressDelete(id)
      
      console.log('删除地址API返回:', result)
      
      if (result.code === 0) {
        wx.showToast({
          title: '删除成功',
          icon: 'success'
        })
        this.loadAddressList()
      } else {
        throw new Error(result.errorMsg || result.message || '删除失败')
      }
    } catch (error) {
      console.error('删除地址失败:', error)
      wx.showToast({
        title: error.message || '网络错误，请重试',
        icon: 'none',
        duration: 3000
      })
    }
  },

  // 选择地址（用于订单页面）
  selectAddress(e) {
    if (!this.data.selectMode) return
    
    const address = e.currentTarget.dataset.address
    const pages = getCurrentPages()
    const prevPage = pages[pages.length - 2]
    
    if (prevPage && prevPage.setSelectedAddress) {
      prevPage.setSelectedAddress(address)
      wx.navigateBack()
    }
  },

  // 下拉刷新
  onPullDownRefresh() {
    this.loadAddressList()
    wx.stopPullDownRefresh()
  }
})