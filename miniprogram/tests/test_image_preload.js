/**
 * 图片预加载功能测试脚本
 * 用于验证图片预加载功能是否正常工作
 */

const { imagePreloader } = require('../utils/imagePreloader')

// 测试图片列表
const testImages = [
  '/images/service/default-service.jpg',
  '/images/empty-state.png',
  '/images/default-avatar.png',
  '/images/service-default.jpg'
]

// 测试单张图片预加载
async function testSingleImagePreload() {
  console.log('=== 测试单张图片预加载 ===')
  
  try {
    const result = await imagePreloader.preloadImage(testImages[0])
    console.log('单张图片预加载成功:', result)
    return true
  } catch (error) {
    console.error('单张图片预加载失败:', error)
    return false
  }
}

// 测试批量图片预加载
async function testBatchImagePreload() {
  console.log('=== 测试批量图片预加载 ===')
  
  try {
    const results = await imagePreloader.preloadImages(
      testImages,
      (progress) => {
        console.log(`预加载进度: ${progress.progress}% (${progress.loaded}/${progress.total})`)
      },
      (result) => {
        console.log('批量预加载完成:', result)
      }
    )
    
    console.log('批量预加载结果:', results)
    return true
  } catch (error) {
    console.error('批量预加载失败:', error)
    return false
  }
}

// 测试缓存功能
function testCacheFunction() {
  console.log('=== 测试缓存功能 ===')
  
  // 获取缓存信息
  const cacheInfo = imagePreloader.getCacheInfo()
  console.log('缓存信息:', cacheInfo)
  
  // 检查图片是否已缓存
  testImages.forEach(src => {
    const isCached = imagePreloader.isCached(src)
    console.log(`图片 ${src} 是否已缓存:`, isCached)
    
    if (isCached) {
      const cachedImage = imagePreloader.getCachedImage(src)
      console.log(`缓存详情:`, cachedImage)
    }
  })
  
  return true
}

// 测试关键图片预加载
async function testCriticalImagePreload() {
  console.log('=== 测试关键图片预加载 ===')
  
  try {
    const results = await imagePreloader.preloadCriticalImages(
      (progress) => {
        console.log(`关键图片预加载进度: ${progress.progress}% (${progress.loaded}/${progress.total})`)
      },
      (result) => {
        console.log('关键图片预加载完成:', result)
      }
    )
    
    console.log('关键图片预加载结果:', results)
    return true
  } catch (error) {
    console.error('关键图片预加载失败:', error)
    return false
  }
}

// 测试首页图片预加载
async function testHomeImagePreload() {
  console.log('=== 测试首页图片预加载 ===')
  
  const mockHomeData = {
    services: [
      { imageUrl: '/images/service/default-service.jpg' },
      { imageUrl: '/images/service/service-1.jpg' },
      { imageUrl: '/images/service/service-2.jpg' }
    ]
  }
  
  try {
    const results = await imagePreloader.preloadHomeImages(
      mockHomeData,
      (progress) => {
        console.log(`首页图片预加载进度: ${progress.progress}% (${progress.loaded}/${progress.total})`)
      },
      (result) => {
        console.log('首页图片预加载完成:', result)
      }
    )
    
    console.log('首页图片预加载结果:', results)
    return true
  } catch (error) {
    console.error('首页图片预加载失败:', error)
    return false
  }
}

// 测试服务页面图片预加载
async function testServiceImagePreload() {
  console.log('=== 测试服务页面图片预加载 ===')
  
  const mockServiceData = {
    imageUrl: '/images/service/service-detail.jpg'
  }
  
  try {
    const results = await imagePreloader.preloadServiceImages(
      mockServiceData,
      (progress) => {
        console.log(`服务页面图片预加载进度: ${progress.progress}% (${progress.loaded}/${progress.total})`)
      },
      (result) => {
        console.log('服务页面图片预加载完成:', result)
      }
    )
    
    console.log('服务页面图片预加载结果:', results)
    return true
  } catch (error) {
    console.error('服务页面图片预加载失败:', error)
    return false
  }
}

// 运行所有测试
async function runAllTests() {
  console.log('🚀 开始图片预加载功能测试...\n')
  
  const tests = [
    { name: '单张图片预加载', fn: testSingleImagePreload },
    { name: '批量图片预加载', fn: testBatchImagePreload },
    { name: '缓存功能', fn: testCacheFunction },
    { name: '关键图片预加载', fn: testCriticalImagePreload },
    { name: '首页图片预加载', fn: testHomeImagePreload },
    { name: '服务页面图片预加载', fn: testServiceImagePreload }
  ]
  
  const results = []
  
  for (const test of tests) {
    console.log(`\n📋 运行测试: ${test.name}`)
    try {
      const success = await test.fn()
      results.push({ name: test.name, success })
      console.log(`${success ? '✅' : '❌'} ${test.name}: ${success ? '通过' : '失败'}`)
    } catch (error) {
      console.error(`❌ ${test.name}: 异常 -`, error)
      results.push({ name: test.name, success: false, error: error.message })
    }
  }
  
  // 输出测试结果摘要
  console.log('\n📊 测试结果摘要:')
  console.log('=' * 50)
  
  const passed = results.filter(r => r.success).length
  const total = results.length
  
  results.forEach(result => {
    const status = result.success ? '✅ 通过' : '❌ 失败'
    console.log(`${status} ${result.name}`)
    if (result.error) {
      console.log(`   错误: ${result.error}`)
    }
  })
  
  console.log(`\n🎯 总计: ${passed}/${total} 个测试通过`)
  
  if (passed === total) {
    console.log('🎉 所有测试都通过了！图片预加载功能工作正常。')
  } else {
    console.log('⚠️  部分测试失败，请检查相关功能。')
  }
  
  return results
}

// 导出测试函数
module.exports = {
  testSingleImagePreload,
  testBatchImagePreload,
  testCacheFunction,
  testCriticalImagePreload,
  testHomeImagePreload,
  testServiceImagePreload,
  runAllTests
}

// 如果直接运行此文件，则执行所有测试
if (typeof wx !== 'undefined') {
  // 在微信小程序环境中运行
  runAllTests().catch(console.error)
} else {
  // 在Node.js环境中运行
  console.log('请在微信开发者工具中运行此测试脚本')
}
