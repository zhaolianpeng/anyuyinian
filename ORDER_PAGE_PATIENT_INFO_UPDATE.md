# 预约页面患者信息更新总结

## 修改内容

### 1. 文本更新
- **就诊人信息** → **患者信息**
- **添加就诊人** → **添加患者**
- **暂无就诊人信息** → **暂无患者信息**
- **请选择就诊人** → **请选择患者**

### 2. 新增字段

#### 患者信息显示
- **患者性别**：显示选中患者的性别（男/女/未知）
- **患者年龄**：根据生日自动计算并显示年龄

#### 预约信息新增字段
- **基础病信息**：文本输入框，最多500字符
- **是否需要助排二便**：单选按钮（需要/不需要）

### 3. 功能增强

#### 年龄计算功能
```javascript
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

#### 新增输入处理函数
```javascript
// 基础病信息输入
onDiseaseInput(e) {
  this.setData({
    'formData.diseaseInfo': e.detail.value
  })
}

// 助排二便选择
onToiletAssistChange(e) {
  const value = e.currentTarget.dataset.value
  this.setData({
    'formData.needToiletAssist': value
  })
}
```

### 4. 数据结构更新

#### 表单数据新增字段
```javascript
formData: {
  appointmentDate: '',
  appointmentTime: '',
  remark: '',
  diseaseInfo: '',        // 新增：基础病信息
  needToiletAssist: ''    // 新增：是否需要助排二便
}
```

#### 订单提交数据
```javascript
const orderData = {
  userId,
  serviceId: this.data.serviceId,
  patientId: this.data.selectedPatient.id,
  addressId: this.data.selectedAddress.id,
  appointmentDate: this.data.formData.appointmentDate,
  appointmentTime: this.data.formData.appointmentTime,
  remark: this.data.formData.remark,
  diseaseInfo: this.data.formData.diseaseInfo,        // 新增
  needToiletAssist: this.data.formData.needToiletAssist, // 新增
  formData: this.data.formData
}
```

### 5. 界面更新

#### 患者信息显示
```xml
<view class="patient-info">
  <text class="patient-name">{{item.name}}</text>
  <text class="patient-gender">{{item.gender === 1 ? '男' : item.gender === 2 ? '女' : '未知'}}</text>
  <text class="patient-age">{{item.birthday ? getAge(item.birthday) + '岁' : ''}}</text>
  <text class="patient-relation">{{item.relation || '本人'}}</text>
  <text class="patient-phone">{{item.phone}}</text>
</view>
```

#### 预约信息新增字段
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

<!-- 基础病信息 -->
<view class="disease-section">
  <text class="label">基础病信息</text>
  <textarea class="disease-input" 
            placeholder="请描述患者的基础病情况（选填）" 
            value="{{formData.diseaseInfo}}"
            bindinput="onDiseaseInput"
            maxlength="500"
            auto-height="true"></textarea>
  <text class="char-count">{{formData.diseaseInfo.length || 0}}/500</text>
</view>

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

### 6. 样式更新

#### 新增CSS样式
```css
/* 患者性别和年龄标签 */
.patient-gender, .patient-age {
  font-size: 22rpx;
  color: #666;
  background-color: #f0f0f0;
  padding: 4rpx 12rpx;
  border-radius: 4rpx;
  margin-right: 10rpx;
  display: inline-block;
  max-width: 100%;
}

/* 患者信息显示 */
.info-section {
  display: flex;
  align-items: center;
  padding: 20rpx 0;
  border-bottom: 1rpx solid #f0f0f0;
  width: 100%;
  box-sizing: border-box;
  max-width: 100%;
}

.info-value {
  flex: 1;
  font-size: 28rpx;
  color: #333;
  text-align: right;
  margin-left: 20rpx;
}

/* 基础病信息 */
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

/* 助排二便选择 */
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

.radio-icon {
  font-size: 32rpx;
  color: #007aff;
  margin-right: 12rpx;
  width: 32rpx;
  text-align: center;
}

.radio-text {
  font-size: 28rpx;
  color: #333;
}
```

## 功能特点

### 1. 自动年龄计算
- 根据患者生日自动计算年龄
- 精确到月份，避免虚岁计算错误
- 处理生日未填写的情况

### 2. 性别显示
- 根据数据库中的性别字段（1=男，2=女）显示中文
- 处理未知性别的情况

### 3. 基础病信息
- 支持多行文本输入
- 自动高度调整
- 字符计数显示（最多500字符）
- 选填字段，不影响订单提交

### 4. 助排二便选择
- 单选按钮设计
- 选中状态视觉反馈
- 支持触摸选择
- 数据格式：'1'=需要，'0'=不需要

### 5. 响应式设计
- 小屏幕设备适配
- 大屏幕设备优化
- 按钮间距自适应

## 数据流程

1. **患者选择** → 自动显示性别和年龄
2. **基础病信息输入** → 实时保存到formData
3. **助排二便选择** → 实时保存到formData
4. **订单提交** → 包含所有新增字段

## 注意事项

1. **年龄计算**：需要确保患者生日数据格式正确（YYYY-MM-DD）
2. **性别显示**：依赖数据库中的性别字段，需要确保数据一致性
3. **字符限制**：基础病信息限制500字符，避免数据过长
4. **必填验证**：新增字段都是选填，不影响订单提交
5. **数据存储**：所有新增字段都会保存到订单数据中

## 测试建议

1. **患者信息显示**：测试不同性别和生日的显示效果
2. **年龄计算**：测试各种生日情况下的年龄计算准确性
3. **基础病信息**：测试长文本输入和字符计数
4. **助排二便选择**：测试单选按钮的交互效果
5. **订单提交**：确认新增字段正确提交到后端 