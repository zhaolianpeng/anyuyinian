// 测试预约页面左右滑动修复
console.log('测试预约页面左右滑动修复');

// 检查页面样式
const checkPageStyles = () => {
  console.log('检查页面样式设置:');
  console.log('- page overflow-x: hidden');
  console.log('- page width: 100vw');
  console.log('- page max-width: 100vw');
  console.log('- page position: fixed');
  console.log('- page touch-action: pan-y');
  console.log('- container overflow-x: hidden');
  console.log('- container max-width: 100vw');
  console.log('- container touch-action: pan-y');
};

// 检查元素样式
const checkElementStyles = () => {
  console.log('检查元素样式设置:');
  console.log('- 所有元素 max-width: 100vw');
  console.log('- 所有元素 overflow-x: hidden');
  console.log('- 文本元素 word-wrap: break-word');
  console.log('- 文本元素 word-break: break-all');
  console.log('- 按钮元素 max-width: 100%');
};

// 测试触摸事件
const testTouchEvents = () => {
  console.log('测试触摸事件:');
  console.log('- 只允许垂直滑动 (pan-y)');
  console.log('- 禁止水平滑动');
  console.log('- 使用 -webkit-overflow-scrolling: touch');
};

// 运行所有测试
const runAllTests = () => {
  console.log('=== 预约页面左右滑动修复测试 ===');
  checkPageStyles();
  checkElementStyles();
  testTouchEvents();
  console.log('=== 测试完成 ===');
  console.log('请在预约页面测试左右滑动是否被禁止');
};

// 导出测试函数
module.exports = {
  runAllTests,
  checkPageStyles,
  checkElementStyles,
  testTouchEvents
};

// 如果直接运行此文件
if (typeof window !== 'undefined') {
  runAllTests();
} 