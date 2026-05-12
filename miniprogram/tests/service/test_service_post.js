// test_service_post.js
// 测试服务详情POST请求的脚本

const { request } = require('./utils/request')

// 测试服务详情POST请求
async function testServiceDetailPost() {
  console.log('=== 测试服务详情POST请求 ===')
  
  const testCases = [
    { serviceId: 7, description: '测试ID为7的服务（数字）' },
    { serviceId: "7", description: '测试ID为7的服务（字符串）' },
    { serviceId: 1, description: '测试ID为1的服务（数字）' },
    { serviceId: "1", description: '测试ID为1的服务（字符串）' },
    { serviceId: 999, description: '测试不存在的服务ID（数字）' },
    { serviceId: "999", description: '测试不存在的服务ID（字符串）' },
    { serviceId: 0, description: '测试无效的服务ID 0' },
    { serviceId: "0", description: '测试无效的服务ID "0"' },
    { serviceId: -1, description: '测试无效的服务ID -1' },
    { serviceId: "-1", description: '测试无效的服务ID "-1"' }
  ]
  
  for (const testCase of testCases) {
    console.log(`\n--- ${testCase.description} ---`)
    
    try {
      const res = await request('/api/service/detail', 'POST', { 
        serviceId: testCase.serviceId 
      })
      
      console.log('请求参数:', { serviceId: testCase.serviceId })
      console.log('响应结果:', res)
      
      if (res.code === 0 && res.data) {
        const service = res.data
        console.log('✅ 成功获取服务信息:')
        console.log(`   - ID: ${service.id}`)
        console.log(`   - 名称: ${service.name}`)
        console.log(`   - 分类: ${service.category}`)
        console.log(`   - 价格: ¥${service.price}`)
        console.log(`   - 状态: ${service.status === 1 ? '上架' : '下架'}`)
      } else {
        console.log('❌ 获取服务失败:', res.errorMsg)
      }
    } catch (error) {
      console.error('❌ 请求失败:', error.message)
    }
  }
}

// 测试批量请求
async function testBatchRequests() {
  console.log('\n=== 测试批量请求 ===')
  
  const serviceIds = [1, 2, 3, 4, 5, 6, 7, 8]
  const results = []
  
  for (const serviceId of serviceIds) {
    try {
      const res = await request('/api/service/detail', 'POST', { serviceId })
      results.push({
        serviceId,
        success: res.code === 0,
        data: res.data,
        error: res.errorMsg
      })
    } catch (error) {
      results.push({
        serviceId,
        success: false,
        error: error.message
      })
    }
  }
  
  console.log('批量请求结果:')
  results.forEach(result => {
    const status = result.success ? '✅' : '❌'
    const info = result.success 
      ? `${result.data.name} (¥${result.data.price})`
      : result.error
    console.log(`${status} ID ${result.serviceId}: ${info}`)
  })
}

// 测试请求参数验证
async function testParameterValidation() {
  console.log('\n=== 测试参数验证 ===')
  
  const invalidRequests = [
    { data: {}, description: '空请求体' },
    { data: { serviceId: 'abc' }, description: '字符串ID' },
    { data: { serviceId: null }, description: 'null ID' },
    { data: { serviceId: undefined }, description: 'undefined ID' },
    { data: { otherField: 7 }, description: '缺少serviceId字段' }
  ]
  
  for (const testCase of invalidRequests) {
    console.log(`\n--- ${testCase.description} ---`)
    
    try {
      const res = await request('/api/service/detail', 'POST', testCase.data)
      console.log('请求参数:', testCase.data)
      console.log('响应结果:', res)
    } catch (error) {
      console.error('请求失败:', error.message)
    }
  }
}

// 主测试函数
async function runPostTests() {
  console.log('开始测试服务详情POST请求...')
  
  // 1. 测试基本功能
  await testServiceDetailPost()
  
  // 2. 测试批量请求
  await testBatchRequests()
  
  // 3. 测试参数验证
  await testParameterValidation()
  
  console.log('\n=== POST请求测试完成 ===')
}

// 运行测试
if (typeof module !== 'undefined' && module.exports) {
  module.exports = {
    testServiceDetailPost,
    testBatchRequests,
    testParameterValidation,
    runPostTests
  }
} else {
  // 在浏览器环境中运行
  runPostTests().then(() => {
    console.log('所有POST测试完成')
  }).catch(error => {
    console.error('测试过程中出现错误:', error)
  })
} 