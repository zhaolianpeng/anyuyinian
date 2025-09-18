# 咨询时间选择修复

## 问题描述
期望咨询时间选择不正确，用户期望从当前时间开始，往后推7天，但实际显示的时间范围不符合预期。

## 问题分析

### 1. 时区问题
- JavaScript的 `new Date()` 可能受到时区影响
- 日期比较时没有统一时区处理
- 导致"今天"的判断不准确

### 2. 日期比较逻辑
- 日历组件的日期比较逻辑需要统一时区
- 最小日期和最大日期的设置需要精确到天

## 修复方案

### 1. 修复咨询时间范围设置

**文件**: `miniprogram/pages/service/detail.js`

#### 修改前
```javascript
initConsultDateRange() {
  const today = new Date()
  const maxDate = new Date(today)
  maxDate.setDate(today.getDate() + 6) // 从今天开始，共7天（今天+6天）
  // ...
}
```

#### 修改后
```javascript
initConsultDateRange() {
  const today = new Date()
  // 设置时间为00:00:00，避免时区问题
  today.setHours(0, 0, 0, 0)
  
  const maxDate = new Date(today)
  maxDate.setDate(today.getDate() + 6) // 从今天开始，共7天（今天+6天）
  
  console.log('设置咨询时间范围:', {
    start: this.formatDate(today),
    end: this.formatDate(maxDate),
    today: today.toISOString(),
    maxDate: maxDate.toISOString()
  })
  // ...
}
```

### 2. 修复日历组件日期比较逻辑

**文件**: `miniprogram/components/calendar-picker/calendar-picker.js`

#### 修改前
```javascript
const today = new Date()
const minDate = this.data.minDate ? new Date(this.data.minDate) : null
const maxDate = this.data.maxDate ? new Date(this.data.maxDate) : null

// 日期比较
const isDisabled = (minDate && date < minDate) || (maxDate && date > maxDate)
```

#### 修改后
```javascript
const today = new Date()
today.setHours(0, 0, 0, 0) // 设置为00:00:00，避免时区问题

const minDate = this.data.minDate ? new Date(this.data.minDate) : null
const maxDate = this.data.maxDate ? new Date(this.data.maxDate) : null

if (minDate) {
  minDate.setHours(0, 0, 0, 0)
}
if (maxDate) {
  maxDate.setHours(0, 0, 0, 0)
}

// 每个日期也设置为00:00:00
date.setHours(0, 0, 0, 0)

// 日期比较
const isDisabled = (minDate && date < minDate) || (maxDate && date > maxDate)
```

## 修复内容

### 1. 时区统一处理
- 所有日期对象都设置为 `00:00:00`，避免时区问题
- 确保日期比较的准确性

### 2. 调试日志增强
- 添加详细的日期范围日志
- 包含ISO格式的时间戳，便于调试

### 3. 日期比较逻辑优化
- 统一所有日期的时区处理
- 确保最小日期和最大日期的比较准确

## 预期效果

### 咨询时间选择
- **开始时间**: 今天（当前日期）
- **结束时间**: 今天 + 6天
- **总天数**: 7天
- **可选日期**: 今天、明天、后天...到第7天

### 预约时间选择
- **开始时间**: 明天（当前日期 + 1天）
- **结束时间**: 今天 + 7天
- **总天数**: 7天
- **可选日期**: 明天、后天...到第8天

## 调试信息

### 控制台日志
1. **咨询时间范围设置**:
   ```
   设置咨询时间范围: {
     start: "2025-09-18",
     end: "2025-09-24",
     today: "2025-09-18T00:00:00.000Z",
     maxDate: "2025-09-24T00:00:00.000Z"
   }
   ```

2. **日历组件日期范围**:
   ```
   日历组件日期范围: {
     minDate: "2025-09-18T00:00:00.000Z",
     maxDate: "2025-09-24T00:00:00.000Z",
     today: "2025-09-18T00:00:00.000Z"
   }
   ```

## 相关文件

- `miniprogram/pages/service/detail.js` - 服务详情页面
- `miniprogram/components/calendar-picker/calendar-picker.js` - 日历选择器组件
- `miniprogram/pages/service/detail.wxml` - 服务详情页面模板

## 测试建议

1. **检查控制台日志**
   - 查看咨询时间范围设置是否正确
   - 确认日期格式和时区处理

2. **测试日期选择**
   - 验证今天是否可以选择
   - 验证7天后的日期是否可以选择
   - 验证超出范围的日期是否被禁用

3. **验证时间范围**
   - 咨询时间：今天到第7天
   - 预约时间：明天到第8天
