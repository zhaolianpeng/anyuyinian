# 安语颐年后端API接口文档

## 接口概述

安语颐年医疗服务平台后端API，提供用户管理、订单管理、服务管理、支付处理、咨询系统等核心功能。

**基础URL**: `https://your-domain.com`

**通用响应格式**:
```json
{
    "code": 0,           // 0表示成功，非0表示失败
    "errorMsg": "",      // 错误信息
    "data": {}           // 响应数据
}
```

---

## 1. 用户认证模块

### 1.1 微信登录
**接口地址**: `POST /api/wx/login`

**请求参数**:
```json
{
    "code": "string",        // 微信登录凭证
    "nickName": "string",    // 用户昵称（可选）
    "avatarUrl": "string",   // 头像URL（可选）
    "gender": 0,             // 性别：0-未知，1-男，2-女（可选）
    "country": "string",     // 国家（可选）
    "province": "string",    // 省份（可选）
    "city": "string",        // 城市（可选）
    "language": "string"     // 语言（可选）
}
```

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "id": 1,
        "userId": "user_123456",
        "openId": "openid_123456",
        "nickName": "用户昵称",
        "avatarUrl": "https://example.com/avatar.jpg",
        "gender": 1,
        "phone": "13800138000",
        "country": "中国",
        "province": "广东省",
        "city": "深圳市",
        "language": "zh_CN",
        "lastLoginAt": "2024-01-01T12:00:00Z",
        "isNewUser": false,
        "token": "token_user_123456_1640995200"
    }
}
```

### 1.2 获取用户信息
**接口地址**: `GET /api/user/info`

**请求参数**:
- `userId` (string, 必需): 用户ID

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "id": 1,
        "userId": "user_123456",
        "openId": "openid_123456",
        "nickName": "用户昵称",
        "avatarUrl": "https://example.com/avatar.jpg",
        "gender": 1,
        "phone": "13800138000",
        "country": "中国",
        "province": "广东省",
        "city": "深圳市",
        "language": "zh_CN"
    }
}
```

### 1.3 绑定手机号
**接口地址**: `POST /api/user/bind_phone`

**请求参数**:
```json
{
    "userId": "string",  // 用户ID
    "phone": "string",   // 手机号
    "code": "string"     // 验证码
}
```

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "userId": "user_123456",
        "phone": "13800138000",
        "message": "手机号绑定成功"
    }
}
```

### 1.4 更新用户信息
**接口地址**: `POST /api/user/update_info`

**请求参数**:
```json
{
    "userId": "string",     // 用户ID
    "nickName": "string",   // 昵称（可选）
    "avatarUrl": "string",  // 头像URL（可选）
    "gender": 1             // 性别（可选）
}
```

### 1.5 解密微信手机号
**接口地址**: `POST /api/user/decrypt_phone`

**请求参数**:
```json
{
    "userId": "string",        // 用户ID
    "encryptedData": "string", // 加密数据
    "iv": "string"             // 初始向量
}
```

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "phoneNumber": "13800138000"
    }
}
```

---

## 2. 地址管理模块

### 2.1 获取地址列表
**接口地址**: `GET /api/user/address`

**请求参数**:
- `userId` (string, 必需): 用户ID

**响应数据**:
```json
{
    "code": 0,
    "data": [
        {
            "id": 1,
            "userId": "user_123456",
            "name": "张三",
            "phone": "13800138000",
            "province": "广东省",
            "city": "深圳市",
            "district": "南山区",
            "address": "科技园南区",
            "isDefault": 1,
            "createdAt": "2024-01-01T12:00:00Z"
        }
    ]
}
```

### 2.2 创建地址
**接口地址**: `POST /api/user/address`

**请求参数**:
```json
{
    "userId": "string",     // 用户ID
    "name": "string",       // 收货人姓名
    "phone": "string",      // 手机号
    "province": "string",   // 省份
    "city": "string",       // 城市
    "district": "string",   // 区县
    "address": "string",    // 详细地址
    "isDefault": true       // 是否默认地址
}
```

### 2.3 更新地址
**接口地址**: `PUT /api/user/address`

