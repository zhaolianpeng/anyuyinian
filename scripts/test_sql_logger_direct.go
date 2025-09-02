package main

import (
	"fmt"
	"log"
	"time"
)

// 模拟SQLLogger结构
type SQLLogger struct {
	operation string
	table     string
	startTime time.Time
	params    interface{}
}

// NewSQLLogger 创建SQL日志记录器
func NewSQLLogger(operation, table string, params interface{}) *SQLLogger {
	logger := &SQLLogger{
		operation: operation,
		table:     table,
		startTime: time.Now(),
		params:    params,
	}

	// 记录SQL操作开始
	logger.logStart()
	return logger
}

// logStart 记录SQL操作开始
func (l *SQLLogger) logStart() {
	if l.params != nil {
		log.Printf("[SQL] 开始 %s 表: %s, 参数: %+v", l.operation, l.table, l.params)
	} else {
		log.Printf("[SQL] 开始 %s 表: %s", l.operation, l.table)
	}
}

// LogQuery 记录查询操作
func (l *SQLLogger) LogQuery(result interface{}, err error) {
	duration := time.Since(l.startTime)

	if err != nil {
		log.Printf("[SQL] 查询 表: %s 失败 (耗时: %v): %v", l.table, duration, err)
	} else {
		if result != nil {
			log.Printf("[SQL] 查询 表: %s 成功 (耗时: %v), 结果: %+v", l.table, duration, result)
		} else {
			log.Printf("[SQL] 查询 表: %s 成功 (耗时: %v)", l.table, duration)
		}
	}
}

// LogInsert 记录插入操作
func (l *SQLLogger) LogInsert(result interface{}, err error) {
	duration := time.Since(l.startTime)

	if err != nil {
		log.Printf("[SQL] 插入 表: %s 失败 (耗时: %v): %v", l.table, duration, err)
	} else {
		log.Printf("[SQL] 插入 表: %s 成功 (耗时: %v)", l.table, duration)
	}
}

// LogUpdate 记录更新操作
func (l *SQLLogger) LogUpdate(result interface{}, err error) {
	duration := time.Since(l.startTime)

	if err != nil {
		log.Printf("[SQL] 更新 表: %s 失败 (耗时: %v): %v", l.table, duration, err)
	} else {
		if result != nil {
			log.Printf("[SQL] 更新 表: %s 成功 (耗时: %v), 影响行数: %v", l.table, duration, result)
		} else {
			log.Printf("[SQL] 更新 表: %s 成功 (耗时: %v)", l.table, duration)
		}
	}
}

// LogCount 记录计数操作
func (l *SQLLogger) LogCount(count int64, err error) {
	duration := time.Since(l.startTime)

	if err != nil {
		log.Printf("[SQL] 计数 表: %s 失败 (耗时: %v): %v", l.table, duration, err)
	} else {
		log.Printf("[SQL] 计数 表: %s 成功 (耗时: %v), 记录数: %d", l.table, duration, count)
	}
}

func main() {
	fmt.Println("=== SQL日志系统直接测试 ===")

	// 测试查询操作
	fmt.Println("\n1. 测试查询操作:")
	logger1 := NewSQLLogger("查询", "ServiceItems", map[string]interface{}{
		"id":     123,
		"status": 1,
	})
	time.Sleep(2 * time.Millisecond) // 模拟查询耗时
	logger1.LogQuery(map[string]interface{}{
		"id":    123,
		"name":  "测试服务",
		"price": 100.0,
	}, nil)

	// 测试插入操作
	fmt.Println("\n2. 测试插入操作:")
	logger2 := NewSQLLogger("插入", "Orders", map[string]interface{}{
		"orderNo": "ORD20240101001",
		"userId":  "user123",
		"status":  0,
	})
	time.Sleep(3 * time.Millisecond) // 模拟插入耗时
	logger2.LogInsert(nil, nil)

	// 测试更新操作
	fmt.Println("\n3. 测试更新操作:")
	logger3 := NewSQLLogger("更新", "Orders", map[string]interface{}{
		"id":     456,
		"status": 1,
	})
	time.Sleep(1 * time.Millisecond) // 模拟更新耗时
	logger3.LogUpdate(int64(1), nil)

	// 测试计数操作
	fmt.Println("\n4. 测试计数操作:")
	logger4 := NewSQLLogger("计数", "ServiceItems", map[string]interface{}{
		"category": "居家照护",
	})
	time.Sleep(1 * time.Millisecond) // 模拟计数耗时
	logger4.LogCount(25, nil)

	// 测试错误情况
	fmt.Println("\n5. 测试错误情况:")
	logger5 := NewSQLLogger("查询", "ServiceItems", map[string]interface{}{
		"id": 999,
	})
	time.Sleep(1 * time.Millisecond)
	logger5.LogQuery(nil, fmt.Errorf("记录不存在"))

	fmt.Println("\n=== 测试完成 ===")
	fmt.Println("如果看到上述SQL日志输出，说明SQL日志系统工作正常")
}
