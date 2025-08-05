-- 为Users表添加userId字段
ALTER TABLE Users ADD COLUMN userId VARCHAR(24) UNIQUE NOT NULL COMMENT '用户唯一标识符' AFTER id;

-- 为现有用户生成userId（如果表中有数据）
-- 注意：这个更新语句需要在应用代码中处理，因为需要生成随机字符串
-- UPDATE Users SET userId = CONCAT('user_', UNIX_TIMESTAMP(), '_', FLOOR(RAND() * 1000000)) WHERE userId IS NULL OR userId = '';

-- 添加索引
CREATE INDEX idx_user_id ON Users(userId); 