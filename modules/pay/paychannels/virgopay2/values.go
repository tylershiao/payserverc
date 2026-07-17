package virgopay2

const (
	payURL           = "/api/v1/payIn/create"
	queryPayURL      = "/api/v1/payIn/query"
	withdrawURL      = "/api/v1/payOut/create"
	queryWithdrawURL = "/api/v1/payOut/query"
)

type PayReq struct {
	MchId     int    `json:"MchId"`
	MOrderId  string `json:"MOrderId"`
	UserId    string `json:"UserId"`
	Amount    int64  `json:"Amount"`
	IP        string `json:"IP"`
	Name      string `json:"Name"`
	Mobile    string `json:"Mobile"`
	Email     string `json:"Email"`
	NotifyURL string `json:"NotifyURL"`
	DeviceID  string `json:"DeviceID"`
	Sign      string `json:"Sign"`
	Extra     string `json:"Extra"`
}

type PayResp struct {
	Code int `json:"code"` // 响应代码
	Data struct {
		MchId    int    // 商户号
		OrderId  string // 平台订单号
		MOrderId string // 商户订单号
		PayURL   string // 支付链接
	} `json:"data"` // 响应数据
	Message string `json:"message"` // 响应消息
}

type PayCallbackReq struct {
	MchId    int    `json:"MchId"`
	OrderId  string `json:"OrderId"`  // 平台订单号
	MOrderId string `json:"MOrderId"` // 商户订单号
	Amount   int64  `json:"Amount"`   // 金额（单位：最小币种，如分）
	Status   int    `json:"Status"`   // 订单状态
	Message  string `json:"Message"`  // 错误信息
	Extra    string `json:"Extra"`
	Sign     string `json:"Sign"` // 签名
}

type WithdrawReq struct {
	MchId     int    `json:"MchId"`
	MOrderId  string `json:"MOrderId"`
	UserId    string `json:"UserId"`
	Amount    int64  `json:"Amount"`
	IP        string `json:"IP"`
	Name      string `json:"Name"`
	Mobile    string `json:"Mobile"`
	Email     string `json:"Email"`
	DeviceID  string `json:"DeviceID"`
	PayType   int    `json:"PayType"` // 代付方式 1bank 2upi
	Account   string `json:"Account"`
	IFSC      string `json:"IFSC"`     // PayType为1时必填
	BankName  string `json:"BankName"` // PayType为1时必填
	NotifyURL string `json:"NotifyURL"`
	Extra     string `json:"Extra"`
	Sign      string `json:"Sign"`
}

type WithdrawResp struct {
	Code    int    `json:"code"`    // 状态码，0 表示提单成功
	Message string `json:"message"` // 状态信息，"SUCCESS" 表示成功，"FAIL" 表示失败
	Data    struct {
		MchId    int    // 商户号
		OrderId  string // 平台订单号
		MOrderId string // 商户订单号
	} `json:"data"` // 数据字段
}

type WithdrawCallbackReq struct {
	MchId    int    `json:"MchId"`
	OrderId  string `json:"OrderId"`  // 平台订单号
	MOrderId string `json:"MOrderId"` // 商户订单号
	Amount   int64  `json:"Amount"`   // 金额（单位：最小币种，如分）
	Status   int    `json:"Status"`   // 订单状态
	Message  string `json:"Message"`  // 错误信息
	Extra    string `json:"Extra"`
	Utr      string `json:"Utr"`
	Sign     string `json:"Sign"` // 签名
}

type QueryWithdrawReq struct {
	MchId    int    `json:"MchId"`
	MOrderId string `json:"MOrderId"`
	Sign     string `json:"Sign"`
}

type QueryWithdrawResp struct {
	Code    int    `json:"code"`    // 成功：0
	Message string `json:"message"` // 响应消息，成功：SUCCESS
	Data    struct {
		MchId    int    // 商户号
		OrderId  string // 平台订单号
		MOrderId string // 商户订单号
		Amount   int64  // 订单金额
		Status   int    // 订单状态
	} `json:"data"` // 回调通知数据
}

type QueryPayReq struct {
	MchId    int    `json:"MchId"`
	MOrderId string `json:"MOrderId"`
	Sign     string `json:"Sign"`
}

type QueryPayResp struct {
	Code    int    `json:"code"`    // 成功：0，失败：其他
	Message string `json:"message"` // 响应消息，成功：SUCCESS
	Data    struct {
		MchId    int    // 商户号
		OrderId  string // 平台订单号
		MOrderId string // 商户订单号
		Amount   int64  // 订单金额
		Status   int    // 订单状态
	} `json:"data"` // 回调通知数据
}
