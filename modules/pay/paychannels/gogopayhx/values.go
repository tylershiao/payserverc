package gogopayhx

import "encoding/json"

const (
	baseURL          = "https://v202509.gogopay.live"
	payURL           = "/payin/unifiedorder"
	queryPayURL      = "/payin/query"
	withdrawURL      = "/payout/unifiedorder"
	queryWithdrawURL = "/payout/query"
	balURL           = "/payout/balance"
	mid              = "T813285283171891"
	key              = "9700889c9e10e52b7039ccf3ea16750b"
)

type PayReq struct {
	MerchantID     string `json:"merchantid"`       // 商户号
	OutTradeNo     string `json:"out_trade_no"`     // 商户订单号
	TotalFee       string `json:"total_fee"`        // 交易金额（保留两位小数）
	NotifyURL      string `json:"notify_url"`       // 异步通知地址
	ReplyType      string `json:"reply_type"`       // 执行方式：FORM 或 URL
	Timestamp      string `json:"timestamp"`        // 时间戳
	CustomerName   string `json:"customer_name"`    // 客户姓名（英文，去除空格）
	CustomerMobile string `json:"customer_mobile"`  // 客户手机号
	CustomerEmail  string `json:"customer_email"`   // 客户邮箱
	Sign           string `json:"sign" form:"sign"` // 签名
}

type PayResp struct {
	Code    int    `json:"code"`    // 返回码
	Message string `json:"message"` // 返回信息
	Data    struct {
		URL         string `json:"url"`          // 收银台跳转地址
		RealFee     string `json:"real_fee"`     // 实际支付金额
		OrderNumber string `json:"order_number"` // 平台订单号
	} `json:"data"` // 业务数据
	Errno int `json:"errno"` // 错误码
}

type PayCallbackReq struct {
	Status     string `json:"status"`       // 状态（payin_ing：交易中；payin_success：交易完成；payin_fail：交易失败）
	MerchantID string `json:"merchantid"`   // 商户号
	OutTradeNo string `json:"out_trade_no"` // 商户订单号
	TotalFee   string `json:"total_fee"`    // 交易金额（用户实际充值金额）
	Poundage   string `json:"poundage"`     // 交易手续费
	AccountFee string `json:"account_fee"`  // 到账金额
	Sign       string `json:"sign"`         // 签名
}

type QueryPayReq struct {
	MerchantID string `json:"merchantid" form:"merchantid"`     // 商户号
	OutTradeNo string `json:"out_trade_no" form:"out_trade_no"` // 商户订单号
	Timestamp  string `json:"timestamp" form:"timestamp"`       // 时间戳
	Sign       string `json:"sign" form:"sign"`                 // 签名
}

type QueryPayResp struct {
	Code    int    `json:"code"`    // 返回码 0 请求成功
	Message string `json:"message"` // 返回信息
	Data    struct {
		Status      string `json:"status"`       // 状态（payin_ing / payin_success / payin_fail）
		MerchantID  string `json:"merchantid"`   // 商户号
		OrderNumber string `json:"order_number"` // 平台订单号
		OutTradeNo  string `json:"out_trade_no"` // 商户订单号
		TotalFee    string `json:"total_fee"`    // 交易金额
		Poundage    string `json:"poundage"`     // 交易手续费
		AccountFee  string `json:"account_fee"`  // 到账金额
		UTR         string `json:"utr"`          // 支付凭证
		Sign        string `json:"sign"`         // 签名
	} `json:"data"` // 业务数据
}

type WithdrawReq struct {
	MerchantID     string `json:"merchantid" form:"merchantid"`           // 商户号
	OutTradeNo     string `json:"out_trade_no" form:"out_trade_no"`       // 商户订单号
	TotalFee       string `json:"total_fee" form:"total_fee"`             // 交易金额（保留两位小数）
	NotifyURL      string `json:"notify_url" form:"notify_url"`           // 异步通知地址
	Timestamp      string `json:"timestamp" form:"timestamp"`             // 时间戳
	PaymentMode    string `json:"payment_mode" form:"payment_mode"`       // 代付模式（IMPS）
	AccountNumber  string `json:"account_number" form:"account_number"`   // 账号
	IFSC           string `json:"ifsc" form:"ifsc"`                       // IFSC
	CustomerName   string `json:"customer_name" form:"customer_name"`     // 客户姓名（英文，去除空格）
	CustomerMobile string `json:"customer_mobile" form:"customer_mobile"` // 客户手机号
	CustomerEmail  string `json:"customer_email" form:"customer_email"`   // 客户邮箱
	Sign           string `json:"sign" form:"sign"`                       // 签名
}

type WithdrawResp struct {
	Code    int             `json:"code"`    // 0 请求成功
	Message string          `json:"message"` // 返回信息
	Data    json.RawMessage `json:"data"`    // 业务数据
}

type WidhrawRespData struct {
	Result      string `json:"result"`      // 请求结果（payout_ing 等）
	OrderNumber string `json:"orderNumber"` // 系统订单号
}

type WithdrawCallbackReq struct {
	OrderNumber string `json:"order_number" form:"order_number"` // 系统订单号
	Status      string `json:"status" form:"status"`             // 状态（payout_ing / payout_success / payout_fail）
	MerchantID  string `json:"merchantid" form:"merchantid"`     // 商户号
	OutTradeNo  string `json:"out_trade_no" form:"out_trade_no"` // 商户订单号
	TotalFee    string `json:"total_fee" form:"total_fee"`       // 交易金额
	Poundage    string `json:"poundage" form:"poundage"`         // 交易手续费
	AccountFee  string `json:"account_fee" form:"account_fee"`   // 到账金额
	FailInfo    string `json:"fail_info" form:"fail_info"`       // 失败信息
	UTR         string `json:"utr" form:"utr"`                   // UTR 回调凭证
	Sign        string `json:"sign" form:"sign"`                 // 签名
}

type QueryWithdrawReq struct {
	MerchantID string `json:"merchantid" form:"merchantid"`     // 商户号
	OutTradeNo string `json:"out_trade_no" form:"out_trade_no"` // 商户订单号
	Timestamp  string `json:"timestamp" form:"timestamp"`       // 时间戳
	Sign       string `json:"sign" form:"sign"`                 // 签名
}

type QueryWithdrawResp struct {
	Code    int    `json:"code"`    // 返回码 0 请求成功
	Message string `json:"message"` // 返回信息
	Data    struct {
		OrderNumber string `json:"order_number"` // 系统订单号
		Status      string `json:"status"`       // 状态（payout_ing / payout_success / payout_fail）
		MerchantID  string `json:"merchantid"`   // 商户号
		OutTradeNo  string `json:"out_trade_no"` // 商户订单号
		TotalFee    string `json:"total_fee"`    // 交易金额
		Poundage    string `json:"poundage"`     // 交易手续费
		AccountFee  string `json:"account_fee"`  // 到账金额
		FailInfo    string `json:"fail_info"`    // 失败信息
		UTR         string `json:"utr"`          // UTR 回调凭证
		Sign        string `json:"sign"`         // 签名
	} `json:"data"` // 业务数据
}
