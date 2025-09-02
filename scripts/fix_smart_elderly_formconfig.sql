-- 修复智慧养老设备formConfig脚本
-- 检查并修复智慧养老设备的表单配置

USE anyuyinian;

-- 1. 检查当前智慧养老设备的formConfig
SELECT 
    id,
    name,
    category,
    formConfig,
    LENGTH(formConfig) as config_length
FROM ServiceItems 
WHERE category = '智慧养老' 
AND status = 1;

-- 2. 更新智慧养老设备的formConfig（如果为空或无效）
UPDATE ServiceItems 
SET formConfig = '{"fields":[{"name":"deliveryAddress","label":"收货地址","type":"textarea","required":true,"placeholder":"请输入详细收货地址"}]}'
WHERE category = '智慧养老' 
AND status = 1 
AND (formConfig IS NULL OR formConfig = '' OR formConfig = 'null' OR formConfig = '{}');

-- 3. 验证更新结果
SELECT 
    id,
    name,
    category,
    formConfig,
    LENGTH(formConfig) as config_length
FROM ServiceItems 
WHERE category = '智慧养老' 
AND status = 1;

-- 4. 检查所有服务的formConfig状态
SELECT 
    category,
    COUNT(*) as total_count,
    SUM(CASE WHEN formConfig IS NULL OR formConfig = '' OR formConfig = 'null' THEN 1 ELSE 0 END) as empty_config_count,
    SUM(CASE WHEN formConfig IS NOT NULL AND formConfig != '' AND formConfig != 'null' THEN 1 ELSE 0 END) as valid_config_count
FROM ServiceItems 
WHERE status = 1
GROUP BY category
ORDER BY category;
