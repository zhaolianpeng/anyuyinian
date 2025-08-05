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
