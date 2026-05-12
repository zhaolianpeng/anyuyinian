# 图片居中修复测试

## 修复内容

### 1. 首页服务列表图片居中修复

**文件**: `pages/index/index.wxss`

**修复前**:
```css
.service-image {
  width: 100px;
  height: 100px;
  flex-shrink: 0;
}
```

**修复后**:
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

**同时优化了服务项容器**:
```css
.service-item {
  display: flex;
  background-color: #ffffff;
  border-radius: 15px;
  margin-bottom: 15px;
  overflow: hidden;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  border: 1px solid #f0f0f0;
  align-items: stretch; /* 新增 */
}
```

### 2. 服务列表页面图片居中修复

**文件**: `pages/service/list.wxss`

**修复前**:
```css
.service-image {
  width: 100%;
  height: 300rpx;
  overflow: hidden;
}

.service-image image {
  width: 100%;
  height: 100%;
}
```

**修复后**:
```css
.service-image {
  width: 100%;
  height: 300rpx;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.service-image image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: center;
}
```

## 修复说明

### 首页服务列表
- 添加了 `object-fit: cover` 确保图片按比例缩放并填满容器
- 添加了 `object-position: center` 确保图片在容器中居中显示
- 添加了 `display: block` 确保图片作为块级元素显示
- 设置 `border-radius: 0` 避免圆角影响图片显示
- 在服务项容器中添加 `align-items: stretch` 确保子元素对齐

### 服务列表页面
- 在图片容器中添加了 `display: flex`、`align-items: center`、`justify-content: center` 确保图片在容器中居中
- 为图片元素添加了 `object-fit: cover` 和 `object-position: center` 确保图片内容居中显示

## 测试方法

1. **首页测试**:
   - 打开小程序首页
   - 查看服务列表中的图片是否居中显示
   - 检查"一般陪诊服务"和"高级陪诊服务"等服务的图片

2. **服务列表页面测试**:
   - 点击"全部服务"按钮进入服务列表页面
   - 查看每个服务项的图片是否居中显示
   - 检查图片是否按比例缩放并填满容器

3. **不同设备测试**:
   - 在不同尺寸的设备上测试
   - 确保图片在各种屏幕尺寸下都能正确居中

## 预期效果

- 所有服务列表中的图片都应该在容器中居中显示
- 图片应该按比例缩放，不会变形
- 图片应该填满整个容器区域
- 在不同设备上都能保持一致的显示效果

## 注意事项

- 确保图片资源存在且可访问
- 如果图片不存在，会显示默认图片
- 图片的 `mode="aspectFill"` 属性确保图片按比例缩放并填满容器
