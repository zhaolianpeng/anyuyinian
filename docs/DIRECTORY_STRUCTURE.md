# 项目目录结构整理

## 后端目录结构

### 文档目录 (`docs/`)
```
docs/
├── README.md                    # 文档总览
├── backend/                     # 后端API文档
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
│   ├── SERVICE_ID_DEBUG_SUMMARY.md
│   ├── SERVICE_ID_FEATURE_IMPLEMENTATION.md
│   ├── SERVICE_ID_IMPLEMENTATION_SUMMARY.md
│   ├── SERVICEITEMID_UPDATE_SUMMARY.md
│   ├── TEST_DATA_CLEANUP_SUMMARY.md
│   ├── USER_CREATION_FIX.md
│   ├── USER_INFO_FIX.md
│   ├── USER_PROFILE_FEATURES.md
│   └── WX_CONFIG_FIX.md
├── frontend/                    # 前端文档
│   ├── NAVIGATION_JUMP_FEATURE.md
│   ├── ORDER_NAVIGATION_FEATURE.md
│   ├── PATIENT_MANAGEMENT_UPDATE.md
│   ├── README.md
│   └── SETUP_PROFILE_DEBUG_GUIDE.md
├── websocket/                   # WebSocket相关文档
│   ├── WEBSOCKET_DIAGNOSIS_SUMMARY.md
│   ├── WEBSOCKET_ERROR_ANALYSIS.md
│   └── WEBSOCKET_TROUBLESHOOTING.md
├── debug/                       # 调试相关文档
│   ├── debug_websocket.js
│   └── debug_websocket_detailed.js
├── build/                       # 构建相关文档
│   └── BUILD_FIX_SUMMARY.md
├── CLOUD_CONTAINER_ANALYSIS.md  # 云托管分析
├── GLOBAL_CLOUD_CONTAINER_USAGE.md
└── DOCUMENT_ORGANIZATION_SUMMARY.md
```

### 测试目录 (`tests/`)
```
tests/
├── README.md                    # 测试说明
├── api/                         # API测试
│   ├── test_build
│   ├── test_config.sh
│   ├── test_home_init.sh
│   ├── test_hospital_detail_api.sh
│   ├── test_hospital_detail_simple.sh
│   ├── test_kefu_hospital.sh
│   ├── test_order.sh
│   ├── test_referral.sh
│   ├── test_service.sh
│   ├── test_upload.sh
│   ├── test_user_info.sh
│   └── test_websocket.sh
├── backend/                     # 后端测试
│   ├── check_order_total_amount.sql
│   ├── check_services.sql
│   ├── debug_service_query.go
│   ├── fix_order_total_amount.sql
│   ├── init_service_data.sql
│   ├── test_api_service_id.sh
│   ├── test_order_api_changes.sh
│   ├── test_order_countdown.sh
│   ├── test_order_detail_appointment.sh
│   └── test_order_total_amount.sh
├── database/                    # 数据库测试
└── frontend/                    # 前端测试
```

## 小程序端目录结构

