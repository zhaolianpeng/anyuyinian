# 订单API修改总结

## 修改概述

根据需求，对订单相关的API进行了以下修改：

1. **订单详情接口改为POST方法**
2. **订单列表接口显示更多信息**
3. **前端页面相应调整**

## 具体修改内容

### 1. 后端修改

#### 1.1 订单详情接口 (`/api/order/detail`)
- **修改前**: GET方法，通过路径参数获取订单ID
- **修改后**: POST方法，通过请求体传递订单号
- **文件**: `anyuyinian/main.go`, `anyuyinian/service/order_service.go`

```go
// 修改路由
http.HandleFunc("/api/order/detail", service.NewLogMiddleware(service.OrderDetailHandler))

// 请求结构
type OrderDetailRequest struct {
    OrderNo string `json:"orderNo"` // 订单号
}
```

#### 1.2 订单列表接口增强 (`/api/order/list`)
- **新增字段**: `ServiceTitle`, `Amount` (兼容字段)
- **显示信息**: 服务名称、预约时间、沟通时间、金额等
- **文件**: `anyuyinian/service/order_service.go`

```go
type OrderListItem struct {
    // ... 原有字段
    ServiceTitle    string    `json:"serviceTitle"`    // 服务标题
    Amount          float64   `json:"amount"`          // 兼容字段
}
```

### 2. 前端修改

#### 2.1 订单列表页面 (`miniprogram/pages/order/list.wxml`)
- **显示服务名称**: 使用 `serviceName` 字段
- **显示预约时间**: 格式化为 "预约时间：YYYY-MM-DD HH:MM"
- **显示沟通时间**: 从 `formData` 中提取的 `consultTime`
- **显示金额**: 使用 `totalAmount` 字段

```xml
<view class="service-details">
  <text class="appointment-info">预约时间：{{item.appointmentDate}} {{item.appointmentTime}}</text>
  <text class="consult-time" wx:if="{{item.consultTime}}">沟通时间：{{item.consultTime}}</text>
</view>
```

#### 2.2 订单列表样式 (`miniprogram/pages/order/list.wxss`)
- **新增样式**: `.service-details`, `.appointment-info`, `.consult-time`
- **布局优化**: 使用flex布局，垂直排列服务详情信息

#### 2.3 订单列表逻辑 (`miniprogram/pages/order/list.js`)
- **状态映射**: 修复状态显示逻辑，使用数字状态码
- **跳转逻辑**: 使用订单号而不是订单ID跳转到详情页
- **新增方法**: `goToService()` 跳转到服务页面

#### 2.4 订单详情页面 (`miniprogram/pages/order/detail.js`)
- **参数兼容**: 支持 `orderNo` 和 `id` 参数
- **错误处理**: 对旧版本ID参数给出友好提示

#### 2.5 请求工具 (`miniprogram/utils/request.js`)
- **订单详情**: 使用POST方法，传递订单号
```javascript
const getOrderDetail = (orderNo) => {
  return request('/api/order/detail', 'POST', { orderNo })
}
```

## 订单列表显示信息

现在订单列表会显示以下信息：

1. **服务名称** (`serviceName`)
2. **预约时间** (`appointmentDate` + `appointmentTime`)
3. **沟通时间** (`consultTime`) - 从表单数据中提取
4. **订单金额** (`totalAmount`)
5. **订单状态** (`statusText`)
6. **创建时间** (`formattedDate`)

## 测试验证

创建了测试脚本 `test_order_api_changes.sh` 来验证修改：

1. **测试订单详情接口**: POST请求，传递订单号
2. **测试订单列表接口**: GET请求，验证返回字段
3. **验证字段完整性**: 确保所有必要字段都正确返回

## 兼容性说明

- **向后兼容**: 订单详情页面支持旧的ID参数
- **数据兼容**: 订单列表保持原有字段，新增字段不影响现有功能
- **API兼容**: 订单列表接口保持GET方法不变

## 部署建议

1. **后端部署**: 重启Go服务以应用路由修改
2. **前端部署**: 重新编译小程序代码
3. **数据验证**: 确保数据库中的订单数据包含必要字段
4. **功能测试**: 测试订单列表和详情页面的完整流程 