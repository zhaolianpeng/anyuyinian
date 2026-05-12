// pages/order/order.js
const { api } = require('../../utils/cloud-container-standard')
const { processImageUrl } = require('../../utils/image')
const { getCurrentUserId, needsUserIdMigration, clearUserId, safeClearUserId } = require('../../utils/user-id-compatibility')

Page({
  data: {
    serviceId: '',
    serviceInfo: null,
    patientList: [],
    addressList: [],
    selectedPatient: null,
    selectedAddress: null,
    appointmentDate: '',
    appointmentTime: '',
    remark: '',
    disease: '',
    toiletAssist: false,
    loading: false,
    submitting: false,
    showCalendar: false,
    appointmentDateRange: {
      start: '',
      end: ''
    }
  },

  onLoad(options) {
    console.log('订单页面加载，参数:', options)
    
    if (options.serviceId) {
      this.setData({ serviceId: options.serviceId })
      this.loadServiceInfo()
    }
    
    this.setAppointmentDateRange()
  },

  onShow() {
    // 每次显示页面时重新加载用户数据
    this.loadUserData()
  },

  // 设置预约日期范围（从明天开始，往后推7天）
  setAppointmentDateRange() {
    const today = new Date()
    const startDate = new Date(today.getTime() + 24 * 60 * 60 * 1000) // 明天
    const endDate = new Date(today.getTime() + 7 * 24 * 60 * 60 * 1000) // 明天+6天，共7天
    
    const formatDate = (date) => {
      const year = date.getFullYear()
      const month = String(date.getMonth() + 1).padStart(2, '0')
      const day = String(date.getDate()).padStart(2, '0')
      return `${year}-${month}-${day}`
    }
    
    console.log('设置预约时间范围:', {
      start: formatDate(startDate),
      end: formatDate(endDate)
    })
    
    this.setData({
      appointmentDateRange: {
        start: formatDate(startDate),
        end: formatDate(endDate)
      }
    })
  },

  // 加载服务信息
  async loadServiceInfo() {
    try {
      this.setData({ loading: true })
      
      const { serviceId } = this.data
      if (!serviceId) {
        throw new Error('缺少服务ID')
      }

      console.log('开始加载服务信息，serviceId:', serviceId)
      const result = await api.serviceDetail({ serviceId })
      
      if (result.code === 0) {
        const serviceInfo = {
          ...result.data,
          imageUrl: processImageUrl(result.data?.imageUrl),
          videoUrl: processImageUrl(result.data?.videoUrl)
        }

        this.setData({ 
          serviceInfo,
          loading: false
        })
        console.log('服务信息加载成功:', serviceInfo)
      } else {
        throw new Error(result.errorMsg || '加载服务信息失败')
      }
    } catch (error) {
      console.error('加载服务信息失败:', error)
      this.setData({ loading: false })
      
      wx.showToast({
        title: error.message || '加载服务信息失败',
        icon: 'none'
      })
    }
  },

  // 加载用户数据
  async loadUserData() {
    try {
      const userId = getCurrentUserId()
      console.log('=== 开始加载用户数据 ===')
      console.log('用户ID:', userId)
      console.log('用户ID类型:', typeof userId)
      
      if (!userId) {
        console.error('❌ 用户未登录，跳转到登录页面')
        wx.navigateTo({ url: '/pages/login/login' })
        return
      }

      // 检查是否需要迁移用户ID
      if (needsUserIdMigration()) {
        console.log('检测到旧格式的用户ID，需要重新登录')
        wx.showModal({
          title: '系统升级',
          content: '系统已升级，需要重新登录以获取新的用户ID',
          showCancel: false,
          success: () => {
            clearUserId()
            wx.navigateTo({ url: '/pages/login/login' })
          }
        })
        return
      }

      console.log('✅ 用户已登录，开始加载数据')

      // 并行加载就诊人和地址列表
      console.log('开始并行请求就诊人和地址数据...')
      const [patientRes, addressRes] = await Promise.all([
        this.getPatientList(userId),
        this.getAddressList(userId)
      ])

      console.log('=== 就诊人API响应 ===')
      console.log('就诊人API响应:', patientRes)
      console.log('就诊人响应code:', patientRes.code)
      console.log('就诊人响应data:', patientRes.data)
      console.log('就诊人响应data类型:', typeof patientRes.data)
      console.log('就诊人响应data是否为数组:', Array.isArray(patientRes.data))
      
      if (patientRes.code === 0) {
        // 确保数据是数组格式
        let patientList = patientRes.data || []
        if (!Array.isArray(patientList)) {
          console.warn('⚠️ 就诊人数据不是数组，转换为数组')
          patientList = []
        }
        console.log('✅ 就诊人数据加载成功')
        console.log('就诊人列表:', patientList)
        console.log('就诊人数量:', patientList.length)
        
        if (patientList.length > 0) {
          console.log('第一个就诊人详情:', patientList[0])
          console.log('就诊人字段检查:')
          console.log('- id:', patientList[0].id)
          console.log('- name:', patientList[0].name)
          console.log('- relation:', patientList[0].relation)
          console.log('- phone:', patientList[0].phone)
          console.log('- isDefault:', patientList[0].isDefault)
        }
        
        const selectedPatient = patientList.find(p => p.isDefault === 1) || (patientList.length > 0 ? patientList[0] : null)
        console.log('选中的就诊人:', selectedPatient)
        
        this.setData({ 
          patientList,
          selectedPatient
        })
        console.log('✅ 就诊人数据设置完成')
        console.log('设置后的就诊人列表:', this.data.patientList)
        console.log('设置后的选中就诊人:', this.data.selectedPatient)
      } else {
        console.error('❌ 就诊人API错误:', patientRes)
        this.setData({ patientList: [] })
      }

      console.log('=== 地址API响应 ===')
      console.log('地址API响应:', addressRes)
      console.log('地址响应code:', addressRes.code)
      console.log('地址响应data:', addressRes.data)
      console.log('地址响应data类型:', typeof addressRes.data)
      console.log('地址响应data是否为数组:', Array.isArray(addressRes.data))
      
      if (addressRes.code === 0) {
        // 确保数据是数组格式
        let addressList = addressRes.data || []
        if (!Array.isArray(addressList)) {
          console.warn('⚠️ 地址数据不是数组，转换为数组')
          addressList = []
        }
        console.log('✅ 地址数据加载成功')
        console.log('地址列表:', addressList)
        console.log('地址数量:', addressList.length)
        
        if (addressList.length > 0) {
          console.log('第一个地址详情:', addressList[0])
          console.log('地址字段检查:')
          console.log('- id:', addressList[0].id)
          console.log('- name:', addressList[0].name)
          console.log('- phone:', addressList[0].phone)
          console.log('- province:', addressList[0].province)
          console.log('- city:', addressList[0].city)
          console.log('- district:', addressList[0].district)
          console.log('- address:', addressList[0].address)
          console.log('- isDefault:', addressList[0].isDefault)
        }
        
        const selectedAddress = addressList.find(a => a.isDefault === 1) || (addressList.length > 0 ? addressList[0] : null)
        console.log('选中的地址:', selectedAddress)
        
        this.setData({ 
          addressList,
          selectedAddress
        })
        console.log('✅ 地址数据设置完成')
        console.log('设置后的地址列表:', this.data.addressList)
        console.log('设置后的选中地址:', this.data.selectedAddress)
      } else {
        console.error('❌ 地址API错误:', addressRes)
        this.setData({ addressList: [] })
      }

      // 强制更新页面显示
      console.log('=== 强制更新页面 ===')
      this.forceUpdate()
      this.checkCanSubmit()
      
      console.log('=== 数据加载完成 ===')
      console.log('最终页面数据状态:')
      console.log('- 就诊人列表长度:', this.data.patientList?.length || 0)
      console.log('- 地址列表长度:', this.data.addressList?.length || 0)
      console.log('- 选中的就诊人:', this.data.selectedPatient)
      console.log('- 选中的地址:', this.data.selectedAddress)
      
    } catch (error) {
      console.error('❌ 加载用户数据失败:', error)
    }
  },

  // 强制更新页面
  forceUpdate() {
    console.log('强制更新页面数据')
    this.setData({
      patientList: this.data.patientList,
      addressList: this.data.addressList,
      selectedPatient: this.data.selectedPatient,
      selectedAddress: this.data.selectedAddress
    })
  },

  // 获取服务详情
  getServiceDetail(serviceId) {
    return new Promise((resolve, reject) => {
      wx.request({
        url: `${app.globalData.baseUrl}/api/service/detail`,
        method: 'POST',
        data: { serviceId },
        success: res => resolve(res.data),
        fail: reject
      })
    })
  },

  // 获取就诊人列表
  getPatientList(userId) {
    console.log('开始获取就诊人列表，userId:', userId)
    
    // 确保userId是字符串类型
    const stringUserId = String(userId)
    console.log('就诊人API调用，用户ID:', stringUserId, '类型:', typeof stringUserId)
    
    return api.userPatient({ 
      userId: stringUserId 
    }).catch(error => {
      console.error('获取就诊人列表失败:', error)
      return { code: -1, errorMsg: error.message }
    })
  },

  // 获取地址列表
  getAddressList(userId) {
    console.log('开始获取地址列表，userId:', userId)
    
    // 确保userId是字符串类型
    const stringUserId = String(userId)
    console.log('地址API调用，用户ID:', stringUserId, '类型:', typeof stringUserId)
    
    return api.userAddress({ 
      userId: stringUserId 
    }).catch(error => {
      console.error('获取地址列表失败:', error)
      return { code: -1, errorMsg: error.message }
    })
  },

  // 选择就诊人
  selectPatient(e) {
    const { patient } = e.currentTarget.dataset
    console.log('选择就诊人:', patient)
    this.setData({ selectedPatient: patient })
    this.checkCanSubmit()
  },

  // 选择地址
  selectAddress(e) {
    const { address } = e.currentTarget.dataset
    console.log('选择地址:', address)
    this.setData({ selectedAddress: address })
    this.checkCanSubmit()
  },

  // 日期选择
  onDateChange(e) {
    console.log('选择日期:', e.detail.value)
    this.setData({
      'formData.appointmentDate': e.detail.value
    })
    this.checkCanSubmit()
  },

  // 显示日历选择器
  showCalendarPicker() {
    console.log('=== 显示日历选择器 ===')
    console.log('当前 showCalendarPicker 值:', this.data.showCalendarPicker)
    console.log('当前 minDate:', this.data.minDate)
    console.log('当前 maxDate:', this.data.maxDate)
    
    this.setData({
      showCalendarPicker: true
    })
    
    console.log('设置 showCalendarPicker 为 true')
    console.log('设置后的 showCalendarPicker 值:', this.data.showCalendarPicker)
  },

  // 日历选择器确认
  onCalendarConfirm(e) {
    const { date, time } = e.detail
    console.log('日历选择确认:', { date, time })
    this.setData({
      'formData.appointmentDate': date,
      'formData.appointmentTime': time,
      showCalendarPicker: false
    })
    this.checkCanSubmit()
  },

  // 日历选择器关闭
  onCalendarClose() {
    console.log('关闭日历选择器')
    this.setData({
      showCalendarPicker: false
    })
  },

  // 备注输入
  onRemarkInput(e) {
    console.log('备注输入:', e.detail.value)
    this.setData({
      'formData.remark': e.detail.value
    })
  },

  // 备注失去焦点
  onRemarkBlur(e) {
    console.log('备注失去焦点:', e.detail.value)
    this.setData({
      'formData.remark': e.detail.value
    })
  },

  // 基础病信息输入
  onDiseaseInput(e) {
    console.log('基础病信息输入:', e.detail.value)
    this.setData({
      'formData.diseaseInfo': e.detail.value
    })
  },

  // 助排二便选择
  onToiletAssistChange(e) {
    const value = e.currentTarget.dataset.value
    console.log('助排二便选择:', value)
    this.setData({
      'formData.needToiletAssist': value
    })
  },

  // 计算年龄
  getAge(birthday) {
    if (!birthday) return ''
    const birthDate = new Date(birthday)
    const today = new Date()
    let age = today.getFullYear() - birthDate.getFullYear()
    const monthDiff = today.getMonth() - birthDate.getMonth()
    
    if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birthDate.getDate())) {
      age--
    }
    
    return age
  },

  // 检查是否可以提交
  checkCanSubmit() {
    const canSubmit = this.data.selectedPatient && 
                     this.data.selectedAddress && 
                     this.data.formData.appointmentDate && 
                     this.data.formData.appointmentTime
    
    console.log('检查提交状态:', {
      selectedPatient: !!this.data.selectedPatient,
      selectedAddress: !!this.data.selectedAddress,
      appointmentDate: !!this.data.formData.appointmentDate,
      appointmentTime: !!this.data.formData.appointmentTime,
      canSubmit
    })
    
    this.setData({ canSubmit })
  },

  // 添加就诊人
  addPatient() {
    wx.navigateTo({
      url: '/pages/user/patient/add'
    })
  },

  // 添加地址
  addAddress() {
    wx.navigateTo({
      url: '/pages/user/address/add'
    })
  },

  // 调试数据
  debugData() {
    console.log('=== 调试页面数据 ===')
    console.log('当前页面数据:', this.data)
    console.log('患者列表:', this.data.patientList)
    console.log('地址列表:', this.data.addressList)
    console.log('选中的患者:', this.data.selectedPatient)
    console.log('选中的地址:', this.data.selectedAddress)
    console.log('预约时间范围:', {
      minDate: this.data.minDate,
      maxDate: this.data.maxDate
    })
    console.log('表单数据:', this.data.formData)
    
    wx.showModal({
      title: '调试信息',
      content: `患者数量: ${this.data.patientList?.length || 0}\n地址数量: ${this.data.addressList?.length || 0}\n预约时间范围: ${this.data.minDate} 至 ${this.data.maxDate}`,
      showCancel: false
    })
  },



  // 提交订单
  async submitOrder() {
    try {
      this.setData({ submitting: true })

      // 验证必填项
      if (!this.validateForm()) {
        return
      }

      const userId = getCurrentUserId()
      if (!userId) {
        wx.showToast({
          title: '请先登录',
          icon: 'none'
        })
        return
      }

      const orderData = {
        userId,
        serviceId: this.data.serviceId,
        patientId: this.data.selectedPatient.id,
        addressId: this.data.selectedAddress.id,
        appointmentDate: this.data.appointmentDate,
        appointmentTime: this.data.appointmentTime,
        remark: this.data.remark,
        diseaseInfo: this.data.disease,
        needToiletAssist: this.data.toiletAssist
      }

      console.log('提交订单数据:', orderData)
      const res = await this.submitOrderToServer(orderData)
      
      if (res.code === 0) {
        wx.showToast({
          title: '订单提交成功',
          icon: 'success'
        })

        // 跳转到支付页面
        setTimeout(() => {
          this.payOrder(res.data.orderId)
        }, 1500)
      } else {
        wx.showToast({
          title: res.errorMsg || '提交失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('提交订单失败:', error)
      wx.showToast({
        title: '网络错误，请重试',
        icon: 'none'
      })
    } finally {
      this.setData({ submitting: false })
    }
  },

  // 验证表单
  validateForm() {
    if (!this.data.selectedPatient) {
      wx.showToast({
        title: '请选择患者',
        icon: 'none'
      })
      return false
    }

    if (!this.data.selectedAddress) {
      wx.showToast({
        title: '请选择服务地址',
        icon: 'none'
      })
      return false
    }

    if (!this.data.appointmentDate) {
      wx.showToast({
        title: '请选择预约日期',
        icon: 'none'
      })
      return false
    }

    if (!this.data.appointmentTime) {
      wx.showToast({
        title: '请选择预约时间',
        icon: 'none'
      })
      return false
    }

    // 验证预约时间是否在允许范围内
    const selectedDateTime = new Date(`${this.data.appointmentDate} ${this.data.appointmentTime}`)
    const tomorrow = new Date()
    tomorrow.setDate(tomorrow.getDate() + 1)
    tomorrow.setHours(0, 0, 0, 0)

    if (selectedDateTime < tomorrow) {
      wx.showToast({
        title: '预约时间不能早于明天',
        icon: 'none'
      })
      return false
    }

    return true
  },

  // 提交订单到服务器
  submitOrderToServer(orderData) {
    return new Promise((resolve, reject) => {
      wx.request({
        url: `${app.globalData.baseUrl}/api/order/submit`,
        method: 'POST',
        data: orderData,
        success: res => resolve(res.data),
        fail: reject
      })
    })
  },

  // 支付订单 - 使用云开发微信支付
  async payOrder(orderId) {
    try {
      // 显示支付加载状态
      wx.showLoading({
        title: '正在调起支付...',
        mask: true
      })

      // 获取用户信息
      const userInfo = wx.getStorageSync('userInfo')
      if (!userInfo || !userInfo.openId) {
        throw new Error('用户信息不完整，请重新登录')
      }

      // 获取订单信息
      const orderInfo = await this.getOrderInfo(orderId)
      if (!orderInfo) {
        throw new Error('订单信息不存在')
      }

      // 构建支付参数
      const orderData = {
        out_trade_no: orderInfo.orderNo,
        body: orderInfo.serviceName || '安语一年服务',
        total_fee: Math.round(orderInfo.totalAmount * 100), // 转换为分
        spbill_create_ip: '127.0.0.1',
        openid: userInfo.openId,
        attach: `order_${orderId}`,
        detail: orderInfo.serviceName || '安语一年服务',
        goods_tag: orderInfo.serviceCategory || '',
        time_start: this.formatTime(new Date()),
        time_expire: this.formatTime(new Date(Date.now() + 30 * 60 * 1000)) // 30分钟后过期
      }

      console.log('调用云开发云函数统一下单:', orderData)
      
      // 调用云开发云函数统一下单
      const result = await wx.cloud.callFunction({
        name: 'payOrder',
        data: {
          orderInfo: orderData
        }
      })

      console.log('云开发云函数返回结果:', result)

      if (result.result && result.result.success) {
        const payment = result.result.data.payment
        console.log('获取到云开发支付参数:', payment)
        
        // 调用微信支付
        wx.requestPayment({
          timeStamp: payment.timeStamp,
          nonceStr: payment.nonceStr,
          package: payment.package,
          signType: payment.signType || 'MD5',
          paySign: payment.paySign,
          success: (res) => {
            console.log('支付成功:', res)
            wx.hideLoading()
            wx.showToast({
              title: '支付成功',
              icon: 'success',
              duration: 2000
            })
            
            // 跳转到订单详情页
            setTimeout(() => {
              wx.navigateTo({
                url: `/pages/order/detail?id=${orderId}`
              })
            }, 2000)
          },
          fail: (err) => {
            console.error('支付失败:', err)
            wx.hideLoading()
            
            // 根据错误类型显示不同提示
            let errorMsg = '支付失败'
            if (err.errMsg) {
              if (err.errMsg.includes('cancel')) {
                errorMsg = '支付已取消'
              } else if (err.errMsg.includes('fail')) {
                errorMsg = '支付失败，请重试'
              }
            }
            
            wx.showToast({
              title: errorMsg,
              icon: 'none',
              duration: 3000
            })
          }
        })
      } else {
        throw new Error(result.result?.error || '获取支付参数失败')
      }
    } catch (error) {
      console.error('云开发支付失败:', error)
      wx.hideLoading()
      wx.showToast({
        title: error.message || '支付失败，请重试',
        icon: 'none',
        duration: 3000
      })
    }
  },

  // 获取订单信息
  async getOrderInfo(orderId) {
    try {
      // 这里应该调用后端API获取订单详情
      // 暂时返回模拟数据
      return {
        orderNo: `ORDER_${Date.now()}`,
        serviceName: this.data.service?.name || '安语一年服务',
        serviceCategory: this.data.service?.category || '',
        totalAmount: this.data.service?.price || 0
      }
    } catch (error) {
      console.error('获取订单信息失败:', error)
      return null
    }
  },

  // 格式化时间
  formatTime(date) {
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const hours = String(date.getHours()).padStart(2, '0')
    const minutes = String(date.getMinutes()).padStart(2, '0')
    const seconds = String(date.getSeconds()).padStart(2, '0')
    return `${year}${month}${day}${hours}${minutes}${seconds}`
  },

  // 支付订单到服务器
  payOrderToServer(paymentData) {
    return new Promise((resolve, reject) => {
      wx.request({
        url: `${app.globalData.baseUrl}/api/order/pay`,
        method: 'POST',
        data: paymentData,
        success: res => resolve(res.data),
        fail: reject
      })
    })
  }
})