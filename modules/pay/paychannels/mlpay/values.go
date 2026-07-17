package mlpay

const (
	baseURL          = "https://api.mlpayment.cc"
	payURL           = "/open-api/trade/collection"
	queryPayURL      = "/open-api/trade/collection/status"
	withdrawURL      = "/open-api/trade/payout"
	upiwithdrawURL   = "/open-api/trade/upi/payout"
	queryWithdrawURL = "/open-api/trade/payout/status"
	balURL           = ""
	mid              = "slot1818"
	// appID = "S1773216155928"
	key  = "SIExMlqndCGxMjEpuxju"
	key2 = "PaBriofUwTgQbsVBaZfd"
)

type PayReq struct {
	AppNo       string `json:"appNo"`       // 商家唯一编码
	CallbackURL string `json:"callbackUrl"` // 收款回调地址
	OrderCode   string `json:"orderCode"`   // 商家订单编号（唯一）
	Amount      string `json:"amount"`      // 交易金额（卢比，保留2位小数）
	Name        string `json:"name"`        // 用户姓名
	// Remark      string `json:"remark"`      // 订单备注
	Sign string `json:"sign"` // 签名
}

type PayResp struct {
	Code int    `json:"code"` // 响应码（200成功）
	Msg  string `json:"msg"`  // 响应描述
	Data struct {
		MerchantOrderCode string `json:"merchantOrderCode"` // 商户订单编号
		OrderCode         string `json:"orderCode"`         // 平台订单编号
		TradeURL          string `json:"tradeUrl"`          // 支付收银台地址
	} `json:"data"`
}

/*
0: 订单已创建，等待支付;
2: 失败;
1: 付款成功
*/
type PayCallbackReq struct {
	Amount       string `json:"amount"`       // 订单金额
	AppNo        string `json:"appNo"`        // 商家编码
	MerchantCode string `json:"merchantCode"` // 商家订单号
	RealAmount   string `json:"realAmount"`   // 实际收款金额
	Sign         string `json:"sign"`         // 签名
	Status       int    `json:"status"`       // 订单状态
	OrderCode    string `json:"orderCode"`    // 平台订单号
}

type QueryPayReq struct {
	AppNo     string `json:"appNo"`     // 商家唯一编码
	OrderCode string `json:"orderCode"` // 商家订单号
	Sign      string `json:"sign"`      // 签名
}

type QueryPayResp struct {
	Code int    `json:"code"` // 响应码：1成功
	Msg  string `json:"msg"`  // 描述信息
	Data struct {
		Amount       string `json:"amount"`       // 订单金额
		AppNo        string `json:"appNo"`        // 商家编码
		MerchantCode string `json:"merchantCode"` // 商家订单号
		RealAmount   string `json:"realAmount"`   // 实际收款金额
		Sign         string `json:"sign"`         // 签名
		Status       int    `json:"status"`       // 订单状态
		OrderCode    string `json:"orderCode"`    // 平台订单号
	} `json:"data"` // 数据体
}

type WithdrawReq struct {
	Amount        float64 `json:"amount"`        // 订单金额（保留2位小数）
	AppNo         string  `json:"appNo"`         // 商户编码
	CallbackURL   string  `json:"callbackUrl"`   // 异步回调地址
	AccountName   string  `json:"accountName"`   // 收款人姓名
	AccountNumber string  `json:"accountNumber"` // 收款银行卡号
	IFSC          string  `json:"ifsc"`          // 印度IFSC（11位）
	UPI           string  `json:"upi"`           // 印度UPI（11位）
	OrderCode     string  `json:"orderCode"`     // 商户订单号
	// Remark        string  `json:"remark"`        // 备注
	Sign string `json:"sign"` // 签名
}

type WithdrawResp struct {
	Code int    `json:"code"` // 响应码（200成功）
	Msg  string `json:"msg"`  // 描述信息
	Data struct {
		MerchantOrderCode string `json:"merchantOrderCode"` // 商户订单号
		OrderCode         string `json:"orderCode"`         // 平台订单号
	} `json:"data"`
}

/*
0:付款中；
1: 付款成功；
2: 失败;
3: 反转;
*/
type WithdrawCallbackReq struct {
	AppNo        string `json:"appNo"`        // 商家编码
	MerchantCode string `json:"merchantCode"` // 商户订单号
	OrderCode    string `json:"orderCode"`    // 平台订单号
	UTR          string `json:"utr"`          // 银行UTR
	Status       int    `json:"status"`       // 订单状态
	Sign         string `json:"sign"`         // 签名
}

type QueryWithdrawReq struct {
	AppNo     string `json:"appNo"`     // 商家唯一编码
	OrderCode string `json:"orderCode"` // 商家订单号
	Sign      string `json:"sign"`      // 签名
}

type QueryWithdrawResp struct {
	Code int    `json:"code"` // 响应码：1成功
	Msg  string `json:"msg"`  // 描述信息
	Data struct {
		AppNo        string `json:"appNo"`        // 商家编码
		MerchantCode string `json:"merchantCode"` // 商户订单号
		OrderCode    string `json:"orderCode"`    // 平台订单号
		Sign         string `json:"sign"`         // 签名
		Status       int    `json:"status"`       // 订单状态
	} `json:"data"` // 数据体
}
