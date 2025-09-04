# 生产环境支付回调配置指南

## 🚨 重要提醒

**云托管域名不能用于生产环境支付回调！**

当前使用的云托管域名 `https://golang-lfwy-176496-6-1353115175.sh.run.tcloudbase.com` 仅适用于开发测试，生产环境必须使用正式域名。

## 🎯 推荐方案

### 方案一：云服务器部署（强烈推荐）

#### 优势
- ✅ 完全控制，稳定可靠
- ✅ 支持HTTPS证书
- ✅ 可配置域名
- ✅ 支持高并发
- ✅ 成本可控

#### 部署步骤
1. **购买云服务器**
   - 阿里云、腾讯云、华为云
   - 配置：2核4G，5Mbps带宽
   - 操作系统：Ubuntu 20.04

2. **服务器配置**
   ```bash
   # 运行配置脚本
   chmod +x scripts/setup_production_server.sh
   ./scripts/setup_production_server.sh
   ```

3. **部署应用**
   ```bash
   # 使用部署脚本
   chmod +x scripts/deploy_production.sh
   ./deploy_production.sh your-domain.com
   ```

4. **配置支付回调**
   - 微信商户平台设置回调地址：`https://your-domain.com/api/payment/notify`
   - 确保HTTPS证书有效
   - 测试回调接口可访问性

#### 成本估算
- 云服务器：约200-500元/月
- 域名：约50-100元/年
- SSL证书：免费（Let's Encrypt）

### 方案二：云函数部署

#### 优势
- ✅ 无需管理服务器
- ✅ 自动扩缩容
- ✅ 按量付费
- ✅ 高可用

#### 部署步骤
1. **构建云函数包**
   ```bash
   chmod +x scripts/build_cloud_function.sh
   ./scripts/build_cloud_function.sh
   ```

2. **配置腾讯云凭证**
   ```bash
   # 安装Serverless Framework
   npm install -g serverless
   
   # 配置凭证
   serverless config credentials --provider tencent --key YOUR_SECRET_ID --secret YOUR_SECRET_KEY
   ```

3. **部署云函数**
   ```bash
   cd cloud_function_package
   ./deploy.sh
   ```

4. **获取云函数URL**
   - 部署完成后获取API网关地址
   - 配置到微信商户平台

#### 成本估算
- 云函数：约100-300元/月
- API网关：约50-150元/月

### 方案三：Docker容器化部署

#### 优势
- ✅ 环境一致
- ✅ 易于扩展
- ✅ 支持多环境
- ✅ 便于维护

#### 部署步骤
1. **准备Docker环境**
   ```bash
   # 安装Docker和Docker Compose
   curl -fsSL https://get.docker.com -o get-docker.sh
   sh get-docker.sh
   ```

2. **配置环境变量**
   ```bash
   # 修改docker-compose.production.yml中的配置
   vim docker-compose.production.yml
   ```

3. **启动服务**
   ```bash
   docker-compose -f docker-compose.production.yml up -d
   ```

4. **配置域名和SSL**
   - 配置域名解析
   - 申请SSL证书
   - 配置Nginx反向代理

## 🔧 支付回调配置

### 1. 微信商户平台配置

1. **登录微信商户平台**
   - 访问：https://pay.weixin.qq.com/
   - 商户号：1726638701

2. **配置支付回调地址**
   - 产品中心 → 开发配置 → 支付配置
   - 回调地址：`https://your-domain.com/api/payment/notify`

3. **验证回调地址**
   - 确保地址可公网访问
   - 确保使用HTTPS协议
   - 确保返回"success"字符串

### 2. 后端代码配置

确保后端代码中的支付回调处理正确：

```go
// 支付回调处理
func HandleWechatPayNotify(w http.ResponseWriter, r *http.Request) {
    // 验证签名
    // 处理支付结果
    // 更新订单状态
    // 返回success
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("success"))
}
```

### 3. 测试支付回调

```bash
# 测试回调接口
curl -X POST https://your-domain.com/api/payment/notify \
  -H "Content-Type: application/xml" \
  -d '<xml>...</xml>'
```

## ⚠️ 注意事项

### 1. 域名要求
- 必须使用HTTPS协议
- 域名必须可公网访问
- 建议使用正式域名，避免使用云托管域名

### 2. 安全要求
- 验证微信支付签名
- 防止重复处理
- 记录详细日志
- 设置超时时间

### 3. 性能要求
- 回调处理时间 < 5秒
- 支持并发处理
- 异常情况重试机制

## 🚀 快速开始

如果您想快速部署到生产环境，推荐使用云服务器方案：

1. **购买云服务器**（阿里云/腾讯云）
2. **运行配置脚本**：`./scripts/setup_production_server.sh`
3. **部署应用**：`./scripts/deploy_production.sh your-domain.com`
4. **配置支付回调**：在微信商户平台设置回调地址
5. **测试支付功能**：创建订单并测试支付

## 📞 技术支持

如果在部署过程中遇到问题，可以：
1. 查看部署日志：`sudo journalctl -u anyuyinian -f`
2. 检查Nginx状态：`sudo systemctl status nginx`
3. 测试API接口：`curl https://your-domain.com/api/health`
4. 查看错误日志：`tail -f /var/log/nginx/error.log`
