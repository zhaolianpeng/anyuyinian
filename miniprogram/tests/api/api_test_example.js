// api_test_example.js
// API测试示例 - 服务详情POST请求

const { request } = require('./utils/request')

// 示例1: 获取ID为7的服务详情
async function getServiceById7() {
  console.log('=== 获取ID为7的服务详情 ===')
  
  try {
    const res = await request('/api/service/detail', 'POST', { serviceId: 7 })
    
    if (res.code === 0 && res.data) {
      const service = res.data
      console.log('✅ 成功获取服务信息:')
      console.log(`   服务名称: ${service.name}`)
      console.log(`   服务分类: ${service.category}`)
      console.log(`   服务价格: ¥${service.price}`)
      console.log(`   原价: ¥${service.originalPrice}`)
      console.log(`   服务描述: ${service.description}`)
      
      // 解析表单配置
      if (service.formConfig) {
        try {
          const formConfig = JSON.parse(service.formConfig)
          console.log(`   表单字段数: ${formConfig.fields ? formConfig.fields.length : 0}`)
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

// 示例2: 批量获取服务详情
async function getMultipleServices() {
  console.log('\n=== 批量获取服务详情 ===')
  
  const serviceIds = [1, 2, 3, 4, 5, 6, 7, 8]
  const results = []
  
  for (const serviceId of serviceIds) {
    try {
      const res = await request('/api/service/detail', 'POST', { serviceId })
      
      if (res.code === 0 && res.data) {
        results.push({
          id: serviceId,
          name: res.data.name,
          price: res.data.price,
          status: 'success'
        })
      } else {
        results.push({
          id: serviceId,
          error: res.errorMsg,
          status: 'failed'
        })
      }
    } catch (error) {
      results.push({
        id: serviceId,
        error: error.message,
        status: 'error'
      })
    }
  }
  
  console.log('批量获取结果:')
  results.forEach(result => {
    if (result.status === 'success') {
      console.log(`✅ ID ${result.id}: ${result.name} (¥${result.price})`)
    } else {
      console.log(`❌ ID ${result.id}: ${result.error}`)
    }
  })
}

// 示例3: 错误处理测试
async function testErrorHandling() {
  console.log('\n=== 错误处理测试 ===')
  
  const testCases = [
    { serviceId: 999, description: '不存在的服务ID' },
    { serviceId: 0, description: '无效的服务ID 0' },
    { serviceId: -1, description: '无效的服务ID -1' },
    { serviceId: 'abc', description: '字符串ID' }
  ]
  
  for (const testCase of testCases) {
    console.log(`\n--- ${testCase.description} ---`)
    
    try {
      const res = await request('/api/service/detail', 'POST', { 
        serviceId: testCase.serviceId 
      })
      
      console.log('请求参数:', { serviceId: testCase.serviceId })
      console.log('响应结果:', res)
    } catch (error) {
      console.error('请求失败:', error.message)
    }
  }
}

// 示例4: 在小程序页面中使用
function exampleInPage() {
  console.log('\n=== 在小程序页面中的使用示例 ===')
  
  // 在页面的onLoad方法中
  const page = {
    data: {
      serviceId: 7,
      service: null,
      loading: false
    },
    
    onLoad(options) {
      const serviceId = options.id || this.data.serviceId
      this.loadServiceDetail(serviceId)
    },
    
    async loadServiceDetail(serviceId) {
      this.setData({ loading: true })
      
      try {
        const res = await request('/api/service/detail', 'POST', { serviceId })
        
        if (res.code === 0 && res.data) {
          this.setData({ 
            service: res.data,
            loading: false 
          })
        } else {
          wx.showToast({
            title: res.errorMsg || '加载失败',
            icon: 'none'
          })
        }
      } catch (error) {
        console.error('加载服务详情失败:', error)
        wx.showToast({
          title: '网络错误，请重试',
          icon: 'none'
        })
      } finally {
        this.setData({ loading: false })
      }
    }
  }
  
  console.log('页面示例代码已生成')
  return page
}

// 运行所有示例
async function runAllExamples() {
  console.log('开始运行API测试示例...')
  
  // 1. 获取ID为7的服务
  await getServiceById7()
  
  // 2. 批量获取服务
  await getMultipleServices()
  
  // 3. 错误处理测试
  await testErrorHandling()
  
  // 4. 页面使用示例
  exampleInPage()
  
  console.log('\n=== 所有示例运行完成 ===')
}

// 导出函数
module.exports = {
  getServiceById7,
  getMultipleServices,
  testErrorHandling,
  exampleInPage,
  runAllExamples
}

// 如果直接运行此文件
if (require.main === module) {
  runAllExamples().then(() => {
    console.log('示例运行完成')
  }).catch(error => {
    console.error('运行示例时出现错误:', error)
  })
} 