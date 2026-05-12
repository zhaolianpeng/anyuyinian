# 开发指南

## 常见问题解决

### 1. 域名配置问题

#### 问题描述
在微信小程序开发中，如果请求的URL不在允许的域名列表中，会出现以下错误：
```
request:fail url not in domain list
```

#### 解决方案

##### 方案一：使用模拟数据（推荐用于开发）
当前项目已配置为开发环境使用模拟数据，无需配置真实域名。

1. **配置说明**
   - 在 `config.js` 中设置 `CURRENT_ENV = ENV.DEV`
   - 在 `config.js` 中设置 `dev.useMockData = true`

2. **模拟数据**
   - 模拟数据在 `utils/request.js` 中的 `mockData` 对象中定义
   - 包含首页数据、服务列表、订单列表等

3. **切换环境**
   ```javascript
   // 在 config.js 中修改
   const CURRENT_ENV = ENV.DEV  // 开发环境
   // const CURRENT_ENV = ENV.PROD  // 生产环境
   ```

##### 方案二：配置真实域名（用于生产环境）

1. **申请小程序AppID**
   - 登录微信公众平台
   - 创建小程序，获取AppID

2. **配置服务器域名**
   - 在微信公众平台 → 开发 → 开发设置 → 服务器域名
   - 添加 `request合法域名`：`https://your-api-domain.com`
   - 添加 `uploadFile合法域名`：`https://your-api-domain.com`
   - 添加 `downloadFile合法域名`：`https://your-api-domain.com`

3. **修改配置文件**
   ```javascript
   // 在 config.js 中修改
   const CURRENT_ENV = ENV.PROD
   const DOMAINS = {
     [ENV.PROD]: 'https://your-api-domain.com'
   }
   ```

4. **配置开发工具**
   - 在微信开发者工具中，点击右上角"详情"
   - 勾选"不校验合法域名、web-view（业务域名）、TLS 版本以及 HTTPS 证书"

### 2. 登录问题

#### 问题描述
如果遇到以下错误：
```
getUserProfile:fail can only be invoked by user TAP gesture.
```

#### 原因分析
`wx.getUserProfile` API 只能通过用户点击事件触发，不能在页面加载时自动调用。

#### 解决方案

1. **修改登录逻辑**
   ```javascript
   // 错误做法：在 onLoad 中自动调用
   onLoad() {
     this.handleLogin()  // 这会导致错误
   }

   // 正确做法：用户点击触发
   onLoginTap() {
     this.handleLogin()  // 用户点击按钮时调用
   }
   ```

2. **添加登录按钮**
   ```xml
   <!-- 在 WXML 中添加登录按钮 -->
   <button bindtap="onLoginTap">微信登录</button>
   ```

3. **完善用户界面**
   - 添加登录按钮和状态提示
   - 包含用户协议和隐私政策
   - 优化登录流程的用户体验

### 3. 图片显示问题

#### 问题描述
如果遇到本地图片不显示的问题，通常是因为：
1. 字段名不匹配
2. 图片路径错误
3. 图片文件不存在

#### 解决方案

##### 方案一：检查字段名映射
确保WXML模板中的字段名与模拟数据中的字段名一致：

```javascript
// 模拟数据中的字段名
{
  banners: [
    {
      id: 1,
      title: '专业护理服务',
      image: '/images/fuwu_1.jpeg',  // 注意这里是 image
      url: '/pages/service/list',    // 注意这里是 url
      sort: 1,
      status: 1
    }
  ]
}
```

```xml
<!-- WXML模板中应该使用相同的字段名 -->
<image src="{{item.image}}" mode="aspectFill" class="banner-image" />
```

##### 方案二：检查图片路径
1. **确保图片文件存在**
   ```bash
   # 检查图片文件是否存在
   ls -la miniprogram/images/
   ```

2. **使用正确的路径格式**
   ```javascript
   // 正确的路径格式
   image: '/images/fuwu_1.jpeg'  // 以 / 开头
   ```

3. **检查文件扩展名**
   - 确保文件扩展名正确（.jpg, .jpeg, .png, .gif）
   - 确保文件名大小写正确

##### 方案三：添加图片加载错误处理
```xml
<image 
  src="{{item.image}}" 
  mode="aspectFill" 
  class="banner-image"
  binderror="onImageError"
  bindload="onImageLoad"
/>
```

```javascript
// 在JS文件中添加错误处理
onImageError(e) {
  console.error('图片加载失败:', e.detail)
  // 可以设置默认图片
  const index = e.currentTarget.dataset.index
  const banners = this.data.banners
  banners[index].image = '/images/default.png'
  this.setData({ banners })
},

onImageLoad(e) {
  console.log('图片加载成功:', e.detail)
}
```

