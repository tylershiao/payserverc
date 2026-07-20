package app

// API 响应码。
// 41xxx: 请求校验；42xxx: 商户；43xxx: 订单；44xxx: 风控；
// 45xxx: 支付通道；50xxx: 平台服务。
const (
	CodeOK                  = 0     // 成功
	CodeParam               = 41001 // 请求参数错误
	CodeSign                = 41002 // 签名校验失败
	CodeIP                  = 41003 // 请求 IP 不在白名单
	CodeTooManyRequests     = 41029 // 请求过于频繁
	CodeInvalidMerchat      = 42001 // 商户不存在
	CodeMerchatClose        = 42002 // 商户不可用
	CodeInsufficientBalance = 42003 // 商户余额不足
	CodeOrderNotExist       = 43001 // 订单不存在
	CodeDuplicateOrder      = 43002 // 商户订单号重复
	CodeRisk                = 44001 // 风控拒绝
	CodeNoChannel           = 45001 // 暂无可用支付通道
	CodeInnerError          = 50001 // 平台内部错误
	CodeServer              = 50002 // 平台维护中
)
