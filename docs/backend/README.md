# 后端文档目录

本目录包含后端相关的技术文档，按功能模块分类组织。

## 目录结构

### admin/ - 管理员功能文档
- `ADMIN_DEPLOYMENT_GUIDE.md` - 管理员功能部署指南
- `ADMIN_FEATURE_SUMMARY.md` - 管理员功能总结
- `ADMIN_FINAL_DEPLOYMENT.md` - 管理员功能最终部署
- `ADMIN_LOGIN_FEATURE_SUMMARY.md` - 管理员登录功能总结
- `ADMIN_TIMEOUT_AMOUNT_FEATURE.md` - 管理员超时金额功能
- `ADMIN_TIMEOUT_AMOUNT_SUMMARY.md` - 管理员超时金额总结
- `ADMIN_UPDATE_ORDER_AMOUNT_FEATURE.md` - 管理员修改订单金额功能

### order/ - 订单功能文档
- `TIMEOUT_AMOUNT_FIX_SUMMARY.md` - 超时金额修复总结
- `TIMEOUT_AMOUNT_ZERO_ANALYSIS.md` - 超时金额为零分析

### payment/ - 支付功能文档
- `PAYMENT_SETUP_GUIDE.md` - 支付设置指南
- `PAYMENT_IMPLEMENTATION_SUMMARY.md` - 支付实现总结
- `REFUND_FEATURE_IMPLEMENTATION.md` - 退款功能实现

### service/ - 服务功能文档
- `qrcode_system.md` - 二维码系统文档
- `promoter_code_system.md` - 推广码系统文档

### user/ - 用户功能文档
- `BACKEND_USER_ID_FIX_SUMMARY.md` - 后端用户ID修复总结
- `BACKEND_USER_ID_MIGRATION_SUMMARY.md` - 后端用户ID迁移总结
- `USER_ID_MIGRATION_SUMMARY.md` - 用户ID迁移总结
- `USER_ID_RECORD_NOT_FOUND_FIX.md` - 用户ID记录未找到修复

### deployment/ - 部署相关文档
- `CLOUD_CONTAINER_ANALYSIS.md` - 云容器分析
- `CLOUD_CONTAINER_DEPLOYMENT.md` - 云容器部署
- `GLOBAL_CLOUD_CONTAINER_USAGE.md` - 全局云容器使用
- `EMERGENCY_FIX_GUIDE.md` - 紧急修复指南

## 通用API文档

根目录下还包含以下通用API文档：
- `service_apis.md` - 服务相关API
- `order_apis.md` - 订单相关API
- `user_apis.md` - 用户相关API
- `referral_apis.md` - 推荐相关API
- `kefu_hospital_apis.md` - 客服医院API
- `config_api.md` - 配置API
- `upload_api.md` - 上传API
- `home_init_api.md` - 首页初始化API
- `wx_login_api.md` - 微信登录API
- `logging_system.md` - 日志系统
- `cos_setup.md` - 对象存储设置

## 文档说明

- 每个功能模块都有对应的文档说明其实现细节
- 部署相关文档包含容器化部署的完整流程
- API文档详细说明了各个接口的使用方法
- 修复文档记录了问题排查和解决方案

## 维护说明

- 新增功能时请同步更新对应文档
- 修复问题时请记录在相应的修复文档中
- 部署变更时请更新部署相关文档 