**请求参数**:
```json
{
    "id": 1,                // 地址ID
    "userId": "string",     // 用户ID
    "name": "string",       // 收货人姓名
    "phone": "string",      // 手机号
    "province": "string",   // 省份
    "city": "string",       // 城市
    "district": "string",   // 区县
    "address": "string",    // 详细地址
    "isDefault": true       // 是否默认地址
}
```

### 2.4 删除地址
**接口地址**: `DELETE /api/user/address`

**请求参数**:
- `id` (int, 必需): 地址ID

---

## 3. 就诊人管理模块

### 3.1 获取就诊人列表
**接口地址**: `GET /api/user/patient`

**请求参数**:
- `userId` (string, 必需): 用户ID

**响应数据**:
```json
{
    "code": 0,
    "data": [
        {
            "id": 1,
            "userId": "user_123456",
            "name": "张三",
            "idCard": "440301199001011234",
            "phone": "13800138000",
            "gender": 1,
            "birthday": "1990-01-01",
            "age": 34,
            "relation": "本人",
            "isDefault": 1,
            "status": 1,
            "createdAt": "2024-01-01T12:00:00Z"
        }
    ]
}
```

### 3.2 创建就诊人
**接口地址**: `POST /api/user/patient`

**请求参数**:
```json
{
    "userId": "string",     // 用户ID
    "name": "string",       // 姓名
    "idCard": "string",     // 身份证号
    "phone": "string",      // 手机号
    "gender": 1,            // 性别：1-男，2-女
    "birthday": "string",   // 生日
    "relation": "string",   // 关系
    "isDefault": true       // 是否默认
}
```

### 3.3 更新就诊人
**接口地址**: `PUT /api/user/patient`

**请求参数**:
```json
{
    "id": 1,                // 就诊人ID
    "userId": "string",     // 用户ID
    "name": "string",       // 姓名
    "idCard": "string",     // 身份证号
    "phone": "string",      // 手机号
    "gender": 1,            // 性别
    "birthday": "string",   // 生日
    "relation": "string",   // 关系
    "isDefault": true       // 是否默认
}
```

### 3.4 删除就诊人
**接口地址**: `DELETE /api/user/patient`

**请求参数**:
- `id` (int, 必需): 就诊人ID

---

## 4. 服务管理模块

### 4.1 获取服务列表
**接口地址**: `GET /api/service/list`

**请求参数**:
- `category` (string, 可选): 服务分类
- `page` (int, 可选): 页码，默认1
- `pageSize` (int, 可选): 每页大小，默认10，最大50

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "list": [
            {
                "id": 1,
                "name": "上门护理服务",
                "description": "专业护理人员上门服务",
                "category": "护理服务",
                "price": 200.00,
                "originalPrice": 250.00,
                "imageUrl": "https://example.com/service1.jpg",
                "detailImages": "[\"url1\", \"url2\"]",
                "formConfig": "{\"fields\": [...]}",
                "videoUrl": "https://example.com/video.mp4",
                "status": 1,
                "sort": 0,
                "createdAt": "2024-01-01T12:00:00Z"
            }
        ],
        "total": 100,
        "page": 1,
        "pageSize": 10,
        "hasMore": true
    }
}
```

### 4.2 获取服务分类
**接口地址**: `GET /api/service/categories`

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "categories": [
            {
                "name": "全部",
                "value": "",
                "count": 0
            },
            {
                "name": "护理服务",
                "value": "护理服务",
                "count": 25
            },
            {
                "name": "康复治疗",
                "value": "康复治疗",
                "count": 15
            }
        ]
    }
}
```

### 4.3 获取服务详情
**接口地址**: `POST /api/service/detail`

**请求参数**:
```json
{
    "serviceId": 1
}
```

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "id": 1,
        "name": "上门护理服务",
        "description": "专业护理人员上门服务",
        "category": "护理服务",
        "price": 200.00,
        "originalPrice": 250.00,
        "imageUrl": "https://example.com/service1.jpg",
        "detailImages": "[\"url1\", \"url2\"]",
        "formConfig": "{\"fields\": [...]}",
        "videoUrl": "https://example.com/video.mp4",
        "status": 1,
        "sort": 0,
        "createdAt": "2024-01-01T12:00:00Z"
    }
}
```

### 4.4 获取服务表单配置
**接口地址**: `GET /api/service/form_config/{serviceId}`

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "fields": [
            {
                "name": "consultTime",
                "label": "服务时间",
                "type": "select",
                "required": true,
                "placeholder": "请选择服务时间",
                "options": [
                    {"label": "上午", "value": "morning"},
                    {"label": "下午", "value": "afternoon"}
                ]
            }
        ]
    }
}
```

