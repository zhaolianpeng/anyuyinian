# 订单相关接口文档

## 接口概览

本文档包含以下订单相关接口：
1. **提交订单** - `POST /api/order/submit`
2. **发起支付** - `POST /api/order/pay/:id`
3. **取消订单** - `POST /api/order/cancel/:id`
4. **申请退款** - `POST /api/order/refund/:id`
5. **订单列表** - `GET /api/order/list`
6. **订单详情** - `GET /api/order/detail/:id`

## 1. 提交订单

### 接口信息
- **接口地址**: `POST /api/order/submit`
- **请求方式**: POST
- **功能**: 提交新订单

### 请求参数
```json
{
  "userId": 1,
  "serviceId": 1,
  "patientId": 1,
  "addressId": 1,
  "appointmentDate": "2024-01-15",
  "appointmentTime": "morning",
  "specialRequirements": "无特殊要求",
  "formData": {
    "patientName": "张三",
    "patientPhone": "13800138000",
    "appointmentDate": "2024-01-15",
    "appointmentTime": "morning",
    "specialRequirements": "无特殊要求"
  }
}
```

### 响应格式
```json
{
  "code": 0,
  "data": {
    "orderId": 1,
    "orderNo": "202401150001",
    "message": "订单提交成功，请及时支付"
  }
}
```

## 2. 发起支付

### 接口信息
- **接口地址**: `POST /api/order/pay/:id`
- **请求方式**: POST
- **功能**: 发起订单支付（微信支付）

### 路径参数
- `id`: 订单ID

### 请求参数
```json
{
  "paymentMethod": "wechat_pay",
  "openId": "用户openId"
}
```

### 响应格式
```json
{
  "code": 0,
  "data": {
    "orderId": 1,
    "orderNo": "202401150001",
    "amount": 299.00,
    "paymentParams": {
      "appId": "wx1234567890abcdef",
      "timeStamp": "1642234567",
      "nonceStr": "randomstring",
      "package": "prepay_id=wx1234567890abcdef",
      "signType": "RSA",
      "paySign": "signature"
    }
  }
}
```

## 3. 取消订单

### 接口信息
- **接口地址**: `POST /api/order/cancel/:id`
- **请求方式**: POST
- **功能**: 取消订单

### 路径参数
- `id`: 订单ID

### 请求参数
```json
{
  "reason": "个人原因取消"
}
```

### 响应格式
```json
{
  "code": 0,
  "data": {
    "message": "订单取消成功"
  }
}
```

## 4. 申请退款

### 接口信息
- **接口地址**: `POST /api/order/refund/:id`
- **请求方式**: POST
- **功能**: 申请订单退款

### 路径参数
- `id`: 订单ID

### 请求参数
```json
{
  "reason": "服务不满意",
  "refundAmount": 299.00
}
```

### 响应格式
```json
{
  "code": 0,
  "data": {
    "refundId": 1,
    "message": "退款申请提交成功"
  }
}
```

## 5. 订单列表

### 接口信息
- **接口地址**: `GET /api/order/list`
- **请求方式**: GET
- **功能**: 获取订单列表

### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| userId | int | 是 | 用户ID |
| status | string | 否 | 订单状态：pending_pay, paid, cancelled, refunded |
| page | int | 否 | 页码，默认1 |
| pageSize | int | 否 | 每页数量，默认10，最大50 |

### 响应格式
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "orderNo": "202401150001",
        "userId": 1,
        "serviceId": 1,
        "serviceTitle": "全面体检套餐",
        "amount": 299.00,
        "status": "pending_pay",
        "statusText": "待支付",
        "appointmentDate": "2024-01-15",
        "appointmentTime": "morning",
        "createdAt": "2024-01-01T12:00:00Z",
        "updatedAt": "2024-01-01T12:00:00Z"
      }
    ],
    "total": 10,
    "page": 1,
    "pageSize": 10,
    "hasMore": false
  }
}
```

## 6. 订单详情

### 接口信息
- **接口地址**: `GET /api/order/detail/:id`
- **请求方式**: GET
- **功能**: 获取订单详细信息

### 路径参数
- `id`: 订单ID

### 响应格式
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "orderNo": "202401150001",
    "userId": 1,
    "serviceId": 1,
    "serviceTitle": "全面体检套餐",
    "serviceDescription": "包含血常规、尿常规、心电图等多项检查",
    "amount": 299.00,
    "status": "paid",
    "statusText": "已支付",
    "appointmentDate": "2024-01-15",
    "appointmentTime": "morning",
    "specialRequirements": "无特殊要求",
    "formData": {
      "patientName": "张三",
      "patientPhone": "13800138000",
      "appointmentDate": "2024-01-15",
      "appointmentTime": "morning",
      "specialRequirements": "无特殊要求"
    },
    "paymentInfo": {
      "paymentMethod": "wechat_pay",
      "transactionId": "wx1234567890abcdef",
      "paidAt": "2024-01-01T12:30:00Z"
    },
    "refundInfo": null,
    "createdAt": "2024-01-01T12:00:00Z",
    "updatedAt": "2024-01-01T12:30:00Z"
  }
}
```

## 订单状态说明

