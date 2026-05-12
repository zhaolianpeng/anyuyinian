// 测试当前配置
const { baseURL, env } = require('./config')

console.log('=== 当前配置测试 ===')
console.log('当前环境:', env)
console.log('baseURL:', baseURL)

// 测试完整的API路径
const testUrl = baseURL + '/api/wx/login'
console.log('完整登录URL:', testUrl)

// 检查是否是有效的URL
const urlPattern = /^https?:\/\/.+/i
if (urlPattern.test(testUrl)) {
  console.log('✅ URL格式正确')
} else {
  console.log('❌ URL格式错误')
}

// 检查域名是否可访问（模拟）
if (baseURL.includes('example.com')) {
  console.log('⚠️  警告：使用的是示例域名，可能无法访问')
} else if (baseURL.includes('golang-lfwy')) {
  console.log('✅ 使用的是生产环境域名')
} else {
  console.log('ℹ️  使用其他域名')
}

console.log('\n=== 测试完成 ===')
console.log('请确保后端服务正在运行，并且域名配置正确。') 