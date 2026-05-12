// 快速图片修复脚本

console.log('=== 快速图片修复 ===')

// 问题诊断
function diagnoseImageIssues() {
  console.log('\n1. 检查配置...')
  
  try {
    const { env, dev, cos } = require('./config')
    console.log(`   当前环境: ${env}`)
    console.log(`   使用模拟数据: ${dev.useMockData}`)
    console.log(`   COS域名: ${cos.bucketDomain}`)
    console.log(`   图片前缀: ${cos.imagePrefix}`)
  } catch (error) {
    console.error('   配置加载失败:', error.message)
  }
}

// 测试图片URL
function testImageUrls() {
  console.log('\n2. 测试图片URL...')
  
  const testUrls = [
    'https://7072-prod-5g94mx7a3d07e78c-1353115175.cos.ap-shanghai.myqcloud.com/static/fuwu_2.jpeg',
    'https://7072-prod-5g94mx7a3d07e78c-1353115175.cos.ap-shanghai.myqcloud.com/static/fuwu_1.jpeg',
    'https://7072-prod-5g94mx7a3d07e78c-1353115175.cos.ap-shanghai.myqcloud.com/static/fuwuyuyue_logo.png'
  ]
  
  testUrls.forEach((url, index) => {
    console.log(`   图片 ${index + 1}: ${url}`)
    
    // 检查URL格式
    if (url.startsWith('https://')) {
      console.log(`     ✅ HTTPS格式正确`)
    } else {
      console.log(`     ❌ 不是HTTPS格式`)
    }
    
    // 检查域名
    if (url.includes('cos.ap-shanghai.myqcloud.com')) {
      console.log(`     ✅ COS域名正确`)
    } else {
      console.log(`     ❌ COS域名错误`)
    }
    
    // 检查路径
    if (url.includes('/static/')) {
      console.log(`     ✅ 路径格式正确`)
    } else {
      console.log(`     ❌ 路径格式错误`)
    }
  })
}

// 提供解决方案
function provideSolutions() {
  console.log('\n3. 解决方案...')
  
  console.log('\n方案1: 检查小程序后台域名配置')
  console.log('   1. 登录微信公众平台')
  console.log('   2. 进入小程序管理后台')
  console.log('   3. 点击"开发" -> "开发管理" -> "开发设置"')
  console.log('   4. 在"服务器域名"中添加:')
  console.log('      downloadFile合法域名:')
  console.log('      https://7072-prod-5g94mx7a3d07e78c-1353115175.cos.ap-shanghai.myqcloud.com')
  
  console.log('\n方案2: 开发环境临时解决')
  console.log('   1. 打开微信开发者工具')
  console.log('   2. 点击右上角"详情"')
  console.log('   3. 勾选"不校验合法域名、web-view（业务域名）、TLS 版本以及 HTTPS 证书"')
  
  console.log('\n方案3: 检查COS桶权限')
  console.log('   1. 登录腾讯云控制台')
  console.log('   2. 进入COS存储桶管理')
  console.log('   3. 确认存储桶权限为"公有读私有写"')
  console.log('   4. 检查文件是否存在')
  
  console.log('\n方案4: 在浏览器中测试图片')
  console.log('   在浏览器中访问以下URL测试:')
  console.log('   https://7072-prod-5g94mx7a3d07e78c-1353115175.cos.ap-shanghai.myqcloud.com/static/fuwu_2.jpeg')
}

// 生成修复代码
function generateFixCode() {
  console.log('\n4. 修复代码...')
  
  console.log('\n在首页模板中添加调试信息:')
  console.log(`
<!-- 调试信息 -->
<view wx:if="{{banners.length > 0}}" style="padding: 10rpx; background: #f0f0f0;">
  <text>轮播图数量: {{banners.length}}</text>
  <view wx:for="{{banners}}" wx:key="id">
    <text>图片 {{index + 1}}: {{item.imageUrl}}</text>
  </view>
</view>
  `)
  
  console.log('\n在首页JS中添加错误处理:')
  console.log(`
// 图片加载错误处理
onImageError(e) {
  console.error('图片加载失败:', e.detail)
  const { index, type } = e.currentTarget.dataset
  console.log('失败的图片:', index, type)
}
  `)
  
  console.log('\n在模板中添加错误处理:')
  console.log(`
<image 
  src="{{item.imageUrl}}" 
  mode="aspectFill" 
  binderror="onImageError"
  data-index="{{index}}"
  data-type="banner"
/>
  `)
}

// 检查当前状态
function checkCurrentStatus() {
  console.log('\n5. 当前状态检查...')
  
  try {
    const { processImageUrl } = require('./utils/image')
    const testUrl = '/static/fuwu_2.jpeg'
    const processed = processImageUrl(testUrl)
    
    console.log(`   测试URL: ${testUrl}`)
    console.log(`   处理后: ${processed}`)
    console.log(`   处理成功: ${processed.startsWith('https://')}`)
    
  } catch (error) {
    console.error('   图片处理工具加载失败:', error.message)
  }
}

// 执行修复
function runQuickFix() {
  console.log('开始快速图片修复...\n')
  
  diagnoseImageIssues()
  testImageUrls()
  provideSolutions()
  generateFixCode()
  checkCurrentStatus()
  
  console.log('\n=== 修复完成 ===')
  console.log('')
  console.log('下一步操作:')
  console.log('1. 检查小程序后台域名配置')
  console.log('2. 在浏览器中测试图片URL')
  console.log('3. 确认COS桶权限设置')
  console.log('4. 在开发环境中勾选"不校验合法域名"')
  console.log('')
  console.log('如果问题仍然存在，请提供以下信息:')
  console.log('- 小程序后台域名配置截图')
  console.log('- 浏览器中图片URL访问结果')
  console.log('- 开发者工具控制台错误信息')
}

// 运行修复
runQuickFix() 