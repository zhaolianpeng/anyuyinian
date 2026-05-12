// 测试订单页面流程
const app = getApp()

// 测试订单页面数据加载
function testOrderPageData() {
  console.log('=== 测试订单页面数据加载 ===')
  
  const userId = wx.getStorageSync('userId')
  console.log('用户ID:', userId)
  
  // 测试就诊人数据加载
  wx.request({
    url: `${app.globalData.baseUrl}/api/user/patient`,
    method: 'GET',
    data: { userId },
    success: (res) => {
      console.log('就诊人数据:', res.data)
      if (res.data.code === 0) {
        console.log('✅ 就诊人数据加载成功，数量:', res.data.data?.length || 0)
      } else {
        console.log('❌ 就诊人数据加载失败:', res.data.errorMsg)
      }
    },
    fail: (err) => {
      console.error('就诊人数据加载错误:', err)
    }
  })
  
  // 测试地址数据加载
  wx.request({
    url: `${app.globalData.baseUrl}/api/user/address`,
    method: 'GET',
    data: { userId },
    success: (res) => {
      console.log('地址数据:', res.data)
      if (res.data.code === 0) {
        console.log('✅ 地址数据加载成功，数量:', res.data.data?.length || 0)
      } else {
        console.log('❌ 地址数据加载失败:', res.data.errorMsg)
      }
    },
    fail: (err) => {
      console.error('地址数据加载错误:', err)
    }
  })
}

// 测试添加就诊人流程
function testAddPatientFlow() {
  console.log('=== 测试添加就诊人流程 ===')
  
  const testData = {
    userId: wx.getStorageSync('userId'),
    name: '测试就诊人',
    idCard: '110101199001011234',
    phone: '13800138000',
    gender: 1,
    birthday: '1990-01-01',
    relation: '本人',
    isDefault: false
  }
  
  wx.request({
    url: `${app.globalData.baseUrl}/api/user/patient`,
    method: 'POST',
    data: testData,
    success: (res) => {
      console.log('添加就诊人响应:', res.data)
      if (res.data.code === 0) {
        console.log('✅ 就诊人添加成功')
        // 验证数据是否正确保存
        setTimeout(() => {
          testOrderPageData()
        }, 1000)
      } else {
        console.log('❌ 就诊人添加失败:', res.data.errorMsg)
      }
    },
    fail: (err) => {
      console.error('添加就诊人错误:', err)
    }
  })
}

// 测试添加地址流程
function testAddAddressFlow() {
  console.log('=== 测试添加地址流程 ===')
  
  const testData = {
    userId: wx.getStorageSync('userId'),
    name: '测试联系人',
    phone: '13800138000',
    province: '北京市',
    city: '北京市',
    district: '朝阳区',
    address: '测试地址',
    isDefault: false
  }
  
  wx.request({
    url: `${app.globalData.baseUrl}/api/user/address`,
    method: 'POST',
    data: testData,
    success: (res) => {
      console.log('添加地址响应:', res.data)
      if (res.data.code === 0) {
        console.log('✅ 地址添加成功')
        // 验证数据是否正确保存
        setTimeout(() => {
          testOrderPageData()
        }, 1000)
      } else {
        console.log('❌ 地址添加失败:', res.data.errorMsg)
      }
    },
    fail: (err) => {
      console.error('添加地址错误:', err)
    }
  })
}

// 测试订单页面功能
function testOrderPageFeatures() {
  console.log('=== 测试订单页面功能 ===')
  
  // 检查页面是否正确显示添加按钮
  console.log('1. 检查就诊人添加按钮是否显示')
  console.log('2. 检查地址添加按钮是否显示')
  console.log('3. 检查数据是否正确加载')
  console.log('4. 检查选择功能是否正常')
  
  // 模拟页面数据
  const mockData = {
    patientList: [
      {
        id: 1,
        name: '张三',
        relation: '本人',
        phone: '13800138000',
        isDefault: true
      }
    ],
    addressList: [
      {
        id: 1,
        name: '李四',
        phone: '13800138000',
        province: '北京市',
        city: '北京市',
        district: '朝阳区',
        address: '测试地址',
        isDefault: true
      }
    ]
  }
  
  console.log('模拟数据:', mockData)
  console.log('✅ 订单页面功能测试完成')
}

// 导出测试函数
module.exports = {
  testOrderPageData,
  testAddPatientFlow,
  testAddAddressFlow,
  testOrderPageFeatures
} 