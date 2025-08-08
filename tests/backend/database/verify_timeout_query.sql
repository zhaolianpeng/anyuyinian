-- 验证修改后的超时未支付查询逻辑

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
ORDER BY createdAt DESC;

-- 2. 查看待支付订单（status = 0 且 payStatus = 0）
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

-- 3. 查看已取消但未支付订单（status = 3 且 payStatus = 0）
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
WHERE status = 3 AND payStatus = 0
ORDER BY createdAt DESC;

-- 4. 查看修改后的超时未支付订单（包含已取消订单）
SELECT 
    id,
    orderNo,
    userId,
    totalAmount,
    status,
    payStatus,
    payDeadline,
    createdAt,
    CASE 
        WHEN status = 0 THEN '待支付'
        WHEN status = 3 THEN '已取消'
        ELSE '其他'
    END as status_text
FROM Orders 
WHERE (status = 0 OR status = 3) 
  AND payStatus = 0 
  AND payDeadline IS NOT NULL 
  AND payDeadline < NOW()
ORDER BY payDeadline DESC;

-- 5. 统计修改后的超时未支付订单总金额
SELECT 
    COUNT(*) as timeout_order_count,
    IFNULL(SUM(totalAmount), 0) as timeout_total_amount,
    COUNT(CASE WHEN status = 0 THEN 1 END) as pending_count,
    COUNT(CASE WHEN status = 3 THEN 1 END) as cancelled_count
FROM Orders 
WHERE (status = 0 OR status = 3) 
  AND payStatus = 0 
  AND payDeadline IS NOT NULL 
  AND payDeadline < NOW();

-- 6. 按状态分组统计超时订单
SELECT 
    CASE 
        WHEN status = 0 THEN '待支付'
        WHEN status = 3 THEN '已取消'
        ELSE '其他'
    END as status_text,
    COUNT(*) as count,
    IFNULL(SUM(totalAmount), 0) as total_amount
FROM Orders 
WHERE (status = 0 OR status = 3) 
  AND payStatus = 0 
  AND payDeadline IS NOT NULL 
  AND payDeadline < NOW()
GROUP BY status
ORDER BY status; 