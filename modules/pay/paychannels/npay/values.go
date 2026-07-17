package npay

const (
	baseURL          = "http://mmisid.nasiwjplaky.top"
	payURL           = "/payin"
	queryPayURL      = "/payinQuery"
	withdrawURL      = "/payout"
	queryWithdrawURL = "/payoutQuery"
	queryBalanceURL  = "/balance"
	mid              = "YM103127"
	payinkey         = "680c4c753aa90304183d915c68107b52"
	payoutkey        = "eb2195c3118ab61d1f7ed759f96f5504"
)

type PayReq struct {
}

type PayResp struct {
}

type PayCallbackReq struct {
}

type QueryPayReq struct {
}

type QueryPayResp struct {
}

type WithdrawReq struct {
	Merchant     string `json:"merchant"`      // merchant
	WithdrawID   string `json:"withdrawId"`    // withdrawId
	DrawAmount   string `json:"drawAmount"`    // drawAmount
	AccountNo    string `json:"drawAccountNo"` // AccountNo
	DrawName     string `json:"drawName"`      // Name
	DrawBankName string `json:"drawBankName"`  // BankName (Fixed: IFSC)
	DrawBankID   string `json:"drawBankId"`    // BankId (Fixed: ifsc)
	IfscNo       string `json:"ifscNo"`        // IFSC Account
	NoticeURL    string `json:"noticeUrl"`     // noticeUrl
	Md5Info      string `json:"md5Info"`       // md5Info
	Sign         string `json:"sign"`          // sign
}

type WithdrawResp struct {
	Merchant   string `json:"merchant"`   // merchant
	WithdrawID string `json:"withdrawId"` // withdrawId
	SerialID   string `json:"serialId"`   // serialId
	Status     string `json:"status"`     // status (2: apply Success)
	Message    string `json:"message"`    // message
	Amount     string `json:"amount"`     // amount
	RecordDate string `json:"recordDate"` // yyyyMMddHHmmss
	Sign       string `json:"sign"`       // sign
}

type WithdrawCallbackReq struct {
	Merchant        string `json:"merchant"`        // merchant
	WithdrawID      string `json:"withdrawId"`      // withdrawId
	SerialID        string `json:"serialId"`        // serialId
	Status          string `json:"status"`          // status (0: failed, 2: successful)
	Amount          string `json:"amount"`          // amount
	RecordDate      string `json:"recordDate"`      // yyyyMMddHHmmss
	UTR             string `json:"utr"`             // UTR
	AccountNo       string `json:"accountNo"`       // accountNo (optional)
	AccountName     string `json:"accountName"`     // accountName (optional)
	AccountBankName string `json:"accountBankName"` // accountBankName (optional)
	Sign            string `json:"sign"`            // sign
}

type QueryWithdrawReq struct {
	Merchant   string `json:"merchant"`   // merchant
	WithdrawID string `json:"withdrawId"` // withdrawId
	Sign       string `json:"sign"`       // sign (to be Uppercase)
}

type QueryWithdrawResp struct {
	Merchant        string `json:"merchant"`        // merchant
	WithdrawID      string `json:"withdrawId"`      // withdrawId
	SerialID        string `json:"serialId"`        // serialId
	Status          string `json:"status"`          // Transaction status (2: successful, 0: failed, 3: processing)
	Amount          string `json:"amount"`          // amount
	RecordDate      string `json:"recordDate"`      // yyyyMMddHHmmss
	AccountNo       string `json:"accountNo"`       // accountNo
	AccountName     string `json:"accountName"`     // accountName
	AccountBankName string `json:"accountBankName"` // accountBankName
	Sign            string `json:"sign"`            // sign (to be Uppercase)
	UTR             string `json:"utr"`             // UTR
}
