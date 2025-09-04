# 微信支付商户密钥设置指南

## 商户信息
- **商户号**: 1726638701
- **小程序AppID**: wx101090677bd5219e

## 获取商户密钥

### 方法一：通过微信商户平台（推荐）

#### 1. 登录商户平台
1. 访问 https://pay.weixin.qq.com/
2. 使用您的微信账号登录
3. 选择商户号 **1726638701**

#### 2. 设置API密钥
1. 进入 **账户中心** → **API安全** → **API密钥**
2. 点击 **设置API密钥** 或 **修改API密钥**
3. 输入32位字符串作为商户密钥
4. 确认设置

#### 3. 记录密钥
- 密钥是32位字符串
- 请妥善保管，不要泄露
- 如果忘记，可以重新设置

### 方法二：使用配置脚本

```bash
# 运行配置脚本
./scripts/set_merchant_key.sh
```

## 配置到代码中

### 1. 设置环境变量

```bash
# 设置商户密钥环境变量
export WECHAT_PAY_MCH_KEY="您的32位商户密钥"

# 验证设置
echo $WECHAT_PAY_MCH_KEY
```

### 2. 创建环境变量文件

```bash
# 创建 .env.wechat_pay 文件
cat > .env.wechat_pay << EOF
# 微信支付配置
WECHAT_PAY_APP_ID=wx101090677bd5219e
WECHAT_PAY_MCH_ID=1726638701
WECHAT_PAY_MCH_KEY=您的32位商户密钥
WECHAT_PAY_NOTIFY_URL=https://your-domain.com/api/payment/notify
WECHAT_PAY_ENVIRONMENT=production
EOF
```

### 3. 加载环境变量

```bash
# 加载环境变量
source .env.wechat_pay

# 或者重启应用
```

## 验证配置

### 1. 检查配置是否正确

```bash
# 运行配置检查脚本
./scripts/setup_wechat_pay.sh
```

### 2. 测试支付功能

```bash
# 运行支付测试
./tests/wechat_pay/test_payment.sh
```

## 常见问题

### Q1: 找不到API密钥设置入口
**解决方案**:
1. 确保已登录正确的商户号 1726638701
2. 检查账号权限，需要管理员权限
3. 尝试刷新页面或重新登录

### Q2: 密钥长度不是32位
**解决方案**:
1. 确保输入的是32位字符串
2. 可以使用在线随机字符串生成器
3. 建议使用字母和数字组合

### Q3: 设置后仍然报错
**解决方案**:
1. 确认环境变量已正确设置
2. 重启应用以加载新配置
3. 检查密钥是否包含特殊字符

## 安全注意事项

1. **保密性**: 商户密钥是敏感信息，不要提交到代码仓库
2. **权限控制**: 只有授权人员才能访问商户平台
3. **定期更换**: 建议定期更换商户密钥
4. **日志脱敏**: 日志中不要记录完整的商户密钥

## 示例配置

### 正确的环境变量设置
```bash
export WECHAT_PAY_APP_ID="wx101090677bd5219e"
export WECHAT_PAY_MCH_ID="1726638701"
export WECHAT_PAY_MCH_KEY="abcdef1234567890abcdef1234567890"
export WECHAT_PAY_NOTIFY_URL="https://your-domain.com/api/payment/notify"
export WECHAT_PAY_ENVIRONMENT="production"
```

### 验证配置的代码
```go
// 在代码中验证配置
paymentConfig := config.GetPaymentConfig()
if len(paymentConfig.WechatPay.MchKey) != 32 {
    log.Fatal("商户密钥长度错误，应为32位")
}
```

## 联系支持

如果遇到问题，请提供：
1. 商户号: 1726638701
2. 错误信息截图
3. 配置步骤描述
