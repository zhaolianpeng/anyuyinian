# 客服、医院相关接口文档

## 接口概览

本文档包含以下客服、医院相关接口：
1. **提交用户咨询问题** - `POST /api/kefu/send_msg`
2. **常见问题列表** - `GET /api/kefu/faq`
3. **获取医院列表** - `GET /api/hospital/list`
4. **获取医院详情** - `GET /api/hospital/detail/:id`

## 1. 提交用户咨询问题

### 接口信息
- **接口地址**: `POST /api/kefu/send_msg`
- **请求方式**: POST
- **功能**: 提交用户咨询问题

### 请求参数
```json
{
  "userId": 1,
  "userName": "张三",
  "userAvatar": "https://example.com/avatar.jpg",
  "content": "我想咨询一下体检套餐的具体内容",
  "images": ["https://example.com/image1.jpg"]
}
```

### 响应格式
```json
{
  "code": 0,
  "data": {
    "messageId": 1,
    "message": "消息发送成功，客服将尽快回复您"
  }
}
```

## 2. 常见问题列表

### 接口信息
- **接口地址**: `GET /api/kefu/faq`
- **请求方式**: GET
- **功能**: 获取常见问题列表，支持分页和分类筛选

### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| category | string | 否 | 问题分类 |
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
        "question": "如何预约体检服务？",
        "answer": "您可以在首页选择体检服务，填写个人信息后提交订单，支付成功后即可预约成功。",
        "category": "预约服务",
        "sort": 1,
        "status": 1,
        "viewCount": 0,
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

## 3. 获取医院列表

### 接口信息
- **接口地址**: `GET /api/hospital/list`
- **请求方式**: GET
- **功能**: 获取可选医院列表，支持位置排序

### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| longitude | float | 否 | 用户经度 |
| latitude | float | 否 | 用户纬度 |
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
        "name": "深圳市人民医院",
        "logo": "https://example.com/hospital1.jpg",
        "address": "深圳市罗湖区东门北路1017号",
        "phone": "0755-25533018",
        "description": "深圳市人民医院是一所集医疗、教学、科研、预防、保健为一体的现代化综合性三级甲等医院。",
        "level": "三级甲等",
        "type": "综合医院",
        "longitude": 114.0579,
        "latitude": 22.5431,
        "sort": 1,
        "status": 1,
        "createdAt": "2024-01-01T12:00:00Z",
        "updatedAt": "2024-01-01T12:00:00Z"
      }
    ],
    "total": 5,
    "page": 1,
    "pageSize": 10,
    "hasMore": false
  }
}
```

## 4. 获取医院详情

### 接口信息
- **接口地址**: `GET /api/hospital/detail/:id`
- **请求方式**: GET
- **功能**: 获取医院详情，包含导航信息

### 路径参数
- `id`: 医院ID

### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| userLongitude | float | 否 | 用户经度 |
| userLatitude | float | 否 | 用户纬度 |

### 响应格式
```json
{
  "code": 0,
  "data": {
    "hospital": {
      "id": 1,
      "name": "深圳市人民医院",
      "logo": "https://example.com/hospital1.jpg",
      "address": "深圳市罗湖区东门北路1017号",
      "phone": "0755-25533018",
      "description": "深圳市人民医院是一所集医疗、教学、科研、预防、保健为一体的现代化综合性三级甲等医院。",
      "level": "三级甲等",
      "type": "综合医院",
      "longitude": 114.0579,
      "latitude": 22.5431,
      "sort": 1,
      "status": 1,
      "createdAt": "2024-01-01T12:00:00Z",
      "updatedAt": "2024-01-01T12:00:00Z"
    },
    "navigation": {
      "distance": 2.5,
      "duration": 5,
      "routeType": "driving",
      "routePoints": [
        {
          "latitude": 22.5431,
          "longitude": 114.0579
        },
        {
          "latitude": 22.5431,
          "longitude": 114.0579
        }
      ]
    }
  }
}
```

## 使用示例

### 微信小程序端

```javascript
// 提交用户咨询问题
wx.request({
  url: 'http://your-server.com/api/kefu/send_msg',
  method: 'POST',
  header: { 'Content-Type': 'application/json' },
  data: {
    userId: 1,
    userName: '张三',
    userAvatar: 'https://example.com/avatar.jpg',
    content: '我想咨询一下体检套餐的具体内容',
    images: ['https://example.com/image1.jpg']
  },
  success: (res) => {
    if (res.data.code === 0) {
      console.log('消息发送成功:', res.data.data);
    }
  }
});

