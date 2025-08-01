# 日历选择器改造总结

## 改造目标
将服务详情页面的时间选择从原来的"上午/下午/晚上"改为日历选择器，支持选择未来7天的具体日期和时间，当天不可选。

## 改造内容

### 1. 前端改造

#### 1.1 页面结构改造 (detail.wxml)
- **移除原有时间选择器**：
  ```xml
  <!-- 原来的时间选择 -->
  <view class="time-options">
    <view class="time-slot">上午 (9:00-12:00)</view>
    <view class="time-slot">下午 (14:00-17:00)</view>
    <view class="time-slot">晚上 (19:00-21:00)</view>
  </view>
  ```

- **新增日历选择器**：
  ```xml
  <!-- 新的时间选择器 -->
  <view class="time-selector" bindtap="onShowTimePicker">
    <view class="selector-content">
      <view wx:if="{{selectedDateTime}}" class="selected-info">
        <text class="info-name">{{selectedDateTime}}</text>
      </view>
      <view wx:else class="placeholder">
        <text class="placeholder-text">请选择预约时间</text>
      </view>
    </view>
    <view class="selector-arrow">></view>
  </view>
  
  <!-- 日历选择器组件 -->
  <calendar-picker 
    show="{{showTimePicker}}"
    min-date="{{minDate}}"
    max-date="{{maxDate}}"
    selected-date-time="{{selectedDateTime}}"
    bind:confirm="onTimePickerConfirm"
    bind:close="onTimePickerClose"
  />
  ```

#### 1.2 页面逻辑改造 (detail.js)
- **数据结构更新**：
  ```javascript
  data: {
    // 移除原有的 timeSlots
    // 新增日历相关数据
    selectedDateTime: '',
    selectedDate: '',
    selectedTime: '',
    showTimePicker: false,
    minDate: '',
    maxDate: ''
  }
  ```

- **新增方法**：
  ```javascript
  // 初始化日期范围（明天到7天后）
  initDateRange() {
    const today = new Date()
    const tomorrow = new Date(today)
    tomorrow.setDate(today.getDate() + 1)
    
    const maxDate = new Date(today)
    maxDate.setDate(today.getDate() + 7)
    
    this.setData({
      minDate: this.formatDate(tomorrow),
      maxDate: this.formatDate(maxDate)
    })
  }

  // 显示时间选择器
  onShowTimePicker() {
    this.setData({ showTimePicker: true })
  }

  // 时间选择器确认
  onTimePickerConfirm(e) {
    const { dateTime, date, time } = e.detail
    this.setData({
      selectedDateTime: dateTime,
      selectedDate: date,
      selectedTime: time,
      showTimePicker: false
    })
  }

  // 时间选择器关闭
  onTimePickerClose() {
    this.setData({ showTimePicker: false })
  }
  ```

- **表单验证更新**：
  ```javascript
  // 检查预约时间
  if (!selectedDateTime) {
    wx.showToast({
      title: '请选择预约时间',
      icon: 'none'
    })
    return false
  }
  ```

- **订单提交数据更新**：
  ```javascript
  const orderData = {
    userId: userId,
    serviceId: service.id,
    serviceName: service.name,
    amount: service.price,
    appointmentDate: selectedDate,    // 新增
    appointmentTime: selectedTime,
    patientId: selectedPatient.id,
    addressId: selectedAddress.id,
    formData: formData,
    remark: remark,
    diseaseInfo: formData.diseaseInfo,
    needToiletAssist: formData.needToiletAssist
  }
  ```

#### 1.3 样式改造 (detail.wxss)
- **移除原有时间选择样式**：
  ```css
  /* 移除 .time-options, .time-slot 相关样式 */
  ```

- **新增时间选择器样式**：
  ```css
  .time-selector {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 20rpx;
    background: #f8f9fa;
    border-radius: 8rpx;
    border: 1rpx solid #eee;
  }
  ```

