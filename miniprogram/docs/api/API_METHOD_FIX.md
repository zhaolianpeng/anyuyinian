# API方法修复总结

> 本文中的接口调用统一按“服务端 API 封装”理解。
> 代码里当时的函数名仍位于 `miniprogram/utils/cloud-container-standard.js`，但这里不再强调历史命名，而强调统一请求层的职责。

## 问题描述

在前后端接口调用过程中，发现部分接口的方法不匹配问题：

1. **订单详情API错误**：`/api/order/detail`接口返回只支持POST请求，但前端使用GET方法
2. **文件删除API错误**：`/api/file/delete`接口要求DELETE方法，但前端使用POST方法
3. **文件权限更新API错误**：`/api/file/permission`接口要求PUT方法，但前端使用POST方法

## 问题分析

### 1. 订单详情API问题
- **原因**：前端使用GET方法调用`/api/order/detail`接口
- **后端要求**：该接口只支持POST请求
- **解决方案**：修改前端请求方法为POST，并调整参数传递方式

### 2. 文件管理API问题
- **原因**：前端使用POST方法调用文件删除和权限更新接口
- **后端要求**：文件删除要求DELETE方法，权限更新要求PUT方法
- **解决方案**：修改前端请求方法为对应的HTTP方法

## 修复方案

### 1. 修复订单详情API请求方法

**文件**：`miniprogram/utils/cloud-container-standard.js`

**修改内容**：
```javascript
// 修改前：统一服务端 API 封装以 GET 方式请求订单详情
orderDetail: (orderNo) => requestApi('/api/order/detail', 'GET', { orderNo }),

// 修改后：统一服务端 API 封装改为 POST 提交参数
orderDetail: (data) => requestApi('/api/order/detail', 'POST', data),
```

**文件**：`miniprogram/pages/order/detail.js`

**修改内容**：
```javascript
// 修改前
const result = await api.orderDetail(orderNo)

// 修改后
const result = await api.orderDetail({ orderNo })
```

### 2. 修复文件管理API请求方法

**文件**：`miniprogram/utils/cloud-container-standard.js`

**修改内容**：
```javascript
// 修改前：统一服务端 API 封装错误地复用了 POST
deleteFile: (data) => requestApi('/api/file/delete', 'POST', data),
updateFilePermission: (data) => requestApi('/api/file/permission', 'POST', data),

// 修改后：统一服务端 API 封装按后端协议使用正确方法
deleteFile: (fileId) => requestApi('/api/file/delete', 'DELETE', {}, { fileId }),
updateFilePermission: (data) => requestApi('/api/file/permission', 'PUT', data),
```

## 完整的前后端接口方法匹配表

### ✅ 已正确匹配的接口：

| 接口路径 | 后端要求 | 前端使用 | 状态 |
|---------|---------|---------|------|
| `/api/wx/login` | POST | POST | ✅ |
| `/api/home/init` | GET/POST | GET | ✅ |
| `/api/user/info` | GET | GET | ✅ |
| `/api/user/bind_phone` | POST | POST | ✅ |
| `/api/user/address` | GET | GET | ✅ |
| `/api/user/patient` | GET | GET | ✅ |
| `/api/service/list` | GET | GET | ✅ |
| `/api/service/detail` | POST | POST | ✅ |
| `/api/service/form_config/:id` | GET | GET | ✅ |
| `/api/order/submit` | POST | POST | ✅ |
| `/api/order/pay/:id` | POST | POST | ✅ |
| `/api/order/cancel/:id` | POST | POST | ✅ |
| `/api/order/refund/:id` | POST | POST | ✅ |
| `/api/order/list` | GET | GET | ✅ |
| `/api/order/detail` | POST | POST | ✅ (已修复) |
| `/api/order/time_slots` | POST | POST | ✅ |
| `/api/upload` | POST | POST | ✅ |
| `/api/files` | GET | GET | ✅ |
| `/api/file/delete` | DELETE | DELETE | ✅ (已修复) |
| `/api/file/permission` | PUT | PUT | ✅ (已修复) |
| `/api/file/permission/get` | GET | GET | ✅ |
| `/api/config` | GET | GET | ✅ |
| `/api/referral/qrcode` | GET | GET | ✅ |
| `/api/referral/report` | GET | GET | ✅ |
| `/api/referral/config` | GET | GET | ✅ |
| `/api/referral/apply_cashout` | POST | POST | ✅ |
| `/api/kefu/send_msg` | POST | POST | ✅ |
| `/api/kefu/faq` | GET | GET | ✅ |
| `/api/hospital/list` | GET | GET | ✅ |
| `/api/hospital/detail/:id` | GET | GET | ✅ |
| `/api/count` | GET/POST | GET/POST | ✅ |

### 🔧 已修复的接口：

1. **`/api/order/detail`** - 从 GET 改为 POST ✅
2. **`/api/file/delete`** - 从 POST 改为 DELETE ✅
3. **`/api/file/permission`** - 从 POST 改为 PUT ✅

## 修复效果

修复后，所有前后端接口的方法都正确匹配，不再出现以下错误：

- `HTTP 405: 只支持POST请求`
- `HTTP 405: 只支持GET请求`
- `HTTP 405: 只支持DELETE请求`
- `HTTP 405: 只支持PUT请求`

## 注意事项

1. **统一封装约定**：所有业务请求都应先经过统一服务端 API 封装，再由封装层处理方法、参数和公共头信息
2. **参数传递方式**：GET请求通过URL参数传递，POST/PUT/DELETE请求通过请求体传递
3. **错误处理**：所有接口调用都应该包含适当的错误处理
4. **日志记录**：建议在关键接口调用时添加日志记录，便于调试
5. **测试验证**：修复后应该进行完整的接口测试，确保功能正常

## 后续维护

1. **新增接口**：添加新接口时，先在统一服务端 API 封装层定义调用方式，再接入业务页面
2. **接口变更**：修改接口方法时，需要同步更新封装层与业务调用代码
3. **文档同步**：保持接口文档与实际实现的一致性