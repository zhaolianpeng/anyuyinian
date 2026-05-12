# 服务下单和预约功能说明

## 功能概述

本系统实现了完整的服务下单和预约功能，包括服务浏览、预约下单、支付管理、订单跟踪等核心功能。

## 主要功能模块

### 1. 服务详情页面 (`pages/service/detail`)

#### 功能特点：
- **服务信息展示**：显示服务名称、描述、价格、分类等基本信息
- **就诊人选择**：支持选择就诊人信息，包括姓名、电话、身份证号
- **地址选择**：支持选择服务地址，包括联系人、电话、详细地址
- **动态表单**：根据服务配置动态生成预约表单字段
- **时间选择**：提供上午、下午、晚上三个时间段选择
- **备注信息**：支持添加特殊需求或备注信息

#### 表单字段类型：
- `text`：文本输入
- `number`：数字输入
- `date`：日期选择
- `radio`：单选框
- `select`：下拉选择
- `textarea`：多行文本

#### 验证规则：
- 必填字段验证
- 就诊人选择验证
- 地址选择验证
- 预约时间验证

### 2. 订单详情页面 (`pages/order/detail`)

#### 功能特点：
- **订单状态显示**：实时显示订单状态和状态描述
- **服务信息展示**：显示服务名称、价格、数量等信息
- **订单信息展示**：订单号、下单时间、支付时间、交易号等
- **预约信息展示**：解析并显示预约表单数据
- **费用明细**：显示服务费用、数量、总计等费用信息
- **操作按钮**：根据订单状态显示不同的操作按钮

#### 订单状态：
- `0`：待支付
- `1`：已支付
- `2`：已完成
- `3`：已取消
- `4`：已退款

#### 操作功能：
- **支付订单**：调用微信支付接口
- **取消订单**：取消待支付的订单
- **申请退款**：申请已支付订单的退款
- **联系客服**：显示客服联系方式
- **复制订单号**：复制订单号到剪贴板

### 3. 支付功能

#### 支付流程：
1. **生成支付参数**：调用后端接口生成微信支付参数
2. **调用微信支付**：使用微信小程序支付API
3. **支付确认**：支付成功后调用确认接口更新订单状态
4. **状态更新**：更新订单支付状态和支付时间

#### 支付状态：
- `0`：未支付
- `1`：已支付

### 4. 后端接口

#### 订单相关接口：
- `POST /api/order/submit`：提交订单
- `POST /api/order/pay/{orderId}`：支付订单
- `POST /api/order/pay_confirm/{orderId}`：支付确认
- `POST /api/order/cancel/{orderId}`：取消订单
- `POST /api/order/refund/{orderId}`：申请退款
- `GET /api/order/list`：获取订单列表
- `GET /api/order/detail/{orderId}`：获取订单详情

#### 服务相关接口：
- `GET /api/service/list`：获取服务列表
- `POST /api/service/detail`：获取服务详情
- `GET /api/service/form_config/{serviceId}`：获取服务表单配置

## 数据库设计

### 订单表 (Orders)
```sql
CREATE TABLE Orders (
  id INT PRIMARY KEY AUTO_INCREMENT,
  orderNo VARCHAR(50) UNIQUE NOT NULL,
  userId INT NOT NULL,
  serviceId INT NOT NULL,
  serviceName VARCHAR(100) NOT NULL,
  price DECIMAL(10,2) NOT NULL,
  quantity INT DEFAULT 1,
  totalAmount DECIMAL(10,2) NOT NULL,
  formData TEXT,
  status INT DEFAULT 0,
  payStatus INT DEFAULT 0,
  payTime DATETIME,
  payMethod VARCHAR(20),
  transactionId VARCHAR(100),
  refundStatus INT DEFAULT 0,
  refundTime DATETIME,
  refundAmount DECIMAL(10,2),
  refundReason TEXT,
  remark TEXT,
  referrerId INT,
  commission DECIMAL(10,2),
  createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
  updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### 服务表 (ServiceItems)
```sql
CREATE TABLE ServiceItems (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  category VARCHAR(50) NOT NULL,
  price DECIMAL(10,2) NOT NULL,
  originalPrice DECIMAL(10,2),
  imageUrl VARCHAR(255),
  detailImages TEXT,
  formConfig TEXT,
  status INT DEFAULT 1,
  sort INT DEFAULT 0,
  createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
  updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

## 使用流程

### 1. 服务预约流程
1. 用户浏览服务列表
2. 点击进入服务详情页
3. 选择就诊人信息
4. 选择服务地址
5. 填写预约表单
6. 选择预约时间
7. 添加备注信息（可选）
8. 提交订单
9. 完成支付

### 2. 订单管理流程
1. 查看订单详情
2. 根据订单状态进行相应操作
3. 支付订单（待支付状态）
4. 取消订单（待支付状态）
5. 申请退款（已支付状态）
6. 联系客服（任何状态）

## 技术特点

### 前端特点：
- **响应式设计**：适配不同屏幕尺寸
- **用户体验优化**：加载状态、错误提示、成功反馈
- **表单验证**：实时验证用户输入
- **状态管理**：根据订单状态显示不同操作
- **支付集成**：集成微信小程序支付API

### 后端特点：
- **RESTful API**：标准的REST接口设计
- **数据验证**：严格的参数验证
- **状态管理**：完整的订单状态流转
- **支付集成**：支持微信支付
- **错误处理**：完善的错误处理机制

## 安全考虑

1. **用户认证**：所有接口都需要用户认证
2. **数据验证**：前后端双重数据验证
3. **支付安全**：使用微信官方支付接口
4. **状态控制**：严格控制订单状态变更
5. **日志记录**：记录关键操作日志

## 扩展功能

### 可扩展的功能：
1. **多种支付方式**：支付宝、银行卡等
2. **优惠券系统**：支持优惠券使用
3. **积分系统**：支持积分抵扣
4. **评价系统**：服务完成后可评价
5. **推荐系统**：基于用户行为的服务推荐
6. **消息通知**：订单状态变更通知
7. **数据统计**：订单数据统计分析

## 测试建议

### 功能测试：
1. 测试完整的预约流程
2. 测试各种表单验证
3. 测试支付流程
4. 测试订单状态变更
5. 测试异常情况处理

### 性能测试：
1. 测试并发下单性能
2. 测试支付接口响应时间
3. 测试数据库查询性能

### 安全测试：
1. 测试接口权限控制
2. 测试数据验证机制
3. 测试支付安全性

## 部署说明

### 环境要求：
- Node.js 14+
- MySQL 5.7+
- 微信小程序开发工具
- HTTPS证书（生产环境）

### 配置项：
- 数据库连接配置
- 微信小程序配置
- 支付配置
- 文件上传配置

### 部署步骤：
1. 部署后端服务
2. 配置数据库
3. 配置微信小程序
4. 测试功能完整性
5. 上线运行 