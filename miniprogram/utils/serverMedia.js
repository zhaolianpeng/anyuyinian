const { api } = require('./cloud-container-standard')

const SERVER_BASE_URL = 'https://api.succ.online/anyuyinian'
const mediaCache = new Map()
const videoCache = new Map()

function isHtmlResponse(payload) {
  return typeof payload === 'string' && payload.trim().toLowerCase().startsWith('<!doctype html')
}

function requestServerMediaImagesDirect(paths) {
  return new Promise((resolve, reject) => {
    wx.request({
      url: `${SERVER_BASE_URL}/api/media/images`,
      method: 'POST',
      data: { paths },
      header: {
        'content-type': 'application/json'
      },
      success: (response) => {
        if (response.statusCode !== 200) {
          reject(new Error(`媒体接口HTTP ${response.statusCode}`))
          return
        }

        if (isHtmlResponse(response.data)) {
          reject(new Error('媒体接口返回HTML而不是JSON'))
          return
        }

        resolve(response.data)
      },
      fail: (error) => {
        reject(error)
      }
    })
  })
}

async function requestMediaImages(paths) {
  try {
    const result = await requestServerMediaImagesDirect(paths)
    if (result && result.code === 0) {
      return result
    }
    throw new Error((result && result.errorMsg) || '直连媒体接口返回异常')
  } catch (directError) {
    console.warn('直连媒体接口失败，回退统一服务端 API:', directError)
    const fallbackResult = await api.mediaImages({ paths })

    if (isHtmlResponse(fallbackResult)) {
      throw new Error('统一服务端 API 返回HTML而不是JSON')
    }

    return fallbackResult
  }
}

function normalizeServerImagePath(input) {
  if (!input || typeof input !== 'string') {
    return ''
  }

  const trimmed = input.trim()
  if (!trimmed) {
    return ''
  }

  if (trimmed.startsWith(SERVER_BASE_URL)) {
    return trimmed.slice(SERVER_BASE_URL.length)
  }

  if (trimmed.startsWith('/static/') || trimmed.startsWith('/images/')) {
    return trimmed
  }

  return ''
}

function getImageExtension(mediaPath, mimeType = '') {
  const match = mediaPath.match(/\.([a-zA-Z0-9]+)(?:\?|$)/)
  if (match) {
    return match[1]
  }

  if (mimeType.includes('png')) {
    return 'png'
  }

  if (mimeType.includes('webp')) {
    return 'webp'
  }

  if (mimeType.includes('gif')) {
    return 'gif'
  }

  return 'jpg'
}

function ensureMediaDir() {
  const fileSystemManager = wx.getFileSystemManager()
  const dirPath = `${wx.env.USER_DATA_PATH}/server-media`

  return new Promise((resolve, reject) => {
    fileSystemManager.mkdir({
      dirPath,
      recursive: true,
      success: () => resolve(dirPath),
      fail: (error) => {
        if (error.errMsg && error.errMsg.includes('file already exists')) {
          resolve(dirPath)
          return
        }

        reject(error)
      }
    })
  })
}

function normalizeServerVideoPath(input) {
  if (!input || typeof input !== 'string') {
    return ''
  }

  const trimmed = input.trim()
  if (!trimmed) {
    return ''
  }

  if (trimmed.startsWith(SERVER_BASE_URL)) {
    return trimmed
  }

  if (trimmed.startsWith('/')) {
    return `${SERVER_BASE_URL}${trimmed}`
  }

  return ''
}

async function fetchServerVideoToLocal(videoUrl) {
  const normalizedVideoUrl = normalizeServerVideoPath(videoUrl)
  if (!normalizedVideoUrl) {
    return ''
  }

  if (videoCache.has(normalizedVideoUrl)) {
    return videoCache.get(normalizedVideoUrl)
  }

  const dirPath = await ensureMediaDir()
  const extension = getImageExtension(normalizedVideoUrl, 'video/mp4')
  const safeName = normalizedVideoUrl.replace(/[^a-zA-Z0-9]/g, '_')
  const targetPath = `${dirPath}/${safeName}.${extension}`

  return new Promise((resolve, reject) => {
    wx.downloadFile({
      url: normalizedVideoUrl,
      filePath: targetPath,
      success: (result) => {
        if (result.statusCode !== 200) {
          reject(new Error(`视频下载失败: HTTP ${result.statusCode}`))
          return
        }

        const resolvedPath = result.tempFilePath || result.filePath || targetPath
        videoCache.set(normalizedVideoUrl, resolvedPath)
        resolve(resolvedPath)
      },
      fail: reject
    })
  })
}

async function writeServerImageToLocal(mediaPath, base64Content, mimeType) {
  if (mediaCache.has(mediaPath)) {
    return mediaCache.get(mediaPath)
  }

  const dirPath = await ensureMediaDir()
  const extension = getImageExtension(mediaPath, mimeType)
  const safeName = mediaPath.replace(/[^a-zA-Z0-9]/g, '_')
  const filePath = `${dirPath}/${safeName}.${extension}`
  const fileSystemManager = wx.getFileSystemManager()
  const binaryData = wx.base64ToArrayBuffer(base64Content)

  return new Promise((resolve, reject) => {
    fileSystemManager.writeFile({
      filePath,
      data: binaryData,
      success: () => {
        wx.getImageInfo({
          src: filePath,
          success: (result) => {
            const resolvedPath = result.path || filePath
            mediaCache.set(mediaPath, resolvedPath)
            resolve(resolvedPath)
          },
          fail: reject
        })
      },
      fail: reject
    })
  })
}

async function fetchServerImages(paths) {
  const normalizedPaths = [...new Set((paths || []).map(path => normalizeServerImagePath(path)).filter(Boolean))]
  const pendingPaths = normalizedPaths.filter(path => !mediaCache.has(path))

  if (pendingPaths.length > 0) {
    const result = await requestMediaImages(pendingPaths)

    if (result.code !== 0) {
      throw new Error(result.errorMsg || '服务器图片获取失败')
    }

    const items = (result.data && result.data.items) || []
    await Promise.all(items.map(async (item) => {
      if (item.status !== 0 || !item.base64) {
        console.warn('服务器图片项不可用:', item)
        return
      }

      await writeServerImageToLocal(item.path, item.base64, item.mimeType || '')
    }))
  }

  return normalizedPaths.reduce((accumulator, path) => {
    if (mediaCache.has(path)) {
      accumulator[path] = mediaCache.get(path)
    }
    return accumulator
  }, {})
}

module.exports = {
  fetchServerImages,
  fetchServerVideoToLocal,
  normalizeServerImagePath,
  normalizeServerVideoPath,
}