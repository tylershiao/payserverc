package tkpay

const (
	baseURL          = "https://seabird.world"
	payURL           = "/api/order/pay/create"
	queryPayURL      = "/api/order/pay/query"
	withdrawURL      = "/api/order/withdraw/create/v2"
	queryWithdrawURL = "/api/order/withdraw/query"
	mid              = "202588027"
	key              = "142ebb6649bc456ab1bc1a02c0e2e5de"
	// mid = "202366100"
	// key = "c385fe7029344aef826d8112625b2625x"
)

var privateKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAqSbw9rY58HuqDZ0tyY2xCwQfRiXUI3EasGAMhLbQiMAMBGSZ
QjcKUXHAnTe2YKJHD/Tn7vQpQtB6/ktuECH60lSVkEph+e6uXgE4nWbRkYur5pko
sEHG/fTE+Upo0E9ID1BgElFQ/3gR9iUgOuEiy2EoyJWNyyCrT6OZB/nKSrl7M7EW
Fr6X84NYZWp1pzKcf8cQIxwzuw9yF3Kwa9L43ANgJItrEqZ4zxeJk0/qJEC20UaE
2AmCTkl/HfY6B6j8mNeHSPSyayVcoDitWHiNeUKE3QL2sXog7G+7pljYFZz6M259
w2r/zm8VkD5XJt3BGtjdZtIN0Ckf3C+ygOMY/wIDAQABAoIBABvsODnZwf2/REak
DDNp5vBlGU6InxjMqnW0HUzpKspwPWuaOaxzDop8RyoM0Nfp4eYx7PmGp6sGZE9H
H1+dZFzFJr1fcgYCjxkohXKCTxUCI/nYjZ1u9kzk04E9dcjxVXllsL5FWyQY64og
/c2KbaZ6HvPhsDr4ZR8r5+6nnuuvpl6R6D5qLq6M+V2N9UnCgS7MvLHDoXBFYY3T
C3KAlRdTJKvYqwVy8Dfe7Kmjix8AHuUHbDOZME+KwSsleAYlSbcmBJND3kKIseCT
qWLOZV5n+XTL6/Y9HQvzXBbLd845iDf6MgLF4LAZkjxSAQU9pPWFEsJ/Cdg5HvPi
/AhsMl0CgYEA7mXbfaN5GJ/Bbkkw745ljYWp98q5nQlVH/NHS8FXO1uTpwtrsqMX
NCIkMN/W0RgeZCuBE42BRhCQrwDAoXtm9IoV8DDWVT26Uf9PfePMI19sqzG3Mld6
i+VRp+L5jV0gOy6x1WG4ifNAPicu/vUCr3m9Bczy6QuCxyPwtPvBdqsCgYEAtaQ3
VCymhqIye2/awazy/e6spdnZt1D2KEa6K5zpbbxAWXa0JpvVmmNa0BqapWSCommc
qlq4amoQrl9tb9s59CuVkkitQODqinBRUOeHIQWbZ86v50ZMB/VTMzg+R4WFSjeF
J/NVsJm8gUFfAqWHfF64kY5ydmNY/9Quken6dv0CgYEAxWUPlfta1LnLFz+BRooY
qlMnxmZUxKOqmXb97MsdKKgOHT3GojaKiv9OCi39mEqpN3YSPZdomBQ1SMKNvegI
+qivu58kiZZABY8qSenp6MwllabnssZeUSHDOAknaxbWDkr+stlkTM1ev1gAgNls
rNFXY0stQyvl0HCA6cqfiVsCgYAPbqREwoZjaeRgGGC9k7QRXvXMIXYd3QMM6PVo
PyyV5iiRcURdA2BifY55n8kDxEZVBiQJVSvV3q0MI49osXd7Bzjt9L7DhFho6qH9
/a29YHyTlTvcqlxLUMyAzYZ9DRsF8zGc80+UcXtSf/tSfj3Ek0ebFj8jT81/7scG
HprpAQKBgFtIDm8Y09pGBw+u+ymFU1iNeDaBYzGhlWTz/PXia0Up14ztK2y8MA7p
vtMdzYtzp/UC3F3MVu3ui8LPB8NQnUkGG4iYEV1wnIn0OQ5DVC1847w/WQY/WKyY
fRi98UBfVpwjH4tu9ytZBOsjqnwFkJAog+UxMMAUNZZo20gwyft1
-----END RSA PRIVATE KEY-----`)

var privateKeyCS8 = []byte(`-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCpJvD2tjnwe6oN
nS3JjbELBB9GJdQjcRqwYAyEttCIwAwEZJlCNwpRccCdN7ZgokcP9Ofu9ClC0Hr+
S24QIfrSVJWQSmH57q5eATidZtGRi6vmmSiwQcb99MT5SmjQT0gPUGASUVD/eBH2
JSA64SLLYSjIlY3LIKtPo5kH+cpKuXszsRYWvpfzg1hlanWnMpx/xxAjHDO7D3IX
crBr0vjcA2Aki2sSpnjPF4mTT+okQLbRRoTYCYJOSX8d9joHqPyY14dI9LJrJVyg
OK1YeI15QoTdAvaxeiDsb7umWNgVnPozbn3Dav/ObxWQPlcm3cEa2N1m0g3QKR/c
L7KA4xj/AgMBAAECggEAG+w4OdnB/b9ERqQMM2nm8GUZToifGMyqdbQdTOkqynA9
a5o5rHMOinxHKgzQ1+nh5jHs+YanqwZkT0cfX51kXMUmvV9yBgKPGSiFcoJPFQIj
+diNnW72TOTTgT11yPFVeWWwvkVbJBjriiD9zYptpnoe8+GwOvhlHyvn7qee66+m
XpHoPmouroz5XY31ScKBLsy8scOhcEVhjdMLcoCVF1Mkq9irBXLwN97sqaOLHwAe
5QdsM5kwT4rBKyV4BiVJtyYEk0PeQoix4JOpYs5lXmf5dMvr9j0dC/NcFst3zjmI
N/oyAsXgsBmSPFIBBT2k9YUSwn8J2Dke8+L8CGwyXQKBgQDuZdt9o3kYn8FuSTDv
jmWNhan3yrmdCVUf80dLwVc7W5OnC2uyoxc0IiQw39bRGB5kK4ETjYFGEJCvAMCh
e2b0ihXwMNZVPbpR/09948wjX2yrMbcyV3qL5VGn4vmNXSA7LrHVYbiJ80A+Jy7+
9QKveb0FzPLpC4LHI/C0+8F2qwKBgQC1pDdULKaGojJ7b9rBrPL97qyl2dm3UPYo
RrornOltvEBZdrQmm9WaY1rQGpqlZIKiaZyqWrhqahCuX21v2zn0K5WSSK1A4OqK
cFFQ54chBZtnzq/nRkwH9VMzOD5HhYVKN4Un81WwmbyBQV8CpYd8XriRjnJ2Y1j/
1C6R6fp2/QKBgQDFZQ+V+1rUucsXP4FGihiqUyfGZlTEo6qZdv3syx0oqA4dPcai
NoqK/04KLf2YSqk3dhI9l2iYFDVIwo296Aj6qK+7nySJlkAFjypJ6enozCWVpuey
xl5RIcM4CSdrFtYOSv6y2WRMzV6/WACA2Wys0VdjSy1DK+XQcIDpyp+JWwKBgA9u
pETChmNp5GAYYL2TtBFe9cwhdh3dAwzo9Wg/LJXmKJFxRF0DYGJ9jnmfyQPERlUG
JAlVK9XerQwjj2ixd3sHOO30vsOEWGjqof39rb1gfJOVO9yqXEtQzIDNhn0NGwXz
MZzzT5Rxe1J/+1J+PcSTR5sWPyNPzX/uxwYemukBAoGAW0gObxjT2kYHD677KYVT
WI14NoFjMaGVZPP89eJrRSnXjO0rbLwwDum+0x3Ni3On9QLcXcxW7e6Lws8Hw1Cd
SQYbiJgRXXCcifQ5DkNULXzjvD9ZBj9YrJh9GL3xQF9WnCMfi273K1kE6yOqfAWQ
kCiD5TEwwBQ1lmjbSDDJ+3U=
-----END PRIVATE KEY-----`)

// var publicKey = []byte(`-----BEGIN RSA PUBLIC KEY-----
// MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAohh2ldSLRmZNU/v/i11prEHU2m7cmDgloHa4Hu1rF6F8BexKWqqeMRw96dB+f7+la12WfND5e3eQWEqJk8ufjHMWudDjQXlon5cdu5hqB8rNiaBeNete4VQ5PRZmggrKH/UQCoVVtN4RO4kEsnOS4nnecnl+euvyKN7Qrn8APp99v0o+usNUWd8gZ+o+K+WT5Q+sdIMbdtDmIxI8AmKr0mlRrCPuzvG+99CWWBoDKuonU/8BG4ADbl+6KmZb7fWbkdlIHrWEQyYYNOIgeBoMJ41KEJIJGwwMkbw8OkF4t20gs4LTjixK5TUSAjbe9L2HAA/X9GbAjS1Owo+AfXNDFwIDAQAB
// -----END RSA PUBLIC KEY-----`)

