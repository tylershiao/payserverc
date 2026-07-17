package mqr

// ==========================
// 基础常量定义
// ==========================
const (
	// API基础URL
	BaseURL = "https://order.octopus-pay.com" // 需要根据实际API地址修改

	// API端点
	PayURL      = "/api/order/create" // 代收下单
	QueryPayURL = "/api/order/status" // 代收查询

	// 商户配置
	MerchantID = "VIRGOPLAY"
	APIKey     = "143d7615597569b5397e10718a6729bb259414d1d9a521a83e834da08ed0c2a76a435d064cd3be71"
	APISecret  = "06e365b0347b2b863d74a091c6f3469ee9e966380f22b458b07ba420fe314cc7273eacaf2e133150"
	SignType   = "HMAC-SHA256" // 签名类型: HMAC-SHA256
	Currency   = "INR"         // 货币类型
)

// ============ 代收（Pay）相关结构体 ============

// PayReq 代收下单请求
type PayReq struct {
	ProjectCode     string `json:"project_code"`      // 项目代码/商户ID
	MerchantOrderNo string `json:"merchant_order_no"` // 商户订单号
	UserUID         string `json:"user_uid"`          // 用户ID
	Amount          string `json:"amount"`            // 金额
	Currency        string `json:"currency"`          // 货币类型
	PaymentMode     int    `json:"payment_mode"`      // 支付模式
	NotifyURL       string `json:"notify_url"`        // 异步通知地址
	ReturnURL       string `json:"return_url"`        // 同步跳转地址
	ExtraData       string `json:"extra_data"`        // 扩展数据
}

// PayResp 代收下单响应
type PayResp struct {
	Code int                `json:"code"` // 返回状态，0-成功 其他-失败
	Msg  string             `json:"msg"`  // 返回信息
	Data *OrderResponseData `json:"data,omitempty"`
}

// PayCallbackReq 代收回调通知
type PayCallbackReq struct {
	OrderNo         string `json:"order_no"`          // 订单号
	MerchantOrderNo string `json:"merchant_order_no"` // 商户订单号
	ProjectCode     string `json:"project_code"`      // 项目代码
	UserUID         string `json:"user_uid"`          // 用户ID
	Amount          string `json:"amount"`            // 金额
	Currency        string `json:"currency"`          // 货币类型
	Status          int    `json:"status"`            // 订单状态
	PaymentMode     int    `json:"payment_mode"`      // 支付模式
	CreatedAt       int64  `json:"created_at"`        // 创建时间
	ExtraData       string `json:"extra_data"`        // 扩展数据
	Signature       string `json:"signature"`         // 签名
}

// QueryPayReq 代收查询请求
type QueryPayReq struct {
	OrderNo string `json:"order_no"` // 订单号
}

// QueryPayResp 代收查询响应
type QueryPayResp struct {
	Code    int              `json:"code"`    // 返回状态，200-成功 其他-失败
	Message string           `json:"message"` // 返回信息
	Data    *OrderStatusData `json:"data,omitempty"`
}
