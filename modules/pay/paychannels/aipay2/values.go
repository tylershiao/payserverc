package aipay2

const (
	baseURL          = "https://top.adkjk.in/aipay-api"
	payURL           = "/order/submit"
	queryPayURL      = "/order/batch/query"
	withdrawURL      = "/payout/submit"
	queryWithdrawURL = "/payout/query"
	mid              = 5066
	key              = "ItWvowD3NFUbDiUk"
)

type PayReq struct {
	MerchantID      int    `json:"merchantId"`      // 商户id (必传)
	MerchantOrderID string `json:"merchantOrderId"` // 商户订单id (必传)
	Amount          string `json:"amount"`          // 金额 (必传，不同国家有不同格式，建议用 string 保留小数/整数格式)
	Timestamp       int64  `json:"timestamp"`       // 请求时间戳（毫秒，必传）
	Sign            string `json:"sign"`            // 签名 (必传，MD5加密后转小写)
	PayType         int    `json:"payType"`         // 支付模式 (必传，不同国家固定值)
	NotifyURL       string `json:"notifyUrl"`       // 异步通知地址 (可选)
	FirstName       string `json:"firstName"`       // 用户名字 (可选，尽量传)
	LastName        string `json:"lastName"`        // 用户姓氏 (可选，尽量传)
	Mobile          string `json:"mobile"`          // 用户电话 (可选，埃及商户必传 01开头11位手机号)
	Remark          string `json:"remark"`          // 订单说明 (必传)
	Email           string `json:"email"`           // 用户 email (可选)
}

type PayResp struct {
	Code  int    `json:"code"`  // 响应码 0 表示成功
	Error string `json:"error"` // 错误消息
	Data  struct {
		H5Url     string `json:"h5Url"`     // 用于支付的H5地址
		UpiIntent string `json:"upiIntent"` // 深度支付链接(若无则使用 h5Url)
		OrderID   string `json:"orderId"`   // 平台订单号
		Remark    string `json:"remark"`    // 备注
	} `json:"data"` // 响应结构体
}

type QueryPayReq struct {
	MerchantID       int32    `json:"merchantId"`       // 商户id (必传)
	MerchantOrderIDs []string `json:"merchantOrderIds"` // 商户订单id列表 (必传)
	Sign             string   `json:"sign"`             // 签名 (必传，MD5小写)
}

type QueryPayResp struct {
	Code  int32  `json:"code"`  // 响应码 0 表示成功
	Error string `json:"error"` // 错误消息
	Data  []struct {
		Status          int    `json:"status"`          // 订单状态 0 待支付 1 已支付 2 已失效
		MerchantOrderID string `json:"merchantOrderId"` // 商户订单id
		OrderID         string `json:"orderId"`         // 平台订单id
		PayAmount       string `json:"payAmount"`       // 实付金额 (失败或待支付时无返回，建议用 string 保留格式)
	} `json:"data"` // 订单数据列表
}

type PayCallbackReq struct {
	MerchantID      int    `json:"merchantId"`      // 商户id
	MerchantOrderID string `json:"merchantOrderId"` // 商户订单id
	Amount          string `json:"amount"`          // 订单金额（带小数）
	PayAmount       string `json:"payAmount"`       // 实付金额（带小数，代付无该字段）
	Timestamp       int64  `json:"timestamp"`       // 时间戳
	OrderID         int64  `json:"orderId"`         // 平台订单id
	Status          int    `json:"status"`          // 订单状态 0 等待 1 成功 2 失败
	Sign            string `json:"sign"`            // 签名
	Msg             string `json:"msg"`             // 错误消息
	PayTime         int64  `json:"payTime"`         // 支付时间(毫秒级，仅成功时有该字段)
}

