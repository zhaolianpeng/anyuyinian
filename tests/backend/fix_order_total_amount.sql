-- 修复订单总金额字段
-- 将总金额为0的订单更新为正确的计算值

-- 首先查看需要修复的订单
SELECT 
    id,
    orderNo,
    serviceName,
    price,
    quantity,
    totalAmount as current_total,
    (price * quantity) as correct_total
FROM Orders 
WHERE totalAmount = 0 
ORDER BY id DESC;

-- 更新总金额为0的订单
UPDATE Orders 
SET totalAmount = price * quantity 
WHERE totalAmount = 0;

-- 验证修复结果
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
        ELSE '错误'
    END as status
FROM Orders 
ORDER BY id DESC 
LIMIT 10; 