package hexpay

const (
	baseURL          = "https://api.hexmeta.xyz"
	payURL           = "/exposed/v1/order/payment"
	queryPayURL      = "/exposed/v1/query/payment"
	withdrawURL      = "/exposed/v1/order/payout"
	queryWithdrawURL = "/exposed/v1/query/payout"
	balURL           = "/exposed/v1/query/balance"
	mid              = "189643432899"
	key              = "ekFUyHH8nd3Wy8Ye8M8dGMxJQJMXgWtl"
)

type PayReq struct {
	AppID     string `json:"appId"`     // 应用ID，创建商户后提供（必填）
	Timestamp string `json:"timestamp"` // Unix时间戳（毫秒，13位），与服务器时间误差不超过5分钟（必填）
	Sign      string `json:"sign"`      // 请求签名（必填）

	OrderNo string `json:"orderNo"` // 商户订单号，在该商户下需唯一（必填）
	Amount  string `json:"amount"`  // 订单金额，单位：分（1元=100分），不支持小数（必填）

	Name  string `json:"name"`  // 付款人姓名，为提高成功率请尽量真实（必填）
	Phone string `json:"phone"` // 付款人手机号，为提高成功率请尽量真实（必填）
	Email string `json:"email"` // 付款人邮箱，为提高成功率请尽量真实（必填）

	UserIP    string `json:"userIP"`    // 用户IP地址（非必填）
	NotifyURL string `json:"notifyUrl"` // 异步回调地址，订单状态变化后通知（非必填）
}

type PayResp struct {
	Code int    `json:"code"` // 请求错误码：1=成功，其他=失败
	Msg  string `json:"msg"`  // 错误说明：OK=成功，其他=失败
	Data struct {
		SysOrderNo string `json:"sysOrderNo"` // 平台订单号
		PayURL     string `json:"payUrl"`     // 收银台地址
	} `json:"data"` // 请求结果
}

type PayCallbackReq struct {
	Timestamp string `json:"timestamp"` // Unix时间戳（毫秒，13位）
	Sign      string `json:"sign"`      // 请求签名

	OrderNo    string `json:"orderNo"`    // 商户订单号
	SysOrderNo string `json:"sysOrderNo"` // 平台订单号

	Status string `json:"status"` // 订单状态：1=处理中，2=成功，3=失败
	Amount string `json:"amount"` // 实收金额，单位：分（1元=100分）
}

type QueryPayReq struct {
	AppID     string `json:"appId"`     // 应用ID，创建商户后提供（必选）
	Timestamp string `json:"timestamp"` // Unix时间戳（毫秒，13位），与服务器时间误差不超过5分钟（必选）
	Sign      string `json:"sign"`      // 请求签名（必选）

	OrderNo string `json:"orderNo"` // 商户订单号（必选）
}

type QueryPayResp struct {
	Code int    `json:"code"` // 请求错误码：1=成功，其他=失败
	Msg  string `json:"msg"`  // 错误说明：OK=成功，其他=失败
	Data struct {
		AppID      string `json:"appId"`      // 应用ID
		OrderNo    string `json:"orderNo"`    // 商户订单号
		SysOrderNo string `json:"sysOrderNo"` // 平台订单号
		Amount     string `json:"amount"`     // 订单金额，单位：分（1元=100分）
		Status     string `json:"status"`     // 订单状态：1=处理中，2=成功，3=失败
	} `json:"data"` // 请求结果
}

type WithdrawReq struct {
	AppID     string `json:"appId"`     // 应用ID，创建商户后提供（必选）
	Timestamp string `json:"timestamp"` // Unix时间戳（毫秒13位），与服务器时间误差不超过5分钟（必选）
	Sign      string `json:"sign"`      // 请求签名（必选）

	OrderNo string `json:"orderNo"` // 商户订单号，在该商户下需唯一（必选）
	WayCode string `json:"wayCode"` // 代付方式：1=银行卡，2=UPI（当前只支持1）（必选）
	Amount  string `json:"amount"`  // 订单金额，单位：分（1元=100分），不支持小数（必选）

	BankName  string `json:"bankName"`  // 开户银行名称，没有则固定传 IndiaBank（必选）
	BankCode  string `json:"bankCode"`  // IFSC代码（4位大写字母 + 0 + 6位数字）（必选）
	AccountNo string `json:"accountNo"` // 收款账号：银行卡号 / UPI账号（必选）

	Name  string `json:"name"`  // 收款人姓名（必选）
	Phone string `json:"phone"` // 收款人手机号（必选）
	Email string `json:"email"` // 收款人邮箱（必选）

	NotifyURL string `json:"notifyUrl"` // 异步回调地址（非必选，但不省略）
}

type WithdrawResp struct {
	Code int    `json:"code"` // 请求错误码：1=成功，其他=失败
	Msg  string `json:"msg"`  // 错误说明：OK=成功，其他=失败
	Data struct {
		SysOrderNo string `json:"sysOrderNo"` // 平台订单号
	} `json:"data"` // 请求结果
}

type WithdrawCallbackReq struct {
	OrderNo     string `json:"orderNo"`     // 商户订单号
	SysOrderNo  string `json:"sysOrderNo"`  // 平台订单号
	Status      string `json:"status"`      // 订单状态：1=处理中，2=成功，3=失败
	Amount      string `json:"amount"`      // 交易金额，单位：分（1元=100分）
	ReferenceNo string `json:"referenceNo"` // UTR
	Timestamp   string `json:"timestamp"`   // Unix时间戳（毫秒）
	Sign        string `json:"sign"`        // 请求签名
}

type QueryWithdrawReq struct {
	AppID     string `json:"appId"`     // 应用ID，创建商户后提供（必选）
	Timestamp string `json:"timestamp"` // Unix时间戳（毫秒13位），与服务器时间误差不可超过5分钟（必选）
	Sign      string `json:"sign"`      // 请求签名（必选）

	OrderNo string `json:"orderNo"` // 商户订单号（必选）
}

type QueryWithdrawResp struct {
	Code int    `json:"code"` // 请求错误码：1=成功，其他=失败
	Msg  string `json:"msg"`  // 错误说明：OK=成功，其他=失败
	Data struct {
		AppID       string `json:"appId"`       // 应用ID
		OrderNo     string `json:"orderNo"`     // 商户订单号
		SysOrderNo  string `json:"sysOrderNo"`  // 平台订单号
		WayCode     string `json:"wayCode"`     // 代付方式：1=银行卡，2=UPI
		Status      string `json:"status"`      // 订单状态：1=处理中，2=成功，3=失败
		Amount      string `json:"amount"`      // 订单金额，单位：分（1元=100分）
		ReferenceNo string `json:"referenceNo"` // UTR / referenceNo
	} `json:"data"` // 请求结果
}
