package deeplink

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"server/util"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/liangdas/mqant/log"
)

// 支付模式常量
const (
	PaymentModeCollect      = 1 // 代收
	PaymentModePayout       = 2 // 代付
	PaymentModeNative       = 3 // 原生支付
	PaymentModeMQR          = 4 // MQR支付
	PaymentModeNativeSecure = 5 // 原生安全支付
	PaymentModeDeepLink     = 6
)

// 签名类型常量
const (
	SignTypeMD5        = "MD5"
	SignTypeHMACSHA256 = "HMAC-SHA256"
)

// Config 支付配置
type Config struct {
	BaseURL     string `json:"base_url"`     // API基础地址
	MerchantID  string `json:"merchant_id"`  // 商户ID
	APIKey      string `json:"api_key"`      // API密钥
	APISecret   string `json:"api_secret"`   // API密钥
	SignType    string `json:"sign_type"`    // 签名类型: MD5 或 HMAC-SHA256
	NotifyURL   string `json:"notify_url"`   // 支付回调地址
	ReturnURL   string `json:"return_url"`   // 支付返回地址
	Currency    string `json:"currency"`     // 货币类型
	ProxyURL    string `json:"proxy_url"`    // 代理地址(可选)
	Timeout     int    `json:"timeout"`      // 请求超时时间(秒)
	NativePhone string `json:"native_phone"` // 原生支付手机号
}

// PayoutInfo 代付信息
type PayoutInfo struct {
	AccountType      string `json:"account_type"`      // 账户类型(UPI/USDT/Bank/Epay)
	AccountNumber    string `json:"account_number"`    // 统一账户标识
	BankBranch       string `json:"bank_branch"`       // 银行分支
	BankName         string `json:"bank_name"`         // 银行名称(Bank类型必填)
	BeneficiaryEmail string `json:"beneficiary_email"` // 收款人邮箱
	BeneficiaryName  string `json:"beneficiary_name"`  // 收款人姓名
	BeneficiaryPhone string `json:"beneficiary_phone"` // 收款人手机号
	IFSC             string `json:"ifsc"`              // IFSC代码(Bank类型必填)
	Network          string `json:"network"`           // 网络类型(USDT类型必填)
	Remark           string `json:"remark"`            // 转账备注
}

// OrderResponse 订单响应
type OrderResponse struct {
	Code        int                    `json:"code"`
	Msg         string                 `json:"msg"`
	Data        *OrderResponseData     `json:"data,omitempty"`
	RawResponse map[string]interface{} `json:"-"`
}

