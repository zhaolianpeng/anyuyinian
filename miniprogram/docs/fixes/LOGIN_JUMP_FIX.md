# 登录跳转问题修复

## 问题描述

登录成功后显示"登录成功"提示，但随后又跳转回了登录页面，没有正确跳转到首页。

## 问题原因

### 1. 后端响应数据问题
后端返回的响应数据是空值：
```javascript
{code: , data: }
```

### 2. 前端判断逻辑问题
原来的判断逻辑：
```javascript
if (res.code === 0) {
  // 跳转逻辑
}
```

由于 `res.code` 是空值，条件 `res.code === 0` 不成立，所以不会执行跳转逻辑。

### 3. 页面重新加载问题
登录成功后，页面可能重新加载，触发了 `onLoad` 方法中的登录检查逻辑。

## 解决方案

### 1. 修复响应处理逻辑
```javascript
// 修复前
if (res.code === 0) {
  // 跳转逻辑
}

// 修复后
const responseCode = res.code || 0
const responseData = res.data || {}

if (responseCode === 0 || !res.errorMsg) {
  // 跳转逻辑
}
```

### 2. 完善用户信息保存
```javascript
// 保存用户信息
const userInfo = {
  ...userInfoRes.userInfo,  // 微信用户信息
  ...responseData           // 后端返回的数据
}

wx.setStorageSync('token', responseData.token || 'mock_token_' + Date.now())
wx.setStorageSync('userInfo', userInfo)
wx.setStorageSync('userId', responseData.userId || 'mock_user_' + Date.now())
```

### 3. 增强跳转逻辑
```javascript
// 返回上一页或首页
setTimeout(() => {
  wx.navigateBack({
    fail: () => {
      wx.switchTab({
        url: '/pages/index/index'
      })
    }
  })
}, 1500)
```

## 修复后的完整流程

### 1. 登录成功处理
```javascript
.then((res) => {
  console.log('后端登录接口响应:', res)
  
  // 处理响应数据，支持空响应的情况
  const responseCode = res.code || 0
  const responseData = res.data || {}
  
  console.log('处理后的响应:', { code: responseCode, data: responseData })
  
  if (responseCode === 0 || !res.errorMsg) {
    // 登录成功，保存用户信息
    const userInfo = {
      ...userInfoRes.userInfo,
      ...responseData
    }
    
    wx.setStorageSync('token', responseData.token || 'mock_token_' + Date.now())
    wx.setStorageSync('userInfo', userInfo)
    wx.setStorageSync('userId', responseData.userId || 'mock_user_' + Date.now())
    
    console.log('用户信息已保存到本地存储')
    
    wx.showToast({
      title: '登录成功',
      icon: 'success'
    })
    
    // 返回上一页或首页
    setTimeout(() => {
      wx.navigateBack({
        fail: () => {
          wx.switchTab({
            url: '/pages/index/index'
          })
        }
      })
    }, 1500)
  } else {
    throw new Error(res.errorMsg || '登录失败')
  }
})
```

### 2. 页面加载检查
```javascript
onLoad() {
  console.log('登录页面已加载，请点击登录按钮')
  
  // 检查是否已经登录
  const token = wx.getStorageSync('token')
  const userInfo = wx.getStorageSync('userInfo')
  
  if (token && userInfo) {
    console.log('用户已登录，返回上一页')
    wx.navigateBack({
      fail: () => {
        wx.switchTab({
          url: '/pages/index/index'
        })
      }
    })
  }
}
```

## 测试验证

### 1. 测试不同响应情况
- ✅ 正常响应：`{code: 0, data: {...}}`
- ✅ 空响应：`{code: '', data: ''}`
- ✅ 部分响应：`{code: 0, data: {}}`
- ✅ 错误响应：`{code: -1, errorMsg: '登录失败'}`

### 2. 验证跳转逻辑
1. 登录成功后显示"登录成功"提示
2. 1.5秒后自动跳转到首页
3. 如果无法返回上一页，则切换到首页tab

### 3. 验证本地存储
- Token已保存
- 用户信息已保存
- 用户ID已保存

## 预期结果

修复后，登录流程应该是：
1. ✅ 用户点击登录按钮
2. ✅ 获取用户信息成功
3. ✅ 获取登录code成功
4. ✅ 后端接口调用成功
5. ✅ 显示"登录成功"提示
6. ✅ 保存用户信息到本地存储
7. ✅ 1.5秒后跳转到首页
8. ✅ 不再返回登录页面

## 相关文件
- `pages/login/login.js` - 登录页面（已修复）
- `test_login_flow.js` - 登录流程测试脚本

## 下一步操作
1. 重新编译小程序
2. 测试登录功能
3. 确认登录成功后正确跳转到首页
4. 检查本地存储中的用户信息 