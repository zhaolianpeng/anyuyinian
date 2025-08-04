/**
 * WebSocket调试工具
 * 用于诊断WebSocket连接问题
 */

const { websocket } = require('./config')

class WebSocketDebugger {
  constructor() {
    this.connectionAttempts = 0
    this.maxAttempts = 3
    this.debugLog = []
  }

  /**
   * 开始调试WebSocket连接
   */
  async startDebug() {
    console.log('=== WebSocket调试开始 ===')
    
    // 1. 检查配置
    this.checkConfig()
    
    // 2. 检查环境
    this.checkEnvironment()
    
    // 3. 尝试连接
    await this.testConnection()
    
    // 4. 输出调试信息
    this.printDebugInfo()
  }

  /**
   * 检查WebSocket配置
   */
  checkConfig() {
    console.log('1. 检查WebSocket配置...')
    
    const config = {
      env: websocket.env,
      service: websocket.service,
      path: websocket.path,
      reconnect: websocket.reconnect
    }
    
    console.log('WebSocket配置:', config)
    
    if (!config.env || !config.service) {
      console.error('❌ WebSocket配置不完整')
      this.debugLog.push('配置错误: env或service未设置')
    } else {
      console.log('✅ WebSocket配置正确')
    }
  }

  /**
   * 检查运行环境
   */
  checkEnvironment() {
    console.log('2. 检查运行环境...')
    
    if (typeof wx === 'undefined') {
      console.error('❌ 不在微信小程序环境中')
      this.debugLog.push('环境错误: wx对象不存在')
      return
    }
    
    if (!wx.cloud) {
      console.error('❌ wx.cloud不存在')
      this.debugLog.push('环境错误: wx.cloud不存在')
      return
    }
    
    if (!wx.cloud.connectContainer) {
      console.error('❌ wx.cloud.connectContainer不存在')
      this.debugLog.push('环境错误: wx.cloud.connectContainer不存在')
      return
    }
    
    console.log('✅ 运行环境正常')
  }

  /**
   * 测试WebSocket连接
   */
  async testConnection() {
    console.log('3. 测试WebSocket连接...')
    
    try {
      const { socketTask } = await wx.cloud.connectContainer({
        config: {
          env: websocket.env,
        },
        service: websocket.service,
        path: websocket.path
      })
      
      console.log('✅ WebSocket连接成功')
      this.debugLog.push('连接成功: WebSocket已建立')
      
      // 设置事件监听
      this.setupEventListeners(socketTask)
      
    } catch (error) {
      console.error('❌ WebSocket连接失败:', error)
      this.debugLog.push(`连接失败: ${error.message}`)
      
      // 分析错误原因
      this.analyzeError(error)
    }
  }

  /**
   * 设置事件监听
   */
  setupEventListeners(socketTask) {
    socketTask.onOpen((res) => {
      console.log('WebSocket连接已打开:', res)
      this.debugLog.push('事件: onOpen触发')
    })

    socketTask.onMessage((res) => {
      console.log('收到WebSocket消息:', res.data)
      this.debugLog.push(`事件: onMessage - ${JSON.stringify(res.data)}`)
    })

    socketTask.onClose((res) => {
      console.log('WebSocket连接已关闭:', res)
      this.debugLog.push(`事件: onClose - ${JSON.stringify(res)}`)
    })

    socketTask.onError((res) => {
      console.error('WebSocket连接错误:', res)
      this.debugLog.push(`事件: onError - ${JSON.stringify(res)}`)
    })
  }

  /**
   * 分析错误原因
   */
  analyzeError(error) {
    const errorMessage = error.message || error.toString()
    
    if (errorMessage.includes('timeout')) {
      console.log('可能原因: 连接超时')
      this.debugLog.push('分析: 连接超时')
    } else if (errorMessage.includes('fail')) {
      console.log('可能原因: 网络请求失败')
      this.debugLog.push('分析: 网络请求失败')
    } else if (errorMessage.includes('service')) {
      console.log('可能原因: 服务不存在或未启动')
      this.debugLog.push('分析: 服务不存在或未启动')
    } else if (errorMessage.includes('env')) {
      console.log('可能原因: 环境ID错误')
      this.debugLog.push('分析: 环境ID错误')
    } else {
      console.log('可能原因: 未知错误')
      this.debugLog.push('分析: 未知错误')
    }
  }

  /**
   * 输出调试信息
   */
  printDebugInfo() {
    console.log('=== 调试信息汇总 ===')
    console.log('调试日志:')
    this.debugLog.forEach((log, index) => {
      console.log(`${index + 1}. ${log}`)
    })
    
    console.log('=== 调试完成 ===')
  }

  /**
   * 获取调试信息
   */
  getDebugInfo() {
    return {
      config: websocket,
      debugLog: this.debugLog,
      timestamp: new Date().toISOString()
    }
  }
}

// 导出调试器
module.exports = WebSocketDebugger 