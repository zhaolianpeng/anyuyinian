// 云开发URL转换测试脚本
const { cos } = require('./config')

// 转换云开发URL为直接访问URL
const convertCloudUrlToDirectUrl = (cloudUrl) => {
  // 云开发URL格式: @cloud://{环境ID}.{存储桶名称}/{文件路径}
  // 例如: @cloud://prod-5g94mx7a3d07e78c.7072-prod-5g94mx7a3d07e78c-1353115175/images/button/daiquyao_logo.png
  
  // 移除 @cloud:// 前缀
  const urlWithoutPrefix = cloudUrl.replace('@cloud://', '')
  
  // 分割环境ID和存储桶名称
  const parts = urlWithoutPrefix.split('/')
  if (parts.length < 2) {
    console.warn('Invalid cloud URL format:', cloudUrl)
    return cloudUrl
  }
  
  // 第一部分包含环境ID和存储桶名称
  const envAndBucket = parts[0]
  const envAndBucketParts = envAndBucket.split('.')
  
  if (envAndBucketParts.length !== 2) {
    console.warn('Invalid cloud URL format:', cloudUrl)
    return cloudUrl
  }
  
  const envId = envAndBucketParts[0]
  const bucketName = envAndBucketParts[1]
  
  // 构建文件路径
  const filePath = parts.slice(1).join('/')
  
  // 构建直接访问URL
  // 格式: https://{存储桶名称}-{APPID}.cos.{地域}.myqcloud.com/{文件路径}
  const directUrl = `https://${bucketName}-wx101090677bd5219e.cos.ap-shanghai.myqcloud.com/${filePath}`
  
  return directUrl
}

// 测试用例
const testUrls = [
  '@cloud://prod-5g94mx7a3d07e78c.7072-prod-5g94mx7a3d07e78c-1353115175/images/button/daiquyao_logo.png',
  '@cloud://prod-5g94mx7a3d07e78c.7072-prod-5g94mx7a3d07e78c-1353115175/images/banner/nursing.jpg',
  '@cloud://prod-5g94mx7a3d07e78c.7072-prod-5g94mx7a3d07e78c-1353115175/images/service/escort.jpg',
  '@cloud://prod-5g94mx7a3d07e78c.7072-prod-5g94mx7a3d07e78c-1353115175/images/hospital/logo.png'
]

console.log('=== 云开发URL转换测试 ===')
console.log('COS配置信息:')
console.log('- 云开发环境ID:', cos.cloudEnvId)
console.log('- 云开发存储桶名称:', cos.cloudBucketName)
console.log('- COS域名:', cos.bucketDomain)
console.log('')

console.log('=== URL转换测试 ===')
testUrls.forEach((cloudUrl, index) => {
  const directUrl = convertCloudUrlToDirectUrl(cloudUrl)
  console.log(`测试 ${index + 1}:`)
  console.log(`  云开发URL: ${cloudUrl}`)
  console.log(`  直接访问URL: ${directUrl}`)
  console.log('')
})

console.log('=== 使用说明 ===')
console.log('1. 在WXML中可以直接使用转换后的URL:')
console.log('   <image src="转换后的URL" mode="aspectFill" />')
console.log('')
console.log('2. 在JS中可以通过图片处理函数自动转换:')
console.log('   const imageUrl = getCosImageUrl(cloudUrl)')
console.log('')
console.log('3. 转换后的URL可以直接在浏览器中访问测试') 