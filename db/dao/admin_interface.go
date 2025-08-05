package dao

import (
	"wxcloudrun-golang/db/model"
)

// AdminInterface 管理员数据访问接口
type AdminInterface interface {
	// 管理员登录
	AdminLogin(username, password string) (*model.UserModel, error)

	// 获取管理员信息
	GetAdminByUserId(userId string) (*model.UserModel, error)

	// 获取所有管理员列表
	GetAllAdmins(page, pageSize int) ([]*model.UserModel, int64, error)

	// 设置用户为管理员
	SetUserAsAdmin(userId string, adminLevel int, parentAdminId string) error

	// 取消用户管理员权限
	RemoveAdmin(userId string) error

	// 获取下级管理员列表
	GetSubAdmins(parentAdminId string, page, pageSize int) ([]*model.UserModel, int64, error)

	// 获取管理员可见的用户列表
	GetVisibleUsers(adminUserId string, page, pageSize int) ([]*model.UserModel, int64, error)

	// 获取管理员可见的订单列表
	GetVisibleOrders(adminUserId string, page, pageSize int) ([]*model.OrderModel, int64, error)

	// 记录管理员登录日志
	LogAdminLogin(log *model.AdminLoginLogModel) error

	// 获取管理员登录日志
	GetAdminLoginLogs(adminUserId string, page, pageSize int) ([]*model.AdminLoginLogModel, int64, error)
}