### 4. 图片加载问题

#### 问题描述
如果遇到图片加载失败的错误：
```
Failed to load image <URL>
net::ERR_NAME_NOT_RESOLVED
```

#### 解决方案

##### 方案一：使用本地图片（推荐）
1. 将图片文件放入 `images/` 目录
2. 在模拟数据中使用相对路径，如：`/images/icon-appointment.png`
3. 确保图片文件存在

##### 方案二：使用可靠的CDN
1. 使用可靠的图片CDN服务
2. 确保CDN域名在微信小程序中可访问
3. 在模拟数据中使用完整的HTTPS URL

##### 方案三：暂时不使用图片
1. 在模拟数据中将图片字段设为空字符串
2. 在页面模板中添加图片加载失败的处理
3. 使用文字或图标字体替代图片

## 开发环境配置

#### 当前配置
```javascript
// config.js
const CURRENT_ENV = ENV.DEV
const DOMAINS = {
  [ENV.DEV]: 'https://api.example.com',  // 开发环境域名
   [ENV.PROD]: 'https://api.succ.online/anyuyinian'  // 生产环境域名
}
```

#### 模拟数据配置
```javascript
// config.js
dev: {
  useMockData: true,    // 使用模拟数据
  mockDelay: 1000       // 模拟请求延迟（毫秒）
}
```

## 生产环境部署

1. **修改环境配置**
   ```javascript
   // config.js
   const CURRENT_ENV = ENV.PROD
   ```

2. **配置真实域名**
   - 确保后端API服务正常运行
   - 在微信公众平台配置合法域名
   - 测试所有接口功能正常

3. **发布小程序**
   - 在微信开发者工具中点击"上传"
   - 在微信公众平台提交审核
   - 审核通过后发布上线

## 调试技巧

#### 1. 查看网络请求
- 在微信开发者工具中打开"调试器"
- 查看"Network"面板的网络请求

#### 2. 查看控制台日志
- 在微信开发者工具中打开"调试器"
- 查看"Console"面板的日志输出

#### 3. 模拟数据调试
- 修改 `utils/request.js` 中的 `mockData` 对象
- 添加或修改模拟数据

#### 4. 切换环境
```javascript
// 开发环境（使用模拟数据）
const CURRENT_ENV = ENV.DEV

// 生产环境（使用真实API）
const CURRENT_ENV = ENV.PROD
```

#### 5. 图片调试
- 检查图片文件是否存在
- 检查字段名是否正确
- 添加图片加载错误处理
- 使用开发者工具的"Sources"面板查看图片路径

## 常见问题

#### Q: 为什么会出现域名限制错误？
A: 微信小程序出于安全考虑，只允许请求在合法域名列表中的URL。

#### Q: 如何快速测试功能？
A: 使用模拟数据，无需配置真实域名，可以快速开发和测试。

#### Q: 如何添加新的模拟数据？
A: 在 `utils/request.js` 的 `mockData` 对象中添加新的接口数据。

#### Q: 如何切换到真实API？
A: 修改 `config.js` 中的 `CURRENT_ENV` 为 `ENV.PROD`，并配置合法域名。

#### Q: 为什么登录按钮点击无效？
A: 确保登录按钮绑定了正确的事件处理函数，并且 `wx.getUserProfile` 只在用户点击时调用。

#### Q: 为什么本地图片不显示？
A: 检查字段名是否匹配、图片路径是否正确、文件是否存在。

#### Q: 如何调试图片加载问题？
A: 在WXML中添加 `binderror` 和 `bindload` 事件，在控制台查看错误信息。

## 注意事项

1. **开发环境**：建议使用模拟数据进行开发和测试
2. **生产环境**：必须配置真实的API域名
3. **域名要求**：必须是HTTPS协议，且已备案
4. **证书要求**：域名必须配置有效的SSL证书
5. **测试环境**：可以在开发者工具中勾选"不校验合法域名"进行测试
6. **API调用**：某些API（如 `wx.getUserProfile`）只能通过用户交互触发
7. **图片路径**：本地图片路径必须以 `/` 开头，相对于小程序根目录
8. **字段名**：确保WXML模板中的字段名与数据源中的字段名一致

## 相关文件

- `config.js` - 环境配置
- `utils/request.js` - 请求工具和模拟数据
- `pages/login/login.js` - 登录页面逻辑
- `pages/login/login.wxml` - 登录页面模板
- `pages/index/index.wxml` - 首页模板
- `pages/index/index.js` - 首页逻辑
- `images/` - 图片资源目录
- `DEVELOPMENT.md` - 开发指南（本文件）
- `README.md` - 项目说明 