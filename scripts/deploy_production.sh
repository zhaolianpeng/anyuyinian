#!/bin/bash

# 生产环境部署脚本
# 使用方法: ./deploy_production.sh your-domain.com

if [ $# -eq 0 ]; then
    echo "❌ 请提供域名参数"
    echo "使用方法: ./deploy_production.sh your-domain.com"
    exit 1
fi

DOMAIN=$1
APP_DIR="/var/www/anyuyinian"

echo "🚀 开始部署到生产环境..."
echo "🌐 域名: $DOMAIN"
echo "📁 应用目录: $APP_DIR"

# 1. 停止服务
echo "⏹️ 停止现有服务..."
sudo systemctl stop anyuyinian

# 2. 备份当前版本
echo "💾 备份当前版本..."
if [ -f "$APP_DIR/app" ]; then
    sudo cp "$APP_DIR/app" "$APP_DIR/app.backup.$(date +%Y%m%d_%H%M%S)"
fi

# 3. 编译新版本
echo "🔨 编译应用..."
go build -o app main.go
if [ $? -ne 0 ]; then
    echo "❌ 编译失败"
    exit 1
fi

# 4. 部署到服务器
echo "📦 部署应用..."
sudo cp app "$APP_DIR/"
sudo cp -r static "$APP_DIR/"
sudo cp -r templates "$APP_DIR/"
sudo cp -r config "$APP_DIR/"
sudo cp -r db "$APP_DIR/"
sudo cp -r service "$APP_DIR/"
sudo cp -r handler "$APP_DIR/"
sudo cp -r middleware "$APP_DIR/"
sudo cp -r utils "$APP_DIR/"
sudo cp go.mod go.sum "$APP_DIR/"

# 5. 设置权限
echo "🔐 设置权限..."
sudo chown -R www-data:www-data "$APP_DIR"
sudo chmod +x "$APP_DIR/app"

# 6. 配置环境变量
echo "⚙️ 配置环境变量..."
sudo tee "$APP_DIR/.env" << EOF
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=anyuyinian

# 微信配置
WX_APP_ID=wx101090677bd5219e
WX_APP_SECRET=your_app_secret

# 微信支付配置
WECHAT_PAY_APP_ID=wx101090677bd5219e
WECHAT_PAY_MCH_ID=1726638701
WECHAT_PAY_MCH_KEY=JQzOCB8doIdgaUjAobELsk9nTyxdKhat
WECHAT_PAY_NOTIFY_URL=https://$DOMAIN/api/payment/notify
WECHAT_PAY_ENVIRONMENT=production

# 服务器配置
GIN_MODE=release
PORT=80
EOF

# 7. 配置Nginx
echo "🌐 配置Nginx..."
sudo cp nginx_config.conf /etc/nginx/sites-available/anyuyinian
sudo sed -i "s/your-domain.com/$DOMAIN/g" /etc/nginx/sites-available/anyuyinian
sudo ln -sf /etc/nginx/sites-available/anyuyinian /etc/nginx/sites-enabled/
sudo nginx -t

# 8. 申请SSL证书
echo "🔒 申请SSL证书..."
sudo certbot --nginx -d $DOMAIN --non-interactive --agree-tos --email admin@$DOMAIN

# 9. 启动服务
echo "▶️ 启动服务..."
sudo systemctl daemon-reload
sudo systemctl enable anyuyinian
sudo systemctl start anyuyinian
sudo systemctl reload nginx

# 10. 检查服务状态
echo "🔍 检查服务状态..."
sleep 5
if systemctl is-active --quiet anyuyinian; then
    echo "✅ 应用服务启动成功"
else
    echo "❌ 应用服务启动失败"
    sudo systemctl status anyuyinian
fi

if systemctl is-active --quiet nginx; then
    echo "✅ Nginx服务运行正常"
else
    echo "❌ Nginx服务异常"
    sudo systemctl status nginx
fi

# 11. 测试API
echo "🧪 测试API..."
curl -s "https://$DOMAIN/api/health" | head -1

echo "🎉 部署完成！"
echo "🌐 访问地址: https://$DOMAIN"
echo "📱 支付回调地址: https://$DOMAIN/api/payment/notify"
echo "🔍 查看日志: sudo journalctl -u anyuyinian -f"
