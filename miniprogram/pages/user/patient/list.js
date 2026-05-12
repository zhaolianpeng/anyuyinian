// pages/user/patient/list.js
const app = getApp()
const { api } = require('../../../utils/cloud-container-standard')

Page({
  data: {
    patientList: [],
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
    this.loadPatientList()
  },

  async loadPatientList() {
    if (!this.data.userId) {
      wx.showToast({
        title: '请先登录',
        icon: 'none'
      })
      return
    }

    this.setData({ loading: true })

    try {
      console.log('开始加载患者列表，userId:', this.data.userId)
      
      const result = await api.userPatient({ userId: this.data.userId })
      
      console.log('患者列表API返回:', result)
      
      if (result.code === 0) {
        this.setData({
          patientList: result.data || []
        })
        console.log('患者列表加载成功，数量:', result.data ? result.data.length : 0)
      } else {
        throw new Error(result.errorMsg || result.message || '获取患者列表失败')
      }
    } catch (error) {
      console.error('获取患者列表失败:', error)
      wx.showToast({
        title: error.message || '网络错误，请重试',
        icon: 'none',
        duration: 3000
      })
    } finally {
      this.setData({ loading: false })
    }
  },

  // 添加就诊人
  addPatient() {
    wx.navigateTo({
      url: '/pages/user/patient/add'
    })
  },

  // 编辑就诊人
  editPatient(e) {
    const patient = e.currentTarget.dataset.patient
    wx.navigateTo({
      url: `/pages/user/patient/add?edit=1&patient=${JSON.stringify(patient)}`
    })
  },

  // 删除患者
  deletePatient(e) {
    const id = e.currentTarget.dataset.id
    
    wx.showModal({
      title: '确认删除',
      content: '确定要删除这个患者吗？',
      success: (res) => {
        if (res.confirm) {
          this.deletePatientById(id)
        }
      }
    })
  },

  // 删除患者API调用
  async deletePatientById(id) {
    try {
      console.log('删除患者，id:', id)
      
      const result = await api.userPatientDelete(id)
      
      console.log('删除患者API返回:', result)
      
      if (result.code === 0) {
        wx.showToast({
          title: '删除成功',
          icon: 'success'
        })
        this.loadPatientList()
      } else {
        throw new Error(result.errorMsg || result.message || '删除失败')
      }
    } catch (error) {
      console.error('删除患者失败:', error)
      wx.showToast({
        title: error.message || '网络错误，请重试',
        icon: 'none',
        duration: 3000
      })
    }
  },

    // 选择患者（用于订单页面）
    selectPatient(e) {
      if (!this.data.selectMode) return
      
      const patient = e.currentTarget.dataset.patient
      const pages = getCurrentPages()
      const prevPage = pages[pages.length - 2]
      
      if (prevPage && prevPage.setSelectedPatient) {
        prevPage.setSelectedPatient(patient)
        wx.navigateBack()
      }
    },

    // 下拉刷新
    onPullDownRefresh() {
      this.loadPatientList()
      wx.stopPullDownRefresh()
    }
  })