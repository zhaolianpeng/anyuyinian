// pages/service/detail.js
const { api } = require('../../utils/cloud-container-standard')
const { processImageUrl } = require('../../utils/image')
const { getSingleServiceImage } = require('../../utils/imageService')
const { fetchServerImages, normalizeServerImagePath, fetchServerVideoToLocal } = require('../../utils/serverMedia')

function normalizeFormField(field = {}) {
  const normalizedName = field.name || field.fieldName || field.key || ''
  const normalizedType = field.type === 'phone' ? 'text' : (field.type || 'text')
  const normalizedLabel = field.label || field.title || normalizedName
  const normalizedPlaceholder = field.placeholder || field.tips || ''
  const rawOptions = Array.isArray(field.options) ? field.options : []
  const normalizedOptions = rawOptions.map(option => {
    if (typeof option === 'string' || typeof option === 'number') {
      const value = String(option)
      return {
        label: value,
        value
      }
    }

    return {
      label: option.label || option.text || option.name || String(option.value || ''),
      value: String(option.value != null ? option.value : (option.label || option.text || option.name || ''))
    }
  })

  return {
    ...field,
    name: normalizedName,
    label: normalizedLabel,
    type: normalizedType,
    placeholder: normalizedPlaceholder,
    options: normalizedOptions
  }
}

function normalizeFormFields(fields = []) {
  return fields
    .filter(field => field && (field.name || field.fieldName || field.key))
    .map(normalizeFormField)
}

async function localizeServiceMedia(service = {}) {
  const localizedService = { ...service }
  const imageCandidates = [localizedService.imageTempUrl, localizedService.imageUrl, localizedService.imageCosId]
  const mainImagePath = normalizeServerImagePath(imageCandidates.find(Boolean))

  if (mainImagePath) {
    try {
      const imageMap = await fetchServerImages([mainImagePath])
      localizedService.imageUrl = imageMap[mainImagePath] || localizedService.imageTempUrl || processImageUrl(localizedService.imageUrl)
    } catch (error) {
      console.warn('服务详情主图本地化失败，使用原图地址:', error)
      localizedService.imageUrl = localizedService.imageTempUrl || processImageUrl(localizedService.imageUrl)
    }
  } else {
    localizedService.imageUrl = localizedService.imageTempUrl || processImageUrl(localizedService.imageUrl)
  }

  if (localizedService.detailImages) {
    try {
      const rawDetailImages = Array.isArray(localizedService.detailImages)
        ? localizedService.detailImages
        : JSON.parse(localizedService.detailImages)
      const detailPaths = rawDetailImages.map(imageUrl => normalizeServerImagePath(imageUrl)).filter(Boolean)
      const detailImageMap = detailPaths.length > 0 ? await fetchServerImages(detailPaths) : {}

      localizedService.detailImages = rawDetailImages.map(imageUrl => {
        const normalizedPath = normalizeServerImagePath(imageUrl)
        return detailImageMap[normalizedPath] || processImageUrl(imageUrl)
      })
    } catch (error) {
      console.error('解析或本地化服务详情图片失败:', error)
      localizedService.detailImages = []
    }
  }

  try {
    localizedService.videoUrl = await fetchServerVideoToLocal(localizedService.videoUrl)
  } catch (error) {
    console.warn('服务详情视频本地化失败，回退原地址:', error)
    localizedService.videoUrl = processImageUrl(localizedService.videoUrl)
  }

  return localizedService
}

