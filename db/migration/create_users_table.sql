-- 创建用户表
CREATE TABLE IF NOT EXISTS Users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    openId VARCHAR(100) UNIQUE NOT NULL COMMENT '微信openid',
    unionId VARCHAR(100) COMMENT '微信unionid',
    nickName VARCHAR(100) COMMENT '用户昵称',
    avatarUrl TEXT COMMENT '用户头像URL',
    gender INT DEFAULT 0 COMMENT '性别：0-未知，1-男，2-女',
    country VARCHAR(50) COMMENT '国家',
    province VARCHAR(50) COMMENT '省份',
    city VARCHAR(50) COMMENT '城市',
    language VARCHAR(20) COMMENT '语言',
    sessionKey VARCHAR(100) COMMENT '微信session_key',
    lastLoginAt DATETIME COMMENT '最后登录时间',
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_openid (openId),
    INDEX idx_unionid (unionId),
    INDEX idx_last_login (lastLoginAt)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='微信小程序用户表'; 