// components/calendar-picker/calendar-picker.js
const { api } = require('../../utils/cloud-container-standard')

Component({
  properties: {
    // 是否显示
    show: {
      type: Boolean,
      value: false,
      observer: function(newVal) {
        console.log('=== 日历组件 show 属性变化 ===')
        console.log('show 值:', newVal)
        if (newVal) {
          console.log('显示日历组件')
          this.initCalendar()
        } else {
          console.log('隐藏日历组件')
        }
      }
    },
    // 最小日期（明天）
    minDate: {
      type: String,
      value: ''
    },
    // 最大日期（7天后）
    maxDate: {
      type: String,
      value: ''
    },
    // 已选择的日期时间
    selectedDateTime: {
      type: String,
      value: ''
    }
  },

  data: {
    // 当前显示的月份
    currentMonth: '',
    // 日历数据
    calendarDays: [],
    // 可选时间段（动态获取）
    timeSlots: [],
    // 选中的日期
    selectedDate: '',
    // 选中的时间
    selectedTime: '',
    // 当前日期（用于禁用今天）
    today: '',
    // 月份导航状态
    canPrevMonth: true,
    canNextMonth: true,
    // 基础URL
    baseUrl: ''
  },

  lifetimes: {
    attached() {
      console.log('日历组件已附加')
      this.initCalendar()
    }
  },

  observers: {
    'show': function(show) {
      console.log('=== 日历组件 show 属性变化 ===')
      console.log('show 值:', show)
      if (show) {
        console.log('显示日历组件')
        this.initCalendar()
      } else {
        console.log('隐藏日历组件')
      }
    },
    'selectedDateTime': function(selectedDateTime) {
      console.log('=== 日历组件 selectedDateTime 属性变化 ===')
      console.log('selectedDateTime 值:', selectedDateTime)
      if (selectedDateTime) {
        // 解析已选择的日期时间
        const [date, time] = selectedDateTime.split(' ')
        this.setData({
          selectedDate: date,
          selectedTime: time || ''
        })
        // 重新生成日历数据以显示选中状态
        this.generateCalendarDays()
      }
    }
  },

  methods: {
    // 初始化日历
    initCalendar() {
      const today = new Date()
      const currentMonth = this.formatDate(today, 'YYYY-MM')
      
      this.setData({ 
        currentMonth,
        selectedDate: this.data.selectedDateTime ? this.data.selectedDateTime.split(' ')[0] : '',
        selectedTime: this.data.selectedDateTime ? this.data.selectedDateTime.split(' ')[1] : ''
      })
      
      this.generateCalendarDays()
      this.updateMonthNavStatus()
    },

    // 生成日历天数
    generateCalendarDays() {
      const [year, month] = this.data.currentMonth.split('-').map(Number)
      const firstDay = new Date(year, month - 1, 1)
      const lastDay = new Date(year, month, 0)
      const startDate = new Date(firstDay)
      startDate.setDate(startDate.getDate() - firstDay.getDay())
      
      const weeks = []
      const today = new Date()
      today.setHours(0, 0, 0, 0) // 设置为00:00:00，避免时区问题
      
      const minDate = this.data.minDate ? new Date(this.data.minDate) : null
      const maxDate = this.data.maxDate ? new Date(this.data.maxDate) : null
      
      if (minDate) {
        minDate.setHours(0, 0, 0, 0)
      }
      if (maxDate) {
        maxDate.setHours(0, 0, 0, 0)
      }
      
      console.log('日历组件日期范围:', {
        minDate: minDate ? minDate.toISOString() : 'null',
        maxDate: maxDate ? maxDate.toISOString() : 'null',
        today: today.toISOString()
      })
      
      // 生成6周的日历数据
      for (let week = 0; week < 6; week++) {
        const weekDays = []
        for (let day = 0; day < 7; day++) {
          const date = new Date(startDate)
          date.setDate(startDate.getDate() + week * 7 + day)
          date.setHours(0, 0, 0, 0) // 设置为00:00:00，避免时区问题
          
          const dateStr = this.formatDate(date)
          const isCurrentMonth = date.getMonth() === month - 1
          const isToday = this.formatDate(date) === this.formatDate(today)
          const isSelected = dateStr === this.data.selectedDate
          const isDisabled = (minDate && date < minDate) || (maxDate && date > maxDate)
          
          weekDays.push({
            date: dateStr,
            day: date.getDate(),
            isCurrentMonth,
            isToday,
            isSelected,
            isDisabled
          })
        }
        weeks.push(weekDays)
      }
      
      this.setData({ calendarDays: weeks })
    },

    // 更新月份导航状态
    updateMonthNavStatus() {
      const [year, month] = this.data.currentMonth.split('-').map(Number)
      const minDate = this.data.minDate ? new Date(this.data.minDate) : null
      const maxDate = this.data.maxDate ? new Date(this.data.maxDate) : null
      
      const prevMonth = new Date(year, month - 2, 1)
      const nextMonth = new Date(year, month, 1)
      
      this.setData({
        canPrevMonth: !minDate || prevMonth >= minDate,
        canNextMonth: !maxDate || nextMonth <= maxDate
      })
    },

    // 格式化日期
    formatDate(date, format = 'YYYY-MM-DD') {
      const year = date.getFullYear()
      const month = String(date.getMonth() + 1).padStart(2, '0')
      const day = String(date.getDate()).padStart(2, '0')
      
      if (format === 'YYYY-MM') {
        return `${year}-${month}`
      }
      return `${year}-${month}-${day}`
    },

    // 选择日期
    selectDate(e) {
      const { date, isDisabled } = e.currentTarget.dataset
      if (!date || isDisabled) return
      
      console.log('选择日期:', date)
      
      this.setData({ 
        selectedDate: date,
        selectedTime: '' // 清空时间选择
      })
      
      // 重新生成日历以更新选中状态
      this.generateCalendarDays()
      
      // 获取该日期的可用时间槽
      this.getAvailableTimeSlots(date)
    },

    // 获取可用时间槽
    async getAvailableTimeSlots(date) {
      try {
        console.log('获取可用时间槽:', date)
        
        const result = await api.orderTimeSlots({ date })
        console.log('时间槽接口返回:', result)
        
        if (result && result.code === 0 && result.data) {
          const timeSlots = result.data.timeSlots.map(time => ({
            time,
            label: time
          }))
          
          console.log('可用时间槽:', timeSlots)
          this.setData({ timeSlots })
        } else {
          console.error('获取时间槽失败:', result)
          // 如果获取失败，使用默认时间槽
          this.setDefaultTimeSlots()
        }
      } catch (error) {
        console.error('获取时间槽请求失败:', error)
        // 如果请求失败，使用默认时间槽
        this.setDefaultTimeSlots()
      }
    },

    // 设置默认时间槽
    setDefaultTimeSlots() {
      const defaultTimeSlots = [
        { time: '08:00', label: '08:00' },
        { time: '09:00', label: '09:00' },
        { time: '10:00', label: '10:00' },
        { time: '11:00', label: '11:00' },
        { time: '14:00', label: '14:00' },
        { time: '15:00', label: '15:00' },
        { time: '16:00', label: '16:00' },
        { time: '17:00', label: '17:00' },
        { time: '18:00', label: '18:00' },
        { time: '19:00', label: '19:00' }
      ]
      this.setData({ timeSlots: defaultTimeSlots })
    },

    // 选择时间
    selectTime(e) {
      const { time } = e.currentTarget.dataset
      console.log('选择时间:', time)
      
      if (time) {
        this.setData({ selectedTime: time })
        console.log('时间设置成功:', time)
        console.log('当前选择状态:', {
          selectedDate: this.data.selectedDate,
          selectedTime: this.data.selectedTime
        })
      }
    },

    // 上个月
    prevMonth() {
      if (!this.data.canPrevMonth) return
      
      const [year, month] = this.data.currentMonth.split('-').map(Number)
      const prevDate = new Date(year, month - 2, 1)
      const prevMonthStr = this.formatDate(prevDate, 'YYYY-MM')
      
      this.setData({ currentMonth: prevMonthStr })
      this.generateCalendarDays()
      this.updateMonthNavStatus()
    },

    // 下个月
    nextMonth() {
      if (!this.data.canNextMonth) return
      
      const [year, month] = this.data.currentMonth.split('-').map(Number)
      const nextDate = new Date(year, month, 1)
      const nextMonthStr = this.formatDate(nextDate, 'YYYY-MM')
      
      this.setData({ currentMonth: nextMonthStr })
      this.generateCalendarDays()
      this.updateMonthNavStatus()
    },

    // 确认选择
    confirmSelection() {
      const { selectedDate, selectedTime } = this.data
      console.log('日历组件确认选择:', { selectedDate, selectedTime })
      
      if (!selectedDate || !selectedTime) {
        wx.showToast({
          title: '请选择日期和时间',
          icon: 'none'
        })
        return
      }
      
      const dateTime = `${selectedDate} ${selectedTime}`
      console.log('生成的日期时间字符串:', dateTime)
      
      const eventData = { dateTime, date: selectedDate, time: selectedTime }
      console.log('触发确认事件，数据:', eventData)
      
      this.triggerEvent('confirm', eventData)
      this.triggerEvent('close')
    },

    // 取消选择
    cancel() {
      console.log('用户点击了取消按钮')
      this.triggerEvent('close')
    },

    // 隐藏选择器
    hide() {
      console.log('用户点击了隐藏区域')
      this.triggerEvent('close')
    },

    // 显示选择器
    show() {
      this.initCalendar()
    },

    // 阻止冒泡
    preventBubble() {
      // 阻止事件冒泡
    }
  }
}) 