# SQL日志系统实现指南

## 概述

本文档详细说明了如何为后端所有接口的SQL操作添加日志记录。通过统一的日志系统，我们可以更好地监控和调试数据库操作。

## 已实现的功能

### 1. SQL日志记录器 (`db/dao/sql_logger.go`)

创建了专门的SQL日志记录器，提供以下功能：

- **操作开始日志**: 记录SQL操作的开始时间、操作类型、表名和参数
- **操作结果日志**: 记录SQL操作的执行结果、耗时和影响行数
- **分类日志方法**: 
  - `LogQuery()` - 查询操作
  - `LogInsert()` - 插入操作
  - `LogUpdate()` - 更新操作
  - `LogDelete()` - 删除操作
  - `LogCount()` - 计数操作

### 2. 日志格式

所有SQL日志都遵循统一的格式：

```
[SQL] 操作类型 表名: 详细信息 (耗时: 时间)
```

示例：
```
[SQL] 开始 查询 表: ServiceItems, 参数: {"id":123,"status":1}
[SQL] 查询 表: ServiceItems 成功 (耗时: 2.5ms), 结果: {"id":123,"name":"测试服务"}
```

## 已完成的DAO文件

### 1. service_dao.go ✅
- `GetServiceById()` - 根据ID获取服务
- `GetServicesByCategory()` - 根据分类获取服务列表
- `GetAllServices()` - 获取所有服务列表
- `GetServiceCategories()` - 获取服务分类列表
- `GetServiceCategoriesWithCount()` - 获取服务分类及数量
- `CreateService()` - 创建服务
- `UpdateService()` - 更新服务
- `DeleteService()` - 删除服务

### 2. order_dao.go 🔄 (部分完成)
- `CreateOrder()` - 创建订单 ✅
- 其他方法待完成

### 3. user_dao.go 🔄 (部分完成)
- `GetUserByOpenId()` - 根据OpenId查询用户 ✅
- 其他方法待完成

## 待完成的DAO文件

需要为以下DAO文件添加SQL日志：

1. **order_dao.go** - 订单相关操作
2. **user_dao.go** - 用户相关操作
3. **home_dao.go** - 首页数据操作
4. **admin_dao.go** - 管理员操作
5. **upload_dao.go** - 文件上传操作
6. **referral_dao.go** - 推荐相关操作
7. **commission_dao.go** - 佣金相关操作
8. **cashout_dao.go** - 提现相关操作
9. **kefu_dao.go** - 客服相关操作
10. **config_dao.go** - 配置相关操作
11. **consultation_dao.go** - 咨询相关操作
12. **user_extend_dao.go** - 用户扩展信息操作

## 添加SQL日志的标准模式

### 1. 查询操作模式

```go
func (imp *ExampleInterfaceImp) GetExampleById(id int32) (*model.ExampleModel, error) {
    var example = new(model.ExampleModel)
    cli := db.Get()
    
    // 记录SQL操作日志
    logger := NewSQLLogger("查询", "ExampleTable", map[string]interface{}{
        "id": id,
    })
    
    err := cli.Table("ExampleTable").Where("id = ?", id).First(example).Error
    logger.LogQuery(example, err)
    
    return example, err
}
```

### 2. 插入操作模式

```go
func (imp *ExampleInterfaceImp) CreateExample(example *model.ExampleModel) error {
    cli := db.Get()
    example.CreatedAt = time.Now()
    example.UpdatedAt = time.Now()
    
    // 记录SQL操作日志
    logger := NewSQLLogger("插入", "ExampleTable", map[string]interface{}{
        "name": example.Name,
        "type": example.Type,
    })
    
    err := cli.Table("ExampleTable").Create(example).Error
    logger.LogInsert(example, err)
    
    return err
}
```

### 3. 更新操作模式

```go
func (imp *ExampleInterfaceImp) UpdateExample(example *model.ExampleModel) error {
    cli := db.Get()
    example.UpdatedAt = time.Now()
    
    // 记录SQL操作日志
    logger := NewSQLLogger("更新", "ExampleTable", map[string]interface{}{
        "id":   example.Id,
        "name": example.Name,
    })
    
    result := cli.Table("ExampleTable").Where("id = ?", example.Id).Updates(example)
    logger.LogUpdate(result.RowsAffected, result.Error)
    
    return result.Error
}
```

### 4. 删除操作模式

