# WebSocket业务集成总结

## 概述

本文档总结了WebSocket功能在各个业务场景中的集成情况，包括订单管理、客服聊天、系统通知等核心功能。

## 集成场景

### 1. 订单管理场景

#### 订单列表页面 (`/pages/order/list`)
**功能特性：**
- 实时接收订单状态更新
- 自动刷新订单列表
- 显示订单状态变更通知
- 支持手动切换实时更新

**WebSocket消息类型：**
```javascript
// 订单更新消息
{
  type: 'order_update',
  orderId: 'ORDER_123',
  orderNo: '202401010001',
  status: 'paid',
  amount: 199.99,
  timestamp: 1640995200000
}
```

**主要方法：**
- `handleOrderUpdate()` - 处理订单更新
- `sendOrderStatusChange()` - 发送订单状态变更
- `showOrderUpdateNotification()` - 显示更新通知

#### 订单详情页面 (`/pages/order/detail`)
**功能特性：**
- 实时更新订单状态
- 自动停止倒计时（支付成功后）
- 显示订单状态变更通知
- 支持手动切换实时更新

**主要方法：**
- `handleOrderUpdate()` - 处理订单更新
- `stopCountdown()` - 停止倒计时
- `sendOrderStatusChange()` - 发送订单状态变更

### 2. 客服聊天场景

#### 客服聊天页面 (`/pages/kefu/chat`)
**功能特性：**
- 实时消息收发
- 输入状态同步
- 系统通知处理
- 支持图片消息

**WebSocket消息类型：**
```javascript
// 聊天消息
{
  type: 'chat_message',
  sender: 'user', // 'user' | 'kefu'
  content: '消息内容',
  messageType: 'text', // 'text' | 'image'
  timestamp: 1640995200000
}

// 输入状态
{
  type: 'typing_status',
  sender: 'user',
  isTyping: true
}

// 系统通知
{
  type: 'system_notification',
  type: 'kefu_online', // 'kefu_online' | 'kefu_offline' | 'session_timeout'
  content: '客服已上线，为您服务'
}
```

**主要方法：**
- `handleChatMessage()` - 处理聊天消息
- `handleTypingStatus()` - 处理输入状态
- `handleSystemNotification()` - 处理系统通知
- `sendChatMessage()` - 发送聊天消息
- `sendTypingStatus()` - 发送输入状态

### 3. 系统通知场景

#### 首页 (`/pages/index/index`)
**功能特性：**
- 接收系统通知
- 显示维护通知
- 显示更新通知
- 显示促销通知
- 订单状态更新通知

**WebSocket消息类型：**
```javascript
// 通知消息
{
  type: 'notification',
  title: '系统通知',
  content: '通知内容',
  type: 'info', // 'info' | 'warning' | 'error' | 'success'
  timestamp: 1640995200000
}

// 系统消息
{
  type: 'system_message',
  type: 'maintenance', // 'maintenance' | 'update' | 'promotion'
  content: '系统维护通知',
  action: 'restart'
}
```

**主要方法：**
- `handleNotification()` - 处理通知消息
- `handleSystemMessage()` - 处理系统消息
- `handleOrderUpdate()` - 处理订单更新
- `showMaintenanceNotice()` - 显示维护通知
- `showUpdateNotice()` - 显示更新通知
- `showPromotionNotice()` - 显示促销通知

## 消息处理机制

### 1. 消息处理器注册
```javascript
// 在页面初始化时注册消息处理器
websocketManager.setMessageHandler('order_update', this.handleOrderUpdate.bind(this))
websocketManager.setMessageHandler('chat_message', this.handleChatMessage.bind(this))
websocketManager.setMessageHandler('notification', this.handleNotification.bind(this))
```

### 2. 消息处理器清理
```javascript
// 在页面卸载时清理消息处理器
onUnload() {
  websocketManager.removeMessageHandler('order_update')
  websocketManager.removeMessageHandler('chat_message')
  websocketManager.removeMessageHandler('notification')
}
```

### 3. 连接状态检查
```javascript
// 检查WebSocket连接状态
checkWebSocketStatus() {
  const isConnected = websocketManager.isConnected()
  this.setData({
    websocketConnected: isConnected
  })
  
  if (!isConnected) {
    websocketManager.reconnect()
  }
}
```

## 业务场景应用

### 1. 订单状态实时更新
**应用场景：**
- 用户支付成功后，订单状态立即更新
- 订单取消或退款后，状态实时同步
- 支付超时后，自动更新订单状态

