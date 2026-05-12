const { api } = require("../../utils/cloud-container-standard")
// pages/admin/home.js
Page({
  data: {
    adminInfo: null,
    stats: {
      totalUsers: 0,
      totalOrders: 0,
      todayOrders: 0,
      totalAmount: 0
    },
    loading: true
  },

  onLoad: function (options) {
    this.checkAdminAuth();
    this.loadAdminInfo();
    this.loadStats();
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

  // 加载管理员信息
  loadAdminInfo: function () {
    const adminInfo = wx.getStorageSync('adminInfo');
    if (adminInfo) {
      this.setData({
        adminInfo: adminInfo
      });
    }
  },

  // 加载统计数据
  loadStats: function () {
    const adminInfo = wx.getStorageSync('adminInfo');
    console.log('当前adminInfo:', adminInfo);
    console.log('调用了 loadStats');
    if (!adminInfo) {
      console.warn('adminInfo 为空，无法请求数据概览');
      return;
    }
    console.log('请求数据概览参数：', { adminUserId: adminInfo.userId });
    const app = getApp();
    
    // 加载基础统计数据
    app.callContainer('/api/admin/stats', 'GET', {
      adminUserId: adminInfo.userId
    }).then(res => {
      console.log('数据概览返回：', res);
      if (res.code === 0) {
        this.setData({
          stats: { ...this.data.stats, ...res.data }
        });
      }
    }).catch(err => {
      console.error('数据概览请求失败：', err);
    });

    // 加载咨询统计数据
    app.callContainer('/api/consultation/stats', 'GET').then(res => {
      console.log('咨询统计返回：', res);
      if (res.code === 0) {
        this.setData({
          stats: {
            ...this.data.stats,
            waitingConsultations: res.data.waitingCount || 0,
            totalConsultations: res.data.totalCount || 0
          }
        });
      }
    }).catch(err => {
      console.error('咨询统计请求失败：', err);
    });

    this.setData({ loading: false });
  },

  // 跳转到订单管理
  onGoToOrders: function () {
    wx.navigateTo({
      url: '/pages/admin/orders'
    });
  },

  // 跳转到用户管理
  onGoToUsers: function () {
    wx.navigateTo({
      url: '/pages/admin/users'
    });
  },

  // 跳转到管理员管理
  onGoToAdmins: function () {
    wx.navigateTo({
      url: '/pages/admin/admins'
    });
  },

  // 跳转到服务管理
  onGoToServices: function () {
    wx.navigateTo({
      url: '/pages/admin/services'
    });
  },

  // 跳转到在线咨询
  onGoToConsultations: function () {
    wx.navigateTo({
      url: '/pages/admin/consultation/list'
    });
  },

  // 退出管理员模式
  onLogout: function () {
    wx.showModal({
      title: '确认退出',
      content: '确定要退出管理员模式吗？',
      success: (res) => {
        if (res.confirm) {
          // 清除管理员信息
          wx.removeStorageSync('adminInfo');
          wx.removeStorageSync('isAdmin');
          
          // 跳转到普通用户首页
          wx.switchTab({
            url: '/pages/index/index'
          });
        }
      }
    });
  }
}); 