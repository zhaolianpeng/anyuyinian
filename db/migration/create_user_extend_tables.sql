-- 创建配置表
CREATE TABLE IF NOT EXISTS Configs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    config_key VARCHAR(100) UNIQUE NOT NULL COMMENT '配置键',
    value TEXT NOT NULL COMMENT '配置值',
    description VARCHAR(500) COMMENT '配置描述',
    type VARCHAR(20) DEFAULT 'string' COMMENT '配置类型：string, number, boolean, json',
    status INT DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_key (config_key),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='平台配置表';

-- 创建用户地址表
CREATE TABLE IF NOT EXISTS UserAddresses (
    id INT AUTO_INCREMENT PRIMARY KEY,
    userId INT NOT NULL COMMENT '用户ID',
    name VARCHAR(100) NOT NULL COMMENT '收货人姓名',
    phone VARCHAR(20) NOT NULL COMMENT '收货人电话',
    province VARCHAR(50) COMMENT '省份',
    city VARCHAR(50) COMMENT '城市',
    district VARCHAR(50) COMMENT '区县',
    address VARCHAR(500) NOT NULL COMMENT '详细地址',
    isDefault INT DEFAULT 0 COMMENT '是否默认地址：1-是，0-否',
    status INT DEFAULT 1 COMMENT '状态：1-正常，0-删除',
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_user_id (userId),
    INDEX idx_status (status),
    INDEX idx_is_default (isDefault)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户地址表';

-- 创建就诊人表
CREATE TABLE IF NOT EXISTS Patients (
    id INT AUTO_INCREMENT PRIMARY KEY,
    userId INT NOT NULL COMMENT '用户ID',
    name VARCHAR(100) NOT NULL COMMENT '就诊人姓名',
    idCard VARCHAR(20) COMMENT '身份证号',
    phone VARCHAR(20) COMMENT '联系电话',
    gender INT DEFAULT 0 COMMENT '性别：0-未知，1-男，2-女',
    birthday VARCHAR(20) COMMENT '出生日期',
    relation VARCHAR(50) COMMENT '与用户关系：本人、父亲、母亲等',
    isDefault INT DEFAULT 0 COMMENT '是否默认就诊人：1-是，0-否',
    status INT DEFAULT 1 COMMENT '状态：1-正常，0-删除',
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_user_id (userId),
    INDEX idx_status (status),
    INDEX idx_is_default (isDefault)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='就诊人信息表';

-- 插入示例配置数据
INSERT INTO Configs (config_key, value, description, type, status) VALUES
('customer_service_phone', '400-123-4567', '客服电话', 'string', 1),
('privacy_policy_url', 'https://example.com/privacy', '隐私政策链接', 'string', 1),
('user_agreement_url', 'https://example.com/agreement', '用户协议链接', 'string', 1),
('about_us_url', 'https://example.com/about', '关于我们链接', 'string', 1),
('help_center_url', 'https://example.com/help', '帮助中心链接', 'string', 1),
('app_version', '1.0.0', '应用版本号', 'string', 1),
('force_update', 'false', '是否强制更新', 'boolean', 1),
('maintenance_mode', 'false', '是否维护模式', 'boolean', 1),
('maintenance_message', '系统维护中，请稍后再试', 'string', 'string', 1);
