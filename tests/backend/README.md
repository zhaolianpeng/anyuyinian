# 后端测试脚本目录

本目录包含后端API的测试脚本，按功能模块分类组织。

## 目录结构

### admin/ - 管理员功能测试
- `test_admin_timeout_amount.sh` - 管理员超时金额设置测试
- `test_update_order_amount.sh` - 管理员修改订单金额测试
- `test_refund_functionality.sh` - 退款功能测试

### order/ - 订单功能测试
- `test_time_slots_api.sh` - 时间段API测试
- `test_order_page_backend.sh` - 订单页面后端测试
- `test_order_status_filter.sh` - 订单状态筛选测试
- `test_order_list_display.sh` - 订单列表显示测试
- `test_order_detail_total_amount.sh` - 订单详情总金额测试
- `test_order_list_price_display.sh` - 订单列表价格显示测试
- `test_order_detail_appointment.sh` - 订单详情预约测试
- `test_order_price_debug.sh` - 订单价格调试测试
- `test_order_list_data.sh` - 订单列表数据测试
- `test_order_countdown.sh` - 订单倒计时测试
- `test_order_api_changes.sh` - 订单API变更测试
- `test_updated_timeout_amount.sh` - 更新超时金额测试
- `debug_timeout_amount.sh` - 超时金额调试测试

### service/ - 服务功能测试
- `test_service_filter.sh` - 服务分类筛选测试
- `test_qrcode.sh` - 二维码功能测试
- `test_promoter_code.sh` - 推广码功能测试
- `test_promoter_center.sh` - 推广中心测试
- `test_api_service_id.sh` - 服务ID API测试
- `test_service_id_feature.sh` - 服务ID功能测试
- `debug_service_query.go` - 服务查询调试

### payment/ - 支付功能测试
- `test_payment_flow.sh` - 支付流程测试
- `test_payment_config.sh` - 支付配置测试

### user/ - 用户功能测试
- `test_user_info.sh` - 用户信息测试
- `test_user_id_migration.sh` - 用户ID迁移测试

### database/ - 数据库测试
- `verify_timeout_query.sql` - 超时查询验证
- `debug_timeout_orders.sql` - 超时订单调试
- `cleanup_test_orders.sql` - 清理测试订单
- `create_test_timeout_orders.sql` - 创建测试超时订单
- `check_orders_simple.sql` - 简单订单检查
- `test_service_query.sql` - 服务查询测试
- `verify_service_id_database.sql` - 服务ID数据库验证
- `fix_order_total_amount.sql` - 修复订单总金额
- `check_order_total_amount.sql` - 检查订单总金额
- `init_service_data.sql` - 初始化服务数据
- `check_services.sql` - 检查服务数据

## 使用方法

1. 确保后端服务正在运行
2. 根据需要修改脚本中的URL和参数
3. 运行对应的测试脚本：

```bash
# 测试服务分类筛选
./service/test_service_filter.sh

# 测试管理员功能
./admin/test_admin_timeout_amount.sh

# 测试订单功能
./order/test_order_list_display.sh
```

## 注意事项

- 测试前请确保数据库中有测试数据
- 某些测试可能需要管理员权限
- 建议在测试环境中运行，避免影响生产数据 