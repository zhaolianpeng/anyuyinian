const { api } = require('../../utils/cloud-container-standard')
const { preloadHomeImages, preloadCriticalImages } = require('../../utils/imagePreloader')
const { getServiceImages } = require('../../utils/imageService')
const { processImageUrl } = require('../../utils/image')
const { fetchServerImages, fetchServerVideoToLocal, normalizeServerImagePath } = require('../../utils/serverMedia')

Page({
  data: {
    banners: [],
    navigations: [],
    services: [],
    hospitals: [],
    loading: true,
    userInfo: null,
    networkError: false,

    // 视频相关
    videoUrl: '',
    videoPoster: '', // 视频封面图（使用网络地址）
    mainBannerImage: '',
    mainBannerSource: '/static/fuwu_1.jpeg',
    videoSource: 'https://api.succ.online/anyuyinian/video/video.mp4',
    videoPosterSource: '/static/Wechat-IMG36.jpg',
    videoLoadError: false, // 视频加载错误状态
    fallbackMode: false, // 是否使用降级模式（只显示封面图）
    isVideoPlaying: false, // 视频是否正在播放
    currentSwiperIndex: 0, // 当前轮播索引

    caregiverServices: [],
    serviceCategories: [
      { name: '居家照护', value: '居家照护' }
    ],
    currentServiceCategory: '居家照护',
    currentServices: [],
    allCaregiverServices: {},

  },

  onLoad() {
    this.initHome()
    this.preloadCriticalImages()
    this.loadServerHeroImages()
    this.loadVideo()
    
    // 创建视频上下文
    this.videoContext = wx.createVideoContext('homeVideo', this)
  },

  onShow() {
    if (this.data.videoLoadError && !this.data.fallbackMode) {
      this.enableFallbackMode()
    }
    
    if (this.data.currentSwiperIndex === 1 && this.data.videoUrl && !this.data.fallbackMode && !this.data.videoLoadError) {
      setTimeout(() => {
        if (this.data.currentSwiperIndex === 1 && !this.data.fallbackMode && !this.data.isVideoPlaying) {
          this.enableFallbackMode()
        }
      }, 8000)
    }
  },

  onUnload() {
    // 清理定时器
    if (this.videoValidationTimeout) {
      clearTimeout(this.videoValidationTimeout)
    }
  },

  // 加载服务分类
  async loadServiceCategories() {
    try {
      const res = await api.serviceCategories()
      
      if (res.code === 0 && res.data && res.data.categories) {
        const categories = res.data.categories
        
        const filteredCategories = categories.filter(cat => cat.value !== '')
        
        this.setData({
          serviceCategories: filteredCategories
        })
        
        // 如果当前选中的分类不在新的分类列表中，则选择第一个分类
        const currentCategoryExists = filteredCategories.some(cat => cat.value === this.data.currentServiceCategory)
        if (!currentCategoryExists && filteredCategories.length > 0) {
          this.setData({
            currentServiceCategory: filteredCategories[0].value
          })
        }
      } else {
        console.error('服务分类加载失败:', res)
      }
    } catch (error) {
      console.error('加载服务分类失败:', error)
    }
  },

  onServiceCategoryTap(e) {
    const { category } = e.currentTarget.dataset
    this.setData({ 
      currentServiceCategory: category 
    })
    
    wx.showLoading({
      title: '加载中...',
      mask: true
    })
    
    this.fetchServicesByCategory(category)
  },

  async fetchServicesByCategory(category) {
    try {
      const result = await api.serviceList({
        category: category,
        page: 1,
        pageSize: 10
      })
      
      if (result.code === 0 && result.data) {
        const services = result.data.list || []
        
        const servicesWithImages = await getServiceImages(services)

        const processedServices = servicesWithImages.map(service => ({
          constDirectImageUrl: service.imageTempUrl || processImageUrl(service.imageUrl || service.imageCosId || ''),
          id: service.id,
          name: service.name,
          description: service.description,
          price: service.price,
          originalPrice: service.originalPrice,
          serverImagePath: normalizeServerImagePath(service.imageTempUrl || service.imageUrl || service.imageCosId),
          imageUrl: '',
          category: service.category,
          serviceitemid: service.id,
          icon: this.getCategoryIcon(service.category),
          linkUrl: `/pages/service/detail?id=${service.id}`,
          sort: service.sort || 0,
          status: service.status || 1,
          detailImages: service.detailImages,
          formConfig: service.formConfig
        }))

        const imageMap = await fetchServerImages(processedServices.map(service => service.serverImagePath))

        const hydratedServices = processedServices.map(service => ({
          ...service,
          imageUrl: imageMap[service.serverImagePath] || service.constDirectImageUrl || ''
        }))
        
        const limitedServices = hydratedServices.slice(0, 5)
        
        this.setData({
          currentServices: limitedServices,
        })
        
        const allServices = { ...this.data.allCaregiverServices }
        allServices[category] = hydratedServices
        this.setData({
          allCaregiverServices: allServices
        })
        
        wx.hideLoading()
        
        if (limitedServices.length === 0) {
          wx.showToast({
            title: '该分类暂无服务',
            icon: 'none',
            duration: 2000
          })
        }
        
      } else {
        throw new Error(result.errorMsg || result.message || '获取服务数据失败')
      }
      
    } catch (error) {
      console.error(`获取${category}分类服务数据失败:`, error)
      
      wx.hideLoading()
      
      wx.showToast({
        title: '加载失败，请重试',
        icon: 'none',
        duration: 2000
      })
      
      this.fallbackToLocalData(category)
    }
  },

  fallbackToLocalData(category) {
    const allServices = this.data.allCaregiverServices
    let currentServices = []
    
    if (allServices && allServices[category]) {
      currentServices = allServices[category].slice(0, 5)
    }
    
    this.setData({
      currentServices: currentServices
    })
  },

  getCategoryIcon(category) {
    const iconMap = {
      '居家照护': '🏠',
      '医院陪诊': '🏥',
      '周期护理': '💊',
      '家政服务': '🧹',
      '预约咨询': '📋',
      '智慧养老': '🤖'
    }
    return iconMap[category] || '🏥'
  },

  onBookService(e) {
    const { service } = e.currentTarget.dataset
    wx.navigateTo({
      url: `/pages/service/detail?id=${service.id}`,
      fail: (error) => {
        console.error('跳转服务详情页面失败:', error)
        wx.showToast({
          title: '页面跳转失败',
          icon: 'none'
        })
      }
    })
  },

  onAllServices() {
    wx.navigateTo({
      url: '/pages/service/list',
      fail: (error) => {
        console.error('跳转服务列表页面失败:', error)
        wx.showToast({
          title: '页面跳转失败',
          icon: 'none'
        })
      }
    })
  },

  onQuickAppointment() {
    this.getAppointmentServiceId().then(serviceId => {
      if (serviceId) {
        wx.navigateTo({
          url: `/pages/service/detail?id=${serviceId}`,
          fail: (error) => {
            console.error('跳转预约咨询服务页面失败:', error)
            wx.showToast({
              title: '页面跳转失败',
              icon: 'none'
            })
          }
        })
      } else {
        wx.showToast({
          title: '预约服务暂不可用',
          icon: 'none'
        })
        console.error('预约咨询服务ID获取失败')
      }
    }).catch(error => {
      console.error('获取预约咨询服务ID失败:', error)
      wx.showToast({
        title: '服务获取失败',
        icon: 'none'
      })
    })
  },

  async getAppointmentServiceId() {
    try {
      if (this.data.services && this.data.services.length > 0) {
        const appointmentService = this.data.services.find(service => 
          service.category === '预约咨询' || service.name === '预约咨询服务'
        )
        
        if (appointmentService) {
          return appointmentService.serviceitemid || 34
        }
      }
      
      if (this.data.caregiverServices && this.data.caregiverServices.length > 0) {
        const appointmentService = this.data.caregiverServices.find(service => 
          service.category === '预约咨询' || service.name === '预约咨询服务'
        )
        
        if (appointmentService) {
          return appointmentService.serviceitemid || 34
        }
      }

      return 34
      
    } catch (error) {
      console.error('获取预约咨询服务ID失败:', error)
      return 34
    }
  },

  onOnlineConsultation() {
    try {
      wx.navigateTo({
        url: '/pages/consultation/consultation',
        fail: (err) => {
          console.error('跳转咨询页面失败:', err)
          wx.showToast({
            title: '跳转失败，请重试',
            icon: 'none'
          })
        }
      })
    } catch (error) {
      console.error('在线咨询跳转异常:', error)
      wx.showToast({
        title: '跳转异常，请重试',
        icon: 'none'
      })
    }
  },

  onLocationTap() {
    // 确保没有loading状态残留
    wx.hideLoading()
    this.showCitySelectionModal()
  },

  onDebugLocation() {
    wx.showToast({
      title: '调试功能开发中',
      icon: 'none'
    })
  },

  getLocation() {
    const userSelectedCity = wx.getStorageSync('userSelectedCity')
    const citySelectionType = wx.getStorageSync('citySelectionType')
    
    if (userSelectedCity && citySelectionType === 'manual') {
      this.setData({
        currentLocation: userSelectedCity
      })
      console.log('使用用户手动选择的城市:', userSelectedCity)
      return
    }

    // 检测是否在微信开发者工具中
    const systemInfo = wx.getSystemInfoSync()
    if (systemInfo.platform === 'devtools') {
      console.log('检测到微信开发者工具环境，跳过真实定位')
      this.setData({
        currentLocation: '北京'
      })
      wx.setStorageSync('userSelectedCity', '北京')
      wx.setStorageSync('citySelectionType', 'devtools')
      return
    }

    // 设置默认城市，避免一直显示"定位中"
    this.setData({
      currentLocation: '北京'
    })

    wx.getSetting({
      success: (res) => {
        if (res.authSetting['scope.userLocation'] === false) {
          this.showLocationAuthModal()
        } else if (res.authSetting['scope.userLocation'] === true) {
          this.startLocation()
        } else {
          this.startLocation()
        }
      },
      fail: () => {
        this.startLocation()
      }
    })
  },

  startLocation() {
    // 确保没有其他loading状态
    wx.hideLoading()
    
    wx.showLoading({
      title: '定位中...'
    })

    // 设置定位超时（10秒）
    const locationTimeout = setTimeout(() => {
      console.log('定位超时，允许手动选择城市')
      wx.hideLoading()
      this.setData({
        currentLocation: '北京'
      })
      wx.showModal({
        title: '定位超时',
        content: '定位服务调用超时10秒，您可以手动选择当前所在城市。',
        confirmText: '手动选择',
        cancelText: '使用默认城市',
        success: (res) => {
          if (res.confirm) {
            this.showCitySelectionModal()
          } else {
            // 用户选择使用默认城市
            this.setData({
              currentLocation: '北京'
            })
            wx.setStorageSync('userSelectedCity', '北京')
            wx.setStorageSync('citySelectionType', 'default')
          }
        }
      })
    }, 10000)

    wx.getLocation({
      type: 'gcj02',
      success: (res) => {
        console.log('定位成功:', res)
        clearTimeout(locationTimeout) // 清除超时定时器
        wx.setStorageSync('citySelectionType', 'auto')
        this.reverseGeocode(res.latitude, res.longitude)
      },
      fail: (err) => {
        console.error('定位失败:', err)
        clearTimeout(locationTimeout) // 清除超时定时器
        wx.hideLoading()
        
        if (err.errMsg.includes('auth deny')) {
          this.showLocationAuthModal()
        } else {
          this.handleLocationFail()
        }
      }
    })
  },

  reverseGeocode(latitude, longitude) {
    // 确保loading状态被清理
    wx.hideLoading()
    
    // 由于微信小程序没有内置的逆地理编码API，我们提供两个选择：
    // 1. 根据经纬度估算城市（简单但不够准确）
    // 2. 让用户手动选择城市（推荐）
    
    // 估算城市（仅作参考）
    let estimatedCity = '北京'
    if (latitude > 30 && latitude < 40 && longitude > 110 && longitude < 130) {
      estimatedCity = '北京'
    } else if (latitude > 30 && latitude < 35 && longitude > 120 && longitude < 125) {
      estimatedCity = '上海'
    } else if (latitude > 22 && latitude < 25 && longitude > 113 && longitude < 115) {
      estimatedCity = '广州'
    } else if (latitude > 22 && latitude < 25 && longitude > 113 && longitude < 115) {
      estimatedCity = '深圳'
    }
    
    // 设置弹窗超时（15秒），避免用户不操作时一直阻塞
    const modalTimeout = setTimeout(() => {
      console.log('城市选择弹窗超时，自动设置默认城市')
      this.setData({
        currentLocation: estimatedCity
      })
      wx.setStorageSync('userSelectedCity', estimatedCity)
      wx.setStorageSync('citySelectionType', 'auto')
      wx.showToast({
        title: `已自动定位到${estimatedCity}`,
        icon: 'success',
        duration: 2000
      })
    }, 15000)
    
    // 显示城市选择弹窗，让用户确认或选择
    wx.showModal({
      title: '定位成功',
      content: `根据您的位置，我们估算您在${estimatedCity}附近。请确认或选择您的实际城市。`,
      confirmText: '确认',
      cancelText: '选择其他城市',
      success: (res) => {
        clearTimeout(modalTimeout) // 清除超时定时器
        if (res.confirm) {
          // 用户确认估算的城市
          this.setData({
            currentLocation: estimatedCity
          })
          wx.setStorageSync('userSelectedCity', estimatedCity)
          wx.setStorageSync('citySelectionType', 'auto')
          
          wx.showToast({
            title: `已定位到${estimatedCity}`,
            icon: 'success',
            duration: 2000
          })
        } else {
          // 用户选择手动选择城市
          this.showCitySelectionModal()
        }
      },
      fail: () => {
        clearTimeout(modalTimeout) // 清除超时定时器
        // 如果showModal失败，也要确保loading状态被清理
        wx.hideLoading()
      }
    })
  },

  showLocationAuthModal() {
    wx.showModal({
      title: '需要定位权限',
      content: '为了为您提供更好的服务，需要获取您的位置信息。您也可以手动选择当前所在城市。',
      confirmText: '去设置',
      cancelText: '手动选择',
      success: (res) => {
        if (res.confirm) {
          wx.openSetting({
            success: (settingRes) => {
              if (settingRes.authSetting['scope.userLocation']) {
                this.startLocation()
              } else {
                this.showCitySelectionModal()
              }
            }
          })
        } else {
          this.showCitySelectionModal()
        }
      },
      fail: () => {
        // 如果showModal失败，也要确保loading状态被清理
        wx.hideLoading()
      }
    })
  },

  showCitySelectionModal() {
    const cityGroups = {
      '热门城市': ['北京', '上海', '广州', '深圳', '杭州', '南京', '武汉', '成都', '重庆', '西安'],
      '华北地区': ['北京', '天津', '石家庄', '太原', '呼和浩特', '沈阳', '长春', '哈尔滨'],
      '华东地区': ['上海', '南京', '杭州', '苏州', '无锡', '宁波', '合肥', '福州', '厦门', '南昌', '济南', '青岛'],
      '华南地区': ['广州', '深圳', '佛山', '东莞', '南宁', '海口', '珠海', '中山'],
      '华中地区': ['武汉', '长沙', '郑州', '南昌'],
      '西南地区': ['成都', '重庆', '昆明', '贵阳', '拉萨'],
      '西北地区': ['西安', '兰州', '西宁', '银川', '乌鲁木齐']
    }

    const actionList = []
    const cityMap = {}
    
    cityGroups['热门城市'].forEach(city => {
      if (!cityMap[city]) {
        cityMap[city] = city
        actionList.push(`🔥 ${city}`)
      }
    })
    
    actionList.push('─────────')
    
    Object.keys(cityGroups).forEach(region => {
      if (region !== '热门城市') {
        cityGroups[region].forEach(city => {
          if (!cityMap[city]) {
            cityMap[city] = city
            actionList.push(`${city}`)
          }
        })
      }
    })

    actionList.unshift('📍 重新定位')
    
    // 设置城市选择弹窗超时（20秒），避免用户不操作时一直阻塞
    const actionSheetTimeout = setTimeout(() => {
      console.log('城市选择弹窗超时，自动设置默认城市')
      this.setData({
        currentLocation: '北京'
      })
      wx.setStorageSync('userSelectedCity', '北京')
      wx.setStorageSync('citySelectionType', 'default')
      wx.showToast({
        title: '已自动设置默认城市：北京',
        icon: 'success',
        duration: 2000
      })
    }, 20000)
    
    wx.showActionSheet({
      itemList: actionList,
      success: (res) => {
        clearTimeout(actionSheetTimeout) // 清除超时定时器
        if (res.tapIndex === 0) {
          this.retryLocation()
        } else if (actionList[res.tapIndex] === '─────────') {
          return
        } else {
          const selectedCity = actionList[res.tapIndex].replace('🔥 ', '')
          this.selectCity(selectedCity)
        }
      },
      fail: () => {
        clearTimeout(actionSheetTimeout) // 清除超时定时器
        console.log('用户取消选择城市')
        // 用户取消选择城市时，确保loading状态被清理
        wx.hideLoading()
      }
    })
  },

  retryLocation() {
    // 确保没有其他loading状态
    wx.hideLoading()
    
    wx.showLoading({
      title: '重新定位中...'
    })

    wx.removeStorageSync('userSelectedCity')
    this.startLocation()
  },

  selectCity(cityName) {
    this.setData({
      currentLocation: cityName
    })
    console.log('用户选择城市:', cityName)
    
    wx.setStorageSync('userSelectedCity', cityName)
    wx.setStorageSync('citySelectionType', 'manual')
    
    wx.showToast({
      title: `已选择${cityName}`,
      icon: 'success',
      duration: 2000
    })
  },

  handleLocationFail() {
    wx.showModal({
      title: '定位失败',
      content: '无法获取您的位置信息，可能是GPS信号较弱或网络问题。您可以手动选择当前所在城市。',
      confirmText: '重新定位',
      cancelText: '手动选择',
      success: (res) => {
        if (res.confirm) {
          this.startLocation()
        } else {
          this.showCitySelectionModal()
        }
      },
      fail: () => {
        // 如果showModal失败，也要确保loading状态被清理
        wx.hideLoading()
      }
    })
  },

  onCallService() {
    wx.makePhoneCall({
      phoneNumber: '13522113924',
      success: () => {
        console.log('拨打电话成功')
      },
      fail: (err) => {
        console.error('拨打电话失败:', err)
        wx.showToast({
          title: '拨打电话失败',
          icon: 'none'
        })
      }
    })
  },

  async loadServerHeroImages() {
    try {
      const imageMap = await fetchServerImages([this.data.mainBannerSource, this.data.videoPosterSource])
      this.setData({
        mainBannerImage: imageMap[this.data.mainBannerSource] || '',
        videoPoster: imageMap[this.data.videoPosterSource] || '',
      })
    } catch (error) {
      console.error('首页头图加载失败:', error)
    }
  },

  showVideoPoster() {
    this.setData({
      videoUrl: '',
      fallbackMode: true,
      videoLoadError: false
    })
  },

  // 加载视频
  async loadVideo() {
    if (!this.data.videoSource) {
      this.showVideoPoster()
      return
    }

    try {
      const localVideoPath = await fetchServerVideoToLocal(this.data.videoSource)
      this.directVideoPlayback(localVideoPath || this.data.videoSource)
    } catch (error) {
      console.error('首页视频下载到本地失败，回退远程地址:', error)
      this.directVideoPlayback(this.data.videoSource)
    }
  },

  onHomeServiceImageError(e) {
    const { index } = e.currentTarget.dataset
    const currentServices = [...this.data.currentServices]
    const fallbackImage = this.data.mainBannerImage || this.data.videoPoster || ''

    if (currentServices[index]) {
      currentServices[index] = {
        ...currentServices[index],
        imageUrl: fallbackImage
      }

      this.setData({ currentServices })
    }
  },

  // 直接视频播放（使用自己的视频文件）
  directVideoPlayback(videoUrl) {
    if (videoUrl === this.data.videoSource) {
      this.videoRetriedWithRemoteSource = true
    } else {
      this.videoRetriedWithRemoteSource = false
    }

    this.setData({
      videoUrl: videoUrl,
      videoLoadError: false,
      fallbackMode: false
    })
    
    this.videoValidationTimeout = setTimeout(() => {
      if (this.data.currentSwiperIndex === 1 && !this.data.videoLoadError && !this.data.fallbackMode && !this.data.isVideoPlaying) {
        this.enableFallbackMode()
      }
    }, 15000)
    
    this.videoLoadStartTime = Date.now()
    this.videoLoadCheckInterval = setInterval(() => {
      Date.now() - this.videoLoadStartTime
    }, 2000)
    
    setTimeout(() => {
      if (this.videoContext && !this.data.fallbackMode && this.data.currentSwiperIndex === 1) {
        try {
          this.videoContext.play()
        } catch (error) {
          console.error('视频播放命令失败:', error)
        }
      }
    }, 3000)
  },

  // 重试加载视频
  async retryLoadVideo() {
    this.loadVideo()
  },

  // 视频播放事件
  onVideoPlay() {
    this.setData({
      fallbackMode: false,
      videoLoadError: false,
      isVideoPlaying: true
    })
    if (this.videoValidationTimeout) {
      clearTimeout(this.videoValidationTimeout)
    }
  },

  // 视频暂停事件
  onVideoPause() {
    this.setData({
      isVideoPlaying: false
    })
  },

  // 视频播放结束事件
  onVideoEnded() {
    this.setData({
      isVideoPlaying: false
    })
  },

  // 轮播切换事件
  onSwiperChange(e) {
    const currentIndex = e.detail.current
    this.setData({
      currentSwiperIndex: currentIndex
    })
    
    if (currentIndex === 1) {
      this.setData({
        isVideoPlaying: false
      })

      if (this.data.fallbackMode) {
        return
      }

      if (!this.data.videoUrl && this.data.videoSource) {
        this.loadVideo()
        return
      }

      if (this.videoContext && this.data.videoUrl) {
        setTimeout(() => {
          if (this.data.currentSwiperIndex === 1 && !this.data.fallbackMode) {
            this.videoContext.play()
          }
        }, 300)
      }
    }
  },

  // 视频错误事件
  onVideoError(e) {
    console.error('视频播放错误:', e.detail)
    console.error('错误详情:', JSON.stringify(e.detail))

    const currentVideoUrl = this.data.videoUrl || ''
    const remoteVideoUrl = this.data.videoSource || ''
    const isLocalVideo = !!currentVideoUrl && currentVideoUrl !== remoteVideoUrl

    if (isLocalVideo && remoteVideoUrl && !this.videoRetriedWithRemoteSource) {
      console.warn('本地视频播放失败，尝试切换到服务端原始视频地址重试:', currentVideoUrl)
      this.videoRetriedWithRemoteSource = true
      this.directVideoPlayback(remoteVideoUrl)
      return
    }
    
    this.setData({
      videoLoadError: true
    })
    
    if (!this.data.fallbackMode) {
      this.enableFallbackMode()
      wx.showToast({
        title: '视频加载失败，已切换到图片模式',
        icon: 'none',
        duration: 2000
      })
    }
  },

  // 视频开始加载
  onVideoLoadStart() {
    if (this.videoValidationTimeout) {
      clearTimeout(this.videoValidationTimeout)
    }
  },

  // 视频数据加载完成
  onVideoLoadedData() {
    this.setData({
      videoLoadError: false
    })
    if (this.videoValidationTimeout) {
      clearTimeout(this.videoValidationTimeout)
    }
  },

  // 视频等待中
  onVideoWaiting() {},

  // 视频可以播放
  onVideoCanPlay() {
    this.setData({
      videoLoadError: false
    })
    if (this.videoValidationTimeout) {
      clearTimeout(this.videoValidationTimeout)
    }
  },

  // 进入画中画模式
  onVideoEnterPictureInPicture() {},

  // 退出画中画模式
  onVideoLeavePictureInPicture() {},


  // 启用降级模式（只显示封面图）
  enableFallbackMode() {
    this.setData({
      fallbackMode: true,
      videoUrl: '',
      videoLoadError: true
    })
  },

  // 禁用降级模式，重新尝试加载视频
  disableFallbackMode() {
    this.loadVideo()
  },




  // 播放视频
  playVideo() {
    if (this.videoContext) {
      this.videoContext.play()
    }
  },

  // 暂停视频
  pauseVideo() {
    if (this.videoContext) {
      this.videoContext.pause()
    }
  },

  // 停止视频
  stopVideo() {
    if (this.videoContext) {
      this.videoContext.stop()
    }
  },

  // 跳转到指定时间
  seekVideo(time) {
    if (this.videoContext) {
      this.videoContext.seek(time)
    }
  },

  onUrgentService() {
    wx.showToast({
      title: '跳转到紧急服务页面',
      icon: 'none'
    })
  },

  async initHome() {
    try {
      this.setData({ loading: true, networkError: false })

      await this.loadServiceCategories()

      const currentCategory = this.data.currentServiceCategory

      await this.fetchServicesByCategory(currentCategory)
      
      this.setData({ loading: false })

      this.preloadHomeImages()
      
    } catch (error) {
      console.error('首页数据加载失败:', error)
      
      this.setData({
        loading: false,
        networkError: true
      })
      
      wx.showToast({
        title: '数据加载失败，请重试',
        icon: 'none',
        duration: 3000
      })
    }
  },

  // 预加载关键图片
  async preloadCriticalImages() {
    try {
      await preloadCriticalImages()
    } catch (error) {
      console.error('关键图片预加载失败:', error)
    }
  },

  // 预加载首页服务图片
  async preloadHomeImages() {
    try {
      const homeData = {
        services: this.data.currentServices
      }

      await preloadHomeImages(homeData)
    } catch (error) {
      console.error('首页图片预加载失败:', error)
    }
  },

  // 分享给好友
  onShareAppMessage(res) {
    if (res.from === 'button') {
      return {
        title: res.target.dataset.title || '安语颐年护理陪诊',
        desc: res.target.dataset.desc || '专业护理陪诊服务，让您就医更安心',
        path: res.target.dataset.path || 'pages/index/index',
        imageUrl: res.target.dataset.imageUrl || ''
      }
    }

    return {
      title: '安语颐年护理陪诊',
      desc: '专业护理陪诊服务，让您就医更安心',
      path: 'pages/index/index'
    }
  },

  // 分享到朋友圈
  onShareTimeline() {
    return {
      title: '安语颐年护理陪诊 - 专业护理陪诊服务',
      query: '',
      imageUrl: ''
    }
  }
})
