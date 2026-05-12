// 测试API响应格式，确保list字段不为null

function testApiResponse() {
  console.log('=== 测试API响应格式 ===')
  
  // 模拟API响应数据
  const mockResponses = {
    // 正常情况：有数据
    withData: {
      code: 0,
      data: {
        list: [
          { id: 1, amount: 100, status: 1 },
          { id: 2, amount: 200, status: 0 }
        ],
        total: 2,
        page: 1,
        pageSize: 20,
        hasMore: false
      }
    },
    
    // 空数据情况：list应该为空数组而不是null
    emptyData: {
      code: 0,
      data: {
        list: [], // 应该是空数组，不是null
        total: 0,
        page: 1,
        pageSize: 20,
        hasMore: false
      }
    },
    
    // 错误情况：list为null（这是我们要修复的问题）
    nullData: {
      code: 0,
      data: {
        list: null, // 这是问题所在
        total: 0,
        page: 1,
        pageSize: 20,
        hasMore: false
      }
    }
  }
  
  // 测试处理函数
  function processApiResponse(response) {
    console.log('处理API响应:', response)
    
    try {
      if (response.code === 0) {
        const { list, hasMore } = response.data
        
        // 安全检查：确保list不为null
        const safeList = list || []
        
        console.log(`✓ 安全处理后的list长度: ${safeList.length}`)
        console.log(`✓ hasMore: ${hasMore}`)
        
        return {
          success: true,
          listLength: safeList.length,
          hasMore: hasMore
        }
      } else {
        console.log('✗ API返回错误码:', response.code)
        return { success: false, error: 'API错误' }
      }
    } catch (error) {
      console.error('✗ 处理响应时出错:', error)
      return { success: false, error: error.message }
    }
  }
  
  // 执行测试
  const testCases = [
    { name: '有数据的情况', response: mockResponses.withData },
    { name: '空数据的情况', response: mockResponses.emptyData },
    { name: 'list为null的情况', response: mockResponses.nullData }
  ]
  
  testCases.forEach((testCase, index) => {
    console.log(`\n测试用例 ${index + 1}: ${testCase.name}`)
    const result = processApiResponse(testCase.response)
    
    if (result.success) {
      console.log(`✅ 测试通过: list长度=${result.listLength}, hasMore=${result.hasMore}`)
    } else {
      console.log(`❌ 测试失败: ${result.error}`)
    }
  })
  
  // 测试前端处理逻辑
  console.log('\n=== 测试前端处理逻辑 ===')
  
  function testFrontendProcessing(response) {
    const { list, hasMore } = response.data
    const safeList = list || []
    const currentList = [] // 模拟当前页面的数据
    const mergedList = [...currentList, ...safeList]
    
    console.log(`原始list: ${JSON.stringify(list)}`)
    console.log(`安全list: ${JSON.stringify(safeList)}`)
    console.log(`合并后list: ${JSON.stringify(mergedList)}`)
    
    return mergedList
  }
  
  // 测试各种情况
  console.log('\n1. 测试有数据的情况:')
  testFrontendProcessing(mockResponses.withData)
  
  console.log('\n2. 测试空数据的情况:')
  testFrontendProcessing(mockResponses.emptyData)
  
  console.log('\n3. 测试null数据的情况:')
  testFrontendProcessing(mockResponses.nullData)
}

// 运行测试
testApiResponse() 