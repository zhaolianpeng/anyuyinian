-- 创建轮播图表
CREATE TABLE IF NOT EXISTS Banners (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(200) COMMENT '轮播图标题',
    imageUrl TEXT NOT NULL COMMENT '轮播图图片URL',
    linkUrl VARCHAR(500) COMMENT '轮播图链接URL',
    sort INT DEFAULT 0 COMMENT '排序',
    status INT DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_status_sort (status, sort),
    INDEX idx_sort (sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='轮播图表';

-- 创建导航表
CREATE TABLE IF NOT EXISTS Navigations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL COMMENT '导航名称',
    icon VARCHAR(200) NOT NULL COMMENT '导航图标',
    linkUrl VARCHAR(500) COMMENT '导航链接URL',
    sort INT DEFAULT 0 COMMENT '排序',
    status INT DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_status_sort (status, sort),
    INDEX idx_sort (sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='导航表';

-- 创建服务项表
CREATE TABLE IF NOT EXISTS Services (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL COMMENT '服务名称',
    description TEXT COMMENT '服务描述',
    icon VARCHAR(200) NOT NULL COMMENT '服务图标',
    imageUrl VARCHAR(500) COMMENT '服务图片URL',
    linkUrl VARCHAR(500) COMMENT '服务链接URL',
    sort INT DEFAULT 0 COMMENT '排序',
    status INT DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_status_sort (status, sort),
    INDEX idx_sort (sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='服务项表';

-- 创建医院表
CREATE TABLE IF NOT EXISTS Hospitals (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(200) NOT NULL COMMENT '医院名称',
    logo VARCHAR(500) COMMENT '医院logo',
    address VARCHAR(500) COMMENT '医院地址',
    phone VARCHAR(50) COMMENT '医院电话',
    description TEXT COMMENT '医院描述',
    level VARCHAR(50) COMMENT '医院等级',
    type VARCHAR(50) COMMENT '医院类型',
    longitude DECIMAL(10, 7) COMMENT '经度',
    latitude DECIMAL(10, 7) COMMENT '纬度',
    sort INT DEFAULT 0 COMMENT '排序',
    status INT DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_status_sort (status, sort),
    INDEX idx_location (longitude, latitude),
    INDEX idx_sort (sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='医院表';

-- 插入示例数据

-- 插入轮播图示例数据
INSERT INTO Banners (title, imageUrl, linkUrl, sort, status) VALUES
('健康体检', 'https://example.com/banner1.jpg', '/pages/health-check', 1, 1),
('在线问诊', 'https://example.com/banner2.jpg', '/pages/online-consultation', 2, 1),
('预约挂号', 'https://example.com/banner3.jpg', '/pages/appointment', 3, 1);

-- 插入导航示例数据
INSERT INTO Navigations (name, icon, linkUrl, sort, status) VALUES
('服务预约', '/images/nav/appointment.png', '/pages/service/list', 1, 1),
('我的订单', '/images/nav/consultation.png', '/pages/order/list', 2, 1),
('个人中心', '/images/nav/health-record.png', '/pages/user/profile', 3, 1),
('客服咨询', '/images/nav/report.png', '/pages/kefu/chat', 4, 1),
('推荐中心', '/images/nav/medicine.png', '/pages/referral/center', 5, 1),
('健康资讯', '/images/nav/news.png', '/pages/webview/webview?url=https://example.com/news', 6, 1),
('附近医院', '/images/nav/hospital.png', '/pages/hospital/list', 7, 1);

-- 插入服务项示例数据
INSERT INTO Services (name, description, icon, imageUrl, linkUrl, sort, status) VALUES
('预约挂号', '快速预约专家门诊', '/images/service/appointment.png', '/images/service/appointment-bg.jpg', '/pages/appointment', 1, 1),
('在线问诊', '足不出户看专家', '/images/service/consultation.png', '/images/service/consultation-bg.jpg', '/pages/consultation', 2, 1),
('健康体检', '专业体检服务', '/images/service/checkup.png', '/images/service/checkup-bg.jpg', '/pages/checkup', 3, 1),
('检查报告', '查看检查结果', '/images/service/report.png', '/images/service/report-bg.jpg', '/pages/report', 4, 1),
('药品配送', '药品配送到家', '/images/service/medicine.png', '/images/service/medicine-bg.jpg', '/pages/medicine', 5, 1),
('健康档案', '个人健康管理', '/images/service/record.png', '/images/service/record-bg.jpg', '/pages/record', 6, 1);

-- 插入医院示例数据
INSERT INTO Hospitals (name, logo, address, phone, description, level, type, longitude, latitude, sort, status) VALUES
('深圳市人民医院', '/images/hospital/rmyy-logo.png', '深圳市罗湖区东门北路1017号', '0755-25533018', '深圳市人民医院是一所集医疗、教学、科研、预防、保健为一体的综合性三级甲等医院', '三级甲等', '综合医院', 114.123456, 22.654321, 1, 1),
('深圳市第二人民医院', '/images/hospital/dermyy-logo.png', '深圳市福田区笋岗西路3002号', '0755-83366388', '深圳市第二人民医院是一所现代化综合性三级甲等医院', '三级甲等', '综合医院', 114.234567, 22.765432, 2, 1),
('深圳市中医院', '/images/hospital/zyy-logo.png', '深圳市福田区福华路1号', '0755-83000111', '深圳市中医院是一所集医疗、教学、科研、预防、保健为一体的综合性三级甲等中医院', '三级甲等', '中医院', 114.345678, 22.876543, 3, 1),
('深圳市儿童医院', '/images/hospital/etyy-logo.png', '深圳市福田区益田路7019号', '0755-83009888', '深圳市儿童医院是一所集医疗、教学、科研、预防、保健为一体的三级甲等儿童专科医院', '三级甲等', '专科医院', 114.456789, 22.987654, 4, 1),
('深圳市妇幼保健院', '/images/hospital/fybjy-logo.png', '深圳市福田区红荔路2004号', '0755-83000111', '深圳市妇幼保健院是一所集医疗、教学、科研、预防、保健为一体的三级甲等妇幼保健院', '三级甲等', '专科医院', 114.567890, 23.098765, 5, 1); 