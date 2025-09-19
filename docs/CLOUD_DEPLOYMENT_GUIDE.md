# 云托管部署指南

## 问题描述
管理员服务价格更新接口 `/api/admin/service/update-price` 返回404错误，提示 "Cannot find path"。

## 问题分析

### 1. 代码状态
- ✅ **路由配置**: 已在 `main.go` 中正确配置
- ✅ **处理器实现**: `UpdateServicePriceHandler` 已完整实现
- ✅ **本地编译**: 代码可以正常编译
- ❌ **云托管部署**: 可能未部署最新代码

### 2. 可能原因
1. **代码未部署**: 云托管可能运行的是旧版本代码
2. **部署失败**: 部署过程中可能出现了错误
3. **服务未重启**: 即使部署成功，服务可能未重启
4. **缓存问题**: 云托管可能存在路由缓存

## 解决步骤

### 步骤1: 确认代码状态
```bash
cd /Users/zhaolianpeng/code/Goproject/src/anyuyinian

# 检查路由配置
grep -n "update-price" main.go

# 检查处理器实现
grep -n "UpdateServicePriceHandler" service/admin_service.go

# 重新构建
./build.sh
```

### 步骤2: 部署到云托管

#### 方法1: 使用微信开发者工具
1. 打开微信开发者工具
2. 选择云托管项目
3. 右键点击 `golang-lfwy` 服务
4. 选择 "上传并部署"
5. 等待部署完成

#### 方法2: 使用命令行工具
```bash
# 确保已安装微信开发者工具命令行
npm install -g @cloudbase/cli

# 登录
tcb login

# 部署服务
tcb cloud deploy --env prod-5g94mx7a3d07e78c --service golang-lfwy
```

### 步骤3: 验证部署
```bash
# 测试服务是否可访问
curl -X GET "https://golang-lfwy-prod-5g94mx7a3d07e78c.service.tcloudbaseapp.com/api/service/categories"

# 测试管理员接口（需要认证）
curl -X POST "https://golang-lfwy-prod-5g94mx7a3d07e78c.service.tcloudbaseapp.com/api/admin/service/update-price" \
  -H "Content-Type: application/json" \
  -d '{"serviceId": 21, "newPrice": 258, "newOriginalPrice": 380, "reason": "测试"}'
```

### 步骤4: 检查服务状态
1. 登录腾讯云控制台
2. 进入云托管服务
3. 检查 `golang-lfwy` 服务状态
4. 查看服务日志
5. 确认服务版本是否为最新

## 部署检查清单

### 代码检查
- [ ] 路由配置正确 (`/api/admin/service/update-price`)
- [ ] 处理器实现完整 (`UpdateServicePriceHandler`)
- [ ] 本地编译成功
- [ ] 代码已提交到版本控制

### 部署检查
- [ ] 云托管服务状态正常
- [ ] 服务版本为最新
- [ ] 服务日志无错误
- [ ] 网络连接正常

### 功能检查
- [ ] 基础API可访问 (`/api/service/categories`)
- [ ] 管理员API可访问 (`/api/admin/service/update-price`)
- [ ] 前端可以正常调用接口

## 常见问题

### 1. 部署后仍然404
**可能原因**: 服务未完全重启
**解决方案**: 等待几分钟后重试，或手动重启服务

### 2. 权限问题
**可能原因**: 缺少 `adminUserId` 参数
**解决方案**: 确保请求包含管理员身份参数

### 3. 网络问题
**可能原因**: 云托管网络配置问题
**解决方案**: 检查云托管网络设置和防火墙规则

## 调试命令

### 检查服务状态
```bash
# 检查服务是否运行
curl -I "https://golang-lfwy-prod-5g94mx7a3d07e78c.service.tcloudbaseapp.com/"

# 检查健康状态
curl -X GET "https://golang-lfwy-prod-5g94mx7a3d07e78c.service.tcloudbaseapp.com/health"
```

### 查看服务日志
1. 登录腾讯云控制台
2. 进入云托管服务
3. 选择 `golang-lfwy` 服务
4. 查看实时日志
5. 查找错误信息

## 联系支持

如果问题仍然存在，请：
1. 提供云托管服务日志
2. 提供具体的错误信息
3. 说明部署步骤和结果
4. 联系腾讯云技术支持

## 相关文件

- `main.go` - 路由配置
- `service/admin_service.go` - 管理员服务处理器
- `deploy.sh` - 部署脚本
- `build.sh` - 构建脚本
