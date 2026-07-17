package wandapay

import (
	"server/config"
	"sync"
)

var (
	// 正式服配置
	prodMerchantNo  = "395"
	prodMerchantKey = "27851755da22416a9cf817e5b3a80285"

	// 测试服配置
	testMerchantNo  = "201"
	testMerchantKey = "1f23c2295a08e214304bf09463d0fbcb"
)

const (
	baseURL          = "https://wandapay.cc"
	payURL           = "/api/payin/request"
	queryPayURL      = "/api/payin/status"
	withdrawURL      = "/api/payout/request"
	queryWithdrawURL = "/api/payout/status"
)

var (
	wandaOnce     sync.Once
	wandaInstance *WandaClient
)

type WandaClient struct {
	MerchantNo  string
	MerchantKey string
}

func getWandaConfig() (merchantNo, merchantKey string) {
	if config.GetBase().Release {
		return prodMerchantNo, prodMerchantKey
	}
	return testMerchantNo, testMerchantKey
}

// 请求结构体
type PayReq struct {
	MerchantNo string `json:"merchant_no"`
	TxnID      string `json:"txn_id"`
	Amount     string `json:"amount"`
	Type       int    `json:"type"`
	Callback   string `json:"callback"`
	Currency   string `json:"currency"`
	Sign       string `json:"sign"`
}

type WithdrawReq struct {
	MerchantNo string `json:"merchant_no"`
	TxnID      string `json:"txn_id"`
	Amount     string `json:"amount"`
	Type       string `json:"type"`
	IFSC       string `json:"ifsc"`
	CardNum    string `json:"card_num"`
	Name       string `json:"name,omitempty"`
	Email      string `json:"email,omitempty"`
	Mobile     string `json:"mobile,omitempty"`
	Callback   string `json:"callback"`
	Currency   string `json:"currency"`
	Sign       string `json:"sign"`
}

type QueryPayReq struct {
	MerchantNo string `json:"merchant_no"`
	TxnID      string `json:"txn_id"`
	Sign       string `json:"sign"`
}

type QueryWithdrawReq struct {
	MerchantNo string `json:"merchant_no"`
	TxnID      string `json:"txn_id"`
	Sign       string `json:"sign"`
}

// 响应结构体
type PayResp struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	MerchantNo string `json:"merchant_no"`
	TxnID      string `json:"txn_id"`
	ID         int64  `json:"id"`
	Amount     string `json:"amount"`
	PayLink    string `json:"pay_link"`
}

func (p *PayResp) IsSuccess() bool {
	return p.Status == "SUCCESS" || p.Message == "TXN CREATED"
}

type WithdrawResp struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	MerchantNo string `json:"merchant_no"`
	TxnID      string `json:"txn_id"`
	ID         int64  `json:"id"`
	Amount     int64  `json:"amount"`
}

func (p *WithdrawResp) IsSuccess() bool {
	return p.Status == "PENDING" || p.Status == "SUCCESS" || p.Message == "TXN CREATED"
}

type QueryPayResp struct {
	Status string `json:"status"`
	TxnID  string `json:"txn_id"`
	Amount string `json:"amount"`
	Code   int    `json:"code,omitempty"`
}

type QueryWithdrawResp struct {
	Status     string `json:"status"`
	TxnID      string `json:"txn_id"`
	FailReason string `json:"fail_reason,omitempty"`
	Amount     string `json:"amount"`
	Message    string `json:"message,omitempty"`
	Code       int    `json:"code,omitempty"`
}

// 回调结构体
type PayCallbackReq struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	ID         string `json:"id"`
	TxnID      string `json:"txn_id"`
	MerchantNo string `json:"merchant_no"`
	Sign       string `json:"sign"`
	UTR        string `json:"utr"`
	Amount     string `json:"amount"`
}

type WithdrawCallbackReq struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	MerchantNo string `json:"merchant_no"`
	ID         string `json:"id"`
	TxnID      string `json:"txn_id"`
	Amount     string `json:"amount"`
	UTR        string `json:"utr"`
	Sign       string `json:"sign"`
}
