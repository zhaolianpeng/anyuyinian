-- 检查服务数据
-- 查询所有服务
SELECT 
    id,
    name,
    category,
    price,
    originalPrice,
    status,
    sort,
    createdAt
FROM ServiceItems 
ORDER BY id;

-- 查询上架的服务
SELECT 
    id,
    name,
    category,
    price,
    originalPrice,
    sort
FROM ServiceItems 
WHERE status = 1 
ORDER BY sort;

-- 统计服务数量
SELECT 
    COUNT(*) as total_services,
    COUNT(CASE WHEN status = 1 THEN 1 END) as active_services,
    COUNT(CASE WHEN status = 0 THEN 1 END) as inactive_services
FROM ServiceItems;

-- 按分类统计服务
SELECT 
    category,
    COUNT(*) as service_count,
    COUNT(CASE WHEN status = 1 THEN 1 END) as active_count
FROM ServiceItems 
GROUP BY category 
ORDER BY service_count DESC; 