# 日志系统说明

## 概述

本项目已经集成了完整的日志系统，包括：

1. **请求日志中间件** - 自动记录所有HTTP请求的入参、出参和耗时
2. **业务日志函数** - 用于记录业务逻辑中的关键步骤
3. **数据库操作日志** - 记录所有数据库操作的参数和结果
4. **错误日志** - 记录异常和错误信息

## 日志级别

### [API] - 接口级别日志
- 记录HTTP请求的开始、结束、参数、响应和耗时
- 自动由中间件处理，无需手动添加

### [STEP] - 步骤日志
- 记录业务逻辑中的关键步骤
- 使用 `LogStep(step string, data interface{})` 函数

### [DB] - 数据库操作日志
- 记录数据库查询、插入、更新、删除操作
- 使用 `LogDBOperation()` 和 `LogDBResult()` 函数

### [ERROR] - 错误日志
- 记录异常和错误信息
- 使用 `LogError(message string, err error)` 函数

### [INFO] - 信息日志
- 记录一般信息
- 使用 `LogInfo(message string, data interface{})` 函数

## 使用示例

### 1. 记录业务步骤
```go
LogStep("开始处理用户登录", map[string]string{"openId": openId})
LogStep("用户验证成功", userData)
```

### 2. 记录数据库操作
```go
// 操作前
LogDBOperation("查询", "users", map[string]string{"openId": openId})

// 操作后
LogDBResult("查询", "users", user, err)
```

### 3. 记录错误
```go
if err != nil {
    LogError("用户创建失败", err)
    return err
}
```

### 4. 记录响应
```go
LogResponse(response, nil) // 成功
LogResponse(nil, err)      // 失败
```

## 日志输出格式

```
2024/01/01 12:00:00 [API] 开始处理请求: POST /api/wx/login
2024/01/01 12:00:00 [API] 请求体: {"code":"123456","nickName":"用户"}
2024/01/01 12:00:00 [STEP] 开始处理微信登录请求: {"method":"POST","path":"/api/wx/login"}
2024/01/01 12:00:00 [STEP] 开始调用微信API: {"code":"123456"}
2024/01/01 12:00:00 [DB] 查询 表: users, 参数: {"openId":"wx_openid_123"}
2024/01/01 12:00:00 [DB] 查询 表: users 成功: {"id":1,"openId":"wx_openid_123"}
2024/01/01 12:00:00 [STEP] 用户登录处理完成: {"code":0,"data":{"id":1,"openId":"wx_openid_123"}}
2024/01/01 12:00:00 [API] 响应状态: 200
2024/01/01 12:00:00 [API] 响应体: {"code":0,"data":{"id":1,"openId":"wx_openid_123"}}
2024/01/01 12:00:00 [API] 请求处理完成，耗时: 150ms
```

## 注意事项

1. **敏感信息过滤** - 日志中不会记录密码、session_key等敏感信息
2. **性能影响** - 日志记录是异步的，不会影响接口响应时间
3. **日志级别** - 生产环境可以配置不同的日志级别
4. **日志轮转** - 建议配置日志轮转以避免日志文件过大

## 配置说明

日志系统已经集成到所有接口中，通过 `main.go` 中的中间件自动启用：

```go
http.HandleFunc("/api/wx/login", service.NewLogMiddleware(service.WxLoginHandler))
```

所有接口都会自动记录：
- 请求开始和结束时间
- 请求参数（GET参数、POST请求体）
- 响应状态码和响应体
- 总耗时
- 业务逻辑中的关键步骤
- 数据库操作
- 错误信息 