/**
 * 图片预加载工具类
 * 提供图片预加载、缓存管理、进度监控等功能
 */

const { imagePathChecker } = require('./imagePathChecker')

function shouldSkipRemotePreload(src) {
  if (!src || typeof src !== 'string') {
    return false
  }

  return src.startsWith('https://api.succ.online/anyuyinian/')
}

class ImagePreloader {
  constructor() {
    this.cache = new Map(); // 图片缓存
    this.loadingQueue = []; // 加载队列
    this.loadedCount = 0; // 已加载数量
    this.totalCount = 0; // 总数量
    this.isLoading = false; // 是否正在加载
    this.maxConcurrent = 3; // 最大并发数
    this.timeout = 10000; // 超时时间（毫秒）
  }

  /**
   * 预加载单张图片
   * @param {string} src 图片地址
   * @param {Object} options 配置选项
   * @returns {Promise} 加载结果
   */
  preloadImage(src, options = {}) {
    return new Promise((resolve, reject) => {
      // 检查缓存
      if (this.cache.has(src)) {
        resolve({
          src,
          success: true,
          fromCache: true,
          timestamp: this.cache.get(src).timestamp
        });
        return;
      }

      // 检查图片是否有效
      if (!src || src === '') {
        reject(new Error('图片地址无效'));
        return;
      }

      if (shouldSkipRemotePreload(src)) {
        resolve({
          src,
          success: false,
          skipped: true,
          fromCache: false,
          error: 'skip remote preload'
        })
        return
      }

      const startTime = Date.now();
      const timeoutId = setTimeout(() => {
        reject(new Error(`图片加载超时: ${src}`));
      }, this.timeout);

      // 使用 wx.getImageInfo 预加载图片
      wx.getImageInfo({
        src: src,
        success: (res) => {
          clearTimeout(timeoutId);
          const loadTime = Date.now() - startTime;
          
          // 缓存图片信息
          this.cache.set(src, {
            width: res.width,
            height: res.height,
            path: res.path,
            timestamp: Date.now(),
            loadTime: loadTime
          });

          resolve({
            src,
            success: true,
            fromCache: false,
            loadTime: loadTime,
            width: res.width,
            height: res.height
          });
        },
        fail: (error) => {
          clearTimeout(timeoutId);
          console.warn('图片预加载失败:', src, error);
          // 对于文件不存在的错误，我们静默处理，不抛出异常
          if (error.errMsg && error.errMsg.includes('file not found')) {
            console.warn(`图片文件不存在，跳过预加载: ${src}`);
            resolve({
              src,
              success: false,
              fromCache: false,
              error: 'file not found',
              loadTime: Date.now() - startTime
            });
          } else {
            reject(new Error(`图片加载失败: ${src} - ${error.errMsg}`));
          }
        }
      });
    });
  }

  /**
   * 批量预加载图片
   * @param {Array} imageList 图片地址数组
   * @param {Function} onProgress 进度回调
   * @param {Function} onComplete 完成回调
   * @returns {Promise} 加载结果
   */
  async preloadImages(imageList, onProgress, onComplete) {
    if (!Array.isArray(imageList) || imageList.length === 0) {
      onComplete && onComplete([]);
      return [];
    }

    // 过滤有效图片地址
    const validImages = imageList.filter(src => src && src.trim() !== '');
    this.totalCount = validImages.length;
    this.loadedCount = 0;
    this.isLoading = true;

    const results = [];
    const errors = [];

    // 分批加载，控制并发数
    for (let i = 0; i < validImages.length; i += this.maxConcurrent) {
      const batch = validImages.slice(i, i + this.maxConcurrent);
      const batchPromises = batch.map(src => 
        this.preloadImage(src).catch(error => ({ 
          src, 
          success: false, 
          error: error.message,
          fromCache: false
        }))
      );

      try {
        const batchResults = await Promise.all(batchPromises);
        results.push(...batchResults);
        
        // 更新进度
        this.loadedCount += batch.length;
        const progress = Math.round((this.loadedCount / this.totalCount) * 100);
        
        onProgress && onProgress({
          loaded: this.loadedCount,
          total: this.totalCount,
          progress: progress,
          current: validImages[Math.min(i + this.maxConcurrent - 1, validImages.length - 1)]
        });

        // 添加小延迟，避免阻塞UI
        if (i + this.maxConcurrent < validImages.length) {
          await new Promise(resolve => setTimeout(resolve, 50));
        }
      } catch (error) {
        console.error('批量预加载出错:', error);
        errors.push(error);
      }
    }

    this.isLoading = false;
    const successResults = results.filter(r => r.success);
    const failedResults = results.filter(r => !r.success);

    onComplete && onComplete({
      success: successResults,
      failed: failedResults,
      total: this.totalCount,
      successCount: successResults.length,
      failedCount: failedResults.length
    });

    return results;
  }

