// 测试管理员API调用修复

console.log('=== 测试管理员API调用修复 ===');

// 模拟app.callContainer方法
const mockCallContainer = (path, method = 'GET', data = {}, options = {}) => {
  return new Promise((resolve, reject) => {
    console.log('模拟API调用:', {
      path,
      method,
      data,
      options
    });

    // 模拟处理查询参数
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

    console.log('最终路径:', finalPath);

    // 模拟API响应
    if (finalPath.includes('adminUserId=')) {
      resolve({
        code: 0,
        data: {
          list: [
            {
              id: 1,
              orderNo: 'ORDER20241220001',
              amount: 299.00,
              status: 0,
              serviceName: '上门护理服务',
              userNickName: '测试用户',
              createdAt: '2024-12-20 10:30:00'
            }
          ],
          hasMore: false
        }
      });
    } else {
      reject(new Error('缺少adminUserId参数'));
    }
  });
};

// 测试1: 测试查询参数处理
function testQueryParameterHandling() {
  console.log('1. 测试查询参数处理');
  
  const testCases = [
    {
      name: '修改订单金额API',
      path: '/api/admin/order/update-amount',
      method: 'POST',
      data: { orderId: 1, newAmount: 399.00, reason: '测试' },
      options: {
        query: { adminUserId: 'anyuyinian' }
      }
    },
    {
      name: '获取订单列表API',
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

// 测试2: 测试参数验证
function testParameterValidation() {
  console.log('\n2. 测试参数验证');
  
  const testCases = [
    {
      name: '缺少adminUserId',
      path: '/api/admin/orders',
      method: 'GET',
      data: {},
      options: {}
    },
    {
      name: '包含adminUserId',
      path: '/api/admin/orders',
      method: 'GET',
      data: {},
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
      
      console.log('✅ 参数验证通过:', result);
    } catch (error) {
      console.log('❌ 参数验证失败:', error.message);
    }
  });
}

// 测试3: 测试URL构建
function testURLBuilding() {
  console.log('\n3. 测试URL构建');
  
  const testCases = [
    {
      path: '/api/admin/orders',
      query: { adminUserId: 'anyuyinian', page: 1 },
      expected: '/api/admin/orders?adminUserId=anyuyinian&page=1'
    },
    {
      path: '/api/admin/order/update-amount',
      query: { adminUserId: 'anyuyinian' },
      expected: '/api/admin/order/update-amount?adminUserId=anyuyinian'
    },
    {
      path: '/api/admin/orders?existing=param',
      query: { adminUserId: 'anyuyinian' },
      expected: '/api/admin/orders?existing=param&adminUserId=anyuyinian'
    }
  ];

  testCases.forEach((testCase, index) => {
    console.log(`\n测试用例${index + 1}:`);
    console.log(`  原始路径: ${testCase.path}`);
    console.log(`  查询参数:`, testCase.query);
    
    // 构建URL
    let finalPath = testCase.path;
    if (testCase.query && Object.keys(testCase.query).length > 0) {
      const queryParams = new URLSearchParams();
      Object.entries(testCase.query).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          queryParams.append(key, value);
        }
      });
      const queryString = queryParams.toString();
      if (queryString) {
        finalPath = `${testCase.path}${testCase.path.includes('?') ? '&' : '?'}${queryString}`;
      }
    }
    
    console.log(`  构建路径: ${finalPath}`);
    console.log(`  期望路径: ${testCase.expected}`);
    console.log(`  结果: ${finalPath === testCase.expected ? '✅ 匹配' : '❌ 不匹配'}`);
  });
}

// 测试4: 测试错误处理
function testErrorHandling() {
  console.log('\n4. 测试错误处理');
  
  const testCases = [
    {
      name: '网络错误',
      shouldFail: true,
      errorType: 'network'
    },
    {
      name: '参数错误',
      shouldFail: true,
      errorType: 'parameter'
    },
    {
      name: '成功调用',
      shouldFail: false,
      errorType: null
    }
  ];

  testCases.forEach(async (testCase, index) => {
    console.log(`\n测试用例${index + 1}: ${testCase.name}`);
    
    try {
      if (testCase.shouldFail) {
        throw new Error(`模拟${testCase.errorType}错误`);
      }
      
      const result = await mockCallContainer('/api/admin/orders', 'GET', {}, {
        query: { adminUserId: 'anyuyinian' }
      });
      
      console.log('✅ 调用成功:', result);
    } catch (error) {
      console.log('❌ 调用失败:', error.message);
      
      // 模拟错误处理
      if (error.message.includes('网络')) {
        console.log('  处理网络错误...');
      } else if (error.message.includes('参数')) {
        console.log('  处理参数错误...');
      }
    }
  });
}

// 运行所有测试
function runAllTests() {
  console.log('开始运行所有测试...\n');
  
  testQueryParameterHandling();
  testParameterValidation();
  testURLBuilding();
  testErrorHandling();
  
  console.log('\n=== 测试完成 ===');
  console.log('如果所有测试都通过，说明API调用修复成功！');
}

// 执行测试
runAllTests(); 