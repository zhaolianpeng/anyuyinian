-- 添加智慧养老服务类型
-- 迁移文件：026_add_smart_elderly_care_services.sql

-- 插入智慧养老设备服务
INSERT INTO ServiceItems (name, description, category, price, originalPrice, imageUrl, detailImages, formConfig, status, sort, createdAt, updatedAt) VALUES
('智能血压计', '24小时血压监测，数据同步到手机APP，异常自动报警', '智慧养老', 299.00, 399.00, '/static/smart-elderly/blood-pressure-monitor.jpg', '["/static/smart-elderly/bp1.jpg","/static/smart-elderly/bp2.jpg","/static/smart-elderly/bp3.jpg"]', '{"fields":[{"name":"deliveryAddress","label":"收货地址","type":"textarea","required":true,"placeholder":"请输入详细收货地址"}]}', 1, 1, NOW(), NOW()),

('智能血糖仪', '无痛采血，快速检测，数据云端存储，支持家人远程查看', '智慧养老', 199.00, 299.00, '/static/smart-elderly/glucose-monitor.jpg', '["/static/smart-elderly/gm1.jpg","/static/smart-elderly/gm2.jpg","/static/smart-elderly/gm3.jpg"]', '{"fields":[{"name":"deliveryAddress","label":"收货地址","type":"textarea","required":true,"placeholder":"请输入详细收货地址"}]}', 1, 2, NOW(), NOW()),

('智能手环', '心率监测、睡眠分析、跌倒检测、GPS定位，紧急求助一键呼叫', '智慧养老', 399.00, 599.00, '/static/smart-elderly/smart-bracelet.jpg', '["/static/smart-elderly/sb1.jpg","/static/smart-elderly/sb2.jpg","/static/smart-elderly/sb3.jpg"]', '{"fields":[{"name":"deliveryAddress","label":"收货地址","type":"textarea","required":true,"placeholder":"请输入详细收货地址"}]}', 1, 3, NOW(), NOW()),

('智能药盒', '定时提醒服药，记录服药情况，支持家人远程监督', '智慧养老', 159.00, 259.00, '/static/smart-elderly/smart-medicine-box.jpg', '["/static/smart-elderly/smb1.jpg","/static/smart-elderly/smb2.jpg","/static/smart-elderly/smb3.jpg"]', '{"fields":[{"name":"deliveryAddress","label":"收货地址","type":"textarea","required":true,"placeholder":"请输入详细收货地址"}]}', 1, 4, NOW(), NOW()),

('智能摄像头', '远程监控，异常行为识别，双向语音通话，夜视功能', '智慧养老', 299.00, 399.00, '/static/smart-elderly/smart-camera.jpg', '["/static/smart-elderly/sc1.jpg","/static/smart-elderly/sc2.jpg","/static/smart-elderly/sc3.jpg"]', '{"fields":[{"name":"deliveryAddress","label":"收货地址","type":"textarea","required":true,"placeholder":"请输入详细收货地址"}]}', 1, 5, NOW(), NOW()),

('智能门锁', '指纹识别、密码开锁、远程授权、异常开锁报警', '智慧养老', 899.00, 1299.00, '/static/smart-elderly/smart-lock.jpg', '["/static/smart-elderly/sl1.jpg","/static/smart-elderly/sl2.jpg","/static/smart-elderly/sl3.jpg"]', '{"fields":[{"name":"deliveryAddress","label":"收货地址","type":"textarea","required":true,"placeholder":"请输入详细收货地址"}]}', 1, 6, NOW(), NOW()),

('智能马桶', '自动冲洗、座圈加热、健康检测、防跌倒设计', '智慧养老', 2999.00, 3999.00, '/static/smart-elderly/smart-toilet.jpg', '["/static/smart-elderly/st1.jpg","/static/smart-elderly/st2.jpg","/static/smart-elderly/st3.jpg"]', '{"fields":[{"name":"deliveryAddress","label":"收货地址","type":"textarea","required":true,"placeholder":"请输入详细收货地址"}]}', 1, 7, NOW(), NOW()),

('智能床垫', '睡眠质量监测、呼吸监测、离床报警、防褥疮设计', '智慧养老', 1999.00, 2999.00, '/static/smart-elderly/smart-mattress.jpg', '["/static/smart-elderly/sm1.jpg","/static/smart-elderly/sm2.jpg","/static/smart-elderly/sm3.jpg"]', '{"fields":[{"name":"deliveryAddress","label":"收货地址","type":"textarea","required":true,"placeholder":"请输入详细收货地址"}]}', 1, 8, NOW(), NOW());

-- 为智慧养老设备添加视频URL字段（如果ServiceItems表没有videoUrl字段，需要先添加）
-- 注意：这里假设ServiceItems表已经有videoUrl字段，如果没有需要先执行ALTER TABLE语句

-- 更新智慧养老设备的视频URL
UPDATE ServiceItems SET videoUrl = 'https://7072-prod-5g94mx7a3d07e78c-1353115175.cos.ap-shanghai.myqcloud.com/static/videos/smart-elderly-intro.mp4' WHERE category = '智慧养老';
