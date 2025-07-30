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
	err := cli.Table(userTableName).Where("openId = ?", openId).First(user).Error
	return user, err
}

// CreateUser 创建用户
func (imp *UserInterfaceImp) CreateUser(user *model.UserModel) error {
	cli := db.Get()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.LastLoginAt = time.Now()
	return cli.Table(userTableName).Create(user).Error
}

// UpdateUser 更新用户信息
func (imp *UserInterfaceImp) UpdateUser(user *model.UserModel) error {
	cli := db.Get()
	user.UpdatedAt = time.Now()
	user.LastLoginAt = time.Now()
	return cli.Table(userTableName).Where("openId = ?", user.OpenId).Updates(user).Error
}

// UpsertUser 更新或创建用户
func (imp *UserInterfaceImp) UpsertUser(user *model.UserModel) error {
	cli := db.Get()
	user.UpdatedAt = time.Now()
	user.LastLoginAt = time.Now()
	return cli.Table(userTableName).Save(user).Error
}
