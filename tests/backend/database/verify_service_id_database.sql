-- 验证serviceitemid字段是否正确添加
USE anyuyinian;

-- 1. 检查Services表结构
DESCRIBE Services;

-- 2. 查看Services表中的数据，包括serviceitemid字段
SELECT id, serviceitemid, name, description, status FROM Services ORDER BY id;

-- 3. 验证serviceitemid字段是否与id字段值相同（初始状态）
SELECT 
    id, 
    serviceitemid, 
    name,
    CASE 
        WHEN id = serviceitemid THEN '✅ 一致'
        ELSE '❌ 不一致'
    END as status
FROM Services 
ORDER BY id;

-- 4. 统计serviceitemid字段的分布
SELECT 
    COUNT(*) as total_services,
    COUNT(serviceitemid) as services_with_serviceitemid,
    COUNT(*) - COUNT(serviceitemid) as services_without_serviceitemid
FROM Services;

-- 5. 显示serviceitemid字段的索引
SHOW INDEX FROM Services WHERE Key_name = 'idx_serviceitemid'; 