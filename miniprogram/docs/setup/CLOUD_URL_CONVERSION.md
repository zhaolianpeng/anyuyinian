# 云开发URL转换说明

## 概述
本小程序已配置支持云开发格式的图片URL，可以自动将云开发格式的链接转换为可直接在 `<image src="...">` 中使用的URL。

## 云开发URL格式
```
@cloud://{环境ID}.{存储桶名称}/{文件路径}
```

### 示例
```
@cloud://prod-5g94mx7a3d07e78c.7072-prod-5g94mx7a3d07e78c-1353115175/images/button/daiquyao_logo.png
```

### 参数说明
- **环境ID**: `prod-5g94mx7a3d07e78c`
- **存储桶名称**: `7072-prod-5g94mx7a3d07e78c-1353115175`
- **文件路径**: `images/button/daiquyao_logo.png`

## 转换后的URL格式
```
https://{存储桶名称}-{APPID}.cos.{地域}.myqcloud.com/{文件路径}
```

### 示例
```
https://7072-prod-5g94mx7a3d07e78c-1353115175-wx101090677bd5219e.cos.ap-shanghai.myqcloud.com/images/button/daiquyao_logo.png
```

## 使用方法

### 1. 在WXML中直接使用
```xml
<!-- 云开发格式 -->
<image src="@cloud://prod-5g94mx7a3d07e78c.7072-prod-5g94mx7a3d07e78c-1353115175/images/button/daiquyao_logo.png" mode="aspectFill" />

<!-- 转换后的格式 -->
<image src="https://7072-prod-5g94mx7a3d07e78c-1353115175-wx101090677bd5219e.cos.ap-shanghai.myqcloud.com/images/button/daiquyao_logo.png" mode="aspectFill" />
```

### 2. 在JS中自动转换
```javascript
// 使用图片处理函数自动转换
const imageUrl = getCosImageUrl('@cloud://prod-5g94mx7a3d07e78c.7072-prod-5g94mx7a3d07e78c-1353115175/images/button/daiquyao_logo.png')

// 在数据中使用
this.setData({
  imageUrl: imageUrl
})
```

### 3. 在模拟数据中使用
```javascript
// 模拟数据中使用云开发格式
const mockData = {
  banners: [
    {
      id: 1,
      title: '专业护理服务',
      image: '@cloud://prod-5g94mx7a3d07e78c.7072-prod-5g94mx7a3d07e78c-1353115175/images/banner/nursing.jpg',
      url: '/pages/service/list'
    }
  ]
}
```

## 配置信息

### COS配置
```javascript
const COS_CONFIG = {
  // 云开发环境ID
  cloudEnvId: 'prod-5g94mx7a3d07e78c',
  // 云开发存储桶名称
  cloudBucketName: '7072-prod-5g94mx7a3d07e78c-1353115175',
  // COS域名
  bucketDomain: 'https://7072-prod-5g94mx7a3d07e78c-1353115175-wx101090677bd5219e.cos.ap-shanghai.myqcloud.com'
}
```

### 转换函数
```javascript
// 转换云开发URL为直接访问URL
const convertCloudUrlToDirectUrl = (cloudUrl) => {
  // 移除 @cloud:// 前缀
  const urlWithoutPrefix = cloudUrl.replace('@cloud://', '')
  
  // 分割环境ID和存储桶名称
  const parts = urlWithoutPrefix.split('/')
  const envAndBucket = parts[0]
  const envAndBucketParts = envAndBucket.split('.')
  
  const envId = envAndBucketParts[0]
  const bucketName = envAndBucketParts[1]
  const filePath = parts.slice(1).join('/')
  
  // 构建直接访问URL
  return `https://${bucketName}-wx101090677bd5219e.cos.ap-shanghai.myqcloud.com/${filePath}`
}
```

## 支持的图片类型

### 1. Banner图片
- 路径: `images/banner/`
- 格式: JPG, PNG
- 建议尺寸: 750x300px

### 2. 按钮图标
- 路径: `images/button/`
- 格式: PNG
- 建议尺寸: 81x81px

### 3. 服务图片
- 路径: `images/service/`
- 格式: JPG, PNG
- 建议尺寸: 400x300px

### 4. 医院图片
- 路径: `images/hospital/`
- 格式: JPG, PNG
- 建议尺寸: 200x200px

## 测试和验证

### 1. 测试转换功能
```javascript
// 运行测试脚本
node test_cloud_url.js
```

### 2. 验证URL可访问性
1. 在浏览器中访问转换后的URL
2. 确认图片能正常显示
3. 检查图片加载速度

### 3. 小程序测试
1. 在微信开发者工具中测试
2. 检查图片是否正常显示
3. 查看控制台是否有错误

## 常见问题

### Q: 转换后的URL无法访问？
A: 
1. 检查云开发环境ID和存储桶名称是否正确
2. 确认图片文件已上传到云开发存储
3. 检查存储桶权限设置

### Q: 图片显示异常？
A:
1. 检查图片文件格式是否支持
2. 确认图片文件大小是否合理
3. 验证图片路径是否正确

### Q: 转换失败？
A:
1. 检查云开发URL格式是否正确
2. 确认环境ID和存储桶名称配置正确
3. 查看控制台错误信息

## 最佳实践

### 1. 图片命名规范
- 使用有意义的文件名
- 避免使用中文和特殊字符
- 保持文件名简洁明了

### 2. 图片优化
- 压缩图片文件大小
- 选择合适的图片格式
- 使用适当的图片尺寸

### 3. 错误处理
- 为图片加载添加错误处理
- 提供默认图片作为备选
- 记录图片加载错误日志

## 相关文件
- `config.js` - COS配置
- `utils/request.js` - 图片处理函数
- `test_cloud_url.js` - 转换测试脚本
- `CLOUD_URL_CONVERSION.md` - 本说明文档 