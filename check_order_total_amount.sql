-- 检查订单表中的总金额字段
SELECT 
    id,
    orderNo,
    serviceName,
    price,
    quantity,
    totalAmount,
    (price * quantity) as calculated_total,
    CASE 
        WHEN totalAmount = (price * quantity) THEN '正确'
        WHEN totalAmount = 0 THEN '错误：总金额为0'
        ELSE '错误：总金额不匹配'
    END as status
FROM Orders 
ORDER BY id DESC 
LIMIT 10;

-- 统计总金额为0的订单数量
SELECT 
    COUNT(*) as total_orders,
    SUM(CASE WHEN totalAmount = 0 THEN 1 ELSE 0 END) as zero_amount_orders,
    SUM(CASE WHEN totalAmount > 0 THEN 1 ELSE 0 END) as valid_amount_orders
FROM Orders;

-- 更新总金额为0的订单（如果需要修复）
-- UPDATE Orders 
-- SET totalAmount = price * quantity 
-- WHERE totalAmount = 0; 