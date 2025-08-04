# WebSocket连接问题诊断总结

## 问题描述
小程序测试页面提示WebSocket连接失败，需要检查后端服务。

## 诊断工具和文件

### 1. 后端测试工具
- **文件**: `test_websocket.sh`
- **功能**: 测试后端服务状态和WebSocket端点
- **使用方法**: `./test_websocket.sh`

### 2. 小程序调试工具
- **文件**: `debug_websocket.js`
- **功能**: 诊断小程序WebSocket连接问题
- **使用方法**: 在小程序中引入并调用

### 3. WebSocket测试页面
- **页面路径**: `miniprogram/pages/test/websocket-test`
- **功能**: 提供完整的WebSocket连接测试界面
- **包含功能**:
  - 连接状态显示
  - 调试信息输出
  - 测试结果记录
  - 配置信息展示

## 可能的问题原因

### 1. 后端服务问题
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

### 2. 云托管配置问题
**症状**: 连接超时或服务不存在错误
**检查项目**:
- 环境ID: `prod-5g94mx7a3d07e78c`
- 服务名: `golang-lfwy`
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

### 4. WebSocket路由配置问题
**症状**: 404错误或路由不存在
**检查项目**:
- 后端是否正确注册了 `/ws` 路由
- WebSocket处理器是否正确实现

## 诊断步骤

### 步骤1: 使用测试页面
1. 在小程序中打开 `pages/test/websocket-test` 页面
2. 点击"开始调试"按钮
3. 查看调试信息和错误详情

### 步骤2: 检查后端服务
```bash
# 运行后端测试脚本
./test_websocket.sh
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

## 新增文件列表

### 后端文件
1. `test_websocket.sh` - 后端WebSocket测试脚本
2. `WEBSOCKET_TROUBLESHOOTING.md` - 问题诊断和解决方案文档

### 小程序文件
1. `debug_websocket.js` - WebSocket调试工具
2. `pages/test/websocket-test.js` - WebSocket测试页面逻辑
3. `pages/test/websocket-test.wxml` - WebSocket测试页面模板
4. `pages/test/websocket-test.wxss` - WebSocket测试页面样式
5. `pages/test/websocket-test.json` - WebSocket测试页面配置

## 使用建议

1. **优先使用测试页面**: 在小程序中打开WebSocket测试页面进行诊断
2. **查看控制台日志**: 在微信开发者工具中查看详细错误信息
3. **检查云托管状态**: 确认云托管服务正常运行
4. **验证网络连接**: 确保小程序有网络访问权限
5. **查看服务日志**: 检查后端服务的运行日志

## 后续优化

1. **添加健康检查**: 定期检查WebSocket连接状态
2. **实现自动重连**: 网络异常时自动恢复连接
3. **优化错误提示**: 提供更友好的错误信息
4. **添加监控**: 记录连接成功率和错误率
5. **性能优化**: 减少不必要的连接尝试 