# 预约页面字段状态总结

## 当前字段状态

### ✅ 已实现的字段

#### 1. 患者信息区域
- **患者选择**: 显示患者列表，支持选择
- **性别**: 自动显示选中患者的性别（男/女）
- **年龄**: 自动计算并显示选中患者的年龄
- **关系**: 显示患者与用户的关系（本人/其他）
- **电话**: 显示患者联系电话

#### 2. 服务地址区域
- **地址选择**: 显示地址列表，支持选择
- **联系人**: 显示地址联系人姓名
- **电话**: 显示地址联系电话
- **详细地址**: 显示完整的地址信息

#### 3. 预约信息区域
- **预约时间**: 日期和时间选择器
- **患者性别**: 显示选中患者的性别（只读）
- **患者年龄**: 显示选中患者的年龄（只读）
- **既往病史**: 文本输入框，支持500字符
- **助排二便**: 单选按钮（需要/不需要）

#### 4. 备注信息区域
- **备注**: 文本输入框，支持200字符

#### 5. 订单信息
- **订单金额**: 显示服务价格

## 字段详细说明

### 性别和年龄字段
```xml
<!-- 患者性别 -->
<view class="info-section" wx:if="{{selectedPatient}}">
  <text class="label">患者性别</text>
  <text class="info-value">{{selectedPatient.gender === 1 ? '男' : selectedPatient.gender === 2 ? '女' : '未知'}}</text>
</view>

<!-- 患者年龄 -->
<view class="info-section" wx:if="{{selectedPatient}}">
  <text class="label">患者年龄</text>
  <text class="info-value">{{selectedPatient.birthday ? getAge(selectedPatient.birthday) + '岁' : '未知'}}</text>
</view>
```

### 既往病史字段
```xml
<!-- 既往病史 -->
<view class="disease-section">
  <text class="label">既往病史</text>
  <textarea class="disease-input" 
            placeholder="请描述患者的既往病史（选填）" 
            value="{{formData.diseaseInfo}}"
            bindinput="onDiseaseInput"
            maxlength="500"
            auto-height="true"></textarea>
  <text class="char-count">{{formData.diseaseInfo.length || 0}}/500</text>
</view>
```

### 助排二便字段
```xml
<!-- 是否需要助排二便 -->
<view class="toilet-section">
  <text class="label">是否需要助排二便</text>
  <view class="radio-group">
    <view class="radio-item {{formData.needToiletAssist === '1' ? 'selected' : ''}}" 
          bindtap="onToiletAssistChange" 
          data-value="1">
      <view class="radio-icon">{{formData.needToiletAssist === '1' ? '●' : '○'}}</view>
      <text class="radio-text">需要</text>
    </view>
    <view class="radio-item {{formData.needToiletAssist === '0' ? 'selected' : ''}}" 
          bindtap="onToiletAssistChange" 
          data-value="0">
      <view class="radio-icon">{{formData.needToiletAssist === '0' ? '●' : '○'}}</view>
      <text class="radio-text">不需要</text>
    </view>
  </view>
</view>
```

## JavaScript 处理函数

### 既往病史输入处理
```javascript
// 既往病史输入
onDiseaseInput(e) {
  console.log('既往病史输入:', e.detail.value)
  this.setData({
    'formData.diseaseInfo': e.detail.value
  })
}
```

### 助排二便选择处理
```javascript
// 助排二便选择
onToiletAssistChange(e) {
  const value = e.currentTarget.dataset.value
  console.log('助排二便选择:', value)
  this.setData({
    'formData.needToiletAssist': value
  })
}
```

### 年龄计算函数
```javascript
// 计算年龄
getAge(birthday) {
  if (!birthday) return ''
  const birthDate = new Date(birthday)
  const today = new Date()
  let age = today.getFullYear() - birthDate.getFullYear()
  const monthDiff = today.getMonth() - birthDate.getMonth()
  
  if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birthDate.getDate())) {
    age--
  }
  
  return age
}
```

## 数据提交

### 表单数据结构
```javascript
formData: {
  appointmentDate: '',
  appointmentTime: '',
  remark: '',
  diseaseInfo: '',        // 既往病史
  needToiletAssist: ''    // 助排二便
}
```

### 订单提交数据
```javascript
const orderData = {
  userId,
  serviceId: this.data.serviceId,
  patientId: this.data.selectedPatient.id,
  addressId: this.data.selectedAddress.id,
  appointmentDate: this.data.formData.appointmentDate,
  appointmentTime: this.data.formData.appointmentTime,
  remark: this.data.formData.remark,
  diseaseInfo: this.data.formData.diseaseInfo,        // 既往病史
  needToiletAssist: this.data.formData.needToiletAssist,  // 助排二便
  formData: this.data.formData
}
```

## 样式支持

### 既往病史样式
```css
.disease-section {
  padding: 20rpx 0;
  border-bottom: 1rpx solid #f0f0f0;
  width: 100%;
  box-sizing: border-box;
  max-width: 100%;
  position: relative;
}

.disease-input {
  width: 100%;
  min-height: 120rpx;
  padding: 20rpx;
  border: 1rpx solid #e0e0e0;
  border-radius: 8rpx;
  font-size: 28rpx;
  line-height: 1.5;
  box-sizing: border-box;
  word-break: break-all;
  max-width: 100%;
  margin-top: 10rpx;
}
```

### 助排二便样式
```css
.toilet-section {
  padding: 20rpx 0;
  border-bottom: 1rpx solid #f0f0f0;
  width: 100%;
  box-sizing: border-box;
  max-width: 100%;
}

.radio-group {
  display: flex;
  gap: 40rpx;
  margin-top: 10rpx;
}

.radio-item {
  display: flex;
  align-items: center;
  padding: 16rpx 24rpx;
  border: 2rpx solid #e0e0e0;
  border-radius: 8rpx;
  transition: all 0.3s;
  cursor: pointer;
}

.radio-item.selected {
  border-color: #007aff;
  background-color: #f0f8ff;
}
```

## 总结

✅ **所有要求的字段都已实现**：
- ✅ 性别（自动显示）
- ✅ 年龄（自动计算显示）
- ✅ 既往病史（可输入）
- ✅ 助排二便（可选择）

✅ **没有营养目标字段**：
- 确认没有营养目标相关的代码需要移除

✅ **字段功能完整**：
- 所有字段都有相应的数据处理
- 所有字段都有合适的样式
- 所有字段都会在订单提交时包含

预约页面已经完全满足您的要求，包含了性别、年龄、既往病史和助排二便字段，并且没有营养目标相关的内容。 