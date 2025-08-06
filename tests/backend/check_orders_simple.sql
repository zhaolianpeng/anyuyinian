-- 简化的订单数据检查脚本

-- 1. 查看订单总数
SELECT COUNT(*) as total_orders FROM Orders;

-- 2. 查看订单状态分布
SELECT 
    status,
    payStatus,
    COUNT(*) as count,
    IFNULL(SUM(totalAmount), 0) as total_amount
FROM Orders 
GROUP BY status, payStatus
ORDER BY status, payStatus;

-- 3. 查看待支付订单
SELECT 
    id,
    orderNo,
    userId,
    totalAmount,
    status,
    payStatus,
    payDeadline,
    createdAt
FROM Orders 
WHERE status = 0 AND payStatus = 0
ORDER BY createdAt DESC;

-- 4. 查看有支付截止时间的订单
SELECT 
    id,
    orderNo,
    userId,
    totalAmount,
    status,
    payStatus,
    payDeadline,
    createdAt
FROM Orders 
WHERE payDeadline IS NOT NULL
ORDER BY payDeadline DESC;

-- 5. 查看超时未支付订单
SELECT 
    id,
    orderNo,
    userId,
    totalAmount,
    status,
    payStatus,
    payDeadline,
    createdAt
FROM Orders 
WHERE status = 0 
  AND payStatus = 0 
  AND payDeadline IS NOT NULL 
  AND payDeadline < NOW()
ORDER BY payDeadline DESC;

-- 6. 统计超时未支付订单
SELECT 
    COUNT(*) as timeout_order_count,
    IFNULL(SUM(totalAmount), 0) as timeout_total_amount
FROM Orders 
WHERE status = 0 
  AND payStatus = 0 
  AND payDeadline IS NOT NULL 
  AND payDeadline < NOW(); 