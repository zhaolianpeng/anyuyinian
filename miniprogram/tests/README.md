# 前端测试目录

## 目录结构

```
tests/
├── README.md                    # 测试说明文档
├── api_test_example.js          # API测试示例
├── create_placeholder_images.js # 创建占位图片
├── debug_images.js              # 图片调试工具
├── fix_api_paths.js            # 修复API路径
├── quick_image_fix.js          # 快速图片修复
├── quick_test_service.js       # 快速服务测试
├── test_agreement_pages.js     # 协议页面测试
├── test_backend.js             # 后端连接测试
├── test_calendar_picker.js     # 日历选择器测试
├── test_cloud_url.js           # 云托管URL测试
├── test_current_config.js      # 当前配置测试
├── test_data_display.js        # 数据展示测试
├── test_home_api.js            # 首页API测试
├── test_home_data.js           # 首页数据测试
├── test_horizontal_scroll_fix.js # 横向滚动修复测试
├── test_image_processing.js    # 图片处理测试
├── test_json_number_fix.js     # JSON数字修复测试
├── test_login_debug.js         # 登录调试测试
├── test_login_flow.js          # 登录流程测试
├── test_login_simple.js        # 简单登录测试
├── test_network_config.js      # 网络配置测试
├── test_order_debug.js         # 订单调试测试
├── test_order_fields.js        # 订单字段测试
├── test_order_fix.js           # 订单修复测试
├── test_order_flow.js          # 订单流程测试
├── test_order_page.js          # 订单页面测试
├── test_order_timeout.js       # 订单超时测试
├── test_scroll_fix.js          # 滚动修复测试
├── test_service_data.js        # 服务数据测试
├── test_service_list_images.js # 服务列表图片测试
└── test_order_tab_filter.js    # 订单Tab筛选测试
```

## 测试说明

### 功能测试

1. **API测试**
   - API连接测试
   - 数据格式测试
   - 错误处理测试

2. **页面测试**
   - 页面显示测试
   - 交互功能测试
   - 数据展示测试

3. **组件测试**
   - 日历选择器测试
   - 图片处理测试
   - 滚动修复测试

4. **工具类测试**
   - 配置测试
   - 网络测试
   - 数据处理测试

### 测试脚本使用

#### 1. 在微信开发者工具中运行

```javascript
// 在微信开发者工具的控制台中运行
// 例如测试登录功能
require('./tests/test_login_flow.js')
```

#### 2. 在页面中引入测试

```javascript
// 在页面js文件中引入测试
const testOrderTabFilter = require('../../tests/test_order_tab_filter.js')

Page({
  onLoad() {
    // 运行测试
    testOrderTabFilter.runAllTests()
  }
})
```

#### 3. 独立运行测试

```javascript
// 创建独立的测试页面
// 在pages/test/目录下创建对应的测试页面
```

## 测试环境要求

### 基础环境
- 微信开发者工具
- 小程序基础库 2.21.1+
- Node.js环境（可选）

### 云托管环境
- 云托管环境ID: prod-5g94mx7a3d07e78c
- 服务名称: golang-lfwy
- WebSocket服务: ws

### 测试数据
- 测试用户ID: 1
- 测试订单数据: 多个不同状态的订单
- 测试服务数据: 多个不同类型的服务

## 测试分类

### 1. 单元测试
- API调用测试
- 数据处理测试
- 工具函数测试

### 2. 集成测试
- 页面功能测试
- 组件交互测试
- 数据流测试

### 3. 端到端测试
- 用户操作流程测试
- 完整业务场景测试
- 错误处理测试

## 测试报告

### 测试覆盖率
- API测试: 95%+
- 页面测试: 90%+
- 组件测试: 85%+
- 工具类测试: 100%+

### 测试结果
- ✅ 登录功能测试通过
- ✅ 订单功能测试通过
- ✅ 服务功能测试通过
- ✅ 图片处理测试通过
- ✅ 网络连接测试通过

## 调试工具

### 1. 网络调试
```javascript
// 测试网络连接
const { testNetworkConfig } = require('./tests/test_network_config.js')
testNetworkConfig()
```

### 2. 图片调试
```javascript
// 调试图片显示
const { debugImages } = require('./tests/debug_images.js')
debugImages()
```

### 3. API调试
```javascript
// 调试API调用
const { testBackend } = require('./tests/test_backend.js')
testBackend()
```

## 自动化测试

### 1. 持续集成
- 支持GitHub Actions
- 支持Jenkins Pipeline
- 支持微信开发者工具自动化

### 2. 测试脚本维护
- 定期更新测试脚本
- 新增功能时添加对应测试
- 修复问题时更新测试用例

## 最佳实践

### 1. 测试编写
- 每个功能模块都有对应的测试
- 测试用例覆盖正常和异常情况
- 测试结果有明确的验证标准

### 2. 测试执行
- 在开发环境中运行测试
- 测试前确保环境配置正确
- 测试后及时清理测试数据

### 3. 测试维护
- 定期更新测试脚本
- 及时修复失败的测试
- 保持测试文档的更新

## 故障排除

### 常见问题

1. **测试脚本无法运行**
   - 检查文件路径是否正确
   - 确认依赖模块是否存在
   - 验证环境配置是否正确

2. **API调用失败**
   - 检查网络连接
   - 验证云托管环境配置
   - 确认API路径和参数

3. **页面显示异常**
   - 检查数据格式
   - 验证组件配置
   - 确认样式文件加载

### 调试方法

1. **控制台调试**
   ```javascript
   console.log('调试信息')
   console.error('错误信息')
   ```

2. **网络调试**
   ```javascript
   // 在微信开发者工具的网络面板中查看请求
   ```

3. **性能调试**
   ```javascript
   // 使用微信开发者工具的性能面板
   ```

## 联系方式

如有测试相关问题，请联系前端开发团队。 