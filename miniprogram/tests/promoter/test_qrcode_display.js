// 测试二维码显示功能

function testQRCodeDisplay() {
  console.log('=== 测试二维码显示功能 ===')
  
  // 模拟推广员信息
  const mockPromoterInfo = {
    userId: '507f1f77bcf86cd799439011',
    promoterCode: 'ABC123',
    nickName: '测试用户',
    avatarUrl: 'https://example.com/avatar.jpg',
    qrCodeUrl: 'https://via.placeholder.com/256x256/CCCCCC/666666?text=QR+Code+ABC123',
    totalIncome: 100.50,
    todayIncome: 10.00,
    monthIncome: 50.00,
    totalOrders: 5,
    todayOrders: 1,
    monthOrders: 3
  }
  
  // 测试二维码URL处理
  function testQRCodeUrl(qrCodeUrl) {
    console.log('测试二维码URL:', qrCodeUrl)
    
    if (!qrCodeUrl) {
      console.log('❌ 二维码URL为空')
      return false
    }
    
    if (qrCodeUrl.includes('example.com')) {
      console.log('⚠️ 使用占位符图片（可能是404错误）')
      return false
    }
    
    if (qrCodeUrl.startsWith('data:image/')) {
      console.log('✅ 使用Base64编码的二维码')
      return true
    }
    
    if (qrCodeUrl.startsWith('http')) {
      console.log('✅ 使用外部图片URL')
      return true
    }
    
    console.log('❌ 无效的二维码URL格式')
    return false
  }
  
  // 测试图片加载
  function testImageLoad(url) {
    return new Promise((resolve) => {
      const img = new Image()
      img.onload = () => {
        console.log('✅ 图片加载成功:', url)
        resolve(true)
      }
      img.onerror = () => {
        console.log('❌ 图片加载失败:', url)
        resolve(false)
      }
      img.src = url
    })
  }
  
  // 测试二维码预览功能
  function testQRCodePreview(qrCodeUrl) {
    console.log('测试二维码预览功能')
    
    if (!qrCodeUrl) {
      console.log('❌ 无法预览：二维码URL为空')
      return false
    }
    
    // 模拟wx.previewImage调用
    console.log('📱 调用wx.previewImage:', {
      urls: [qrCodeUrl],
      current: qrCodeUrl
    })
    
    return true
  }
  
  // 执行测试
  console.log('1. 测试二维码URL格式:')
  const urlValid = testQRCodeUrl(mockPromoterInfo.qrCodeUrl)
  
  console.log('\n2. 测试图片加载:')
  testImageLoad(mockPromoterInfo.qrCodeUrl).then(success => {
    console.log('图片加载测试结果:', success ? '成功' : '失败')
  })
  
  console.log('\n3. 测试二维码预览:')
  const previewValid = testQRCodePreview(mockPromoterInfo.qrCodeUrl)
  
  // 测试不同的二维码URL格式
  console.log('\n4. 测试不同格式的二维码URL:')
  const testUrls = [
    'https://via.placeholder.com/256x256/CCCCCC/666666?text=QR+Code+ABC123',
    'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==',
    'https://your-domain.com/static/qrcode/promoter_ABC123.png',
    '',
    null
  ]
  
  testUrls.forEach((url, index) => {
    console.log(`\n测试URL ${index + 1}:`)
    testQRCodeUrl(url)
  })
  
  // 测试错误处理
  console.log('\n5. 测试错误处理:')
  
  function testErrorHandling() {
    const promoterInfo = mockPromoterInfo
    
    // 模拟二维码URL为空的情况
    if (!promoterInfo.qrCodeUrl) {
      console.log('⚠️ 二维码URL为空，显示占位符')
      return 'https://via.placeholder.com/256x256/CCCCCC/666666?text=QR+Code+Error'
    }
    
    // 模拟二维码URL无效的情况
    if (promoterInfo.qrCodeUrl.includes('example.com')) {
      console.log('⚠️ 检测到占位符URL，可能需要重新生成二维码')
      return 'https://via.placeholder.com/256x256/CCCCCC/666666?text=QR+Code+Error'
    }
    
    return promoterInfo.qrCodeUrl
  }
  
  const fallbackUrl = testErrorHandling()
  console.log('错误处理后的URL:', fallbackUrl)
  
  console.log('\n=== 测试完成 ===')
}

// 运行测试
testQRCodeDisplay() 