# 服务列表布局更新 - 图片居中修复

## 问题描述
首页服务列表中的图片是左对齐的，用户希望图片能够居中显示。

## 解决方案
将服务列表的布局从水平布局（图片在左，内容在右）改为垂直布局（图片在上，内容在下），并确保图片在容器中居中显示。

## 修改内容

### 1. WXML结构修改
**文件**: `pages/index/index.wxml`

**修改前**:
```xml
<view class="service-item">
  <image src="{{item.imageUrl}}" class="service-image" mode="aspectFill" />
  <view class="service-info">
    <!-- 服务信息 -->
  </view>
</view>
```

**修改后**:
```xml
<view class="service-item">
  <view class="service-image-container">
    <image src="{{item.imageUrl}}" class="service-image" mode="aspectFill" />
  </view>
  <view class="service-info">
    <!-- 服务信息 -->
  </view>
</view>
```

### 2. CSS样式修改
**文件**: `pages/index/index.wxss`

#### 服务项容器
**修改前**:
```css
.service-item {
  display: flex;
  background-color: #ffffff;
  border-radius: 15px;
  margin-bottom: 15px;
  overflow: hidden;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  border: 1px solid #f0f0f0;
  align-items: stretch;
}
```

**修改后**:
```css
.service-item {
  display: flex;
  flex-direction: column; /* 改为垂直布局 */
  background-color: #ffffff;
  border-radius: 15px;
  margin-bottom: 15px;
  overflow: hidden;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  border: 1px solid #f0f0f0;
}
```

#### 图片容器
**新增样式**:
```css
.service-image-container {
  width: 100%;
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f8f9fa;
  overflow: hidden;
}
```

#### 图片样式
**修改前**:
```css
.service-image {
  width: 100px;
  height: 100px;
  flex-shrink: 0;
  object-fit: cover;
  object-position: center;
  display: block;
  border-radius: 0;
}
```

**修改后**:
```css
.service-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: center;
  display: block;
}
```

#### 服务信息容器
**修改前**:
```css
.service-info {
  flex: 1;
  padding: 15px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}
```

**修改后**:
```css
.service-info {
  flex: 1;
  padding: 20px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  min-height: 120px;
}
```

#### 预约按钮
**修改前**:
```css
.book-btn {
  align-self: flex-end;
  background: linear-gradient(135deg, #28a745, #20c997);
  color: white;
  border: none;
  border-radius: 20px;
  padding: 8px 16px;
  font-size: 12px;
}
```

**修改后**:
```css
.book-btn {
  align-self: flex-end;
  background: linear-gradient(135deg, #28a745, #20c997);
  color: white;
  border: none;
  border-radius: 20px;
  padding: 10px 20px;
  font-size: 14px;
  margin-top: 10px;
}
```

## 布局变化

### 修改前（水平布局）
```
┌─────────────────────────────────┐
│ [图片] 服务名称                  │
│ 100x100  服务描述               │
│          ¥150元  [预约服务]     │
└─────────────────────────────────┘
```

### 修改后（垂直布局）
```
┌─────────────────────────────────┐
│         [图片]                  │
│      200px高度                  │
├─────────────────────────────────┤
│ 服务名称                        │
│ 服务描述                        │
│ ¥150元        [预约服务]        │
└─────────────────────────────────┘
```

## 优势

1. **图片居中**: 图片现在在容器中完全居中显示
2. **更好的视觉效果**: 垂直布局提供更好的视觉层次
3. **更大的图片显示区域**: 图片从100x100px增加到200px高度
4. **更好的内容组织**: 服务信息有更多的空间显示
5. **响应式友好**: 垂直布局在不同屏幕尺寸下表现更好

## 测试要点

1. **图片居中**: 确认图片在容器中居中显示
2. **图片比例**: 确认图片按比例缩放，不变形
3. **内容布局**: 确认服务信息正确显示
4. **按钮位置**: 确认预约按钮位置正确
5. **响应式**: 在不同设备上测试布局效果

## 注意事项

- 图片容器设置了背景色 `#f8f9fa`，在图片加载失败时提供视觉反馈
- 图片使用 `object-fit: cover` 确保按比例缩放并填满容器
- 服务信息容器设置了最小高度，确保布局一致性
- 预约按钮增加了上边距，提供更好的视觉分离
