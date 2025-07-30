# 推荐系统接口文档

## 接口概览

本文档包含以下推荐系统相关接口：
1. **获取推广二维码** - `GET /api/referral/qrcode`
2. **获取推荐报告** - `GET /api/referral/report`
3. **获取推荐配置** - `GET /api/referral/config`
4. **申请佣金提现** - `POST /api/referral/apply_cashout`

## 1. 获取推广二维码

### 接口信息
- **接口地址**: `GET /api/referral/qrcode`
- **请求方式**: GET
- **功能**: 获取用户专属推广二维码

### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| userId | int | 是 | 用户ID |

### 响应格式
```json
{
  "code": 0,
  "data": {
    "qrcodeUrl": "https://example.com/qrcode/user_1.png",
    "qrcodeContent": "https://example.com/register?ref=user_1",
    "shareTitle": "邀请好友注册，获得佣金奖励",
    "shareDesc": "注册即送50元优惠券，下单还有佣金返现",
    "shareImage": "https://example.com/share_image.jpg"
  }
}
```

## 2. 获取推荐报告

### 接口信息
- **接口地址**: `GET /api/referral/report`
- **请求方式**: GET
- **功能**: 获取用户的推荐人及下单记录

### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| userId | int | 是 | 用户ID |
| page | int | 否 | 页码，默认1 |
| pageSize | int | 否 | 每页数量，默认10，最大50 |

### 响应格式
```json
{
  "code": 0,
  "data": {
    "referrer": {
      "userId": 2,
      "nickName": "推荐人昵称",
      "avatarUrl": "https://example.com/avatar.jpg"
    },
    "commissionStats": {
      "totalCommission": 150.00,
      "availableCommission": 100.00,
      "frozenCommission": 50.00,
      "totalOrders": 5,
      "totalAmount": 1500.00
    },
    "orderList": [
      {
        "id": 1,
        "orderNo": "202401150001",
        "referredUserId": 3,
        "referredUserNickName": "被推荐人昵称",
        "serviceTitle": "全面体检套餐",
        "amount": 299.00,
        "commissionAmount": 29.90,
        "commissionRate": 0.1,
        "status": "paid",
        "createdAt": "2024-01-01T12:00:00Z"
      }
    ],
    "total": 5,
    "page": 1,
    "pageSize": 10,
    "hasMore": false
  }
}
```

## 3. 获取推荐配置

### 接口信息
- **接口地址**: `GET /api/referral/config`
- **请求方式**: GET
- **功能**: 获取推荐返佣规则说明

### 请求参数
无

### 响应格式
```json
{
  "code": 0,
  "data": {
    "commissionRate": 0.1,
    "minCashoutAmount": 50.00,
    "maxCashoutAmount": 1000.00,
    "cashoutFee": 0.02,
    "rules": [
      "推荐好友注册，获得50元优惠券",
      "好友下单，获得订单金额10%的佣金",
      "佣金满50元可申请提现",
      "提现手续费2%",
      "单次提现最高1000元"
    ],
    "levels": [
      {
        "level": 1,
        "name": "普通推荐人",
        "commissionRate": 0.1,
        "description": "推荐好友下单，获得10%佣金"
      },
      {
        "level": 2,
        "name": "高级推荐人",
        "commissionRate": 0.15,
        "description": "累计推荐10人，佣金提升至15%"
      }
    ]
  }
}
```

## 4. 申请佣金提现

### 接口信息
- **接口地址**: `POST /api/referral/apply_cashout`
- **请求方式**: POST
- **功能**: 申请佣金提现

### 请求参数
```json
{
  "userId": 1,
  "amount": 100.00,
  "bankName": "中国银行",
  "bankAccount": "6222021234567890123",
  "accountName": "张三",
  "phone": "13800138000"
}
```

### 响应格式
```json
{
  "code": 0,
  "data": {
    "cashoutId": 1,
    "cashoutNo": "TX202401150001",
    "amount": 100.00,
    "fee": 2.00,
    "actualAmount": 98.00,
    "status": "pending",
    "message": "提现申请提交成功，预计1-3个工作日到账"
  }
}
```

## 佣金计算规则

### 基础规则
- 推荐好友注册：获得50元优惠券
- 好友下单：获得订单金额10%的佣金
- 佣金满50元可申请提现
- 提现手续费2%
- 单次提现最高1000元

### 等级规则
| 等级 | 名称 | 佣金比例 | 条件 |
|------|------|----------|------|
| 1 | 普通推荐人 | 10% | 无 |
| 2 | 高级推荐人 | 15% | 累计推荐10人 |
| 3 | 金牌推荐人 | 20% | 累计推荐50人 |

## 使用示例

### 微信小程序端

