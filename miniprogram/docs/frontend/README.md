# 安语颐年小程序

## 项目结构

```
miniprogram/
├── app.js                 # 小程序入口文件
├── app.json              # 小程序配置文件
├── app.wxss              # 全局样式文件
├── config.js             # 配置文件
├── project.config.json   # 项目配置文件
├── sitemap.json          # 站点地图配置
├── components/           # 组件目录
├── pages/               # 页面目录
├── images/              # 图片资源目录
├── utils/               # 工具函数目录
├── docs/                # 文档目录
└── tests/               # 测试文件目录
```

## 目录说明

### 📁 docs/ - 文档目录
包含所有项目相关的文档文件，按功能分类：

#### 功能文档
- `AGREEMENT_PAGES_SUMMARY.md` - 协议页面功能总结
- `ALL_API_FIX_SUMMARY.md` - API修复总结
- `API_FIX_SUMMARY.md` - API修复文档
- `API_METHOD_FIX.md` - API方法修复
- `API_PATH_FIX.md` - API路径修复

#### 登录相关文档
- `LOGIN_BACKEND_FIX.md` - 登录后端修复
- `LOGIN_FIX_DESC.md` - 登录修复描述
- `LOGIN_FIX_DIRECT.md` - 登录直接修复
- `LOGIN_FIX_FINAL.md` - 登录最终修复
- `LOGIN_FIX_SIMPLE.md` - 登录简单修复
- `LOGIN_FIX_SYNC.md` - 登录同步修复
- `LOGIN_FIX.md` - 登录修复
- `LOGIN_JUMP_FIX.md` - 登录跳转修复
- `LOGIN_SUCCESS_SUMMARY.md` - 登录成功总结

#### 订单相关文档
- `ORDER_BOOKING_FEATURES.md` - 订单预约功能
- `ORDER_BUTTON_FIX.md` - 订单按钮修复
- `ORDER_PAGE_FIELDS_DEBUG.md` - 订单页面字段调试
- `ORDER_PAGE_FIELDS_SUMMARY.md` - 订单页面字段总结
- `ORDER_PAGE_FIX.md` - 订单页面修复
- `ORDER_PAGE_IMPLEMENTATION.md` - 订单页面实现
- `ORDER_PAGE_SCROLL_FIX.md` - 订单页面滚动修复

#### 服务相关文档
- `SERVICE_DATA_UPDATE.md` - 服务数据更新
- `SERVICE_ID_ERROR_FIX.md` - 服务ID错误修复
- `SERVICE_LIST_IMAGE_FIX.md` - 服务列表图片修复
- `SERVICE_NAVIGATION_FIX.md` - 服务导航修复

#### 用户相关文档
- `USER_CREATION_FIX.md` - 用户创建修复
- `USER_INFO_FIX.md` - 用户信息修复
- `USER_PROFILE_FEATURES.md` - 用户资料功能

#### 其他功能文档
- `HOME_DATA_FIX.md` - 首页数据修复
- `IMAGE_ERROR_FIX.md` - 图片错误修复
- `DATA_DISPLAY_FIX.md` - 数据显示修复
- `FEATURE_IMPLEMENTATION_SUMMARY.md` - 功能实现总结

#### 配置相关文档
- `COS_DOMAIN_SETUP.md` - COS域名设置
- `COS_IMAGE_FIX.md` - COS图片修复
- `COS_SETUP.md` - COS设置
- `COS_VERIFICATION.md` - COS验证
- `CLOUD_URL_CONVERSION.md` - 云URL转换
- `check_domain_config.md` - 域名配置检查

#### 开发文档
- `CHANGELOG.md` - 更新日志
- `DEVELOPMENT.md` - 开发文档
- `FINAL_SOLUTION.md` - 最终解决方案

### 📁 tests/ - 测试文件目录
包含所有测试相关的文件：

#### API测试
- `api_test_example.js` - API测试示例
- `test_api_fix.js` - API修复测试
- `test_backend.js` - 后端测试
- `test_cloud_url.js` - 云URL测试
- `test_cos.js` - COS测试

#### 功能测试
- `test_agreement_pages.js` - 协议页面测试
- `test_data_display.js` - 数据显示测试
- `test_fixes.js` - 修复测试
- `test_home_api.js` - 首页API测试
- `test_home_data.js` - 首页数据测试
- `test_image_processing.js` - 图片处理测试
- `test_json_number_fix.js` - JSON数字修复测试
- `test_order_debug.js` - 订单调试测试
- `test_order_fields.js` - 订单字段测试
- `test_order_fix.js` - 订单修复测试
- `test_order_flow.js` - 订单流程测试
- `test_order_page.js` - 订单页面测试
- `test_scroll_fix.js` - 滚动修复测试
- `test_service_data.js` - 服务数据测试
- `test_service_list_images.js` - 服务列表图片测试
- `test_service_navigation.js` - 服务导航测试
- `test_service_post.js` - 服务POST测试
- `test_updated_service_data.js` - 更新服务数据测试
- `test_user_status.js` - 用户状态测试
- `test_wxss_fix.js` - WXSS修复测试

#### 登录测试
- `test_login_debug.js` - 登录调试测试
- `test_login_flow.js` - 登录流程测试
- `test_login_response.js` - 登录响应测试
- `test_login_simple.js` - 简单登录测试

#### 工具测试
- `create_placeholder_images.js` - 创建占位图片
- `debug_images.js` - 图片调试
- `fix_api_paths.js` - 修复API路径
- `quick_image_fix.js` - 快速图片修复
- `quick_test_service.js` - 快速服务测试

## 使用说明

### 查看文档
所有项目文档都整理在 `docs/` 目录中，按功能分类便于查找。

### 运行测试
所有测试文件都整理在 `tests/` 目录中，可以根据需要运行相应的测试。

### 开发流程
1. 查看相关文档了解功能实现
2. 运行对应测试验证功能
3. 根据测试结果进行调试和修复

## 注意事项

- 文档按功能分类，便于查找和维护
- 测试文件按功能分类，便于运行和调试
- 保持文档和测试文件的同步更新
- 新增功能时请同时更新相关文档和测试文件 