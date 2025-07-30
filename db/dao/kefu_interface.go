package dao

import (
	"wxcloudrun-golang/db/model"
)

// KefuInterface 客服数据接口
type KefuInterface interface {
	CreateMessage(message *model.KefuMessageModel) error
	GetMessagesByUserId(userId int32, page, pageSize int) ([]*model.KefuMessageModel, int64, error)
	GetMessageById(id int32) (*model.KefuMessageModel, error)
	UpdateMessage(message *model.KefuMessageModel) error
	UpdateMessageStatus(id int32, status int) error
	ReplyMessage(id int32, replyContent string, replyUserId int32) error

	// FAQ相关
	GetFaqsByCategory(category string, page, pageSize int) ([]*model.FaqModel, int64, error)
	GetAllFaqs(page, pageSize int) ([]*model.FaqModel, int64, error)
	GetFaqById(id int32) (*model.FaqModel, error)
	GetFaqCategories() ([]string, error)
	IncrementFaqViewCount(id int32) error
}

// KefuInterfaceImp 客服数据实现
type KefuInterfaceImp struct{}

// KefuImp 客服实现实例
var KefuImp KefuInterface = &KefuInterfaceImp{}
