-- 为现有用户生成UserId的迁移脚本
-- 注意：这个脚本需要在应用代码中执行，因为需要生成随机字符串

-- 首先检查是否有用户没有userId
-- SELECT COUNT(*) FROM Users WHERE userId IS NULL OR userId = '';

-- 为没有userId的用户生成新的userId
-- 这个更新需要在Go代码中执行，因为需要调用utils.GenerateUserID()
-- UPDATE Users SET userId = 'generated_user_id' WHERE userId IS NULL OR userId = '';

-- 验证所有用户都有userId
-- SELECT COUNT(*) FROM Users WHERE userId IS NULL OR userId = ''; 