# 二维码系统说明

## 概述

二维码系统用于生成推广员专属的推广二维码，用户扫描二维码可以直接进入小程序并自动填入推广码。

## 功能特性

### 1. 二维码生成
- 基于推广码生成唯一二维码
- 支持URL和Base64两种格式
- 自动创建小程序页面链接

### 2. 二维码存储
- 支持文件系统存储
- 支持Base64内联显示
- 自动创建存储目录

### 3. 二维码访问
- 提供公共访问URL
- 支持图片预览和下载
- 错误处理和占位符

## API接口

### 1. 生成二维码URL
```
GET /api/qrcode/generate?promoterCode={promoterCode}
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "promoterCode": "ABC123",
    "qrCodeUrl": "https://your-domain.com/static/qrcode/promoter_ABC123.png",
    "pageURL": "https://your-domain.com/pages/index/index?promoterCode=ABC123"
  }
}
```

### 2. 生成Base64编码的二维码
```
GET /api/qrcode/generate_base64?promoterCode={promoterCode}
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "promoterCode": "ABC123",
    "base64QRCode": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...",
    "pageURL": "https://your-domain.com/pages/index/index?promoterCode=ABC123"
  }
}
```

## 技术实现

### 1. 二维码生成库
使用 `github.com/skip2/go-qrcode` 库生成二维码：
```go
import "github.com/skip2/go-qrcode"

// 生成二维码
qr, err := qrcode.New(pageURL, qrcode.Medium)
if err != nil {
    return "", fmt.Errorf("创建二维码失败: %v", err)
}
```

### 2. 文件存储
```go
// 确保输出目录存在
if err := os.MkdirAll(s.OutputDir, 0755); err != nil {
    return "", fmt.Errorf("创建输出目录失败: %v", err)
}

// 保存二维码图片
filename := fmt.Sprintf("promoter_%s.png", promoterCode)
filepath := filepath.Join(s.OutputDir, filename)
if err := os.WriteFile(filepath, qrCode, 0644); err != nil {
    return "", fmt.Errorf("保存二维码图片失败: %v", err)
}
```

### 3. Base64编码
```go
// 编码为PNG
var buf bytes.Buffer
if err := png.Encode(&buf, qr.Image(256)); err != nil {
    return "", fmt.Errorf("编码PNG失败: %v", err)
}

// 转换为Base64
base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
dataURL := fmt.Sprintf("data:image/png;base64,%s", base64Str)
```

## 配置说明

### 1. 服务配置
```go
type QRCodeService struct {
    BaseURL    string // 小程序页面基础URL
    OutputDir  string // 二维码图片输出目录
    PublicURL  string // 公共访问URL
}
```

### 2. 默认配置
```go
func NewQRCodeService() *QRCodeService {
    return &QRCodeService{
        BaseURL:   "https://your-domain.com/pages/index/index?promoterCode=",
        OutputDir: "./static/qrcode",
        PublicURL: "https://your-domain.com/static/qrcode",
    }
}
```

## 使用流程

### 1. 推广员注册
1. 用户注册成为推广员
2. 系统自动生成推广码
3. 基于推广码生成二维码
4. 保存二维码URL到数据库

### 2. 二维码分享
1. 推广员获取专属二维码
2. 分享二维码给潜在用户
3. 用户扫描二维码进入小程序
4. 自动填入推广码建立推荐关系

### 3. 佣金计算
1. 被推荐用户下单
2. 系统根据推荐关系计算佣金
3. 推广员获得相应佣金

## 错误处理

### 1. 生成失败
- 返回占位符图片URL
- 记录错误日志
- 不影响其他功能

### 2. 存储失败
- 使用Base64编码作为备选方案
- 提供错误提示
- 自动重试机制

### 3. 访问失败
- 提供占位符图片
- 显示错误信息
- 支持手动输入推广码

## 测试验证

### 1. 功能测试
```bash
# 运行二维码测试
cd tests/backend
./test_qrcode.sh
```

### 2. 测试用例
- 有效推广码生成二维码
- 无效推广码错误处理
- 空推广码参数验证
- Base64编码功能测试

### 3. 性能测试
- 二维码生成速度
- 文件存储性能
- 并发访问测试

## 部署说明

### 1. 依赖安装
```bash
go get github.com/skip2/go-qrcode
```

### 2. 目录创建
```bash
mkdir -p ./static/qrcode
chmod 755 ./static/qrcode
```

### 3. 配置更新
- 更新BaseURL为实际域名
- 配置PublicURL为可访问地址
- 设置正确的输出目录

### 4. 静态文件服务
确保静态文件目录可以通过Web服务器访问：
```nginx
location /static/qrcode/ {
    alias /path/to/your/app/static/qrcode/;
    expires 1y;
    add_header Cache-Control "public, immutable";
}
```

## 注意事项

### 1. 安全性
- 验证推广码格式
- 限制文件访问权限
- 防止恶意文件上传

### 2. 性能优化
- 缓存已生成的二维码
- 异步生成大文件
- 定期清理过期文件

### 3. 兼容性
- 支持多种二维码格式
- 兼容不同设备扫描
- 提供多种访问方式

## 总结

二维码系统提供了完整的推广码二维码生成和管理功能，支持文件存储和Base64编码两种方式，具有良好的错误处理和性能优化机制。 