package epayglobal

const (
	baseURL          = "https://order.8ae.one/payApi/PayApi"
	payURL           = "/CreateOrder"
	queryPayURL      = "/QueryOrder"
	withdrawURL      = "/CreatePayOutOrder"
	queryWithdrawURL = "/QueryOrder"
	balURL           = "/Balance"
	mid              = 2524
	appID            = 1716
	key              = "6599ac69a0ba64db920c9801f3d9c80d"
)

type PayReq struct {
	TradeNo int    `json:"trade_no"` // 商户号
	AppID   int    `json:"app_id"`   // 商户 appId
	PayCode int    `json:"pay_code"` // 产品编码
	Price   int    `json:"price"`    // 下单金额（单位：分）
	OrderNo string `json:"order_no"` // 商户订单号
	// SuccessURL   string `json:"success_url"`    // 支付成功跳转 URL
	// FailURL      string `json:"fail_url"`       // 支付失败跳转 URL
	PayNoticeURL string `json:"pay_notice_url"` // 支付成功通知 URL
	UserID       string `json:"user_id"`        // 用户 ID
	UserIP       string `json:"user_ip"`        // 付款人 IP
	// Attach       string `json:"attach"`         // 附加参数
	Sign string `json:"sign"` // 签名
}

type PayResp struct {
	Code       int    `json:"code"`         // 200: 下单成功，其他: 下单失败
	Msg        string `json:"msg"`          // 失败原因
	PayURL     string `json:"pay_url"`      // 付款链接
	QRCode     string `json:"qr_code"`      // 二维码字符串
	OrderNo    string `json:"order_no"`     // 商户订单号
	DisOrderNo string `json:"dis_order_no"` // 平台订单号
	CreateTime int    `json:"create_time"`  // 创建时间（时间戳）
	Sign       string `json:"sign"`         // 签名
}

type PayCallbackReq struct {
	TradeNo    int    `json:"trade_no"`     // 商户号
	Status     int    `json:"status"`       // 订单状态：2 成功，3 失败
	OrderNo    string `json:"order_no"`     // 商户订单号
	DisOrderNo string `json:"dis_order_no"` // 平台订单号
	OrderPrice int    `json:"order_price"`  // 订单金额（分）
	RealPrice  int    `json:"real_price"`   // 实际支付金额（分）
	Fee        int    `json:"fee"`          // 手续费（分）
	NtiTime    int    `json:"nti_time"`     // 通知时间
	Payer      string `json:"payer"`        // 付款人信息（JSON 字符串）
	Attach     string `json:"attach"`       // 附加参数
	CreateTime int    `json:"create_time"`  // 创建时间
	Sign       string `json:"sign"`         // 签名
}

type QueryPayReq struct {
	OrderType   string `json:"order_type"`    // pay_out: 代付, pay_in: 代收
	TradeNo     int    `json:"trade_no"`      // 商户号
	OrderNo     string `json:"order_no"`      // 商户订单号（可选）
	BankOrderNo string `json:"bank_order_no"` // 代收 UTR
	Sign        string `json:"sign"`          // 签名
}

type QueryPayResp struct {
	Code        int    `json:"code"`          // 200: 查询成功，其他: 失败
	Msg         string `json:"msg"`           // 查询失败原因
	TradeNo     int    `json:"trade_no"`      // 商户号
	RealPrice   int    `json:"real_price"`    // 实际付款金额（分）
	Status      int    `json:"status"`        // 1未支付 / 2成功 / 3失败 / 7驳回 / 8冲正
	SuccessTime int    `json:"success_time"`  // 成功时间戳
	OrderNo     string `json:"order_no"`      // 商户订单号
	DisOrderNo  string `json:"dis_order_no"`  // 平台订单号
	BankOrderNo string `json:"bank_order_no"` // 银行订单号
	Remark      string `json:"remark"`        // 失败原因
	Fee         int    `json:"fee"`           // 手续费（分）
	CreateTime  int    `json:"create_time"`   // 创建时间
	Payer       string `json:"payer"`         // 付款人信息（JSON 字符串）
	Sign        string `json:"sign"`          // 签名
}

