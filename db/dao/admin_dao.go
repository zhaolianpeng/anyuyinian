package dao

import (
	"fmt"
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

// AdminImp 管理员数据访问实现
type AdminImp struct{}

// AdminLogin 管理员登录
func (a *AdminImp) AdminLogin(username, password string) (*model.UserModel, error) {
	var user model.UserModel
	cli := db.Get()

	err := cli.Table("Users").Where("adminUsername = ? AND adminPassword = ? AND isAdmin = 1", username, password).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("管理员登录失败: %v", err)
	}

	return &user, nil
}

// GetAdminByUserId 获取管理员信息
func (a *AdminImp) GetAdminByUserId(userId string) (*model.UserModel, error) {
	var user model.UserModel
	cli := db.Get()

	err := cli.Table("Users").Where("userId = ? AND isAdmin = 1", userId).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("获取管理员信息失败: %v", err)
	}

	return &user, nil
}

// GetAllAdmins 获取所有管理员列表
func (a *AdminImp) GetAllAdmins(page, pageSize int) ([]*model.UserModel, int64, error) {
	var users []*model.UserModel
	var total int64
	cli := db.Get()

	offset := (page - 1) * pageSize

	// 获取总数
	err := cli.Table("Users").Where("isAdmin = 1").Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("获取管理员总数失败: %v", err)
	}

	// 获取列表
	err = cli.Table("Users").Where("isAdmin = 1").
		Order("adminLevel DESC, adminCreatedAt DESC").
		Offset(offset).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, fmt.Errorf("获取管理员列表失败: %v", err)
	}

	return users, total, nil
}

// SetUserAsAdmin 设置用户为管理员
func (a *AdminImp) SetUserAsAdmin(userId string, adminLevel int, parentAdminId string) error {
	cli := db.Get()
	now := time.Now()

	updates := map[string]interface{}{
		"isAdmin":        1,
		"adminLevel":     adminLevel,
		"parentAdminId":  parentAdminId,
		"adminCreatedAt": now,
		"updatedAt":      now,
	}

	err := cli.Table("Users").Where("userId = ?", userId).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("设置用户为管理员失败: %v", err)
	}

	return nil
}

// RemoveAdmin 取消用户管理员权限
func (a *AdminImp) RemoveAdmin(userId string) error {
	cli := db.Get()
	now := time.Now()

	updates := map[string]interface{}{
		"isAdmin":        0,
		"adminLevel":     0,
		"parentAdminId":  nil,
		"adminCreatedAt": nil,
		"updatedAt":      now,
	}

	err := cli.Table("Users").Where("userId = ?", userId).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("取消管理员权限失败: %v", err)
	}

	return nil
}

// GetSubAdmins 获取下级管理员列表
func (a *AdminImp) GetSubAdmins(parentAdminId string, page, pageSize int) ([]*model.UserModel, int64, error) {
	var users []*model.UserModel
	var total int64
	cli := db.Get()

	offset := (page - 1) * pageSize

	// 获取总数
	err := cli.Table("Users").Where("parentAdminId = ? AND isAdmin = 1", parentAdminId).Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("获取下级管理员总数失败: %v", err)
	}

	// 获取列表
	err = cli.Table("Users").Where("parentAdminId = ? AND isAdmin = 1", parentAdminId).
		Order("adminCreatedAt DESC").
		Offset(offset).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, fmt.Errorf("获取下级管理员列表失败: %v", err)
	}

	return users, total, nil
}

