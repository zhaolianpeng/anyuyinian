// 图片处理工具函数

const { cos } = require('../config')

function buildStaticUrl(assetPath) {
  if (!assetPath) {
    return ''
  }

  if (assetPath.startsWith('/')) {
    return `${cos.bucketDomain}${assetPath}`
  }

  return `${cos.bucketDomain}/${assetPath}`
}

function convertCloudPathToStaticUrl(cloudUrl) {
  if (!cloudUrl) {
    return ''
  }

  const normalized = cloudUrl.replace('@cloud://', '').replace('cloud://', '')
  const pathParts = normalized.split('/')

  if (pathParts.length < 2) {
    return cloudUrl
  }

  return buildStaticUrl(pathParts.slice(1).join('/'))
}

/**
 * 处理图片URL，确保能正常显示
 * @param {string} imageUrl - 原始图片URL
 * @returns {string} - 处理后的图片URL
 */
function processImageUrl(imageUrl) {
  if (!imageUrl) {
    return ''
  }

  if (imageUrl.startsWith('@cloud://') || imageUrl.startsWith('cloud://')) {
    return convertCloudPathToStaticUrl(imageUrl)
  }
  
  // 如果是完整的HTTPS URL，直接返回
  if (imageUrl.startsWith('http://') || imageUrl.startsWith('https://')) {
    return imageUrl
  }
  
  // 如果是相对路径，添加基础URL
  if (imageUrl.startsWith('/')) {
    return buildStaticUrl(imageUrl)
  }

  if (imageUrl.includes('/')) {
    return buildStaticUrl(imageUrl)
  }
  
  // 如果是其他格式，直接返回
  return buildStaticUrl(`${cos.imagePrefix}${imageUrl}`)
}

/**
 * 批量处理图片URL数组
 * @param {Array} items - 包含图片URL的对象数组
 * @param {string} imageField - 图片字段名，默认为'imageUrl'
 * @returns {Array} - 处理后的数组
 */
function processImageUrls(items, imageField = 'imageUrl') {
  if (!Array.isArray(items)) {
    return items
  }
  
  return items.map(item => {
    if (item && item[imageField]) {
      return {
        ...item,
        [imageField]: processImageUrl(item[imageField])
      }
    }
    return item
  })
}

/**
 * 处理首页数据中的图片URL
 * @param {Object} homeData - 首页数据
 * @returns {Object} - 处理后的首页数据
 */
function processHomeDataImages(homeData) {
  if (!homeData || !homeData.data) {
    return homeData
  }
  
  const processedData = { ...homeData }
  
  // 处理轮播图
  if (processedData.data.banners) {
    processedData.data.banners = processImageUrls(processedData.data.banners, 'imageUrl')
  }
  
  // 处理导航菜单图标
  if (processedData.data.navigations) {
    processedData.data.navigations = processImageUrls(processedData.data.navigations, 'icon')
  }
  
  // 处理服务列表图片
  if (processedData.data.services) {
    processedData.data.services = processImageUrls(processedData.data.services, 'imageUrl')
    processedData.data.services = processImageUrls(processedData.data.services, 'icon')
  }
  
  // 处理医院列表logo
  if (processedData.data.hospitals) {
    processedData.data.hospitals = processImageUrls(processedData.data.hospitals, 'logo')
  }
  
  return processedData
}

/**
 * 检查图片是否加载成功
 * @param {string} imageUrl - 图片URL
 * @returns {Promise<boolean>} - 是否加载成功
 */
function checkImageLoad(imageUrl) {
  return new Promise((resolve) => {
    if (!imageUrl) {
      resolve(false)
      return
    }
    
    const img = new Image()
    img.onload = () => resolve(true)
    img.onerror = () => resolve(false)
    img.src = imageUrl
  })
}

/**
 * 获取图片加载状态
 * @param {string} imageUrl - 图片URL
 * @param {Function} onSuccess - 成功回调
 * @param {Function} onError - 错误回调
 */
function loadImage(imageUrl, onSuccess, onError) {
  if (!imageUrl) {
    onError && onError('图片URL为空')
    return
  }
  
  const processedUrl = processImageUrl(imageUrl)
  
  wx.getImageInfo({
    src: processedUrl,
    success: (res) => {
      onSuccess && onSuccess(res)
    },
    fail: (err) => {
      console.error('图片加载失败:', processedUrl, err)
      onError && onError(err)
    }
  })
}

/**
 * 处理医院图片URL
 * @param {string} logoUrl 后端返回的logo URL
 * @param {string} defaultUrl 默认图片URL
 * @returns {string} 处理后的图片URL
 */
function processHospitalLogo(logoUrl, defaultUrl = '/images/hospital-default.jpg') {
  if (!logoUrl) {
    return defaultUrl
  }

  return processImageUrl(logoUrl)
}

/**
 * 处理医院数据中的图片
 * @param {Array} hospitals 医院数据数组
 * @returns {Array} 处理后的医院数据
 */
function processHospitalImages(hospitals) {
  if (!hospitals || !Array.isArray(hospitals)) {
    return []
  }
  
  return hospitals.map(hospital => {
    if (hospital && typeof hospital === 'object') {
      return {
        ...hospital,
        logo: processHospitalLogo(hospital.logo)
      }
    }
    return hospital
  })
}

/**
 * 处理单个医院图片
 * @param {Object} hospital 医院数据对象
 * @returns {Object} 处理后的医院数据
 */
function processSingleHospitalImage(hospital) {
  if (!hospital || typeof hospital !== 'object') {
    return hospital
  }
  
  return {
    ...hospital,
    logo: processHospitalLogo(hospital.logo)
  }
}

module.exports = {
  buildStaticUrl,
  convertCloudPathToStaticUrl,
  processImageUrl,
  processImageUrls,
  processHomeDataImages,
  checkImageLoad,
  loadImage,
  processHospitalLogo,
  processHospitalImages,
  processSingleHospitalImage
} 