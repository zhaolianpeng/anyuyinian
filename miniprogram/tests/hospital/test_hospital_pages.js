// 医院信息页面测试脚本
const { api } = require('../utils/cloud-container-standard')

/**
 * 测试医院列表API
 */
async function testHospitalListAPI() {
  console.log('=== 测试医院列表API ===')
  
  try {
    const params = {
      longitude: 121.4737,
      latitude: 31.2304,
      page: 1,
      pageSize: 10
    }
    
    console.log('请求参数:', params)
    const result = await api.hospitalList(params)
    console.log('医院列表API返回:', result)
    
    if (result.code === 0 && result.data) {
      console.log('✅ 医院列表API测试成功')
      console.log('医院数量:', result.data.list?.length || 0)
      console.log('是否有更多:', result.data.hasMore)
      
      // 显示医院信息
      if (result.data.list && result.data.list.length > 0) {
        console.log('医院信息示例:')
        const hospital = result.data.list[0]
        console.log('- 名称:', hospital.name)
        console.log('- 等级:', hospital.level)
        console.log('- 类型:', hospital.type)
        console.log('- 地址:', hospital.address)
        console.log('- 电话:', hospital.phone)
        console.log('- 坐标:', hospital.longitude, hospital.latitude)
      }
    } else {
      console.log('❌ 医院列表API测试失败:', result.message)
    }
  } catch (error) {
    console.error('❌ 医院列表API测试异常:', error)
  }
}

/**
 * 测试医院详情API
 */
async function testHospitalDetailAPI() {
  console.log('\n=== 测试医院详情API ===')
  
  try {
    const hospitalId = 1 // 测试第一个医院
    const params = {
      userLongitude: 121.4737,
      userLatitude: 31.2304
    }
    
    console.log('请求参数:', { hospitalId, ...params })
    const result = await api.hospitalDetail(hospitalId, params)
    console.log('医院详情API返回:', result)
    
    if (result.code === 0 && result.data) {
      console.log('✅ 医院详情API测试成功')
      
      const { hospital, navigation } = result.data
      console.log('医院信息:')
      console.log('- 名称:', hospital.name)
      console.log('- 等级:', hospital.level)
      console.log('- 类型:', hospital.type)
      console.log('- 地址:', hospital.address)
      console.log('- 电话:', hospital.phone)
      console.log('- 描述:', hospital.description)
      
      if (navigation) {
        console.log('导航信息:')
        console.log('- 距离:', navigation.distance, '公里')
        console.log('- 预计时间:', navigation.duration, '分钟')
        console.log('- 路线类型:', navigation.routeType)
      } else {
        console.log('⚠️ 无导航信息')
      }
    } else {
      console.log('❌ 医院详情API测试失败:', result.message)
    }
  } catch (error) {
    console.error('❌ 医院详情API测试异常:', error)
  }
}

/**
 * 测试首页医院数据
 */
async function testHomeHospitalData() {
  console.log('\n=== 测试首页医院数据 ===')
  
  try {
    const params = {
      longitude: 121.4737,
      latitude: 31.2304,
      limit: 5
    }
    
    console.log('请求参数:', params)
    const result = await api.homeInit(params)
    console.log('首页初始化API返回:', result)
    
    if (result.code === 0 && result.data) {
      console.log('✅ 首页医院数据测试成功')
      
      const hospitals = result.data.hospitals || []
      console.log('医院数量:', hospitals.length)
      
      if (hospitals.length > 0) {
        console.log('医院列表:')
        hospitals.forEach((hospital, index) => {
          console.log(`${index + 1}. ${hospital.name} (${hospital.level})`)
        })
      }
    } else {
      console.log('❌ 首页医院数据测试失败:', result.message)
    }
  } catch (error) {
    console.error('❌ 首页医院数据测试异常:', error)
  }
}

/**
 * 运行所有测试
 */
async function runAllTests() {
  console.log('🏥 开始医院信息页面功能测试\n')
  
  await testHospitalListAPI()
  await testHospitalDetailAPI()
  await testHomeHospitalData()
  
  console.log('\n🏥 医院信息页面功能测试完成')
}

// 如果直接运行此脚本
if (typeof module !== 'undefined' && module.exports) {
  module.exports = {
    testHospitalListAPI,
    testHospitalDetailAPI,
    testHomeHospitalData,
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