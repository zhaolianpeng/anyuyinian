# 管理员API标准化改造

## 问题描述
管理员页面使用了 `app.callContainer` 调用方式，而其他页面使用标准的 `callContainer` 函数。为了保持一致性，需要将所有管理员页面改为使用标准API调用。

## 已完成的修改

### 1. 添加标准API导入
已为所有管理员页面添加了标准API导入：
```javascript
const { api } = require('../../utils/cloud-container-standard')
```

### 2. 添加管理员API到标准调用
已在 `cloud-container-standard.js` 中添加了所有管理员相关API：
```javascript
// 管理员相关接口
adminLogin: (data) => callContainer('/api/admin/login', 'POST', data),
adminCheckStatus: (data) => callContainer('/api/admin/check-status', 'POST', data),
adminUsers: (params) => callContainer('/api/admin/users', 'GET', params),
adminOrders: (params) => callContainer('/api/admin/orders', 'GET', params),
adminStats: (params) => callContainer('/api/admin/stats', 'GET', params),
adminAdmins: (params) => callContainer('/api/admin/admins', 'GET', params),
adminSetAdmin: (data) => callContainer('/api/admin/set-admin', 'POST', data),
adminRemoveAdmin: (data) => callContainer('/api/admin/remove-admin', 'POST', data),
adminUpdateOrderAmount: (data) => callContainer('/api/admin/order/update-amount', 'POST', data),
adminRefundOrder: (data) => callContainer('/api/admin/order/refund', 'POST', data),
adminServices: (params) => callContainer('/api/admin/services', 'GET', params),
adminUpdateServicePrice: (data) => callContainer('/api/admin/service/updateprice', 'POST', data),
consultationStats: () => callContainer('/api/consultation/stats', 'GET')
```

### 3. 已修改的文件
- `pages/admin/services.js` - 已完全修改为使用标准API调用

## 需要手动修改的文件

### 1. pages/admin/login.js
```javascript
// 修改前
app.callContainer('/api/admin/login', 'POST', {
  username: superUsername,
  password: superPassword
})

// 修改后
api.adminLogin({
  username: superUsername,
  password: superPassword
})
```

### 2. pages/admin/home.js
```javascript
// 修改前
app.callContainer('/api/admin/stats', 'GET', {
  adminUserId: adminInfo.userId
})

// 修改后
api.adminStats({
  adminUserId: adminInfo.userId
})

// 修改前
app.callContainer('/api/consultation/stats', 'GET')

// 修改后
api.consultationStats()
```

### 3. pages/admin/admins.js
```javascript
// 修改前
app.callContainer('/api/admin/admins', 'GET', {
  adminUserId: adminInfo.userId,
  page: this.data.page,
  pageSize: this.data.pageSize
})

// 修改后
api.adminAdmins({
  adminUserId: adminInfo.userId,
  page: this.data.page,
  pageSize: this.data.pageSize
})

// 修改前
app.callContainer('/api/admin/set-admin', 'POST', {
  userId: userId,
  adminLevel: 1,
  parentAdminId: adminInfo.userId,
  adminUsername: username,
  adminPassword: password
})

// 修改后
api.adminSetAdmin({
  userId: userId,
  adminLevel: 1,
  parentAdminId: adminInfo.userId,
  adminUsername: username,
  adminPassword: password
})

// 修改前
app.callContainer('/api/admin/remove-admin', 'POST', {
  userId: userId
})

// 修改后
api.adminRemoveAdmin({
  userId: userId
})
```

### 4. pages/admin/orders.js
```javascript
// 修改前
app.callContainer('/api/admin/orders', 'GET', {}, {
  query: {
    adminUserId: adminInfo.userId,
    page: this.data.page,
    pageSize: this.data.pageSize
  }
})

// 修改后
api.adminOrders({
  adminUserId: adminInfo.userId,
  page: this.data.page,
  pageSize: this.data.pageSize
})

// 修改前
app.callContainer('/api/admin/order/update-amount', 'POST', {
  orderId: orderId,
  newAmount: newAmount,
  reason: '管理员手动修改'
}, {
  query: {
    adminUserId: adminInfo.userId
  }
})

// 修改后
api.adminUpdateOrderAmount({
  orderId: orderId,
  newAmount: newAmount,
  reason: '管理员手动修改',
  adminUserId: adminInfo.userId
})

// 修改前
app.callContainer('/api/admin/order/refund', 'POST', {
  orderId: orderId,
  refundAmount: refundAmount,
  reason: reason,
  refundStatus: refundStatus
}, {
  query: {
    adminUserId: adminInfo.userId
  }
})

// 修改后
api.adminRefundOrder({
  orderId: orderId,
  refundAmount: refundAmount,
  reason: reason,
  refundStatus: refundStatus,
  adminUserId: adminInfo.userId
})
```

### 5. pages/admin/users.js
```javascript
// 修改前
app.callContainer('/api/admin/users', 'GET', {
  adminUserId: adminInfo.userId,
  page: this.data.page,
  pageSize: this.data.pageSize
})

// 修改后
api.adminUsers({
  adminUserId: adminInfo.userId,
  page: this.data.page,
  pageSize: this.data.pageSize
})

// 其他调用方式同 admins.js
```

## 修改原则

1. **移除 app.callContainer 调用**：删除所有 `app.callContainer` 调用
2. **使用 api.xxx 调用**：替换为对应的 `api.xxx` 调用
3. **参数传递**：将 query 参数合并到主参数对象中
4. **保持功能不变**：确保修改后功能完全一致

## 验证方法

修改完成后，测试以下功能：
1. 管理员登录
2. 管理员首页数据加载
3. 用户管理功能
4. 订单管理功能
5. 管理员管理功能
6. 服务价格更新功能

## 相关文件
- `utils/cloud-container-standard.js` - 标准API调用定义
- `pages/admin/*.js` - 所有管理员页面
- `app.js` - 原有的 callContainer 实现
