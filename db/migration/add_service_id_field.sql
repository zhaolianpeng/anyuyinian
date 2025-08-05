-- 为Services表添加serviceitemid字段
ALTER TABLE Services ADD COLUMN serviceitemid INT COMMENT '服务项目ID，用于前端跳转' AFTER id;

-- 更新现有数据的serviceitemid字段，使其与id字段值相同
UPDATE Services SET serviceitemid = id WHERE serviceitemid IS NULL;

-- 为serviceitemid字段添加索引
ALTER TABLE Services ADD INDEX idx_serviceitemid (serviceitemid); 