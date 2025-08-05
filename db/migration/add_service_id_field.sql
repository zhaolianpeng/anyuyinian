-- 为Services表添加serviceId字段
ALTER TABLE Services ADD COLUMN serviceId INT COMMENT '服务ID，用于前端跳转' AFTER id;

-- 更新现有数据的serviceId字段，使其与id字段值相同
UPDATE Services SET serviceId = id WHERE serviceId IS NULL;

-- 为serviceId字段添加索引
ALTER TABLE Services ADD INDEX idx_service_id (serviceId); 