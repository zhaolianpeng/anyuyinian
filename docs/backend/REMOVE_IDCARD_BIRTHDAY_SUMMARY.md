# 去掉身份证号和出生日期输入框总结

## 需求描述
在服务详情页预约信息部分去掉身份证号、出生日期的输入框。

## 当前状态分析

### 1. 服务详情页面现状
经过检查，服务详情页面的预约信息部分已经去掉了身份证号和出生日期的输入框，当前包含的字段有：

#### 固定字段（非输入框）
- **患者性别**：显示从后端获取的患者性别信息
- **患者年龄**：显示从后端根据身份证自动计算的年龄

#### 输入字段
- **基础病信息**：文本域输入框
- **是否需要助排二便**：单选框
- **期望咨询时间**：日历选择器

#### 动态表单字段
- 根据服务配置动态生成的表单字段

### 2. 身份证号和出生日期显示位置
身份证号和出生日期信息仅在以下位置显示（非输入框）：

1. **患者选择区域**：
   ```xml
   <text class="info-detail">{{selectedPatient.phone}} | {{selectedPatient.idCard}}</text>
   ```

2. **患者选择弹窗**：
   ```xml
   <text class="item-detail">{{item.phone}} | {{item.idCard}}</text>
   ```

这些显示仅用于展示已选择的患者信息，不是输入框。

## 修改内容

### 1. 动态表单字段过滤
为了确保动态表单字段中不会出现身份证号和出生日期相关的输入框，在`loadServiceDetail`函数中添加了过滤逻辑：

```javascript
// 过滤掉身份证号和出生日期相关的字段
formFields = allFields.filter(field => {
  const fieldName = field.name || ''
  const fieldLabel = field.label || ''
  
  // 过滤掉包含身份证号、出生日期等关键词的字段
  const excludeKeywords = ['身份证', 'idCard', 'birthday', '出生', '生日']
  const shouldExclude = excludeKeywords.some(keyword => 
    fieldName.toLowerCase().includes(keyword.toLowerCase()) ||
    fieldLabel.toLowerCase().includes(keyword.toLowerCase())
  )
  
  return !shouldExclude
})
```

### 2. 过滤关键词
过滤掉包含以下关键词的字段：
- `身份证`
- `idCard`
- `birthday`
- `出生`
- `生日`

### 3. 过滤逻辑
- 检查字段名称（`field.name`）和标签（`field.label`）
- 不区分大小写进行匹配
- 包含任一关键词的字段都会被过滤掉

## 修改文件

### 1. miniprogram/pages/service/detail.js
**修改位置**：`loadServiceDetail`函数中的表单配置解析部分

**修改内容**：
- 添加了动态表单字段过滤逻辑
- 过滤掉身份证号和出生日期相关的字段
- 添加了调试日志输出

**修改前**：
```javascript
// 解析表单配置
let formFields = []
try {
  const formConfig = JSON.parse(service.formConfig)
  formFields = formConfig.fields || []
} catch (error) {
  console.error('解析表单配置失败:', error)
}
```

**修改后**：
```javascript
// 解析表单配置
let formFields = []
try {
  const formConfig = JSON.parse(service.formConfig)
  const allFields = formConfig.fields || []
  
  // 过滤掉身份证号和出生日期相关的字段
  formFields = allFields.filter(field => {
    const fieldName = field.name || ''
    const fieldLabel = field.label || ''
    
    // 过滤掉包含身份证号、出生日期等关键词的字段
    const excludeKeywords = ['身份证', 'idCard', 'birthday', '出生', '生日']
    const shouldExclude = excludeKeywords.some(keyword => 
      fieldName.toLowerCase().includes(keyword.toLowerCase()) ||
      fieldLabel.toLowerCase().includes(keyword.toLowerCase())
    )
    
    return !shouldExclude
  })
  
  console.log('过滤后的表单字段:', formFields)
} catch (error) {
  console.error('解析表单配置失败:', error)
}
```

## 功能验证

### 1. 静态字段验证
- ✅ 患者性别：显示，非输入框
- ✅ 患者年龄：显示，非输入框
- ✅ 基础病信息：输入框，保留
- ✅ 是否需要助排二便：单选框，保留
- ✅ 期望咨询时间：选择器，保留

### 2. 动态字段验证
- ✅ 身份证号相关字段：被过滤
- ✅ 出生日期相关字段：被过滤
- ✅ 其他业务字段：正常显示

### 3. 显示字段验证
- ✅ 患者选择区域：显示身份证号（仅显示，非输入）
- ✅ 患者选择弹窗：显示身份证号（仅显示，非输入）

## 测试场景

### 1. 正常场景
- 服务配置中包含身份证号字段 → 被过滤掉
- 服务配置中包含出生日期字段 → 被过滤掉
- 服务配置中包含其他业务字段 → 正常显示

### 2. 边界场景
- 字段名称为空 → 正常处理
- 字段标签为空 → 正常处理
- 关键词部分匹配 → 被过滤掉
- 关键词大小写不同 → 被过滤掉

### 3. 异常场景
- 表单配置解析失败 → 正常错误处理
- 字段配置格式错误 → 正常错误处理

## 用户体验

### 1. 简化表单
- 去掉了身份证号和出生日期的输入框
- 减少了用户填写表单的复杂度
- 避免了重复输入患者信息

### 2. 数据一致性
- 身份证号和出生日期信息从患者选择中获取
- 确保数据的一致性和准确性
- 避免用户输入错误

### 3. 流程优化
- 用户只需选择患者，无需重复输入基本信息
- 身份证号和出生日期自动从患者信息中获取
- 年龄自动根据身份证号计算

## 总结

通过以上修改，服务详情页预约信息部分已经完全去掉了身份证号和出生日期的输入框：

1. **静态字段**：患者性别和年龄只显示，不提供输入框
2. **动态字段**：通过过滤逻辑确保不会出现身份证号和出生日期相关的输入框
3. **显示字段**：身份证号仅在患者信息展示中显示，不提供输入功能

所有修改都已完成，用户体验得到优化，表单填写更加简洁高效。 