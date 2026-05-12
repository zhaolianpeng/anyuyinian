// pages/order/detail.js
const { api } = require('../../utils/cloud-container-standard')
const { processImageUrl } = require('../../utils/image')

Page({
  data: {
    orderId: null,
    orderNo: null,
    order: null,
    loading: false,
    countdown: '',
    countdownTimer: null,
    statusMap: {
      0: { text: '待支付', color: '#ff6b35' },
      1: { text: '已支付', color: '#52c41a' },
      2: { text: '已完成', color: '#1890ff' },
      3: { text: '已取消', color: '#999999' },
      4: { text: '已退款', color: '#ff4d4f' }
    },
    payStatusMap: {
      0: { text: '未支付', color: '#ff6b35' },
      1: { text: '已支付', color: '#52c41a' }
    }
  },

  onLoad(options) {
    const { orderNo, id } = options
    if (orderNo) {
      this.setData({ orderNo: orderNo })
      this.loadOrderDetail(orderNo)
    } else if (id) {
      // 通过订单ID获取订单详情
      this.setData({ orderId: id })
      this.loadOrderDetailById(id)
    }
  },

  onShow() {
    // 页面显示时刷新订单详情
    const { orderNo } = this.data
    if (orderNo) {
      this.loadOrderDetail(orderNo)
    }
  },

  onUnload() {
    // 清除定时器
    if (this.data.countdownTimer) {
      clearInterval(this.data.countdownTimer)
    }
  },

  /**
   * 加载订单详情
   */
  async loadOrderDetail(orderNo) {
    try {
      this.setData({ loading: true })
      
      const result = await api.orderDetail({ orderNo })
      
      if (result.code === 0 && result.data) {
        const order = {
          ...result.data,
          serviceImageUrl: processImageUrl(result.data.serviceImageUrl)
        }
        
        this.setData({ 
          order,
          loading: false
        })
        
        // 如果订单待支付，开始倒计时
        if (order.status === 0) {
          this.startCountdown()
        }
        
        console.log('订单详情加载成功:', order)
      } else {
        throw new Error(result.message || '获取订单详情失败')
      }
      
    } catch (error) {
      console.error('加载订单详情失败:', error)
      this.setData({ loading: false })
      
      wx.showToast({
        title: '加载失败，请重试',
        icon: 'none'
      })
    }
  },

  /**
   * 通过订单ID获取订单详情
   */
  async loadOrderDetailById(orderId) {
    try {
      console.log('通过订单ID获取订单详情:', orderId)
      
      // 先通过订单ID获取订单号
      const result = await api.orderDetailById(orderId)
      
      if (result.code === 0 && result.data) {
        const order = result.data
        console.log('通过ID获取订单成功:', order)
        
        // 设置订单号并重新加载详情
        this.setData({ orderNo: order.orderNo })
        this.loadOrderDetail(order.orderNo)
      } else {
        throw new Error(result.errorMsg || '获取订单详情失败')
      }
    } catch (error) {
      console.error('通过ID获取订单详情失败:', error)
      this.setData({ loading: false })
      
      wx.showToast({
        title: '获取订单详情失败',
        icon: 'none'
      })
    }
  },

  /**
   * 支付确认
   */
  async confirmPayment(orderId, transactionId) {
    try {
      console.log('开始支付确认:', { orderId, transactionId })
      
      const result = await api.orderPayConfirm(orderId, {
        transactionId: transactionId || '',
        payMethod: 'wechat'
      })
      
      if (result.code === 0) {
        console.log('支付确认成功:', result.data)
        return result.data
      } else {
        throw new Error(result.errorMsg || '支付确认失败')
      }
    } catch (error) {
      console.error('支付确认失败:', error)
      throw error
    }
  },

  /**
   * 支付订单
   */
  async onPayOrder() {
    const { order } = this.data
    if (!order) return
    
    // 检查订单状态
    if (order.status !== 0) {
      let message = ''
      switch (order.status) {
        case 1:
          message = '订单已支付，无需重复支付'
          break
        case 2:
          message = '订单已完成'
          break
        case 3:
          message = '订单已取消'
          break
        case 4:
          message = '订单已退款'
          break
        default:
          message = '订单状态异常，无法支付'
      }
      
      wx.showToast({
        title: message,
        icon: 'none',
        duration: 2000
      })
      return
    }
    
    try {
      wx.showLoading({ title: '正在支付...' })
      
      // 获取用户OpenID
      const userInfo = wx.getStorageSync('userInfo')
      const openId = userInfo?.openId
      
      if (!openId) {
        throw new Error('用户OpenID不存在，请重新登录')
      }
      
      const result = await api.orderPay(order.id, {
        orderNo: order.orderNo,
        amount: order.totalAmount,
        openId: openId
      })
      
      if (result.code === 0) {
        wx.hideLoading()
        
        // 获取支付参数
        const paymentParams = result.data.paymentParams
        if (!paymentParams) {
          throw new Error('支付参数获取失败')
        }
        
        console.log('支付参数:', paymentParams)
        
        // 调起微信支付
        wx.requestPayment({
          timeStamp: paymentParams.timeStamp,
          nonceStr: paymentParams.nonceStr,
          package: paymentParams.package,
          signType: paymentParams.signType,
          paySign: paymentParams.paySign,
          success: async (res) => {
            console.log('支付成功:', res)
            console.log('准备调用支付确认接口，订单ID:', order.id, '交易ID:', res.transaction_id)
            
            try {
              // 调用支付确认接口
              console.log('开始调用支付确认接口...')
              await this.confirmPayment(order.id, res.transaction_id)
              console.log('支付确认接口调用成功')
              
              wx.showToast({
                title: '支付成功',
                icon: 'success'
              })
              
              // 重新加载订单详情
              console.log('重新加载订单详情...')
              this.loadOrderDetail(order.orderNo)
            } catch (error) {
              console.error('支付确认失败:', error)
              wx.showToast({
                title: '支付确认失败，请刷新页面',
                icon: 'none'
              })
            }
          },
          fail: (err) => {
            console.error('支付失败:', err)
            if (err.errMsg === 'requestPayment:fail cancel') {
              wx.showToast({
                title: '支付已取消',
                icon: 'none'
              })
            } else {
              wx.showToast({
                title: '支付失败，请重试',
                icon: 'none'
              })
            }
          }
        })
      } else {
        throw new Error(result.message || '支付失败')
      }
      
    } catch (error) {
      console.error('支付失败:', error)
      wx.hideLoading()
      
      wx.showToast({
        title: error.message || '支付失败，请重试',
        icon: 'none'
      })
    }
  },

  /**
   * 取消订单
   */
  async onCancelOrder() {
    const { order } = this.data
    if (!order) return
    
    wx.showModal({
      title: '确认取消',
      content: '确定要取消这个订单吗？',
      success: async (res) => {
        if (res.confirm) {
          try {
            wx.showLoading({ title: '正在取消...' })
            
            const result = await api.orderCancel(order.id, {
              orderNo: order.orderNo,
              reason: '用户主动取消'
            })
            
            if (result.code === 0) {
              wx.hideLoading()
              wx.showToast({
                title: '取消成功',
                icon: 'success'
              })
              
              // 重新加载订单详情
              this.loadOrderDetail(order.orderNo)
            } else {
              throw new Error(result.message || '取消失败')
            }
            
          } catch (error) {
            console.error('取消失败:', error)
            wx.hideLoading()
            
            wx.showToast({
              title: error.message || '取消失败，请重试',
              icon: 'none'
            })
          }
        }
      }
    })
  },

  /**
   * 申请退款
   */
  async onRefundOrder() {
    const { order } = this.data
    if (!order) return
    
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
            
            const result = await api.orderRefund(order.id, {
              orderId: order.id,
              refundAmount: order.totalAmount,
              reason: '用户申请退款'
            })
            
            if (result.code === 0) {
              wx.hideLoading()
              wx.showToast({
                title: '退款申请已提交',
                icon: 'success'
              })
              
              // 重新加载订单详情
              this.loadOrderDetail(order.orderNo)
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
   * 联系客服
   */
  onContactService() {
    wx.navigateTo({
      url: '/pages/kefu/chat'
    })
  },

  /**
   * 复制订单号
   */
  onCopyOrderNo() {
    const { order } = this.data
    if (!order || !order.orderNo) return
    
    wx.setClipboardData({
      data: order.orderNo,
      success: () => {
        wx.showToast({
          title: '订单号已复制',
          icon: 'success'
        })
      }
    })
  },

  /**
   * 下拉刷新
   */
  async onPullDownRefresh() {
    const { orderNo } = this.data
    if (orderNo) {
      await this.loadOrderDetail(orderNo)
    }
    wx.stopPullDownRefresh()
  },

  /**
   * 开始倒计时
   */
  startCountdown() {
    const { order } = this.data
    if (!order || !order.payDeadline) return
    
    const deadline = new Date(order.payDeadline).getTime()
    const now = Date.now()
    
    if (deadline <= now) {
      this.setData({ countdown: '支付超时' })
      return
    }
    
    const timer = setInterval(() => {
      this.updateCountdown()
    }, 1000)
    
    this.setData({ countdownTimer: timer })
    this.updateCountdown()
  },

  /**
   * 更新倒计时
   */
  updateCountdown() {
    const { order } = this.data
    if (!order || !order.payDeadline) return
    
    const deadline = new Date(order.payDeadline).getTime()
    const now = Date.now()
    
    if (deadline <= now) {
      this.setData({ countdown: '支付超时' })
      this.stopCountdown()
      return
    }
    
    const diff = deadline - now
    const hours = Math.floor(diff / (1000 * 60 * 60))
    const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
    const seconds = Math.floor((diff % (1000 * 60)) / 1000)
    
    const countdown = `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`
    this.setData({ countdown })
  },

  /**
   * 停止倒计时
   */
  stopCountdown() {
    if (this.data.countdownTimer) {
      clearInterval(this.data.countdownTimer)
      this.setData({ countdownTimer: null })
    }
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
  }
})