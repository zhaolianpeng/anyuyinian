const { api } = require("../../utils/cloud-container-standard")
// pages/admin/admins.js
Page({
  data: {
    adminInfo: null,
    admins: [],
    loading: true,
    hasMore: true,
    page: 1,
    pageSize: 20
  },

  onLoad: function (options) {
    this.checkAdminAuth();
    this.loadAdmins();
  },

  onShow: function () {
    this.checkAdminAuth();
  },

  // 检查管理员权限
  checkAdminAuth: function () {
    const adminInfo = wx.getStorageSync('adminInfo');
    const isAdmin = wx.getStorageSync('isAdmin');
    
    if (!adminInfo || !isAdmin || adminInfo.adminLevel !== 2) {
      wx.showToast({
        title: '权限不足',
        icon: 'none'
      });
      setTimeout(() => {
        wx.navigateBack();
      }, 1500);
      return;
    }

    this.setData({
      adminInfo: adminInfo
    });
  },

  // 加载管理员列表
  loadAdmins: function (refresh = false) {
    const adminInfo = wx.getStorageSync('adminInfo');
    if (!adminInfo) return;

    if (refresh) {
      this.setData({
        page: 1,
        admins: [],
        hasMore: true
      });
    }

    if (!this.data.hasMore) return;

    const app = getApp();
    app.callContainer('/api/admin/admins', 'GET', {
      adminUserId: adminInfo.userId,
      page: this.data.page,
      pageSize: this.data.pageSize
    }).then(res => {
      this.setData({ loading: false });
      console.log('管理员列表返回：', res);
      if (res.code === 0) {
        const newAdmins = res.data.list || [];
        const hasMore = res.data.hasMore;
        
        // 调试：打印每个管理员的adminCreatedAt
        newAdmins.forEach((admin, index) => {
          console.log(`管理员${index + 1}:`, {
            nickName: admin.nickName,
            adminCreatedAt: admin.adminCreatedAt,
            formattedTime: this.formatTime(admin.adminCreatedAt)
          });
          
          // 预先格式化时间
          admin.formattedAdminTime = this.formatTime(admin.adminCreatedAt);
        });
        
        this.setData({
          admins: refresh ? newAdmins : [...this.data.admins, ...newAdmins],
          hasMore: hasMore,
          page: this.data.page + 1
        });
      } else {
        wx.showToast({
          title: res.errorMsg || '获取管理员列表失败',
          icon: 'none'
        });
      }
    }).catch(err => {
      this.setData({ loading: false });
      wx.showToast({
        title: '网络错误，请重试',
        icon: 'none'
      });
      console.error('获取管理员列表失败:', err);
    });
  },

  // 下拉刷新
  onPullDownRefresh: function () {
    this.loadAdmins(true);
    wx.stopPullDownRefresh();
  },

  // 上拉加载更多
  onReachBottom: function () {
    if (this.data.hasMore && !this.data.loading) {
      this.loadAdmins();
    }
  },

  // 设置用户为管理员
  onSetAsAdmin: function (e) {
    const userId = e.currentTarget.dataset.userId;
    const user = this.data.admins.find(u => u.userId === userId);
    
    if (!user) return;

    wx.showModal({
      title: '设置管理员',
      content: `确定要将用户"${user.nickName}"设置为管理员吗？`,
      success: (res) => {
        if (res.confirm) {
          this.showSetAdminModal(user);
        }
      }
    });
  },

  // 显示设置管理员弹窗
  showSetAdminModal: function (user) {
    wx.showModal({
      title: '设置管理员',
      content: '请输入管理员用户名和密码',
      editable: true,
      placeholderText: '用户名',
      success: (res) => {
        if (res.confirm && res.content) {
          const username = res.content;
          wx.showModal({
            title: '设置密码',
            content: '请输入管理员密码',
            editable: true,
            placeholderText: '密码',
            success: (res2) => {
              if (res2.confirm && res2.content) {
                const password = res2.content;
                this.setUserAsAdmin(user.userId, username, password);
              }
            }
          });
        }
      }
    });
  },

  // 设置用户为管理员
  setUserAsAdmin: function (userId, username, password) {
    const adminInfo = wx.getStorageSync('adminInfo');
    
    const app = getApp();
    app.callContainer('/api/admin/set-admin', 'POST', {
      userId: userId,
      adminLevel: 1, // 一级管理员
      parentAdminId: adminInfo.userId,
      adminUsername: username,
      adminPassword: password
    }).then(res => {
      if (res.code === 0) {
        wx.showToast({
          title: '设置成功',
          icon: 'success'
        });
        // 刷新管理员列表
        this.loadAdmins(true);
      } else {
        wx.showToast({
          title: res.errorMsg || '设置失败',
          icon: 'none'
        });
      }
    }).catch(err => {
      wx.showToast({
        title: '网络错误，请重试',
        icon: 'none'
      });
      console.error('设置管理员失败:', err);
    });
  },

  // 取消管理员权限
  onRemoveAdmin: function (e) {
    const userId = e.currentTarget.dataset.userId;
    const user = this.data.admins.find(u => u.userId === userId);
    
    if (!user) return;

    wx.showModal({
      title: '取消管理员权限',
      content: `确定要取消用户"${user.nickName}"的管理员权限吗？`,
      success: (res) => {
        if (res.confirm) {
          this.removeAdmin(userId);
        }
      }
    });
  },

  // 取消管理员权限
  removeAdmin: function (userId) {
    const app = getApp();
    app.callContainer('/api/admin/remove-admin', 'POST', {
      userId: userId
    }).then(res => {
      if (res.code === 0) {
        wx.showToast({
          title: '取消成功',
          icon: 'success'
        });
        // 刷新管理员列表
        this.loadAdmins(true);
      } else {
        wx.showToast({
          title: res.errorMsg || '取消失败',
          icon: 'none'
        });
      }
    }).catch(err => {
      wx.showToast({
        title: '网络错误，请重试',
        icon: 'none'
      });
      console.error('取消管理员权限失败:', err);
    });
  },

  // 获取管理员级别文本
  getAdminLevelText: function (level) {
    const levelMap = {
      0: '普通用户',
      1: '一级管理员',
      2: '超级管理员'
    };
    return levelMap[level] || '未知';
  },

  // 格式化时间
  formatTime: function (timeStr) {
    if (!timeStr) return '';
    
    try {
      const date = new Date(timeStr);
      
      // 检查日期是否有效
      if (isNaN(date.getTime())) {
        return timeStr;
      }
      
      const now = new Date();
      const diff = now - date;
      
      // 如果是今天
      if (diff < 24 * 60 * 60 * 1000 && date.toDateString() === now.toDateString()) {
        return '今天 ' + date.toLocaleTimeString('zh-CN', { 
          hour: '2-digit', 
          minute: '2-digit' 
        });
      }
      
      // 如果是昨天
      const yesterday = new Date(now);
      yesterday.setDate(yesterday.getDate() - 1);
      if (date.toDateString() === yesterday.toDateString()) {
        return '昨天 ' + date.toLocaleTimeString('zh-CN', { 
          hour: '2-digit', 
          minute: '2-digit' 
        });
      }
      
      // 其他时间
      return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
      });
      
    } catch (error) {
      return timeStr;
    }
  }
}); 