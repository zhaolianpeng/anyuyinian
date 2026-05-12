// 医院图片处理测试脚本
const { processHospitalLogo, processHospitalImages, processSingleHospitalImage } = require('../utils/image')

/**
 * 测试医院图片URL处理
 */
function testHospitalLogoProcessing() {
  console.log('=== 测试医院图片URL处理 ===')
  
  const testCases = [
    {
      input: '/images/hospital/rmyy-logo.png',
      expected: '/images/hospital/rmyy-logo.png',
      description: '相对路径'
    },
    {
      input: 'https://example.com/hospital-logo.jpg',
      expected: 'https://example.com/hospital-logo.jpg',
      description: '完整URL'
    },
    {
      input: '',
      expected: '/images/hospital-default.jpg',
      description: '空字符串'
    },
    {
      input: null,
      expected: '/images/hospital-default.jpg',
      description: 'null值'
    },
    {
      input: undefined,
      expected: '/images/hospital-default.jpg',
      description: 'undefined值'
    },
    {
      input: 'invalid-url',
      expected: '/images/hospital-default.jpg',
      description: '无效URL'
    }
  ]
  
  testCases.forEach((testCase, index) => {
    const result = processHospitalLogo(testCase.input)
    const passed = result === testCase.expected
    
    console.log(`测试 ${index + 1} (${testCase.description}):`)
    console.log(`  输入: ${testCase.input}`)
    console.log(`  期望: ${testCase.expected}`)
    console.log(`  实际: ${result}`)
    console.log(`  结果: ${passed ? '✅ 通过' : '❌ 失败'}`)
    console.log('')
  })
}

/**
 * 测试医院数据图片处理
 */
function testHospitalDataProcessing() {
  console.log('=== 测试医院数据图片处理 ===')
  
  const testHospitals = [
    {
      id: 1,
      name: '深圳市人民医院',
      logo: '/images/hospital/rmyy-logo.png',
      address: '深圳市罗湖区东门北路1017号'
    },
    {
      id: 2,
      name: '深圳市第二人民医院',
      logo: 'https://example.com/hospital-logo.jpg',
      address: '深圳市福田区笋岗西路3002号'
    },
    {
      id: 3,
      name: '深圳市中医院',
      logo: '',
      address: '深圳市福田区福华路1号'
    },
    {
      id: 4,
      name: '深圳市儿童医院',
      logo: null,
      address: '深圳市福田区益田路7019号'
    }
  ]
  
  console.log('原始数据:')
  testHospitals.forEach(hospital => {
    console.log(`- ${hospital.name}: ${hospital.logo}`)
  })
  
  const processedHospitals = processHospitalImages(testHospitals)
  
  console.log('\n处理后数据:')
  processedHospitals.forEach(hospital => {
    console.log(`- ${hospital.name}: ${hospital.logo}`)
  })
}

/**
 * 测试单个医院图片处理
 */
function testSingleHospitalProcessing() {
  console.log('\n=== 测试单个医院图片处理 ===')
  
  const testHospital = {
    id: 1,
    name: '深圳市人民医院',
    logo: '/images/hospital/rmyy-logo.png',
    address: '深圳市罗湖区东门北路1017号'
  }
  
  console.log('原始数据:')
  console.log(`- ${testHospital.name}: ${testHospital.logo}`)
  
  const processedHospital = processSingleHospitalImage(testHospital)
  
  console.log('处理后数据:')
  console.log(`- ${processedHospital.name}: ${processedHospital.logo}`)
}

/**
 * 运行所有测试
 */
function runAllTests() {
  console.log('🏥 开始医院图片处理功能测试\n')
  
  testHospitalLogoProcessing()
  testHospitalDataProcessing()
  testSingleHospitalProcessing()
  
  console.log('🏥 医院图片处理功能测试完成')
}

// 如果直接运行此脚本
if (typeof module !== 'undefined' && module.exports) {
  module.exports = {
    testHospitalLogoProcessing,
    testHospitalDataProcessing,
    testSingleHospitalProcessing,
    runAllTests
  }
}

// 在微信小程序环境中运行测试
if (typeof wx !== 'undefined') {
  // 延迟执行，确保页面加载完成
  setTimeout(() => {
    runAllTests()
  }, 1000)
} 