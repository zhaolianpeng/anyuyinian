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

	// 启动SSE管理器（替代WebSocket）
	go service.SSEManagerInstance.Start()

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
	http.HandleFunc("/api/user/update_info", service.NewLogMiddleware(service.UpdateUserInfoHandler))
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

	// 支付相关接口
	http.HandleFunc("/api/payment/notify", service.NewLogMiddleware(service.HandleWechatPayNotify))

	// 订单超时相关接口
	http.HandleFunc("/api/order/check_expired", service.NewLogMiddleware(service.CheckExpiredOrdersHandler))
	http.HandleFunc("/api/order/expired_count", service.NewLogMiddleware(service.GetExpiredOrdersCountHandler))

	// 推荐相关接口
	http.HandleFunc("/api/referral/qrcode", service.NewLogMiddleware(service.ReferralQrCodeHandler))
	http.HandleFunc("/api/referral/report", service.NewLogMiddleware(service.ReferralReportHandler))
	http.HandleFunc("/api/referral/config", service.NewLogMiddleware(service.ReferralConfigHandler))
	http.HandleFunc("/api/referral/apply_cashout", service.NewLogMiddleware(service.ApplyCashoutHandler))

	// 推广中心相关接口
	http.HandleFunc("/api/promoter/info", service.NewLogMiddleware(service.GetPromoterInfoHandler))
	http.HandleFunc("/api/promoter/commission_list", service.NewLogMiddleware(service.GetCommissionListHandler))
	http.HandleFunc("/api/promoter/cashout_list", service.NewLogMiddleware(service.GetCashoutListHandler))
	http.HandleFunc("/api/promoter/find_user", service.NewLogMiddleware(service.GetUserByPromoterCodeHandler))
	http.HandleFunc("/api/promoter/generate_codes", service.NewLogMiddleware(service.GeneratePromoterCodesHandler))

	// 二维码相关接口
	http.HandleFunc("/api/qrcode/generate", service.NewLogMiddleware(service.GenerateQRCodeHandler))
	http.HandleFunc("/api/qrcode/generate_base64", service.NewLogMiddleware(service.GenerateQRCodeBase64Handler))

	// 客服相关接口
	http.HandleFunc("/api/kefu/send_msg", service.NewLogMiddleware(service.SendMessageHandler))
	http.HandleFunc("/api/kefu/faq", service.NewLogMiddleware(service.FaqHandler))

	// 医院相关接口
	http.HandleFunc("/api/hospital/list", service.NewLogMiddleware(service.HospitalListHandler))
	http.HandleFunc("/api/hospital/detail/", service.NewLogMiddleware(service.HospitalDetailHandler))

	// SSE路由（替代WebSocket）
	http.HandleFunc("/sse", service.NewLogMiddleware(service.SSEHandler))

	// 迁移服务相关接口
	http.HandleFunc("/api/migration/generate_user_ids", service.NewLogMiddleware(service.GenerateUserIdHandler))
	http.HandleFunc("/api/migration/migrate_users", service.NewLogMiddleware(service.MigrateUsersHandler))
	http.HandleFunc("/api/migration/migrate_all_tables", service.NewLogMiddleware(service.MigrateAllTablesUserIdHandler))
	http.HandleFunc("/api/migration/validate", service.NewLogMiddleware(service.ValidateUserIdsHandler))

	// 紧急修复相关接口
	http.HandleFunc("/api/emergency/fix_user_ids", service.NewLogMiddleware(service.EmergencyFixUserIdsHandler))
	http.HandleFunc("/api/emergency/test_user_info", service.NewLogMiddleware(service.TestUserInfoHandler))
	http.HandleFunc("/api/emergency/user_status", service.NewLogMiddleware(service.GetUserStatusHandler))

	// 管理员相关接口
	http.HandleFunc("/api/admin/login", service.NewLogMiddleware(service.AdminLoginHandler))
	http.HandleFunc("/api/admin/check-status", service.NewLogMiddleware(service.CheckAdminStatusHandler))
	http.HandleFunc("/api/admin/users", service.NewLogMiddleware(service.GetAdminUsersHandler))
	http.HandleFunc("/api/admin/orders", service.NewLogMiddleware(service.GetAdminOrdersHandler))
	http.HandleFunc("/api/admin/set-admin", service.NewLogMiddleware(service.SetAdminHandler))
	http.HandleFunc("/api/admin/remove-admin", service.NewLogMiddleware(service.RemoveAdminHandler))
	http.HandleFunc("/api/admin/stats", service.NewLogMiddleware(service.AdminStatsHandler))
	http.HandleFunc("/api/admin/admins", service.NewLogMiddleware(service.AdminAdminsHandler))
	http.HandleFunc("/api/admin/order/update-amount", service.NewLogMiddleware(service.UpdateOrderAmountHandler))
    http.HandleFunc("/api/admin/order/refund", service.NewLogMiddleware(service.AdminRefundOrderHandler))

    // 管理员服务管理相关接口
    http.HandleFunc("/api/admin/services", service.NewLogMiddleware(service.GetAdminServicesHandler))
    http.HandleFunc("/api/admin/service/update-price", service.NewLogMiddleware(service.UpdateServicePriceHandler))

	log.Fatal(http.ListenAndServe(":80", nil))
}
