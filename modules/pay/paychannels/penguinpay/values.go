package penguinpay

const (
	baseURL          = "https://api.penguinpay-app.com"
	payURL           = "/api/pay/create/order"
	queryPayURL      = "/api/pay/query/order"
	withdrawURL      = "/api/pay/create/order"
	queryWithdrawURL = "/api/pay/query/order"
	balURL           = "/api/pay/merchant/balance"
	mid              = 199
	// appID = "S1773216155928"
	key = "7df580b6bb130d8bf62f66864ed9bc11"
)

type PayReq struct {
	Type            int     `json:"type"`
	MerchantID      int     `json:"merchantId"`
	MerchantOrderNo string  `json:"merchantOrderNo"`
	OrderAmount     float64 `json:"orderAmount"`
	NotifyURL       string  `json:"notifyUrl"`
	// ReturnURL       string  `json:"returnUrl"`
	// Extra           string  `json:"extra"`
	Sign string `json:"sign"`
}

type PayResp struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Data    PaymentResponseData `json:"data"`
}

type PaymentResponseData struct {
	OrderNo string `json:"orderNo"`
	// OrderAmount int    `json:"orderAmount"`
	PayURL   string `json:"payUrl"`
	UPI      string `json:"upi"`
	Deeplink string `json:"deeplink"`
}

type PayCallbackReq struct {
	Code           int     `json:"code"`
	MchId          int     `json:"mchId"`
	MchOrderNo     string  `json:"mchOrderNo"`
	OrderAmount    float64 `json:"orderAmount"`
	Utr            string  `json:"utr"`
	OrderNo        string  `json:"orderNo"`
	PaySuccessTime string  `json:"paySuccessTime"`
	Message        string  `json:"message"`
	Extra          string  `json:"extra"`
	Sign           string  `json:"sign"`
}

type QueryPayReq struct {
	Type            int    `json:"type"`
	MerchantId      int    `json:"merchantId"`
	MerchantOrderNo string `json:"merchantOrderNo"`
	Sign            string `json:"sign"`
}

/*
状态	描述	备注
0	预生成	商户方推单、系统接收后的初始状态, 此时未为该订单分配UPI
1	支付中	当该订单对应的收银台地址被打开，为该订单分配UPI
2	支付成功	系统确认已收款，订单完成
3	超时	订单未完成并且已经超时
4	处理成功	在2的基础上回调成功（商户方确认订单完成后）
*/
type QueryPayResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		MerchantID int `json:"merchantId"`
		// OrderAmount     int    `json:"orderAmount"`
		MerchantOrderNo string `json:"merchantOrderNo"`
		Status          int    `json:"status"`
		CreateTime      string `json:"createTime"`
		PayTime         string `json:"payTime"`
		NotifyTime      string `json:"notifyTime"`
	} `json:"data"`
}

type WithdrawReq struct {
	Type            int     `json:"type"`
	MerchantID      int     `json:"merchantId"`
	MerchantOrderNo string  `json:"merchantOrderNo"`
	OrderAmount     float64 `json:"orderAmount"`
	AccountName     string  `json:"accountName"`
	CardNumber      string  `json:"cardNumber"`
	IFSC            string  `json:"ifsc"`
	BankName        string  `json:"bankName"`
	UPI             string  `json:"upi"`
	NotifyURL       string  `json:"notifyUrl"`
	Phone           string  `json:"phone"`
	Email           string  `json:"email"`
	// Extra           string  `json:"extra"`
	Sign string `json:"sign"`
}

type WithdrawResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		OrderNo string `json:"orderNo"`
		// OrderAmount float64 `json:"orderAmount"`
	} `json:"data"`
}

type WithdrawCallbackReq struct {
	Code       int    `json:"code"`
	MchId      int    `json:"mchId"`
	MchOrderNo string `json:"mchOrderNo"`
	// OrderAmount    float64 `json:"orderAmount"`
	Utr            string `json:"utr"`
	OrderNo        string `json:"orderNo"`
	PaySuccessTime string `json:"paySuccessTime"`
	Message        string `json:"message"`
	Extra          string `json:"extra"`
	Sign           string `json:"sign"`
}

type QueryWithdrawReq struct {
	Type            int    `json:"type"`
	MerchantId      int    `json:"merchantId"`
	MerchantOrderNo string `json:"merchantOrderNo"`
	Sign            string `json:"sign"`
}

/*
状态	描述	备注
0	支付中	商户方推单、系统接收后的初始状态
1	支付成功	出款完成、订单完成
2	支付失败	支付失败(此状态为设计预留，目前不会返回)
3	驳回	订单被驳回
4	处理成功	在1的基础上回调成功（商户方确认订单完成后）
*/
type QueryWithdrawResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		MerchantID int `json:"merchantId"`
		// OrderAmount     int    `json:"orderAmount"`
		MerchantOrderNo string `json:"merchantOrderNo"`
		Status          int    `json:"status"`
		CreateTime      string `json:"createTime"`
		PayTime         string `json:"payTime"`
		NotifyTime      string `json:"notifyTime"`
	} `json:"data"`
}
