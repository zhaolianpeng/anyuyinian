# 安语颐年后端开发文档

## 项目概述

安语颐年是一个基于微信小程序的医疗服务平台，后端采用Go语言开发，提供预约服务、订单管理、用户管理、支付处理、咨询系统等核心功能。

## 技术栈

- **语言**: Go 1.16+
- **数据库**: MySQL 5.7+
- **ORM**: GORM v1.21.16
- **架构**: MVC分层架构
- **API风格**: RESTful API
- **部署**: 微信云托管
- **文件存储**: 腾讯云COS
- **支付**: 微信支付

## 项目结构

```
anyuyinian/
├── main.go                    # 主程序入口
├── go.mod                     # Go模块依赖
├── go.sum                     # 依赖版本锁定
├── config/                    # 配置管理
│   ├── wx_config.go          # 微信配置
│   ├── payment_config.go     # 支付配置
│   ├── cos_config.go         # 腾讯云COS配置
│   └── cos_secrets.go        # COS密钥配置
├── db/                        # 数据访问层
│   ├── init.go               # 数据库初始化
│   ├── dao/                  # 数据访问对象
│   │   ├── user_dao.go       # 用户DAO
│   │   ├── order_dao.go      # 订单DAO
│   │   ├── service_dao.go    # 服务DAO
│   │   └── ...               # 其他DAO
│   ├── model/                # 数据模型
│   │   ├── user.go           # 用户模型
│   │   ├── order.go          # 订单模型
│   │   ├── service.go        # 服务模型
│   │   └── ...               # 其他模型
│   └── migration/            # 数据库迁移脚本
├── service/                   # 业务逻辑层
│   ├── user_service.go       # 用户服务
│   ├── order_service.go      # 订单服务
│   ├── wechat_pay_service.go # 微信支付服务
│   ├── log_middleware.go     # 日志中间件
│   └── ...                   # 其他服务
├── handler/                   # 处理器层
│   └── consultation_handler.go
├── middleware/                # 中间件
├── utils/                     # 工具类
│   ├── id_generator.go       # ID生成器
│   └── promoter_code.go      # 推广码工具
├── static/                    # 静态文件
├── templates/                 # 模板文件
├── tests/                     # 测试文件
├── scripts/                   # 脚本文件
└── docs/                      # 文档
```

## 核心功能模块

### 1. 用户管理模块

#### 数据模型
```go
type UserModel struct {
    Id          int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
    UserId      string    `gorm:"column:userId;uniqueIndex;type:varchar(24);not null" json:"userId"`
    OpenId      string    `gorm:"column:openId;uniqueIndex;not null" json:"openId"`
    UnionId     string    `gorm:"column:unionId" json:"unionId"`
    NickName    string    `gorm:"column:nickName" json:"nickName"`
    AvatarUrl   string    `gorm:"column:avatarUrl" json:"avatarUrl"`
    Gender      int       `gorm:"column:gender" json:"gender"`
    Phone       string    `gorm:"column:phone" json:"phone"`
    // ... 其他字段
}
```

#### 主要接口
- `GET /api/user/info` - 获取用户信息
- `POST /api/user/bind_phone` - 绑定手机号
- `POST /api/user/update_info` - 更新用户信息
- `POST /api/user/decrypt_phone` - 解密微信手机号
- `GET /api/user/address` - 获取地址列表
- `POST /api/user/address` - 创建地址
- `PUT /api/user/address` - 更新地址
- `DELETE /api/user/address` - 删除地址
- `GET /api/user/patient` - 获取就诊人列表
- `POST /api/user/patient` - 创建就诊人
- `PUT /api/user/patient` - 更新就诊人
- `DELETE /api/user/patient` - 删除就诊人

### 2. 订单管理模块

#### 数据模型
```go
type OrderModel struct {
    Id               int32      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
    OrderNo          string     `gorm:"column:orderNo;uniqueIndex;not null" json:"orderNo"`
    UserId           string     `gorm:"column:userId;not null;type:varchar(24)" json:"userId"`
    ServiceId        int32      `gorm:"column:serviceId;not null" json:"serviceId"`
    PatientId        *int32     `gorm:"column:patientId" json:"patientId"`
    AddressId        int32      `gorm:"column:addressId;not null" json:"addressId"`
    AppointmentDate  *string    `gorm:"column:appointmentDate" json:"appointmentDate"`
    AppointmentTime  *string    `gorm:"column:appointmentTime" json:"appointmentTime"`
    ServiceName      string     `gorm:"column:serviceName;not null" json:"serviceName"`
    Price            float64    `gorm:"column:price;not null" json:"price"`
    TotalAmount      float64    `gorm:"column:totalAmount;not null" json:"totalAmount"`
    Status           int        `gorm:"column:status;default:0" json:"status"`
    PayStatus        int        `gorm:"column:payStatus;default:0" json:"payStatus"`
    // ... 其他字段
}
```

