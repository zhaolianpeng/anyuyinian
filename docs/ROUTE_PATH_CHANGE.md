# 路由路径修改说明

## 修改内容
将管理员服务价格更新接口路径从：
```
/api/admin/service/update-price
```
修改为：
```
/api/admin/service/updateprice
```

## 修改原因
1. **避免路由冲突**: 原路径包含连字符，可能与某些路由匹配规则冲突
2. **简化路径**: 去除连字符，使路径更简洁
3. **提高兼容性**: 避免特殊字符可能导致的解析问题

## 修改文件

### 后端文件
- `main.go`: 修改路由注册路径
  ```go
  // 修改前
  http.HandleFunc("/api/admin/service/update-price", service.NewLogMiddleware(service.UpdateServicePriceHandler))
  
  // 修改后
  http.HandleFunc("/api/admin/service/updateprice", service.NewLogMiddleware(service.UpdateServicePriceHandler))
  ```

### 前端文件
- `pages/admin/services.js`: 修改API调用路径
  ```javascript
  // 修改前
  app.callContainer('/api/admin/service/update-price', 'POST', {
  
  // 修改后
  app.callContainer('/api/admin/service/updateprice', 'POST', {
  ```

### 测试文件
- `test_fix.sh`: 更新测试脚本中的API路径
- `test_api_calls.sh`: 更新测试脚本中的API路径

## 部署步骤
1. **构建项目**: `./deploy.sh` (已完成)
2. **手动部署**: 使用微信开发者工具上传并部署
3. **测试验证**: 运行 `./test_fix.sh` 验证功能

## 验证方法
```bash
# 测试新的API路径
curl -X POST "https://golang-lfwy-prod-5g94mx7a3d07e78c.service.tcloudbaseapp.com/api/admin/service/updateprice" \
  -H "Content-Type: application/json" \
  -H "X-WX-SERVICE: golang-lfwy" \
  -d '{"serviceId": 21, "newPrice": 280, "newOriginalPrice": 200, "reason": "测试"}'
```

## 注意事项
1. **前端缓存**: 确保小程序重新编译，清除缓存
2. **API调用**: 所有调用该接口的地方都需要更新路径
3. **向后兼容**: 旧路径将不再可用，需要确保所有调用都已更新

## 相关文件
- `main.go` - 后端路由配置
- `pages/admin/services.js` - 前端API调用
- `test_fix.sh` - 测试脚本
- `test_api_calls.sh` - API测试脚本
