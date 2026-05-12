const { api } = require("../../utils/cloud-container-standard")
// pages/admin/orders.js
Page({
  data: {
    adminInfo: null,
    orders: [],
    loading: true,
    hasMore: true,
    page: 1,
    pageSize: 20
  },

  onLoad: function (options) {
    // 确保数据初始化
    this.setData({
      orders: [],
      loading: true,
      hasMore: true,
      page: 1
    });
    
    this.checkAdminAuth();
    this.loadOrders();
  },

  onShow: function () {
    this.checkAdminAuth();
  },

  // 检查管理员权限
  checkAdminAuth: function () {
    const adminInfo = wx.getStorageSync('adminInfo');
    const isAdmin = wx.getStorageSync('isAdmin');
    
    if (!adminInfo || !isAdmin) {
      wx.redirectTo({
        url: '/pages/admin/login'
      });
      return;
    }

    this.setData({
      adminInfo: adminInfo
    });
  },

  // 加载订单列表
  loadOrders: function (refresh = false) {
    const adminInfo = wx.getStorageSync('adminInfo');
    if (!adminInfo) return;

    if (refresh) {
      this.setData({
        page: 1,
        orders: [],
        hasMore: true
      });
    }

    if (!this.data.hasMore) return;

    const app = getApp();
    app.callContainer('/api/admin/orders', 'GET', {}, {
      query: {
        adminUserId: adminInfo.userId,
        page: this.data.page,
        pageSize: this.data.pageSize
      }
    }).then(res => {
      this.setData({ loading: false });
      console.log('订单返回：', res);
      console.log('当前orders数据：', this.data.orders);
      if (res.code === 0) {
        const newOrders = res.data.list || [];
        const hasMore = res.data.hasMore;
        
        // 处理订单数据，添加状态文本和颜色
        const processedOrders = newOrders.map(order => ({
          ...order,
          statusText: this.getOrderStatusText(order.status),
          statusColor: this.getOrderStatusColor(order.status)
        }));
        
        this.setData({
          orders: refresh ? processedOrders : [...this.data.orders, ...processedOrders],
          hasMore: hasMore,
          page: this.data.page + 1
        });
      } else {
        wx.showToast({
          title: res.errorMsg || '获取订单失败',
          icon: 'none'
        });
      }
    }).catch(err => {
      this.setData({ 
        loading: false,
        orders: this.data.orders || [] // 确保orders始终是数组
      });
      wx.showToast({
        title: '网络错误，请重试',
        icon: 'none'
      });
      console.error('获取订单列表失败:', err);
    });
  },

  // 下拉刷新
  onPullDownRefresh: function () {
    this.loadOrders(true);
    wx.stopPullDownRefresh();
  },

  // 上拉加载更多
  onReachBottom: function () {
    if (this.data.hasMore && !this.data.loading) {
      this.loadOrders();
    }
  },

  // 查看订单详情
  onOrderDetail: function (e) {
    const orderId = e.currentTarget.dataset.orderId;
    wx.navigateTo({
      url: `/pages/order/detail?id=${orderId}`
    });
  },

  // 获取订单状态文本
  getOrderStatusText: function (status) {
    const statusMap = {
      0: '待支付',
      1: '已支付',
      2: '已完成',
      3: '已取消',
      4: '已退款'
    };
    return statusMap[status] || '未知';
  },

  // 获取订单状态颜色
  getOrderStatusColor: function (status) {
    const colorMap = {
      0: '#ff9500',
      1: '#007aff',
      2: '#ff3b30',
      3: '#34c759'
    };
    return colorMap[status] || '#999999';
  },

  // 修改订单金额
  onEditAmount: function (e) {
    const order = e.currentTarget.dataset.order;
    if (!order) {
      wx.showToast({
        title: '订单数据无效',
        icon: 'none'
      });
      return;
    }
    
    const adminInfo = wx.getStorageSync('adminInfo');
    
    if (!adminInfo || adminInfo.adminLevel !== 2) {
      wx.showToast({
        title: '权限不足',
        icon: 'none'
      });
      return;
    }

    wx.showModal({
      title: '修改订单金额',
      content: '',
      editable: true,
      placeholderText: '请输入新金额',
      success: (res) => {
        if (res.confirm) {
          const newAmount = parseFloat(res.content);
          if (isNaN(newAmount) || newAmount <= 0) {
            wx.showToast({
              title: '请输入有效金额',
              icon: 'none'
            });
            return;
          }

          this.updateOrderAmount(order.id, newAmount);
        }
      }
    });
  },

  // 更新订单金额
  updateOrderAmount: function (orderId, newAmount) {
    const adminInfo = wx.getStorageSync('adminInfo');
    
    wx.showLoading({
      title: '正在修改...'
    });

    const app = getApp();
    app.callContainer('/api/admin/order/update-amount', 'POST', {
      orderId: orderId,
      newAmount: newAmount,
      reason: '管理员手动修改'
    }, {
      query: {
        adminUserId: adminInfo.userId
      }
    }).then(res => {
      wx.hideLoading();
      if (res.code === 0) {
        wx.showToast({
          title: '修改成功',
          icon: 'success'
        });
        // 刷新订单列表
        this.loadOrders(true);
      } else {
        wx.showToast({
          title: res.errorMsg || '修改失败',
          icon: 'none'
        });
      }
    }).catch(err => {
      wx.hideLoading();
      wx.showToast({
        title: '网络错误，请重试',
        icon: 'none'
      });
      console.error('修改订单金额失败:', err);
    });
  },

  // 处理退款
  onRefundOrder: function (e) {
    const order = e.currentTarget.dataset.order;
    if (!order) {
      wx.showToast({
        title: '订单数据无效',
        icon: 'none'
      });
      return;
    }
    
    const adminInfo = wx.getStorageSync('adminInfo');
    
    if (!adminInfo || adminInfo.adminLevel !== 2) {
      wx.showToast({
        title: '权限不足',
        icon: 'none'
      });
      return;
    }

    // 检查订单状态
    if (order.status !== 1) {
      wx.showToast({
        title: '只有已支付的订单可以退款',
        icon: 'none'
      });
      return;
    }

    wx.showActionSheet({
      itemList: ['设置为退款中', '设置为已退款'],
      success: (res) => {
        const refundStatus = res.tapIndex === 0 ? 1 : 2;
        const statusText = refundStatus === 1 ? '退款中' : '已退款';
        
        wx.showModal({
          title: '处理退款',
          content: `确定要将订单设置为${statusText}吗？\n退款金额：¥${order.amount}`,
          editable: true,
          placeholderText: '请输入退款原因',
          success: (modalRes) => {
            if (modalRes.confirm) {
              const reason = modalRes.content || '管理员处理退款';
              this.processRefund(order.id, order.amount, reason, refundStatus);
            }
          }
        });
      }
    });
  },

  // 处理退款请求
  processRefund: function (orderId, refundAmount, reason, refundStatus) {
    const adminInfo = wx.getStorageSync('adminInfo');
    
    wx.showLoading({
      title: '正在处理退款...'
    });

    const app = getApp();
    app.callContainer('/api/admin/order/refund', 'POST', {
      orderId: orderId,
      refundAmount: refundAmount,
      reason: reason,
      refundStatus: refundStatus
    }, {
      query: {
        adminUserId: adminInfo.userId
      }
    }).then(res => {
      wx.hideLoading();
      if (res.code === 0) {
        wx.showToast({
          title: '退款处理成功',
          icon: 'success'
        });
        // 刷新订单列表
        this.loadOrders(true);
      } else {
        wx.showToast({
          title: res.errorMsg || '退款处理失败',
          icon: 'none'
        });
      }
    }).catch(err => {
      wx.hideLoading();
      wx.showToast({
        title: '网络错误，请重试',
        icon: 'none'
      });
      console.error('退款处理失败:', err);
    });
  }
}); 