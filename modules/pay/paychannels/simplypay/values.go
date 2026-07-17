package simplypay

const (
	baseURL          = "https://api.paysimply.net"
	payURL           = "/api/v2/payment/order/create"
	queryPayURL      = "/api/v2/payment/order/query"
	withdrawURL      = "/api/v2/payout/order/create"
	queryWithdrawURL = "/api/v2/payout/order/query"
	balURL           = "/api/v2/merchant/balance"
	appID            = "0885b2143f163782c2aec4d93323a98d"
	key              = "032bcfbd6262cfa3729035fceeacbe07"
)

type PayReq struct {
	AppId      string `json:"appId"`
	MerOrderNo string `json:"merOrderNo"`
	Currency   string `json:"currency"`
	Amount     string `json:"amount"`
	// ReturnUrl  string `json:"returnUrl"`
	NotifyUrl string `json:"notifyUrl"`
	// Attach     string `json:"attach"`
	Sign  string `json:"sign"`
	Extra struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		Mobile string `json:"mobile"`
		// CashierTheme string `json:"cashierTheme"`
	} `json:"extra"`
}

type PayResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		OrderStatus int     `json:"orderStatus"`
		OrderNo     string  `json:"orderNo"`
		MerOrderNo  string  `json:"merOrderNo"`
		Amount      float64 `json:"amount"`
		Currency    string  `json:"currency"`
		CreateTime  int64   `json:"createTime"`
		UpdateTime  int64   `json:"updateTime"`
		Message     string  `json:"message"`
		Params      struct {
			// 支付数据，自由扩展
			// 如 paymentLink / deepLink 等
			PaymentLink string `json:"paymentLink"`
			DeepLink    string `json:"deepLink"`
		} `json:"params"`
		Sign   string `json:"sign"`
		Attach string `json:"attach"`
		Extra  struct {
			Name   string `json:"name"`
			Email  string `json:"email"`
			Mobile string `json:"mobile"`
		} `json:"extra"`
		PayTime int64 `json:"payTime"`
	} `json:"data"`
}

type PayCallbackReq struct {
	OrderStatus int    `json:"orderStatus"` // 订单状态（支付中 0 1 -4, 支付成功 2 3, 支付失败 -1 -2 -3）
	OrderNo     string `json:"orderNo"`     // 平台交易订单号
	MerOrderNo  string `json:"merOrderNo"`  // 商户订单号
	Amount      int64  `json:"amount"`      // 订单金额
	Currency    string `json:"currency"`    // 金额币种
	Attach      string `json:"attach"`      // 附加信息（原样返回）
	CreateTime  int64  `json:"createTime"`  // 创建时间（时间戳毫秒）
	UpdateTime  int64  `json:"updateTime"`  // 更新时间（时间戳毫秒）
	Message     string `json:"message"`     // 描述信息（异常返回）
	Sign        string `json:"sign"`        // 签名
}

type QueryPayReq struct {
	AppId      string `json:"appId"`
	MerOrderNo string `json:"merOrderNo"`
	OrderNo    string `json:"orderNo"`
	Sign       string `json:"sign"`
}

