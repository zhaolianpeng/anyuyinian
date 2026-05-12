# 视频降级模式调试指南

## 当前问题
视频播放失败，出现 `net::ERR_FAILED` 错误，但降级模式可能没有正确触发。

## 调试步骤

### 1. 检查控制台日志
在开发者工具控制台中查看以下日志：
- `视频播放错误:` - 视频错误事件是否触发
- `触发视频错误处理，准备启用降级模式` - 错误处理是否执行
- `启用降级模式，只显示封面图` - 降级模式是否启用
- `当前状态:` - 检查各个状态值

### 2. 检查页面状态
在控制台中执行以下代码检查当前状态：
```javascript
// 获取当前页面实例
const pages = getCurrentPages()
const currentPage = pages[pages.length - 1]

// 检查视频相关状态
console.log('视频状态检查:', {
  videoUrl: currentPage.data.videoUrl,
  videoPoster: currentPage.data.videoPoster,
  fallbackMode: currentPage.data.fallbackMode,
  videoLoadError: currentPage.data.videoLoadError
})
```

### 3. 手动触发降级模式
如果自动降级没有触发，可以手动触发：
```javascript
// 获取当前页面实例
const pages = getCurrentPages()
const currentPage = pages[pages.length - 1]

// 手动触发降级模式
currentPage.triggerFallbackMode()
```

### 4. 使用调试按钮
在视频Banner上有一个红色的"测试降级模式"按钮，点击可以手动触发降级模式。

## 预期行为

### 正常情况
1. 视频开始加载
2. 视频加载失败，触发错误事件
3. 自动启用降级模式
4. 显示封面图片和操作按钮

### 降级模式显示
- 封面图片全屏显示
- 半透明遮罩层
- "专业护理服务"标题
- "用心呵护每一位长者"副标题
- "重新加载视频"按钮
- "联系客服"按钮

## 可能的问题

### 1. 视频错误事件未触发
**原因**：视频组件可能没有正确绑定错误事件
**解决**：检查WXML中的 `binderror="onVideoError"` 是否正确

### 2. 降级模式状态未更新
**原因**：setData可能没有正确执行
**解决**：检查控制台日志，确认状态更新

### 3. 封面图片未显示
**原因**：videoPoster为空或图片加载失败
**解决**：检查封面图片URL是否正确

## 调试代码

### 检查视频组件状态
```javascript
// 在控制台执行
const videoContext = wx.createVideoContext('homeVideo')
console.log('视频上下文:', videoContext)
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

### 强制启用降级模式
```javascript
// 获取当前页面实例
const pages = getCurrentPages()
const currentPage = pages[pages.length - 1]

// 强制设置降级模式
currentPage.setData({
  fallbackMode: true,
  videoUrl: '',
  videoLoadError: true
})

console.log('强制启用降级模式完成')
```

## 解决方案

### 1. 如果降级模式没有自动触发
- 检查视频错误事件绑定
- 手动触发降级模式
- 检查状态更新逻辑

### 2. 如果降级模式显示异常
- 检查封面图片URL
- 检查样式是否正确应用
- 检查条件渲染逻辑

### 3. 如果功能按钮不工作
- 检查事件绑定
- 检查方法实现
- 检查状态管理

## 测试验证

### 1. 自动降级测试
1. 刷新页面
2. 等待视频加载失败
3. 检查是否自动切换到降级模式

### 2. 手动降级测试
1. 点击"测试降级模式"按钮
2. 检查是否切换到降级模式
3. 检查降级模式UI是否正常

### 3. 重试功能测试
1. 在降级模式下点击"重新加载视频"
2. 检查是否重新尝试加载视频
3. 检查状态是否正确更新

## 相关文件
- `pages/index/index.js` - 主要逻辑
- `pages/index/index.wxml` - 页面结构
- `pages/index/index.wxss` - 样式定义
- `utils/videoService.js` - 视频服务工具

## 注意事项
1. 调试按钮仅用于测试，正式版本需要移除
2. 控制台日志有助于定位问题
3. 状态检查是调试的关键步骤
4. 手动触发可以验证功能是否正常
