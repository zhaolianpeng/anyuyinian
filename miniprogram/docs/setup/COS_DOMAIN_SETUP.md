# COS存储桶域名配置指南

## 获取COS存储桶域名

### 方法一：从腾讯云控制台获取

1. **登录腾讯云控制台**
   - 访问 https://console.cloud.tencent.com/
   - 登录您的腾讯云账号

2. **进入对象存储COS**
   - 在控制台首页搜索"对象存储"或"COS"
   - 点击进入对象存储COS服务

3. **查看存储桶列表**
   - 在左侧菜单选择"存储桶列表"
   - 找到您要使用的存储桶

4. **获取域名信息**
   - 点击存储桶名称进入详情页
   - 在"基础配置"中找到"访问域名"
   - 复制完整的域名（格式如下）

### 方法二：手动构建域名

COS存储桶域名格式：
```
https://{存储桶名称}-{APPID}.cos.{地域}.myqcloud.com
```

#### 参数说明：
- **存储桶名称**: 您创建的存储桶名称
- **APPID**: 您的腾讯云账号APPID（12位数字）
- **地域**: 存储桶所在地域代码

#### 地域代码对照表：
| 地域名称 | 地域代码 |
|---------|---------|
| 北京 | ap-beijing |
| 上海 | ap-shanghai |
| 广州 | ap-guangzhou |
| 深圳 | ap-shenzhen |
| 成都 | ap-chengdu |
| 重庆 | ap-chongqing |
| 天津 | ap-tianjin |
| 南京 | ap-nanjing |
| 武汉 | ap-wuhan |
| 西安 | ap-xian |
| 杭州 | ap-hangzhou |
| 福州 | ap-fuzhou |
| 厦门 | ap-xiamen |
| 长沙 | ap-changsha |
| 济南 | ap-jinan |
| 青岛 | ap-qingdao |
| 大连 | ap-dalian |
| 石家庄 | ap-shijiazhuang |
| 郑州 | ap-zhengzhou |
| 太原 | ap-taiyuan |
| 沈阳 | ap-shenyang |
| 哈尔滨 | ap-haerbin |
| 长春 | ap-changchun |
| 合肥 | ap-hefei |
| 南昌 | ap-nanchang |
| 昆明 | ap-kunming |
| 贵阳 | ap-guiyang |
| 南宁 | ap-nanning |
| 海口 | ap-haikou |
| 兰州 | ap-lanzhou |
| 西宁 | ap-xining |
| 银川 | ap-yinchuan |
| 乌鲁木齐 | ap-urumqi |
| 拉萨 | ap-lhasa |
| 香港 | ap-hongkong |
| 新加坡 | ap-singapore |
| 东京 | ap-tokyo |
| 首尔 | ap-seoul |
| 孟买 | ap-mumbai |
| 曼谷 | ap-bangkok |
| 法兰克福 | eu-frankfurt |
| 莫斯科 | eu-moscow |
| 硅谷 | na-siliconvalley |
| 弗吉尼亚 | na-ashburn |
| 多伦多 | na-toronto |
| 圣保罗 | sa-saopaulo |

## 配置步骤

### 1. 获取必要信息
您需要准备以下信息：
- **存储桶名称**: 例如 `my-bucket`
- **APPID**: 您的腾讯云账号APPID（12位数字）
- **地域**: 存储桶所在地域，例如 `ap-beijing`

### 2. 构建域名
使用以下格式构建域名：
```
https://{存储桶名称}-{APPID}.cos.{地域}.myqcloud.com
```

**示例**：
- 存储桶名称: `my-bucket`
- APPID: `1234567890`
- 地域: `ap-beijing`
- 完整域名: `https://my-bucket-1234567890.cos.ap-beijing.myqcloud.com`

### 3. 更新配置文件
在 `config.js` 中更新 `bucketDomain`：

```javascript
const COS_CONFIG = {
  // 替换为您的实际域名
  bucketDomain: 'https://my-bucket-1234567890.cos.ap-beijing.myqcloud.com',
  imagePrefix: '/images/',
  // ... 其他配置
}
```

## 验证配置

### 1. 测试域名访问
在浏览器中访问您的COS域名，应该能看到存储桶的XML列表页面。

### 2. 测试图片访问
上传一张测试图片到存储桶，然后通过以下URL访问：
```
https://{您的域名}/images/test.jpg
```

### 3. 小程序测试
在微信开发者工具中测试图片加载是否正常。

## 常见问题

### Q: 如何找到我的APPID？
A: 在腾讯云控制台右上角点击头像，选择"账号信息"，在"基本信息"中可以看到APPID。

### Q: 存储桶名称包含特殊字符怎么办？
A: 存储桶名称中的特殊字符在域名中会被URL编码，但建议使用简单的英文名称。

### Q: 域名访问失败怎么办？
A: 
1. 检查存储桶权限设置是否为"公有读私有写"
2. 确认域名格式是否正确
3. 检查网络连接是否正常

### Q: 图片加载失败怎么办？
A:
1. 确认图片文件已上传到正确路径
2. 检查图片文件格式是否支持
3. 验证域名配置是否正确

## 安全建议

### 1. 权限设置
- 建议设置为"公有读私有写"
- 避免设置为"公有读写"

### 2. 防盗链配置
- 可以配置防盗链来限制访问来源
- 添加微信小程序的域名到白名单

### 3. 访问控制
- 定期检查访问日志
- 监控异常访问行为

## 相关文件
- `config.js` - 主配置文件
- `COS_SETUP.md` - COS完整配置说明
- `COS_DOMAIN_SETUP.md` - 本域名配置指南 