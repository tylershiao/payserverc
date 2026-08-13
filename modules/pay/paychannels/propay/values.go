package propay

// ==========================
// 基础常量定义
// ==========================
const (
	// API基础URL
	baseURL = "https://gtw.propay.one"

	// API端点
	payURL           = "/api/pay/pay"      // 代收下单
	queryPayURL      = "/api/pay/query"    // 代收查询
	withdrawURL      = "/api/payout/pay"   // 代付下单
	queryWithdrawURL = "/api/payout/query" // 代付查询

	// 商户配置 - 需要根据实际配置修改
	mchNo = "M1756377521"                                                                                                                      // 商户号
	appId = "68b031b1e4b09f5f71f96d79"                                                                                                         // 应用编号
	key   = "lrszu4M3uJhl5TBqdy9weYd6hoe3cwL8fwjN1fwSr3djjoxG3q5bLXf41LYbyLU5lrh0TWke9Hx6nvwgW1adAVQXvlX64q0nqm3MmlLL6Czf701PkUs1h8H60pY7mS2P" // 密钥
)

// ============ 代收（Pay）相关结构体 ============

// PayReq 代收下单请求
type PayReq struct {
	MchNo         string `json:"mchNo"`         // 商户号
	AppId         string `json:"appId"`         // 应用编号
	MchOrderNo    string `json:"mchOrderNo"`    // 商户订单号
	Amount        int64  `json:"amount"`        // 代收金额（单位：分）
	CustomerName  string `json:"customerName"`  // 客户姓名
	CustomerEmail string `json:"customerEmail"` // 客户邮箱
	CustomerPhone string `json:"customerPhone"` // 客户手机号
	NotifyUrl     string `json:"notifyUrl"`     // 异步通知地址
	ReturnUrl     string `json:"returnUrl"`     // 同步跳转的地址
	Sign          string `json:"sign"`          // 签名
}

// PayResp 代收下单响应
type PayResp struct {
	Code int    `json:"code"` // 返回状态，0-成功 其他-失败
	Msg  string `json:"msg"`  // 返回信息
	Sign string `json:"sign"` // 签名信息
	Data struct {
		PayOrderId  string `json:"payOrderId"`  // 代收订单号
		MchOrderNo  string `json:"mchOrderNo"`  // 商户订单号
		OrderState  int    `json:"orderState"`  // 订单状态
		PayDataType string `json:"payDataType"` // 支付参数类型
		PayData     string `json:"payData"`     // 收银台地址
		UpiUrl      string `json:"upiUrl"`      // UPI协议
		ErrCode     string `json:"errCode"`     // 上游错误码
		ErrMsg      string `json:"errMsg"`      // 上游错误描述
	} `json:"data"`
}

// PayCallbackReq 代收回调通知
type PayCallbackReq struct {
	PayOrderId     string `json:"payOrderId"`     // 代收订单号
	MchNo          string `json:"mchNo"`          // 商户号
	AppId          string `json:"appId"`          // 应用编号
	WayCode        string `json:"wayCode"`        // 支付方式
	IfCode         string `json:"ifCode"`         // 支付编码
	ClientIp       string `json:"clientIp"`       // 客户端IP
	MchOrderNo     string `json:"mchOrderNo"`     // 商户订单号
	ChannelOrderNo string `json:"channelOrderNo"` // 渠道订单号
	Amount         int64  `json:"amount"`         // 支付金额（单位：分）
	Currency       string `json:"currency"`       // 货币代码，固定值 INR
	CustomerName   string `json:"customerName"`   // 客户姓名
	CustomerEmail  string `json:"customerEmail"`  // 客户邮箱
	CustomerPhone  string `json:"customerPhone"`  // 客户手机号
	State          int    `json:"state"`          // 订单状态 2-成功 3-失败
	ErrCode        string `json:"errCode"`        // 渠道错误码
	ErrMsg         string `json:"errMsg"`         // 渠道错误描述
	CreatedAt      int64  `json:"createdAt"`      // 创建时间，13位时间戳
	ReqTime        int64  `json:"reqTime"`        // 通知请求时间，13位时间戳
	ChannelExtra   string `json:"channelExtra"`
	Sign           string `json:"sign"` // 签名
}

