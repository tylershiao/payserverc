package eastpay

const (
	baseURL          = "https://api.eastpay-app.com"
	payURL           = "/api/pay/create/order"
	queryPayURL      = "/api/pay/query/order"
	withdrawURL      = "/api/pay/create/order"
	queryWithdrawURL = "/api/pay/query/order"
	balURL           = "/api/pay/merchant/balance"
	mid              = 249
	// appID = "S1768487216092"
	key = "b3eae4bc4ed834174441707cb2d12a01"
)

type PayReq struct {
	Type            int     `json:"type"`            // 业务类型：代收-0，代付-1（必填）
	MerchantID      int64   `json:"merchantId"`      // 商户ID（数字）（必填）
	MerchantOrderNo string  `json:"merchantOrderNo"` // 商户订单号（必填）
	OrderAmount     float64 `json:"orderAmount"`     // 订单金额（必填，1:1传值）
	NotifyURL       string  `json:"notifyUrl"`       // 回调地址（必填）
	// ReturnURL       string  `json:"returnUrl"`       // 页面跳转地址（可选）
	// Extra           string  `json:"extra"`           // 额外信息（可选，全流程透传）
	Sign string `json:"sign"` // 签名（必填）
}

type PayResp struct {
	Code    int    `json:"code"`    // 状态码（200 表示成功）
	Message string `json:"message"` // success
	Data    struct {
		OrderNo     string  `json:"orderNo"`     // 平台订单号
		OrderAmount float64 `json:"orderAmount"` // 金额
		PayURL      string  `json:"payUrl"`      // 收银台地址
		UPI         string  `json:"upi"`         // UPI
		Deeplink    string  `json:"deeplink"`    // deeplink
	} `json:"data"` // 数据对象
}

type PayCallbackReq struct {
	Code       int    `json:"code"`       // 1=成功(支付成功/代付成功)，2=驳回(代付驳回)
	MchID      int64  `json:"mchId"`      // 商户ID（数字）
	MchOrderNo string `json:"mchOrderNo"` // 商户订单号
	// OrderAmount    float64 `json:"orderAmount"`    // 订单金额（保留2位小数）
	// UTR            string `json:"utr"`            // 代付完成UTR（代收不返回）
	OrderNo        string `json:"orderNo"`        // 平台订单号
	PaySuccessTime string `json:"paySuccessTime"` // 支付成功时间（字符串时间）
	Message        string `json:"message"`        // success / 失败原因
	Extra          string `json:"extra"`          // 扩展字段（透传）
	Sign           string `json:"sign"`           // 签名
}

type QueryPayReq struct {
	Type            int    `json:"type"`            // 业务类型：0=代收，1=代付
	MerchantID      int64  `json:"merchantId"`      // 商户ID（纯数字）
	MerchantOrderNo string `json:"merchantOrderNo"` // 商户订单号
	Sign            string `json:"sign"`            // 签名
}

type QueryPayResp struct {
	Code    int    `json:"code"`    // 状态码（200成功）
	Message string `json:"message"` // 描述信息
	Data    struct {
		MerchantID      int64   `json:"merchantId"`      // 商户ID
		OrderAmount     float64 `json:"orderAmount"`     // 订单金额
		MerchantOrderNo string  `json:"merchantOrderNo"` // 商户订单号
		Status          int     `json:"status"`          // 订单状态
		CreateTime      string  `json:"createTime"`      // 创建时间
		PayTime         string  `json:"payTime"`         // 支付时间
		NotifyTime      string  `json:"notifyTime"`      // 回调时间
	} `json:"data"` // 数据对象
}

type WithdrawReq struct {
	Type            int     `json:"type"`            // 业务类型：0=代收，1=代付（这里固定为1）
	MerchantID      int64   `json:"merchantId"`      // 商户ID（纯数字）
	MerchantOrderNo string  `json:"merchantOrderNo"` // 商户订单号
	OrderAmount     float64 `json:"orderAmount"`     // 订单金额（1:1，不换算）
	AccountName     string  `json:"accountName"`     // 收款银行账户名称
	CardNumber      string  `json:"cardNumber"`      // 收款银行卡号
	IFSC            string  `json:"ifsc"`            // 收款银行IFSC
	BankName        string  `json:"bankName"`        // 收款银行名称
	UPI             string  `json:"upi,omitempty"`   // 收款方UPI（可选）
	NotifyURL       string  `json:"notifyUrl"`       // 回调地址
	Phone           string  `json:"phone"`           // 收款人手机号（可选）
	Email           string  `json:"email"`           // 收款人邮箱（可选）
	// Extra           string  `json:"extra"`           // 扩展字段（透传）
	Sign string `json:"sign"` // 签名
}

type WithdrawResp struct {
	Code    int    `json:"code"`    // 状态码
	Message string `json:"message"` // 描述信息
	Data    struct {
		OrderNo     string  `json:"orderNo"`     // 平台订单号
		OrderAmount float64 `json:"orderAmount"` // 订单金额
	} `json:"data"` // 数据对象
}

type WithdrawCallbackReq struct {
	Code       int    `json:"code"`       // 1=成功(支付成功/代付成功)，2=驳回(代付驳回)
	MchID      int64  `json:"mchId"`      // 商户ID（数字）
	MchOrderNo string `json:"mchOrderNo"` // 商户订单号
	// OrderAmount    float64 `json:"orderAmount"`    // 订单金额（保留2位小数）
	UTR            string `json:"utr"`            // 代付完成UTR（代收不返回）
	OrderNo        string `json:"orderNo"`        // 平台订单号
	PaySuccessTime string `json:"paySuccessTime"` // 支付成功时间（字符串时间）
	Message        string `json:"message"`        // success / 失败原因
	Extra          string `json:"extra"`          // 扩展字段（透传）
	Sign           string `json:"sign"`           // 签名
}

type QueryWithdrawReq struct {
	Type            int    `json:"type"`            // 业务类型：0=代收，1=代付
	MerchantID      int64  `json:"merchantId"`      // 商户ID（纯数字）
	MerchantOrderNo string `json:"merchantOrderNo"` // 商户订单号
	Sign            string `json:"sign"`            // 签名
}

type QueryWithdrawResp struct {
	Code    int    `json:"code"`    // 状态码（200成功）
	Message string `json:"message"` // 描述信息
	Data    struct {
		MerchantID      int64   `json:"merchantId"`      // 商户ID
		OrderAmount     float64 `json:"orderAmount"`     // 订单金额
		MerchantOrderNo string  `json:"merchantOrderNo"` // 商户订单号
		Status          int     `json:"status"`          // 订单状态
		CreateTime      string  `json:"createTime"`      // 创建时间
		PayTime         string  `json:"payTime"`         // 支付时间
		NotifyTime      string  `json:"notifyTime"`      // 回调时间
	} `json:"data"` // 数据对象
}