type WithdrawReq struct {
	TradeNo      int    `json:"trade_no"`       // 商户号
	OrderNo      string `json:"order_no"`       // 商户订单号
	AppID        int    `json:"app_id"`         // 商户 appId
	PayCode      int    `json:"pay_code"`       // 产品编码
	Price        int    `json:"price"`          // 下单金额（分）
	AccountNo    string `json:"account_no"`     // 出款账号
	AccountType  string `json:"account_type"`   // 账号类型：BANK / UPI
	AccountName  string `json:"account_name"`   // 姓名
	BankCode     string `json:"bank_code"`      // 银行编码
	IdentifyNum  string `json:"identify_num"`   // IFSC 号码
	PayNoticeURL string `json:"pay_notice_url"` // 代付成功通知 URL
	// Attach       string `json:"attach"`         // 附加参数
	Sign string `json:"sign"` // 签名
}

type WithdrawResp struct {
	Code       int    `json:"code"`         // 200: 成功，其他: 失败
	Msg        string `json:"msg"`          // 失败原因
	DisOrderNo string `json:"dis_order_no"` // 平台订单号
	OrderNo    string `json:"order_no"`     // 商户订单号
	Status     int    `json:"status"`       // 2成功 / 3失败 / 7驳回 / 8冲正 / 其他处理中
	CreateTime int    `json:"create_time"`  // 创建时间
	Sign       string `json:"sign"`         // 签名
}

type WithdrawCallbackReq struct {
	TradeNo    int    `json:"trade_no"`     // 商户号
	OrderNo    string `json:"order_no"`     // 商户订单号
	DisOrderNo string `json:"dis_order_no"` // 平台订单号
	Price      int    `json:"price"`        // 订单金额（分）
	Fee        int    `json:"fee"`          // 手续费（分）
	Status     int    `json:"status"`       // 2成功 / 3失败 / 7驳回 / 8冲正
	NtiTime    int    `json:"nti_time"`     // 通知时间
	Attach     string `json:"attach"`       // 附加参数
	Remark     string `json:"remark"`       // 失败原因
	CreateTime int    `json:"create_time"`  // 创建时间
	Sign       string `json:"sign"`         // 签名
}

type QueryWithdrawReq struct {
	OrderType   string `json:"order_type"`    // pay_out: 代付, pay_in: 代收
	TradeNo     int    `json:"trade_no"`      // 商户号
	OrderNo     string `json:"order_no"`      // 商户订单号（可选）
	BankOrderNo string `json:"bank_order_no"` // 代收 UTR
	Sign        string `json:"sign"`          // 签名
}

type QueryWithdrawResp struct {
	Code        int    `json:"code"`          // 200: 查询成功，其他: 失败
	Msg         string `json:"msg"`           // 查询失败原因
	TradeNo     int    `json:"trade_no"`      // 商户号
	RealPrice   int    `json:"real_price"`    // 实际付款金额（分）
	Status      int    `json:"status"`        // 1未支付 / 2成功 / 3失败 / 7驳回 / 8冲正
	SuccessTime int    `json:"success_time"`  // 成功时间戳
	OrderNo     string `json:"order_no"`      // 商户订单号
	DisOrderNo  string `json:"dis_order_no"`  // 平台订单号
	BankOrderNo string `json:"bank_order_no"` // 银行订单号
	Remark      string `json:"remark"`        // 失败原因
	Fee         int    `json:"fee"`           // 手续费（分）
	CreateTime  int    `json:"create_time"`   // 创建时间
	Payer       string `json:"payer"`         // 付款人信息（JSON 字符串）
	Sign        string `json:"sign"`          // 签名
}
