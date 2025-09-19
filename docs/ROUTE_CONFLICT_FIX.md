# 路由冲突修复文档

## 问题描述
管理员服务价格更新接口 `/api/admin/service/update-price` 返回404错误，但其他接口正常工作。

## 问题根因
**路由冲突**：Go的`http.HandleFunc`按照最长匹配原则，`/api/admin/services`会匹配`/api/admin/service/update-price`，因为`/api/admin/service`是`/api/admin/services`的前缀。

## 问题分析

### 原始路由注册顺序
```go
// 管理员服务管理相关接口
http.HandleFunc("/api/admin/services", service.NewLogMiddleware(service.GetAdminServicesHandler))
http.HandleFunc("/api/admin/service/update-price", service.NewLogMiddleware(service.UpdateServicePriceHandler))
```

### 问题说明
1. **路由匹配规则**: Go的`http.HandleFunc`会匹配最长的路径前缀
2. **冲突路径**: `/api/admin/services` 会匹配 `/api/admin/service/update-price`
3. **结果**: 所有对 `/api/admin/service/update-price` 的请求都被 `/api/admin/services` 处理
4. **错误**: `GetAdminServicesHandler` 只处理GET请求，POST请求返回404

## 解决方案

### 修复方法：调整路由注册顺序
```go
// 管理员服务管理相关接口
http.HandleFunc("/api/admin/service/update-price", service.NewLogMiddleware(service.UpdateServicePriceHandler))
http.HandleFunc("/api/admin/services", service.NewLogMiddleware(service.GetAdminServicesHandler))
```

### 修复原理
1. **具体路径优先**: 将更具体的路径 `/api/admin/service/update-price` 放在前面
2. **通用路径在后**: 将更通用的路径 `/api/admin/services` 放在后面
3. **避免前缀冲突**: 确保具体路径不会被通用路径覆盖

## 验证修复

### 1. 检查路由注册顺序
```bash
grep -n "admin.*service" main.go
```

### 2. 测试接口访问
```bash
# 测试管理员API
curl -X POST "https://golang-lfwy-prod-5g94mx7a3d07e78c.service.tcloudbaseapp.com/api/admin/service/update-price" \
  -H "Content-Type: application/json" \
  -d '{"serviceId": 21, "newPrice": 258, "newOriginalPrice": 380, "reason": "测试"}'
```

### 3. 检查后端日志
确认请求到达了正确的处理器：
```
[INFO] 开始处理修改服务价格请求
```

## 预防措施

### 1. 路由命名规范
- 使用具体的路径名称
- 避免路径前缀冲突
- 按具体程度排序（具体路径在前）

### 2. 路由测试
- 为每个路由编写测试用例
- 验证路由匹配正确性
- 检查HTTP方法支持

### 3. 代码审查
- 检查路由注册顺序
- 验证路径命名规范
- 确认处理器实现

## 相关文件

- `main.go` - 路由配置
- `service/admin_service.go` - 处理器实现
- `ROUTE_CONFLICT_FIX.md` - 本文档

## 修复状态

- ✅ **问题识别**: 路由冲突导致404错误
- ✅ **根因分析**: `/api/admin/services` 覆盖了 `/api/admin/service/update-price`
- ✅ **修复实施**: 调整路由注册顺序
- ✅ **代码构建**: 重新构建成功
- ⏳ **部署验证**: 需要重新部署到云托管

## 下一步操作

1. **重新部署服务**到云托管
2. **验证接口访问**是否正常
3. **测试管理员功能**是否可用
4. **检查后端日志**确认请求到达正确处理器

---

**⚠️ 重要提醒**: 路由冲突是常见的Go Web开发问题，需要仔细规划路由注册顺序！
