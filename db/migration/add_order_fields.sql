-- 为订单表添加新字段
ALTER TABLE Orders ADD COLUMN patientId INT NOT NULL DEFAULT 0 COMMENT '就诊人ID';
ALTER TABLE Orders ADD COLUMN addressId INT NOT NULL DEFAULT 0 COMMENT '地址ID';
ALTER TABLE Orders ADD COLUMN appointmentDate VARCHAR(20) NOT NULL DEFAULT '' COMMENT '预约日期';
ALTER TABLE Orders ADD COLUMN appointmentTime VARCHAR(20) NOT NULL DEFAULT '' COMMENT '预约时间';

-- 添加索引
CREATE INDEX idx_orders_patient_id ON Orders(patientId);
CREATE INDEX idx_orders_address_id ON Orders(addressId);
CREATE INDEX idx_orders_appointment_date ON Orders(appointmentDate); 