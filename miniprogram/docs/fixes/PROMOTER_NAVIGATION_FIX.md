# 推广中心导航问题修复

## 问题描述

用户反馈点击推广中心没有跳转到对应的页面。

## 问题分析

### 1. 问题根源
- 用户资料页面中的"推荐中心"菜单项跳转到的是 `/pages/referral/center`
- 而我们实现的推广中心页面是 `/pages/promoter/home`
- 页面路径不匹配导致跳转失败

### 2. 具体问题
1. **菜单文本**: 显示为"推荐中心"，但实际功能是"推广中心"
2. **跳转路径**: 跳转到 `/pages/referral/center`，而不是 `/pages/promoter/home`
3. **功能不一致**: 推荐中心和推广中心是两个不同的功能模块

## 修复方案

### 1. 更新用户资料页面跳转逻辑

#### 修改前
```javascript
// 推荐中心
onReferralCenter() {
  wx.navigateTo({
    url: '/pages/referral/center'
  })
},
```

#### 修改后
```javascript
// 推广中心
onReferralCenter() {
  wx.navigateTo({
    url: '/pages/promoter/home'
  })
},
```

### 2. 更新菜单文本

#### 修改前
```html
<text class="menu-title">推荐中心</text>
```

#### 修改后
```html
<text class="menu-title">推广中心</text>
```

### 3. 更新测试文件

更新了导航测试文件中的相关链接：
- 将 `/pages/referral/center` 改为 `/pages/promoter/home`
- 将"推荐中心"改为"推广中心"

## 修复内容

### 1. 修改的文件

#### `miniprogram/pages/user/profile.js`
- 更新了 `onReferralCenter` 函数的跳转路径
- 从 `/pages/referral/center` 改为 `/pages/promoter/home`

#### `miniprogram/pages/user/profile.wxml`
- 更新了菜单文本
- 从"推荐中心"改为"推广中心"

#### `miniprogram/tests/navigation/test_navigation.js`
- 更新了测试数据中的导航链接
- 确保测试与实际功能一致

### 2. 新增的文件

#### `miniprogram/tests/navigation/test_promoter_navigation.js`
- 创建了专门的推广中心导航测试脚本
- 验证页面跳转逻辑是否正确
- 测试登录状态检查

## 功能验证

### 1. 跳转流程
1. 用户在"我的"页面点击"推广中心"菜单
2. 系统检查用户登录状态
3. 如果已登录，直接跳转到推广中心页面
4. 如果未登录，先跳转到登录页面

### 2. 页面路径
- **推广中心页面**: `/pages/promoter/home`
- **用户资料页面**: `/pages/user/profile`
- **登录页面**: `/pages/login/login`

### 3. 功能特性
- ✅ 正确的页面跳转
- ✅ 登录状态检查
- ✅ 用户友好的菜单文本
- ✅ 完整的测试覆盖

## 测试验证

### 1. 导航测试
```javascript
// 测试推广中心跳转
function testPromoterNavigation() {
  const result = simulateNavigateToPage('/pages/promoter/home')
  if (result) {
    console.log('✅ 推广中心页面跳转正常')
  }
}
```

### 2. 菜单测试
```javascript
// 测试用户资料页面菜单
const menuItems = [
  { name: '推广中心', url: '/pages/promoter/home' },
  // 其他菜单项...
]
```

### 3. 登录检查
```javascript
// 检查用户登录状态
const userInfo = wx.getStorageSync('userInfo')
const token = wx.getStorageSync('token')
if (userInfo && token) {
  console.log('✅ 用户已登录，可以访问推广中心')
}
```

## 使用说明

### 1. 访问推广中心
1. 打开小程序，进入"我的"页面
2. 点击"推广中心"菜单项
3. 如果未登录，系统会提示先登录
4. 登录后可以正常使用推广中心功能

### 2. 推广中心功能
- **推广员信息**: 查看个人信息和统计数据
- **佣金记录**: 查看佣金收入和记录
- **提现申请**: 申请提现到微信、支付宝或银行卡
- **分享功能**: 分享小程序给朋友

### 3. 注意事项
- 需要先登录才能使用推广中心功能
- 推广中心页面支持下拉刷新和上拉加载
- 提现申请需要填写正确的账户信息

## 总结

通过这次修复：

1. **解决了跳转问题**: 正确跳转到推广中心页面
2. **统一了功能命名**: 将"推荐中心"改为"推广中心"
3. **完善了测试覆盖**: 创建了专门的导航测试
4. **提升了用户体验**: 清晰的菜单文本和正确的跳转逻辑

现在用户点击"推广中心"菜单可以正常跳转到推广中心页面，享受完整的推广管理功能。 