package app

// 错误码
const (
	CodeOK                  = iota
	CodeInnerError          = iota + 5500 // 内部错误
	CodeOrderNotExist                     // 订单不存在
	CodeDuplicateOrder                    // 订单号重复
	CodeTooManyRequests                   // 请求过于频繁
	CodeParam                             // 参数有误
	CodeInvalidMerchat                    // 非法商户号
	CodeServer                            // 服务器维护中
	CodeSign                              // 签名错误
	CodeRisk                              // 风险账号
	CodeNoChannel                         // 无可用通道
	CodeMerchatClose                      // 商户关闭
	CodeInsufficientBalance               // 余额不足
	CodeIP                                // 非法ip
)
