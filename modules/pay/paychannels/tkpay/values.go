package tkpay

const (
	baseURL          = "https://seabird.world"
	payURL           = "/api/order/pay/create"
	queryPayURL      = "/api/order/pay/query"
	withdrawURL      = "/api/order/withdraw/create/v2"
	queryWithdrawURL = "/api/order/withdraw/query"
	mid              = "202688046"
	key              = "520700c9f1cc4b9ea4f3fcab2b77b498"
	// mid              = "202366100"
	// key              = "c385fe7029344aef826d8112625b2625x"
)

var privateKey = []byte(`-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDcg+uGgNLbFJbW
Hqfj7XfXWB7LI4AeE7StzlDLi4SFm6XhkK39xZ1EwtlF0URTWnW4uI1q2jKqOIPj
//0MY8fwJrmfQXqAkuCQm8nM2F2gJbKbTFNSph+fk12r49FjrKnYMUID3Qjep5qT
qSJdA7KfnF/SnWNDLX32OCm9ly0wkmKJNajKPsBRxCOts7xudt9BP03Z4Edlh3ux
JIZQao6L/A74GZspH1EweH15uz6StEhVMieizzbSD6Lprpbls3lb02jADVfloDbm
AFdXdbk+Nhgm3lv72Mn5Fw14AY56vnQi4DyvM7TAZHCaa0FTjR2zGT6gGc1YQS5t
1JQSMR6NAgMBAAECggEALfLvKN+/EBeCNqu8oCnjRuP9xqZnE91XUMAk5cT0oEKG
ljHnj3GWC7XuQDYS3quCzDLohNNiq0WLdf7pY6OHDeuD0kN6mz//uyI7MowDNlOk
m+eTlS04ea/jXOXjw7k+yuW4s8VjC8Iw24PqvScAGXIkJWbWQDFChgO97tuJdw0m
owiK9dfwj6C+KdEYmfHED+WSqYPLWVdB/4rBRaq9vyMonul4bIFbpUSCF7WdAVs7
yOmXNHK9muJyrjx/2DxN+Nhh5/Mez6o9t985GGzEFlPPX0/SAtDtjlWi+VtYNjXC
TbNIHU3FdDtAZNr7hfJ2jUBgCExeWecSqymJowreAQKBgQDxnC2TvySMKe5+r6hk
HAqzlJ6R/XmdhxuawnwmFpss41KOYNxXfxVkJDrc4xhnAiXE+M8Yfh+/Hpsq6So/
aP4SlPVSR6/Sr1GZ1VPkZo8lL0cKhMZTXjV1MjMO159daThl+yw5b6CEoYIsSKpt
0zabeuzlEaTn+dbdfk/Qokh+sQKBgQDpphxdqNEyqfwUEpqr8tjJaY9wcVnTKhuU
62tSkhFN+jqtsXHqnpuLcVFveAwu3e/l0F0O4U/l0vKjx9fpUeOgRN+vlYPvo52y
eKkhqdrLyhljW59uskdkY0ANV63lwQcFPfGxDayhnGVnj3pVNWPr9EGyU5XqM0Wj
OkZvhIUsnQKBgQC3DyT9N6rCRBb+h6rnRtPY4ncOl+WoNDD3HfoghapG5WmslLj9
1jzn5gdvSuU3Tso9JmMK5Ky7Cr54Ubd/2dX/UDtbfxM5ZzpFevuYJxi4DM0frB/0
FIIebM2ApmCJzuo7YoZMckZFl3VghiITGh/ax9VrbeXrGUHyCf0BT6Pr0QKBgQCc
U/j8rXwzAsrXf08fpx9nO4queBjakXzc2dr35yVO4BKsrmIC0dZOqSYMh+zW2sWt
JQlucFAt/EM+U/LheFntnMdxSn95MibMB7FaBbSPz4vHD8GPphbgUwvDIZwHcmxJ
FJYGmZiclPDAJp56LnvJnqbqMsTGa0G3ocEI0kjsfQKBgF261WB+uuSkIBT9smTP
xz9XBaZ3tEVVgW2OTzpPL9MrONPpUb1xulF2QQ9eB48Z/zntxviO6E785NbDcJ5m
IhE6u8zFufE3mOJj+cb7bos9eXt3F5Iv1gw7TP+Z75OHG1v2FFJScjZqjzW/mlcH
qmeDKtB5TvnUT0LE9Eyp0gkx
-----END PRIVATE KEY-----`)

