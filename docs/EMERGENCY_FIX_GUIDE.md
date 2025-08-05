# 紧急修复指南 - 解决"record not found"问题

## 问题描述

从日志中可以看到持续的"record not found"错误：
```
[ERROR] 数据库查询用户信息失败: record not found
[2.227ms] [rows:0] SELECT * FROM `Users` WHERE userId = '1' ORDER BY `Users`.`id` LIMIT 1
```

## 根本原因

1. **数据库迁移未完成**: 现有用户的`userId`字段为空
2. **前端使用旧格式**: 传递`userId: "1"`（数字格式）
3. **查询失败**: 数据库中找不到`userId = '1'`的用户

## 解决方案

### 1. 部署修复代码

首先需要部署包含紧急修复API的后端代码：

```bash
# 构建代码
go build -o main main.go

# 部署到云托管
# (根据你的部署流程)
```

### 2. 执行紧急修复

部署完成后，执行紧急修复脚本：

```bash
# 执行紧急修复
./scripts/emergency_fix.sh

# 或者手动执行API调用
curl -X POST "https://prod-5g94mx7a3d07e78c.service.tcloudbase.com/api/emergency/fix_user_ids"
```

### 3. 验证修复结果

```bash
# 检查用户状态
curl -X GET "https://prod-5g94mx7a3d07e78c.service.tcloudbase.com/api/emergency/user_status"

# 测试用户信息API
curl -X GET "https://prod-5g94mx7a3d07e78c.service.tcloudbase.com/api/emergency/test_user_info?userId=1"
```

## 新增的API端点

### 1. 紧急修复API
- **路径**: `/api/emergency/fix_user_ids`
- **方法**: `POST`
- **功能**: 为所有没有`userId`的用户生成新的`userId`

### 2. 用户状态检查API
- **路径**: `/api/emergency/user_status`
- **方法**: `GET`
- **功能**: 检查当前用户状态，显示需要修复的用户数量

### 3. 用户信息测试API
- **路径**: `/api/emergency/test_user_info`
- **方法**: `GET`
- **参数**: `userId` (可选，默认为"1")
- **功能**: 测试用户信息API是否正常工作

## 修复流程

### 步骤1: 检查当前状态
```bash
curl -X GET "https://prod-5g94mx7a3d07e78c.service.tcloudbase.com/api/emergency/user_status"
```

预期响应：
```json
{
  "code": 0,
  "data": {
    "totalUsers": 10,
    "usersWithUserId": 5,
    "usersWithoutUserId": 5,
    "fixNeeded": true,
    "sampleUsers": [...]
  }
}
```

### 步骤2: 执行修复
```bash
curl -X POST "https://prod-5g94mx7a3d07e78c.service.tcloudbase.com/api/emergency/fix_user_ids"
```

预期响应：
```json
{
  "code": 0,
  "data": {
    "message": "紧急修复完成",
    "fixedCount": 5,
    "totalUsers": 10,
    "usersWithUserId": 10,
    "success": true
  }
}
```

### 步骤3: 验证修复结果
```bash
curl -X GET "https://prod-5g94mx7a3d07e78c.service.tcloudbase.com/api/emergency/test_user_info?userId=1"
```

## 前端兼容性

前端已经实现了兼容性处理：

### 1. 自动检测旧格式
```javascript
// miniprogram/utils/user-id-compatibility.js
const needsUserIdMigration = (userId) => {
  return typeof userId === 'number' || 
         (typeof userId === 'string' && /^\d+$/.test(userId));
};
```

### 2. 自动处理错误
```javascript
// 在API调用失败时自动清除本地存储
if (errorMsg.includes("record not found")) {
  clearUserId();
  wx.navigateTo({ url: '/pages/login/login' });
}
```

### 3. 用户重新登录
- 检测到旧格式`userId`时提示重新登录
- 自动清除本地存储
- 跳转到登录页面

## 验证方法

### 1. 后端验证
```bash
# 检查用户状态
curl -X GET "https://prod-5g94mx7a3d07e78c.service.tcloudbase.com/api/emergency/user_status"

# 测试用户信息API
curl -X GET "https://prod-5g94mx7a3d07e78c.service.tcloudbase.com/api/user/info?userId=1"
```

### 2. 前端验证
- 清除本地存储的`userId`
- 重新登录获取新`userId`
- 验证所有功能正常工作

### 3. 数据库验证
```sql
-- 检查所有用户都有userId
SELECT COUNT(*) as total_users,
       COUNT(userId) as users_with_userid,
       COUNT(CASE WHEN userId IS NOT NULL AND userId != '' AND LENGTH(userId) = 24 THEN 1 END) as valid_userids
FROM Users;
```

## 预期结果

### ✅ 修复后
- 所有用户都有24位MongoDB风格的`userId`
- 用户信息API正常工作
- 前端自动处理兼容性
- 消除"record not found"错误

### 🔄 用户操作
- 用户可能需要重新登录
- 前端会自动检测并提示
- 或者自动清除本地存储

## 监控要点

### 1. 错误监控
- 监控"record not found"错误数量
- 监控用户信息API失败率
- 监控用户重新登录频率

### 2. 数据质量
- 确保所有用户都有有效的`userId`
- 监控`userId`格式正确性
- 定期检查数据一致性

### 3. 用户体验
- 监控用户登录成功率
- 监控功能使用成功率
- 收集用户反馈

## 总结

通过这个紧急修复方案，我们可以：

1. **快速解决当前问题**: 为现有用户生成`userId`
2. **确保数据一致性**: 所有用户都有正确的`userId`格式
3. **保持用户体验**: 前端自动处理兼容性
4. **预防未来问题**: 新用户自动生成正确的`userId`

执行修复后，系统应该能够正常运行，用户无需手动干预即可正常使用所有功能。 