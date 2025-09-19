# 订单状态显示完整修复

## 问题描述
历史退款订单在订单管理页面显示为"未知"状态，而不是正确的"已退款"。

## 问题根因
存在多个地方的状态映射不一致：

### 1. 后端状态映射（已修复）
`service/admin_service.go` 中的 `getOrderStatusText` 函数状态映射不正确。

### 2. 前端状态映射（已修复）
`pages/admin/orders.js` 中的 `getOrderStatusText` 函数状态映射也不正确。

### 3. 测试文件状态映射（已修复）
`tests/frontend/test_admin_orders_fix.js` 中的状态映射也不正确。

## 修复内容

### 1. 后端修复（service/admin_service.go）
```go
func getOrderStatusText(status int) string {
    switch status {
    case 0: return "待支付"
    case 1: return "已支付"
    case 2: return "已完成"  // ✅ 修正
    case 3: return "已取消"  // ✅ 修正
    case 4: return "已退款"  // ✅ 新增
    default: return "未知"
    }
}
```

### 2. 前端修复（pages/admin/orders.js）
```javascript
getOrderStatusText: function (status) {
  const statusMap = {
    0: '待支付',
    1: '已支付',
    2: '已完成',  // ✅ 修正
    3: '已取消',  // ✅ 修正
    4: '已退款'   // ✅ 新增
  };
  return statusMap[status] || '未知';
}
```

### 3. 测试文件修复（tests/frontend/test_admin_orders_fix.js）
```javascript
const getOrderStatusText = (status) => {
  const statusMap = { 0: '待支付', 1: '已支付', 2: '已完成', 3: '已取消', 4: '已退款' };
  return statusMap[status] || '未知';
};
```

## 正确的状态映射
根据数据库定义，订单状态应该是：
- `0`: 待支付
- `1`: 已支付
- `2`: 已完成
- `3`: 已取消
- `4`: 已退款

## 修复的文件列表
1. `service/admin_service.go` - 后端管理员订单状态处理
2. `pages/admin/orders.js` - 前端管理员订单页面
3. `tests/frontend/test_admin_orders_fix.js` - 测试文件

## 验证方法
1. **重新部署后端服务**
2. **重新编译小程序**
3. **检查订单管理页面**：
   - 状态为4的订单应该显示"已退款"
   - 状态为2的订单应该显示"已完成"
   - 状态为3的订单应该显示"已取消"
   - 不再有订单显示"未知"状态

## 注意事项
- 此修复同时涉及前端和后端
- 需要重新部署后端服务
- 需要重新编译小程序
- 修复后历史订单状态将正确显示

## 相关文件
- `service/admin_service.go` - 后端状态处理
- `pages/admin/orders.js` - 前端状态显示
- `tests/frontend/test_admin_orders_fix.js` - 测试文件
- `db/model/order.go` - 订单模型定义
