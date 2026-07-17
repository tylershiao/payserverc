package decentpay

const (
	baseURL          = "https://api-in.paydecent.com/pay"
	payURL           = "/collection/direct-payment/create"
	queryPayURL      = "/collection/direct-payment/"
	withdrawURL      = "/payment/disbursement/create"
	queryWithdrawURL = "/payment/disbursement/"
	balURL           = "/tenant/query/balance"
	mid              = "1177"
	appID            = "117701"
	key              = "1a92r20Mo1JUI0CELBFFIcUvSXqLKhxdvkqBCPstwsEnITVqxqGcV24V8NXnSTIiNSMIjalHb2n/VE0DQHwL0Q=="
)

type PayReq struct {
	Amount          float64 `json:"amount"`          // 金额
	NotifyURL       string  `json:"notifyUrl"`       // 回调地址（可选）
	MerchantOrderNo string  `json:"merchantOrderNo"` // 商户订单号
	PayProduct      string  `json:"payProduct"`      // 支付产品 (默认 india-upi-h5)
	Description     string  `json:"description"`     // 商品描述，可选
	UserID          string  `json:"userId"`          // 用户ID
	FullName        string  `json:"fullName"`        // 用户姓名
	DateOfReg       string  `json:"dateOfReg"`       // 注册日期 yyyy-MM-dd HH:mm:ss
}

type PayResp struct {
	Success bool   `json:"success"` // 是否成功
	ErrCode string `json:"errCode"` // 错误码
	ErrMsg  string `json:"errMsg"`  // 错误信息
	Data    struct {
		MerchantOrderNo string  `json:"merchantOrderNo"` // 商户订单号
		OrderNo         string  `json:"orderNo"`         // 订单号
		Amount          float64 `json:"amount"`          // 金额
		Fee             float64 `json:"fee"`             // 手续费
		Status          string  `json:"status"`          // processing / wait_received / fail / received
		Target          string  `json:"target"`          // 客户重定向付款页面的 URL
		FinishAt        string  `json:"finishAt"`        // 完成时间 localDateTime
		Deeplink        string  `json:"deeplink"`        // upi deeplink（需添加前缀 upi://pay?）
	} `json:"data"` // 返回数据
}

type PayCallbackReq struct {
	MerchantOrderNo string  `json:"merchantOrderNo"` // 商户订单号
	OrderNo         string  `json:"orderNo"`         // 订单号
	Amount          float64 `json:"amount"`          // 金额(两位小数)
	Fee             float64 `json:"fee"`             // 手续费(两位小数)
	SettlementFee   string  `json:"settlementFee"`   // 结算手续费（可选）
	Status          string  `json:"status"`          // processing / wait_received / fail / received
	FinishAt        int     `json:"finishAt"`        // 完成时间（localDateTime，可选）
	TargetStatus    string  `json:"targetStatus"`    // creating / valid / invalid
	PlatMsg         string  `json:"platMsg"`         // 失败原因（可选）
}

type QueryPayReq struct {
	// MerchantID      int32  `json:"merchantId"`      // 商户号
	// MerchantOrderID string `json:"merchantOrderId"` // 商户订单
	// Sign            string `json:"sign"`            // 请求签名
}

type QueryPayResp struct {
	Success bool   `json:"success"` // 是否成功
	ErrCode string `json:"errCode"` // 错误码
	ErrMsg  string `json:"errMsg"`  // 错误消息
	Data    struct {
		MerchantOrderNo string  `json:"merchantOrderNo"` // 商户订单号
		OrderNo         string  `json:"orderNo"`         // 订单号
		Amount          float64 `json:"amount"`          // 金额
		Fee             float64 `json:"fee"`             // 手续费
		Status          string  `json:"status"`          // processing / wait_received / fail / received
		Target          string  `json:"target"`          // 重定向 URL（付款页面）
		FinishAt        string  `json:"finishAt"`        // 完成时间（localDateTime）
	} `json:"data"` // 业务数据
}

