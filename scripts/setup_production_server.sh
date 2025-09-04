#!/bin/bash

# 生产环境服务器配置脚本
# 适用于 Ubuntu 20.04

echo "🚀 开始配置生产环境服务器..."

# 1. 更新系统
echo "📦 更新系统包..."
sudo apt update && sudo apt upgrade -y

# 2. 安装Go环境
echo "🔧 安装Go环境..."
wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 3. 安装MySQL
echo "🗄️ 安装MySQL..."
sudo apt install mysql-server -y
sudo mysql_secure_installation

# 4. 安装Nginx
echo "🌐 安装Nginx..."
sudo apt install nginx -y

# 5. 安装Certbot (SSL证书)
echo "🔒 安装SSL证书工具..."
sudo apt install certbot python3-certbot-nginx -y

# 6. 创建应用目录
echo "📁 创建应用目录..."
sudo mkdir -p /var/www/anyuyinian
sudo chown -R $USER:$USER /var/www/anyuyinian

# 7. 创建系统服务
echo "⚙️ 创建系统服务..."
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
Environment=GIN_MODE=release

[Install]
WantedBy=multi-user.target
EOF

echo "✅ 服务器环境配置完成！"
echo "📝 下一步："
echo "1. 将代码上传到 /var/www/anyuyinian"
echo "2. 编译应用: go build -o app main.go"
echo "3. 配置数据库连接"
echo "4. 配置Nginx反向代理"
echo "5. 申请SSL证书"
