package toplemonpay

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
	// baseURL = "http://apiwk.sifang99.com"
	baseURL          = "https://adm.toplemonpay.com"
	payURL           = "/payment/createOrder"
	queryPayURL      = "/payment/orderStatus"
	withdrawURL      = "/payout/createOrder"
	queryWithdrawURL = "/payout/orderStatus"
	queryBalanceURL  = "/payout/balance"
	mid              = "slot666ly"
	key              = "Wlqct8jIta1od0JhDWTy2H3CeOm2wFVYE6qWh9w4GYvn0k3ev2i8hnMSGuFl"
)

type PayReq struct {
	MerchantId string `json:"merchantId"`
	Version    string `json:"version"`
	Channel    string `json:"channel"`
	OrderNo    string `json:"orderNo"`
	Amount     string `json:"amount"`
	Phone      string `json:"phone"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Sign       string `json:"sign"`
	NotifyUrl  string `json:"notifyUrl"`
	ReturnUrl  string `json:"returnUrl"`
	NativeApp  string `json:"nativeApp"`
}

type PayResp struct {
	Code int    `json:"code"` // 通用的返回状态码
	Msg  string `json:"msg"`  // 通用的返回消息
	Data struct {
		RealAmount string `json:"realAmount"` // 实际充值金额
		PayNo      string `json:"payNo"`      // 平台订单号
		Amount     string `json:"amount"`     // 发起的金额
		OrderNo    string `json:"orderNo"`    // 商户的订单号码
		MerchantId string `json:"merchantId"` // 商户号
		Fee        string `json:"fee"`        // 手续费
		Sign       string `json:"sign"`       // 返回的数据签名
		Currency   string `json:"currency"`   // 币种
		Extended   string `json:"extended"`   // 扩展字段（JSON 字符串）
		JumpUrl    string `json:"jumpUrl"`    // 跳转地址
	} `json:"data"` // 可变的data字段
}

type PayCallbackReq struct {
	MerchantId string `json:"merchantId"` // 商户号
	OrderNo    string `json:"orderNo"`    // 商户的订单号码
	PayNo      string `json:"payNo"`      // 平台订单号
	Amount     string `json:"amount"`     // 发起的金额
	AmountReal string `json:"amountReal"` // 实际充值金额，一般情况下与 Amount 一致
	Utr        string `json:"utr"`        // UTR (Unique Transaction Reference)
	Status     string `json:"status"`     // 订单状态: 0 - 发起, 1 - 交易中, 2 - 成功, 3 - 失败, 4 - 失效, 5 - 退款
	FailMsg    string `json:"failMsg"`    // 失败原因
	Sign       string `json:"sign"`       // 加密的sign，使用方式参加加密规则
}

type WithdrawReq struct {
	MerchantId string `json:"merchantId"` // 商户号
	Version    string `json:"version"`    // 版本号，传 V3
	Channel    string `json:"channel"`    // 支出类型 banks
	OrderNo    string `json:"orderNo"`    // 订单号
	Amount     string `json:"amount"`     // 金额 保留两位小数 例如:100.01
	Name       string `json:"name"`       // 用户真实姓名
	Phone      string `json:"phone"`      // 用户真实手机号，不需要国际区号
	Email      string `json:"email"`      // 用户真实邮箱
	Ifsc       string `json:"ifsc"`       // ifsc, 一般是由数字与字母组成
	AccountNum string `json:"accountNum"` // 银行卡号
	NotifyUrl  string `json:"notifyUrl"`  // 回调地址
	Sign       string `json:"sign"`       // 加密的sign
}

type WithdrawResp struct {
	Code int    `json:"code"` // 通用的返回状态码
	Msg  string `json:"msg"`  // 通用的返回消息
	Data struct {
		MerchantId string `json:"merchantId"` // 商户号
		Amount     string `json:"amount"`     // 发起的金额
		PayNo      string `json:"payNo"`      // 平台订单号，记录一下，方便后续查账
		OrderNo    string `json:"orderNo"`    // 商户的订单号码
		Fee        string `json:"fee"`        // 手续费
		Channel    string `json:"channel"`    // 通道
		Sign       string `json:"sign"`       // 返回的数据签名
		Currency   string `json:"currency"`   // 币种
	} `json:"data"` // 可变的data字段
}

type WithdrawCallbackReq struct {
	MerchantId string `json:"merchantId"` // 商户号
	OrderNo    string `json:"orderNo"`    // 商户的订单号码
	PayNo      string `json:"payNo"`      // 平台订单号
	Channel    string `json:"channel"`    // 支出类型 (e.g., "banks", "paytm", "upi")
	Amount     string `json:"amount"`     // 发起的金额
	AmountReal string `json:"amountReal"` // 实际充值金额，一般情况下与 Amount 一致
	Utr        string `json:"utr"`        // UTR (Unique Transaction Reference)
	Status     string `json:"status"`     // 订单状态: "1" - 交易中, "2" - 成功, "3" - 失败, "5" - 退款
	FailMsg    string `json:"failMsg"`    // 失败原因
	Sign       string `json:"sign"`       // 加密的sign，使用方式参加加密规则
}

type QueryWithdrawReq struct {
	MerchantId string `json:"merchantId"` // 商户号
	Timestamp  string `json:"timestamp"`  // 时间戳
	PayNo      string `json:"payNo"`      // 平台订单号，downNo为空时必传
	DownNo     string `json:"downNo"`     // 商户订单号，payNo为空时必传
	Sign       string `json:"sign"`       // 加密的sign
}

type QueryWithdrawResp struct {
	Code int    `json:"code"` // 通用的返回状态码
	Msg  string `json:"msg"`  // 通用的返回消息
	Data struct {
		MerchantId string `json:"merchantId"` // 商户号
		OrderNo    string `json:"orderNo"`    // 商户订单号
		PayNo      string `json:"payNo"`      // 平台订单号
		Amount     string `json:"amount"`     // 金额
		AmountReal string `json:"amountReal"` // 实际金额
		Utr        string `json:"utr"`        // UTR (Unique Transaction Reference)
		Sign       string `json:"sign"`       // 签名
		Status     string `json:"status"`     // 状态：0 - 发起, 1 - 交易中, 2 - 成功, 3 - 失败, 4 - 失效, 5 - 退款
	} `json:"data"` // 可变的data字段
}

type QueryPayReq struct {
	MerchantId string `json:"merchantId"` // 商户号
	Timestamp  string `json:"timestamp"`  // 时间戳
	PayNo      string `json:"payNo"`      // 平台订单号，downNo为空时必传
	DownNo     string `json:"downNo"`     // 商户订单号，payNo为空时必传
	Sign       string `json:"sign"`       // 加密的sign，使用方式参考加密规则
}

type QueryPayResp struct {
	Code int    `json:"code"` // 通用的返回状态码
	Msg  string `json:"msg"`  // 通用的返回消息
	Data struct {
		MerchantId string `json:"merchantId"` // 商户号
		PayNo      string `json:"payNo"`      // 平台订单号
		Amount     string `json:"amount"`     // 发起的金额
		OrderNo    string `json:"orderNo"`    // 商户的订单号码
		AmountReal string `json:"amountReal"` // 实际充值金额，一般情况下与 Amount 一致
		Utr        string `json:"utr"`        // UTR (Unique Transaction Reference)
		Status     string `json:"status"`     // 订单状态: 0 - 发起, 1 - 交易中, 2 - 成功, 3 - 失败, 4 - 失效, 5 - 退款
		FailMsg    string `json:"failMsg"`    // 失败原因
	} `json:"data"` // 可变的data字段
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
		if fieldName == "sign" {
			continue
		}
		if field.Kind() == reflect.String && field.String() != "" {
			paramMap[fieldName] = field.String()
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

	// fmt.Println(data)
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
