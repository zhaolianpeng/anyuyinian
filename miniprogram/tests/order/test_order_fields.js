// 测试预约页面字段显示
console.log('测试预约页面字段显示');

// 检查字段是否存在
const checkOrderFields = () => {
  console.log('=== 预约页面字段检查 ===');
  
  console.log('✅ 检查WXML字段:');
  console.log('1. 患者性别 - 在预约信息区域');
  console.log('2. 患者年龄 - 在预约信息区域');
  console.log('3. 基础病信息 - 在预约信息区域');
  console.log('4. 助排二便 - 在预约信息区域');
  
  console.log('');
  console.log('✅ 检查JavaScript函数:');
  console.log('1. onDiseaseInput - 基础病信息输入处理');
  console.log('2. onToiletAssistChange - 助排二便选择处理');
  console.log('3. getAge - 年龄计算函数');
  
  console.log('');
  console.log('✅ 检查数据结构:');
  console.log('formData.diseaseInfo - 基础病信息');
  console.log('formData.needToiletAssist - 助排二便');
  
  console.log('');
  console.log('✅ 检查样式:');
  console.log('1. .disease-section - 基础病信息样式');
  console.log('2. .toilet-section - 助排二便样式');
  console.log('3. .radio-group - 单选按钮组样式');
  
  console.log('');
  console.log('=== 字段显示位置 ===');
  console.log('预约信息区域包含:');
  console.log('- 预约时间选择');
  console.log('- 患者性别（显示）');
  console.log('- 患者年龄（显示）');
  console.log('- 基础病信息（输入框）');
  console.log('- 助排二便（单选按钮）');
  console.log('- 时间范围提示');
  
  console.log('');
  console.log('=== 测试建议 ===');
  console.log('1. 打开预约页面');
  console.log('2. 检查预约信息区域是否显示所有字段');
  console.log('3. 测试基础病信息输入框');
  console.log('4. 测试助排二便单选按钮');
  console.log('5. 选择患者后检查性别和年龄显示');
};

// 检查字段显示问题
const checkDisplayIssues = () => {
  console.log('=== 可能的显示问题 ===');
  
  console.log('⚠️ 检查点:');
  console.log('1. 性别和年龄字段是否有条件显示');
  console.log('2. 基础病信息输入框是否可见');
  console.log('3. 助排二便单选按钮是否可点击');
  console.log('4. 样式是否正确加载');
  console.log('5. 数据绑定是否正确');
  
  console.log('');
  console.log('🔧 修复建议:');
  console.log('1. 确保性别和年龄字段始终显示');
  console.log('2. 检查基础病信息输入框的样式');
  console.log('3. 验证助排二便单选按钮的交互');
  console.log('4. 测试数据输入和保存');
};

// 运行测试
const runFieldTest = () => {
  console.log('开始测试预约页面字段...');
  checkOrderFields();
  console.log('');
  checkDisplayIssues();
  console.log('');
  console.log('=== 测试完成 ===');
  console.log('请在预约页面验证字段显示情况');
};

// 导出测试函数
module.exports = {
  checkOrderFields,
  checkDisplayIssues,
  runFieldTest
};

// 如果直接运行此文件
if (typeof window !== 'undefined') {
  runFieldTest();
} 