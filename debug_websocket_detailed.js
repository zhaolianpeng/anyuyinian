/**
 * 详细WebSocket诊断工具
 * 提供更深入的WebSocket连接问题诊断
 */

const { websocket } = require('./config')

class DetailedWebSocketDebugger {
  constructor() {
    this.debugLog = []
    this.connectionAttempts = 0
    this.maxAttempts = 3
  }

  /**
   * 开始详细诊断
   */
  async startDetailedDebug() {
    console.log('=== 详细WebSocket诊断开始 ===')
    
    // 1. 环境检查
    await this.checkEnvironment()
    
    // 2. 配置检查
    await this.checkConfiguration()
    
    // 3. 网络检查
    await this.checkNetwork()
    
    // 4. WebSocket连接测试
    await this.testWebSocketConnection()
    
    // 5. 输出诊断报告
    this.generateDiagnosticReport()
  }

  /**
   * 检查运行环境
   */
  async checkEnvironment() {
    console.log('1. 检查运行环境...')
    
    const envChecks = {
      wx: typeof wx !== 'undefined',
      wxCloud: typeof wx !== 'undefined' && !!wx.cloud,
      connectContainer: typeof wx !== 'undefined' && !!wx.cloud && !!wx.cloud.connectContainer,
      config: !!websocket,
      env: !!websocket.env,
      service: !!websocket.service,
      path: !!websocket.path
    }
    
    console.log('环境检查结果:', envChecks)
    
    if (!envChecks.wx) {
      this.debugLog.push('❌ 错误: 不在微信小程序环境中')
      return false
    }
    
    if (!envChecks.wxCloud) {
      this.debugLog.push('❌ 错误: wx.cloud不存在')
      return false
    }
    
    if (!envChecks.connectContainer) {
      this.debugLog.push('❌ 错误: wx.cloud.connectContainer不存在')
      return false
    }
    
    if (!envChecks.config) {
      this.debugLog.push('❌ 错误: WebSocket配置不存在')
      return false
    }
    
    this.debugLog.push('✅ 环境检查通过')
    return true
  }

  /**
   * 检查配置
   */
  async checkConfiguration() {
    console.log('2. 检查WebSocket配置...')
    
    const config = {
      env: websocket.env,
      service: websocket.service,
      path: websocket.path,
      reconnect: websocket.reconnect
    }
    
    console.log('WebSocket配置:', config)
    
    if (!config.env) {
      this.debugLog.push('❌ 错误: 环境ID未配置')
      return false
    }
    
    if (!config.service) {
      this.debugLog.push('❌ 错误: 服务名未配置')
      return false
    }
    
    if (!config.path) {
      this.debugLog.push('❌ 错误: WebSocket路径未配置')
      return false
    }
    
    this.debugLog.push('✅ 配置检查通过')
    this.debugLog.push(`配置详情: env=${config.env}, service=${config.service}, path=${config.path}`)
    return true
  }

  /**
   * 检查网络连接
   */
  async checkNetwork() {
    console.log('3. 检查网络连接...')
    
    try {
      // 测试基本网络连接
      const networkTest = await this.testBasicNetwork()
      if (networkTest) {
        this.debugLog.push('✅ 基本网络连接正常')
      } else {
        this.debugLog.push('❌ 基本网络连接失败')
        return false
      }
      
      // 测试云托管连接
      const containerTest = await this.testContainerConnection()
      if (containerTest) {
        this.debugLog.push('✅ 云托管连接正常')
      } else {
        this.debugLog.push('❌ 云托管连接失败')
        return false
      }
      
      return true
    } catch (error) {
      this.debugLog.push(`❌ 网络检查失败: ${error.message}`)
      return false
    }
  }

  /**
   * 测试基本网络连接
   */
  async testBasicNetwork() {
    return new Promise((resolve) => {
      wx.request({
        url: 'https://www.baidu.com',
        method: 'GET',
        timeout: 5000,
        success: () => {
          resolve(true)
        },
        fail: () => {
          resolve(false)
        }
      })
    })
  }

  /**
   * 测试云托管连接
   */
  async testContainerConnection() {
    return new Promise((resolve) => {
      // 尝试调用一个简单的API来测试云托管连接
      wx.cloud.callContainer({
        config: {
          env: websocket.env,
        },
        service: websocket.service,
        path: '/',
        method: 'GET',
        timeout: 5000,
        success: () => {
          resolve(true)
        },
        fail: () => {
          resolve(false)
        }
      })
    })
  }

