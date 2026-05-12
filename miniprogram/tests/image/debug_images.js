// 图片调试脚本

const { processImageUrl, processHomeDataImages } = require('./utils/image')

console.log('=== 图片调试 ===')

// 测试数据
const testData = {
  "code": 0,
  "data": {
    "banners": [
      {
        "id": 1,
        "imageUrl": "https://7072-prod-5g94mx7a3d07e78c-1353115175.cos.ap-shanghai.myqcloud.com/static/fuwu_2.jpeg",
        "linkUrl": "/static/fuwu_2.jpeg",
        "sort": 1,
        "title": "专业护理服务"
      },
      {
        "id": 2,
        "imageUrl": "https://7072-prod-5g94mx7a3d07e78c-1353115175.cos.ap-shanghai.myqcloud.com/static/fuwu_1.jpeg",
        "linkUrl": "/static/fuwu_1.jpeg",
        "sort": 2,
        "title": "医院陪诊套餐"
      }
    ],
    "navigations": [
      {
        "icon": "/static/fuwuyuyue_logo.png",
        "id": 1,
        "linkUrl": "/pages/service/list",
        "name": "服务预约",
        "sort": 1
      },
      {
        "icon": "/static/wodedingdan_logo.png",
        "id": 2,
        "linkUrl": "/pages/order/list",
        "name": "我的订单",
        "sort": 2
      }
    ],
    "services": [
      {
        "description": "专业护工上门提供护理服务",
        "icon": "/static/service_pic_1.png",
        "id": 1,
        "imageUrl": "/static/service_pic_1.png",
        "linkUrl": "/pages/service/list",
        "name": "上门护理服务",
        "sort": 1
      },
      {
        "description": "专业陪诊师全程陪同就医，提供挂号、排队、取药等服务",
        "icon": "/static/service_pic_2.png",
        "id": 2,
        "imageUrl": "/static/service_pic_2.png",
        "linkUrl": "/pages/service/list",
        "name": "专业陪诊服务",
        "sort": 2
      }
    ],
    "hospitals": [
      {
        "address": "深圳市福田区红荔路2004号",
        "description": "深圳市妇幼保健院是一所集医疗、教学、科研、预防、保健为一体的三级甲等妇幼保健院",
        "id": 5,
        "latitude": 23.098765,
        "level": "三级甲等",
        "logo": "/static/logo.jpeg",
        "longitude": 114.56789,
        "name": "深圳市妇幼保健院",
        "phone": "0755-83000111",
        "sort": 5,
        "type": "专科医院"
      }
    ]
  }
}

// 调试单个图片URL
function debugSingleImage() {
  console.log('\n调试单个图片URL...')
  
  const testUrls = [
    'https://7072-prod-5g94mx7a3d07e78c-1353115175.cos.ap-shanghai.myqcloud.com/static/fuwu_2.jpeg',
    '/static/fuwu_2.jpeg',
    '/static/fuwuyuyue_logo.png',
    '/static/service_pic_1.png',
    '/static/logo.jpeg'
  ]
  
  testUrls.forEach((url, index) => {
    const processed = processImageUrl(url)
    console.log(`  ${index + 1}. 原始URL: ${url}`)
    console.log(`     处理后: ${processed}`)
    console.log(`     是否HTTPS: ${processed.startsWith('https://')}`)
    console.log('')
  })
}

// 调试首页数据处理
function debugHomeData() {
  console.log('\n调试首页数据处理...')
  
  const processed = processHomeDataImages(testData)
  
  console.log('轮播图:')
  processed.data.banners.forEach((banner, index) => {
    console.log(`  Banner ${index + 1}:`)
    console.log(`    原始: ${testData.data.banners[index].imageUrl}`)
    console.log(`    处理后: ${banner.imageUrl}`)
    console.log('')
  })
  
  console.log('导航菜单:')
  processed.data.navigations.forEach((nav, index) => {
    console.log(`  Nav ${index + 1}:`)
    console.log(`    原始: ${testData.data.navigations[index].icon}`)
    console.log(`    处理后: ${nav.icon}`)
    console.log('')
  })
  
  console.log('服务列表:')
  processed.data.services.forEach((service, index) => {
    console.log(`  Service ${index + 1}:`)
    console.log(`    原始imageUrl: ${testData.data.services[index].imageUrl}`)
    console.log(`    处理后imageUrl: ${service.imageUrl}`)
    console.log(`    原始icon: ${testData.data.services[index].icon}`)
    console.log(`    处理后icon: ${service.icon}`)
    console.log('')
  })
  
  console.log('医院列表:')
  processed.data.hospitals.forEach((hospital, index) => {
    console.log(`  Hospital ${index + 1}:`)
    console.log(`    原始: ${testData.data.hospitals[index].logo}`)
    console.log(`    处理后: ${hospital.logo}`)
    console.log('')
  })
}

// 检查配置
function checkConfig() {
  console.log('\n检查配置...')
  
  try {
    const { cos } = require('./config')
    console.log('COS配置:')
    console.log(`  域名: ${cos.bucketDomain}`)
    console.log(`  前缀: ${cos.imagePrefix}`)
    console.log(`  环境ID: ${cos.cloudEnvId}`)
    console.log(`  存储桶: ${cos.cloudBucketName}`)
  } catch (error) {
    console.error('配置加载失败:', error.message)
  }
}

// 测试图片加载
function testImageLoad() {
  console.log('\n测试图片加载...')
  
  const testImages = [
    'https://7072-prod-5g94mx7a3d07e78c-1353115175.cos.ap-shanghai.myqcloud.com/static/fuwu_2.jpeg',
    'https://7072-prod-5g94mx7a3d07e78c-1353115175.cos.ap-shanghai.myqcloud.com/static/fuwu_1.jpeg',
    'https://7072-prod-5g94mx7a3d07e78c-1353115175.cos.ap-shanghai.myqcloud.com/static/fuwuyuyue_logo.png'
  ]
  
  testImages.forEach((url, index) => {
    console.log(`  图片 ${index + 1}: ${url}`)
    
    // 模拟图片加载测试
    if (url.startsWith('https://')) {
      console.log(`    ✅ HTTPS URL格式正确`)
    } else {
      console.log(`    ❌ URL格式错误`)
    }
    
    // 检查域名
    if (url.includes('cos.ap-shanghai.myqcloud.com')) {
      console.log(`    ✅ COS域名正确`)
    } else {
      console.log(`    ❌ COS域名错误`)
    }
    
    console.log('')
  })
}

// 执行调试
function runDebug() {
  console.log('开始图片调试...\n')
  
  checkConfig()
  debugSingleImage()
  debugHomeData()
  testImageLoad()
  
  console.log('=== 调试完成 ===')
  console.log('')
  console.log('可能的问题:')
  console.log('1. 检查COS域名是否正确')
  console.log('2. 检查图片路径是否存在')
  console.log('3. 检查网络连接是否正常')
  console.log('4. 检查小程序域名配置')
  console.log('')
  console.log('建议操作:')
  console.log('1. 在浏览器中直接访问图片URL测试')
  console.log('2. 检查小程序后台域名配置')
  console.log('3. 确认COS桶权限设置')
}

// 运行调试
runDebug() 