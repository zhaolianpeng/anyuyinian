# 推广码系统说明

## 概述

推广码系统已经更新为使用六位随机字符串，而不是用户ID。每个用户都有唯一的六位推广码，格式为：大写字母和数字的组合（如：ABC123）。

## 数据库结构

### Referrals表
- `promoterCode` VARCHAR(6) UNIQUE - 六位推广码，唯一索引
- 其他字段保持不变

## API接口

### 1. 获取推广员信息
```
GET /api/promoter/info?userId={userId}
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "userId": "user_123",
    "promoterCode": "ABC123",
    "nickName": "张三",
    "avatarUrl": "https://...",
    "qrCodeUrl": "https://...",
    "totalIncome": 100.50,
    "todayIncome": 10.00,
    "monthIncome": 50.00,
    "totalOrders": 5,
    "todayOrders": 1,
    "monthOrders": 3
  }
}
```

### 2. 通过推广码查找用户
```
GET /api/promoter/find_user?promoterCode={promoterCode}
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "userId": "user_123",
    "nickName": "张三",
    "avatarUrl": "https://...",
    "promoterCode": "ABC123"
  }
}
```

### 3. 批量生成推广码
```
POST /api/promoter/generate_codes
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "totalCount": 10,
    "successCount": 8,
    "failedCount": 2,
    "results": [
      {
        "userId": "user_123",
        "promoterCode": "ABC123",
        "status": "成功"
      },
      {
        "userId": "user_456",
        "promoterCode": "",
        "status": "失败",
        "error": "数据库错误"
      }
    ]
  }
}
```

## 推广码生成规则

### 格式要求
- 长度：6位字符
- 字符集：大写字母（A-Z）和数字（0-9）
- 示例：ABC123, XYZ789, 123ABC

### 生成逻辑
1. 使用加密安全的随机数生成器
2. 从字符集 `ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789` 中随机选择
3. 检查数据库中是否已存在，如果存在则重新生成
4. 最大尝试次数：100次，避免无限循环

### 唯一性保证
- 数据库层面：`promoterCode` 字段有唯一索引
- 应用层面：生成时检查是否已存在
- 冲突处理：如果100次尝试后仍有冲突，使用时间戳作为备选方案

## 使用流程

### 1. 新用户注册
1. 用户注册时自动创建推荐关系记录
2. 系统自动生成唯一的六位推广码
3. 推广码与用户绑定，不可更改

### 2. 推广码使用
1. 用户分享自己的推广码给其他人
2. 其他用户输入推广码时，系统通过 `/api/promoter/find_user` 接口查找推广员
3. 建立推荐关系，记录推荐人信息

### 3. 佣金计算
1. 被推荐用户下单时，系统根据推荐关系计算佣金
2. 佣金记录与推广员关联
3. 推广员可以在推广中心查看佣金明细

## 迁移说明

### 现有用户处理
1. 运行数据库迁移脚本：`db/migration/generate_promoter_codes.sql`
2. 调用批量生成接口：`POST /api/promoter/generate_codes`
3. 验证所有用户都有有效的推广码

### 数据验证
```sql
-- 检查推广码格式
SELECT 
    id,
    userId,
    promoterCode,
    CASE 
        WHEN LENGTH(promoterCode) != 6 THEN '长度错误'
        WHEN promoterCode REGEXP '^[A-Z0-9]{6}$' = 0 THEN '格式错误'
        ELSE '格式正确'
    END as validation_result
FROM Referrals 
WHERE promoterCode IS NOT NULL AND promoterCode != '';
```

## 测试

### 运行测试脚本
```bash
cd tests/backend
./test_promoter_code.sh
```

### 手动测试
1. 生成推广码：`POST /api/promoter/generate_codes`
2. 获取推广员信息：`GET /api/promoter/info?userId={userId}`
3. 通过推广码查找：`GET /api/promoter/find_user?promoterCode={code}`

## 注意事项

1. **推广码唯一性**：每个推广码在系统中是唯一的，不能重复
2. **格式验证**：所有推广码必须符合6位大写字母+数字的格式
3. **不可修改**：推广码生成后不可修改，确保推广关系的稳定性
4. **性能考虑**：推广码查找使用数据库索引，查询性能良好
5. **安全性**：使用加密安全的随机数生成器，避免可预测性

## 错误处理

### 常见错误码
- `-1`：一般错误，查看 `errorMsg` 字段获取详细信息
- `400`：请求参数错误
- `405`：请求方法不支持

### 错误场景
1. 推广码不存在：返回"推广码不存在或已失效"
2. 推广码格式错误：返回"推广码格式无效"
3. 缺少参数：返回"缺少promoterCode参数" 