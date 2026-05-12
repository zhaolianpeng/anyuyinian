# 视频封面图片设置指南

## 问题说明
小程序视频组件的 `poster` 属性只支持网络地址，不支持本地路径。因此需要将封面图片上传到腾讯云对象存储。

## 解决方案

### 1. 上传封面图片到腾讯云对象存储

#### 方法一：通过微信开发者工具
1. 打开微信开发者工具
2. 进入云开发控制台
3. 选择"存储" -> "文件管理"
4. 上传封面图片到路径：`static/video-poster.jpg`
5. 确保图片尺寸为 750x400px，格式为 JPG 或 PNG

#### 方法二：通过代码上传
```javascript
// 在开发者工具控制台执行
wx.cloud.uploadFile({
  cloudPath: 'static/video-poster.jpg',
  filePath: '/path/to/your/poster.jpg', // 本地图片路径
  success: res => {
    console.log('封面图片上传成功:', res.fileID)
  },
  fail: err => {
    console.error('封面图片上传失败:', err)
  }
})
```

### 2. 更新配置

封面图片上传成功后，系统会自动使用以下COS ID：
```
cloud://prod-5g94mx7a3d07e78c.7072-prod-5g94mx7a3d07e78c-1353115175/static/video-poster.jpg
```

### 3. 封面图片要求

- **尺寸**：750x400px（推荐）
- **格式**：JPG 或 PNG
- **大小**：建议小于 500KB
- **内容**：护理服务相关图片，与视频内容匹配

### 4. 临时解决方案

如果暂时没有封面图片，可以：
1. 移除 `poster` 属性，让视频显示第一帧
2. 或者使用现有的服务图片作为封面

### 5. 测试验证

上传完成后，重新加载小程序，检查：
1. 视频是否正常加载
2. 封面图片是否显示
3. 视频播放是否正常

## 相关文件
- `pages/index/index.js` - 视频配置
- `pages/index/index.wxml` - 视频组件
- `pages/index/index.wxss` - 视频样式
