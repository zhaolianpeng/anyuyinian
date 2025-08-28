-- 创建咨询会话表
CREATE TABLE IF NOT EXISTS `consultations` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` varchar(100) NOT NULL COMMENT '用户ID',
  `user_name` varchar(100) NOT NULL COMMENT '用户姓名',
  `user_phone` varchar(20) DEFAULT NULL COMMENT '用户手机号',
  `status` varchar(20) NOT NULL DEFAULT 'waiting' COMMENT '状态：waiting-等待回复,chatting-咨询中,closed-已结束',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `last_message` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后一条消息时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_last_message` (`last_message`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='咨询会话表';

-- 创建咨询消息表
CREATE TABLE IF NOT EXISTS `consultation_messages` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `consultation_id` int(11) NOT NULL COMMENT '咨询会话ID',
  `sender_type` varchar(20) NOT NULL COMMENT '发送者类型：user-用户,admin-管理员',
  `content` text NOT NULL COMMENT '消息内容',
  `is_read` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否已读',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_consultation_id` (`consultation_id`),
  KEY `idx_sender_type` (`sender_type`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `fk_consultation_messages_consultation` FOREIGN KEY (`consultation_id`) REFERENCES `consultations` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='咨询消息表';

-- 创建咨询通知表
CREATE TABLE IF NOT EXISTS `consultation_notifications` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `consultation_id` int(11) NOT NULL COMMENT '咨询会话ID',
  `type` varchar(20) NOT NULL COMMENT '通知类型：new_message-新消息,status_change-状态变更',
  `title` varchar(200) NOT NULL COMMENT '通知标题',
  `content` text COMMENT '通知内容',
  `is_read` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否已读',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_consultation_id` (`consultation_id`),
  KEY `idx_type` (`type`),
  KEY `idx_is_read` (`is_read`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `fk_consultation_notifications_consultation` FOREIGN KEY (`consultation_id`) REFERENCES `consultations` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='咨询通知表';

-- 插入一些测试数据
INSERT INTO `consultations` (`user_id`, `user_name`, `user_phone`, `status`) VALUES
('test_user_001', '张三', '13800138001', 'waiting'),
('test_user_002', '李四', '13800138002', 'chatting'),
('test_user_003', '王五', '13800138003', 'closed');

-- 插入测试消息
INSERT INTO `consultation_messages` (`consultation_id`, `sender_type`, `content`) VALUES
(1, 'admin', '您好！欢迎使用在线咨询服务，请问有什么可以帮助您的吗？'),
(1, 'user', '我想了解一下护工服务的价格'),
(2, 'admin', '您好！欢迎使用在线咨询服务，请问有什么可以帮助您的吗？'),
(2, 'user', '护工服务的时间安排是怎样的？'),
(2, 'admin', '护工服务时间比较灵活，可以根据您的需求来安排'),
(3, 'admin', '您好！欢迎使用在线咨询服务，请问有什么可以帮助您的吗？'),
(3, 'user', '谢谢，我已经了解了');

-- 插入测试通知
INSERT INTO `consultation_notifications` (`consultation_id`, `type`, `title`, `content`) VALUES
(1, 'new_message', '新用户消息', '用户 张三 发送新消息: 我想了解一下护工服务的价格'),
(2, 'new_message', '新用户消息', '用户 李四 发送新消息: 护工服务的时间安排是怎样的？'),
(3, 'status_change', '咨询会话已关闭', '用户 王五 的咨询会话已关闭');
