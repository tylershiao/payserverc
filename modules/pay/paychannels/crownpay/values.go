package crownpay

const (
	baseURL          = "https://apiv2.crownpay8888.com"
	payURL           = "/api/v2/payment/order/create"
	queryPayURL      = "/api/v2/payment/order/query"
	withdrawURL      = "/api/v2/payout/order/create"
	queryWithdrawURL = "/api/v2/payout/order/query"
	balURL           = "/api/v2/merchant/balance"
	// mid              = 2313
	appID = "0b813240c1890a16e6c113ac44f348c9"
	key   = "8f50a09417c529f6c3c7751f347dc33a"
)

type PayReq struct {
	AppID      string `json:"appId"`      // 应用号（必需）
	MerOrderNo string `json:"merOrderNo"` // 商户订单号（必需）
	Currency   string `json:"currency"`   // 币种，3位字符（必需）
	Amount     string `json:"amount"`     // 订单金额，字符串，保留2位小数（必需）
	// ReturnURL  string `json:"returnUrl"`  // 同步通知地址（可选）
	NotifyURL string `json:"notifyUrl"` // 异步通知地址（必需）
	// Attach     string `json:"attach"`     // 附加信息（可选）
	Sign  string `json:"sign"` // 签名（必需）
	Extra struct {
		Name   string `json:"name"`   // 客户姓名（必需）
		Email  string `json:"email"`  // 客户邮箱（必需）
		Mobile string `json:"mobile"` // 客户电话（必需）
	} `json:"extra"` // 扩展信息（必需）
}

type PayResp struct {
	Code int `json:"code"` // 响应状态码，0 为成功（必需）
	Data struct {
		OrderStatus int    `json:"orderStatus"` // 订单状态（必需）
		OrderNo     string `json:"orderNo"`     // 交易订单号（必需）
		MerOrderNo  string `json:"merOrderNo"`  // 商户订单号（必需）
		// Amount      float64 `json:"amount"`            // 金额（必需）
		// Currency    string  `json:"currency"`          // 币种（必需）
		// CreateTime  int64   `json:"createTime"`        // 创建时间（必需）
		// UpdateTime  int64   `json:"updateTime"`        // 更新时间（必需）
		// PayTime     int64   `json:"payTime"`           // 支付时间（必需）
		Message string `json:"message"` // 异常描述（可选，异常时返回）
		Params  struct {
			PaymentLink string `json:"paymentLink"`
			DeepLink    string `json:"deepLink"`
		} `json:"params"` // 支付数据（必需）
		// Attach string `json:"attach"` // 附加信息（必需）
		// Extra  struct {
		// 	Name   string `json:"name"`   // 客户姓名
		// 	Email  string `json:"email"`  // 客户邮箱
		// 	Mobile string `json:"mobile"` // 客户电话
		// } `json:"extra"` // 扩展信息（必需）
	} `json:"data"` // 返回数据（必需）
	Msg  string `json:"msg"`  // 返回文字描述（必需）
	Sign string `json:"sign"` // 签名（必需）
}

type PayCallbackReq struct {
	Sign           string  `json:"sign"`        // 签名
	Amount         float64 `json:"amount"`      // 金额
	Attach         string  `json:"attach"`      // 附加信息
	OrderNo        string  `json:"orderNo"`     // 交易订单号
	MerOrderNo     string  `json:"merOrderNo"`  // 商户订单号
	Currency       string  `json:"currency"`    // 币种
	OrderStatus    int     `json:"orderStatus"` // 订单状态（支付中 0 1 -4, 支付成功 2 3, 支付失败 -1 -2 -3 ）
	CreateTime     int64   `json:"createTime"`  // 创建时间（毫秒）
	UpdateTime     int64   `json:"updateTime"`  // 更新时间（毫秒）
	ChannelTradeNo string  `json:"channelTradeNo"`
}

type QueryPayReq struct {
	AppID      string `json:"appId"`      // 应用号（必需）
	MerOrderNo string `json:"merOrderNo"` // 商户订单号（二选一）
	OrderNo    string `json:"orderNo"`    // 交易订单号（二选一）
	Sign       string `json:"sign"`       // 签名（必需）
}

type QueryPayResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		OrderStatus int    `json:"orderStatus"` // 订单状态
		OrderNo     string `json:"orderNo"`     // 交易订单号
		MerOrderNo  string `json:"merOrderNo"`  // 商户订单号
		// Amount      float64   `json:"amount"`      // 金额
		// Currency    string    `json:"currency"`    // 币种
		// Attach      string    `json:"attach"`      // 附加信息
		// Extra       ExtraInfo `json:"extra"`       // 客户信息
		// CreateTime  int64     `json:"createTime"`  // 创建时间（毫秒）
		// UpdateTime  int64     `json:"updateTime"`  // 更新时间（毫秒）
		Sign string `json:"sign"` // 签名 支付中 0 1 -4 支付成功 2 3 支付失败 -1 -2 -3
	} `json:"data"`
}

