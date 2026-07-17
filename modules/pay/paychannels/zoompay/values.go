package zoompay

const (
	baseURL          = "https://api.zoom-pay.vip"
	payURL           = "/api/pay/order/payment/create"
	queryPayURL      = "/api/pay/order/query"
	withdrawURL      = "/api/pay/order/payout/create"
	queryWithdrawURL = "/api/pay/order/query"
	balURL           = "/api/pay/mch/balance/query"
	mid              = 1999
	key              = "cde43c6fa9c8421481c8d2ccc374bc7b"
)

type PayReq struct {
	MchId          int64  `json:"mchId"`          // 商户ID（必选）
	MchOrderNo     string `json:"mchOrderNo"`     // 商户订单号（必选）
	ProductId      int64  `json:"productId"`      // 支付产品ID（必选）
	OrderAmount    int64  `json:"orderAmount"`    // 订单金额，单位分（必选）
	NotifyUrl      string `json:"notifyUrl"`      // 回调地址（必选）
	CustomerName   string `json:"customerName"`   // 玩家姓名（必选）
	CustomerMobile string `json:"customerMobile"` // 电话号码 10 位（必选）
	CustomerEmail  string `json:"customerEmail"`  // 邮箱地址（必选）
	// Remark         string `json:"remark"`         // 附加参数（可选）
	Sign string `json:"sign"` // 签名（必选）
}

type PayResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		PayOrderId  string `json:"payOrderId"`  // 平台单号
		OrderAmount int64  `json:"orderAmount"` // 订单金额
		PayUrl      string `json:"payUrl"`      // 支付链接
	} `json:"data"`
}

type PayCallbackReq struct {
	Code           int64  `json:"code"`           // 回调结果代码 1=成功 2=驳回
	MchId          int64  `json:"mchId"`          // 商户ID
	MchOrderNo     string `json:"mchOrderNo"`     // 商户订单号
	ProductId      int64  `json:"productId"`      // 支付产品ID
	OrderAmount    int64  `json:"orderAmount"`    // 订单金额(单位分)
	PayOrderId     string `json:"payOrderId"`     // 平台订单号
	PaySuccessTime string `json:"paySuccessTime"` // 成功时间戳(毫秒) code=1 时才有
	Message        string `json:"message"`        // 成功固定 success，失败返回原因
	Remark         string `json:"remark"`         // 备注（可空）
	Sign           string `json:"sign"`           // 签名
	Utr            string `json:"utr"`            // UTR，可空
}

type QueryPayReq struct {
	MchId      int64  `json:"mchId"`      // 商户ID
	MchOrderNo string `json:"mchOrderNo"` // 商户系统订单号
	Sign       string `json:"sign"`       // 签名
}

type QueryPayResp struct {
	Code int `json:"code"`
	Data struct {
		Msg         string `json:"msg"`
		PayOrderId  string `json:"payOrderId"`
		MchId       int64  `json:"mchId"`
		OrderAmount int64  `json:"orderAmount"`
		ProductId   int64  `json:"productId"`
		CreateTime  string `json:"createTime"`
		MchOrderNo  string `json:"mchOrderNo"`
		PaySuccTime int64  `json:"paySuccTime"`
		Status      string `json:"status"`
	} `json:"data"`
	Message string `json:"message"`
}

type WithdrawReq struct {
	MchId          int64  `json:"mchId"`          // 商户ID
	MchOrderNo     string `json:"mchOrderNo"`     // 商户订单号
	ProductId      int64  `json:"productId"`      // 支付产品ID
	OrderAmount    int64  `json:"orderAmount"`    // 订单金额（分）
	NotifyUrl      string `json:"notifyUrl"`      // 回调地址
	CustomerName   string `json:"customerName"`   // 玩家姓名
	CustomerMobile string `json:"customerMobile"` // 电话号码
	CustomerEmail  string `json:"customerEmail"`  // 邮箱
	AccountName    string `json:"accountName"`    // 银行开户姓名
	CardNumber     string `json:"cardNumber"`     // 银行卡号
	IFSC           string `json:"ifsc"`           // 印度银行编码
	// Remark         string `json:"remark"`         // 备注(可选)
	BankName string `json:"bankName"` // 银行名称(可选)
	Sign     string `json:"sign"`     // 签名
}

type WithdrawResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		PayOrderId  string `json:"payOrderId"`  // 平台单号
		OrderAmount int64  `json:"orderAmount"` // 订单金额
	} `json:"data"`
}

type WithdrawCallbackReq struct {
	Code           int64  `json:"code"`           // 回调结果代码 1=成功 2=驳回
	MchId          int64  `json:"mchId"`          // 商户ID
	MchOrderNo     string `json:"mchOrderNo"`     // 商户订单号
	ProductId      int64  `json:"productId"`      // 支付产品ID
	OrderAmount    int64  `json:"orderAmount"`    // 订单金额(单位分)
	PayOrderId     string `json:"payOrderId"`     // 平台订单号
	PaySuccessTime string `json:"paySuccessTime"` // 成功时间戳(毫秒) code=1 时才有
	Message        string `json:"message"`        // 成功固定 success，失败返回原因
	Remark         string `json:"remark"`         // 备注（可空）
	Sign           string `json:"sign"`           // 签名
	Utr            string `json:"utr"`            // UTR，可空
}

type QueryWithdrawReq struct {
	MchId      int64  `json:"mchId"`      // 商户ID
	MchOrderNo string `json:"mchOrderNo"` // 商户系统订单号
	Sign       string `json:"sign"`       // 签名
}

type QueryWithdrawResp struct {
	Code int `json:"code"`
	Data struct {
		Msg         string `json:"msg"`
		PayOrderId  string `json:"payOrderId"`
		MchId       int64  `json:"mchId"`
		OrderAmount int64  `json:"orderAmount"`
		ProductId   int64  `json:"productId"`
		CreateTime  string `json:"createTime"`
		MchOrderNo  string `json:"mchOrderNo"`
		PaySuccTime int64  `json:"paySuccTime"`
		Status      string `json:"status"`
	} `json:"data"`
	Message string `json:"message"`
}