---

## 5. 订单管理模块

### 5.1 提交订单
**接口地址**: `POST /api/order/submit`

**请求参数**:
```json
{
    "userId": "string",           // 用户ID
    "serviceId": 1,               // 服务ID
    "patientId": 1,               // 就诊人ID
    "addressId": 1,               // 地址ID
    "appointmentDate": "2024-01-02", // 预约日期
    "appointmentTime": "09:00",   // 预约时间
    "quantity": 1,                // 数量
    "formData": {                 // 表单数据
        "consultTime": "morning"
    },
    "referrerId": 1,              // 推荐人ID（可选）
    "remark": "string",           // 备注（可选）
    "diseaseInfo": "string",      // 既往病史（可选）
    "needToiletAssist": "1"       // 是否需要助排二便（可选）
}
```

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "orderId": 1,
        "orderNo": "ORDER20240101123456",
        "totalAmount": 200.00
    }
}
```

### 5.2 智慧养老设备订单
**接口地址**: `POST /api/order/smart-elderly`

**请求参数**:
```json
{
    "userId": "string",           // 用户ID
    "serviceId": 1,               // 服务ID
    "addressId": 1,               // 地址ID
    "quantity": 1,                // 数量
    "formData": {                 // 表单数据
        "deviceType": "血压计"
    },
    "referrerId": 1,              // 推荐人ID（可选）
    "remark": "string"            // 备注（可选）
}
```

### 5.3 支付订单
**接口地址**: `POST /api/order/pay/{orderId}`

**请求参数**:
```json
{
    "payMethod": "wechat",        // 支付方式
    "openId": "string"            // 用户openID
}
```

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "orderId": 1,
        "orderNo": "ORDER20240101123456",
        "totalAmount": 200.00,
        "paymentParams": {
            "timeStamp": "1640995200",
            "nonceStr": "abc123",
            "package": "prepay_id=wx123456",
            "signType": "MD5",
            "paySign": "signature"
        }
    }
}
```

### 5.4 支付确认
**接口地址**: `POST /api/order/pay_confirm/{orderId}`

**请求参数**:
```json
{
    "transactionId": "string",    // 交易号
    "payMethod": "wechat"         // 支付方式
}
```

### 5.5 取消订单
**接口地址**: `POST /api/order/cancel/{orderId}`

**请求参数**:
```json
{
    "orderId": 1,                 // 订单ID
    "reason": "string"            // 取消原因
}
```

### 5.6 申请退款
**接口地址**: `POST /api/order/refund/{orderId}`

**请求参数**:
```json
{
    "orderId": 1,                 // 订单ID
    "refundAmount": 200.00,       // 退款金额
    "reason": "string"            // 退款原因
}
```

### 5.7 获取订单列表
**接口地址**: `GET /api/order/list`

**请求参数**:
- `userId` (string, 必需): 用户ID
- `page` (int, 可选): 页码，默认1
- `pageSize` (int, 可选): 每页大小，默认10
- `status` (string, 可选): 订单状态筛选

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "list": [
            {
                "id": 1,
                "orderNo": "ORDER20240101123456",
                "serviceName": "上门护理服务",
                "serviceTitle": "上门护理服务",
                "appointmentDate": "2024-01-02",
                "appointmentTime": "09:00",
                "consultTime": "morning",
                "price": 200.00,
                "totalAmount": 200.00,
                "status": 1,
                "payStatus": 1,
                "createdAt": "2024-01-01T12:00:00Z",
                "statusText": "已支付",
                "payStatusText": "已支付",
                "formattedAmount": "¥200.00",
                "formattedDate": "2024-01-01T12:00:00Z"
            }
        ],
        "total": 100,
        "page": 1,
        "pageSize": 10,
        "hasMore": true
    }
}
```

### 5.8 获取订单详情
**接口地址**: `POST /api/order/detail`

**请求参数**:
```json
{
    "orderNo": "ORDER20240101123456"
}
```

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "id": 1,
        "orderNo": "ORDER20240101123456",
        "userId": "user_123456",
        "serviceId": 1,
        "patientId": 1,
        "addressId": 1,
        "appointmentDate": "2024-01-02",
        "appointmentTime": "09:00",
        "serviceName": "上门护理服务",
        "price": 200.00,
        "totalAmount": 200.00,
        "status": 1,
        "payStatus": 1,
        "createdAt": "2024-01-01T12:00:00Z",
        "patientName": "张三",
        "patientPhone": "13800138000",
        "addressInfo": "广东省深圳市南山区科技园南区",
        "serviceTitle": "上门护理服务",
        "formattedPrice": "200.00"
    }
}
```

