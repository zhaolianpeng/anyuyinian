# ServiceId功能实现文档

## 功能概述

在首页点击热门服务时，跳转到对应的服务详情页并带上serviceId参数，请求具体的服务内容。

## 实现内容

### 1. 后端修改

#### 1.1 数据库迁移
- 文件：`anyuyinian/db/migration/add_service_id_field.sql`
- 功能：为Services表添加serviceId字段

```sql
-- 为Services表添加serviceId字段
ALTER TABLE Services ADD COLUMN serviceId INT COMMENT '服务ID，用于前端跳转' AFTER id;

-- 更新现有数据的serviceId字段，使其与id字段值相同
UPDATE Services SET serviceId = id WHERE serviceId IS NULL;

-- 为serviceId字段添加索引
ALTER TABLE Services ADD INDEX idx_service_id (serviceId);
```

#### 1.2 更新服务模型
- 文件：`anyuyinian/db/model/home.go`
- 修改：`ServiceModel` 结构体
- 添加：`ServiceId` 字段

```go
// ServiceModel 服务项模型
type ServiceModel struct {
	Id          int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ServiceId   int32     `gorm:"column:serviceId" json:"serviceId"` // 服务ID，用于前端跳转
	Name        string    `gorm:"column:name;not null" json:"name"`
	Description string    `gorm:"column:description" json:"description"`
	Icon        string    `gorm:"column:icon;not null" json:"icon"`
	ImageUrl    string    `gorm:"column:imageUrl" json:"imageUrl"`
	LinkUrl     string    `gorm:"column:linkUrl" json:"linkUrl"`
	Sort        int       `gorm:"column:sort;default:0" json:"sort"`
	Status      int       `gorm:"column:status;default:1" json:"status"` // 1-启用，0-禁用
	CreatedAt   time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}
```

#### 1.3 修改服务数据转换函数
- 文件：`anyuyinian/service/home_init_service.go`
- 修改：`convertServicesToInterface` 函数
- 使用：数据库中的 `serviceId` 字段

```go
// convertServicesToInterface 转换服务项数据
func convertServicesToInterface(services []*model.ServiceModel) []interface{} {
	result := make([]interface{}, len(services))
	for i, service := range services {
		result[i] = map[string]interface{}{
			"id":          service.Id,
			"serviceId":   service.ServiceId, // 使用数据库中的serviceId字段
			"name":        service.Name,
			"description": service.Description,
			"icon":        service.Icon,
			"imageUrl":    service.ImageUrl,
			"linkUrl":     service.LinkUrl,
			"sort":        service.Sort,
		}
	}
	return result
}
```

### 2. 前端修改

#### 2.1 修改服务点击处理逻辑
- 文件：`miniprogram/pages/index/index.js`
- 修改：`onServiceTap` 函数
- 功能：使用 `serviceId` 进行跳转

```javascript
/**
 * 服务点击
 */
onServiceTap(e) {
  const { serviceId } = e.currentTarget.dataset
  // 使用serviceId进行跳转
  if (serviceId) {
    wx.navigateTo({
      url: `/pages/service/detail?id=${serviceId}`
    })
  } else {
    console.error('serviceId不存在，无法跳转到服务详情页')
    wx.showToast({
      title: '服务信息不完整',
      icon: 'none'
    })
  }
},
```

#### 2.2 修改服务列表WXML
- 文件：`miniprogram/pages/index/index.wxml`
- 修改：服务项的数据绑定
- 使用：`data-service-id="{{item.serviceId}}"` 属性

```xml
<!-- 服务列表 -->
<view wx:if="{{services.length > 0}}" class="service-section">
  <view class="section-title">热门服务</view>
  <view class="service-grid">
    <view wx:for="{{services}}" wx:key="id" class="service-item"
          bindtap="onServiceTap" data-service-id="{{item.serviceId}}">
      <image src="{{item.imageUrl}}" class="service-bg" mode="aspectFill" />
      <view class="service-content">
        <view class="service-info">
          <text class="service-name">{{item.name}}</text>
          <text class="service-desc">{{item.description}}</text>
        </view>
      </view>
    </view>
  </view>
</view>
```

## 数据流程

1. **数据库存储**：
   - Services表包含 `serviceId` 字段，类型为INT
   - `serviceId` 字段用于前端跳转，与 `id` 字段分离

2. **后端返回数据**：
   - 首页init接口返回的services数组中，每个服务项包含 `serviceId` 字段
   - `serviceId` 的值来自数据库Services表的serviceId字段

3. **前端处理**：
   - 首页接收到服务数据，包含 `serviceId` 字段
   - 用户点击服务项时，传递 `serviceId` 参数
   - 跳转到服务详情页：`/pages/service/detail?id=${serviceId}`

4. **服务详情页**：
   - 接收 `id` 参数（实际是serviceId）
   - 调用服务详情接口获取具体服务内容
   - 显示服务详情信息

## 测试验证

### 后端测试
```bash
# 运行后端测试脚本
./tests/backend/test_service_id_feature.sh
```

### 前端测试
```javascript
// 在小程序开发工具中运行
const { testServiceClickFeature } = require('./tests/test_service_click.js')
testServiceClickFeature()
```

## 兼容性说明

- 严格使用serviceId：只使用serviceId进行跳转，不使用id字段
- 数据一致性：`serviceId` 和 `id` 的值相同，确保数据一致性
- 错误处理：如果serviceId不存在，会显示错误提示并阻止跳转

## 注意事项

1. 确保后端服务正常运行
2. 确保数据库中的服务数据完整
3. 前端需要重新编译和部署
4. 测试时注意检查网络连接状态 