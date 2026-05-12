// 创建占位图片的脚本
// 这个脚本会创建简单的占位图片来解决图片加载错误

const fs = require('fs')
const path = require('path')

// 需要创建的图片列表
const imagesToCreate = [
  // 导航图标
  'images/nav/appointment.png',
  'images/nav/consultation.png',
  'images/nav/health-record.png',
  'images/nav/report.png',
  'images/nav/medicine.png',
  'images/nav/news.png',
  
  // 服务图片
  'images/service/appointment.png',
  'images/service/consultation.png',
  'images/service/checkup.png',
  'images/service/report.png',
  'images/service/medicine.png',
  'images/service/record.png',
  
  // 医院图片
  'images/hospital/rmyy-logo.png',
  'images/hospital/dermyy-logo.png',
  'images/hospital/zyy-logo.png',
  'images/hospital/etyy-logo.png',
  'images/hospital/fybjy-logo.png',
  
  // 其他需要的图片
  'images/service-default.jpg',
  'images/empty-state.png',
  'images/default-avatar.png'
]

console.log('=== 创建占位图片 ===')

// 创建目录
const dirs = [
  'images/nav',
  'images/service', 
  'images/hospital'
]

dirs.forEach(dir => {
  if (!fs.existsSync(dir)) {
    fs.mkdirSync(dir, { recursive: true })
    console.log(`✅ 创建目录: ${dir}`)
  }
})

// 创建占位图片文件
imagesToCreate.forEach(imagePath => {
  const fullPath = path.join(__dirname, imagePath)
  const dir = path.dirname(fullPath)
  
  // 确保目录存在
  if (!fs.existsSync(dir)) {
    fs.mkdirSync(dir, { recursive: true })
  }
  
  // 创建占位图片文件（1x1像素的透明PNG）
  const placeholderData = Buffer.from([
    0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG header
    0x00, 0x00, 0x00, 0x0D, // IHDR chunk length
    0x49, 0x48, 0x44, 0x52, // IHDR
    0x00, 0x00, 0x00, 0x01, // width: 1
    0x00, 0x00, 0x00, 0x01, // height: 1
    0x08, 0x06, 0x00, 0x00, 0x00, // bit depth, color type, compression, filter, interlace
    0x1F, 0x15, 0xC4, 0x89, // CRC
    0x00, 0x00, 0x00, 0x0C, // IDAT chunk length
    0x49, 0x44, 0x41, 0x54, // IDAT
    0x08, 0x99, 0x01, 0x01, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x02, 0x00, 0x01, // IDAT data
    0x00, 0x00, 0x00, 0x00, // IEND chunk length
    0x49, 0x45, 0x4E, 0x44, // IEND
    0xAE, 0x42, 0x60, 0x82  // CRC
  ])
  
  fs.writeFileSync(fullPath, placeholderData)
  console.log(`✅ 创建占位图片: ${imagePath}`)
})

console.log('')
console.log('=== 完成 ===')
console.log('所有占位图片已创建完成！')
console.log('现在可以重新编译小程序，图片加载错误应该会消失。')
console.log('')
console.log('注意：这些是1x1像素的透明图片，仅用于解决加载错误。')
console.log('在实际使用中，请替换为真实的图片文件。') 