-- 为现有用户生成推广码
-- 这个脚本需要配合应用代码执行，因为需要生成随机字符串

-- 首先确保推广码字段存在
ALTER TABLE Referrals ADD COLUMN IF NOT EXISTS promoterCode VARCHAR(6) UNIQUE COMMENT '六位推广码' AFTER referrerId;

-- 为没有推广码的记录生成推广码
-- 注意：这个更新需要在应用代码中处理，因为需要生成随机字符串
-- 可以使用以下API接口来生成：
-- POST /api/promoter/generate_codes

-- 添加索引（如果不存在）
CREATE INDEX IF NOT EXISTS idx_promoter_code ON Referrals(promoterCode);

-- 验证推广码格式
-- 确保所有推广码都是6位字符，只包含大写字母和数字
SELECT 
    id,
    userId,
    promoterCode,
    CASE 
        WHEN LENGTH(promoterCode) != 6 THEN '长度错误'
        WHEN promoterCode REGEXP '^[A-Z0-9]{6}$' = 0 THEN '格式错误'
        ELSE '格式正确'
    END as validation_result
FROM Referrals 
WHERE promoterCode IS NOT NULL AND promoterCode != ''; 