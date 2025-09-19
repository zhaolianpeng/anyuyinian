# 订单状态显示修复

## 问题描述
订单管理页面中，已退款的订单状态显示为"未知"，而不是正确的"已退款"。

## 问题根因
`admin_service.go` 中的 `getOrderStatusText` 函数状态映射不正确：

### 错误的映射（修复前）
```go
func getOrderStatusText(status int) string {
    switch status {
    case 0: return "待支付"
    case 1: return "已支付"
    case 2: return "已取消"  // ❌ 错误
    case 3: return "已完成"  // ❌ 错误
    default: return "未知"   // ❌ 状态4会显示未知
    }
}
```

### 正确的映射（修复后）
```go
func getOrderStatusText(status int) string {
    switch status {
    case 0: return "待支付"
    case 1: return "已支付"
    case 2: return "已完成"  // ✅ 正确
    case 3: return "已取消"  // ✅ 正确
    case 4: return "已退款"  // ✅ 添加缺失的状态
    default: return "未知"
    }
}
```

## 数据库状态定义
根据数据库表定义，订单状态应该是：
- `0`: 待支付
- `1`: 已支付
- `2`: 已完成
- `3`: 已取消
- `4`: 已退款

## 修复内容
1. **修正状态映射顺序**：将"已完成"和"已取消"的状态码交换
2. **添加缺失状态**：添加状态4"已退款"的处理

## 验证方法
部署后，检查订单管理页面：
1. 状态为4的订单应该显示"已退款"
2. 状态为2的订单应该显示"已完成"
3. 状态为3的订单应该显示"已取消"
4. 不再有订单显示"未知"状态

## 相关文件
- `service/admin_service.go` - 管理员订单状态处理
- `service/order_service.go` - 用户订单状态处理（已正确）
- `db/model/order.go` - 订单模型定义
- `db/migration/create_service_order_referral_tables.sql` - 数据库表定义

## 注意事项
- 此修复只影响管理员页面的订单状态显示
- 用户页面的订单状态显示已经是正确的
- 修复后需要重新部署服务才能生效
