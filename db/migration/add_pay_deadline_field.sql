-- 为订单表添加支付截止时间字段
ALTER TABLE Orders ADD COLUMN payDeadline DATETIME COMMENT '支付截止时间';

-- 为现有订单设置支付截止时间（创建时间+30分钟）
UPDATE Orders SET payDeadline = DATE_ADD(createdAt, INTERVAL 30 MINUTE) WHERE payDeadline IS NULL;

-- 为payDeadline字段添加索引以提高查询性能
CREATE INDEX idx_pay_deadline ON Orders(payDeadline);

-- 为状态和支付截止时间创建复合索引
CREATE INDEX idx_status_pay_deadline ON Orders(status, payDeadline); 