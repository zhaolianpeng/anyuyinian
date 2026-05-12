# 预约页面字段显示调试指南

## 问题描述
用户反馈预约页面没有正常展示要求的字段内容（性别、年龄、既往病史、助排二便）。

## 当前字段状态

### ✅ 已确认存在的字段

#### 1. 患者性别
- **位置**: 预约信息区域
- **显示条件**: 已修改为始终显示
- **数据来源**: 从选中患者获取
- **显示格式**: 男/女/未知

#### 2. 患者年龄
- **位置**: 预约信息区域
- **显示条件**: 已修改为始终显示
- **数据来源**: 从选中患者生日计算
- **显示格式**: XX岁

#### 3. 既往病史
- **位置**: 预约信息区域
- **类型**: 文本输入框
- **字符限制**: 500字符
- **占位符**: "请描述患者的基础病情况（选填）"

#### 4. 助排二便
- **位置**: 预约信息区域
- **类型**: 单选按钮
- **选项**: 需要/不需要
- **默认值**: 无

## 可能的问题原因

### 1. 显示条件问题
**问题**: 性别和年龄字段之前有条件显示 `wx:if="{{selectedPatient}}"`
**解决**: 已修改为始终显示，未选择患者时显示提示文字

### 2. 样式问题
**检查点**:
- 字段是否被CSS隐藏
- 字段是否超出容器范围
- 字段是否被其他元素遮挡

### 3. 数据绑定问题
**检查点**:
- formData是否正确初始化
- 数据更新是否正确
- 事件处理函数是否正常工作

### 4. 页面结构问题
**检查点**:
- 字段是否在正确的容器内
- 字段是否被其他元素覆盖
- 页面布局是否正确

## 调试步骤

### 步骤1: 检查页面结构
```xml
<!-- 预约信息区域 -->
<view class="section">
  <view class="section-header">
    <text class="section-title">预约信息</text>
    <text class="required">*</text>
  </view>
  
  <view class="appointment-section">
    <!-- 预约时间选择 -->
    <view class="datetime-section">...</view>
    
    <!-- 患者性别 -->
    <view class="info-section">
      <text class="label">患者性别</text>
      <text class="info-value" wx:if="{{selectedPatient}}">{{selectedPatient.gender === 1 ? '男' : selectedPatient.gender === 2 ? '女' : '未知'}}</text>
      <text class="info-value" wx:else>请先选择患者</text>
    </view>
    
    <!-- 患者年龄 -->
    <view class="info-section">
      <text class="label">患者年龄</text>
      <text class="info-value" wx:if="{{selectedPatient}}">{{selectedPatient.birthday ? getAge(selectedPatient.birthday) + '岁' : '未知'}}</text>
      <text class="info-value" wx:else>请先选择患者</text>
    </view>
    
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
  </view>
</view>
```

### 步骤2: 检查数据初始化
```javascript
data: {
  formData: {
    appointmentDate: '',
    appointmentTime: '',
    remark: '',
    diseaseInfo: '',        // 既往病史
    needToiletAssist: ''    // 助排二便
  }
}
```

### 步骤3: 检查事件处理函数
```javascript
// 既往病史输入
onDiseaseInput(e) {
  console.log('既往病史输入:', e.detail.value)
  this.setData({
    'formData.diseaseInfo': e.detail.value
  })
}

// 助排二便选择
onToiletAssistChange(e) {
  const value = e.currentTarget.dataset.value
  console.log('助排二便选择:', value)
  this.setData({
    'formData.needToiletAssist': value
  })
}
```

### 步骤4: 检查样式
```css
/* 既往病史 */
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
```

## 测试建议

### 1. 页面加载测试
- 打开预约页面
- 检查预约信息区域是否显示
- 检查既往病史输入框是否可见
- 检查助排二便单选按钮是否可见

### 2. 交互测试
- 点击既往病史输入框，测试输入功能
- 点击助排二便单选按钮，测试选择功能
- 选择患者后，检查性别和年龄是否正确显示

### 3. 数据测试
- 输入既往病史，检查字符计数
- 选择助排二便选项，检查选中状态
- 提交订单时检查数据是否正确包含

## 常见问题解决

### 问题1: 字段不显示
**原因**: 可能是CSS隐藏或条件显示问题
**解决**: 检查字段的显示条件和CSS样式

### 问题2: 字段无法交互
**原因**: 可能是事件绑定问题或样式问题
**解决**: 检查事件处理函数和CSS样式

### 问题3: 数据不保存
**原因**: 可能是数据绑定或更新问题
**解决**: 检查数据初始化和更新逻辑

## 验证清单

- [ ] 预约信息区域显示
- [ ] 患者性别字段显示
- [ ] 患者年龄字段显示
- [ ] 既往病史输入框显示
- [ ] 助排二便单选按钮显示
- [ ] 既往病史可以输入
- [ ] 助排二便可以选择
- [ ] 选择患者后性别年龄正确显示
- [ ] 数据正确保存和提交

如果以上检查都通过，但用户仍然看不到字段，可能是：
1. 页面缓存问题 - 清除缓存重新加载
2. 样式冲突问题 - 检查是否有其他CSS影响
3. 数据加载问题 - 检查网络和数据加载状态 