**实现方式：**
```javascript
// 接收订单更新消息
handleOrderUpdate(message) {
  const { orderId, orderNo, status, amount } = message
  
  // 更新本地订单数据
  const orders = [...this.data.orders]
  const orderIndex = orders.findIndex(order => 
    order.id === orderId || order.orderNo === orderNo
  )
  
  if (orderIndex !== -1) {
    orders[orderIndex].status = status
    this.setData({ orders })
    
    // 显示更新通知
    this.showOrderUpdateNotification(message)
  }
}
```

### 2. 实时客服聊天
**应用场景：**
- 用户与客服实时对话
- 显示客服输入状态
- 系统消息推送

**实现方式：**
```javascript
// 发送聊天消息
sendChatMessage(content, messageType = 'text') {
  if (websocketManager.isConnected()) {
    const messageData = {
      type: 'chat_message',
      sender: 'user',
      content: content,
      messageType: messageType,
      timestamp: Date.now()
    }
    
    websocketManager.send(messageData, 'chat_message')
  }
}

// 处理接收到的消息
handleChatMessage(message) {
  const { sender, content, messageType, timestamp } = message
  
  this.addMessage({
    type: sender === 'user' ? 1 : 0,
    content: content,
    messageType: messageType,
    timestamp: timestamp
  })
}
```

### 3. 系统通知推送
**应用场景：**
- 系统维护通知
- 版本更新提醒
- 促销活动推送
- 订单状态变更通知

**实现方式：**
```javascript
// 处理系统消息
handleSystemMessage(message) {
  const { type, content } = message
  
  switch (type) {
    case 'maintenance':
      this.showMaintenanceNotice(content)
      break
    case 'update':
      this.showUpdateNotice(content)
      break
    case 'promotion':
      this.showPromotionNotice(content)
      break
  }
}
```

## 错误处理机制

### 1. 连接失败处理
```javascript
// 自动重连机制
if (!isConnected) {
  console.log('WebSocket未连接，尝试重连...')
  websocketManager.reconnect()
}
```

### 2. 消息发送失败处理
```javascript
// 发送消息时的错误处理
sendMessage() {
  if (!websocketManager.isConnected()) {
    wx.showToast({
      title: 'WebSocket未连接',
      icon: 'none'
    })
    return
  }
  
  // 发送消息逻辑
}
```

### 3. 降级处理
```javascript
// 当WebSocket不可用时，使用模拟数据
if (!websocketManager.isConnected()) {
  this.simulateKefuReply(content)
}
```

## 性能优化

### 1. 消息去重
```javascript
// 避免重复处理相同消息
handleOrderUpdate(message) {
  const messageId = `${message.orderId}_${message.status}_${message.timestamp}`
  if (this.processedMessages.has(messageId)) {
    return
  }
  this.processedMessages.add(messageId)
  
  // 处理消息逻辑
}
```

### 2. 批量更新
```javascript
// 批量更新UI，减少重绘次数
updateOrders(updates) {
  const orders = [...this.data.orders]
  
  updates.forEach(update => {
    const index = orders.findIndex(order => order.id === update.orderId)
    if (index !== -1) {
      orders[index] = { ...orders[index], ...update }
    }
  })
  
  this.setData({ orders })
}
```

### 3. 内存管理
```javascript
// 限制通知数量，避免内存泄漏
handleNotification(message) {
  const notifications = [...this.data.notifications]
  notifications.unshift(newNotification)
  
  // 限制通知数量
  if (notifications.length > 10) {
    notifications.splice(10)
  }
  
  this.setData({ notifications })
}
```

## 测试验证

### 1. 功能测试
- 使用测试页面验证WebSocket连接
- 测试各种消息类型的发送和接收
- 验证错误处理机制

### 2. 性能测试
- 测试大量消息的处理性能
- 验证内存使用情况
- 测试网络异常时的表现

### 3. 用户体验测试
- 测试实时更新的响应速度
- 验证通知显示的及时性
- 测试不同网络环境下的表现

## 总结

WebSocket功能已成功集成到以下业务场景：

1. **订单管理** - 实时状态更新，提升用户体验
2. **客服聊天** - 实时消息收发，支持多种消息类型
3. **系统通知** - 及时推送重要信息，增强用户粘性

通过合理的错误处理、性能优化和降级机制，确保了WebSocket功能的稳定性和可靠性。用户可以在各个页面享受到实时、流畅的交互体验。 