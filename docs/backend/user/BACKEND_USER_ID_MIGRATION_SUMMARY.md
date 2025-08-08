# 后端UserId迁移总结

## 概述

将后端数据库中所有表的userId字段从INT类型改为VARCHAR(24)类型，使用MongoDB风格的随机字符串作为用户唯一标识符。

## 修改的文件

### 1. 模型文件 (Model)

#### `anyuyinian/db/model/order.go`
- **修改**: `UserId`字段从`int32`改为`string`
- **GORM标签**: 添加`type:varchar(24)`
- **影响**: 订单表中的用户ID字段

#### `anyuyinian/db/model/kefu.go`
- **修改**: `UserId`和`ReplyUserId`字段从`int32`改为`string`
- **GORM标签**: 添加`type:varchar(24)`
- **影响**: 客服消息表中的用户ID和回复用户ID字段

#### `anyuyinian/db/model/referral.go`
- **修改**: 所有模型中的`UserId`字段从`int32`改为`string`
- **影响**: 
  - `ReferralModel`: 推荐关系表
  - `CommissionModel`: 佣金记录表
  - `CashoutModel`: 提现记录表

#### `anyuyinian/db/model/upload.go`
- **修改**: `UserId`字段从`int32`改为`string`
- **GORM标签**: 添加`type:varchar(24)`
- **影响**: 文件上传表中的用户ID字段

### 2. DAO接口文件

#### `anyuyinian/db/dao/order_interface.go`
- **修改**: 
  - `GetOrdersByUserId(userId int32, page, pageSize int)` → `GetOrdersByUserId(userId string, page, pageSize int)`
  - `GetOrdersByStatusAndUserId(status int, userId int32, page, pageSize int)` → `GetOrdersByStatusAndUserId(status int, userId string, page, pageSize int)`

#### `anyuyinian/db/dao/order_dao.go`
- **修改**: 实现方法中的userId参数类型从`int32`改为`string`

#### `anyuyinian/db/dao/kefu_interface.go`
- **修改**:
  - `GetMessagesByUserId(userId int32, page, pageSize int)` → `GetMessagesByUserId(userId string, page, pageSize int)`
  - `ReplyMessage(id int32, replyContent string, replyUserId int32)` → `ReplyMessage(id int32, replyContent string, replyUserId string)`

#### `anyuyinian/db/dao/kefu_dao.go`
- **修改**: 实现方法中的userId和replyUserId参数类型从`int32`改为`string`

#### `anyuyinian/db/dao/referral_interface.go`
- **修改**:
  - `GetReferralByUserId(userId int32)` → `GetReferralByUserId(userId string)`
  - `GetCommissionsByUserId(userId int32, page, pageSize int)` → `GetCommissionsByUserId(userId string, page, pageSize int)`
  - `GetCashoutsByUserId(userId int32, page, pageSize int)` → `GetCashoutsByUserId(userId string, page, pageSize int)`

#### `anyuyinian/db/dao/referral_dao.go`
- **修改**: 实现方法中的userId参数类型从`int32`改为`string`

#### `anyuyinian/db/dao/upload_interface.go`
- **修改**: `GetFilesByUserId(userId int32, limit int)` → `GetFilesByUserId(userId string, limit int)`

#### `anyuyinian/db/dao/upload_dao.go`
- **修改**: 实现方法中的userId参数类型从`int32`改为`string`

### 3. 服务文件

#### `anyuyinian/service/order_service.go`
- **修改**:
  - `SubmitOrderRequest.UserId`从`int32`改为`string`
  - `OrderListRequest.UserId`从`int32`改为`string`
  - 所有使用userId的地方都改为字符串类型
  - 修复CommissionModel创建时的UserId字段类型转换

### 4. 迁移文件

#### `anyuyinian/db/migration/update_all_userid_fields.sql`
- **新增**: 数据库迁移脚本，将所有表的userId字段从INT改为VARCHAR(24)

#### `anyuyinian/service/migration_service.go`
- **扩展**: 添加`MigrateAllTablesUserId()`方法，处理所有表的userId字段迁移

