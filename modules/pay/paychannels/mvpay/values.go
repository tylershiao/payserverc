package mvpay

const (
	baseURL          = "https://www.mv-pay.org"
	payURL           = "/Transfer/index"
	queryPayURL      = "/Transfer/getorderstatus"
	withdrawURL      = "/Transfer/replay"
	queryWithdrawURL = "/Transfer/repayorderstatus"
	mid              = "M265"
	key              = "87be2e87e7e01ad6e03b56323bca6b5f"
	// mid = "202366100"
	// key = "c385fe7029344aef826d8112625b2625x"
)

type PayReq struct {
	MerchantID string `json:"merchant_id"` // 商户号
	No         string `json:"no"`          // 订单号
	Amount     int    `json:"amount"`      // 订单金额, 只能为正整数
	Sign       string `json:"sign"`        // 签名 (无具体说明，可能是空值)
}

type PayResp struct {
	Code int    `json:"code"` // 0:业务请求成功,别的都是失败
	Msg  string `json:"msg"`  // 提示消息
	Data struct {
		No         string `json:"no"`          // 商户订单号
		Tx         string `json:"tx"`          // 平台订单号
		MerchantID string `json:"merchant_id"` // 商户号
		PayURL     string `json:"pay_url"`     // 收银台地址
		PayLink    string `json:"pay_link"`    // upi 付款链接
	} `json:"data"`
}

type PayCallbackReq struct {
	MerchantID  string `json:"merchant_id"`  // 商户号
	Tx          string `json:"tx"`           // 平台唯一订单号
	No          string `json:"no"`           // 商户自己上传的订单号
	Status      int    `json:"status"`       // 支付状态：0:付款中，1：付款成功
	SuccessTime string `json:"success_time"` // 付款时间
	Amount      string `json:"amount"`       // 订单金额，下单时传的值
	UTR         string `json:"utr"`          // upr付款订单号
	Sign        string `json:"sign"`         // 签名，除去sign字段外所有字段md签名
}

type QueryPayReq struct {
	MerchantID string `json:"merchant_id"` // 商户号
	No         string `json:"no"`          // 商户订单号 或者 用户UTR编号
	Sign       string `json:"sign"`        // 签名
}

type QueryPayResp struct {
	Code int    `json:"code"` // 0:业务请求成功，别的都失败
	Msg  string `json:"msg"`  // 提示消息，提示失败
	Data struct {
		Tx          string `json:"tx"`           // 平台订单号
		No          string `json:"no"`           // 商户订单号
		MerchantID  string `json:"merchant_id"`  // 商户号
		Status      string `json:"status"`       // 支付状态：0：付款中, 1：成功
		UTR         string `json:"utr"`          // 无
		Amount      string `json:"amount"`       // 订单金额
		SuccessTime string `json:"success_time"` // 付款成功时间
	} `json:"data"` // 无
}

type WithdrawReq struct {
	MerchantID string `json:"merchant_id"` // 商户号
	No         string `json:"no"`          // 代付订单号
	Amount     string `json:"amount"`      // 收款金额，必须为正整数
	Account    string `json:"account"`     // 收款人帐号
	IfscCode   string `json:"ifsc_code"`   // 银行编码
	Name       string `json:"name"`        // 收款人姓名
	Sign       string `json:"sign"`        // 签名
}

type WithdrawResp struct {
	Code int    `json:"code"` // 0：请求成功，别的失败
	Msg  string `json:"msg"`  // 提示消息
	Data struct {
		Tx     string `json:"tx"`     // 平台订单号
		No     string `json:"no"`     // 商户订单号
		Status int    `json:"status"` // 代付状态：1: 代付中, 2: 处理中, 3: 成功, 4: 失败, 5: 反转
	} `json:"data"`
}

type WithdrawCallbackReq struct {
	MerchantID  string `json:"merchant_id"`  // 商户号
	Tx          string `json:"tx"`           // 平台唯一订单号
	No          string `json:"no"`           // 商户自己上传的订单号
	Amount      string `json:"amount"`       // 订单金额，下单时传的值
	Status      int    `json:"status"`       // 支付状态：1:代付, 2:处理中, 3:成功, 4:失败, 5:反转
	SuccessTime string `json:"success_time"` // 付款时间
	Remark      string `json:"remark"`       // 代付失败时的原因
	UTR         string `json:"utr"`          // UTR号
	Sign        string `json:"sign"`         // 签名，除去sign字段外所有字段md签名
}

type QueryWithdrawReq struct {
	MerchantID string `json:"merchant_id"` // 商户号
	No         string `json:"no"`          // 商户自己的订单号
	Sign       string `json:"sign"`        // 签名
}

type QueryWithdrawResp struct {
	Code int    `json:"code"` // 0:业务请求成功，别的都失败. 不能根据这个判断代付订单状态
	Msg  string `json:"msg"`  // 提示消息，提示失败
	Data struct {
		Tx         string `json:"tx"`          // 平台订单号
		No         string `json:"no"`          // 商户订单号
		MerchantID string `json:"merchant_id"` // 商户号
		PayNo      string `json:"pay_no"`      // UTR
		Status     string `json:"status"`      // 代付状态：1：代付中, 2: 处理中, 3: 成功, 4: 失败, 5: 反转
	} `json:"data"` // 无
}
