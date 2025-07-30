# 服务相关接口文档

## 接口概览

本文档包含以下服务相关接口：
1. **服务列表** - `GET /api/service/list`
2. **服务详情** - `GET /api/service/detail/:id`
3. **服务表单配置** - `GET /api/service/form_config/:id`

## 1. 服务列表接口

### 接口信息
- **接口地址**: `GET /api/service/list`
- **请求方式**: GET
- **功能**: 获取服务列表，支持分页和分类筛选

### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| category | string | 否 | 服务分类 |
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
        "title": "全面体检套餐",
        "description": "包含血常规、尿常规、心电图等多项检查",
        "price": 299.00,
        "originalPrice": 399.00,
        "category": "体检套餐",
        "images": ["https://example.com/service1.jpg"],
        "formConfig": "{\"fields\":[...]}",
        "status": 1,
        "sort": 1,
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

## 2. 服务详情接口

### 接口信息
- **接口地址**: `GET /api/service/detail/:id`
- **请求方式**: GET
- **功能**: 获取服务详细信息

### 路径参数
- `id`: 服务ID

### 响应格式
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "title": "全面体检套餐",
    "description": "包含血常规、尿常规、心电图等多项检查",
    "detail": "详细的服务介绍内容...",
    "price": 299.00,
    "originalPrice": 399.00,
    "category": "体检套餐",
    "images": [
      "https://example.com/service1.jpg",
      "https://example.com/service1_detail.jpg"
    ],
    "formConfig": "{\"fields\":[...]}",
    "status": 1,
    "sort": 1,
    "viewCount": 100,
    "createdAt": "2024-01-01T12:00:00Z",
    "updatedAt": "2024-01-01T12:00:00Z"
  }
}
```

## 3. 服务表单配置接口

### 接口信息
- **接口地址**: `GET /api/service/form_config/:id`
- **请求方式**: GET
- **功能**: 获取服务的动态表单配置

### 路径参数
- `id`: 服务ID

### 响应格式
```json
{
  "code": 0,
  "data": {
    "serviceId": 1,
    "formConfig": {
      "fields": [
        {
          "name": "patientName",
          "label": "就诊人姓名",
          "type": "text",
          "required": true,
          "placeholder": "请输入就诊人姓名"
        },
        {
          "name": "patientPhone",
          "label": "联系电话",
          "type": "phone",
          "required": true,
          "placeholder": "请输入联系电话"
        },
        {
          "name": "appointmentDate",
          "label": "预约日期",
          "type": "date",
          "required": true,
          "minDate": "2024-01-01",
          "maxDate": "2024-12-31"
        },
        {
          "name": "appointmentTime",
          "label": "预约时间",
          "type": "select",
          "required": true,
          "options": [
            {"label": "上午", "value": "morning"},
            {"label": "下午", "value": "afternoon"}
          ]
        },
        {
          "name": "specialRequirements",
          "label": "特殊要求",
          "type": "textarea",
          "required": false,
          "placeholder": "如有特殊要求请说明"
        }
      ],
      "submitText": "提交预约",
      "submitUrl": "/api/order/submit"
    }
  }
}
```

## 表单字段类型说明

| 类型 | 说明 | 示例 |
|------|------|------|
| text | 文本输入框 | 姓名、备注等 |
| phone | 手机号输入框 | 联系电话 |
| email | 邮箱输入框 | 邮箱地址 |
| number | 数字输入框 | 年龄、数量等 |
| date | 日期选择器 | 预约日期 |
| time | 时间选择器 | 预约时间 |
| select | 下拉选择框 | 预约时段、性别等 |
| radio | 单选按钮 | 性别、是否等 |
| checkbox | 多选框 | 检查项目等 |
| textarea | 多行文本 | 特殊要求、备注等 |
| file | 文件上传 | 检查报告、身份证等 |

## 使用示例

### 微信小程序端

```javascript
// 获取服务列表
wx.request({
  url: 'http://your-server.com/api/service/list',
  method: 'GET',
  data: {
    category: '体检套餐',
    page: 1,
    pageSize: 10
  },
  success: (res) => {
    if (res.data.code === 0) {
      console.log('服务列表:', res.data.data);
    }
  }
});

// 获取服务详情
wx.request({
  url: 'http://your-server.com/api/service/detail/1',
  method: 'GET',
  success: (res) => {
    if (res.data.code === 0) {
      console.log('服务详情:', res.data.data);
    }
  }
});

// 获取服务表单配置
wx.request({
  url: 'http://your-server.com/api/service/form_config/1',
  method: 'GET',
  success: (res) => {
    if (res.data.code === 0) {
      const formConfig = res.data.data.formConfig;
      console.log('表单配置:', formConfig);
      
      // 根据配置动态生成表单
      this.generateForm(formConfig);
    }
  }
});

// 动态生成表单
generateForm(formConfig) {
  const fields = formConfig.fields;
  let formHtml = '';
  
  fields.forEach(field => {
    switch (field.type) {
      case 'text':
        formHtml += `<input type="text" name="${field.name}" placeholder="${field.placeholder}" ${field.required ? 'required' : ''} />`;
        break;
      case 'phone':
        formHtml += `<input type="tel" name="${field.name}" placeholder="${field.placeholder}" ${field.required ? 'required' : ''} />`;
        break;
      case 'select':
        let options = '';
        field.options.forEach(option => {
          options += `<option value="${option.value}">${option.label}</option>`;
        });
        formHtml += `<select name="${field.name}" ${field.required ? 'required' : ''}>${options}</select>`;
        break;
      // 其他字段类型...
    }
  });
  
  // 渲染表单
  this.setData({
    formHtml: formHtml
  });
}
```

## 数据库表结构

### 服务项目表 (ServiceItems)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | INT | 主键，自增 |
| title | VARCHAR(200) | 服务标题 |
| description | VARCHAR(500) | 服务描述 |
| detail | TEXT | 服务详情 |
| price | DECIMAL(10,2) | 服务价格 |
| originalPrice | DECIMAL(10,2) | 原价 |
| category | VARCHAR(100) | 服务分类 |
| images | TEXT | 服务图片（JSON数组） |
| formConfig | TEXT | 表单配置（JSON） |
| status | INT | 状态：1-启用，0-禁用 |
| sort | INT | 排序 |
| viewCount | INT | 查看次数 |
| createdAt | DATETIME | 创建时间 |
| updatedAt | DATETIME | 更新时间 |

## 服务分类说明

| 分类名 | 说明 |
|--------|------|
| 体检套餐 | 各种体检套餐服务 |
| 单项检查 | 单项医疗检查 |
| 健康咨询 | 健康咨询服务 |
| 疫苗接种 | 疫苗接种服务 |
| 其他服务 | 其他医疗服务 |

## 注意事项

1. **分页处理**: 支持分页查询，默认每页10条，最大50条
2. **分类筛选**: 支持按服务分类筛选
3. **价格显示**: 显示当前价格和原价，支持优惠信息
4. **图片管理**: 服务图片支持多张，以JSON数组存储
5. **表单配置**: 动态表单配置支持多种字段类型
6. **状态控制**: 通过status字段控制服务是否可用
7. **排序规则**: 按sort字段升序排列
8. **查看统计**: 记录服务查看次数
9. **数据验证**: 表单字段支持必填、格式验证等
10. **响应式设计**: 表单配置支持不同设备适配 