### 5.9 获取可用时间槽
**接口地址**: `POST /api/order/time_slots`

**请求参数**:
```json
{
    "date": "2024-01-02"
}
```

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "date": "2024-01-02",
        "timeSlots": ["08:00", "09:00", "10:00", "11:00", "14:00", "15:00", "16:00", "17:00", "18:00", "19:00"]
    }
}
```

---

## 6. 支付模块

### 6.1 微信支付回调
**接口地址**: `POST /api/payment/notify`

**说明**: 微信支付结果通知接口，由微信服务器调用

---

## 7. 咨询模块

### 7.1 创建咨询
**接口地址**: `POST /api/consultation/create`

**请求参数**:
```json
{
    "userId": "string",       // 用户ID
    "userName": "string",     // 用户姓名
    "userPhone": "string"     // 用户电话
}
```

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "consultationId": 1,
        "status": "waiting"
    }
}
```

### 7.2 获取咨询消息
**接口地址**: `GET /api/consultation/messages`

**请求参数**:
- `consultationId` (int, 必需): 咨询ID

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "messages": [
            {
                "id": 1,
                "consultationId": 1,
                "senderType": "user",
                "content": "您好，我想咨询一下护理服务",
                "isRead": true,
                "createdAt": "2024-01-01T12:00:00Z"
            }
        ]
    }
}
```

### 7.3 发送咨询消息
**接口地址**: `POST /api/consultation/send`

**请求参数**:
```json
{
    "consultationId": 1,      // 咨询ID
    "content": "string",      // 消息内容
    "senderType": "user"      // 发送者类型：user-用户，admin-管理员
}
```

### 7.4 获取咨询状态
**接口地址**: `GET /api/consultation/status`

**请求参数**:
- `consultationId` (int, 必需): 咨询ID

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "status": "waiting"    // waiting-等待回复，chatting-咨询中，closed-已结束
    }
}
```

### 7.5 关闭咨询
**接口地址**: `POST /api/consultation/close`

**请求参数**:
```json
{
    "consultationId": 1
}
```

### 7.6 获取活跃咨询列表
**接口地址**: `GET /api/consultation/active`

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "consultations": [
            {
                "id": 1,
                "userId": "user_123456",
                "userName": "张三",
                "userPhone": "13800138000",
                "status": "waiting",
                "createdAt": "2024-01-01T12:00:00Z",
                "lastMessage": "2024-01-01T12:00:00Z"
            }
        ]
    }
}
```

### 7.7 获取咨询统计
**接口地址**: `GET /api/consultation/stats`

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "totalConsultations": 100,
        "activeConsultations": 5,
        "closedConsultations": 95,
        "todayConsultations": 10
    }
}
```

### 7.8 获取未读通知
**接口地址**: `GET /api/consultation/notifications`

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "notifications": [
            {
                "id": 1,
                "consultationId": 1,
                "type": "new_message",
                "title": "新消息",
                "content": "您有新的咨询消息",
                "isRead": false,
                "createdAt": "2024-01-01T12:00:00Z"
            }
        ]
    }
}
```

### 7.9 标记通知为已读
**接口地址**: `POST /api/consultation/notification/read`

**请求参数**:
```json
{
    "notificationId": 1
}
```

---

## 8. 推荐系统模块

### 8.1 生成推荐二维码
**接口地址**: `GET /api/referral/qrcode`

**请求参数**:
- `userId` (string, 必需): 用户ID

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "qrcodeUrl": "https://example.com/qrcode.png",
        "promoterCode": "PROMO123456"
    }
}
```

