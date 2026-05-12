const { api } = require("../../utils/cloud-container-standard")
// pages/admin/users.js
Page({
  data: {
    adminInfo: null,
    users: [],
    loading: true,
    hasMore: true,
    page: 1,
    pageSize: 20
  },

  onLoad: function (options) {
    this.checkAdminAuth();
    this.loadUsers();
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

  // 加载用户列表
  loadUsers: function (refresh = false) {
    const adminInfo = wx.getStorageSync('adminInfo');
    if (!adminInfo) return;

    if (refresh) {
      this.setData({
        page: 1,
        users: [],
        hasMore: true
      });
    }

    if (!this.data.hasMore) return;

    const app = getApp();
    app.callContainer('/api/admin/users', 'GET', {
      adminUserId: adminInfo.userId,
      page: this.data.page,
      pageSize: this.data.pageSize
    }).then(res => {
      this.setData({ loading: false });
      console.log('用户返回：', res);
      if (res.code === 0) {
        const newUsers = res.data.list || [];
        const hasMore = res.data.hasMore;
        this.setData({
          users: refresh ? newUsers : [...this.data.users, ...newUsers],
          hasMore: hasMore,
          page: this.data.page + 1
        });
      } else {
        wx.showToast({
          title: res.errorMsg || '获取用户失败',
          icon: 'none'
        });
      }
    }).catch(err => {
      this.setData({ loading: false });
      wx.showToast({
        title: '网络错误，请重试',
        icon: 'none'
      });
      console.error('获取用户列表失败:', err);
    });
  },

  // 下拉刷新
  onPullDownRefresh: function () {
    this.loadUsers(true);
    wx.stopPullDownRefresh();
  },

  // 上拉加载更多
  onReachBottom: function () {
    if (this.data.hasMore && !this.data.loading) {
      this.loadUsers();
    }
  },

  // 设置用户为管理员
  onSetAsAdmin: function (e) {
    const userId = e.currentTarget.dataset.userId;
    const user = this.data.users.find(u => u.userId === userId);
    
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
        // 刷新用户列表
        this.loadUsers(true);
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
    const user = this.data.users.find(u => u.userId === userId);
    
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
        // 刷新用户列表
        this.loadUsers(true);
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
  }
}); 