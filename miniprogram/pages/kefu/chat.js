// pages/kefu/chat.js
const app = getApp()
const { api } = require('../../utils/cloud-container-standard')
const { getCurrentUserId } = require('../../utils/user-id-compatibility')

Page({
  data: {
    messageList: [],
    inputMessage: '',
    canSend: false,
    scrollToView: '',
    userInfo: {},
    quickReplies: [
      '您好，请问有什么可以帮助您的？',
      '请稍等，我为您查询一下',
      '还有其他问题吗？',
      '感谢您的咨询'
    ],
    showQuickReply: false,
    selectedImages: [],
    loading: false
  },

  onLoad(options) {
    this.loadUserInfo()
    this.loadMessageHistory()
    
    // 显示欢迎消息
    this.addWelcomeMessage()
  },

  onShow() {
    this.loadUserInfo()
  },

  /**
   * 加载用户信息
   */
  loadUserInfo() {
    const userInfo = wx.getStorageSync('userInfo') || {}
    const userId = getCurrentUserId()
    this.setData({
      userInfo: {
        ...userInfo,
        userId: userInfo.userId || userId || ''
      }
    })
  },

  /**
   * 加载消息历史
   */
  loadMessageHistory() {
    // 从本地存储加载消息历史
    const messageList = wx.getStorageSync('chatHistory') || []
    this.setData({ messageList })
  },

  /**
   * 添加欢迎消息
   */
  addWelcomeMessage() {
    const welcomeMessage = {
      type: 0, // 客服消息
      content: '您好！欢迎咨询客服，请问有什么可以帮助您的？',
      messageType: 'text',
      timestamp: Date.now(),
      sender: 'kefu'
    }
    
    this.addMessage(welcomeMessage)
  },

  /**
   * 输入框内容变化
   */
  onInputChange(e) {
    const value = e.detail.value || ''
    this.setData({
      inputMessage: value,
      canSend: value.trim().length > 0
    })
  },

  ensureLoggedIn() {
    const currentUserId = getCurrentUserId() || this.data.userInfo.userId
    if (currentUserId) {
      return true
    }

    const redirect = encodeURIComponent('/pages/kefu/chat')
    wx.showModal({
      title: '登录提示',
      content: '您还没有登录，请先登录后再发送咨询消息',
      showCancel: false,
      confirmText: '去登录',
      success: () => {
        wx.navigateTo({
          url: `/pages/login/login?redirect=${redirect}`
        })
      }
    })
    return false
  },

  /**
   * 发送消息
   */
  async sendMessage() {
    const { inputMessage } = this.data
    
    if (!inputMessage.trim()) {
      return
    }

    if (!this.ensureLoggedIn()) {
      return
    }

    // 添加用户消息
    const userMessage = {
      type: 1, // 用户消息
      content: inputMessage,
      messageType: 'text',
      timestamp: Date.now(),
      sender: 'user'
    }
    
    this.addMessage(userMessage)
    this.setData({ inputMessage: '', canSend: false })
    
    // 发送消息到服务器
    await this.sendMessageToServer(userMessage)
    
    // 模拟客服回复
    this.simulateKefuReply(inputMessage)
  },

  /**
   * 添加消息到列表
   */
  addMessage(message) {
    const messageList = [...this.data.messageList, message]
    this.setData({ messageList })
    
    // 保存到本地存储
    wx.setStorageSync('chatHistory', messageList)
    
    // 滚动到底部
    this.scrollToBottom()
  },

  /**
   * 模拟客服回复
   */
  simulateKefuReply(userMessage) {
    // 简单的自动回复逻辑
    let reply = '感谢您的咨询，我们会尽快为您处理。'
    
    if (userMessage.includes('价格') || userMessage.includes('费用')) {
      reply = '具体价格请查看服务详情页面，或联系客服咨询。'
    } else if (userMessage.includes('时间') || userMessage.includes('预约')) {
      reply = '您可以在服务详情页面选择合适的时间进行预约。'
    } else if (userMessage.includes('退款') || userMessage.includes('取消')) {
      reply = '退款和取消政策请查看订单详情页面，或联系客服处理。'
    }
    
    // 延迟回复，模拟真实场景
    setTimeout(() => {
      const kefuMessage = {
        type: 0, // 客服消息
        content: reply,
        messageType: 'text',
        timestamp: Date.now(),
        sender: 'kefu'
      }
      
      this.addMessage(kefuMessage)
    }, 1000)
  },

  /**
   * 选择图片
   */
  chooseImage() {
    wx.chooseImage({
      count: 1,
      sizeType: ['compressed'],
      sourceType: ['album', 'camera'],
      success: (res) => {
        const tempFilePath = res.tempFilePaths[0]
        
        // 添加图片消息
        const imageMessage = {
          type: 1, // 用户消息
          content: tempFilePath,
          messageType: 'image',
          timestamp: Date.now(),
          sender: 'user'
        }
        
        this.addMessage(imageMessage)
      }
    })
  },

  /**
   * 预览图片
   */
  previewImage(e) {
    const { url } = e.currentTarget.dataset
    wx.previewImage({
      urls: [url]
    })
  },

  /**
   * 选择快捷回复
   */
  selectQuickReply(e) {
    const { reply } = e.currentTarget.dataset
    this.setData({
      inputMessage: reply,
      canSend: !!(reply && reply.trim())
    })
  },

  /**
   * 滚动到底部
   */
  scrollToBottom() {
    this.setData({
      scrollToView: `msg-${this.data.messageList.length - 1}`
    })
  },

  /**
   * 格式化时间
   */
  formatTime(date) {
    const d = new Date(date)
    return `${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}`
  },

  /**
   * 发送消息到服务器
   */
  async sendMessageToServer(messageData) {
    try {
      this.setData({ loading: true })

      const currentUserId = getCurrentUserId() || this.data.userInfo.userId
      if (!currentUserId) {
        throw new Error('用户未登录或缺少userId')
      }

      const { userInfo } = this.data
      
      const result = await api.kefuSendMsg({
        userId: String(currentUserId),
        userName: userInfo.nickName || userInfo.nickname || userInfo.userName || '用户',
        userAvatar: userInfo.avatarUrl || userInfo.avatar || '',
        content: messageData.content,
        images: []
      })
      
      console.log('消息发送成功:', result)
      
    } catch (error) {
      console.error('消息发送失败:', error)
      wx.showToast({
        title: '发送失败，请重试',
        icon: 'none'
      })
    } finally {
      this.setData({ loading: false })
    }
  }
})