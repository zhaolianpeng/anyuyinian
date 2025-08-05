# WebSocket连接错误分析

## 问题现象
根据测试结果：
- ✅ API调用测试成功 - 后端服务正常运行
- ❌ WebSocket连接测试失败 - WebSocket服务有问题

## 问题分析

### 1. 后端服务状态
- **HTTP API正常**: 说明后端服务已启动并正常运行
- **WebSocket服务异常**: 可能是WebSocket路由或处理器有问题

### 2. 可能的原因

#### 2.1 WebSocket路由问题
**检查项目**:
- `/ws` 路由是否正确注册
- WebSocketHandler是否正确实现
- 是否有中间件影响WebSocket升级

#### 2.2 云托管配置问题
**检查项目**:
- 云托管服务是否支持WebSocket
- 环境配置是否正确
- 服务是否已正确部署

#### 2.3 网络连接问题
**检查项目**:
- 小程序网络权限
- 云托管网络配置
- 防火墙设置

## 诊断工具

### 1. 详细诊断工具
**文件**: `debug_websocket_detailed.js`
**功能**:
- 环境检查
- 配置验证
- 网络测试
- WebSocket连接测试
- 详细错误分析

### 2. 测试页面
**页面**: `pages/test/websocket-test`
**功能**:
- 基础调试
- 详细诊断
- 连接测试
- 结果展示

## 诊断步骤

### 步骤1: 使用详细诊断
1. 在小程序中打开测试页面
2. 点击"详细诊断"按钮
3. 等待诊断完成
4. 查看诊断报告

### 步骤2: 检查后端日志
```bash
# 查看后端服务日志
docker logs <container_id>

# 或者查看云托管服务日志
# 在微信云托管控制台中查看
```

### 步骤3: 验证WebSocket端点
```bash
# 测试WebSocket端点
curl -i -N -H "Connection: Upgrade" \
     -H "Upgrade: websocket" \
     -H "Sec-WebSocket-Version: 13" \
     -H "Sec-WebSocket-Key: x3JJHMbDL1EzLkh9GBhXDw==" \
     http://localhost:80/ws
```

## 解决方案

### 方案1: 检查WebSocket路由注册
确保在 `main.go` 中正确注册了WebSocket路由：

```go
// WebSocket路由
http.HandleFunc("/ws", service.NewLogMiddleware(service.WebSocketHandler))
```

### 方案2: 检查WebSocket处理器
确保 `WebSocketHandler` 函数正确实现：

```go
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("WebSocket升级失败: %v", err)
        return
    }
    // ... 处理WebSocket连接
}
```

### 方案3: 检查中间件
确保中间件不影响WebSocket升级：

```go
// 可能需要为WebSocket路由跳过某些中间件
http.HandleFunc("/ws", service.WebSocketHandler) // 不使用日志中间件
```

### 方案4: 验证云托管配置
1. 确认云托管环境ID正确
2. 确认服务名正确
3. 检查服务是否支持WebSocket
4. 重新部署服务

### 方案5: 添加调试日志
在WebSocket处理器中添加更多日志：

```go
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("收到WebSocket连接请求: %s", r.URL.Path)
    
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("WebSocket升级失败: %v", err)
        return
    }
    
    log.Printf("WebSocket连接成功建立")
    // ... 处理连接
}
```

## 常见错误及解决方案

### 错误1: "WebSocket升级失败"
**原因**: 路由未正确注册或处理器有问题
**解决**: 检查路由注册和处理器实现

### 错误2: "连接超时"
**原因**: 网络问题或服务未响应
**解决**: 检查网络连接和服务状态

### 错误3: "服务不存在"
**原因**: 云托管配置错误
**解决**: 检查环境ID和服务名

### 错误4: "路由不存在"
**原因**: WebSocket路由未注册
**解决**: 确保在main.go中注册了 `/ws` 路由

## 验证步骤

### 1. 本地测试
```bash
# 构建并运行服务
./build.sh
./main

# 测试WebSocket端点
curl -i -N -H "Connection: Upgrade" \
     -H "Upgrade: websocket" \
     -H "Sec-WebSocket-Version: 13" \
     http://localhost:80/ws
```

### 2. 云托管测试
1. 部署到云托管
2. 使用小程序测试页面进行测试
3. 查看云托管服务日志

### 3. 小程序测试
1. 打开WebSocket测试页面
2. 运行详细诊断
3. 查看诊断报告

## 预防措施

1. **添加健康检查**: 定期检查WebSocket服务状态
2. **完善错误处理**: 提供详细的错误信息
3. **添加监控**: 记录连接成功率和错误率
4. **优化重连机制**: 网络异常时自动恢复连接
5. **定期测试**: 定期运行WebSocket连接测试

## 相关文件

- `debug_websocket_detailed.js` - 详细诊断工具
- `pages/test/websocket-test.*` - 测试页面
- `service/websocket_service.go` - WebSocket服务实现
- `main.go` - 路由注册
- `WEBSOCKET_TROUBLESHOOTING.md` - 问题诊断文档 