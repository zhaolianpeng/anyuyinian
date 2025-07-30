# 用户信息相关接口文档

## 接口概览

本文档包含以下用户相关接口：
1. **获取用户信息** - `GET /api/user/info`
2. **绑定手机号** - `POST /api/user/bind_phone`
3. **地址管理** - `GET/POST/PUT/DELETE /api/user/address`
4. **就诊人管理** - `GET/POST/PUT/DELETE /api/user/patient`

## 1. 获取用户信息

### 接口信息
- **接口地址**: `GET /api/user/info`
- **请求方式**: GET
- **功能**: 获取用户个人信息

### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| userId | int | 是 | 用户ID |

### 响应格式
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "openId": "用户openId",
    "unionId": "用户unionId",
    "nickName": "用户昵称",
    "avatarUrl": "用户头像URL",
    "gender": 1,
    "country": "国家",
    "province": "省份",
    "city": "城市",
    "language": "语言",
    "phone": "13800138000",
    "lastLoginAt": "2024-01-01T12:00:00Z",
    "createdAt": "2024-01-01T12:00:00Z",
    "updatedAt": "2024-01-01T12:00:00Z"
  }
}
```

## 2. 绑定手机号

### 接口信息
- **接口地址**: `POST /api/user/bind_phone`
- **请求方式**: POST
- **功能**: 用户绑定手机号

### 请求参数
```json
{
  "userId": 1,
  "phone": "13800138000",
  "code": "123456"
}
```

### 响应格式
```json
{
  "code": 0,
  "data": {
    "message": "手机号绑定成功"
  }
}
```

## 3. 地址管理

### 3.1 获取地址列表
- **接口地址**: `GET /api/user/address`
- **请求参数**: `userId=1`

### 3.2 添加地址
- **接口地址**: `POST /api/user/address`
- **请求参数**:
```json
{
  "userId": 1,
  "name": "张三",
  "phone": "13800138000",
  "province": "广东省",
  "city": "深圳市",
  "district": "罗湖区",
  "address": "东门北路1017号",
  "isDefault": true
}
```

### 3.3 更新地址
- **接口地址**: `PUT /api/user/address`
- **请求参数**:
```json
{
  "id": 1,
  "userId": 1,
  "name": "张三",
  "phone": "13800138000",
  "province": "广东省",
  "city": "深圳市",
  "district": "罗湖区",
  "address": "东门北路1017号",
  "isDefault": true
}
```

### 3.4 删除地址
- **接口地址**: `DELETE /api/user/address`
- **请求参数**: `id=1&userId=1`

### 地址响应格式
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "userId": 1,
        "name": "张三",
        "phone": "13800138000",
        "province": "广东省",
        "city": "深圳市",
        "district": "罗湖区",
        "address": "东门北路1017号",
        "isDefault": true,
        "createdAt": "2024-01-01T12:00:00Z",
        "updatedAt": "2024-01-01T12:00:00Z"
      }
    ]
  }
}
```

## 4. 就诊人管理

### 4.1 获取就诊人列表
- **接口地址**: `GET /api/user/patient`
- **请求参数**: `userId=1`

### 4.2 添加就诊人
- **接口地址**: `POST /api/user/patient`
- **请求参数**:
```json
{
  "userId": 1,
  "name": "张三",
  "idCard": "440301199001011234",
  "phone": "13800138000",
  "gender": 1,
  "birthday": "1990-01-01",
  "relationship": "本人",
  "isDefault": true
}
```

### 4.3 更新就诊人
- **接口地址**: `PUT /api/user/patient`
- **请求参数**:
```json
{
  "id": 1,
  "userId": 1,
  "name": "张三",
  "idCard": "440301199001011234",
  "phone": "13800138000",
  "gender": 1,
  "birthday": "1990-01-01",
  "relationship": "本人",
  "isDefault": true
}
```

### 4.4 删除就诊人
- **接口地址**: `DELETE /api/user/patient`
- **请求参数**: `id=1&userId=1`

