package epay

const (
	baseURL          = "https://api.epay178.com/api/manager/epayController"
	payURL           = "/payJsonV2/create"
	queryPayURL      = "/payJsonV2/get"
	withdrawURL      = "/withdrawJsonV2/create"
	queryWithdrawURL = "/withdrawJsonV2/get"
	balURL           = "/merchantJsonV2/getBalance"
	mid              = "1360930223675162624"
	key              = "BE00F75AE9EDA3CE4008782FBE2036FA"
)

type PayReq struct {
	MchID      string `json:"mch_id"`       // 商户ID
	OutTradeNo string `json:"out_trade_no"` // 商户单号
	Type       string `json:"type"`         // 支付类型（1004）
	Amount     int64  `json:"amount"`       // 订单金额，单位：分
	Currency   string `json:"currency"`     // 货币类型（INR）
	Subject    string `json:"subject"`      // 商品标题
	Extra      string `json:"extra"`        // 代收参数

	// ReturnURL string `json:"return_url"` // 前端跳转地址
	NotifyURL string `json:"notify_url"` // 回调通知地址
	Nonce     string `json:"nonce"`      // 随机数
	Sign      string `json:"sign"`       // 签名字符串
}

type PayExtra struct {
	Name  string `json:"name"`  // 姓名
	Phone string `json:"phone"` // 电话
	Email string `json:"email"` // 邮箱
}

type PayResp struct {
	ReasonCode   string `json:"reasoncode"`   // 返回码：000000 成功
	ReasonRemark string `json:"reasonremark"` // 返回信息
	Data         struct {
		Object        string `json:"object"`         // 固定为 charge
		TransactionID string `json:"transaction_id"` // 平台订单号
		OutTradeNo    string `json:"out_trade_no"`   // 商户单号
		Amount        string `json:"amount"`         // 金额，单位：分
		Currency      string `json:"currency"`       // 货币单位
		CreateAt      string `json:"create_at"`      // 订单创建时间
		PayURL        string `json:"pay_url"`        // 付款链接
	} `json:"data"` // 返回数据
	Nonce string `json:"nonce"` // 随机数
	Sign  string `json:"sign"`  // 签名字符串
}

type PayCallbackReq struct {
	Object        string `json:"object"`         // 通知类型，固定为 charge
	OutTradeNo    string `json:"out_trade_no"`   // 商户单号
	TransactionID string `json:"transaction_id"` // 平台订单号
	Amount        string `json:"amount"`         // 订单金额，单位：分
	Currency      string `json:"currency"`       // 货币类型（INR）
	CreateAt      string `json:"create_at"`      // 订单创建时间
	Paid          string `json:"paid"`           // 付款状态：1成功，2失败，3退款
	PaidAt        string `json:"paid_at"`        // 订单支付时间（未付款可能为空）
	UTR           string `json:"utr"`            // UTR
	Nonce         string `json:"nonce"`          // 随机数
	Sign          string `json:"sign"`           // 签名字符串
}

type QueryPayReq struct {
	MchID         string `json:"mch_id"`         // 商户ID（必填）
	OutTradeNo    string `json:"out_trade_no"`   // 商户单号（必填）
	TransactionID string `json:"transaction_id"` // 平台单号（必填）
	Nonce         string `json:"nonce"`          // 随机数（必填）
	Sign          string `json:"sign"`           // 签名字符串（必填）
}

type QueryPayResp struct {
	ReasonCode   string `json:"reasoncode"`   // 返回码：000000 成功
	ReasonRemark string `json:"reasonremark"` // 返回信息
	Data         struct {
		Object        string `json:"object"`         // 固定为 charge
		OutTradeNo    string `json:"out_trade_no"`   // 商户单号
		TransactionID string `json:"transaction_id"` // 平台订单号
		Amount        string `json:"amount"`         // 金额，单位：分
		Currency      string `json:"currency"`       // 货币单位
		CreateAt      string `json:"create_at"`      // 订单创建时间
		Paid          string `json:"paid"`           // 付款状态：0待处理，1成功，2失败，3退款
		PaidAt        string `json:"paid_at"`        // 订单支付时间（未付款可能为空）
		UTR           string `json:"utr"`            // UTR
	} `json:"data"` // 订单数据
	Nonce string `json:"nonce"` // 随机数
	Sign  string `json:"sign"`  // 签名字符串
}

