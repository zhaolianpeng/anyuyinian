# 微信支付配置指南

## 🎯 概述

本指南将帮助您完成微信支付的完整配置，包括：
- 微信商户平台配置
- 小程序后台配置
- 后端API实现
- 前端支付集成

## 🔧 配置步骤

### 1. 微信商户平台配置

#### 1.1 登录微信商户平台
- 访问：https://pay.weixin.qq.com/
- 使用您的微信账号登录
- 选择商户号：**1726638701**

#### 1.2 配置支付回调地址
1. 进入 **产品中心** → **开发配置** → **支付配置**
2. 设置回调地址：`https://golang-lfwy-176496-6-1353115175.sh.run.tcloudbase.com/api/payment/notify`
3. 点击 **保存**

#### 1.3 配置API密钥
1. 进入 **账户中心** → **API安全** → **API密钥**
2. 设置密钥：`JQzOCB8doIdgaUjAobELsk9nTyxdKhat`
3. 点击 **确认**

### 2. 小程序后台配置

#### 2.1 登录微信公众平台
- 访问：https://mp.weixin.qq.com/
- 登录您的小程序账号
- 小程序AppID：**wx101090677bd5219e**

#### 2.2 配置服务器域名
1. 进入 **开发** → **开发管理** → **服务器域名**
2. 添加request合法域名：`https://golang-lfwy-176496-6-1353115175.sh.run.tcloudbase.com`
3. 点击 **保存**

### 3. 后端API实现

#### 3.1 支付相关API
- ✅ **统一下单**：`POST /api/order/pay/{orderId}`
- ✅ **支付回调**：`POST /api/payment/notify`
- ✅ **查询订单**：`GET /api/order/detail?orderNo={orderNo}`
- ✅ **申请退款**：`POST /api/order/refund/{orderId}`

#### 3.2 支付配置
```go
// 支付配置
WechatPayConfig{
    AppID:       "wx101090677bd5219e",
    MchID:       "1726638701",
    MchKey:      "JQzOCB8doIdgaUjAobELsk9nTyxdKhat",
    NotifyURL:   "https://golang-lfwy-176496-6-1353115175.sh.run.tcloudbase.com/api/payment/notify",
    Environment: "production",
}
```

### 4. 前端支付集成

#### 4.1 小程序端调用
```javascript
// 统一下单
const result = await wx.request({
  url: 'https://golang-lfwy-176496-6-1353115175.sh.run.tcloudbase.com/api/order/pay/123',
  method: 'POST',
  data: {
    openId: 'user_openid',
    paymentMethod: 'wechat_pay'
  }
})

// 获取支付参数
const paymentParams = result.data.paymentParams

// 调起微信支付
wx.requestPayment({
  timeStamp: paymentParams.timeStamp,
  nonceStr: paymentParams.nonceStr,
  package: paymentParams.package,
  signType: paymentParams.signType,
  paySign: paymentParams.paySign,
  success: (res) => {
    console.log('支付成功:', res)
  },
  fail: (err) => {
    console.error('支付失败:', err)
  }
})
```

## 🧪 测试验证

### 1. 本地测试
```bash
# 启动后端服务
go run main.go

# 运行支付配置测试
./tests/wechat_pay/test_payment_config.sh

# 运行支付流程测试
./tests/wechat_pay/test_payment_flow.sh
```

### 2. 真机测试
1. 在微信开发者工具中预览
2. 在真机上测试支付功能
3. 验证支付结果通知

### 3. 生产环境测试
1. 部署到云托管
2. 配置域名白名单
3. 测试完整支付流程

## 📋 配置检查清单

### 微信商户平台
- [ ] 登录微信商户平台
- [ ] 配置支付回调地址
- [ ] 设置API密钥
- [ ] 验证商户号信息

### 小程序后台
- [ ] 登录微信公众平台
- [ ] 配置服务器域名
- [ ] 验证小程序AppID
- [ ] 检查权限设置

### 后端服务
- [ ] 支付API实现
- [ ] 支付配置正确
- [ ] 回调处理逻辑
- [ ] 错误处理机制

### 前端集成
- [ ] 支付页面集成
- [ ] 支付参数处理
- [ ] 支付结果处理
- [ ] 错误提示显示

## ⚠️ 注意事项

### 1. 域名配置
- 确保回调地址可公网访问
- 使用HTTPS协议
- 配置正确的域名白名单

### 2. 支付安全
- 验证微信支付签名
- 防止重复支付
- 记录详细日志

### 3. 错误处理
- 完善的错误处理机制
- 用户友好的错误提示
- 支付状态同步

### 4. 测试环境
- 使用沙箱环境测试
- 验证支付回调
- 测试异常情况

## 🚀 快速开始

### 1. 检查配置
```bash
# 运行配置验证脚本
./tests/wechat_pay/test_payment_config.sh
```

### 2. 测试支付
```bash
# 运行支付流程测试
./tests/wechat_pay/test_payment_flow.sh
```

### 3. 部署服务
```bash
# 部署到云托管
# 配置域名白名单
# 测试生产环境
```

## 📞 技术支持

如果遇到问题：
1. 检查配置是否正确
2. 查看服务日志
3. 验证网络连接
4. 测试API接口

## 🔗 相关链接

- [微信支付开发文档](https://pay.weixin.qq.com/wiki/doc/api/index.html)
- [小程序支付文档](https://developers.weixin.qq.com/miniprogram/dev/api/payment/wx.requestPayment.html)
- [云托管部署指南](https://cloud.tencent.com/document/product/876)