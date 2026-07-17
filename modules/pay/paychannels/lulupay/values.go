package lulupay

const (
	baseURL          = "https://api.lupay.vip"
	payURL           = "/api/v3/deposits"
	queryPayURL      = "/api/v3/deposits/query"
	withdrawURL      = "/api/v3/transfers"
	queryWithdrawURL = "/api/v3/transfers/query"
	// balURL           = "/api/v3/balance"
	mid    = 13
	key    = "356ce124bd6b78c7"
	secret = "3fcbe3fdd41e7214b27b1e6fde6db826"
)

type PayReq struct {
	MID       int    `json:"mid"`        // 商户ID
	Amount    string `json:"amount"`     // 订单金额（元，最多2位小数）
	OrderNo   string `json:"order_no"`   // 商户订单号
	Gateway   string `json:"gateway"`    // 通道代码
	IP        string `json:"ip"`         // 客户 IP（支持 IPv6）
	NotifyURL string `json:"notify_url"` // 异步通知地址
	// ReturnURL string `json:"return_url"` // 支付成功跳转地址
	// Name      string `json:"name"`       // 付款人姓名
	// BankName  string `json:"bank_name"`  // 汇款银行
	// CardNo    string `json:"card_no"`    // 汇款银行账号
	// CardName  string `json:"card_name"`  // 汇款人实名
	Sign string `json:"sign"` // 签名
}

type PayResp struct {
	Code    int    `json:"code"`    // 返回状态码（200 成功）
	Message string `json:"message"` // 返回信息
	Data    struct {
		URL        string `json:"url"`        // 付款链接
		No         string `json:"no"`         // 系统订单号
		OrderNo    string `json:"order_no"`   // 商户订单号
		BankCode   string `json:"bankCode"`   // 银行代码
		BankName   string `json:"bankName"`   // 银行名称
		BankBranch string `json:"bankBranch"` // 银行支行
		CardNo     string `json:"cardNo"`     // 银行卡号 / 支付宝 / 微信 / 数字人民币账号
		CardName   string `json:"cardName"`   // 持卡人姓名
		QRCodeURL  string `json:"qrCodeUrl"`  // 二维码链接
		PayAmount  string `json:"pay_amount"` // 付款金额
		AppURL     string `json:"appUrl"`     // App 链接
	} `json:"data"` // 请求成功时返回
}

type PayCallbackReq struct {
	Code    int    `json:"code"`    // 200 成功，-9999 失败
	Message string `json:"message"` // 返回信息
	UTR     string `json:"utr"`     // UTR（可选）
	Data    struct {
		MID          int    `json:"mid"`           // 商户ID
		No           string `json:"no"`            // 系统订单号
		OrderNo      string `json:"order_no"`      // 商户订单号
		Amount       string `json:"amount"`        // 订单金额（元，最多2位小数）
		ActualAmount string `json:"actual_amount"` // 实际付款金额
		Fee          string `json:"fee"`           // 订单手续费
		// CreatedTime  int64  `json:"created_time"`  // 订单创建时间
		// DepositTime  int64  `json:"deposit_time"`  // 订单付款时间
		// NotifyTime   int64  `json:"notify_time"`   // 订单回调时间
		Status     string `json:"status"`      // inprogress / succeeded / failed
		Extra      string `json:"extra"`       // 穿透参数
		OrigAmount string `json:"orig_amount"` // 原始订单金额
		PayerName  string `json:"payer_name"`  // 付款人姓名
		Sign       string `json:"sign"`        // 签名
	} `json:"data"` // 请求成功时返回
}

type QueryPayReq struct {
	MID     int    `json:"mid"`      // 商户ID
	OrderNo string `json:"order_no"` // 商户订单号
	Sign    string `json:"sign"`     // 签名
}

