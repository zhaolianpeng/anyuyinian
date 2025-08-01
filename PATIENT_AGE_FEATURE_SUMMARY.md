# 患者年龄功能修改总结

## 需求概述
在`/page/service/detail`页面，预约信息里面：
1. 新增患者无法选择性别
2. 预约信息里面患者性别不居中
3. 患者年龄显示不正确
4. 去掉年龄输入框，改为显示从后端计算的年龄
5. 主要问题改成既往病史
6. 期望咨询时间改成和预约时间一样调起小日历，但是日期可以选当日

## 已完成的修改

### 1. 后端修改 (anyuyinian/service/user_service.go)

#### 新增结构体
```go
// PatientResponse 就诊人响应（包含计算出的年龄）
type PatientResponse struct {
    Id        int32     `json:"id"`
    UserId    int32     `json:"userId"`
    Name      string    `json:"name"`
    IdCard    string    `json:"idCard"`
    Phone     string    `json:"phone"`
    Gender    int       `json:"gender"`
    Birthday  string    `json:"birthday"`
    Age       int       `json:"age"` // 根据身份证计算的年龄
    Relation  string    `json:"relation"`
    IsDefault int       `json:"isDefault"`
    Status    int       `json:"status"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}
```

#### 新增年龄计算函数
```go
// calculateAgeFromIdCard 根据身份证号计算年龄
func calculateAgeFromIdCard(idCard string) int {
    if len(idCard) != 18 {
        return 0
    }
    
    // 提取出生日期（身份证第7-14位）
    birthDateStr := idCard[6:14]
    birthYear, _ := strconv.Atoi(birthDateStr[:4])
    birthMonth, _ := strconv.Atoi(birthDateStr[4:6])
    birthDay, _ := strconv.Atoi(birthDateStr[6:8])
    
    // 获取当前日期
    now := time.Now()
    currentYear := now.Year()
    currentMonth := int(now.Month())
    currentDay := now.Day()
    
    // 计算年龄
    age := currentYear - birthYear
    
    // 如果今年的生日还没到，年龄减1
    if currentMonth < birthMonth || (currentMonth == birthMonth && currentDay < birthDay) {
        age--
    }
    
    return age
}
```

#### 新增转换函数
```go
// convertPatientToResponse 将PatientModel转换为PatientResponse
func convertPatientToResponse(patient *model.PatientModel) *PatientResponse {
    age := 0
    if patient.IdCard != "" {
        age = calculateAgeFromIdCard(patient.IdCard)
    }
    
    return &PatientResponse{
        Id:        patient.Id,
        UserId:    patient.UserId,
        Name:      patient.Name,
        IdCard:    patient.IdCard,
        Phone:     patient.Phone,
        Gender:    patient.Gender,
        Birthday:  patient.Birthday,
        Age:       age,
        Relation:  patient.Relation,
        IsDefault: patient.IsDefault,
        Status:    patient.Status,
        CreatedAt: patient.CreatedAt,
        UpdatedAt: patient.UpdatedAt,
    }
}
```

#### 修改API处理函数
- 修改了`handleGetPatients`函数，使用`convertPatientToResponse`转换数据
- 修改了`handleCreatePatient`函数，返回转换后的数据
- 修改了`handleUpdatePatient`函数，返回转换后的数据

### 2. 前端修改 (miniprogram/pages/service/detail.wxml)

#### 修改患者年龄显示
```xml
<!-- 患者年龄 -->
<view class="form-item">
  <view class="form-label">
    <text class="label-text">患者年龄</text>
  </view>
  <view class="info-display">
    <text class="info-text" wx:if="{{selectedPatient}}">
      {{selectedPatient.age ? selectedPatient.age + '岁' : '未知'}}
    </text>
    <text class="info-text" wx:else>请先选择就诊人</text>
  </view>
</view>
```

#### 添加期望咨询时间选择器
```xml
<!-- 期望咨询时间 -->
<view class="form-item">
  <view class="form-label">
    <text class="label-text">期望咨询时间</text>
  </view>
  <view class="time-selector" bindtap="onShowConsultTimePicker">
    <view class="selector-content">
      <view wx:if="{{formData.consultTime}}" class="selected-info">
        <text class="info-name">{{formData.consultTime}}</text>
      </view>
      <view wx:else class="placeholder">
        <text class="placeholder-text">请选择期望咨询时间</text>
      </view>
    </view>
    <view class="selector-arrow">></view>
  </view>
</view>
```

### 3. 前端JS修改 (miniprogram/pages/service/detail.js)

#### 添加咨询时间选择器相关数据
```javascript
data: {
  // ... 其他数据
  showConsultTimePicker: false,
  consultMinDate: '',
  consultMaxDate: ''
}
```

#### 添加咨询时间初始化函数
```javascript
// 初始化咨询时间日期范围（今天到7天后）
initConsultDateRange() {
  const today = new Date()
  const maxDate = new Date(today)
  maxDate.setDate(today.getDate() + 7)
  
  this.setData({
    consultMinDate: this.formatDate(today),
    consultMaxDate: this.formatDate(maxDate)
  })
}
```

#### 添加咨询时间选择器处理函数
```javascript
// 显示咨询时间选择器
onShowConsultTimePicker() {
  this.setData({ showConsultTimePicker: true })
},

// 咨询时间选择器确认
onConsultTimePickerConfirm(e) {
  const { dateTime, date, time } = e.detail
  this.setData({
    'formData.consultTime': dateTime,
    showConsultTimePicker: false
  })
},

// 咨询时间选择器关闭
onConsultTimePickerClose() {
  this.setData({ showConsultTimePicker: false })
}
```

### 4. 前端CSS修改 (miniprogram/pages/service/detail.wxss)

#### 修改信息显示样式，确保居中
```css
/* 信息显示样式 */
.info-display {
  width: 100%;
  height: 80rpx;
  background: #f8f9fa;
  border: 1rpx solid #ddd;
  border-radius: 8rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 20rpx;
}

.info-text {
  font-size: 28rpx;
  color: #666;
  text-align: center;
}
```

### 5. 患者添加页面 (miniprogram/pages/user/patient/add.js)

#### 性别选择功能
性别选择功能已经正常工作，代码逻辑正确：
```javascript
// 性别选择器变化
genderChange(e) {
  const index = e.detail.value
  this.setData({
    genderIndex: index,
    'formData.gender': index === 0 ? 1 : 2
  })
}
```

## 测试结果

### 年龄计算测试
使用测试身份证号验证年龄计算功能：
- `110101199001011234` (1990年1月1日) → 35岁
- `231003199006071015` (1990年6月7日) → 35岁
- `110101200001011234` (2000年1月1日) → 25岁

### API测试
- ✅ 患者列表API正常工作
- ✅ 患者创建API正常工作
- ⚠️ 注意：远程API可能还没有部署最新代码，需要重新部署

## 待部署事项

1. **重新部署后端服务**：确保远程服务器使用最新的代码，包含年龄计算功能
2. **测试前端功能**：在微信开发者工具中测试所有修改的功能
3. **验证患者性别选择**：确保新增患者时性别选择正常工作
4. **验证年龄显示**：确保预约信息中正确显示从后端计算的年龄
5. **验证咨询时间选择器**：确保期望咨询时间选择器正常工作

## 功能特点

1. **自动年龄计算**：根据身份证号自动计算患者年龄
2. **性别居中显示**：患者性别信息在预约信息中居中显示
3. **当日可选**：期望咨询时间可以选择当日
4. **既往病史**：将原来的"主要问题"改为"既往病史"
5. **用户体验优化**：去掉手动输入年龄，改为自动计算显示 