  /**
   * 预加载首页图片
   * @param {Object} homeData 首页数据
   * @param {Function} onProgress 进度回调
   * @param {Function} onComplete 完成回调
   */
  async preloadHomeImages(homeData, onProgress, onComplete) {
    const imageList = [];

    // 添加服务列表图片
    if (homeData.services && Array.isArray(homeData.services)) {
      homeData.services.forEach(service => {
        if (service.imageUrl) {
          imageList.push(service.imageUrl);
        }
      });
    }

    // 添加默认图片
    const defaultImages = [
      '/images/service-default.jpg',
      '/images/empty-state.png',
      '/images/default-avatar.png'
    ];
    
    // 过滤有效的默认图片
    const validDefaultImages = await imagePathChecker.filterValidPaths(defaultImages);
    imageList.push(...validDefaultImages);

    return this.preloadImages(imageList, onProgress, onComplete);
  }

  /**
   * 预加载服务页面图片
   * @param {Object} serviceData 服务数据
   * @param {Function} onProgress 进度回调
   * @param {Function} onComplete 完成回调
   */
  async preloadServiceImages(serviceData, onProgress, onComplete) {
    const imageList = [];

    // 添加服务详情图片
    if (serviceData.imageUrl) {
      imageList.push(serviceData.imageUrl);
    }

    // 添加默认图片
    const defaultImages = [
      '/images/service-default.jpg',
      '/images/empty-state.png'
    ];
    
    // 过滤有效的默认图片
    const validDefaultImages = await imagePathChecker.filterValidPaths(defaultImages);
    imageList.push(...validDefaultImages);

    return this.preloadImages(imageList, onProgress, onComplete);
  }

  /**
   * 预加载关键图片（应用启动时）
   * @param {Function} onProgress 进度回调
   * @param {Function} onComplete 完成回调
   */
  async preloadCriticalImages(onProgress, onComplete) {
    // 获取实际存在的关键图片
    const allCriticalImages = [
      '/images/service-default.jpg',
      '/images/empty-state.png',
      '/images/default-avatar.png',
      '/images/hospital-default.jpg'
    ];

    // 过滤有效的图片路径
    const validImages = await imagePathChecker.filterValidPaths(allCriticalImages);
    
    if (validImages.length === 0) {
      console.warn('没有找到有效的关键图片路径');
      onComplete && onComplete({ success: [], failed: [], total: 0, successCount: 0, failedCount: 0 });
      return [];
    }

    return this.preloadImages(validImages, onProgress, onComplete);
  }

  /**
   * 获取缓存信息
   * @returns {Object} 缓存统计信息
   */
  getCacheInfo() {
    return {
      size: this.cache.size,
      keys: Array.from(this.cache.keys()),
      memory: this.cache.size * 0.1 // 估算内存使用（MB）
    };
  }

  /**
   * 清理缓存
   * @param {number} maxAge 最大缓存时间（毫秒）
   */
  clearCache(maxAge = 24 * 60 * 60 * 1000) { // 默认24小时
    const now = Date.now();
    const keysToDelete = [];

    this.cache.forEach((value, key) => {
      if (now - value.timestamp > maxAge) {
        keysToDelete.push(key);
      }
    });

    keysToDelete.forEach(key => this.cache.delete(key));
    
    console.log(`清理了 ${keysToDelete.length} 个过期缓存图片`);
  }

  /**
   * 清空所有缓存
   */
  clearAllCache() {
    this.cache.clear();
    console.log('已清空所有图片缓存');
  }

  /**
   * 检查图片是否已缓存
   * @param {string} src 图片地址
   * @returns {boolean} 是否已缓存
   */
  isCached(src) {
    return this.cache.has(src);
  }

  /**
   * 获取图片缓存信息
   * @param {string} src 图片地址
   * @returns {Object|null} 缓存信息
   */
  getCachedImage(src) {
    return this.cache.get(src) || null;
  }
}

// 创建全局实例
const imagePreloader = new ImagePreloader();

// 导出工具函数
module.exports = {
  imagePreloader,
  
  // 便捷方法
  preloadImage: (src, options) => imagePreloader.preloadImage(src, options),
  preloadImages: (imageList, onProgress, onComplete) => 
    imagePreloader.preloadImages(imageList, onProgress, onComplete),
  preloadHomeImages: (homeData, onProgress, onComplete) => 
    imagePreloader.preloadHomeImages(homeData, onProgress, onComplete),
  preloadServiceImages: (serviceData, onProgress, onComplete) => 
    imagePreloader.preloadServiceImages(serviceData, onProgress, onComplete),
  preloadCriticalImages: (onProgress, onComplete) => 
    imagePreloader.preloadCriticalImages(onProgress, onComplete),
  
  // 缓存管理
  getCacheInfo: () => imagePreloader.getCacheInfo(),
  clearCache: (maxAge) => imagePreloader.clearCache(maxAge),
  clearAllCache: () => imagePreloader.clearAllCache(),
  isCached: (src) => imagePreloader.isCached(src),
  getCachedImage: (src) => imagePreloader.getCachedImage(src)
};
