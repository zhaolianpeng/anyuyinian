// 前端医院详情API测试脚本
const { api } = require('../utils/cloud-container-standard')

/**
 * 测试医院详情API调用
 */
async function testHospitalDetailAPI() {
  console.log('=== 测试医院详情API调用 ===')
  
  const testCases = [
    { id: 1, description: '第一个医院' },
    { id: 5, description: '第五个医院' },
    { id: 999, description: '不存在的医院' }
  ]
  
  for (const testCase of testCases) {
    console.log(`\n测试: ${testCase.description} (ID: ${testCase.id})`)
    
    try {
      const result = await api.hospitalDetail(testCase.id)
      console.log('✅ API调用成功')
      console.log('返回结果:', result)
      
      if (result.code === 0 && result.data) {
        const { hospital, navigation } = result.data
        console.log('医院信息:')
        console.log('- 名称:', hospital.name)
        console.log('- 等级:', hospital.level)
        console.log('- 类型:', hospital.type)
        console.log('- 地址:', hospital.address)
        console.log('- 电话:', hospital.phone)
        
        if (navigation) {
          console.log('导航信息:')
          console.log('- 距离:', navigation.distance, '公里')
          console.log('- 预计时间:', navigation.duration, '分钟')
        }
      } else {
        console.log('❌ API返回错误:', result.message)
      }
    } catch (error) {
      console.log('❌ API调用失败:', error.message)
    }
  }
}

/**
 * 测试医院列表API调用
 */
async function testHospitalListAPI() {
  console.log('\n=== 测试医院列表API调用 ===')
  
  try {
    const params = {
      longitude: 121.4737,
      latitude: 31.2304,
      page: 1,
      pageSize: 10
    }
    
    console.log('请求参数:', params)
    const result = await api.hospitalList(params)
    
    if (result.code === 0 && result.data) {
      console.log('✅ 医院列表API调用成功')
      console.log('医院数量:', result.data.list?.length || 0)
      
      if (result.data.list && result.data.list.length > 0) {
        console.log('医院列表:')
        result.data.list.forEach((hospital, index) => {
          console.log(`${index + 1}. ${hospital.name} (ID: ${hospital.id})`)
        })
      }
    } else {
      console.log('❌ 医院列表API返回错误:', result.message)
    }
  } catch (error) {
    console.log('❌ 医院列表API调用失败:', error.message)
  }
}

/**
 * 运行所有测试
 */
async function runAllTests() {
  console.log('🏥 开始医院详情API前端测试\n')
  
  await testHospitalListAPI()
  await testHospitalDetailAPI()
  
  console.log('\n🏥 医院详情API前端测试完成')
}

// 如果直接运行此脚本
if (typeof module !== 'undefined' && module.exports) {
  module.exports = {
    testHospitalDetailAPI,
    testHospitalListAPI,
    runAllTests
  }
}

// 在微信小程序环境中运行测试
if (typeof wx !== 'undefined') {
  // 延迟执行，确保页面加载完成
  setTimeout(() => {
    runAllTests()
  }, 1000)
} 