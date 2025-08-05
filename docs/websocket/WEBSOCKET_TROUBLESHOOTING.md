# WebSocket连接问题诊断和解决方案

## 问题描述
小程序测试页面提示WebSocket连接失败，需要检查后端服务。

## 可能的原因分析

### 1. 后端服务未启动
**症状**: 无法连接到WebSocket端点
**检查方法**:
```bash
# 检查服务是否运行
curl -s http://localhost:80/ > /dev/null && echo "服务运行中" || echo "服务未运行"

# 检查WebSocket端点
curl -i -N -H "Connection: Upgrade" -H "Upgrade: websocket" -H "Sec-WebSocket-Version: 13" http://localhost:80/ws
```

**解决方案**:
- 确保后端服务已启动
- 检查服务是否监听在正确的端口(80)
- 查看服务日志确认WebSocket路由已注册

### 2. 云托管环境配置问题
**症状**: 连接超时或服务不存在错误
**检查项目**:
- 环境ID是否正确: `prod-5g94mx7a3d07e78c`
- 服务名是否正确: `golang-lfwy`
- 云托管服务是否已部署并运行

**解决方案**:
- 确认云托管环境配置正确
- 重新部署后端服务
- 检查云托管服务状态

### 3. 网络连接问题
**症状**: 网络请求失败或超时
**检查方法**:
```javascript
// 在小程序中测试网络连接
wx.request({
  url: 'https://www.baidu.com',
  success: (res) => console.log('网络正常'),
  fail: (err) => console.log('网络异常:', err)
})
```

**解决方案**:
- 检查网络连接
- 确认小程序有网络权限
- 检查防火墙设置

### 4. WebSocket路由配置问题
**症状**: 404错误或路由不存在
**检查项目**:
- 后端是否正确注册了 `/ws` 路由
- WebSocket处理器是否正确实现

**解决方案**:
- 确认main.go中已注册WebSocket路由
- 检查WebSocketHandler函数实现

## 诊断步骤

### 步骤1: 检查后端服务状态
```bash
# 运行测试脚本
./test_websocket.sh
```

### 步骤2: 检查小程序配置
```javascript
// 在小程序中使用调试工具
const WebSocketDebugger = require('./debug_websocket.js')
const debugger = new WebSocketDebugger()
debugger.startDebug()
```

### 步骤3: 检查云托管服务
- 登录微信云托管控制台
- 检查服务 `golang-lfwy` 是否运行
- 查看服务日志

### 步骤4: 验证配置
确认以下配置正确:
```javascript
// miniprogram/config.js
const WEBSOCKET_CONFIG = {
  env: 'prod-5g94mx7a3d07e78c',  // 云托管环境ID
  service: 'golang-lfwy',         // 服务名
  path: '/ws',                    // WebSocket路径
  reconnect: {
    enabled: true,
    maxAttempts: 5,
    interval: 3000
  }
}
```

## 解决方案

### 方案1: 重启后端服务
```bash
# 重新构建并启动服务
./build.sh
# 部署到云托管
```

### 方案2: 检查云托管配置
1. 确认云托管环境ID正确
2. 确认服务名正确
3. 检查服务是否已部署

### 方案3: 更新小程序配置
1. 确认config.js中的WebSocket配置正确
2. 清除小程序缓存
3. 重新编译小程序

### 方案4: 添加错误处理
在小程序中添加更详细的错误处理:
```javascript
// 在app.js中添加
handleWebSocketError(error) {
  console.error('WebSocket错误详情:', error)
  
  // 根据错误类型提供具体建议
  if (error.message.includes('timeout')) {
    wx.showToast({
      title: '连接超时，请检查网络',
      icon: 'none'
    })
  } else if (error.message.includes('service')) {
    wx.showToast({
      title: '服务暂时不可用',
      icon: 'none'
    })
  } else {
    wx.showToast({
      title: '连接失败，请稍后重试',
      icon: 'none'
    })
  }
}
```

## 调试工具

### 1. 后端测试脚本
```bash
# test_websocket.sh
./test_websocket.sh
```

### 2. 小程序调试工具
```javascript
// debug_websocket.js
const debugger = new WebSocketDebugger()
debugger.startDebug()
```

### 3. 日志检查
```bash
# 检查后端日志
docker logs <container_id>

# 检查小程序控制台日志
# 在微信开发者工具中查看
```

## 预防措施

1. **定期检查服务状态**
2. **添加健康检查端点**
3. **实现自动重连机制**
4. **添加详细的错误日志**
5. **定期更新依赖包**

## 相关文件

- `service/websocket_service.go` - WebSocket服务实现
- `main.go` - 路由注册
- `miniprogram/config.js` - WebSocket配置
- `miniprogram/app.js` - WebSocket连接管理
- `test_websocket.sh` - 后端测试脚本
- `debug_websocket.js` - 小程序调试工具 