type WithdrawExtra struct {
	Name  string `json:"name"`  // 收款账户姓名
	Phone string `json:"phone"` // 收款账户电话
	Email string `json:"email"` // 收款账户邮箱

	TransType string `json:"transType"` // BANK_TRANSFER 或 UPI
	TransMode string `json:"transMode"` // IMPS / NEFT / RTGS（BANK_TRANSFER 时必填）
	Account   string `json:"account"`   // 收款账户号（BANK_TRANSFER 时必填）
	IFSC      string `json:"ifsc"`      // 收款账户 IFSC（BANK_TRANSFER 时必填）
	BankName  string `json:"bankName"`  // 收款账户所在银行名称
	VPA       string `json:"vpa"`       // 用户 VPA（UPI 时必填）
}

type WithdrawReq struct {
	MchID      string `json:"mch_id"`       // 商户ID
	OutTradeNo string `json:"out_trade_no"` // 商户单号
	Type       string `json:"type"`         // 代付类型（2001）
	Amount     int64  `json:"amount"`       // 订单金额，单位：分
	Currency   string `json:"currency"`     // 货币类型（INR）
	Subject    string `json:"subject"`      // 用途
	Extra      string `json:"extra"`        // 代付参数

	NotifyURL string `json:"notify_url"` // 回调通知地址

	Nonce string `json:"nonce"` // 随机数
	Sign  string `json:"sign"`  // 签名字符串
}

type WithdrawResp struct {
	ReasonCode   string `json:"reasoncode"`   // 返回码：000000 成功（超时需按成功处理，等回调）
	ReasonRemark string `json:"reasonremark"` // 返回信息
	Data         struct {
		Object        string `json:"object"`         // 固定为 withdraw
		TransactionID string `json:"transaction_id"` // 平台订单号
		OutTradeNo    string `json:"out_trade_no"`   // 商户单号
		Amount        string `json:"amount"`         // 金额，单位：分
		Currency      string `json:"currency"`       // 货币类型（INR）
		CreateAt      string `json:"create_at"`      // 订单创建时间
	} `json:"data"` // 返回数据
	Nonce string `json:"nonce"` // 随机数
	Sign  string `json:"sign"`  // 签名字符串
}

type WithdrawCallbackReq struct {
	Object        string `json:"object"`         // 通知类型，固定为 charge
	MchID         string `json:"mch_id"`         // 商户ID
	OutTradeNo    string `json:"out_trade_no"`   // 商户单号
	TransactionID string `json:"transaction_id"` // 平台订单号
	Amount        string `json:"amount"`         // 订单金额，单位：分
	Currency      string `json:"currency"`       // 货币类型（INR）
	Paid          string `json:"paid"`           // 已转账金额，单位：分（可能为空）
	Result        string `json:"result"`         // 转账结果：1成功，2失败，4回退
	UTR           string `json:"utr"`            // UTR
	Remark        string `json:"remark"`         // 备注（失败原因等）
	CreateAt      string `json:"create_at"`      // 订单创建时间
	PaidAt        string `json:"paid_at"`        // 转账时间
	Nonce         string `json:"nonce"`          // 随机数
	Sign          string `json:"sign"`           // 签名字符串
}

type QueryWithdrawReq struct {
	MchID         string `json:"mch_id"`         // 商户ID（必填）
	OutTradeNo    string `json:"out_trade_no"`   // 商户单号（必填）
	TransactionID string `json:"transaction_id"` // 平台订单号（必填）
	Nonce         string `json:"nonce"`          // 随机数（必填）
	Sign          string `json:"sign"`           // 签名字符串（必填）
}

type QueryWithdrawResp struct {
	ReasonCode   string `json:"reasoncode"`   // 返回码：000000 成功
	ReasonRemark string `json:"reasonremark"` // 返回信息
	Data         struct {
		Object        string `json:"object"`         // 固定为 charge
		MchID         string `json:"mch_id"`         // 商户ID
		OutTradeNo    string `json:"out_trade_no"`   // 商户单号
		TransactionID string `json:"transaction_id"` // 平台订单号
		Amount        string `json:"amount"`         // 金额，单位：分
		Currency      string `json:"currency"`       // 货币类型（INR）
		Paid          string `json:"paid"`           // 已转账金额，单位：分
		Result        string `json:"result"`         // 转账结果：0待处理，1成功，2失败，4回退
		PaidAt        string `json:"paid_at"`        // 转账时间
		CreateAt      string `json:"create_at"`      // 订单创建时间
		UTR           string `json:"utr"`            // UTR
	} `json:"data"` // 返回数据
	Nonce string `json:"nonce"` // 随机数
	Sign  string `json:"sign"`  // 签名字符串
}
