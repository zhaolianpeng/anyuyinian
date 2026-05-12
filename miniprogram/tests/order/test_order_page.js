// 测试预约页面功能
const app = getApp()

// 测试就诊人API
function testPatientAPI() {
  const userId = wx.getStorageSync('userId')
  console.log('测试就诊人API，用户ID:', userId)
  
  wx.request({
    url: `${app.globalData.baseUrl}/api/user/patient`,
    method: 'GET',
    data: { userId },
    success: (res) => {
      console.log('就诊人API响应:', res.data)
      if (res.data.code === 0) {
        const patientList = res.data.data || []
        console.log('就诊人列表:', patientList)
        console.log('就诊人数量:', patientList.length)
        
        if (patientList.length > 0) {
          console.log('第一个就诊人:', patientList[0])
        }
      } else {
        console.error('就诊人API错误:', res.data)
      }
    },
    fail: (err) => {
      console.error('就诊人API请求失败:', err)
    }
  })
}

// 测试地址API
function testAddressAPI() {
  const userId = wx.getStorageSync('userId')
  console.log('测试地址API，用户ID:', userId)
  
  wx.request({
    url: `${app.globalData.baseUrl}/api/user/address`,
    method: 'GET',
    data: { userId },
    success: (res) => {
      console.log('地址API响应:', res.data)
      if (res.data.code === 0) {
        const addressList = res.data.data || []
        console.log('地址列表:', addressList)
        console.log('地址数量:', addressList.length)
        
        if (addressList.length > 0) {
          console.log('第一个地址:', addressList[0])
        }
      } else {
        console.error('地址API错误:', res.data)
      }
    },
    fail: (err) => {
      console.error('地址API请求失败:', err)
    }
  })
}

// 测试预约时间范围
function testAppointmentDateRange() {
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
  
  console.log('预约时间范围测试:')
  console.log('今天:', formatDate(today))
  console.log('明天:', minDate)
  console.log('7天后:', maxDateStr)
  console.log('时间范围:', `${minDate} 至 ${maxDateStr}`)
}

// 测试页面数据
function testPageData() {
  const pageData = {
    serviceId: '1',
    serviceInfo: {
      id: 1,
      name: '测试服务',
      price: 100,
      description: '测试服务描述'
    },
    patientList: [
      {
        id: 1,
        name: '张三',
        relation: '本人',
        phone: '13800138000',
        isDefault: 1
      }
    ],
    addressList: [
      {
        id: 1,
        name: '张三',
        phone: '13800138000',
        province: '广东省',
        city: '深圳市',
        district: '罗湖区',
        address: '东门北路1017号',
        isDefault: 1
      }
    ],
    selectedPatient: null,
    selectedAddress: null,
    formData: {
      appointmentDate: '',
      appointmentTime: '',
      remark: ''
    }
  }
  
  console.log('页面数据测试:', pageData)
  console.log('就诊人数量:', pageData.patientList.length)
  console.log('地址数量:', pageData.addressList.length)
}

// 运行所有测试
function runAllTests() {
  console.log('=== 开始测试预约页面功能 ===')
  
  testPatientAPI()
  testAddressAPI()
  testAppointmentDateRange()
  testPageData()
  
  console.log('=== 测试完成 ===')
}

// 导出测试函数
module.exports = {
  testPatientAPI,
  testAddressAPI,
  testAppointmentDateRange,
  testPageData,
  runAllTests
}

// 如果直接运行此文件
if (typeof wx !== 'undefined') {
  // 在微信小程序环境中
  setTimeout(() => {
    runAllTests()
  }, 1000)
} 