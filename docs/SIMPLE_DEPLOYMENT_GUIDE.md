# 简化版本部署指南

## 问题描述
原版本服务返回418错误，可能由于数据库连接失败导致服务启动失败。

## 解决方案
使用不依赖数据库的简化版本进行快速验证和部署。

## 部署步骤

### 步骤1: 备份原版本
```bash
# 备份原Dockerfile
mv Dockerfile Dockerfile.original
mv main main.original
```

### 步骤2: 使用简化版本
```bash
# 使用简化版Dockerfile
mv Dockerfile.simple Dockerfile
```

### 步骤3: 重新部署
1. 打开微信开发者工具
2. 选择云托管项目
3. 右键点击 `golang-lfwy` 服务
4. 选择 "上传并部署"
5. 等待部署完成

### 步骤4: 验证部署
```bash
# 检查服务状态
curl -k "https://golang-lfwy-prod-5g94mx7a3d07e78c.service.tcloudbaseapp.com/health"

# 检查管理员API
curl -k -X POST "https://golang-lfwy-prod-5g94mx7a3d07e78c.service.tcloudbaseapp.com/api/admin/service/update-price" \
  -H "Content-Type: application/json" \
  -d '{"serviceId": 21, "newPrice": 280, "newOriginalPrice": 200, "reason": "测试"}'
```

## 简化版本功能

### 已实现的接口
- `/health` - 健康检查
- `/` - 根路径
- `/api/test` - 测试API
- `/api/admin/service/update-price` - 管理员API（测试版本）
- `/api/service/categories` - 分类API（测试版本）

### 特点
- 不依赖数据库
- 快速启动
- 包含详细日志
- 返回测试数据

## 验证成功标准

### 1. 服务启动成功
- 健康检查返回200状态码
- 根路径返回HTML页面

### 2. 管理员API正常
- POST请求返回200状态码
- 返回JSON响应
- 后端有请求日志

### 3. 前端功能正常
- 管理员页面可以正常加载
- 价格更新功能可以正常使用
- 无404错误

## 后续步骤

### 如果简化版本成功
1. 确认问题出在数据库连接
2. 修复数据库配置
3. 恢复原版本部署

### 如果简化版本失败
1. 检查云托管配置
2. 检查网络连接
3. 联系腾讯云支持

## 相关文件

- `simple_main.go` - 简化版主程序
- `Dockerfile.simple` - 简化版Dockerfile
- `main.go` - 原版主程序
- `Dockerfile.original` - 原版Dockerfile

## 注意事项

- 简化版本只用于测试，不包含真实业务逻辑
- 部署前请备份原版本文件
- 测试完成后请恢复原版本
