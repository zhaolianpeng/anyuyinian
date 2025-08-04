# 微信配置问题解决方案

## 问题描述
后端日志显示微信API调用失败，错误信息：`invalid appid, rid: 688b379c-6aa34941-7c47450d`

## 问题原因
后端使用的微信配置是占位符：
- `appid`: `your_app_id`
- `secret`: `your_app_secret`

这些不是真实的微信小程序配置，导致微信API返回 `invalid appid` 错误。

## 解决方案

### 1. 设置环境变量
需要在运行后端服务前设置正确的微信小程序配置：

```bash
# 设置微信小程序AppID和AppSecret
export WX_APP_ID="你的真实微信小程序AppID"
export WX_APP_SECRET="你的真实微信小程序AppSecret"
```

### 2. 获取微信小程序配置
1. 登录微信公众平台：https://mp.weixin.qq.com/
2. 选择你的小程序
3. 进入"开发" -> "开发管理" -> "开发设置"
4. 复制 `AppID(小程序ID)` 和 `AppSecret(小程序密钥)`

### 3. 验证配置
设置环境变量后，可以通过以下方式验证：

```bash
# 检查环境变量是否设置成功
echo $WX_APP_ID
echo $WX_APP_SECRET
```

### 4. 重启后端服务
设置环境变量后，需要重启后端服务：

```bash
# 停止当前服务
# 重新启动服务
go run main.go
```

## 配置示例

### 正确的环境变量设置
```bash
export WX_APP_ID="wx1234567890abcdef"
export WX_APP_SECRET="abcdef1234567890abcdef1234567890"
```

### 错误的配置（当前状态）
```bash
export WX_APP_ID="your_app_id"
export WX_APP_SECRET="your_app_secret"
```

## 验证步骤

### 1. 检查环境变量
```bash
# 检查环境变量是否正确设置
env | grep WX_APP
```

### 2. 测试微信API调用
设置正确的环境变量后，重新测试登录功能：

1. 在小程序中点击"微信登录"
2. 查看后端日志，确认不再出现 `invalid appid` 错误
3. 确认登录流程正常完成

### 3. 查看后端日志
正确的日志应该显示：
```
[STEP] 微信配置获取成功: {"appID":"wx1234567890abcdef"}
[STEP] 微信API调用成功: {"openid":"xxx","session_key":"xxx"}
```

而不是：
```
[STEP] 微信配置获取成功: {"appID":"your_app_id"}
[ERROR] 微信API返回错误: 错误码: 40013, 错误信息: invalid appid
```

## 注意事项

### 1. 安全性
- `AppSecret` 是敏感信息，不要提交到代码仓库
- 建议使用环境变量或配置文件管理
- 生产环境应该使用更安全的配置管理方式

### 2. 环境变量持久化
如果需要在每次启动时自动设置环境变量，可以：

#### 方法1：创建环境变量文件
```bash
# 创建 .env 文件
echo "export WX_APP_ID=你的AppID" >> ~/.bashrc
echo "export WX_APP_SECRET=你的AppSecret" >> ~/.bashrc
source ~/.bashrc
```

#### 方法2：使用启动脚本
```bash
#!/bin/bash
# start.sh
export WX_APP_ID="你的AppID"
export WX_APP_SECRET="你的AppSecret"
go run main.go
```

### 3. 开发环境配置
开发环境可以使用测试账号的配置，但生产环境必须使用正式的小程序配置。

## 相关文件
- `config/wx_config.go` - 微信配置读取逻辑
- `service/wx_login_service.go` - 微信登录服务
- `CONFIG.md` - 配置说明文档

## 下一步操作
1. 获取真实的微信小程序AppID和AppSecret
2. 设置环境变量
3. 重启后端服务
4. 测试登录功能
5. 确认问题解决 