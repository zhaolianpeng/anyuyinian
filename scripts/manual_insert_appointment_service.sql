-- 手动插入预约咨询服务
-- 如果自动迁移脚本有问题，可以使用此脚本手动插入数据

-- 1. 先检查是否已存在
SELECT '检查现有数据' as step;

SELECT 
    'ServiceItems' as table_name,
    COUNT(*) as count
FROM ServiceItems 
WHERE name = '预约咨询服务' OR category = '预约咨询';

SELECT 
    'Services' as table_name,
    COUNT(*) as count
FROM Services 
WHERE name = '预约咨询服务' OR category = '预约咨询';

-- 2. 如果不存在，则插入数据
-- 插入ServiceItems表
INSERT INTO ServiceItems (
    name, 
    description, 
    category, 
    price, 
    originalPrice, 
    imageUrl, 
    detailImages, 
    formConfig, 
    status, 
    sort, 
    createdAt, 
    updatedAt
) 
SELECT 
    '预约咨询服务',
    '专业护工服务预约咨询，了解服务详情、价格、流程等信息，为您提供最适合的护工服务方案。',
    '预约咨询',
    0.01,
    0.01,
    '/images/service/appointment-consultation.jpg',
    '["/images/service/appointment-consultation-detail1.jpg","/images/service/appointment-consultation-detail2.jpg"]',
    '{"fields":[{"name":"姓名","type":"text","required":true,"placeholder":"请输入您的姓名"},{"name":"手机号","type":"phone","required":true,"placeholder":"请输入您的手机号"},{"name":"服务类型","type":"select","required":true,"options":["居家照护","医院陪诊","周期护理","家政服务"],"placeholder":"请选择服务类型"},{"name":"服务时间","type":"date","required":true,"placeholder":"请选择服务开始时间"},{"name":"服务地址","type":"text","required":true,"placeholder":"请输入服务地址"},{"name":"特殊需求","type":"textarea","required":false,"placeholder":"请描述您的特殊需求（可选）"}]}',
    1,
    999,
    NOW(),
    NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM ServiceItems 
    WHERE name = '预约咨询服务' AND category = '预约咨询'
);

-- 插入Services表
INSERT INTO Services (
    name, 
    description, 
    icon, 
    imageUrl, 
    linkUrl, 
    price, 
    category, 
    sort, 
    status, 
    createdAt, 
    updatedAt
)
SELECT 
    '预约咨询服务',
    '专业护工服务预约咨询，了解服务详情、价格、流程等信息，为您提供最适合的护工服务方案。',
    '📋',
    '/images/service/appointment-consultation.jpg',
    '/pages/service/detail?id=999',
    0.01,
    '预约咨询',
    999,
    1,
    NOW(),
    NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM Services 
    WHERE name = '预约咨询服务' AND category = '预约咨询'
);

-- 3. 更新Services表的serviceitemid字段，建立关联关系
UPDATE Services s
JOIN ServiceItems si ON s.name = si.name AND s.category = si.category
SET s.serviceitemid = si.id
WHERE s.name = '预约咨询服务' AND s.category = '预约咨询'
AND s.serviceitemid IS NULL;

-- 4. 验证插入结果
SELECT '验证插入结果' as step;

SELECT 
    'ServiceItems' as table_name,
    id,
    name,
    category,
    price,
    status,
    createdAt
FROM ServiceItems 
WHERE name = '预约咨询服务' AND category = '预约咨询';

SELECT 
    'Services' as table_name,
    id,
    name,
    category,
    price,
    serviceitemid,
    status,
    createdAt
FROM Services 
WHERE name = '预约咨询服务' AND category = '预约咨询';

-- 5. 检查关联关系
SELECT '检查关联关系' as step;

SELECT 
    s.id as service_id,
    s.name as service_name,
    s.serviceitemid,
    si.id as serviceitem_id,
    si.name as serviceitem_name,
    CASE 
        WHEN s.serviceitemid = si.id THEN '✓ 关联正确'
        ELSE '✗ 关联错误'
    END as relationship_status
FROM Services s
JOIN ServiceItems si ON s.serviceitemid = si.id
WHERE s.name = '预约咨询服务' AND s.category = '预约咨询';

-- 6. 获取实际的服务ID（用于前端调用）
SELECT '获取服务ID信息' as step;

SELECT 
    '前端调用信息' as info,
    CONCAT('服务详情页面URL: /pages/service/detail?id=', s.id) as detail_url,
    CONCAT('服务项目ID: ', s.serviceitemid) as serviceitem_id,
    CONCAT('服务ID: ', s.id) as service_id,
    s.name as service_name,
    s.price as price,
    s.status as status
FROM Services s
WHERE s.name = '预约咨询服务' AND s.category = '预约咨询';
