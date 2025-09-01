# 智慧养老分类显示逻辑说明

## 问题描述

用户反馈服务列表页面看不到新增的智慧养老分类，需要实现动态显示逻辑：**只有在后端有智慧养老服务项时，才在服务列表页面显示智慧养老分类，否则默认隐藏**。

## 解决方案

### 1. 后端实现

#### 1.1 新增分类API端点
- **端点**: `GET /api/service/categories`
- **功能**: 动态返回可用的服务分类及其数量
- **逻辑**: 只返回有服务的分类，空分类自动过滤

#### 1.2 数据库查询优化
```go
// 新增方法：GetServiceCategoriesWithCount
func (imp *ServiceInterfaceImp) GetServiceCategoriesWithCount() ([]CategoryCount, error) {
    var categories []CategoryCount
    cli := db.Get()
    
    // 查询每个分类的服务数量
    err := cli.Table(serviceTableName).
        Where("status = ?", 1).
        Select("category, COUNT(*) as count").
        Group("category").
        Order("count DESC").
        Scan(&categories).Error
    
    // 只返回有服务的分类
    for i := range categories {
        if categories[i].Count > 0 {
            // 设置分类显示名称
            switch categories[i].Category {
            case "智慧养老":
                categories[i].Name = "智慧养老"
            // ... 其他分类
            }
        }
    }
    
    return categories, nil
}
```

#### 1.3 分类过滤逻辑
```go
// 在ServiceCategoriesHandler中
for _, cat := range categories {
    if cat.Count > 0 { // 只显示有服务的分类
        categoryList = append(categoryList, ServiceCategory{
            Name:  cat.Name,
            Value: cat.Category,
            Count: cat.Count,
        })
    }
}
```

### 2. 前端实现

#### 2.1 动态加载分类
```javascript
// 修改服务列表页面
data: {
  categories: [
    { name: '全部', value: '' } // 只保留"全部"分类
  ]
},

// 新增loadCategories方法
async loadCategories() {
  try {
    const res = await api.serviceCategories()
    if (res.code === 0 && res.data && res.data.categories) {
      this.setData({
        categories: res.data.categories
      })
    }
  } catch (error) {
    console.error('加载分类失败:', error)
  }
}
```

#### 2.2 API配置
```javascript
// 在cloud-container-standard.js中添加
serviceCategories: () => callContainer('/api/service/categories', 'GET'),
```

### 3. 显示逻辑

#### 3.1 智慧养老分类显示条件
- ✅ **有智慧养老服务**: 分类显示在服务列表页面
- ❌ **无智慧养老服务**: 分类自动隐藏，不显示在分类标签中

#### 3.2 其他分类显示条件
- ✅ **有服务**: 分类正常显示
- ❌ **无服务**: 分类自动隐藏

#### 3.3 分类排序
- 按服务数量降序排列
- "全部"分类始终排在第一位

## 测试验证

### 1. 测试脚本
```bash
# 测试分类API
./tests/smart_elderly/test_categories_api.sh
```

### 2. 测试场景

#### 场景1: 有智慧养老服务
1. 数据库中有智慧养老设备数据
2. 分类API返回智慧养老分类
3. 前端服务列表页面显示智慧养老分类标签
4. 用户可以点击智慧养老分类查看相关服务

#### 场景2: 无智慧养老服务
1. 数据库中没有智慧养老设备数据
2. 分类API不返回智慧养老分类
3. 前端服务列表页面不显示智慧养老分类标签
4. 用户看不到智慧养老分类

### 3. 验证要点
- ✅ 分类API正常工作
- ✅ 智慧养老分类数量匹配数据库
- ✅ 空分类被正确过滤
- ✅ 前端动态加载分类
- ✅ 分类显示逻辑符合预期

## 文件清单

### 后端文件
- `db/dao/service_dao.go` - 新增GetServiceCategoriesWithCount方法
- `db/dao/service_interface.go` - 更新接口定义
- `service/service_service.go` - 新增ServiceCategoriesHandler
- `main.go` - 注册分类API端点

### 前端文件
- `pages/service/list.js` - 修改为动态加载分类
- `utils/cloud-container-standard.js` - 添加serviceCategories API

### 测试文件
- `tests/smart_elderly/test_categories_api.sh` - 分类API测试脚本

## 部署说明

1. **数据库迁移**: 确保智慧养老设备数据已插入
2. **后端部署**: 重启服务以加载新的API端点
3. **前端部署**: 更新小程序代码
4. **测试验证**: 运行分类API测试脚本

## 注意事项

1. **缓存考虑**: 分类数据变化不频繁，可以考虑添加缓存
2. **性能优化**: 分类查询使用GROUP BY，确保索引优化
3. **错误处理**: 分类加载失败时，前端显示默认分类
4. **向后兼容**: 保持现有分类显示逻辑不变

## 总结

通过实现动态分类加载机制，智慧养老分类现在会根据后端服务数据自动显示或隐藏，确保用户界面的一致性和准确性。这种设计既满足了功能需求，又保持了系统的灵活性和可维护性。
