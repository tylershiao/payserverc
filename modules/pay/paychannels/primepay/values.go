package primepay

const (
	baseURL          = "https://payapi.primecash.net/api"
	payURL           = "/payin/create"
	queryPayURL      = "/payin/order/query"
	withdrawURL      = "/payout/create"
	queryWithdrawURL = "/payout/order/query"
	balURL           = "/exposed/v1/query/balance"
	mid              = "M231182"
	key              = "lhYWkwdZDdYrvoV"
	secret           = "WrYRaac9YGQUUay9Kv1u3DoqzNE1sM"
)

type PayReq struct {
	MerchantOrderId     string `json:"merchantOrderId"`     // 商户订单号，不能重复，≤40位
	Amount              string `json:"amount"`              // 订单金额，浮点数字符串，保留两位小数
	Currency            string `json:"currency"`            // 货币种类，三位大写字母，如 INR
	NotificationURL     string `json:"notificationURL"`     // 支付成功通知地址
	MerchantRedirectUrl string `json:"merchantRedirectUrl"` // 支付后跳转链接
	TransType           string `json:"transType"`           // 收银台类型：1=NATIVE，2=WAKEUP
	UUID                string `json:"uuid"`                // 客户在商户系统的账户ID
	IPAddress           string `json:"ipAddress"`           // 客户真实IP
}

type PayResp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		PrimecashOrder string  `json:"primecashOrder"` // primecash系统订单号
		MerchantOrder  string  `json:"merchantOrder"`  // 商户订单号
		OrderStatus    string  `json:"orderStatus"`    // 订单状态，如 PENDING
		PayURL         string  `json:"payUrl"`         // 支付跳转链接
		Currency       string  `json:"currency"`       // 交易币种，如 INR
		Amount         float64 `json:"amount"`         // 支付金额
	} `json:"data"`
}

type PayCallbackReq struct {
	MerchantNo     string  `json:"merchantNo"`     // 商户编号
	MerchantName   string  `json:"merchantName"`   // 商户名称
	RespCode       string  `json:"respCode"`       // 响应码
	RespMsg        string  `json:"respMsg"`        // 响应码说明
	PrimecashOrder string  `json:"primecashOrder"` // PRIMECASH订单号
	MerchantOrder  string  `json:"merchantOrder"`  // 商户订单号
	Amount         float64 `json:"amount"`         // 订单金额
	Fee            float64 `json:"fee"`            // 订单手续费
	Currency       string  `json:"currency"`       // 订单币种
	GmtCreate      string  `json:"gmtCreate"`      // 订单创建时间 yyyy-MM-dd HH:mm:ss
	GmtPay         string  `json:"gmtPay"`         // 订单支付时间 yyyy-MM-dd HH:mm:ss
	OrderStatus    string  `json:"orderStatus"`    // 订单状态，如 SUCCESSFUL
	TransType      string  `json:"transType"`      // 交易类型，如 NATIVE
	Sign           string  `json:"sign"`           // 签名
}

type QueryPayReq struct {
}

type QueryPayResp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		PrimecashOrder string  `json:"primecashOrder"` // PRIMECASH订单号
		OrderStatus    string  `json:"orderStatus"`    // 订单状态，如 SUCCESSFUL
		PayURL         string  `json:"payUrl"`         // 支付链接
		Currency       string  `json:"currency"`       // 订单币种，如 SAR、AED
		Amount         float64 `json:"amount"`         // 订单金额
		Fee            float64 `json:"fee"`            // 订单手续费
		CreateTime     string  `json:"createTime"`     // 订单创建时间 yyyy-MM-dd HH:mm:ss
		PayTime        string  `json:"payTime"`        // 订单支付时间 yyyy-MM-dd HH:mm:ss
		TransType      string  `json:"transType"`      // 支付类型，如 NATIVE
		MerchantOrder  string  `json:"merchantOrder"`  // 商户订单号
	} `json:"data"`
}

type WithdrawReq struct {
	MerchantOrderId string  `json:"merchantOrderId"` // 商户代付订单号
	PayType         string  `json:"payType"`         // 付款类别：1=一类原生，2=一类唤醒，3=二类唤醒
	Amount          float64 `json:"amount"`          // 订单金额
	MerchantURL     string  `json:"merchanturl"`     // 付款成功异步回调地址
	AccNumber       string  `json:"accNumber"`       // 收款账号
	AccName         string  `json:"accName"`         // 收款姓名
	IFSC            string  `json:"ifsc"`            // IFSC 编码
	IPAddress       string  `json:"ipAddress"`       // 客户真实 IP
	UUID            string  `json:"uuid"`            // 客户在商户系统的账户ID
}

type WithdrawResp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		PrimecashOrder string  `json:"primecashOrder"` // PrimeCash 付款订单号
		MerchantOrder  string  `json:"merchantOrder"`  // 商户付款订单号
		Amount         float64 `json:"amount"`         // 金额
		Fee            float64 `json:"fee"`            // 手续费
		CreateAt       string  `json:"createAt"`       // 订单创建时间 yyyy-MM-dd HH:mm:ss
	} `json:"data"`
}

type WithdrawCallbackReq struct {
	PrimecashOrder  string  `json:"primecashOrder"`  // PrimeCash 付款订单号
	MerchantOrder   string  `json:"merchantOrder"`   // 商户付款订单号
	Fee             float64 `json:"fee"`             // 手续费
	Amount          float64 `json:"amount"`          // 订单金额
	UTR             string  `json:"utr"`             // 付款凭证，测试环境固定返回 MOCK_UTR_FOR_TESTING
	PayeeStatus     int     `json:"payeeStatus"`     // 付款状态码
	PayeeStatusName string  `json:"payeeStatusName"` // 状态描述，如 submit
	ConfirmTime     string  `json:"confirmtime"`     // 打款完成时间，Unix 时间戳
	CreatedAt       string  `json:"createdAt"`       // 创建时间 yyyy-MM-dd HH:mm:ss
	UpdatedAt       string  `json:"updatedAt"`       // 更新时间 yyyy-MM-dd HH:mm:ss
	Sign            string  `json:"sign"`            // 签名
}

type QueryWithdrawReq struct {
	// AppID     string `json:"appId"`     // 应用ID，创建商户后提供（必选）
	// Timestamp string `json:"timestamp"` // Unix时间戳（毫秒13位），与服务器时间误差不可超过5分钟（必选）
	// Sign      string `json:"sign"`      // 请求签名（必选）

	// OrderNo string `json:"orderNo"` // 商户订单号（必选）
}

type QueryWithdrawResp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Amount         float64 `json:"amount"`         // 付款金额
		Fee            float64 `json:"fee"`            // 手续费
		PrimecashOrder string  `json:"primecashOrder"` // PrimeCash 付款订单号
		PaymentMethod  string  `json:"paymentmethod"`  // 付款方式：NATIVE, WAKEUP
		CreateAt       string  `json:"createAt"`       // 创建时间 yyyy-MM-dd HH:mm:ss
		Status         int     `json:"status"`         // 付款状态 -1=Initialize, 0=Pending, 1=Transferred, 2=Error
		ConfirmTime    string  `json:"confirmtime"`    // 打款完成时间 yyyy-MM-dd HH:mm:ss
		UTR            string  `json:"utr"`            // 付款凭证
	} `json:"data"`
}
