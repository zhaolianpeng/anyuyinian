-- 创建文件表
CREATE TABLE IF NOT EXISTS Files (
    id INT AUTO_INCREMENT PRIMARY KEY,
    fileName VARCHAR(255) NOT NULL COMMENT '文件名',
    originalName VARCHAR(255) NOT NULL COMMENT '原始文件名',
    filePath VARCHAR(500) NOT NULL COMMENT '文件路径',
    fileUrl VARCHAR(500) NOT NULL COMMENT '文件访问URL',
    fileSize BIGINT COMMENT '文件大小（字节）',
    fileType VARCHAR(50) COMMENT '文件类型',
    mimeType VARCHAR(100) COMMENT 'MIME类型',
    category VARCHAR(50) COMMENT '文件分类',
    description TEXT COMMENT '文件描述',
    userId INT COMMENT '上传用户ID',
    status INT DEFAULT 1 COMMENT '状态：1-正常，0-删除',
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_user_id (userId),
    INDEX idx_category (category),
    INDEX idx_status (status),
    INDEX idx_created_at (createdAt),
    INDEX idx_file_type (fileType)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件上传表'; 