```go
func (imp *ExampleInterfaceImp) DeleteExample(id int32) error {
    cli := db.Get()
    
    // 记录SQL操作日志
    logger := NewSQLLogger("删除", "ExampleTable", map[string]interface{}{
        "id": id,
    })
    
    result := cli.Table("ExampleTable").Where("id = ?", id).Delete(&model.ExampleModel{})
    logger.LogDelete(result.RowsAffected, result.Error)
    
    return result.Error
}
```

### 5. 分页查询模式

```go
func (imp *ExampleInterfaceImp) GetExamplesByType(typeName string, page, pageSize int) ([]*model.ExampleModel, int64, error) {
    var examples []*model.ExampleModel
    var total int64
    cli := db.Get()

    // 记录SQL操作日志
    logger := NewSQLLogger("查询", "ExampleTable", map[string]interface{}{
        "type":     typeName,
        "page":     page,
        "pageSize": pageSize,
    })

    // 获取总数
    err := cli.Table("ExampleTable").Where("type = ?", typeName).Count(&total).Error
    if err != nil {
        logger.LogCount(0, err)
        return nil, 0, err
    }
    logger.LogCount(total, nil)

    // 获取分页数据
    offset := (page - 1) * pageSize
    err = cli.Table("ExampleTable").
        Where("type = ?", typeName).
        Order("createdAt DESC").
        Offset(offset).
        Limit(pageSize).
        Find(&examples).Error

    logger.LogQuery(examples, err)
    return examples, total, err
}
```

## 实施步骤

### 步骤1: 检查现有日志系统
确保 `service/log_middleware.go` 中的日志系统正常工作。

### 步骤2: 为每个DAO文件添加日志
按照上述模式，为每个DAO方法添加SQL日志：

1. 在方法开始处创建 `NewSQLLogger`
2. 在SQL操作后调用相应的日志方法
3. 确保记录关键参数和结果

### 步骤3: 测试日志输出
运行应用程序并检查日志输出：

```bash
# 启动应用
go run main.go

# 在另一个终端测试API
curl -X GET "http://localhost:8080/api/service/list?page=1&pageSize=10"
```

### 步骤4: 验证日志格式
确保日志输出格式正确，包含：
- 操作类型
- 表名
- 参数信息
- 执行结果
- 耗时信息

## 日志级别和过滤

### 日志级别
- `[SQL]` - SQL操作日志
- `[API]` - API请求日志
- `[STEP]` - 中间步骤日志
- `[ERROR]` - 错误日志
- `[INFO]` - 信息日志

### 日志过滤
可以通过grep命令过滤特定类型的日志：

```bash
# 只查看SQL日志
tail -f app.log | grep "\[SQL\]"

# 只查看错误日志
tail -f app.log | grep "\[ERROR\]"

# 查看特定表的操作
tail -f app.log | grep "ServiceItems"
```

## 性能考虑

### 1. 日志性能影响
- SQL日志记录会增加少量性能开销
- 建议在生产环境中根据需要调整日志级别

### 2. 日志存储
- 考虑使用日志轮转避免日志文件过大
- 可以集成ELK等日志分析系统

### 3. 敏感信息
- 避免在日志中记录敏感信息（如密码、token等）
- 对敏感字段进行脱敏处理

## 监控和告警

### 1. 慢查询监控
通过日志可以识别执行时间较长的SQL操作：

```bash
# 查找执行时间超过100ms的SQL操作
grep "\[SQL\].*耗时.*[1-9][0-9][0-9]ms" app.log
```

### 2. 错误率监控
监控SQL操作的错误率：

```bash
# 统计SQL错误数量
grep "\[SQL\].*失败" app.log | wc -l
```

### 3. 操作频率监控
统计各种SQL操作的频率：

```bash
# 统计查询操作数量
grep "\[SQL\].*查询.*成功" app.log | wc -l

# 统计插入操作数量
grep "\[SQL\].*插入.*成功" app.log | wc -l
```

## 总结

通过实施SQL日志系统，我们可以：

1. **提高可观测性**: 清楚了解每个SQL操作的执行情况
2. **快速定位问题**: 通过日志快速定位数据库相关问题
3. **性能优化**: 识别慢查询和性能瓶颈
4. **审计追踪**: 记录所有数据库操作的详细信息
5. **监控告警**: 基于日志数据建立监控和告警机制

建议按优先级逐步为所有DAO文件添加SQL日志，优先处理核心业务模块（如订单、用户、服务等）。
