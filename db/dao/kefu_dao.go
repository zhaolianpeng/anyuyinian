package dao

import (
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"

	"gorm.io/gorm"
)

const kefuMessageTableName = "KefuMessages"
const faqTableName = "Faqs"

// 客服消息相关方法

// CreateMessage 创建客服消息
func (imp *KefuInterfaceImp) CreateMessage(message *model.KefuMessageModel) error {
	cli := db.Get()
	message.CreatedAt = time.Now()
	message.UpdatedAt = time.Now()
	return cli.Table(kefuMessageTableName).Create(message).Error
}

// GetMessagesByUserId 根据用户ID获取消息列表（分页）
func (imp *KefuInterfaceImp) GetMessagesByUserId(userId string, page, pageSize int) ([]*model.KefuMessageModel, int64, error) {
	var messages []*model.KefuMessageModel
	var total int64
	cli := db.Get()

	// 获取总数
	err := cli.Table(kefuMessageTableName).Where("userId = ?", userId).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = cli.Table(kefuMessageTableName).
		Where("userId = ?", userId).
		Order("createdAt DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&messages).Error

	return messages, total, err
}

// GetMessageById 根据ID获取消息
func (imp *KefuInterfaceImp) GetMessageById(id int32) (*model.KefuMessageModel, error) {
	var message = new(model.KefuMessageModel)
	cli := db.Get()
	err := cli.Table(kefuMessageTableName).Where("id = ?", id).First(message).Error
	return message, err
}

// UpdateMessage 更新消息
func (imp *KefuInterfaceImp) UpdateMessage(message *model.KefuMessageModel) error {
	cli := db.Get()
	message.UpdatedAt = time.Now()
	return cli.Table(kefuMessageTableName).Where("id = ?", message.Id).Updates(message).Error
}

// UpdateMessageStatus 更新消息状态
func (imp *KefuInterfaceImp) UpdateMessageStatus(id int32, status int) error {
	cli := db.Get()
	return cli.Table(kefuMessageTableName).Where("id = ?", id).Updates(map[string]interface{}{
		"status":    status,
		"updatedAt": time.Now(),
	}).Error
}

// ReplyMessage 回复消息
func (imp *KefuInterfaceImp) ReplyMessage(id int32, replyContent string, replyUserId string) error {
	cli := db.Get()
	replyTime := time.Now()
	return cli.Table(kefuMessageTableName).Where("id = ?", id).Updates(map[string]interface{}{
		"replyContent": replyContent,
		"replyTime":    replyTime,
		"replyUserId":  replyUserId,
		"status":       2, // 已回复
		"updatedAt":    replyTime,
	}).Error
}

// FAQ相关方法

// GetFaqsByCategory 根据分类获取FAQ列表（分页）
func (imp *KefuInterfaceImp) GetFaqsByCategory(category string, page, pageSize int) ([]*model.FaqModel, int64, error) {
	var faqs []*model.FaqModel
	var total int64
	cli := db.Get()

	// 获取总数
	err := cli.Table(faqTableName).Where("category = ? AND status = ?", category, 1).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = cli.Table(faqTableName).
		Where("category = ? AND status = ?", category, 1).
		Order("sort ASC, createdAt DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&faqs).Error

	return faqs, total, err
}

// GetAllFaqs 获取所有FAQ列表（分页）
func (imp *KefuInterfaceImp) GetAllFaqs(page, pageSize int) ([]*model.FaqModel, int64, error) {
	var faqs []*model.FaqModel
	var total int64
	cli := db.Get()

	// 获取总数
	err := cli.Table(faqTableName).Where("status = ?", 1).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = cli.Table(faqTableName).
		Where("status = ?", 1).
		Order("sort ASC, createdAt DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&faqs).Error

	return faqs, total, err
}

// GetFaqById 根据ID获取FAQ
func (imp *KefuInterfaceImp) GetFaqById(id int32) (*model.FaqModel, error) {
	var faq = new(model.FaqModel)
	cli := db.Get()
	err := cli.Table(faqTableName).Where("id = ? AND status = ?", id, 1).First(faq).Error
	return faq, err
}

// GetFaqCategories 获取FAQ分类列表
func (imp *KefuInterfaceImp) GetFaqCategories() ([]string, error) {
	var categories []string
	cli := db.Get()
	err := cli.Table(faqTableName).
		Where("status = ?", 1).
		Distinct("category").
		Pluck("category", &categories).Error
	return categories, err
}

// IncrementFaqViewCount 增加FAQ查看次数
func (imp *KefuInterfaceImp) IncrementFaqViewCount(id int32) error {
	cli := db.Get()
	return cli.Table(faqTableName).Where("id = ?", id).Update("viewCount", gorm.Expr("viewCount + ?", 1)).Error
}
