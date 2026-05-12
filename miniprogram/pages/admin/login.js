const { api } = require("../../utils/cloud-container-standard")
// pages/admin/login.js
Page({
  data: {
    // 超级管理员登录
    superUsername: '',
    superPassword: '',
    // 普通管理员登录
    isNormalLogin: false,
    loading: false
  },

  onLoad: function (options) {
    // 检查是否已经有管理员登录状态
    this.checkAdminStatus();
  },

  onShow: function () {
    this.checkAdminStatus();
  },

  // 检查管理员状态
  checkAdminStatus: function () {
    const adminInfo = wx.getStorageSync('adminInfo');
    const isAdmin = wx.getStorageSync('isAdmin');
    
    if (adminInfo && isAdmin) {
      // 已经是管理员登录状态，直接跳转到管理员首页
      wx.switchTab({
        url: '/pages/admin/home'
      });
    }
  },

  // 输入超级管理员用户名
  onSuperUsernameInput: function (e) {
    this.setData({
      superUsername: e.detail.value
    });
  },

  // 输入超级管理员密码
  onSuperPasswordInput: function (e) {
    this.setData({
      superPassword: e.detail.value
    });
  },

  // 超级管理员登录
  onSuperAdminLogin: function () {
    const { superUsername, superPassword } = this.data;
    
    if (!superUsername || !superPassword) {
      wx.showToast({
        title: '请输入用户名和密码',
        icon: 'none'
      });
      return;
    }

    // 超级管理员只能登录anyuyinian账号
    if (superUsername !== 'anyuyinian') {
      wx.showToast({
        title: '超级管理员账号错误',
        icon: 'none'
      });
      return;
    }

    this.setData({ loading: true });

    // 统一调用app.js中的callContainer方法
    const app = getApp();
    app.callContainer('/api/admin/login', 'POST', {
      username: superUsername,
      password: superPassword
    }).then(res => {
      this.setData({ loading: false });
      console.log('超级管理员登录返回：', res);
      if (res.code === 0) {
        wx.setStorageSync('adminInfo', res.data);
        wx.setStorageSync('isAdmin', true);
        wx.setStorageSync('isSuperAdmin', true);
        wx.showToast({
          title: '超级管理员登录成功',
          icon: 'success'
        });
        setTimeout(() => {
          wx.redirectTo({
            url: '/pages/admin/home'
          });
        }, 1500);
      } else {
        wx.showModal({
          title: '登录失败',
          content: res.message || res.errorMsg || JSON.stringify(res) || '登录失败',
          showCancel: false
        });
      }
    }).catch(err => {
      this.setData({ loading: false });
      console.error('超级管理员登录失败：', err);
      wx.showModal({
        title: '网络错误',
        content: (err && err.errMsg ? err.errMsg : '') + '\n' + JSON.stringify(err),
        showCancel: false
      });
    });
  },

  // 普通管理员登录
  onNormalAdminLogin: function () {
    // 检查当前用户是否已登录
    const userId = wx.getStorageSync('userId');
    
    if (userId) {
      // 用户已登录，检查是否为管理员
      this.checkUserAdminStatus(userId);
    } else {
      // 用户未登录，调起微信登录
      this.wxLogin();
    }
  },

  // 微信登录
  wxLogin: function () {
    wx.showLoading({
      title: '登录中...'
    });

    wx.login({
      success: async (res) => {
        if (!res.code) {
          wx.hideLoading();
          wx.showToast({
            title: '微信登录失败',
            icon: 'none'
          });
          return;
        }

        try {
          const loginRes = await api.wxLogin({
            code: res.code
          });
          wx.hideLoading();

          if (loginRes.code === 0) {
            wx.setStorageSync('userId', loginRes.data.userId);
            wx.setStorageSync('userInfo', loginRes.data);
            this.checkUserAdminStatus(loginRes.data.userId);
          } else {
            wx.showToast({
              title: '微信登录失败',
              icon: 'none'
            });
          }
        } catch (err) {
          wx.hideLoading();
          wx.showToast({
            title: '网络错误，请重试',
            icon: 'none'
          });
          console.error('微信登录失败:', err);
        }
      },
      fail: (err) => {
        wx.hideLoading();
        wx.showToast({
          title: '微信登录失败',
          icon: 'none'
        });
        console.error('微信登录失败:', err);
      }
    });
  },

  // 检查用户管理员状态
  checkUserAdminStatus: function (userId) {
    const app = getApp();
    app.callContainer('/api/admin/check-status', 'GET', {
      userId: userId
    }).then((res) => {
      if (res.code === 0) {
        if (res.data.isAdmin) {
          wx.setStorageSync('adminInfo', res.data.adminInfo);
          wx.setStorageSync('isAdmin', true);
          wx.setStorageSync('isSuperAdmin', false);

          wx.showToast({
            title: '管理员登录成功',
            icon: 'success'
          });

          setTimeout(() => {
            wx.switchTab({
              url: '/pages/admin/home'
            });
          }, 1500);
        } else {
          wx.showModal({
            title: '提示',
            content: '当前用户非管理员账号，无法使用管理员功能',
            showCancel: false,
            success: () => {
              wx.switchTab({
                url: '/pages/index/index'
              });
            }
          });
        }
      } else {
        wx.showToast({
          title: res.errorMsg || '检查管理员状态失败',
          icon: 'none'
        });
      }
    }).catch((err) => {
      wx.showToast({
        title: '网络错误，请重试',
        icon: 'none'
      });
      console.error('检查管理员状态失败:', err);
    });
  }
  ,
  // 取消/返回
  onCancelTap: function () {
    if (this.data.loading) return;
    wx.navigateBack({
      delta: 1,
      fail: () => {
        wx.switchTab({
          url: '/pages/index/index'
        });
      }
    });
  }
}); 