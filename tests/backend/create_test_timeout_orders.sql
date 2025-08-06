-- 创建超时未支付订单的测试数据

-- 1. 插入一些超时未支付的订单（支付截止时间设置为过去的时间）
INSERT INTO Orders (
    orderNo, 
    userId, 
    serviceId, 
    patientId, 
    addressId, 
    appointmentDate, 
    appointmentTime, 
    serviceName, 
    price, 
    quantity, 
    totalAmount, 
    status, 
    payStatus, 
    payDeadline, 
    createdAt, 
    updatedAt
) VALUES 
-- 订单1：超时30分钟
('ORDER_TEST_TIMEOUT_001', '507f1f77bcf86cd799439011', 1, 1, 1, '2024-12-20', '09:00', '测试服务1', 299.00, 1, 299.00, 0, 0, DATE_SUB(NOW(), INTERVAL 30 MINUTE), NOW(), NOW()),

-- 订单2：超时2小时
('ORDER_TEST_TIMEOUT_002', '507f1f77bcf86cd799439011', 1, 1, 1, '2024-12-20', '10:00', '测试服务2', 150.00, 1, 150.00, 0, 0, DATE_SUB(NOW(), INTERVAL 2 HOUR), NOW(), NOW()),

-- 订单3：超时1天
('ORDER_TEST_TIMEOUT_003', '507f1f77bcf86cd799439011', 1, 1, 1, '2024-12-19', '14:00', '测试服务3', 599.00, 1, 599.00, 0, 0, DATE_SUB(NOW(), INTERVAL 1 DAY), NOW(), NOW()),

-- 订单4：超时3天
('ORDER_TEST_TIMEOUT_004', '507f1f77bcf86cd799439011', 1, 1, 1, '2024-12-17', '16:00', '测试服务4', 199.00, 1, 199.00, 0, 0, DATE_SUB(NOW(), INTERVAL 3 DAY), NOW(), NOW()),

-- 订单5：超时1周
('ORDER_TEST_TIMEOUT_005', '507f1f77bcf86cd799439011', 1, 1, 1, '2024-12-13', '11:00', '测试服务5', 399.00, 1, 399.00, 0, 0, DATE_SUB(NOW(), INTERVAL 1 WEEK), NOW(), NOW());

-- 2. 验证插入的数据
SELECT 
    id,
    orderNo,
    userId,
    totalAmount,
    status,
    payStatus,
    payDeadline,
    createdAt,
    TIMESTAMPDIFF(MINUTE, payDeadline, NOW()) as minutes_overdue
FROM Orders 
WHERE orderNo LIKE 'ORDER_TEST_TIMEOUT_%'
ORDER BY payDeadline DESC;

-- 3. 统计超时未支付订单
SELECT 
    COUNT(*) as timeout_order_count,
    IFNULL(SUM(totalAmount), 0) as timeout_total_amount
FROM Orders 
WHERE status = 0 
  AND payStatus = 0 
  AND payDeadline IS NOT NULL 
  AND payDeadline < NOW(); 