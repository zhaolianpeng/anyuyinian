# 设置页面调试指南

## 问题描述

用户反馈点击"完善资料"按钮没有反应。

## 可能的原因

1. **页面未注册**：`setup-profile` 页面可能未在 `app.json` 中注册
2. **路由未配置**：后端可能缺少更新用户信息的路由
3. **JavaScript错误**：页面代码可能存在语法错误
4. **权限问题**：微信API权限可能未正确配置

## 已修复的问题

### 1. 页面注册问题 ✅
- 已将 `pages/user/setup-profile` 添加到 `app.json` 的 pages 数组中

### 2. 后端路由问题 ✅
- 已在 `main.go` 中添加了 `/api/user/update_info` 路由
- 路由指向 `service.UpdateUserInfoHandler`

### 3. 调试信息增强 ✅
- 在 `profile.js` 中添加了详细的调试日志
- 添加了页面跳转的成功/失败回调

### 4. 微信API兼容性问题 ✅
- 修复了 `wx.getPhoneNumber is not a function` 错误
- 在开发环境中提供模拟数据
- 添加了API可用性检查
- 确保在微信开发者工具中可以正常测试

## 调试步骤

### 1. 检查控制台日志
在用户资料页面点击"完善资料"按钮，查看控制台输出：
```
点击完善资料按钮
准备跳转到设置页面
跳转成功/跳转失败
```

### 2. 使用调试页面
访问调试页面：`/pages/test/setup-profile-debug`
- 测试微信API可用性
- 测试后端API连接
- 查看详细调试信息

### 3. 检查用户登录状态
确保用户已登录，`userId` 存在于本地存储中。

## 测试方法

### 1. 基本功能测试
```javascript
// 在用户资料页面控制台运行
const { testPageNavigation } = require('./tests/test_page_navigation.js')
testPageNavigation()
```

### 2. 页面跳转测试
```javascript
// 测试页面跳转
wx.navigateTo({
  url: '/pages/user/setup-profile',
  success: () => console.log('跳转成功'),
  fail: (error) => console.error('跳转失败:', error)
})
```

### 3. API测试
```javascript
// 测试更新用户信息API
const { api } = require('../../utils/cloud-container-standard')
const result = await api.updateUserInfo({
  userId: wx.getStorageSync('userId'),
  nickName: '测试昵称'
})
console.log('API结果:', result)
```

## 常见问题解决

### 1. 页面跳转失败
**症状**：点击按钮后没有任何反应
**解决**：
- 检查 `app.json` 中是否包含 `pages/user/setup-profile`
- 检查页面文件是否存在且语法正确
- 查看控制台是否有JavaScript错误

### 2. API调用失败
**症状**：页面可以跳转，但获取微信信息失败
**解决**：
- 检查微信开发者工具是否开启了"不校验合法域名"
- 检查后端服务是否正常运行
- 检查网络连接是否正常

### 3. 微信API不可用
**症状**：`wx.getUserProfile` 或 `wx.getPhoneNumber` 不可用
**解决**：
- 检查微信开发者工具版本
- 检查小程序基础库版本
- 检查是否在真机上测试
- **已修复**：在开发环境中提供模拟数据，确保功能可以正常测试

## 验证清单

- [ ] `app.json` 中包含 `pages/user/setup-profile`
- [ ] 后端 `main.go` 中包含 `/api/user/update_info` 路由
- [ ] 用户已登录，`userId` 存在
- [ ] 微信开发者工具控制台无错误
- [ ] 网络连接正常
- [ ] 后端服务正常运行

## 下一步操作

1. **重新编译小程序**：在微信开发者工具中重新编译
2. **清除缓存**：清除小程序缓存后重新测试
3. **真机测试**：在真机上测试功能
4. **查看日志**：检查后端日志是否有错误信息

## 联系支持

如果问题仍然存在，请提供以下信息：
1. 微信开发者工具版本
2. 小程序基础库版本
3. 控制台错误信息
4. 网络请求日志
5. 后端服务日志 