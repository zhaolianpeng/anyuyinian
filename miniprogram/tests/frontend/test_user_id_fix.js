/**
 * 用户ID修复测试脚本
 * 测试用户ID兼容性工具的各种功能
 */

const {
  isValidMongoId,
  isNumericId,
  formatUserId,
  validateUserId,
  getCurrentUserId,
  setUserId,
  clearUserId,
  needsUserIdMigration,
  safeClearUserId
} = require('../../utils/user-id-compatibility')

// 测试数据
const testCases = [
  // 有效的MongoDB ID
  '507f1f77bcf86cd799439011',
  '507f1f77bcf86cd799439012',
  '507f1f77bcf86cd799439013',
  
  // 数字ID
  '1',
  '123',
  '999',
  
  // 无效ID
  'invalid_id',
  '',
  null,
  undefined
]

// 测试函数
function testUserIdCompatibility() {
  console.log('🧪 开始用户ID兼容性测试')
  
  // 测试MongoDB ID验证
  console.log('\n📋 测试MongoDB ID验证:')
  testCases.forEach(id => {
    const isValid = isValidMongoId(id)
    console.log(`${isValid ? '✅' : '❌'} ${id}: ${isValid}`)
  })
  
  // 测试数字ID检测
  console.log('\n📋 测试数字ID检测:')
  testCases.forEach(id => {
    const isNumeric = isNumericId(id)
    console.log(`${isNumeric ? '✅' : '❌'} ${id}: ${isNumeric}`)
  })
  
  // 测试ID格式化
  console.log('\n📋 测试ID格式化:')
  testCases.forEach(id => {
    const formatted = formatUserId(id)
    console.log(`${id} -> ${formatted}`)
  })
  
  // 测试ID验证
  console.log('\n📋 测试ID验证:')
  const validId = '507f1f77bcf86cd799439011'
  const invalidId = 'invalid_id'
  console.log('✅ 有效ID验证:', validateUserId(validId))
  console.log('❌ 无效ID验证:', validateUserId(invalidId))
  
  // 测试用户ID设置和获取
  console.log('\n📋 测试用户ID设置和获取:')
  const testUserId = '507f1f77bcf86cd799439011'
  setUserId(testUserId)
  const retrievedUserId = getCurrentUserId()
  console.log(`设置: ${testUserId}`)
  console.log(`获取: ${retrievedUserId}`)
  console.log(`匹配: ${testUserId === retrievedUserId ? '✅' : '❌'}`)
  
  // 测试迁移检测
  console.log('\n📋 测试迁移检测:')
  console.log('数字ID需要迁移:', needsUserIdMigration())
  
  // 测试安全清除
  console.log('\n📋 测试安全清除:')
  console.log('清除前用户ID:', getCurrentUserId())
  const cleared = safeClearUserId('测试清除')
  console.log('安全清除结果:', cleared)
  console.log('清除后用户ID:', getCurrentUserId())
  
  // 测试强制清除
  console.log('\n📋 测试强制清除:')
  setUserId('507f1f77bcf86cd799439011')
  console.log('清除前用户ID:', getCurrentUserId())
  clearUserId('测试强制清除')
  console.log('清除后用户ID:', getCurrentUserId())
  
  console.log('\n🎉 用户ID兼容性测试完成')
}

// 运行测试
testUserIdCompatibility()

module.exports = {
  testUserIdCompatibility
} 