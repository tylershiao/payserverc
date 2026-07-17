package hdpay

const (
	baseURL          = "https://dd1688.cc"
	payURL           = "/api/payin/submit"
	queryPayURL      = "/api/payin/status"
	withdrawURL      = "/api/payout/submit"
	queryWithdrawURL = "/api/payout/query"
	balURL           = "/api/payout/balance"
	mid              = 1090
	key              = "DsXopIVzVvdQCRwlV1nYxByd3Z02x4K8"
)

type PayReq struct {
	MerchantID      int32  `json:"merchantId"`      // 商户号
	MerchantOrderID string `json:"merchantOrderId"` // 商户订单
	Amount          string `json:"amount"`          // 金额
	NotifyURL       string `json:"notifyUrl"`       // 通知地址
	Name            string `json:"name"`            // 用户姓名
	Mobile          string `json:"mobile"`          // 用户电话
	Email           string `json:"email"`           // 邮箱
	Sign            string `json:"sign"`            // 请求签名
	// Deeplink        bool   `json:"deeplink"`        // 是否返回 UPI deeplink
	// Mode            string `json:"mode"`            // UPI通道模式 fz|ptm
}

type PayResp struct {
	Code int32 `json:"code"` // 状态码
	Data struct {
		PayURL    string `json:"payUrl"`    // 支付链接
		OrderID   string `json:"orderId"`   // 平台订单号
		Message   string `json:"message"`   // 信息
		Status    string `json:"status"`    // 订单状态 0等待 1成功 2失败
		ChannelID int32  `json:"channelId"` // 通道
		Deeplink  string `json:"deeplink"`  // UPI deeplink
	} `json:"data"` // 数据
	Msg string `json:"msg"` // 返回消息
}

type PayCallbackReq struct {
	MerchantID      int32  `json:"merchantId"`      // 商户
	MerchantOrderID string `json:"merchantOrderId"` // 商户订单号
	Amount          string `json:"amount"`          // 订单金额
	PayAmount       string `json:"payAmount"`       // 实付金额
	OrderID         string `json:"orderId"`         // 平台订单
	Status          string `json:"status"`          // 订单状态 0等待 1成功 2失败
	Msg             string `json:"msg"`             // 错误消息
	PayTime         string `json:"payTime"`         // 支付时间
	Sign            string `json:"sign"`            // 签名
}

type QueryPayReq struct {
	MerchantID      int32  `json:"merchantId"`      // 商户号
	MerchantOrderID string `json:"merchantOrderId"` // 商户订单
	Sign            string `json:"sign"`            // 请求签名
}

type QueryPayResp struct {
	Code int32 `json:"code"` // 状态码
	Data struct {
		OrderID         string `json:"orderId"`         // 平台订单号
		MerchantOrderID string `json:"merchantOrderId"` // 商户订单号
		Message         string `json:"message"`         // 消息
		Status          string `json:"status"`          // 订单状态 0待支付 1已支付 2已失效
	} `json:"data"` // 数据
	Msg string `json:"msg"` // 返回消息
}

type WithdrawReq struct {
	MerchantID       int32  `json:"merchantId"`        // 商户号
	MerchantPayoutID string `json:"merchantPayoutId"`  // 商户订单号
	Amount           string `json:"amount"`            // 金额
	NotifyURL        string `json:"notifyUrl"`         // 通知地址
	Name             string `json:"name"`              // 姓名
	Type             string `json:"type"`              // 出款类型 0银行卡 1UPI
	Account          string `json:"account,omitempty"` // 卡号 type=0 必填
	IFSC             string `json:"ifsc,omitempty"`    // ifsc type=0 必填
	UPI              string `json:"upi,omitempty"`     // upi type=1 必填
	Sign             string `json:"sign"`              // 请求签名
}

type WithdrawResp struct {
	Code int32 `json:"code"` // 状态码
	Data struct {
		PayoutID         string `json:"payoutId"`         // 平台订单号
		MerchantPayoutID string `json:"merchantPayoutId"` // 商户订单号
		Type             string `json:"type"`             // 用户提现类型
		Status           string `json:"status"`           // 提现状态 0等待 1成功 2失败
		Msg              string `json:"msg"`              // 错误消息
		Amount           string `json:"amount"`           // 订单金额
		UTR              string `json:"utr"`              // utr
	} `json:"data"` // 数据
	Msg string `json:"msg"` // 返回消息
}

type WithdrawCallbackReq struct {
	MerchantID       int32  `json:"merchantId"`       // 商户
	MerchantPayoutID string `json:"merchantPayoutId"` // 商户订单
	Amount           string `json:"amount"`           // 金额
	PayoutID         string `json:"payoutId"`         // 平台订单号
	Status           string `json:"status"`           // 订单状态 1成功 2失败
	Msg              string `json:"msg"`              // 错误消息
	Fee              string `json:"fee"`              // 手续费
	SingleFee        string `json:"singleFee"`        // 单笔
	UTR              string `json:"utr"`              // utr
	Sign             string `json:"sign"`             // 签名
}

type QueryWithdrawReq struct {
	MerchantID       int32  `json:"merchantId"`       // 商户id
	MerchantPayoutID string `json:"merchantPayoutId"` // 商户订单id
	Sign             string `json:"sign"`             // 请求签名
}

type QueryWithdrawResp struct {
	Code int32 `json:"code"` // 状态码
	Data struct {
		PayoutID         string `json:"payoutId"`         // 平台订单号
		MerchantPayoutID string `json:"merchantPayoutId"` // 商户订单号
		Type             string `json:"type"`             // 用户提现类型
		Status           string `json:"status"`           // 提现状态 0等待 1成功 2失败
		Msg              string `json:"msg"`              // 错误消息
		Amount           string `json:"amount"`           // 订单金额
		UTR              string `json:"utr"`              // utr
	} `json:"data"` // 数据
	Msg string `json:"msg"` // 返回消息
}
