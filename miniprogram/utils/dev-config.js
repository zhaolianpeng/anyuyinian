// 开发环境配置文件
// 用于管理微信开发工具中的各种设置

const devConfig = {
  // 是否启用开发模式
  enableDevMode: true,
  
  // 开发环境下的默认城市
  defaultCity: '北京',
  
  // 是否跳过定位功能
  skipLocation: true,
  
  // 是否模拟定位成功
  simulateLocationSuccess: false,
  
  // 定位超时时间（毫秒）
  locationTimeout: 10000,
  
  // 是否显示调试按钮
  showDebugButton: true,
  
  // 开发环境检测
  isDevelopmentEnvironment() {
    if (typeof wx !== 'undefined' && wx.getSystemInfoSync) {
      try {
        const systemInfo = wx.getSystemInfoSync()
        return systemInfo.platform === 'devtools' || 
               systemInfo.version.includes('devtools') ||
               systemInfo.brand === 'devtools'
      } catch (e) {
        console.log('获取系统信息失败:', e)
        return false
      }
    }
    return false
  },
  
  // 获取开发环境配置
  getConfig() {
    if (this.isDevelopmentEnvironment()) {
      return {
        ...this,
        isDev: true
      }
    }
    return {
      ...this,
      isDev: false,
      skipLocation: false,
      showDebugButton: false
    }
  },
  
  // 开发环境日志
  log(message, data = null) {
    if (this.isDevelopmentEnvironment()) {
      console.log(`[DEV] ${message}`, data)
    }
  },
  
  // 开发环境警告
  warn(message, data = null) {
    if (this.isDevelopmentEnvironment()) {
      console.warn(`[DEV] ${message}`, data)
    }
  },
  
  // 开发环境错误
  error(message, data = null) {
    if (this.isDevelopmentEnvironment()) {
      console.error(`[DEV] ${message}`, data)
    }
  }
}

module.exports = devConfig
