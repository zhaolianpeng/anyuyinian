package main

import (
	"fmt"
	"log"
	"net/http"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/service"
)

func main() {
	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("mysql init failed with %+v", err))
	}

	// 基础页面和统计接口
	http.HandleFunc("/", service.IndexHandler)
	http.HandleFunc("/api/count", service.CounterHandler)

	// 微信登录相关接口
	http.HandleFunc("/api/wx/login", service.WxLoginHandler)

	// 首页初始化接口
	http.HandleFunc("/api/home/init", service.HomeInitHandler)

	// 文件上传和管理接口
	http.HandleFunc("/api/upload", service.UploadHandler)
	http.HandleFunc("/api/files", service.GetFileListHandler)

	// 系统配置接口
	http.HandleFunc("/api/config", service.ConfigHandler)

	// 用户相关接口
	http.HandleFunc("/api/user/info", service.GetUserInfoHandler)
	http.HandleFunc("/api/user/bind_phone", service.BindPhoneHandler)
	http.HandleFunc("/api/user/address", service.AddressHandler)
	http.HandleFunc("/api/user/patient", service.PatientHandler)

	// 服务相关接口
	http.HandleFunc("/api/service/list", service.ServiceListHandler)
	http.HandleFunc("/api/service/detail/", service.ServiceDetailHandler)
	http.HandleFunc("/api/service/form_config/", service.ServiceFormConfigHandler)

	// 订单相关接口
	http.HandleFunc("/api/order/submit", service.SubmitOrderHandler)
	http.HandleFunc("/api/order/pay/", service.PayOrderHandler)
	http.HandleFunc("/api/order/cancel/", service.CancelOrderHandler)
	http.HandleFunc("/api/order/refund/", service.RefundOrderHandler)
	http.HandleFunc("/api/order/list", service.OrderListHandler)
	http.HandleFunc("/api/order/detail/", service.OrderDetailHandler)

	// 推荐相关接口
	http.HandleFunc("/api/referral/qrcode", service.ReferralQrCodeHandler)
	http.HandleFunc("/api/referral/report", service.ReferralReportHandler)
	http.HandleFunc("/api/referral/config", service.ReferralConfigHandler)
	http.HandleFunc("/api/referral/apply_cashout", service.ApplyCashoutHandler)

	// 客服相关接口
	http.HandleFunc("/api/kefu/send_msg", service.SendMessageHandler)
	http.HandleFunc("/api/kefu/faq", service.FaqHandler)

	// 医院相关接口
	http.HandleFunc("/api/hospital/list", service.HospitalListHandler)
	http.HandleFunc("/api/hospital/detail/", service.HospitalDetailHandler)

	log.Fatal(http.ListenAndServe(":80", nil))
}
