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

	// 初始化订单超时处理服务
	service.InitOrderTimeoutService()

	// 启动WebSocket管理器
	go service.WsManager.Start()

	// 基础页面和统计接口
	http.HandleFunc("/", service.NewLogMiddleware(service.IndexHandler))
	http.HandleFunc("/api/count", service.NewLogMiddleware(service.CounterHandler))

	// 微信登录相关接口
	http.HandleFunc("/api/wx/login", service.NewLogMiddleware(service.WxLoginHandler))

	// 首页初始化接口
	http.HandleFunc("/api/home/init", service.NewLogMiddleware(service.HomeInitHandler))

	// 文件上传和管理接口
	http.HandleFunc("/api/upload", service.NewLogMiddleware(service.UploadHandler))
	http.HandleFunc("/api/files", service.NewLogMiddleware(service.GetFileListHandler))
	http.HandleFunc("/api/file/delete", service.NewLogMiddleware(service.DeleteFileHandler))
	http.HandleFunc("/api/file/permission", service.NewLogMiddleware(service.UpdateFilePermissionHandler))
	http.HandleFunc("/api/file/permission/get", service.NewLogMiddleware(service.GetFilePermissionHandler))

	// 系统配置接口
	http.HandleFunc("/api/config", service.NewLogMiddleware(service.ConfigHandler))

	// 用户相关接口
	http.HandleFunc("/api/user/info", service.NewLogMiddleware(service.GetUserInfoHandler))
	http.HandleFunc("/api/user/bind_phone", service.NewLogMiddleware(service.BindPhoneHandler))
	http.HandleFunc("/api/user/address", service.NewLogMiddleware(service.AddressHandler))
	http.HandleFunc("/api/user/patient", service.NewLogMiddleware(service.PatientHandler))

	// 服务相关接口
	http.HandleFunc("/api/service/list", service.NewLogMiddleware(service.ServiceListHandler))
	http.HandleFunc("/api/service/detail", service.NewLogMiddleware(service.ServiceDetailHandler))
	http.HandleFunc("/api/service/form_config/", service.NewLogMiddleware(service.ServiceFormConfigHandler))

	// 订单相关接口
	http.HandleFunc("/api/order/submit", service.NewLogMiddleware(service.SubmitOrderHandler))
	http.HandleFunc("/api/order/pay/", service.NewLogMiddleware(service.PayOrderHandler))
	http.HandleFunc("/api/order/pay_confirm/", service.NewLogMiddleware(service.PayConfirmHandler))
	http.HandleFunc("/api/order/cancel/", service.NewLogMiddleware(service.CancelOrderHandler))
	http.HandleFunc("/api/order/refund/", service.NewLogMiddleware(service.RefundOrderHandler))
	http.HandleFunc("/api/order/list", service.NewLogMiddleware(service.OrderListHandler))
	http.HandleFunc("/api/order/detail", service.NewLogMiddleware(service.OrderDetailHandler))
	http.HandleFunc("/api/order/time_slots", service.NewLogMiddleware(service.GetAvailableTimeSlotsHandler))

	// 订单超时相关接口
	http.HandleFunc("/api/order/check_expired", service.NewLogMiddleware(service.CheckExpiredOrdersHandler))
	http.HandleFunc("/api/order/expired_count", service.NewLogMiddleware(service.GetExpiredOrdersCountHandler))

	// 推荐相关接口
	http.HandleFunc("/api/referral/qrcode", service.NewLogMiddleware(service.ReferralQrCodeHandler))
	http.HandleFunc("/api/referral/report", service.NewLogMiddleware(service.ReferralReportHandler))
	http.HandleFunc("/api/referral/config", service.NewLogMiddleware(service.ReferralConfigHandler))
	http.HandleFunc("/api/referral/apply_cashout", service.NewLogMiddleware(service.ApplyCashoutHandler))

	// 客服相关接口
	http.HandleFunc("/api/kefu/send_msg", service.NewLogMiddleware(service.SendMessageHandler))
	http.HandleFunc("/api/kefu/faq", service.NewLogMiddleware(service.FaqHandler))

	// 医院相关接口
	http.HandleFunc("/api/hospital/list", service.NewLogMiddleware(service.HospitalListHandler))
	http.HandleFunc("/api/hospital/detail/", service.NewLogMiddleware(service.HospitalDetailHandler))

	// WebSocket路由
	http.HandleFunc("/ws", service.NewLogMiddleware(service.WebSocketHandler))

	log.Fatal(http.ListenAndServe(":80", nil))
}
