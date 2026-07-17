package sypay

const (
	baseURL          = "https://in-sypay.com"
	payURL           = "/sy/pay"
	queryPayURL      = "/sy/qorders"
	withdrawURL      = "/sy/repay"
	queryWithdrawURL = "/sy/qorders"
	balURL           = "/Balance"
	// mid              = 2313
	appID = "S1776148625210"
	key   = "f46df666048776d03cfd93233d20d95c"
)

type PayReq struct {
	AppID       string `json:"app_id"`       // 平台端分配的应用唯一标识
	OrderNo     string `json:"order_no"`     // 商户订单号
	OrderAmount int    `json:"order_amount"` // 金额（整数，单位以接口约定为准）
	NotifyURL   string `json:"notify_url"`   // 异步通知地址
	Sign        string `json:"sign"`         // 数据校验串（签名）
}

type PayResp struct {
	Code    any    `json:"code"`    // 返回码：0=成功，其他为失败码
	Msg     string `json:"msg"`     // 返回说明（如：Ok）
	PayURL  string `json:"pay_url"` // 支付链接
	OrderNo string `json:"orderno"` // 商户订单号
}

type PayCallbackReq struct {
	AppID          string `json:"app_id"`           // 商户号
	ChannelOrderNo string `json:"channel_order_no"` // 平台订单号
	PltOrderNo     string `json:"plt_order_no"`     // 平台订单号
	OrderAmount    int    `json:"order_amount"`     // 订单金额
	RealAmount     int    `json:"real_amount"`      // 实际支付金额
	Status         int    `json:"status"`           // 状态：1=成功，-1=失败
	Sign           string `json:"sign"`             // 签名
}

type QueryPayReq struct {
	AppID   string `json:"app_id"`   // 平台端分配的应用唯一标识
	OrderNo string `json:"order_no"` // 用户自己订单号
	Sign    string `json:"sign"`     // 签名
}

type QueryPayResp struct {
	Code   int    `json:"code"`   // 返回码：0=成功，其他为失败码
	Msg    string `json:"msg"`    // 返回说明（如：Ok）
	States int    `json:"states"` // 订单状态：0=待支付，1=已支付未通知，2=已支付已通知
}

type WithdrawReq struct {
	AppID       string `json:"app_id"`       // 平台端分配的应用唯一标识
	OrderNo     string `json:"order_no"`     // 用户自己订单号
	OrderAmount int    `json:"order_amount"` // 金额（整数，单位以接口约定为准）
	PayType     string `json:"pay_type"`     // india-bank-repay / india-upi-repay
	BankName    string `json:"bank_name"`    // 银行名称（UPI 代付不需要）
	BankCard    string `json:"bank_card"`    // 银行账号 或 UPI 账号
	BankBranch  string `json:"bank_branch"`  // IFSC Code（UPI 代付不需要）
	UserName    string `json:"user_name"`    // 收款人姓名
	NotifyURL   string `json:"notify_url"`   // 异步通知地址
	Sign        string `json:"sign"`         // 数据校验串（签名）
}

type WithdrawResp struct {
	Code    any    `json:"code"`    // 返回码：0=成功，其他为失败码
	Msg     string `json:"msg"`     // 返回说明（如：提交成功）
	OrderNo string `json:"orderno"` // 平台订单号
}

type WithdrawCallbackReq struct {
	AppID          string `json:"app_id"`           // 商户号
	ChannelOrderNo string `json:"channel_order_no"` // 平台订单号
	PltOrderNo     string `json:"plt_order_no"`     // 平台订单号
	OrderAmount    string `json:"order_amount"`     // 金额
	Status         int    `json:"status"`           // 状态：1=成功，-1=失败
	Sign           string `json:"sign"`             // 签名
}

type QueryWithdrawReq struct {
	AppID   string `json:"app_id"`   // 平台端分配的应用唯一标识
	OrderNo string `json:"order_no"` // 用户自己订单号
	Sign    string `json:"sign"`     // 签名
}

type QueryWithdrawResp struct {
	Code   int    `json:"code"`   // 返回码：0=成功，其他为失败码
	Msg    string `json:"msg"`    // 返回说明（如：Ok）
	States int    `json:"states"` // 订单状态：0=待支付，1=已支付未通知，2=已支付已通知
}
