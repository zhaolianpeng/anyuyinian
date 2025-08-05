package dao

import (
	"wxcloudrun-golang/db/model"
)

// UserExtendInterface 用户扩展数据接口
type UserExtendInterface interface {
	// 地址相关
	CreateAddress(address *model.UserAddressModel) error
	GetAddressById(id int32) (*model.UserAddressModel, error)
	GetAddressesByUserId(userId string) ([]*model.UserAddressModel, error)
	UpdateAddress(address *model.UserAddressModel) error
	DeleteAddress(id int32) error
	SetDefaultAddress(userId string, addressId int32) error

	// 就诊人相关
	CreatePatient(patient *model.PatientModel) error
	GetPatientById(id int32) (*model.PatientModel, error)
	GetPatientsByUserId(userId string) ([]*model.PatientModel, error)
	UpdatePatient(patient *model.PatientModel) error
	DeletePatient(id int32) error
	SetDefaultPatient(userId string, patientId int32) error
}

// UserExtendInterfaceImp 用户扩展数据实现
type UserExtendInterfaceImp struct{}

// UserExtendImp 用户扩展实现实例
var UserExtendImp UserExtendInterface = &UserExtendInterfaceImp{}
