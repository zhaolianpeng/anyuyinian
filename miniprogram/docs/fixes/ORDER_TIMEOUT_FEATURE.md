# 订单超时功能说明

## 功能概述

订单超时功能实现了未支付订单在30分钟后自动取消的机制，确保订单系统的时效性和资源管理。

## 功能特性

### 1. 自动超时取消
- 订单创建后30分钟内未支付，自动取消
- 后端定时任务每分钟检查超时订单
- 批量处理，提高系统性能

### 2. 前端倒计时显示
- 实时显示支付剩余时间
- 倒计时格式：HH:MM:SS
- 超时后自动更新订单状态

### 3. 支付截止时间显示
- 订单列表显示支付截止时间
- 订单详情页面显示完整截止时间
- 格式化显示，用户友好

## 技术实现

### 后端实现

#### 1. 数据库字段
```sql
-- 添加支付截止时间字段
ALTER TABLE Orders ADD COLUMN payDeadline DATETIME COMMENT '支付截止时间';

-- 创建索引提高查询性能
CREATE INDEX idx_pay_deadline ON Orders(payDeadline);
CREATE INDEX idx_status_pay_deadline ON Orders(status, payDeadline);
```

#### 2. 订单创建时设置截止时间
```go
// 设置支付截止时间（30分钟后）
payDeadline := time.Now().Add(30 * time.Minute)

order := &model.OrderModel{
    // ... 其他字段
    PayDeadline: &payDeadline, // 支付截止时间
}
```

#### 3. 定时任务服务
```go
// 每分钟检查一次超时订单
func (s *OrderTimeoutService) Start() {
    s.ticker = time.NewTicker(1 * time.Minute)
    
    go func() {
        for {
            select {
            case <-s.ticker.C:
                s.checkAndCancelExpiredOrders()
            case <-s.done:
                return
            }
        }
    }()
}
```

#### 4. 批量取消超时订单
```go
func (imp *OrderInterfaceImp) BatchCancelExpiredOrders() (int64, error) {
    cli := db.Get()
    now := time.Now()
    
    result := cli.Table(orderTableName).
        Where("status = ? AND payStatus = ? AND payDeadline < ?", 0, 0, now).
        Updates(map[string]interface{}{
            "status":    3, // 已取消
            "updatedAt": now,
        })
    
    return result.RowsAffected, result.Error
}
```

### 前端实现

#### 1. 倒计时显示
```javascript
updateCountdown() {
    const orders = this.data.orders.map(order => {
        if (order.status === 'pending_pay' && order.payDeadline) {
            const now = new Date().getTime()
            const deadline = new Date(order.payDeadline).getTime()
            const diff = deadline - now
            
            if (diff > 0) {
                const hours = Math.floor(diff / (1000 * 60 * 60))
                const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
                const seconds = Math.floor((diff % (1000 * 60)) / 1000)
                
                order.countdown = `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`
            } else {
                order.countdown = '已过期'
                order.status = 'cancelled'
            }
        }
        return order
    })
    
    this.setData({ orders })
}
```

#### 2. 支付截止时间格式化
```javascript
// 格式化支付截止时间显示
const deadlineDate = new Date(order.payDeadline)
order.formattedPayDeadline = deadlineDate.toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
})
```

## API接口

### 1. 手动检查超时订单
```
POST /api/order/check_expired
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "message": "超时订单检查完成",
    "expiredCount": 2
  }
}
```

### 2. 获取超时订单数量
```
GET /api/order/expired_count
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "expiredCount": 2
  }
}
```

## 配置说明

### 1. 超时时间配置
- 默认超时时间：30分钟
- 可在订单创建时修改：`time.Now().Add(30 * time.Minute)`

### 2. 检查频率配置
- 默认检查频率：每分钟
- 可在服务启动时修改：`time.NewTicker(1 * time.Minute)`

### 3. 数据库索引
- `idx_pay_deadline`：支付截止时间索引
- `idx_status_pay_deadline`：状态和截止时间复合索引

## 测试验证

### 1. 功能测试
```javascript
// 使用测试脚本验证功能
const OrderTimeoutTester = require('./tests/test_order_timeout.js')
const tester = new OrderTimeoutTester()
tester.startTest()
```

### 2. 手动测试
1. 创建测试订单
2. 等待30分钟或修改数据库中的支付截止时间
3. 调用检查接口验证自动取消功能

### 3. 性能测试
- 测试大量订单的批量处理性能
- 验证定时任务的稳定性
- 检查内存使用情况

## 监控和日志

### 1. 服务日志
```
订单超时处理服务已启动
开始检查超时订单...
发现 2 个超时订单
成功取消 2 个超时订单
订单 202401010001 因超时未支付已自动取消
```

### 2. 监控指标
- 超时订单数量
- 取消订单数量
- 处理时间
- 错误率

## 注意事项

### 1. 时区处理
- 确保服务器时区设置正确
- 前端显示时考虑用户时区

### 2. 并发安全
- 使用数据库事务确保数据一致性
- 避免重复处理同一订单

### 3. 性能优化
- 批量处理提高效率
- 合理设置检查频率
- 使用数据库索引优化查询

### 4. 错误处理
- 网络异常时的重试机制
- 数据库连接失败的处理
- 日志记录便于问题排查

## 扩展功能

### 1. 可配置超时时间
- 不同服务类型设置不同超时时间
- 支持动态配置

### 2. 通知功能
- 超时前发送提醒通知
- 支持短信、微信等多种通知方式

### 3. 统计报表
- 超时订单统计
- 取消原因分析
- 用户行为分析

## 总结

订单超时功能通过前后端配合，实现了完整的30分钟超时自动取消机制：

1. **后端**：定时任务 + 批量处理 + 数据库优化
2. **前端**：实时倒计时 + 状态同步 + 用户友好显示
3. **监控**：日志记录 + 性能监控 + 错误处理

该功能确保了订单系统的时效性，提升了用户体验，同时优化了系统资源管理。 