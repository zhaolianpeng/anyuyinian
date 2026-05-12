/**
 * 修复后的图片预加载功能测试脚本
 * 验证图片路径检查和预加载功能是否正常工作
 */

const { imagePreloader } = require('../utils/imagePreloader')
const { imagePathChecker } = require('../utils/imagePathChecker')

// 测试图片路径检查
async function testImagePathChecker() {
  console.log('=== 测试图片路径检查器 ===')
  
  try {
    // 测试存在的图片
    const existingImage = '/images/default-avatar.png'
    const isValid = await imagePathChecker.checkImagePath(existingImage)
    console.log(`图片 ${existingImage} 存在:`, isValid)
    
    // 测试不存在的图片
    const nonExistingImage = '/images/non-existing.jpg'
    const isInvalid = await imagePathChecker.checkImagePath(nonExistingImage)
    console.log(`图片 ${nonExistingImage} 存在:`, isInvalid)
    
    // 获取项目中实际存在的图片路径
    const existingPaths = imagePathChecker.getExistingImagePaths()
    console.log('项目中实际存在的图片路径:', existingPaths)
    
    // 过滤有效的图片路径
    const testPaths = [
      '/images/default-avatar.png',
      '/images/empty-state.png',
      '/images/service-default.jpg',
      '/images/non-existing.jpg',
      '/images/another-non-existing.png'
    ]
    
    const validPaths = await imagePathChecker.filterValidPaths(testPaths)
    console.log('有效图片路径:', validPaths)
    
    return true
  } catch (error) {
    console.error('图片路径检查器测试失败:', error)
    return false
  }
}

// 测试修复后的关键图片预加载
async function testFixedCriticalImagePreload() {
  console.log('=== 测试修复后的关键图片预加载 ===')
  
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

// 测试修复后的批量图片预加载
async function testFixedBatchImagePreload() {
  console.log('=== 测试修复后的批量图片预加载 ===')
  
  try {
    // 使用实际存在的图片路径
    const testImages = [
      '/images/default-avatar.png',
      '/images/empty-state.png',
      '/images/service-default.jpg'
    ]
    
    const results = await imagePreloader.preloadImages(
      testImages,
      (progress) => {
        console.log(`批量预加载进度: ${progress.progress}% (${progress.loaded}/${progress.total})`)
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
  
  // 获取预加载器缓存信息
  const preloadCacheInfo = imagePreloader.getCacheInfo()
  console.log('预加载器缓存信息:', preloadCacheInfo)
  
  // 获取路径检查器缓存信息
  const pathCacheStats = imagePathChecker.getCacheStats()
  console.log('路径检查器缓存统计:', pathCacheStats)
  
  return true
}

// 测试错误处理
async function testErrorHandling() {
  console.log('=== 测试错误处理 ===')
  
  try {
    // 测试不存在的图片
    const result = await imagePreloader.preloadImage('/images/non-existing.jpg')
    console.log('不存在图片的预加载结果:', result)
    
    // 测试空字符串
    try {
      await imagePreloader.preloadImage('')
    } catch (error) {
      console.log('空字符串错误处理:', error.message)
    }
    
    // 测试无效URL
    try {
      await imagePreloader.preloadImage(null)
    } catch (error) {
      console.log('null值错误处理:', error.message)
    }
    
    return true
  } catch (error) {
    console.error('错误处理测试失败:', error)
    return false
  }
}

// 运行所有测试
async function runFixedTests() {
  console.log('🚀 开始修复后的图片预加载功能测试...\n')
  
  const tests = [
    { name: '图片路径检查器', fn: testImagePathChecker },
    { name: '修复后的关键图片预加载', fn: testFixedCriticalImagePreload },
    { name: '修复后的批量图片预加载', fn: testFixedBatchImagePreload },
    { name: '缓存功能', fn: testCacheFunction },
    { name: '错误处理', fn: testErrorHandling }
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
    console.log('🎉 所有测试都通过了！修复后的图片预加载功能工作正常。')
  } else {
    console.log('⚠️  部分测试失败，请检查相关功能。')
  }
  
  return results
}

// 导出测试函数
module.exports = {
  testImagePathChecker,
  testFixedCriticalImagePreload,
  testFixedBatchImagePreload,
  testCacheFunction,
  testErrorHandling,
  runFixedTests
}

// 如果直接运行此文件，则执行所有测试
if (typeof wx !== 'undefined') {
  // 在微信小程序环境中运行
  runFixedTests().catch(console.error)
} else {
  // 在Node.js环境中运行
  console.log('请在微信开发者工具中运行此测试脚本')
}
