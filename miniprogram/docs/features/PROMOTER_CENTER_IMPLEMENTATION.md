# 推广中心功能实现

## 功能概述

推广中心是一个完整的推广员管理系统，包含推广员信息展示、佣金记录管理、提现申请等功能。

## 后端实现

### 1. 数据库设计

#### 推广相关表结构
- **Referrals表**: 推荐关系表
- **Commissions表**: 佣金记录表  
- **Cashouts表**: 提现记录表

#### 主要字段
```sql
-- 推荐关系表
CREATE TABLE Referrals (
    id INT AUTO_INCREMENT PRIMARY KEY,
    userId VARCHAR(24) UNIQUE NOT NULL,
    referrerId VARCHAR(24) NOT NULL,
    qrCodeUrl VARCHAR(500),
    status INT DEFAULT 1,
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 佣金记录表
CREATE TABLE Commissions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    userId VARCHAR(24) NOT NULL,
    orderId INT NOT NULL,
    orderNo VARCHAR(50) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    rate DECIMAL(5,4) NOT NULL,
    status INT DEFAULT 0,
    cashoutTime DATETIME,
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 提现记录表
CREATE TABLE Cashouts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    userId VARCHAR(24) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    method VARCHAR(20) NOT NULL,
    account VARCHAR(200) NOT NULL,
    status INT DEFAULT 0,
    remark VARCHAR(500),
    processTime DATETIME,
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### 2. 数据访问层 (DAO)

#### 佣金DAO接口
```go
type CommissionInterface interface {
    CreateCommission(commission *model.CommissionModel) error
    GetCommissionById(id int32) (*model.CommissionModel, error)
    GetCommissionsByUserId(userId string, page, pageSize int) ([]*model.CommissionModel, int64, error)
    GetCommissionsByOrderId(orderId int32) (*model.CommissionModel, error)
    UpdateCommission(commission *model.CommissionModel) error
    DeleteCommission(id int32) error
    GetCommissionsByStatus(status int, page, pageSize int) ([]*model.CommissionModel, int64, error)
}
```

#### 提现DAO接口
```go
type CashoutInterface interface {
    CreateCashout(cashout *model.CashoutModel) error
    GetCashoutById(id int32) (*model.CashoutModel, error)
    GetCashoutsByUserId(userId string, page, pageSize int) ([]*model.CashoutModel, int64, error)
    UpdateCashout(cashout *model.CashoutModel) error
    DeleteCashout(id int32) error
    GetCashoutsByStatus(status int, page, pageSize int) ([]*model.CashoutModel, int64, error)
}
```

### 3. 业务逻辑层 (Service)

#### 推广服务 (`promoter_service.go`)
- **GetPromoterInfoHandler**: 获取推广员信息
- **GetCommissionListHandler**: 获取佣金记录列表
- **GetCashoutListHandler**: 获取提现记录列表
- **ApplyCashoutHandler**: 申请提现

#### 主要功能
1. **推广员信息管理**
   - 获取用户基本信息
   - 生成专属推广二维码
   - 统计推广数据

2. **佣金记录管理**
   - 分页查询佣金记录
   - 按状态筛选佣金
   - 佣金统计计算

3. **提现申请管理**
   - 提现申请验证
   - 可提现金额检查
   - 提现记录创建

### 4. API接口

#### 推广中心接口
```
GET /api/promoter/info - 获取推广员信息
GET /api/promoter/commission_list - 获取佣金记录列表
GET /api/promoter/cashout_list - 获取提现记录列表
POST /api/referral/apply_cashout - 申请提现
```

## 前端实现

### 1. 页面结构

#### 推广中心页面 (`pages/promoter/home`)
- **推广信息展示**: 用户信息、统计数据、推广二维码
- **标签页切换**: 推广信息、佣金记录、提现记录
- **佣金记录列表**: 分页展示佣金记录
- **提现记录列表**: 分页展示提现记录
- **提现申请弹窗**: 提现申请表单

### 2. 主要功能

#### 推广信息展示
```javascript
// 加载推广员信息
async loadPromoterInfo() {
  const userId = getCurrentUserId()
  const result = await api.promoterInfo({ userId })
  this.setData({ promoterInfo: result.data })
}
```

#### 佣金记录管理
```javascript
// 加载佣金记录列表
async loadCommissionList(refresh = false) {
  const page = refresh ? 1 : this.data.commissionPage
  const result = await api.commissionList({ 
    userId, 
    page, 
    pageSize: 20 
  })
  // 处理分页数据
}
```

#### 提现申请
```javascript
// 提交提现申请
async onSubmitCashout() {
  const { amount, method, account } = this.data.cashoutForm
  const result = await api.applyCashout({
    userId,
    amount: parseFloat(amount),
    method,
    account
  })
}
```

### 3. UI组件

#### 统计数据展示
```html
<view class="stats-grid">
  <view class="stat-item">
    <text class="stat-value">{{promoterInfo.totalIncome || 0}}</text>
    <text class="stat-label">总收入(元)</text>
  </view>
  <!-- 其他统计项 -->
