// test_json_number_fix.js
// 测试JSON Number类型修复的脚本

const { request } = require('./utils/request')

// 测试不同格式的serviceId
async function testDifferentServiceIdFormats() {
  console.log('=== 测试不同格式的serviceId ===')
  
  const testCases = [
    { serviceId: 7, description: '数字类型 7' },
    { serviceId: "7", description: '字符串类型 "7"' },
    { serviceId: "007", description: '字符串类型 "007"' },
    { serviceId: 1, description: '数字类型 1' },
    { serviceId: "1", description: '字符串类型 "1"' },
    { serviceId: 999, description: '不存在的数字ID' },
    { serviceId: "999", description: '不存在的字符串ID' }
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
        console.log(`   - 价格: ¥${service.price}`)
      } else {
        console.log('❌ 获取服务失败:', res.errorMsg)
      }
    } catch (error) {
      console.error('❌ 请求失败:', error.message)
    }
  }
}

// 测试边界情况
async function testEdgeCases() {
  console.log('\n=== 测试边界情况 ===')
  
  const edgeCases = [
    { serviceId: 0, description: 'ID为0' },
    { serviceId: "0", description: '字符串ID为0' },
    { serviceId: -1, description: '负数ID' },
    { serviceId: "-1", description: '字符串负数ID' },
    { serviceId: "abc", description: '非数字字符串' },
    { serviceId: "", description: '空字符串' },
    { serviceId: null, description: 'null值' },
    { serviceId: undefined, description: 'undefined值' }
  ]
  
  for (const testCase of edgeCases) {
    console.log(`\n--- ${testCase.description} ---`)
    
    try {
      const res = await request('/api/service/detail', 'POST', { 
        serviceId: testCase.serviceId 
      })
      
      console.log('请求参数:', { serviceId: testCase.serviceId })
      console.log('响应结果:', res)
    } catch (error) {
      console.error('❌ 请求失败:', error.message)
    }
  }
}

// 测试批量请求
async function testBatchRequests() {
  console.log('\n=== 测试批量请求 ===')
  
  const serviceIds = [
    1, "1", 2, "2", 3, "3", 4, "4", 
    5, "5", 6, "6", 7, "7", 8, "8"
  ]
  
  const results = []
  
  for (const serviceId of serviceIds) {
    try {
      const res = await request('/api/service/detail', 'POST', { serviceId })
      
      if (res.code === 0 && res.data) {
        results.push({
          serviceId,
          success: true,
          name: res.data.name,
          price: res.data.price
        })
      } else {
        results.push({
          serviceId,
          success: false,
          error: res.errorMsg
        })
      }
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
      ? `${result.name} (¥${result.price})`
      : result.error
    console.log(`${status} ID ${result.serviceId}: ${info}`)
  })
}

// 测试前端实际使用场景
async function testFrontendUsage() {
  console.log('\n=== 测试前端实际使用场景 ===')
  
  // 模拟前端页面中的使用
  const page = {
    data: {
      serviceId: 7
    },
    
    async loadServiceDetail() {
      console.log('模拟页面加载服务详情...')
      
      try {
        // 测试数字类型
        console.log('1. 测试数字类型serviceId')
        const res1 = await request('/api/service/detail', 'POST', { 
          serviceId: this.data.serviceId 
        })
        console.log('数字类型结果:', res1)
        
        // 测试字符串类型
        console.log('2. 测试字符串类型serviceId')
        const res2 = await request('/api/service/detail', 'POST', { 
          serviceId: String(this.data.serviceId) 
        })
        console.log('字符串类型结果:', res2)
        
        // 测试从URL参数获取的ID
        console.log('3. 测试URL参数ID')
        const urlParamId = "7" // 模拟从URL获取的字符串ID
        const res3 = await request('/api/service/detail', 'POST', { 
          serviceId: urlParamId 
        })
        console.log('URL参数结果:', res3)
        
      } catch (error) {
        console.error('页面加载失败:', error.message)
      }
    }
  }
  
  await page.loadServiceDetail()
}

// 主测试函数
async function runJsonNumberTests() {
  console.log('开始测试JSON Number类型修复...')
  
  // 1. 测试不同格式的serviceId
  await testDifferentServiceIdFormats()
  
  // 2. 测试边界情况
  await testEdgeCases()
  
  // 3. 测试批量请求
  await testBatchRequests()
  
  // 4. 测试前端实际使用场景
  await testFrontendUsage()
  
  console.log('\n=== JSON Number类型修复测试完成 ===')
}

// 运行测试
if (typeof module !== 'undefined' && module.exports) {
  module.exports = {
    testDifferentServiceIdFormats,
    testEdgeCases,
    testBatchRequests,
    testFrontendUsage,
    runJsonNumberTests
  }
} else {
  // 在浏览器环境中运行
  runJsonNumberTests().then(() => {
    console.log('所有JSON Number测试完成')
  }).catch(error => {
    console.error('测试过程中出现错误:', error)
  })
} 