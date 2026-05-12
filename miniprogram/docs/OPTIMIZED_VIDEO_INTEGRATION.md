# 优化视频集成指南

## 概述

本指南介绍如何将FFmpeg优化后的视频集成到小程序中，解决视频播放兼容性问题。

## 1. 视频配置更新

### 1.1 更新视频路径
```javascript
data: {
  // 优化后的视频路径
  videoCosId: 'cloud://prod-5g94mx7a3d07e78c.7072-prod-5g94mx7a3d07e78c-1353115175/video/video.mp4',
  videoPosterCosId: 'cloud://prod-5g94mx7a3d07e78c.7072-prod-5g94mx7a3d07e78c-1353115175/static/Wechat-IMG36.jpg'
}
```

### 1.2 视频格式规格
- **容器格式**: MP4
- **视频编码**: H.264 (AVC) - Main Profile, Level 3.1
- **音频编码**: AAC - 128kbps, 44.1kHz, 立体声
- **分辨率**: 1280x720 (720p)
- **码率**: 1.2Mbps (视频) + 128kbps (音频)
- **优化**: `+faststart` 标志，优化网络播放

## 2. 代码实现

### 2.1 视频加载逻辑
```javascript
async loadVideo() {
  try {
    console.log('开始加载优化后的视频和封面...')
    
    // 加载优化后的视频
    const videoUrl = await getVideoUrl(this.data.videoCosId)
    
    // 加载封面图
    const posterUrl = await getVideoUrl(this.data.videoPosterCosId)
    
    this.setData({
      videoUrl: videoUrl || '',
      videoPoster: posterUrl || '',
      videoLoadError: false
    })
    
    if (videoUrl) {
      console.log('优化后的视频加载成功:', videoUrl)
      console.log('视频格式: 小程序兼容格式 (H.264 + AAC)')
      // 使用优化后的验证方法
      this.validateOptimizedVideoUrl(videoUrl)
    }
  } catch (error) {
    console.error('优化后的视频/封面加载失败:', error)
    this.setData({ videoLoadError: true })
  }
}
```

### 2.2 优化视频验证
```javascript
validateOptimizedVideoUrl(videoUrl) {
  if (!videoUrl) return
  
  console.log('开始验证优化后的视频URL:', videoUrl)
  
  // 给优化后的视频更多验证时间
  this.videoValidationTimeout = setTimeout(() => {
    if (!this.data.videoLoadError && !this.data.fallbackMode) {
      console.warn('优化后的视频URL验证超时，启用降级模式')
      this.enableFallbackMode()
    }
  }, 8000) // 8秒超时
  
  // 验证URL可访问性
  wx.request({
    url: videoUrl,
    method: 'HEAD',
    success: (res) => {
      console.log('优化后的视频URL验证通过，状态码:', res.statusCode)
      if (this.videoValidationTimeout) {
        clearTimeout(this.videoValidationTimeout)
      }
      // 给更多时间检测播放状态
      setTimeout(() => {
        if (!this.data.videoLoadError && !this.data.fallbackMode) {
          console.log('优化后的视频可能仍有播放问题，启用降级模式')
          this.enableFallbackMode()
        }
      }, 5000) // 5秒检测时间
    },
    fail: (error) => {
      console.log('优化后的视频URL验证失败:', error)
      this.enableFallbackMode()
    }
  })
}
```

## 3. 视频播放器配置

### 3.1 优化后的播放器
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
  style="width: 100%; height: 100%;">
</video>
```

### 3.2 关键配置说明
- `controls="{{true}}"`: 显示播放控制条
- `show-progress="{{true}}"`: 显示进度条
- `enable-progress-gesture="{{true}}"`: 支持手势控制
- `object-fit="cover"`: 视频填充模式

## 4. 测试和验证

### 4.1 功能测试
```javascript
// 测试视频加载
testOptimizedVideo() {
  console.log('测试优化后的视频...')
  this.loadVideo()
}

// 测试视频播放
testVideoPlayback() {
  if (this.videoContext) {
    this.videoContext.play()
    console.log('开始播放优化后的视频')
  }
}

// 检查视频状态
checkVideoStatus() {
  console.log('当前视频状态:', {
    videoUrl: this.data.videoUrl,
    videoPoster: this.data.videoPoster,
    fallbackMode: this.data.fallbackMode,
    videoLoadError: this.data.videoLoadError
  })
}
```

### 4.2 预期结果
1. **视频加载成功**: 控制台显示"优化后的视频加载成功"
2. **URL验证通过**: 状态码200
3. **视频正常播放**: 无`net::ERR_FAILED`错误
4. **降级模式不触发**: 视频能正常播放

## 5. 性能优化

### 5.1 加载优化
- 使用`+faststart`标志优化网络播放
- 合理的码率设置(1.2Mbps)
- 适当的文件大小控制

### 5.2 播放优化
- 静音播放避免自动播放限制
- 循环播放提供持续体验
- 进度条支持用户控制

## 6. 错误处理

### 6.1 降级机制
```javascript
// 如果优化后的视频仍有问题，启用降级模式
enableFallbackMode() {
  this.setData({
    fallbackMode: true,
    videoUrl: '',
    videoLoadError: true
  })
  console.log('启用降级模式，显示封面图片')
}
```

### 6.2 用户反馈
- 显示加载状态
- 提供重试功能
- 保持视觉一致性

## 7. 监控和调试

### 7.1 日志监控
```javascript
// 详细的加载日志
console.log('优化后的视频加载成功:', videoUrl)
console.log('视频格式: 小程序兼容格式 (H.264 + AAC)')

// 验证日志
console.log('优化后的视频URL验证通过，状态码:', res.statusCode)

// 播放状态日志
onVideoCanPlay() {
  console.log('优化后的视频可以播放！')
}
```

### 7.2 性能监控
- 监控加载时间
- 检查播放成功率
- 跟踪用户交互

## 8. 部署和更新

### 8.1 视频文件更新
1. 使用FFmpeg优化视频
2. 上传到云存储
3. 更新小程序配置
4. 测试播放功能

### 8.2 版本控制
- 保留原视频作为备份
- 记录优化参数
- 监控播放效果

## 9. 最佳实践

### 9.1 视频准备
- 使用推荐的编码参数
- 控制文件大小
- 测试兼容性

### 9.2 代码实现
- 实现完善的错误处理
- 提供降级方案
- 优化用户体验

### 9.3 测试验证
- 多设备测试
- 网络环境测试
- 用户反馈收集

## 10. 总结

通过使用FFmpeg优化视频并更新小程序配置，可以解决视频播放兼容性问题：

1. **视频格式优化**: H.264 + AAC，小程序兼容
2. **网络播放优化**: `+faststart`标志
3. **代码逻辑优化**: 针对优化视频的验证机制
4. **用户体验优化**: 完善的错误处理和降级机制

这样既解决了技术问题，又保证了用户体验的完整性。

## 11. 下一步

1. **测试播放效果**: 验证优化后的视频是否正常播放
2. **性能监控**: 监控加载和播放性能
3. **用户反馈**: 收集用户使用反馈
4. **持续优化**: 根据反馈进一步优化
