-- 创建客服消息表
CREATE TABLE IF NOT EXISTS KefuMessages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    userId INT NOT NULL COMMENT '用户ID',
    userName VARCHAR(100) COMMENT '用户姓名',
    userAvatar VARCHAR(500) COMMENT '用户头像',
    type INT DEFAULT 1 COMMENT '消息类型：1-用户消息，2-客服回复',
    content TEXT NOT NULL COMMENT '消息内容',
    images TEXT COMMENT '图片（JSON数组）',
    status INT DEFAULT 0 COMMENT '状态：0-未读，1-已读，2-已回复',
    replyContent TEXT COMMENT '回复内容',
    replyTime DATETIME COMMENT '回复时间',
    replyUserId INT COMMENT '回复人ID',
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_user_id (userId),
    INDEX idx_status (status),
    INDEX idx_type (type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='客服消息表';

-- 创建常见问题表
CREATE TABLE IF NOT EXISTS Faqs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    question VARCHAR(500) NOT NULL COMMENT '问题',
    answer TEXT NOT NULL COMMENT '答案',
    category VARCHAR(100) NOT NULL COMMENT '分类',
    sort INT DEFAULT 0 COMMENT '排序',
    status INT DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
    viewCount INT DEFAULT 0 COMMENT '查看次数',
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_category (category),
    INDEX idx_status (status),
    INDEX idx_sort (sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='常见问题表';

-- 插入示例FAQ数据
INSERT INTO Faqs (question, answer, category, sort, status) VALUES
('如何预约体检服务？', '您可以在首页选择体检服务，填写个人信息后提交订单，支付成功后即可预约成功。', '预约服务', 1, 1),
('体检前需要注意什么？', '体检前请保持空腹8-12小时，避免剧烈运动，保持良好的睡眠。', '体检须知', 2, 1),
('体检报告多久能出来？', '一般体检报告在3-5个工作日内出具，特殊情况会提前通知。', '体检报告', 3, 1),
('可以取消预约吗？', '可以，在体检前24小时可以免费取消预约，24小时内取消将收取一定费用。', '预约服务', 4, 1),
('体检费用如何支付？', '支持微信支付、支付宝等多种支付方式，支付成功后即可预约成功。', '支付问题', 5, 1),
('体检项目可以自选吗？', '可以，您可以根据需要选择不同的体检套餐或单项检查。', '体检项目', 6, 1),
('体检结果异常怎么办？', '如发现异常，我们会及时通知您，并提供专业的健康建议和进一步检查建议。', '体检报告', 7, 1),
('可以代他人预约吗？', '可以，请填写被预约人的真实信息，并确保信息准确无误。', '预约服务', 8, 1),
('体检当天需要带什么？', '请携带身份证、预约单，穿着宽松舒适的衣服，女性避开生理期。', '体检须知', 9, 1),
('体检后多久能拿到报告？', '一般体检后3-5个工作日可以拿到报告，我们会通过短信通知您。', '体检报告', 10, 1); 