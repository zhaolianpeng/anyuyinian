# 小程序对象存储视频播放最佳实践指南

## 概述

本指南介绍如何在小程序中使用腾讯云对象存储（COS）生成的临时视频地址来播放视频，包括最佳实践、常见问题和解决方案。

## 1. 视频播放器配置

### 1.1 基本配置
```xml
<video 
  id="homeVideo"
  class="banner-video"
  src="{{videoUrl}}"
  poster="{{videoPoster}}"
  controls="{{true}}"
  autoplay="{{false}}"
  loop="{{true}}"
  muted="{{true}}"
  show-play-btn="{{true}}"
  show-center-play-btn="{{true}}"
  show-fullscreen-btn="{{false}}"
  show-progress="{{true}}"
  show-loading="{{true}}"
  enable-progress-gesture="{{true}}"
  object-fit="cover"
  style="width: 100%; height: 100%;"
  bindplay="onVideoPlay"
  bindpause="onVideoPause"
  binderror="onVideoError"
  bindloadstart="onVideoLoadStart"
  bindloadeddata="onVideoLoadedData"
  bindwaiting="onVideoWaiting"
  bindcanplay="onVideoCanPlay">
</video>
```

### 1.2 关键属性说明
- `controls="{{true}}"`: 显示播放控制条
- `autoplay="{{false}}"`: 禁止自动播放（小程序限制）
- `loop="{{true}}"`: 循环播放
- `muted="{{true}}"`: 静音播放（避免自动播放限制）
- `object-fit="cover"`: 视频填充模式
- `show-progress="{{true}}"`: 显示进度条
- `enable-progress-gesture="{{true}}"`: 允许手势控制进度

## 2. 对象存储视频URL获取

### 2.1 获取临时视频地址
```javascript
// 获取视频临时URL
async getVideoUrl(cosId) {
  try {
    const result = await wx.cloud.callFunction({
      name: 'getTempFileURL',
      data: {
        fileList: [cosId]
      }
    })
    
    if (result.result && result.result.fileList && result.result.fileList.length > 0) {
      const fileInfo = result.result.fileList[0]
      if (fileInfo.status === 0) {
        return fileInfo.tempFileURL
      } else {
        throw new Error(`获取视频URL失败: ${fileInfo.errMsg}`)
      }
    } else {
      throw new Error('获取视频URL失败: 无返回数据')
    }
  } catch (error) {
    console.error('获取视频URL失败:', error)
    throw error
  }
}
```

### 2.2 获取封面图临时地址
```javascript
// 获取封面图临时URL
async getVideoPosterUrl(cosId) {
  try {
    const result = await wx.cloud.callFunction({
      name: 'getTempFileURL',
      data: {
        fileList: [cosId]
      }
    })
    
    if (result.result && result.result.fileList && result.result.fileList.length > 0) {
      const fileInfo = result.result.fileList[0]
      if (fileInfo.status === 0) {
        return fileInfo.tempFileURL
      } else {
        throw new Error(`获取封面图URL失败: ${fileInfo.errMsg}`)
      }
    } else {
      throw new Error('获取封面图URL失败: 无返回数据')
    }
  } catch (error) {
    console.error('获取封面图URL失败:', error)
    throw error
  }
}
```

## 3. 视频加载和验证

### 3.1 视频加载流程
```javascript
async loadVideo() {
  try {
    console.log('开始加载视频和封面...')
    
    // 先尝试加载视频
    const videoUrl = await getVideoUrl(this.data.videoCosId)
    
    // 尝试加载封面图，如果失败则使用空字符串
    let posterUrl = ''
    try {
      posterUrl = await getVideoUrl(this.data.videoPosterCosId)
    } catch (error) {
      console.warn('封面图加载失败，将不显示封面:', error)
    }
    
    this.setData({
      videoUrl: videoUrl || '',
      videoPoster: posterUrl || '',
      videoLoadError: false
    })
    
    if (videoUrl) {
      console.log('视频加载成功:', videoUrl)
      // 验证视频URL是否可访问
      this.validateVideoUrl(videoUrl)
    } else {
      console.warn('视频加载失败')
      this.setData({
        videoLoadError: true
      })
      // 如果视频加载失败，尝试重试
      setTimeout(() => {
        this.retryLoadVideo()
      }, 3000)
    }
  } catch (error) {
    console.error('视频/封面加载失败:', error)
    this.setData({
      videoLoadError: true
    })
    // 延迟重试
    setTimeout(() => {
      this.retryLoadVideo()
    }, 5000)
  }
}
```

### 3.2 视频URL验证
```javascript
validateVideoUrl(videoUrl) {
  if (!videoUrl) {
    console.log('视频URL为空，跳过验证')
    return
  }
  
  console.log('开始验证视频URL:', videoUrl)
  
  // 设置验证超时
  this.videoValidationTimeout = setTimeout(() => {
    if (!this.data.videoLoadError && !this.data.fallbackMode) {
      console.warn('视频URL验证超时，可能无法播放，启用降级模式')
      this.enableFallbackMode()
    }
  }, 5000) // 5秒超时
  
  // 使用wx.request检查URL是否可访问
  wx.request({
    url: videoUrl,
    method: 'HEAD',
    success: (res) => {
      console.log('视频URL验证通过，状态码:', res.statusCode)
      if (this.videoValidationTimeout) {
        clearTimeout(this.videoValidationTimeout)
      }
      // URL可访问，不需要启用降级模式
    },
    fail: (error) => {
      console.log('视频URL验证失败:', error)
      if (this.videoValidationTimeout) {
        clearTimeout(this.videoValidationTimeout)
      }
      // 延迟启用降级模式，避免立即切换
      setTimeout(() => {
        this.enableFallbackMode()
      }, 1000)
    }
  })
}
```

