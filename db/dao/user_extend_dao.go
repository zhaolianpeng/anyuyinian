package dao

import (
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

const addressTableName = "UserAddresses"
const patientTableName = "Patients"

// 地址相关方法

// CreateAddress 创建地址
func (imp *UserExtendInterfaceImp) CreateAddress(address *model.UserAddressModel) error {
	cli := db.Get()
	address.CreatedAt = time.Now()
	address.UpdatedAt = time.Now()
	return cli.Table(addressTableName).Create(address).Error
}

// GetAddressById 根据ID获取地址
func (imp *UserExtendInterfaceImp) GetAddressById(id int32) (*model.UserAddressModel, error) {
	var address = new(model.UserAddressModel)
	cli := db.Get()
	err := cli.Table(addressTableName).Where("id = ? AND status = ?", id, 1).First(address).Error
	return address, err
}

// GetAddressesByUserId 根据用户ID获取地址列表
func (imp *UserExtendInterfaceImp) GetAddressesByUserId(userId int32) ([]*model.UserAddressModel, error) {
	var addresses []*model.UserAddressModel
	cli := db.Get()
	err := cli.Table(addressTableName).
		Where("userId = ? AND status = ?", userId, 1).
		Order("isDefault DESC, createdAt DESC").
		Find(&addresses).Error
	return addresses, err
}

// UpdateAddress 更新地址
func (imp *UserExtendInterfaceImp) UpdateAddress(address *model.UserAddressModel) error {
	cli := db.Get()
	address.UpdatedAt = time.Now()
	return cli.Table(addressTableName).Where("id = ?", address.Id).Updates(address).Error
}

// DeleteAddress 删除地址（软删除）
func (imp *UserExtendInterfaceImp) DeleteAddress(id int32) error {
	cli := db.Get()
	return cli.Table(addressTableName).Where("id = ?", id).Update("status", 0).Error
}

// SetDefaultAddress 设置默认地址
func (imp *UserExtendInterfaceImp) SetDefaultAddress(userId int32, addressId int32) error {
	cli := db.Get()

	// 先取消该用户的所有默认地址
	err := cli.Table(addressTableName).
		Where("userId = ? AND status = ?", userId, 1).
		Update("isDefault", 0).Error
	if err != nil {
		return err
	}

	// 设置指定地址为默认
	return cli.Table(addressTableName).
		Where("id = ? AND userId = ?", addressId, userId).
		Update("isDefault", 1).Error
}

// 就诊人相关方法

// CreatePatient 创建就诊人
func (imp *UserExtendInterfaceImp) CreatePatient(patient *model.PatientModel) error {
	cli := db.Get()
	patient.CreatedAt = time.Now()
	patient.UpdatedAt = time.Now()
	return cli.Table(patientTableName).Create(patient).Error
}

// GetPatientById 根据ID获取就诊人
func (imp *UserExtendInterfaceImp) GetPatientById(id int32) (*model.PatientModel, error) {
	var patient = new(model.PatientModel)
	cli := db.Get()
	err := cli.Table(patientTableName).Where("id = ? AND status = ?", id, 1).First(patient).Error
	return patient, err
}

// GetPatientsByUserId 根据用户ID获取就诊人列表
func (imp *UserExtendInterfaceImp) GetPatientsByUserId(userId int32) ([]*model.PatientModel, error) {
	var patients []*model.PatientModel
	cli := db.Get()
	err := cli.Table(patientTableName).
		Where("userId = ? AND status = ?", userId, 1).
		Order("isDefault DESC, createdAt DESC").
		Find(&patients).Error
	return patients, err
}

// UpdatePatient 更新就诊人
func (imp *UserExtendInterfaceImp) UpdatePatient(patient *model.PatientModel) error {
	cli := db.Get()
	patient.UpdatedAt = time.Now()
	return cli.Table(patientTableName).Where("id = ?", patient.Id).Updates(patient).Error
}

// DeletePatient 删除就诊人（软删除）
func (imp *UserExtendInterfaceImp) DeletePatient(id int32) error {
	cli := db.Get()
	return cli.Table(patientTableName).Where("id = ?", id).Update("status", 0).Error
}

// SetDefaultPatient 设置默认就诊人
func (imp *UserExtendInterfaceImp) SetDefaultPatient(userId int32, patientId int32) error {
	cli := db.Get()

	// 先取消该用户的所有默认就诊人
	err := cli.Table(patientTableName).
		Where("userId = ? AND status = ?", userId, 1).
		Update("isDefault", 0).Error
	if err != nil {
		return err
	}

	// 设置指定就诊人为默认
	return cli.Table(patientTableName).
		Where("id = ? AND userId = ?", patientId, userId).
		Update("isDefault", 1).Error
}
