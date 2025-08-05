# 订单导航功能

## 功能概述

在用户资料页面（"我的"tab）中，点击"我的订单"按钮可以跳转到订单列表页面，查看用户的所有订单。

## 功能特点

1. **用户登录检查**：跳转前检查用户是否已登录
2. **页面跳转**：使用 `wx.navigateTo` 跳转到订单列表页面
3. **错误处理**：提供详细的错误提示和调试信息
4. **状态保持**：订单列表页面保持用户的筛选状态

## 页面结构

### 1. 用户资料页面 (`/pages/user/profile`)
- **位置**：功能菜单中的"我的订单"项
- **触发**：点击"我的订单"按钮
- **功能**：跳转到订单列表页面

### 2. 订单列表页面 (`/pages/order/list`)
- **功能**：显示用户的所有订单
- **筛选**：支持按订单状态筛选
- **分页**：支持分页加载更多订单
- **操作**：支持订单详情、支付、取消等操作

## 技术实现

### 前端实现

#### 1. 用户资料页面
```javascript
// 订单列表跳转
onOrderList() {
  console.log('点击我的订单按钮')
  
  // 检查用户登录状态
  const userId = wx.getStorageSync('userId')
  if (!userId) {
    wx.showToast({
      title: '请先登录',
      icon: 'none'
    })
    return
  }
  
  console.log('准备跳转到订单列表页面')
  wx.navigateTo({
    url: '/pages/order/list',
    success: () => {
      console.log('跳转到订单列表页面成功')
    },
    fail: (error) => {
      console.error('跳转到订单列表页面失败:', error)
      wx.showToast({
        title: '页面跳转失败',
        icon: 'none'
      })
    }
  })
}
```

#### 2. 订单列表页面
```javascript
// 加载订单列表
async loadOrders(isRefresh = false) {
  try {
    // 获取用户ID
    const userId = wx.getStorageSync('userId')
    if (!userId) {
      console.log('用户未登录，跳转到登录页面')
      wx.navigateTo({
        url: '/pages/login/login'
      })
      return
    }
    
    const params = {
      page: this.data.page,
      pageSize: this.data.pageSize,
      status: this.data.status,
      userId: userId
    }
    
    const result = await api.orderList(params)
    
    if (result.code === 0 && result.data) {
      const { list, total, hasMore } = result.data
      const orders = list || []
      
      // 处理订单数据
      const processedOrders = orders.map(order => ({
        ...order,
        statusText: this.getStatusText(order.status),
        statusClass: this.getStatusClass(order.status),
        formattedAmount: this.formatAmount(order.totalAmount),
        formattedCreatedAt: this.formatTime(order.createdAt)
      }))
      
      this.setData({ 
        orders: isRefresh ? processedOrders : [...this.data.orders, ...processedOrders],
        hasMore: hasMore,
        loading: false 
      })
    }
  } catch (error) {
    console.error('加载订单列表失败:', error)
    this.setData({ loading: false })
  }
}
```

### 后端实现

#### 1. API路由
```go
// main.go
http.HandleFunc("/api/order/list", service.NewLogMiddleware(service.OrderListHandler))
```

#### 2. 订单列表处理
```go
// OrderListHandler 获取订单列表接口
func OrderListHandler(w http.ResponseWriter, r *http.Request) {
    // 解析请求参数
    // 验证用户身份
    // 查询订单列表
    // 返回订单数据
}
```

## 用户体验流程

### 1. 进入用户资料页面
- 用户点击底部tab的"我的"
- 显示用户资料和功能菜单

### 2. 点击我的订单
- 用户点击"我的订单"按钮
- 系统检查用户登录状态
- 跳转到订单列表页面

### 3. 查看订单列表
- 显示用户的所有订单
- 支持按状态筛选订单
- 支持分页加载更多

### 4. 订单操作
- 点击订单查看详情
- 支持支付、取消、退款等操作

## 错误处理

### 1. 用户未登录
- 显示"请先登录"提示
- 阻止页面跳转

### 2. 页面跳转失败
- 显示"页面跳转失败"提示
- 记录详细错误信息

### 3. API调用失败
- 显示加载失败提示
- 提供重试机制

## 测试验证

### 1. 功能测试
```javascript
// 在用户资料页面控制台运行
const { runOrderTests } = require('./tests/test_order_navigation.js')
runOrderTests()
```

### 2. 页面跳转测试
```javascript
// 测试页面跳转
wx.navigateTo({
  url: '/pages/order/list',
  success: () => console.log('跳转成功'),
  fail: (error) => console.error('跳转失败:', error)
})
```

### 3. API测试
```javascript
// 测试订单列表API
const { api } = require('../utils/cloud-container-standard')
const result = await api.orderList({
  userId: wx.getStorageSync('userId'),
  page: 1,
  pageSize: 10,
  status: ''
})
console.log('API结果:', result)
```

## 配置检查

### 1. 页面注册
- ✅ `pages/order/list` 已在 `app.json` 中注册
- ✅ 订单列表页面文件存在且完整

### 2. TabBar配置
- ✅ 订单页面已在 `app.json` 的 tabBar 中配置
- ✅ 用户可以通过底部tab直接访问订单页面

### 3. 后端API
- ✅ `/api/order/list` 路由已配置
- ✅ `OrderListHandler` 处理函数已实现

## 注意事项

1. **用户登录状态**：确保用户已登录才能查看订单
2. **页面权限**：订单列表页面需要用户身份验证
3. **数据安全**：只显示当前用户的订单
4. **性能优化**：使用分页加载，避免一次性加载大量数据

## 后续优化

1. **订单状态实时更新**：使用WebSocket或轮询更新订单状态
2. **订单搜索功能**：支持按订单号、服务名称搜索
3. **订单导出功能**：支持导出订单数据
4. **订单统计功能**：显示订单数量、金额统计
5. **订单提醒功能**：订单状态变更时发送通知 