package megipay

const (
	baseURL         = "https://in1-api.megipay.com"
	payURL          = "/api/v1/payApi/CreatePayInOrder"
	withdrawURL     = "/api/v1/payApi/CreatePayOutOrder"
	queryURL        = "/api/v1/payApi/QueryOrder"
	balanceURL      = "/api/v1/payApi/QueryBalance"
	mid             = 10371
	appID           = 10371
	payCodePayIn    = 0
	payCodePayOut   = 1
	payMethod       = "india"
	bankPayType     = "1"
	upiPayType      = "8"
	bankAccountType = "BANK"
	upiAccountType  = "UPI"
	bankCode        = "INR_BANK"
	upiBankCode     = "INR_UPI"
	merchantSignKey = "hyQPTJ2RyXNr7rfQ"
)

type PayReq struct {
	TradeNo      int    `json:"trade_no"`
	AppID        int    `json:"app_id"`
	PayCode      int    `json:"pay_code"`
	PayMethod    string `json:"pay_method"`
	Price        int64  `json:"price"`
	OrderNo      string `json:"order_no"`
	SuccessURL   string `json:"success_url"`
	FailURL      string `json:"fail_url"`
	PayNoticeURL string `json:"pay_notice_url"`
	UserID       string `json:"user_id"`
	UserIP       string `json:"user_ip"`
	Attach       string `json:"attach"`
	Sign         string `json:"sign"`
	Timestamp    string `json:"timestamp"`
}

type PayAttach struct {
	AccountNo    string `json:"account_no,omitempty"`
	AccountType  string `json:"account_type,omitempty"`
	Name         string `json:"name,omitempty"`
	IdentifyType string `json:"identify_type,omitempty"`
	IdentifyNum  string `json:"identify_num,omitempty"`
}

type PayResp struct {
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
	PayURL     string `json:"pay_url"`
	QRCode     string `json:"qr_code"`
	OrderNo    string `json:"order_no"`
	DisOrderNo string `json:"dis_order_no"`
	CreateTime int64  `json:"create_time"`
	PayInfo    string `json:"pay_info"`
	Sign       string `json:"sign"`
}

type PayCallbackReq struct {
	TradeNo    int    `json:"trade_no"`
	Status     int    `json:"status"`
	OrderNo    string `json:"order_no"`
	DisOrderNo string `json:"dis_order_no"`
	OrderPrice int64  `json:"order_price"`
	RealPrice  int64  `json:"real_price"`
	NtiTime    int64  `json:"nti_time"`
	Payer      string `json:"payer"`
	PayInfo    string `json:"pay_info"`
	CreateTime int64  `json:"create_time"`
	Sign       string `json:"sign"`
}

type WithdrawReq struct {
	TradeNo      int    `json:"trade_no"`
	OrderNo      string `json:"order_no"`
	AppID        int    `json:"app_id"`
	PayCode      int    `json:"pay_code"`
	Price        int64  `json:"price"`
	AccountNo    string `json:"account_no"`
	AccountType  string `json:"account_type"`
	AccountName  string `json:"account_name"`
	BankCode     string `json:"bank_code"`
	IdentifyType string `json:"identify_type"`
	IdentifyNum  string `json:"identify_num"`
	PayNoticeURL string `json:"pay_notice_url"`
	Attach       string `json:"attach"`
	UserIP       string `json:"user_ip"`
	Sign         string `json:"sign"`
	Timestamp    string `json:"timestamp"`
}

type WithdrawAttach struct {
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	BankName string `json:"bank_name,omitempty"`
	PayType  string `json:"pay_type,omitempty"`
}

type WithdrawResp struct {
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
	DisOrderNo string `json:"dis_order_no"`
	OrderNo    string `json:"order_no"`
	Status     int    `json:"status"`
	CreateTime int64  `json:"create_time"`
	Sign       string `json:"sign"`
}

type WithdrawCallbackReq struct {
	TradeNo    int    `json:"trade_no"`
	OrderNo    string `json:"order_no"`
	DisOrderNo string `json:"dis_order_no"`
	OrderPrice int64  `json:"order_price"`
	Fee        int64  `json:"fee"`
	Status     int    `json:"status"`
	PayInfo    string `json:"pay_info"`
	Remark     string `json:"remark"`
	CreateTime int64  `json:"create_time"`
	Sign       string `json:"sign"`
	NtiTime    int64  `json:"nti_time"`
}

type QueryReq struct {
	OrderType  string `json:"order_type"`
	TradeNo    int    `json:"trade_no"`
	AppID      int    `json:"app_id"`
	DisOrderNo string `json:"dis_order_no"`
	OrderNo    string `json:"order_no"`
	Sign       string `json:"sign"`
}

type QueryResp struct {
	Code        int    `json:"code"`
	Msg         string `json:"msg"`
	TradeNo     int    `json:"trade_no"`
	RealPrice   int64  `json:"real_price"`
	Status      int    `json:"status"`
	SuccessTime int64  `json:"success_time"`
	OrderNo     string `json:"order_no"`
	DisOrderNo  string `json:"dis_order_no"`
	Remark      string `json:"remark"`
	Fee         int64  `json:"fee"`
	CreateTime  int64  `json:"create_time"`
	Payer       string `json:"payer"`
	PayInfo     string `json:"pay_info"`
	Sign        string `json:"sign"`
	UTR2        string `json:"utr2"`
}

type BalanceReq struct {
	TradeNo int    `json:"trade_no"`
	AppID   int    `json:"app_id"`
	Sign    string `json:"sign"`
}

type BalanceResp struct {
	Code          int    `json:"code"`
	Msg           string `json:"msg"`
	Balance       int64  `json:"balance"`
	BalanceFrozen int64  `json:"balance_frozen"`
	Sign          string `json:"sign"`
}
