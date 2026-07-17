package kingpay

const (
	baseURL          = "https://gateway.king-pay.net/service"
	payURL           = "/order/payin/create"
	queryPayURL      = "/order/payin/query"
	withdrawURL      = "/order/payout/create"
	queryWithdrawURL = "/order/payout/query"
	// appID            = "309582512"
	// mid   = "202501243500001"
	KeyID = "5OS315dXk22Vwwlb"
	key   = "1a62301374de4951b4c5161c120c2824"
)

type PayReq struct {
	Amount      string `json:"amount"`
	DownOrderNo string `json:"downOrderNo"`
	KeyId       string `json:"keyId"`
	NotifyUrl   string `json:"notifyUrl"`
	Sign        string `json:"sign"`
}

type PayResp struct {
	Code int    `json:"code"` // 响应代码
	Msg  string `json:"msg"`  // 请求失败的消息，失败时有值
	Data struct {
		OrderNo string `json:"orderNo"` // 平台唯一订单号
		PayUrl  string `json:"payUrl"`  // 收银台地址
	} `json:"data"` // 成功时返回的数据，可能是任意类型
}

type PayCallbackReq struct {
	Amount      string `json:"amount"`            // 订单金额
	DownOrderNo string `json:"downOrderNo"`       // 商户唯一订单号
	OrderNo     string `json:"orderNo"`           // 平台订单号
	OrderStatus int    `json:"orderStatus"`       // 订单支付状态：1=支付成功，2=支付失败
	PayTime     int    `json:"payTime,omitempty"` // 支付时间（十位时间戳），失败时可能不返回
	Sign        string `json:"sign"`              // 签名
}

type QueryPayReq struct {
	DownOrderNo string `json:"downOrderNo"` // 商户唯一订单号
	KeyId       string `json:"keyId"`       // 商户key
	Sign        string `json:"sign"`        // 签名
}

type QueryPayResp struct {
	Code int    `json:"code"` // 响应代码
	Msg  string `json:"msg"`  // 请求失败的消息，失败时有值
	Data struct {
		OrderNo        string `json:"orderNo"`        // 商户唯一订单号
		PaymentOrderNo string `json:"paymentOrderNo"` // 平台唯一订单号
		CreateTime     string `json:"createTime"`     // 订单创建时间
		OrderStatus    int    `json:"orderStatus"`    // 支付状态：0=未支付，1=已支付，2=支付失败
		PayTime        string `json:"payTime"`        // 订单支付时间
	} `json:"data"` // 成功时返回的数据，可能是任意类型
}

// 定义主请求结构体
type WithdrawReq struct {
	AccountName string `json:"accountName"` // 账号名称
	AccountNo   string `json:"accountNo"`   // 银行账号
	Ifsc        string `json:"ifsc"`        // 银行IFSC标识符
	DownOrderNo string `json:"downOrderNo"` // 商户唯一订单号
	Amount      string `json:"amount"`      // 订单金额，支持两位小数
	KeyId       string `json:"keyId"`       // 商户key
	NotifyUrl   string `json:"notifyUrl"`   // 订单回调地址
	Sign        string `json:"sign"`        // 签名
}

type WithdrawResp struct {
	Code int    `json:"code"` // 响应代码
	Msg  string `json:"msg"`  // 请求失败的消息，失败时有值
	Data struct {
		OrderNo string `json:"orderNo"` // 订单支付时间
	} `json:"data"` // 成功时返回的数据，可能是任意类型
}

type WithdrawCallback struct {
	Amount      string `json:"amount"`      // 订单金额
	DownOrderNo string `json:"downOrderNo"` // 商户唯一订单号
	OrderNo     string `json:"orderNo"`     // 平台订单号
	OrderStatus int    `json:"orderStatus"` // 订单代付状态：1=代付成功，2=代付失败
	Sign        string `json:"sign"`        // 签名
}

type QueryWithdrawReq struct {
	DownOrderNo string `json:"downOrderNo"` // 商户唯一订单号
	KeyId       string `json:"keyId"`       // 商户key
	Sign        string `json:"sign"`        // 签名
}

type QueryWithdrawResp struct {
	Code int    `json:"code"` // 响应代码
	Msg  string `json:"msg"`  // 请求失败的消息，失败时有值
	Data struct {
		OrderNo        string `json:"orderNo"`        // 商户唯一订单号
		PaymentOrderNo string `json:"paymentOrderNo"` // 平台唯一订单号
		Desc           string `json:"desc"`           // 描述信息
		OrderStatus    int    `json:"orderStatus"`    // 代付状态：1=处理中，2=代付成功，3=代付失败
	} `json:"data"` // 成功时返回的数据，可能是任意类型
}
