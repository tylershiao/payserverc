package worth_pay

const (
	baseURL     = "https://openapi.linkpayment.net/inr/v2"
	payURL      = "/collection"
	withdrawURL = "/payment"
	mid         = "virgo"
	key         = "QOF2omwhcuavid5m"
)

// 单笔代收下单 https://docs.linkpayment.net/#india-collection-order
type PayReq struct {
	CustomerName    string `json:"customer_name"`     // 客户姓名
	CustomerPhone   string `json:"customer_phone"`    // 客户手机号
	CustomerEmail   string `json:"customer_email"`    // 客户电子邮箱
	MerchantOrderId string `json:"merchant_order_id"` // 订单号唯一
	Amount          int64  `json:"amount"`            // 订单金额，单位（分）
	Purpose         string `json:"purpose"`           // 对于产品运营商，可在商户后台设置运营产品名称，系统自动生成支付是由。若未设置，该项必填；对于代理平台，该项必填。该信息会在用户界面显示。例如：Payment for PlayStation Game
	ReturnUrl       string `json:"return_url"`        // 支付成功后跳转的页面地址
	NotifyUrl       string `json:"notify_url"`        // 支付状态变更通知回调地址，默认地址可在商户后台设置
	OrderNote       string `json:"order_note"`        // 备注信息
}

type PayResp struct {
	Code            int    `json:"code"`              // 0 表示交易成功。 其他表示失败。详见错误代码定义
	Message         string `json:"message"`           // 返回代码描述
	MerchantOrderId string `json:"merchant_order_id"` // 商户订单号
	TxId            string `json:"tx_id"`             // 平台订单号
	PaymentLink     string `json:"payment_link"`      // 支付页面地址
	Reference       string `json:"reference"`         // 商户可以根据该值自己生成二维支付码
}

// 代收状态变更通知 https://docs.linkpayment.net/#india-collection-notify
type PayCallbackReq struct {
	MerchantOrderId string `json:"merchant_order_id"` // 商户订单号
	TxId            string `json:"tx_id"`             // 平台订单号
	OrderCurrency   string `json:"order_currency"`    // 交易货币
	Amount          int64  `json:"amount"`            // 单位 分
	OrderStatus     string `json:"order_status"`      // 0 支付成功，1 支付失败，2 待支付
	OrderNote       string `json:"order_note"`        // 备注信息
	PaymentLink     string `json:"payment_link"`      // 支付页面地址
	ErrorMessage    string `json:"error_message"`     // 错误信息
}

// 单笔代付下单 https://docs.linkpayment.net/#india-payment-order
type WithdrawReq struct {
	MerchantOrderId string `json:"merchant_order_id"` // 商户订单号
	Amount          int64  `json:"amount"`            // 订单金额，单位（分）
	TransferMode    string `json:"transfer_mode"`     // 转账模式。填写BANKTRANSFER或UPI
	AccountNumber   string `json:"account_number"`    // 账户号码。卡号/UPI号
	NotifyUrl       string `json:"notify_url"`        // 支付状态变更通知地址，默认地址可在商户后台设置
	BeneName        string `json:"bene_name"`         // 收款人姓名
	BenePhone       string `json:"bene_phone"`        // 手机号
	BeneEmail       string `json:"bene_email"`        // 邮箱
	BeneIfsc        string `json:"bene_ifsc"`         // IFSC码
	Remarks         string `json:"remarks"`           // 备注信息
}

type WithdrawResp struct {
	Code            int    `json:"code"`              // 0 表示交易成功。 其他表示失败。详见错误代码定义
	Message         string `json:"message"`           // 返回代码描述
	MerchantOrderId string `json:"merchant_order_id"` // 商户订单号
	TxId            string `json:"tx_id"`             // 平台订单号
}

