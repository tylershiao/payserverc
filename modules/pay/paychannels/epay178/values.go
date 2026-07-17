package epay178

const (
	baseURL          = "https://api.epay178.com/api/manager/epayController"
	payURL           = "/payJsonV2/create"
	queryPayURL      = "/payJsonV2/get"
	withdrawURL      = "/withdrawJsonV2/create"
	queryWithdrawURL = "/withdrawJsonV2/get"
	balURL           = "/merchantJsonV2/getBalance"
	mid              = "1409459818309107712"
	key              = "B2D0FDD39B08521FA01C2ADE9BBFB577"
)

type PayReq struct {
	MchID      string `json:"mch_id"`
	OutTradeNo string `json:"out_trade_no"`
	Type       string `json:"type"`
	Amount     int64  `json:"amount"`
	Currency   string `json:"currency"`
	Subject    string `json:"subject"`
	Extra      string `json:"extra"`
	NotifyURL  string `json:"notify_url"`
	Nonce      string `json:"nonce"`
	Sign       string `json:"sign"`
}

type PayExtra struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

type PayResp struct {
	ReasonCode   string `json:"reasoncode"`
	ReasonRemark string `json:"reasonremark"`
	Data         struct {
		Object        string `json:"object"`
		TransactionID string `json:"transaction_id"`
		OutTradeNo    string `json:"out_trade_no"`
		Amount        string `json:"amount"`
		Currency      string `json:"currency"`
		CreateAt      string `json:"create_at"`
		PayURL        string `json:"pay_url"`
	} `json:"data"`
	Nonce string `json:"nonce"`
	Sign  string `json:"sign"`
}

type PayCallbackReq struct {
	Object        string `json:"object"`
	OutTradeNo    string `json:"out_trade_no"`
	TransactionID string `json:"transaction_id"`
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	CreateAt      string `json:"create_at"`
	Paid          string `json:"paid"`
	PaidAt        string `json:"paid_at"`
	UTR           string `json:"utr"`
	Nonce         string `json:"nonce"`
	Sign          string `json:"sign"`
}

type QueryPayReq struct {
	MchID         string `json:"mch_id"`
	OutTradeNo    string `json:"out_trade_no"`
	TransactionID string `json:"transaction_id"`
	Nonce         string `json:"nonce"`
	Sign          string `json:"sign"`
}

type QueryPayResp struct {
	ReasonCode   string `json:"reasoncode"`
	ReasonRemark string `json:"reasonremark"`
	Data         struct {
		Object        string `json:"object"`
		OutTradeNo    string `json:"out_trade_no"`
		TransactionID string `json:"transaction_id"`
		Amount        string `json:"amount"`
		Currency      string `json:"currency"`
		CreateAt      string `json:"create_at"`
		Paid          string `json:"paid"`
		PaidAt        string `json:"paid_at"`
		UTR           string `json:"utr"`
	} `json:"data"`
	Nonce string `json:"nonce"`
	Sign  string `json:"sign"`
}

type WithdrawExtra struct {
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	TransType string `json:"transType"`
	TransMode string `json:"transMode"`
	Account   string `json:"account"`
	IFSC      string `json:"ifsc"`
	BankName  string `json:"bankName"`
	VPA       string `json:"vpa"`
}

type WithdrawReq struct {
	MchID      string `json:"mch_id"`
	OutTradeNo string `json:"out_trade_no"`
	Type       string `json:"type"`
	Amount     int64  `json:"amount"`
	Currency   string `json:"currency"`
	Subject    string `json:"subject"`
	Extra      string `json:"extra"`
	NotifyURL  string `json:"notify_url"`
	Nonce      string `json:"nonce"`
	Sign       string `json:"sign"`
}

type WithdrawResp struct {
	ReasonCode   string `json:"reasoncode"`
	ReasonRemark string `json:"reasonremark"`
	Data         struct {
		Object        string `json:"object"`
		TransactionID string `json:"transaction_id"`
		OutTradeNo    string `json:"out_trade_no"`
		Amount        string `json:"amount"`
		Currency      string `json:"currency"`
		CreateAt      string `json:"create_at"`
	} `json:"data"`
	Nonce string `json:"nonce"`
	Sign  string `json:"sign"`
}

type WithdrawCallbackReq struct {
	Object        string `json:"object"`
	MchID         string `json:"mch_id"`
	OutTradeNo    string `json:"out_trade_no"`
	TransactionID string `json:"transaction_id"`
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	Paid          string `json:"paid"`
	Result        string `json:"result"`
	UTR           string `json:"utr"`
	Remark        string `json:"remark"`
	CreateAt      string `json:"create_at"`
	PaidAt        string `json:"paid_at"`
	Nonce         string `json:"nonce"`
	Sign          string `json:"sign"`
}

type QueryWithdrawReq struct {
	MchID         string `json:"mch_id"`
	OutTradeNo    string `json:"out_trade_no"`
	TransactionID string `json:"transaction_id"`
	Nonce         string `json:"nonce"`
	Sign          string `json:"sign"`
}

type QueryWithdrawResp struct {
	ReasonCode   string `json:"reasoncode"`
	ReasonRemark string `json:"reasonremark"`
	Data         struct {
		Object        string `json:"object"`
		MchID         string `json:"mch_id"`
		OutTradeNo    string `json:"out_trade_no"`
		TransactionID string `json:"transaction_id"`
		Amount        string `json:"amount"`
		Currency      string `json:"currency"`
		Paid          string `json:"paid"`
		Result        string `json:"result"`
		PaidAt        string `json:"paid_at"`
		CreateAt      string `json:"create_at"`
		UTR           string `json:"utr"`
	} `json:"data"`
	Nonce string `json:"nonce"`
	Sign  string `json:"sign"`
}
