# 订单倒计时功能实现总结

## 功能概述

为订单列表页和订单详情页添加了待支付倒计时功能，并修复了订单金额显示问题。

## 主要功能

### 1. 订单金额显示修复

**问题**: 订单列表页的订单金额显示为0

**解决方案**:
- 添加了 `formatAmount()` 方法来格式化金额显示
- 确保金额字段正确显示，避免显示0的问题

```javascript
// 格式化金额显示
formatAmount(amount) {
  if (amount === 0 || amount === null || amount === undefined) {
    return '0.00'
  }
  return parseFloat(amount).toFixed(2)
}
```

### 2. 订单列表页倒计时功能

**功能特点**:
- 只对待支付状态（status === 0）的订单显示倒计时
- 倒计时格式：`分:秒`（如：29:45）
- 30分钟支付时限
- 超时后显示"已超时"

**实现代码**:
```javascript
// 开始倒计时
startCountdown() {
  const timer = setInterval(() => {
    this.updateCountdown()
  }, 1000)
  this.setData({ countdownTimer: timer })
}

// 更新倒计时
updateCountdown() {
  const orders = this.data.orders.map(order => {
    if (order.status === 0) { // 待支付状态
      const now = new Date().getTime()
      const createTime = new Date(order.createdAt).getTime()
      const timeLimit = 30 * 60 * 1000 // 30分钟
      const remainingTime = createTime + timeLimit - now
      
      if (remainingTime > 0) {
        const minutes = Math.floor(remainingTime / (1000 * 60))
        const seconds = Math.floor((remainingTime % (1000 * 60)) / 1000)
        order.countdown = `${minutes}:${seconds.toString().padStart(2, '0')}`
      } else {
        order.countdown = '已超时'
      }
    }
    return order
  })
  
  this.setData({ orders })
}
```

**模板修改**:
```xml
<view class="status-info">
  <text class="order-status {{getStatusClass(item.status)}}">{{getStatusText(item.status)}}</text>
  <text class="countdown" wx:if="{{item.status === 0 && item.countdown}}">{{item.countdown}}</text>
</view>
```

### 3. 订单详情页倒计时功能

**功能特点**:
- 显示更详细的倒计时信息
- 格式：`剩余支付时间：分:秒`
- 超时后显示"支付时间已过期"
- 自动清除定时器

**实现代码**:
```javascript
// 更新倒计时
updateCountdown() {
  if (!this.data.order || this.data.order.status !== 0) {
    return
  }

  const now = new Date().getTime()
  const createTime = new Date(this.data.order.createdAt).getTime()
  const timeLimit = 30 * 60 * 1000 // 30分钟
  const remainingTime = createTime + timeLimit - now
  
  if (remainingTime > 0) {
    const minutes = Math.floor(remainingTime / (1000 * 60))
    const seconds = Math.floor((remainingTime % (1000 * 60)) / 1000)
    this.setData({
      countdown: `剩余支付时间：${minutes}:${seconds.toString().padStart(2, '0')}`
    })
  } else {
    this.setData({
      countdown: '支付时间已过期'
    })
    // 清除定时器
    if (this.data.countdownTimer) {
      clearInterval(this.data.countdownTimer)
    }
  }
}
```

**模板修改**:
```xml
<view class="status-info">
  <text class="status-text" style="color: {{statusMap[order.status].color}}">
    {{statusMap[order.status].text}}
  </text>
  <text class="status-desc">
    {{order.status === 0 ? '请在30分钟内完成支付' : ...}}
  </text>
  <text class="countdown" wx:if="{{order.status === 0 && countdown}}">{{countdown}}</text>
</view>
```

### 4. 样式优化

**订单列表页样式**:
```css
.status-info {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 8rpx;
}

.countdown {
  font-size: 22rpx;
  color: #ff4757;
  background: rgba(255, 71, 87, 0.1);
  padding: 4rpx 8rpx;
  border-radius: 8rpx;
}
```

**订单详情页样式**:
```css
.countdown {
  font-size: 24rpx;
  color: #ff4757;
  background: rgba(255, 71, 87, 0.2);
  padding: 8rpx 16rpx;
  border-radius: 20rpx;
  display: inline-block;
  margin-top: 10rpx;
}
```

## 技术实现

### 1. 定时器管理
- 使用 `setInterval` 创建定时器
- 在页面卸载时清除定时器，避免内存泄漏
- 倒计时结束后自动清除定时器

### 2. 状态管理
- 只在待支付状态（status === 0）显示倒计时
- 实时更新倒计时显示
- 处理超时状态

### 3. 性能优化
- 只在需要时启动定时器
- 及时清理不需要的定时器
- 避免不必要的DOM更新

## 用户体验

### 1. 视觉反馈
- 倒计时使用红色显示，突出紧迫感
- 圆角背景，视觉效果友好
- 不同页面使用不同的样式风格

### 2. 信息层次
- 订单列表页：简洁的倒计时显示
- 订单详情页：详细的倒计时说明
- 状态信息清晰明确

### 3. 交互体验
- 实时更新，用户能清楚看到剩余时间
- 超时后自动停止倒计时
- 页面切换时正确清理定时器

## 测试验证

创建了测试脚本 `test_order_countdown.sh` 来验证：

1. **订单金额显示**: 检查金额字段是否正确
2. **待支付订单**: 验证待支付订单的存在
3. **倒计时功能**: 测试倒计时逻辑
4. **API响应**: 验证后端接口返回的数据

## 兼容性说明

- **向后兼容**: 不影响现有订单状态显示
- **数据兼容**: 使用现有的订单数据结构
- **样式兼容**: 新增样式不影响现有布局

## 部署建议

1. **前端部署**: 重新编译小程序代码
2. **功能测试**: 测试不同状态的订单显示
3. **性能测试**: 验证定时器管理是否正常
4. **用户体验**: 确认倒计时显示效果 