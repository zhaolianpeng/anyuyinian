# 预约页面调试指南

## 问题描述

用户反馈预约页面无法显示就诊人和地址信息，但接口返回了正确的数据。

## 接口返回数据

### 就诊人API返回
```json
{
  "code": 0,
  "data": [
    {
      "id": 1,
      "userId": 1,
      "name": "赵连鹏",
      "idCard": "231003199006071015",
      "phone": "13691028481",
      "gender": 1,
      "birthday": "1990-06-07",
      "relation": "",
      "isDefault": 0,
      "status": 1,
      "createdAt": "2025-08-01T03:23:39Z",
      "updatedAt": "2025-08-01T03:23:39Z"
    }
  ]
}
```

### 地址API返回
```json
{
  "code": 0,
  "data": [
    {
      "id": 1,
      "userId": 1,
      "name": "赵连鹏",
      "phone": "13691028481",
      "province": "北京市",
      "city": "北京市",
      "district": "西城区",
      "address": "复兴门北大街9号楼",
      "isDefault": 0,
      "status": 1,
      "createdAt": "2025-08-01T03:16:21Z",
      "updatedAt": "2025-08-01T03:16:21Z"
    }
  ]
}
```

## 可能的问题原因

### 1. 用户ID类型问题
- 存储的userId可能是字符串类型，但API需要数字类型
- 解决方案：在API调用时转换为数字类型

### 2. 页面加载时机问题
- 数据加载可能在页面渲染之前完成
- 解决方案：延迟加载数据，确保页面完全加载

### 3. 数据格式问题
- 就诊人的relation字段为空字符串，可能导致显示问题
- 解决方案：在模板中提供默认值

### 4. 页面更新问题
- setData可能没有正确触发页面更新
- 解决方案：添加强制更新机制

## 已实施的修复

### 1. 增强调试信息
- 添加了详细的控制台日志
- 在页面上显示调试信息
- 添加了测试按钮

### 2. 修复用户ID类型
```javascript
// 确保userId是数字类型
const numericUserId = parseInt(userId) || userId
```

### 3. 修复页面加载时机
```javascript
// 延迟加载用户数据，确保页面完全加载
setTimeout(() => {
  this.loadUserData()
}, 500)
```

### 4. 添加强制更新
```javascript
forceUpdate() {
  console.log('强制更新页面数据')
  this.setData({
    patientList: this.data.patientList,
    addressList: this.data.addressList,
    selectedPatient: this.data.selectedPatient,
    selectedAddress: this.data.selectedAddress
  })
}
```

### 5. 添加模拟数据测试
```javascript
testMockData() {
  // 设置模拟数据来测试页面显示
  const mockPatientList = [...]
  const mockAddressList = [...]
  this.setData({
    patientList: mockPatientList,
    addressList: mockAddressList,
    selectedPatient: mockPatientList[0],
    selectedAddress: mockAddressList[0]
  })
}
```

## 调试步骤

### 1. 检查用户登录状态
- 确认userId是否正确存储
- 确认用户信息是否完整

### 2. 检查API调用
- 查看控制台日志中的API响应
- 确认数据格式是否正确

### 3. 检查页面数据
- 使用"调试数据"按钮查看页面状态
- 使用"测试模拟数据"按钮验证页面显示

### 4. 检查模板渲染
- 查看页面上的调试信息
- 确认数据是否正确传递到模板

## 测试方法

### 1. 使用调试按钮
- 点击"调试数据"查看当前页面状态
- 点击"测试模拟数据"验证页面显示功能

### 2. 查看控制台日志
- 打开开发者工具的控制台
- 查看详细的数据加载日志

### 3. 检查网络请求
- 在开发者工具的网络面板中查看API请求
- 确认请求参数和响应数据

## 预期结果

修复后，预约页面应该能够：
1. 正确显示就诊人信息（姓名、关系、电话）
2. 正确显示地址信息（姓名、电话、详细地址）
3. 支持选择就诊人和地址
4. 显示正确的预约时间范围
5. 支持备注输入

## 如果问题仍然存在

如果修复后问题仍然存在，请：
1. 查看控制台日志，确认数据加载过程
2. 使用"测试模拟数据"按钮验证页面显示功能
3. 检查用户登录状态和用户ID
4. 确认网络连接和API服务状态 