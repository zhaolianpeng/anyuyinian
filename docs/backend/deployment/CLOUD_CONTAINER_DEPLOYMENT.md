# 微信云托管部署指南

## 部署步骤

### 1. 准备代码
确保所有代码已提交到Git仓库，包括：
- 后端Go代码
- 小程序前端代码
- 数据库迁移脚本

### 2. 微信云托管部署

#### 2.1 登录微信云托管控制台
1. 访问 [微信云托管控制台](https://cloud.weixin.qq.com/)
2. 选择环境：`prod-5g94mx7a3d07e78c`
3. 选择服务：`golang-lfwy`

#### 2.2 上传代码
1. 在服务详情页面，点击"版本管理"
2. 选择"上传代码包"
3. 上传包含以下文件的代码包：
   - `main.go`
   - `Dockerfile`
   - `container.config.json`
   - `db/` 目录
   - `service/` 目录
   - `config/` 目录

#### 2.3 配置环境变量
在服务设置中配置以下环境变量：
```
DB_HOST=10.3.110.11
DB_PORT=3306
DB_USER=root
DB_PASSWORD=bU4X6cFW
DB_NAME=anyuyinian
```

#### 2.4 部署服务
1. 点击"部署"按钮
2. 等待部署完成（通常需要3-5分钟）
3. 检查服务状态为"运行中"

### 3. 数据库迁移

#### 3.1 执行管理员数据库迁移
```bash
# 连接到云托管数据库
mysql -h 10.3.110.11 -u root -p anyuyinian

# 执行迁移脚本
source db/migration/add_admin_fields.sql
```

#### 3.2 验证数据库表结构
```sql
-- 检查管理员字段是否添加成功
DESCRIBE Users;

-- 检查默认超级管理员账号
SELECT userId, nickName, isAdmin, adminLevel, adminUsername 
FROM Users 
WHERE isAdmin = 1;
```

### 4. 小程序配置

#### 4.1 确保小程序已配置云托管
在 `miniprogram/app.json` 中添加：
```json
{
  "cloud": true
}
```

#### 4.2 初始化云托管
在 `miniprogram/app.js` 中添加：
```javascript
App({
  onLaunch: function () {
    if (!wx.cloud) {
      console.error('请使用 2.2.3 或以上的基础库以使用云能力')
    } else {
      wx.cloud.init({
        env: 'prod-5g94mx7a3d07e78c',
        traceUser: true,
      })
    }
  }
})
```

### 5. 测试部署

#### 5.1 测试管理员登录
1. 在微信开发者工具中打开小程序
2. 进入"我的"页面
3. 点击"管理员入口"
4. 测试超级管理员登录和普通管理员登录

#### 5.2 测试API接口
```bash
# 测试管理员登录
curl -X POST https://golang-lfwy-176496-6-1353115175.sh.run.tcloudbase.com/api/admin/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "anyuyinian",
    "password": "000000"
  }'

# 测试检查管理员状态
curl -X GET "https://golang-lfwy-176496-6-1353115175.sh.run.tcloudbase.com/api/admin/check-status?userId=test_user_id"
```

## 配置说明

### 云托管配置
- **环境ID**: `prod-5g94mx7a3d07e78c`
- **服务名**: `golang-lfwy`
- **容器端口**: 80
- **CPU**: 0.25核
- **内存**: 0.5GB
- **最小实例数**: 0
- **最大实例数**: 10

### 数据库配置
- **主机**: 10.3.110.11
- **端口**: 3306
- **数据库**: anyuyinian
- **用户名**: root
- **密码**: bU4X6cFW

### 小程序配置
- **AppID**: wx101090677bd5219e
- **云环境**: prod-5g94mx7a3d07e78c

## 监控和维护

### 1. 服务监控
- 在云托管控制台查看服务运行状态
- 监控CPU和内存使用情况
- 查看访问日志和错误日志

### 2. 数据库监控
- 监控数据库连接数
- 检查慢查询日志
- 定期备份数据库

### 3. 小程序监控
- 在微信公众平台查看小程序使用情况
- 监控API调用次数和成功率
- 查看用户反馈和错误报告

## 故障排除

### 1. 服务无法启动
- 检查Dockerfile配置
- 验证环境变量设置
- 查看容器启动日志

### 2. 数据库连接失败
- 检查数据库服务器状态
- 验证网络连接
- 确认数据库用户权限

### 3. API调用失败
- 检查云托管服务状态
- 验证API路径配置
- 查看请求日志

### 4. 小程序无法连接
- 确认云托管环境配置
- 检查小程序AppID设置
- 验证云托管服务域名

## 安全建议

### 1. 数据库安全
- 定期更换数据库密码
- 限制数据库访问IP
- 启用数据库审计日志

### 2. API安全
- 添加请求频率限制
- 实现API签名验证
- 记录异常访问日志

### 3. 管理员安全
- 定期更换管理员密码
- 限制管理员登录IP
- 记录管理员操作日志

---

**部署完成后，管理员功能将完全通过微信云托管服务运行，确保安全性和稳定性。** 