const { api } = require('../../utils/cloud-container-standard')

Page({
  data: {
    orders: [],
    loading: true,
    hasMore: true,
    page: 1,
    pageSize: 10,
    status: '', // 确保初始状态为空
    countdownTimer: null,
    needLogin: false,
    statusOptions: [
      { name: '全部', value: '', count: 0 },
      { name: '待支付', value: '0', count: 0 },
      { name: '已支付', value: '1', count: 0 },
      { name: '已完成', value: '2', count: 0 },
      { name: '已取消', value: '3', count: 0 },
      { name: '已退款', value: '4', count: 0 }
    ],
    currentTab: 0, // 当前选中的tab索引
    tabLoading: false // tab切换时的加载状态
  },

  onLoad(options) {
    console.log('订单列表页面加载，初始状态:', this.data.status)
    console.log('初始 currentTab:', this.data.currentTab)
    console.log('页面参数 options:', options)
    
    // 检查页面参数是否传递了状态
    if (options && options.status) {
      console.log('页面参数传递了状态:', options.status)
      this.setData({ 
        status: options.status,
        currentTab: this.getTabIndexByStatus(options.status)
      })
    } else {
      // 如果没有传递状态参数，默认显示全部订单
      console.log('没有传递状态参数，默认显示全部订单')
      this.setData({ 
        status: '',
        currentTab: 0
      })
    }
    
    this.loadOrders()
    this.startCountdown()
  },

  onUnload() {
    // 清除定时器
    if (this.data.countdownTimer) {
      clearInterval(this.data.countdownTimer)
    }
  },

  onShow() {
    // 每次显示页面时刷新订单列表，但保持当前状态筛选
    console.log('订单列表页面显示，当前状态:', this.data.status)
    console.log('当前 currentTab:', this.data.currentTab)
    
    // 只刷新数据，不重置状态
    this.loadOrders(true)
  },

  // 去登录
  onGoLogin() {
    const redirect = '/pages/order/list'
    wx.navigateTo({ url: `/pages/login/login?redirect=${encodeURIComponent(redirect)}` })
  },

  /**
   * 加载订单列表
   */
  async loadOrders(isRefresh = false) {
    try {
      if (isRefresh) {
        this.setData({ 
          page: 1, 
          hasMore: true,
          loading: true 
        })
      } else {
        this.setData({ loading: true })
      }
      
      // 获取用户ID
      const userId = wx.getStorageSync('userId')
      if (!userId) {
        console.log('用户未登录，显示登录门')
        this.setData({ loading: false, needLogin: true, orders: [], hasMore: false })
        return
      }
      
      console.log('当前用户ID:', userId)
      console.log('当前状态筛选:', this.data.status)
      
      const params = {
        page: this.data.page,
        pageSize: this.data.pageSize,
        status: this.data.status,
        userId: userId
      }
      
      console.log('请求订单列表参数:', params)
      
      const result = await api.orderList(params)
      
      console.log('订单列表接口返回:', result)
      
      if (result.code === 0 && result.data) {
        // 后端返回的是 list，前端期望的是 orders
        const { list, total, hasMore } = result.data
        const orders = list || [] // 将 list 映射为 orders
        
        console.log('解析的订单数据:', { orders, total, hasMore })
        
        // 确保 orders 是数组
        if (!orders || !Array.isArray(orders)) {
          console.log('订单数据为空或格式错误，设置为空数组')
          this.setData({ 
            orders: [],
            hasMore: false,
            loading: false 
          })
          return
        }
        
        // 处理订单数据
        const processedOrders = orders.map(order => {
          console.log('处理订单数据:', {
            orderId: order.id,
            orderNo: order.orderNo,
            totalAmount: order.totalAmount,
            formattedAmount: order.formattedAmount,
            price: order.price
          })
          
          return {
            ...order,
            statusText: this.getStatusText(order.status),
            statusClass: this.getStatusClass(order.status),
            formattedAmount: order.formattedAmount || this.formatAmount(order.totalAmount),
            formattedCreatedAt: this.formatTime(order.createdAt),
            formattedDate: this.formatTime(order.createdAt) // 确保使用正确的时间格式化
          }
        })
        
        console.log('处理后的订单数据:', processedOrders)
        console.log('订单数量:', processedOrders.length)
        
        if (isRefresh) {
          this.setData({ 
            orders: processedOrders,
            hasMore: hasMore,
            loading: false,
            needLogin: false
          })
        } else {
          this.setData({ 
            orders: [...this.data.orders, ...processedOrders],
            hasMore: hasMore,
            loading: false,
            needLogin: false
          })
        }
        
        console.log('设置到页面的订单数据:', this.data.orders)
        
        // 更新状态计数
        this.updateStatusCounts(processedOrders)
        
        console.log('订单列表加载成功:', processedOrders.length)
      } else {
        throw new Error(result.message || '获取订单列表失败')
      }
      
    } catch (error) {
      console.error('加载订单列表失败:', error)
      this.setData({ loading: false })
      
      wx.showToast({
        title: '加载失败，请重试',
        icon: 'none'
      })
    }
  },

  /**
   * 下拉刷新
   */
  async onPullDownRefresh() {
    await this.loadOrders(true)
    wx.stopPullDownRefresh()
  },

  /**
   * 上拉加载更多
   */
  onReachBottom() {
    if (this.data.hasMore && !this.data.loading) {
      this.setData({ page: this.data.page + 1 })
      this.loadOrders()
    }
  },

  /**
   * 根据状态值获取对应的 tab 索引
   */
  getTabIndexByStatus(status) {
    const statusOptions = this.data.statusOptions
    for (let i = 0; i < statusOptions.length; i++) {
      if (statusOptions[i].value === status) {
        return i
      }
    }
    return 0 // 默认返回第一个 tab
  },

  /**
   * 选择状态筛选
   */
  onStatusSelect(e) {
    const { index } = e.currentTarget.dataset
    const status = this.data.statusOptions[index].value
    
    console.log('状态选择:', { index, status, statusName: this.data.statusOptions[index].name })
    
    // 设置tab加载状态
    this.setData({ 
      tabLoading: true,
      status,
      currentTab: index,
      page: 1,
      hasMore: true
    })
    
    console.log('设置状态后:', this.data.status)
    
    // 加载对应状态的订单
    this.loadOrders(true).then(() => {
      // 加载完成后关闭tab加载状态
      this.setData({ tabLoading: false })
    }).catch(() => {
      // 加载失败后也要关闭tab加载状态
      this.setData({ tabLoading: false })
    })
  },

  /**
   * 点击订单
   */
  onOrderTap(e) {
    const { orderNo } = e.currentTarget.dataset
    console.log('点击订单，订单号:', orderNo)
    console.log('事件数据:', e.currentTarget.dataset)
    
    if (!orderNo) {
      console.error('订单号为空')
      wx.showToast({
        title: '订单号缺失',
        icon: 'none'
      })
      return
    }
    
    wx.navigateTo({
      url: `/pages/order/detail?orderNo=${orderNo}`
    })
  },

  /**
   * 取消订单
   */
  async onCancelOrder(e) {
    const { orderId, orderNo } = e.currentTarget.dataset
    
    console.log('取消订单，参数:', { orderId, orderNo })
    
    wx.showModal({
      title: '确认取消',
      content: '确定要取消这个订单吗？',
      success: async (res) => {
        if (res.confirm) {
          try {
            wx.showLoading({ title: '正在取消...' })
            
            console.log('开始调用取消订单API，orderId:', orderId)
            
            const result = await api.orderCancel(orderId, {
              orderId: parseInt(orderId),
              reason: '用户主动取消'
            })
            
            console.log('取消订单API返回:', result)
            
            if (result.code === 0) {
              wx.hideLoading()
              wx.showToast({
                title: '取消成功',
                icon: 'success'
              })
              
              // 刷新订单列表
              this.loadOrders(true)
            } else {
              throw new Error(result.errorMsg || result.message || '取消失败')
            }
            
          } catch (error) {
            console.error('取消失败:', error)
            wx.hideLoading()
            
            wx.showToast({
              title: error.message || '取消失败，请重试',
              icon: 'none',
              duration: 3000
            })
          }
        }
      }
    })
  },

  /**
   * 支付订单
   */
  onPayOrder(e) {
    const { orderNo } = e.currentTarget.dataset
    console.log('支付订单，orderNo:', orderNo)
    wx.navigateTo({
      url: `/pages/order/detail?orderNo=${orderNo}`
    })
  },

  /**
   * 申请退款
   */
  async onRefundOrder(e) {
    const { orderId, orderNo } = e.currentTarget.dataset
    console.log('申请退款，参数:', { orderId, orderNo })
    
    // 获取订单信息
    const order = this.data.orders.find(o => o.id == orderId)
    if (!order) {
      wx.showToast({
        title: '订单信息不存在',
        icon: 'none'
      })
      return
    }
    
    // 检查订单状态
    if (order.status !== 1) {
      wx.showToast({
        title: '只有已支付的订单可以申请退款',
        icon: 'none'
      })
      return
    }
    
    // 检查是否已经申请过退款
    if (order.refundStatus > 0) {
      wx.showToast({
        title: '订单已申请退款，请勿重复申请',
        icon: 'none'
      })
      return
    }
    
    wx.showModal({
      title: '申请退款',
      content: `确定要申请退款吗？\n退款金额：¥${order.totalAmount}`,
      success: async (res) => {
        if (res.confirm) {
          try {
            wx.showLoading({ title: '正在申请退款...' })
            
            const result = await api.orderRefund(orderId, {
              orderId: parseInt(orderId),
              refundAmount: order.totalAmount,
              reason: '用户申请退款'
            })
            
            if (result.code === 0) {
              wx.hideLoading()
              wx.showToast({
                title: '退款申请已提交',
                icon: 'success'
              })
              
              // 刷新订单列表
              this.loadOrders(true)
            } else {
              throw new Error(result.errorMsg || '退款申请失败')
            }
            
          } catch (error) {
            console.error('退款申请失败:', error)
            wx.hideLoading()
            
            wx.showToast({
              title: error.message || '退款申请失败，请重试',
              icon: 'none'
            })
          }
        }
      }
    })
  },

  /**
   * 获取状态文本
   */
  getStatusText(status) {
    const statusMap = {
      0: '待支付',
      1: '已支付',
      2: '已完成',
      3: '已取消',
      4: '已退款'
    }
    return statusMap[status] || status
  },

  /**
   * 获取状态样式类
   */
  getStatusClass(status) {
    const classMap = {
      0: 'status-pending',
      1: 'status-paid',
      2: 'status-completed',
      3: 'status-cancelled',
      4: 'status-refunded'
    }
    return classMap[status] || 'status-default'
  },

  /**
   * 跳转到服务页面
   */
  goToService() {
    wx.navigateTo({
      url: '/pages/service/list'
    })
  },

  /**
   * 开始倒计时
   */
  startCountdown() {
    const timer = setInterval(() => {
      this.updateCountdown()
    }, 1000)
    
    this.setData({ countdownTimer: timer })
  },

  /**
   * 更新倒计时
   */
  updateCountdown() {
    const orders = this.data.orders.map(order => {
      if (order.status === 0 && order.payDeadline) {
        const deadline = new Date(order.payDeadline).getTime()
        const now = Date.now()
        
        if (deadline > now) {
          const diff = deadline - now
          const hours = Math.floor(diff / (1000 * 60 * 60))
          const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
          const seconds = Math.floor((diff % (1000 * 60)) / 1000)
          
          order.countdown = `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`
        } else {
          order.countdown = '支付超时'
          order.status = 3 // 自动取消
          order.statusText = this.getStatusText(3)
          order.statusClass = this.getStatusClass(3)
        }
      }
      return order
    })
    
    this.setData({ orders })
  },

  /**
   * 格式化金额
   */
  formatAmount(amount) {
    if (!amount) return '¥0.00'
    return `¥${parseFloat(amount).toFixed(2)}`
  },

  /**
   * 格式化时间
   */
  formatTime(time) {
    if (!time) return ''
    
    console.log('格式化时间，原始时间:', time)
    
    // 创建日期对象
    const date = new Date(time)
    
    // 检查日期是否有效
    if (isNaN(date.getTime())) {
      console.warn('无效的时间格式:', time)
      return ''
    }
    
    console.log('解析后的日期对象:', date)
    console.log('UTC时间:', date.toISOString())
    console.log('本地时间:', date.toLocaleString('zh-CN'))
    
    // 获取本地时区偏移量（分钟）
    const timezoneOffset = date.getTimezoneOffset()
    console.log('时区偏移量（分钟）:', timezoneOffset)
    
    // 调整到本地时间
    const localDate = new Date(date.getTime() - (timezoneOffset * 60 * 1000))
    console.log('调整后的本地时间:', localDate)
    
    // 格式化为本地时间字符串
    const year = localDate.getFullYear()
    const month = String(localDate.getMonth() + 1).padStart(2, '0')
    const day = String(localDate.getDate()).padStart(2, '0')
    const hours = String(localDate.getHours()).padStart(2, '0')
    const minutes = String(localDate.getMinutes()).padStart(2, '0')
    
    const formattedTime = `${year}-${month}-${day} ${hours}:${minutes}`
    console.log('格式化后的时间:', formattedTime)
    
    return formattedTime
  },

  /**
   * 更新状态计数
   */
  updateStatusCounts(orders) {
    const statusCounts = {
      '': orders.length,
      '0': 0,
      '1': 0,
      '2': 0,
      '3': 0,
      '4': 0
    }
    
    orders.forEach(order => {
      const status = order.status.toString()
      if (statusCounts.hasOwnProperty(status)) {
        statusCounts[status]++
      }
    })
    
    const statusOptions = this.data.statusOptions.map(option => ({
      ...option,
      count: statusCounts[option.value] || 0
    }))
    
    this.setData({ statusOptions })
  },

  // 分享给好友
  onShareAppMessage(res) {
    console.log('订单列表分享被触发:', res)
    
    return {
      title: '安语颐年护理陪诊 - 我的订单',
      desc: '查看我的护理陪诊订单',
      path: 'pages/order/list'
    }
  },

  // 分享到朋友圈
  onShareTimeline(res) {
    console.log('订单列表分享到朋友圈被触发:', res)
    
    return {
      title: '安语颐年护理陪诊 - 我的订单',
      query: '',
      imageUrl: ''
    }
  }
}) 