## 数据库表修改

### 需要修改的表和字段

1. **Users表**
   - `userId`: INT → VARCHAR(24)

2. **Orders表**
   - `userId`: INT → VARCHAR(24)

3. **KefuMessages表**
   - `userId`: INT → VARCHAR(24)
   - `replyUserId`: INT → VARCHAR(24)

4. **Referrals表**
   - `userId`: INT → VARCHAR(24)

5. **Commissions表**
   - `userId`: INT → VARCHAR(24)

6. **Cashouts表**
   - `userId`: INT → VARCHAR(24)

7. **Files表**
   - `userId`: INT → VARCHAR(24)

8. **UserAddresses表**
   - `userId`: INT → VARCHAR(24)

9. **Patients表**
   - `userId`: INT → VARCHAR(24)

## 迁移步骤

### 1. 数据库结构迁移

```sql
-- 执行update_all_userid_fields.sql中的SQL语句
ALTER TABLE Orders MODIFY COLUMN userId VARCHAR(24) NOT NULL;
ALTER TABLE KefuMessages MODIFY COLUMN userId VARCHAR(24) NOT NULL;
-- ... 其他表的修改
```

### 2. 数据迁移

使用Go代码执行数据迁移：

```go
migrationService := NewMigrationService()

// 迁移用户表的userId
err := migrationService.MigrateExistingUsers()
if err != nil {
    log.Fatal("用户userId迁移失败:", err)
}

// 迁移所有表的userId字段
err = migrationService.MigrateAllTablesUserId()
if err != nil {
    log.Fatal("所有表userId迁移失败:", err)
}

// 验证迁移结果
err = migrationService.ValidateUserIds()
if err != nil {
    log.Fatal("userId验证失败:", err)
}
```

### 3. 索引创建

```sql
-- 为所有表添加userId索引
CREATE INDEX idx_user_id ON Orders(userId);
CREATE INDEX idx_user_id ON KefuMessages(userId);
CREATE INDEX idx_user_id ON Referrals(userId);
CREATE INDEX idx_user_id ON Commissions(userId);
CREATE INDEX idx_user_id ON Cashouts(userId);
CREATE INDEX idx_user_id ON Files(userId);
```

## 兼容性处理

### 1. 类型转换
- 在需要的地方使用`fmt.Sprintf("%d", intValue)`将int32转换为string
- 在API参数解析时直接使用字符串类型的userId

### 2. 错误处理
- 添加了完善的错误处理和日志记录
- 迁移过程中如果某个记录失败，会继续处理其他记录

### 3. 验证机制
- 提供验证方法确保所有用户都有有效的userId
- 检查userId格式是否为24位十六进制字符串

## 部署注意事项

### 1. 备份数据
在迁移前务必备份所有相关表的数据

### 2. 停机时间
建议在低峰期进行迁移，避免影响用户使用

### 3. 回滚方案
如果迁移失败，可以：
- 恢复数据库备份
- 回滚代码到之前的版本
- 重新执行迁移

### 4. 监控
- 监控迁移过程中的错误日志
- 验证迁移后的数据完整性
- 检查API调用是否正常

## 测试验证

### 1. 单元测试
- 测试所有修改的DAO方法
- 验证模型字段类型正确性
- 检查服务层逻辑

### 2. 集成测试
- 测试完整的API调用流程
- 验证数据库查询和更新操作
- 检查前端和后端的数据交互

### 3. 性能测试
- 验证索引是否正常工作
- 检查查询性能是否有影响
- 测试大数据量下的表现

## 总结

这次迁移涉及了后端系统的核心用户标识符，需要谨慎处理：

1. **数据完整性**: 确保所有相关表的数据都正确迁移
2. **系统稳定性**: 保证迁移过程中系统的可用性
3. **性能影响**: 新索引的创建和查询性能的优化
4. **兼容性**: 确保前后端API调用的兼容性

通过这次迁移，系统将使用更安全和可扩展的MongoDB风格用户ID，为未来的系统扩展奠定基础。 