## 4. 错误处理和降级模式

### 4.1 视频播放错误处理
```javascript
onVideoError(e) {
  console.error('视频播放错误:', e.detail)
  this.setData({
    videoLoadError: true
  })
  
  // 延迟启用降级模式，避免立即切换
  setTimeout(() => {
    if (!this.data.fallbackMode) {
      this.enableFallbackMode()
    }
  }, 500)
}
```

### 4.2 降级模式实现
```javascript
enableFallbackMode() {
  console.log('启用降级模式')
  this.setData({
    fallbackMode: true,
    videoUrl: '',
    videoLoadError: true
  })
  
  // 清除验证超时定时器
  if (this.videoValidationTimeout) {
    clearTimeout(this.videoValidationTimeout)
  }
}
```

## 5. 样式配置

### 5.1 视频播放器样式
```css
.banner-video {
  width: 100%;
  height: 100%;
  object-fit: cover;
  background-color: #000;
}

.video-banner {
  position: relative;
  width: 100%;
  height: 400rpx;
  overflow: hidden;
  border-radius: 20rpx;
}
```

### 5.2 降级模式样式
```css
.fallback-banner {
  position: relative;
  width: 100%;
  height: 400rpx;
  overflow: hidden;
  border-radius: 20rpx;
}

.fallback-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.fallback-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(0,0,0,0.3), rgba(0,0,0,0.6));
  display: flex;
  align-items: center;
  justify-content: center;
}
```

## 6. 常见问题和解决方案

### 6.1 视频无法播放
**问题**: 视频URL获取成功但无法播放
**解决方案**:
1. 检查视频格式是否支持（推荐MP4）
2. 检查视频编码格式（推荐H.264）
3. 检查视频文件大小（建议小于50MB）
4. 检查网络连接状态

### 6.2 视频加载缓慢
**问题**: 视频加载时间过长
**解决方案**:
1. 优化视频文件大小
2. 使用CDN加速
3. 添加加载状态提示
4. 实现预加载机制

### 6.3 视频播放卡顿
**问题**: 视频播放过程中卡顿
**解决方案**:
1. 降低视频分辨率
2. 优化视频编码参数
3. 检查设备性能
4. 实现自适应码率

### 6.4 临时URL过期
**问题**: 临时URL过期导致视频无法播放
**解决方案**:
1. 实现URL刷新机制
2. 添加过期时间检查
3. 实现自动重新获取URL
4. 添加降级模式

## 7. 性能优化建议

### 7.1 视频文件优化
- 使用H.264编码
- 分辨率控制在720p以内
- 文件大小控制在50MB以内
- 使用适当的帧率（24-30fps）

### 7.2 加载优化
- 实现预加载机制
- 添加加载状态提示
- 实现渐进式加载
- 优化网络请求

### 7.3 用户体验优化
- 添加播放控制
- 实现全屏播放
- 添加播放进度显示
- 实现播放状态记忆

## 8. 测试和调试

### 8.1 功能测试
```javascript
// 测试视频加载
testVideoLoading() {
  console.log('测试视频加载...')
  this.loadVideo()
}

// 测试降级模式
testFallbackMode() {
  console.log('测试降级模式...')
  this.enableFallbackMode()
}

// 测试重试机制
testRetryMechanism() {
  console.log('测试重试机制...')
  this.retryLoadVideo()
}
```

### 8.2 调试工具
- 使用微信开发者工具
- 查看网络请求日志
- 检查视频播放状态
- 监控内存使用情况

## 9. 最佳实践总结

1. **视频格式**: 使用MP4格式，H.264编码
2. **文件大小**: 控制在50MB以内
3. **分辨率**: 720p以内，避免过高分辨率
4. **错误处理**: 实现完善的错误处理和降级机制
5. **用户体验**: 添加加载状态和播放控制
6. **性能优化**: 实现预加载和缓存机制
7. **测试验证**: 充分测试各种网络环境下的表现

## 10. 代码示例

### 10.1 完整的视频播放组件
```javascript
// 页面数据
data: {
  videoUrl: '',
  videoPoster: '',
  videoCosId: 'video/example.mp4',
  videoPosterCosId: 'images/poster.jpg',
  videoLoadError: false,
  fallbackMode: false,
  videoValidationTimeout: null
}

// 页面加载时初始化
onLoad() {
  this.loadVideo()
}

// 视频加载
async loadVideo() {
  // 实现视频加载逻辑
}

// 视频验证
validateVideoUrl(videoUrl) {
  // 实现视频URL验证逻辑
}

// 错误处理
onVideoError(e) {
  // 实现错误处理逻辑
}

// 降级模式
enableFallbackMode() {
  // 实现降级模式逻辑
}
```

这个指南提供了在小程序中使用对象存储视频的完整解决方案，包括配置、实现、优化和调试等各个方面。通过遵循这些最佳实践，可以确保视频播放功能的稳定性和用户体验。
