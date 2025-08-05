package service

import (
	"fmt"
	"log"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"
)

// MigrationService 数据库迁移服务
type MigrationService struct{}

// NewMigrationService 创建迁移服务实例
func NewMigrationService() *MigrationService {
	return &MigrationService{}
}

// MigrateExistingUsers 为现有用户生成UserId
func (m *MigrationService) MigrateExistingUsers() error {
	log.Println("开始为现有用户生成UserId...")

	// 查询所有没有userId的用户
	var users []*model.UserModel
	cli := db.Get()
	err := cli.Table("Users").Where("userId IS NULL OR userId = ''").Find(&users).Error
	if err != nil {
		return fmt.Errorf("查询用户失败: %v", err)
	}

	log.Printf("找到 %d 个需要生成UserId的用户", len(users))

	// 为每个用户生成UserId
	for _, user := range users {
		// 生成新的UserId
		newUserId := utils.GenerateUserID()

		// 更新用户记录
		err := cli.Table("Users").Where("id = ?", user.Id).Update("userId", newUserId).Error
		if err != nil {
			log.Printf("更新用户 %d 的UserId失败: %v", user.Id, err)
			continue
		}

		log.Printf("用户 %d 的UserId已更新为: %s", user.Id, newUserId)
	}

	log.Println("UserId迁移完成")
	return nil
}

// MigrateAllTablesUserId 迁移所有表的userId字段
func (m *MigrationService) MigrateAllTablesUserId() error {
	log.Println("开始迁移所有表的userId字段...")

	cli := db.Get()

	// 1. 迁移Orders表
	log.Println("迁移Orders表...")
	var orders []*model.OrderModel
	err := cli.Table("Orders").Find(&orders).Error
	if err != nil {
		return fmt.Errorf("查询Orders失败: %v", err)
	}

	for _, order := range orders {
		// 检查userId是否为数字格式
		if order.UserId == "" || len(order.UserId) < 24 {
			// 生成新的userId
			newUserId := utils.GenerateUserID()
			err := cli.Table("Orders").Where("id = ?", order.Id).Update("userId", newUserId).Error
			if err != nil {
				log.Printf("更新Order %d 的userId失败: %v", order.Id, err)
				continue
			}
			log.Printf("Order %d 的userId已更新为: %s", order.Id, newUserId)
		}
	}

	// 2. 迁移KefuMessages表
	log.Println("迁移KefuMessages表...")
	var messages []*model.KefuMessageModel
	err = cli.Table("KefuMessages").Find(&messages).Error
	if err != nil {
		return fmt.Errorf("查询KefuMessages失败: %v", err)
	}

	for _, message := range messages {
		// 检查userId是否为数字格式
		if message.UserId == "" || len(message.UserId) < 24 {
			newUserId := utils.GenerateUserID()
			err := cli.Table("KefuMessages").Where("id = ?", message.Id).Update("userId", newUserId).Error
			if err != nil {
				log.Printf("更新KefuMessage %d 的userId失败: %v", message.Id, err)
				continue
			}
			log.Printf("KefuMessage %d 的userId已更新为: %s", message.Id, newUserId)
		}
	}

	// 3. 迁移Referrals表
	log.Println("迁移Referrals表...")
	var referrals []*model.ReferralModel
	err = cli.Table("Referrals").Find(&referrals).Error
	if err != nil {
		return fmt.Errorf("查询Referrals失败: %v", err)
	}

	for _, referral := range referrals {
		if referral.UserId == "" || len(referral.UserId) < 24 {
			newUserId := utils.GenerateUserID()
			err := cli.Table("Referrals").Where("id = ?", referral.Id).Update("userId", newUserId).Error
			if err != nil {
				log.Printf("更新Referral %d 的userId失败: %v", referral.Id, err)
				continue
			}
			log.Printf("Referral %d 的userId已更新为: %s", referral.Id, newUserId)
		}
	}

	// 4. 迁移Commissions表
	log.Println("迁移Commissions表...")
	var commissions []*model.CommissionModel
	err = cli.Table("Commissions").Find(&commissions).Error
	if err != nil {
		return fmt.Errorf("查询Commissions失败: %v", err)
	}

	for _, commission := range commissions {
		if commission.UserId == "" || len(commission.UserId) < 24 {
			newUserId := utils.GenerateUserID()
			err := cli.Table("Commissions").Where("id = ?", commission.Id).Update("userId", newUserId).Error
			if err != nil {
				log.Printf("更新Commission %d 的userId失败: %v", commission.Id, err)
				continue
			}
			log.Printf("Commission %d 的userId已更新为: %s", commission.Id, newUserId)
		}
	}

	// 5. 迁移Cashouts表
	log.Println("迁移Cashouts表...")
	var cashouts []*model.CashoutModel
	err = cli.Table("Cashouts").Find(&cashouts).Error
	if err != nil {
		return fmt.Errorf("查询Cashouts失败: %v", err)
	}

	for _, cashout := range cashouts {
		if cashout.UserId == "" || len(cashout.UserId) < 24 {
			newUserId := utils.GenerateUserID()
			err := cli.Table("Cashouts").Where("id = ?", cashout.Id).Update("userId", newUserId).Error
			if err != nil {
				log.Printf("更新Cashout %d 的userId失败: %v", cashout.Id, err)
				continue
			}
			log.Printf("Cashout %d 的userId已更新为: %s", cashout.Id, newUserId)
		}
	}

	// 6. 迁移Files表
	log.Println("迁移Files表...")
	var files []*model.FileModel
	err = cli.Table("Files").Find(&files).Error
	if err != nil {
		return fmt.Errorf("查询Files失败: %v", err)
	}

	for _, file := range files {
		if file.UserId == "" || len(file.UserId) < 24 {
			newUserId := utils.GenerateUserID()
			err := cli.Table("Files").Where("id = ?", file.Id).Update("userId", newUserId).Error
			if err != nil {
				log.Printf("更新File %d 的userId失败: %v", file.Id, err)
				continue
			}
			log.Printf("File %d 的userId已更新为: %s", file.Id, newUserId)
		}
	}

	log.Println("所有表的userId字段迁移完成")
	return nil
}

// ValidateUserIds 验证所有用户都有有效的UserId
func (m *MigrationService) ValidateUserIds() error {
	log.Println("验证所有用户的UserId...")

	var count int64
	cli := db.Get()
	err := cli.Table("Users").Where("userId IS NULL OR userId = ''").Count(&count).Error
	if err != nil {
		return fmt.Errorf("验证UserId失败: %v", err)
	}

	if count > 0 {
		return fmt.Errorf("还有 %d 个用户没有UserId", count)
	}

	log.Println("所有用户都有有效的UserId")
	return nil
}

// GetUserByUserId 通过UserId获取用户信息
func (m *MigrationService) GetUserByUserId(userId string) (*model.UserModel, error) {
	return dao.UserImp.GetUserByUserId(userId)
}

// GetUserById 通过数据库ID获取用户信息
func (m *MigrationService) GetUserById(id int32) (*model.UserModel, error) {
	return dao.UserImp.GetUserById(id)
}
