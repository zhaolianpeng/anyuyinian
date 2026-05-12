# COS 配置说明

## 概述
本小程序已配置为使用腾讯云对象存储(COS)来存储和拉取图片资源。所有图片请求都会通过COS服务进行加载。

## 配置步骤

### 1. 创建COS存储桶
1. 登录腾讯云控制台
2. 进入对象存储COS服务
3. 创建新的存储桶
4. 记录存储桶域名

### 2. 配置COS域名
在 `config.js` 中更新COS配置：

```javascript
const COS_CONFIG = {
  // 替换为您的COS存储桶域名
  // 格式: https://{存储桶名称}-{APPID}.cos.{地域}.myqcloud.com
  bucketDomain: 'https://your-bucket-name-1234567890.cos.ap-beijing.myqcloud.com',
  imagePrefix: '/images/',
  defaultImages: {
    service: '/images/default-service.png',
    hospital: '/images/default-hospital.png',
    user: '/images/default-user.png'
  }
}
```

**详细配置步骤请参考**: [COS域名配置指南](./COS_DOMAIN_SETUP.md)

### 3. 上传图片到COS
将以下图片上传到COS存储桶的 `/images/` 目录：

#### Banner图片
- `/images/banner-nursing.jpg` - 护理服务banner
- `/images/banner-escort.jpg` - 陪诊服务banner

#### 导航图标
- `/images/icon-service.png` - 服务预约图标
- `/images/icon-order.png` - 我的订单图标
- `/images/icon-hospital.png` - 医院信息图标
- `/images/icon-user.png` - 个人中心图标

#### 服务图片
- `/images/service-nursing.jpg` - 上门护理服务
- `/images/service-escort.jpg` - 专业陪诊服务
- `/images/service-care.jpg` - 生活照护服务
- `/images/service-rehab.jpg` - 康复护理服务

#### 医院图片
- `/images/hospital-1.jpg` - 深圳市人民医院
- `/images/hospital-2.jpg` - 深圳市第二人民医院

#### 默认图片
- `/images/default-service.png` - 默认服务图片
- `/images/default-hospital.png` - 默认医院图片
- `/images/default-user.png` - 默认用户图片

### 4. 图片规格建议
- **Banner图片**: 750x300px, JPG格式
- **导航图标**: 81x81px, PNG格式, 透明背景
- **服务图片**: 400x300px, JPG格式
- **医院图片**: 200x200px, JPG格式
- **文件大小**: 建议小于100KB

### 5. 权限配置
确保COS存储桶的访问权限设置为：
- **公有读私有写**: 允许公开读取图片
- **防盗链**: 可选择性配置

## 本地调试配置

### 1. 环境设置
当前配置为使用生产环境：
```javascript
const CURRENT_ENV = ENV.PROD
```

### 2. 模拟数据
在开发环境中，图片会通过COS处理函数自动转换为完整URL：
```javascript
// 相对路径
'/images/service-nursing.jpg'

// 转换为完整URL
'https://your-bucket-name.cos.ap-region.myqcloud.com/images/service-nursing.jpg'
```

### 3. 图片加载处理
系统会自动处理以下图片字段：
- `image` - 图片
- `icon` - 图标
- `logo` - 标志
- `avatar` - 头像
- `banner` - 横幅
- `images` - 图片数组

## 故障排除

### 1. 图片加载失败
- 检查COS存储桶域名是否正确
- 确认图片文件已上传到正确路径
- 检查网络连接和COS服务状态

### 2. 域名配置问题
- 确保COS域名在微信小程序中可访问
- 检查COS存储桶的访问权限设置

### 3. 图片显示异常
- 检查图片文件格式是否支持
- 确认图片文件大小是否合理
- 验证图片路径是否正确

## 开发建议

### 1. 图片优化
- 使用适当的图片格式(JPG/PNG)
- 压缩图片文件大小
- 选择合适的图片尺寸

### 2. 缓存策略
- 利用COS的CDN加速
- 设置合适的缓存头
- 考虑图片懒加载

### 3. 错误处理
- 添加图片加载失败的处理
- 提供默认图片作为备选
- 记录图片加载错误日志

## 相关文件
- `config.js` - COS配置
- `utils/request.js` - 图片处理函数
- `COS_SETUP.md` - 本配置文档 