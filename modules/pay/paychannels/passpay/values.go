package passpay

const (
	baseURL          = "https://api.merchant.passpay.cc"
	payURL           = "/api/developer/order/create"
	queryPayURL      = "/api/developer/order/query"
	withdrawURL      = "/api/developer/payout/create"
	queryWithdrawURL = "/api/developer/payout/query"
	mid              = "14203622"
	key              = "yrBH0YbUBK0byECaXWlR2zfA651bPV6K"
	payID            = "12"
	// mid = "202366100"
	// key = "c385fe7029344aef826d8112625b2625x"
)

type PayReq struct {
	Mchid      string `json:"mchid"`        // 由我方分配
	PayID      string `json:"pay_id"`       // 由我方分配
	OutTradeNo string `json:"out_trade_no"` // 商户传入订单号，唯一
	Amount     string `json:"amount"`       // 订单金额：100.00、500.00
	NotifyURL  string `json:"notify_url"`   // 回调地址
	Sign       string `json:"sign"`         // 签名
}

type PayResp struct {
	Data struct {
		Status     int    `json:"status"`       // 订单状态：1-订单创建成功 其他-失败
		Msg        string `json:"msg"`          // 状态描述
		PayID      string `json:"pay_id"`       // pay_id 测试时填：11 ，如果提示通道关闭，请联系运营进行配置
		TradeNo    string `json:"trade_no"`     // 我方订单号，唯一
		OutTradeNo string `json:"out_trade_no"` // 商户传入订单号，唯一
		Amount     string `json:"amount"`       // amount
		PayLink    string `json:"pay_link"`     // 落地页链接
	} `json:"data"`
	RCode int    `json:"rCode"` // http状态码 这只是Http状态
	RMsg  string `json:"rMsg"`  // http描述
}

type PayCallbackReq struct {
	Mchid      string `json:"mchid"`        // 由我方分配
	PayID      string `json:"pay_id"`       // 由我方分配
	TradeNo    string `json:"trade_no"`     // 我方单号，唯一
	OutTradeNo string `json:"out_trade_no"` // 商户传入订单号，唯一
	UTR        string `json:"utr"`          // UTR值，可能为空
	Amount     string `json:"amount"`       // 订单金额：100.00 与 real_amount 一样
	RealAmount string `json:"real_amount"`  // 实际支付金额：100.00 与 amount 一样
	Status     int    `json:"status"`       // 交易状态：3-交易中 5-交易成功 6-交易失败
	Msg        string `json:"msg"`          // 订单状态描述
	Sign       string `json:"sign"`         // 签名
}

type QueryPayReq struct {
	Mchid      string `json:"mchid"`        // 由我方分配
	TradeNo    string `json:"trade_no"`     // 我方单号，唯一
	OutTradeNo string `json:"out_trade_no"` // 商户传入的订单号，唯一
	Sign       string `json:"sign"`         // 签名
}

type QueryPayResp struct {
	Data struct {
		Status     int    `json:"status"`       // 交易状态：3-交易中 5-交易成功 6-交易失败
		Msg        string `json:"msg"`          // 状态描述
		TradeNo    string `json:"trade_no"`     // 我方订单号，唯一
		OutTradeNo string `json:"out_trade_no"` // 商户传入订单号，唯一
		UTR        string `json:"utr"`          // UTR值，可能为空
		Amount     string `json:"amount"`       // 订单金额 200.00 此为结算金额
		RealAmount string `json:"real_amount"`  // 实际支付金额 200.00
	} `json:"data"`
	RCode int    `json:"rCode"` // http状态码
	RMsg  string `json:"rMsg"`  // http描述
}

type WithdrawReq struct {
	Mchid         string `json:"mchid"`           // 由我方分配
	OutTradeNo    string `json:"out_trade_no"`    // 商户传入订单号，唯一
	Amount        string `json:"amount"`          // 订单金额 200.00、100.00
	ToHolder      string `json:"to_holder"`       // 持卡人姓名
	ToAccount     string `json:"to_account"`      // 银行卡卡号
	ToIfsc        string `json:"to_ifsc"`         // IFSC
	ToPhoneNumber string `json:"to_phone_number"` // 手机号
	NotifyURL     string `json:"notify_url"`      // 回调地址
	Sign          string `json:"sign"`            // 签名
}

type WithdrawResp struct {
	Data struct {
		Status        int    `json:"status"`          // 创建状态：1-创建成功 其他-创建失败
		Msg           string `json:"msg"`             // 状态描述
		TradeNo       string `json:"trade_no"`        // 我方订单号，唯一
		OutTradeNo    string `json:"out_trade_no"`    // 商户传入订单号，唯一
		Amount        string `json:"amount"`          // amount 100.00
		ToHolder      string `json:"to_holder"`       // 持卡人姓名
		ToAccount     string `json:"to_account"`      // 银行卡卡号
		ToIfsc        string `json:"to_ifsc"`         // IFSC
		ToPhoneNumber string `json:"to_phone_number"` // 手机号
	} `json:"data"`
	RCode int    `json:"rCode"` // http状态码
	RMsg  string `json:"rMsg"`  // http描述
}

type WithdrawCallbackReq struct {
	Mchid         string `json:"mchid"`           // 商户号，我方分配
	Status        int    `json:"status"`          // 代付交易状态：3-交易中 5-代付交易成功 6-代付交易失败
	Msg           string `json:"msg"`             // 状态描述
	TradeNo       string `json:"trade_no"`        // 我方单号，唯一
	OutTradeNo    string `json:"out_trade_no"`    // 商户传入订单号，唯一
	UTR           string `json:"utr"`             // UTR值，可能为空
	Amount        string `json:"amount"`          // 代付订单金额
	ToHolder      string `json:"to_holder"`       // 持卡人姓名
	ToAccount     string `json:"to_account"`      // 银行卡卡号
	ToIfsc        string `json:"to_ifsc"`         // IFSC
	ToPhoneNumber string `json:"to_phone_number"` // 手机号
	Sign          string `json:"sign"`            // 签名
}

type QueryWithdrawReq struct {
	Mchid      string `json:"mchid"`        // 由我方分配
	TradeNo    string `json:"trade_no"`     // 我方单号，唯一
	OutTradeNo string `json:"out_trade_no"` // 商户传入的订单号，唯一
	Sign       string `json:"sign"`         // 签名
}

type QueryWithdrawResp struct {
	Data struct {
		Status        int    `json:"status"`          // 代付交易状态：3-交易中 5-交易成功 6-代付交易失败
		Msg           string `json:"msg"`             // 状态描述
		TradeNo       string `json:"trade_no"`        // 我方订单号，唯一
		OutTradeNo    string `json:"out_trade_no"`    // 商户传入订单号，唯一
		UTR           string `json:"utr"`             // UTR值，可能为空
		Amount        string `json:"amount"`          // 订单金额
		ToHolder      string `json:"to_holder"`       // 持卡人姓名
		ToAccount     string `json:"to_account"`      // 银行卡卡号
		ToIfsc        string `json:"to_ifsc"`         // IFSC
		ToPhoneNumber string `json:"to_phone_number"` // 手机号
		PayoutCertURL string `json:"payout_cert_url"` // 代付成功凭证，仅代付成功有值
	} `json:"data"`
	RCode int    `json:"rCode"` // http状态码
	RMsg  string `json:"rMsg"`  // http描述
}