### 文档目录 (`docs/`)
```
docs/
├── README.md                    # 文档总览
├── CHANGELOG.md                 # 更新日志
├── api/                         # API相关文档
│   ├── API_FIX_SUMMARY.md
│   ├── API_METHOD_FIX.md
│   └── API_PATH_FIX.md
├── features/                    # 功能特性文档
│   ├── AGREEMENT_PAGES_SUMMARY.md
│   ├── BUSINESS_INTEGRATION_SUMMARY.md
│   ├── HOSPITAL_PAGES_IMPLEMENTATION.md
│   ├── NAVIGATION_FIX_SUMMARY.md
│   ├── REALTIME_COMMUNICATION_SUMMARY.md
│   └── WEBSOCKET_IMPLEMENTATION.md
├── fixes/                       # 修复相关文档
│   ├── ALL_API_FIX_SUMMARY.md
│   ├── DATA_DISPLAY_FIX.md
│   ├── FINAL_MIGRATION_SUMMARY.md
│   ├── FINAL_SOLUTION.md
│   ├── HOME_DATA_FIX.md
│   ├── IMAGE_ERROR_FIX.md
│   ├── LOGIN_BACKEND_FIX.md
│   ├── LOGIN_FIX.md
│   ├── LOGIN_FIX_DESC.md
│   ├── LOGIN_FIX_DIRECT.md
│   ├── LOGIN_FIX_FINAL.md
│   ├── LOGIN_FIX_SIMPLE.md
│   ├── LOGIN_FIX_SYNC.md
│   ├── LOGIN_JUMP_FIX.md
│   ├── LOGIN_SUCCESS_SUMMARY.md
│   ├── NETWORK_FIX_SUMMARY.md
│   ├── ORDER_BOOKING_FEATURES.md
│   ├── ORDER_BUTTON_FIX.md
│   ├── ORDER_PAGE_FIELDS_DEBUG.md
│   ├── ORDER_PAGE_FIELDS_SUMMARY.md
│   ├── ORDER_PAGE_FIX.md
│   ├── ORDER_PAGE_IMPLEMENTATION.md
│   ├── ORDER_PAGE_SCROLL_FIX.md
│   ├── ORDER_TIMEOUT_FEATURE.md
│   ├── QUICK_FIX_GUIDE.md
│   ├── SERVICE_DETAIL_FIELDS_UPDATE.md
│   ├── SERVICE_ID_ERROR_FIX.md
│   ├── SERVICE_LIST_IMAGE_FIX.md
│   └── SERVICE_NAVIGATION_FIX.md
└── setup/                       # 设置配置文档
    ├── CLOUD_CONTAINER_MIGRATION_SUMMARY.md
    ├── CLOUD_CONTAINER_USAGE.md
    ├── CLOUD_HOSTING_LIMITATIONS.md
    ├── CLOUD_INIT_FIX.md
    ├── CLOUD_URL_CONVERSION.md
    ├── CONTAINER_API_GUIDE.md
    ├── CONTAINER_USAGE_GUIDE.md
    ├── COS_DOMAIN_SETUP.md
    ├── COS_SETUP.md
    ├── COS_VERIFICATION.md
    ├── DEPLOYMENT_GUIDE.md
    ├── DEVELOPMENT.md
    ├── DOMAIN_CONFIGURATION.md
    ├── check_domain_config.md
    ├── CALL_CONTAINER_NETWORK_FIX.md
    ├── NETWORK_ERROR_QUICK_FIX.md
    ├── QUICK_NETWORK_FIX.md
    └── WEBSOCKET_IMPLEMENTATION_GUIDE.md
```

### 测试目录 (`tests/`)
```
tests/
├── README.md                    # 测试说明
├── agreement/                   # 协议页面测试
├── api/                         # API测试
├── cloud/                       # 云服务测试
├── fixes/                       # 修复测试
├── home/                        # 首页测试
├── hospital/                    # 医院页面测试
├── image/                       # 图片处理测试
├── login/                       # 登录测试
├── navigation/                  # 导航测试
├── network/                     # 网络测试
├── order/                       # 订单测试
└── service/                     # 服务测试
```

## 整理说明

### 整理原则
1. **按功能分类**: 将相关功能的文档和测试文件归类到同一目录
2. **按类型分类**: 区分文档、测试、配置等不同类型的文件
3. **保持层次**: 使用子目录进一步细分，便于查找
4. **保留原文件**: 所有文件都保留，只是重新组织位置

### 主要改进
1. **文档结构清晰**: 按功能模块组织文档
2. **测试文件分类**: 按测试类型和功能模块分类
3. **便于维护**: 相关文件集中存放，便于查找和维护
4. **支持扩展**: 目录结构支持后续功能扩展

### 使用建议
1. **新增文档**: 根据功能模块放入对应目录
2. **新增测试**: 根据测试类型放入对应目录
3. **查找文件**: 根据功能或类型在对应目录查找
4. **维护更新**: 定期整理新增的文件到对应目录 