var privateKeyCS8 = []byte(`-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDcg+uGgNLbFJbW
Hqfj7XfXWB7LI4AeE7StzlDLi4SFm6XhkK39xZ1EwtlF0URTWnW4uI1q2jKqOIPj
//0MY8fwJrmfQXqAkuCQm8nM2F2gJbKbTFNSph+fk12r49FjrKnYMUID3Qjep5qT
qSJdA7KfnF/SnWNDLX32OCm9ly0wkmKJNajKPsBRxCOts7xudt9BP03Z4Edlh3ux
JIZQao6L/A74GZspH1EweH15uz6StEhVMieizzbSD6Lprpbls3lb02jADVfloDbm
AFdXdbk+Nhgm3lv72Mn5Fw14AY56vnQi4DyvM7TAZHCaa0FTjR2zGT6gGc1YQS5t
1JQSMR6NAgMBAAECggEALfLvKN+/EBeCNqu8oCnjRuP9xqZnE91XUMAk5cT0oEKG
ljHnj3GWC7XuQDYS3quCzDLohNNiq0WLdf7pY6OHDeuD0kN6mz//uyI7MowDNlOk
m+eTlS04ea/jXOXjw7k+yuW4s8VjC8Iw24PqvScAGXIkJWbWQDFChgO97tuJdw0m
owiK9dfwj6C+KdEYmfHED+WSqYPLWVdB/4rBRaq9vyMonul4bIFbpUSCF7WdAVs7
yOmXNHK9muJyrjx/2DxN+Nhh5/Mez6o9t985GGzEFlPPX0/SAtDtjlWi+VtYNjXC
TbNIHU3FdDtAZNr7hfJ2jUBgCExeWecSqymJowreAQKBgQDxnC2TvySMKe5+r6hk
HAqzlJ6R/XmdhxuawnwmFpss41KOYNxXfxVkJDrc4xhnAiXE+M8Yfh+/Hpsq6So/
aP4SlPVSR6/Sr1GZ1VPkZo8lL0cKhMZTXjV1MjMO159daThl+yw5b6CEoYIsSKpt
0zabeuzlEaTn+dbdfk/Qokh+sQKBgQDpphxdqNEyqfwUEpqr8tjJaY9wcVnTKhuU
62tSkhFN+jqtsXHqnpuLcVFveAwu3e/l0F0O4U/l0vKjx9fpUeOgRN+vlYPvo52y
eKkhqdrLyhljW59uskdkY0ANV63lwQcFPfGxDayhnGVnj3pVNWPr9EGyU5XqM0Wj
OkZvhIUsnQKBgQC3DyT9N6rCRBb+h6rnRtPY4ncOl+WoNDD3HfoghapG5WmslLj9
1jzn5gdvSuU3Tso9JmMK5Ky7Cr54Ubd/2dX/UDtbfxM5ZzpFevuYJxi4DM0frB/0
FIIebM2ApmCJzuo7YoZMckZFl3VghiITGh/ax9VrbeXrGUHyCf0BT6Pr0QKBgQCc
U/j8rXwzAsrXf08fpx9nO4queBjakXzc2dr35yVO4BKsrmIC0dZOqSYMh+zW2sWt
JQlucFAt/EM+U/LheFntnMdxSn95MibMB7FaBbSPz4vHD8GPphbgUwvDIZwHcmxJ
FJYGmZiclPDAJp56LnvJnqbqMsTGa0G3ocEI0kjsfQKBgF261WB+uuSkIBT9smTP
xz9XBaZ3tEVVgW2OTzpPL9MrONPpUb1xulF2QQ9eB48Z/zntxviO6E785NbDcJ5m
IhE6u8zFufE3mOJj+cb7bos9eXt3F5Iv1gw7TP+Z75OHG1v2FFJScjZqjzW/mlcH
qmeDKtB5TvnUT0LE9Eyp0gkx
-----END PRIVATE KEY-----`)

// var publicKey = []byte(`-----BEGIN RSA PUBLIC KEY-----
// MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAohh2ldSLRmZNU/v/i11prEHU2m7cmDgloHa4Hu1rF6F8BexKWqqeMRw96dB+f7+la12WfND5e3eQWEqJk8ufjHMWudDjQXlon5cdu5hqB8rNiaBeNete4VQ5PRZmggrKH/UQCoVVtN4RO4kEsnOS4nnecnl+euvyKN7Qrn8APp99v0o+usNUWd8gZ+o+K+WT5Q+sdIMbdtDmIxI8AmKr0mlRrCPuzvG+99CWWBoDKuonU/8BG4ADbl+6KmZb7fWbkdlIHrWEQyYYNOIgeBoMJ41KEJIJGwwMkbw8OkF4t20gs4LTjixK5TUSAjbe9L2HAA/X9GbAjS1Owo+AfXNDFwIDAQAB
// -----END RSA PUBLIC KEY-----`)

var myPublicKeyBase = `MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqSbw9rY58HuqDZ0tyY2xCwQfRiXUI3EasGAMhLbQiMAMBGSZQjcKUXHAnTe2YKJHD/Tn7vQpQtB6/ktuECH60lSVkEph+e6uXgE4nWbRkYur5pkosEHG/fTE+Upo0E9ID1BgElFQ/3gR9iUgOuEiy2EoyJWNyyCrT6OZB/nKSrl7M7EWFr6X84NYZWp1pzKcf8cQIxwzuw9yF3Kwa9L43ANgJItrEqZ4zxeJk0/qJEC20UaE2AmCTkl/HfY6B6j8mNeHSPSyayVcoDitWHiNeUKE3QL2sXog7G+7pljYFZz6M259w2r/zm8VkD5XJt3BGtjdZtIN0Ckf3C+ygOMY/wIDAQAB`

var myPublicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3IPrhoDS2xSW1h6n4+13
11geyyOAHhO0rc5Qy4uEhZul4ZCt/cWdRMLZRdFEU1p1uLiNatoyqjiD4//9DGPH
8Ca5n0F6gJLgkJvJzNhdoCWym0xTUqYfn5Ndq+PRY6yp2DFCA90I3qeak6kiXQOy
n5xf0p1jQy199jgpvZctMJJiiTWoyj7AUcQjrbO8bnbfQT9N2eBHZYd7sSSGUGqO
i/wO+BmbKR9RMHh9ebs+krRIVTInos820g+i6a6W5bN5W9NowA1X5aA25gBXV3W5
PjYYJt5b+9jJ+RcNeAGOer50IuA8rzO0wGRwmmtBU40dsxk+oBnNWEEubdSUEjEe
jQIDAQAB
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
		ID int `json:"id"` // 平台订单号
		// MerID   string `json:"merchant_id"`   // 商户号
		// OrderID string `json:"order_id"` // 商户订单号
		// Amount  int    `json:"amount"`  // 金额
		// Fee     int    `json:"fee"`     // 手续费
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
