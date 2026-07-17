package dkpay

const (
	baseURL          = "https://gepopaypro.com"
	payURL           = "/open/nax/payin"
	withdrawURL      = "/open/nax/payout"
	queryWithdrawURL = "/open/nax/payout/findByNo"
	queryPayURL      = "/open/nax/payin/findByNo"
	// mid              = "105"
	// privateKey       = "MIIBSwIBADCCASwGByqGSM44BAEwggEfAoGBAP1/U4EddRIpUt9KnC7s5Of2EbdSPO9EAMMeP4C2USZpRV1AIlH7WT2NWPq/xfW6MPbLm1Vs14E7gB00b/JmYLdrmVClpJ+f6AR7ECLCT7up1/63xhv4O1fnxqimFQ8E+4P208UewwI1VBNaFpEy9nXzrith1yrv8iIDGZ3RSAHHAhUAl2BQjxUjC8yykrmCouuEC/BYHPUCgYEA9+GghdabPd7LvKtcNrhXuXmUr7v6OuqC+VdMCz0HgmdRWVeOutRZT+ZxBxCBgLRJFnEj6EwoFhO3zwkyjMim4TwWeotUfI0o4KOuHiuzpnWRbqN/C/ohNWLx+2J6ASQ7zKTxvqhRkImog9/hWuWfBpKLZl6Ae1UlZAFMO/7PSSoEFgIUHhk46MaXhu9X287/tCNZpMAGPno="
	mid        = "767"
	privateKey = "MIIBSwIBADCCASwGByqGSM44BAEwggEfAoGBAO9B32KYsobruHJoMFiLwE5I2+Z5ZIQ40ZVn3YtGXngBp+vIq1otL/QFPsFmsk1eXLSyXjnaBND4xEOR2oQOdHCH+dJAYsX4R+m3uGpZYbG4QHNsQJFsadeMNwKVA7+nN9ts0LPEmbkkSE5ykBq5Qkqvy+AyXzipzqUD9Z7N0wULAhUAnFSOQnU7NPPlyLAakH6N8QGT/YMCgYEAq0XhnQn28nVJc7bqfGU/T9132xcJRZx2VHOylxov3v/EidH740f0tZnM9/VH/n+uh+K5icF01iP8E7E+8W7gnWzwWO6JeNlAzc1qMC9Y8La8y9seVhLGSc/jhbNPPuFEdoBjeL1fDYzEMSkHZxYuptaYR1dL+/lUYchPnmqJokMEFgIUGu8bkEURXskCiczN65HeEQksqbM="
	publicKey  = "MIIBuDCCASwGByqGSM44BAEwggEfAoGBAP1/U4EddRIpUt9KnC7s5Of2EbdSPO9EAMMeP4C2USZpRV1AIlH7WT2NWPq/xfW6MPbLm1Vs14E7gB00b/JmYLdrmVClpJ+f6AR7ECLCT7up1/63xhv4O1fnxqimFQ8E+4P208UewwI1VBNaFpEy9nXzrith1yrv8iIDGZ3RSAHHAhUAl2BQjxUjC8yykrmCouuEC/BYHPUCgYEA9+GghdabPd7LvKtcNrhXuXmUr7v6OuqC+VdMCz0HgmdRWVeOutRZT+ZxBxCBgLRJFnEj6EwoFhO3zwkyjMim4TwWeotUfI0o4KOuHiuzpnWRbqN/C/ohNWLx+2J6ASQ7zKTxvqhRkImog9/hWuWfBpKLZl6Ae1UlZAFMO/7PSSoDgYUAAoGBAKl3DccrLszwepoc3uy+kzQywHcoOJ6UPsGsuGYIjVfFn50BotjGo6niEMBHlrp39cGcZJrhHYyishS1Z2rNNAKx3UvEPwkdeJSiNpx2xSk+uM3x9srfk7uIoZUnv3680Koua6EFFSYmSzotYE0RZChx1wQ5K+xG1F3qcQF68K6d"
)

