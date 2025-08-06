# 微信支付功能部署和配置指南

## 功能概述

已实现完整的微信支付功能，包括：
- 微信支付参数生成
- 支付结果通知处理
- 订单状态更新
- 支付配置管理

## 系统架构

### 后端组件
1. **支付配置管理** (`config/payment_config.go`)
2. **微信支付服务** (`service/wechat_pay_service.go`)
3. **订单支付处理** (`service/order_service.go`)
4. **支付通知处理** (`service/wechat_pay_service.go`)
请输入新金额：
### 前端组件
1. **支付调用** (`miniprogram/pages/order/order.js`)
2. **支付确认** (`miniprogram/pages/order/detail.js`)

## 配置要求

### 1. 微信小程序配置
- **AppID**: 微信小程序的AppID
- **AppSecret**: 微信小程序的AppSecret

### 2. 微信支付配置
- **商户号 (MchID)**: 微信支付商户号
- **商户密钥 (MchKey)**: 微信支付商户密钥
- **通知地址 (NotifyURL)**: 支付结果通知地址
- **环境设置**: 沙箱环境或生产环境

### 3. 环境变量配置

```bash
# 微信小程序配置
export WX_APP_ID="你的微信小程序AppID"
export WX_APP_SECRET="你的微信小程序AppSecret"

# 微信支付配置
export WECHAT_PAY_APP_ID="你的微信小程序AppID"
export WECHAT_PAY_MCH_ID="你的微信支付商户号"
export WECHAT_PAY_MCH_KEY="你的微信支付商户密钥"
export WECHAT_PAY_NOTIFY_URL="https://your-domain.com/api/payment/notify"
export WECHAT_PAY_ENVIRONMENT="sandbox"  # 或 "production"
```

## 部署步骤

### 1. 环境准备

#### 获取微信支付配置
1. 登录微信商户平台：https://pay.weixin.qq.com/
2. 获取商户号 (MchID)
3. 设置商户密钥 (MchKey)
4. 配置支付通知地址

#### 获取微信小程序配置
1. 登录微信公众平台：https://mp.weixin.qq.com/
2. 选择你的小程序
3. 进入"开发" -> "开发管理" -> "开发设置"
4. 复制AppID和AppSecret

### 2. 后端部署

#### 设置环境变量
```bash
# 在服务器上设置环境变量
export WX_APP_ID="wx101090677bd5219e"
export WX_APP_SECRET="042ff9921818ada9336df6e91fc2287e"
export WECHAT_PAY_APP_ID="wx101090677bd5219e"
export WECHAT_PAY_MCH_ID="你的商户号"
export WECHAT_PAY_MCH_KEY="你的商户密钥"
export WECHAT_PAY_NOTIFY_URL="https://your-domain.com/api/payment/notify"
export WECHAT_PAY_ENVIRONMENT="sandbox"
```

#### 编译和启动
```bash
# 编译后端服务
go build -o main main.go

# 启动服务
./main
```

### 3. 前端配置

#### 更新支付调用
确保前端支付调用包含正确的参数：

```javascript
// 支付订单
async payOrder(orderId) {
  try {
    const userInfo = wx.getStorageSync('userInfo')
    const paymentData = {
      orderId,
      paymentMethod: 'wechat_pay',
      openId: userInfo.openId  // 确保传递openId
    }

    const res = await this.payOrderToServer(paymentData)
    
    if (res.code === 0) {
      const paymentParams = res.data.paymentParams
      
      // 调用微信支付
      wx.requestPayment({
        timeStamp: paymentParams.timeStamp,
        nonceStr: paymentParams.nonceStr,
        package: paymentParams.package,
        signType: paymentParams.signType,
        paySign: paymentParams.paySign,
        success: () => {
          wx.showToast({
            title: '支付成功',
            icon: 'success'
          })
          
          // 跳转到订单详情页
          setTimeout(() => {
            wx.navigateTo({
              url: `/pages/order/detail?id=${orderId}`
            })
          }, 1500)
        },
        fail: (err) => {
          console.error('支付失败:', err)
          wx.showToast({
            title: '支付失败',
            icon: 'none'
          })
        }
      })
    }
  } catch (error) {
    console.error('支付失败:', error)
    wx.showToast({
      title: '支付失败，请重试',
      icon: 'none'
    })
  }
}
```

## 测试验证

### 1. 配置测试
```bash
# 运行配置测试
./tests/backend/test_payment_config.sh
```

### 2. 支付流程测试
```bash
# 运行支付流程测试
./tests/backend/test_payment_flow.sh
```

### 3. 手动测试步骤
1. 创建测试订单
2. 调用支付接口
3. 验证支付参数
4. 测试支付通知

## 支付流程

### 1. 发起支付
```
用户点击支付 -> 前端调用支付接口 -> 后端生成支付参数 -> 返回给前端
```

### 2. 微信支付
```
前端调用wx.requestPayment -> 微信支付界面 -> 用户完成支付 -> 微信返回结果
```

### 3. 支付确认
```
微信发送通知 -> 后端接收通知 -> 验证签名 -> 更新订单状态 -> 返回确认
```

## 错误处理

### 常见错误及解决方案

#### 1. 配置错误
- **错误**: "微信支付配置不完整"
- **解决**: 检查环境变量是否正确设置

#### 2. 签名错误
- **错误**: "签名验证失败"
- **解决**: 检查商户密钥是否正确

#### 3. 网络错误
- **错误**: "HTTP请求失败"
- **解决**: 检查网络连接和防火墙设置

#### 4. 订单状态错误
- **错误**: "订单状态不正确"
- **解决**: 确保订单状态为待支付

## 监控和日志

### 日志记录
支付相关的所有操作都会记录详细日志：
- 支付参数生成
- 微信API调用
- 支付通知处理
- 错误信息

### 监控指标
- 支付成功率
- 支付响应时间
- 错误率统计

## 安全考虑

### 1. 签名验证
- 所有支付请求都进行签名验证
- 支付通知必须验证签名

### 2. 参数验证
- 验证订单金额
- 验证订单状态
- 验证用户权限

### 3. 环境隔离
- 开发环境使用沙箱
- 生产环境使用正式配置

## 故障排除

### 1. 支付调起失败
- 检查微信支付配置
- 验证支付参数完整性
- 查看后端日志

### 2. 支付通知失败
- 检查通知地址配置
- 验证签名算法
- 查看网络连接

### 3. 订单状态未更新
- 检查支付通知处理
- 验证数据库连接
- 查看错误日志

## 联系支持

如果遇到问题，请提供：
1. 错误日志
2. 支付配置信息（隐藏敏感信息）
3. 测试步骤
4. 环境信息 