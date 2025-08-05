-- 更新导航数据脚本
-- 将导航链接更新为实际存在的页面

-- 清空现有导航数据
DELETE FROM Navigations;

-- 插入新的导航数据，匹配实际返回的数据
INSERT INTO Navigations (name, icon, linkUrl, sort, status) VALUES
('服务预约', 'https://i.postimg.cc/p5W10Vw7/fuwuyuyue-logo.png', '/pages/service/list', 1, 1),
('我的订单', 'https://i.postimg.cc/phY6Y56z/wodedingdan-logo.png', '/pages/order/list', 2, 1),
('医院信息', 'https://i.postimg.cc/BLGzKBMJ/yiyuanxinxi-logo.png', '/pages/hospital/list', 3, 1),
('个人中心', 'https://i.postimg.cc/XZ5hqzK2/gerenzhongxin-logo.png', '/pages/user/profile', 4, 1);

-- 验证更新结果
SELECT * FROM Navigations ORDER BY sort; 