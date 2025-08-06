package service

import (
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"wxcloudrun-golang/config"
	"wxcloudrun-golang/db/model"
)

// WechatPayRequest 微信支付请求参数
type WechatPayRequest struct {
	AppID          string `xml:"appid"`
	MchID          string `xml:"mch_id"`
	NonceStr       string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
	Body           string `xml:"body"`
	OutTradeNo     string `xml:"out_trade_no"`
	TotalFee       int    `xml:"total_fee"`
	SpbillCreateIP string `xml:"spbill_create_ip"`
	NotifyURL      string `xml:"notify_url"`
	TradeType      string `xml:"trade_type"`
	OpenID         string `xml:"openid"`
}

// WechatPayResponse 微信支付响应
type WechatPayResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	ResultCode string `xml:"result_code"`
	PrepayID   string `xml:"prepay_id"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
}

// WechatPayNotifyResponse 微信支付通知响应
type WechatPayNotifyResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
}

// GenerateWechatPayParams 生成微信支付参数
func GenerateWechatPayParams(order *model.OrderModel, openID string) (map[string]interface{}, error) {
	LogStep("开始生成微信支付参数", map[string]interface{}{
		"orderId": order.Id,
		"orderNo": order.OrderNo,
		"amount":  order.TotalAmount,
		"openID":  openID,
	})

	// 获取支付配置
	paymentConfig := config.GetPaymentConfig()
	wechatConfig := paymentConfig.WechatPay

	// 验证配置
	if wechatConfig.MchID == "" || wechatConfig.MchKey == "" {
		LogError("微信支付配置不完整", fmt.Errorf("商户号或商户密钥未配置"))
		return nil, fmt.Errorf("微信支付配置不完整")
	}

	// 生成随机字符串
	nonceStr := generateNonceStr()

	// 构建请求参数
	request := &WechatPayRequest{
		AppID:          wechatConfig.AppID,
		MchID:          wechatConfig.MchID,
		NonceStr:       nonceStr,
		Body:           fmt.Sprintf("订单支付-%s", order.ServiceName),
		OutTradeNo:     order.OrderNo,
		TotalFee:       int(order.TotalAmount * 100), // 转换为分
		SpbillCreateIP: "127.0.0.1",                  // 客户端IP，实际应该从请求中获取
		NotifyURL:      wechatConfig.NotifyURL,
		TradeType:      "JSAPI", // 小程序支付
		OpenID:         openID,
	}

	// 生成签名
	request.Sign = generateWechatPaySign(request, wechatConfig.MchKey)

	LogStep("微信支付请求参数构建完成", map[string]interface{}{
		"appID":      request.AppID,
		"mchID":      request.MchID,
		"outTradeNo": request.OutTradeNo,
		"totalFee":   request.TotalFee,
		"tradeType":  request.TradeType,
	})

	// 调用微信支付统一下单接口
	response, err := callWechatPayUnifiedOrder(request)
	if err != nil {
		LogError("调用微信支付接口失败", err)
		return nil, fmt.Errorf("调用微信支付接口失败: %v", err)
	}

	// 检查响应
	if response.ReturnCode != "SUCCESS" {
		LogError("微信支付返回错误", fmt.Errorf("return_code: %s, return_msg: %s", response.ReturnCode, response.ReturnMsg))
		return nil, fmt.Errorf("微信支付返回错误: %s", response.ReturnMsg)
	}

	if response.ResultCode != "SUCCESS" {
		LogError("微信支付业务失败", fmt.Errorf("result_code: %s", response.ResultCode))
		return nil, fmt.Errorf("微信支付业务失败")
	}

	// 生成小程序支付参数
	payParams := generateMiniProgramPayParams(response.PrepayID, wechatConfig.AppID, wechatConfig.MchKey)

	LogStep("微信支付参数生成成功", map[string]interface{}{
		"prepayID":  response.PrepayID,
		"timeStamp": payParams["timeStamp"],
	})

	return payParams, nil
}

// callWechatPayUnifiedOrder 调用微信支付统一下单接口
func callWechatPayUnifiedOrder(request *WechatPayRequest) (*WechatPayResponse, error) {
	// 确定API地址
	paymentConfig := config.GetPaymentConfig()
	apiURL := "https://api.mch.weixin.qq.com/pay/unifiedorder"
	if paymentConfig.WechatPay.Environment == "sandbox" {
		apiURL = "https://api.mch.weixin.qq.com/sandboxnew/pay/unifiedorder"
	}

	// 将请求转换为XML
	xmlData, err := xml.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("XML序列化失败: %v", err)
	}

	LogStep("发送微信支付请求", map[string]interface{}{
		"url":  apiURL,
		"data": string(xmlData),
	})

	// 发送HTTP请求
	resp, err := http.Post(apiURL, "application/xml", strings.NewReader(string(xmlData)))
	if err != nil {
		return nil, fmt.Errorf("HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	LogStep("收到微信支付响应", map[string]interface{}{
		"statusCode": resp.StatusCode,
		"response":   string(body),
	})

	// 解析响应
	var response WechatPayResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("XML解析失败: %v", err)
	}

	return &response, nil
}

// generateMiniProgramPayParams 生成小程序支付参数
func generateMiniProgramPayParams(prepayID, appID, mchKey string) map[string]interface{} {
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonceStr := generateNonceStr()
	packageStr := "prepay_id=" + prepayID
	signType := "MD5"

	// 构建签名字符串
	signParams := map[string]string{
		"appId":     appID,
		"timeStamp": timeStamp,
		"nonceStr":  nonceStr,
		"package":   packageStr,
		"signType":  signType,
	}

	// 生成签名
	paySign := generateWechatPaySign(signParams, mchKey)

	return map[string]interface{}{
		"timeStamp": timeStamp,
		"nonceStr":  nonceStr,
		"package":   packageStr,
		"signType":  signType,
		"paySign":   paySign,
	}
}

// generateWechatPaySign 生成微信支付签名
func generateWechatPaySign(params interface{}, mchKey string) string {
	var paramMap map[string]string

	// 根据参数类型处理
	switch v := params.(type) {
	case *WechatPayRequest:
		paramMap = map[string]string{
			"appid":            v.AppID,
			"mch_id":           v.MchID,
			"nonce_str":        v.NonceStr,
			"body":             v.Body,
			"out_trade_no":     v.OutTradeNo,
			"total_fee":        strconv.Itoa(v.TotalFee),
			"spbill_create_ip": v.SpbillCreateIP,
			"notify_url":       v.NotifyURL,
			"trade_type":       v.TradeType,
			"openid":           v.OpenID,
		}
	case map[string]string:
		paramMap = v
	default:
		return ""
	}

	// 过滤空值
	filteredParams := make(map[string]string)
	for k, v := range paramMap {
		if v != "" && k != "sign" {
			filteredParams[k] = v
		}
	}

	// 按键排序
	var keys []string
	for k := range filteredParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 构建签名字符串
	var signStr strings.Builder
	for _, k := range keys {
		if signStr.Len() > 0 {
			signStr.WriteString("&")
		}
		signStr.WriteString(k)
		signStr.WriteString("=")
		signStr.WriteString(filteredParams[k])
	}

	// 添加商户密钥
	signStr.WriteString("&key=")
	signStr.WriteString(mchKey)

	// 计算MD5
	hash := md5.Sum([]byte(signStr.String()))
	return strings.ToUpper(fmt.Sprintf("%x", hash))
}

// generateNonceStr 生成随机字符串
func generateNonceStr() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// HandleWechatPayNotify 处理微信支付通知
func HandleWechatPayNotify(w http.ResponseWriter, r *http.Request) {
	LogStep("收到微信支付通知", nil)

	// 读取请求体
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		LogError("读取支付通知失败", err)
		http.Error(w, "读取请求失败", http.StatusBadRequest)
		return
	}

	LogStep("支付通知内容", map[string]interface{}{
		"body": string(body),
	})

	// 解析XML
	var notifyData map[string]string
	if err := xml.Unmarshal(body, &notifyData); err != nil {
		LogError("解析支付通知XML失败", err)
		http.Error(w, "XML解析失败", http.StatusBadRequest)
		return
	}

	// 验证签名
	paymentConfig := config.GetPaymentConfig()
	expectedSign := generateWechatPaySign(notifyData, paymentConfig.WechatPay.MchKey)
	if expectedSign != notifyData["sign"] {
		LogError("支付通知签名验证失败", fmt.Errorf("expected: %s, actual: %s", expectedSign, notifyData["sign"]))
		http.Error(w, "签名验证失败", http.StatusBadRequest)
		return
	}

	// 检查支付结果
	if notifyData["return_code"] != "SUCCESS" || notifyData["result_code"] != "SUCCESS" {
		LogError("支付失败", fmt.Errorf("return_code: %s, result_code: %s", notifyData["return_code"], notifyData["result_code"]))
		// 返回成功响应给微信
		response := &WechatPayNotifyResponse{
			ReturnCode: "SUCCESS",
			ReturnMsg:  "OK",
		}
		w.Header().Set("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(response)
		return
	}

	// 处理支付成功
	orderNo := notifyData["out_trade_no"]
	transactionId := notifyData["transaction_id"]
	totalFee := notifyData["total_fee"]

	LogStep("支付成功", map[string]interface{}{
		"orderNo":       orderNo,
		"transactionId": transactionId,
		"totalFee":      totalFee,
	})

	// TODO: 更新订单状态
	// 这里应该调用订单服务更新订单状态为已支付

	// 返回成功响应给微信
	response := &WechatPayNotifyResponse{
		ReturnCode: "SUCCESS",
		ReturnMsg:  "OK",
	}
	w.Header().Set("Content-Type", "application/xml")
	xml.NewEncoder(w).Encode(response)

	LogStep("支付通知处理完成", nil)
}
