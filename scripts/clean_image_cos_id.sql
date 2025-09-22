-- 清理 imageCosId 字段中的无效字符
-- 移除文件名末尾的制表符、空格等无效字符

UPDATE ServiceItems 
SET imageCosId = TRIM(REPLACE(REPLACE(REPLACE(imageCosId, '\t', ''), '\n', ''), '\r', ''))
WHERE imageCosId IS NOT NULL 
  AND imageCosId != ''
  AND (imageCosId LIKE '%\t%' OR imageCosId LIKE '%\n%' OR imageCosId LIKE '%\r%');

-- 查看清理结果
SELECT 
  id, 
  name, 
  imageCosId,
  LENGTH(imageCosId) as cos_id_length,
  CHAR_LENGTH(imageCosId) as char_length
FROM ServiceItems 
WHERE imageCosId IS NOT NULL 
  AND imageCosId != ''
ORDER BY id;
