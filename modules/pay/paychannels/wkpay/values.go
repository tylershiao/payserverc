package wkpay

const (
	baseURL          = "https://api.okty875.vip"
	payURL           = "/Pay_Index.html"
	queryPayURL      = "/Pay_Trade_query.html"
	withdrawURL      = "/Payment_Dfpay_add.html"
	queryWithdrawURL = "/Payment_Dfpay_query.html"
	queryBalanceURL  = "/Pay_Trade_balance.html"
	mid              = "251142958"
	key              = "fpz3zpyqraem9y8elwpprpkcu7nuhsz2"
	// mid = "202366100"
	// key = "c385fe7029344aef826d8112625b2625x"
)

type PayReq struct {
	PayMemberID    string `json:"pay_memberid"`    // 商户号（平台分配）
	PayOrderID     string `json:"pay_orderid"`     // 订单号，上送订单号唯一，长度20
	PayApplyDate   string `json:"pay_applydate"`   // 提交时间，格式：2016-12-26 18:18:18
	PayBankCode    string `json:"pay_bankcode"`    // 银行编码，参考后续说明
	PayNotifyURL   string `json:"pay_notifyurl"`   // 服务端通知地址（POST 返回数据）
	PayCallbackURL string `json:"pay_callbackurl"` // 页面跳转通知地址（POST 返回数据）
	PayAmount      string `json:"pay_amount"`      // 订单金额（商品金额）
	PayMd5Sign     string `json:"pay_md5sign"`     // MD5 签名（不参与签名）
	PayAttach      string `json:"pay_attach"`      // 附加字段，返回时原样返回（中文需URL编码）
	Email          string `json:"email"`           // 付款客户邮箱（必填真实邮箱，影响成功率）
	IP             string `json:"IP"`              // 付款客户 IP，如 127.0.0.1
	Phone          string `json:"phone"`           // 付款客户手机号码（必填真实号码，影响成功率）
	PayProductName string `json:"pay_productname"` // 商品名称
	PayProductNum  string `json:"pay_productnum"`  // 商户品数量
	PayProductDesc string `json:"pay_productdesc"` // 商品描述
	PayProductURL  string `json:"pay_producturl"`  // 商户链接地址
	FirstName      string `json:"firstName"`       // 名字（必填真实付款名字）
	LastName       string `json:"lastName"`        // 姓（必填真实付款姓）
	ReturnJSON     string `json:"return_json"`     // 是否返回跳转链接地址（1返回，0不返回）
}

type PayResp struct {
	Code int    `json:"code"` // 业务响应码，0表示成功，其它为失败
	Msg  string `json:"msg"`  // 提示信息
	Data struct {
		PayURL string `json:"pay_url"` // 收银台地址
	} `json:"data"` // 返回的数据对象
}

type PayCallbackReq struct {
	MemberID      string `json:"memberid"`       // 商户编号（平台分配）
	OrderID       string `json:"orderid"`        // 订单号
	Amount        string `json:"amount"`         // 订单金额
	TransactionID string `json:"transaction_id"` // 交易流水号
	Datetime      string `json:"datetime"`       // 交易时间
	ReturnCode    string `json:"returncode"`     // 交易状态：“00”为成功
	Attach        string `json:"attach"`         // 扩展返回（商户附加数据返回）
	Sign          string `json:"sign"`           // 签名（不参与签名验证）
}

type QueryPayReq struct {
	PayMemberID string `json:"pay_memberid"` // 商户号（平台分配）
	PayOrderID  string `json:"pay_orderid"`  // 订单号，上送订单号唯一，长度20
	PayMd5Sign  string `json:"pay_md5sign"`  // 签名（不参与签名）
}

