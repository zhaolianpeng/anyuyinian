# 推广员页面Tab切换错误修复

## 问题描述

在推广员页面（`pages/promoter/home.js`）中，当用户切换tab时出现以下错误：

```
TypeError: Cannot read property 'length' of null
    at li.onTabChange (home.js:80)
```

## 错误原因

在`onTabChange`方法中，代码尝试访问`this.data.commissionList.length`和`this.data.cashoutList.length`，但在某些情况下这些数组可能为null，导致无法读取length属性。

## 修复方案

### 1. 修复onTabChange方法

**修复前：**
```javascript
onTabChange(e) {
  const activeTab = parseInt(e.currentTarget.dataset.index)
  console.log('切换标签页:', activeTab)
  this.setData({ activeTab })
  
  if (activeTab === 1 && this.data.commissionList.length === 0) {
    this.loadCommissionList(true)
  } else if (activeTab === 2 && this.data.cashoutList.length === 0) {
    this.loadCashoutList(true)
  }
}
```

**修复后：**
```javascript
onTabChange(e) {
  const activeTab = parseInt(e.currentTarget.dataset.index)
  console.log('切换标签页:', activeTab)
  this.setData({ activeTab })
  
  // 添加安全检查，确保数组存在且不为null
  if (activeTab === 1 && (!this.data.commissionList || this.data.commissionList.length === 0)) {
    this.loadCommissionList(true)
  } else if (activeTab === 2 && (!this.data.cashoutList || this.data.cashoutList.length === 0)) {
    this.loadCashoutList(true)
  }
}
```

### 2. 修复loadCommissionList方法

**修复前：**
```javascript
const commissionList = refresh ? list : [...this.data.commissionList, ...list]
```

**修复后：**
```javascript
// 添加安全检查，确保this.data.commissionList存在
const currentList = this.data.commissionList || []
const commissionList = refresh ? list : [...currentList, ...list]
```

### 3. 修复loadCashoutList方法

**修复前：**
```javascript
const cashoutList = refresh ? list : [...this.data.cashoutList, ...list]
```

**修复后：**
```javascript
// 添加安全检查，确保this.data.cashoutList存在
const currentList = this.data.cashoutList || []
const cashoutList = refresh ? list : [...currentList, ...list]
```

## 修复原理

### 空值检查
使用逻辑或操作符（`||`）来提供默认值：
- `!this.data.commissionList` 检查数组是否为null或undefined
- `this.data.commissionList || []` 如果数组为null，则使用空数组作为默认值

### 防御性编程
在访问对象属性之前，先检查对象是否存在，避免运行时错误。

## 测试验证

### 测试用例
1. **正常情况**：数组已正确初始化
2. **异常情况1**：commissionList为null
3. **异常情况2**：cashoutList为null
4. **异常情况3**：两个数组都为null

### 测试脚本
运行 `tests/promoter/test_tab_switch.js` 来验证修复效果。

## 预防措施

### 1. 数据初始化
确保在Page的data中正确初始化数组：
```javascript
data: {
  commissionList: [],
  cashoutList: [],
  // ... 其他数据
}
```

### 2. 错误处理
在API调用失败时，不要将数组设置为null，而是保持为空数组或添加错误状态。

### 3. 代码审查
在代码审查时，特别注意对数组和对象属性的访问，确保有适当的空值检查。

## 相关文件

- `pages/promoter/home.js` - 主要修复文件
- `tests/promoter/test_tab_switch.js` - 测试脚本
- `docs/fixes/PROMOTER_TAB_SWITCH_FIX.md` - 本文档

## 总结

通过添加适当的空值检查，成功修复了推广员页面tab切换时的JavaScript错误。这个修复提高了代码的健壮性，避免了因数据状态异常导致的运行时错误。 