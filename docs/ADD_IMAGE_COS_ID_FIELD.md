# 为 ServiceItems 表添加 imageCosId 字段

## 概述
为 `ServiceItems` 表新增 `imageCosId` 字段，用于存储腾讯云对象存储的图片ID，以便更好地管理服务图片资源。

## 修改内容

### 1. 数据库模型更新
**文件**: `db/model/service.go`
```go
type ServiceItemModel struct {
    // ... 其他字段
    ImageUrl      string    `gorm:"column:imageUrl" json:"imageUrl"`
    ImageCosId    string    `gorm:"column:imageCosId" json:"imageCosId"`     // 腾讯云对象存储ID
    DetailImages  string    `gorm:"column:detailImages" json:"detailImages"` // JSON数组
    // ... 其他字段
}
```

### 2. 管理员服务信息结构更新
**文件**: `service/admin_service.go`
```go
type AdminServiceInfo struct {
    // ... 其他字段
    ImageUrl      string    `json:"imageUrl"`
    ImageCosId    string    `json:"imageCosId"`    // 腾讯云对象存储ID
    // ... 其他字段
}
```

### 3. 数据库迁移文件
**文件**: `db/migration/add_image_cos_id_to_service_items.sql`
```sql
-- 添加字段
ALTER TABLE ServiceItems 
ADD COLUMN imageCosId VARCHAR(255) COMMENT '腾讯云对象存储图片ID' AFTER imageUrl;

-- 添加索引
CREATE INDEX idx_service_items_image_cos_id ON ServiceItems(imageCosId);
```

### 4. 迁移执行脚本
**文件**: `scripts/add_image_cos_id.sh`
- 自动执行数据库迁移
- 验证字段添加成功
- 提供详细的执行日志

## 字段说明

### imageCosId 字段
- **类型**: VARCHAR(255)
- **位置**: 在 imageUrl 字段之后
- **用途**: 存储腾讯云对象存储的图片ID
- **索引**: 已添加索引 `idx_service_items_image_cos_id`
- **注释**: 腾讯云对象存储图片ID

## 使用方法

### 1. 执行数据库迁移
```bash
# 方法1: 使用脚本（推荐）
./scripts/add_image_cos_id.sh

# 方法2: 手动执行SQL
mysql -u用户名 -p密码 数据库名 < db/migration/add_image_cos_id_to_service_items.sql
```

### 2. 在代码中使用
```go
// 创建服务时设置 imageCosId
service := &model.ServiceItemModel{
    Name:       "服务名称",
    ImageUrl:   "https://example.com/image.jpg",
    ImageCosId: "cos_id_123456",  // 腾讯云对象存储ID
    // ... 其他字段
}

// 查询时使用 imageCosId
var services []*model.ServiceItemModel
db.Where("imageCosId = ?", cosId).Find(&services)
```

### 3. API 响应中包含 imageCosId
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "name": "服务名称",
        "imageUrl": "https://example.com/image.jpg",
        "imageCosId": "cos_id_123456",
        // ... 其他字段
      }
    ]
  }
}
```

## 注意事项

1. **向后兼容**: 现有代码不会受到影响，imageCosId 字段为可选
2. **数据迁移**: 现有记录的 imageCosId 字段为空，需要根据业务需求进行数据填充
3. **索引性能**: 已添加索引，查询性能良好
4. **字段长度**: VARCHAR(255) 足够存储腾讯云对象存储ID

## 相关文件
- `db/model/service.go` - 数据库模型
- `service/admin_service.go` - 管理员服务接口
- `db/migration/add_image_cos_id_to_service_items.sql` - 数据库迁移
- `scripts/add_image_cos_id.sh` - 迁移执行脚本

## 验证方法
1. 执行迁移脚本
2. 检查数据库表结构
3. 测试 API 接口返回
4. 验证索引创建成功
