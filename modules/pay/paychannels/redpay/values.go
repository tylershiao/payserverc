package redpay

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/liangdas/mqant/log"
)

const (
	baseURL          = "https://api.redpayapp.com"
	payURL           = "/api/payin/order"
	queryPayURL      = "/api/payin/query"
	withdrawURL      = "/api/payout/order"
	queryWithdrawURL = "/api/payout/query"
	queryBalanceURL  = "/api/balance/query"
	// mid              = "MT00001"
	// key              = "76A5DC53940B05009EA6A240823A771B"
	mid = "R0077"
	key = "0A087679BE2F250EEFBEF0D00D763F41"
)

type PayReq struct {
	MerchantCid     string `json:"merchantCid"`     // 商户cid，可以在管理平台找到
	MerchantOrderNo string `json:"merchantOrderNo"` // 商户订单号，需唯一
	Country         string `json:"country"`         // 国家代码，例如：IN
	Currency        string `json:"currency"`        // 币种，例如：INR
	Amount          string `json:"amount"`          // 订单金额，单位为分，不带小数
	Channel         string `json:"channel"`         // 通道类型，印度原生=1001，印度唤醒=1002
	Name            string `json:"name"`            // 用户真实姓名
	Mobile          string `json:"mobile"`          // 用户真实手机号，不需要国际区号
	Email           string `json:"email"`           // 用户真实邮箱
	NotifyUrl       string `json:"notifyUrl"`       // 异步回调地址
	ReturnUrl       string `json:"returnUrl"`       // 同步回调地址，支付完成后跳转，必须为https
	Subject         string `json:"subject"`         // 商品名称或订单标题
	Body            string `json:"body"`            // 商品或订单的详细描述
	Timestamp       string `json:"timestamp"`       // 请求提交时间戳，毫秒
	Sign            string `json:"sign"`            // 加密签名，详见 Signature
}

type PayResp struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
	Data struct {
		PlatformOrderID string `json:"platformOrderId"` // 平台订单号，用于后续查账
		Amount          string `json:"amount"`          // 发起的金额，单位分
		MerchantOrderNo string `json:"merchantOrderNo"` // 商户订单号
		Currency        string `json:"currency"`        // 币种
		Country         string `json:"country"`         // 国家
		PayURL          string `json:"payUrl"`          // 收银台地址
	} `json:"data"`
}

type PayCallbackReq struct {
	MerchantCid     string `json:"merchantCid"`     // 商户cid
	MerchantOrderNo string `json:"merchantOrderNo"` // 商户订单号
	PlatformOrderID string `json:"platformOrderId"` // 平台订单号
	Amount          string `json:"amount"`          // 发起金额，单位分
	PayAmount       string `json:"payAmount"`       // 实际充值金额，一般与 amount 一致
	Status          string `json:"status"`          // 订单状态：0-交易中，1-成功，2-失败
	PayTime         string `json:"payTime"`         // 支付完成时间，单位秒
	Sign            string `json:"sign"`            // 签名，需要URL解码后验证
	UTR             string `json:"utr"`             // 印度银行付款凭证，不参与签名
}

type QueryPayReq struct {
	MerchantCid     string `json:"merchantCid"`     // 商户cid，可以在管理平台找到
	Timestamp       string `json:"timestamp"`       // 时间戳，毫秒
	MerchantOrderNo string `json:"merchantOrderNo"` // 商户订单号
	Sign            string `json:"sign"`            // 加密签名，参考加密规则
}

type QueryPayResp struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
	Data struct {
		Status          string `json:"status"`          // 订单状态：0-交易中，1-成功，2-失败
		PlatformOrderID string `json:"platformOrderId"` // 平台订单号
		Amount          string `json:"amount"`          // 发起金额，单位分
		MerchantOrderNo string `json:"merchantOrderNo"` // 商户订单号
		PayAmount       string `json:"payAmount"`       // 实际充值金额，一般与 Amount 一致，单位分
		PayTime         string `json:"payTime"`         // 支付完成时间，单位秒
		FailMsg         string `json:"failMsg"`         // 失败原因（当 status=2 时返回）
	} `json:"data"`
}

type WithdrawReq struct {
	MerchantCid     string `json:"merchantCid"`     // 商户cid，可以在管理平台找到
	MerchantOrderNo string `json:"merchantOrderNo"` // 商户订单号，需要唯一
	Country         string `json:"country"`         // 国家代码，例如：IN
	Currency        string `json:"currency"`        // 币种，例如：INR
	Amount          int    `json:"amount"`          // 订单金额，单位分，不带小数
	Channel         string `json:"channel"`         // 支付类型，印度原生=1001，印度唤醒=1002
	Name            string `json:"name"`            // 收款人姓名
	Mobile          string `json:"mobile"`          // 收款人手机号，不需要国际区号
	Email           string `json:"email"`           // 收款人邮箱
	BankCode        string `json:"bankCode"`        // 收款人账户的IFSC code
	BankName        string `json:"bankName"`        // 收款人银行名称，印度固定为IMPS
	BankCardNum     string `json:"bankCardNum"`     // 收款人账户或银行卡号
	NotifyUrl       string `json:"notifyUrl"`       // 异步回调地址
	Timestamp       string `json:"timestamp"`       // 请求提交时间戳，毫秒
	Sign            string `json:"sign"`            // 加密签名，详见 Signature
}

