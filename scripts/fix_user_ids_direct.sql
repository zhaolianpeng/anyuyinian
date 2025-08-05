-- 直接修复现有用户的userId
-- 这个脚本可以直接在数据库中执行

-- 步骤1: 检查当前状态
SELECT '=== 检查当前用户状态 ===' as step;
SELECT id, openId, userId, 
       CASE 
           WHEN userId IS NULL OR userId = '' THEN '需要修复'
           ELSE '正常'
       END as status
FROM Users 
ORDER BY id;

-- 步骤2: 为没有userId的用户生成userId
SELECT '=== 开始修复userId ===' as step;
UPDATE Users 
SET userId = CONCAT(
    LPAD(UNIX_TIMESTAMP(), 10, '0'),
    LPAD(FLOOR(RAND() * 1000000), 6, '0'),
    LPAD(FLOOR(RAND() * 1000000), 8, '0')
)
WHERE userId IS NULL OR userId = '';

-- 步骤3: 验证修复结果
SELECT '=== 验证修复结果 ===' as step;
SELECT id, openId, userId, LENGTH(userId) as userId_length,
       CASE 
           WHEN LENGTH(userId) = 24 THEN '✅ 正确'
           ELSE '❌ 错误'
       END as validation
FROM Users 
ORDER BY id;

-- 步骤4: 检查重复userId
SELECT '=== 检查重复userId ===' as step;
SELECT userId, COUNT(*) as count 
FROM Users 
GROUP BY userId 
HAVING COUNT(*) > 1;

-- 步骤5: 最终验证
SELECT '=== 最终验证 ===' as step;
SELECT 
    COUNT(*) as total_users,
    COUNT(userId) as users_with_userid,
    COUNT(CASE WHEN userId IS NOT NULL AND userId != '' AND LENGTH(userId) = 24 THEN 1 END) as valid_userids,
    CASE 
        WHEN COUNT(*) = COUNT(CASE WHEN userId IS NOT NULL AND userId != '' AND LENGTH(userId) = 24 THEN 1 END) 
        THEN '✅ 所有用户都有有效的userId'
        ELSE '❌ 还有用户没有有效的userId'
    END as result
FROM Users; 