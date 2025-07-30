# 首页初始化接口文档

## 接口信息

- **接口地址**: `GET/POST /api/home/init`
- **请求方式**: GET/POST
- **功能**: 获取首页初始化数据（轮播图、导航、服务项、医院列表）

## 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| longitude | float | 否 | 用户经度，用于医院距离排序 |
| latitude | float | 否 | 用户纬度，用于医院距离排序 |
| limit | int | 否 | 医院列表限制数量，默认10 |

## 响应格式

### 成功响应

```json
{
  "code": 0,
  "data": {
    "banners": [
      {
        "id": 1,
        "title": "轮播图标题",
        "image": "https://example.com/banner1.jpg",
        "url": "https://example.com/link1",
        "sort": 1,
        "status": 1
      }
    ],
    "navigations": [
      {
        "id": 1,
        "title": "导航标题",
        "icon": "https://example.com/icon1.png",
        "url": "https://example.com/nav1",
        "sort": 1,
        "status": 1
      }
    ],
    "services": [
      {
        "id": 1,
        "title": "服务标题",
        "description": "服务描述",
        "icon": "https://example.com/service1.png",
        "url": "https://example.com/service1",
        "sort": 1,
        "status": 1
      }
    ],
    "hospitals": [
      {
        "id": 1,
        "name": "深圳市人民医院",
        "logo": "https://example.com/hospital1.jpg",
        "address": "深圳市罗湖区东门北路1017号",
        "phone": "0755-25533018",
        "description": "医院描述",
        "level": "三级甲等",
        "type": "综合医院",
        "longitude": 114.0579,
        "latitude": 22.5431,
        "sort": 1,
        "status": 1
      }
    ]
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
// 获取首页数据
wx.request({
  url: 'http://your-server.com/api/home/init',
  method: 'GET',
  data: {
    longitude: 114.0579,
    latitude: 22.5431,
    limit: 10
  },
  success: (res) => {
    if (res.data.code === 0) {
      const data = res.data.data;
      console.log('轮播图:', data.banners);
      console.log('导航:', data.navigations);
      console.log('服务:', data.services);
      console.log('医院:', data.hospitals);
    }
  }
});
```

## 数据库表结构

### 轮播图表 (Banners)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | INT | 主键，自增 |
| title | VARCHAR(200) | 轮播图标题 |
| image | VARCHAR(500) | 轮播图片URL |
| url | VARCHAR(500) | 跳转链接 |
| sort | INT | 排序 |
| status | INT | 状态：1-启用，0-禁用 |
| createdAt | DATETIME | 创建时间 |
| updatedAt | DATETIME | 更新时间 |

### 导航表 (Navigations)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | INT | 主键，自增 |
| title | VARCHAR(100) | 导航标题 |
| icon | VARCHAR(500) | 导航图标URL |
| url | VARCHAR(500) | 跳转链接 |
| sort | INT | 排序 |
| status | INT | 状态：1-启用，0-禁用 |
| createdAt | DATETIME | 创建时间 |
| updatedAt | DATETIME | 更新时间 |

### 服务表 (Services)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | INT | 主键，自增 |
| title | VARCHAR(100) | 服务标题 |
| description | VARCHAR(500) | 服务描述 |
| icon | VARCHAR(500) | 服务图标URL |
| url | VARCHAR(500) | 跳转链接 |
| sort | INT | 排序 |
| status | INT | 状态：1-启用，0-禁用 |
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

## 注意事项

1. **位置排序**: 当提供用户位置时，医院列表会按距离排序
2. **数据缓存**: 建议前端对首页数据进行缓存，减少请求频率
3. **图片优化**: 轮播图和图标建议使用CDN加速
4. **状态过滤**: 只返回状态为启用的数据
5. **排序规则**: 按sort字段升序排列
6. **医院距离**: 使用Haversine公式计算用户与医院的距离 