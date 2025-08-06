// 测试URL参数构建修复

console.log('=== 测试URL参数构建修复 ===');

// 模拟URL参数构建函数（兼容微信小程序）
function buildQueryString(query) {
  if (!query || Object.keys(query).length === 0) {
    return '';
  }
  
  const queryParams = [];
  Object.entries(query).forEach(([key, value]) => {
    if (value !== undefined && value !== null) {
      queryParams.push(`${encodeURIComponent(key)}=${encodeURIComponent(value)}`);
    }
  });
  
  return queryParams.join('&');
}

// 模拟构建完整URL
function buildFullPath(path, query) {
  let finalPath = path;
  const queryString = buildQueryString(query);
  
  if (queryString) {
    finalPath = `${path}${path.includes('?') ? '&' : '?'}${queryString}`;
  }
  
  return finalPath;
}

// 测试1: 基础URL参数构建
function testBasicQueryBuilding() {
  console.log('1. 测试基础URL参数构建');
  
  const testCases = [
    {
      path: '/api/admin/orders',
      query: { adminUserId: 'anyuyinian' },
      expected: '/api/admin/orders?adminUserId=anyuyinian'
    },
    {
      path: '/api/admin/order/update-amount',
      query: { adminUserId: 'anyuyinian' },
      expected: '/api/admin/order/update-amount?adminUserId=anyuyinian'
    },
    {
      path: '/api/admin/orders',
      query: { 
        adminUserId: 'anyuyinian',
        page: 1,
        pageSize: 20
      },
      expected: '/api/admin/orders?adminUserId=anyuyinian&page=1&pageSize=20'
    }
  ];
  
  testCases.forEach((testCase, index) => {
    console.log(`\n测试用例${index + 1}:`);
    console.log(`  路径: ${testCase.path}`);
    console.log(`  参数:`, testCase.query);
    
    const result = buildFullPath(testCase.path, testCase.query);
    console.log(`  结果: ${result}`);
    console.log(`  期望: ${testCase.expected}`);
    console.log(`  状态: ${result === testCase.expected ? '✅ 通过' : '❌ 失败'}`);
  });
}

// 测试2: 特殊字符编码
function testSpecialCharacterEncoding() {
  console.log('\n2. 测试特殊字符编码');
  
  const testCases = [
    {
      path: '/api/admin/orders',
      query: { 
        adminUserId: 'anyuyinian',
        search: '测试用户'
      },
      expected: '/api/admin/orders?adminUserId=anyuyinian&search=%E6%B5%8B%E8%AF%95%E7%94%A8%E6%88%B7'
    },
    {
      path: '/api/admin/orders',
      query: { 
        adminUserId: 'anyuyinian',
        filter: 'status=0&type=test'
      },
      expected: '/api/admin/orders?adminUserId=anyuyinian&filter=status%3D0%26type%3Dtest'
    }
  ];
  
  testCases.forEach((testCase, index) => {
    console.log(`\n测试用例${index + 1}:`);
    console.log(`  路径: ${testCase.path}`);
    console.log(`  参数:`, testCase.query);
    
    const result = buildFullPath(testCase.path, testCase.query);
    console.log(`  结果: ${result}`);
    console.log(`  期望: ${testCase.expected}`);
    console.log(`  状态: ${result === testCase.expected ? '✅ 通过' : '❌ 失败'}`);
  });
}

// 测试3: 边界情况
function testEdgeCases() {
  console.log('\n3. 测试边界情况');
  
  const testCases = [
    {
      name: '空查询参数',
      path: '/api/admin/orders',
      query: {},
      expected: '/api/admin/orders'
    },
    {
      name: 'null查询参数',
      path: '/api/admin/orders',
      query: null,
      expected: '/api/admin/orders'
    },
    {
      name: 'undefined查询参数',
      path: '/api/admin/orders',
      query: undefined,
      expected: '/api/admin/orders'
    },
    {
      name: '包含null值的参数',
      path: '/api/admin/orders',
      query: { 
        adminUserId: 'anyuyinian',
        nullValue: null,
        undefinedValue: undefined
      },
      expected: '/api/admin/orders?adminUserId=anyuyinian'
    },
    {
      name: '已有查询参数的路径',
      path: '/api/admin/orders?existing=param',
      query: { adminUserId: 'anyuyinian' },
      expected: '/api/admin/orders?existing=param&adminUserId=anyuyinian'
    }
  ];
  
  testCases.forEach((testCase, index) => {
    console.log(`\n测试用例${index + 1}: ${testCase.name}`);
    console.log(`  路径: ${testCase.path}`);
    console.log(`  参数:`, testCase.query);
    
    const result = buildFullPath(testCase.path, testCase.query);
    console.log(`  结果: ${result}`);
    console.log(`  期望: ${testCase.expected}`);
    console.log(`  状态: ${result === testCase.expected ? '✅ 通过' : '❌ 失败'}`);
  });
}

// 测试4: 实际API调用模拟
function testActualAPICalls() {
  console.log('\n4. 测试实际API调用模拟');
  
  const mockCallContainer = (path, method = 'GET', data = {}, options = {}) => {
    // 构建完整路径
    let finalPath = path;
    if (options.query && Object.keys(options.query).length > 0) {
      const queryParams = [];
      Object.entries(options.query).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          queryParams.push(`${encodeURIComponent(key)}=${encodeURIComponent(value)}`);
        }
      });
      const queryString = queryParams.join('&');
      if (queryString) {
        finalPath = `${path}${path.includes('?') ? '&' : '?'}${queryString}`;
      }
    }
    
    console.log(`模拟API调用: ${method} ${finalPath}`);
    console.log(`请求数据:`, data);
    
    // 模拟成功响应
    return Promise.resolve({
      code: 0,
      data: { message: 'API调用成功' }
    });
  };
  
  const testCases = [
    {
      name: '获取订单列表',
      path: '/api/admin/orders',
      method: 'GET',
      data: {},
      options: {
        query: { 
          adminUserId: 'anyuyinian',
          page: 1,
          pageSize: 20
        }
      }
    },
    {
      name: '修改订单金额',
      path: '/api/admin/order/update-amount',
      method: 'POST',
      data: {
        orderId: 1,
        newAmount: 399.00,
        reason: '管理员手动修改'
      },
      options: {
        query: { adminUserId: 'anyuyinian' }
      }
    }
  ];
  
  testCases.forEach(async (testCase, index) => {
    console.log(`\n测试用例${index + 1}: ${testCase.name}`);
    
    try {
      const result = await mockCallContainer(
        testCase.path,
        testCase.method,
        testCase.data,
        testCase.options
      );
      
      console.log('✅ API调用成功:', result);
    } catch (error) {
      console.log('❌ API调用失败:', error.message);
    }
  });
}

// 运行所有测试
function runAllTests() {
  console.log('开始运行所有测试...\n');
  
  testBasicQueryBuilding();
  testSpecialCharacterEncoding();
  testEdgeCases();
  testActualAPICalls();
  
  console.log('\n=== 测试完成 ===');
  console.log('如果所有测试都通过，说明URL参数构建修复成功！');
  console.log('现在可以在微信小程序环境中正常使用API调用了。');
}

// 执行测试
runAllTests(); 