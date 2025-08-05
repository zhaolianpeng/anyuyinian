# 全局云托管调用使用指南

## 概述

现在整个项目已经统一使用标准化的云托管调用方式。所有调用都通过统一的配置和错误处理机制。

## 配置文件

### 1. 全局配置 (`utils/cloud-config.js`)
```javascript
const CLOUD_CONFIG = {
  env: 'prod-5g94mx7a3d07e78c',  // 云托管环境ID
  service: 'golang-lfwy',         // 服务名称
  appId: 'wx101090677bd5219e',    // 小程序AppID
  timeout: 10000,                 // 默认超时时间
  websocketPath: '/ws'           // WebSocket路径
}
```

### 2. 标准调用工具 (`utils/cloud-container-standard.js`)
提供标准化的API调用方法。

## 使用方式

### 1. 直接使用 wx.cloud.callContainer

```javascript
// 基础调用
wx.cloud.callContainer({
  config: {
    env: 'prod-5g94mx7a3d07e78c'
  },
  path: '/api/count',
  method: 'POST',
  data: { action: 'inc' },
  header: {
    'X-WX-SERVICE': 'golang-lfwy',
    'content-type': 'application/json'
  },
  timeout: 10000
}).then(result => {
  console.log('调用成功:', result)
}).catch(error => {
  console.error('调用失败:', error)
})
```

### 2. 使用封装的API

```javascript
const { api } = require('./utils/cloud-container-standard')

// 计数器操作
api.count.increment().then(result => {
  console.log('计数器增加:', result)
}).catch(error => {
  console.error('操作失败:', error)
})

// 获取计数器值
api.count.get().then(result => {
  console.log('当前计数:', result)
}).catch(error => {
  console.error('获取失败:', error)
})
```

### 3. 使用全局配置

```javascript
const { getCallConfig, handleCloudError } = require('./utils/cloud-config')

// 构建请求配置
const requestConfig = getCallConfig('/api/count', 'POST', { action: 'inc' })

// 发送请求
wx.cloud.callContainer(requestConfig)
  .then(result => {
    console.log('调用成功:', result)
  })
  .catch(error => {
    handleCloudError(error)
  })
```

### 4. WebSocket连接

```javascript
const { getWebSocketConfig } = require('./utils/cloud-config')

// 连接WebSocket
const connectConfig = getWebSocketConfig('/ws')
wx.cloud.connectContainer(connectConfig)
  .then(({ socketTask }) => {
    console.log('WebSocket连接成功')
    
    // 设置事件监听
    socketTask.onOpen(() => {
      console.log('连接已打开')
    })
    
    socketTask.onMessage((res) => {
      console.log('收到消息:', res.data)
    })
    
    socketTask.onClose(() => {
      console.log('连接已关闭')
    })
    
    socketTask.onError((error) => {
      console.error('连接错误:', error)
    })
  })
  .catch(error => {
    console.error('WebSocket连接失败:', error)
  })
```

## 错误处理

### 1. 统一错误处理
```javascript
const { handleCloudError } = require('./utils/cloud-config')

try {
  const result = await wx.cloud.callContainer(requestConfig)
  console.log('调用成功:', result)
} catch (error) {
  handleCloudError(error) // 自动显示错误提示
}
```

### 2. 自定义错误处理
```javascript
const { handleCloudError } = require('./utils/cloud-config')

try {
  const result = await wx.cloud.callContainer(requestConfig)
  console.log('调用成功:', result)
} catch (error) {
  const handledError = handleCloudError(error)
  // 可以在这里添加自定义的错误处理逻辑
  console.error('自定义错误处理:', handledError)
}
```

## 重试机制

```javascript
const { callWithRetry } = require('./utils/cloud-config')

// 带重试的API调用
const apiCall = () => wx.cloud.callContainer(requestConfig)
const result = await callWithRetry(apiCall, 3, 1000) // 最多重试3次，间隔1秒
```

## 网络检查

```javascript
const { checkNetworkStatus } = require('./utils/cloud-config')

// 检查网络状态
const isNetworkAvailable = await checkNetworkStatus()
if (!isNetworkAvailable) {
  wx.showToast({
    title: '网络连接不可用',
    icon: 'none'
  })
  return
}

// 继续执行API调用
const result = await wx.cloud.callContainer(requestConfig)
```

## 在页面中使用

### 1. 在页面JS文件中
```javascript
// pages/example/example.js
const { api } = require('../../utils/cloud-container-standard')

Page({
  data: {
    count: 0
  },

  onLoad() {
    this.loadCount()
  },

  async loadCount() {
    try {
      const result = await api.count.get()
      this.setData({ count: result.count })
    } catch (error) {
      console.error('加载计数失败:', error)
    }
  },

  async incrementCount() {
    try {
      const result = await api.count.increment()
      this.setData({ count: result.count })
      wx.showToast({
        title: '计数增加成功',
        icon: 'success'
      })
    } catch (error) {
      console.error('增加计数失败:', error)
    }
  }
})
```

### 2. 在app.js中使用
```javascript
// app.js
const { getCallConfig, handleCloudError } = require('./utils/cloud-config')

App({
  async call(obj) {
    try {
      const requestConfig = getCallConfig(obj.path, obj.method || 'GET', obj.data || {}, obj)
      const result = await wx.cloud.callContainer(requestConfig)
      
      if (result.statusCode === 200) {
        return result.data
      } else {
        throw new Error(`HTTP ${result.statusCode}: ${result.data?.message || '请求失败'}`)
      }
    } catch (error) {
      throw handleCloudError(error)
    }
  }
})
```

## 配置管理

### 1. 修改环境配置
在 `utils/cloud-config.js` 中修改配置：
```javascript
const CLOUD_CONFIG = {
  env: 'your-env-id',        // 修改为您的环境ID
  service: 'your-service',   // 修改为您的服务名
  timeout: 15000,           // 修改超时时间
  // ... 其他配置
}
```

### 2. 添加新的API接口
在 `utils/cloud-container-standard.js` 中添加：
```javascript
const api = {
  // 现有接口...
  
  // 新增接口
  newApi: (params) => callContainer('/api/new', 'POST', params),
  getData: (id) => callContainer(`/api/data/${id}`, 'GET'),
  updateData: (id, data) => callContainer(`/api/data/${id}`, 'PUT', data)
}
```

## 最佳实践

1. **统一使用封装的API**: 使用 `api` 对象中的方法
2. **添加错误处理**: 所有调用都应该有错误处理
3. **使用重试机制**: 对于重要操作使用重试机制
4. **检查网络状态**: 在调用前检查网络连接
5. **记录详细日志**: 便于调试和问题排查

## 调试工具

### 1. 使用测试页面
- 打开 `pages/test/websocket-test`
- 运行详细诊断
- 查看诊断报告

### 2. 查看控制台日志
- 在微信开发者工具中查看控制台
- 关注网络请求和错误信息

### 3. 检查云托管状态
- 登录微信云托管控制台
- 查看服务运行状态和日志

## 总结

现在整个项目已经统一使用标准化的云托管调用方式，具有以下优势：

1. **配置统一**: 所有配置集中管理
2. **错误处理统一**: 统一的错误处理和用户提示
3. **代码复用**: 减少重复代码
4. **易于维护**: 修改配置只需要改一个地方
5. **调试友好**: 详细的日志和错误信息

建议在开发中优先使用封装的API方法，这样可以获得更好的开发体验和错误处理。 