package dao

import (
	"encoding/json"
	"log"
	"time"
)

// SQLLogger SQL操作日志记录器
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
		paramsJSON, _ := json.Marshal(l.params)
		log.Printf("[SQL] 开始 %s 表: %s, 参数: %s", l.operation, l.table, string(paramsJSON))
	} else {
		log.Printf("[SQL] 开始 %s 表: %s", l.operation, l.table)
	}
}

// LogResult 记录SQL操作结果
func (l *SQLLogger) LogResult(result interface{}, err error) {
	duration := time.Since(l.startTime)

	if err != nil {
		log.Printf("[SQL] %s 表: %s 失败 (耗时: %v): %v", l.operation, l.table, duration, err)
	} else {
		if result != nil {
			// 根据结果类型记录不同的信息
			switch v := result.(type) {
			case int64:
				log.Printf("[SQL] %s 表: %s 成功 (耗时: %v), 影响行数: %d", l.operation, l.table, duration, v)
			case []interface{}:
				log.Printf("[SQL] %s 表: %s 成功 (耗时: %v), 返回记录数: %d", l.operation, l.table, duration, len(v))
			default:
				resultJSON, _ := json.Marshal(result)
				log.Printf("[SQL] %s 表: %s 成功 (耗时: %v), 结果: %s", l.operation, l.table, duration, string(resultJSON))
			}
		} else {
			log.Printf("[SQL] %s 表: %s 成功 (耗时: %v)", l.operation, l.table, duration)
		}
	}
}

// LogQuery 记录查询操作
func (l *SQLLogger) LogQuery(result interface{}, err error) {
	duration := time.Since(l.startTime)

	if err != nil {
		log.Printf("[SQL] 查询 表: %s 失败 (耗时: %v): %v", l.table, duration, err)
	} else {
		if result != nil {
			// 根据结果类型记录不同的信息
			switch v := result.(type) {
			case int64:
				log.Printf("[SQL] 查询 表: %s 成功 (耗时: %v), 记录数: %d", l.table, duration, v)
			case []interface{}:
				log.Printf("[SQL] 查询 表: %s 成功 (耗时: %v), 返回记录数: %d", l.table, duration, len(v))
			default:
				resultJSON, _ := json.Marshal(result)
				log.Printf("[SQL] 查询 表: %s 成功 (耗时: %v), 结果: %s", l.table, duration, string(resultJSON))
			}
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
			if affectedRows, ok := result.(int64); ok {
				log.Printf("[SQL] 更新 表: %s 成功 (耗时: %v), 影响行数: %d", l.table, duration, affectedRows)
			} else {
				log.Printf("[SQL] 更新 表: %s 成功 (耗时: %v)", l.table, duration)
			}
		} else {
			log.Printf("[SQL] 更新 表: %s 成功 (耗时: %v)", l.table, duration)
		}
	}
}

// LogDelete 记录删除操作
func (l *SQLLogger) LogDelete(result interface{}, err error) {
	duration := time.Since(l.startTime)

	if err != nil {
		log.Printf("[SQL] 删除 表: %s 失败 (耗时: %v): %v", l.table, duration, err)
	} else {
		if result != nil {
			if affectedRows, ok := result.(int64); ok {
				log.Printf("[SQL] 删除 表: %s 成功 (耗时: %v), 影响行数: %d", l.table, duration, affectedRows)
			} else {
				log.Printf("[SQL] 删除 表: %s 成功 (耗时: %v)", l.table, duration)
			}
		} else {
			log.Printf("[SQL] 删除 表: %s 成功 (耗时: %v)", l.table, duration)
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
