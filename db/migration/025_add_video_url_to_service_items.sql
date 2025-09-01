-- 为ServiceItems表添加videoUrl字段
-- 迁移文件：025_add_video_url_to_service_items.sql

-- 添加videoUrl字段到ServiceItems表
ALTER TABLE ServiceItems ADD COLUMN videoUrl VARCHAR(500) COMMENT '视频URL，用于智慧养老设备的产品介绍视频';
