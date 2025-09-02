# 智慧养老功能实现文档

## 概述

本文档描述了为小程序添加"智慧养老"服务类型的完整实现，包括智能养老设备的购买流程、视频播放功能等。

## 功能特性

### 1. 智慧养老设备类型
- **智能血压计**: 24小时血压监测，数据同步到手机APP
- **智能血糖仪**: 无痛采血，快速检测，数据云端存储
- **智能手环**: 心率监测、睡眠分析、跌倒检测、GPS定位
- **智能药盒**: 定时提醒服药，记录服药情况
- **智能摄像头**: 远程监控，异常行为识别，双向语音通话
- **智能门锁**: 指纹识别、密码开锁、远程授权
- **智能马桶**: 自动冲洗、座圈加热、健康检测
- **智能床垫**: 睡眠质量监测、呼吸监测、离床报警

### 2. 核心功能
- ✅ 智慧养老设备展示
- ✅ 产品介绍视频播放
- ✅ 极简购买流程（只需选择收货地址，无需填写任何表单）
- ✅ 收货地址管理
- ✅ 订单创建和支付

## 技术实现

### 数据库变更

#### 1. 添加videoUrl字段
```sql
-- 迁移文件: 025_add_video_url_to_service_items.sql
ALTER TABLE ServiceItems ADD COLUMN videoUrl VARCHAR(500) COMMENT '视频URL，用于智慧养老设备的产品介绍视频';
```

#### 2. 插入智慧养老设备数据
```sql
-- 迁移文件: 026_add_smart_elderly_care_services.sql
INSERT INTO ServiceItems (name, description, category, price, originalPrice, imageUrl, detailImages, formConfig, status, sort, createdAt, updatedAt) VALUES
('智能血压计', '24小时血压监测，数据同步到手机APP，异常自动报警', '智慧养老', 299.00, 399.00, '/static/smart-elderly/blood-pressure-monitor.jpg', '["/static/smart-elderly/bp1.jpg","/static/smart-elderly/bp2.jpg","/static/smart-elderly/bp3.jpg"]', '{"fields":[{"name":"deliveryAddress","label":"收货地址","type":"textarea","required":true,"placeholder":"请输入详细收货地址"}]}', 1, 1, NOW(), NOW()),
-- ... 其他设备数据
```

### 后端API实现

#### 1. 新增智慧养老订单API
- **端点**: `POST /api/order/smart-elderly`
- **文件**: `service/smart_elderly_service.go`
- **功能**: 专门处理智慧养老设备的订单创建

#### 2. 订单模型支持
- 修改 `OrderModel` 以支持可选的 `PatientId`、`AppointmentDate`、`AppointmentTime`
- 智慧养老设备订单这些字段为 `nil`

#### 3. 首页API集成
- 在 `home_init_service.go` 中添加智慧养老服务查询
- 自动在首页展示智慧养老设备

### 4. 服务分类API
- 新增 `/api/service/categories` 接口，动态返回可用的服务分类
- 只显示有服务的分类，智慧养老分类只在有对应服务时才显示
- 前端服务列表页面动态加载分类，不再硬编码分类列表
- 前端首页页面动态加载分类，支持智慧养老分类显示

### 前端实现

#### 1. 服务详情页面改造
- **文件**: `pages/service/detail.js` 和 `pages/service/detail.wxml`
- **功能**:
  - 根据服务类型动态显示/隐藏患者信息选择
  - 根据服务类型动态显示/隐藏预约时间选择
  - 根据服务类型动态显示/隐藏购买信息表单（智慧养老设备不显示）
  - 添加视频播放组件
  - 修改表单验证逻辑（智慧养老设备只需验证地址）
  - 修改订单提交逻辑

#### 2. 视频播放功能
```xml
<!-- 智慧养老设备视频播放 -->
<view class="video-section" wx:if="{{service.category === '智慧养老' && service.videoUrl}}">
  <view class="section-title">产品介绍视频</view>
  <video 
    class="product-video" 
    src="{{service.videoUrl}}" 
    controls="{{true}}"
    show-center-play-btn="{{true}}"
    show-play-btn="{{true}}"
    show-fullscreen-btn="{{true}}"
    show-progress="{{true}}"
    enable-progress-gesture="{{true}}"
    object-fit="contain"
  ></video>
</view>
```

