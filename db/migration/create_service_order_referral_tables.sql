-- 创建服务项目表
CREATE TABLE IF NOT EXISTS ServiceItems (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(200) NOT NULL COMMENT '服务名称',
    description TEXT COMMENT '服务描述',
    category VARCHAR(100) NOT NULL COMMENT '服务分类',
    price DECIMAL(10,2) NOT NULL COMMENT '服务价格',
    originalPrice DECIMAL(10,2) COMMENT '原价',
    imageUrl VARCHAR(500) COMMENT '服务图片',
    detailImages TEXT COMMENT '详情图片（JSON数组）',
    formConfig TEXT COMMENT '表单配置（JSON格式）',
    status INT DEFAULT 1 COMMENT '状态：1-上架，0-下架',
    sort INT DEFAULT 0 COMMENT '排序',
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_category (category),
    INDEX idx_status (status),
    INDEX idx_sort (sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='服务项目表';

-- 创建订单表
CREATE TABLE IF NOT EXISTS Orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    orderNo VARCHAR(50) UNIQUE NOT NULL COMMENT '订单号',
    userId VARCHAR(24) NOT NULL COMMENT '用户ID',
    serviceId INT NOT NULL COMMENT '服务ID',
    serviceName VARCHAR(200) NOT NULL COMMENT '服务名称',
    price DECIMAL(10,2) NOT NULL COMMENT '单价',
    quantity INT DEFAULT 1 COMMENT '数量',
    totalAmount DECIMAL(10,2) NOT NULL COMMENT '总金额',
    formData TEXT COMMENT '表单数据（JSON格式）',
    status INT DEFAULT 0 COMMENT '订单状态：0-待支付，1-已支付，2-已完成，3-已取消，4-已退款',
    payStatus INT DEFAULT 0 COMMENT '支付状态：0-未支付，1-已支付',
    payTime DATETIME COMMENT '支付时间',
    payMethod VARCHAR(20) COMMENT '支付方式',
    transactionId VARCHAR(100) COMMENT '第三方支付交易号',
    refundStatus INT DEFAULT 0 COMMENT '退款状态：0-未退款，1-退款中，2-已退款',
    refundTime DATETIME COMMENT '退款时间',
    refundAmount DECIMAL(10,2) COMMENT '退款金额',
    refundReason VARCHAR(500) COMMENT '退款原因',
    remark VARCHAR(500) COMMENT '备注',
    referrerId VARCHAR(24) COMMENT '推荐人ID',
    commission DECIMAL(10,2) DEFAULT 0 COMMENT '佣金金额',
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_user_id (userId),
    INDEX idx_order_no (orderNo),
    INDEX idx_status (status),
    INDEX idx_pay_status (payStatus),
    INDEX idx_referrer_id (referrerId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单表';

-- 创建推荐关系表
CREATE TABLE IF NOT EXISTS Referrals (
    id INT AUTO_INCREMENT PRIMARY KEY,
    userId VARCHAR(24) UNIQUE NOT NULL COMMENT '用户ID',
    referrerId VARCHAR(24) NOT NULL COMMENT '推荐人ID',
    qrCodeUrl VARCHAR(500) COMMENT '专属二维码URL',
    status INT DEFAULT 1 COMMENT '状态：1-正常，0-禁用',
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_user_id (userId),
    INDEX idx_referrer_id (referrerId),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='推荐关系表';

-- 创建佣金记录表
CREATE TABLE IF NOT EXISTS Commissions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    userId VARCHAR(24) NOT NULL COMMENT '用户ID',
    orderId INT NOT NULL COMMENT '订单ID',
    orderNo VARCHAR(50) NOT NULL COMMENT '订单号',
    amount DECIMAL(10,2) NOT NULL COMMENT '佣金金额',
    rate DECIMAL(5,4) NOT NULL COMMENT '佣金比例',
    status INT DEFAULT 0 COMMENT '状态：0-待结算，1-已结算，2-已提现',
    cashoutTime DATETIME COMMENT '提现时间',
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_user_id (userId),
    INDEX idx_order_id (orderId),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='佣金记录表';

-- 创建提现记录表
CREATE TABLE IF NOT EXISTS Cashouts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    userId VARCHAR(24) NOT NULL COMMENT '用户ID',
    amount DECIMAL(10,2) NOT NULL COMMENT '提现金额',
    method VARCHAR(20) NOT NULL COMMENT '提现方式：wechat, alipay, bank',
    account VARCHAR(200) NOT NULL COMMENT '提现账户',
    status INT DEFAULT 0 COMMENT '状态：0-待审核，1-已通过，2-已拒绝，3-已到账',
    remark VARCHAR(500) COMMENT '备注',
    processTime DATETIME COMMENT '处理时间',
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_user_id (userId),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='提现记录表';

-- 插入示例服务数据
INSERT INTO ServiceItems (name, description, category, price, originalPrice, imageUrl, formConfig, status, sort) VALUES
('健康体检套餐A', '基础健康体检，包含血常规、尿常规、心电图等检查项目', '体检', 299.00, 399.00, 'https://example.com/images/health_check_a.jpg', '{"fields":[{"name":"name","label":"姓名","type":"text","required":true,"placeholder":"请输入姓名"},{"name":"phone","label":"手机号","type":"text","required":true,"placeholder":"请输入手机号"},{"name":"idCard","label":"身份证号","type":"text","required":true,"placeholder":"请输入身份证号"},{"name":"birthday","label":"出生日期","type":"date","required":true},{"name":"gender","label":"性别","type":"radio","required":true,"options":[{"label":"男","value":"1"},{"label":"女","value":"2"}]}]}', 1, 1),
('健康体检套餐B', '全面健康体检，包含血常规、尿常规、心电图、B超等检查项目', '体检', 599.00, 799.00, 'https://example.com/images/health_check_b.jpg', '{"fields":[{"name":"name","label":"姓名","type":"text","required":true,"placeholder":"请输入姓名"},{"name":"phone","label":"手机号","type":"text","required":true,"placeholder":"请输入手机号"},{"name":"idCard","label":"身份证号","type":"text","required":true,"placeholder":"请输入身份证号"},{"name":"birthday","label":"出生日期","type":"date","required":true},{"name":"gender","label":"性别","type":"radio","required":true,"options":[{"label":"男","value":"1"},{"label":"女","value":"2"}]},{"name":"address","label":"详细地址","type":"textarea","required":true,"placeholder":"请输入详细地址"}]}', 1, 2),
('心理咨询服务', '专业心理咨询，一对一心理疏导，帮助解决心理问题', '心理咨询', 200.00, 300.00, 'https://example.com/images/psychology.jpg', '{"fields":[{"name":"name","label":"姓名","type":"text","required":true,"placeholder":"请输入姓名"},{"name":"phone","label":"手机号","type":"text","required":true,"placeholder":"请输入手机号"},{"name":"age","label":"年龄","type":"text","required":true,"placeholder":"请输入年龄"},{"name":"problem","label":"主要问题","type":"textarea","required":true,"placeholder":"请描述您的主要问题"},{"name":"preferredTime","label":"期望咨询时间","type":"text","required":true,"placeholder":"请输入期望的咨询时间"}]}', 1, 3),
('营养师咨询', '专业营养师一对一咨询，制定个性化营养方案', '营养咨询', 150.00, 200.00, 'https://example.com/images/nutrition.jpg', '{"fields":[{"name":"name","label":"姓名","type":"text","required":true,"placeholder":"请输入姓名"},{"name":"phone","label":"手机号","type":"text","required":true,"placeholder":"请输入手机号"},{"name":"height","label":"身高(cm)","type":"text","required":true,"placeholder":"请输入身高"},{"name":"weight","label":"体重(kg)","type":"text","required":true,"placeholder":"请输入体重"},{"name":"goal","label":"营养目标","type":"select","required":true,"options":[{"label":"减重","value":"weight_loss"},{"label":"增重","value":"weight_gain"},{"label":"健康管理","value":"health_management"}]}]}', 1, 4); 