// 代付状态变更通知 https://docs.linkpayment.net/#india-payment-notify
type WithdrawCallbackReq struct {
	MerchantOrderId string `json:"merchant_order_id"` // 商户订单号
	TxId            string `json:"tx_id"`             // 平台订单号
	OrderCurrency   string `json:"order_currency"`    // 交易货币
	Amount          int64  `json:"amount"`            // 订单金额，单位（分）
	OrderStatus     string `json:"order_status"`      // 0 支付成功，1 支付失败，2 待支付
	TransferMode    string `json:"transfer_mode"`     // 转账模式。填写BANKTRANSFER或UPI
	Remarks         string `json:"remarks"`           // 备注信息
	ErrorMessage    string `json:"error_message"`     // 错误信息
}

// 单笔代付订单查询 https://docs.linkpayment.net/#india-payment-query
type QueryWithdrawReq struct {
	MerchantOrderId string `json:"merchant_order_id"` // 商户订单号
	TxId            string `json:"tx_id"`             // 平台订单号
	Nonce           string `json:"nonce"`             // 随机字符串
}

type QueryWithdrawResp struct {
	Code            int    `json:"code"`              // 0 表示交易成功。 其他表示失败。详见错误代码定义
	Message         string `json:"message"`           // 返回代码描述
	MerchantOrderId string `json:"merchant_order_id"` // 商户订单号
	TxId            string `json:"tx_id"`             // 平台订单号
	OrderCurrency   string `json:"order_currency"`    // 交易货币
	Amount          int64  `json:"amount"`            // 单位 分
	OrderStatus     string `json:"order_status"`      // 0 支付成功，1 支付失败，2 资金回退，3 处理中
	BeneName        string `json:"bene_name"`         // 收款人姓名
	BenePhone       string `json:"bene_phone"`        // 客户手机号
	BeneEmail       string `json:"bene_email"`        // 收款人电子邮箱
	BeneBankAcct    string `json:"bene_bank_acct"`    // 	收款人账号
	BeneIfsc        string `json:"bene_ifsc"`         // 	收款人银行代码
	BeneVpa         string `json:"bene_vpa"`          // UPI模式必填
	BeneCardNo      string `json:"bene_card_no"`      // 收款人卡号，只有数字。从4或5开始(16个字符限制)
	BeneAddress     string `json:"bene_address"`      // 收款人地址
	Utr             string `json:"utr"`               // 	utr
	TransferMode    string `json:"transfer_mode"`     // 转账模式。允许的值是BANKTRANSFER。
	NotifyUrl       string `json:"notify_url"`        // 支付状态变更通知地址，默认地址可在商户后台设置
	Remarks         string `json:"remarks"`           // 备注信息
	ErrorMessage    string `json:"error_message"`     // 错误信息
}

// 单笔代收订单查询 https://docs.linkpayment.net/#india-collection-query
type QueryPayReq struct {
	MerchantOrderId string `json:"merchant_order_id"` // 商户订单号
	TxId            string `json:"tx_id"`             // 平台订单号
	Nonce           string `json:"nonce"`             // 随机字符串
}

type QueryPayResp struct {
	Code            int    `json:"code"`              // 0 表示交易成功。 其他表示失败。详见错误代码定义
	Message         string `json:"message"`           // 返回代码描述
	MerchantOrderId string `json:"merchant_order_id"` // 商户订单号
	TxId            string `json:"tx_id"`             // 平台订单号
	OrderCurrency   string `json:"order_currency"`    // 交易货币
	Amount          int64  `json:"amount"`            // 单位 分
	CustomerEmail   string `json:"customer_email"`    // 客户姓名
	CustomerPhone   string `json:"customer_phone"`    // 客户手机号
	CustomerName    string `json:"customer_name"`     // 客户电子邮箱
	Purpose         string `json:"purpose"`           // 支付事由
	OrderStatus     string `json:"order_status"`      // 0 支付成功，1 支付失败，2 待支付
	OrderNote       string `json:"order_note"`        // 订单备注信息
	PaymentLink     string `json:"payment_link"`      // 支付页面地址
	ReturnUrl       string `json:"return_url"`        // 支付成功后跳转的页面地址
	NotifyUrl       string `json:"notify_url"`        // 支付状态变更通知地址，默认地址可在商户后台
}
