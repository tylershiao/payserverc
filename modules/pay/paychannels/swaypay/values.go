package swaypay

const (
	baseURL          = "https://hub.swaypay.xin"
	payURL           = "/api/v1/payin/create"
	queryPayURL      = "/api/v1/payin/query"
	withdrawURL      = "/api/v1/payout/unifiedOrder"
	queryWithdrawURL = "/api/v1/payout/payOutInfo"
	balURL           = "/api/v1/payout/getBalance"
	mid              = "slot8899"
	key              = "bbc6f5a9e03c5c2070597da1a50aa146873299b5c74d1fbe5e20d21c9cc162c3"
)

type PayReq struct {
	Amount     string `json:"amount"`
	Currency   string `json:"currency"`
	MchNo      string `json:"mchNo"`
	MchOrderNo string `json:"mchOrderNo"`
	// Mark       string `json:"mark"`
	NotifyURL  string `json:"notifyUrl"`
	PayinEmail string `json:"payinEmail"`
	PayinName  string `json:"payinName"`
	PayinPhone string `json:"payinPhone"`
	ReturnURL  string `json:"returnUrl"`
	// Version    string `json:"version"`
	Sign string `json:"sign"`
}

type PayResp struct {
	Code int    `json:"code"` // 提交状态
	Msg  string `json:"msg"`  // 提示信息
	Data struct {
		MchNo         string `json:"mchNo"`         // 商户号
		Version       string `json:"version"`       // 版本号
		MchOrderNo    string `json:"mchOrderNo"`    // 商户订单号
		SystemOrderNo string `json:"systemOrderNo"` // 平台订单号
		Amount        string `json:"amount"`        // 金额（两位小数）
		Currency      string `json:"currency"`      // 币种
		PayinName     string `json:"payinName"`     // 充值人姓名
		Upi           string `json:"upi"`           // 收款 UPI
		PayinState    int    `json:"payinState"`    // 代收状态
		PayinH5Url    string `json:"payinH5Url"`    // 支付链接
		Sign          string `json:"sign"`          // 签名
	} `json:"data"`
}

type PayCallbackReq struct {
	Code int    `json:"code"` // 提交状态
	Msg  string `json:"msg"`  // 提示信息
	Data struct {
		MchNo              string `json:"mchNo"`              // 商户号
		Version            string `json:"version"`            // 版本号
		Sign               string `json:"sign"`               // 签名
		MchOrderNo         string `json:"mchOrderNo"`         // 商户订单号
		SystemOrderNo      string `json:"systemOrderNo"`      // 平台订单号
		Amount             string `json:"amount"`             // 充值金额
		PayinRate          string `json:"payinRate"`          // 手续费
		PayinState         int    `json:"payinState"`         // 代收状态
		ActualAmount       string `json:"actualAmount"`       // 实际充值金额
		PayinTransactionId string `json:"payinTransactionId"` // 支付流水 ID
		Utr                string `json:"utr"`                // 订单 UTR
		PayinSettleAmount  string `json:"payinSettleAmount"`  // 结算金额
		SuccessTime        string `json:"successTime"`        // 成功时间
		Mark               string `json:"mark"`               // 订单备注
		PayinName          string `json:"payinName"`          // 收款人姓名
		Upi                string `json:"upi"`                // 收款 UPI
	} `json:"data"`
}

type QueryPayReq struct {
	MchNo      string `json:"mchNo"`      // 商户号
	Version    string `json:"version"`    // 版本号（非必填，默认 V1）
	MchOrderNo string `json:"mchOrderNo"` // 商户订单号（唯一）
	Sign       string `json:"sign"`       // 签名
}

type QueryPayResp struct {
	Code int    `json:"code"` // 提交状态
	Msg  string `json:"msg"`  // 提示信息
	Data struct {
		MchNo              string `json:"mchNo"`              // 商户号
		Version            string `json:"version"`            // 版本号
		MchOrderNo         string `json:"mchOrderNo"`         // 商户订单号
		SystemOrderNo      string `json:"systemOrderNo"`      // 平台订单号
		Amount             string `json:"amount"`             // 金额（两位小数）
		PayinRate          string `json:"payinRate"`          // 手续费
		PayinState         int    `json:"payinState"`         // 代收状态 (1:交易中/2:交易成功/3:部分支付/4:超时未付/-1:未知状态)
		ActualAmount       string `json:"actualAmount"`       // 实际到账金额
		PayinSettleAmount  string `json:"payinSettleAmount"`  // 结算金额(收款金额-手续费)
		Utr                string `json:"utr"`                // UTR
		PayinName          string `json:"payinName"`          // 充值人姓名
		Upi                string `json:"upi"`                // 收款 UPI
		PayinTransactionId string `json:"payinTransactionId"` // 支付流水 ID
		SuccessTime        string `json:"successTime"`        // 成功时间
		Mark               string `json:"mark"`               // 备注
		Sign               string `json:"sign"`               // 签名
	} `json:"data"`
}