```javascript
// 获取推广二维码
wx.request({
  url: 'http://your-server.com/api/referral/qrcode',
  method: 'GET',
  data: { userId: 1 },
  success: (res) => {
    if (res.data.code === 0) {
      const data = res.data.data;
      console.log('推广二维码:', data);
      
      // 显示二维码
      this.setData({
        qrcodeUrl: data.qrcodeUrl,
        shareTitle: data.shareTitle,
        shareDesc: data.shareDesc,
        shareImage: data.shareImage
      });
    }
  }
});

// 获取推荐报告
wx.request({
  url: 'http://your-server.com/api/referral/report',
  method: 'GET',
  data: {
    userId: 1,
    page: 1,
    pageSize: 10
  },
  success: (res) => {
    if (res.data.code === 0) {
      const data = res.data.data;
      console.log('推荐报告:', data);
      
      // 显示佣金统计
      this.setData({
        commissionStats: data.commissionStats,
        orderList: data.orderList
      });
    }
  }
});

// 获取推荐配置
wx.request({
  url: 'http://your-server.com/api/referral/config',
  method: 'GET',
  success: (res) => {
    if (res.data.code === 0) {
      const config = res.data.data;
      console.log('推荐配置:', config);
      
      // 显示规则说明
      this.setData({
        rules: config.rules,
        levels: config.levels
      });
    }
  }
});

// 申请佣金提现
wx.request({
  url: 'http://your-server.com/api/referral/apply_cashout',
  method: 'POST',
  header: { 'Content-Type': 'application/json' },
  data: {
    userId: 1,
    amount: 100.00,
    bankName: '中国银行',
    bankAccount: '6222021234567890123',
    accountName: '张三',
    phone: '13800138000'
  },
  success: (res) => {
    if (res.data.code === 0) {
      console.log('提现申请成功:', res.data.data);
      wx.showToast({
        title: '提现申请提交成功',
        icon: 'success'
      });
    }
  }
});

// 分享推广链接
onShareAppMessage() {
  return {
    title: this.data.shareTitle,
    desc: this.data.shareDesc,
    path: '/pages/register?ref=' + this.data.userId,
    imageUrl: this.data.shareImage
  };
}
```

## 数据库表结构

### 推荐关系表 (Referrals)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | INT | 主键，自增 |
| referrerUserId | INT | 推荐人用户ID |
| referredUserId | INT | 被推荐人用户ID |
| level | INT | 推荐等级 |
| status | INT | 状态：1-正常，0-已取消 |
| createdAt | DATETIME | 创建时间 |
| updatedAt | DATETIME | 更新时间 |

### 佣金记录表 (Commissions)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | INT | 主键，自增 |
| referrerUserId | INT | 推荐人用户ID |
| orderId | INT | 订单ID |
| orderNo | VARCHAR(50) | 订单号 |
| amount | DECIMAL(10,2) | 订单金额 |
| commissionAmount | DECIMAL(10,2) | 佣金金额 |
| commissionRate | DECIMAL(5,4) | 佣金比例 |
| status | VARCHAR(20) | 状态：pending, confirmed, cancelled |
| createdAt | DATETIME | 创建时间 |
| updatedAt | DATETIME | 更新时间 |

### 提现记录表 (Cashouts)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | INT | 主键，自增 |
| cashoutNo | VARCHAR(50) | 提现单号 |
| userId | INT | 用户ID |
| amount | DECIMAL(10,2) | 提现金额 |
| fee | DECIMAL(10,2) | 手续费 |
| actualAmount | DECIMAL(10,2) | 实际到账金额 |
| bankName | VARCHAR(100) | 银行名称 |
| bankAccount | VARCHAR(50) | 银行账号 |
| accountName | VARCHAR(100) | 账户姓名 |
| phone | VARCHAR(20) | 联系电话 |
| status | VARCHAR(20) | 状态：pending, approved, rejected, completed |
| rejectReason | VARCHAR(500) | 拒绝原因 |
| createdAt | DATETIME | 创建时间 |
| updatedAt | DATETIME | 更新时间 |

## 注意事项

1. **推荐关系**: 一个用户只能有一个推荐人
2. **佣金计算**: 按订单金额和佣金比例计算
3. **提现限制**: 最低50元，最高1000元，手续费2%
4. **状态管理**: 佣金状态：待确认 -> 已确认 -> 可提现
5. **等级升级**: 根据推荐人数自动升级等级
6. **数据安全**: 银行账号等敏感信息需要加密存储
7. **审核机制**: 提现申请需要后台审核
8. **通知机制**: 佣金到账、提现成功等会发送通知
9. **防刷机制**: 防止恶意刷佣金的行为
10. **数据统计**: 记录推荐数据，便于分析 