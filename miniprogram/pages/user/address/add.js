// pages/user/address/add.js
const app = getApp()
const { api } = require('../../../utils/cloud-container-standard')

Page({
  data: {
    isEdit: false,
    addressId: null,
    userId: null,
    formData: {
      name: '',
      phone: '',
      province: '',
      city: '',
      district: '',
      address: '',
      isDefault: false
    },
    region: [],
    submitting: false
  },

  onLoad(options) {
    const userId = wx.getStorageSync('userId')
    this.setData({ userId })

    // 如果是编辑模式
    if (options.edit === '1' && options.address) {
      const address = JSON.parse(options.address)
      this.setData({
        isEdit: true,
        addressId: address.id,
        formData: {
          name: address.name,
          phone: address.phone,
          province: address.province,
          city: address.city,
          district: address.district,
          address: address.address,
          isDefault: address.isDefault === 1
        },
        region: [address.province, address.city, address.district]
      })
    }
  },

  // 地区选择器变化
  regionChange(e) {
    const region = e.detail.value
    this.setData({
      region,
      'formData.province': region[0],
      'formData.city': region[1],
      'formData.district': region[2]
    })
  },

  // 默认地址开关变化
  switchChange(e) {
    this.setData({
      'formData.isDefault': e.detail.value
    })
  },

  // 表单提交
  submitForm(e) {
    const formData = e.detail.value
    
    // 验证表单
    if (!this.validateForm(formData)) {
      return
    }

    // 合并表单数据
    const submitData = {
      ...this.data.formData,
      ...formData,
      userId: this.data.userId
    }

    if (this.data.isEdit) {
      submitData.id = this.data.addressId
    }

    this.submitAddress(submitData)
  },

  // 验证表单
  validateForm(formData) {
    if (!formData.name.trim()) {
      wx.showToast({
        title: '请输入联系人姓名',
        icon: 'none'
      })
      return false
    }

    if (!formData.phone.trim()) {
      wx.showToast({
        title: '请输入手机号',
        icon: 'none'
      })
      return false
    }

    // 简单的手机号验证
    const phoneRegex = /^1[3-9]\d{9}$/
    if (!phoneRegex.test(formData.phone)) {
      wx.showToast({
        title: '请输入正确的手机号',
        icon: 'none'
      })
      return false
    }

    if (!this.data.region[0]) {
      wx.showToast({
        title: '请选择所在地区',
        icon: 'none'
      })
      return false
    }

    if (!formData.address.trim()) {
      wx.showToast({
        title: '请输入详细地址',
        icon: 'none'
      })
      return false
    }

    return true
  },

  // 提交地址
  async submitAddress(data) {
    this.setData({ submitting: true })

    try {
      console.log('提交地址数据:', data)
      
      let result
      if (this.data.isEdit) {
        // 编辑模式
        result = await api.userAddressUpdate(data)
      } else {
        // 添加模式
        result = await api.userAddressAdd(data)
      }
      
      console.log('地址操作API返回:', result)
      
      if (result.code === 0) {
        wx.showToast({
          title: this.data.isEdit ? '修改成功' : '添加成功',
          icon: 'success'
        })
        
        // 返回上一页并刷新
        setTimeout(() => {
          wx.navigateBack()
        }, 1500)
      } else {
        throw new Error(result.errorMsg || result.message || '操作失败')
      }
    } catch (error) {
      console.error('提交地址失败:', error)
      wx.showToast({
        title: error.message || '网络错误，请重试',
        icon: 'none',
        duration: 3000
      })
    } finally {
      this.setData({ submitting: false })
    }
  }
})