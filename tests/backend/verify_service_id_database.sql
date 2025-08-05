-- 验证serviceId字段是否正确添加
USE anyuyinian;

-- 1. 检查Services表结构
DESCRIBE Services;

-- 2. 查看Services表中的数据，包括serviceId字段
SELECT id, serviceId, name, description, status FROM Services ORDER BY id;

-- 3. 验证serviceId字段是否与id字段值相同（初始状态）
SELECT 
    id, 
    serviceId, 
    name,
    CASE 
        WHEN id = serviceId THEN '✅ 一致'
        ELSE '❌ 不一致'
    END as status
FROM Services 
ORDER BY id;

-- 4. 统计serviceId字段的分布
SELECT 
    COUNT(*) as total_services,
    COUNT(serviceId) as services_with_serviceId,
    COUNT(*) - COUNT(serviceId) as services_without_serviceId
FROM Services;

-- 5. 显示serviceId字段的索引
SHOW INDEX FROM Services WHERE Key_name = 'idx_service_id'; 