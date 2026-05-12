// 首页API测试脚本
const { api, baseURL } = require('./config')

console.log('=== 首页API测试 ===')

// 测试配置
console.log('当前环境baseURL:', baseURL)
console.log('首页初始化API路径:', api.homeInit)

// 测试完整的API路径
const testPaths = [
  api.homeInit,
  api.wxLogin,
  api.userInfo,
  api.serviceList,
  api.orderList,
  api.hospitalList
]

console.log('\n完整API路径测试:')
testPaths.forEach(path => {
  const fullUrl = baseURL + path
  console.log(`${path} -> ${fullUrl}`)
})

// 验证首页API路径
const homeInitUrl = baseURL + api.homeInit
console.log('\n首页初始化完整URL:', homeInitUrl)

// 检查是否包含重复的 /api
if (homeInitUrl.includes('/api/api/')) {
  console.log('❌ 错误：路径中包含重复的 /api')
} else {
  console.log('✅ 正确：路径格式正确')
}

// 模拟首页初始化请求参数
const mockParams = {
  longitude: 114.0579,
  latitude: 22.5431,
  limit: 10
}

console.log('\n模拟请求参数:', mockParams)
console.log('完整请求URL:', homeInitUrl + '?' + new URLSearchParams(mockParams).toString())

console.log('\n=== 测试完成 ===')
console.log('请重新编译小程序并测试首页功能。') 