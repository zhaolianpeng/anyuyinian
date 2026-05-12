const { baseURL, errorCode, env, dev, cos } = require('../config')
const { callContainer, api } = require('./cloud-container-standard')

// 模拟数据
const mockData = {
  '/home/init': {
    code: 0,
    data: {
      banners: [
        {
          id: 1,
          title: '专业护理服务',
          image: '/static/fuwu_1.jpeg',
          url: '/pages/service/list',
          sort: 1,
          status: 1
        },
        {
          id: 2,
          title: '医院陪诊套餐',
          image: '/static/fuwu_2.jpeg',
          url: '/pages/service/list?category=医院陪诊',
          sort: 2,
          status: 1
        }
      ],
      navigations: [
        {
          id: 1,
          title: '服务预约',
          icon: '/static/fuwuyuyue_logo.png',
          url: '/pages/service/list',
          sort: 1,
          status: 1
        },
        {
          id: 2,
          title: '我的订单',
          icon: '/static/wodedingdan_logo.png',
          url: '/pages/order/list',
          sort: 2,
          status: 1
        },
        {
          id: 3,
          title: '医院信息',
          icon: '/static/yiyuanxinxi_logo.png',
          url: '/pages/hospital/list',
          sort: 3,
          status: 1
        },
        {
          id: 4,
          title: '个人中心',
          icon: '/static/gerenzhongxin_logo.png',
          url: '/pages/user/profile',
          sort: 4,
          status: 1
        }
      ],
      services: [
        {
          id: 1,
          title: '上门护理服务',
          description: '专业护工上门提供护理服务',
          icon: '/images/service-nursing.jpg',  
          url: '/pages/service/list',
          sort: 1,
          status: 1
        },
        {
          id: 2,
          title: '专业陪诊服务',
          description: '专业陪诊师全程陪同就医，提供挂号、排队、取药等服务',
          icon: '/images/service-escort.jpg',  
          url: '/pages/service/list',
          sort: 2,
          status: 1
        }
      ],
      hospitals: [
        {
          id: 1,
          name: '深圳市人民医院',
          logo: '/images/hospital-1.jpg',  
          address: '深圳市罗湖区东门北路1017号',
          phone: '0755-25533018',
          description: '三级甲等综合医院',
          level: '三级甲等',
          type: '综合医院',
          longitude: 114.0579,
          latitude: 22.5431,
          sort: 1,
          status: 1
        },
        {
          id: 2,
          name: '深圳市第二人民医院',
          logo: '/images/hospital-2.jpg',  
          address: '深圳市福田区笋岗西路3002号',
          phone: '0755-83366388',
          description: '三级甲等综合医院',
          level: '三级甲等',
          type: '综合医院',
          longitude: 114.0579,
          latitude: 22.5431,
          sort: 2,
          status: 1
        }
      ]
    }
  },
  '/service/list': {
    code: 0,
    data: {
      list: [
        {
          id: 1,
          title: '专业陪诊套餐',
          description: '专业陪诊师全程陪同就医，提供挂号、排队、取药等服务',
          price: 299.00,
          originalPrice: 399.00,
          category: '医院陪诊',
          images: ['/images/service-escort.jpg'],  // 使用COS图片
          status: 1,
          sort: 1
        },
        {
          id: 2,
          title: '上门护理服务',
          description: '专业护士上门提供护理服务',
          price: 199.00,
          originalPrice: 299.00,
          category: '专业护理',
          images: ['/images/service-nursing.jpg'],  // 使用COS图片
          status: 1,
          sort: 2
        },
        {
          id: 3,
          title: '生活照护服务',
          description: '专业护工提供日常生活照料，包括饮食、清洁、陪伴等服务',
          price: 150.00,
          originalPrice: 200.00,
          category: '生活照护',
          images: ['/images/service-care.jpg'],  // 使用COS图片
          status: 1,
          sort: 3
        },
        {
          id: 4,
          title: '康复护理服务',
          description: '专业康复师提供术后康复、理疗、功能训练等服务',
          price: 250.00,
          originalPrice: 350.00,
          category: '专业护理',
          images: ['/images/service-rehab.jpg'],  // 使用COS图片
          status: 1,
          sort: 4
        }
      ],
      total: 2,
      page: 1,
      pageSize: 10,
      hasMore: false
    }
  },
  '/order/list': {
    code: 0,
    data: {
      list: [
        {
          id: 1,
          orderNo: '202401150001',
          userId: 1,
          serviceId: 1,
          serviceTitle: '专业陪诊套餐',
          amount: 299.00,
          status: 'pending_pay',
          statusText: '待支付',
          appointmentDate: '2024-01-15',
          appointmentTime: 'morning',
          createdAt: '2024-01-01T12:00:00Z'
        }
      ],
      total: 1,
      page: 1,
      pageSize: 10,
      hasMore: false
    }
  }
}

// 静态资源图片处理函数
const getCosImageUrl = (imagePath) => {
  if (!imagePath) return ''
  
  // 如果是完整URL，直接返回
  if (imagePath.startsWith('http')) {
    return imagePath
  }
  
  // 如果是云开发格式的链接，转换为当前服务器静态资源URL
  if (imagePath.startsWith('@cloud://') || imagePath.startsWith('cloud://')) {
    return convertCloudUrlToDirectUrl(imagePath)
  }
  
  // 如果是相对路径，拼接COS域名
  if (imagePath.startsWith('/')) {
    return cos.bucketDomain + imagePath
  }
  
  // 如果是文件名，添加前缀
  return cos.bucketDomain + cos.imagePrefix + imagePath
}

