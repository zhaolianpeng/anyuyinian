-- 为Referrals表添加推广码字段
ALTER TABLE Referrals ADD COLUMN promoterCode VARCHAR(6) UNIQUE COMMENT '六位推广码' AFTER referrerId;

-- 为现有记录生成推广码（如果表中有数据）
-- 注意：这个更新语句需要在应用代码中处理，因为需要生成随机字符串
-- UPDATE Referrals SET promoterCode = CONCAT(
--   CHAR(65 + FLOOR(RAND() * 26)), 
--   CHAR(65 + FLOOR(RAND() * 26)), 
--   CHAR(65 + FLOOR(RAND() * 26)),
--   FLOOR(RAND() * 10),
--   FLOOR(RAND() * 10),
--   FLOOR(RAND() * 10)
-- ) WHERE promoterCode IS NULL OR promoterCode = '';

-- 添加索引
CREATE INDEX idx_promoter_code ON Referrals(promoterCode); 