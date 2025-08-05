-- 添加管理员相关字段到Users表
ALTER TABLE Users ADD COLUMN isAdmin TINYINT(1) DEFAULT 0 COMMENT '是否为管理员 0-否 1-是';
ALTER TABLE Users ADD COLUMN adminLevel TINYINT(1) DEFAULT 0 COMMENT '管理员级别 0-普通用户 1-一级管理员 2-超级管理员';
ALTER TABLE Users ADD COLUMN adminPassword VARCHAR(255) DEFAULT NULL COMMENT '管理员密码';
ALTER TABLE Users ADD COLUMN adminUsername VARCHAR(50) DEFAULT NULL COMMENT '管理员用户名';
ALTER TABLE Users ADD COLUMN parentAdminId VARCHAR(24) DEFAULT NULL COMMENT '上级管理员ID';
ALTER TABLE Users ADD COLUMN adminCreatedAt TIMESTAMP NULL COMMENT '成为管理员时间';

-- 创建管理员登录记录表
CREATE TABLE IF NOT EXISTS AdminLoginLogs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    adminUserId VARCHAR(24) NOT NULL COMMENT '管理员用户ID',
    loginTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
    loginIp VARCHAR(45) COMMENT '登录IP',
    userAgent TEXT COMMENT '用户代理',
    status TINYINT(1) DEFAULT 1 COMMENT '登录状态 0-失败 1-成功',
    remark TEXT COMMENT '备注',
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_adminUserId (adminUserId),
    INDEX idx_loginTime (loginTime)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员登录记录表';

-- 插入默认超级管理员账号
INSERT INTO Users (userId, openId, nickName, avatarUrl, isAdmin, adminLevel, adminUsername, adminPassword, adminCreatedAt, createdAt, updatedAt) 
VALUES ('admin_super', 'admin_super_openid', '超级管理员', '/images/default-avatar.png', 1, 2, 'anyuyinian', '000000', NOW(), NOW(), NOW())
ON DUPLICATE KEY UPDATE 
    isAdmin = 1, 
    adminLevel = 2, 
    adminUsername = 'anyuyinian', 
    adminPassword = '000000',
    adminCreatedAt = NOW(); 