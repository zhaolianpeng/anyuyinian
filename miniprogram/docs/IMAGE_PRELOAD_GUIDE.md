# 图片预加载功能使用指南

## 功能概述

图片预加载功能可以显著提升小程序的用户体验，通过提前加载图片资源，减少用户等待时间，让图片显示更加流畅。

## 主要特性

### 🚀 核心功能
- **智能预加载**：自动识别需要预加载的图片
- **批量处理**：支持批量预加载多张图片
- **进度监控**：实时显示预加载进度
- **缓存管理**：智能缓存机制，避免重复加载
- **错误处理**：完善的错误处理和重试机制

### 📱 用户体验
- **静默预加载**：后台静默预加载，不显示进度提示
- **非阻塞**：预加载不阻塞页面渲染
- **智能调度**：控制并发数量，避免性能问题

## 使用方法

### 1. 基础用法

```javascript
const { imagePreloader } = require('../../utils/imagePreloader')

// 预加载单张图片
try {
  const result = await imagePreloader.preloadImage('/images/example.jpg')
  console.log('预加载成功:', result)
} catch (error) {
  console.error('预加载失败:', error)
}
```

### 2. 批量预加载

```javascript
const imageList = [
  '/images/image1.jpg',
  '/images/image2.jpg',
  '/images/image3.jpg'
]

await imagePreloader.preloadImages(
  imageList,
  (progress) => {
    // 进度回调
    console.log(`进度: ${progress.progress}%`)
  },
  (result) => {
    // 完成回调
    console.log('预加载完成:', result)
  }
)
```

### 3. 页面集成

#### 首页预加载
```javascript
// 在页面数据加载完成后
this.preloadHomeImages()

// 预加载方法
async preloadHomeImages() {
  const homeData = { services: this.data.currentServices }
  
await preloadHomeImages(
  homeData,
  (progress) => {
    console.log(`首页服务图片预加载进度: ${progress.progress}% (${progress.loaded}/${progress.total})`)
  },
  (result) => {
    console.log('首页图片预加载完成:', result)
  }
)
}
```

#### 服务页面预加载
```javascript
// 在服务数据加载完成后
this.preloadServiceImages(services)

// 预加载方法
async preloadServiceImages(services) {
  const imageUrls = services
    .map(service => service.imageUrl)
    .filter(url => url && url.trim() !== '')

await preloadImages(
  imageUrls,
  (progress) => {
    console.log(`服务图片预加载进度: ${progress.progress}% (${progress.loaded}/${progress.total})`)
  },
  (result) => {
    console.log('服务图片预加载完成:', result)
  }
)
}
```

### 4. 全局预加载

在 `app.js` 中预加载关键图片：

```javascript
const { preloadCriticalImages } = require('./utils/imagePreloader')

App({
  onLaunch() {
    // 预加载关键图片
    this.preloadCriticalImages()
  },

  async preloadCriticalImages() {
    await preloadCriticalImages(
      (progress) => {
        console.log(`关键图片预加载进度: ${progress.progress}%`)
      },
      (result) => {
        console.log('关键图片预加载完成:', result)
      }
    )
  }
})
```

## 配置选项

### 预加载器配置

```javascript
// 修改配置
imagePreloader.maxConcurrent = 5  // 最大并发数
imagePreloader.timeout = 15000    // 超时时间（毫秒）
```

### 缓存管理

```javascript
// 获取缓存信息
const cacheInfo = imagePreloader.getCacheInfo()
console.log('缓存统计:', cacheInfo)

// 清理过期缓存（默认24小时）
imagePreloader.clearCache()

// 清空所有缓存
imagePreloader.clearAllCache()

// 检查图片是否已缓存
const isCached = imagePreloader.isCached('/images/example.jpg')
```

## 静默预加载

图片预加载现在采用静默模式，在后台进行，不会显示任何进度提示。预加载过程完全透明，用户不会感知到预加载的存在，但能享受到更快的图片加载体验。

### 预加载日志

预加载过程会在控制台输出日志，方便开发调试：

```javascript
// 控制台输出示例
关键图片预加载进度: 25% (1/4)
关键图片预加载进度: 50% (2/4)
关键图片预加载进度: 75% (3/4)
关键图片预加载进度: 100% (4/4)
关键图片预加载完成: {success: [...], failed: [...], total: 4, successCount: 4, failedCount: 0}
```

## 性能优化建议

### 1. 图片优化
- 使用适当的图片格式（WebP、JPEG、PNG）
- 控制图片大小，避免过大的图片
- 使用 CDN 加速图片加载

### 2. 预加载策略
- 只预加载用户可能看到的图片
- 避免预加载过多图片导致内存问题
- 根据网络状况调整预加载策略

### 3. 缓存管理
- 定期清理过期缓存
- 监控缓存使用情况
- 避免缓存过多图片

## 测试

运行测试脚本验证功能：

```javascript
const { runAllTests } = require('./tests/test_image_preload')

// 运行所有测试
runAllTests().then(results => {
  console.log('测试完成:', results)
})
```

## 常见问题

### Q: 预加载失败怎么办？
A: 检查图片URL是否正确，网络是否正常，图片是否存在。预加载失败不会影响页面正常显示。

### Q: 如何控制预加载的图片数量？
A: 可以通过过滤图片列表来控制预加载数量，只预加载必要的图片。

### Q: 预加载会影响页面性能吗？
A: 预加载采用分批处理，控制并发数量，不会显著影响页面性能。

### Q: 如何清理缓存？
A: 使用 `imagePreloader.clearCache()` 清理过期缓存，或 `imagePreloader.clearAllCache()` 清空所有缓存。

### Q: 预加载时出现"file not found"错误怎么办？
A: 这是正常现象，表示某些图片文件不存在。系统会自动跳过这些图片，不会影响其他图片的预加载。可以通过图片路径检查器来验证图片是否存在。

### Q: 如何添加新的图片到预加载列表？
A: 将图片文件放在 `images` 目录下，然后在预加载方法中添加对应的路径。系统会自动检查图片是否存在。

## 更新日志

### v1.2.0
- 移除预加载进度提示，改为静默预加载
- 预加载过程完全在后台进行，不干扰用户体验
- 保留控制台日志输出，方便开发调试

### v1.1.0
- 修复图片路径不存在导致的预加载失败问题
- 添加图片路径检查器，自动过滤无效路径
- 改进错误处理，静默处理文件不存在的错误
- 优化预加载策略，只预加载实际存在的图片

### v1.0.0
- 实现基础图片预加载功能
- 添加进度监控和缓存管理
- 支持首页和服务页面预加载
- 添加美观的进度提示UI

---

如有问题或建议，请联系开发团队。
