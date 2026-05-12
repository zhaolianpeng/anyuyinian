/**
 * 用户ID兼容性处理工具
 * 处理从数字ID到字符串ID的迁移
 */

// 检查是否为有效的MongoDB ID格式
const isValidMongoId = (id) => {
  if (!id || typeof id !== 'string') return false
  return /^[0-9a-fA-F]{24}$/.test(id)
}

// 检查是否为数字ID
const isNumericId = (id) => {
  if (!id) return false
  return !isNaN(Number(id)) && Number.isInteger(Number(id))
}

// 格式化用户ID，确保是字符串格式
const formatUserId = (userId) => {
  if (!userId) return null
  
  // 如果已经是有效的MongoDB ID格式，直接返回
  if (isValidMongoId(userId)) {
    return userId
  }
  
  // 如果是数字ID，转换为字符串
  if (isNumericId(userId)) {
    return String(userId)
  }
  
  // 其他情况，转换为字符串
  return String(userId)
}

// 验证用户ID格式
const validateUserId = (userId) => {
  const formatted = formatUserId(userId)
  if (!formatted) {
    console.warn('无效的用户ID:', userId)
    return false
  }
  
  // 检查是否为有效的MongoDB ID格式
  if (!isValidMongoId(formatted)) {
    console.warn('用户ID不是有效的MongoDB格式:', formatted)
    return false
  }
  
  return true
}

// 获取当前用户ID
const getCurrentUserId = () => {
  const userId = wx.getStorageSync('userId')
  const formatted = formatUserId(userId)
  console.log('获取当前用户ID:', { original: userId, formatted: formatted })
  return formatted
}

// 设置用户ID
const setUserId = (userId) => {
  const formatted = formatUserId(userId)
  if (formatted) {
    wx.setStorageSync('userId', formatted)
    console.log('用户ID已设置:', formatted)
    return true
  }
  console.warn('设置用户ID失败，无效的userId:', userId)
  return false
}

// 清除用户ID（添加保护机制）
const clearUserId = (reason = '未知原因') => {
  const currentUserId = wx.getStorageSync('userId')
  console.log(`准备清除用户ID，原因: ${reason}，当前用户ID:`, currentUserId)
  
  // 如果是有效的MongoDB ID，记录警告但不清除
  if (currentUserId && isValidMongoId(currentUserId)) {
    console.warn('检测到有效的用户ID被清除，这可能是一个错误:', currentUserId)
    console.warn('清除原因:', reason)
  }
  
  wx.removeStorageSync('userId')
  console.log('用户ID已清除，原因:', reason)
}

// 检查是否需要迁移用户ID
const needsUserIdMigration = () => {
  const userId = wx.getStorageSync('userId')
  const needsMigration = isNumericId(userId)
  console.log('检查用户ID迁移需求:', { userId, needsMigration })
  return needsMigration
}

// 安全清除用户ID（仅在特定情况下使用）
const safeClearUserId = (reason = '安全清除') => {
  const currentUserId = wx.getStorageSync('userId')
  
  // 只有在用户ID无效或为空时才清除
  if (!currentUserId || !isValidMongoId(currentUserId)) {
    clearUserId(reason)
    return true
  }
  
  console.log('跳过清除有效用户ID:', currentUserId)
  return false
}

module.exports = {
  isValidMongoId,
  isNumericId,
  formatUserId,
  validateUserId,
  getCurrentUserId,
  setUserId,
  clearUserId,
  needsUserIdMigration,
  safeClearUserId
} 