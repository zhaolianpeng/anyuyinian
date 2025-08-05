package main

import (
	"fmt"
	"wxcloudrun-golang/utils"
)

func main() {
	fmt.Println("=== 推广码生成测试 ===")

	// 测试生成多个推广码
	fmt.Println("生成10个推广码:")
	for i := 0; i < 10; i++ {
		code := utils.GeneratePromoterCode()
		fmt.Printf("推广码 %d: %s (格式验证: %t)\n", i+1, code, utils.ValidatePromoterCode(code))
	}

	// 测试格式化
	fmt.Println("\n测试格式化:")
	testCodes := []string{"ABC123", "XYZ789", "123ABC"}
	for _, code := range testCodes {
		formatted := utils.FormatPromoterCode(code)
		fmt.Printf("原始: %s -> 格式化: %s\n", code, formatted)
	}

	// 测试验证
	fmt.Println("\n测试验证:")
	testCases := []string{
		"ABC123",    // 有效
		"XYZ789",    // 有效
		"123ABC",    // 有效
		"ABC12",     // 无效：长度不足
		"ABC1234",   // 无效：长度过长
		"ABC12a",    // 无效：包含小写字母
		"ABC12@",    // 无效：包含特殊字符
		"",          // 无效：空字符串
	}

	for _, testCase := range testCases {
		isValid := utils.ValidatePromoterCode(testCase)
		fmt.Printf("'%s' -> 有效: %t\n", testCase, isValid)
	}

	// 测试唯一性检查
	fmt.Println("\n测试唯一性生成:")
	existingCodes := map[string]bool{
		"ABC123": true,
		"XYZ789": true,
		"123ABC": true,
	}

	checkExists := func(code string) bool {
		return existingCodes[code]
	}

	for i := 0; i < 5; i++ {
		uniqueCode := utils.GenerateUniquePromoterCode(checkExists)
		fmt.Printf("唯一推广码 %d: %s\n", i+1, uniqueCode)
	}

	fmt.Println("\n=== 测试完成 ===")
} 