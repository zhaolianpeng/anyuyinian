# ServiceId功能实现文档

## 功能概述

在首页点击热门服务时，跳转到对应的服务详情页并带上serviceId参数，请求具体的服务内容。

## 实现内容

### 1. 后端修改

#### 1.1 修改服务数据转换函数
- 文件：`anyuyinian/service/home_init_service.go`
- 修改：`convertServicesToInterface` 函数
- 添加：`serviceId` 字段，值为服务的 `id`

```go
// convertServicesToInterface 转换服务项数据
func convertServicesToInterface(services []*model.ServiceModel) []interface{} {
	result := make([]interface{}, len(services))
	for i, service := range services {
		result[i] = map[string]interface{}{
			"id":          service.Id,
			"serviceId":   service.Id, // 添加serviceId字段，用于前端跳转
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
- 功能：优先使用 `serviceId`，如果没有则使用 `id`

```javascript
/**
 * 服务点击
 */
onServiceTap(e) {
  const { id, serviceId } = e.currentTarget.dataset
  // 优先使用serviceId，如果没有则使用id
  const targetId = serviceId || id
  wx.navigateTo({
    url: `/pages/service/detail?id=${targetId}`
  })
},
```

#### 2.2 修改服务列表WXML
- 文件：`miniprogram/pages/index/index.wxml`
- 修改：服务项的数据绑定
- 添加：`data-service-id="{{item.serviceId}}"` 属性

```xml
<!-- 服务列表 -->
<view wx:if="{{services.length > 0}}" class="service-section">
  <view class="section-title">热门服务</view>
  <view class="service-grid">
    <view wx:for="{{services}}" wx:key="id" class="service-item"
          bindtap="onServiceTap" data-id="{{item.id}}" data-service-id="{{item.serviceId}}">
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

1. **后端返回数据**：
   - 首页init接口返回的services数组中，每个服务项包含 `serviceId` 字段
   - `serviceId` 的值与服务项的 `id` 相同

2. **前端处理**：
   - 首页接收到服务数据，包含 `serviceId` 字段
   - 用户点击服务项时，传递 `serviceId` 参数
   - 跳转到服务详情页：`/pages/service/detail?id=${serviceId}`

3. **服务详情页**：
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

- 向后兼容：如果 `serviceId` 不存在，会使用 `id` 字段
- 数据一致性：`serviceId` 和 `id` 的值相同，确保数据一致性
- 错误处理：如果服务详情接口调用失败，会显示相应的错误提示

## 注意事项

1. 确保后端服务正常运行
2. 确保数据库中的服务数据完整
3. 前端需要重新编译和部署
4. 测试时注意检查网络连接状态 