type WithdrawResp struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
	Data struct {
		PlatformOrderID string      `json:"platformOrderId"` // 平台订单号，用于记录和后续查账
		Amount          interface{} `json:"amount"`          // 发起的金额，单位分
		MerchantOrderNo string      `json:"merchantOrderNo"` // 商户订单号
		Currency        string      `json:"currency"`        // 币种
		Country         string      `json:"country"`         // 国家
	} `json:"data"`
}

type WithdrawCallbackReq struct {
	MerchantCid     string `json:"merchantCid"`     // 商户cid
	MerchantOrderNo string `json:"merchantOrderNo"` // 商户订单号
	PlatformOrderID string `json:"platformOrderId"` // 平台订单号
	Amount          string `json:"amount"`          // 发起金额，单位分
	PayAmount       string `json:"payAmount"`       // 实际代付金额，与 Amount 一致
	Status          string `json:"status"`          // 订单状态：0=交易中，1=成功，2=失败，5=冲正，6=驳回
	PayTime         string `json:"payTime"`         // 支付完成时间，单位秒
	Sign            string `json:"sign"`            // 加密签名，需要URL解码后验证
	UTR             string `json:"utr"`             // 印度银行付款凭证，不参与签名
}

type QueryWithdrawReq struct {
	MerchantCid     string `json:"merchantCid"`     // 商户cid，可以在管理平台找到
	Timestamp       string `json:"timestamp"`       // 时间戳，毫秒
	MerchantOrderNo string `json:"merchantOrderNo"` // 商户订单号
	Sign            string `json:"sign"`            // 加密签名
}

type QueryWithdrawResp struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
	Data struct {
		Status          string `json:"status"`          // 订单状态：0=交易中，1=成功，2=失败，5=冲正
		PlatformOrderID string `json:"platformOrderId"` // 平台订单号
		Amount          string `json:"amount"`          // 发起金额，单位分
		MerchantOrderNo string `json:"merchantOrderNo"` // 商户订单号
		PayAmount       string `json:"payAmount"`       // 实际充值金额，一般与 Amount 一致
		PayTime         string `json:"payTime"`         // 支付完成时间，单位秒
		FailMsg         string `json:"failMsg"`         // 失败原因（当 status=2 时返回）
	} `json:"data"`
}

type QueryBalanceReq struct {
	MerchantId string `json:"merchantId"` // 商户号
	Timestamp  string `json:"timestamp"`  // 时间戳
	Sign       string `json:"sign"`       // 加密的sign
}

type QueryBalanceResp struct {
	Balance    string `json:"balance"`    // 商户余额
	WaitSettle string `json:"waitSettle"` // 待结算
	Frozen     string `json:"frozen"`     // 冻结资金
}

func GenSign(params interface{}) (string, error) {
	// Use reflection to extract the struct fields
	val := reflect.ValueOf(params)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return "", fmt.Errorf("params should be a struct")
	}

	// Create a map to hold the field names and values
	paramMap := make(map[string]string)

	// Iterate through the struct fields
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := val.Type().Field(i).Tag.Get("json")
		if fieldName == "sign" || fieldName == "utr" {
			continue
		}
		switch field.Kind() {
		case reflect.String:
			if field.String() != "" {
				paramMap[fieldName] = field.String()
			}
		case reflect.Int:
			if field.Int() != 0 {
				paramMap[fieldName] = fmt.Sprintf("%d", field.Int())
			}
		}
	}

	// Sort the keys in the map
	var keys []string
	for k := range paramMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// Prepare the data string by joining the sorted key-value pairs
	var dataParts []string
	for _, k := range keys {
		dataParts = append(dataParts, k+"="+paramMap[k])
	}
	data := strings.Join(dataParts, "&")

	log.Debug("sign data:%v", data)
	// Create the HMAC-SHA1 signer using the key
	signingKey := []byte(key)
	mac := hmac.New(sha1.New, signingKey)

	// Write the data to the MAC
	_, err := mac.Write([]byte(data))
	if err != nil {
		log.Error("err:%v", err)
		return "", err
	}

	// Get the resulting HMAC and encode it in Base64
	rawHmac := mac.Sum(nil)
	encoded := base64.StdEncoding.EncodeToString(rawHmac)

	return encoded, nil
}
