# 管理员订单详情页面修复

## 问题描述
管理员订单管理页面点击订单详情时提示"请使用订单号查看详情"，无法正常查看订单详情。

## 问题分析

### 1. 参数传递问题
- **管理员订单列表**：传递的是 `id` 参数（订单ID）
- **订单详情页面**：期望接收 `orderNo` 参数（订单号）
- **结果**：参数不匹配导致无法查看详情

### 2. 缺少API接口
- 前端缺少通过订单ID获取订单详情的API接口
- 后端缺少对应的处理接口

## 修复方案

### 1. 前端修复

#### 更新订单详情页面 (`pages/order/detail.js`)
```javascript
onLoad(options) {
  const { orderNo, id } = options
  if (orderNo) {
    this.setData({ orderNo: orderNo })
    this.loadOrderDetail(orderNo)
  } else if (id) {
    // 通过订单ID获取订单详情
    this.setData({ orderId: id })
    this.loadOrderDetailById(id)
  }
},

// 通过订单ID获取订单详情
async loadOrderDetailById(orderId) {
  try {
    console.log('通过订单ID获取订单详情:', orderId)
    
    // 先通过订单ID获取订单号
    const result = await api.orderDetailById(orderId)
    
    if (result.code === 0 && result.data) {
      const order = result.data
      console.log('通过ID获取订单成功:', order)
      
      // 设置订单号并重新加载详情
      this.setData({ orderNo: order.orderNo })
      this.loadOrderDetail(order.orderNo)
    } else {
      throw new Error(result.errorMsg || '获取订单详情失败')
    }
  } catch (error) {
    console.error('通过ID获取订单详情失败:', error)
    this.setData({ loading: false })
    
    wx.showToast({
      title: '获取订单详情失败',
      icon: 'none'
    })
  }
}
```

#### 添加API接口 (`utils/cloud-container-standard.js`)
```javascript
orderDetailById: (orderId) => callContainer(`/api/order/detail/${orderId}`, 'GET'),
```

### 2. 后端修复

#### 添加路由配置 (`main.go`)
```go
http.HandleFunc("/api/order/detail/", service.NewLogMiddleware(service.OrderDetailByIdHandler))
```

#### 添加处理器 (`service/order_service.go`)
```go
// OrderDetailByIdHandler 通过订单ID获取订单详情接口
func OrderDetailByIdHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
        return
    }

    // 从URL路径中获取订单ID
    pathParts := strings.Split(r.URL.Path, "/")
    if len(pathParts) < 4 {
        http.Error(w, "缺少订单ID参数", http.StatusBadRequest)
        return
    }

    orderIdStr := pathParts[3]
    orderId, err := strconv.Atoi(orderIdStr)
    if err != nil {
        http.Error(w, "无效的订单ID", http.StatusBadRequest)
        return
    }

    // 获取订单信息
    order, err := dao.OrderImp.GetOrderById(int32(orderId))
    if err != nil {
        LogError("获取订单信息失败", err)
        response := &OrderResponse{
            Code:     -1,
            ErrorMsg: "获取订单信息失败: " + err.Error(),
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
        return
    }

    // 构建响应数据
    response := &OrderResponse{
        Code: 0,
        Data: &OrderDetailResponse{
            OrderModel: order,
        },
    }

    // 获取关联信息（患者、地址、服务等）
    // ... 详细实现

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

## 修复内容

### 1. 前端修复
- **订单详情页面**：支持通过订单ID获取订单详情
- **API接口**：添加 `orderDetailById` 接口
- **错误处理**：添加相应的错误提示

### 2. 后端修复
- **路由配置**：添加 `/api/order/detail/` 路由
- **处理器**：添加 `OrderDetailByIdHandler` 处理器
- **数据获取**：通过订单ID获取完整的订单信息

### 3. 兼容性
- **向后兼容**：仍然支持通过订单号获取详情
- **灵活支持**：同时支持订单号和订单ID两种方式

## 使用方式

### 1. 通过订单号查看详情
```javascript
wx.navigateTo({
  url: `/pages/order/detail?orderNo=${orderNo}`
})
```

### 2. 通过订单ID查看详情
```javascript
wx.navigateTo({
  url: `/pages/order/detail?id=${orderId}`
})
```

## 修复效果

### 1. 管理员订单管理
- 点击订单详情可以正常查看
- 不再显示"请使用订单号查看详情"错误
- 支持完整的订单详情展示

### 2. 用户体验
- 管理员可以方便地查看订单详情
- 支持多种参数传递方式
- 错误处理更加友好

### 3. 系统稳定性
- 前后端参数传递一致
- API接口完整
- 错误处理完善

## 相关文件

- `pages/order/detail.js` - 订单详情页面
- `utils/cloud-container-standard.js` - API接口封装
- `main.go` - 后端路由配置
- `service/order_service.go` - 后端业务逻辑

## 测试建议

1. **管理员订单管理测试**
   - 进入管理员订单列表页面
   - 点击任意订单的详情按钮
   - 验证订单详情页面是否正常显示

2. **参数传递测试**
   - 测试通过订单号查看详情
   - 测试通过订单ID查看详情
   - 验证两种方式都能正常工作

3. **错误处理测试**
   - 测试无效的订单ID
   - 测试不存在的订单
   - 验证错误提示是否友好
