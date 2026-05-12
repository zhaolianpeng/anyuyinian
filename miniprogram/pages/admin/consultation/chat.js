const { consultationAPI } = require('../../../utils/request.js')

Page({
  data: {
    consultationId: '', // 咨询ID
    consultation: null, // 咨询信息
    messages: [], // 聊天消息
    inputMessage: '', // 输入框内容
    scrollToMessage: '', // 滚动到指定消息
    loading: false,
    messagePollingTimer: null // 消息轮询定时器
  },

  onLoad(options) {
    console.log('管理员聊天页面参数:', options)
    const { consultationId } = options
    
    if (consultationId) {
      console.log('设置consultationId:', consultationId)
      this.setData({ consultationId })
      this.loadConsultationDetail()
      this.loadMessages()
      this.startMessagePolling()
    } else {
      console.error('缺少consultationId参数')
      wx.showToast({
        title: '参数错误',
        icon: 'none'
      })
      setTimeout(() => {
        wx.navigateBack()
      }, 1500)
    }
  },

  onShow() {
    // 页面显示时刷新消息（如果有consultationId）
    if (this.data.consultationId) {
      this.loadMessages()
    }
  },

  onUnload() {
    // 页面卸载时清理资源
    this.cleanup()
  },

  // 加载咨询详情
  async loadConsultationDetail() {
    try {
      const res = await consultationAPI.getConsultationDetail(this.data.consultationId)
      
      if (res.code === 0 && res.data) {
        this.setData({ consultation: res.data.consultation })
      }
    } catch (error) {
      console.error('加载咨询详情失败:', error)
    }
  },

  // 加载消息记录
  async loadMessages() {
    if (!this.data.consultationId) return

    try {
      const res = await consultationAPI.getConsultationMessages(this.data.consultationId)

      if (res.code === 0 && res.data) {
        const messages = res.data.messages.map(msg => ({
          ...msg,
          createTime: this.formatTime(new Date(msg.createdAt))
        }))

        this.setData({ messages })
        this.scrollToBottom()
      }
    } catch (error) {
      console.error('加载消息记录失败:', error)
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

    // 在函数开始就定义adminMessage，确保作用域正确
    let adminMessage = null

    try {
      // 添加管理员消息到界面
      adminMessage = {
        id: Date.now().toString(),
        content: message,
        senderType: 'admin',
        createTime: this.formatTime(new Date()),
        isPending: true
      }

      this.setData({
        messages: [...this.data.messages, adminMessage],
        inputMessage: ''
      })

      this.scrollToBottom()

      // 发送消息到后端
      console.log('准备发送消息，参数:', {
        consultationId: this.data.consultationId,
        content: message,
        senderType: 'admin'
      })
      
      const res = await consultationAPI.sendConsultationMessage(
        this.data.consultationId,
        message,
        'admin'
      )

      if (res.code === 0) {
        // 发送成功，移除pending状态
        const updatedMessages = this.data.messages.map(msg => 
          msg.id === adminMessage.id ? { ...msg, isPending: false } : msg
        )
        this.setData({ messages: updatedMessages })

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

      // 移除失败的消息（如果adminMessage存在）
      if (adminMessage) {
        const updatedMessages = this.data.messages.filter(msg => msg.id !== adminMessage.id)
        this.setData({ messages: updatedMessages })
      }
    }
  },

  // 开始消息轮询
  startMessagePolling() {
    if (!this.data.consultationId) {
      console.log('没有consultationId，跳过消息轮询')
      return
    }
    
    console.log('开始消息轮询，consultationId:', this.data.consultationId)
    this.data.messagePollingTimer = setInterval(() => {
      if (this.data.consultationId) {
        console.log('轮询获取消息，consultationId:', this.data.consultationId)
        this.loadMessages()
      } else {
        console.log('轮询跳过，consultationId不存在')
      }
    }, 3000) // 3秒轮询一次消息
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
    const now = new Date()
    const diff = now - date
    
    if (diff < 60000) { // 1分钟内
      return '刚刚'
    } else if (diff < 3600000) { // 1小时内
      return `${Math.floor(diff / 60000)}分钟前`
    } else if (diff < 86400000) { // 24小时内
      return `${date.getHours()}:${date.getMinutes().toString().padStart(2, '0')}`
    } else {
      return `${date.getMonth() + 1}月${date.getDate()}日 ${date.getHours()}:${date.getMinutes().toString().padStart(2, '0')}`
    }
  },

  // 清理资源
  cleanup() {
    if (this.data.messagePollingTimer) {
      clearInterval(this.data.messagePollingTimer)
    }
  },

  // 返回上一页
  onBack() {
    wx.navigateBack()
  }
})
