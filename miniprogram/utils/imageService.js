/**
 * 图片服务工具类
 * 统一将历史 fileID / cloud 路径转换为服务器静态资源地址
 */

const { processImageUrl } = require('./image')

/**
 * 将 fileID 或历史路径转换为可直接访问的静态资源地址
 * @param {string|Array} fileID 对象存储 fileID 或路径
 * @returns {Promise} 返回直接访问地址
 */
async function getTempFile(fileID) {
  try {
    const fileList = (typeof fileID === 'string' ? [fileID] : fileID).map(item => {
      return {
        fileID: item,
        status: 0,
        tempFileURL: processImageUrl(item)
      }
    })

    const result = { fileList }
    
    console.log('转换静态图片URL成功:', result)
    return result
  } catch (error) {
    console.error('转换静态图片URL失败:', error)
    throw error
  }
}

/**
 * 批量获取服务图片的临时访问地址
 * @param {Array} services 服务列表，每个服务包含 imageCosId 字段
 * @param {number} time 有效时间，单位秒，默认24小时
 * @returns {Promise<Array>} 返回处理后的服务列表，包含临时访问地址
 */
async function getServiceImages(services, time = 86400) {
  if (!services || services.length === 0) {
    return services
  }
  
  try {
    // 收集所有需要转换的 imageCosId
    const imageCosIds = services
      .map(service => service.imageCosId)
      .filter(cosId => cosId && cosId.trim() !== '')
    
    if (imageCosIds.length === 0) {
      console.log('没有需要处理的 imageCosId，直接使用现有图片地址')
      return services.map(service => ({
        ...service,
        imageTempUrl: processImageUrl(service.imageUrl)
      }))
    }
    
    // 批量转换为静态URL
    const tempFileResult = await getTempFile(imageCosIds, time)
    
    // 创建 imageCosId 到静态URL 的映射
    const cosIdToUrlMap = {}
    tempFileResult.fileList.forEach(item => {
      if (item.status === 0) {
        cosIdToUrlMap[item.fileID] = item.tempFileURL
      } else {
        console.warn('转换静态URL失败:', item.fileID, item.errMsg)
      }
    })
    
    // 为每个服务添加临时访问地址
    return services.map(service => {
      const tempUrl = service.imageCosId ? cosIdToUrlMap[service.imageCosId] : null
      return {
        ...service,
        imageTempUrl: tempUrl || processImageUrl(service.imageUrl),
        imageCosId: service.imageCosId
      }
    })
    
  } catch (error) {
    console.error('批量获取服务图片失败:', error)
    return services.map(service => ({
      ...service,
      imageTempUrl: processImageUrl(service.imageUrl)
    }))
  }
}

/**
 * 获取单个服务图片的临时访问地址
 * @param {Object} service 服务对象，包含 imageCosId 字段
 * @param {number} time 有效时间，单位秒，默认24小时
 * @returns {Promise<Object>} 返回包含临时访问地址的服务对象
 */
async function getSingleServiceImage(service, time = 86400) {
  if (!service || !service.imageCosId) {
    return {
      ...service,
      imageTempUrl: processImageUrl((service && service.imageUrl) || '')
    }
  }
  
  try {
    const result = await getTempFile(service.imageCosId, time)
    const tempUrl = (result.fileList[0] && result.fileList[0].tempFileURL) || service.imageUrl
    
    return {
      ...service,
      imageTempUrl: tempUrl || processImageUrl(service.imageUrl)
    }
  } catch (error) {
    console.error('获取单个服务图片失败:', error)
    return {
      ...service,
      imageTempUrl: processImageUrl(service.imageUrl)
    }
  }
}

/**
 * 预加载服务图片
 * @param {Array} services 服务列表
 * @param {number} time 有效时间，单位秒，默认24小时
 * @returns {Promise<Array>} 返回预加载后的服务列表
 */
async function preloadServiceImages(services, time = 86400) {
  try {
    const processedServices = await getServiceImages(services, time)
    
    // 预加载图片到本地缓存
    const preloadPromises = processedServices.map(async (service) => {
      if (service.imageTempUrl) {
        try {
          await wx.getImageInfo({
            src: service.imageTempUrl
          })
          console.log('图片预加载成功:', service.imageTempUrl)
        } catch (error) {
          console.warn('图片预加载失败:', service.imageTempUrl, error)
        }
      }
      return service
    })
    
    await Promise.all(preloadPromises)
    return processedServices
    
  } catch (error) {
    console.error('预加载服务图片失败:', error)
    return services
  }
}

module.exports = {
  getTempFile,
  getServiceImages,
  getSingleServiceImage,
  preloadServiceImages
}
