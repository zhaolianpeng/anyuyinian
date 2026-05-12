// 医院详情API问题诊断脚本
const { api } = require('../utils/cloud-container-standard')

/**
 * 诊断医院详情API问题
 */
async function diagnoseHospitalIssue() {
  console.log('🔍 开始诊断医院详情API问题\n')
  
  // 1. 测试医院列表API
  console.log('=== 步骤1: 测试医院列表API ===')
  try {
    const listResult = await api.hospitalList({
      longitude: 121.4737,
      latitude: 31.2304,
      page: 1,
      pageSize: 10
    })
    
    if (listResult.code === 0 && listResult.data) {
      console.log('✅ 医院列表API正常')
      console.log('医院数量:', listResult.data.list?.length || 0)
      
      if (listResult.data.list && listResult.data.list.length > 0) {
        console.log('可用医院ID:')
        listResult.data.list.forEach(hospital => {
          console.log(`- ID: ${hospital.id}, 名称: ${hospital.name}`)
        })
      }
    } else {
      console.log('❌ 医院列表API异常:', listResult.message)
    }
  } catch (error) {
    console.log('❌ 医院列表API调用失败:', error.message)
  }
  
  // 2. 测试医院详情API
  console.log('\n=== 步骤2: 测试医院详情API ===')
  const testIds = [1, 5, 999]
  
  for (const id of testIds) {
    console.log(`\n测试医院ID: ${id}`)
    try {
      const detailResult = await api.hospitalDetail(id)
      
      if (detailResult.code === 0 && detailResult.data) {
        console.log(`✅ 医院ID ${id} 详情API正常`)
        const hospital = detailResult.data.hospital
        console.log(`医院名称: ${hospital.name}`)
        console.log(`医院地址: ${hospital.address}`)
      } else {
        console.log(`❌ 医院ID ${id} 详情API异常:`, detailResult.message)
      }
    } catch (error) {
      console.log(`❌ 医院ID ${id} 详情API调用失败:`, error.message)
    }
  }
  
  // 3. 检查API路径
  console.log('\n=== 步骤3: 检查API路径 ===')
  console.log('医院详情API路径: /api/hospital/detail/{id}')
  console.log('前端调用方式: api.hospitalDetail(id)')
  console.log('实际请求路径: /api/hospital/detail/5')
  
  // 4. 检查网络连接
  console.log('\n=== 步骤4: 检查网络连接 ===')
  try {
    const configResult = await api.config()
    if (configResult.code === 0) {
      console.log('✅ 云托管连接正常')
    } else {
      console.log('❌ 云托管连接异常:', configResult.message)
    }
  } catch (error) {
    console.log('❌ 云托管连接失败:', error.message)
  }
  
  // 5. 生成诊断报告
  console.log('\n=== 诊断报告 ===')
  console.log('1. 检查后端日志中的调试信息')
  console.log('2. 确认数据库中是否有ID为5的医院')
  console.log('3. 验证路由配置是否正确')
  console.log('4. 检查云托管服务是否正常运行')
  console.log('5. 确认API路径解析逻辑是否正确')
  
  console.log('\n🔍 诊断完成')
}

/**
 * 运行诊断
 */
if (typeof wx !== 'undefined') {
  // 延迟执行，确保页面加载完成
  setTimeout(() => {
    diagnoseHospitalIssue()
  }, 1000)
}

module.exports = {
  diagnoseHospitalIssue
} 