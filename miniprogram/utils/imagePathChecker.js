/**
 * 图片路径检查工具
 * 用于验证图片文件是否存在，避免预加载不存在的图片
 */

class ImagePathChecker {
  constructor() {
    this.validPaths = new Set(); // 已验证的路径缓存
    this.invalidPaths = new Set(); // 已确认无效的路径缓存
  }

  /**
   * 检查图片路径是否有效
   * @param {string} path 图片路径
   * @returns {Promise<boolean>} 是否有效
   */
  async checkImagePath(path) {
    // 检查缓存
    if (this.validPaths.has(path)) {
      return true;
    }
    if (this.invalidPaths.has(path)) {
      return false;
    }

    return new Promise((resolve) => {
      wx.getImageInfo({
        src: path,
        success: () => {
          this.validPaths.add(path);
          resolve(true);
        },
        fail: () => {
          this.invalidPaths.add(path);
          resolve(false);
        }
      });
    });
  }

  /**
   * 过滤有效的图片路径
   * @param {Array} paths 图片路径数组
   * @returns {Promise<Array>} 有效的图片路径数组
   */
  async filterValidPaths(paths) {
    const validPaths = [];
    
    for (const path of paths) {
      if (await this.checkImagePath(path)) {
        validPaths.push(path);
      }
    }
    
    return validPaths;
  }

  /**
   * 获取项目中实际存在的图片路径
   * @returns {Array} 存在的图片路径数组
   */
  getExistingImagePaths() {
    return [
      // 基础图片
      '/images/default-avatar.png',
      '/images/empty-state.png',
      '/images/service-default.jpg',
      '/images/hospital-default.jpg',
      
      // 导航图标
      '/images/nav/appointment.png',
      '/images/nav/consultation.png',
      '/images/nav/health-record.png',
      '/images/nav/hospital.png',
      '/images/nav/medicine.png',
      '/images/nav/news.png',
      '/images/nav/report.png',
      
      // 服务图标
      '/images/service/appointment.png',
      '/images/service/checkup.png',
      '/images/service/consultation.png',
      '/images/service/medicine.png',
      '/images/service/record.png',
      '/images/service/report.png',
      
      // 医院图标
      '/images/hospital/dermyy-logo.png',
      '/images/hospital/etyy-logo.png',
      '/images/hospital/fybjy-logo.png',
      '/images/hospital/rmyy-logo.png',
      '/images/hospital/zyy-logo.png'
    ];
  }

  /**
   * 清理缓存
   */
  clearCache() {
    this.validPaths.clear();
    this.invalidPaths.clear();
  }

  /**
   * 获取缓存统计
   */
  getCacheStats() {
    return {
      validCount: this.validPaths.size,
      invalidCount: this.invalidPaths.size,
      validPaths: Array.from(this.validPaths),
      invalidPaths: Array.from(this.invalidPaths)
    };
  }
}

// 创建全局实例
const imagePathChecker = new ImagePathChecker();

module.exports = {
  imagePathChecker,
  checkImagePath: (path) => imagePathChecker.checkImagePath(path),
  filterValidPaths: (paths) => imagePathChecker.filterValidPaths(paths),
  getExistingImagePaths: () => imagePathChecker.getExistingImagePaths(),
  clearCache: () => imagePathChecker.clearCache(),
  getCacheStats: () => imagePathChecker.getCacheStats()
};
