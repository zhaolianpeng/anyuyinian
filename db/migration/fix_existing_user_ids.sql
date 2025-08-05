-- 修复现有用户的userId字段
-- 为没有userId的用户生成MongoDB风格的userId

-- 1. 检查当前状态
SELECT '当前用户状态:' as info;
SELECT id, openId, userId FROM Users ORDER BY id;

-- 2. 为没有userId的用户生成userId
-- 使用时间戳和随机数生成24位字符串
UPDATE Users 
SET userId = CONCAT(
    LPAD(UNIX_TIMESTAMP(), 10, '0'),
    LPAD(FLOOR(RAND() * 1000000), 6, '0'),
    LPAD(FLOOR(RAND() * 1000000), 8, '0')
)
WHERE userId IS NULL OR userId = '';

-- 3. 验证更新结果
SELECT '更新后的用户状态:' as info;
SELECT id, openId, userId, LENGTH(userId) as userId_length FROM Users ORDER BY id;

-- 4. 检查是否有重复的userId
SELECT '检查重复userId:' as info;
SELECT userId, COUNT(*) as count 
FROM Users 
GROUP BY userId 
HAVING COUNT(*) > 1;

-- 5. 如果发现重复，为重复的用户重新生成userId
-- 这里需要手动处理重复的情况
-- 可以运行以下查询来查看重复的用户：
-- SELECT * FROM Users WHERE userId IN (
--     SELECT userId FROM Users GROUP BY userId HAVING COUNT(*) > 1
-- );

-- 6. 确保userId字段有正确的约束
-- 如果还没有添加UNIQUE约束，可以运行：
-- ALTER TABLE Users ADD UNIQUE INDEX idx_user_id_unique (userId);

-- 7. 验证所有用户都有userId
SELECT '验证结果:' as info;
SELECT 
    COUNT(*) as total_users,
    COUNT(userId) as users_with_userid,
    COUNT(CASE WHEN userId IS NOT NULL AND userId != '' THEN 1 END) as valid_userids
FROM Users; 