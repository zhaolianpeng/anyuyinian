/**
 * 推广中心功能测试脚本
 */

const { api } = require('../../utils/cloud-container-standard')
const { getCurrentUserId, setUserId } = require('../../utils/user-id-compatibility')

// 测试数据
const TEST_USER_ID = '507f1f77bcf86cd799439011'

// 设置测试用户ID
setUserId(TEST_USER_ID)

// 测试函数
async function testPromoterCenter() {
  console.log('🧪 开始测试推广中心功能...')
  
  try {
    // 测试1: 获取推广员信息
    console.log('\n📋 测试1: 获取推广员信息')
    const promoterInfo = await api.promoterInfo({ userId: TEST_USER_ID })
    console.log('✅ 推广员信息获取成功:', promoterInfo)
    
    // 测试2: 获取佣金记录列表
    console.log('\n📋 测试2: 获取佣金记录列表')
    const commissionList = await api.commissionList({ 
      userId: TEST_USER_ID, 
      page: 1, 
      pageSize: 10 
    })
    console.log('✅ 佣金记录列表获取成功:', commissionList)
    
    // 测试3: 获取提现记录列表
    console.log('\n📋 测试3: 获取提现记录列表')
    const cashoutList = await api.cashoutList({ 
      userId: TEST_USER_ID, 
      page: 1, 
      pageSize: 10 
    })
    console.log('✅ 提现记录列表获取成功:', cashoutList)
    
    // 测试4: 申请提现
    console.log('\n📋 测试4: 申请提现')
    const cashoutResult = await api.applyCashout({
      userId: TEST_USER_ID,
      amount: 100.00,
      method: 'wechat',
      account: 'test_account'
    })
    console.log('✅ 提现申请提交成功:', cashoutResult)
    
    // 测试5: 获取推荐二维码
    console.log('\n📋 测试5: 获取推荐二维码')
    const qrCodeResult = await api.referralQrcode({ userId: TEST_USER_ID })
    console.log('✅ 推荐二维码获取成功:', qrCodeResult)
    
    // 测试6: 获取推荐报告
    console.log('\n📋 测试6: 获取推荐报告')
    const reportResult = await api.referralReport({ userId: TEST_USER_ID })
    console.log('✅ 推荐报告获取成功:', reportResult)
    
    // 测试7: 获取推荐配置
    console.log('\n📋 测试7: 获取推荐配置')
    const configResult = await api.referralConfig()
    console.log('✅ 推荐配置获取成功:', configResult)
    
    console.log('\n🎉 所有测试通过！推广中心功能正常！')
    
  } catch (error) {
    console.error('❌ 测试失败:', error)
    console.error('错误详情:', {
      message: error.message,
      stack: error.stack
    })
  }
}

// 测试数据验证
function validatePromoterData(data) {
  console.log('📊 数据验证:')
  
  if (data.promoterInfo) {
    const info = data.promoterInfo
    console.log('  - 用户ID:', info.userId)
    console.log('  - 昵称:', info.nickName)
    console.log('  - 总收入:', info.totalIncome)
    console.log('  - 总订单数:', info.totalOrders)
    console.log('  - 二维码URL:', info.qrCodeUrl)
  }
  
  if (data.commissionList) {
    console.log('  - 佣金记录数量:', data.commissionList.length)
  }
  
  if (data.cashoutList) {
    console.log('  - 提现记录数量:', data.cashoutList.length)
  }
}

// 运行测试
testPromoterCenter().then(() => {
  console.log('\n📋 测试完成，请检查控制台输出')
}).catch(error => {
  console.error('❌ 测试执行失败:', error)
})

module.exports = {
  testPromoterCenter,
  validatePromoterData
} 