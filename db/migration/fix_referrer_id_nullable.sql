-- 修复 referrerId 字段的 NOT NULL 约束
-- 允许 referrerId 为空，因为用户可能没有推荐人

-- 修改 Referrals 表的 referrerId 字段，允许为空
ALTER TABLE Referrals MODIFY COLUMN referrerId VARCHAR(24) NULL COMMENT '推荐人ID';

-- 更新现有的空字符串记录为 NULL
UPDATE Referrals SET referrerId = NULL WHERE referrerId = '';

-- 修改 Orders 表的 referrerId 字段，允许为空
ALTER TABLE Orders MODIFY COLUMN referrerId VARCHAR(24) NULL COMMENT '推荐人ID';

-- 更新现有的空字符串记录为 NULL
UPDATE Orders SET referrerId = NULL WHERE referrerId = ''; 