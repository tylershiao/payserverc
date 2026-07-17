package velopay

const (
	baseURL          = "https://velopay.ptasm.online"
	payURL           = "/api/payin/request"
	queryPayURL      = "/api/payin/status"
	withdrawURL      = "/api/payout/request"
	queryWithdrawURL = "/api/payout/status"
	balURL           = "/api/merchant/balance"
	mid              = "550"
	key              = "b881dc77ccfc4857b32b3f4a0318f0da"
)

type PayReq struct {
	MerchantNo string `json:"merchant_no"` // 商户编号
	TxnID      string `json:"txn_id"`      // 商户订单号
	Amount     int64  `json:"amount"`      // 金额（如：10000）
	Type       int    `json:"type"`        // 1=返回HTML5，2=返回UPI账号（钱包支付）
	Callback   string `json:"callback"`    // 异步回调地址
	Currency   string `json:"currency"`    // 货币（INR）
	Sign       string `json:"sign"`        // 验签值
}

type PayResp struct {
	Status     string `json:"status"`      // 状态码（SUCCESS）
	Message    string `json:"message"`     // 说明（TXN CREATED）
	MerchantNo string `json:"merchant_no"` // 商户编码
	TxnID      string `json:"txn_id"`      // 商户订单号
	ID         int    `json:"id"`          // 平台流水号
	Amount     int64  `json:"amount"`      // 金额
	PayLink    string `json:"pay_link"`    // 支付链接 / UPI账号（取决于type）
}

type PayCallbackReq struct {
	Status  string `json:"status"`  // 状态：SUCCESS=订单支付成功
	Message string `json:"message"` // 说明（TXN SUCCESS）
	ID      string `json:"id"`      // 平台流水号
	TxnID   string `json:"txn_id"`  // 商户订单号
	// MerchantNo string `json:"merchant_no"` // 商户编号
	Sign string `json:"sign"` // 签名
	UTR  string `json:"utr"`  // UTR（银行唯一交易号）
	// Amount     string `json:"amount"`      // 金额
}

type QueryPayReq struct {
	MerchantNo string `json:"merchant_no"` // 商户编码
	TxnID      string `json:"txn_id"`      // 商户订单号
	Sign       string `json:"sign"`        // 验签值
}

type QueryPayResp struct {
	Status string `json:"status"` // 订单状态：PENDING / SUCCESS
	TxnID  string `json:"txn_id"` // 商户订单号
	Amount int64  `json:"amount"` // 金额（单位：分）
}

type WithdrawReq struct {
	MerchantNo string `json:"merchant_no"`      // 商户编号
	TxnID      string `json:"txn_id"`           // 商户订单号
	Amount     int64  `json:"amount"`           // 金额（整数，100 = 10000）
	Type       string `json:"type"`             // 类型（固定 1）
	IFSC       string `json:"ifsc"`             // IFSC（11位）
	CardNum    string `json:"card_num"`         // 收款账号
	Name       string `json:"name,omitempty"`   // 姓名（可选）
	Email      string `json:"email,omitempty"`  // 邮箱（可选）
	Mobile     string `json:"mobile,omitempty"` // 手机号（可选）
	Callback   string `json:"callback"`         // 异步回调地址
	Currency   string `json:"currency"`         // 货币（INR）
	Sign       string `json:"sign"`             // 加密验签值
}

type WithdrawResp struct {
	Status     string `json:"status"`      // 状态：PENDING / SUCCESS / FAILED（以实际文档为准）
	Message    string `json:"message"`     // 说明（TXN CREATED）
	MerchantNo string `json:"merchant_no"` // 商户编码
	TxnID      string `json:"txn_id"`      // 商户订单号
	ID         int    `json:"id"`          // 平台流水号
	Amount     int64  `json:"amount"`      // 金额
}

type WithdrawCallbackReq struct {
	Status  string `json:"status"`  // 状态：SUCCESS=提现成功，FAILED=提现失败
	Message string `json:"message"` // 说明（TXN SUCCESS）
	// MerchantNo string `json:"merchant_no"` // 商户编号
	ID    string `json:"id"`     // 平台流水号
	TxnID string `json:"txn_id"` // 商户订单号
	// Amount     string `json:"amount"`      // 金额
	UTR  string `json:"utr"`  // UTR
	Sign string `json:"sign"` // 签名
}

type QueryWithdrawReq struct {
	MerchantNo string `json:"merchant_no"` // 商户编码
	TxnID      string `json:"txn_id"`      // 商户订单号
	Sign       string `json:"sign"`        // 验签值
}

type QueryWithdrawResp struct {
	Status     string `json:"status"`      // 状态：PENDING / SUCCESS / FAILED
	TxnID      string `json:"txn_id"`      // 商户订单号
	FailReason string `json:"fail_reason"` // 失败原因（SUCCESS / PENDING 时可能为空）
	Amount     int64  `json:"amount"`      // 金额（单位：分）
}
