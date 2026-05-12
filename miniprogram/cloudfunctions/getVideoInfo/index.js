// 云函数：获取视频信息
const cloud = require('wx-server-sdk')

cloud.init({
  env: cloud.DYNAMIC_CURRENT_ENV
})

exports.main = async (event, context) => {
  const { videoUrl } = event
  
  if (!videoUrl) {
    return {
      success: false,
      error: '视频URL不能为空'
    }
  }
  
  try {
    // 使用云函数环境发起请求，避免跨域问题
    const result = await cloud.callFunction({
      name: 'httpRequest',
      data: {
        url: videoUrl,
        method: 'HEAD'
      }
    })
    
    if (result.result && result.result.statusCode === 200) {
      const headers = result.result.headers || {}
      const contentType = headers['content-type'] || headers['Content-Type']
      const contentLength = headers['content-length'] || headers['Content-Length']
      
      return {
        success: true,
        data: {
          contentType: contentType,
          contentLength: contentLength,
          fileSizeMB: contentLength ? (parseInt(contentLength) / 1024 / 1024).toFixed(2) + ' MB' : '未知',
          lastModified: headers['last-modified'] || headers['Last-Modified'],
          etag: headers['etag'] || headers['ETag'],
          statusCode: result.result.statusCode
        }
      }
    } else {
      return {
        success: false,
        error: '视频URL验证失败',
        statusCode: result.result ? result.result.statusCode : '未知'
      }
    }
  } catch (error) {
    console.error('获取视频信息失败:', error)
    return {
      success: false,
      error: error.message || '获取视频信息失败'
    }
  }
}
