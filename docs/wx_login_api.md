# 微信登录接口文档

## 接口信息

- **接口地址**: `POST /api/wx/login`
- **请求方式**: POST
- **功能**: 微信小程序登录，获取用户信息并保存到数据库

## 请求参数

```json
{
  "code": "微信小程序登录code",
  "userInfo": {
    "nickName": "用户昵称",
    "avatarUrl": "用户头像URL",
    "gender": 1,
    "country": "国家",
    "province": "省份",
    "city": "城市",
    "language": "语言"
  }
}
```

### 参数说明

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| code | string | 是 | 微信小程序登录code |
| userInfo.nickName | string | 否 | 用户昵称 |
| userInfo.avatarUrl | string | 否 | 用户头像URL |
| userInfo.gender | int | 否 | 性别：0-未知，1-男，2-女 |
| userInfo.country | string | 否 | 国家 |
| userInfo.province | string | 否 | 省份 |
| userInfo.city | string | 否 | 城市 |
| userInfo.language | string | 否 | 语言 |

## 响应格式

### 成功响应

```json
{
  "code": 0,
  "data": {
    "userId": 1,
    "openId": "用户openId",
    "unionId": "用户unionId",
    "nickName": "用户昵称",
    "avatarUrl": "用户头像URL",
    "gender": 1,
    "country": "国家",
    "province": "省份",
    "city": "城市",
    "language": "语言",
    "lastLoginAt": "2024-01-01T12:00:00Z",
    "createdAt": "2024-01-01T12:00:00Z",
    "updatedAt": "2024-01-01T12:00:00Z"
  }
}
```

### 失败响应

```json
{
  "code": -1,
  "errorMsg": "错误信息"
}
```

## 使用示例

### 微信小程序端

```javascript
// 调用微信登录
wx.login({
  success: (res) => {
    if (res.code) {
      // 获取用户信息
      wx.getUserInfo({
        success: (userInfo) => {
          // 发送登录请求
          wx.request({
            url: 'http://your-server.com/api/wx/login',
            method: 'POST',
            header: { 'Content-Type': 'application/json' },
            data: {
              code: res.code,
              userInfo: userInfo.userInfo
            },
            success: (response) => {
              if (response.data.code === 0) {
                console.log('登录成功:', response.data.data);
                // 保存用户信息到本地
                wx.setStorageSync('userInfo', response.data.data);
              } else {
                console.error('登录失败:', response.data.errorMsg);
              }
            }
          });
        }
      });
    }
  }
});
```

## 数据库表结构

### 用户表 (Users)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | INT | 主键，自增 |
| openId | VARCHAR(100) | 微信openId，唯一索引 |
| unionId | VARCHAR(100) | 微信unionId |
| nickName | VARCHAR(100) | 用户昵称 |
| avatarUrl | VARCHAR(500) | 用户头像URL |
| gender | INT | 性别：0-未知，1-男，2-女 |
| country | VARCHAR(50) | 国家 |
| province | VARCHAR(50) | 省份 |
| city | VARCHAR(50) | 城市 |
| language | VARCHAR(20) | 语言 |
| sessionKey | VARCHAR(100) | 微信sessionKey |
| lastLoginAt | DATETIME | 最后登录时间 |
| createdAt | DATETIME | 创建时间 |
| updatedAt | DATETIME | 更新时间 |

## 注意事项

1. **首次登录**: 新用户首次登录时会创建用户记录
2. **重复登录**: 已存在用户会更新登录时间和用户信息
3. **SessionKey**: 系统会保存微信返回的sessionKey用于后续接口调用
4. **数据安全**: openId和sessionKey等敏感信息不会返回给前端
5. **错误处理**: 包含微信API调用失败、数据库操作失败等错误处理
6. **配置要求**: 需要在环境变量中配置微信AppID和AppSecret 