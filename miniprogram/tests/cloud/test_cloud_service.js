// 测试云托管服务连接状态
const { callContainer } = require('../utils/cloud-container-standard')

async function testCloudService() {
  console.log('开始测试云托管服务连接...')
  
  try {
    // 测试基础HTTP连接
    const result = await callContainer('/api/count', 'GET')
    
    console.log('✅ 云托管服务HTTP连接正常:', result)
    
    // 测试WebSocket配置
    const wsConfig = {
      env: 'prod-5g94mx7a3d07e78c',
      service: 'golang-lfwy',
      path: '/ws'
    }
    
    console.log('WebSocket配置:', wsConfig)
    console.log('WebSocket URL: wss://golang-lfwy-prod-5g94mx7a3d07e78c-1353115175.ap-shanghai.run.wxcloudrun.com/ws')
    
    return {
      success: true,
      httpStatus: '正常',
      wsConfig: wsConfig
    }
    
  } catch (error) {
    console.error('❌ 云托管服务连接失败:', error)
    return {
      success: false,
      error: error.message
    }
  }
}

// 导出测试函数
module.exports = {
  testCloudService
}

// 如果直接运行此文件
if (typeof wx !== 'undefined') {
  testCloudService().then(result => {
    console.log('测试结果:', result)
  })
} 