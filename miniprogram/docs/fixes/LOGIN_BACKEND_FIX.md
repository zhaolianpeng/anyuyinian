# 微信小程序登录后端服务问题解决方案

## 问题描述
`getUserProfile` 调用成功，但后端登录接口调用失败，错误信息：`请求失败`

## 问题分析
从日志可以看出：
1. ✅ `wx.getUserProfile` 调用成功
2. ✅ `wx.login` 调用成功
3. ❌ 后端登录接口调用失败

## 可能的原因

### 1. 域名限制
真机调试时可能无法访问生产环境的域名，需要在微信小程序后台配置合法域名。

### 2. 网络问题
后端服务可能不可用或网络连接有问题。

### 3. 请求格式问题
后端可能期望不同的请求格式。

## 解决方案

### 1. 切换到开发环境使用模拟数据
临时切换到开发环境，使用模拟数据进行测试：

```javascript
// 在 config.js 中修改
const CURRENT_ENV = ENV.DEV  // 改为开发环境
```

### 2. 添加详细的错误处理和调试信息
在登录代码中添加更详细的日志：

```javascript
// 调用后端登录接口
console.log('调用后端登录接口...')
console.log('请求参数:', { code: loginRes.code, userInfo: userInfoRes.userInfo })

wxLogin(loginRes.code, userInfoRes.userInfo)
  .then((res) => {
    console.log('后端登录接口响应:', res)
    // 处理成功响应...
  })
  .catch((error) => {
    console.error('登录失败:', error)
    console.error('错误详情:', {
      message: error.message,
      errMsg: error.errMsg,
      stack: error.stack
    })
    
    // 针对网络错误的特殊处理
    if (error.message && error.message.includes('网络错误')) {
      wx.showModal({
        title: '网络错误',
        content: '请检查网络连接，或稍后重试',
        showCancel: false
      })
    } else if (error.message && error.message.includes('url not in domain list')) {
      wx.showModal({
        title: '域名限制',
        content: '当前环境无法访问后端服务，请使用模拟数据',
        showCancel: false
      })
    } else {
      wx.showToast({
        title: error.errMsg || error.message || '登录失败，请重试',
        icon: 'none'
      })
    }
  })
```

### 3. 使用测试脚本验证后端服务
创建测试脚本来验证后端服务是否可用：

```javascript
// 在控制台调用测试
const testBackend = require('./test_backend.js')
testBackend.runFullTest()
```

## 测试步骤

### 1. 切换到开发环境
1. 修改 `config.js` 中的 `CURRENT_ENV` 为 `ENV.DEV`
2. 重新编译项目

### 2. 测试模拟数据
1. 在开发者工具中测试登录功能
2. 确认使用模拟数据能正常工作

### 3. 测试后端服务
1. 使用测试脚本验证后端服务
2. 检查网络连接和后端服务状态

### 4. 真机调试
1. 在真机上测试登录功能
2. 查看详细的错误日志

## 域名配置

如果需要在真机上访问后端服务，需要在微信小程序后台配置合法域名：

1. 登录微信公众平台
2. 进入"开发" -> "开发管理" -> "开发设置"
3. 在"服务器域名"中添加后端服务域名
4. 域名必须是 HTTPS 且已备案

## 错误处理优化

### 1. 网络错误处理
```javascript
if (error.message && error.message.includes('网络错误')) {
  wx.showModal({
    title: '网络错误',
    content: '请检查网络连接，或稍后重试',
    showCancel: false
  })
}
```

### 2. 域名限制处理
```javascript
if (error.message && error.message.includes('url not in domain list')) {
  wx.showModal({
    title: '域名限制',
    content: '当前环境无法访问后端服务，请使用模拟数据',
    showCancel: false
  })
}
```

### 3. 通用错误处理
```javascript
wx.showToast({
  title: error.errMsg || error.message || '登录失败，请重试',
  icon: 'none'
})
```

## 注意事项

1. **域名配置**：真机调试需要配置合法域名
2. **环境切换**：开发时可以使用模拟数据
3. **错误处理**：添加详细的错误处理和用户提示
4. **日志输出**：添加详细的日志便于调试

## 相关文件
- `pages/login/login.js` - 登录页面逻辑（已优化错误处理）
- `config.js` - 配置文件（已切换到开发环境）
- `test_backend.js` - 后端服务测试脚本
- `utils/request.js` - 请求工具（包含模拟数据） 