#### 3. 表单验证逻辑
```javascript
// 判断是否为智慧养老设备
const isSmartElderly = service.category === '智慧养老'

if (isSmartElderly) {
  // 智慧养老设备只需要验证地址
  if (!selectedAddress) {
    wx.showToast({ title: '请选择收货地址', icon: 'none' })
    return false
  }
  
  if (!formData.deliveryAddress) {
    wx.showToast({ title: '请填写详细收货地址', icon: 'none' })
    return false
  }
  
  return true
} else {
  // 普通服务需要验证患者信息和预约时间
  // ... 原有验证逻辑
}
```

#### 4. API调用
```javascript
// 根据服务类型选择不同的API端点
let result
if (isSmartElderly) {
  result = await api.smartElderlyOrderSubmit(orderData)
} else {
  result = await api.orderSubmit(orderData)
}
```

## 文件清单

### 后端文件
- `db/migration/025_add_video_url_to_service_items.sql` - 添加视频URL字段
- `db/migration/026_add_smart_elderly_care_services.sql` - 插入智慧养老设备数据
- `db/model/service.go` - 更新ServiceItemModel
- `db/model/order.go` - 更新OrderModel
- `service/smart_elderly_service.go` - 智慧养老订单服务
- `service/home_init_service.go` - 首页服务集成
- `main.go` - 注册新API端点

### 前端文件
- `pages/service/detail.js` - 服务详情页面逻辑
- `pages/service/detail.wxml` - 服务详情页面模板
- `pages/service/detail.wxss` - 服务详情页面样式
- `pages/service/list.js` - 服务列表页面逻辑（动态加载分类）
- `pages/index/index.js` - 首页页面逻辑（动态加载分类）
- `utils/cloud-container-standard.js` - API配置

### 测试文件
- `tests/smart_elderly/test_smart_elderly_services.sh` - 功能测试脚本
- `tests/smart_elderly/test_categories_api.sh` - 分类API测试脚本
- `tests/smart_elderly/test_homepage_categories.sh` - 首页分类加载测试脚本

## 使用流程

### 1. 用户浏览智慧养老设备
1. 用户进入小程序首页
2. 在服务列表中看到智慧养老设备
3. 点击设备查看详情

### 2. 查看设备详情和视频
1. 查看设备基本信息（价格、描述等）
2. 观看产品介绍视频
3. 了解设备功能特点

### 3. 购买设备
1. 选择收货地址
2. 填写详细收货地址
3. 点击"立即购买"
4. 完成支付

### 4. 订单管理
1. 在订单列表中查看智慧养老设备订单
2. 订单状态跟踪
3. 收货确认

## 测试验证

运行测试脚本验证功能：
```bash
cd anyuyinian/tests/smart_elderly/

# 测试智慧养老功能
chmod +x test_smart_elderly_services.sh
./test_smart_elderly_services.sh

# 测试分类API
chmod +x test_categories_api.sh
./test_categories_api.sh

# 测试首页分类加载
chmod +x test_homepage_categories.sh
./test_homepage_categories.sh

# 测试兼容性
chmod +x test_compatibility.sh
./test_compatibility.sh
```

## 兼容性保证

### 1. 向后兼容性
- ✅ **普通服务订单流程**: 完全不受影响，保持原有逻辑
- ✅ **订单列表API**: 兼容智慧养老设备订单，正确处理nil值
- ✅ **订单详情API**: 兼容智慧养老设备订单，正确处理可选字段
- ✅ **服务详情API**: 兼容智慧养老设备，支持videoUrl字段

### 2. API隔离性
- ✅ **智慧养老设备**: 使用专门的 `/api/order/smart-elderly` 接口
- ✅ **普通服务**: 继续使用原有的 `/api/order/submit` 接口
- ✅ **服务类型验证**: 确保智慧养老设备不能使用普通API，普通服务不能使用智慧养老API

### 3. 数据库兼容性
- ✅ **可选字段**: PatientId、AppointmentDate、AppointmentTime 改为可选
- ✅ **订单号区分**: 智慧养老设备订单号以"SE"开头
- ✅ **数据完整性**: 保持现有订单数据的完整性

## 注意事项

1. **数据库迁移**: 需要先执行数据库迁移脚本
2. **视频资源**: 需要上传智慧养老设备的介绍视频到云存储
3. **图片资源**: 需要上传设备图片到云存储
4. **订单号格式**: 智慧养老设备订单号以"SE"开头，便于区分
5. **兼容性测试**: 建议运行兼容性测试脚本验证功能

## 扩展功能

未来可以考虑添加的功能：
- 设备使用教程视频
- 设备安装指导
- 设备维护提醒
- 设备数据同步
- 家人远程监控功能