| 状态 | 状态文本 | 说明 |
|------|----------|------|
| pending_pay | 待支付 | 订单已创建，等待支付 |
| paid | 已支付 | 订单已支付成功 |
| cancelled | 已取消 | 订单已取消 |
| refunded | 已退款 | 订单已退款 |

## 使用示例

### 微信小程序端

```javascript
// 提交订单
wx.request({
  url: 'http://your-server.com/api/order/submit',
  method: 'POST',
  header: { 'Content-Type': 'application/json' },
  data: {
    userId: 1,
    serviceId: 1,
    patientId: 1,
    addressId: 1,
    appointmentDate: '2024-01-15',
    appointmentTime: 'morning',
    specialRequirements: '无特殊要求',
    formData: {
      patientName: '张三',
      patientPhone: '13800138000',
      appointmentDate: '2024-01-15',
      appointmentTime: 'morning',
      specialRequirements: '无特殊要求'
    }
  },
  success: (res) => {
    if (res.data.code === 0) {
      console.log('订单提交成功:', res.data.data);
      // 跳转到支付页面
      this.payOrder(res.data.data.orderId);
    }
  }
});

// 发起支付
payOrder(orderId) {
  wx.request({
    url: `http://your-server.com/api/order/pay/${orderId}`,
    method: 'POST',
    header: { 'Content-Type': 'application/json' },
    data: {
      paymentMethod: 'wechat_pay',
      openId: '用户openId'
    },
    success: (res) => {
      if (res.data.code === 0) {
        const paymentParams = res.data.data.paymentParams;
        
        // 调用微信支付
        wx.requestPayment({
          timeStamp: paymentParams.timeStamp,
          nonceStr: paymentParams.nonceStr,
          package: paymentParams.package,
          signType: paymentParams.signType,
          paySign: paymentParams.paySign,
          success: () => {
            console.log('支付成功');
            // 跳转到订单详情页
            wx.navigateTo({
              url: `/pages/order/detail?id=${orderId}`
            });
          },
          fail: () => {
            console.log('支付失败');
          }
        });
      }
    }
  });
}

// 获取订单列表
wx.request({
  url: 'http://your-server.com/api/order/list',
  method: 'GET',
  data: {
    userId: 1,
    status: 'pending_pay',
    page: 1,
    pageSize: 10
  },
  success: (res) => {
    if (res.data.code === 0) {
      console.log('订单列表:', res.data.data);
    }
  }
});

// 获取订单详情
wx.request({
  url: 'http://your-server.com/api/order/detail/1',
  method: 'GET',
  success: (res) => {
    if (res.data.code === 0) {
      console.log('订单详情:', res.data.data);
    }
  }
});

// 取消订单
wx.request({
  url: 'http://your-server.com/api/order/cancel/1',
  method: 'POST',
  header: { 'Content-Type': 'application/json' },
  data: {
    reason: '个人原因取消'
  },
  success: (res) => {
    if (res.data.code === 0) {
      console.log('订单取消成功:', res.data.data);
    }
  }
});

// 申请退款
wx.request({
  url: 'http://your-server.com/api/order/refund/1',
  method: 'POST',
  header: { 'Content-Type': 'application/json' },
  data: {
    reason: '服务不满意',
    refundAmount: 299.00
  },
  success: (res) => {
    if (res.data.code === 0) {
      console.log('退款申请成功:', res.data.data);
    }
  }
});
```

## 数据库表结构

### 订单表 (Orders)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | INT | 主键，自增 |
| orderNo | VARCHAR(50) | 订单号，唯一 |
| userId | INT | 用户ID |
| serviceId | INT | 服务ID |
| patientId | INT | 就诊人ID |
| addressId | INT | 地址ID |
| amount | DECIMAL(10,2) | 订单金额 |
| status | VARCHAR(20) | 订单状态 |
| appointmentDate | DATE | 预约日期 |
| appointmentTime | VARCHAR(20) | 预约时间 |
| specialRequirements | TEXT | 特殊要求 |
| formData | TEXT | 表单数据（JSON） |
| paymentMethod | VARCHAR(20) | 支付方式 |
| transactionId | VARCHAR(100) | 交易ID |
| paidAt | DATETIME | 支付时间 |
| refundAmount | DECIMAL(10,2) | 退款金额 |
| refundReason | VARCHAR(500) | 退款原因 |
| refundedAt | DATETIME | 退款时间 |
| referralUserId | INT | 推荐人用户ID |
| commissionAmount | DECIMAL(10,2) | 佣金金额 |
| createdAt | DATETIME | 创建时间 |
| updatedAt | DATETIME | 更新时间 |

## 注意事项

1. **订单号生成**: 使用时间戳+随机数生成唯一订单号
2. **状态管理**: 订单状态流转：待支付 -> 已支付 -> 已完成/已取消/已退款
3. **支付集成**: 支持微信支付等第三方支付方式
4. **退款处理**: 支持部分退款和全额退款
5. **数据验证**: 订单提交前会验证用户、服务、就诊人等信息
6. **推荐系统**: 支持推荐人佣金计算
7. **时间限制**: 订单支付有时间限制，超时自动取消
8. **通知机制**: 订单状态变更会发送通知给用户
9. **数据安全**: 敏感信息（如支付参数）需要加密处理
10. **日志记录**: 记录订单操作日志，便于问题排查 