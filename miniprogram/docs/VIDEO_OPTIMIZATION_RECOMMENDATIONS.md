# 视频优化建议

## 当前状态分析

从您的日志可以看出：
- ✅ 视频URL可访问（状态码200）
- ✅ Content-Type: `video/mp4`
- ✅ 视频格式正确
- ❌ 但8秒内未开始播放

## 可能的原因

1. **视频文件过大**：文件大小可能超过小程序推荐范围
2. **网络加载缓慢**：视频文件需要更长时间加载
3. **编码参数问题**：虽然格式正确，但编码参数可能仍需优化
4. **小程序环境限制**：某些视频参数可能不被支持

## 优化建议

### 1. 文件大小优化

#### 目标文件大小
- **推荐大小**：5-10MB
- **最大大小**：20MB
- **最小大小**：1-2MB

#### FFmpeg优化命令
```bash
# 超小文件版本（5MB以下）
ffmpeg -i input.mp4 \
  -c:v libx264 \
  -profile:v baseline \
  -level 3.0 \
  -b:v 500k \
  -maxrate 800k \
  -bufsize 1M \
  -c:a aac \
  -b:a 64k \
  -ar 22050 \
  -ac 1 \
  -movflags +faststart \
  -preset ultrafast \
  -crf 30 \
  -vf "scale=480:270:force_original_aspect_ratio=decrease,pad=480:270:(ow-iw)/2:(oh-ih)/2" \
  -t 30 \
  video_small.mp4

# 标准文件版本（10MB以下）
ffmpeg -i input.mp4 \
  -c:v libx264 \
  -profile:v main \
  -level 3.1 \
  -b:v 800k \
  -maxrate 1.2M \
  -bufsize 2M \
  -c:a aac \
  -b:a 96k \
  -ar 44100 \
  -ac 2 \
  -movflags +faststart \
  -preset fast \
  -crf 25 \
  -vf "scale=640:360:force_original_aspect_ratio=decrease,pad=640:360:(ow-iw)/2:(oh-ih)/2" \
  -t 45 \
  video_standard.mp4
```

### 2. 编码参数优化

#### 关键参数说明
```bash
# 视频编码
-c:v libx264          # H.264编码器
-profile:v baseline   # 最兼容的Profile
-level 3.0           # 较低的Level，兼容性更好
-b:v 500k            # 较低码率
-maxrate 800k        # 最大码率限制
-bufsize 1M          # 缓冲区大小

# 音频编码
-c:a aac             # AAC编码
-b:a 64k             # 较低音频码率
-ar 22050            # 较低采样率
-ac 1                # 单声道

# 优化参数
-movflags +faststart # 关键！优化网络播放
-preset ultrafast   # 最快编码速度
-crf 30             # 较高压缩率
```

### 3. 分辨率优化

#### 推荐分辨率
- **超小文件**：480x270 (16:9)
- **标准文件**：640x360 (16:9)
- **高质量文件**：720x405 (16:9)

#### 分辨率选择建议
```bash
# 根据文件大小选择分辨率
if (fileSize < 5MB) {
  resolution = "480x270"
} else if (fileSize < 10MB) {
  resolution = "640x360"
} else {
  resolution = "720x405"
}
```

### 4. 时长优化

#### 推荐时长
- **Banner视频**：15-30秒
- **介绍视频**：30-60秒
- **演示视频**：60-120秒

#### 时长控制
```bash
# 限制视频时长
-t 30              # 30秒
-t 60              # 60秒
-t 120             # 120秒
```

## 测试不同版本

### 1. 创建多个测试版本

```bash
#!/bin/bash
# 创建多个测试版本

INPUT_FILE="input.mp4"

# 版本1：超小文件（3MB以下）
ffmpeg -i "$INPUT_FILE" \
  -c:v libx264 -profile:v baseline -level 3.0 \
  -b:v 400k -maxrate 600k -bufsize 800k \
  -c:a aac -b:a 48k -ar 22050 -ac 1 \
  -movflags +faststart -preset ultrafast -crf 32 \
  -vf "scale=480:270:force_original_aspect_ratio=decrease,pad=480:270:(ow-iw)/2:(oh-ih)/2" \
  -t 20 \
  video_tiny.mp4

# 版本2：小文件（5MB以下）
ffmpeg -i "$INPUT_FILE" \
  -c:v libx264 -profile:v baseline -level 3.0 \
  -b:v 600k -maxrate 900k -bufsize 1.2M \
  -c:a aac -b:a 64k -ar 22050 -ac 1 \
  -movflags +faststart -preset ultrafast -crf 30 \
  -vf "scale=640:360:force_original_aspect_ratio=decrease,pad=640:360:(ow-iw)/2:(oh-ih)/2" \
  -t 30 \
  video_small.mp4

# 版本3：标准文件（10MB以下）
ffmpeg -i "$INPUT_FILE" \
  -c:v libx264 -profile:v main -level 3.1 \
  -b:v 800k -maxrate 1.2M -bufsize 2M \
  -c:a aac -b:a 96k -ar 44100 -ac 2 \
  -movflags +faststart -preset fast -crf 25 \
  -vf "scale=720:405:force_original_aspect_ratio=decrease,pad=720:405:(ow-iw)/2:(oh-ih)/2" \
  -t 45 \
  video_standard.mp4

echo "测试版本创建完成！"
```

### 2. 上传测试

```bash
# 上传到云存储
coscmd upload video_tiny.mp4 video/video_tiny.mp4
coscmd upload video_small.mp4 video/video_small.mp4
coscmd upload video_standard.mp4 video/video_standard.mp4
```

### 3. 小程序测试

```javascript
// 测试不同版本
const testVersions = [
  'cloud://prod-5g94mx7a3d07e78c.7072-prod-5g94mx7a3d07e78c-1353115175/video/video_tiny.mp4',
  'cloud://prod-5g94mx7a3d07e78c.7072-prod-5g94mx7a3d07e78c-1353115175/video/video_small.mp4',
  'cloud://prod-5g94mx7a3d07e78c.7072-prod-5g94mx7a3d07e78c-1353115175/video/video_standard.mp4'
]

// 逐个测试
testVersions.forEach((videoId, index) => {
  console.log(`测试版本 ${index + 1}:`, videoId)
  // 更新videoCosId并测试
})
```

## 监控和调试

### 1. 文件大小监控
```javascript
// 在控制台查看文件大小
const checkFileSize = async (videoUrl) => {
  const result = await wx.request({
    url: videoUrl,
    method: 'HEAD'
  })
  
  const contentLength = result.header['Content-Length']
  const sizeMB = (parseInt(contentLength) / 1024 / 1024).toFixed(2)
  
  console.log('文件大小:', sizeMB + ' MB')
  return sizeMB
}
```

### 2. 加载时间监控
```javascript
// 监控视频加载时间
const startTime = Date.now()

onVideoLoadedData() {
  const loadTime = Date.now() - startTime
  console.log('视频加载时间:', loadTime + 'ms')
}
```

### 3. 播放成功率监控
```javascript
// 监控播放成功率
const playSuccess = () => {
  console.log('✅ 视频播放成功')
  // 记录成功统计
}

const playFailed = () => {
  console.log('❌ 视频播放失败')
  // 记录失败统计
}
```

## 总结

通过以下步骤可以解决视频播放问题：

1. **优化文件大小**：控制在5-10MB以内
2. **使用兼容编码**：H.264 Baseline Profile
3. **降低分辨率**：480x270或640x360
4. **限制时长**：15-30秒
5. **测试多版本**：找到最佳平衡点

建议先创建超小文件版本进行测试，如果成功再逐步增加质量和大小。
