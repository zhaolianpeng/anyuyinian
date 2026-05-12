const { api } = require('../../utils/cloud-container-standard')
const { processImageUrl } = require('../../utils/image')
const { preloadImages } = require('../../utils/imagePreloader')
const { getServiceImages } = require('../../utils/imageService')
const { fetchServerImages, normalizeServerImagePath } = require('../../utils/serverMedia')

Page({
  data: {
    services: [],
    loading: true,
    hasMore: true,
    page: 1,
    pageSize: 10,
    category: '',
    categories: [
      { name: '全部', value: '' }
    ]
  },

  onLoad(options) {
    const { category } = options
    if (category) {
      this.setData({ category })
    }
    this.loadCategories()
    this.loadServices()
  },

  // 加载分类列表
  async loadCategories() {
    try {
      const res = await api.serviceCategories()
      if (res.code === 0 && res.data && res.data.categories) {
        this.setData({
          categories: res.data.categories
        })
        console.log('分类加载成功:', res.data.categories)
      } else {
        console.error('分类加载失败:', res)
      }
    } catch (error) {
      console.error('加载分类失败:', error)
    }
  },

  // 加载服务列表
  async loadServices(isRefresh = false) {
    try {
      if (isRefresh) {
        this.setData({ 
          page: 1, 
          hasMore: true,
          services: []
        })
      }

      if (!this.data.hasMore) return

      this.setData({ loading: true })

      const params = {
        page: this.data.page,
        pageSize: this.data.pageSize
      }

      if (this.data.category) {
        params.category = this.data.category
      }

      const res = await api.serviceList(params)
      
      if (res.code === 0) {
        const newServices = res.data.list || []
        
        // 兼容历史 imageCosId，统一转换为静态资源地址
        const processedServices = await getServiceImages(newServices)
        const serviceImagePaths = processedServices.map(service =>
          normalizeServerImagePath(service.imageTempUrl || service.imageUrl || service.imageCosId)
        )
        const imageMap = await fetchServerImages(serviceImagePaths)
        
        // 优先使用本地化后的服务图片，避免渲染层直接请求远程静态图
        const finalServices = processedServices.map(service => ({
          ...service,
          imageUrl:
            imageMap[normalizeServerImagePath(service.imageTempUrl || service.imageUrl || service.imageCosId)] ||
            service.imageTempUrl ||
            processImageUrl(service.imageUrl)
        }))
        
        const services = isRefresh ? finalServices : [...this.data.services, ...finalServices]
        
        this.setData({
          services,
          hasMore: res.data.hasMore,
          page: this.data.page + 1,
          loading: false
        })
        
        // 预加载服务图片
        this.preloadServiceImages(finalServices)
      } else {
        this.setData({ loading: false })
      }
    } catch (error) {
      console.error('加载服务列表失败:', error)
      this.setData({ loading: false })
      wx.showToast({
        title: '加载失败，请重试',
        icon: 'none'
      })
    }
  },

  // 下拉刷新
  onPullDownRefresh() {
    this.loadServices(true).then(() => {
      wx.stopPullDownRefresh()
    })
  },

  // 上拉加载更多
  onReachBottom() {
    if (this.data.hasMore && !this.data.loading) {
      this.loadServices()
    }
  },

  // 选择分类
  onCategorySelect(e) {
    const { category } = e.currentTarget.dataset
    this.setData({ category })
    this.loadServices(true)
  },

  // 点击服务项
  onServiceTap(e) {
    const { serviceId } = e.currentTarget.dataset
    wx.navigateTo({
      url: `/pages/service/detail?id=${serviceId}`
    })
  },

  // 搜索服务
  onSearch(e) {
    const { value } = e.detail
    // 这里可以实现搜索功能
    console.log('搜索:', value)
  },

  // 分享
  onShareAppMessage() {
    return {
      title: '专业医疗服务，健康生活从这里开始',
      path: '/pages/service/list'
    }
  },

  // 预加载服务图片
  async preloadServiceImages(services) {
    if (!services || services.length === 0) return

    try {
      // 提取图片URL
      const imageUrls = services
        .map(service => service.imageUrl)
        .filter(url => url && url.trim() !== '')

      if (imageUrls.length === 0) return

      await preloadImages(
        imageUrls,
        (progress) => {
          console.log(`服务图片预加载进度: ${progress.progress}% (${progress.loaded}/${progress.total})`)
        },
        (result) => {
          console.log('服务图片预加载完成:', result)
        }
      )
    } catch (error) {
      console.error('服务图片预加载失败:', error)
    }
  }
}) 