#### 订单状态
- `0` - 待支付
- `1` - 已支付
- `2` - 已完成
- `3` - 已取消
- `4` - 已退款

#### 主要接口
- `POST /api/order/submit` - 提交订单
- `POST /api/order/smart-elderly` - 智慧养老设备订单
- `POST /api/order/pay/{orderId}` - 支付订单
- `POST /api/order/pay_confirm/{orderId}` - 支付确认
- `POST /api/order/cancel/{orderId}` - 取消订单
- `POST /api/order/refund/{orderId}` - 申请退款
- `GET /api/order/list` - 获取订单列表
- `POST /api/order/detail` - 获取订单详情
- `POST /api/order/time_slots` - 获取可用时间槽

### 3. 服务管理模块

#### 数据模型
```go
type ServiceItemModel struct {
    Id            int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
    Name          string    `gorm:"column:name;not null" json:"name"`
    Description   string    `gorm:"column:description" json:"description"`
    Category      string    `gorm:"column:category;not null" json:"category"`
    Price         float64   `gorm:"column:price;not null" json:"price"`
    OriginalPrice float64   `gorm:"column:originalPrice" json:"originalPrice"`
    ImageUrl      string    `gorm:"column:imageUrl" json:"imageUrl"`
    FormConfig    string    `gorm:"column:formConfig" json:"formConfig"`
    Status        int       `gorm:"column:status;default:1" json:"status"`
    // ... 其他字段
}
```

#### 主要接口
- `GET /api/service/list` - 获取服务列表
- `GET /api/service/categories` - 获取服务分类
- `GET /api/service/detail` - 获取服务详情
- `GET /api/service/form_config/{serviceId}` - 获取服务表单配置

### 4. 支付模块

#### 微信支付配置
```go
type WechatPayConfig struct {
    AppID       string `json:"appId"`
    MchID       string `json:"mchId"`
    MchKey      string `json:"mchKey"`
    NotifyURL   string `json:"notifyUrl"`
    Environment string `json:"environment"`
}
```

#### 主要接口
- `POST /api/payment/notify` - 微信支付回调

### 5. 咨询模块

#### 数据模型
```go
type Consultation struct {
    ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
    UserID      string    `json:"userId" gorm:"type:varchar(100);not null;index"`
    UserName    string    `json:"userName" gorm:"type:varchar(100);not null"`
    UserPhone   string    `json:"userPhone" gorm:"type:varchar(20)"`
    Status      string    `json:"status" gorm:"type:varchar(20);default:'waiting'"`
    CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
    UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
    // ... 关联关系
}
```

#### 主要接口
- `POST /api/consultation/create` - 创建咨询
- `GET /api/consultation/messages` - 获取咨询消息
- `POST /api/consultation/send` - 发送消息
- `GET /api/consultation/status` - 获取咨询状态
- `POST /api/consultation/close` - 关闭咨询

### 6. 推荐系统模块

#### 主要接口
- `GET /api/referral/qrcode` - 生成推荐二维码
- `GET /api/referral/report` - 获取推荐报告
- `GET /api/referral/config` - 获取推荐配置
- `POST /api/referral/apply_cashout` - 申请提现

### 7. 管理员模块

#### 主要接口
- `POST /api/admin/login` - 管理员登录
- `GET /api/admin/check-status` - 检查管理员状态
- `GET /api/admin/users` - 获取用户列表
- `GET /api/admin/orders` - 获取订单列表
- `POST /api/admin/set-admin` - 设置管理员
- `POST /api/admin/remove-admin` - 移除管理员
- `GET /api/admin/stats` - 获取统计数据

## 数据库设计

### 核心表结构

#### Users表
- 存储用户基本信息
- 支持微信登录
- 包含管理员权限字段

#### Orders表
- 存储订单信息
- 支持多种服务类型
- 包含支付和退款状态

#### ServiceItems表
- 存储服务项目信息
- 支持分类和表单配置
- 包含价格和状态管理

#### Consultations表
- 存储咨询会话信息
- 支持消息管理
- 包含状态跟踪

### 数据库连接配置

```go
// 数据库连接配置
source := "%s:%s@tcp(%s)/%s?readTimeout=1500ms&writeTimeout=1500ms&charset=utf8&loc=Local&&parseTime=true"
user := "root"
pwd := "bU4X6cFW"
addr := "10.3.110.11:3306"
dataBase := "anyuyinian"
```

### 连接池配置
```go
sqlDB.SetMaxIdleConns(100)    // 最大空闲连接数
sqlDB.SetMaxOpenConns(200)    // 最大打开连接数
sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生存时间
```

## 配置管理

### 微信配置
```go
type WxConfig struct {
    AppID     string
    AppSecret string
}
```

### 支付配置
```go
type PaymentConfig struct {
    WechatPay WechatPayConfig `json:"wechatPay"`
}
```

