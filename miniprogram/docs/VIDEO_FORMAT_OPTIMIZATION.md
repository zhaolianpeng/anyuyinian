# 小程序视频格式优化指南

## 问题分析

从您的日志可以看出：
- ✅ 视频URL获取成功
- ✅ 视频URL验证通过（状态码200）
- ❌ 视频播放时出现 `net::ERR_FAILED` 错误

这是一个典型的视频格式兼容性问题。

## 问题原因

1. **视频编码格式不兼容**：小程序对视频编码有严格要求
2. **视频文件过大**：可能导致加载超时
3. **视频分辨率过高**：可能超出小程序限制
4. **视频容器格式问题**：某些MP4容器格式不被支持

## 解决方案

### 1. 视频格式要求

#### 推荐格式
- **容器格式**：MP4
- **视频编码**：H.264
- **音频编码**：AAC
- **分辨率**：720p (1280x720) 或以下
- **帧率**：24-30fps
- **文件大小**：建议小于20MB

#### 编码参数
```
视频编码：H.264 (AVC)
- Profile: Main
- Level: 3.1
- 码率: 1-2 Mbps

音频编码：AAC
- 采样率: 44.1kHz 或 48kHz
- 码率: 128kbps
- 声道: 立体声
```

### 2. 视频优化工具

#### 使用FFmpeg优化视频
```bash
# 基本优化命令
ffmpeg -i input.mp4 -c:v libx264 -profile:v main -level 3.1 -b:v 1.5M -c:a aac -b:a 128k -movflags +faststart output.mp4

# 针对小程序的优化
ffmpeg -i input.mp4 \
  -c:v libx264 \
  -profile:v main \
  -level 3.1 \
  -b:v 1.2M \
  -maxrate 1.5M \
  -bufsize 3M \
  -c:a aac \
  -b:a 128k \
  -ar 44100 \
  -ac 2 \
  -movflags +faststart \
  -preset fast \
  -crf 23 \
  output.mp4
```

#### 参数说明
- `-c:v libx264`: 使用H.264编码
- `-profile:v main`: 使用Main Profile
- `-level 3.1`: 使用Level 3.1
- `-b:v 1.2M`: 视频码率1.2Mbps
- `-c:a aac`: 使用AAC音频编码
- `-b:a 128k`: 音频码率128kbps
- `-movflags +faststart`: 优化网络播放
- `-preset fast`: 编码速度优化
- `-crf 23`: 质量平衡

### 3. 在线优化工具

#### 推荐工具
1. **腾讯云点播**：提供视频转码服务
2. **阿里云视频点播**：支持小程序格式优化
3. **HandBrake**：免费的视频转换工具
4. **CloudConvert**：在线视频转换

### 4. 小程序视频播放优化

#### 当前实现优化
```javascript
// 视频播放器配置
<video 
  src="{{videoUrl}}"
  poster="{{videoPoster}}"
  controls="{{false}}"
  autoplay="{{false}}"
  loop="{{true}}"
  muted="{{true}}"
  show-play-btn="{{true}}"
  show-center-play-btn="{{true}}"
  show-progress="{{false}}"
  show-loading="{{true}}"
  enable-progress-gesture="{{false}}"
  object-fit="cover">
</video>
```

#### 错误处理优化
```javascript
onVideoError(e) {
  console.error('视频播放错误:', e.detail)
  
  // 立即启用降级模式
  if (!this.data.fallbackMode) {
    this.enableFallbackMode()
  }
}
```

### 5. 测试和验证

#### 视频格式测试
1. **本地测试**：在微信开发者工具中测试
2. **真机测试**：在真实设备上测试
3. **网络测试**：在不同网络环境下测试
4. **兼容性测试**：在不同设备上测试

#### 测试步骤
1. 上传优化后的视频到对象存储
2. 获取临时URL
3. 在小程序中测试播放
4. 检查控制台错误信息
5. 验证降级模式是否正常工作

### 6. 降级模式优化

#### 当前降级模式
- 显示封面图片
- 提供重试功能
- 显示联系客服按钮
- 保持视觉一致性

#### 进一步优化
```javascript
// 添加视频格式检测
checkVideoFormat(videoUrl) {
  const video = document.createElement('video')
  video.src = videoUrl
  video.addEventListener('loadedmetadata', () => {
    console.log('视频信息:', {
      duration: video.duration,
      videoWidth: video.videoWidth,
      videoHeight: video.videoHeight,
      codec: video.canPlayType('video/mp4; codecs="avc1.42E01E"')
    })
  })
}
```

### 7. 监控和日志

#### 添加详细日志
```javascript
onVideoError(e) {
  console.error('视频播放错误详情:', {
    error: e.detail,
    videoUrl: this.data.videoUrl,
    userAgent: navigator.userAgent,
    timestamp: new Date().toISOString()
  })
}
```

#### 错误分类
1. **网络错误**：`net::ERR_FAILED`
2. **格式错误**：`MEDIA_ERR_SRC_NOT_SUPPORTED`
3. **解码错误**：`MEDIA_ERR_DECODE`
4. **超时错误**：`MEDIA_ERR_NETWORK`

### 8. 最佳实践

#### 视频准备
1. 使用推荐的编码格式
2. 控制文件大小
3. 优化分辨率
4. 测试兼容性

#### 代码实现
1. 实现完善的错误处理
2. 提供降级方案
3. 添加用户反馈
4. 监控播放状态

#### 用户体验
1. 显示加载状态
2. 提供重试功能
3. 保持视觉一致性
4. 提供替代方案

## 总结

视频播放失败的主要原因是格式兼容性问题。通过：

1. **优化视频格式**：使用H.264编码，控制文件大小
2. **完善错误处理**：立即启用降级模式
3. **提供替代方案**：封面图片 + 操作按钮
4. **监控和调试**：详细日志和错误分类

可以确保用户获得良好的体验，即使视频无法播放也能通过降级模式获得完整的功能。

## 下一步行动

1. 使用FFmpeg优化视频格式
2. 重新上传优化后的视频
3. 测试播放功能
4. 验证降级模式
5. 监控用户反馈
