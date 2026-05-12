// quick_test_service.js
// 快速测试服务详情接口

const { request } = require('./utils/request')

// 快速测试ID为7的服务
async function quickTestService7() {
  console.log('=== 快速测试ID为7的服务 ===')
  
  const testCases = [
    { serviceId: 7, type: '数字' },
    { serviceId: "7", type: '字符串' },
    { serviceId: "007", type: '带前导零字符串' }
  ]
  
  for (const testCase of testCases) {
    console.log(`\n--- 测试${testCase.type}类型: ${testCase.serviceId} ---`)
    
    try {
      const startTime = Date.now()
      const res = await request('/api/service/detail', 'POST', { 
        serviceId: testCase.serviceId 
      })
      const endTime = Date.now()
      
      console.log('请求参数:', { serviceId: testCase.serviceId })
      console.log('响应时间:', `${endTime - startTime}ms`)
      console.log('响应结果:', res)
      
      if (res.code === 0 && res.data) {
        const service = res.data
        console.log('✅ 成功获取服务信息:')
        console.log(`   - ID: ${service.id}`)
        console.log(`   - 名称: ${service.name}`)
        console.log(`   - 分类: ${service.category}`)
        console.log(`   - 价格: ¥${service.price}`)
        console.log(`   - 原价: ¥${service.originalPrice}`)
        console.log(`   - 状态: ${service.status === 1 ? '上架' : '下架'}`)
        
        // 解析表单配置
        if (service.formConfig) {
          try {
            const formConfig = JSON.parse(service.formConfig)
            console.log(`   - 表单字段数: ${formConfig.fields ? formConfig.fields.length : 0}`)
          } catch (error) {
            console.error('解析表单配置失败:', error)
          }
        }
      } else {
        console.log('❌ 获取服务失败:', res.errorMsg)
      }
    } catch (error) {
      console.error('❌ 请求失败:', error.message)
    }
  }
}

// 测试错误情况
async function testErrorCases() {
  console.log('\n=== 测试错误情况 ===')
  
  const errorCases = [
    { serviceId: "abc", description: '非数字字符串' },
    { serviceId: "", description: '空字符串' },
    { serviceId: null, description: 'null值' },
    { serviceId: undefined, description: 'undefined值' }
  ]
  
  for (const testCase of errorCases) {
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

// 主测试函数
async function runQuickTest() {
  console.log('开始快速测试...')
  
  // 1. 测试ID为7的服务
  await quickTestService7()
  
  // 2. 测试错误情况
  await testErrorCases()
  
  console.log('\n=== 快速测试完成 ===')
}

// 运行测试
if (typeof module !== 'undefined' && module.exports) {
  module.exports = {
    quickTestService7,
    testErrorCases,
    runQuickTest
  }
} else {
  // 在浏览器环境中运行
  runQuickTest().then(() => {
    console.log('快速测试完成')
  }).catch(error => {
    console.error('测试过程中出现错误:', error)
  })
} 