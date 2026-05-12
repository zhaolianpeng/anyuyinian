/**
 * API调用示例
 * 展示如何使用统一的服务端调用方式
 */

const { api } = require('./cloud-container-standard')

/**
 * 示例1: 计数器相关API
 */
const counterExample = {
  // 获取计数器
  async getCount() {
    try {
      const result = await api.count.get()
      console.log('计数器值:', result)
      return result
    } catch (error) {
      console.error('获取计数器失败:', error)
      throw error
    }
  },

  // 增加计数器
  async incrementCount() {
    try {
      const result = await api.count.increment()
      console.log('计数器增加成功:', result)
      return result
    } catch (error) {
      console.error('增加计数器失败:', error)
      throw error
    }
  },

  // 减少计数器
  async decrementCount() {
    try {
      const result = await api.count.decrement()
      console.log('计数器减少成功:', result)
      return result
    } catch (error) {
      console.error('减少计数器失败:', error)
      throw error
    }
  }
}

/**
 * 示例2: 用户相关API
 */
const userExample = {
  // 微信登录
  async wxLogin(code) {
    try {
      const result = await api.wxLogin({ code })
      console.log('微信登录成功:', result)
      return result
    } catch (error) {
      console.error('微信登录失败:', error)
      throw error
    }
  },

  // 获取用户信息
  async getUserInfo() {
    try {
      const result = await api.userInfo()
      console.log('用户信息:', result)
      return result
    } catch (error) {
      console.error('获取用户信息失败:', error)
      throw error
    }
  },

  // 绑定手机号
  async bindPhone(phoneData) {
    try {
      const result = await api.bindPhone(phoneData)
      console.log('绑定手机号成功:', result)
      return result
    } catch (error) {
      console.error('绑定手机号失败:', error)
      throw error
    }
  }
}

/**
 * 示例3: 首页相关API
 */
const homeExample = {
  // 首页初始化
  async initHome(location) {
    try {
      const result = await api.homeInit({
        longitude: location.longitude,
        latitude: location.latitude,
        limit: 10
      })
      console.log('首页数据:', result)
      return result
    } catch (error) {
      console.error('获取首页数据失败:', error)
      throw error
    }
  }
}

/**
 * 示例4: 服务相关API
 */
const serviceExample = {
  // 获取服务列表
  async getServiceList(params = {}) {
    try {
      const result = await api.serviceList(params)
      console.log('服务列表:', result)
      return result
    } catch (error) {
      console.error('获取服务列表失败:', error)
      throw error
    }
  },

  // 获取服务详情
  async getServiceDetail(serviceId) {
    try {
      const result = await api.serviceDetail({ id: serviceId })
      console.log('服务详情:', result)
      return result
    } catch (error) {
      console.error('获取服务详情失败:', error)
      throw error
    }
  }
}

/**
 * 示例5: 订单相关API
 */
const orderExample = {
  // 提交订单
  async submitOrder(orderData) {
    try {
      const result = await api.orderSubmit(orderData)
      console.log('订单提交成功:', result)
      return result
    } catch (error) {
      console.error('订单提交失败:', error)
      throw error
    }
  },

  // 获取订单列表
  async getOrderList(params = {}) {
    try {
      const result = await api.orderList(params)
      console.log('订单列表:', result)
      return result
    } catch (error) {
      console.error('获取订单列表失败:', error)
      throw error
    }
  },

  // 获取订单详情
  async getOrderDetail(orderNo) {
    try {
      const result = await api.orderDetail(orderNo)
      console.log('订单详情:', result)
      return result
    } catch (error) {
      console.error('获取订单详情失败:', error)
      throw error
    }
  }
}

/**
 * 示例6: 直接使用callContainer方法
 */
const directCallExample = {
  // 直接调用统一服务端 API
  async callContainer(path, method = 'GET', data = {}, options = {}) {
    const { callContainer } = require('./cloud-container-standard')
    
    try {
      const result = await callContainer(path, method, data, options)
      console.log('直接调用成功:', result)
      return result
    } catch (error) {
      console.error('直接调用失败:', error)
      throw error
    }
  },

  // 自定义API调用示例
  async customApiCall() {
    try {
      // 示例：调用自定义接口
      const result = await this.callContainer('/api/custom/endpoint', 'POST', {
        param1: 'value1',
        param2: 'value2'
      }, {
        header: {
          'X-Custom-Header': 'custom-value'
        }
      })
      
      console.log('自定义API调用成功:', result)
      return result
    } catch (error) {
      console.error('自定义API调用失败:', error)
      throw error
    }
  }
}

module.exports = {
  counterExample,
  userExample,
  homeExample,
  serviceExample,
  orderExample,
  directCallExample
} 