type QueryPayResp struct {
	Code    int    `json:"code"`    // 返回状态码（200 成功）
	Message string `json:"message"` // 返回信息
	Data    struct {
		MID          string `json:"mid"`           // 商户ID
		No           string `json:"no"`            // 系统订单号
		OrderNo      string `json:"order_no"`      // 商户订单号
		Amount       string `json:"amount"`        // 订单金额（元，最多2位小数）
		ActualAmount string `json:"actual_amount"` // 实际付款金额
		Fee          string `json:"fee"`           // 订单手续费
		CreatedTime  string `json:"created_time"`  // 订单创建时间
		DepositTime  string `json:"deposit_time"`  // 订单付款时间
		NotifyTime   string `json:"notify_time"`   // 订单回调时间
		Status       string `json:"status"`        // inprogress / succeeded / failed
	} `json:"data"` // 请求成功时返回
}

type WithdrawReq struct {
	MID        int    `json:"mid"`         // 商户ID
	Amount     string `json:"amount"`      // 订单金额（元，最多2位小数）
	OrderNo    string `json:"order_no"`    // 商户订单号
	IP         string `json:"ip"`          // 客户IP（支持 IPv6）
	NotifyURL  string `json:"notify_url"`  // 异步通知地址
	BankCode   string `json:"bank_code"`   // 银行代码
	CardNo     string `json:"card_no"`     // 银行卡号 / 支付账号
	HolderName string `json:"holder_name"` // 持卡人姓名
	BankName   string `json:"bank_name"`   // 银行名称
	BankBranch string `json:"bank_branch"` // 银行支行
	// IdentityNo string `json:"identity_no"` // 身份证号
	// WithdrawQueryURL string `json:"withdrawQueryUrl"` // 反查链接（可选）
	Sign string `json:"sign"` // 签名
}

type WithdrawResp struct {
	Code    int    `json:"code"`    // 返回状态码（200 成功）
	Message string `json:"message"` // 返回信息
	Data    struct {
		No      string `json:"no"`       // 系统订单号
		OrderNo string `json:"order_no"` // 商户订单号
	} `json:"data"` // 请求成功时返回
}

type WithdrawCallbackReq struct {
	Code       int    `json:"code"`        // 200 成功，-9999 失败
	Message    string `json:"message"`     // 返回信息
	FailReason string `json:"fail_reason"` // 代付失败原因
	UTR        string `json:"utr"`         // UTR（可选）
	Data       struct {
		MID          int    `json:"mid"`           // 商户ID
		No           string `json:"no"`            // 系统订单号
		OrderNo      string `json:"order_no"`      // 商户订单号
		Amount       string `json:"amount"`        // 订单金额（元，最多2位小数）
		ActualAmount string `json:"actual_amount"` // 实际付款金额
		Fee          string `json:"fee"`           // 订单手续费
		// CreatedTime  int64  `json:"created_time"`  // 订单创建时间
		// TransferTime int64  `json:"transfer_time"` // 订单付款时间
		// NotifyTime   int64  `json:"notify_time"`   // 订单回调时间
		Status     string `json:"status"`       // inprogress / succeeded / failed
		Extra      string `json:"extra"`        // 穿透参数
		FromCardNo string `json:"from_card_no"` // 转账来源
		Sign       string `json:"sign"`         // 签名
	} `json:"data"` // 请求成功时返回
}

type QueryWithdrawReq struct {
	MID     int    `json:"mid"`      // 商户ID
	OrderNo string `json:"order_no"` // 商户订单号
	Sign    string `json:"sign"`     // 签名
}

type QueryWithdrawResp struct {
	Code    int    `json:"code"`    // 返回状态码（200 成功）
	Message string `json:"message"` // 返回信息
	Data    struct {
		MID          string `json:"mid"`           // 商户ID
		No           string `json:"no"`            // 系统订单号
		OrderNo      string `json:"order_no"`      // 商户订单号
		Amount       string `json:"amount"`        // 订单金额（元，最多2位小数）
		ActualAmount string `json:"actual_amount"` // 实际付款金额
		Fee          string `json:"fee"`           // 订单手续费
		CreatedTime  string `json:"created_time"`  // 订单创建时间
		DepositTime  string `json:"deposit_time"`  // 订单付款时间
		NotifyTime   string `json:"notify_time"`   // 订单回调时间
		Status       string `json:"status"`        // inprogress / succeeded / failed
	} `json:"data"` // 请求成功时返回
}
