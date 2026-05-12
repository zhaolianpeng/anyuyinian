// pages/admin/consultation/list.js
const { consultationAPI } = require('../../../utils/request.js')

Page({
  data: {
    consultations: [], // 咨询列表
    stats: {
      totalCount: 0,
      waitingCount: 0,
      chattingCount: 0,
      todayCount: 0
    },
    filterTabs: [
      { label: '全部', value: 'all' },
      { label: '等待回复', value: 'waiting' },
      { label: '咨询中', value: 'chatting' },
      { label: '已结束', value: 'closed' }
    ],
    currentFilter: 'all',
    loading: false,
    hasMore: true,
    page: 1,
    pageSize: 20
  },

  onLoad() {
    this.loadConsultationStats()
    this.loadConsultations()
  },

  onShow() {
    // 页面显示时刷新数据
    this.refreshData()
  },

  onPullDownRefresh() {
    this.refreshData()
  },

  onReachBottom() {
    this.loadMore()
  },

  // 加载咨询统计
  async loadConsultationStats() {
    try {
      const res = await consultationAPI.getConsultationStats()
      if (res.code === 0 && res.data) {
        this.setData({ stats: res.data })
      }
    } catch (error) {
      console.error('加载咨询统计失败:', error)
    }
  },

  // 加载咨询列表
  async loadConsultations(refresh = false) {
    if (this.data.loading) return

    try {
      this.setData({ loading: true })

      let consultations = []
      if (this.data.currentFilter === 'all') {
        // 获取所有活跃咨询
        const res = await consultationAPI.getActiveConsultations()
        if (res.code === 0 && res.data) {
          consultations = res.data.consultations || []
        }
      } else {
        // 根据状态筛选
        const res = await consultationAPI.getConsultationsByStatus(this.data.currentFilter)
        if (res.code === 0 && res.data) {
          consultations = res.data.consultations || []
        }
      }

      // 处理数据
      const processedConsultations = consultations.map(item => {
        try {
          return {
            ...item,
            createTime: this.formatTime(item.createdAt),
            lastMessage: this.getLastMessage(item.messages),
            messageCount: item.messages && Array.isArray(item.messages) ? item.messages.length : 0
          }
        } catch (error) {
          console.error('处理咨询项失败:', error, item)
          return {
            ...item,
            createTime: '时间未知',
            lastMessage: '暂无消息',
            messageCount: 0
          }
        }
      })

      if (refresh) {
        this.setData({
          consultations: processedConsultations,
          page: 1,
          hasMore: processedConsultations.length >= this.data.pageSize
        })
      } else {
        this.setData({
          consultations: [...this.data.consultations, ...processedConsultations],
          hasMore: processedConsultations.length >= this.data.pageSize
        })
      }

    } catch (error) {
      console.error('加载咨询列表失败:', error)
      wx.showToast({
        title: '加载失败，请重试',
        icon: 'none'
      })
    } finally {
      this.setData({ loading: false })
      wx.stopPullDownRefresh()
    }
  },

  // 筛选标签点击
  onFilterTabTap(e) {
    const value = e.currentTarget.dataset.value
    if (value === this.data.currentFilter) return

    this.setData({
      currentFilter: value,
      consultations: [],
      page: 1,
      hasMore: true
    })

    this.loadConsultations(true)
  },

  // 咨询项点击
  onConsultationTap(e) {
    const consultationId = e.currentTarget.dataset.id
    wx.navigateTo({
      url: `/pages/admin/consultation/chat?consultationId=${consultationId}`
    })
  },

  // 回复按钮点击
  onReplyTap(e) {
    const consultationId = e.currentTarget.dataset.id
    wx.navigateTo({
      url: `/pages/admin/consultation/chat?consultationId=${consultationId}`
    })
  },

  // 关闭咨询
  async onCloseTap(e) {
    const consultationId = e.currentTarget.dataset.id
    
    wx.showModal({
      title: '确认关闭',
      content: '确定要关闭这个咨询会话吗？关闭后将无法继续对话。',
      confirmText: '确认关闭',
      cancelText: '取消',
      success: async (res) => {
        if (res.confirm) {
          try {
            const result = await consultationAPI.closeConsultation({
              consultationId: consultationId
            })
            
            if (result.code === 0) {
              wx.showToast({
                title: '咨询已关闭',
                icon: 'success'
              })
              
              // 刷新数据
              this.refreshData()
            } else {
              throw new Error(result.errorMsg || '关闭失败')
            }
          } catch (error) {
            console.error('关闭咨询失败:', error)
            wx.showToast({
              title: '关闭失败，请重试',
              icon: 'none'
            })
          }
        }
      }
    })
  },

  // 加载更多
  loadMore() {
    if (!this.data.hasMore || this.data.loading) return
    
    this.setData({
      page: this.data.page + 1
    })
    
    this.loadConsultations()
  },

  // 刷新数据
  refreshData() {
    this.loadConsultationStats()
    this.loadConsultations(true)
  },

  // 获取最后一条消息
  getLastMessage(messages) {
    try {
      if (!messages || !Array.isArray(messages) || messages.length === 0) {
        return '暂无消息'
      }
      
      const lastMessage = messages[messages.length - 1]
      return lastMessage.content || '暂无消息'
    } catch (error) {
      console.error('获取最后一条消息失败:', error)
      return '暂无消息'
    }
  },

  // 格式化时间
  formatTime(dateString) {
    if (!dateString) return ''
    
    const date = new Date(dateString)
    const now = new Date()
    const diff = now - date
    
    if (diff < 60000) { // 1分钟内
      return '刚刚'
    } else if (diff < 3600000) { // 1小时内
      return `${Math.floor(diff / 60000)}分钟前`
    } else if (diff < 86400000) { // 24小时内
      return `${Math.floor(diff / 3600000)}小时前`
    } else {
      return `${date.getMonth() + 1}月${date.getDate()}日`
    }
  }
})