-- 清理测试订单数据

-- 删除测试超时订单
DELETE FROM Orders WHERE orderNo LIKE 'ORDER_TEST_TIMEOUT_%';

-- 验证清理结果
SELECT COUNT(*) as remaining_test_orders FROM Orders WHERE orderNo LIKE 'ORDER_TEST_TIMEOUT_%'; 