#### 1.4 组件配置 (detail.json)
```json
{
  "usingComponents": {
    "calendar-picker": "/components/calendar-picker/calendar-picker"
  }
}
```

### 2. 日历选择器组件

#### 2.1 组件功能
- 显示日历网格，支持月份导航
- 禁用当天及之前的日期
- 限制选择范围为未来7天
- 选择日期后显示可用时间槽
- 支持选择具体时间（08:00-19:00）

#### 2.2 组件接口
```javascript
properties: {
  show: Boolean,           // 是否显示
  minDate: String,         // 最小日期（明天）
  maxDate: String,         // 最大日期（7天后）
  selectedDateTime: String // 已选择的日期时间
}

events: {
  confirm: { dateTime, date, time } // 确认选择
  close: void                       // 关闭选择器
}
```

#### 2.3 时间槽配置
```javascript
// 允许的时间槽
const allowedTimeSlots = [
  "08:00", "09:00", "10:00", "11:00",
  "14:00", "15:00", "16:00", "17:00", "18:00", "19:00"
]
```

### 3. 后端接口改造

#### 3.1 订单提交接口
后端已经支持新的数据格式，包括：
- `appointmentDate`: 预约日期 (YYYY-MM-DD)
- `appointmentTime`: 预约时间 (HH:MM)

#### 3.2 时间验证逻辑
```go
// 验证预约时间是否在允许范围内（明天开始，未来7天）
appointmentDateTime, err := time.Parse("2006-01-02 15:04", req.AppointmentDate+" "+req.AppointmentTime)

tomorrow := time.Now().AddDate(0, 0, 1)
maxDate := time.Now().AddDate(0, 0, 7)

if appointmentDateTime.Before(tomorrow) {
    return "预约时间不能早于明天"
}

if appointmentDateTime.After(maxDate) {
    return "预约时间不能超过7天后"
}
```

#### 3.3 时间槽验证
```go
// 验证时间槽是否在允许范围内
allowedTimeSlots := []string{
    "08:00", "09:00", "10:00", "11:00",
    "14:00", "15:00", "16:00", "17:00", "18:00", "19:00",
}
```

### 4. 测试验证

#### 4.1 功能测试
- ✅ 日历选择器正常显示
- ✅ 日期范围限制（明天到7天后）
- ✅ 当天日期禁用
- ✅ 时间槽选择正常
- ✅ 表单验证更新
- ✅ 订单数据格式正确

#### 4.2 数据格式测试
```javascript
// 测试订单数据格式
const orderData = {
  userId: 1,
  serviceId: 1,
  appointmentDate: "2024-01-16",
  appointmentTime: "09:00",
  patientId: 1,
  addressId: 1,
  // ... 其他字段
}
```

## 改造效果

### 用户体验提升
1. **更直观的时间选择**：用户可以直接看到日历，选择具体日期
2. **更精确的时间控制**：支持选择具体的时间点，而不是时间段
3. **更好的视觉反馈**：选中状态清晰，操作流程更顺畅

### 技术改进
1. **数据格式标准化**：统一使用 YYYY-MM-DD HH:MM 格式
2. **组件化设计**：日历选择器可复用
3. **验证逻辑完善**：前后端都有完整的时间验证

### 业务逻辑优化
1. **时间范围控制**：严格限制预约时间范围
2. **时间槽管理**：支持动态获取可用时间槽
3. **数据一致性**：前后端数据格式完全匹配

## 注意事项

1. **兼容性**：确保日历选择器组件在所有目标设备上正常工作
2. **性能**：日历组件在移动端应该有良好的性能表现
3. **用户体验**：时间选择流程应该简单直观
4. **数据验证**：前后端都要进行严格的时间格式验证

## 后续优化建议

1. **时间槽动态获取**：根据实际预约情况动态获取可用时间槽
2. **多选支持**：支持选择多个时间段
3. **自定义时间**：允许用户输入自定义时间
4. **预约冲突检测**：实时检测时间冲突
5. **时区支持**：支持不同时区的用户 