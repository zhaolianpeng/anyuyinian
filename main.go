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

	http.HandleFunc("/", service.IndexHandler)
	http.HandleFunc("/api/count", service.CounterHandler)
	http.HandleFunc("/api/wx/login", service.WxLoginHandler)
	http.HandleFunc("/api/home/init", service.HomeInitHandler)
	http.HandleFunc("/api/upload", service.UploadHandler)
	http.HandleFunc("/api/files", service.GetFileListHandler)
	http.HandleFunc("/api/config", service.ConfigHandler)
	http.HandleFunc("/api/user/info", service.GetUserInfoHandler)
	http.HandleFunc("/api/user/bind_phone", service.BindPhoneHandler)
	http.HandleFunc("/api/user/address", service.AddressHandler)
	http.HandleFunc("/api/user/patient", service.PatientHandler)
	http.HandleFunc("/api/service/list", service.ServiceListHandler)
	http.HandleFunc("/api/service/detail/", service.ServiceDetailHandler)
	http.HandleFunc("/api/service/form_config/", service.ServiceFormConfigHandler)
	http.HandleFunc("/api/order/submit", service.SubmitOrderHandler)
	http.HandleFunc("/api/order/pay/", service.PayOrderHandler)
	http.HandleFunc("/api/order/cancel/", service.CancelOrderHandler)
	http.HandleFunc("/api/order/refund/", service.RefundOrderHandler)
	http.HandleFunc("/api/order/list", service.OrderListHandler)
	http.HandleFunc("/api/order/detail/", service.OrderDetailHandler)
	http.HandleFunc("/api/referral/qrcode", service.ReferralQrCodeHandler)
	http.HandleFunc("/api/referral/report", service.ReferralReportHandler)
	http.HandleFunc("/api/referral/config", service.ReferralConfigHandler)
	http.HandleFunc("/api/referral/apply_cashout", service.ApplyCashoutHandler)
	http.HandleFunc("/api/kefu/send_msg", service.SendMessageHandler)
	http.HandleFunc("/api/kefu/faq", service.FaqHandler)
	http.HandleFunc("/api/hospital/list", service.HospitalListHandler)
	http.HandleFunc("/api/hospital/detail/", service.HospitalDetailHandler)

	log.Fatal(http.ListenAndServe(":80", nil))
}
