package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// 数据库配置
const (
	DBHost     = "localhost"
	DBPort     = 3306
	DBUser     = "root"
	DBPassword = "your_password"
	DBName     = "your_database"
)

// 生成MongoDB风格的UserId
func generateUserID() string {
	// 生成24位字符串
	const charset = "0123456789abcdef"
	result := make([]byte, 24)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// 检查数据库连接
func checkDBConnection(db *sql.DB) error {
	return db.Ping()
}

// 检查当前用户状态
func checkCurrentUsers(db *sql.DB) error {
	fmt.Println("=== 检查当前用户状态 ===")
	
	rows, err := db.Query("SELECT id, openId, userId FROM Users ORDER BY id")
	if err != nil {
		return fmt.Errorf("查询用户失败: %v", err)
	}
	defer rows.Close()

	fmt.Printf("%-5s %-20s %-30s %s\n", "ID", "OpenId", "UserId", "Status")
	fmt.Println("--------------------------------------------------------------------------------")

	for rows.Next() {
		var id int
		var openId, userId sql.NullString
		err := rows.Scan(&id, &openId, &userId)
		if err != nil {
			return fmt.Errorf("扫描用户数据失败: %v", err)
		}

		status := "正常"
		if !userId.Valid || userId.String == "" {
			status = "需要修复"
		}

		fmt.Printf("%-5d %-20s %-30s %s\n", 
			id, 
			openId.String, 
			userId.String, 
			status)
	}

	return nil
}

// 修复用户UserId
func fixUserIds(db *sql.DB) error {
	fmt.Println("\n=== 开始修复UserId ===")
	
	// 查询需要修复的用户
	rows, err := db.Query("SELECT id FROM Users WHERE userId IS NULL OR userId = ''")
	if err != nil {
		return fmt.Errorf("查询需要修复的用户失败: %v", err)
	}
	defer rows.Close()

	var userIds []int
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return fmt.Errorf("扫描用户ID失败: %v", err)
		}
		userIds = append(userIds, id)
	}

	fmt.Printf("找到 %d 个需要修复的用户\n", len(userIds))

	// 为每个用户生成新的UserId
	for _, id := range userIds {
		newUserId := generateUserID()
		
		_, err := db.Exec("UPDATE Users SET userId = ? WHERE id = ?", newUserId, id)
		if err != nil {
			return fmt.Errorf("更新用户 %d 的UserId失败: %v", id, err)
		}
		
		fmt.Printf("用户 %d 的UserId已更新为: %s\n", id, newUserId)
	}

	return nil
}

// 验证修复结果
func validateFix(db *sql.DB) error {
	fmt.Println("\n=== 验证修复结果 ===")
	
	rows, err := db.Query(`
		SELECT id, openId, userId, LENGTH(userId) as userId_length
		FROM Users 
		ORDER BY id
	`)
	if err != nil {
		return fmt.Errorf("查询用户验证数据失败: %v", err)
	}
	defer rows.Close()

	fmt.Printf("%-5s %-20s %-30s %-10s %s\n", "ID", "OpenId", "UserId", "Length", "Status")
	fmt.Println("--------------------------------------------------------------------------------")

	for rows.Next() {
		var id int
		var openId, userId sql.NullString
		var length int
		err := rows.Scan(&id, &openId, &userId, &length)
		if err != nil {
			return fmt.Errorf("扫描验证数据失败: %v", err)
		}

		status := "❌ 错误"
		if length == 24 {
			status = "✅ 正确"
		}

		fmt.Printf("%-5d %-20s %-30s %-10d %s\n", 
			id, 
			openId.String, 
			userId.String, 
			length, 
			status)
	}

	return nil
}

// 检查重复UserId
func checkDuplicateUserIds(db *sql.DB) error {
	fmt.Println("\n=== 检查重复UserId ===")
	
	rows, err := db.Query(`
		SELECT userId, COUNT(*) as count 
		FROM Users 
		GROUP BY userId 
		HAVING COUNT(*) > 1
	`)
	if err != nil {
		return fmt.Errorf("查询重复UserId失败: %v", err)
	}
	defer rows.Close()

	hasDuplicates := false
	for rows.Next() {
		var userId string
		var count int
		err := rows.Scan(&userId, &count)
		if err != nil {
			return fmt.Errorf("扫描重复数据失败: %v", err)
		}
		
		fmt.Printf("发现重复UserId: %s (出现 %d 次)\n", userId, count)
		hasDuplicates = true
	}

	if !hasDuplicates {
		fmt.Println("✅ 没有发现重复的UserId")
	}

	return nil
}

// 最终验证
func finalValidation(db *sql.DB) error {
	fmt.Println("\n=== 最终验证 ===")
	
	var totalUsers, usersWithUserId, validUserIds int
	
	err := db.QueryRow("SELECT COUNT(*) FROM Users").Scan(&totalUsers)
	if err != nil {
		return fmt.Errorf("查询总用户数失败: %v", err)
	}

	err = db.QueryRow("SELECT COUNT(userId) FROM Users").Scan(&usersWithUserId)
	if err != nil {
		return fmt.Errorf("查询有UserId的用户数失败: %v", err)
	}

	err = db.QueryRow(`
		SELECT COUNT(*) 
		FROM Users 
		WHERE userId IS NOT NULL AND userId != '' AND LENGTH(userId) = 24
	`).Scan(&validUserIds)
	if err != nil {
		return fmt.Errorf("查询有效UserId的用户数失败: %v", err)
	}

	fmt.Printf("总用户数: %d\n", totalUsers)
	fmt.Printf("有UserId的用户数: %d\n", usersWithUserId)
	fmt.Printf("有效UserId的用户数: %d\n", validUserIds)

	if totalUsers == validUserIds {
		fmt.Println("✅ 所有用户都有有效的UserId")
	} else {
		fmt.Println("❌ 还有用户没有有效的UserId")
	}

	return nil
}

func main() {
	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())

	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		DBUser, DBPassword, DBHost, DBPort, DBName)
	
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	// 检查数据库连接
	if err := checkDBConnection(db); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	fmt.Println("✅ 数据库连接成功")

	// 检查当前用户状态
	if err := checkCurrentUsers(db); err != nil {
		log.Fatalf("检查用户状态失败: %v", err)
	}

	// 修复用户UserId
	if err := fixUserIds(db); err != nil {
		log.Fatalf("修复UserId失败: %v", err)
	}

	// 验证修复结果
	if err := validateFix(db); err != nil {
		log.Fatalf("验证修复结果失败: %v", err)
	}

	// 检查重复UserId
	if err := checkDuplicateUserIds(db); err != nil {
		log.Fatalf("检查重复UserId失败: %v", err)
	}

	// 最终验证
	if err := finalValidation(db); err != nil {
		log.Fatalf("最终验证失败: %v", err)
	}

	fmt.Println("\n🎉 UserId修复完成！")
} 