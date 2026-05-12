// 批量修复API路径的脚本
const fs = require('fs')
const path = require('path')

// 需要修改的路径映射
const pathMappings = [
  { from: "/user/info", to: "/api/user/info" },
  { from: "/user/bind_phone", to: "/api/user/bind_phone" },
  { from: "/user/address", to: "/api/user/address" },
  { from: "/user/patient", to: "/api/user/patient" },
  { from: "/service/list", to: "/api/service/list" },
  { from: "/service/detail/", to: "/api/service/detail/" },
  { from: "/service/form_config/", to: "/api/service/form_config/" },
  { from: "/order/submit", to: "/api/order/submit" },
  { from: "/order/pay/", to: "/api/order/pay/" },
  { from: "/order/cancel/", to: "/api/order/cancel/" },
  { from: "/order/refund/", to: "/api/order/refund/" },
  { from: "/order/list", to: "/api/order/list" },
  { from: "/order/detail/", to: "/api/order/detail/" },
  { from: "/referral/qrcode", to: "/api/referral/qrcode" },
  { from: "/referral/list", to: "/api/referral/list" },
  { from: "/commission/list", to: "/api/commission/list" },
  { from: "/cashout/submit", to: "/api/cashout/submit" },
  { from: "/kefu/message", to: "/api/kefu/message" },
  { from: "/kefu/faq", to: "/api/kefu/faq" },
  { from: "/hospital/list", to: "/api/hospital/list" },
  { from: "/hospital/detail/", to: "/api/hospital/detail/" },
  { from: "/upload/list", to: "/api/upload/list" },
  { from: "/config", to: "/api/config" }
]

// 读取文件内容
const filePath = path.join(__dirname, 'utils/request.js')
let content = fs.readFileSync(filePath, 'utf8')

console.log('开始修复API路径...')

// 批量替换路径
pathMappings.forEach(mapping => {
  const oldPattern = `return request('${mapping.from}`
  const newPattern = `return request('${mapping.to}`
  
  if (content.includes(oldPattern)) {
    content = content.replace(new RegExp(oldPattern.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'g'), newPattern)
    console.log(`✅ 已修复: ${mapping.from} -> ${mapping.to}`)
  } else {
    console.log(`⚠️  未找到: ${mapping.from}`)
  }
})

// 写回文件
fs.writeFileSync(filePath, content, 'utf8')

console.log('\nAPI路径修复完成！')
console.log('请重新编译小程序并测试登录功能。') 