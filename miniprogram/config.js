// 环境配置
const ENV = {
  DEV: 'dev',
  PROD: 'prod'
}

// 当前环境 - 请根据实际情况修改
const CURRENT_ENV = ENV.PROD

// 服务端直连配置
// 当前小程序统一通过 wx.request 直连用户自己的服务端。
// 以下 DOMAINS 和 container 配置仅为历史兼容保留，运行时主链不再依赖它们。

// 历史域名配置（仅兼容保留）
const DOMAINS = {
  [ENV.DEV]: 'https://golang-lfwy-176496-6-1353115175.sh.run.tcloudbase.com',  // 开发环境域名
  [ENV.PROD]: 'https://golang-lfwy-176496-6-1353115175.sh.run.tcloudbase.com'  // 生产环境域名
}

// 历史配置（兼容保留）
const CONTAINER_CONFIG = {
  // 历史环境ID
  envId: 'prod-5g94mx7a3d07e78c',
  // 历史服务名称
  serviceName: 'golang-lfwy',
  // 小程序AppID
  appId: 'wx101090677bd5219e'
}

// SSE配置
const SSE_CONFIG = {
  // 云托管环境ID
  env: 'prod-5g94mx7a3d07e78c',
  // SSE服务名 - 使用与HTTP相同的服务名
  service: 'golang-lfwy',
  // SSE路径
  path: '/sse',
  // 重连配置
  reconnect: {
    enabled: true,
    maxAttempts: 5,
    interval: 3000
  }
}

// 静态资源配置
const COS_CONFIG = {
  // 静态资源基础地址
  bucketDomain: 'https://api.succ.online/anyuyinian',
  // 图片路径前缀
  imagePrefix: '/static/',
  // 默认图片
  defaultImages: {
    service: '/static/default-service.png',
    hospital: '/static/default-hospital.png',
    user: '/static/default-user.png'
  }
}

module.exports = {
  // 当前环境
  env: CURRENT_ENV,
  
  // 历史基础URL（兼容保留）
  baseURL: DOMAINS[CURRENT_ENV],
  
  // 历史配置（兼容保留）
  container: CONTAINER_CONFIG,
  
  // SSE配置
  sse: SSE_CONFIG,
  
  // 静态资源配置
  cos: COS_CONFIG,
  
  // API接口路径
  api: {
    // 用户认证
    wxLogin: '/api/wx/login',
    
    // 首页相关
    homeInit: '/api/home/init',
    
    // 用户管理
    userInfo: '/api/user/info',
    bindPhone: '/api/user/bind_phone',
    userAddress: '/api/user/address',
    userPatient: '/api/user/patient',
    
    // 服务管理
    serviceList: '/api/service/list',
    serviceDetail: '/api/service/detail',
    serviceFormConfig: '/api/service/form_config',
    
    // 订单管理
    orderSubmit: '/api/order/submit',
    orderPay: '/api/order/pay',
    orderCancel: '/api/order/cancel',
    orderRefund: '/api/order/refund',
    orderList: '/api/order/list',
    orderDetail: '/api/order/detail',
    
    // 推荐系统
    referralQrcode: '/api/referral/qrcode',
    referralList: '/api/referral/list',
    commissionList: '/api/commission/list',
    cashoutSubmit: '/api/cashout/submit',
    
    // 客服医院
    kefuMessage: '/api/kefu/message',
    kefuFaq: '/api/kefu/faq',
    hospitalList: '/api/hospital/list',
    hospitalDetail: '/api/hospital/detail',
    
    // 文件管理
    uploadFile: '/api/upload/file',
    fileList: '/api/upload/list',
    
    // 系统配置
    config: '/api/config'
  },
  
  // 错误码定义
  errorCode: {
    SUCCESS: 0,
    PARAM_ERROR: -1,
    AUTH_ERROR: -2,
    SERVER_ERROR: -3,
    NETWORK_ERROR: -4
  },
  
  // 订单状态
  orderStatus: {
    PENDING_PAY: 'pending_pay',
    PAID: 'paid',
    CANCELLED: 'cancelled',
    REFUNDED: 'refunded'
  },
  
  // 订单状态文本
  orderStatusText: {
    pending_pay: '待支付',
    paid: '已支付',
    cancelled: '已取消',
    refunded: '已退款'
  },
  
  // 开发环境配置
  dev: {
    // 开发环境可以使用模拟数据
    useMockData: false,
    // 模拟数据延迟
    mockDelay: 1000
  }
}