// test_service_data.js
// 测试服务数据的脚本

const { request } = require('./utils/request')

// 测试服务列表
async function testServiceList() {
  console.log('=== 测试服务列表 ===')
  
  try {
    const res = await request('/api/service/list', 'GET')
    console.log('服务列表响应:', res)
    
    if (res.code === 0 && res.data && res.data.list) {
      console.log(`找到 ${res.data.list.length} 个服务`)
      res.data.list.forEach((service, index) => {
        console.log(`${index + 1}. ID: ${service.id}, 名称: ${service.name}, 价格: ¥${service.price}`)
      })
      return res.data.list
    } else {
      console.error('获取服务列表失败:', res.errorMsg)
      return []
    }
  } catch (error) {
    console.error('测试服务列表失败:', error)
    return []
  }
}

// 测试服务详情
async function testServiceDetail(serviceId) {
  console.log(`\n=== 测试服务详情 ID: ${serviceId} ===`)
  
  try {
    const res = await request('/api/service/detail', 'POST', { serviceId })
    console.log('服务详情响应:', res)
    
    if (res.code === 0 && res.data) {
      const service = res.data
      console.log('服务信息:')
      console.log(`- ID: ${service.id}`)
      console.log(`- 名称: ${service.name}`)
      console.log(`- 分类: ${service.category}`)
      console.log(`- 价格: ¥${service.price}`)
      console.log(`- 原价: ¥${service.originalPrice}`)
      console.log(`- 状态: ${service.status === 1 ? '上架' : '下架'}`)
      
      // 解析表单配置
      if (service.formConfig) {
        try {
          const formConfig = JSON.parse(service.formConfig)
          console.log(`- 表单字段数: ${formConfig.fields ? formConfig.fields.length : 0}`)
        } catch (error) {
          console.error('解析表单配置失败:', error)
        }
      }
      
      return service
    } else {
      console.error('获取服务详情失败:', res.errorMsg)
      return null
    }
  } catch (error) {
    console.error('测试服务详情失败:', error)
    return null
  }
}

// 测试所有服务详情
async function testAllServiceDetails() {
  console.log('\n=== 测试所有服务详情 ===')
  
  const services = await testServiceList()
  
  for (const service of services) {
    await testServiceDetail(service.id)
  }
}

// 测试无效服务ID
async function testInvalidServiceId() {
  console.log('\n=== 测试无效服务ID ===')
  
  const invalidIds = [999, 0, -1]
  
  for (const id of invalidIds) {
    console.log(`测试无效ID: ${id}`)
    try {
      const res = await request('/api/service/detail', 'POST', { serviceId: id })
      console.log('响应:', res)
    } catch (error) {
      console.error(`ID ${id} 测试失败:`, error.message)
    }
  }
}

// 主测试函数
async function runServiceTests() {
  console.log('开始测试服务数据...')
  
  // 1. 测试服务列表
  await testServiceList()
  
  // 2. 测试所有服务详情
  await testAllServiceDetails()
  
  // 3. 测试无效服务ID
  await testInvalidServiceId()
  
  console.log('\n=== 服务数据测试完成 ===')
}

// 运行测试
if (typeof module !== 'undefined' && module.exports) {
  module.exports = {
    testServiceList,
    testServiceDetail,
    testAllServiceDetails,
    testInvalidServiceId,
    runServiceTests
  }
} else {
  // 在浏览器环境中运行
  runServiceTests().then(() => {
    console.log('所有测试完成')
  }).catch(error => {
    console.error('测试过程中出现错误:', error)
  })
} 