# 腾讯云COS配置说明

## 概述

本项目已集成腾讯云对象存储(COS)功能，用于文件上传和存储。

## 配置信息

- **存储桶名称**: `7072-prod-5g94mx7a3d07e78c-1353115175`
- **存储桶地域**: `ap-shanghai` (上海)
- **访问域名**: `https://7072-prod-5g94mx7a3d07e78c-1353115175.cos.ap-shanghai.myqcloud.com`
- **权限策略**: 所有用户可读，仅创建者可读写

## 环境变量配置

### 开发环境

在开发环境中，需要设置以下环境变量：

```bash
export COS_SECRET_ID="your_actual_secret_id"
export COS_SECRET_KEY="your_actual_secret_key"
```

### 生产环境

在生产环境中，建议通过容器编排工具（如Docker、Kubernetes）设置环境变量：

```yaml
# docker-compose.yml 示例
environment:
  - COS_SECRET_ID=your_actual_secret_id
  - COS_SECRET_KEY=your_actual_secret_key
```

## 获取腾讯云密钥

1. 登录腾讯云控制台
2. 进入 [访问管理](https://console.cloud.tencent.com/cam)
3. 选择 [API密钥管理](https://console.cloud.tencent.com/cam/capi)
4. 创建新的密钥或使用现有密钥
5. 复制 SecretId 和 SecretKey

## 权限配置

### 腾讯云账号权限
确保您的腾讯云账号具有以下权限：

- `cos:PutObject` - 上传对象
- `cos:GetObject` - 下载对象
- `cos:DeleteObject` - 删除对象
- `cos:ListBucket` - 列出存储桶内容
- `cos:PutObjectACL` - 设置对象ACL
- `cos:GetObjectACL` - 获取对象ACL

### 对象权限策略
本项目采用以下权限策略：

- **上传时**: 自动设置为 `public-read`（所有用户可读，仅创建者可读写）
- **可修改为**: `private`（仅创建者可读写）
- **权限管理**: 支持动态修改文件权限

## 功能特性

### 1. 文件上传
- 支持多种文件格式（图片、文档、PDF等）
- 自动生成唯一文件名
- 文件大小限制：10MB
- 自动分类存储

### 2. 文件管理
- 文件列表查询
- 按用户ID查询
- 按分类查询
- 文件信息记录

### 3. 安全特性
- 文件类型验证
- 文件大小限制
- 用户权限验证
- 敏感信息过滤

## API接口

### 上传文件
```
POST /api/upload
Content-Type: multipart/form-data

参数:
- file: 文件对象
- userId: 用户ID
- category: 文件分类（可选）
- description: 文件描述（可选）
```

### 获取文件列表
```
GET /api/files?userId=123&limit=20
GET /api/files?category=image&limit=20
```

### 删除文件
```
DELETE /api/file/delete?fileId=123
```

### 更新文件权限
```
PUT /api/file/permission
Content-Type: application/json

{
  "fileId": 123,
  "acl": "public-read"  // "public-read" 或 "private"
}
```

### 获取文件权限
```
GET /api/file/permission/get?fileId=123
```

## 日志记录

所有文件操作都会记录详细日志：

```
[API] 开始处理请求: POST /api/upload
[STEP] 开始处理文件上传请求: {"method":"POST","path":"/api/upload"}
[STEP] 开始验证文件: {"fileName":"test.jpg","fileSize":1024000,"contentType":"image/jpeg"}
[STEP] 使用腾讯云COS上传文件: {"fileName":"1234567890_abcdef12.jpg"}
[STEP] COS上传成功: {"fileUrl":"https://...","acl":"public-read"}
[STEP] 设置对象为公共读取权限: {"objectKey":"1234567890_abcdef12.jpg"}
[STEP] 对象公共读取权限设置成功: {"objectKey":"1234567890_abcdef12.jpg"}
[DB] 创建 表: files, 参数: {"fileName":"1234567890_abcdef12.jpg","fileUrl":"https://..."}
[STEP] 文件保存完成: {"fileId":1,"fileUrl":"https://..."}
```

## 故障排除

### 1. 上传失败
- 检查网络连接
- 验证密钥是否正确
- 确认存储桶权限
- 查看日志错误信息

### 2. 权限错误
- 确认SecretId和SecretKey正确
- 检查存储桶权限设置
- 验证API密钥是否启用

### 3. 网络问题
- 检查防火墙设置
- 确认网络连接正常
- 验证域名解析

## 性能优化

1. **并发上传**: 支持多文件并发上传
2. **断点续传**: 大文件支持断点续传
3. **CDN加速**: 可配置CDN加速访问
4. **压缩传输**: 自动压缩传输数据

## 监控告警

建议配置以下监控指标：

- 上传成功率
- 上传响应时间
- 存储空间使用量
- 错误率统计

## 成本控制

1. **存储费用**: 按实际存储量计费
2. **流量费用**: 按下载流量计费
3. **请求费用**: 按API请求次数计费
4. **CDN费用**: 如启用CDN加速

## 安全建议

1. **密钥管理**: 定期轮换API密钥
2. **权限最小化**: 仅授予必要权限
3. **访问控制**: 配置IP白名单
4. **日志审计**: 定期检查访问日志 