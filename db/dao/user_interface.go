package dao

import (
	"wxcloudrun-golang/db/model"
)

// UserInterface 用户数据模型接口
type UserInterface interface {
	GetUserByOpenId(openId string) (*model.UserModel, error)
	GetUserById(id int32) (*model.UserModel, error)
	CreateUser(user *model.UserModel) error
	UpdateUser(user *model.UserModel) error
	UpsertUser(user *model.UserModel) error
}

// UserInterfaceImp 用户数据模型实现
type UserInterfaceImp struct{}

// UserImp 用户实现实例
var UserImp UserInterface = &UserInterfaceImp{}
