# 域名配置说明

## 重要提醒

当前小程序已经统一切换到自有服务端直连。

当前生产基础地址应以 `https://api.succ.online/anyuyinian` 为准，不再使用历史的 tcloudbase 云托管域名。

## 需要配置的正式域名

### 1. API服务器域名
您需要将后端服务部署到正式的服务器上，并配置域名。

**当前配置：**
```javascript
[ENV.PROD]: 'https://api.succ.online/anyuyinian'
```

### 2. COS存储域名
COS域名可以继续使用，但建议也配置到您自己的域名下。

**当前配置：**
```javascript
bucketDomain: 'https://7072-prod-5g94mx7a3d07e78c-1353115175.cos.ap-shanghai.myqcloud.com'
```

## 部署步骤

### 1. 后端服务部署
1. 将Go后端代码部署到您的服务器
2. 配置HTTPS证书
3. 配置域名解析
4. 测试API接口可访问性

### 2. 域名配置
1. 修改 `miniprogram/config.js` 中的生产环境域名
2. 在微信小程序后台配置新的域名
3. 测试网络请求

### 3. 微信小程序后台配置
在微信公众平台配置以下域名：

**request合法域名：**
```
https://api.yourcompany.com  // 您的正式API域名
```

**uploadFile合法域名：**
```
https://api.yourcompany.com  // 您的正式API域名
```

**downloadFile合法域名：**
```
https://7072-prod-5g94mx7a3d07e78c-1353115175.cos.ap-shanghai.myqcloud.com
```

## 环境切换

### 开发环境
```javascript
const CURRENT_ENV = ENV.DEV
```

### 生产环境
```javascript
const CURRENT_ENV = ENV.PROD
```

## 测试验证

### 1. 开发环境测试
```bash
# 在微信开发者工具中测试
# 确保勾选"不校验合法域名"
```

### 2. 生产环境测试
```bash
# 在真机上测试
# 确保域名已配置到小程序后台
```

### 3. 网络测试
```bash
# 测试API域名
curl -I https://api.succ.online/anyuyinian/api/home/init

# 测试COS域名
curl -I https://7072-prod-5g94mx7a3d07e78c-1353115175.cos.ap-shanghai.myqcloud.com
```

## 常见问题

### 1. 域名未备案
**问题**：域名未备案导致无法使用
**解决**：确保域名已通过ICP备案

### 2. HTTPS证书问题
**问题**：HTTPS证书无效
**解决**：配置有效的SSL证书

### 3. 服务器配置问题
**问题**：服务器无法访问
**解决**：检查服务器防火墙和网络配置

## 配置清单

### 服务器配置
- [ ] 部署Go后端服务
- [ ] 配置HTTPS证书
- [ ] 配置域名解析
- [ ] 测试API接口

### 小程序配置
- [ ] 修改config.js中的域名
- [ ] 配置微信小程序后台域名
- [ ] 测试网络请求
- [ ] 上传体验版测试

### 域名要求
- [ ] 域名已备案
- [ ] 使用HTTPS协议
- [ ] 证书有效
- [ ] 服务器稳定

## 联系支持

如果您需要帮助配置正式域名，请提供：
1. 您的正式域名
2. 服务器部署情况
3. 域名备案信息 