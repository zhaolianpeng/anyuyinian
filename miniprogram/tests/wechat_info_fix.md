# 微信用户信息获取修复

## 问题描述
完善资料页面点击"获取微信信息"按钮后，获取到的昵称和手机号都是错误的：
- 昵称显示为"微信用户"（默认昵称）
- 手机号显示为"13800138000"（模拟手机号）

## 问题原因
1. **API不可用**: 在开发环境中，`wx.getUserProfile` 和 `wx.getPhoneNumber` API可能不可用
2. **使用模拟数据**: 代码中使用了硬编码的模拟数据
3. **手机号获取方式错误**: 手机号需要通过button的`open-type="getPhoneNumber"`获取，而不是直接调用API

## 修复方案

### 1. 优化微信用户信息获取
**文件**: `pages/user/setup-profile.js`

**修复前**:
```javascript
// 直接使用模拟数据
const mockUserInfo = {
  nickName: '微信用户',
  avatarUrl: 'https://...',
  // ...
}
resolve(mockUserInfo)
```

**修复后**:
```javascript
// 尝试多种方式获取用户信息
if (typeof wx.getUserProfile !== 'function') {
  // 尝试使用wx.getUserInfo作为备选方案
  if (typeof wx.getUserInfo === 'function') {
    wx.getUserInfo({
      success: (res) => {
        resolve(res.userInfo)
      },
      fail: (err) => {
        // 提供默认数据
        resolve(defaultUserInfo)
      }
    })
  }
}
```

### 2. 修复手机号获取方式
**文件**: `pages/user/setup-profile.js`

**修复前**:
```javascript
// 直接调用wx.getPhoneNumber API
wx.getPhoneNumber({
  success: (res) => {
    resolve('13800138000') // 硬编码模拟数据
  }
})
```

**修复后**:
```javascript
// 通过button的open-type获取
getWxPhone() {
  return new Promise((resolve, reject) => {
    // 手机号需要通过button的open-type="getPhoneNumber"触发
    console.log('手机号需要通过按钮点击获取')
    resolve('')
  })
}

// 处理手机号获取按钮点击
onGetPhoneNumber(e) {
  if (e.detail.errMsg === 'getPhoneNumber:ok') {
    const { encryptedData, iv } = e.detail
    if (encryptedData && iv) {
      this.decryptPhoneNumber(encryptedData, iv)
    }
  }
}
```

### 3. 添加手机号解密功能
**文件**: `pages/user/setup-profile.js`

**新增功能**:
```javascript
// 解密手机号
async decryptPhoneNumber(encryptedData, iv) {
  try {
    const result = await api.decryptPhoneNumber({
      userId: userId,
      encryptedData: encryptedData,
      iv: iv
    })

    if (result.code === 0 && result.data && result.data.phoneNumber) {
      const phoneNumber = result.data.phoneNumber
      this.setData({ phone: phoneNumber })
    }
  } catch (error) {
    console.error('解密手机号失败:', error)
  }
}
```

### 4. 更新WXML界面
**文件**: `pages/user/setup-profile.wxml`

**修复前**:
```xml
<text class="value">{{phone || '未获取到'}}</text>
<button class="phone-btn" bindtap="onShowPhoneModal" wx:if="{{!phone}}">
  手动输入
</button>
```

**修复后**:
```xml
<text class="value">{{phone || '未获取到'}}</text>
<button 
  class="phone-btn" 
  open-type="getPhoneNumber" 
  bindgetphonenumber="onGetPhoneNumber"
  wx:if="{{!phone}}"
>
  获取手机号
</button>
<button class="phone-btn" bindtap="onShowPhoneModal" wx:if="{{!phone}}">
  手动输入
</button>
```

### 5. 添加API接口
**文件**: `utils/cloud-container-standard.js`

**新增接口**:
```javascript
// 用户管理
decryptPhoneNumber: (data) => callContainer('/api/user/decrypt_phone', 'POST', data),
```

## 修复效果

### 修复前
- 昵称: "微信用户"（硬编码）
- 手机号: "13800138000"（硬编码）
- 无法获取真实用户信息

### 修复后
- 昵称: 尝试获取真实微信昵称，失败时显示默认值
- 手机号: 通过按钮点击获取真实手机号，支持手动输入备选方案
- 提供更好的用户体验和错误处理

## 测试要点

1. **用户信息获取**:
   - 在真实微信环境中测试获取用户昵称
   - 在开发环境中测试降级处理

2. **手机号获取**:
   - 测试"获取手机号"按钮功能
   - 测试"手动输入"备选方案
   - 测试手机号解密功能

3. **错误处理**:
   - 测试API不可用时的降级处理
   - 测试用户拒绝授权时的提示
   - 测试网络错误时的处理

## 注意事项

1. **微信API限制**: 某些API在开发环境中可能不可用
2. **用户授权**: 需要用户主动授权才能获取信息
3. **后端支持**: 需要后端提供手机号解密接口
4. **隐私保护**: 确保用户信息的安全处理

## 后续优化

1. **添加更多降级方案**: 提供更多获取用户信息的方式
2. **改进错误提示**: 提供更友好的错误信息
3. **缓存机制**: 缓存已获取的用户信息
4. **数据验证**: 添加用户输入数据的验证
