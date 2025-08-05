# 用户ID迁移总结

## 概述
本次改动将系统的用户ID从数据库自增ID改为MongoDB风格的随机字符串，同时保留原有的数据库ID字段。

## 主要改动

### 1. 数据库结构变更
- **Users表**：新增`userId`字段（VARCHAR(24)），用于存储随机字符串ID
- **UserAddresses表**：将`userId`字段类型从INT改为VARCHAR(24)
- **Patients表**：将`userId`字段类型从INT改为VARCHAR(24)

### 2. 后端代码变更

#### 新增文件
- `utils/id_generator.go`：ID生成工具
- `service/migration_service.go`：数据库迁移服务
- `service/test_service.go`：测试服务
- `db/migration/add_user_id_field.sql`：数据库迁移脚本

#### 修改文件
- `db/model/user.go`：新增UserId字段
- `db/model/user_extend.go`：修改UserId字段类型
- `db/dao/user_interface.go`：新增GetUserByUserId方法
- `db/dao/user_dao.go`：实现GetUserByUserId方法
- `db/dao/user_extend_interface.go`：修改UserId参数类型
- `db/dao/user_extend_dao.go`：修改UserId参数类型
- `service/wx_login_service.go`：在创建用户时生成UserId
- `service/user_service.go`：支持通过UserId查询用户

### 3. 前端代码变更
- 登录成功后保存新的UserId格式
- API调用使用新的UserId格式
- 所有涉及用户ID的地方都使用字符串格式

## 迁移步骤

### 1. 数据库迁移
```sql
-- 为Users表添加userId字段
ALTER TABLE Users ADD COLUMN userId VARCHAR(24) UNIQUE NOT NULL COMMENT '用户唯一标识符' AFTER id;

-- 为UserAddresses表修改userId字段类型
ALTER TABLE UserAddresses MODIFY COLUMN userId VARCHAR(24) NOT NULL;

-- 为Patients表修改userId字段类型
ALTER TABLE Patients MODIFY COLUMN userId VARCHAR(24) NOT NULL;

-- 添加索引
CREATE INDEX idx_user_id ON Users(userId);
CREATE INDEX idx_user_id ON UserAddresses(userId);
CREATE INDEX idx_user_id ON Patients(userId);
```

### 2. 代码部署
1. 部署新的后端代码
2. 运行数据库迁移脚本
3. 执行用户数据迁移（为现有用户生成UserId）

### 3. 数据迁移
使用迁移服务为现有用户生成UserId：
```bash
curl -X POST http://localhost:8080/api/test/migrate-userid
```

## API变更

### 登录接口
**响应格式变更**：
```json
{
  "code": 0,
  "data": {
    "id": 1,           // 数据库自增ID（保留）
    "userId": "507f1f77bcf86cd799439011",  // 新的随机字符串ID
    "openId": "wx_openid",
    "nickName": "用户昵称",
    // ... 其他字段
  }
}
```

### 用户信息接口
**请求参数**：
```
GET /api/user/info?userId=507f1f77bcf86cd799439011
```

**响应格式**：
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "userId": "507f1f77bcf86cd799439011",
    "nickName": "用户昵称",
    // ... 其他字段
  }
}
```

## 兼容性说明

### 向后兼容
- 保留原有的数据库ID字段，确保现有功能不受影响
- 新用户自动生成UserId，老用户通过迁移服务生成UserId
- API支持两种ID格式，优先使用UserId

### 前端兼容
- 前端存储使用新的UserId格式
- 所有API调用使用UserId
- 保持用户体验的一致性

## 测试验证

### 1. 单元测试
- ID生成工具测试
- 数据库查询测试
- API接口测试

### 2. 集成测试
```bash
# 运行测试脚本
./tests/backend/test_user_id_migration.sh
```

### 3. 功能测试
- 新用户注册流程
- 老用户登录流程
- 地址管理功能
- 就诊人管理功能

## 注意事项

1. **数据一致性**：确保所有相关表都使用相同的UserId格式
2. **性能考虑**：UserId字段添加了索引，查询性能良好
3. **安全性**：使用加密随机数生成UserId，确保唯一性
4. **回滚方案**：保留原有ID字段，如需回滚可快速切换

## 部署检查清单

- [ ] 数据库迁移脚本已执行
- [ ] 后端代码已部署
- [ ] 用户数据迁移已完成
- [ ] API接口测试通过
- [ ] 前端代码已更新
- [ ] 功能测试通过
- [ ] 性能测试通过

## 监控指标

1. **迁移进度**：已迁移用户数量
2. **API性能**：使用UserId查询的响应时间
3. **错误率**：UserId相关的错误数量
4. **用户体验**：登录和功能使用是否正常 