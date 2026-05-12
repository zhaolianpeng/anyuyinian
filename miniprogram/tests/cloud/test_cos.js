// COS配置测试脚本
const { cos } = require('./config')

// 测试图片路径
const testImages = [
  '/images/banner-nursing.jpg',
  '/images/banner-escort.jpg',
  '/images/icon-service.png',
  '/images/icon-order.png',
  '/images/service-nursing.jpg',
  '/images/service-escort.jpg',
  '/images/hospital-1.jpg',
  '/images/hospital-2.jpg'
]

console.log('=== COS配置测试 ===')
console.log('存储桶域名:', cos.bucketDomain)
console.log('图片前缀:', cos.imagePrefix)
console.log('')

console.log('=== 图片URL转换测试 ===')
testImages.forEach(imagePath => {
  const fullUrl = cos.bucketDomain + imagePath
  console.log(`${imagePath} -> ${fullUrl}`)
})

console.log('')
console.log('=== 测试访问 ===')
console.log('请在浏览器中访问以下URL测试图片加载:')
console.log(`${cos.bucketDomain}/images/test.jpg`)
console.log('')
console.log('=== 小程序测试 ===')
console.log('1. 在微信开发者工具中打开小程序')
console.log('2. 检查首页图片是否正常加载')
console.log('3. 检查控制台是否有图片加载错误') 