# 视频文件大小优化指南

## 当前问题
- **文件大小**: 14.16 MB
- **推荐大小**: 5-10 MB
- **问题**: 文件过大导致加载缓慢，8秒内无法开始播放

## 优化方案

### 方案1：超小文件版本（推荐）
```bash
# 目标：3-5MB，快速加载
ffmpeg -i input.mp4 \
  -c:v libx264 \
  -profile:v baseline \
  -level 3.0 \
  -b:v 400k \
  -maxrate 600k \
  -bufsize 800k \
  -c:a aac \
  -b:a 48k \
  -ar 22050 \
  -ac 1 \
  -movflags +faststart \
  -preset ultrafast \
  -crf 32 \
  -vf "scale=480:270:force_original_aspect_ratio=decrease,pad=480:270:(ow-iw)/2:(oh-ih)/2" \
  -t 30 \
  video_optimized_small.mp4
```

### 方案2：标准文件版本
```bash
# 目标：5-8MB，平衡质量和大小
ffmpeg -i input.mp4 \
  -c:v libx264 \
  -profile:v baseline \
  -level 3.0 \
  -b:v 600k \
  -maxrate 900k \
  -bufsize 1.2M \
  -c:a aac \
  -b:a 64k \
  -ar 22050 \
  -ac 1 \
  -movflags +faststart \
  -preset fast \
  -crf 30 \
  -vf "scale=640:360:force_original_aspect_ratio=decrease,pad=640:360:(ow-iw)/2:(oh-ih)/2" \
  -t 45 \
  video_optimized_standard.mp4
```

### 方案3：高质量小文件版本
```bash
# 目标：8-10MB，保持较好质量
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
  -vf "scale=720:405:force_original_aspect_ratio=decrease,pad=720:405:(ow-iw)/2:(oh-ih)/2" \
  -t 60 \
  video_optimized_quality.mp4
```

## 参数说明

### 关键优化参数
- `-b:v 400k`: 视频码率400kbps（原文件可能>1Mbps）
- `-maxrate 600k`: 最大码率600kbps
- `-bufsize 800k`: 缓冲区800KB
- `-crf 32`: 较高压缩率（原文件可能23-25）
- `-vf "scale=480:270"`: 降低分辨率到480x270
- `-t 30`: 限制时长30秒
- `-preset ultrafast`: 最快编码速度

### 文件大小对比
- **原文件**: 14.16 MB
- **方案1**: ~3-4 MB (减少70%+)
- **方案2**: ~5-6 MB (减少60%+)
- **方案3**: ~8-9 MB (减少40%+)

## 上传和测试

### 1. 上传优化后的视频
```bash
# 上传到云存储
coscmd upload video_optimized_small.mp4 video/video_optimized_small.mp4
coscmd upload video_optimized_standard.mp4 video/video_optimized_standard.mp4
coscmd upload video_optimized_quality.mp4 video/video_optimized_quality.mp4
```

### 2. 更新小程序配置
```javascript
// 测试不同版本
data: {
  // 方案1：超小文件
  videoCosId: 'cloud://prod-5g94mx7a3d07e78c.7072-prod-5g94mx7a3d07e78c-1353115175/video/video_optimized_small.mp4',
  
  // 方案2：标准文件
  // videoCosId: 'cloud://prod-5g94mx7a3d07e78c.7072-prod-5g94mx7a3d07e78c-1353115175/video/video_optimized_standard.mp4',
  
  // 方案3：高质量文件
  // videoCosId: 'cloud://prod-5g94mx7a3d07e78c.7072-prod-5g94mx7a3d07e78c-1353115175/video/video_optimized_quality.mp4',
}
```

### 3. 测试播放效果
```javascript
// 在控制台测试
const pages = getCurrentPages()
const currentPage = pages[pages.length - 1]

// 检查文件大小
currentPage.checkVideoInfo()

// 测试播放
currentPage.testVideoPlayback()
```

## 预期结果

优化后应该看到：
- ✅ 文件大小：3-10MB
- ✅ 加载时间：<3秒
- ✅ 播放成功：无降级模式
- ✅ 用户体验：流畅播放

## 建议

1. **优先使用方案1**：超小文件版本，确保快速加载
2. **如果质量不够**：再尝试方案2或方案3
3. **监控加载时间**：目标<3秒开始播放
4. **保持降级模式**：作为备用方案

通过文件大小优化，应该能解决视频播放问题！