// 转换历史云开发URL为当前服务器静态资源URL
const convertCloudUrlToDirectUrl = (cloudUrl) => {
  const urlWithoutPrefix = cloudUrl.replace('@cloud://', '').replace('cloud://', '')
  const parts = urlWithoutPrefix.split('/')
  if (parts.length < 2) {
    console.warn('Invalid cloud URL format:', cloudUrl)
    return cloudUrl
  }

  const filePath = parts.slice(1).join('/')
  const directUrl = `${cos.bucketDomain}/${filePath}`
  
  console.log('Cloud URL converted:', cloudUrl, '->', directUrl)
  return directUrl
}

// 处理图片数组
const processImages = (images) => {
  if (!images || !Array.isArray(images)) return []
  return images.map(img => getCosImageUrl(img))
}

// 处理对象中的图片字段
const processImageFields = (obj) => {
  if (!obj) return obj
  
  const processed = { ...obj }
  
  // 处理常见图片字段
  const imageFields = ['image', 'icon', 'logo', 'avatar', 'banner']
  imageFields.forEach(field => {
    if (processed[field]) {
      processed[field] = getCosImageUrl(processed[field])
    }
  })
  
  // 处理图片数组
  if (processed.images) {
    processed.images = processImages(processed.images)
  }
  
  return processed
}

// 递归处理数据中的图片
const processDataImages = (data) => {
  if (!data) return data
  
  if (Array.isArray(data)) {
    return data.map(item => processDataImages(item))
  }
  
  if (typeof data === 'object') {
    const processed = processImageFields(data)
    Object.keys(processed).forEach(key => {
      if (typeof processed[key] === 'object') {
        processed[key] = processDataImages(processed[key])
      }
    })
    return processed
  }
  
  return data
}

// 模拟请求函数
const mockRequest = (url, method = 'GET', data = {}) => {
  return new Promise((resolve, reject) => {
    setTimeout(() => {
      const mockResponse = mockData[url]
      if (mockResponse) {
        // 处理响应数据中的图片
        const processedResponse = {
          ...mockResponse,
          data: processDataImages(mockResponse.data)
        }
        resolve(processedResponse)
      } else {
        reject(new Error('接口不存在'))
      }
    }, dev.mockDelay)
  })
}

// 统一请求函数
const request = (url, method = "GET", data = {}, options = {}) => {
  return new Promise((resolve, reject) => {
    // 开发环境使用模拟数据
    if (env === 'dev' && dev.useMockData) {
      return mockRequest(url, method, data)
        .then(resolve)
        .catch(reject)
    }

    // 使用云托管API调用
    return callContainer(url, method, data, options)
      .then(resolve)
      .catch(reject)
  })
}

// 网络状态检查
const checkNetworkStatus = () => {
  return new Promise((resolve) => {
    wx.getNetworkType({
      success: (res) => {
        console.log('网络类型:', res.networkType)
        resolve({
          isConnected: res.networkType !== 'none',
          networkType: res.networkType
        })
      },
      fail: () => {
        console.error('获取网络状态失败')
        resolve({
          isConnected: false,
          networkType: 'unknown'
        })
      }
    })
  })
}

// 错误处理
const handleError = (error) => {
  console.error('请求错误:', error)
  
  let message = '请求失败'
  if (error.message) {
    if (error.message.includes('timeout')) {
      message = '请求超时，请检查网络'
    } else if (error.message.includes('fail')) {
      message = '服务暂时不可用，请稍后重试'
    } else if (error.message.includes('service')) {
      message = '云托管服务异常，请稍后重试'
    } else if (error.message.includes('network')) {
      message = '网络连接异常，请检查网络'
    } else {
      message = error.message
    }
  }
  
  // 显示错误提示
  if (typeof wx !== 'undefined') {
    wx.showToast({
      title: message,
      icon: 'none',
      duration: 2000
    })
  }
  
  return error
}

// 咨询相关API
const consultationAPI = {
  // 创建咨询会话
  createConsultation: (data) => {
    return request('/api/consultation/create', 'POST', data)
  },

  // 获取咨询消息
  getConsultationMessages: (consultationId) => {
    return request('/api/consultation/messages', 'GET', { consultationId })
  },

  // 发送咨询消息
  sendConsultationMessage: (consultationId, content, senderType) => {
    return request('/api/consultation/send', 'POST', { consultationId, content, senderType })
  },

  // 获取咨询状态
  getConsultationStatus: (consultationId) => {
    return request('/api/consultation/status', 'GET', { consultationId })
  },

  // 获取咨询详情
  getConsultationDetail: (consultationId) => {
    return request('/api/consultation/detail', 'GET', { consultationId })
  },

  // 关闭咨询会话
  closeConsultation: (consultationId) => {
    return request('/api/consultation/close', 'POST', { consultationId })
  },

  // 获取活跃咨询列表
  getActiveConsultations: () => {
    return request('/api/consultation/active', 'GET')
  },

  // 获取未读通知
  getUnreadNotifications: () => {
    return request('/api/consultation/notifications', 'GET')
  },

  // 标记通知为已读
  markNotificationAsRead: (data) => {
    return request('/api/consultation/notification/read', 'POST', data)
  },

  // 获取咨询统计
  getConsultationStats: () => {
    return request('/api/consultation/stats', 'GET')
  }
}

module.exports = { 
  request,
  checkNetworkStatus,
  handleError,
  api, // 导出云托管API
  consultationAPI // 导出咨询API
}