// GetVisibleUsers 获取管理员可见的用户列表
func (a *AdminImp) GetVisibleUsers(adminUserId string, page, pageSize int) ([]*model.UserModel, int64, error) {
	var users []*model.UserModel
	var total int64
	cli := db.Get()

	offset := (page - 1) * pageSize

	// 先获取管理员信息
	admin, err := a.GetAdminByUserId(adminUserId)
	if err != nil {
		return nil, 0, err
	}

	if admin.AdminLevel == 2 { // 超级管理员可以看到所有用户
		// 获取总数
		err = cli.Table("Users").Count(&total).Error
		if err != nil {
			return nil, 0, fmt.Errorf("获取用户总数失败: %v", err)
		}

		// 获取列表
		err = cli.Table("Users").Order("createdAt DESC").Offset(offset).Limit(pageSize).Find(&users).Error
		if err != nil {
			return nil, 0, fmt.Errorf("获取用户列表失败: %v", err)
		}
	} else { // 一级管理员只能看到通过自己推广码注册的用户
		// 通过推广关系查找下级用户
		subQuery := cli.Table("Referrals").Select("userId").Where("referrerId = ?", adminUserId)

		// 获取总数
		err = cli.Table("Users").Where("userId IN (?)", subQuery).Count(&total).Error
		if err != nil {
			return nil, 0, fmt.Errorf("获取用户总数失败: %v", err)
		}

		// 获取列表
		err = cli.Table("Users").Where("userId IN (?)", subQuery).Order("createdAt DESC").Offset(offset).Limit(pageSize).Find(&users).Error
		if err != nil {
			return nil, 0, fmt.Errorf("获取用户列表失败: %v", err)
		}
	}
	if err != nil {
		return nil, 0, fmt.Errorf("获取用户列表失败: %v", err)
	}

	return users, total, nil
}

// GetVisibleOrders 获取管理员可见的订单列表
func (a *AdminImp) GetVisibleOrders(adminUserId string, page, pageSize int) ([]*model.OrderModel, int64, error) {
	var orders []*model.OrderModel
	var total int64
	cli := db.Get()

	offset := (page - 1) * pageSize

	// 先获取管理员信息
	admin, err := a.GetAdminByUserId(adminUserId)
	if err != nil {
		return nil, 0, err
	}

	if admin.AdminLevel == 2 { // 超级管理员可以看到所有订单
		// 获取总数
		err = cli.Table("Orders").Count(&total).Error
		if err != nil {
			return nil, 0, fmt.Errorf("获取订单总数失败: %v", err)
		}

		// 获取列表
		err = cli.Table("Orders").Order("createdAt DESC").Offset(offset).Limit(pageSize).Find(&orders).Error
		if err != nil {
			return nil, 0, fmt.Errorf("获取订单列表失败: %v", err)
		}
	} else { // 一级管理员只能看到通过自己推广码注册用户的订单
		// 通过推广关系查找下级用户的订单
		subQuery := cli.Table("Referrals").Select("userId").Where("referrerId = ?", adminUserId)

		// 获取总数
		err = cli.Table("Orders").Where("userId IN (?)", subQuery).Count(&total).Error
		if err != nil {
			return nil, 0, fmt.Errorf("获取订单总数失败: %v", err)
		}

		// 获取列表
		err = cli.Table("Orders").Where("userId IN (?)", subQuery).Order("createdAt DESC").Offset(offset).Limit(pageSize).Find(&orders).Error
		if err != nil {
			return nil, 0, fmt.Errorf("获取订单列表失败: %v", err)
		}
	}
	if err != nil {
		return nil, 0, fmt.Errorf("获取订单列表失败: %v", err)
	}

	return orders, total, nil
}

// LogAdminLogin 记录管理员登录日志
func (a *AdminImp) LogAdminLogin(log *model.AdminLoginLogModel) error {
	cli := db.Get()
	return cli.Table("AdminLoginLogs").Create(log).Error
}

// GetAdminLoginLogs 获取管理员登录日志
func (a *AdminImp) GetAdminLoginLogs(adminUserId string, page, pageSize int) ([]*model.AdminLoginLogModel, int64, error) {
	var logs []*model.AdminLoginLogModel
	var total int64
	cli := db.Get()

	offset := (page - 1) * pageSize

	// 获取总数
	err := cli.Table("AdminLoginLogs").Where("adminUserId = ?", adminUserId).Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("获取登录日志总数失败: %v", err)
	}

	// 获取列表
	err = cli.Table("AdminLoginLogs").Where("adminUserId = ?", adminUserId).
		Order("loginTime DESC").Offset(offset).Limit(pageSize).Find(&logs).Error
	if err != nil {
		return nil, 0, fmt.Errorf("获取登录日志列表失败: %v", err)
	}

	return logs, total, nil
}
