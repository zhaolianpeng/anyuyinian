-- 更新Services表字段名从serviceId改为serviceitemid
USE anyuyinian;

-- 如果存在serviceId字段，先删除
ALTER TABLE Services DROP COLUMN IF EXISTS serviceId;

-- 如果存在serviceId索引，先删除
DROP INDEX IF EXISTS idx_service_id ON Services;

-- 添加serviceitemid字段
ALTER TABLE Services ADD COLUMN serviceitemid INT COMMENT '服务项目ID，用于前端跳转' AFTER id;

-- 更新现有数据的serviceitemid字段，使其与id字段值相同
UPDATE Services SET serviceitemid = id WHERE serviceitemid IS NULL;

-- 为serviceitemid字段添加索引
ALTER TABLE Services ADD INDEX idx_serviceitemid (serviceitemid);

-- 验证字段是否正确添加
SELECT id, serviceitemid, name FROM Services ORDER BY id; 