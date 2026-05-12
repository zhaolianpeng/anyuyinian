// 测试数据加载和显示
const app = getApp()

// 测试地址数据加载
function testAddressData() {
  const userId = wx.getStorageSync('userId')
  console.log('=== 测试地址数据加载 ===')
  console.log('用户ID:', userId)
  
  wx.request({
    url: `${app.globalData.baseUrl}/api/user/address`,
    method: 'GET',
    data: { userId },
    success: (res) => {
      console.log('地址API完整响应:', res)
      console.log('地址API数据:', res.data)
      
      if (res.data.code === 0) {
        const addressList = res.data.data || []
        console.log('地址列表:', addressList)
        console.log('地址数量:', addressList.length)
        
        if (addressList.length > 0) {
          console.log('第一个地址:', addressList[0])
          console.log('isDefault值:', addressList[0].isDefault)
          console.log('isDefault类型:', typeof addressList[0].isDefault)
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

// 测试就诊人数据加载
function testPatientData() {
  const userId = wx.getStorageSync('userId')
  console.log('=== 测试就诊人数据加载 ===')
  console.log('用户ID:', userId)
  
  wx.request({
    url: `${app.globalData.baseUrl}/api/user/patient`,
    method: 'GET',
    data: { userId },
    success: (res) => {
      console.log('就诊人API完整响应:', res)
      console.log('就诊人API数据:', res.data)
      
      if (res.data.code === 0) {
        const patientList = res.data.data || []
        console.log('就诊人列表:', patientList)
        console.log('就诊人数量:', patientList.length)
        
        if (patientList.length > 0) {
          console.log('第一个就诊人:', patientList[0])
          console.log('isDefault值:', patientList[0].isDefault)
          console.log('isDefault类型:', typeof patientList[0].isDefault)
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

// 测试页面数据状态
function testPageDataState() {
  console.log('=== 测试页面数据状态 ===')
  
  // 获取当前页面实例
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1]
  
  if (currentPage && currentPage.data) {
    console.log('当前页面数据:', currentPage.data)
    console.log('就诊人列表:', currentPage.data.patientList)
    console.log('地址列表:', currentPage.data.addressList)
    console.log('选中的就诊人:', currentPage.data.selectedPatient)
    console.log('选中的地址:', currentPage.data.selectedAddress)
  } else {
    console.log('无法获取页面数据')
  }
}

// 模拟数据设置
function simulateDataSetting() {
  console.log('=== 模拟数据设置 ===')
  
  const mockAddressData = {
    code: 0,
    data: [
      {
        id: 1,
        userId: 1,
        name: "赵连鹏",
        phone: "13691028481",
        province: "北京市",
        city: "北京市",
        district: "西城区",
        address: "复兴门北大街9号楼",
        isDefault: 0,
        status: 1,
        createdAt: "2025-08-01T03:16:21Z",
        updatedAt: "2025-08-01T03:16:21Z"
      }
    ]
  }
  
  const mockPatientData = {
    code: 0,
    data: [
      {
        id: 1,
        userId: 1,
        name: "测试就诊人",
        relation: "本人",
        phone: "13800138000",
        idCard: "110101199001011234",
        gender: 1,
        birthday: "1990-01-01",
        isDefault: 0,
        status: 1
      }
    ]
  }
  
  console.log('模拟地址数据:', mockAddressData)
  console.log('模拟就诊人数据:', mockPatientData)
  
  // 模拟数据设置逻辑
  if (mockAddressData.code === 0) {
    const addressList = mockAddressData.data || []
    console.log('设置地址列表:', addressList)
    console.log('默认地址:', addressList.find(a => a.isDefault === 1))
    console.log('第一个地址:', addressList[0])
  }
  
  if (mockPatientData.code === 0) {
    const patientList = mockPatientData.data || []
    console.log('设置就诊人列表:', patientList)
    console.log('默认就诊人:', patientList.find(p => p.isDefault === 1))
    console.log('第一个就诊人:', patientList[0])
  }
}

// 导出测试函数
module.exports = {
  testAddressData,
  testPatientData,
  testPageDataState,
  simulateDataSetting
} 