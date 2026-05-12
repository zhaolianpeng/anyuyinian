# 视频调试命令指南

## 当前状态
- ✅ 文件大小：1.16 MB（已优化）
- ✅ 文件格式：video/mp4
- ❌ 8秒内未开始播放

## 调试命令

### 1. 检查视频文件信息
```javascript
// 在控制台执行
const pages = getCurrentPages()
const currentPage = pages[pages.length - 1]
currentPage.checkVideoInfo()
```

### 2. 检查视频元素状态
```javascript
// 检查视频元素是否正常渲染
const pages = getCurrentPages()
const currentPage = pages[pages.length - 1]
currentPage.checkVideoElement()
```

### 3. 手动测试视频播放
```javascript
// 清除所有状态，重新测试
const pages = getCurrentPages()
const currentPage = pages[pages.length - 1]
currentPage.testVideoPlayback()
```

### 4. 强制播放视频
```javascript
// 强制调用播放方法
const pages = getCurrentPages()
const currentPage = pages[pages.length - 1]
currentPage.forcePlayVideo()
```

### 5. 检查当前状态
```javascript
// 查看所有视频相关状态
const pages = getCurrentPages()
const currentPage = pages[pages.length - 1]
console.log('当前视频状态:', {
  videoUrl: currentPage.data.videoUrl,
  videoPoster: currentPage.data.videoPoster,
  fallbackMode: currentPage.data.fallbackMode,
  videoLoadError: currentPage.data.videoLoadError,
  videoContext: currentPage.videoContext ? '存在' : '不存在'
})
```

## 可能的问题和解决方案

### 问题1：视频元素未正确渲染
**症状**：视频URL正确但无法播放
**解决方案**：
```javascript
// 检查视频元素
currentPage.checkVideoElement()
```

### 问题2：视频上下文问题
**症状**：videoContext不存在
**解决方案**：
```javascript
// 重新创建视频上下文
currentPage.videoContext = wx.createVideoContext('homeVideo', currentPage)
```

### 问题3：视频格式兼容性
**症状**：文件格式正确但小程序不支持
**解决方案**：
```bash
# 使用更兼容的编码参数
ffmpeg -i input.mp4 \
  -c:v libx264 \
  -profile:v baseline \
  -level 3.0 \
  -b:v 300k \
  -maxrate 500k \
  -bufsize 600k \
  -c:a aac \
  -b:a 32k \
  -ar 22050 \
  -ac 1 \
  -movflags +faststart \
  -preset ultrafast \
  -crf 35 \
  -vf "scale=320:180" \
  -t 15 \
  video_ultra_compatible.mp4
```

### 问题4：网络加载问题
**症状**：文件小但加载慢
**解决方案**：
```javascript
// 检查网络状态
wx.getNetworkType({
  success: (res) => {
    console.log('网络类型:', res.networkType)
    console.log('网络状态:', res.isConnected)
  }
})
```

## 测试步骤

1. **检查文件信息**
```javascript
currentPage.checkVideoInfo()
```

2. **检查元素状态**
```javascript
currentPage.checkVideoElement()
```

3. **手动测试播放**
```javascript
currentPage.testVideoPlayback()
```

4. **强制播放**
```javascript
currentPage.forcePlayVideo()
```

5. **查看详细日志**
观察控制台输出，特别关注：
- 视频加载状态
- 元素渲染信息
- 播放事件触发

## 预期结果

如果一切正常，应该看到：
- ✅ 视频元素信息正常
- ✅ 视频开始加载
- ✅ 视频数据加载完成
- ✅ 视频可以播放
- ✅ 视频开始播放

如果仍有问题，请提供详细的日志输出！
