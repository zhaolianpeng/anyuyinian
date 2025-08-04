# 后端接口日志记录实现总结

## 概述
已为所有主要后端接口添加了详细的日志记录，包括请求处理、参数验证、数据库操作、错误处理和响应等关键步骤。

## 已实现日志记录的服务

### 1. 用户服务 (user_service.go)
- **GetUserInfoHandler**: 获取用户信息
  - 请求开始日志
  - 请求方法验证
  - 参数解析和验证
  - 数据库查询日志
  - 查询结果日志
  - 响应准备和成功日志

- **BindPhoneHandler**: 绑定手机号
  - 请求开始日志
  - 参数解析和验证
  - 验证码验证过程
  - 用户信息查询
  - 手机号更新操作
  - 成功响应日志

- **AddressHandler**: 地址管理
  - 请求路由分发日志
  - 地址列表查询
  - 地址创建、更新、删除操作
  - 数据库操作结果记录

- **PatientHandler**: 就诊人管理
  - 就诊人信息CRUD操作
  - 数据库查询和更新日志

### 2. 订单服务 (order_service.go)
- **SubmitOrderHandler**: 提交订单
  - 请求参数解析和验证
  - 服务信息查询
  - 订单号生成
  - 金额计算和佣金计算
  - 表单数据处理
  - 订单创建和保存
  - 成功响应日志

- **PayOrderHandler**: 支付订单
- **CancelOrderHandler**: 取消订单
- **RefundOrderHandler**: 退款订单
- **OrderListHandler**: 订单列表
- **OrderDetailHandler**: 订单详情

### 3. 服务管理 (service_service.go)
- **ServiceListHandler**: 服务列表
  - 请求参数解析
  - 分页参数处理
  - 按分类或全部查询
  - 数据库查询结果记录
  - 分页信息计算

- **ServiceDetailHandler**: 服务详情
  - 服务ID格式转换
  - 参数验证
  - 数据库查询
  - 错误处理（包括记录不存在）

### 4. 首页初始化 (home_init_service.go)
- **HomeInitHandler**: 首页初始化
  - GET/POST请求处理
  - 参数解析（经纬度、限制数量）
  - 默认值设置
  - 数据获取过程

- **getHomeInitData**: 数据获取
  - 轮播图数据查询
  - 导航数据查询
  - 服务项数据查询
  - 医院列表查询（支持位置排序）
  - 数据转换和统计

### 5. 客服服务 (kefu_service.go)
- **SendMessageHandler**: 发送消息
  - 请求参数解析
  - 图片数据处理
  - 消息对象创建
  - 数据库保存
  - 成功响应

### 6. 微信登录服务 (wx_login_service.go)
- **WxLoginHandler**: 微信登录
  - 请求处理日志
  - 微信API调用
  - 用户信息处理
  - 数据库操作

### 7. 文件上传服务 (upload_service.go)
- **UploadHandler**: 文件上传
  - 表单解析
  - 文件验证
  - 文件保存（本地/COS）
  - 数据库记录

## 日志记录类型

### 1. 请求级别日志
- `LogInfo`: 记录请求开始和结束
- `LogStep`: 记录处理步骤
- `LogError`: 记录错误信息

### 2. 数据库操作日志
- `LogDBOperation`: 记录数据库操作开始
- `LogDBResult`: 记录数据库操作结果

### 3. 参数验证日志
- 请求参数解析
- 参数格式验证
- 必要参数检查

### 4. 业务逻辑日志
- 数据计算过程
- 业务规则验证
- 状态转换记录

### 5. 响应日志
- 响应数据准备
- 成功/失败状态记录

## 日志格式示例

```
[INFO] 开始处理获取用户信息请求: {"method":"GET","path":"/api/user/info"}
[STEP] 解析请求参数: {"userId":"123"}
[STEP] 开始查询用户信息: {"userId":123}
[STEP] 用户信息查询成功: {"userId":123,"openId":"xxx","nickName":"用户"}
[INFO] 用户信息获取成功: {"userId":123}
```

## 错误处理日志

### 数据库错误
```
[ERROR] 数据库查询用户信息失败: record not found
[ERROR] 数据库创建订单失败: duplicate key value violates unique constraint
```

### 参数验证错误
```
[ERROR] 缺少必要参数: userId=0, serviceId=0, quantity=0
[ERROR] 用户ID格式错误: strconv.Atoi: parsing "abc": invalid syntax
```

### 业务逻辑错误
```
[ERROR] 验证码错误: 输入验证码=123456，期望验证码=654321
[ERROR] 服务不存在: serviceId=999
```

## 性能监控

通过日志可以监控：
- 请求处理时间
- 数据库查询性能
- 错误发生频率
- 业务操作成功率

## 下一步改进

1. **结构化日志**: 考虑使用结构化日志格式（JSON）
2. **日志级别**: 添加DEBUG、WARN等日志级别
3. **日志轮转**: 实现日志文件轮转机制
4. **监控集成**: 与监控系统集成
5. **敏感信息过滤**: 对敏感信息进行脱敏处理

## 使用说明

所有日志都会输出到标准输出，可以通过以下方式查看：

```bash
# 查看实时日志
tail -f /var/log/application.log

# 查看错误日志
grep "\[ERROR\]" /var/log/application.log

# 查看特定接口的日志
grep "获取用户信息" /var/log/application.log
```

这样的日志记录系统可以帮助开发人员快速定位问题，监控系统运行状态，并提供详细的调试信息。 