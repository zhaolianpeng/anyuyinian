-- 护工服务平台首页数据迁移文件
-- 使用现有的Services和ServiceItems表，添加护工服务相关数据

-- 1. 为Services表添加价格和分类字段（如果不存在）
ALTER TABLE `Services` 
ADD COLUMN `price` decimal(10,2) COMMENT '服务价格' AFTER `linkUrl`,
ADD COLUMN `category` varchar(100) COMMENT '服务分类' AFTER `price`;

-- 2. 为ServiceItems表添加护工服务相关字段（如果不存在）
ALTER TABLE `ServiceItems` 
ADD COLUMN `duration` varchar(100) COMMENT '服务时长，如"26天/月"' AFTER `price`,
ADD COLUMN `features` text COMMENT '服务特色，JSON格式' AFTER `duration`;

-- 3. 插入护工服务数据到Services表
INSERT INTO `Services` (`name`, `description`, `icon`, `imageUrl`, `linkUrl`, `price`, `category`, `sort`, `status`, `createdAt`, `updatedAt`) VALUES
('慢病照护', '生活支援,守护健康', '🏥', '/images/service/chronic-care.jpg', '/pages/service/detail?id=1', 4880.00, '护工服务', 1, 1, NOW(), NOW()),
('居家术后照护', '省心省力、全天照护', '🏠', '/images/service/post-surgery.jpg', '/pages/service/detail?id=2', 5580.00, '护工服务', 2, 1, NOW(), NOW()),
('康复照护', '偏瘫、肢体康复训练', '🦽', '/images/service/rehabilitation.jpg', '/pages/service/detail?id=3', 6280.00, '护工服务', 3, 1, NOW(), NOW()),
('认知症照护', '认知症(阿尔兹海默病)习惯培养、守护健康', '🧠', '/images/service/dementia-care.jpg', '/pages/service/detail?id=4', 4980.00, '护工服务', 4, 1, NOW(), NOW()),
('肌无力照护', '行动支持、安全防护', '🦴', '/images/service/muscle-weakness.jpg', '/pages/service/detail?id=5', 5380.00, '护工服务', 5, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE 
`description` = VALUES(`description`),
`price` = VALUES(`price`),
`category` = VALUES(`category`),
`updatedAt` = NOW();

-- 4. 插入对应的ServiceItems数据
INSERT INTO `ServiceItems` (`name`, `description`, `category`, `price`, `duration`, `features`, `imageUrl`, `status`, `sort`, `createdAt`, `updatedAt`) VALUES
('慢病照护', '生活支援,守护健康', '护工服务', 4880.00, '26天/月', '["专业护理","健康监测","用药提醒"]', '/images/service/chronic-care.jpg', 1, 1, NOW(), NOW()),
('居家术后照护', '省心省力、全天照护', '护工服务', 5580.00, '26天/月', '["伤口护理","康复训练","生活照料"]', '/images/service/post-surgery.jpg', 1, 2, NOW(), NOW()),
('康复照护', '偏瘫、肢体康复训练', '护工服务', 6280.00, '26天/月', '["专业康复","运动训练","功能恢复"]', '/images/service/rehabilitation.jpg', 1, 3, NOW(), NOW()),
('认知症照护', '认知症(阿尔兹海默病)习惯培养、守护健康', '护工服务', 4980.00, '26天/月', '["认知训练","行为管理","安全监护"]', '/images/service/dementia-care.jpg', 1, 4, NOW(), NOW()),
('肌无力照护', '行动支持、安全防护', '护工服务', 5380.00, '26天/月', '["行动辅助","安全防护","生活照料"]', '/images/service/muscle-weakness.jpg', 1, 5, NOW(), NOW())
ON DUPLICATE KEY UPDATE 
`description` = VALUES(`description`),
`price` = VALUES(`price`),
`duration` = VALUES(`duration`),
`features` = VALUES(`features`),
`updatedAt` = NOW();

-- 5. 更新现有Services记录的category字段（如果为空）
UPDATE `Services` SET `category` = '护工服务' WHERE (`category` IS NULL OR `category` = '') AND `name` LIKE '%照护%';
UPDATE `Services` SET `category` = '护工服务' WHERE (`category` IS NULL OR `category` = '') AND `name` LIKE '%护理%';
UPDATE `Services` SET `category` = '护工服务' WHERE (`category` IS NULL OR `category` = '') AND `name` LIKE '%陪护%';

-- 6. 为现有Services记录添加默认价格（如果为空）
UPDATE `Services` SET `price` = 5000.00 WHERE (`price` IS NULL OR `price` = 0) AND `category` = '护工服务';

-- 7. 创建索引优化查询性能（使用标准MySQL语法）
-- 注意：如果索引已存在，这些语句会报错，可以忽略错误
CREATE INDEX `idx_services_category` ON `Services` (`category`);
CREATE INDEX `idx_services_price` ON `Services` (`price`);
CREATE INDEX `idx_serviceitems_category` ON `ServiceItems` (`category`);
