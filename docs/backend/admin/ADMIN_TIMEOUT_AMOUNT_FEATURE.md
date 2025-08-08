# 管理员首页超时未支付金额功能

## 功能概述

在管理员首页数据概览中新增了"超时未支付金额"统计项，用于显示已超过支付截止时间但尚未支付的订单总金额。

## 实现细节

### 后端修改

#### 1. 数据库查询逻辑
- **超级管理员**：统计所有超时未支付订单
- **一级管理员**：只统计通过自己推广码注册用户的超时未支付订单

#### 2. 查询条件
```sql
-- 超时未支付订单查询条件
status = 0 AND payStatus = 0 AND payDeadline IS NOT NULL AND payDeadline < NOW()
```

- `status = 0`：订单状态为待支付
- `payStatus = 0`：支付状态为未支付
- `payDeadline IS NOT NULL`：有支付截止时间
- `payDeadline < NOW()`：支付截止时间已过期

#### 3. 修改文件
- `anyuyinian/service/admin_service.go`：在 `AdminStatsHandler` 函数中添加超时未支付金额统计

### 前端修改

#### 1. 页面显示
- 在管理员首页数据概览网格中新增"超时未支付金额"卡片
- 显示格式：`¥{{stats.timeoutUnpaidAmount || 0}}`

#### 2. 样式设计
- 使用橙色 (`#fa8c16`) 突出显示超时金额
- 与其他金额统计保持一致的卡片样式

#### 3. 修改文件
- `miniprogram/pages/admin/home.wxml`：添加超时未支付金额显示
- `miniprogram/pages/admin/home.wxss`：添加超时金额样式

## 数据字段

### 新增字段
```json
{
  "timeoutUnpaidAmount": 0.00  // 超时未支付金额（元）
}
```

### 完整响应示例
```json
{
  "code": 0,
  "data": {
    "totalUsers": 100,
    "totalOrders": 50,
    "todayOrders": 5,
    "totalAmount": 5000.00,
    "paidAmount": 3000.00,
    "unpaidAmount": 1500.00,
    "refundAmount": 500.00,
    "timeoutUnpaidAmount": 800.00
  }
}
```

## 权限控制

### 超级管理员 (adminLevel = 2)
- 可以看到所有用户的超时未支付金额
- 统计范围：全平台所有订单

### 一级管理员 (adminLevel = 1)
- 只能看到通过自己推广码注册用户的超时未支付金额
- 统计范围：自己推广的用户订单

## 测试验证

### 测试脚本
- `tests/backend/test_admin_timeout_amount.sh`：验证后端API功能

### 测试要点
1. ✅ 接口返回包含 `timeoutUnpaidAmount` 字段
2. ✅ 超时未支付金额计算正确
3. ✅ 权限控制生效（不同级别管理员看到不同数据）
4. ✅ 前端正确显示超时金额

## 部署说明

### 后端部署
1. 更新 `admin_service.go` 文件
2. 重新编译并部署后端服务
3. 无需数据库迁移（使用现有字段）

### 前端部署
1. 更新管理员首页相关文件
2. 重新编译小程序
3. 上传到微信开发者工具

## 注意事项

1. **数据准确性**：超时未支付金额基于 `payDeadline` 字段计算，确保订单创建时正确设置支付截止时间
2. **性能考虑**：查询条件包含时间比较，建议在 `payDeadline` 字段上建立索引
3. **用户体验**：超时金额使用橙色显示，提醒管理员关注此类订单
4. **权限安全**：一级管理员只能看到自己推广用户的超时金额，确保数据隔离

## 后续优化建议

1. **实时更新**：考虑添加WebSocket推送，实时更新超时金额
2. **筛选功能**：在订单管理页面添加超时订单筛选
3. **自动处理**：考虑添加超时订单自动取消功能
4. **统计报表**：添加超时订单趋势分析 