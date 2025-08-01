package model

import "time"

// OrderModel 订单模型
type OrderModel struct {
	Id               int32      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OrderNo          string     `gorm:"column:orderNo;uniqueIndex;not null" json:"orderNo"`
	UserId           int32      `gorm:"column:userId;not null" json:"userId"`
	ServiceId        int32      `gorm:"column:serviceId;not null" json:"serviceId"`
	PatientId        int32      `gorm:"column:patientId;not null" json:"patientId"`                // 患者ID
	AddressId        int32      `gorm:"column:addressId;not null" json:"addressId"`                // 地址ID
	AppointmentDate  string     `gorm:"column:appointmentDate;not null" json:"appointmentDate"`    // 预约日期
	AppointmentTime  string     `gorm:"column:appointmentTime;not null" json:"appointmentTime"`    // 预约时间
	DiseaseInfo      string     `gorm:"column:diseaseInfo" json:"diseaseInfo"`                     // 既往病史
	NeedToiletAssist int        `gorm:"column:needToiletAssist;default:0" json:"needToiletAssist"` // 是否需要助排二便：0-不需要，1-需要
	ServiceName      string     `gorm:"column:serviceName;not null" json:"serviceName"`
	Price            float64    `gorm:"column:price;not null" json:"price"`
	Quantity         int        `gorm:"column:quantity;default:1" json:"quantity"`
	TotalAmount      float64    `gorm:"column:totalAmount;not null" json:"totalAmount"`
	FormData         string     `gorm:"column:formData" json:"formData"`             // JSON格式的表单数据
	Status           int        `gorm:"column:status;default:0" json:"status"`       // 0-待支付，1-已支付，2-已完成，3-已取消，4-已退款
	PayStatus        int        `gorm:"column:payStatus;default:0" json:"payStatus"` // 0-未支付，1-已支付
	PayTime          *time.Time `gorm:"column:payTime" json:"payTime"`
	PayMethod        string     `gorm:"column:payMethod" json:"payMethod"`                 // 支付方式：wechat, alipay等
	TransactionId    string     `gorm:"column:transactionId" json:"transactionId"`         // 第三方支付交易号
	RefundStatus     int        `gorm:"column:refundStatus;default:0" json:"refundStatus"` // 0-未退款，1-退款中，2-已退款
	RefundTime       *time.Time `gorm:"column:refundTime" json:"refundTime"`
	RefundAmount     float64    `gorm:"column:refundAmount" json:"refundAmount"`
	RefundReason     string     `gorm:"column:refundReason" json:"refundReason"`
	Remark           string     `gorm:"column:remark" json:"remark"`
	ReferrerId       int32      `gorm:"column:referrerId" json:"referrerId"` // 推荐人ID
	Commission       float64    `gorm:"column:commission" json:"commission"` // 佣金金额
	CreatedAt        time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
}

// TableName 指定表名
func (OrderModel) TableName() string {
	return "Orders"
}
