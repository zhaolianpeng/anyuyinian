// 测试患者管理功能
const { api } = require('../utils/cloud-container-standard')

async function testPatientManagement() {
  console.log('=== 测试患者管理功能 ===')
  
  try {
    // 1. 检查用户登录状态
    const userId = wx.getStorageSync('userId')
    console.log('当前用户ID:', userId || '未登录')
    
    if (!userId) {
      console.log('❌ 用户未登录')
      wx.showToast({
        title: '请先登录',
        icon: 'none'
      })
      return
    }

    // 2. 测试患者列表API
    console.log('测试患者列表API...')
    const patientListResult = await api.userPatient({ userId })
    console.log('患者列表API结果:', patientListResult)
    
    if (patientListResult.code === 0) {
      console.log('✅ 患者列表API调用成功')
      const patients = patientListResult.data || []
      console.log('患者数量:', patients.length)
      
      if (patients.length > 0) {
        console.log('第一个患者信息:', {
          id: patients[0].id,
          name: patients[0].name,
          relation: patients[0].relation,
          isDefault: patients[0].isDefault
        })
      }
    } else {
      console.log('❌ 患者列表API调用失败:', patientListResult.message)
    }

    // 3. 测试页面跳转
    console.log('测试患者管理页面跳转...')
    wx.navigateTo({
      url: '/pages/user/patient/list',
      success: () => {
        console.log('✅ 患者管理页面跳转成功')
      },
      fail: (error) => {
        console.error('❌ 患者管理页面跳转失败:', error)
      }
    })

  } catch (error) {
    console.error('患者管理测试失败:', error)
  }
}

// 测试患者添加功能
async function testPatientAdd() {
  console.log('=== 测试患者添加功能 ===')
  
  try {
    const userId = wx.getStorageSync('userId')
    if (!userId) {
      console.log('❌ 用户未登录')
      return
    }

    // 测试患者添加API
    console.log('测试患者添加API...')
    const testPatientData = {
      userId: userId,
      name: '测试患者',
      idCard: '110101199001011234',
      phone: '13800138000',
      gender: 1,
      birthday: '1990-01-01',
      relation: '本人',
      isDefault: false
    }

    const addResult = await api.userPatientAdd(testPatientData)
    console.log('患者添加API结果:', addResult)
    
    if (addResult.code === 0) {
      console.log('✅ 患者添加API调用成功')
    } else {
      console.log('❌ 患者添加API调用失败:', addResult.message)
    }

  } catch (error) {
    console.error('患者添加测试失败:', error)
  }
}

// 测试患者编辑功能
async function testPatientEdit() {
  console.log('=== 测试患者编辑功能 ===')
  
  try {
    const userId = wx.getStorageSync('userId')
    if (!userId) {
      console.log('❌ 用户未登录')
      return
    }

    // 先获取患者列表
    const patientListResult = await api.userPatient({ userId })
    if (patientListResult.code === 0 && patientListResult.data && patientListResult.data.length > 0) {
      const firstPatient = patientListResult.data[0]
      
      // 测试患者编辑API
      console.log('测试患者编辑API...')
      const editData = {
        id: firstPatient.id,
        userId: userId,
        name: firstPatient.name + '_编辑',
        idCard: firstPatient.idCard,
        phone: firstPatient.phone,
        gender: firstPatient.gender,
        birthday: firstPatient.birthday,
        relation: firstPatient.relation,
        isDefault: firstPatient.isDefault
      }

      const editResult = await api.userPatientUpdate(editData)
      console.log('患者编辑API结果:', editResult)
      
      if (editResult.code === 0) {
        console.log('✅ 患者编辑API调用成功')
      } else {
        console.log('❌ 患者编辑API调用失败:', editResult.message)
      }
    } else {
      console.log('❌ 没有患者数据可供编辑')
    }

  } catch (error) {
    console.error('患者编辑测试失败:', error)
  }
}

// 检查页面配置
function checkPatientPageConfig() {
  console.log('=== 检查患者管理页面配置 ===')
  
  try {
    // 检查患者列表页面是否存在
    wx.navigateTo({
      url: '/pages/user/patient/list',
      success: () => {
        console.log('✅ 患者列表页面可以正常访问')
        // 立即返回
        setTimeout(() => {
          wx.navigateBack()
        }, 100)
      },
      fail: (error) => {
        console.error('❌ 患者列表页面访问失败:', error)
      }
    })
    
    // 检查患者添加页面是否存在
    setTimeout(() => {
      wx.navigateTo({
        url: '/pages/user/patient/add',
        success: () => {
          console.log('✅ 患者添加页面可以正常访问')
          // 立即返回
          setTimeout(() => {
            wx.navigateBack()
          }, 100)
        },
        fail: (error) => {
          console.error('❌ 患者添加页面访问失败:', error)
        }
      })
    }, 200)
    
  } catch (error) {
    console.error('页面配置检查失败:', error)
  }
}

// 运行所有测试
function runPatientTests() {
  console.log('🚀 开始运行患者管理功能测试...')
  
  // 检查页面配置
  checkPatientPageConfig()
  
  // 测试患者管理
  setTimeout(() => {
    testPatientManagement()
  }, 500)
  
  // 测试患者添加
  setTimeout(() => {
    testPatientAdd()
  }, 1000)
  
  // 测试患者编辑
  setTimeout(() => {
    testPatientEdit()
  }, 1500)
}

// 导出测试函数
module.exports = {
  testPatientManagement,
  testPatientAdd,
  testPatientEdit,
  checkPatientPageConfig,
  runPatientTests
}

// 如果直接运行此文件，执行测试
if (typeof wx !== 'undefined') {
  console.log('患者管理测试脚本已加载')
  console.log('运行 runPatientTests() 执行所有测试')
  console.log('运行 testPatientManagement() 测试患者管理')
  console.log('运行 testPatientAdd() 测试患者添加')
  console.log('运行 testPatientEdit() 测试患者编辑')
} else {
  console.log('请在微信小程序环境中运行此测试')
} 