/*
支付中 0 1 -4
支付成功 2 3
支付失败 -1 -2 -3
*/
type QueryPayResp struct {
	Code int `json:"code"`
	Data struct {
		OrderStatus int     `json:"orderStatus"`
		OrderNo     string  `json:"orderNo"`
		MerOrderNo  string  `json:"merOrderNo"`
		Amount      float64 `json:"amount"`
		Currency    string  `json:"currency"`
		CreateTime  int64   `json:"createTime"`
		UpdateTime  int64   `json:"updateTime"`
		Message     string  `json:"message"`
		Sign        string  `json:"sign"`
		Attach      string  `json:"attach"`
		Extra       struct {
			Name   string `json:"name"`
			Email  string `json:"email"`
			Mobile string `json:"mobile"`
		} `json:"extra"`
		Params struct {
			// 代收信息，不定义字段，可动态接收
		} `json:"params"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type WithdrawReq struct {
	AppID      string `json:"appId"`      // 应用号
	MerOrderNo string `json:"merOrderNo"` // 商户订单号
	Currency   string `json:"currency"`   // 金额币种，3位字符
	Amount     string `json:"amount"`     // 金额（字符串，保留两位小数）
	NotifyURL  string `json:"notifyUrl"`  // 异步通知地址
	Extra      struct {
		PayoutType string `json:"payoutType"` // 代付类型，固定值[IFSC]
		IFSC       string `json:"ifsc"`       // IFSC，代付类型为 IFSC 时必填
		Account    string `json:"account"`    // 收款账号
		Name       string `json:"name"`       // 客户名称
		Email      string `json:"email"`      // 客户邮箱
		Mobile     string `json:"mobile"`     // 客户手机
	} `json:"extra"` // 扩展信息
	// Attach string `json:"attach"` // 附加信息，原样返回
	Sign string `json:"sign"` // 签名
}

type WithdrawResp struct {
	Code int `json:"code"` // 响应状态码，0为成功
	Data struct {
		OrderStatus int     `json:"orderStatus"` // 订单状态: 支付中 0 1 -4, 支付成功 2 3, 支付失败 -1 -2 -3
		OrderNo     string  `json:"orderNo"`     // 交易订单号
		MerOrderNo  string  `json:"merOrderNo"`  // 商户订单号
		Amount      float64 `json:"amount"`      // 订单金额
		Currency    string  `json:"currency"`    // 金额币种
		CreateTime  int64   `json:"createTime"`  // 创建时间
		UpdateTime  int64   `json:"updateTime"`  // 更新时间
		Message     string  `json:"message"`     // 描述信息（订单异常时返回）
		Sign        string  `json:"sign"`        // 签名
		Attach      string  `json:"attach"`      // 附加信息
		Extra       struct {
			Msg string `json:"msg"` // 返回文字描述
		} `json:"extra"` // 客户信息
	} `json:"data"`
}

type WithdrawCallbackReq struct {
	OrderStatus int     `json:"orderStatus"` // 订单状态：支付中 0 1 -4, 支付成功 2 3, 支付失败 -1 -2 -3
	OrderNo     string  `json:"orderNo"`     // 交易订单号
	MerOrderNo  string  `json:"merOrderNo"`  // 商户订单号
	Amount      float64 `json:"amount"`      // 订单金额
	Currency    string  `json:"currency"`    // 金额币种
	Attach      string  `json:"attach"`      // 附加信息（商户附加信息，原样返回）
	ReceiptURL  string  `json:"receiptUrl"`  // 凭证信息（动态数据，可能为空）
	CreateTime  int64   `json:"createTime"`  // 创建时间
	UpdateTime  int64   `json:"updateTime"`  // 更新时间
	Message     string  `json:"message"`     // 描述信息（订单异常时返回）
	Sign        string  `json:"sign"`        // 签名
}

type QueryWithdrawReq struct {
	AppID      string `json:"appId"`      // 应用号
	MerOrderNo string `json:"merOrderNo"` // 商户订单号（可选，orderNo 二选一）
	OrderNo    string `json:"orderNo"`    // 交易订单号（可选，merOrderNo 二选一）
	Sign       string `json:"sign"`       // 签名
}

type QueryWithdrawResp struct {
	Code int `json:"code"` // 响应状态码，必需
	Data struct {
		OrderStatus int     `json:"orderStatus"` // 订单状态: 支付中 0 1 -4, 支付成功 2 3, 支付失败 -1 -2 -3
		OrderNo     string  `json:"orderNo"`     // 交易订单号
		MerOrderNo  string  `json:"merOrderNo"`  // 商户订单号
		Amount      float64 `json:"amount"`      // 金额
		Currency    string  `json:"currency"`    // 金额币种
		EndToEndID  string  `json:"endToEndId"`  // UTR，可选
		CreateTime  int64   `json:"createTime"`  // 创建时间
		UpdateTime  int64   `json:"updateTime"`  // 更新时间
		Message     string  `json:"message"`     // 描述信息（订单异常时返回）
		Sign        string  `json:"sign"`        // 签名
		Extra       struct {
			Msg string `json:"msg"` // 返回文字描述
		} `json:"extra"` // 客户信息
	} `json:"data"`
}