</view>
```

#### 标签页切换
```html
<view class="tabs">
  <view class="tab-item {{activeTab === 0 ? 'active' : ''}}" data-index="0" bindtap="onTabChange">
    <text>推广信息</text>
  </view>
  <!-- 其他标签页 -->
</view>
```

#### 提现申请弹窗
```html
<view wx:if="{{showCashoutModal}}" class="modal-overlay">
  <view class="modal-content">
    <!-- 提现表单 -->
  </view>
</view>
```

## 功能特性

### 1. 推广员信息
- **用户信息展示**: 头像、昵称、推广码
- **统计数据**: 总收入、今日收入、本月收入、订单数
- **推广二维码**: 专属推广二维码生成和展示
- **分享功能**: 小程序分享功能

### 2. 佣金管理
- **佣金记录**: 分页展示佣金记录
- **状态管理**: 待结算、已结算、已提现
- **统计计算**: 按时间维度统计佣金
- **下拉刷新**: 支持下拉刷新和上拉加载

### 3. 提现管理
- **提现申请**: 在线申请提现
- **提现记录**: 查看提现历史记录
- **状态跟踪**: 提现状态实时更新
- **多种方式**: 支持微信、支付宝、银行卡提现

### 4. 用户体验
- **响应式设计**: 适配不同屏幕尺寸
- **加载状态**: 完善的加载状态提示
- **错误处理**: 友好的错误提示
- **数据缓存**: 合理的数据缓存策略

## 技术栈

### 后端
- **语言**: Go
- **框架**: 原生HTTP包
- **数据库**: MySQL
- **ORM**: GORM
- **日志**: 自定义日志系统

### 前端
- **框架**: 微信小程序
- **样式**: WXSS
- **模板**: WXML
- **逻辑**: JavaScript

## 部署说明

### 1. 数据库迁移
```bash
# 执行数据库迁移脚本
mysql -u username -p database_name < create_service_order_referral_tables.sql
```

### 2. 后端部署
```bash
# 编译Go程序
go build -o main .

# 运行服务
./main
```

### 3. 前端部署
- 将小程序代码上传到微信开发者工具
- 配置服务端域名与请求环境
- 发布小程序

## 测试验证

### 1. 功能测试
- [ ] 推广员信息获取
- [ ] 佣金记录查询
- [ ] 提现申请提交
- [ ] 提现记录查询
- [ ] 分享功能测试

### 2. 性能测试
- [ ] 接口响应时间
- [ ] 数据库查询性能
- [ ] 并发处理能力

### 3. 兼容性测试
- [ ] 不同微信版本兼容性
- [ ] 不同设备适配性
- [ ] 网络异常处理

## 总结

推广中心功能已经完整实现，包含：

1. **完整的后端API**: 支持推广员信息、佣金记录、提现申请等所有功能
2. **美观的前端界面**: 现代化的UI设计，良好的用户体验
3. **完善的数据管理**: 支持分页查询、状态管理、数据统计
4. **安全的业务逻辑**: 参数验证、权限控制、错误处理

该功能可以为推广员提供完整的推广管理工具，帮助他们更好地进行推广活动和收益管理。 