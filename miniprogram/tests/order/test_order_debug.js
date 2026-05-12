// 调试预约页面数据加载问题
const app = getApp()

// 测试预约页面数据加载
function testOrderPageData() {
  console.log('=== 开始调试预约页面数据加载 ===')
  
  const userId = wx.getStorageSync('userId')
  console.log('当前用户ID:', userId)
  
  if (!userId) {
    console.error('❌ 用户未登录，无法获取用户ID')
    return
  }

  // 测试就诊人数据
  console.log('1. 测试就诊人数据加载...')
  wx.request({
    url: `${app.globalData.baseUrl}/api/user/patient`,
    method: 'GET',
    data: { userId },
    success: (res) => {
      console.log('就诊人API响应:', res.data)
      if (res.data.code === 0) {
        const patientList = res.data.data || []
        console.log('✅ 就诊人数据加载成功')
        console.log('就诊人数量:', patientList.length)
        console.log('就诊人列表:', patientList)
        
        if (patientList.length > 0) {
          console.log('第一个就诊人:', patientList[0])
          console.log('就诊人字段检查:')
          console.log('- id:', patientList[0].id)
          console.log('- name:', patientList[0].name)
          console.log('- relation:', patientList[0].relation)
          console.log('- phone:', patientList[0].phone)
          console.log('- isDefault:', patientList[0].isDefault)
        }
      } else {
        console.error('❌ 就诊人数据加载失败:', res.data)
      }
    },
    fail: (err) => {
      console.error('❌ 就诊人API请求失败:', err)
    }
  })

  // 测试地址数据
  console.log('2. 测试地址数据加载...')
  wx.request({
    url: `${app.globalData.baseUrl}/api/user/address`,
    method: 'GET',
    data: { userId },
    success: (res) => {
      console.log('地址API响应:', res.data)
      if (res.data.code === 0) {
        const addressList = res.data.data || []
        console.log('✅ 地址数据加载成功')
        console.log('地址数量:', addressList.length)
        console.log('地址列表:', addressList)
        
        if (addressList.length > 0) {
          console.log('第一个地址:', addressList[0])
          console.log('地址字段检查:')
          console.log('- id:', addressList[0].id)
          console.log('- name:', addressList[0].name)
          console.log('- phone:', addressList[0].phone)
          console.log('- province:', addressList[0].province)
          console.log('- city:', addressList[0].city)
          console.log('- district:', addressList[0].district)
          console.log('- address:', addressList[0].address)
          console.log('- isDefault:', addressList[0].isDefault)
        }
      } else {
        console.error('❌ 地址数据加载失败:', res.data)
      }
    },
    fail: (err) => {
      console.error('❌ 地址API请求失败:', err)
    }
  })

  // 测试预约时间范围
  console.log('3. 测试预约时间范围...')
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
  
  console.log('✅ 预约时间范围设置成功')
  console.log('今天:', formatDate(today))
  console.log('明天:', minDate)
  console.log('7天后:', maxDateStr)
  console.log('时间范围:', `${minDate} 至 ${maxDateStr}`)
}

// 测试页面状态
function testPageState() {
  console.log('4. 测试页面状态...')
  
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1]
  
  if (currentPage && currentPage.data) {
    console.log('✅ 获取到当前页面数据')
    console.log('页面数据:', currentPage.data)
    console.log('就诊人列表长度:', currentPage.data.patientList?.length || 0)
    console.log('地址列表长度:', currentPage.data.addressList?.length || 0)
    console.log('选中的就诊人:', currentPage.data.selectedPatient)
    console.log('选中的地址:', currentPage.data.selectedAddress)
    console.log('预约时间范围:', {
      minDate: currentPage.data.minDate,
      maxDate: currentPage.data.maxDate
    })
    console.log('表单数据:', currentPage.data.formData)
  } else {
    console.error('❌ 无法获取当前页面数据')
  }
}

// 模拟数据设置
function simulateDataSetting() {
  console.log('5. 模拟数据设置...')
  
  const mockData = {
    patientList: [
      {
        id: 1,
        name: '赵连鹏',
        relation: '本人',
        phone: '13691028481',
        isDefault: 0
      }
    ],
    addressList: [
      {
        id: 1,
        name: '赵连鹏',
        phone: '13691028481',
        province: '北京市',
        city: '北京市',
        district: '西城区',
        address: '复兴门北大街9号楼',
        isDefault: 0
      }
    ]
  }
  
  console.log('模拟数据:', mockData)
  console.log('就诊人列表长度 > 0:', mockData.patientList.length > 0)
  console.log('地址列表长度 > 0:', mockData.addressList.length > 0)
  
  // 检查默认项
  const defaultPatient = mockData.patientList.find(p => p.isDefault === 1)
  const defaultAddress = mockData.addressList.find(a => a.isDefault === 1)
  
  console.log('默认就诊人:', defaultPatient)
  console.log('默认地址:', defaultAddress)
  
  // 选择第一个作为默认
  const selectedPatient = defaultPatient || mockData.patientList[0]
  const selectedAddress = defaultAddress || mockData.addressList[0]
  
  console.log('选中的就诊人:', selectedPatient)
  console.log('选中的地址:', selectedAddress)
}

// 运行所有测试
function runAllTests() {
  console.log('=== 开始调试预约页面 ===')
  
  testOrderPageData()
  
  // 延迟执行页面状态测试
  setTimeout(() => {
    testPageState()
    simulateDataSetting()
    console.log('=== 调试完成 ===')
  }, 2000)
}

// 导出测试函数
module.exports = {
  testOrderPageData,
  testPageState,
  simulateDataSetting,
  runAllTests
}

// 如果直接运行此文件
if (typeof wx !== 'undefined') {
  // 在微信小程序环境中
  setTimeout(() => {
    runAllTests()
  }, 1000)
} 