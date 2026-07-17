package wheatpay

const (
	baseURL          = "https://api.nicepayhope.com"
	payURL           = "/collection/upi/v1/initial"
	queryPayURL      = "/collection/upi/check/status"
	withdrawURL      = "/api/payout/v2/transfer"
	queryWithdrawURL = "/api/payout/v1/check/status"
	// balURL           = "/api/payout/v1/get/balance"
	mid    = "M231182"
	appID  = "240001861"
	key    = "mf3EfwB8LqMMFRdQ"
	secret = "RWDHHuVLQaCKodEyECmqHaIb7BiFUzET"
)

type PayReq struct {
	AppID       string  `json:"appId"`       // appId（后台获取）
	TransferID  string  `json:"transferId"`  // 商户订单号
	Currency    string  `json:"currency"`    // 币种：INR
	Amount      float64 `json:"amount"`      // 金额（保留两位小数）
	PayerName   string  `json:"payerName"`   // 昵称
	PayerEmail  string  `json:"payerEmail"`  // 邮箱
	PayerPhone  string  `json:"payerPhone"`  // 手机号
	CallbackURL string  `json:"callbackUrl"` // 回调地址（可选）
	Hash        string  `json:"hash"`        // HmacSHA256 签名
}

type PayResp struct {
	Status  int    `json:"status"`  // 200: 请求成功
	Message string `json:"message"` // 失败原因（status != 200）
	Data    struct {
		TransferID   string  `json:"transferId"`   // 商户订单号
		TxnReference string  `json:"txnReference"` // 平台交易号
		Amount       float64 `json:"amount"`       // 订单金额
		AppID        string  `json:"appId"`        // appId
		PaymentPgURL string  `json:"paymentPgUrl"` // 支付链接
	} `json:"data"` // 成功时返回
}

type PayCallbackReq struct {
	AppID           string  `json:"appId"`           // appId
	TransferID      string  `json:"transferId"`      // 商户订单号
	TxnReference    string  `json:"txnReference"`    // 平台交易号
	Status          string  `json:"status"`          // SUCCESS / FAILED / PENDING
	TxnAmount       float64 `json:"txnAmount"`       // 实际支付金额
	Amount          float64 `json:"amount"`          // 下单金额
	TransferMessage string  `json:"transferMessage"` // 支付消息
	UTR             string  `json:"utr"`             // 支付 UTR
	Hash            string  `json:"hash"`            // HmacSHA256 签名
}

type QueryPayReq struct {
	AppID      string `json:"appId"`      // appId（后台获取）
	TransferID string `json:"transferId"` // 商户订单号
	Hash       string `json:"hash"`       // HmacSHA256 签名
}

type QueryPayResp struct {
	Status  int    `json:"status"`  // 200: 请求成功（不代表订单成功）
	Message string `json:"message"` // 接口消息
	Data    struct {
		AppID           string  `json:"appId"`           // appId
		TransferID      string  `json:"transferId"`      // 商户订单号
		TxnReference    string  `json:"txnReference"`    // 平台交易号
		Status          string  `json:"status"`          // SUCCESS / FAILED / PENDING
		TxnAmount       float64 `json:"txnAmount"`       // 实际支付金额
		Amount          float64 `json:"amount"`          // 下单金额
		TransferMessage string  `json:"transferMessage"` // 订单消息
		UTR             string  `json:"utr"`             // 支付 UTR
	} `json:"data"` // status == 200 时返回
}

type WithdrawReq struct {
	AppID       string `json:"appId"`       // appId（后台获取）
	TransferID  string `json:"transferId"`  // 商户订单号
	Currency    string `json:"currency"`    // 币种：INR
	Amount      string `json:"amount"`      // 金额（保留两位小数）
	PayerName   string `json:"payerName"`   // 昵称
	PayerEmail  string `json:"payerEmail"`  // 邮箱
	PayerPhone  string `json:"payerPhone"`  // 手机
	BankAccount string `json:"bankAccount"` // 银行卡号
	IFSC        string `json:"ifsc"`        // IFSC code
	CallbackURL string `json:"callbackUrl"` // 回调地址（可选）
	Hash        string `json:"hash"`        // HmacSHA256 签名
}

type WithdrawResp struct {
	Status  int    `json:"status"`  // 200: 请求成功，其他失败
	Message string `json:"message"` // 接口消息
	Data    struct {
		AppID        string  `json:"appId"`        // appId
		TransferID   string  `json:"transferId"`   // 商户订单号
		TxnReference string  `json:"txnReference"` // 平台交易号
		Amount       float64 `json:"amount"`       // 金额（单位:分）
		Message      string  `json:"message"`      // 订单消息
		Status       string  `json:"status"`       // PENDING / SUCCESS / FAILED
	} `json:"data"` // status == 200 时返回
}

type WithdrawCallbackReq struct {
	AppID           string  `json:"appId"`           // appId
	TransferID      string  `json:"transferId"`      // 商户订单号
	TxnReference    string  `json:"txnReference"`    // 平台交易号
	Status          string  `json:"status"`          // SUCCESS / FAILED / PENDING
	Amount          float64 `json:"amount"`          // 订单金额
	TransferMessage string  `json:"transferMessage"` // 订单消息
	UTR             string  `json:"utr"`             // 付款 UTR
	Hash            string  `json:"hash"`            // HmacSHA256 签名
}

type QueryWithdrawReq struct {
	AppID      string `json:"appId"`      // appId
	TransferID string `json:"transferId"` // 商户订单号
	Hash       string `json:"hash"`       // HmacSHA256 签名
}

type QueryWithdrawResp struct {
	Status  int    `json:"status"`  // 200: 请求成功，不代表订单成功
	Message string `json:"message"` // 接口消息
	Data    struct {
		AppID        string  `json:"appId"`        // appId
		TransferID   string  `json:"transferId"`   // 商户订单号
		TxnReference string  `json:"txnReference"` // 平台交易号
		Amount       float64 `json:"amount"`       // 订单金额
		Message      string  `json:"message"`      // 订单消息
		UTR          string  `json:"utr"`          // 订单 utr
		Status       string  `json:"status"`       // SUCCESS / FAILED / PENDING
	} `json:"data"` // status == 200 时返回
}
