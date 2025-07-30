# 配置说明

## 环境变量配置

在运行服务前，请设置以下环境变量：

### 微信小程序配置

```bash
# 设置微信小程序AppID和AppSecret
export WX_APP_ID="你的微信小程序AppID"
export WX_APP_SECRET="你的微信小程序AppSecret"
```

### 数据库配置

确保你的数据库配置正确，可以参考 `db/init.go` 文件中的配置。

## 数据库表创建

### 用户表

运行以下SQL语句创建用户表：

```sql
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
```

### 首页相关表

运行以下SQL语句创建首页相关表：

```sql
-- 创建轮播图表
CREATE TABLE IF NOT EXISTS Banners (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(200) COMMENT '轮播图标题',
    imageUrl TEXT NOT NULL COMMENT '轮播图图片URL',
    linkUrl VARCHAR(500) COMMENT '轮播图链接URL',
    sort INT DEFAULT 0 COMMENT '排序',
    status INT DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_status_sort (status, sort),
    INDEX idx_sort (sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='轮播图表';

-- 创建导航表
CREATE TABLE IF NOT EXISTS Navigations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL COMMENT '导航名称',
    icon VARCHAR(200) NOT NULL COMMENT '导航图标',
    linkUrl VARCHAR(500) COMMENT '导航链接URL',
    sort INT DEFAULT 0 COMMENT '排序',
    status INT DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_status_sort (status, sort),
    INDEX idx_sort (sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='导航表';

-- 创建服务项表
CREATE TABLE IF NOT EXISTS Services (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL COMMENT '服务名称',
    description TEXT COMMENT '服务描述',
    icon VARCHAR(200) NOT NULL COMMENT '服务图标',
    imageUrl VARCHAR(500) COMMENT '服务图片URL',
    linkUrl VARCHAR(500) COMMENT '服务链接URL',
    sort INT DEFAULT 0 COMMENT '排序',
    status INT DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_status_sort (status, sort),
    INDEX idx_sort (sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='服务项表';

-- 创建医院表
CREATE TABLE IF NOT EXISTS Hospitals (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(200) NOT NULL COMMENT '医院名称',
    logo VARCHAR(500) COMMENT '医院logo',
    address VARCHAR(500) COMMENT '医院地址',
    phone VARCHAR(50) COMMENT '医院电话',
    description TEXT COMMENT '医院描述',
    level VARCHAR(50) COMMENT '医院等级',
    type VARCHAR(50) COMMENT '医院类型',
    longitude DECIMAL(10, 7) COMMENT '经度',
    latitude DECIMAL(10, 7) COMMENT '纬度',
    sort INT DEFAULT 0 COMMENT '排序',
    status INT DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_status_sort (status, sort),
    INDEX idx_location (longitude, latitude),
    INDEX idx_sort (sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='医院表';
```

### 文件上传表

运行以下SQL语句创建文件上传表：

```sql
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
```

### 用户扩展表

运行以下SQL语句创建用户扩展相关表：

```sql
-- 创建配置表
CREATE TABLE IF NOT EXISTS Configs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    `key` VARCHAR(100) UNIQUE NOT NULL COMMENT '配置键',
    value TEXT NOT NULL COMMENT '配置值',
    description VARCHAR(500) COMMENT '配置描述',
    type VARCHAR(20) DEFAULT 'string' COMMENT '配置类型：string, number, boolean, json',
    status INT DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_key (key),
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
```

### 服务、订单、推荐表

运行以下SQL语句创建服务、订单、推荐相关表：

```sql
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
    userId INT NOT NULL COMMENT '用户ID',
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
    referrerId INT COMMENT '推荐人ID',
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
    userId INT UNIQUE NOT NULL COMMENT '用户ID',
    referrerId INT NOT NULL COMMENT '推荐人ID',
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
    userId INT NOT NULL COMMENT '用户ID',
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
    userId INT NOT NULL COMMENT '用户ID',
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
```

## 运行服务

1. 设置环境变量
2. 创建数据库表
3. 运行服务：

```bash
go run main.go
```

服务将在 `:80` 端口启动。

## 接口测试

所有测试脚本已拆分到 `tests/` 目录下，可以使用以下命令进行测试：

### 给测试脚本添加执行权限
```bash
chmod +x tests/*.sh
```

### 运行单个测试脚本
```bash
# 测试微信登录接口
./tests/test_wx_login.sh

# 测试首页初始化接口
./tests/test_home_init.sh

# 测试文件上传接口
./tests/test_upload.sh

# 测试用户信息接口
./tests/test_user_info.sh

# 测试用户地址管理接口
./tests/test_user_address.sh

# 测试用户就诊人管理接口
./tests/test_user_patient.sh

# 测试服务相关接口
./tests/test_service.sh

# 测试订单相关接口
./tests/test_order.sh

# 测试推荐系统接口
./tests/test_referral.sh

# 测试客服医院接口
./tests/test_kefu_hospital.sh
```

### 运行所有测试脚本
```bash
for script in tests/*.sh; do
    echo "运行测试脚本: $script"
    ./$script
    echo "----------------------------------------"
done
```

注意：测试脚本中的code是模拟的，实际使用时需要从小程序获取真实的code。

## 文档说明

所有接口文档已拆分到 `docs/` 目录下：

- `docs/README.md` - 接口文档总览
- `docs/wx_login_api.md` - 微信登录接口文档
- `docs/home_init_api.md` - 首页初始化接口文档
- `docs/upload_api.md` - 文件上传接口文档
- `docs/config_api.md` - 平台配置接口文档
- `docs/user_apis.md` - 用户信息相关接口文档
- `docs/service_apis.md` - 服务相关接口文档
- `docs/order_apis.md` - 订单相关接口文档
- `docs/referral_apis.md` - 推荐系统接口文档
- `docs/kefu_hospital_apis.md` - 客服、医院相关接口文档 