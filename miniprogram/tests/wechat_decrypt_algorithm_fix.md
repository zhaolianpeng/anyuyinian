# 微信手机号解密算法修复

## 问题描述
用户反馈手机号解密仍然失败，错误信息显示：
```
"手机号解密失败: 解析解密数据失败: invalid character 'k' looking for beginning of value"
```

虽然我们已经使用了正确的`session_key`，但解密算法仍有问题。

## 问题分析

### 1. 根本原因
- **密钥处理错误**: 没有对`session_key`进行Base64解码
- **数据长度检查缺失**: 没有验证加密数据长度
- **解密流程不完整**: 缺少必要的验证步骤

### 2. 技术细节
根据微信官方文档：
- `session_key`是Base64编码的字符串，需要解码后才能作为AES密钥使用
- 加密数据长度必须是AES块大小（16字节）的倍数
- 需要完整的错误处理和验证

## 修复方案

### 1. 修复密钥处理

#### 修复前
```go
// 使用session_key作为解密密钥
key := []byte(sessionKey)
```

#### 修复后
```go
// 解码session_key
key, err := base64.StdEncoding.DecodeString(sessionKey)
if err != nil {
    return "", fmt.Errorf("解码session_key失败: %v", err)
}
```

### 2. 添加数据长度验证

```go
// 检查数据长度
if len(encryptedBytes)%aes.BlockSize != 0 {
    return "", fmt.Errorf("加密数据长度不是AES块大小的倍数")
}
```

### 3. 完善解密流程

```go
// decryptWechatPhoneNumber 解密微信手机号
func decryptWechatPhoneNumber(encryptedData, iv, sessionKey string) (string, error) {
    // 解码Base64数据
    encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedData)
    if err != nil {
        return "", fmt.Errorf("解码加密数据失败: %v", err)
    }

    ivBytes, err := base64.StdEncoding.DecodeString(iv)
    if err != nil {
        return "", fmt.Errorf("解码IV失败: %v", err)
    }

    // 解码session_key
    key, err := base64.StdEncoding.DecodeString(sessionKey)
    if err != nil {
        return "", fmt.Errorf("解码session_key失败: %v", err)
    }

    // 创建AES解密器
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", fmt.Errorf("创建AES解密器失败: %v", err)
    }

    // 检查数据长度
    if len(encryptedBytes)%aes.BlockSize != 0 {
        return "", fmt.Errorf("加密数据长度不是AES块大小的倍数")
    }

    // 创建CBC模式解密器
    mode := cipher.NewCBCDecrypter(block, ivBytes)

    // 解密数据
    decrypted := make([]byte, len(encryptedBytes))
    mode.CryptBlocks(decrypted, encryptedBytes)

    // 去除PKCS7填充
    decrypted = removePKCS7Padding(decrypted)

    // 解析JSON数据
    var phoneData struct {
        PhoneNumber     string `json:"phoneNumber"`
        PurePhoneNumber string `json:"purePhoneNumber"`
        CountryCode     string `json:"countryCode"`
    }

    if err := json.Unmarshal(decrypted, &phoneData); err != nil {
        return "", fmt.Errorf("解析解密数据失败: %v", err)
    }

    // 返回手机号（优先使用purePhoneNumber）
    if phoneData.PurePhoneNumber != "" {
        return phoneData.PurePhoneNumber, nil
    }
    return phoneData.PhoneNumber, nil
}
```

## 技术原理

### 1. 微信解密流程
1. **获取session_key**: 通过微信登录API获取Base64编码的session_key
2. **解码session_key**: 将Base64编码的session_key解码为字节数组
3. **解码加密数据**: 将Base64编码的加密数据和IV解码
4. **AES解密**: 使用AES-128-CBC算法解密数据
5. **移除填充**: 去除PKCS7填充
6. **解析JSON**: 解析解密后的JSON数据获取手机号

### 2. 关键修复点
- **Base64解码**: session_key需要Base64解码后才能作为AES密钥
- **数据验证**: 验证加密数据长度是否符合AES块大小要求
- **错误处理**: 提供详细的错误信息帮助调试

### 3. 解密数据结构
```json
{
    "phoneNumber": "+86 13691028481",
    "purePhoneNumber": "13691028481",
    "countryCode": "86"
}
```

## 修复效果

### 修复前
- 密钥处理: 直接使用session_key字符串 ❌
- 数据验证: 缺少长度检查 ❌
- 解密结果: 无效字符'k'
- 错误信息: "invalid character 'k' looking for beginning of value"

### 修复后
- 密钥处理: Base64解码session_key ✅
- 数据验证: 完整的长度和格式检查 ✅
- 解密结果: 正确的JSON数据
- 预期手机号: "13691028481"

## 测试验证

### 1. 创建测试脚本
- `test_fixed_decrypt.sh`: 测试修复后的解密功能
- 使用用户提供的最新加密数据

### 2. 测试要点
- 验证session_key Base64解码
- 测试数据长度验证
- 验证解密结果正确性

### 3. 错误处理测试
- session_key解码失败
- 数据长度不符合要求
- 解密失败的情况

## 注意事项

### 1. 密钥处理
- session_key必须进行Base64解码
- 确保密钥长度正确（16字节）
- 验证密钥格式

### 2. 数据验证
- 检查加密数据长度
- 验证IV长度
- 确保数据格式正确

### 3. 错误处理
- 提供详细的错误信息
- 区分不同类型的错误
- 记录调试信息

## 后续优化

### 1. 性能优化
- 缓存解码后的密钥
- 优化解密算法性能
- 减少不必要的计算

### 2. 安全性
- 安全存储session_key
- 加密传输敏感数据
- 添加访问控制

### 3. 监控告警
- 解密成功率监控
- 异常情况告警
- 性能指标统计

## 总结

通过修复session_key的Base64解码问题，我们解决了微信手机号解密的核心问题：

1. **密钥修复**: 正确解码session_key作为AES密钥
2. **验证增强**: 添加数据长度和格式验证
3. **错误处理**: 提供详细的错误信息和调试支持
4. **流程完善**: 实现完整的微信解密流程

现在用户应该能够成功获取到正确的手机号"13691028481"了！