// QueryPayReq 代收查询请求
type QueryPayReq struct {
	MchNo      string `json:"mchNo"`                // 商户号
	AppId      string `json:"appId"`                // 应用编号
	PayOrderId string `json:"payOrderId,omitempty"` // 代收订单号（与mchOrderNo二选一）
	MchOrderNo string `json:"mchOrderNo,omitempty"` // 商户订单号（与payOrderId二选一）
	Sign       string `json:"sign"`                 // 签名
}

// QueryPayResp 代收查询响应
type QueryPayResp struct {
	Code int    `json:"code"` // 返回状态，0-成功 其他-失败
	Msg  string `json:"msg"`  // 返回信息
	Sign string `json:"sign"` // 签名信息
	Data struct {
		PayOrderId     string `json:"payOrderId"`     // 代收订单号
		MchNo          string `json:"mchNo"`          // 商户号
		AppId          string `json:"appId"`          // 应用编号
		MchOrderNo     string `json:"mchOrderNo"`     // 商户订单号
		Amount         int64  `json:"amount"`         // 支付金额（单位：分）
		Currency       string `json:"currency"`       // 货币代码，固定值 INR
		State          int    `json:"state"`          // 订单状态 1-支付中 2-支付成功 3-支付失败 6-订单关闭
		CustomerName   string `json:"customerName"`   // 客户姓名
		CustomerEmail  string `json:"customerEmail"`  // 客户邮箱
		CustomerPhone  string `json:"customerPhone"`  // 客户手机号
		ChannelOrderNo string `json:"channelOrderNo"` // 渠道订单号
		ErrCode        string `json:"errCode"`        // 渠道错误码
		ErrMsg         string `json:"errMsg"`         // 渠道错误描述
		ExtParam       string `json:"extParam"`       // 商户扩展参数
		CreatedAt      int64  `json:"createdAt"`      // 创建时间，13位时间戳
		SuccessTime    int64  `json:"successTime"`    // 成功时间，13位时间戳
	} `json:"data"`
}

// ============ 代付（Withdraw/Payout）相关结构体 ============

// WithdrawReq 代付下单请求
type WithdrawReq struct {
	MchNo        string `json:"mchNo"`                  // 商户号
	AppId        string `json:"appId"`                  // 应用编号
	MchOrderNo   string `json:"mchOrderNo"`             // 商户订单号
	Amount       int64  `json:"amount"`                 // 代付金额（单位：分）
	EntryType    string `json:"entryType"`              // 入账方式 IMPS或UPI（目前暂支持IMPS）
	AccountNo    string `json:"accountNo"`              // 收款账号或银行卡号
	AccountCode  string `json:"accountCode"`            // 银行编码（IMPS时为IFSC代码）
	BankName     string `json:"bankName,omitempty"`     // 收款账号开户行名称
	AccountName  string `json:"accountName"`            // 收款人姓名
	AccountEmail string `json:"accountEmail"`           // 收款人邮箱
	AccountPhone string `json:"accountPhone"`           // 收款人手机号
	NotifyUrl    string `json:"notifyUrl"`              // 异步通知地址
	TransferDesc string `json:"transferDesc,omitempty"` // 备注信息
	Sign         string `json:"sign"`                   // 签名
}

// WithdrawResp 代付下单响应
type WithdrawResp struct {
	Code int    `json:"code"` // 返回状态，0-成功 其他-失败
	Msg  string `json:"msg"`  // 返回信息
	Sign string `json:"sign"` // 签名信息
	Data struct {
		TransferID   string `json:"transferId"`   // 代付订单号
		MchOrderNo   string `json:"mchOrderNo"`   // 商户订单号
		Amount       int64  `json:"amount"`       // 到账金额（扣除手续费后，单位：分）
		MchFeeAmount int64  `json:"mchFeeAmount"` // 手续费（单位：分）
		AmountTo     int64  `json:"amountTo"`     // 交易总金额（含手续费，单位：分）
		AccountNo    string `json:"accountNo"`    // 收款账号
		AccountName  string `json:"accountName"`  // 收款人姓名
		State        int    `json:"state"`        // 转账状态 0-订单生成 1-转账中 2-转账成功 3-转账失败 4-转账关闭
		ErrCode      string `json:"errCode"`      // 渠道错误码
		ErrMsg       string `json:"errMsg"`       // 渠道错误描述
	} `json:"data"`
}

