# 视频降级模式最终调试指南

## 当前状态
视频仍然出现 `net::ERR_FAILED` 错误，但降级模式应该已经自动启用。

## 调试步骤

### 1. 检查当前状态
在开发者工具控制台中执行：
```javascript
// 获取当前页面实例
const pages = getCurrentPages()
const currentPage = pages[pages.length - 1]

// 检查视频相关状态
console.log('当前视频状态:', {
  videoUrl: currentPage.data.videoUrl,
  videoPoster: currentPage.data.videoPoster,
  fallbackMode: currentPage.data.fallbackMode,
  videoLoadError: currentPage.data.videoLoadError
})
```

### 2. 手动触发降级模式
如果降级模式没有自动启用，可以手动触发：

#### 方法一：使用调试按钮
- 在视频Banner上点击红色的"强制降级"按钮
- 这会立即启用降级模式

#### 方法二：控制台执行
```javascript
// 获取当前页面实例
const pages = getCurrentPages()
const currentPage = pages[pages.length - 1]

// 强制启用降级模式
currentPage.forceFallbackMode()
```

#### 方法三：手动设置状态
```javascript
// 获取当前页面实例
const pages = getCurrentPages()
const currentPage = pages[pages.length - 1]

// 手动设置降级模式
currentPage.setData({
  fallbackMode: true,
  videoUrl: '',
  videoLoadError: true
})

console.log('手动启用降级模式完成')
```

### 3. 检查降级模式显示
降级模式启用后，应该看到：
- 封面图片全屏显示
- 半透明遮罩层
- "专业护理服务"标题
- "用心呵护每一位长者"副标题
- "重新加载视频"按钮
- "联系客服"按钮

### 4. 测试重试功能
在降级模式下：
1. 点击"重新加载视频"按钮
2. 检查是否重新尝试加载视频
3. 如果仍然失败，应该再次回到降级模式

## 预期行为

### 自动降级触发条件
1. **视频错误事件**：`onVideoError` 被触发
2. **测试视频失败**：`onTestVideoError` 被触发
3. **验证超时**：3秒后自动触发
4. **页面显示检查**：4秒后主动检查

### 降级模式状态
- `fallbackMode: true`
- `videoUrl: ""`
- `videoLoadError: true`
- `videoPoster: "封面图URL"`

## 可能的问题

### 1. 降级模式没有自动触发
**原因**：视频错误事件可能没有正确绑定
**解决**：使用手动触发方法

### 2. 降级模式显示异常
**原因**：封面图片URL可能有问题
**解决**：检查 `videoPoster` 是否正确

### 3. 重试功能不工作
**原因**：事件绑定或方法实现问题
**解决**：检查控制台错误信息

## 调试代码

### 检查视频组件
```javascript
// 检查视频上下文
const videoContext = wx.createVideoContext('homeVideo')
console.log('视频上下文:', videoContext)

// 检查测试视频上下文
const testVideoContext = wx.createVideoContext('testVideo')
console.log('测试视频上下文:', testVideoContext)
```

### 检查网络请求
```javascript
// 检查视频URL是否可访问
wx.request({
  url: 'https://7072-prod-5g94mx7a3d07e78c-1353115175.tcb.qcloud.la/video/光厂_5440331_照顾瘫痪老人按摩理发.mp4',
  method: 'HEAD',
  success: (res) => {
    console.log('视频URL请求成功:', res)
  },
  fail: (err) => {
    console.log('视频URL请求失败:', err)
  }
})
```

### 检查封面图片
```javascript
// 检查封面图片URL
wx.request({
  url: 'https://7072-prod-5g94mx7a3d07e78c-1353115175.tcb.qcloud.la/static/Wechat-IMG36.jpg',
  method: 'HEAD',
  success: (res) => {
    console.log('封面图片请求成功:', res)
  },
  fail: (err) => {
    console.log('封面图片请求失败:', err)
  }
})
```

## 解决方案

### 1. 立即启用降级模式
```javascript
// 获取当前页面实例
const pages = getCurrentPages()
const currentPage = pages[pages.length - 1]

// 强制启用降级模式
currentPage.forceFallbackMode()
```

### 2. 检查页面显示
刷新小程序页面，查看：
1. 是否显示封面图片
2. 是否有操作按钮
3. 是否有错误提示

### 3. 测试功能
1. 点击"重新加载视频"按钮
2. 点击"联系客服"按钮
3. 检查轮播切换是否正常

## 相关文件
- `pages/index/index.js` - 主要逻辑
- `pages/index/index.wxml` - 页面结构
- `pages/index/index.wxss` - 样式定义
- `utils/videoService.js` - 视频服务工具

## 注意事项
1. 调试按钮仅用于测试，正式版本需要移除
2. 控制台日志有助于定位问题
3. 手动触发可以验证功能是否正常
4. 降级模式应该提供完整的用户体验
