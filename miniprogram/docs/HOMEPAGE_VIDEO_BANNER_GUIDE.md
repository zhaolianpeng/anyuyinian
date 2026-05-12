# 首页视频播放器轮播Banner功能指南

## 功能概述
在首页添加了轮播Banner功能，包含两个Banner：
1. **第一个Banner**：原有的图片内容（24小时服务电话）
2. **第二个Banner**：视频播放器，播放腾讯云对象存储中的护理服务视频

## 技术实现

### 1. 文件结构
```
miniprogram/
├── pages/index/
│   ├── index.wxml          # 更新了轮播Banner结构
│   ├── index.wxss          # 添加了视频播放器样式
│   └── index.js            # 添加了视频加载和播放逻辑
├── utils/
│   └── videoService.js     # 视频服务工具类
└── images/
    └── video-poster.jpg    # 视频封面图片（需要替换）
```

### 2. 核心功能

#### 轮播Banner组件
- 使用 `swiper` 组件实现轮播
- 自动播放：5秒间隔
- 循环播放：`circular="{{true}}"`
- 指示器：显示当前位置

#### 视频播放器
- 视频源：腾讯云对象存储
- 自动静音：`muted="{{true}}"`
- 循环播放：`loop="{{true}}"`
- 显示播放按钮：`show-play-btn="{{true}}"`
- 居中播放按钮：`show-center-play-btn="{{true}}"`
- 对象适配：`object-fit="cover"`

### 3. 视频服务工具类

#### `videoService.js` 功能
- `getVideoUrl(videoCosId, time)`: 获取视频临时访问地址
- `getTempVideoFile(fileID, time)`: 批量获取临时文件URL
- `preloadVideo(videoUrl)`: 预加载视频

#### 使用示例
```javascript
const { getVideoUrl } = require('../../utils/videoService')

// 获取视频URL
const videoUrl = await getVideoUrl('cloud://xxx/video/xxx.mp4')
```

### 4. 样式设计

#### 轮播容器
```css
.banner-swiper {
  height: 200px;
  margin: 20px;
  border-radius: 20px;
  overflow: hidden;
}
```

#### 视频Banner
```css
.video-banner {
  position: relative;
  height: 100%;
  border-radius: 20px;
  overflow: hidden;
  background: #000;
}

.banner-video {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
```

#### 视频覆盖层
```css
.video-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(0,0,0,0.3) 0%, rgba(0,0,0,0.1) 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}
```

## 配置说明

### 1. 视频配置
在 `pages/index/index.js` 中配置：
```javascript
data: {
  // 视频相关
  videoUrl: '',                    // 视频播放地址（动态获取）
  videoPoster: '/images/video-poster.jpg',  // 视频封面图
  videoCosId: 'cloud://prod-5g94mx7a3d07e78c.7072-prod-5g94mx7a3d07e78c-1353115175/video/光厂_5440331_照顾瘫痪老人按摩理发.mp4'
}
```

### 2. 轮播配置
在 `index.wxml` 中配置：
```xml
<swiper class="banner-swiper" 
        indicator-dots="{{true}}" 
        autoplay="{{true}}" 
        interval="{{5000}}" 
        duration="{{500}}"
        circular="{{true}}"
        indicator-color="rgba(255,255,255,0.5)"
        indicator-active-color="#ffffff">
```

## 使用方法

### 1. 替换视频封面
将实际的视频封面图片文件放在：
```
miniprogram/images/video-poster.jpg
```
建议尺寸：750x400px，格式：JPG或PNG

### 2. 更换视频
修改 `videoCosId` 为新的视频文件ID：
```javascript
videoCosId: 'cloud://your-env/your-video-path/video.mp4'
```

### 3. 调整轮播时间
修改 `interval` 属性：
```xml
interval="{{3000}}"  <!-- 3秒切换 -->
```

## 事件处理

### 1. 视频播放事件
```javascript
onVideoPlay() {
  console.log('视频开始播放')
}
```

### 2. 视频暂停事件
```javascript
onVideoPause() {
  console.log('视频暂停')
}
```

### 3. 视频错误处理
```javascript
onVideoError(e) {
  console.error('视频播放错误:', e.detail)
  wx.showToast({
    title: '视频加载失败',
    icon: 'none'
  })
}
```

## 性能优化

### 1. 视频预加载
- 页面加载时自动获取视频临时URL
- 支持视频预加载到本地缓存

### 2. 错误处理
- 视频加载失败时显示封面图
- 网络错误时提供友好提示

### 3. 用户体验
- 自动静音避免打扰用户
- 循环播放提升观看体验
- 覆盖层文字说明视频内容

## 注意事项

1. **视频格式**：建议使用MP4格式，兼容性最好
2. **文件大小**：建议控制在10MB以内，避免加载过慢
3. **网络环境**：视频加载依赖网络，建议添加加载状态提示
4. **存储费用**：使用腾讯云对象存储会产生流量费用
5. **临时URL**：视频URL有效期为24小时，需要定期刷新

## 扩展功能

### 1. 添加更多Banner
在 `swiper-item` 中添加更多轮播项

### 2. 视频控制
可以添加播放/暂停按钮、进度条等控制元素

### 3. 数据统计
可以添加视频播放次数、观看时长等统计功能

## 相关文件
- `pages/index/index.wxml` - 轮播Banner结构
- `pages/index/index.wxss` - 视频播放器样式
- `pages/index/index.js` - 视频加载逻辑
- `utils/videoService.js` - 视频服务工具类
- `images/video-poster.jpg` - 视频封面图片
