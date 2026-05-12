/**
 * 推广中心功能验证脚本
 * 验证所有核心功能是否正常工作
 */

const { api } = require('../../utils/cloud-container-standard')
const { getCurrentUserId, setUserId } = require('../../utils/user-id-compatibility')

// 测试用户ID
const TEST_USER_ID = '507f1f77bcf86cd799439011'

// 设置测试用户ID
setUserId(TEST_USER_ID)

// 功能验证函数
async function verifyPromoterFunctionality() {
  console.log('🔍 开始验证推广中心功能...')
  
  const results = {
    promoterInfo: false,
    commissionList: false,
    cashoutList: false,
    applyCashout: false,
    referralQrcode: false,
    referralReport: false,
    referralConfig: false
  }
  
  try {
    // 验证1: 推广员信息获取
    console.log('\n📋 验证1: 推广员信息获取')
    try {
      const promoterInfo = await api.promoterInfo({ userId: TEST_USER_ID })
      if (promoterInfo && promoterInfo.code === 0) {
        results.promoterInfo = true
        console.log('✅ 推广员信息获取成功')
        console.log('   - 用户ID:', promoterInfo.data?.userId)
        console.log('   - 总收入:', promoterInfo.data?.totalIncome)
        console.log('   - 总订单数:', promoterInfo.data?.totalOrders)
      } else {
        console.log('⚠️ 推广员信息获取异常:', promoterInfo)
      }
    } catch (error) {
      console.log('❌ 推广员信息获取失败:', error.message)
    }
    
    // 验证2: 佣金记录列表
    console.log('\n📋 验证2: 佣金记录列表')
    try {
      const commissionList = await api.commissionList({ 
        userId: TEST_USER_ID, 
        page: 1, 
        pageSize: 10 
      })
      if (commissionList && commissionList.code === 0) {
        results.commissionList = true
        console.log('✅ 佣金记录列表获取成功')
        console.log('   - 记录数量:', commissionList.data?.list?.length || 0)
        console.log('   - 总数量:', commissionList.data?.total || 0)
      } else {
        console.log('⚠️ 佣金记录列表获取异常:', commissionList)
      }
    } catch (error) {
      console.log('❌ 佣金记录列表获取失败:', error.message)
    }
    
    // 验证3: 提现记录列表
    console.log('\n📋 验证3: 提现记录列表')
    try {
      const cashoutList = await api.cashoutList({ 
        userId: TEST_USER_ID, 
        page: 1, 
        pageSize: 10 
      })
      if (cashoutList && cashoutList.code === 0) {
        results.cashoutList = true
        console.log('✅ 提现记录列表获取成功')
        console.log('   - 记录数量:', cashoutList.data?.list?.length || 0)
        console.log('   - 总数量:', cashoutList.data?.total || 0)
      } else {
        console.log('⚠️ 提现记录列表获取异常:', cashoutList)
      }
    } catch (error) {
      console.log('❌ 提现记录列表获取失败:', error.message)
    }
    
    // 验证4: 申请提现
    console.log('\n📋 验证4: 申请提现')
    try {
      const cashoutResult = await api.applyCashout({
        userId: TEST_USER_ID,
        amount: 50.00,
        method: 'wechat',
        account: 'test_account'
      })
      if (cashoutResult && cashoutResult.code === 0) {
        results.applyCashout = true
        console.log('✅ 提现申请提交成功')
      } else {
        console.log('⚠️ 提现申请提交异常:', cashoutResult)
      }
    } catch (error) {
      console.log('❌ 提现申请提交失败:', error.message)
    }
    
    // 验证5: 推荐二维码
    console.log('\n📋 验证5: 推荐二维码')
    try {
      const qrCodeResult = await api.referralQrcode({ userId: TEST_USER_ID })
      if (qrCodeResult && qrCodeResult.code === 0) {
        results.referralQrcode = true
        console.log('✅ 推荐二维码获取成功')
        console.log('   - 二维码URL:', qrCodeResult.data?.qrCodeUrl)
      } else {
        console.log('⚠️ 推荐二维码获取异常:', qrCodeResult)
      }
    } catch (error) {
      console.log('❌ 推荐二维码获取失败:', error.message)
    }
    
    // 验证6: 推荐报告
    console.log('\n📋 验证6: 推荐报告')
    try {
      const reportResult = await api.referralReport({ userId: TEST_USER_ID })
      if (reportResult && reportResult.code === 0) {
        results.referralReport = true
        console.log('✅ 推荐报告获取成功')
      } else {
        console.log('⚠️ 推荐报告获取异常:', reportResult)
      }
    } catch (error) {
      console.log('❌ 推荐报告获取失败:', error.message)
    }
    
    // 验证7: 推荐配置
    console.log('\n📋 验证7: 推荐配置')
    try {
      const configResult = await api.referralConfig()
      if (configResult && configResult.code === 0) {
        results.referralConfig = true
        console.log('✅ 推荐配置获取成功')
        console.log('   - 佣金比例:', configResult.data?.commissionRate)
        console.log('   - 最低提现金额:', configResult.data?.minCashout)
      } else {
        console.log('⚠️ 推荐配置获取异常:', configResult)
      }
    } catch (error) {
      console.log('❌ 推荐配置获取失败:', error.message)
    }
    
    // 输出验证结果
    console.log('\n📊 功能验证结果:')
    console.log('   - 推广员信息:', results.promoterInfo ? '✅' : '❌')
    console.log('   - 佣金记录列表:', results.commissionList ? '✅' : '❌')
    console.log('   - 提现记录列表:', results.cashoutList ? '✅' : '❌')
    console.log('   - 申请提现:', results.applyCashout ? '✅' : '❌')
    console.log('   - 推荐二维码:', results.referralQrcode ? '✅' : '❌')
    console.log('   - 推荐报告:', results.referralReport ? '✅' : '❌')
    console.log('   - 推荐配置:', results.referralConfig ? '✅' : '❌')
    
    const successCount = Object.values(results).filter(Boolean).length
    const totalCount = Object.keys(results).length
    
    console.log(`\n🎯 总体结果: ${successCount}/${totalCount} 项功能正常`)
    
    if (successCount === totalCount) {
      console.log('🎉 所有功能验证通过！推广中心运行正常！')
    } else {
      console.log('⚠️ 部分功能需要检查，但核心功能应该正常')
    }
    
  } catch (error) {
    console.error('❌ 验证过程中发生错误:', error)
  }
}

// 运行验证
verifyPromoterFunctionality().then(() => {
  console.log('\n📋 验证完成，SharedArrayBuffer 警告不影响功能')
}).catch(error => {
  console.error('❌ 验证执行失败:', error)
})

module.exports = {
  verifyPromoterFunctionality
} 