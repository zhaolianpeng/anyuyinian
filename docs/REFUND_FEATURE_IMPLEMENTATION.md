# 退款功能实现文档

## 功能概述

退款功能允许用户申请退款，管理员处理退款，包括完整的权限控制、状态管理和错误处理。

## 功能特性

### 1. 用户退款申请
- ✅ 只有已支付的订单可以申请退款
- ✅ 防止重复申请退款
- ✅ 退款金额验证（不能超过订单金额）
- ✅ 完整的错误处理和用户提示

### 2. 管理员退款处理
- ✅ 只有超级管理员可以处理退款
- ✅ 支持设置退款状态（退款中/已退款）
- ✅ 自动更新订单状态
- ✅ 完整的权限验证

### 3. 状态管理
- ✅ 退款状态：0-未退款，1-退款中，2-已退款
- ✅ 订单状态联动：已退款时自动更新订单状态为4
- ✅ 状态变更记录和日志

## 后端实现

### 1. 用户退款接口

**接口地址**: `POST /api/order/refund/:id`

**请求参数**:
```json
{
  "orderId": 1,
  "refundAmount": 299.00,
  "reason": "用户申请退款"
}
```

**响应格式**:
```json
{
  "code": 0,
  "data": {
    "orderId": 1,
    "orderNo": "202401150001",
    "refundAmount": 299.00,
    "reason": "用户申请退款",
    "message": "退款申请提交成功"
  }
}
```

**验证规则**:
- 订单必须存在
- 订单状态必须为已支付（status=1, payStatus=1）
- 退款金额必须大于0且不超过订单金额
- 不能重复申请退款

### 2. 管理员退款接口

**接口地址**: `POST /api/admin/order/refund?adminUserId=xxx`

**请求参数**:
```json
{
  "orderId": 1,
  "refundAmount": 299.00,
  "reason": "管理员处理退款",
  "refundStatus": 2
}
```

**响应格式**:
```json
{
  "code": 0,
  "data": {
    "orderId": 1,
    "orderNo": "202401150001",
    "refundAmount": 299.00,
    "reason": "管理员处理退款",
    "refundStatus": 2,
    "adminId": "anyuyinian",
    "message": "退款处理成功"
  }
}
```

**权限要求**:
- 必须是有效的管理员账号
- 必须是超级管理员（adminLevel=2）

## 前端实现

### 1. 用户端退款功能

#### 订单详情页面 (`pages/order/detail.js`)
```javascript
async onRefundOrder() {
  // 检查订单状态
  if (order.status !== 1) {
    wx.showToast({ title: '只有已支付的订单可以申请退款' });
    return;
  }
  
  // 检查是否已申请退款
  if (order.refundStatus > 0) {
    wx.showToast({ title: '订单已申请退款，请勿重复申请' });
    return;
  }
  
  // 显示确认弹窗
  wx.showModal({
    title: '申请退款',
    content: `确定要申请退款吗？\n退款金额：¥${order.totalAmount}`,
    success: async (res) => {
      if (res.confirm) {
        // 调用退款API
        const result = await api.orderRefund(order.id, {
          orderId: order.id,
          refundAmount: order.totalAmount,
          reason: '用户申请退款'
        });
      }
    }
  });
}
```

#### 订单列表页面 (`pages/order/list.js`)
- 在订单列表中显示退款按钮
- 只有已支付且未申请退款的订单显示退款按钮
- 点击退款按钮直接申请退款

### 2. 管理员端退款功能

#### 管理员订单页面 (`pages/admin/orders.js`)
```javascript
onRefundOrder: function (e) {
  const order = e.currentTarget.dataset.order;
  
  // 权限检查
  if (!adminInfo || adminInfo.adminLevel !== 2) {
    wx.showToast({ title: '权限不足' });
    return;
  }
  
  // 显示操作选项
  wx.showActionSheet({
    itemList: ['设置为退款中', '设置为已退款'],
    success: (res) => {
      const refundStatus = res.tapIndex === 0 ? 1 : 2;
      
      // 显示退款确认弹窗
      wx.showModal({
        title: '处理退款',
        content: `确定要将订单设置为${statusText}吗？`,
        editable: true,
        placeholderText: '请输入退款原因',
        success: (modalRes) => {
          if (modalRes.confirm) {
            this.processRefund(order.id, order.amount, modalRes.content, refundStatus);
          }
        }
      });
    }
  });
}
```

## 数据库设计

### 订单表字段
```sql
-- 退款相关字段
refundStatus INT DEFAULT 0 COMMENT '退款状态：0-未退款，1-退款中，2-已退款',
refundTime DATETIME COMMENT '退款时间',
refundAmount DECIMAL(10,2) COMMENT '退款金额',
refundReason VARCHAR(500) COMMENT '退款原因',
```

### 状态说明
- **refundStatus**: 
  - 0: 未退款
  - 1: 退款中（用户申请或管理员设置）
  - 2: 已退款（管理员确认退款完成）

## 错误处理

### 常见错误情况
1. **订单不存在**: 返回"订单不存在"
2. **订单状态错误**: 返回"只有已支付的订单可以申请退款"
3. **退款金额无效**: 返回"退款金额必须大于0"
4. **退款金额超限**: 返回"退款金额不能超过订单金额"
5. **重复申请**: 返回"订单已申请退款，请勿重复申请"
6. **权限不足**: 返回"只有超级管理员可以处理退款"

### 日志记录
- 所有退款操作都有详细的日志记录
- 包括操作时间、操作人、订单信息、退款金额等

## 测试

### 测试脚本
运行 `tests/backend/test_refund_functionality.sh` 进行完整的功能测试：

```bash
./tests/backend/test_refund_functionality.sh
```

### 测试内容
1. ✅ 用户申请退款
2. ✅ 管理员处理退款
3. ✅ 权限控制测试
4. ✅ 参数验证测试
5. ✅ 错误情况测试

## 使用流程

### 用户申请退款流程
1. 用户在订单详情页或订单列表页点击"申请退款"
2. 系统检查订单状态和退款状态
3. 显示确认弹窗，显示退款金额
4. 用户确认后提交退款申请
5. 系统更新退款状态为"退款中"
6. 显示申请成功提示

### 管理员处理退款流程
1. 超级管理员在订单列表页点击"处理退款"
2. 选择退款状态（退款中/已退款）
3. 输入退款原因
4. 确认处理退款
5. 系统更新退款状态和订单状态
6. 显示处理成功提示

## 安全考虑

1. **权限控制**: 只有超级管理员可以处理退款
2. **参数验证**: 严格的参数验证防止恶意请求
3. **状态检查**: 多重状态检查确保操作合法性
4. **日志记录**: 完整的操作日志便于审计
5. **金额验证**: 防止退款金额超过订单金额

## 扩展功能

### 未来可能的扩展
1. **部分退款**: 支持部分金额退款
2. **退款通知**: 退款状态变更时发送通知
3. **退款统计**: 退款金额统计和报表
4. **自动退款**: 特定条件下自动处理退款
5. **退款审核**: 多级审核流程

## 总结

退款功能已经完整实现，包括：
- ✅ 用户端退款申请
- ✅ 管理员端退款处理
- ✅ 完整的权限控制
- ✅ 严格的状态管理
- ✅ 详细的错误处理
- ✅ 完整的测试覆盖

功能已经可以投入使用！ 