var myPublicKeyBase = `MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqSbw9rY58HuqDZ0tyY2xCwQfRiXUI3EasGAMhLbQiMAMBGSZQjcKUXHAnTe2YKJHD/Tn7vQpQtB6/ktuECH60lSVkEph+e6uXgE4nWbRkYur5pkosEHG/fTE+Upo0E9ID1BgElFQ/3gR9iUgOuEiy2EoyJWNyyCrT6OZB/nKSrl7M7EWFr6X84NYZWp1pzKcf8cQIxwzuw9yF3Kwa9L43ANgJItrEqZ4zxeJk0/qJEC20UaE2AmCTkl/HfY6B6j8mNeHSPSyayVcoDitWHiNeUKE3QL2sXog7G+7pljYFZz6M259w2r/zm8VkD5XJt3BGtjdZtIN0Ckf3C+ygOMY/wIDAQAB`

var myPublicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqSbw9rY58HuqDZ0tyY2x
CwQfRiXUI3EasGAMhLbQiMAMBGSZQjcKUXHAnTe2YKJHD/Tn7vQpQtB6/ktuECH6
0lSVkEph+e6uXgE4nWbRkYur5pkosEHG/fTE+Upo0E9ID1BgElFQ/3gR9iUgOuEi
y2EoyJWNyyCrT6OZB/nKSrl7M7EWFr6X84NYZWp1pzKcf8cQIxwzuw9yF3Kwa9L4
3ANgJItrEqZ4zxeJk0/qJEC20UaE2AmCTkl/HfY6B6j8mNeHSPSyayVcoDitWHiN
eUKE3QL2sXog7G+7pljYFZz6M259w2r/zm8VkD5XJt3BGtjdZtIN0Ckf3C+ygOMY
/wIDAQAB
-----END PUBLIC KEY-----`)

type PayReq struct {
	MerchantID string `json:"merchant_id"` // 商户编号id
	OrderID    string `json:"order_id"`    // 商户订单号，订单唯一
	PayType    string `json:"pay_type"`    // 默认：1
	NotifyUrl  string `json:"notify_url"`  // 商户接受回调通知地址
	ReturnUrl  string `json:"return_url"`  // 订单成功收银台跳转地址
	Sign       string `json:"sign"`        // 验签值
	Amount     int64  `json:"amount"`      // 订单金额，单位（分）
	Currency   string `json:"currency"`    // 货币类型，默认 INR
	Name       string `json:"name"`        // 姓名，强烈建议真实
	Phone      string `json:"phone"`       // 手机号，强烈建议真实
	Email      string `json:"email"`       // 邮箱，强烈建议真实
}

type PayResp struct {
	Code int `json:"code"` // 响应代码
	Data struct {
		Amount     int64  `json:"amount"`      // 订单金额（单位：分）
		Fee        int64  `json:"fee"`         // 手续费（单位：分）
		ID         int64  `json:"id"`          // 订单 ID
		MerchantID string `json:"merchant_id"` // 商户 ID
		OrderID    string `json:"order_id"`    // 商户订单号
		PayUrl     string `json:"pay_url"`     // 支付链接
	} `json:"data"` // 响应数据
	Message string `json:"message"` // 响应消息
}

type PayCallbackReq struct {
	Code    int    `json:"code"`    // 成功：200 失败：400
	Message string `json:"message"` // 成功：SUCCESS 失败：FAILED
	Data    struct {
		ID          string `json:"id"`           // 平台订单编号id
		MerchantID  string `json:"merchant_id"`  // 商户编号id
		Amount      string `json:"amount"`       // 金额（分）
		Fee         string `json:"fee"`          // 手续费（分）
		Sign        string `json:"sign"`         // 验签值，只有data中的参数参与验签
		OrderID     string `json:"order_id"`     // 商户订单编号id
		Currency    string `json:"currency"`     // 货币
		OperatorNum string `json:"operator_num"` // 操作员编号
		UTR         string `json:"utr"`          // 交易流水号
	} `json:"data"` // 回调通知数据
}

type WithdrawReq struct {
	MerchantID   string `json:"merchant_id"`   // 商户编号id
	OrderID      string `json:"order_id"`      // 商户订单号，订单唯一
	Amount       int64  `json:"amount"`        // 订单金额，单位（分）
	WithdrawType string `json:"withdraw_type"` // 提现类型：1=银行转账，2=UPI转账
	NotifyUrl    string `json:"notify_url"`    // 商户接受回调通知地址
	Name         string `json:"name"`          // 姓名
	Phone        string `json:"phone"`         // 手机号
	Email        string `json:"email"`         // 邮箱
	BankCode     string `json:"bank_code"`     // 银行IFSC，UPI代付可以随便填
	Currency     string `json:"currency"`      // 货币类型，默认为 INR
	Sign         string `json:"sign"`          // 验签值
	Account      string `json:"account"`       // 银行账号或者UPI账号
}

type WithdrawResp struct {
	Code    int    `json:"code"`    // 状态码，200 表示提单成功，400 表示提单失败
	Message string `json:"message"` // 状态信息，"SUCCESS" 表示成功，"FAIL" 表示失败
	Data    struct {
		ID      int    `json:"id"`      // 平台订单号
		MerID   string `json:"merId"`   // 商户号
		OrderID string `json:"orderId"` // 商户订单号
		Amount  int    `json:"amount"`  // 金额
		Fee     int    `json:"fee"`     // 手续费
	} `json:"data"` // 数据字段
}

type WithdrawCallbackReq struct {
	Code    int    `json:"code"`    // 成功：200 失败：400
	Message string `json:"message"` // 成功：SUCCESS 失败：FAILED
	Data    struct {
		ID              string `json:"id"`                // 平台订单编号id
		PlatformOrderID string `json:"platform_order_id"` // 平台订单编号id
		MerchantID      string `json:"merchant_id"`       // 商户编号id
		OrderID         string `json:"order_id"`          // 商户订单编号id
		Amount          string `json:"amount"`            // 金额（分）
		Fee             string `json:"fee"`               // 手续费（分）
		Sign            string `json:"sign"`              // 验签值，只有data中的参数参与验签
		Currency        string `json:"currency"`          // 货币
		UTR             string `json:"utr"`               // 交易流水号
	} `json:"data"` // 回调通知数据
}

type QueryWithdrawReq struct {
	MerchantID string `json:"merchant_id"` // 商户编号id，联系运营人员获取
	OrderID    string `json:"order_id"`    // 商户需要查询的订单号
	Sign       string `json:"sign"`        // 验签值
}

type QueryWithdrawResp struct {
	Code    int    `json:"code"`    // 成功：200，失败：其他
	Message string `json:"message"` // 响应消息，成功：SUCCESS，失败：FAILED
	Data    struct {
		ID         string `json:"id"`          // 平台单号
		MerchantID string `json:"merchant_id"` // 商户编号
		OrderID    string `json:"order_id"`    // 商户订单号
		Status     int    `json:"status"`      // 订单状态：2=成功，-2=失败，1=处理中
		UTR        string `json:"utr"`         // 成功订单 utr
	} `json:"data"` // 回调通知数据
}

type QueryPayReq struct {
	MerchantID string `json:"merchant_id"` // 商户编号id，联系运营人员获取
	OrderID    string `json:"order_id"`    // 商户需要查询的订单号
	Sign       string `json:"sign"`        // 验签值
}

type QueryPayResp struct {
	Code    int    `json:"code"`    // 成功：200，失败：其他
	Message string `json:"message"` // 响应消息，成功：SUCCESS，失败：FAILED
	Data    struct {
		ID         string `json:"id"`          // 平台单号
		MerchantID string `json:"merchant_id"` // 商户编号
		OrderID    string `json:"order_id"`    // 商户订单号
		Status     int    `json:"status"`      // 订单状态：1=成功，0=处理中
	} `json:"data"` // 回调通知数据
}
