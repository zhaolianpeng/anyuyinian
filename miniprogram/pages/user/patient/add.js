// pages/user/patient/add.js
const app = getApp()
const { api } = require('../../../utils/cloud-container-standard')

Page({
  data: {
    isEdit: false,
    patientId: null,
    userId: null,
    formData: {
      name: '',
      idCard: '',
      phone: '',
      gender: 1,
      birthday: '',
      relation: '',
      isDefault: false
    },
    genderOptions: ['男', '女'],
    genderIndex: 0,
    relationOptions: ['本人', '父亲', '母亲', '儿子', '女儿', '配偶', '其他'],
    relationIndex: 0,
    submitting: false
  },

  onLoad(options) {
    const userId = wx.getStorageSync('userId')
    this.setData({ userId })

    // 如果是编辑模式
    if (options.edit === '1' && options.patient) {
      const patient = JSON.parse(options.patient)
      this.setData({
        isEdit: true,
        patientId: patient.id,
        formData: {
          name: patient.name,
          idCard: patient.idCard,
          phone: patient.phone,
          gender: patient.gender,
          birthday: patient.birthday,
          relation: patient.relation,
          isDefault: patient.isDefault === 1
        },
        genderIndex: patient.gender === 1 ? 0 : 1,
        relationIndex: this.getRelationIndex(patient.relation)
      })
    }
  },

  // 获取关系索引
  getRelationIndex(relation) {
    const index = this.data.relationOptions.findIndex(item => item === relation)
    return index >= 0 ? index : 0
  },

  // 性别选择器变化
  genderChange(e) {
    const index = e.detail.value
    this.setData({
      genderIndex: index,
      'formData.gender': index === 0 ? 1 : 2
    })
  },

  // 出生日期选择器变化
  birthdayChange(e) {
    this.setData({
      'formData.birthday': e.detail.value
    })
  },

  // 关系选择器变化
  relationChange(e) {
    const index = e.detail.value
    this.setData({
      relationIndex: index,
      'formData.relation': this.data.relationOptions[index]
    })
  },

  // 默认就诊人开关变化
  switchChange(e) {
    this.setData({
      'formData.isDefault': e.detail.value
    })
  },

  // 表单提交
  submitForm(e) {
    const formData = e.detail.value
    
    // 合并表单数据，确保生日数据正确
    const submitData = {
      ...this.data.formData,
      ...formData,
      birthday: this.data.formData.birthday, // 确保使用正确的生日数据
      userId: this.data.userId
    }

    // 验证表单
    if (!this.validateForm(submitData)) {
      return
    }

    if (this.data.isEdit) {
      submitData.id = this.data.patientId
    }

    this.submitPatient(submitData)
  },

  // 验证表单
  validateForm(formData) {
    if (!formData.name.trim()) {
      wx.showToast({
        title: '请输入患者姓名',
        icon: 'none'
      })
      return false
    }

    if (!formData.idCard.trim()) {
      wx.showToast({
        title: '请输入身份证号',
        icon: 'none'
      })
      return false
    }

    // 身份证号验证
    const idCardRegex = /^[1-9]\d{5}(18|19|20)\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$/
    if (!idCardRegex.test(formData.idCard)) {
      wx.showToast({
        title: '请输入正确的身份证号',
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

    // 手机号验证
    const phoneRegex = /^1[3-9]\d{9}$/
    if (!phoneRegex.test(formData.phone)) {
      wx.showToast({
        title: '请输入正确的手机号',
        icon: 'none'
      })
      return false
    }

    if (!formData.birthday) {
      wx.showToast({
        title: '请选择出生日期',
        icon: 'none'
      })
      return false
    }

    return true
  },

  // 提交患者信息
  async submitPatient(data) {
    this.setData({ submitting: true })

    try {
      console.log('提交患者数据:', data)
      
      let result
      if (this.data.isEdit) {
        // 编辑模式
        result = await api.userPatientUpdate(data)
      } else {
        // 添加模式
        result = await api.userPatientAdd(data)
      }
      
      console.log('患者操作API返回:', result)
      
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
      console.error('提交患者信息失败:', error)
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