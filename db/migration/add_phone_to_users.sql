-- 添加phone字段到Users表
ALTER TABLE Users ADD COLUMN phone VARCHAR(20) DEFAULT '' COMMENT '用户手机号';

-- 为phone字段添加索引（可选）
CREATE INDEX idx_users_phone ON Users(phone); 