### COS配置
```go
type COSConfig struct {
    SecretID  string
    SecretKey string
    Region    string
    Bucket    string
    Domain    string
    ACL       string
}
```

## 日志系统

### 日志中间件
项目使用自定义的日志中间件记录所有API请求和响应：

```go
func NewLogMiddleware(handler http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 记录请求信息
        // 记录响应信息
        // 记录处理时间
    }
}
```

### 日志类型
- **API日志**: 记录HTTP请求和响应
- **数据库日志**: 记录SQL操作
- **步骤日志**: 记录业务处理步骤
- **错误日志**: 记录错误信息
- **信息日志**: 记录一般信息

## 开发规范

### 代码结构
1. **分层架构**: 严格按照MVC分层
2. **接口设计**: 使用接口定义DAO层
3. **错误处理**: 统一的错误处理机制
4. **日志记录**: 完整的日志记录系统

### 命名规范
1. **包名**: 小写字母，简短有意义
2. **结构体**: 大驼峰命名
3. **方法**: 大驼峰命名
4. **变量**: 小驼峰命名
5. **常量**: 全大写，下划线分隔

### 数据库规范
1. **表名**: 复数形式，首字母大写
2. **字段名**: 小驼峰命名
3. **索引**: 合理设置索引
4. **外键**: 正确设置外键关系

## 部署说明

### 环境要求
- Go 1.16+
- MySQL 5.7+
- 微信云托管环境

### 部署步骤
1. 配置数据库连接
2. 配置微信小程序参数
3. 配置支付参数
4. 配置COS存储
5. 运行数据库迁移
6. 启动服务

### 环境变量
```bash
# 微信配置
WX_APP_ID=wx101090677bd5219e
WX_APP_SECRET=042ff9921818ada9336df6e91fc2287e

# 支付配置
WECHAT_PAY_APP_ID=wx101090677bd5219e
WECHAT_PAY_MCH_ID=1726638701
WECHAT_PAY_MCH_KEY=JQzOCB8doIdgaUjAobELsk9nTyxdKhat
WECHAT_PAY_NOTIFY_URL=https://your-domain.com/api/payment/notify
WECHAT_PAY_ENVIRONMENT=production

# COS配置
COS_SECRET_ID=your_secret_id
COS_SECRET_KEY=your_secret_key
```

## API文档

### 通用响应格式
```json
{
    "code": 0,           // 0表示成功，非0表示失败
    "errorMsg": "",      // 错误信息
    "data": {}           // 响应数据
}
```

### 分页响应格式
```json
{
    "code": 0,
    "data": {
        "list": [],      // 数据列表
        "total": 100,    // 总数量
        "page": 1,       // 当前页
        "pageSize": 10,  // 每页大小
        "hasMore": true  // 是否有更多数据
    }
}
```

## 测试

### 测试结构
```
tests/
├── api/                    # API测试
├── backend/               # 后端测试
│   ├── admin/            # 管理员测试
│   ├── database/         # 数据库测试
│   ├── order/            # 订单测试
│   ├── payment/          # 支付测试
│   └── service/          # 服务测试
├── frontend/             # 前端测试
└── README.md             # 测试说明
```

### 运行测试
```bash
# 运行所有测试
go test ./...

# 运行特定包测试
go test ./service

# 运行特定测试
go test -run TestUserService
```

## 常见问题

### 1. 数据库连接问题
- 检查数据库配置
- 确认网络连接
- 验证用户权限

### 2. 微信登录问题
- 检查AppID和AppSecret
- 确认域名配置
- 验证签名算法

### 3. 支付问题
- 检查商户配置
- 确认回调地址
- 验证签名算法

### 4. 文件上传问题
- 检查COS配置
- 确认权限设置
- 验证文件格式

## 版本历史

### v1.3.0 (当前版本)
- ✅ 使用官方推荐的云托管API调用方式
- ✅ 优化错误处理和重试机制
- ✅ 完善订单Tab筛选功能
- ✅ 增强WebSocket连接稳定性

### v1.2.0
- ✅ 优化错误处理和重试机制
- ✅ 添加WebSocket支持
- ✅ 完善订单管理功能

### v1.1.0
- ✅ 添加WebSocket支持
- ✅ 优化用户界面
- ✅ 完善订单功能

### v1.0.0
- ✅ 初始版本发布
- ✅ 基础功能实现
- ✅ 用户认证功能

## 贡献指南

### 开发流程
1. Fork项目到个人仓库
2. 创建功能分支
3. 开发新功能或修复问题
4. 编写测试用例
5. 提交Pull Request

### 代码审查
- 所有代码变更需要经过审查
- 确保代码质量和安全性
- 保持代码风格的一致性

## 联系方式

- **项目维护**: 开发团队
- **技术支持**: 技术团队
- **问题反馈**: 通过GitHub Issues提交

## 许可证

本项目采用 MIT 许可证，详见 [LICENSE](LICENSE) 文件。
