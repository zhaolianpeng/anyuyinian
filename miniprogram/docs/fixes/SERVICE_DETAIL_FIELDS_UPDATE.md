# 服务详情页字段修改总结

## 修改内容

### 1. 删除营养目标字段
- **位置**: `miniprogram/pages/service/detail.js`
- **修改**: 在表单字段过滤逻辑中添加了对营养目标字段的过滤
- **代码变更**:
  ```javascript
  // 过滤掉包含身份证号、出生日期、营养目标等关键词的字段
  const excludeKeywords = ['身份证', 'idCard', 'birthday', '出生', '生日', '营养目标', 'nutritionGoal']
  ```

### 2. 调整疾病史字段位置
- **位置**: `miniprogram/pages/service/detail.js`
- **修改**: 添加了字段排序逻辑，将疾病史字段移到体重字段后面
- **代码变更**:
  ```javascript
  // 重新排序字段：将疾病史字段移到体重字段后面
  formFields = formFields.sort((a, b) => {
    const aName = a.name || ''
    const bName = b.name || ''
    
    // 定义字段顺序
    const fieldOrder = {
      'patientName': 1,
      'patientPhone': 2,
      'patientHeight': 3,
      'patientWeight': 4,
      'medicalHistory': 5, // 疾病史移到体重后面
      'dietaryRestrictions': 6
    }
    
    const aOrder = fieldOrder[aName] || 999
    const bOrder = fieldOrder[bName] || 999
    
    return aOrder - bOrder
  })
  ```

### 3. 修改基础病信息为既往病史
- **位置**: `miniprogram/pages/service/detail.wxml`
- **修改**: 将"基础病信息"字段改为"既往病史"，并调整位置到体重字段下面
- **代码变更**:
  ```xml
  <!-- 既往病史 -->
  <view class="form-item">
    <view class="form-label">
      <text class="label-text">既往病史</text>
    </view>
    <textarea 
      class="form-textarea"
      name="diseaseInfo"
      placeholder="请描述患者的既往病史（选填）"
      value="{{formData.diseaseInfo}}"
      bindinput="onFormInput"
      data-field="diseaseInfo"
      maxlength="500"
    />
    <text class="char-count">{{formData.diseaseInfo.length || 0}}/500</text>
  </view>
  ```

## 修改效果

### 修改前
营养师咨询服务的表单字段顺序：
1. 咨询者姓名
2. 联系电话
3. 身高(cm)
4. 体重(kg)
5. **营养目标** ← 将被删除
6. 疾病史
7. 饮食禁忌
8. **基础病信息** ← 将被修改

### 修改后
营养师咨询服务的表单字段顺序：
1. 咨询者姓名
2. 联系电话
3. 身高(cm)
4. 体重(kg)
5. **疾病史** ← 移到体重后面
6. 饮食禁忌
7. **既往病史** ← 移到体重字段下面，名称改为既往病史

## 技术实现

### 字段过滤逻辑
- 使用关键词匹配来识别需要过滤的字段
- 支持中英文关键词匹配
- 不区分大小写

### 字段排序逻辑
- 使用预定义的字段顺序映射
- 未在映射中的字段按原始顺序排列
- 确保疾病史字段在体重字段之后

### 兼容性
- 保持了对其他服务类型的兼容性
- 不影响现有的表单验证逻辑
- 保持了数据提交的完整性

## 测试验证

### 测试文件
- `miniprogram/tests/test_service_detail_fields.js` - JavaScript测试脚本
- `miniprogram/tests/test_service_detail_fields.html` - HTML测试页面

### 验证内容
1. ✅ 营养目标字段被成功过滤掉
2. ✅ 疾病史字段正确移动到体重字段后面
3. ✅ 基础病信息字段改为既往病史
4. ✅ 既往病史字段移到体重字段下面
5. ✅ 其他字段顺序保持不变
6. ✅ 表单功能正常工作

## 部署说明

### 小程序端
1. 修改已应用到 `miniprogram/pages/service/detail.js`
2. 无需修改数据库配置
3. 无需修改后端API

### 验证步骤
1. 打开小程序，进入营养师咨询服务详情页
2. 查看预约表单，确认营养目标字段已消失
3. 确认疾病史字段在体重字段后面
4. 确认基础病信息字段已改为既往病史
5. 确认既往病史字段显示在体重字段下面
6. 测试表单填写和提交功能

## 注意事项

1. **数据库兼容性**: 修改仅影响前端显示，数据库中的原始配置保持不变
2. **其他服务**: 此修改仅影响营养师咨询服务，其他服务的表单字段不受影响
3. **字段名称**: 过滤逻辑基于字段名称和标签，确保关键词匹配准确
4. **排序逻辑**: 字段排序基于预定义的顺序映射，新增字段需要更新映射

## 后续维护

1. **新增字段**: 如需添加新字段，需要更新 `fieldOrder` 映射
2. **过滤规则**: 如需过滤其他字段，在 `excludeKeywords` 数组中添加关键词
3. **测试验证**: 每次修改后应运行测试脚本验证功能正确性 