Page({
  data: {
    serviceId: null,
    service: {},
    formFields: [],
    formData: {},
    selectedDateTime: '',
    selectedDate: '',
    selectedTime: '',
    consultDateTime: '',  // 咨询时间
    consultDate: '',      // 咨询日期
    consultTime: '',      // 咨询时间
    remark: '',
    loading: false,
    userInfo: null,
    patientList: [],
    addressList: [],
    selectedPatient: null,
    selectedAddress: null,
    showPatientModal: false,
    showAddressModal: false,
    showTimePicker: false,
    showConsultTimePicker: false,
    minDate: '',
    maxDate: '',
    consultMinDate: '',
    consultMaxDate: ''
  },

  onLoad(options) {
    const { id } = options
    if (id) {
      const numId = Number(id)
      this.setData({ serviceId: numId })
      this.loadServiceDetail(numId)
    }
    
    // 获取用户信息
    const userInfo = wx.getStorageSync('userInfo')
    this.setData({ userInfo })
    
    // 初始化日期范围（明天到7天后）
    this.initDateRange()
    
    // 初始化咨询时间日期范围（今天到7天后）
    this.initConsultDateRange()
    
    // 加载用户数据
    this.loadUserData()
    
    // 调试信息：显示初始状态
    console.log('页面加载完成，初始数据状态:', {
      selectedDateTime: this.data.selectedDateTime,
      selectedDate: this.data.selectedDate,
      selectedTime: this.data.selectedTime,
      consultDateTime: this.data.consultDateTime,
      consultDate: this.data.consultDate,
      consultTime: this.data.consultTime
    })
  },

  onShow() {
    // 页面显示时刷新患者和地址列表
    this.loadUserData()
  },

  // 初始化预约时间日期范围（从明天开始，往后推7天）
  initDateRange() {
    const today = new Date()
    const tomorrow = new Date(today)
    tomorrow.setDate(today.getDate() + 1)
    
    const maxDate = new Date(today)
    maxDate.setDate(today.getDate() + 7) // 从明天开始，共7天（明天+6天）
    
    console.log('设置预约时间范围:', {
      start: this.formatDate(tomorrow),
      end: this.formatDate(maxDate)
    })
    
    this.setData({
      minDate: this.formatDate(tomorrow),
      maxDate: this.formatDate(maxDate)
    })
  },

  // 初始化咨询时间日期范围（从当前时间开始，往后推7天）
  initConsultDateRange() {
    const today = new Date()
    // 设置时间为00:00:00，避免时区问题
    today.setHours(0, 0, 0, 0)
    
    const maxDate = new Date(today)
    maxDate.setDate(today.getDate() + 6) // 从今天开始，共7天（今天+6天）
    
    console.log('设置咨询时间范围:', {
      start: this.formatDate(today),
      end: this.formatDate(maxDate),
      today: today.toISOString(),
      maxDate: maxDate.toISOString()
    })
    
    this.setData({
      consultMinDate: this.formatDate(today),
      consultMaxDate: this.formatDate(maxDate)
    })
  },

  // 格式化日期
  formatDate(date) {
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    return `${year}-${month}-${day}`
  },

  // 加载用户数据
  async loadUserData() {
    try {
      const userId = wx.getStorageSync('userId')
      console.log('加载用户数据，userId:', userId, '类型:', typeof userId)
      
      if (!userId) {
        console.log('用户未登录，无法加载用户数据')
        return
      }
      
      // 确保userId是字符串格式
      const strUserId = String(userId)
      if (!strUserId || strUserId === 'undefined' || strUserId === 'null') {
        console.error('userId无效:', userId)
        return
      }

      // 并行加载就诊人和地址列表
      const [patientList, addressList] = await Promise.all([
        this.getPatientList(strUserId),
        this.getAddressList(strUserId)
      ])

      console.log('获取到的患者列表:', patientList)
      console.log('获取到的地址列表:', addressList)

      this.setData({
        patientList: patientList || [],
        addressList: addressList || []
      })

      console.log('设置后的 patientList:', this.data.patientList)
      console.log('设置后的 addressList:', this.data.addressList)

    } catch (error) {
      console.error('加载用户数据失败:', error)
    }
  },

  // 获取就诊人列表
  async getPatientList(userId) {
    try {
      const result = await api.userPatient({ userId })
      return result.code === 0 ? result.data : []
    } catch (error) {
      console.error('获取就诊人列表失败:', error)
      return []
    }
  },

  // 获取地址列表
  async getAddressList(userId) {
    try {
      const result = await api.userAddress({ userId })
      return result.code === 0 ? result.data : []
    } catch (error) {
      console.error('获取地址列表失败:', error)
      return []
    }
  },

  // 加载服务详情
  async loadServiceDetail(serviceId) {
    try {
      this.setData({ loading: true })
      
      // 确保 serviceId 是数字类型
      const numId = Number(serviceId)
      console.log('请求服务详情，serviceId:', numId, '类型:', typeof numId)
      
      const result = await api.serviceDetail({ serviceId: numId })
      console.log('服务详情接口返回:', result)
      
      if (result.code === 0 && result.data) {
        const service = result.data
        
        // 兼容历史 imageCosId，统一转换为静态资源地址
        const processedService = await getSingleServiceImage(service)
        const localizedService = await localizeServiceMedia(processedService)
        
        // 解析表单配置
        let formFields = []
        try {
          // 检查formConfig是否存在且不为空
          if (!localizedService.formConfig || localizedService.formConfig === '' || localizedService.formConfig === 'null') {
            console.log('服务没有表单配置，使用默认配置')
            formFields = []
          } else {
            const formConfig = JSON.parse(localizedService.formConfig)
            const allFields = normalizeFormFields(formConfig.fields || [])
          
          // 对于预约咨询服务，显示所有字段
          if (localizedService.category === '预约咨询') {
            formFields = allFields
            console.log('预约咨询服务，显示所有字段:', formFields)
          } else {
            // 其他服务类型，过滤掉身份证号、出生日期和营养目标相关的字段
            formFields = allFields.filter(field => {
              const fieldName = field.name || ''
              const fieldLabel = field.label || ''
              
              // 过滤掉包含身份证号、出生日期、营养目标等关键词的字段
              const excludeKeywords = ['身份证', 'idCard', 'birthday', '出生', '生日', '营养目标', 'nutritionGoal']
              const shouldExclude = excludeKeywords.some(keyword => 
                fieldName.toLowerCase().includes(keyword.toLowerCase()) ||
                fieldLabel.toLowerCase().includes(keyword.toLowerCase())
              )
              
              return !shouldExclude
            })
          
          // 重新排序字段：将疾病史字段移到体重字段后面
          formFields = formFields.sort((a, b) => {
            const aName = a.name || ''
            const bName = b.name || ''
            
            // 定义字段顺序
            const fieldOrder = {
              'patientName': 1,
              'patientPhone': 2,
              'patientHeight': 3,
              'patientWeight': 4,
              'medicalHistory': 5, // 疾病史移到体重后面
              'dietaryRestrictions': 6
            }
            
            const aOrder = fieldOrder[aName] || 999
            const bOrder = fieldOrder[bName] || 999
            
            return aOrder - bOrder
          })
          
            console.log('过滤并重新排序后的表单字段:', formFields)
          }
        }
        } catch (error) {
          console.error('解析表单配置失败:', error)
        }
        
        // 初始化表单数据
        const formData = {
          diseaseInfo: '', // 既往病史
          needToiletAssist: ''
        }
        formFields.forEach(field => {
          formData[field.name] = ''
        })
        
        console.log('初始化后的 formData:', formData)
        console.log('formFields 中的字段名:', formFields.map(f => f.name))
        
        this.setData({
          service: localizedService,
          formFields,
          formData,
          loading: false
        })
        
        console.log('设置后的 formData:', this.data.formData)
      } else {
        throw new Error(result.message || '加载失败')
      }
    } catch (error) {
      console.error('加载服务详情失败:', error)
      this.setData({ loading: false })
      
      // 处理特定错误
      let errorMessage = error.message || '加载失败'
      
      // 如果是服务不存在错误，提供更友好的提示
      if (errorMessage.includes('不存在')) {
        errorMessage = '该服务不存在或已下架，请选择其他服务'
        
        // 显示错误提示并跳转到服务列表
        wx.showModal({
          title: '服务不存在',
          content: errorMessage,
          showCancel: true,
          cancelText: '返回',
          confirmText: '查看其他服务',
          success: (res) => {
            if (res.confirm) {
              wx.navigateTo({
                url: '/pages/service/list'
              })
            } else {
              wx.navigateBack()
            }
          }
        })
        return
      }
      
      wx.showToast({
        title: errorMessage,
        icon: 'none'
      })
    }
  },

  // 表单输入处理
  onFormInput(e) {
    const { field } = e.currentTarget.dataset
    const { value } = e.detail
    
    console.log('表单输入:', { field, value })
    
    this.setData({
      [`formData.${field}`]: value
    })
    
    // 检查设置后的数据
    setTimeout(() => {
      console.log('设置后的 formData:', this.data.formData)
    }, 100)
  },

  // 日期选择
  onDateChange(e) {
    this.setData({
      selectedDate: e.detail.value
    })
  },

  // 单选按钮选择
  onRadioChange(e) {
    const { field } = e.currentTarget.dataset
    const { value } = e.detail
    
    this.setData({
      [`formData.${field}`]: value
    })
  },

  // 下拉选择
  onSelectChange(e) {
    const { field } = e.currentTarget.dataset
    const { options = [] } = e.currentTarget.dataset
    const selectedIndex = Number(e.detail.value)
    const selectedOption = options[selectedIndex]
    const selectedValue = selectedOption && selectedOption.value != null
      ? String(selectedOption.value)
      : ''
    
    this.setData({
      [`formData.${field}`]: selectedValue
    })
  },

  // 获取选中选项的标签
  getSelectedOptionLabel(fieldName) {
    const field = this.data.formFields.find(f => f.name === fieldName)
    if (!field || !field.options) return ''
    
    const selectedValue = this.data.formData[fieldName]
    const option = field.options.find(opt => opt.value === selectedValue)
    return option ? option.label : ''
  },

  // 计算年龄
  getAge(birthday) {
    if (!birthday) return ''
    const today = new Date()
    const birthDate = new Date(birthday)
    let age = today.getFullYear() - birthDate.getFullYear()
    const monthDiff = today.getMonth() - birthDate.getMonth()
    
    if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birthDate.getDate())) {
      age--
    }
    
    return age + '岁'
  },

  // 显示时间选择器
  onShowTimePicker() {
    this.setData({ showTimePicker: true })
  },

  // 时间选择器确认
  onTimePickerConfirm(e) {
    console.log('时间选择器确认事件:', e)
    console.log('选择的时间:', e.detail)
    
    const { dateTime, date, time } = e.detail
    console.log('解析的数据:', { dateTime, date, time })
    
    if (dateTime) {
      this.setData({
        selectedDateTime: dateTime,
        selectedDate: date,
        selectedTime: time,
        showTimePicker: false
      })
      console.log('日期时间设置成功:', dateTime)
      
      // 检查设置后的数据
      setTimeout(() => {
        console.log('设置后的数据状态:', {
          selectedDateTime: this.data.selectedDateTime,
          selectedDate: this.data.selectedDate,
          selectedTime: this.data.selectedTime
        })
      }, 100)
    } else {
      console.warn('日期时间值为空，不进行设置')
      this.setData({ showTimePicker: false })
    }
  },

  // 时间选择器关闭
  onTimePickerClose() {
    this.setData({ showTimePicker: false })
  },

  // 显示咨询时间选择器
  onShowConsultTimePicker() {
    this.setData({ showConsultTimePicker: true })
  },

  // 咨询时间选择器确认
  onConsultTimePickerConfirm(e) {
    console.log('咨询时间选择器确认事件:', e)
    console.log('选择的数据:', e.detail)
    
    const { dateTime, date, time } = e.detail
    console.log('解析的数据:', { dateTime, date, time })
    
    if (dateTime) {
      this.setData({
        consultDateTime: dateTime,
        consultDate: date,
        consultTime: time,
        showConsultTimePicker: false
      })
      console.log('咨询时间设置成功:', dateTime)
      
      // 检查设置后的数据
      setTimeout(() => {
        console.log('设置后的咨询时间数据状态:', {
          consultDateTime: this.data.consultDateTime,
          consultDate: this.data.consultDate,
          consultTime: this.data.consultTime
        })
      }, 100)
    } else {
      console.warn('咨询时间值为空，不进行设置')
      this.setData({ showConsultTimePicker: false })
    }
  },

  // 咨询时间选择器关闭
  onConsultTimePickerClose() {
    this.setData({ showConsultTimePicker: false })
  },

  // 备注输入
  onRemarkInput(e) {
    this.setData({ remark: e.detail.value })
  },

  // 选择就诊人
  onPatientSelect(e) {
    console.log('就诊人选择事件:', e)
    console.log('patientList 数据:', this.data.patientList)
    console.log('patientList 长度:', this.data.patientList.length)
    console.log('currentTarget dataset:', e.currentTarget.dataset)
    console.log('target dataset:', e.target.dataset)
    
    // 直接获取索引，并转换为数字
    let index = e.currentTarget.dataset.index
    console.log('原始索引值:', index, '类型:', typeof index)
    
    // 如果索引是字符串，转换为数字
    if (typeof index === 'string') {
      index = parseInt(index, 10)
    }
    
    console.log('转换后的索引:', index, '类型:', typeof index)
    
    if (index === undefined || index === null || isNaN(index)) {
      console.warn('无法获取有效的就诊人索引')
      return
    }
    
    const patient = this.data.patientList[index]
    console.log('就诊人列表:', this.data.patientList)
    console.log('索引对应的就诊人:', patient)
    
    // 确保 patient 存在且不为 undefined
    if (patient) {
      this.setData({ selectedPatient: patient })
      console.log('选择就诊人成功:', patient)
    } else {
      console.warn('就诊人数据不存在，index:', index, '列表长度:', this.data.patientList.length)
    }
  },

  // 选择地址
  onAddressSelect(e) {
    console.log('地址选择事件:', e)
    
    // 直接获取索引，并转换为数字
    let index = e.currentTarget.dataset.index
    console.log('原始索引值:', index, '类型:', typeof index)
    
    // 如果索引是字符串，转换为数字
    if (typeof index === 'string') {
      index = parseInt(index, 10)
    }
    
    console.log('转换后的索引:', index, '类型:', typeof index)
    
    if (index === undefined || index === null || isNaN(index)) {
      console.warn('无法获取有效的地址索引')
      return
    }
    
    const address = this.data.addressList[index]
    console.log('地址列表:', this.data.addressList)
    console.log('索引对应的地址:', address)
    
    // 确保 address 存在且不为 undefined
    if (address) {
      this.setData({ selectedAddress: address })
      console.log('选择地址成功:', address)
    } else {
      console.warn('地址数据不存在，index:', index, '列表长度:', this.data.addressList.length)
    }
  },

  // 阻止事件冒泡
  stopPropagation(e) {
    // 空方法，用于阻止事件冒泡
  },

  // 显示就诊人选择弹窗
  onShowPatientModal() {
    this.setData({ showPatientModal: true })
  },

  // 显示地址选择弹窗
  onShowAddressModal() {
    this.setData({ showAddressModal: true })
  },

  // 添加就诊人
  onAddPatient() {
    wx.navigateTo({
      url: '/pages/user/patient/add'
    })
  },

  // 添加地址
  onAddAddress() {
    wx.navigateTo({
      url: '/pages/user/address/add'
    })
  },

  // 关闭弹窗
  onCloseModal() {
    this.setData({
      showPatientModal: false,
      showAddressModal: false
    })
  },

  // 验证表单
  validateForm() {
    const { formData, selectedPatient, selectedAddress, selectedDateTime, consultDateTime, service } = this.data
    
    console.log('表单验证开始，当前数据:', {
      formData,
      selectedPatient,
      selectedAddress,
      selectedDateTime,
      consultDateTime,
      serviceCategory: processedService.category
    })
    
    // 判断是否为智慧养老设备
    const isSmartElderly = processedService.category === '智慧养老'
    
    if (isSmartElderly) {
      // 智慧养老设备只需要验证地址
      if (!selectedAddress) {
        console.log('未选择收货地址')
        wx.showToast({
          title: '请选择收货地址',
          icon: 'none'
        })
        return false
      }
      
      console.log('智慧养老设备表单验证通过')
      return true
    } else {
      // 其他服务需要验证患者信息和预约时间
      // 检查必填字段
      const requiredFields = ['patientName', 'patientPhone']
      for (const field of requiredFields) {
        console.log(`检查字段 ${field}:`, formData[field])
        if (!formData[field]) {
          console.log('必填字段缺失:', field)
          
          // 如果选择了就诊人，尝试从就诊人中获取信息
          if (selectedPatient && field === 'patientName' && selectedPatient.name) {
            console.log('从就诊人中获取患者姓名:', selectedPatient.name)
            this.setData({
              [`formData.${field}`]: selectedPatient.name
            })
            continue
          }
          
          if (selectedPatient && field === 'patientPhone' && selectedPatient.phone) {
            console.log('从就诊人中获取患者电话:', selectedPatient.phone)
            this.setData({
              [`formData.${field}`]: selectedPatient.phone
            })
            continue
          }
          
          wx.showToast({
            title: '请填写完整信息',
            icon: 'none'
          })
          return false
        }
      }
      
      // 检查是否选择了就诊人
      if (!selectedPatient) {
        console.log('未选择就诊人')
        wx.showToast({
          title: '请选择就诊人',
          icon: 'none'
        })
        return false
      }
      
      // 检查是否选择了地址
      if (!selectedAddress) {
        console.log('未选择服务地址')
        wx.showToast({
          title: '请选择服务地址',
          icon: 'none'
        })
        return false
      }
      
      // 检查是否选择了预约时间
      if (!selectedDateTime) {
        console.log('未选择预约时间')
        wx.showToast({
          title: '请选择预约时间',
          icon: 'none'
        })
        return false
      }
      
      console.log('普通服务表单验证通过')
      return true
    }
  },

  // 提交订单
  async onSubmitOrder() {
    console.log('开始提交订单')
    if (!this.validateForm()) {
      console.log('表单验证失败')
      return
    }
    try {
      wx.showLoading({ title: '正在提交...' })
      const { service, formData, selectedPatient, selectedAddress, selectedDateTime, remark } = this.data
      
      // 获取用户ID
      const userId = wx.getStorageSync('userId')
      console.log('获取到的userId:', userId, '类型:', typeof userId)
      
      if (!userId) {
        throw new Error('用户未登录')
      }
      
      // 确保userId是有效的字符串格式
      const strUserId = String(userId)
      if (!strUserId || strUserId === 'undefined' || strUserId === 'null') {
        console.error('userId转换失败:', userId)
        throw new Error('用户ID无效，请重新登录')
      }
      
      // 判断是否为智慧养老设备
      const isSmartElderly = processedService.category === '智慧养老'
      
      let orderData = {
        userId: strUserId,
        serviceId: processedService.id,
        addressId: selectedAddress.id,
        quantity: 1, // 默认数量为1
        remark: remark || '',
        formData: formData // 保持原有的formData结构
      }
      
      if (isSmartElderly) {
        // 智慧养老设备不需要患者信息和预约时间
        orderData.patientId = null
        orderData.appointmentDate = null
        orderData.appointmentTime = null
        orderData.diseaseInfo = ''
        orderData.needToiletAssist = ''
        
        console.log('智慧养老设备订单数据:', orderData)
      } else {
        // 普通服务需要患者信息和预约时间
        // 解析日期时间
        const [appointmentDate, appointmentTime] = selectedDateTime.split(' ')
        
        orderData.patientId = selectedPatient.id
        orderData.appointmentDate = appointmentDate
        orderData.appointmentTime = appointmentTime
        orderData.diseaseInfo = formData.diseaseInfo || ''
        orderData.needToiletAssist = formData.needToiletAssist || ''
        
        console.log('普通服务订单数据:', orderData)
      }
      
      console.log('发送到后端的订单数据:', orderData)
      
      // 根据服务类型选择不同的API端点
      let result
      if (isSmartElderly) {
        result = await api.smartElderlyOrderSubmit(orderData)
      } else {
        result = await api.orderSubmit(orderData)
      }
      if (result.code === 0) {
        wx.hideLoading()
        wx.showToast({
          title: '订单提交成功',
          icon: 'success'
        })
        // 跳转到订单详情页
        setTimeout(() => {
          wx.redirectTo({
            url: `/pages/order/detail?orderNo=${result.data.orderNo}`
          })
        }, 1500)
      } else {
        throw new Error(result.message || '提交失败')
      }
    } catch (error) {
      console.error('订单提交失败:', error)
      wx.hideLoading()
      wx.showToast({
        title: error.message || '提交失败，请重试',
        icon: 'none'
      })
    }
  },

  // 联系客服
  onContactService() {
    wx.navigateTo({
      url: '/pages/kefu/chat'
    })
  },

  // 分享
  onShareAppMessage() {
    const { service } = this.data
    return {
      title: service.title || '专业护理服务',
      path: `/pages/service/detail?id=${service.id}`,
      imageUrl: service.imageUrl
    }
  },

  // 下拉刷新
  async onPullDownRefresh() {
    const { serviceId } = this.data
    if (serviceId) {
      await this.loadServiceDetail(serviceId)
    }
    await this.loadUserData()
    wx.stopPullDownRefresh()
  }
})