# 重复支付问题修复

## 问题描述
用户尝试对已支付的订单进行重复支付，导致微信支付返回错误：`ORDERPAID - 该订单已支付`

## 错误信息
```
云托管调用成功: {
  statusCode: 200,
  data: {
    code: -1,
    data: null,
    errorMsg: "生成支付参数失败: 微信支付业务失败: ORDERPAID - 该订单已支付"
  }
}
```

## 问题原因
1. **前端缺少状态检查**: 没有在支付前检查订单状态
2. **后端错误信息不友好**: 直接返回微信支付原始错误信息
3. **页面状态不同步**: 订单状态更新后页面没有及时刷新

## 解决方案

### 1. 前端状态检查

#### 添加订单状态验证
**文件**: `miniprogram/pages/order/detail.js`

```javascript
async onPayOrder() {
  const { order } = this.data
  if (!order) return
  
  // 检查订单状态
  if (order.status !== 0) {
    let message = ''
    switch (order.status) {
      case 1:
        message = '订单已支付，无需重复支付'
        break
      case 2:
        message = '订单已完成'
        break
      case 3:
        message = '订单已取消'
        break
      case 4:
        message = '订单已退款'
        break
      default:
        message = '订单状态异常，无法支付'
    }
    
    wx.showToast({
      title: message,
      icon: 'none',
      duration: 2000
    })
    return
  }
  
  // 继续支付流程...
}
```

#### 添加页面刷新机制
```javascript
onShow() {
  // 页面显示时刷新订单详情
  const { orderNo } = this.data
  if (orderNo) {
    this.loadOrderDetail(orderNo)
  }
}
```

### 2. 后端错误处理优化

#### 友好错误信息
**文件**: `service/wechat_pay_service.go`

```go
if response.ResultCode != "SUCCESS" {
    LogError("微信支付业务失败", fmt.Errorf("result_code: %s, err_code: %s, err_code_des: %s", 
        response.ResultCode, response.ErrCode, response.ErrCodeDes))
    
    // 根据错误代码返回更友好的错误信息
    var friendlyError string
    switch response.ErrCode {
    case "ORDERPAID":
        friendlyError = "该订单已支付，无需重复支付"
    case "ORDERCLOSED":
        friendlyError = "订单已关闭，无法支付"
    case "SYSTEMERROR":
        friendlyError = "系统错误，请稍后重试"
    case "PARAM_ERROR":
        friendlyError = "参数错误，请检查订单信息"
    default:
        friendlyError = fmt.Sprintf("支付失败: %s", response.ErrCodeDes)
    }
    
    return nil, fmt.Errorf(friendlyError)
}
```

## 订单状态说明

### 状态码定义
- **0**: 待支付 - 可以支付
- **1**: 已支付 - 不能重复支付
- **2**: 已完成 - 订单完成
- **3**: 已取消 - 订单取消
- **4**: 已退款 - 订单退款

### 状态检查逻辑
1. **支付前检查**: 只有状态为0的订单才能支付
2. **状态提示**: 根据订单状态显示相应的提示信息
3. **页面刷新**: 页面显示时自动刷新订单状态

## 用户体验改进

### 1. 状态提示优化
- **已支付**: "订单已支付，无需重复支付"
- **已完成**: "订单已完成"
- **已取消**: "订单已取消"
- **已退款**: "订单已退款"

### 2. 页面刷新机制
- **onShow**: 页面显示时自动刷新订单状态
- **支付成功**: 支付完成后刷新订单详情
- **状态同步**: 确保页面状态与服务器状态一致

### 3. 错误处理
- **友好提示**: 将技术错误转换为用户友好的提示
- **状态引导**: 根据订单状态引导用户进行相应操作
- **避免重复**: 防止用户对已支付订单进行重复支付

## 测试场景

### 1. 正常支付流程
1. 创建待支付订单
2. 点击支付按钮
3. 完成支付
4. 验证订单状态更新

### 2. 重复支付测试
1. 对已支付订单点击支付
2. 验证显示"订单已支付，无需重复支付"
3. 验证不会调起微信支付

### 3. 状态同步测试
1. 在其他页面完成支付
2. 返回订单详情页面
3. 验证订单状态自动刷新

## 相关文件

- `miniprogram/pages/order/detail.js` - 订单详情页面
- `service/wechat_pay_service.go` - 微信支付服务
- `service/order_service.go` - 订单服务

## 注意事项

1. **状态一致性**: 确保前端状态与后端状态保持一致
2. **用户体验**: 提供清晰的状态提示和操作引导
3. **错误处理**: 将技术错误转换为用户友好的信息
4. **性能优化**: 避免不必要的重复请求
