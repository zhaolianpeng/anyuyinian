-- 添加预约咨询服务
-- 创建时间：2024-12-20
-- 描述：为首页一键预约功能创建预约咨询服务

-- 1. 在ServiceItems表中添加预约咨询服务项目
INSERT INTO `ServiceItems` (
    `name`, 
    `description`, 
    `category`, 
    `price`, 
    `originalPrice`, 
    `imageUrl`, 
    `detailImages`, 
    `formConfig`, 
    `status`, 
    `sort`, 
    `createdAt`, 
    `updatedAt`
) VALUES (
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
);

-- 2. 在Services表中添加预约咨询服务
INSERT INTO `Services` (
    `name`, 
    `description`, 
    `icon`, 
    `imageUrl`, 
    `linkUrl`, 
    `price`, 
    `category`, 
    `sort`, 
    `status`, 
    `createdAt`, 
    `updatedAt`
) VALUES (
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
);

-- 3. 更新Services表的serviceitemid字段，建立关联关系
UPDATE `Services` 
SET `serviceitemid` = (SELECT `id` FROM `ServiceItems` WHERE `name` = '预约咨询服务' AND `category` = '预约咨询' LIMIT 1)
WHERE `name` = '预约咨询服务' AND `category` = '预约咨询';

-- 4. 创建索引（如果不存在）
-- 注意：如果索引已存在，会报错，可以忽略
CREATE INDEX idx_serviceitems_appointment_category ON `ServiceItems` (`category`);
CREATE INDEX idx_services_appointment_category ON `Services` (`category`);

-- 5. 验证数据插入
SELECT 
    'ServiceItems' as table_name,
    id,
    name,
    category,
    price,
    status
FROM `ServiceItems` 
WHERE `name` = '预约咨询服务' AND `category` = '预约咨询';

SELECT 
    'Services' as table_name,
    id,
    name,
    category,
    price,
    serviceitemid,
    status
FROM `Services` 
WHERE `name` = '预约咨询服务' AND `category` = '预约咨询';
