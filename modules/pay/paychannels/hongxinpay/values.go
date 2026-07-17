package hongxinpay

const (
	baseURL          = "https://pay.lbspay.online"
	payURL           = "/api/open/payment/payin/create"
	queryPayURL      = "/api/open/payment/payin/queryOrderStatus"
	withdrawURL      = "/api/open/payment/payout/create"
	queryWithdrawURL = "/api/open/payment/payout/queryOrderStatus"

	merchantCode = "R9GJW"
	key          = "WYVK9T19MXQ9945YUA4GBA8PQYKAGXMJSVP44UXA5G37SHUM3NPQQBK1VW2G8IW7"
)

type PayReq struct {
	MerchantCode    string  `json:"merchantCode"`
	MerchantOrderNo string  `json:"merchantOrderNo"`
	OrderAmount     float64 `json:"orderAmount"`
	BuyerPhone      string  `json:"buyerPhone"`
	BuyerName       string  `json:"buyerName"`
	BuyerEmail      string  `json:"buyerEmail"`
	NotifyUrl       string  `json:"notifyUrl"`
	Sign            string  `json:"sign"`
}

type PayResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		MerchantOrderNo string  `json:"merchantOrderNo"`
		OrderNo         string  `json:"orderNo"`
		OrderAmount     float64 `json:"orderAmount"`
		PayUrl          string  `json:"payUrl"`
	} `json:"data"`
}

type PayCallbackReq struct {
	OrderStatus     string  `json:"orderStatus"`
	OrderAmount     float64 `json:"orderAmount"`
	MerchantOrderNo string  `json:"merchantOrderNo"`
	Fee             float64 `json:"fee"`
	OrderNo         string  `json:"orderNo"`
	CreateTime      string  `json:"createTime"`
	Sign            string  `json:"sign"`
}

type QueryPayReq struct {
	MerchantCode    string `json:"merchantCode"`
	MerchantOrderNo string `json:"merchantOrderNo"`
	Sign            string `json:"sign"`
}

type QueryPayResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		OrderStatus     string  `json:"orderStatus"`
		OrderAmount     float64 `json:"orderAmount"`
		MerchantOrderNo string  `json:"merchantOrderNo"`
		Fee             float64 `json:"fee"`
		OrderNo         string  `json:"orderNo"`
		CreateTime      string  `json:"createTime"`
		Sign            string  `json:"sign"`
	} `json:"data"`
}

type PayoutReq struct {
	MerchantCode    string  `json:"merchantCode"`
	MerchantOrderNo string  `json:"merchantOrderNo"`
	OrderAmount     float64 `json:"orderAmount"`
	BuyerPhone      string  `json:"buyerPhone"`
	BuyerName       string  `json:"buyerName"`
	PayoutType      string  `json:"payoutType"`
	BuyerEmail      string  `json:"buyerEmail"`
	BuyerBankNumber string  `json:"buyerBankNumber"`
	BuyerIfscCode   string  `json:"buyerIfscCode"`
	NotifyUrl       string  `json:"notifyUrl"`
	Sign            string  `json:"sign"`
}

type PayoutResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		MerchantOrderNo string  `json:"merchantOrderNo"`
		OrderNo         string  `json:"orderNo"`
		OrderAmount     float64 `json:"orderAmount"`
		OrderStatus     string  `json:"orderStatus"`
	} `json:"data"`
}

type PayoutCallbackReq struct {
	OrderStatus     string  `json:"orderStatus"`
	OrderAmount     float64 `json:"orderAmount"`
	MerchantOrderNo string  `json:"merchantOrderNo"`
	OrderNo         string  `json:"orderNo"`
	Fee             float64 `json:"fee"`
	CreateTime      string  `json:"createTime"`
	Utr             string  `json:"utr"`
	Sign            string  `json:"sign"`
}

type QueryPayoutReq struct {
	MerchantCode    string `json:"merchantCode"`
	MerchantOrderNo string `json:"merchantOrderNo"`
	Sign            string `json:"sign"`
}

type QueryPayoutResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		OrderStatus     string  `json:"orderStatus"`
		OrderAmount     float64 `json:"orderAmount"`
		MerchantOrderNo string  `json:"merchantOrderNo"`
		OrderNo         string  `json:"orderNo"`
		Fee             float64 `json:"fee"`
		CreateTime      string  `json:"createTime"`
		Utr             string  `json:"utr"`
		Sign            string  `json:"sign"`
	} `json:"data"`
}
