# 项目文档目录

## 目录结构

```
docs/
├── backend/          # 后端相关文档
│   ├── API_DATA_STRUCTURE_FIX.md
│   ├── BACKEND_API_CHANGES.md
│   ├── BACKEND_FIELDS_SYNC_UPDATE.md
│   ├── BACKEND_FIX_SUMMARY.md
│   ├── BACKEND_WX_FIX.md
│   ├── CALENDAR_PICKER_FINAL_SUMMARY.md
│   ├── CALENDAR_PICKER_FIX_SUMMARY.md
│   ├── CALENDAR_PICKER_IMPLEMENTATION.md
│   ├── CLOUD_CONTAINER_IMPLEMENTATION_SUMMARY.md
│   ├── CONFIG.md
│   ├── HORIZONTAL_SCROLL_FIX_SUMMARY.md
│   ├── LOCATION_PERMISSION_UPDATE.md
│   ├── LOGGING_IMPLEMENTATION_SUMMARY.md
│   ├── MOCK_LOGIN_FIX.md
│   ├── NETWORK_ERROR_FIX_GUIDE.md
│   ├── ORDER_API_CHANGES_SUMMARY.md
│   ├── ORDER_COUNTDOWN_FEATURES.md
│   ├── ORDER_DATA_FIX_SUMMARY.md
│   ├── ORDER_DETAIL_APPOINTMENT_ENHANCEMENT.md
│   ├── ORDER_DETAIL_APPOINTMENT_FIX.md
│   ├── ORDER_FIELDS_SYNC_CHECK.md
│   ├── ORDER_LIST_DISPLAY_FIX.md
│   ├── ORDER_LIST_PRICE_DISPLAY_FIX.md
│   ├── ORDER_PAGE_DEBUG_GUIDE.md
│   ├── ORDER_PAGE_FIX_SUMMARY.md
│   ├── ORDER_PAGE_PATIENT_INFO_UPDATE.md
│   ├── ORDER_PAGE_TROUBLESHOOTING.md
│   ├── ORDER_PAGE_WIDTH_FIX.md
│   ├── ORDER_PRICE_DISPLAY_FIX.md
│   ├── ORDER_PRICE_FIELD_UNIFICATION.md
│   ├── ORDER_TAB_FILTER_IMPLEMENTATION.md
│   ├── ORDER_TOTAL_AMOUNT_FIX.md
│   ├── PATIENT_ADDRESS_ENHANCEMENT_SUMMARY.md
│   ├── PATIENT_AGE_FEATURE_SUMMARY.md
│   ├── REMOVE_IDCARD_BIRTHDAY_SUMMARY.md
│   ├── TEST_DATA_CLEANUP_SUMMARY.md
│   ├── USER_CREATION_FIX.md
│   ├── USER_INFO_FIX.md
│   ├── USER_PROFILE_FEATURES.md
│   └── WX_CONFIG_FIX.md
└── frontend/         # 前端相关文档
    ├── FILE_ORGANIZATION_SUMMARY.md
    └── README.md
```

## 文档说明

### 后端文档 (backend/)

包含所有后端相关的实现文档，主要涵盖：

1. **API相关**
   - API数据结构修复
   - 后端API变更
   - 订单API变更

2. **功能实现**
   - 云托管实现
   - 订单功能实现
   - 用户功能实现
   - 日历选择器实现

3. **问题修复**
   - 网络错误修复
   - 订单显示修复
   - 价格显示修复
   - 患者信息更新

4. **配置相关**
   - 微信配置修复
   - 系统配置说明

### 前端文档 (frontend/)

包含所有前端相关的实现文档，主要涵盖：

1. **项目组织**
   - 文件组织结构
   - 项目README

2. **功能实现**
   - 页面实现
   - 组件实现
   - 工具类实现

## 使用说明

1. **查找文档**: 根据功能模块在对应目录中查找相关文档
2. **更新文档**: 新增功能时请在对应目录创建文档
3. **文档命名**: 使用功能名称 + 类型的方式命名，如 `ORDER_TAB_FILTER_IMPLEMENTATION.md`

## 文档模板

### 实现文档模板
```markdown
# 功能名称实现总结

## 概述
简要描述功能实现的目的和范围

## 实现内容
详细描述实现的具体内容

## 技术细节
描述技术实现的关键点

## 测试验证
描述测试方法和结果

## 使用说明
描述如何使用该功能

## 注意事项
描述使用时的注意事项
```

### 修复文档模板
```markdown
# 问题修复总结

## 问题描述
描述遇到的问题

## 原因分析
分析问题的根本原因

## 解决方案
详细描述解决方案

## 验证结果
描述修复后的验证结果

## 预防措施
描述如何避免类似问题
``` 