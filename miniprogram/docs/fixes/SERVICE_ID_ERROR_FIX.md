# 服务ID错误问题解决指南

## 问题描述

当访问服务详情页面时，出现"无效的服务ID"错误，通常是因为数据库中不存在对应的服务记录。

## 问题原因

1. **数据库中没有服务数据**
2. **服务ID不存在或已被删除**
3. **服务状态为下架**
4. **数据库连接问题**

## 解决步骤

### 1. 检查数据库连接

确保数据库服务正常运行，并且应用能够连接到数据库。

### 2. 初始化数据库数据

#### 方法一：使用初始化脚本

```bash
# 进入后端目录
cd anyuyinian

# 修改数据库配置
vim init_database.sh

# 执行初始化脚本
chmod +x init_database.sh
./init_database.sh
```

#### 方法二：手动执行SQL

```sql
-- 连接到数据库
mysql -u root -p

-- 选择数据库
USE anyuyinian;

-- 检查现有服务
SELECT * FROM ServiceItems;

-- 如果没有数据，执行初始化脚本
source init_service_data.sql;
```

### 3. 验证服务数据

执行以下SQL查询验证服务数据：

```sql
-- 查看所有服务
SELECT id, name, category, price, status FROM ServiceItems ORDER BY id;

-- 查看上架的服务
SELECT id, name, category, price FROM ServiceItems WHERE status = 1 ORDER BY sort;

-- 统计服务数量
SELECT 
    COUNT(*) as total_services,
    COUNT(CASE WHEN status = 1 THEN 1 END) as active_services
FROM ServiceItems;
```

### 4. 测试服务接口

使用提供的测试脚本验证服务接口：

```javascript
// 在小程序开发工具中运行
const { runServiceTests } = require('./test_service_data.js')
runServiceTests()
```

### 5. 检查服务ID

确保访问的服务ID在数据库中存在：

```sql
-- 检查特定服务ID
SELECT * FROM ServiceItems WHERE id = 7;
```

## 可用的服务ID

根据初始化脚本，以下服务ID应该可用：

| ID | 服务名称 | 分类 | 价格 |
|----|----------|------|------|
| 1 | 专业陪诊套餐 | 医院陪诊 | ¥299 |
| 2 | 上门护理服务 | 专业护理 | ¥199 |
| 3 | 生活照护服务 | 生活照护 | ¥150 |
| 4 | 康复护理服务 | 专业护理 | ¥250 |
| 5 | 健康体检套餐A | 体检 | ¥299 |
| 6 | 健康体检套餐B | 体检 | ¥599 |
| 7 | 心理咨询服务 | 心理咨询 | ¥200 |
| 8 | 营养师咨询 | 营养咨询 | ¥150 |

## 常见问题

### Q1: 为什么服务ID 7不存在？

**A1**: 可能是因为：
- 数据库中没有执行初始化脚本
- 服务被手动删除
- 数据库连接配置错误

**解决方案**：
1. 执行初始化脚本
2. 检查数据库连接
3. 验证服务数据

### Q2: 如何添加新的服务？

**A2**: 可以通过以下方式添加：

```sql
INSERT INTO ServiceItems (name, description, category, price, originalPrice, imageUrl, formConfig, status, sort) VALUES
('新服务名称', '服务描述', '服务分类', 价格, 原价, '图片URL', '表单配置JSON', 1, 排序值);
```

### Q3: 如何修改服务状态？

**A3**: 使用以下SQL：

```sql
-- 上架服务
UPDATE ServiceItems SET status = 1 WHERE id = 服务ID;

-- 下架服务
UPDATE ServiceItems SET status = 0 WHERE id = 服务ID;
```

## 调试技巧

### 1. 检查后端日志

查看后端服务日志，确认：
- 数据库连接是否正常
- SQL查询是否成功
- 错误信息详情

### 2. 使用测试脚本

运行测试脚本验证各个接口：

```javascript
// 测试服务列表
const { testServiceList } = require('./test_service_data.js')
testServiceList()

// 测试特定服务
const { testServiceDetail } = require('./test_service_data.js')
testServiceDetail(7)
```

### 3. 检查网络请求

在小程序开发工具中查看网络请求：
- 确认请求URL正确
- 检查响应状态码
- 查看响应内容

## 预防措施

1. **定期备份数据**：定期备份服务数据
2. **监控服务状态**：监控服务接口的可用性
3. **数据验证**：在添加服务时进行数据验证
4. **错误处理**：完善错误处理机制

## 联系支持

如果问题仍然存在，请提供以下信息：
- 错误信息截图
- 数据库查询结果
- 后端服务日志
- 网络请求详情 