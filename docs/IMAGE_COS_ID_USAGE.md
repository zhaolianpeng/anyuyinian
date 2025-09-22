# imageCosId 字段使用指南

## 概述
`imageCosId` 字段用于存储腾讯云对象存储的图片ID，小程序可以通过此字段获取图片的临时访问地址，实现图片资源的统一管理。

## 实现功能

### 1. 图片服务工具类
**文件**: `utils/imageService.js`

提供以下功能：
- `getTempFile(fileID, time)` - 获取单个或多个文件的临时访问地址
- `getServiceImages(services, time)` - 批量处理服务列表的图片
- `getSingleServiceImage(service, time)` - 处理单个服务的图片
- `preloadServiceImages(services, time)` - 预加载服务图片到本地缓存

### 2. 页面集成

#### 服务列表页面 (`pages/service/list.js`)
```javascript
// 使用 imageCosId 获取临时访问地址
const processedServices = await getServiceImages(newServices)

// 处理图片URL（兼容旧版本）
const finalServices = processedServices.map(service => ({
  ...service,
  imageUrl: service.imageTempUrl || processImageUrl(service.imageUrl)
}))
```

#### 服务详情页面 (`pages/service/detail.js`)
```javascript
// 使用 imageCosId 获取临时访问地址
const processedService = await getSingleServiceImage(service)

// 处理图片URL（兼容旧版本）
processedService.imageUrl = processedService.imageTempUrl || processImageUrl(processedService.imageUrl)
```

#### 首页服务展示 (`pages/index/index.js`)
```javascript
// 使用 imageCosId 获取临时访问地址
const servicesWithImages = await getServiceImages(services)

const processedServices = servicesWithImages.map(service => ({
  ...service,
  imageUrl: service.imageTempUrl || service.imageUrl || '/images/service/default-service.jpg'
}))
```

## 使用示例

### 1. 基本使用
```javascript
const { getTempFile, getServiceImages } = require('../../utils/imageService')

// 获取单个文件的临时访问地址
const result = await getTempFile('cloud://test.png', 86400)
console.log(result.fileList[0].tempFileURL)

// 批量处理服务图片
const services = [
  { id: 1, name: '服务1', imageCosId: 'cloud://service1.jpg' },
  { id: 2, name: '服务2', imageCosId: 'cloud://service2.jpg' }
]
const processedServices = await getServiceImages(services)
```

### 2. 在页面中使用
```javascript
Page({
  async onLoad() {
    // 获取服务列表
    const result = await api.serviceList({ page: 1, pageSize: 10 })
    
    if (result.code === 0) {
      // 使用 imageCosId 获取临时访问地址
      const servicesWithImages = await getServiceImages(result.data.list)
      
      this.setData({
        services: servicesWithImages
      })
    }
  }
})
```

### 3. 预加载图片
```javascript
// 预加载图片到本地缓存，提高显示速度
const processedServices = await preloadServiceImages(services)
```

## 字段说明

### imageCosId 字段
- **类型**: `VARCHAR(255)`
- **用途**: 存储腾讯云对象存储的图片ID
- **格式**: `cloud://bucket-name/path/to/image.jpg`
- **示例**: `cloud://my-bucket/images/service1.jpg`

### 临时访问地址
- **有效期**: 默认24小时（86400秒）
- **格式**: `https://temp-file-url`
- **用途**: 用于在小程序中显示图片

## 兼容性处理

### 1. 向后兼容
- 如果 `imageCosId` 为空，使用原始的 `imageUrl`
- 如果获取临时地址失败，回退到原始 `imageUrl`
- 保持现有代码的稳定性

### 2. 错误处理
```javascript
// 自动处理错误，确保图片正常显示
const processedService = await getSingleServiceImage(service)
// 如果获取失败，processedService.imageTempUrl 为 null
// 会使用 service.imageUrl 作为备用
```

## 性能优化

### 1. 批量处理
- 一次性获取多个图片的临时访问地址
- 减少网络请求次数
- 提高加载速度

### 2. 预加载
- 提前加载图片到本地缓存
- 减少用户等待时间
- 提升用户体验

### 3. 缓存机制
- 临时访问地址有24小时有效期
- 避免频繁请求
- 降低服务器压力

## 注意事项

1. **网络依赖**: 需要网络连接才能获取临时访问地址
2. **有效期限制**: 临时访问地址有24小时有效期
3. **错误处理**: 需要处理获取失败的情况
4. **性能考虑**: 大量图片时建议使用预加载

## 相关文件
- `utils/imageService.js` - 图片服务工具类
- `pages/service/list.js` - 服务列表页面
- `pages/service/detail.js` - 服务详情页面
- `pages/index/index.js` - 首页
- `db/model/service.go` - 数据库模型
- `service/admin_service.go` - 管理员服务接口