  /**
   * 测试WebSocket连接
   */
  async testWebSocketConnection() {
    console.log('4. 测试WebSocket连接...')
    
    for (let attempt = 1; attempt <= this.maxAttempts; attempt++) {
      this.debugLog.push(`尝试连接 #${attempt}...`)
      
      try {
        const result = await this.attemptWebSocketConnection(attempt)
        if (result.success) {
          this.debugLog.push(`✅ WebSocket连接成功 (尝试 #${attempt})`)
          return true
        } else {
          this.debugLog.push(`❌ WebSocket连接失败 (尝试 #${attempt}): ${result.error}`)
        }
      } catch (error) {
        this.debugLog.push(`❌ WebSocket连接异常 (尝试 #${attempt}): ${error.message}`)
      }
      
      // 等待一段时间再重试
      if (attempt < this.maxAttempts) {
        await this.sleep(2000)
      }
    }
    
    this.debugLog.push('❌ 所有WebSocket连接尝试都失败了')
    return false
  }

  /**
   * 尝试WebSocket连接
   */
  async attemptWebSocketConnection(attempt) {
    return new Promise((resolve) => {
      const startTime = Date.now()
      
      wx.cloud.connectContainer({
        config: {
          env: websocket.env,
        },
        service: websocket.service,
        path: websocket.path,
        timeout: 10000
      }).then(({ socketTask }) => {
        const duration = Date.now() - startTime
        
        // 设置事件监听
        let connected = false
        let error = null
        
        socketTask.onOpen(() => {
          connected = true
          this.debugLog.push(`WebSocket连接已打开 (耗时: ${duration}ms)`)
        })
        
        socketTask.onError((err) => {
          error = err
          this.debugLog.push(`WebSocket连接错误: ${JSON.stringify(err)}`)
        })
        
        socketTask.onClose((res) => {
          this.debugLog.push(`WebSocket连接已关闭: ${JSON.stringify(res)}`)
        })
        
        // 等待一段时间看是否成功连接
        setTimeout(() => {
          if (connected) {
            socketTask.close()
            resolve({ success: true, duration })
          } else {
            socketTask.close()
            resolve({ success: false, error: error || '连接超时' })
          }
        }, 3000)
        
      }).catch((error) => {
        const duration = Date.now() - startTime
        this.debugLog.push(`WebSocket连接失败: ${error.message} (耗时: ${duration}ms)`)
        resolve({ success: false, error: error.message })
      })
    })
  }

  /**
   * 生成诊断报告
   */
  generateDiagnosticReport() {
    console.log('=== WebSocket诊断报告 ===')
    
    const report = {
      timestamp: new Date().toISOString(),
      config: websocket,
      debugLog: this.debugLog,
      summary: this.generateSummary()
    }
    
    console.log('诊断报告:', report)
    
    // 显示总结
    console.log('=== 诊断总结 ===')
    console.log(report.summary)
    
    return report
  }

  /**
   * 生成总结
   */
  generateSummary() {
    const errorCount = this.debugLog.filter(log => log.includes('❌')).length
    const successCount = this.debugLog.filter(log => log.includes('✅')).length
    
    let summary = `诊断完成: ${successCount} 项通过, ${errorCount} 项失败\n\n`
    
    if (errorCount === 0) {
      summary += '🎉 所有检查都通过了！WebSocket应该可以正常工作。'
    } else {
      summary += '⚠️ 发现了一些问题，请根据上面的错误信息进行修复。\n\n'
      
      // 提供具体的修复建议
      if (this.debugLog.some(log => log.includes('wx.cloud.connectContainer不存在'))) {
        summary += '建议: 检查微信小程序基础库版本，确保支持云托管功能\n'
      }
      
      if (this.debugLog.some(log => log.includes('云托管连接失败'))) {
        summary += '建议: 检查云托管服务是否正常运行，确认环境ID和服务名正确\n'
      }
      
      if (this.debugLog.some(log => log.includes('WebSocket连接失败'))) {
        summary += '建议: 检查后端WebSocket服务是否正常启动，确认/ws路由已注册\n'
      }
    }
    
    return summary
  }

  /**
   * 延时函数
   */
  sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms))
  }

  /**
   * 获取诊断信息
   */
  getDiagnosticInfo() {
    return {
      config: websocket,
      debugLog: this.debugLog,
      timestamp: new Date().toISOString()
    }
  }
}

// 导出调试器
module.exports = DetailedWebSocketDebugger 