type WithdrawReq struct {
	MchNo string `json:"mchNo"` // 商户号
	// Version    string `json:"version"`    // 版本号（非必填，默认 V1）
	MchOrderNo string `json:"mchOrderNo"` // 商户订单号（唯一）
	Amount     string `json:"amount"`     // 金额（两位小数）
	CashName   string `json:"cashName"`   // 收款人姓名
	CashPhone  string `json:"cashPhone"`  // 收款人真实手机号
	Email      string `json:"email"`      // 邮箱（可选）
	IFSC       string `json:"ifsc"`       // IFSC
	BankName   string `json:"bankName"`   // 收款银行名称
	CashNumber string `json:"cashNumber"` // 收款人账号（9–18 位数字）
	NotifyUrl  string `json:"notifyUrl"`  // 回调地址
	Sign       string `json:"sign"`       // 签名
}

type WithdrawResp struct {
	Code int    `json:"code"` // 提交状态
	Msg  string `json:"msg"`  // 提示信息
	Data struct {
		MchNo           string `json:"mchNo"`           // 商户号
		Version         string `json:"version"`         // 版本号
		MchOrderNo      string `json:"mchOrderNo"`      // 商户订单号
		SystemOrderNo   string `json:"systemOrderNo"`   // 平台订单号
		Amount          string `json:"amount"`          // 金额（两位小数）
		Rate            string `json:"rate"`            // 总手续费
		MerchantRate    string `json:"merchantRate"`    // 商家费率
		TransactionRate string `json:"transactionRate"` // 单笔手续费
		SettleAmount    string `json:"settleAmount"`    // 结算金额
		NotifyUrl       string `json:"notifyUrl"`       // 回调地址
		CashState       int    `json:"cashState"`       // 代付状态
		CashName        string `json:"cashName"`        // 收款人姓名
		CashNumber      string `json:"cashNumber"`      // 收款人账号
		CashPhone       string `json:"cashPhone"`       // 收款用户手机号
		BankName        string `json:"bankName"`        // 收款银行名称
		Email           string `json:"email"`           // 邮箱
		IFSC            string `json:"ifsc"`            // IFSC
		CreatedAt       string `json:"createdAt"`       // 创建时间
		SuccessTime     string `json:"successTime"`     // 成功时间
		Sign            string `json:"sign"`            // 签名
	} `json:"data"`
}

type WithdrawCallbackReq struct {
	Code int    `json:"code"` // 提交状态
	Msg  string `json:"msg"`  // 提示信息
	Data struct {
		MchNo           string `json:"mchNo"`           // 商户号
		Version         string `json:"version"`         // 版本号
		MchOrderNo      string `json:"mchOrderNo"`      // 商户订单号
		SystemOrderNo   string `json:"systemOrderNo"`   // 平台订单号
		Amount          string `json:"amount"`          // 金额（两位小数）
		Rate            string `json:"rate"`            // 总手续费
		MerchantRate    string `json:"merchantRate"`    // 商家费率
		TransactionRate string `json:"transactionRate"` // 单笔手续费
		SettleAmount    string `json:"settleAmount"`    // 结算金额
		NotifyUrl       string `json:"notifyUrl"`       // 回调地址
		CashState       int    `json:"cashState"`       // 代付状态
		CashName        string `json:"cashName"`        // 收款人姓名
		CashNumber      string `json:"cashNumber"`      // 收款人账号
		CashPhone       string `json:"cashPhone"`       // 收款用户手机号
		BankName        string `json:"bankName"`        // 收款银行名称
		Email           string `json:"email"`           // 邮箱
		IFSC            string `json:"ifsc"`            // IFSC
		Utr             string `json:"utr"`             // 订单 UTR
		CreatedAt       string `json:"createdAt"`       // 创建时间
		SuccessTime     string `json:"successTime"`     // 成功时间
		Sign            string `json:"sign"`            // 签名
	} `json:"data"`
}

type QueryWithdrawReq struct {
	MchNo      string `json:"mchNo"`      // 商户号
	MchOrderNo string `json:"mchOrderNo"` // 商户订单号（唯一）
	Sign       string `json:"sign"`       // 签名
}

type QueryWithdrawResp struct {
	Code int    `json:"code"` // 提交状态
	Msg  string `json:"msg"`  // 提示信息
	Data struct {
		MchNo           string `json:"mchNo"`           // 商户号
		Version         string `json:"version"`         // 版本号
		MchOrderNo      string `json:"mchOrderNo"`      // 商户订单号
		SystemOrderNo   string `json:"systemOrderNo"`   // 平台订单号
		Amount          string `json:"amount"`          // 金额（两位小数）
		Rate            string `json:"rate"`            // 总手续费
		MerchantRate    string `json:"merchantRate"`    // 商家费率
		TransactionRate string `json:"transactionRate"` // 单笔手续费
		SettleAmount    string `json:"settleAmount"`    // 结算金额
		NotifyUrl       string `json:"notifyUrl"`       // 回调地址
		CashState       int    `json:"cashState"`       // 代付状态:1 订单创建/2 支付中/3 支付失败/4 支付成功/5 拒绝/6 锁单中
		CashName        string `json:"cashName"`        // 收款人姓名
		CashNumber      string `json:"cashNumber"`      // 收款人账号
		CashPhone       string `json:"cashPhone"`       // 收款用户手机号
		BankName        string `json:"bankName"`        // 收款银行名称
		Email           string `json:"email"`           // 邮箱
		IFSC            string `json:"ifsc"`            // IFSC
		Utr             string `json:"utr"`             // 订单 UTR
		CreatedAt       string `json:"createdAt"`       // 创建时间
		SuccessTime     string `json:"successTime"`     // 成功时间
		Sign            string `json:"sign"`            // 签名
	} `json:"data"`
}