type WithdrawReq struct {
	MerchantID      int    `json:"merchantId"`      // 商户id (必传)
	MerchantOrderID string `json:"merchantOrderId"` // 商户订单id (必传)
	Amount          string `json:"amount"`          // 金额 (不同国家要求不同格式，必传，保留字符串避免精度问题)
	Timestamp       int64  `json:"timestamp"`       // 请求时间戳 (必传)
	Sign            string `json:"sign"`            // 签名 (必传)
	NotifyURL       string `json:"notifyUrl"`       // 异步通知地址 (可选)
	FundAccount     struct {
		AccountType string `json:"accountType"` // 账户类型 (必传)
		Contact     struct {
			Name    string `json:"name"`    // 收款人姓名 (必传)
			Email   string `json:"email"`   // 邮箱 (可选)
			Mobile  string `json:"mobile"`  // 手机 (可选)
			Address string `json:"address"` // 地址 (可选)
		} `json:"contact"` // 收款人信息 (必传)
		BankAccount struct {
			Name          string `json:"name"`          // 银行名称 (无则传 Bank)
			IFSC          string `json:"ifsc"`          // 银行编码
			AccountNumber string `json:"accountNumber"` // 银行卡号
		} `json:"bankAccount"` // 银行卡信息 (accountType=bank_account 时必传)
		VPA struct {
			Address string `json:"address"` // UPI 地址
		} `json:"vpa"`
	} `json:"fundAccount"` // 收款信息 (必传)
}

type WithdrawResp struct {
	Code  int    `json:"code"`  // 响应码 0 表示成功
	Error string `json:"error"` // 错误消息
	Data  struct {
		PayoutID          string  `json:"payoutId"`          // 平台提现id
		MerchantPayoutID  string  `json:"merchantPayoutId"`  // 商户订单号
		PayoutAccountType string  `json:"payoutAccountType"` // 用户提现类型
		Status            int     `json:"status"`            // 提现状态 0 等待 1 成功 2 失败
		Msg               string  `json:"msg"`               // 错误消息
		Amount            float64 `json:"amount"`            // 订单金额（₹，保留 string 避免精度问题）
		Fee               string  `json:"fee"`               // 手续费（₹）
		Tax               string  `json:"tax"`               // 税费（₹）
		CreateTime        string  `json:"createTime"`        // 订单时间
		PayoutTime        string  `json:"payoutTime"`        // 提现到账时间
		UTR               string  `json:"utr"`               // UTR
	} `json:"data"` // 响应数据
}

type WithdrawCallbackReq struct {
	MerchantID      int    `json:"merchantId"`      // 商户id
	MerchantOrderID string `json:"merchantOrderId"` // 商户订单id
	Amount          string `json:"amount"`          // 订单金额（带小数）
	PayAmount       string `json:"payAmount"`       // 实付金额（带小数，代付无该字段）
	Timestamp       int64  `json:"timestamp"`       // 时间戳
	OrderID         int64  `json:"orderId"`         // 平台订单id
	Status          int    `json:"status"`          // 订单状态 0 等待 1 成功 2 失败
	Sign            string `json:"sign"`            // 签名（MD5小写）
	Msg             string `json:"msg"`             // 错误消息
	PayTime         int64  `json:"payTime"`         // 支付时间(毫秒级，仅成功时有)
}

type QueryWithdrawReq struct {
	MerchantID      int    `json:"merchantId"`      // 商户id (必传)
	MerchantOrderID string `json:"merchantOrderId"` // 商户订单id (必传)
	Timestamp       int64  `json:"timestamp"`       // 请求时间戳 (必传)
	Sign            string `json:"sign"`            // 签名 (必传，MD5加密后小写)
}

type QueryWithdrawResp struct {
	Code  int    `json:"code"`  // 响应码 0 表示成功
	Error string `json:"error"` // 错误消息
	Data  struct {
		PayoutID          string `json:"payoutId"`          // 平台提现id
		MerchantPayoutID  string `json:"merchantPayoutId"`  // 商户订单号
		PayoutAccountType string `json:"payoutAccountType"` // 用户提现类型
		Status            int    `json:"status"`            // 提现状态 0 等待 1 成功 2 失败
		Msg               string `json:"msg"`               // 错误消息
		Amount            string `json:"amount"`            // 订单金额（₹，建议用 string 避免精度问题）
		Fee               string `json:"fee"`               // 手续费（₹，建议用 string）
		Tax               string `json:"tax"`               // 税费（₹，建议用 string）
		CreateTime        string `json:"createTime"`        // 订单时间
		PayoutTime        string `json:"payoutTime"`        // 提现到账时间
		UTR               string `json:"utr"`               // UTR
	} `json:"data"` // 响应数据
}
