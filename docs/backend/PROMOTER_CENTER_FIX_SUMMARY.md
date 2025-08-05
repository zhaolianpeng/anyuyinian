# 推广中心页面修复总结

## 问题描述

用户在点击推广中心页面时出现以下错误：

```
Error 1366: Incorrect integer value: '' for column 'referrerId' at row 1
```

## 问题原因

1. **数据库字段约束问题**: `Referrals` 表中的 `referrerId` 字段被定义为 `NOT NULL`，但代码中尝试插入空字符串 `""`
2. **数据类型不匹配**: 数据库期望整数类型，但接收到空字符串
3. **模型定义不一致**: Go 结构体中的字段类型与数据库定义不匹配

## 修复方案

### 1. 数据库迁移修复

创建了新的迁移文件 `fix_referrer_id_nullable.sql`：

```sql
-- 修改 Referrals 表的 referrerId 字段，允许为空
ALTER TABLE Referrals MODIFY COLUMN referrerId VARCHAR(24) NULL COMMENT '推荐人ID';

-- 更新现有的空字符串记录为 NULL
UPDATE Referrals SET referrerId = NULL WHERE referrerId = '';

-- 修改 Orders 表的 referrerId 字段，允许为空
ALTER TABLE Orders MODIFY COLUMN referrerId VARCHAR(24) NULL COMMENT '推荐人ID';

-- 更新现有的空字符串记录为 NULL
UPDATE Orders SET referrerId = NULL WHERE referrerId = '';
```

### 2. 模型定义修复

修改了 `ReferralModel` 结构体：

```go
// 修改前
ReferrerId string `gorm:"column:referrerId;not null;type:varchar(24)" json:"referrerId"`

// 修改后
ReferrerId *string `gorm:"column:referrerId;type:varchar(24)" json:"referrerId"` // 推荐人ID，可为空
```

### 3. 代码逻辑修复

修改了 `promoter_service.go` 和 `referral_service.go` 中的相关代码：

```go
// 修改前
ReferrerId: "", // 空字符串

// 修改后
ReferrerId: nil, // 使用 nil 表示没有推荐人
```

### 4. 空值检查修复

修改了推荐人信息获取逻辑：

```go
// 修改前
if referral.ReferrerId != "" {
    referrer, _ = dao.UserImp.GetUserByUserId(referral.ReferrerId)
}

// 修改后
if referral.ReferrerId != nil && *referral.ReferrerId != "" {
    referrer, _ = dao.UserImp.GetUserByUserId(*referral.ReferrerId)
}
```

## 修复文件列表

1. `db/migration/fix_referrer_id_nullable.sql` - 数据库迁移修复
2. `db/migration/create_service_order_referral_tables.sql` - 原始迁移文件修复
3. `db/model/referral.go` - 模型定义修复
4. `service/promoter_service.go` - 推广员服务修复
5. `service/referral_service.go` - 推荐服务修复

## 部署步骤

1. 运行数据库迁移：
   ```bash
   ./run_migration.sh
   ```

2. 重新构建项目：
   ```bash
   ./build.sh
   ```

3. 部署到云托管环境

## 验证方法

1. 访问推广中心页面
2. 检查是否不再出现数据库错误
3. 验证推广员信息正常显示
4. 测试推荐关系创建功能

## 注意事项

1. **数据兼容性**: 修改后的代码兼容现有的 NULL 值记录
2. **类型安全**: 使用指针类型确保类型安全
3. **向后兼容**: 保持 API 接口的向后兼容性

## 相关文档

- [数据库迁移指南](../migration/README.md)
- [推广中心功能说明](../features/PROMOTER_CENTER_FEATURES.md)
- [API 接口文档](../api/promoter_apis.md) 