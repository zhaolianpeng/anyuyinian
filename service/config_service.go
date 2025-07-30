package service

import (
	"encoding/json"
	"net/http"

	"wxcloudrun-golang/db/dao"
)

// ConfigResponse 配置响应
type ConfigResponse struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

// PlatformConfig 平台配置
type PlatformConfig struct {
	CustomerServicePhone string `json:"customerServicePhone"`
	PrivacyPolicyUrl     string `json:"privacyPolicyUrl"`
	UserAgreementUrl     string `json:"userAgreementUrl"`
	AboutUsUrl           string `json:"aboutUsUrl"`
	HelpCenterUrl        string `json:"helpCenterUrl"`
	AppVersion           string `json:"appVersion"`
	ForceUpdate          bool   `json:"forceUpdate"`
	MaintenanceMode      bool   `json:"maintenanceMode"`
	MaintenanceMessage   string `json:"maintenanceMessage"`
}

// ConfigHandler 获取平台配置接口
func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	// 获取所有配置
	configs, err := dao.ConfigImp.GetAllConfigs()
	if err != nil {
		response := &ConfigResponse{
			Code:     -1,
			ErrorMsg: "获取配置失败: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 构建平台配置
	platformConfig := &PlatformConfig{}

	for _, config := range configs {
		switch config.Key {
		case "customer_service_phone":
			platformConfig.CustomerServicePhone = config.Value
		case "privacy_policy_url":
			platformConfig.PrivacyPolicyUrl = config.Value
		case "user_agreement_url":
			platformConfig.UserAgreementUrl = config.Value
		case "about_us_url":
			platformConfig.AboutUsUrl = config.Value
		case "help_center_url":
			platformConfig.HelpCenterUrl = config.Value
		case "app_version":
			platformConfig.AppVersion = config.Value
		case "force_update":
			platformConfig.ForceUpdate = config.Value == "true"
		case "maintenance_mode":
			platformConfig.MaintenanceMode = config.Value == "true"
		case "maintenance_message":
			platformConfig.MaintenanceMessage = config.Value
		}
	}

	response := &ConfigResponse{
		Code: 0,
		Data: platformConfig,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
