# 时间槽查询接口修复

## 问题描述
咨询时间选择器允许选择今天，但后端时间槽查询接口只允许查询"明天开始的日期"，导致选择今天时出现错误：

```
HTTP 400: 只能查询明天开始的日期
```

## 问题分析

### 1. 前后端不一致
- **前端**: 咨询时间选择器允许选择今天到未来7天
- **后端**: 时间槽查询接口只允许查询明天到未来7天

### 2. 业务逻辑冲突
- 咨询时间应该支持今天选择（紧急咨询需求）
- 但时间槽查询接口限制了今天的选择

## 修复方案

### 修改后端时间槽查询接口

**文件**: `service/order_service.go`

#### 修改前
```go
// 验证日期范围（明天开始，未来7天）
tomorrow := time.Now().AddDate(0, 0, 1)
tomorrow = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())

maxDate := time.Now().AddDate(0, 0, 7)
maxDate = time.Date(maxDate.Year(), maxDate.Month(), maxDate.Day(), 23, 59, 59, 0, maxDate.Location())

if requestDateTime.Before(tomorrow) {
    LogError("请求日期过早", fmt.Errorf("requestDate=%v, tomorrow=%v", requestDateTime, tomorrow))
    http.Error(w, "只能查询明天开始的日期", http.StatusBadRequest)
    return
}
```

#### 修改后
```go
// 验证日期范围（今天开始，未来7天）
today := time.Now()
today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

maxDate := time.Now().AddDate(0, 0, 6) // 今天+6天，共7天
maxDate = time.Date(maxDate.Year(), maxDate.Month(), maxDate.Day(), 23, 59, 59, 0, maxDate.Location())

if requestDateTime.Before(today) {
    LogError("请求日期过早", fmt.Errorf("requestDate=%v, today=%v", requestDateTime, today))
    http.Error(w, "只能查询今天开始的日期", http.StatusBadRequest)
    return
}
```

## 修复内容

### 1. 日期范围调整
- **开始时间**: 从明天改为今天
- **结束时间**: 从今天+7天改为今天+6天（保持7天总长度）
- **错误信息**: 从"只能查询明天开始的日期"改为"只能查询今天开始的日期"

### 2. 业务逻辑统一
- 前端咨询时间选择器：今天到未来7天
- 后端时间槽查询接口：今天到未来7天
- 保持前后端逻辑一致

### 3. 时间范围说明

| 时间类型 | 开始时间 | 结束时间 | 总天数 | 说明 |
|---------|---------|---------|--------|------|
| **咨询时间** | 今天 | 今天+6天 | 7天 | 支持当天咨询 |
| **预约时间** | 明天 | 今天+7天 | 7天 | 从第二天开始预约 |
| **时间槽查询** | 今天 | 今天+6天 | 7天 | 支持当天时间槽查询 |

## 修复效果

### 1. 咨询时间选择
- 用户可以选择今天进行咨询
- 时间槽查询接口支持今天的时间槽
- 前后端逻辑完全一致

### 2. 预约时间选择
- 用户仍然只能从明天开始预约
- 保持原有的业务逻辑

### 3. 错误处理
- 不再出现"只能查询明天开始的日期"错误
- 支持今天的时间槽查询

## 相关文件

- `service/order_service.go` - 后端时间槽查询接口
- `miniprogram/components/calendar-picker/calendar-picker.js` - 前端日历选择器
- `miniprogram/pages/service/detail.js` - 服务详情页面

## 测试建议

### 1. 咨询时间测试
- 选择今天作为咨询时间
- 验证时间槽查询成功
- 确认可以正常选择时间

### 2. 预约时间测试
- 选择明天作为预约时间
- 验证时间槽查询成功
- 确认可以正常选择时间

### 3. 边界测试
- 选择今天+6天（最后一天）
- 选择今天+7天（超出范围）
- 验证错误处理正确

## 注意事项

1. **业务逻辑**: 咨询时间支持当天，预约时间从明天开始
2. **时间槽**: 所有时间槽查询都支持今天到未来7天
3. **错误处理**: 超出范围的日期会返回相应的错误信息
4. **一致性**: 前后端时间范围逻辑完全一致
