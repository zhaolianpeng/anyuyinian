/**
 * 视频服务工具类
 * 统一将历史 fileID / cloud 路径转换为服务器静态资源地址
 */

const { processImageUrl } = require('./image')

/**
 * 将视频 fileID 或路径转换为可直接访问的静态资源地址
 * @param {string|Array} fileID 对象存储 fileID 或路径
 * @returns {Promise} 返回直接访问地址
 */
async function getTempVideoFile(fileID) {
  try {
    const fileList = (typeof fileID === 'string' ? [fileID] : fileID).map(item => {
      return {
        fileID: item,
        status: 0,
        tempFileURL: processImageUrl(item)
      }
    })

    const result = { fileList }
    
    console.log('转换静态视频URL成功:', result)
    return result
  } catch (error) {
    console.error('转换静态视频URL失败:', error)
    throw error
  }
}

/**
 * 获取单个视频的临时访问地址
 * @param {string} videoCosId 视频的COS ID
 * @param {number} time 有效时间，单位秒，默认24小时
 * @returns {Promise<string>} 返回临时访问地址
 */
async function getVideoUrl(videoCosId, time = 86400) {
  if (!videoCosId || videoCosId.trim() === '') {
    return null
  }
  
  try {
    const result = await getTempVideoFile(videoCosId, time)
    const tempUrl = result.fileList[0]?.tempFileURL || processImageUrl(videoCosId)
    
    if (!tempUrl) {
      console.warn('获取视频临时URL失败:', videoCosId)
    }
    
    return tempUrl
  } catch (error) {
    console.error('获取视频URL失败:', error)
    return null
  }
}

/**
 * 预加载视频
 * @param {string} videoUrl 视频URL
 * @returns {Promise<boolean>} 返回是否预加载成功
 */
async function preloadVideo(videoUrl) {
  if (!videoUrl) {
    return false
  }
  
  try {
    // 创建视频上下文进行预加载
    const videoContext = wx.createVideoContext('homeVideo')
    if (videoContext) {
      // 这里可以添加预加载逻辑
      console.log('视频预加载:', videoUrl)
      return true
    }
    return false
  } catch (error) {
    console.error('视频预加载失败:', error)
    return false
  }
}

module.exports = {
  getTempVideoFile,
  getVideoUrl,
  preloadVideo
}
