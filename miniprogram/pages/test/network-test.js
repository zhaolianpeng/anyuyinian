// 网络连接测试页面
const { api } = require('../../utils/cloud-container-standard')

Page({
  data: {
    testResults: [],
    isTesting: false,
    networkStatus: null
  },

  onLoad() {
    this.checkNetwork()
  },

  // 检查网络状态
  async checkNetwork() {
    try {
      const status = await this.getNetworkStatus()
      this.setData({
        networkStatus: status
      })
      console.log('网络状态:', status)
    } catch (error) {
      console.error('检查网络状态失败:', error)
    }
  },

  // 获取网络状态
  getNetworkStatus() {
    return new Promise((resolve) => {
      wx.getNetworkType({
        success: (res) => {
          console.log('网络类型:', res.networkType)
          resolve({
            isConnected: res.networkType !== 'none',
            networkType: res.networkType
          })
        },
        fail: () => {
          console.error('获取网络状态失败')
          resolve({
            isConnected: false,
            networkType: 'unknown'
          })
        }
      })
    })
  },

  // 测试服务端连接
  async testCloudContainer() {
    this.setData({
      isTesting: true,
      testResults: []
    })

    const results = []

    try {
      // 测试1: 基础连接测试
      console.log('开始测试服务端基础连接...')
      results.push({
        name: '基础连接测试',
        status: 'testing',
        message: '正在测试...'
      })
      this.setData({ testResults: results })

      const app = getApp()
      const testResult = await app.callContainer('/api/count', 'GET')

      results[0] = {
        name: '基础连接测试',
        status: 'success',
        message: '连接成功',
        data: testResult
      }

    } catch (error) {
      results[0] = {
        name: '基础连接测试',
        status: 'error',
        message: `连接失败: ${error.message}`,
        error: error
      }
    }

    // 测试2: API调用测试
    try {
      console.log('开始测试API调用...')
      results.push({
        name: 'API调用测试',
        status: 'testing',
        message: '正在测试...'
      })
      this.setData({ testResults: results })

      const homeData = await api.homeInit({
        longitude: 121.4737,
        latitude: 31.2304,
        limit: 5
      })
      
      results[1] = {
        name: 'API调用测试',
        status: 'success',
        message: 'API调用成功',
        data: homeData
      }

    } catch (error) {
      results[1] = {
        name: 'API调用测试',
        status: 'error',
        message: `API调用失败: ${error.message}`,
        error: error
      }
    }

    // 测试3: 计数器API测试
    try {
      console.log('开始测试计数器API...')
      results.push({
        name: '计数器API测试',
        status: 'testing',
        message: '正在测试...'
      })
      this.setData({ testResults: results })

      const countResult = await api.count.get()
      
      results[2] = {
        name: '计数器API测试',
        status: 'success',
        message: '计数器API调用成功',
        data: countResult
      }

    } catch (error) {
      results[2] = {
        name: '计数器API测试',
        status: 'error',
        message: `计数器API调用失败: ${error.message}`,
        error: error
      }
    }

    // 测试4: 错误处理测试
    try {
      console.log('开始测试错误处理...')
      results.push({
        name: '错误处理测试',
        status: 'testing',
        message: '正在测试...'
      })
      this.setData({ testResults: results })

      // 故意调用一个不存在的API
      await api.orderDetail('INVALID_ORDER_NO')
      
      results[3] = {
        name: '错误处理测试',
        status: 'error',
        message: '错误处理测试失败 - 应该返回错误',
        error: '未按预期返回错误'
      }

    } catch (error) {
      results[3] = {
        name: '错误处理测试',
        status: 'success',
        message: '错误处理正常',
        data: error.message
      }
    }

    this.setData({
      testResults: results,
      isTesting: false
    })

    console.log('网络测试完成:', results)
  },

  // 重新测试
  retest() {
    this.testCloudContainer()
  },

  // 查看错误详情
  viewError(e) {
    const { index } = e.currentTarget.dataset
    const result = this.data.testResults[index]
    
    if (result.error) {
      wx.showModal({
        title: '错误详情',
        content: JSON.stringify(result.error, null, 2),
        showCancel: false
      })
    }
  },

  // 复制测试结果
  copyResults() {
    const results = this.data.testResults.map(r => 
      `${r.name}: ${r.status === 'success' ? '✅' : '❌'} ${r.message}`
    ).join('\n')
    
    wx.setClipboardData({
      data: results,
      success: () => {
        wx.showToast({
          title: '结果已复制',
          icon: 'success'
        })
      }
    })
  }
}) 