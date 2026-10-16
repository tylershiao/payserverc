package masterpay

const (
	baseURL          = "https://api.masterpay88.in/app-api"
	payURL           = "/v1.0/api/order/create"
	queryPayURL      = "/v1.0/api/order/query"
	withdrawURL      = "/v1.0/api/order/create"
	queryWithdrawURL = "/v1.0/api/order/query"
	mid              = "10431"
	key              = "G9BBJJTIUMNS2BCJCNHAVZI3ZLOPBMI4"
)

type PayReq struct {
	MchID      string `json:"mchId"`      // 商户ID
	ProductID  string `json:"productId"`  // 通道ID，用于区分代收和代付
	MchOrderNo string `json:"mchOrderNo"` // 商户订单号
	Amount     int    `json:"amount"`     // 订单金额，单位分，不支持小数点
	ClientIP   string `json:"clientIp"`   // 会员IP，可固定填：0.0.0.0
	NotifyURL  string `json:"notifyUrl"`  // 支付结果后台回调URL
	Sign       string `json:"sign"`       // 签名值
	// ReturnURL        string `json:"returnUrl"`        // （可选）支付结果前端跳转URL
	// Subject          string `json:"subject"`          // （可选）商品主题
	// Body             string `json:"body"`             // （可选）商品描述信息
	// Param1           string `json:"param1"`           // （可选）扩展参数1，回调时原样返回
	// Param2           string `json:"param2"`           // （可选）扩展参数2，回调时原样返回
	// ValidateUserName string `json:"validateUserName"` // （可选）会员姓名，实名制模式需填
	RequestCardInfo bool `json:"requestCardInfo"` // （可选）是否需要提供分配的卡信息，默认为 false
}

type PayResp struct {
	RetCode    string      `json:"retCode"`    // 返回状态码，成功返回 "SUCCESS"，失败返回 "FAIL"
	RetMsg     string      `json:"retMsg"`     // 返回信息，错误原因（如有）
	Sign       string      `json:"sign"`       // 签名，仅在 retCode 为 "SUCCESS" 时返回
	PayOrderID string      `json:"payOrderId"` // 平台订单号
	PayURL     string      `json:"payUrl"`     // 支付跳转地址
	PayParams  interface{} `json:"payParams"`  // 支付参数（通常为对象或空）
	CardInfo   *CardInfo   `json:"cardInfo"`   // 卡信息对象
}

type CardInfo struct {
	BankCard         string `json:"bankCard"`         // 银行卡号
	UserName         string `json:"userName"`         // 收款人姓名
	BankName         string `json:"bankName"`         // 银行名称（可为空）
	ValidateUserName string `json:"validateUserName"` // 验证用户名（可为空）
	Amount           int    `json:"amount"`           // 金额，单位分
	TimeoutMinutes   int    `json:"timeoutMinutes"`   // 超时时间（分钟）
	Link             string `json:"link"`             // UPI 链接
}

type PayCallbackReq struct {
	PayOrderId     string `json:"payOrderId"`     // 平台订单号
	MchId          string `json:"mchId"`          // 商户ID
	ProductId      string `json:"productId"`      // 通道ID
	MchOrderNo     string `json:"mchOrderNo"`     // 商户订单号
	Amount         int64  `json:"amount"`         // 订单金额，单位分
	RealAmount     int64  `json:"realAmount"`     // 实际要求付款金额（浮动金额），单位分
	Income         int64  `json:"income"`         // 入账金额，单位分，商户入账金额（扣除手续费后）
	Status         int    `json:"status"`         // 支付状态：0-支付中，1-已完成，3-已超时，8-用户退款
	PaySuccessTime int64  `json:"paySuccessTime"` // 支付成功时间，精确到毫秒
	Sign           string `json:"sign"`           // 签名
	Param1         string `json:"param1"`         // 扩展参数1，回调时原样返回
	Param2         string `json:"param2"`         // 扩展参数2，回调时原样返回
	UTR            string `json:"utr"`            // 银行交易参考号（UTR），仅在支付成功时返回
}

type QueryPayReq struct {
	MchId      string `json:"mchId"`      // 商户ID
	MchOrderNo string `json:"mchOrderNo"` // 商户订单号
	Sign       string `json:"sign"`       // 签名
}

