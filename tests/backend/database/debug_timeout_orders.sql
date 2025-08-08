-- 检查超时未支付订单的SQL脚本

-- 1. 查看所有订单的基本信息
SELECT 
    id,
    orderNo,
    userId,
    totalAmount,
    status,
    payStatus,
    payDeadline,
    createdAt,
    updatedAt
FROM Orders 
ORDER BY createdAt DESC 
LIMIT 20;

-- 2. 查看所有待支付订单（status = 0 且 payStatus = 0）
SELECT 
    id,
    orderNo,
    userId,
    totalAmount,
    status,
    payStatus,
    payDeadline,
    createdAt,
    updatedAt
FROM Orders 
WHERE status = 0 AND payStatus = 0
ORDER BY createdAt DESC;

-- 3. 查看有支付截止时间的订单
SELECT 
    id,
    orderNo,
    userId,
    totalAmount,
    status,
    payStatus,
    payDeadline,
    createdAt,
    updatedAt
FROM Orders 
WHERE payDeadline IS NOT NULL
ORDER BY payDeadline DESC;

-- 4. 查看超时未支付订单（核心查询）
SELECT 
    id,
    orderNo,
    userId,
    totalAmount,
    status,
    payStatus,
    payDeadline,
    createdAt,
    updatedAt,
    NOW() as current_datetime,
    TIMESTAMPDIFF(MINUTE, payDeadline, NOW()) as minutes_overdue
FROM Orders 
WHERE status = 0 
  AND payStatus = 0 
  AND payDeadline IS NOT NULL 
  AND payDeadline < NOW()
ORDER BY payDeadline DESC;

-- 5. 统计超时未支付订单总金额
SELECT 
    COUNT(*) as timeout_order_count,
    IFNULL(SUM(totalAmount), 0) as timeout_total_amount
FROM Orders 
WHERE status = 0 
  AND payStatus = 0 
  AND payDeadline IS NOT NULL 
  AND payDeadline < NOW();

-- 6. 查看订单状态分布
SELECT 
    status,
    payStatus,
    COUNT(*) as count,
    IFNULL(SUM(totalAmount), 0) as total_amount
FROM Orders 
GROUP BY status, payStatus
ORDER BY status, payStatus;

-- 7. 查看支付截止时间分布
SELECT 
    CASE 
        WHEN payDeadline IS NULL THEN '无截止时间'
        WHEN payDeadline < NOW() THEN '已超时'
        WHEN payDeadline > NOW() THEN '未超时'
    END as deadline_status,
    COUNT(*) as count,
    IFNULL(SUM(totalAmount), 0) as total_amount
FROM Orders 
GROUP BY 
    CASE 
        WHEN payDeadline IS NULL THEN '无截止时间'
        WHEN payDeadline < NOW() THEN '已超时'
        WHEN payDeadline > NOW() THEN '未超时'
    END; 