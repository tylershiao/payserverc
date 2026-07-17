package richpay

import (
	"sync"
	"time"
)

const (
	baseURL          = "https://uat-pf.dkpaypro.com"
	tokenURL         = "/auth/def/oauth/token"
	refreshTokenURL  = "/auth/def/oauth/refresh_token"
	payURL           = "/open/payin"
	withdrawURL      = "/open/payout"
	queryWithdrawURL = "/open/payout/findByNo"
	queryPayURL      = "/open/payin/findByNo"
	mid              = "236"
	privateKey       = "MIIBSgIBADCCASsGByqGSM44BAEwggEeAoGBANjZuuQIywjCbV2OilNYQjW618jcwicZBycv9e8Q/ikHoiy1ShvVlKmZWwhSZSz4dAXlniHcDfE6Nt4fPPoBxfwfzgy0T1IfsdrMmpy02tXO+ZyLKSRk6bNHD4y+9HIqeaR4QbADR0/TdSVpTRcL9S3qV8UsZqBmu+xnyzsVFy6NAhUA8GZwmQl14WTu6iGJeMsJvg0WErUCgYBXpHCOx74DgMGh7gGnxVXcX9C1xb3HxSLsixsTfrxGzLBHL1Jsa+wE4Y/Ff1Zga8YGLTcEwCN9YO5fRa17Ry1TPsCnxoobxdPK8+gNcjP8u1DYvT9HHeh9NwbsCX0gs77JhvgDLD8TXkE8wnt7CwZZK6mEJkoLu5rJgPcMyqDQdwQWAhQ8ItMfC4ofwIyZzdg++oYplXPrNQ=="
	publicKey        = "MIIBuDCCASwGByqGSM44BAEwggEfAoGBAP1/U4EddRIpUt9KnC7s5Of2EbdSPO9EAMMeP4C2USZpRV1AIlH7WT2NWPq/xfW6MPbLm1Vs14E7gB00b/JmYLdrmVClpJ+f6AR7ECLCT7up1/63xhv4O1fnxqimFQ8E+4P208UewwI1VBNaFpEy9nXzrith1yrv8iIDGZ3RSAHHAhUAl2BQjxUjC8yykrmCouuEC/BYHPUCgYEA9+GghdabPd7LvKtcNrhXuXmUr7v6OuqC+VdMCz0HgmdRWVeOutRZT+ZxBxCBgLRJFnEj6EwoFhO3zwkyjMim4TwWeotUfI0o4KOuHiuzpnWRbqN/C/ohNWLx+2J6ASQ7zKTxvqhRkImog9/hWuWfBpKLZl6Ae1UlZAFMO/7PSSoDgYUAAoGBAKl3DccrLszwepoc3uy+kzQywHcoOJ6UPsGsuGYIjVfFn50BotjGo6niEMBHlrp39cGcZJrhHYyishS1Z2rNNAKx3UvEPwkdeJSiNpx2xSk+uM3x9srfk7uIoZUnv3680Koua6EFFSYmSzotYE0RZChx1wQ5K+xG1F3qcQF68K6d"
	secret           = "0nhH9O8X"
)

var (
	accessToken    = ""
	refreshToken   = ""
	tokenTimer     *time.Timer
	tokenTimerLock = &sync.Mutex{}
)

type TokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Code         int    `json:"code"`
}

type PayReq struct {
	MerchantID      string `json:"merchant_id"`       // 商户号
	MerchantOrderNo string `json:"merchant_order_no"` // 商户订单号
	OrderAmount     string `json:"order_amount"`      // 订单金额（卢比）
	Email           string `json:"email"`             // 邮箱
	Name            string `json:"name"`              // 姓名（英文）
	Phone           string `json:"phone"`             // 手机号
	DeeplinkSwitch  string `json:"deeplink_switch"`   // 深度链接开关（0 开，1 关）
	NotifyURL       string `json:"notify_url"`        // 回调地址（可选）
	NonceStr        string `json:"noncestr"`          // 随机数（6位）
	Action          string `json:"action"`            // 操作类型
	Timestamp       string `json:"timestamp"`         // 时间戳（秒）
	Sign            string `json:"sign"`              // DSA签名
}

type PayResp struct {
	Code    int    `json:"code"` // 请求返回码（200：成功）
	BusCode string `json:"busCode"`
	Success bool   `json:"success"` // 请求是否成功
	Msg     string `json:"msg"`     // 请求提示
	Data    struct {
		OrderNo         string  `json:"order_no"`          // 平台单号
		MerchantOrderNo string  `json:"merchant_order_no"` // 商户单号
		OrderAmount     float64 `json:"order_amount"`      // 订单金额
		PayURL          string  `json:"pay_url"`           // 支付 UPI 链接
		PayURLH5        string  `json:"pay_url_H5"`        // 支付 H5 链接
	} `json:"data"` // 返回数据
}