### 8.2 获取推荐报告
**接口地址**: `GET /api/referral/report`

**请求参数**:
- `userId` (string, 必需): 用户ID

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "totalReferrals": 10,
        "totalCommission": 500.00,
        "pendingCommission": 100.00,
        "paidCommission": 400.00
    }
}
```

### 8.3 获取推荐配置
**接口地址**: `GET /api/referral/config`

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "commissionRate": 0.05,
        "minCashoutAmount": 100.00,
        "cashoutEnabled": true
    }
}
```

### 8.4 申请提现
**接口地址**: `POST /api/referral/apply_cashout`

**请求参数**:
```json
{
    "userId": "string",       // 用户ID
    "amount": 100.00,         // 提现金额
    "bankAccount": "string",  // 银行账户
    "bankName": "string"      // 银行名称
}
```

---

## 9. 管理员模块

### 9.1 管理员登录
**接口地址**: `POST /api/admin/login`

**请求参数**:
```json
{
    "username": "string",     // 用户名
    "password": "string"      // 密码
}
```

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "userId": "user_123456",
        "nickName": "管理员",
        "avatarUrl": "https://example.com/avatar.jpg",
        "adminLevel": 1,
        "adminUsername": "admin"
    }
}
```

### 9.2 检查管理员状态
**接口地址**: `GET /api/admin/check-status`

**请求参数**:
- `userId` (string, 必需): 用户ID

### 9.3 获取用户列表
**接口地址**: `GET /api/admin/users`

**请求参数**:
- `page` (int, 可选): 页码
- `pageSize` (int, 可选): 每页大小
- `keyword` (string, 可选): 搜索关键词

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "list": [
            {
                "userId": "user_123456",
                "nickName": "用户昵称",
                "avatarUrl": "https://example.com/avatar.jpg",
                "phone": "13800138000",
                "isAdmin": 0,
                "adminLevel": 0,
                "adminUsername": "",
                "parentAdminId": "",
                "adminCreatedAt": null,
                "createdAt": "2024-01-01T12:00:00Z"
            }
        ],
        "total": 100,
        "page": 1,
        "pageSize": 10,
        "hasMore": true
    }
}
```

### 9.4 获取订单列表
**接口地址**: `GET /api/admin/orders`

**请求参数**:
- `page` (int, 可选): 页码
- `pageSize` (int, 可选): 每页大小
- `status` (int, 可选): 订单状态筛选

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "list": [
            {
                "id": 1,
                "orderNo": "ORDER20240101123456",
                "userId": "user_123456",
                "userNickName": "用户昵称",
                "serviceId": 1,
                "serviceName": "上门护理服务",
                "amount": 200.00,
                "status": 1,
                "statusText": "已支付",
                "createdAt": "2024-01-01T12:00:00Z"
            }
        ],
        "total": 100,
        "page": 1,
        "pageSize": 10,
        "hasMore": true
    }
}
```

### 9.5 设置管理员
**接口地址**: `POST /api/admin/set-admin`

**请求参数**:
```json
{
    "userId": "string",       // 用户ID
    "adminLevel": 1,          // 管理员级别
    "adminUsername": "string", // 管理员用户名
    "adminPassword": "string"  // 管理员密码
}
```

### 9.6 移除管理员
**接口地址**: `POST /api/admin/remove-admin`

**请求参数**:
```json
{
    "userId": "string"        // 用户ID
}
```

### 9.7 获取统计数据
**接口地址**: `GET /api/admin/stats`

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "totalUsers": 1000,
        "totalOrders": 500,
        "totalRevenue": 100000.00,
        "todayUsers": 10,
        "todayOrders": 5,
        "todayRevenue": 1000.00
    }
}
```

### 9.8 修改订单金额
**接口地址**: `POST /api/admin/order/update-amount`

**请求参数**:
```json
{
    "orderId": 1,             // 订单ID
    "newAmount": 150.00,      // 新金额
    "reason": "string"        // 修改原因
}
```

### 9.9 管理员退款
**接口地址**: `POST /api/admin/order/refund`

**请求参数**:
```json
{
    "orderId": 1,             // 订单ID
    "refundAmount": 200.00,   // 退款金额
    "reason": "string",       // 退款原因
    "refundStatus": 2         // 退款状态：1-退款中，2-已退款
}
```