type WithdrawReq struct {
	AppID      string `json:"appId"`      // 应用号（必需）
	MerOrderNo string `json:"merOrderNo"` // 商户订单号（必需）
	Currency   string `json:"currency"`   // 币种（必需）
	Amount     string `json:"amount"`     // 金额，string（必需，避免精度/签名问题）
	NotifyURL  string `json:"notifyUrl"`  // 异步通知地址（必需）
	// Attach     string       `json:"attach"`     // 附加信息（可选）
	Sign  string `json:"sign"` // 签名（必需）
	Extra struct {
		PayoutType string `json:"payoutType"` // 固定值[IFSC]
		IFSC       string `json:"ifsc"`       // IFSC Code
		Account    string `json:"account"`    // 银行账户
		Name       string `json:"name"`       // 收款人姓名
		Email      string `json:"email"`      // 收款人邮箱
		Mobile     string `json:"mobile"`     // 收款人手机号
	} `json:"extra"` // 代付信息（必需）
}

type WithdrawResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		OrderStatus int    `json:"orderStatus"` // 订单状态
		OrderNo     string `json:"orderNo"`     // 交易订单号
		MerOrderNo  string `json:"merOrderNo"`  // 商户订单号
		// Amount      float64 `json:"amount"`      // 金额
		// Currency    string  `json:"currency"`    // 币种
		// Attach      string  `json:"attach"`      // 附加信息
		// Extra       struct {
		// 	PayoutType string `json:"payoutType"` // IFSC
		// 	Account    string `json:"account"`    // 银行账户
		// 	IFSC       string `json:"ifsc"`       // IFSC Code
		// 	Name       string `json:"name"`       // 账户名称
		// 	Email      string `json:"email"`      // 邮箱
		// 	Mobile     string `json:"mobile"`     // 手机号
		// } `json:"extra"` // 代付信息
		// CreateTime int64  `json:"createTime"` // 创建时间（毫秒）
		// UpdateTime int64  `json:"updateTime"` // 更新时间（毫秒）
		Sign string `json:"sign"` // 签名
	} `json:"data"`
}

type WithdrawCallbackReq struct {
	Sign           string  `json:"sign"`        // 签名
	Amount         float64 `json:"amount"`      // 金额
	Attach         string  `json:"attach"`      // 附加信息
	OrderNo        string  `json:"orderNo"`     // 平台订单号
	MerOrderNo     string  `json:"merOrderNo"`  // 商户订单号
	Currency       string  `json:"currency"`    // 币种
	OrderStatus    int     `json:"orderStatus"` // 订单状态（支付中 0 1 -4, 支付成功 2 3, 支付失败 -1 -2 -3 ）
	CreateTime     int64   `json:"createTime"`  // 创建时间（毫秒）
	UpdateTime     int64   `json:"updateTime"`  // 更新时间（毫秒）
	Message        string  `json:"message"`
	ChannelTradeNo string  `json:"channelTradeNo"`
}

type QueryWithdrawReq struct {
	AppID      string `json:"appId"`                // 应用号（必需）
	MerOrderNo string `json:"merOrderNo,omitempty"` // 商户订单号（二选一）
	OrderNo    string `json:"orderNo,omitempty"`    // 平台订单号（二选一）
	Sign       string `json:"sign"`                 // 签名（必需）
}

type QueryWithdrawResp struct {
	Code int    `json:"code"` // 响应状态码
	Msg  string `json:"msg"`  // 返回文字描述
	Data struct {
		OrderStatus int    `json:"orderStatus"` // 订单状态
		OrderNo     string `json:"orderNo"`     // 平台交易订单号
		MerOrderNo  string `json:"merOrderNo"`  // 商户订单号
		// Amount      float64 `json:"amount"`      // 金额
		// Currency    string  `json:"currency"`    // 币种
		// Extra       struct {
		// 	PayoutType string `json:"payoutType"` // 代付类型，如 IFSC
		// 	Account    string `json:"account"`    // 银行账户
		// 	IFSC       string `json:"ifsc"`       // IFSC Code
		// 	Name       string `json:"name"`       // 账户名称
		// 	Email      string `json:"email"`      // 邮箱
		// 	Mobile     string `json:"mobile"`     // 手机号
		// } `json:"extra"` // 代付收款人信息
		// CreateTime int64  `json:"createTime"` // 创建时间（毫秒）
		// UpdateTime int64  `json:"updateTime"` // 更新时间（毫秒）
		Message string `json:"message"` // 异常描述（必需，订单异常时返回）
		Sign    string `json:"sign"`    // 签名
	} `json:"data"` // 返回数据
}
