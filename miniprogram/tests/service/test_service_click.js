// 测试服务点击功能
const { api } = require('../utils/cloud-container-standard')

async function testServiceClickFeature() {
  console.log('=== 测试服务点击功能 ===')
  
  try {
    // 1. 测试首页init接口，获取服务数据
    console.log('1. 获取首页服务数据...')
    const result = await api.homeInit({
      longitude: 121.4737,
      latitude: 31.2304,
      limit: 10
    })
    
    if (result.code === 0 && result.data && result.data.services) {
      const services = result.data.services
      console.log('获取到服务数量:', services.length)
      
      // 2. 检查服务数据结构
      if (services.length > 0) {
        const firstService = services[0]
        console.log('第一个服务数据结构:', {
          id: firstService.id,
          serviceitemid: firstService.serviceitemid,
          name: firstService.name,
          description: firstService.description
        })
        
        // 3. 验证serviceitemid字段是否存在
        if (firstService.serviceitemid) {
          console.log('✅ serviceitemid字段存在:', firstService.serviceitemid)
        } else {
          console.log('❌ serviceitemid字段不存在')
        }
        
        // 4. 验证serviceitemid是否存在且有效
        if (firstService.serviceitemid && firstService.serviceitemid > 0) {
          console.log('✅ serviceitemid存在且有效:', firstService.serviceitemid)
        } else {
          console.log('❌ serviceitemid不存在或无效')
        }
        
        // 5. 测试服务详情接口
        console.log('2. 测试服务详情接口...')
        const serviceDetailResult = await api.serviceDetail({ 
          serviceId: firstService.serviceitemid 
        })
        
        if (serviceDetailResult.code === 0) {
          console.log('✅ 服务详情接口调用成功')
          console.log('服务详情:', {
            id: serviceDetailResult.data.id,
            name: serviceDetailResult.data.name,
            description: serviceDetailResult.data.description
          })
        } else {
          console.log('❌ 服务详情接口调用失败:', serviceDetailResult.message)
        }
        
      } else {
        console.log('❌ 没有获取到服务数据')
      }
      
    } else {
      console.log('❌ 首页init接口调用失败:', result.message)
    }
    
  } catch (error) {
    console.error('测试失败:', error)
  }
  
  console.log('=== 测试完成 ===')
}

// 导出测试函数
module.exports = {
  testServiceClickFeature
}

// 如果直接运行此文件，执行测试
if (typeof wx !== 'undefined') {
  // 在小程序环境中
  testServiceClickFeature()
} else {
  // 在Node.js环境中
  console.log('请在微信小程序环境中运行此测试')
} 