---

## 10. 其他模块

### 10.1 首页初始化
**接口地址**: `GET /api/home/init`

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "banners": [
            {
                "id": 1,
                "title": "专业护理服务",
                "imageUrl": "https://example.com/banner1.jpg",
                "linkUrl": "/service/1"
            }
        ],
        "services": [
            {
                "id": 1,
                "name": "上门护理服务",
                "price": 200.00,
                "imageUrl": "https://example.com/service1.jpg"
            }
        ],
        "announcements": [
            {
                "id": 1,
                "title": "系统维护通知",
                "content": "系统将于今晚进行维护",
                "createdAt": "2024-01-01T12:00:00Z"
            }
        ]
    }
}
```

### 10.2 获取医院列表
**接口地址**: `GET /api/hospital/list`

**请求参数**:
- `city` (string, 可选): 城市筛选

**响应数据**:
```json
{
    "code": 0,
    "data": [
        {
            "id": 1,
            "name": "深圳市人民医院",
            "address": "深圳市罗湖区东门北路1017号",
            "phone": "0755-25533018",
            "imageUrl": "https://example.com/hospital1.jpg",
            "level": "三甲",
            "city": "深圳市"
        }
    ]
}
```

### 10.3 获取医院详情
**接口地址**: `GET /api/hospital/detail/{hospitalId}`

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "id": 1,
        "name": "深圳市人民医院",
        "address": "深圳市罗湖区东门北路1017号",
        "phone": "0755-25533018",
        "imageUrl": "https://example.com/hospital1.jpg",
        "level": "三甲",
        "city": "深圳市",
        "description": "医院简介",
        "departments": [
            {
                "name": "内科",
                "doctors": 50
            }
        ]
    }
}
```

### 10.4 文件上传
**接口地址**: `POST /api/upload`

**请求参数**:
- `file` (file, 必需): 上传文件
- `type` (string, 可选): 文件类型

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "fileId": "file_123456",
        "fileName": "image.jpg",
        "fileUrl": "https://example.com/files/image.jpg",
        "fileSize": 1024000,
        "fileType": "image/jpeg"
    }
}
```

### 10.5 获取文件列表
**接口地址**: `GET /api/files`

**请求参数**:
- `page` (int, 可选): 页码
- `pageSize` (int, 可选): 每页大小
- `type` (string, 可选): 文件类型筛选

### 10.6 删除文件
**接口地址**: `POST /api/file/delete`

**请求参数**:
```json
{
    "fileId": "string"        // 文件ID
}
```

### 10.7 生成二维码
**接口地址**: `POST /api/qrcode/generate`

**请求参数**:
```json
{
    "content": "string",      // 二维码内容
    "size": 200               // 二维码尺寸
}
```

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "qrcodeUrl": "https://example.com/qrcode.png"
    }
}
```

### 10.8 生成Base64二维码
**接口地址**: `POST /api/qrcode/generate_base64`

**请求参数**:
```json
{
    "content": "string",      // 二维码内容
    "size": 200               // 二维码尺寸
}
```

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "qrcodeBase64": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."
    }
}
```

---

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| -1 | 通用错误 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 注意事项

1. 所有时间格式均为ISO 8601格式（UTC时间）
2. 金额单位为元，保留两位小数
3. 分页参数：page从1开始，pageSize最大50
4. 文件上传支持图片、文档等格式
5. 微信相关接口需要正确的AppID和AppSecret配置
6. 支付接口需要配置微信支付商户信息
7. 所有接口都支持CORS跨域请求
8. 建议在生产环境中使用HTTPS协议

---

## 更新日志

### v1.3.0 (2024-01-01)
- 新增智慧养老设备订单接口
- 优化订单状态管理
- 完善支付流程
- 增强错误处理机制

### v1.2.0 (2023-12-01)
- 新增咨询系统模块
- 完善推荐系统功能
- 优化管理员权限管理
- 增强文件上传功能

### v1.1.0 (2023-11-01)
- 新增订单管理模块
- 完善用户管理功能
- 集成微信支付
- 优化API响应格式

### v1.0.0 (2023-10-01)
- 初始版本发布
- 基础用户认证功能
- 服务管理模块
- 基础API框架