// 获取常见问题列表
wx.request({
  url: 'http://your-server.com/api/kefu/faq',
  method: 'GET',
  data: {
    category: '预约服务',
    page: 1,
    pageSize: 10
  },
  success: (res) => {
    if (res.data.code === 0) {
      console.log('FAQ列表:', res.data.data);
    }
  }
});

// 获取医院列表
wx.request({
  url: 'http://your-server.com/api/hospital/list',
  method: 'GET',
  data: {
    longitude: 114.0579,
    latitude: 22.5431,
    page: 1,
    pageSize: 10
  },
  success: (res) => {
    if (res.data.code === 0) {
      console.log('医院列表:', res.data.data);
    }
  }
});

// 获取医院详情
wx.request({
  url: 'http://your-server.com/api/hospital/detail/1',
  method: 'GET',
  data: {
    userLongitude: 114.0579,
    userLatitude: 22.5431
  },
  success: (res) => {
    if (res.data.code === 0) {
      console.log('医院详情:', res.data.data);
      
      // 显示导航信息
      const navigation = res.data.data.navigation;
      if (navigation.distance) {
        console.log(`距离: ${navigation.distance}公里`);
        console.log(`预计时间: ${navigation.duration}分钟`);
      }
    }
  }
});
```

## 数据库表结构

### 客服消息表 (KefuMessages)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | INT | 主键，自增 |
| userId | INT | 用户ID |
| userName | VARCHAR(100) | 用户姓名 |
| userAvatar | VARCHAR(500) | 用户头像 |
| type | INT | 消息类型：1-用户消息，2-客服回复 |
| content | TEXT | 消息内容 |
| images | TEXT | 图片（JSON数组） |
| status | INT | 状态：0-未读，1-已读，2-已回复 |
| replyContent | TEXT | 回复内容 |
| replyTime | DATETIME | 回复时间 |
| replyUserId | INT | 回复人ID |
| createdAt | DATETIME | 创建时间 |
| updatedAt | DATETIME | 更新时间 |

### 常见问题表 (Faqs)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | INT | 主键，自增 |
| question | VARCHAR(500) | 问题 |
| answer | TEXT | 答案 |
| category | VARCHAR(100) | 分类 |
| sort | INT | 排序 |
| status | INT | 状态：1-启用，0-禁用 |
| viewCount | INT | 查看次数 |
| createdAt | DATETIME | 创建时间 |
| updatedAt | DATETIME | 更新时间 |

### 医院表 (Hospitals)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | INT | 主键，自增 |
| name | VARCHAR(200) | 医院名称 |
| logo | VARCHAR(500) | 医院logo |
| address | VARCHAR(500) | 医院地址 |
| phone | VARCHAR(50) | 联系电话 |
| description | TEXT | 医院描述 |
| level | VARCHAR(50) | 医院等级 |
| type | VARCHAR(50) | 医院类型 |
| longitude | DECIMAL(10,6) | 经度 |
| latitude | DECIMAL(10,6) | 纬度 |
| sort | INT | 排序 |
| status | INT | 状态：1-启用，0-禁用 |
| createdAt | DATETIME | 创建时间 |
| updatedAt | DATETIME | 更新时间 |

## FAQ分类说明

| 分类名 | 说明 |
|--------|------|
| 预约服务 | 预约相关常见问题 |
| 体检须知 | 体检前注意事项 |
| 体检报告 | 体检报告相关问题 |
| 支付问题 | 支付相关常见问题 |
| 体检项目 | 体检项目相关问题 |

## 医院等级说明

| 等级 | 说明 |
|------|------|
| 三级甲等 | 最高等级医院 |
| 三级乙等 | 三级医院 |
| 二级甲等 | 二级医院 |
| 二级乙等 | 二级医院 |
| 一级甲等 | 一级医院 |
| 一级乙等 | 一级医院 |

## 医院类型说明

| 类型 | 说明 |
|------|------|
| 综合医院 | 综合性医疗机构 |
| 专科医院 | 专科医疗机构 |
| 中医医院 | 中医医疗机构 |
| 妇幼保健院 | 妇幼保健机构 |
| 社区卫生服务中心 | 社区医疗机构 |

## 注意事项

1. **消息状态**: 0-未读，1-已读，2-已回复
2. **消息类型**: 1-用户消息，2-客服回复
3. **FAQ分类**: 支持按分类筛选常见问题
4. **医院等级**: 按医院等级分类管理
5. **位置排序**: 根据用户位置计算距离，按距离排序
6. **导航信息**: 包含距离、预计时间、路线类型、路线点等信息
7. **距离计算**: 使用Haversine公式计算两点间距离
8. **数据验证**: 所有接口都包含基本的数据验证
9. **分页处理**: 支持分页查询，默认每页10条，最大50条
10. **图片支持**: 客服消息支持图片上传 