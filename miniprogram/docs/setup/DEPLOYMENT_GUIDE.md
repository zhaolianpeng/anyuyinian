# 后端服务部署指南

## 概述

当前小程序已经统一切换到自有服务端直连，生产环境请以自有服务器 + Nginx 反向代理部署为准。

## 部署选项

### 1. 云服务器部署（推荐）

#### 1.1 购买云服务器
- 阿里云、腾讯云、华为云等
- 建议配置：2核4G内存，带宽5Mbps以上
- 操作系统：Ubuntu 20.04 或 CentOS 8

#### 1.2 服务器环境准备
```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装Go环境
wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 安装MySQL
sudo apt install mysql-server -y
sudo mysql_secure_installation

# 安装Nginx
sudo apt install nginx -y
```

#### 1.3 部署后端服务
```bash
# 克隆代码
git clone <your-repo-url>
cd anyuyinian

# 编译
go build -o app main.go

# 创建服务文件
sudo tee /etc/systemd/system/anyuyinian.service << EOF
[Unit]
Description=Anyuyinian API Server
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/www/anyuyinian
ExecStart=/var/www/anyuyinian/app
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# 启动服务
sudo systemctl enable anyuyinian
sudo systemctl start anyuyinian
```

#### 1.4 配置Nginx反向代理
```bash
# 创建Nginx配置
sudo tee /etc/nginx/sites-available/anyuyinian << EOF
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        proxy_pass http://localhost:80;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}
EOF

# 启用配置
sudo ln -s /etc/nginx/sites-available/anyuyinian /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

#### 1.5 配置SSL证书
```bash
# 安装Certbot
sudo apt install certbot python3-certbot-nginx -y

# 获取SSL证书
sudo certbot --nginx -d your-domain.com

# 自动续期
sudo crontab -e
# 添加：0 12 * * * /usr/bin/certbot renew --quiet
```

### 2. 容器化部署

#### 2.1 创建Dockerfile
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/app .
EXPOSE 80

CMD ["./app"]
```

#### 2.2 创建docker-compose.yml
```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "80:80"
    environment:
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=root
      - DB_PASSWORD=your_password
      - DB_NAME=anyuyinian
    depends_on:
      - mysql
    restart: always

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: your_password
      MYSQL_DATABASE: anyuyinian
    volumes:
      - mysql_data:/var/lib/mysql
    restart: always

volumes:
  mysql_data:
```

#### 2.3 部署
```bash
# 构建和启动
docker-compose up -d

# 查看日志
docker-compose logs -f
```

### 3. 云函数部署

#### 3.1 腾讯云函数
```bash
# 安装Serverless Framework
npm install -g serverless

# 创建serverless.yml
cat > serverless.yml << EOF
service: anyuyinian-api

provider:
  name: tencent
  runtime: Go1
  region: ap-shanghai
  credentials:
    secretId: YOUR_SECRET_ID
    secretKey: YOUR_SECRET_KEY

functions:
  api:
    handler: app
    events:
      - http:
          path: /{proxy+}
          method: ANY
EOF

# 部署
serverless deploy
```

## 数据库配置

### 1. 本地MySQL配置
```sql
-- 创建数据库
CREATE DATABASE anyuyinian CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建用户
CREATE USER 'anyuyinian'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON anyuyinian.* TO 'anyuyinian'@'localhost';
FLUSH PRIVILEGES;
```

### 2. 云数据库配置
- 使用阿里云RDS、腾讯云CDB等
- 配置安全组允许应用服务器访问
- 设置数据库连接参数

## 域名配置

### 1. 域名解析
```bash
# A记录
your-domain.com -> 服务器IP

# CNAME记录
api.your-domain.com -> your-domain.com
```

### 2. HTTPS证书
- 使用Let's Encrypt免费证书
- 或购买商业SSL证书
- 配置自动续期

## 监控和日志

### 1. 日志配置
```bash
# 创建日志目录
sudo mkdir -p /var/log/anyuyinian
sudo chown www-data:www-data /var/log/anyuyinian

# 配置日志轮转
sudo tee /etc/logrotate.d/anyuyinian << EOF
/var/log/anyuyinian/*.log {
    daily
    missingok
    rotate 7
    compress
    notifempty
    create 644 www-data www-data
}
EOF
```

### 2. 监控配置
```bash
# 安装监控工具
sudo apt install htop iotop -y

# 配置系统监控
sudo apt install prometheus node-exporter -y
```

## 安全配置

### 1. 防火墙配置
```bash
# 配置UFW
sudo ufw allow 22
sudo ufw allow 80
sudo ufw allow 443
sudo ufw enable
```

### 2. 数据库安全
```sql
-- 删除匿名用户
DELETE FROM mysql.user WHERE User='';

-- 删除测试数据库
DROP DATABASE IF EXISTS test;
DELETE FROM mysql.db WHERE Db='test' OR Db='test\\_%';

-- 刷新权限
FLUSH PRIVILEGES;
```

## 测试验证

### 1. API测试
```bash
# 测试健康检查
curl -I https://your-domain.com/

# 测试API接口
curl -X GET https://your-domain.com/api/home/init
```

### 2. 性能测试
```bash
# 安装ab工具
sudo apt install apache2-utils -y

# 压力测试
ab -n 1000 -c 10 https://your-domain.com/api/home/init
```

## 故障排除

### 1. 常见问题
- 服务无法启动：检查端口占用和权限
- 数据库连接失败：检查网络和认证
- SSL证书问题：检查证书有效期和配置

### 2. 日志查看
```bash
# 查看应用日志
sudo journalctl -u anyuyinian -f

# 查看Nginx日志
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log
```

## 备份策略

### 1. 数据库备份
```bash
# 创建备份脚本
cat > /root/backup.sh << 'EOF'
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
mysqldump -u root -p anyuyinian > /backup/db_$DATE.sql
find /backup -name "db_*.sql" -mtime +7 -delete
EOF

# 添加到定时任务
chmod +x /root/backup.sh
crontab -e
# 添加：0 2 * * * /root/backup.sh
```

### 2. 代码备份
```bash
# 使用Git进行版本控制
git add .
git commit -m "Backup $(date)"
git push origin main
```

## 联系支持

如果您需要帮助部署，请提供：
1. 服务器配置信息
2. 域名信息
3. 数据库配置
4. 具体错误信息 