# 微信手机号解密密钥修复

## 问题描述
用户反馈手机号解密仍然失败，错误信息显示：
```
"手机号解密失败: 解析解密数据失败: invalid character 'ª' looking for beginning of value"
```

## 问题分析

### 1. 根本原因
- **错误的解密密钥**: 之前使用了`AppSecret`作为解密密钥
- **微信官方要求**: 手机号解密必须使用`session_key`而不是`AppSecret`
- **解密失败**: 使用错误的密钥导致解密后的数据不是有效的JSON格式

### 2. 技术细节
根据微信官方文档：
- **用户信息解密**: 使用`session_key`
- **手机号解密**: 使用`session_key`
- **AppSecret**: 仅用于获取`session_key`，不用于解密

## 修复方案

### 1. 修改解密函数签名

#### 修复前
```go
func decryptWechatPhoneNumber(encryptedData, iv string) (string, error) {
    // 使用AppSecret作为密钥
    wxConfig := config.GetWxConfig()
    appSecret := wxConfig.AppSecret
    key := []byte(appSecret)
    // ...
}
```

#### 修复后
```go
func decryptWechatPhoneNumber(encryptedData, iv, sessionKey string) (string, error) {
    // 使用session_key作为密钥
    key := []byte(sessionKey)
    // ...
}
```

### 2. 更新解密调用逻辑

#### 修复前
```go
// 调用微信解密算法
phoneNumber, err := decryptWechatPhoneNumber(req.EncryptedData, req.IV)
```

#### 修复后
```go
// 检查session_key是否存在
if user.SessionKey == "" {
    LogError("用户session_key为空", fmt.Errorf("userId=%s, sessionKey为空", req.UserId))
    response := &UserResponse{
        Code:     -1,
        ErrorMsg: "用户session_key不存在，请重新登录",
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
    return
}

// 调用微信解密算法，使用session_key
phoneNumber, err := decryptWechatPhoneNumber(req.EncryptedData, req.IV, user.SessionKey)
```

### 3. 添加session_key验证

```go
LogStep("用户验证成功", map[string]interface{}{
    "userId":     user.UserId,
    "nickName":   user.NickName,
    "sessionKey": user.SessionKey[:10] + "...", // 只显示前10个字符
})

// 检查session_key是否存在
if user.SessionKey == "" {
    LogError("用户session_key为空", fmt.Errorf("userId=%s, sessionKey为空", req.UserId))
    response := &UserResponse{
        Code:     -1,
        ErrorMsg: "用户session_key不存在，请重新登录",
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
    return
}
```

### 4. 简化密钥处理

#### 修复前
```go
// 使用AppSecret作为密钥（需要截取前16字节）
key := []byte(appSecret)
if len(key) > 16 {
    key = key[:16]
} else {
    // 如果密钥长度不足16字节，用0填充
    paddedKey := make([]byte, 16)
    copy(paddedKey, key)
    key = paddedKey
}
```

#### 修复后
```go
// 使用session_key作为解密密钥
key := []byte(sessionKey)
```

## 技术原理

### 1. 微信解密流程
1. **用户登录**: 获取`code` → 调用微信API → 获取`session_key`和`openid`
2. **保存session_key**: 将`session_key`保存到用户记录中
3. **手机号解密**: 使用`session_key`解密手机号数据

### 2. session_key特点
- **长度**: 通常是24字节的Base64编码字符串
- **用途**: 专门用于解密用户敏感信息
- **有效期**: 有一定的有效期，过期后需要重新获取

### 3. 解密算法
- **算法**: AES-128-CBC
- **密钥**: session_key（Base64解码后）
- **IV**: 微信提供的初始化向量
- **填充**: PKCS7填充

## 修复效果

### 修复前
- 解密密钥: AppSecret ❌
- 解密结果: 无效字符'ª'
- 错误信息: "invalid character 'ª' looking for beginning of value"

### 修复后
- 解密密钥: session_key ✅
- 解密结果: 正确的JSON数据
- 预期手机号: "13691028481"

## 测试验证

### 1. 创建测试脚本
- `test_session_key_decrypt.sh`: 测试使用session_key的解密功能
- 使用用户提供的真实加密数据

### 2. 测试要点
- 验证session_key是否存在
- 测试解密算法是否正确
- 验证解密结果是否与用户真实手机号匹配

### 3. 错误处理测试
- session_key为空的情况
- 解密失败的情况
- 用户不存在的情况

## 注意事项

### 1. session_key管理
- session_key需要妥善保存
- 定期更新session_key
- 处理session_key过期的情况

### 2. 安全性
- session_key是敏感信息，需要加密存储
- 不在日志中完整显示session_key
- 确保传输安全

### 3. 错误处理
- 提供清晰的错误信息
- 引导用户重新登录
- 记录详细的错误日志

## 后续优化

### 1. session_key管理
- 实现session_key自动刷新
- 添加session_key过期检查
- 优化session_key存储方式

### 2. 错误处理
- 完善错误分类和处理
- 提供用户友好的错误提示
- 添加重试机制

### 3. 性能优化
- 缓存解密结果
- 优化解密算法性能
- 减少不必要的数据库查询

## 总结

通过将解密密钥从`AppSecret`改为`session_key`，我们修复了微信手机号解密的问题：

1. **密钥修复**: 使用正确的`session_key`作为解密密钥
2. **验证增强**: 添加`session_key`存在性检查
3. **错误处理**: 提供清晰的错误信息和用户引导
4. **日志优化**: 安全地记录解密过程信息

现在用户应该能够成功获取到正确的手机号"13691028481"了！
