# 文件上传接口文档

## 接口信息

- **接口地址**: `POST /api/upload`
- **请求方式**: POST
- **功能**: 上传图片或资料文件

## 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| file | file | 是 | 要上传的文件 |
| userId | int | 否 | 用户ID |
| category | string | 否 | 文件分类 |

### 支持的文件类型

- 图片: jpg, jpeg, png, gif, webp
- 文档: pdf, doc, docx, xls, xlsx
- 其他: txt, zip, rar

### 文件大小限制

- 单个文件最大: 10MB

## 响应格式

### 成功响应

```json
{
  "code": 0,
  "data": {
    "fileId": 1,
    "fileName": "example.jpg",
    "fileUrl": "https://example.com/uploads/example.jpg",
    "fileSize": 1024000,
    "fileType": "image",
    "category": "home_environment",
    "message": "文件上传成功"
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
// 选择文件
wx.chooseImage({
  count: 1,
  sizeType: ['original', 'compressed'],
  sourceType: ['album', 'camera'],
  success: (res) => {
    const tempFilePath = res.tempFilePaths[0];
    
    // 上传文件
    wx.uploadFile({
      url: 'http://your-server.com/api/upload',
      filePath: tempFilePath,
      name: 'file',
      formData: {
        userId: 1,
        category: 'home_environment'
      },
      success: (uploadRes) => {
        const data = JSON.parse(uploadRes.data);
        if (data.code === 0) {
          console.log('上传成功:', data.data);
        } else {
          console.error('上传失败:', data.errorMsg);
        }
      }
    });
  }
});
```

### 获取文件列表

```javascript
// 获取文件列表
wx.request({
  url: 'http://your-server.com/api/files',
  method: 'GET',
  data: {
    userId: 1,
    category: 'home_environment',
    page: 1,
    pageSize: 10
  },
  success: (res) => {
    if (res.data.code === 0) {
      console.log('文件列表:', res.data.data);
    }
  }
});
```

## 数据库表结构

### 文件表 (Files)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | INT | 主键，自增 |
| fileName | VARCHAR(200) | 文件名 |
| originalName | VARCHAR(200) | 原始文件名 |
| filePath | VARCHAR(500) | 文件路径 |
| fileUrl | VARCHAR(500) | 文件访问URL |
| fileSize | BIGINT | 文件大小（字节） |
| fileType | VARCHAR(50) | 文件类型 |
| fileExt | VARCHAR(20) | 文件扩展名 |
| category | VARCHAR(100) | 文件分类 |
| userId | INT | 上传用户ID |
| status | INT | 状态：1-正常，0-已删除 |
| createdAt | DATETIME | 创建时间 |
| updatedAt | DATETIME | 更新时间 |

## 文件分类说明

| 分类名 | 说明 |
|--------|------|
| home_environment | 居家环境图 |
| medical_report | 医疗报告 |
| id_card | 身份证 |
| other | 其他文件 |

## 文件存储说明

1. **存储路径**: 文件保存在 `./uploads` 目录下
2. **文件命名**: 使用时间戳+MD5的方式生成唯一文件名
3. **目录结构**: 按年月日创建子目录
4. **访问URL**: 通过 `/uploads/` 路径访问文件

## 注意事项

1. **文件验证**: 上传前会验证文件类型和大小
2. **安全处理**: 文件名会进行安全处理，防止路径遍历攻击
3. **软删除**: 删除文件时采用软删除，不会物理删除文件
4. **文件类型**: 自动识别文件类型和扩展名
5. **用户关联**: 文件与用户关联，支持按用户查询
6. **分类管理**: 支持按分类管理文件
7. **大小限制**: 单个文件最大10MB，总上传大小无限制
8. **格式支持**: 支持常见图片、文档格式 