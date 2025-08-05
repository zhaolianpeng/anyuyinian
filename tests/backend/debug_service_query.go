package main

import (
	"encoding/json"
	"fmt"
	"log"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/dao"
)

func main() {
	fmt.Println("=== 调试Service查询问题 ===")

	// 初始化数据库连接
	err := db.Init()
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 测试GetServices查询
	fmt.Println("1. 测试GetServices查询...")
	services, err := dao.HomeImp.GetServices()
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}

	fmt.Printf("查询到 %d 个服务\n", len(services))

	// 打印每个服务的详细信息
	for i, service := range services {
		fmt.Printf("\n服务 %d:\n", i+1)
		fmt.Printf("  ID: %d\n", service.Id)
		fmt.Printf("  ServiceItemId: %d\n", service.ServiceItemId)
		fmt.Printf("  Name: %s\n", service.Name)
		fmt.Printf("  Description: %s\n", service.Description)

		// 检查serviceItemId是否正确
		if service.ServiceItemId == 0 {
			fmt.Printf("  ❌ ServiceItemId为0，GORM没有正确映射serviceitemid字段\n")
		} else if service.ServiceItemId == service.Id {
			fmt.Printf("  ✅ ServiceItemId与Id相同: %d\n", service.ServiceItemId)
		} else {
			fmt.Printf("  ✅ ServiceItemId与Id不同: ServiceItemId=%d, Id=%d\n", service.ServiceItemId, service.Id)
		}
	}

	// 测试JSON序列化
	fmt.Println("\n2. 测试JSON序列化...")
	for i, service := range services {
		jsonData, err := json.Marshal(service)
		if err != nil {
			log.Printf("JSON序列化失败: %v", err)
			continue
		}
		fmt.Printf("服务 %d JSON: %s\n", i+1, string(jsonData))
	}

	// 测试API响应格式
	fmt.Println("\n3. 测试API响应格式...")
	for i, service := range services {
		apiResponse := map[string]interface{}{
			"id":          service.Id,
			"serviceId":   service.ServiceItemId,
			"name":        service.Name,
			"description": service.Description,
		}
		jsonData, _ := json.Marshal(apiResponse)
		fmt.Printf("API响应 %d: %s\n", i+1, string(jsonData))
	}

	fmt.Println("\n=== 调试完成 ===")
}
