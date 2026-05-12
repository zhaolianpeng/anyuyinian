// pages/hospital/list.js
const { api } = require('../../utils/cloud-container-standard')
const { processHospitalImages } = require('../../utils/image')

Page({
  data: {
    hospitals: [],
    loading: true,
    userLocation: null,
    page: 1,
    pageSize: 10,
    hasMore: true,
    refreshing: false
  },

  onLoad(options) {
    console.log('医院列表页面加载')
    this.loadHospitals()
  },

  onShow() {
    // 每次显示页面时刷新数据
    this.loadHospitals()
  },

  // 加载医院列表
  async loadHospitals(refresh = false) {
    try {
      if (refresh) {
        this.setData({ 
          refreshing: true,
          page: 1,
          hasMore: true
        })
      } else {
        this.setData({ loading: true })
      }
      
      // 获取用户位置
      let userLocation = null
      try {
        userLocation = await this.getUserLocation()
        console.log('用户位置:', userLocation)
      } catch (error) {
        console.log('获取用户位置失败，使用默认位置')
        userLocation = {
          longitude: 121.4737,
          latitude: 31.2304
        }
      }

      // 构建请求参数
      const params = {
        longitude: userLocation.longitude,
        latitude: userLocation.latitude,
        page: this.data.page,
        pageSize: this.data.pageSize
      }

      console.log('请求医院列表参数:', params)
      
      const result = await api.hospitalList(params)
      
      console.log('医院列表API返回:', result)
      
      if (result.code === 0 && result.data) {
        const { list, hasMore } = result.data
        
        // 处理医院图片
        const processedList = processHospitalImages(list || [])
        
        if (refresh) {
          this.setData({
            hospitals: processedList,
            hasMore,
            refreshing: false
          })
        } else {
          this.setData({
            hospitals: this.data.page === 1 ? processedList : [...this.data.hospitals, ...processedList],
            hasMore,
            loading: false,
            userLocation
          })
        }
        
        console.log('医院列表加载成功，数量:', processedList.length)
      } else {
        throw new Error(result.message || '获取医院列表失败')
      }
    } catch (error) {
      console.error('加载医院列表失败:', error)
      this.setData({ 
        loading: false,
        refreshing: false
      })
      
      wx.showToast({
        title: error.message || '加载失败，请重试',
        icon: 'none'
      })
    }
  },

  // 获取用户位置
  getUserLocation() {
    return new Promise((resolve, reject) => {
      wx.getLocation({
        type: 'gcj02',
        success: (res) => {
          resolve({
            latitude: res.latitude,
            longitude: res.longitude
          })
        },
        fail: (err) => {
          console.error('获取位置失败:', err)
          reject(err)
        }
      })
    })
  },

  // 跳转到医院详情
  onHospitalDetail(e) {
    const { id } = e.currentTarget.dataset
    if (id) {
      wx.navigateTo({
        url: `/pages/hospital/detail?id=${id}`
      })
    }
  },

  // 拨打电话
  onCallPhone(e) {
    const { phone } = e.currentTarget.dataset
    if (phone) {
      wx.makePhoneCall({
        phoneNumber: phone,
        success: () => {
          console.log('拨打电话成功')
        },
        fail: (err) => {
          console.error('拨打电话失败:', err)
          wx.showToast({
            title: '拨号失败',
            icon: 'none'
          })
        }
      })
    } else {
      wx.showToast({
        title: '暂无联系电话',
        icon: 'none'
      })
    }
  },

  // 导航到医院
  onNavigateToHospital(e) {
    const { hospital } = e.currentTarget.dataset
    if (hospital && hospital.latitude && hospital.longitude) {
      wx.openLocation({
        latitude: hospital.latitude,
        longitude: hospital.longitude,
        name: hospital.name,
        address: hospital.address,
        success: () => {
          console.log('打开导航成功')
        },
        fail: (err) => {
          console.error('打开导航失败:', err)
          wx.showToast({
            title: '导航失败',
            icon: 'none'
          })
        }
      })
    } else {
      wx.showToast({
        title: '暂无位置信息',
        icon: 'none'
      })
    }
  },

  // 下拉刷新
  async onPullDownRefresh() {
    await this.loadHospitals(true)
    wx.stopPullDownRefresh()
  },

  // 上拉加载更多
  async onReachBottom() {
    if (this.data.hasMore && !this.data.loading) {
      this.setData({
        page: this.data.page + 1
      })
      await this.loadHospitals()
    }
  },

  // 分享给朋友
  onShareAppMessage() {
    return {
      title: '附近医院 - 专业医疗服务',
      path: '/pages/hospital/list'
    }
  },

  // 分享到朋友圈
  onShareTimeline() {
    return {
      title: '附近医院 - 专业医疗服务'
    }
  },

  // 图片加载错误处理
  onImageError(e) {
    const { index } = e.currentTarget.dataset
    console.log('医院图片加载失败，索引:', index)
    
    // 设置默认图片
    const hospitals = this.data.hospitals
    if (hospitals && hospitals[index]) {
      hospitals[index].logo = '/images/hospital-default.jpg'
      this.setData({
        hospitals: hospitals
      })
    }
  }
})