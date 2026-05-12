// pages/hospital/detail.js
const { api } = require('../../utils/cloud-container-standard')
const { processSingleHospitalImage } = require('../../utils/image')

Page({
  data: {
    hospital: null,
    navigation: null,
    loading: true,
    hospitalId: null,
    userLocation: null
  },

  onLoad(options) {
    console.log('医院详情页面加载，参数:', options)
    
    const hospitalId = options.id || options.hospitalId
    if (!hospitalId) {
      wx.showToast({
        title: '缺少医院ID',
        icon: 'none'
      })
      return
    }

    this.setData({ hospitalId })
    this.loadHospitalDetail()
  },

  onShow() {
    // 每次显示页面时刷新数据
    if (this.data.hospitalId) {
      this.loadHospitalDetail()
    }
  },

  // 加载医院详情
  async loadHospitalDetail() {
    try {
      this.setData({ loading: true })
      
      console.log('开始加载医院详情，ID:', this.data.hospitalId)
      
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
        userLongitude: userLocation.longitude,
        userLatitude: userLocation.latitude
      }

      console.log('请求医院详情参数:', params)
      
      const result = await api.hospitalDetail(this.data.hospitalId, params)
      
      console.log('医院详情API返回:', result)
      
      if (result.code === 0 && result.data) {
        const { hospital, navigation } = result.data
        
        // 处理医院图片
        const processedHospital = processSingleHospitalImage(hospital)
        
        this.setData({
          hospital: processedHospital,
          navigation,
          userLocation,
          loading: false
        })
        
        console.log('医院详情加载成功:', {
          hospitalName: processedHospital.name,
          hasNavigation: !!navigation
        })
      } else {
        throw new Error(result.message || '获取医院详情失败')
      }
    } catch (error) {
      console.error('加载医院详情失败:', error)
      this.setData({ loading: false })
      
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

  // 拨打电话
  onCallPhone() {
    const { hospital } = this.data
    if (hospital && hospital.phone) {
      wx.makePhoneCall({
        phoneNumber: hospital.phone,
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
  onNavigateToHospital() {
    const { hospital } = this.data
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

  // 查看医院服务
  onViewServices() {
    const { hospital } = this.data
    if (hospital) {
      wx.navigateTo({
        url: `/pages/service/list?hospitalId=${hospital.id}`
      })
    }
  },

  // 分享医院信息
  onShareHospital() {
    const { hospital } = this.data
    if (hospital) {
      wx.showShareMenu({
        withShareTicket: true,
        menus: ['shareAppMessage', 'shareTimeline']
      })
    }
  },

  // 下拉刷新
  async onPullDownRefresh() {
    await this.loadHospitalDetail()
    wx.stopPullDownRefresh()
  },

  // 分享给朋友
  onShareAppMessage() {
    const { hospital } = this.data
    if (hospital) {
      return {
        title: `${hospital.name} - 专业医疗服务`,
        path: `/pages/hospital/detail?id=${hospital.id}`,
        imageUrl: hospital.logo || '/images/hospital-default.jpg'
      }
    }
    return {
      title: '医院详情',
      path: '/pages/hospital/detail'
    }
  },

  // 分享到朋友圈
  onShareTimeline() {
    const { hospital } = this.data
    if (hospital) {
      return {
        title: `${hospital.name} - 专业医疗服务`,
        query: `id=${hospital.id}`,
        imageUrl: hospital.logo || '/images/hospital-default.jpg'
      }
    }
    return {
      title: '医院详情'
    }
  },

  // 图片加载错误处理
  onImageError(e) {
    console.log('医院图片加载失败')
    
    // 设置默认图片
    const hospital = this.data.hospital
    if (hospital) {
      hospital.logo = '/images/hospital-default.jpg'
      this.setData({
        hospital: hospital
      })
    }
  }
})