type QueryPayResp struct {
	RetCode        string `json:"retCode"`        // 返回状态码，成功返回 "SUCCESS"，失败返回 "FAIL"
	RetMsg         string `json:"retMsg"`         // 返回信息，错误原因（如有）
	MchId          string `json:"mchId"`          // 商户ID
	ProductId      int    `json:"productId"`      // 通道ID
	PayOrderId     string `json:"payOrderId"`     // 平台订单号
	MchOrderNo     string `json:"mchOrderNo"`     // 商户订单号
	Amount         int64  `json:"amount"`         // 订单金额，单位分
	RealAmount     int64  `json:"realAmount"`     // 实际要求付款金额，单位分
	Status         int    `json:"status"`         // 订单状态：0-支付中，1-已完成，3-已超时，5-驳回中，7-已驳回
	UTR            string `json:"utr"`            // 银行交易参考号
	PaySuccessTime int64  `json:"paySuccessTime"` // 支付成功时间，精确到秒
}

type WithdrawReq struct {
	MchID      string `json:"mchId"`      // 商户ID
	ProductID  string `json:"productId"`  // 通道ID，用于区分代收和代付
	MchOrderNo string `json:"mchOrderNo"` // 商户订单号
	Amount     int    `json:"amount"`     // 订单金额，单位分
	ClientIP   string `json:"clientIp"`   // 会员IP，可固定填：0.0.0.0
	NotifyURL  string `json:"notifyUrl"`  // 支付结果后台回调URL
	Sign       string `json:"sign"`       // 签名值
	UserName   string `json:"userName"`   // 收款人姓名
	CardNumber string `json:"cardNumber"` // 收款人银行卡号
	// PayAccountType int    `json:"payAccountType"` // 收款账号类型（墨西哥支付必填，40=CLABE）
	IfscCode string `json:"ifscCode"` // IFSC代码（印度），或银行编码（印尼）
	BankName string `json:"bankName"` // 收款人开户行名称（可为空但需传递）
	// Param1         string `json:"param1"`         // 扩展参数1，回调时原样返回
	// Param2         string `json:"param2"`         // 扩展参数2，回调时原样返回
}

type WithdrawResp struct {
	RetCode    string `json:"retCode"`    // 返回状态码：SUCCESS 表示成功，FAIL 表示失败
	RetMsg     string `json:"retMsg"`     // 返回信息，错误原因（如有）
	Sign       string `json:"sign"`       // 签名，仅在 retCode=SUCCESS 时返回
	PayOrderID string `json:"payOrderId"` // 支付订单号，仅在 retCode=SUCCESS 时返回
	Status     int    `json:"status"`     // 订单状态：0=支付中，1=已完成，3=已超时，5=驳回中，7=已驳回
}

type WithdrawCallbackReq struct {
	PayOrderID     string `json:"payOrderId"`     // 平台订单号
	MchID          string `json:"mchId"`          // 商户ID
	ProductID      string `json:"productId"`      // 通道ID
	MchOrderNo     string `json:"mchOrderNo"`     // 商户订单号
	Amount         int    `json:"amount"`         // 订单金额，单位分
	RealAmount     int    `json:"realAmount"`     // 实际要求付款金额（浮动金额），单位分
	Income         int    `json:"income"`         // 入账金额，单位分（商户入账金额，浮动后扣完手续费）
	Status         int    `json:"status"`         // 支付状态：1=成功，3=已超时，5=驳回中，7=已驳回，9=冲正（5和7视为失败）
	PaySuccessTime int64  `json:"paySuccessTime"` // 支付成功时间，精确到毫秒
	Sign           string `json:"sign"`           // 签名
	RejectReason   string `json:"rejectReason"`   // 驳回理由（仅代付订单被驳回时返回）
	Param1         string `json:"param1"`         // 扩展参数1，回调时原样返回
	Param2         string `json:"param2"`         // 扩展参数2，回调时原样返回
	UTR            string `json:"utr"`            // 银行交易参考号（UTR），仅在支付成功时返回
}

type QueryWithdrawReq struct {
	MchID      string `json:"mchId"`      // 商户ID
	MchOrderNo string `json:"mchOrderNo"` // 商户订单号
	Sign       string `json:"sign"`       // 签名
}

type QueryWithdrawResp struct {
	RetCode        string `json:"retCode"`        // 返回状态码：SUCCESS 表示成功，FAIL 表示失败
	RetMsg         string `json:"retMsg"`         // 返回信息，错误原因
	MchID          string `json:"mchId"`          // 商户ID
	ProductID      string `json:"productId"`      // 通道ID
	PayOrderID     string `json:"payOrderId"`     // 平台订单号
	MchOrderNo     string `json:"mchOrderNo"`     // 商户订单号
	Amount         int    `json:"amount"`         // 订单金额，单位分
	RealAmount     int    `json:"realAmount"`     // 实际要求付款金额，单位分
	Status         int    `json:"status"`         // 订单状态：0=支付中，1=已完成，3=已超时，5=驳回中，7=已驳回
	UTR            string `json:"utr"`            // 银行交易凭证号
	PaySuccessTime int64  `json:"paySuccessTime"` // 支付成功时间，毫秒级时间戳
}
