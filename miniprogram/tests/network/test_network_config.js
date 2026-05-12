// 网络配置测试脚本
const { baseURL, cos } = require('../config')

console.log('=== 网络配置测试 ===')

// 测试API域名
console.log('API域名:', baseURL)

// 测试COS域名
console.log('COS域名:', cos.bucketDomain)

// 测试域名格式
const testDomains = [
  baseURL,
  cos.bucketDomain
]

testDomains.forEach((domain, index) => {
  console.log(`域名 ${index + 1}:`, domain)
  
  // 检查是否为HTTPS
  if (!domain.startsWith('https://')) {
    console.error('❌ 域名必须使用HTTPS协议:', domain)
  } else {
    console.log('✅ HTTPS协议正确')
  }
  
  // 检查是否包含端口号
  if (domain.includes(':') && !domain.includes('://')) {
    console.error('❌ 域名不能包含端口号:', domain)
  } else {
    console.log('✅ 域名格式正确')
  }
  
  // 检查是否为IP地址
  const domainWithoutProtocol = domain.replace('https://', '')
  const ipRegex = /^\d+\.\d+\.\d+\.\d+$/
  if (ipRegex.test(domainWithoutProtocol)) {
    console.error('❌ 域名不能使用IP地址:', domain)
  } else {
    console.log('✅ 域名格式正确')
  }
})

// 需要在微信小程序后台配置的域名列表
console.log('\n=== 需要在微信小程序后台配置的域名 ===')
console.log('1. request合法域名:')
console.log('   ', baseURL)
console.log('\n2. uploadFile合法域名:')
console.log('   ', baseURL)
console.log('\n3. downloadFile合法域名:')
console.log('   ', cos.bucketDomain)

// 配置检查清单
console.log('\n=== 配置检查清单 ===')
console.log('□ 登录微信公众平台')
console.log('□ 进入"开发" -> "开发管理" -> "开发设置"')
console.log('□ 在"服务器域名"部分添加上述域名')
console.log('□ 确保域名使用HTTPS协议')
console.log('□ 保存配置')
console.log('□ 等待配置生效（通常几分钟）')
console.log('□ 重新上传体验版测试')

// 常见问题排查
console.log('\n=== 常见问题排查 ===')
console.log('1. 域名未配置: 检查微信小程序后台域名配置')
console.log('2. HTTPS证书问题: 检查服务器SSL证书')
console.log('3. 网络超时: 检查服务器响应时间')
console.log('4. 防火墙拦截: 检查服务器防火墙设置')
console.log('5. DNS解析问题: 检查域名解析是否正确')

module.exports = {
  testDomains,
  baseURL,
  cosDomain: cos.bucketDomain
} 