# 视频调试工具指南

## 问题诊断

从您的日志分析，优化后的视频仍然无法播放，可能的原因：

1. **视频格式问题**：虽然URL可访问，但视频编码可能仍不兼容
2. **网络问题**：视频文件过大或网络不稳定
3. **小程序限制**：某些视频参数仍不符合小程序要求
4. **缓存问题**：浏览器或小程序缓存了旧版本

## 调试工具

### 1. 视频格式检测
```javascript
// 在控制台执行以下代码检测视频格式
const checkVideoFormat = async (videoUrl) => {
  try {
    const result = await wx.request({
      url: videoUrl,
      method: 'HEAD'
    })
    
    console.log('视频响应头:', result.header)
    console.log('Content-Type:', result.header['Content-Type'])
    console.log('Content-Length:', result.header['Content-Length'])
    console.log('Accept-Ranges:', result.header['Accept-Ranges'])
    
    return result.header
  } catch (error) {
    console.error('检测失败:', error)
  }
}

// 使用示例
checkVideoFormat('https://7072-prod-5g94mx7a3d07e78c-1353115175.tcb.qcloud.la/video/video.mp4')
```

### 2. 视频播放测试
```javascript
// 测试视频播放能力
const testVideoPlayback = () => {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1]
  
  console.log('当前视频状态:', {
    videoUrl: currentPage.data.videoUrl,
    videoPoster: currentPage.data.videoPoster,
    fallbackMode: currentPage.data.fallbackMode,
    videoLoadError: currentPage.data.videoLoadError
  })
  
  // 强制禁用降级模式测试
  currentPage.setData({
    fallbackMode: false,
    videoLoadError: false
  })
  
  console.log('已禁用降级模式，测试视频播放')
}

// 使用示例
testVideoPlayback()
```

### 3. 视频URL验证
```javascript
// 验证视频URL的详细信息
const validateVideoUrl = async (videoUrl) => {
  console.log('开始验证视频URL:', videoUrl)
  
  try {
    // 检查URL格式
    const url = new URL(videoUrl)
    console.log('URL解析:', {
      protocol: url.protocol,
      hostname: url.hostname,
      pathname: url.pathname,
      search: url.search
    })
    
    // 检查文件扩展名
    const extension = url.pathname.split('.').pop()
    console.log('文件扩展名:', extension)
    
    // 发送HEAD请求检查
    const result = await new Promise((resolve, reject) => {
      wx.request({
        url: videoUrl,
        method: 'HEAD',
        success: resolve,
        fail: reject
      })
    })
    
    console.log('HTTP响应:', {
      statusCode: result.statusCode,
      header: result.header
    })
    
    return result
  } catch (error) {
    console.error('URL验证失败:', error)
  }
}

// 使用示例
validateVideoUrl('https://7072-prod-5g94mx7a3d07e78c-1353115175.tcb.qcloud.la/video/video.mp4')
```

## 解决方案

### 1. 重新优化视频

如果视频格式仍有问题，使用更严格的FFmpeg参数：

```bash
# 超严格的小程序兼容格式
ffmpeg -i input.mp4 \
  -c:v libx264 \
  -profile:v baseline \
  -level 3.0 \
  -b:v 800k \
  -maxrate 1M \
  -bufsize 2M \
  -c:a aac \
  -b:a 96k \
  -ar 44100 \
  -ac 2 \
  -movflags +faststart \
  -preset ultrafast \
  -crf 28 \
  -vf "scale=640:360:force_original_aspect_ratio=decrease,pad=640:360:(ow-iw)/2:(oh-ih)/2" \
  -t 30 \
  output_ultra_compatible.mp4
```

### 2. 测试不同格式

创建多个版本的视频进行测试：

```bash
# 版本1：超小文件
ffmpeg -i input.mp4 -c:v libx264 -profile:v baseline -level 3.0 -b:v 500k -c:a aac -b:a 64k -ar 22050 -ac 1 -movflags +faststart -preset ultrafast -crf 30 -vf "scale=480:270" -t 15 test_small.mp4

# 版本2：标准格式
ffmpeg -i input.mp4 -c:v libx264 -profile:v main -level 3.1 -b:v 1M -c:a aac -b:a 128k -ar 44100 -ac 2 -movflags +faststart -preset fast -crf 23 -vf "scale=720:405" -t 30 test_standard.mp4

# 版本3：高质量格式
ffmpeg -i input.mp4 -c:v libx264 -profile:v high -level 4.0 -b:v 2M -c:a aac -b:a 192k -ar 48000 -ac 2 -movflags +faststart -preset medium -crf 20 -vf "scale=1280:720" -t 60 test_high.mp4
```

### 3. 使用云存储的原始视频

如果优化后的视频仍有问题，可以尝试使用云存储的原始视频：

```javascript
// 临时测试原始视频
data: {
  videoCosId: 'cloud://prod-5g94mx7a3d07e78c.7072-prod-5g94mx7a3d07e78c-1353115175/video/光厂_5440331_照顾瘫痪老人按摩理发.mp4',
  // 其他配置...
}
```

### 4. 添加视频预加载

```javascript
// 添加视频预加载功能
preloadVideo() {
  if (this.data.videoUrl) {
    const video = wx.createVideoContext('homeVideo', this)
    // 预加载视频
    video.play()
    setTimeout(() => {
      video.pause()
    }, 1000)
  }
}
```

## 测试步骤

1. **清除缓存**：重新编译小程序
2. **检查日志**：查看详细的调试信息
3. **测试播放**：手动测试视频播放
4. **格式检测**：使用调试工具检测视频格式
5. **降级测试**：验证降级模式是否正常工作

## 预期结果

如果优化成功，您应该看到：
- `🎬 优化后的视频开始加载`
- `🎬 优化后的视频数据加载完成`
- `🎉 优化后的视频可以播放！`
- `🎬 优化后的视频开始播放！`

如果仍然失败，降级模式会正常工作，显示封面图片。

## 总结

通过详细的调试工具和多种解决方案，我们可以：
1. 诊断视频播放问题的根本原因
2. 测试不同格式的视频文件
3. 验证小程序视频播放的兼容性
4. 确保降级模式正常工作

请使用这些调试工具来诊断问题，并告诉我结果！
