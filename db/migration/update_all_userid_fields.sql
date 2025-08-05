-- 将所有表的userId字段从INT改为VARCHAR(24)
-- 这个脚本需要在应用代码中执行，因为需要生成随机字符串

-- 1. 更新Orders表的userId字段
ALTER TABLE Orders MODIFY COLUMN userId VARCHAR(24) NOT NULL;

-- 2. 更新KefuMessages表的userId字段
ALTER TABLE KefuMessages MODIFY COLUMN userId VARCHAR(24) NOT NULL;
ALTER TABLE KefuMessages MODIFY COLUMN replyUserId VARCHAR(24);

-- 3. 更新Referrals表的userId字段
ALTER TABLE Referrals MODIFY COLUMN userId VARCHAR(24) NOT NULL;

-- 4. 更新Commissions表的userId字段
ALTER TABLE Commissions MODIFY COLUMN userId VARCHAR(24) NOT NULL;

-- 5. 更新Cashouts表的userId字段
ALTER TABLE Cashouts MODIFY COLUMN userId VARCHAR(24) NOT NULL;

-- 6. 更新Files表的userId字段
ALTER TABLE Files MODIFY COLUMN userId VARCHAR(24);

-- 7. 更新UserAddresses表的userId字段（如果还没有更新）
ALTER TABLE UserAddresses MODIFY COLUMN userId VARCHAR(24) NOT NULL;

-- 8. 更新Patients表的userId字段（如果还没有更新）
ALTER TABLE Patients MODIFY COLUMN userId VARCHAR(24) NOT NULL;

-- 添加索引
CREATE INDEX idx_user_id ON Orders(userId);
CREATE INDEX idx_user_id ON KefuMessages(userId);
CREATE INDEX idx_user_id ON Referrals(userId);
CREATE INDEX idx_user_id ON Commissions(userId);
CREATE INDEX idx_user_id ON Cashouts(userId);
CREATE INDEX idx_user_id ON Files(userId);

-- 注意：这个脚本需要在Go代码中执行，因为需要为现有数据生成新的userId
-- 现有数据的userId需要从数字ID转换为MongoDB风格的字符串ID 