type WithdrawReq struct {
	AccountNumber   string  `json:"accountNumber"`   // 收款方账号
	NotifyURL       string  `json:"notifyUrl"`       // 回调地址（可选）
	Amount          float64 `json:"amount"`          // 收款金额
	BankID          int64   `json:"bankId"`          // 银行ID（未传默认 1）
	BankCode        string  `json:"bankCode"`        // 银行 code（IFSC 对应 code）
	CustomerName    string  `json:"customerName"`    // 收款人姓名
	Description     string  `json:"description"`     // 转账描述
	MerchantOrderNo string  `json:"merchantOrderNo"` // 商户订单号
}

type WithdrawResp struct {
	Success *bool  `json:"success"` // 是否成功
	ErrCode string `json:"errCode"` // 错误码
	ErrMsg  string `json:"errMsg"`  // 错误消息
	Data    struct {
		OrderNo         string  `json:"orderNo"`         // 订单号
		BankID          int64   `json:"bankId"`          // 银行ID
		BankCode        string  `json:"bankCode"`        // 银行代码
		AccountNumber   string  `json:"accountNumber"`   // 账号号码
		CustomerName    string  `json:"customerName"`    // 客户姓名
		Amount          float64 `json:"amount"`          // 金额
		Fee             float64 `json:"fee"`             // 手续费
		Description     string  `json:"description"`     // 描述
		MerchantOrderNo string  `json:"merchantOrderNo"` // 商户订单号
		Status          string  `json:"status"`          // 状态 creating/processing/success/fail
		PlatMsg         string  `json:"platMsg"`         // 失败原因
		CreatedAt       int64   `json:"createdAt"`       // 创建时间（毫秒时间戳）
	} `json:"data"` // 返回数据
}

type WithdrawCallbackReq struct {
	OrderNo         string  `json:"orderNo"`         // 订单号
	BankID          int64   `json:"bankId"`          // 银行ID
	BankCode        string  `json:"bankCode"`        // 银行代码
	AccountNumber   string  `json:"accountNumber"`   // 账号号码
	CustomerName    string  `json:"customerName"`    // 客户姓名
	Amount          float64 `json:"amount"`          // 金额(两位小数)
	Fee             float64 `json:"fee"`             // 手续费(两位小数)
	Description     string  `json:"description"`     // 描述
	MerchantOrderNo string  `json:"merchantOrderNo"` // 商户订单号
	Status          string  `json:"status"`          // creating/processing/success/fail
	CreatedAt       int64   `json:"createdAt"`       // 创建时间(毫秒)
	TargetStatus    string  `json:"targetStatus"`    // 收款目标：creating/valid/invalid
	PlatMsg         string  `json:"platMsg"`         // 失败原因 (可选)
}

type QueryWithdrawReq struct {
	// MerchantID       int32  `json:"merchantId"`       // 商户id
	// MerchantPayoutID string `json:"merchantPayoutId"` // 商户订单id
	// Sign             string `json:"sign"`             // 请求签名
}

type QueryWithdrawResp struct {
	Success bool   `json:"success"` // 是否成功
	ErrCode string `json:"errCode"` // 错误码
	ErrMsg  string `json:"errMsg"`  // 错误消息
	Data    struct {
		OrderNo         string  `json:"orderNo"`         // 订单号
		BankID          int64   `json:"bankId"`          // 银行ID
		BankCode        string  `json:"bankCode"`        // 银行代码
		AccountNumber   string  `json:"accountNumber"`   // 账号号码
		CustomerName    string  `json:"customerName"`    // 客户姓名
		Amount          float64 `json:"amount"`          // 金额
		Fee             float64 `json:"fee"`             // 手续费
		Description     string  `json:"description"`     // 描述
		MerchantOrderNo string  `json:"merchantOrderNo"` // 商户订单号
		Status          string  `json:"status"`          // 状态：creating/processing/success/fail
		PlatMsg         string  `json:"platMsg"`         // 失败原因
		CreatedAt       int64   `json:"createdAt"`       // 创建时间（毫秒时间戳）
	} `json:"data"` // 返回数据
}
