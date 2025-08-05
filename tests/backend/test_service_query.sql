-- 测试Services表查询
USE anyuyinian;

-- 1. 查看Services表结构
DESCRIBE Services;

-- 2. 查看Services表中的所有数据
SELECT id, serviceitemid, name, description, status FROM Services ORDER BY id;

-- 3. 测试查询语句（模拟DAO查询）
SELECT id, serviceitemid, name, description, icon, imageUrl, linkUrl, sort, status, createdAt, updatedAt 
FROM Services 
WHERE status = 1 
ORDER BY sort ASC, id DESC;

-- 4. 验证serviceitemid字段的值
SELECT 
    id, 
    serviceitemid, 
    name,
    CASE 
        WHEN serviceitemid IS NOT NULL THEN '✅ 有值'
        ELSE '❌ 空值'
    END as serviceitemid_status,
    CASE 
        WHEN id = serviceitemid THEN '✅ 与id相同'
        WHEN serviceitemid IS NOT NULL THEN '⚠️ 与id不同'
        ELSE '❌ 无serviceitemid'
    END as comparison
FROM Services 
ORDER BY id; 