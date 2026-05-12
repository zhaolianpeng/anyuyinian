/**
 * 服务端调用测试
 * 用于验证统一服务端 API 和基础接口是否正常工作
 */

const { api } = require('../utils/cloud-container-standard')

/**
 * 测试服务端连接
 */
async function testCloudConnection() {
  console.log('=== 开始测试服务端连接 ===')
  
  try {
    // 测试计数器API
    const countResult = await api.count.get()
    console.log('✅ 计数器API测试成功:', countResult)
    
    return true
  } catch (error) {
    console.error('❌ 服务端连接测试失败:', error)
    return false
  }
}

/**
 * 测试首页API
 */
async function testHomeAPI() {
  console.log('=== 开始测试首页API ===')
  
  try {
    const result = await api.homeInit({
      longitude: 121.4737,
      latitude: 31.2304,
      limit: 5
    })
    
    console.log('✅ 首页API测试成功:', result)
    return true
  } catch (error) {
    console.error('❌ 首页API测试失败:', error)
    return false
  }
}

/**
 * 测试用户API
 */
async function testUserAPI() {
  console.log('=== 开始测试用户API ===')
  
  try {
    const result = await api.userInfo()
    console.log('✅ 用户API测试成功:', result)
    return true
  } catch (error) {
    console.error('❌ 用户API测试失败:', error)
    return false
  }
}

/**
 * 测试服务API
 */
async function testServiceAPI() {
  console.log('=== 开始测试服务API ===')
  
  try {
    const result = await api.serviceList({ limit: 5 })
    console.log('✅ 服务API测试成功:', result)
    return true
  } catch (error) {
    console.error('❌ 服务API测试失败:', error)
    return false
  }
}

/**
 * 测试订单API
 */
async function testOrderAPI() {
  console.log('=== 开始测试订单API ===')
  
  try {
    const result = await api.orderList({ limit: 5 })
    console.log('✅ 订单API测试成功:', result)
    return true
  } catch (error) {
    console.error('❌ 订单API测试失败:', error)
    return false
  }
}

/**
 * 运行所有测试
 */
async function runAllTests() {
  console.log('🚀 开始运行服务端API测试...')
  
  const tests = [
    { name: '服务端连接', test: testCloudConnection },
    { name: '首页API', test: testHomeAPI },
    { name: '用户API', test: testUserAPI },
    { name: '服务API', test: testServiceAPI },
    { name: '订单API', test: testOrderAPI }
  ]
  
  const results = []
  
  for (const test of tests) {
    console.log(`\n📋 测试: ${test.name}`)
    const success = await test.test()
    results.push({ name: test.name, success })
  }
  
  // 输出测试结果
  console.log('\n📊 测试结果汇总:')
  results.forEach(result => {
    const status = result.success ? '✅ 通过' : '❌ 失败'
    console.log(`${status} ${result.name}`)
  })
  
  const passedCount = results.filter(r => r.success).length
  const totalCount = results.length
  
  console.log(`\n🎯 测试完成: ${passedCount}/${totalCount} 通过`)
  
  if (passedCount === totalCount) {
    console.log('🎉 所有测试通过！云托管连接正常')
  } else {
    console.log('⚠️  部分测试失败，请检查云托管配置')
  }
  
  return results
}

/**
 * 快速连接测试
 */
async function quickTest() {
  console.log('⚡ 快速连接测试...')
  
  try {
    const result = await api.count.get()
    console.log('✅ 连接正常:', result)
    return true
  } catch (error) {
    console.error('❌ 连接失败:', error)
    return false
  }
}

// 导出测试函数
module.exports = {
  testCloudConnection,
  testHomeAPI,
  testUserAPI,
  testServiceAPI,
  testOrderAPI,
  runAllTests,
  quickTest
}

// 如果直接运行此文件，执行快速测试
if (typeof wx !== 'undefined') {
  // 在小程序环境中，可以调用测试
  console.log('小程序环境检测到，可以运行测试')
} else {
  console.log('非小程序环境，测试函数已导出')
} 