type QueryPayResp struct {
	MemberID      string `json:"memberid"`       // 商户编号（平台分配）
	OrderID       string `json:"orderid"`        // 订单号
	Amount        string `json:"amount"`         // 订单金额
	TimeEnd       string `json:"time_end"`       // 支付成功时间
	TransactionID string `json:"transaction_id"` // 交易流水号
	ReturnCode    string `json:"returncode"`     // 交易状态：“00”为成功
	TradeState    string `json:"trade_state"`    // 交易状态：NOTPAY 未支付，SUCCESS 已支付
	Sign          string `json:"sign"`           // 签名（不参与签名计算
}

type WithdrawReq struct {
	MchID       string `json:"mchid"`        // 商户号（平台分配商户号）
	OutTradeNo  string `json:"out_trade_no"` // 订单号（保证唯一）
	Money       string `json:"money"`        // 金额，单位：元
	BankName    string `json:"bankname"`     // 开户行名称
	SubBranch   string `json:"subbranch"`    // 支行名称
	AccountName string `json:"accountname"`  // 开户名
	CardNumber  string `json:"cardnumber"`   // 银行卡号
	Province    string `json:"province"`     // 省份
	City        string `json:"city"`         // 城市
	Extends     string `json:"extends"`      // 附加字段（数组转JSON后再Base64编码）
	NotifyURL   string `json:"notifyurl"`    // 回调地址
	PayMd5Sign  string `json:"pay_md5sign"`  // MD5签名（不参与签名）
}

type ExtendsData struct {
	Phone string `json:"phone"`
	IFSC  string `json:"ifsc"`
}

type WithdrawResp struct {
	IP  string `json:"ip"`  // 请求来源IP
	Msg string `json:"msg"` // 返回消息，例如："代付申请成功"
	// Data          string `json:"data"`           // 可能为内部处理编号
	Status        string `json:"status"`         // 请求状态，例如："success"
	OrderID       string `json:"OrderId"`        // 商户订单号
	TransactionID string `json:"transaction_id"` // 平台交易号
}

/*
refCode
1 成功
2 失败
3 处理中
4 待处理
5 审核驳回
6 待审核
7 交易不存在
8 未知状态
*/
type WithdrawCallbackReq struct {
	Msg           string `json:"msg"`            // 状态描述（如："请求成功"）
	TransactionID string `json:"transaction_id"` // 平台交易订单号
	Amount        string `json:"amount"`         // 订单金额
	UTR           string `json:"utr"`            // UTR号（可能为空）
	PayMd5Sign    string `json:"pay_md5sign"`    // 签名
	MchID         string `json:"mchid"`          // 商户号
	OutTradeNo    string `json:"out_trade_no"`   // 商户订单号
	RefMsg        string `json:"refMsg"`         // 审核备注（例如：审核驳回）
	RefCode       string `json:"refCode"`        // 审核状态码，例如："5" 失败,"1" 成功
	Status        string `json:"status"`         // 状态：success（成功）、error（失败）
}

type QueryWithdrawReq struct {
	OutTradeNo string `json:"out_trade_no"` // 商户订单号（必填）
	MchID      string `json:"mchid"`        // 商户号（必填）
	PayMd5Sign string `json:"pay_md5sign"`  // 签名（不参与签名）
}

type QueryWithdrawResp struct {
	Status        string `json:"status"`         // 状态：success=请求成功（不代表业务成功），error=请求失败
	Msg           string `json:"msg"`            // 状态描述
	MchID         string `json:"mchid"`          // 商户号（status=success 时返回）
	OutTradeNo    string `json:"out_trade_no"`   // 商户订单号（status=success 时返回）
	Amount        string `json:"amount"`         // 金额（status=success 时返回）
	TransactionID string `json:"transaction_id"` // 平台流水号（status=success 时返回）
	RefCode       string `json:"refCode"`        // 业务状态（status=success 时返回）
	RefMsg        string `json:"refMsg"`         // 业务描述（status=success 时返回）
	SuccessTime   string `json:"success_time"`   // 成功时间（status=success 且 refCode=1 时返回）
	PayMd5Sign    string `json:"pay_md5sign"`    // 签名（不参与签名）
}