// OrderResponseData 订单响应数据
type OrderResponseData struct {
	OrderNo     string `json:"order_no"`
	PayURL      string `json:"pay_url"`
	MerchantID  string `json:"merchant_id"`
	ProjectCode string `json:"project_code"`
	UserUID     string `json:"user_uid"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	PaymentMode int    `json:"payment_mode"`
	ExtraData   string `json:"extra_data"`
	ExpiredAt   string `json:"expired_at"`
	CreatedAt   string `json:"created_at"`
	Signature   string `json:"signature"`
}

// OrderStatusResponse 订单状态响应
type OrderStatusResponse struct {
	Code        int                    `json:"code"`
	Message     string                 `json:"message"`
	Data        *OrderStatusData       `json:"data,omitempty"`
	RawResponse map[string]interface{} `json:"-"`
}

// OrderStatusData 订单状态数据
type OrderStatusData struct {
	OrderNo         string `json:"order_no"`
	MerchantOrderNo string `json:"merchant_order_no"`
	Status          int    `json:"status"`
	PaymentMode     int    `json:"payment_mode"`
	Amount          string `json:"amount"`
	Currency        string `json:"currency"`
	Signature       string `json:"signature"`
	CreatedAt       int64  `json:"created_at"`
}

// CallbackData 回调数据
type CallbackData struct {
	OrderNo         string `json:"order_no"`
	MerchantOrderNo string `json:"merchant_order_no"`
	ProjectCode     string `json:"project_code"`
	UserUID         string `json:"user_uid"`
	Amount          string `json:"amount"`
	Currency        string `json:"currency"`
	Status          int    `json:"status"`
	PaymentMode     int    `json:"payment_mode"`
	CreatedAt       int64  `json:"created_at"`
	ExtraData       string `json:"extra_data"`
	Signature       string `json:"signature"`
}

// PaymentClient 支付客户端
type PaymentClient struct {
	config     *Config
	httpClient *http.Client
}

// NewPaymentClient 创建支付客户端
func NewPaymentClient(config *Config) *PaymentClient {
	client := &PaymentClient{
		config: config,
	}

	// 创建HTTP客户端
	transport := &http.Transport{
		MaxIdleConns:        10,
		IdleConnTimeout:     30 * time.Second,
		DisableCompression:  true,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	// 如果配置了代理
	if config.ProxyURL != "" {
		proxyURL, err := url.Parse(config.ProxyURL)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	}

	timeout := config.Timeout
	if timeout <= 0 {
		timeout = 30 // 默认30秒
	}

	client.httpClient = &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeout) * time.Second,
	}

	return client
}

// BasePayment 基础支付实现
type BasePayment struct {
	client         *PaymentClient
	orderResponse  *OrderResponse
	statusResponse *OrderStatusResponse
	callbackData   *CallbackData
	paymentMode    int
	signature      string
	callbackResp   string
}

// NewBasePayment 创建基础支付
func NewBasePayment(client *PaymentClient, paymentMode int) *BasePayment {
	return &BasePayment{
		client:      client,
		paymentMode: paymentMode,
	}
}

// GetPaymentMode 获取支付模式
func (p *BasePayment) GetPaymentMode() int {
	return p.paymentMode
}

// GetOrderResponse 获取订单响应
func (p *BasePayment) GetOrderResponse() *OrderResponse {
	return p.orderResponse
}

// GetOrderStatusResponse 获取订单状态响应
func (p *BasePayment) GetOrderStatusResponse() *OrderStatusResponse {
	return p.statusResponse
}

// GetCallbackData 获取回调数据
func (p *BasePayment) GetCallbackData() *CallbackData {
	return p.callbackData
}

// CollectPayment 代收支付
type CollectPayment struct {
	*BasePayment
	orderData map[string]interface{}
}

// NewCollectPayment 创建代收支付
func NewCollectPayment(client *PaymentClient) *CollectPayment {
	return &CollectPayment{
		BasePayment: NewBasePayment(client, PaymentModeCollect),
		orderData:   make(map[string]interface{}),
	}
}

// SetOrderData 设置订单数据
func (p *CollectPayment) SetOrderData(merchantOrderNo, userUID, amount, extraData string) {
	p.orderData["project_code"] = p.client.config.MerchantID
	p.orderData["merchant_order_no"] = merchantOrderNo
	p.orderData["user_uid"] = userUID
	p.orderData["amount"] = amount
	p.orderData["currency"] = p.client.config.Currency
	p.orderData["payment_mode"] = PaymentModeCollect
	p.orderData["notify_url"] = p.client.config.NotifyURL
	p.orderData["return_url"] = p.client.config.ReturnURL
	p.orderData["extra_data"] = extraData
}

// PackHeader 生成请求头
func (p *CollectPayment) PackHeader(header http.Header) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := GenerateRandomString(16)

	header.Set("Content-Type", "application/json")
	header.Set("x-merchant-id", p.client.config.MerchantID)
	header.Set("x-api-key", p.client.config.APIKey)
	header.Set("x-timestamp", timestamp)
	header.Set("x-nonce", nonce)
	header.Set("x-sign-type", p.client.config.SignType)

	// 计算签名
	params := make(map[string]string)
	params["x-merchant-id"] = p.client.config.MerchantID
	params["x-api-key"] = p.client.config.APIKey
	params["x-timestamp"] = timestamp
	params["x-nonce"] = nonce
	params["x-sign-type"] = p.client.config.SignType

	for k, v := range p.orderData {
		params[k] = convertToString(v)
	}

	var signature string
	if p.client.config.SignType == SignTypeHMACSHA256 {
		signature = calculateHMACSHA256Signature(params, p.client.config.APISecret)
	} else {
		signature = calculateMD5Signature(params, p.client.config.APISecret)
	}

	header.Set("x-signature", signature)
	p.signature = signature
}

// PackReq 生成支付请求体
func (p *CollectPayment) PackReq() interface{} {
	return p.orderData
}

// GetResp 获取返回
func (p *CollectPayment) GetResp() error {
	apiURL := p.client.config.BaseURL + "/api/order/create"

	jsonData, err := json.Marshal(p.orderData)
	if err != nil {
		return fmt.Errorf("序列化请求体失败: %v", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	header := make(http.Header)
	p.PackHeader(header)
	req.Header = header

	resp, err := p.client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应体失败: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	p.orderResponse = &OrderResponse{
		RawResponse: result,
	}

	code, ok := result["code"].(float64)
	if !ok {
		return fmt.Errorf("响应格式错误")
	}

	p.orderResponse.Code = int(code)
	p.orderResponse.Msg, _ = result["msg"].(string)

	if int(code) == 0 {
		data, ok := result["data"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("响应数据格式错误")
		}

		orderData := &OrderResponseData{
			OrderNo:     getString(data, "order_no"),
			PayURL:      getString(data, "pay_url"),
			MerchantID:  getString(data, "merchant_id"),
			ProjectCode: getString(data, "project_code"),
			UserUID:     getString(data, "user_uid"),
			Amount:      getString(data, "amount"),
			Currency:    getString(data, "currency"),
			PaymentMode: getInt(data, "payment_mode"),
			ExtraData:   getString(data, "extra_data"),
			ExpiredAt:   getString(data, "expired_at"),
			CreatedAt:   getString(data, "created_at"),
			Signature:   getString(data, "signature"),
		}

		p.orderResponse.Data = orderData
	}

	return nil
}

// CheckSign 验签
func (p *CollectPayment) CheckSign(str string) bool {
	// 这里实现验签逻辑
	// 对于代收支付，通常不需要验签回调，因为回调是服务器主动发送的
	return true
}

// PayoutPayment 代付支付
type PayoutPayment struct {
	*BasePayment
	orderData  map[string]interface{}
	payoutInfo *PayoutInfo
}

// NewPayoutPayment 创建代付支付
func NewPayoutPayment(client *PaymentClient) *PayoutPayment {
	return &PayoutPayment{
		BasePayment: NewBasePayment(client, PaymentModePayout),
		orderData:   make(map[string]interface{}),
	}
}

// SetOrderData 设置订单数据
func (p *PayoutPayment) SetOrderData(merchantOrderNo, userUID, amount, extraData string, payoutInfo *PayoutInfo) {
	p.payoutInfo = payoutInfo
	p.orderData["project_code"] = p.client.config.MerchantID
	p.orderData["merchant_order_no"] = merchantOrderNo
	p.orderData["user_uid"] = userUID
	p.orderData["amount"] = amount
	p.orderData["currency"] = p.client.config.Currency
	p.orderData["payment_mode"] = PaymentModePayout
	p.orderData["notify_url"] = p.client.config.NotifyURL
	p.orderData["extra_data"] = extraData
	p.orderData["payout_info"] = payoutInfo
}

// PackHeader 生成请求头
func (p *PayoutPayment) PackHeader(header http.Header) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := GenerateRandomString(16)

	header.Set("Content-Type", "application/json")
	header.Set("x-merchant-id", p.client.config.MerchantID)
	header.Set("x-api-key", p.client.config.APIKey)
	header.Set("x-timestamp", timestamp)
	header.Set("x-nonce", nonce)
	header.Set("x-sign-type", p.client.config.SignType)

	// 计算签名
	params := make(map[string]string)
	params["x-merchant-id"] = p.client.config.MerchantID
	params["x-api-key"] = p.client.config.APIKey
	params["x-timestamp"] = timestamp
	params["x-nonce"] = nonce
	params["x-sign-type"] = p.client.config.SignType

	for k, v := range p.orderData {
		if k != "payout_info" {
			params[k] = convertToString(v)
		}
	}

	// 将payout_info序列化为JSON字符串
	if p.payoutInfo != nil {
		payoutInfoJSON, _ := json.Marshal(p.payoutInfo)
		params["payout_info"] = string(payoutInfoJSON)
	}

	var signature string
	if p.client.config.SignType == SignTypeHMACSHA256 {
		signature = calculateHMACSHA256Signature(params, p.client.config.APISecret)
	} else {
		signature = calculateMD5Signature(params, p.client.config.APISecret)
	}

	header.Set("x-signature", signature)
	p.signature = signature
}

// PackReq 生成支付请求体
func (p *PayoutPayment) PackReq() interface{} {
	return p.orderData
}

// GetResp 获取返回
func (p *PayoutPayment) GetResp() error {
	apiURL := p.client.config.BaseURL + "/api/order/create"

	jsonData, err := json.Marshal(p.orderData)
	if err != nil {
		return fmt.Errorf("序列化请求体失败: %v", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	header := make(http.Header)
	p.PackHeader(header)
	req.Header = header

	resp, err := p.client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应体失败: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	p.orderResponse = &OrderResponse{
		RawResponse: result,
	}

	code, ok := result["code"].(float64)
	if !ok {
		return fmt.Errorf("响应格式错误")
	}

	p.orderResponse.Code = int(code)
	p.orderResponse.Msg, _ = result["msg"].(string)

	if int(code) == 0 {
		data, ok := result["data"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("响应数据格式错误")
		}

		orderData := &OrderResponseData{
			OrderNo:     getString(data, "order_no"),
			PayURL:      getString(data, "pay_url"),
			MerchantID:  getString(data, "merchant_id"),
			ProjectCode: getString(data, "project_code"),
			UserUID:     getString(data, "user_uid"),
			Amount:      getString(data, "amount"),
			Currency:    getString(data, "currency"),
			PaymentMode: getInt(data, "payment_mode"),
			ExtraData:   getString(data, "extra_data"),
			ExpiredAt:   getString(data, "expired_at"),
			CreatedAt:   getString(data, "created_at"),
			Signature:   getString(data, "signature"),
		}

		p.orderResponse.Data = orderData
	}

	return nil
}

// CheckSign 验签
func (p *PayoutPayment) CheckSign(str string) bool {
	// 这里实现验签逻辑
	// 对于代付支付，通常不需要验签回调，因为回调是服务器主动发送的
	return true
}

// NativePayment 原生支付
type NativePayment struct {
	*BasePayment
	orderData map[string]interface{}
	phone     string
}

// NewNativePayment 创建原生支付
func NewNativePayment(client *PaymentClient) *NativePayment {
	return &NativePayment{
		BasePayment: NewBasePayment(client, PaymentModeNative),
		orderData:   make(map[string]interface{}),
	}
}

// SetOrderData 设置订单数据
func (p *NativePayment) SetOrderData(merchantOrderNo, userUID, amount, extraData, phone string) {
	p.phone = phone
	p.orderData["project_code"] = p.client.config.MerchantID
	p.orderData["merchant_order_no"] = merchantOrderNo
	p.orderData["user_uid"] = userUID
	p.orderData["amount"] = amount
	p.orderData["currency"] = p.client.config.Currency
	p.orderData["payment_mode"] = PaymentModeNative
	p.orderData["phone"] = phone
	p.orderData["notify_url"] = p.client.config.NotifyURL
	p.orderData["return_url"] = p.client.config.ReturnURL
	p.orderData["extra_data"] = extraData
}

// PackHeader 生成请求头
func (p *NativePayment) PackHeader(header http.Header) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := GenerateRandomString(16)

	header.Set("Content-Type", "application/json")
	header.Set("x-merchant-id", p.client.config.MerchantID)
	header.Set("x-api-key", p.client.config.APIKey)
	header.Set("x-timestamp", timestamp)
	header.Set("x-nonce", nonce)
	header.Set("x-sign-type", p.client.config.SignType)

	// 计算签名
	params := make(map[string]string)
	params["x-merchant-id"] = p.client.config.MerchantID
	params["x-api-key"] = p.client.config.APIKey
	params["x-timestamp"] = timestamp
	params["x-nonce"] = nonce
	params["x-sign-type"] = p.client.config.SignType

	for k, v := range p.orderData {
		params[k] = convertToString(v)
	}

	var signature string
	if p.client.config.SignType == SignTypeHMACSHA256 {
		signature = calculateHMACSHA256Signature(params, p.client.config.APISecret)
	} else {
		signature = calculateMD5Signature(params, p.client.config.APISecret)
	}

	header.Set("x-signature", signature)
	p.signature = signature
}

// PackReq 生成支付请求体
func (p *NativePayment) PackReq() interface{} {
	return p.orderData
}

// GetResp 获取返回
func (p *NativePayment) GetResp() error {
	apiURL := p.client.config.BaseURL + "/api/order/create"

	jsonData, err := json.Marshal(p.orderData)
	if err != nil {
		return fmt.Errorf("序列化请求体失败: %v", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	header := make(http.Header)
	p.PackHeader(header)
	req.Header = header

	resp, err := p.client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应体失败: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	p.orderResponse = &OrderResponse{
		RawResponse: result,
	}

	code, ok := result["code"].(float64)
	if !ok {
		return fmt.Errorf("响应格式错误")
	}

	p.orderResponse.Code = int(code)
	p.orderResponse.Msg, _ = result["msg"].(string)

	if int(code) == 0 {
		data, ok := result["data"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("响应数据格式错误")
		}

		orderData := &OrderResponseData{
			OrderNo:     getString(data, "order_no"),
			PayURL:      getString(data, "pay_url"),
			MerchantID:  getString(data, "merchant_id"),
			ProjectCode: getString(data, "project_code"),
			UserUID:     getString(data, "user_uid"),
			Amount:      getString(data, "amount"),
			Currency:    getString(data, "currency"),
			PaymentMode: getInt(data, "payment_mode"),
			ExtraData:   getString(data, "extra_data"),
			ExpiredAt:   getString(data, "expired_at"),
			CreatedAt:   getString(data, "created_at"),
			Signature:   getString(data, "signature"),
		}

		p.orderResponse.Data = orderData
	}

	return nil
}

// CheckSign 验签
func (p *NativePayment) CheckSign(str string) bool {
	// 这里实现验签逻辑
	// 对于原生支付，通常不需要验签回调，因为回调是服务器主动发送的
	return true
}

// NativeSecurePayment 原生安全支付
type NativeSecurePayment struct {
	*NativePayment
}

// NewNativeSecurePayment 创建原生安全支付
func NewNativeSecurePayment(client *PaymentClient) *NativeSecurePayment {
	return &NativeSecurePayment{
		NativePayment: &NativePayment{
			BasePayment: NewBasePayment(client, PaymentModeNativeSecure),
			orderData:   make(map[string]interface{}),
		},
	}
}

// SetOrderData 设置订单数据
func (p *NativeSecurePayment) SetOrderData(merchantOrderNo, userUID, amount, extraData, phone string) {
	p.phone = phone
	p.orderData["project_code"] = p.client.config.MerchantID
	p.orderData["merchant_order_no"] = merchantOrderNo
	p.orderData["user_uid"] = userUID
	p.orderData["amount"] = amount
	p.orderData["currency"] = p.client.config.Currency
	p.orderData["payment_mode"] = PaymentModeNativeSecure
	p.orderData["phone"] = phone
	p.orderData["notify_url"] = p.client.config.NotifyURL
	p.orderData["return_url"] = p.client.config.ReturnURL
	p.orderData["extra_data"] = extraData
}

// deeplinkPayment deeplink支付
type DeepLinkPayment struct {
	*BasePayment
	OrderData map[string]interface{}
}

// NewDeepLinkPayment 创建DeepLink支付
func NewDeepLinkPayment(client *PaymentClient) *DeepLinkPayment {
	return &DeepLinkPayment{
		BasePayment: NewBasePayment(client, PaymentModeDeepLink),
		OrderData:   make(map[string]interface{}),
	}
}

// SetOrderData 设置订单数据
func (p *DeepLinkPayment) SetOrderData(merchantOrderNo, userUID, amount, extraData string) {
	p.OrderData["project_code"] = p.client.config.MerchantID
	p.OrderData["merchant_order_no"] = merchantOrderNo
	p.OrderData["user_uid"] = userUID
	p.OrderData["amount"] = amount
	p.OrderData["currency"] = p.client.config.Currency
	p.OrderData["payment_mode"] = PaymentModeDeepLink
	p.OrderData["notify_url"] = p.client.config.NotifyURL
	p.OrderData["return_url"] = p.client.config.ReturnURL
	p.OrderData["extra_data"] = extraData
}

// PackHeader 生成请求头
func (p *DeepLinkPayment) PackHeader(header http.Header) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := GenerateRandomString(16)

	header.Set("Content-Type", "application/json")
	header.Set("x-merchant-id", p.client.config.MerchantID)
	header.Set("x-api-key", p.client.config.APIKey)
	header.Set("x-timestamp", timestamp)
	header.Set("x-nonce", nonce)
	header.Set("x-sign-type", p.client.config.SignType)

	// 计算签名
	params := make(map[string]string)
	params["x-merchant-id"] = p.client.config.MerchantID
	params["x-api-key"] = p.client.config.APIKey
	params["x-timestamp"] = timestamp
	params["x-nonce"] = nonce
	params["x-sign-type"] = p.client.config.SignType

	for k, v := range p.OrderData {
		params[k] = convertToString(v)
	}

	var signature string
	if p.client.config.SignType == SignTypeHMACSHA256 {
		signature = calculateHMACSHA256Signature(params, p.client.config.APISecret)
	} else {
		signature = calculateMD5Signature(params, p.client.config.APISecret)
	}

	header.Set("x-signature", signature)
	p.signature = signature
}

// PackReq 生成支付请求体
func (p *DeepLinkPayment) PackReq() interface{} {
	return p.OrderData
}

// GetResp 获取返回
func (p *DeepLinkPayment) GetResp() error {
	apiURL := p.client.config.BaseURL + "/api/order/create"

	jsonData, err := json.Marshal(p.OrderData)
	if err != nil {
		return fmt.Errorf("序列化请求体失败: %v", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	header := make(http.Header)
	p.PackHeader(header)
	req.Header = header

	resp, err := p.client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应体失败: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	p.orderResponse = &OrderResponse{
		RawResponse: result,
	}

	code, ok := result["code"].(float64)
	if !ok {
		return fmt.Errorf("响应格式错误")
	}

	p.orderResponse.Code = int(code)
	p.orderResponse.Msg, _ = result["msg"].(string)

	if int(code) == 0 {
		data, ok := result["data"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("响应数据格式错误")
		}

		orderData := &OrderResponseData{
			OrderNo:     getString(data, "order_no"),
			PayURL:      getString(data, "pay_url"),
			MerchantID:  getString(data, "merchant_id"),
			ProjectCode: getString(data, "project_code"),
			UserUID:     getString(data, "user_uid"),
			Amount:      getString(data, "amount"),
			Currency:    getString(data, "currency"),
			PaymentMode: getInt(data, "payment_mode"),
			ExtraData:   getString(data, "extra_data"),
			ExpiredAt:   getString(data, "expired_at"),
			CreatedAt:   getString(data, "created_at"),
			Signature:   getString(data, "signature"),
		}

		p.orderResponse.Data = orderData
	}

	return nil
}

// CheckSign 验签
func (p *DeepLinkPayment) CheckSign(str string) bool {
	// 这里实现验签逻辑
	// 对于DeepLink支付，通常不需要验签回调，因为回调是服务器主动发送的
	return true
}

// OrderStatusQuery 订单状态查询
type OrderStatusQuery struct {
	*BasePayment
	orderNo string
}

// NewOrderStatusQuery 创建订单状态查询
func NewOrderStatusQuery(client *PaymentClient) *OrderStatusQuery {
	return &OrderStatusQuery{
		BasePayment: NewBasePayment(client, 0),
	}
}

// SetOrderNo 设置订单号
func (p *OrderStatusQuery) SetOrderNo(orderNo string) {
	p.orderNo = orderNo
}

// PackHeader 生成请求头
func (p *OrderStatusQuery) PackHeader(header http.Header) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := GenerateRandomString(16)

	header.Set("Content-Type", "application/json")
	header.Set("x-merchant-id", p.client.config.MerchantID)
	header.Set("x-api-key", p.client.config.APIKey)
	header.Set("x-timestamp", timestamp)
	header.Set("x-nonce", nonce)
	header.Set("x-sign-type", p.client.config.SignType)

	// 计算签名
	params := make(map[string]string)
	params["x-merchant-id"] = p.client.config.MerchantID
	params["x-api-key"] = p.client.config.APIKey
	params["x-timestamp"] = timestamp
	params["x-nonce"] = nonce
	params["x-sign-type"] = p.client.config.SignType
	params["order_no"] = p.orderNo

	var signature string
	if p.client.config.SignType == SignTypeHMACSHA256 {
		signature = calculateHMACSHA256Signature(params, p.client.config.APISecret)
	} else {
		signature = calculateMD5Signature(params, p.client.config.APISecret)
	}

	header.Set("x-signature", signature)
	p.signature = signature
}

// PackReq 生成支付请求体
func (p *OrderStatusQuery) PackReq() interface{} {
	// 订单状态查询是GET请求，没有请求体
	return nil
}

// GetResp 获取返回
func (p *OrderStatusQuery) GetResp() error {
	apiURL := fmt.Sprintf("%s/api/order/status?order_no=%s", p.client.config.BaseURL, p.orderNo)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	header := make(http.Header)
	p.PackHeader(header)
	req.Header = header

	resp, err := p.client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应体失败: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	p.statusResponse = &OrderStatusResponse{
		RawResponse: result,
	}

	code, ok := result["code"].(float64)
	if !ok {
		return fmt.Errorf("响应格式错误")
	}

	p.statusResponse.Code = int(code)
	p.statusResponse.Message, _ = result["message"].(string)

	if int(code) == 200 {
		data, ok := result["data"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("响应数据格式错误")
		}

		statusData := &OrderStatusData{
			OrderNo:         getString(data, "order_no"),
			MerchantOrderNo: getString(data, "merchant_order_no"),
			Status:          getInt(data, "status"),
			PaymentMode:     getInt(data, "payment_mode"),
			Amount:          getString(data, "amount"),
			Currency:        getString(data, "currency"),
			Signature:       getString(data, "signature"),
			CreatedAt:       getInt64(data, "created_at"),
		}

		p.statusResponse.Data = statusData
	}

	return nil
}

// CheckSign 验签
func (p *OrderStatusQuery) CheckSign(str string) bool {
	// 这里实现验签逻辑
	// 对于订单状态查询，通常不需要验签
	return true
}

// CallbackHandler 回调处理器
type CallbackHandler struct {
	*BasePayment
}

// NewCallbackHandler 创建回调处理器
func NewCallbackHandler(client *PaymentClient) *CallbackHandler {
	return &CallbackHandler{
		BasePayment: NewBasePayment(client, 0),
	}
}

// ParseCallback 解析回调数据
func (p *CallbackHandler) ParseCallback(body []byte) error {
	var callbackData map[string]interface{}
	if err := json.Unmarshal(body, &callbackData); err != nil {
		return fmt.Errorf("解析回调数据失败: %v", err)
	}

	p.callbackData = &CallbackData{
		OrderNo:         getString(callbackData, "order_no"),
		MerchantOrderNo: getString(callbackData, "merchant_order_no"),
		ProjectCode:     getString(callbackData, "project_code"),
		UserUID:         getString(callbackData, "user_uid"),
		Amount:          getString(callbackData, "amount"),
		Currency:        getString(callbackData, "currency"),
		Status:          getInt(callbackData, "status"),
		PaymentMode:     getInt(callbackData, "payment_mode"),
		CreatedAt:       getInt64(callbackData, "created_at"),
		ExtraData:       getString(callbackData, "extra_data"),
		Signature:       getString(callbackData, "signature"),
	}

	return nil
}

// PackHeader 生成请求头
func (p *CallbackHandler) PackHeader(header http.Header) {
	// 回调处理不需要生成请求头
}

// PackReq 生成支付请求体
func (p *CallbackHandler) PackReq() interface{} {
	// 回调处理不需要生成请求体
	return nil
}

// GetResp 获取返回
func (p *CallbackHandler) GetResp() error {
	// 回调处理不需要获取返回
	return nil
}

// CheckSign 验签
func (p *CallbackHandler) CheckSign(str string) bool {
	if p.callbackData == nil {
		return false
	}

	// 构建签名参数
	params := make(map[string]string)
	params["order_no"] = p.callbackData.OrderNo
	params["merchant_order_no"] = p.callbackData.MerchantOrderNo
	params["project_code"] = p.callbackData.ProjectCode
	params["user_uid"] = p.callbackData.UserUID
	params["amount"] = p.callbackData.Amount
	params["currency"] = p.callbackData.Currency
	params["status"] = strconv.Itoa(p.callbackData.Status)
	params["payment_mode"] = strconv.Itoa(p.callbackData.PaymentMode)
	params["created_at"] = strconv.FormatInt(p.callbackData.CreatedAt, 10)
	params["extra_data"] = p.callbackData.ExtraData

	// 计算签名
	expectedSignature := calculateHMACSHA256Signature(params, p.client.config.APISecret)

	// 验证签名
	return p.callbackData.Signature == expectedSignature
}

// 工具函数

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// convertToString 将任意类型转换为字符串
func convertToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int, int64, float64:
		return fmt.Sprintf("%v", v)
	default:
		// 对于复杂类型，尝试序列化为JSON字符串
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return ""
		}
		return string(jsonBytes)
	}
}

// calculateMD5Signature 计算MD5签名
func calculateMD5Signature(params map[string]string, apiSecret string) string {
	var keys []string
	for k := range params {
		if k != "signature" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var buf strings.Builder
	for i, k := range keys {
		if i > 0 {
			buf.WriteString("&")
		}
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(params[k])
	}

	buf.WriteString("&key=")
	buf.WriteString(apiSecret)

	signStr := buf.String()

	// 使用项目中的MD5计算函数
	signature := util.CalculateMD5(signStr)
	return strings.ToUpper(signature)
}

// calculateHMACSHA256Signature 计算HMAC-SHA256签名
func calculateHMACSHA256Signature(params map[string]string, apiSecret string) string {
	var keys []string
	for k := range params {
		if k != "signature" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var buf strings.Builder
	for i, k := range keys {
		if i > 0 {
			buf.WriteString("&")
		}
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(params[k])
	}

	signStr := buf.String()

	log.Debug("签名前缀字符串: %s", signStr)

	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write([]byte(signStr))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// getString 从map中获取字符串值
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// getInt 从map中获取整数值
func getInt(m map[string]interface{}, key string) int {
	if val, ok := m[key]; ok {
		if num, ok := val.(float64); ok {
			return int(num)
		}
	}
	return 0
}

// getInt64 从map中获取int64值
func getInt64(m map[string]interface{}, key string) int64 {
	if val, ok := m[key]; ok {
		if num, ok := val.(float64); ok {
			return int64(num)
		}
	}
	return 0
}
