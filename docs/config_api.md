# 平台配置接口文档

## 接口信息

- **接口地址**: `GET /api/config`
- **请求方式**: GET
- **功能**: 获取平台配置信息（客服电话、协议链接等）

## 请求参数

无

## 响应格式

### 成功响应

```json
{
  "code": 0,
  "data": {
    "customer_service_phone": "400-123-4567",
    "customer_service_qq": "123456789",
    "customer_service_wechat": "kefu123",
    "privacy_policy_url": "https://example.com/privacy.html",
    "user_agreement_url": "https://example.com/agreement.html",
    "about_us_url": "https://example.com/about.html",
    "app_version": "1.0.0",
    "force_update": false,
    "update_url": "https://example.com/download",
    "maintenance_mode": false,
    "maintenance_message": "",
    "wechat_appid": "wx1234567890abcdef",
    "wechat_secret": "abcdef1234567890",
    "payment_config": {
      "wechat_pay": {
        "mch_id": "1234567890",
        "notify_url": "https://example.com/api/payment/notify"
      }
    }
  }
}
```

### 失败响应

```json
{
  "code": -1,
  "errorMsg": "错误信息"
}
```

## 配置项说明

| 配置项 | 类型 | 说明 |
|--------|------|------|
| customer_service_phone | string | 客服电话 |
| customer_service_qq | string | 客服QQ |
| customer_service_wechat | string | 客服微信 |
| privacy_policy_url | string | 隐私政策链接 |
| user_agreement_url | string | 用户协议链接 |
| about_us_url | string | 关于我们链接 |
| app_version | string | 应用版本号 |
| force_update | boolean | 是否强制更新 |
| update_url | string | 更新下载链接 |
| maintenance_mode | boolean | 是否维护模式 |
| maintenance_message | string | 维护提示信息 |
| wechat_appid | string | 微信AppID |
| wechat_secret | string | 微信AppSecret |
| payment_config | object | 支付配置 |

## 使用示例

### 微信小程序端

```javascript
// 获取平台配置
wx.request({
  url: 'http://your-server.com/api/config',
  method: 'GET',
  success: (res) => {
    if (res.data.code === 0) {
      const config = res.data.data;
      
      // 检查是否需要更新
      if (config.force_update) {
        wx.showModal({
          title: '版本更新',
          content: '发现新版本，请立即更新',
          showCancel: false,
          success: () => {
            wx.navigateTo({
              url: config.update_url
            });
          }
        });
      }
      
      // 检查是否维护模式
      if (config.maintenance_mode) {
        wx.showModal({
          title: '系统维护',
          content: config.maintenance_message,
          showCancel: false
        });
      }
      
      // 保存配置到本地
      wx.setStorageSync('appConfig', config);
    }
  }
});
```

## 数据库表结构

### 配置表 (Configs)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | INT | 主键，自增 |
| key | VARCHAR(100) | 配置键，唯一索引 |
| value | TEXT | 配置值 |
| description | VARCHAR(500) | 配置描述 |
| type | VARCHAR(20) | 配置类型：string, number, boolean, json |
| status | INT | 状态：1-启用，0-禁用 |
| createdAt | DATETIME | 创建时间 |
| updatedAt | DATETIME | 更新时间 |

## 默认配置数据

```sql
INSERT INTO Configs (key, value, description, type, status) VALUES
('customer_service_phone', '400-123-4567', '客服电话', 'string', 1),
('customer_service_qq', '123456789', '客服QQ', 'string', 1),
('customer_service_wechat', 'kefu123', '客服微信', 'string', 1),
('privacy_policy_url', 'https://example.com/privacy.html', '隐私政策链接', 'string', 1),
('user_agreement_url', 'https://example.com/agreement.html', '用户协议链接', 'string', 1),
('about_us_url', 'https://example.com/about.html', '关于我们链接', 'string', 1),
('app_version', '1.0.0', '应用版本号', 'string', 1),
('force_update', 'false', '是否强制更新', 'boolean', 1),
('update_url', 'https://example.com/download', '更新下载链接', 'string', 1),
('maintenance_mode', 'false', '是否维护模式', 'boolean', 1),
('maintenance_message', '', '维护提示信息', 'string', 1),
('wechat_appid', 'wx1234567890abcdef', '微信AppID', 'string', 1),
('wechat_secret', 'abcdef1234567890', '微信AppSecret', 'string', 1),
('payment_config', '{"wechat_pay":{"mch_id":"1234567890","notify_url":"https://example.com/api/payment/notify"}}', '支付配置', 'json', 1);
```

## 注意事项

1. **配置缓存**: 建议前端缓存配置信息，减少请求频率
2. **敏感信息**: 某些配置项（如AppSecret）不应返回给前端
3. **配置类型**: 支持string、number、boolean、json四种类型
4. **状态控制**: 通过status字段控制配置是否生效
5. **版本管理**: 支持应用版本控制和强制更新
6. **维护模式**: 支持系统维护模式，可显示维护信息
7. **支付配置**: 包含微信支付等第三方支付配置
8. **协议链接**: 包含隐私政策、用户协议等重要链接 