type PayCallbackReq struct {
	Status          string `json:"status"`            // 0：待支付；1：处理中；2：成功；3：失败；4：拒绝；5：撤销
	Message         string `json:"message"`           // 信息提示
	OrderAmount     string `json:"order_amount"`      // 实际支付金额（单位：卢比）
	OrderNo         string `json:"order_no"`          // 平台订单号
	MerchantOrderNo string `json:"merchant_order_no"` // 商户订单号
	UTR             string `json:"utr"`               // 银行交易号（UTR）
	NonceStr        string `json:"noncestr"`          // 随机码
	Timestamp       string `json:"timestamp"`         // 时间戳
	Sign            string `json:"sign"`              // DSA签名
}

type QueryPayReq struct {
	MerchantOrderNo string `json:"merchant_order_no"` // 商户订单号
}

type QueryPayResp struct {
	Code    int    `json:"code"` // 返回码（200表示成功）
	BusCode string `json:"busCode"`
	Success bool   `json:"success"` // 是否成功
	Msg     string `json:"msg"`     // 提示信息
	Data    struct {
		UTR             string  `json:"utr"`               // 银行交易参考号
		MerchantOrderNo string  `json:"merchant_order_no"` // 商户订单号
		OrderAmount     float64 `json:"order_amount"`      // 订单金额
		Status          string  `json:"status"`            // 订单状态：0：待支付；1：处理中；2：成功；3：失败；4：拒绝；5：撤销
	} `json:"data"` // 返回数据
}

type WithdrawReq struct {
	MerchantID      string `json:"merchant_id"`       // 商户号
	MerchantOrderNo string `json:"merchant_order_no"` // 商户订单号，保证唯一性
	OrderAmount     int64  `json:"order_amount"`      // 订单金额（单位：卢比）
	Name            string `json:"name"`              // 姓名（英文，5~20个字符）
	Email           string `json:"email"`             // 邮箱（推荐使用 Gmail）
	Phone           string `json:"phone"`             // 手机号（真实或符合格式）
	NotifyURL       string `json:"notify_url"`        // 回调地址（可选）
	NonceStr        string `json:"noncestr"`          // 随机数（6位）
	Action          string `json:"action"`            // 操作类型
	Timestamp       string `json:"timestamp"`         // 时间戳（秒）
	Sign            string `json:"sign"`              // DSA 签名
	Type            string `json:"type"`              // 支付类型：BANK/UPI（可选）
	VPA             string `json:"vpa"`               // UPI ID（type=UPI 时必填）
	Account         string `json:"account"`           // 银行账号（type=BANK 时必填）
	IFSC            string `json:"ifsc"`              // 银行 IFSC（type=BANK 时必填）
}

type WithdrawResp struct {
	Code    int    `json:"code"` // 请求返回码，200 表示成功
	BusCode string `json:"busCode"`
	Success bool   `json:"success"` // 请求是否成功
	Msg     string `json:"msg"`     // 提示信息
	Data    struct {
		OrderNo         string  `json:"order_no"`          // 平台订单号
		MerchantOrderNo string  `json:"merchant_order_no"` // 商户订单号
		OrderAmount     float64 `json:"order_amount"`      // 订单金额
		Status          string  `json:"status"`            // 订单状态
	} `json:"data"` // 返回数据（可能为空）
}

type WithdrawCallbackReq struct {
	Status          string `json:"status"`            // 订单状态：0：待支付；1：处理中；2：成功；3：失败；4：拒绝；5：撤销
	Message         string `json:"message"`           // 提示信息
	OrderAmount     string `json:"order_amount"`      // 订单金额（单位：卢比，字符串格式）
	OrderNo         string `json:"order_no"`          // 平台订单号
	MerchantOrderNo string `json:"merchant_order_no"` // 商户订单号
	UTR             string `json:"utr"`               // 银行交易参考号
	NonceStr        string `json:"noncestr"`          // 随机码
	Timestamp       string `json:"timestamp"`         // 时间戳
	Sign            string `json:"sign"`              // DSA 签名
}

type QueryWithdrawReq struct {
	MerchantOrderNo string `json:"merchant_order_no"` // 商户订单号
}

type QueryWithdrawResp struct {
	Code    int    `json:"code"` // 请求返回码，200表示成功
	BusCode string `json:"busCode"`
	Success bool   `json:"success"` // 请求是否成功
	Msg     string `json:"msg"`     // 请求提示
	Data    struct {
		UTR             string  `json:"utr"`               // 银行交易参考号
		MerchantOrderNo string  `json:"merchant_order_no"` // 商户订单号
		OrderAmount     float64 `json:"order_amount"`      // 订单金额
		Status          string  `json:"status"`            // 订单状态：0：待支付；1：处理中；2：成功；3：失败；4：拒绝；5：撤销
	} `json:"data"` // 请求返回数据
}