// WithdrawCallbackReq 代付回调通知
type WithdrawCallbackReq struct {
	MchNo          string `json:"mchNo"`          // 商户号
	AppId          string `json:"appId"`          // 应用编号
	MchOrderNo     string `json:"mchOrderNo"`     // 商户转账订单号
	TransferID     string `json:"transferId"`     // 代付订单号
	Amount         int64  `json:"amount"`         // 到账金额（不含手续费，单位：分）
	MchFeeAmount   int64  `json:"mchFeeAmount"`   // 手续费（单位：分）
	AmountTo       int64  `json:"amountTo"`       // 交易总金额（含手续费，单位：分）
	Currency       string `json:"currency"`       // 货币代码，固定值 INR
	EntryType      string `json:"entryType"`      // 入账方式，固定值 IMPS
	State          int    `json:"state"`          // 转账状态 2-成功 3-失败
	AccountNo      string `json:"accountNo"`      // 收款账号
	AccountName    string `json:"accountName"`    // 收款人姓名
	AccountEmail   string `json:"accountEmail"`   // 收款人邮箱
	AccountPhone   string `json:"accountPhone"`   // 收款人手机号
	BankName       string `json:"bankName"`       // 收款人开户行名称
	ChannelOrderNo string `json:"channelOrderNo"` // 渠道订单号
	ChannelExtra   string `json:"channelExtra"`   // UTR
	IfCode         string `json:"ifCode"`         // 支付接口
	ErrCode        string `json:"errCode"`        // 渠道错误码
	WayCode        string `json:"wayCode"`        // 支付编码
	ErrMsg         string `json:"errMsg"`         // 渠道错误描述
	ExtraParam     string `json:"extraParam"`     // 商户扩展参数
	CreatedAt      int64  `json:"createdAt"`      // 创建时间，13位时间戳
	SuccessTime    int64  `json:"successTime"`    // 成功时间，13位时间戳
	Sign           string `json:"sign"`           // 签名
	ReqTime        int64  `json:"reqTime"`        // 通知请求时间，13位时间戳
}

// QueryWithdrawReq 代付查询请求
type QueryWithdrawReq struct {
	MchNo      string `json:"mchNo"`                // 商户号
	AppId      string `json:"appId"`                // 应用编号
	TransferID string `json:"transferId,omitempty"` // 代付订单号（与mchOrderNo二选一）
	MchOrderNo string `json:"mchOrderNo,omitempty"` // 商户订单号（与transferId二选一）
	Sign       string `json:"sign"`                 // 签名
}

// QueryWithdrawResp 代付查询响应
type QueryWithdrawResp struct {
	Code int    `json:"code"` // 返回状态，0-成功 其他-失败
	Msg  string `json:"msg"`  // 返回信息
	Sign string `json:"sign"` // 签名信息
	Data struct {
		MchNo          string `json:"mchNo"`          // 商户号
		AppId          string `json:"appId"`          // 应用编号
		MchOrderNo     string `json:"mchOrderNo"`     // 商户订单号
		TransferID     string `json:"transferId"`     // 代付订单号
		Amount         int64  `json:"amount"`         // 到账金额（扣除手续费后，单位：分）
		MchFeeAmount   int64  `json:"mchFeeAmount"`   // 手续费（单位：分）
		AmountTo       int64  `json:"amountTo"`       // 交易总金额（含手续费，单位：分）
		Currency       string `json:"currency"`       // 货币代码，固定值 INR
		EntryType      string `json:"entryType"`      // 入账方式，固定值 IMPS
		State          int    `json:"state"`          // 转账状态 0-订单生成 1-转账中 2-转账成功 3-转账失败 4-转账关闭
		AccountNo      string `json:"accountNo"`      // 收款账号
		AccountName    string `json:"accountName"`    // 收款人姓名
		AccountEmail   string `json:"accountEmail"`   // 收款人邮箱
		AccountPhone   string `json:"accountPhone"`   // 收款人手机号
		BankName       string `json:"bankName"`       // 收款人开户行名称
		TransferDesc   string `json:"transferDesc"`   // 备注信息
		ChannelOrderNo string `json:"channelOrderNo"` // 渠道订单号
		ErrCode        string `json:"errCode"`        // 渠道错误码
		ErrMsg         string `json:"errMsg"`         // 渠道错误描述
		ExtraParam     string `json:"extraParam"`     // 商户扩展参数
		CreatedAt      int64  `json:"createdAt"`      // 创建时间，13位时间戳
		SuccessTime    int64  `json:"successTime"`    // 成功时间，13位时间戳
	} `json:"data"`
}
