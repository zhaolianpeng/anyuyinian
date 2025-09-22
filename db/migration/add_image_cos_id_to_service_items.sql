-- 为 ServiceItems 表添加 imageCosId 字段
-- 用于存储腾讯云对象存储的图片ID

ALTER TABLE ServiceItems 
ADD COLUMN imageCosId VARCHAR(255) COMMENT '腾讯云对象存储图片ID' AFTER imageUrl;

-- 添加索引以提高查询性能
CREATE INDEX idx_service_items_image_cos_id ON ServiceItems(imageCosId);

-- 更新现有记录的 imageCosId 字段（可选）
-- 如果现有记录有 imageUrl，可以尝试从 URL 中提取 COS ID
-- UPDATE ServiceItems 
-- SET imageCosId = SUBSTRING_INDEX(SUBSTRING_INDEX(imageUrl, '/', -1), '.', 1)
-- WHERE imageUrl IS NOT NULL AND imageUrl != '' AND imageCosId IS NULL;
