package deepay

const (
	baseURL          = "https://api.deepay.ink"
	payURL           = "/cashin/create"
	queryPayURL      = "/cashin/query"
	withdrawURL      = "/cashout/create"
	queryWithdrawURL = "/cashout/query"
	mid              = 12469
	currency         = "INR"
	payTradeCode     = "UPI"
	bankTradeCode    = "BANK"
	upiTradeCode     = "UPI"
	remarks          = "product info"
)

const privateKey = `-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQC29J4f7OEiWQ0B
nYGLIDFs7873DOeENfK9XyX/ShclIuHZQGM2oY6/FHDGsrOgxChocT7+NzJw+K8G
vVHO0xCKbD7Uw4UxePItHSwKr+/9EmkVfMVnVZ/gkrBtDi0EP4jMQP8IFRq3PT2X
TlV6d5Iw8mRYUT59FYtzMUXLB5bPaLzX4UbqYv7LQfrkn5EvsprmePQMf6lKlTxA
LXipmpROt/o1umxq68r7S21o8hT1rHKjguS0DAarBGkvFR5UkRaFHL8KrsTkKxEk
cg5m1MLZst21TXRcGUW4K7nk2L0D0IMDRA82jk1TcJtH6NLSZHeD+gsNnGDQYfyC
DAcpqT3tAgMBAAECggEADAzyBVmi5+eZqwuRfB/ItvYsrv+97Wq9lv+XZHNK3E//
NEIWHC5pVT38F6UESTM3viQHaX3VPV+ul5u3XGTGkszXw/KE6WH1VvIJnpfsD207
rVtPMqPkffRA2U4KZQc89ByBL9UZXLcUlsy6u0WjmW8TG5tuVEArs+YA74M2dGlt
q2l5Gc4+XTJZ4accflEfghS3UNLGiZc9M8Tho5qX52TTjBxpSCQho7xxJicdJO1f
W0Ty2UU3vsTUjqhAAwqoRBKbGdl7y1Rjr0ZB9bnsPJVMrom+9BljTJOyjJnEbTAB
OcxYvKtpb0rZWFnQPnUQM9vlsQ+k61CaHrTLOQwoOQKBgQDwmIRpMaq7fiOggSNH
4F2ahqTbLBg/7qR8A1dN3QRa8Wut8OEI2567N2gG9qEh5Gxb7mlCLRU9Q8TnHV8W
6wVwzJjIgzOeTvGqrE6NaQ9vnZJ/aoL1T/TOMTsSULrZu7OABqp4qIeTSolj8GZf
aRZbPP/9puR0EMKYnhM/VLinBwKBgQDCq1lHAs+UWFgXfgjo39Jh9WHpV11J3cCq
tZ2+cJxXiBzMT+wB1PW88ZZ9BrioBjp/SI3JYPcfMt4M89Affhq7bDKRAvV/kj+C
AbrRDdRkDCgK4R+mHm18v39a0fkOy2OZg7jekeRc9pS0cZ04ru5GIaXx5+yWvSkC
6JU7P6+iawKBgQDtGsFTEflZ8QOXR+zHBtv4TzzkKd50EvDSrlA5kdia71TZBIwh
OJ6DjiTcx+0gA4auSpQ8/UTzLXg883skkNS5x0fCQ1p5Zv9O3aRERAdMThNzsjoB
G5z8HxRCBoE91ncdcnRJmHyTsV8YuG+LE3tOyxzUksaNg37i/NWnEg0vvQKBgQC1
t+/x6rzNFJ6+50b3459nFAGfKY+K2Xb1o+2pjUNULln7mwuLNjlic/5vnWuYZaAL
mzXH6z0kIX7IhnwcOcJIfJrDZQi4WlPwTx8N8SP3dKWYakf58ONxlO0UYXdliRTj
g+N52KRc500vF/OUqy7/xn86zq26lL+I2IYdtbiFGwKBgQC5mvu4ZVuPUzM4A6NE
/18gGHP6c6KF8Sr2IHm1EK9p6R7ZIQM+eUaFct2eN6hVCcOE5RuPM4FsbfW+K3i5
GIf2wF9hTefr3Kw47xdzYj5boSjF5oyqNCdcM3PRQ++g/5UhLQv/6k2XIyBe+NuR
GJyuX+/Wc0U801DPCRf0pGoyZQ==
-----END PRIVATE KEY-----`

const publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2fJbv9eTMjgKnPvfAYHC
GnCkyDlWcmy0nRRMQYbarfDSFeLTWhjY6QjvGjLlw31fnkTpnCN6N034rADgcwG/
EHfGqc+SQONwvPq8FRRuZ8M0m/09mhIZljmitBublouZF5mV0Ky1hjN2ZBJ7F22m
oU/J+m9N3mmjucNHZWecGqQcM+Ks6d+tb6JlFyvW5W6W13bdVre6t79T5Rd69fPx
4lSyk80zMEfj+63I68kh/+RUOUgeAiIQ5C8xqLL6xd0AvTxihuoN0/cqfAukK+HF
G/kzyIXlEqXnp04wj+vYd+un4NhipjwjoNP90qAMpAYKBYXQ5ouFkGPcgX9VaoiT
HwIDAQAB
-----END PUBLIC KEY-----`

type PayReq struct {
	TradeSn   string  `json:"tradeSn"`
	Currency  string  `json:"currency"`
	Amount    float64 `json:"amount"`
	Name      string  `json:"name,omitempty"`
	Phone     string  `json:"phone,omitempty"`
	Email     string  `json:"email,omitempty"`
	Timestamp int64   `json:"timestamp"`
	Mid       int     `json:"mid"`
	Remarks   string  `json:"remarks"`
	NotifyURL string  `json:"notifyUrl"`
	ReturnURL string  `json:"returnUrl,omitempty"`
	TradeCode string  `json:"tradeCode,omitempty"`
}

type PayResp struct {
	Status int          `json:"status"`
	Msg    string       `json:"msg"`
	Data   *PayRespData `json:"data"`
}

type PayRespData struct {
	Link         string  `json:"link"`
	TradeSn      string  `json:"tradeSn"`
	PlatOrderSn  string  `json:"platOrderSn"`
	Amount       float64 `json:"amount"`
	ActualAmount float64 `json:"actualAmount"`
	DeepLink     string  `json:"deepLink"`
	CreateTime   int64   `json:"createTime"`
	ExpTime      int64   `json:"expTime"`
}

type PayCallbackReq struct {
	Mid          int     `json:"mid"`
	TradeSn      string  `json:"tradeSn"`
	PlatOrderSn  string  `json:"platOrderSn"`
	OrderStatus  string  `json:"orderStatus"`
	Amount       float64 `json:"amount"`
	ActualAmount float64 `json:"actualAmount"`
	Fee          float64 `json:"fee"`
	ErrCode      string  `json:"errCode"`
	ErrMsg       string  `json:"errMsg"`
	FinishTime   int64   `json:"finishTime"`
	BankRefNo    string  `json:"bankRefNo"`
	Timestamp    int64   `json:"timestamp"`
}

type QueryPayReq struct {
	Mid       int    `json:"mid"`
	Timestamp int64  `json:"timestamp"`
	TradeSn   string `json:"tradeSn"`
}

type QueryPayResp struct {
	Status int           `json:"status"`
	Msg    string        `json:"msg"`
	Data   *QueryPayData `json:"data"`
}

type QueryPayData struct {
	Mid          int     `json:"mid"`
	TradeSn      string  `json:"tradeSn"`
	PlatOrderSn  string  `json:"platOrderSn"`
	OrderStatus  string  `json:"orderStatus"`
	Amount       float64 `json:"amount"`
	ActualAmount float64 `json:"actualAmount"`
	Fee          float64 `json:"fee"`
	ErrCode      string  `json:"errCode"`
	ErrMsg       string  `json:"errMsg"`
	FinishTime   int64   `json:"finishTime"`
	BankRefNo    string  `json:"bankRefNo"`
}

type WithdrawReq struct {
	TradeSn   string  `json:"tradeSn"`
	Currency  string  `json:"currency"`
	Amount    float64 `json:"amount"`
	TradeCode string  `json:"tradeCode"`
	Account   string  `json:"account"`
	NotifyURL string  `json:"notifyUrl"`
	Name      string  `json:"name,omitempty"`
	Phone     string  `json:"phone,omitempty"`
	Email     string  `json:"email,omitempty"`
	Timestamp int64   `json:"timestamp"`
	Mid       int     `json:"mid"`
	Remarks   string  `json:"remarks"`
	IFSC      string  `json:"ifsc,omitempty"`
}

type WithdrawResp struct {
	Status int               `json:"status"`
	Msg    string            `json:"msg"`
	Data   *WithdrawRespData `json:"data"`
}

type WithdrawRespData struct {
	TradeSn     string  `json:"tradeSn"`
	PlatOrderSn string  `json:"platOrderSn"`
	Amount      float64 `json:"amount"`
	OrderStatus string  `json:"orderStatus"`
	CreateTime  int64   `json:"createTime"`
}

type WithdrawCallbackReq struct {
	Mid          int     `json:"mid"`
	TradeSn      string  `json:"tradeSn"`
	PlatOrderSn  string  `json:"platOrderSn"`
	OrderStatus  string  `json:"orderStatus"`
	Amount       float64 `json:"amount"`
	ActualAmount float64 `json:"actualAmount"`
	Fee          float64 `json:"fee"`
	ErrCode      string  `json:"errCode"`
	ErrMsg       string  `json:"errMsg"`
	FinishTime   int64   `json:"finishTime"`
	BankRefNo    string  `json:"bankRefNo"`
	Timestamp    int64   `json:"timestamp"`
}

type QueryWithdrawReq struct {
	Mid       int    `json:"mid"`
	TradeSn   string `json:"tradeSn"`
	Timestamp int64  `json:"timestamp"`
}

type QueryWithdrawResp struct {
	Status int                `json:"status"`
	Msg    string             `json:"msg"`
	Data   *QueryWithdrawData `json:"data"`
}

type QueryWithdrawData struct {
	Mid          int     `json:"mid"`
	TradeSn      string  `json:"tradeSn"`
	PlatOrderSn  string  `json:"platOrderSn"`
	OrderStatus  string  `json:"orderStatus"`
	Amount       float64 `json:"amount"`
	ActualAmount float64 `json:"actualAmount"`
	Fee          float64 `json:"fee"`
	ErrCode      string  `json:"errCode"`
	ErrMsg       string  `json:"errMsg"`
	FinishTime   int64   `json:"finishTime"`
	BankRefNo    string  `json:"bankRefNo"`
}
