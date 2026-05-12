// 测试预约页面修复效果
const app = getApp()

// 测试数据加载
function testDataLoading() {
  console.log('=== 测试预约页面数据加载 ===')
  
  const userId = wx.getStorageSync('userId')
  console.log('用户ID:', userId)
  
  if (!userId) {
    console.error('用户未登录')
    return
  }

  // 测试就诊人数据
  wx.request({
    url: `${app.globalData.baseUrl}/api/user/patient`,
    method: 'GET',
    data: { userId },
    success: (res) => {
      console.log('就诊人API响应:', res.data)
      if (res.data.code === 0) {
        const patientList = res.data.data || []
        console.log('✅ 就诊人数据加载成功，数量:', patientList.length)
        
        if (patientList.length > 0) {
          console.log('第一个就诊人:', patientList[0])
          console.log('字段检查:')
          console.log('- name:', patientList[0].name)
          console.log('- relation:', patientList[0].relation)
          console.log('- phone:', patientList[0].phone)
          console.log('- isDefault:', patientList[0].isDefault)
        }
      }
    }
  })

  // 测试地址数据
  wx.request({
    url: `${app.globalData.baseUrl}/api/user/address`,
    method: 'GET',
    data: { userId },
    success: (res) => {
      console.log('地址API响应:', res.data)
      if (res.data.code === 0) {
        const addressList = res.data.data || []
        console.log('✅ 地址数据加载成功，数量:', addressList.length)
        
        if (addressList.length > 0) {
          console.log('第一个地址:', addressList[0])
          console.log('字段检查:')
          console.log('- name:', addressList[0].name)
          console.log('- phone:', addressList[0].phone)
          console.log('- province:', addressList[0].province)
          console.log('- city:', addressList[0].city)
          console.log('- district:', addressList[0].district)
          console.log('- address:', addressList[0].address)
          console.log('- isDefault:', addressList[0].isDefault)
        }
      }
    }
  })
}

// 测试页面显示
function testPageDisplay() {
  console.log('=== 测试页面显示 ===')
  
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1]
  
  if (currentPage && currentPage.data) {
    const data = currentPage.data
    console.log('页面数据状态:')
    console.log('- 就诊人列表长度:', data.patientList?.length || 0)
    console.log('- 地址列表长度:', data.addressList?.length || 0)
    console.log('- 选中的就诊人:', data.selectedPatient)
    console.log('- 选中的地址:', data.selectedAddress)
    console.log('- 预约时间范围:', {
      minDate: data.minDate,
      maxDate: data.maxDate
    })
    console.log('- 表单数据:', data.formData)
    
    // 检查数据是否正确设置
    if (data.patientList && data.patientList.length > 0) {
      console.log('✅ 就诊人数据已加载')
    } else {
      console.log('❌ 就诊人数据未加载')
    }
    
    if (data.addressList && data.addressList.length > 0) {
      console.log('✅ 地址数据已加载')
    } else {
      console.log('❌ 地址数据未加载')
    }
    
    if (data.selectedPatient) {
      console.log('✅ 已选中就诊人')
    } else {
      console.log('❌ 未选中就诊人')
    }
    
    if (data.selectedAddress) {
      console.log('✅ 已选中地址')
    } else {
      console.log('❌ 未选中地址')
    }
  }
}

// 测试预约时间范围
function testAppointmentTime() {
  console.log('=== 测试预约时间范围 ===')
  
  const today = new Date()
  const tomorrow = new Date(today)
  tomorrow.setDate(today.getDate() + 1)
  
  const maxDate = new Date(today)
  maxDate.setDate(today.getDate() + 7)
  
  const formatDate = (date) => {
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    return `${year}-${month}-${day}`
  }
  
  const minDate = formatDate(tomorrow)
  const maxDateStr = formatDate(maxDate)
  
  console.log('预约时间范围:')
  console.log('- 今天:', formatDate(today))
  console.log('- 明天:', minDate)
  console.log('- 7天后:', maxDateStr)
  console.log('- 可预约范围:', `${minDate} 至 ${maxDateStr}`)
  
  // 验证时间范围是否正确
  const tomorrowDate = new Date(tomorrow)
  const maxDateObj = new Date(maxDate)
  const diffDays = Math.ceil((maxDateObj - tomorrowDate) / (1000 * 60 * 60 * 24))
  
  console.log('时间范围验证:')
  console.log('- 从明天开始:', tomorrowDate > today)
  console.log('- 7天范围:', diffDays === 7)
}

// 测试表单验证
function testFormValidation() {
  console.log('=== 测试表单验证 ===')
  
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1]
  
  if (currentPage && currentPage.data) {
    const data = currentPage.data
    const canSubmit = data.selectedPatient && 
                     data.selectedAddress && 
                     data.formData.appointmentDate && 
                     data.formData.appointmentTime
    
    console.log('表单验证状态:')
    console.log('- 选中就诊人:', !!data.selectedPatient)
    console.log('- 选中地址:', !!data.selectedAddress)
    console.log('- 选择日期:', !!data.formData.appointmentDate)
    console.log('- 选择时间:', !!data.formData.appointmentTime)
    console.log('- 可以提交:', canSubmit)
  }
}

// 运行所有测试
function runAllTests() {
  console.log('=== 开始测试预约页面修复效果 ===')
  
  testDataLoading()
  testAppointmentTime()
  
  // 延迟执行页面相关测试
  setTimeout(() => {
    testPageDisplay()
    testFormValidation()
    console.log('=== 测试完成 ===')
  }, 2000)
}

// 导出测试函数
module.exports = {
  testDataLoading,
  testPageDisplay,
  testAppointmentTime,
  testFormValidation,
  runAllTests
}

// 如果直接运行此文件
if (typeof wx !== 'undefined') {
  setTimeout(() => {
    runAllTests()
  }, 1000)
} 