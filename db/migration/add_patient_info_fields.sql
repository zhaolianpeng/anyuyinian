-- 为订单表添加患者信息相关字段
-- 既往病史字段
ALTER TABLE Orders ADD COLUMN diseaseInfo TEXT COMMENT '既往病史';

-- 是否需要助排二便字段
ALTER TABLE Orders ADD COLUMN needToiletAssist TINYINT DEFAULT 0 COMMENT '是否需要助排二便：0-不需要，1-需要';

-- 添加索引
CREATE INDEX idx_orders_need_toilet_assist ON Orders(needToiletAssist); 