# 时间范围选择修复

## 问题描述
用户希望修改订单预约页面的时间选择逻辑：
- **咨询时间**：从当前时间开始，往后推7天
- **预约时间**：从当前时间的第二天开始，往后推7天

## 修改内容

### 1. 订单页面预约时间范围修改

**文件**: `miniprogram/pages/order/order.js`

#### 修改前
```javascript
// 设置预约日期范围
setAppointmentDateRange() {
  const today = new Date()
  const startDate = new Date(today.getTime() + 24 * 60 * 60 * 1000) // 明天
  const endDate = new Date(today.getTime() + 30 * 24 * 60 * 60 * 1000) // 30天后
  // ...
}
```

#### 修改后
```javascript
// 设置预约日期范围（从明天开始，往后推7天）
setAppointmentDateRange() {
  const today = new Date()
  const startDate = new Date(today.getTime() + 24 * 60 * 60 * 1000) // 明天
  const endDate = new Date(today.getTime() + 7 * 24 * 60 * 60 * 1000) // 明天+6天，共7天
  
  console.log('设置预约时间范围:', {
    start: formatDate(startDate),
    end: formatDate(endDate)
  })
  // ...
}
```

### 2. 服务详情页面咨询时间范围修改

**文件**: `miniprogram/pages/service/detail.js`

#### 修改前
```javascript
// 初始化咨询时间日期范围（今天到7天后）
initConsultDateRange() {
  const today = new Date()
  const maxDate = new Date(today)
  maxDate.setDate(today.getDate() + 6) // 从今天开始，共7天（今天+6天）
  // ...
}
```

#### 修改后
```javascript
// 初始化咨询时间日期范围（从当前时间开始，往后推7天）
initConsultDateRange() {
  const today = new Date()
  const maxDate = new Date(today)
  maxDate.setDate(today.getDate() + 6) // 从今天开始，共7天（今天+6天）
  
  console.log('设置咨询时间范围:', {
    start: this.formatDate(today),
    end: this.formatDate(maxDate)
  })
  // ...
}
```

### 3. 服务详情页面预约时间范围修改

**文件**: `miniprogram/pages/service/detail.js`

#### 修改前
```javascript
// 初始化日期范围
initDateRange() {
  const today = new Date()
  const tomorrow = new Date(today)
  tomorrow.setDate(today.getDate() + 1)
  
  const maxDate = new Date(today)
  maxDate.setDate(today.getDate() + 7) // 从明天开始，共7天（明天+6天）
  // ...
}
```

#### 修改后
```javascript
// 初始化预约时间日期范围（从明天开始，往后推7天）
initDateRange() {
  const today = new Date()
  const tomorrow = new Date(today)
  tomorrow.setDate(today.getDate() + 1)
  
  const maxDate = new Date(today)
  maxDate.setDate(today.getDate() + 7) // 从明天开始，共7天（明天+6天）
  
  console.log('设置预约时间范围:', {
    start: this.formatDate(tomorrow),
    end: this.formatDate(maxDate)
  })
  // ...
}
```

## 时间范围说明

### 咨询时间范围
- **开始时间**: 当前时间（今天）
- **结束时间**: 当前时间 + 6天
- **总天数**: 7天
- **用途**: 用户选择期望的咨询时间

### 预约时间范围
- **开始时间**: 当前时间的第二天（明天）
- **结束时间**: 当前时间 + 7天
- **总天数**: 7天
- **用途**: 用户选择具体的预约时间

## 修改效果

### 1. 咨询时间选择
- 用户可以选择今天到未来6天内的任意时间进行咨询
- 时间范围更加灵活，支持当天咨询

### 2. 预约时间选择
- 用户可以选择明天到未来7天内的任意时间进行预约
- 避免当天预约，给服务提供方准备时间

### 3. 调试日志
- 添加了控制台日志，方便调试时间范围设置
- 可以在开发者工具中查看时间范围是否正确

## 相关文件

- `miniprogram/pages/order/order.js` - 订单页面
- `miniprogram/pages/service/detail.js` - 服务详情页面
- `miniprogram/pages/order/order.wxml` - 订单页面模板
- `miniprogram/pages/service/detail.wxml` - 服务详情页面模板

## 测试建议

1. **咨询时间测试**
   - 进入服务详情页面
   - 点击"期望咨询时间"选择器
   - 验证可选择的时间范围是今天到未来6天

2. **预约时间测试**
   - 进入订单页面或服务详情页面
   - 点击"预约时间"选择器
   - 验证可选择的时间范围是明天到未来7天

3. **时间范围验证**
   - 检查控制台日志输出
   - 验证时间范围计算是否正确
   - 确认日期格式正确

## 注意事项

1. **时区处理**: 所有时间计算都基于本地时区
2. **日期格式**: 统一使用 `YYYY-MM-DD` 格式
3. **边界处理**: 确保时间范围计算准确
4. **用户体验**: 时间选择器界面保持友好