### 就诊人响应格式
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "userId": 1,
        "name": "张三",
        "idCard": "440301199001011234",
        "phone": "13800138000",
        "gender": 1,
        "birthday": "1990-01-01",
        "relationship": "本人",
        "isDefault": true,
        "createdAt": "2024-01-01T12:00:00Z",
        "updatedAt": "2024-01-01T12:00:00Z"
      }
    ]
  }
}
```

## 使用示例

### 微信小程序端

```javascript
// 获取用户信息
wx.request({
  url: 'http://your-server.com/api/user/info',
  method: 'GET',
  data: { userId: 1 },
  success: (res) => {
    if (res.data.code === 0) {
      console.log('用户信息:', res.data.data);
    }
  }
});

// 绑定手机号
wx.request({
  url: 'http://your-server.com/api/user/bind_phone',
  method: 'POST',
  header: { 'Content-Type': 'application/json' },
  data: {
    userId: 1,
    phone: '13800138000',
    code: '123456'
  },
  success: (res) => {
    if (res.data.code === 0) {
      console.log('绑定成功:', res.data.data);
    }
  }
});

// 获取地址列表
wx.request({
  url: 'http://your-server.com/api/user/address',
  method: 'GET',
  data: { userId: 1 },
  success: (res) => {
    if (res.data.code === 0) {
      console.log('地址列表:', res.data.data);
    }
  }
});

// 添加地址
wx.request({
  url: 'http://your-server.com/api/user/address',
  method: 'POST',
  header: { 'Content-Type': 'application/json' },
  data: {
    userId: 1,
    name: '张三',
    phone: '13800138000',
    province: '广东省',
    city: '深圳市',
    district: '罗湖区',
    address: '东门北路1017号',
    isDefault: true
  },
  success: (res) => {
    if (res.data.code === 0) {
      console.log('添加成功:', res.data.data);
    }
  }
});

// 获取就诊人列表
wx.request({
  url: 'http://your-server.com/api/user/patient',
  method: 'GET',
  data: { userId: 1 },
  success: (res) => {
    if (res.data.code === 0) {
      console.log('就诊人列表:', res.data.data);
    }
  }
});

// 添加就诊人
wx.request({
  url: 'http://your-server.com/api/user/patient',
  method: 'POST',
  header: { 'Content-Type': 'application/json' },
  data: {
    userId: 1,
    name: '张三',
    idCard: '440301199001011234',
    phone: '13800138000',
    gender: 1,
    birthday: '1990-01-01',
    relationship: '本人',
    isDefault: true
  },
  success: (res) => {
    if (res.data.code === 0) {
      console.log('添加成功:', res.data.data);
    }
  }
});
```

## 数据库表结构

### 用户地址表 (UserAddresses)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | INT | 主键，自增 |
| userId | INT | 用户ID |
| name | VARCHAR(100) | 收货人姓名 |
| phone | VARCHAR(20) | 联系电话 |
| province | VARCHAR(50) | 省份 |
| city | VARCHAR(50) | 城市 |
| district | VARCHAR(50) | 区县 |
| address | VARCHAR(500) | 详细地址 |
| isDefault | BOOLEAN | 是否默认地址 |
| status | INT | 状态：1-正常，0-已删除 |
| createdAt | DATETIME | 创建时间 |
| updatedAt | DATETIME | 更新时间 |

### 就诊人表 (Patients)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | INT | 主键，自增 |
| userId | INT | 用户ID |
| name | VARCHAR(100) | 就诊人姓名 |
| idCard | VARCHAR(20) | 身份证号 |
| phone | VARCHAR(20) | 联系电话 |
| gender | INT | 性别：1-男，2-女 |
| birthday | DATE | 出生日期 |
| relationship | VARCHAR(50) | 与用户关系 |
| isDefault | BOOLEAN | 是否默认就诊人 |
| status | INT | 状态：1-正常，0-已删除 |
| createdAt | DATETIME | 创建时间 |
| updatedAt | DATETIME | 更新时间 |

## 注意事项

1. **默认设置**: 每个用户只能有一个默认地址和一个默认就诊人
2. **数据验证**: 身份证号、手机号等字段会进行格式验证
3. **软删除**: 删除操作采用软删除，不会物理删除数据
4. **权限控制**: 用户只能操作自己的地址和就诊人信息
5. **数据关联**: 地址和就诊人与用户ID关联
6. **字段必填**: 姓名、电话等关键字段为必填项
7. **格式要求**: 身份证号、手机号等需要符合格式要求
8. **关系类型**: 就诊人关系包括：本人、父母、子女、配偶、其他 