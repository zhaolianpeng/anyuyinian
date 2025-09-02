package dao

import (
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

const userTableName = "Users"

// GetUserByOpenId 根据OpenId查询用户
func (imp *UserInterfaceImp) GetUserByOpenId(openId string) (*model.UserModel, error) {
	var user = new(model.UserModel)
	cli := db.Get()

	// 记录SQL操作日志
	logger := NewSQLLogger("查询", userTableName, map[string]interface{}{
		"openId": openId,
	})

	err := cli.Table(userTableName).Where("openId = ?", openId).First(user).Error
	logger.LogQuery(user, err)

	return user, err
}

// GetUserById 根据用户ID查询用户
func (imp *UserInterfaceImp) GetUserById(id int32) (*model.UserModel, error) {
	var user = new(model.UserModel)
	cli := db.Get()

	// 记录SQL操作日志
	logger := NewSQLLogger("查询", userTableName, map[string]interface{}{
		"id": id,
	})

	err := cli.Table(userTableName).Where("id = ?", id).First(user).Error
	logger.LogQuery(user, err)

	return user, err
}

// GetUserByUserId 根据UserId查询用户
func (imp *UserInterfaceImp) GetUserByUserId(userId string) (*model.UserModel, error) {
	var user = new(model.UserModel)
	cli := db.Get()

	// 记录SQL操作日志
	logger := NewSQLLogger("查询", userTableName, map[string]interface{}{
		"userId": userId,
	})

	err := cli.Table(userTableName).Where("userId = ?", userId).First(user).Error
	logger.LogQuery(user, err)

	return user, err
}

// CreateUser 创建用户
func (imp *UserInterfaceImp) CreateUser(user *model.UserModel) error {
	cli := db.Get()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.LastLoginAt = time.Now()

	// 记录SQL操作日志
	logger := NewSQLLogger("插入", userTableName, map[string]interface{}{
		"openId": user.OpenId,
		"userId": user.UserId,
	})

	err := cli.Table(userTableName).Create(user).Error
	logger.LogInsert(user, err)

	return err
}

// UpdateUser 更新用户信息
func (imp *UserInterfaceImp) UpdateUser(user *model.UserModel) error {
	cli := db.Get()
	user.UpdatedAt = time.Now()
	user.LastLoginAt = time.Now()

	// 记录SQL操作日志
	logger := NewSQLLogger("更新", userTableName, map[string]interface{}{
		"openId": user.OpenId,
		"userId": user.UserId,
	})

	err := cli.Table(userTableName).Where("openId = ?", user.OpenId).Updates(user).Error
	logger.LogUpdate(user, err)

	return err
}

// UpsertUser 更新或创建用户
func (imp *UserInterfaceImp) UpsertUser(user *model.UserModel) error {
	cli := db.Get()
	user.UpdatedAt = time.Now()
	user.LastLoginAt = time.Now()

	// 记录SQL操作日志
	logger := NewSQLLogger("更新或插入", userTableName, map[string]interface{}{
		"openId": user.OpenId,
		"userId": user.UserId,
	})

	err := cli.Table(userTableName).Save(user).Error
	logger.LogUpdate(user, err)

	return err
}
