const { consultationAPI, api } = require('../../utils/request.js')

Page({
  data: {
    consultationId: '', // 咨询ID
    consultationStatus: 'waiting', // 咨询状态：waiting-等待回复, chatting-咨询中, closed-已结束
    messages: [], // 聊天消息
    inputMessage: '', // 输入框内容
    scrollToMessage: '', // 滚动到指定消息
    userInfo: {}, // 用户信息
    quickQuestions: [
      '护工服务价格是多少？',
      '如何预约护工服务？',
      '护工服务范围包括哪些？',
      '服务时间如何安排？',
      '如何评价护工服务？'
    ],
    websocket: null, // WebSocket连接
    reconnectTimer: null, // 重连定时器
    heartbeatTimer: null, // 心跳定时器
    isConnected: false // 连接状态
  },

  onLoad(options) {
    // 获取用户信息
    this.getUserInfo()
    
    // 创建或获取咨询会话
    this.initConsultation()
    
    // 建立WebSocket连接
    this.connectWebSocket()
  },

  onShow() {
    // 页面显示时刷新消息（如果有咨询会话ID）
    if (this.data.consultationId) {
      this.loadMessages()
    }
  },

  onHide() {
    // 页面隐藏时暂停WebSocket
    this.pauseWebSocket()
  },

  onUnload() {
    // 页面卸载时清理资源
    this.cleanup()
  },

  // 检查用户登录状态
  async checkUserLogin() {
    try {
      // 首先尝试从存储中获取
      let userInfo = wx.getStorageSync('userInfo')
      
      // 如果没有存储的用户信息，或者信息不完整，则重新获取
      if (!userInfo || !userInfo.userId) {
        console.log('从存储获取的用户信息不完整，重新获取...')
        const res = await api.userInfo()
        if (res.code === 0 && res.data) {
          userInfo = res.data
          this.setData({ userInfo })
          wx.setStorageSync('userInfo', userInfo)
          console.log('成功获取用户信息:', userInfo)
          
          // 用户信息获取成功，继续初始化
          this.initConsultation()
        } else {
          throw new Error('获取用户信息失败: ' + (res.errorMsg || '未知错误'))
        }
      } else {
        console.log('使用存储的用户信息:', userInfo)
        this.setData({ userInfo })
        
        // 用户信息存在，继续初始化
        this.initConsultation()
      }
    } catch (error) {
      console.error('获取用户信息失败:', error)
      wx.showModal({
        title: '登录提示',
        content: '您还没有登录，请先登录后再使用咨询功能',
        showCancel: false,
        confirmText: '去登录',
        success: () => {
          wx.navigateTo({
            url: '/pages/login/login'
          })
        }
      })
    }
  },

  // 获取用户信息
  async getUserInfo() {
    try {
      // 首先尝试从存储中获取
      let userInfo = wx.getStorageSync('userInfo')
      
      // 如果没有存储的用户信息，或者信息不完整，则重新获取
      if (!userInfo || !userInfo.userId) {
        console.log('从存储获取的用户信息不完整，重新获取...')
        const res = await api.userInfo()
        if (res.code === 0 && res.data) {
          userInfo = res.data
          this.setData({ userInfo })
          wx.setStorageSync('userInfo', userInfo)
          console.log('成功获取用户信息:', userInfo)
        } else {
          throw new Error('获取用户信息失败: ' + (res.errorMsg || '未知错误'))
        }
      } else {
        console.log('使用存储的用户信息:', userInfo)
        this.setData({ userInfo })
      }
    } catch (error) {
      console.error('获取用户信息失败:', error)
      wx.showToast({
        title: '获取用户信息失败，请重新登录',
        icon: 'none'
      })
    }
  },

  // 初始化咨询会话
  async initConsultation() {
    try {
      const consultationId = wx.getStorageSync('currentConsultationId')
      
      if (consultationId) {
        // 使用现有咨询会话
        console.log('使用现有咨询会话:', consultationId)
        this.setData({ consultationId })
        await this.loadMessages()
        await this.getConsultationStatus()
      } else {
        // 创建新的咨询会话
        console.log('创建新的咨询会话')
        await this.createConsultation()
      }
    } catch (error) {
      console.error('初始化咨询会话失败:', error)
      wx.showToast({
        title: '初始化失败，请重试',
        icon: 'none'
      })
    }
  },

  // 创建新的咨询会话
  async createConsultation() {
    try {
      // 确保有有效的用户信息
      if (!this.data.userInfo || !this.data.userInfo.userId) {
        throw new Error('用户信息不完整，请先登录')
      }

      // 添加调试日志
      console.log('准备创建咨询会话，用户信息:', {
        userId: this.data.userInfo.userId,
        userName: this.data.userInfo.nickName,
        userPhone: this.data.userInfo.phone
      })

      const res = await consultationAPI.createConsultation({
        userId: this.data.userInfo.userId,
        userName: this.data.userInfo.nickName || '用户',
        userPhone: this.data.userInfo.phone || ''
      })

      console.log('创建咨询会话响应:', res)

      if (res.code === 0 && res.data) {
        const consultationId = res.data.consultationId
        console.log('成功创建咨询会话，ID:', consultationId)
        
        this.setData({ consultationId })
        wx.setStorageSync('currentConsultationId', consultationId)
        
        // 发送欢迎消息
        await this.sendWelcomeMessage()
        
        wx.showToast({
          title: '咨询会话已创建',
          icon: 'success'
        })
      } else {
        throw new Error(res.errorMsg || '创建咨询会话失败')
      }
    } catch (error) {
      console.error('创建咨询会话失败:', error)
      wx.showToast({
        title: '创建咨询会话失败',
        icon: 'none'
      })
    }
  },

  // 发送欢迎消息
  async sendWelcomeMessage() {
    const now = new Date()
    const welcomeMessage = {
      id: Date.now().toString(),
      content: '您好！欢迎使用在线咨询服务，请问有什么可以帮助您的吗？',
      senderType: 'admin',
      createTime: this.formatTime(now),
      isWelcome: true
    }

    this.setData({
      messages: [welcomeMessage],
      consultationStatus: 'waiting'
    })

    // 滚动到底部
    this.scrollToBottom()
  },

  // 加载消息记录
  async loadMessages() {
    if (!this.data.consultationId) return

    try {
      const res = await consultationAPI.getConsultationMessages(this.data.consultationId)

      if (res.code === 0 && res.data) {
        const messages = res.data.messages.map(msg => {
          try {
            // 后端返回的时间字段是 createdAt
            const timeField = msg.createdAt
            console.log('处理消息时间，原始时间字段:', timeField, '消息对象:', msg)
            
            // 确保时间字符串能正确转换为Date对象
            const msgDate = new Date(timeField)
            return {
              ...msg,
              createTime: this.formatTime(msgDate)
            }
          } catch (error) {
            console.error('处理消息时间失败:', error, msg)
            return {
              ...msg,
              createTime: '时间未知'
            }
          }
        })

        this.setData({ messages })
        this.scrollToBottom()
      }
    } catch (error) {
      console.error('加载消息记录失败:', error)
    }
  },

  // 获取咨询状态
  async getConsultationStatus() {
    if (!this.data.consultationId) return

    try {
      const res = await consultationAPI.getConsultationStatus(this.data.consultationId)

      if (res.code === 0 && res.data) {
        this.setData({ consultationStatus: res.data.status })
      }
    } catch (error) {
      console.error('获取咨询状态失败:', error)
    }
  },

  // 建立WebSocket连接
  connectWebSocket() {
    // 暂时禁用WebSocket功能，避免连接错误
    console.log('WebSocket功能暂时禁用，使用轮询方式获取消息')
    this.setData({ isConnected: false })
    
    // 设置定时器定期刷新消息
    this.startMessagePolling()
  },

  // 处理WebSocket消息
  handleWebSocketMessage(data) {
    try {
      const message = JSON.parse(data)
      
      if (message.type === 'new_message') {
        // 收到新消息
        const timeField = message.data.createdAt
        const newMessage = {
          id: message.data.id,
          content: message.data.content,
          senderType: message.data.senderType,
          createTime: this.formatTime(new Date(timeField))
        }

        this.setData({
          messages: [...this.data.messages, newMessage]
        })

        this.scrollToBottom()

        // 如果是管理员回复，更新状态
        if (message.data.senderType === 'admin') {
          this.setData({ consultationStatus: 'chatting' })
        }

        // 显示通知
        if (message.data.senderType === 'admin') {
          wx.showToast({
            title: '收到回复',
            icon: 'success'
          })
        }
      } else if (message.type === 'status_update') {
        // 状态更新
        this.setData({ consultationStatus: message.data.status })
      }
    } catch (error) {
      console.error('处理WebSocket消息失败:', error)
    }
  },

  // 开始消息轮询
  startMessagePolling() {
    this.data.messagePollingTimer = setInterval(() => {
      if (this.data.consultationId) {
        console.log('轮询获取消息，consultationId:', this.data.consultationId)
        this.loadMessages()
      } else {
        console.log('轮询跳过，consultationId不存在')
      }
    }, 5000) // 5秒轮询一次消息
  },

  // 开始心跳
  startHeartbeat() {
    this.data.heartbeatTimer = setInterval(() => {
      if (this.data.isConnected && this.data.websocket) {
        this.data.websocket.send({
          data: JSON.stringify({ type: 'heartbeat' })
        })
      }
    }, 30000) // 30秒发送一次心跳
  },

  // 停止心跳
  stopHeartbeat() {
    if (this.data.heartbeatTimer) {
      clearInterval(this.data.heartbeatTimer)
      this.data.heartbeatTimer = null
    }
  },

  // 安排重连
  scheduleReconnect() {
    if (this.data.reconnectTimer) {
      clearTimeout(this.data.reconnectTimer)
    }

    this.data.reconnectTimer = setTimeout(() => {
      console.log('尝试重新连接WebSocket...')
      this.connectWebSocket()
    }, 5000) // 5秒后重连
  },

  // 暂停WebSocket
  pauseWebSocket() {
    if (this.data.websocket) {
      this.data.websocket.close()
    }
  },

  // 清理资源
  cleanup() {
    if (this.data.websocket) {
      this.data.websocket.close()
    }
    
    if (this.data.reconnectTimer) {
      clearTimeout(this.data.reconnectTimer)
    }
    
    if (this.data.heartbeatTimer) {
      clearInterval(this.data.heartbeatTimer)
    }
    
    if (this.data.messagePollingTimer) {
      clearInterval(this.data.messagePollingTimer)
    }
  },

  // 输入框内容变化
  onInputChange(e) {
    this.setData({
      inputMessage: e.detail.value
    })
  },

  // 发送消息
  async sendMessage() {
    const message = this.data.inputMessage.trim()
    if (!message) return

    if (!this.data.consultationId) {
      wx.showToast({
        title: '咨询会话未创建，请重试',
        icon: 'none'
      })
      return
    }

    try {
      // 添加用户消息到界面
      const now = new Date()
      const userMessage = {
        id: Date.now().toString(),
        content: message,
        senderType: 'user',
        createTime: this.formatTime(now),
        isPending: true
      }

      this.setData({
        messages: [...this.data.messages, userMessage],
        inputMessage: ''
      })

      this.scrollToBottom()

      // 发送消息到后端
      const res = await consultationAPI.sendConsultationMessage(
        this.data.consultationId,
        message,
        'user'
      )

      if (res.code === 0) {
        // 发送成功，移除pending状态
        const updatedMessages = this.data.messages.map(msg => 
          msg.id === userMessage.id ? { ...msg, isPending: false } : msg
        )
        this.setData({ messages: updatedMessages })

        // 更新咨询状态
        this.setData({ consultationStatus: 'waiting' })

        // 通过WebSocket发送消息
        if (this.data.isConnected && this.data.websocket) {
          this.data.websocket.send({
            data: JSON.stringify({
              type: 'send_message',
              data: {
                consultationId: this.data.consultationId,
                content: message,
                senderType: 'user'
              }
            })
          })
        }

        wx.showToast({
          title: '消息已发送',
          icon: 'success'
        })
      } else {
        throw new Error(res.errorMsg || '发送失败')
      }
    } catch (error) {
      console.error('发送消息失败:', error)
      
      // 发送失败，显示错误提示
      wx.showToast({
        title: '发送失败，请重试',
        icon: 'none'
      })

      // 移除失败的消息
      const updatedMessages = this.data.messages.filter(msg => msg.id !== userMessage.id)
      this.setData({ messages: updatedMessages })
    }
  },

  // 快捷问题点击
  onQuickQuestionTap(e) {
    const question = e.currentTarget.dataset.question
    this.setData({ inputMessage: question })
    this.sendMessage()
  },

  // 返回上一页
  onBack() {
    wx.navigateBack()
  },

  // 重新咨询
  restartConsultation() {
    // 清除当前咨询会话
    wx.removeStorageSync('currentConsultationId')
    
    // 重新初始化
    this.initConsultation()
  },

  // 滚动到底部
  scrollToBottom() {
    if (this.data.messages.length > 0) {
      const lastMessage = this.data.messages[this.data.messages.length - 1]
      this.setData({
        scrollToMessage: `msg-${lastMessage.id}`
      })
    }
  },

  // 格式化时间
  formatTime(date) {
    try {
      console.log('formatTime 输入:', date, '类型:', typeof date)
      
      // 确保date是有效的Date对象
      if (!(date instanceof Date) || isNaN(date.getTime())) {
        console.error('无效的日期对象:', date)
        return '时间未知'
      }

      const now = new Date()
      const diff = now.getTime() - date.getTime()
      
      if (diff < 60000) { // 1分钟内
        return '刚刚'
      } else if (diff < 3600000) { // 1小时内
        return `${Math.floor(diff / 60000)}分钟前`
      } else if (diff < 86400000) { // 24小时内
        return `${Math.floor(diff / 3600000)}小时前`
      } else {
        // 格式化月日时分
        const month = (date.getMonth() + 1).toString().padStart(2, '0')
        const day = date.getDate().toString().padStart(2, '0')
        const hours = date.getHours().toString().padStart(2, '0')
        const minutes = date.getMinutes().toString().padStart(2, '0')
        
        return `${month}月${day}日 ${hours}:${minutes}`
      }
    } catch (error) {
      console.error('时间格式化失败:', error, date)
      return '时间未知'
    }
  }
})
