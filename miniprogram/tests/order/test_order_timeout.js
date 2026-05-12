/**
 * 订单超时功能测试脚本
 * 用于验证30分钟超时自动取消功能
 */

const app = getApp()

class OrderTimeoutTester {
  constructor() {
    this.testResults = []
  }

  /**
   * 开始测试
   */
  async startTest() {
    console.log('开始订单超时功能测试...')
    
    // 测试创建订单
    await this.testCreateOrder()
    
    // 测试检查超时订单
    await this.testCheckExpiredOrders()
    
    // 测试获取超时订单数量
    await this.testGetExpiredCount()
    
    // 输出测试结果
    this.printTestResults()
  }

  /**
   * 测试创建订单
   */
  async testCreateOrder() {
    console.log('测试创建订单...')
    
    try {
      const orderData = {
        userId: 1,
        serviceId: 1,
        patientId: 1,
        addressId: 1,
        appointmentDate: '2024-01-15',
        appointmentTime: '09:00',
        quantity: 1,
        formData: {
          consultTime: '30分钟',
          specialRequirements: '无特殊要求'
        },
        diseaseInfo: '无',
        needToiletAssist: '0',
        remark: '测试订单'
      }

      const result = await app.call({
        path: '/api/order/submit',
        method: 'POST',
        data: orderData
      })

      if (result.code === 0) {
        this.addTestResult('创建订单测试', {
          success: true,
          orderId: result.data.orderId,
          orderNo: result.data.orderNo,
          message: '订单创建成功'
        })
        
        // 保存订单信息用于后续测试
        this.testOrderId = result.data.orderId
        this.testOrderNo = result.data.orderNo
        
        console.log('订单创建成功:', result.data)
      } else {
        this.addTestResult('创建订单测试', {
          success: false,
          error: result.message || '创建订单失败'
        })
      }
      
    } catch (error) {
      console.error('创建订单测试失败:', error)
      this.addTestResult('创建订单测试', {
        success: false,
        error: error.message
      })
    }
  }

  /**
   * 测试检查超时订单
   */
  async testCheckExpiredOrders() {
    console.log('测试检查超时订单...')
    
    try {
      const result = await app.call({
        path: '/api/order/check_expired',
        method: 'POST',
        data: {}
      })

      if (result.code === 0) {
        this.addTestResult('检查超时订单测试', {
          success: true,
          expiredCount: result.data.expiredCount,
          message: result.data.message
        })
        
        console.log('超时订单检查完成:', result.data)
      } else {
        this.addTestResult('检查超时订单测试', {
          success: false,
          error: result.message || '检查超时订单失败'
        })
      }
      
    } catch (error) {
      console.error('检查超时订单测试失败:', error)
      this.addTestResult('检查超时订单测试', {
        success: false,
        error: error.message
      })
    }
  }

  /**
   * 测试获取超时订单数量
   */
  async testGetExpiredCount() {
    console.log('测试获取超时订单数量...')
    
    try {
      const result = await app.call({
        path: '/api/order/expired_count',
        method: 'GET',
        data: {}
      })

      if (result.code === 0) {
        this.addTestResult('获取超时订单数量测试', {
          success: true,
          expiredCount: result.data.expiredCount,
          message: '获取超时订单数量成功'
        })
        
        console.log('超时订单数量:', result.data.expiredCount)
      } else {
        this.addTestResult('获取超时订单数量测试', {
          success: false,
          error: result.message || '获取超时订单数量失败'
        })
      }
      
    } catch (error) {
      console.error('获取超时订单数量测试失败:', error)
      this.addTestResult('获取超时订单数量测试', {
        success: false,
        error: error.message
      })
    }
  }

  /**
   * 测试订单详情（验证支付截止时间）
   */
  async testOrderDetail() {
    if (!this.testOrderNo) {
      console.log('跳过订单详情测试，没有测试订单号')
      return
    }
    
    console.log('测试订单详情...')
    
    try {
      const result = await app.call({
        path: '/api/order/detail',
        method: 'POST',
        data: {
          orderNo: this.testOrderNo
        }
      })

      if (result.code === 0 && result.data) {
        const order = result.data
        
        // 检查是否有支付截止时间
        if (order.payDeadline) {
          this.addTestResult('订单详情测试', {
            success: true,
            payDeadline: order.payDeadline,
            message: '订单包含支付截止时间'
          })
          
          console.log('订单支付截止时间:', order.payDeadline)
        } else {
          this.addTestResult('订单详情测试', {
            success: false,
            error: '订单缺少支付截止时间'
          })
        }
      } else {
        this.addTestResult('订单详情测试', {
          success: false,
          error: result.message || '获取订单详情失败'
        })
      }
      
    } catch (error) {
      console.error('订单详情测试失败:', error)
      this.addTestResult('订单详情测试', {
        success: false,
        error: error.message
      })
    }
  }

  /**
   * 添加测试结果
   */
  addTestResult(testName, result) {
    this.testResults.push({
      name: testName,
      result: result,
      timestamp: Date.now()
    })
  }

  /**
   * 输出测试结果
   */
  printTestResults() {
    console.log('\n=== 订单超时功能测试结果 ===')
    
    let successCount = 0
    let totalCount = this.testResults.length
    
    this.testResults.forEach((test, index) => {
      const status = test.result.success ? '✅' : '❌'
      console.log(`${index + 1}. ${status} ${test.name}`)
      
      if (test.result.success) {
        successCount++
        if (test.result.message) {
          console.log(`   信息: ${test.result.message}`)
        }
        if (test.result.expiredCount !== undefined) {
          console.log(`   超时订单数量: ${test.result.expiredCount}`)
        }
        if (test.result.payDeadline) {
          console.log(`   支付截止时间: ${test.result.payDeadline}`)
        }
      } else {
        console.log(`   错误: ${test.result.error}`)
      }
    })
    
    console.log(`\n测试总结: ${successCount}/${totalCount} 通过`)
    console.log(`成功率: ${((successCount / totalCount) * 100).toFixed(1)}%`)
    
    if (successCount === totalCount) {
      console.log('🎉 所有测试通过！订单超时功能正常工作')
    } else {
      console.log('⚠️  部分测试失败，请检查配置和网络连接')
    }
  }

  /**
   * 模拟超时订单（用于测试）
   */
  async simulateExpiredOrder() {
    console.log('模拟超时订单...')
    
    // 这里可以创建一个特殊的测试订单，设置较早的支付截止时间
    // 实际测试中，可以等待30分钟或修改数据库中的支付截止时间
    console.log('注意：实际测试需要等待30分钟或手动修改数据库')
  }
}

// 导出测试器
module.exports = OrderTimeoutTester

// 如果直接运行此文件，执行测试
if (typeof module !== 'undefined' && module.exports) {
  const tester = new OrderTimeoutTester()
  tester.startTest().catch(console.error)
} 