type PayReq struct {
	MerchantID      int64  `json:"merchant_id"`
	MerchantOrderNo string `json:"merchant_order_no"`
	OrderAmount     string `json:"order_amount"`
	Email           string `json:"email"`
	Name            string `json:"name"`
	Phone           string `json:"phone"`
	DeeplinkSwitch  string `json:"deeplink_switch"`
	NotifyURL       string `json:"notify_url"`
	NonceStr        string `json:"noncestr"`
	Action          string `json:"action"`
	Timestamp       string `json:"timestamp"`
	Sign            string `json:"sign"`
}

type PayResp struct {
	Code    int    `json:"code"`
	BusCode string `json:"busCode"`
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Data    struct {
		OrderNo         string  `json:"order_no"`
		MerchantOrderNo string  `json:"merchant_order_no"`
		OrderAmount     float64 `json:"order_amount"`
		PayURL          string  `json:"pay_url"`
		PayURLH5        string  `json:"pay_url_H5"`
	} `json:"data"`
}

type PayCallbackReq struct {
	Status          string `json:"status"`
	Message         string `json:"message"`
	OrderAmount     string `json:"order_amount"`
	OrderNo         string `json:"order_no"`
	MerchantOrderNo string `json:"merchant_order_no"`
	UTR             string `json:"utr"`
	NonceStr        string `json:"noncestr"`
	Timestamp       string `json:"timestamp"`
	Sign            string `json:"sign"`
}

type QueryPayReq struct {
	MerchantID      int64  `json:"merchant_id"`
	MerchantOrderNo string `json:"merchant_order_no"`
	Sign            string `json:"sign"`
}

type QueryPayResp struct {
	Code    int    `json:"code"`
	BusCode string `json:"busCode"`
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Data    struct {
		UTR             string  `json:"utr"`
		MerchantOrderNo string  `json:"merchant_order_no"`
		OrderAmount     float64 `json:"order_amount"`
		PaymentAmount   float64 `json:"payment_amount"`
		Status          string  `json:"status"`
	} `json:"data"`
}

type WithdrawReq struct {
	MerchantID      int64  `json:"merchant_id"`
	MerchantOrderNo string `json:"merchant_order_no"`
	OrderAmount     int64  `json:"order_amount"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	NotifyURL       string `json:"notify_url"`
	NonceStr        string `json:"noncestr"`
	Action          string `json:"action"`
	Timestamp       string `json:"timestamp"`
	Sign            string `json:"sign"`
	Type            string `json:"type"`
	VPA             string `json:"vpa"`
	Account         string `json:"account"`
	IFSC            string `json:"ifsc"`
}

type WithdrawResp struct {
	Code    int    `json:"code"`
	BusCode string `json:"busCode"`
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Data    struct {
		OrderNo         string  `json:"order_no"`
		MerchantOrderNo string  `json:"merchant_order_no"`
		OrderAmount     float64 `json:"order_amount"`
		Status          string  `json:"status"`
	} `json:"data"`
}

type WithdrawCallbackReq struct {
	Status          string `json:"status"`
	Message         string `json:"message"`
	OrderAmount     string `json:"order_amount"`
	OrderNo         string `json:"order_no"`
	MerchantOrderNo string `json:"merchant_order_no"`
	UTR             string `json:"utr"`
	NonceStr        string `json:"noncestr"`
	Timestamp       string `json:"timestamp"`
	Sign            string `json:"sign"`
}

type QueryWithdrawReq struct {
	MerchantID      int64  `json:"merchant_id"`
	MerchantOrderNo string `json:"merchant_order_no"`
	Sign            string `json:"sign"`
}

type QueryWithdrawResp struct {
	Code    int    `json:"code"`
	BusCode string `json:"busCode"`
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Data    struct {
		UTR             string  `json:"utr"`
		MerchantOrderNo string  `json:"merchant_order_no"`
		OrderAmount     float64 `json:"order_amount"`
		PaymentAmount   float64 `json:"payment_amount"`
		Status          string  `json:"status"`
	} `json:"data"`
}
