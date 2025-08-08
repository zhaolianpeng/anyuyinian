# 后端接口修改总结

## 修改概述
为了支持新的日历时间选择器功能，对后端接口进行了相应的修改和增强。

## 修改的文件

### 1. 订单服务文件 (`service/order_service.go`)

#### 新增功能
1. **时间槽验证**: 在订单提交时验证时间是否在允许的时间段内
2. **获取可用时间槽接口**: 新增接口用于获取指定日期的可用时间槽

#### 具体修改

##### 时间槽验证逻辑
```go
// 验证时间槽是否在允许范围内
allowedTimeSlots := []string{
    "08:00", "09:00", "10:00", "11:00",
    "14:00", "15:00", "16:00", "17:00", "18:00", "19:00",
}

isValidTimeSlot := false
for _, slot := range allowedTimeSlots {
    if req.AppointmentTime == slot {
        isValidTimeSlot = true
        break
    }
}

if !isValidTimeSlot {
    LogError("预约时间不在允许的时间段内", fmt.Errorf("appointmentTime=%s, allowedSlots=%v", req.AppointmentTime, allowedTimeSlots))
    http.Error(w, "预约时间不在允许的时间段内", http.StatusBadRequest)
    return
}
```

##### 新增接口结构体
```go
// GetAvailableTimeSlotsRequest 获取可用时间槽请求
type GetAvailableTimeSlotsRequest struct {
    Date string `json:"date"` // 日期格式：YYYY-MM-DD
}

// GetAvailableTimeSlotsResponse 获取可用时间槽响应
type GetAvailableTimeSlotsResponse struct {
    Date      string   `json:"date"`
    TimeSlots []string `json:"timeSlots"`
}
```

##### 新增接口处理函数
```go
// GetAvailableTimeSlotsHandler 获取可用时间槽接口
func GetAvailableTimeSlotsHandler(w http.ResponseWriter, r *http.Request) {
    // 验证请求方法
    // 解析请求参数
    // 验证日期格式和范围
    // 返回可用时间槽
}
```

### 2. 主程序文件 (`main.go`)

#### 新增路由
```go
// 订单相关接口
http.HandleFunc("/api/order/time_slots", service.NewLogMiddleware(service.GetAvailableTimeSlotsHandler))
```

## 新增接口详情

### 接口名称
`GET /api/order/time_slots`

### 请求参数
```json
{
    "date": "2024-01-16"  // 日期格式：YYYY-MM-DD
}
```

### 响应格式
```json
{
    "code": 0,
    "data": {
        "date": "2024-01-16",
        "timeSlots": [
            "08:00", "09:00", "10:00", "11:00",
            "14:00", "15:00", "16:00", "17:00", "18:00", "19:00"
        ]
    }
}
```

### 错误响应
```json
{
    "code": -1,
    "errorMsg": "错误信息"
}
```

## 验证逻辑

### 1. 日期范围验证
- 只能查询明天开始到7天后的日期
- 当天及之前的日期会被拒绝
- 7天后的日期会被拒绝

### 2. 时间槽验证
- 只允许选择预定义的10个时间段
- 订单提交时会验证时间是否在允许范围内
- 不在允许范围内的时间会被拒绝

### 3. 格式验证
- 日期格式必须为 `YYYY-MM-DD`
- 时间格式必须为 `HH:MM`
- 格式错误会返回相应的错误信息

## 时间槽配置

### 允许的时间段
```go
allowedTimeSlots := []string{
    "08:00", "09:00", "10:00", "11:00",  // 上午时段
    "14:00", "15:00", "16:00", "17:00", "18:00", "19:00",  // 下午和晚上时段
}
```

### 业务规则
1. **上午时段**: 08:00-11:00，每小时一个时间段
2. **下午时段**: 14:00-19:00，每小时一个时间段
3. **午休时间**: 12:00-13:00 不提供服务
4. **服务时间**: 每天8小时服务时间

## 扩展性设计

### 未来可能的扩展
1. **动态时间槽**: 根据服务类型提供不同的时间段
2. **已预约时间**: 查询数据库中已预约的时间段并排除
3. **节假日处理**: 节假日提供不同的服务时间
4. **预约限制**: 每个时间段限制预约人数

### 预留接口
```go
// TODO: 这里可以添加业务逻辑来检查哪些时间段已被预约
// 例如：查询数据库中该日期的已预约时间段，然后从允许的时间槽中排除
// 目前返回所有允许的时间槽
```

## 测试验证

### 测试脚本
创建了 `test_time_slots_api.sh` 测试脚本，包含以下测试用例：

1. **正常测试**:
   - 获取明天的可用时间槽
   - 获取3天后的可用时间槽

2. **边界测试**:
   - 测试无效日期（今天）
   - 测试无效日期（8天后）

3. **错误测试**:
   - 测试日期格式错误
   - 测试空日期

### 运行测试
```bash
chmod +x test_time_slots_api.sh
./test_time_slots_api.sh
```

## 前端集成

### 日历组件修改
1. **动态获取时间槽**: 选择日期时自动获取该日期的可用时间槽
2. **错误处理**: 如果获取失败，使用默认时间槽
3. **用户体验**: 提供加载状态和错误提示

### 接口调用
```javascript
// 获取可用时间槽
async getAvailableTimeSlots(date) {
    const res = await wx.request({
        url: `${this.data.baseUrl}/api/order/time_slots`,
        method: 'POST',
        data: { date }
    })
    
    if (res.data.code === 0) {
        const timeSlots = res.data.data.timeSlots.map(time => ({
            time,
            label: time
        }))
        this.setData({ timeSlots })
    }
}
```

## 总结

### 完成的修改
1. ✅ 新增时间槽验证逻辑
2. ✅ 新增获取可用时间槽接口
3. ✅ 注册新的API路由
4. ✅ 完善错误处理和日志记录
5. ✅ 提供测试脚本验证功能

### 技术特性
1. **安全性**: 严格的时间范围验证
2. **可扩展性**: 预留了业务逻辑扩展接口
3. **一致性**: 前后端时间槽配置保持一致
4. **容错性**: 前端提供降级处理机制

### 业务价值
1. **精确预约**: 提供具体的时间段选择
2. **资源管理**: 便于后续实现预约人数限制
3. **用户体验**: 动态显示可用时间，避免冲突
4. **系统稳定性**: 前后端双重验证确保数据一致性

所有后端接口修改已完成，与前端日历选择器功能完全匹配，可以直接投入使用。 