# 小程序分享功能修复

## 问题描述
小程序无法直接分享，用户无法通过右上角菜单或按钮分享小程序给好友或朋友圈。

## 修复方案

### 1. 添加全局分享配置

#### app.json 配置
```json
{
  "share": {
    "title": "安语颐年护理陪诊",
    "desc": "专业护理陪诊服务，让您就医更安心",
    "path": "pages/index/index"
  }
}
```

#### app.js 全局分享方法
```javascript
// 全局分享配置
onShareAppMessage(res) {
  console.log('全局分享被触发:', res)
  
  // 根据来源页面返回不同的分享内容
  if (res.from === 'button') {
    // 来自按钮点击
    return {
      title: res.target.dataset.title || '安语颐年护理陪诊',
      desc: res.target.dataset.desc || '专业护理陪诊服务，让您就医更安心',
      path: res.target.dataset.path || 'pages/index/index',
      imageUrl: res.target.dataset.imageUrl || ''
    }
  } else {
    // 来自右上角菜单
    return {
      title: '安语颐年护理陪诊',
      desc: '专业护理陪诊服务，让您就医更安心',
      path: 'pages/index/index'
    }
  }
},

// 分享到朋友圈
onShareTimeline(res) {
  console.log('分享到朋友圈被触发:', res)
  
  return {
    title: '安语颐年护理陪诊 - 专业护理陪诊服务',
    query: '',
    imageUrl: ''
  }
}
```

### 2. 页面级分享功能

#### 首页分享 (pages/index/index.js)
```javascript
// 分享给好友
onShareAppMessage(res) {
  console.log('首页分享被触发:', res)
  
  if (res.from === 'button') {
    // 来自按钮点击
    return {
      title: res.target.dataset.title || '安语颐年护理陪诊',
      desc: res.target.dataset.desc || '专业护理陪诊服务，让您就医更安心',
      path: res.target.dataset.path || 'pages/index/index',
      imageUrl: res.target.dataset.imageUrl || ''
    }
  } else {
    // 来自右上角菜单
    return {
      title: '安语颐年护理陪诊',
      desc: '专业护理陪诊服务，让您就医更安心',
      path: 'pages/index/index'
    }
  }
},

// 分享到朋友圈
onShareTimeline(res) {
  console.log('首页分享到朋友圈被触发:', res)
  
  return {
    title: '安语颐年护理陪诊 - 专业护理陪诊服务',
    query: '',
    imageUrl: ''
  }
}
```

#### 订单列表页面分享 (pages/order/list.js)
```javascript
// 分享给好友
onShareAppMessage(res) {
  console.log('订单列表分享被触发:', res)
  
  return {
    title: '安语颐年护理陪诊 - 我的订单',
    desc: '查看我的护理陪诊订单',
    path: 'pages/order/list'
  }
},

// 分享到朋友圈
onShareTimeline(res) {
  console.log('订单列表分享到朋友圈被触发:', res)
  
  return {
    title: '安语颐年护理陪诊 - 我的订单',
    query: '',
    imageUrl: ''
  }
}
```

#### 个人中心页面分享 (pages/user/profile.js)
```javascript
// 分享给好友
onShareAppMessage(res) {
  console.log('个人中心分享被触发:', res)
  
  return {
    title: '安语颐年护理陪诊 - 个人中心',
    desc: '专业护理陪诊服务，让您就医更安心',
    path: 'pages/index/index'
  }
},

// 分享到朋友圈
onShareTimeline(res) {
  console.log('个人中心分享到朋友圈被触发:', res)
  
  return {
    title: '安语颐年护理陪诊 - 专业护理陪诊服务',
    query: '',
    imageUrl: ''
  }
}
```

## 分享功能说明

### 1. 分享方式
- **右上角菜单分享**: 用户点击右上角"..."菜单，选择"转发"或"分享到朋友圈"
- **按钮分享**: 页面中的分享按钮，通过 `open-type="share"` 触发

### 2. 分享内容
- **标题**: 根据页面不同显示相应的标题
- **描述**: 简要介绍小程序功能
- **路径**: 分享后用户点击进入的页面路径
- **图片**: 可选的分享图片

### 3. 分享目标
- **微信好友**: 通过 `onShareAppMessage` 方法处理
- **朋友圈**: 通过 `onShareTimeline` 方法处理

## 已添加分享功能的页面

1. **首页** (`pages/index/index`) - 主要入口页面
2. **订单列表** (`pages/order/list`) - 用户订单管理
3. **个人中心** (`pages/user/profile`) - 用户个人信息
4. **服务详情** (`pages/service/detail`) - 已有分享功能
5. **服务列表** (`pages/service/list`) - 已有分享功能
6. **其他页面** - 已有分享功能

## 使用方法

### 1. 右上角菜单分享
用户在任何页面点击右上角"..."菜单，选择"转发"即可分享给好友。

### 2. 按钮分享
在页面中添加分享按钮：
```xml
<button open-type="share" 
        data-title="自定义标题" 
        data-desc="自定义描述" 
        data-path="pages/target/page">
  分享给好友
</button>
```

### 3. 分享到朋友圈
用户点击右上角"..."菜单，选择"分享到朋友圈"即可。

## 测试建议

1. **右上角菜单测试**
   - 在各个页面点击右上角"..."菜单
   - 验证"转发"和"分享到朋友圈"选项是否显示
   - 测试分享内容是否正确

2. **按钮分享测试**
   - 测试页面中的分享按钮
   - 验证自定义分享内容是否正确

3. **分享接收测试**
   - 分享给其他用户
   - 验证分享链接是否能正常打开
   - 检查分享内容显示是否正确

## 注意事项

1. **分享权限**: 确保小程序已通过微信审核，支持分享功能
2. **分享内容**: 分享内容应符合微信小程序规范
3. **路径正确性**: 确保分享路径指向正确的页面
4. **用户体验**: 分享内容应简洁明了，吸引用户点击

## 相关文件

- `app.json` - 全局分享配置
- `app.js` - 全局分享方法
- `pages/index/index.js` - 首页分享功能
- `pages/order/list.js` - 订单列表分享功能
- `pages/user/profile.js` - 个人中心分享功能
