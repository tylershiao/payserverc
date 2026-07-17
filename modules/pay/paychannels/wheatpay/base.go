package wheatpay

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"server/common"
	"server/modules/pay/base"
	"server/sign"
	"sort"
	"strings"

	"github.com/liangdas/mqant/log"
)

// NewSub 初始化子支付模块
// 根据不同的操作类型（支付、提现、回调等）设置相应的请求URL、响应结构体和签名参数
func NewSub(b *base.Base) {
	sub := &Sub{
		Base: b,
	}
	// b.SignKey = key
	b.HttpType = base.HttpTypeJson
	// b.ShouldSignUpper = true
	b.CallbackResp.Msg = "success"
	b.SignPassStr = []string{"hash"}
	b.SignKey = key
	b.KeyName = "key"
	// b.PassKeyName = true
	// b.ShouldSignUpper = true
	switch b.Opt {
	case base.OPTPay:
		b.Resp = new(PayResp)
		b.ReqURL = baseURL + payURL
	case base.OPTWithdraw:
		b.Resp = new(WithdrawResp)
		b.ReqURL = baseURL + withdrawURL
	case base.OPTPayCB:
		b.CallbackReq = new(PayCallbackReq)
	case base.OPTWithdrawCB:
		b.CallbackReq = new(WithdrawCallbackReq)
	case base.OPTQueryWithdraw:
		b.Resp = new(QueryWithdrawResp)
		b.ReqURL = baseURL + queryWithdrawURL
	case base.OPTQueryPay:
		b.Resp = new(QueryPayResp)
		b.ReqURL = baseURL + queryPayURL
	}
	b.Sub = sub
}

// Sub 支付模块的子结构体
type Sub struct {
	Base *base.Base
}

// PackHeader 用于打包请求头信息
// 当前实现为空，可根据需要扩展
func (s *Sub) PackHeader(header http.Header) {
	header.Add("app-id", appID)
	header.Add("app-key", key)
}

// PackReq 根据操作类型打包请求数据
// 根据不同的操作类型（支付、提现、查询等）返回对应的请求结构体
func (s *Sub) PackReq() interface{} {
	switch s.Base.Opt {
	case base.OPTPay:
		return s.PackPayReq()
	case base.OPTWithdraw:
		return s.PackWithdrawReq()
	case base.OPTQueryWithdraw:
		return s.PackQueryWithdrawReq()
	case base.OPTQueryPay:
		return s.PackQueryPayReq()
	}
	return nil
}

// GetResp 处理响应数据
// 根据不同的操作类型处理响应结果，设置相应的返回值
func (s *Sub) GetResp() error {
	switch s.Base.Opt {
	case base.OPTPay:
		resp := s.Base.Resp.(*PayResp)
		if resp.Status != 200 || resp.Data.PaymentPgURL == "" {
			return errors.New("pay fail")
		}
		s.Base.PayResp = &common.PayResp{APIOrderID: resp.Data.TxnReference, URL: resp.Data.PaymentPgURL, Channel: int(s.Base.Channel), Message: resp.Message}
	case base.OPTQueryPay:
		resp := s.Base.Resp.(*QueryPayResp)
		log.Debug("QueryPayResp:%+v", resp)
		s.Base.QueryPayResp.OrderID = resp.Data.TransferID
		s.Base.QueryPayResp.APIOrderID = resp.Data.TxnReference
		if s.Base.Status == 0 {
			switch resp.Data.Status {
			case "SUCCESS":
				s.Base.QueryPayResp.Status = base.QueryStatusOrderSuccess
			case "FAILED":
				s.Base.QueryPayResp.Status = base.QueryStatusFail
			}
		}
		if s.Base.QueryPayResp.Status == 0 {
			s.Base.QueryPayResp.Status = base.QueryStatusUnknown
		}
	case base.OPTWithdraw:
		resp := s.Base.Resp.(*WithdrawResp)
		if s.Base.Status == 0 && resp.Status != 200 {
			return errors.New("withdraw fail")
		}
		s.Base.WithdrawResp = &common.WithdrawResp{APIOrderID: resp.Data.TxnReference, Channel: int(s.Base.Channel), Message: resp.Message}
	case base.OPTQueryWithdraw:
		resp := s.Base.Resp.(*QueryWithdrawResp)
		log.Debug("QueryWithdrawResp:%+v", resp)
		s.Base.QueryWithdrawResp.OrderID = s.Base.QueryWithdrawReq.OrderID
		s.Base.QueryWithdrawResp.APIOrderID = resp.Data.TxnReference
		if s.Base.Status == 0 {
			// if resp.Code == 1 {
			switch resp.Data.Status {
			case "SUCCESS":
				s.Base.QueryWithdrawResp.Status = base.QueryStatusOrderSuccess
			case "PENDING":
				s.Base.QueryWithdrawResp.Status = base.QueryStatusSuccess
			case "FAILED":
				s.Base.QueryWithdrawResp.Status = base.QueryStatusFail
			default:
				s.Base.QueryWithdrawResp.Status = base.QueryStatusUnknown
			}
			// } else {
			// if resp.Data.Status == "3" {
			// s.Base.QueryWithdrawResp.Status = base.QueryStatusFail
			// }
			// }
		}
		if s.Base.QueryWithdrawResp.Status == 0 {
			s.Base.QueryWithdrawResp.Status = base.QueryStatusUnknown
		}
	default:
		return errors.New("unknown opt")
	}
	return nil
}

// PackPayReq 打包支付请求数据
// 构造支付请求结构体，包括商户ID、订单号、支付类型、通知地址等信息，并生成签名
func (s *Sub) PackPayReq() interface{} {
	r := s.Base.PayReq
	send := &PayReq{
		AppID:       appID,
		TransferID:  r.OrderID,
		Currency:    "INR",
		Amount:      r.Amount,
		PayerName:   r.Name,
		PayerEmail:  r.Email,
		PayerPhone:  r.Mobile,
		CallbackURL: s.Base.GetPayCallbackURL(),
	}
	signStr := sign.GetSignStrl(send)
	send.Hash = base.SignWithHmacSHA256([]byte(signStr), secret)
	return send
}

// PackWithdrawReq 打包提现请求数据
// 构造提现请求结构体，根据支付类型（银行、UPI）设置不同的字段，并生成签名
func (s *Sub) PackWithdrawReq() interface{} {
	r := s.Base.WithdrawReq
	send := &WithdrawReq{
		AppID:      appID,
		TransferID: r.OrderID,
		Currency:   "INR",
		Amount:     fmt.Sprintf("%.2f", r.Amount),
		PayerName:  r.Name,
		PayerEmail: r.Email,
		PayerPhone: r.Mobile,
		// BankAccount: ,
		// IFSC: ,
		CallbackURL: s.Base.GetWithdrawCallbackURL(),
	}
	switch r.PayType {
	case common.WithdrawTypeBank:
		send.BankAccount = r.Account
		send.IFSC = r.IFSC
	// case common.WithdrawTypeUPI:
	// 	send.WayCode = "2"
	// 	send.AccountNo = r.Account
	default:
		return nil
	}
	signStr := sign.GetSignStrl(send)
	send.Hash = base.SignWithHmacSHA256([]byte(signStr), secret)
	return send
}

// PackQueryWithdrawReq 打包查询提现请求数据
// 构造查询提现请求结构体，包括商户ID和订单号，并生成签名
func (s *Sub) PackQueryWithdrawReq() interface{} {
	r := s.Base.QueryWithdrawReq
	send := &QueryWithdrawReq{
		AppID:      appID,
		TransferID: r.OrderID,
	}
	signStr := sign.GetSignStrl(send)
	send.Hash = base.SignWithHmacSHA256([]byte(signStr), secret)
	return send
}

// PackQueryPayReq 打包查询支付请求数据
// 构造查询支付请求结构体，包括商户ID和订单号，并生成签名
func (s *Sub) PackQueryPayReq() interface{} {
	r := s.Base.QueryWithdrawReq
	send := &QueryPayReq{
		AppID:      appID,
		TransferID: r.OrderID,
	}
	signStr := sign.GetSignStrl(send)
	send.Hash = base.SignWithHmacSHA256([]byte(signStr), secret)
	return send
}

// CheckSign 验证签名
// 根据回调类型验证签名是否正确，设置回调响应信息
func (s *Sub) CheckSign(str string) bool {
	checkSign := ""
	signStr := ""
	switch s.Base.Opt {
	case base.OPTPayCB:
		req := s.Base.CallbackReq.(*PayCallbackReq)
		log.Debug("checkSign pay:%+v", *req)
		switch req.Status {
		case "SUCCESS":
			s.Base.CallbackResp.Success = true
		case "FAILED":
			s.Base.CallbackResp.Success = false
		default:
			s.Base.CallbackStatus = base.CallbackStatusPass
			return false
		}
		checkSign = req.Hash
		s.Base.CallbackResp.OrderID = req.TransferID
		s.Base.CallbackResp.FailMessage = req.TransferMessage
		signStr = GetSignStrNull(str, s.Base.SignPassStr...)
	case base.OPTWithdrawCB:
		req := s.Base.CallbackReq.(*WithdrawCallbackReq)
		log.Debug("checkSign withdraw:%+v", *req)
		switch req.Status {
		case "SUCCESS":
			s.Base.CallbackResp.Success = true
		case "FAILED":
			s.Base.CallbackResp.Success = false
		default:
			s.Base.CallbackStatus = base.CallbackStatusPass
			return false
		}
		checkSign = req.Hash
		s.Base.CallbackResp.OrderID = req.TransferID
		s.Base.CallbackResp.Utr = req.UTR
		s.Base.CallbackResp.FailMessage = req.TransferMessage
		signStr = GetSignStrNull(str, s.Base.SignPassStr...)
	default:
		return false
	}
	mySign := base.SignWithHmacSHA256([]byte(signStr), secret)
	pass := mySign == checkSign
	if !pass {
		log.Error("check sign fail:mySign:%s,checkSign:%s,signStr:%s", mySign, checkSign, signStr)
	}
	return pass
}

// 根据body里的字段直接拼接出签名字符串(去除null)
func GetSignStrNull(str string, pass ...string) string {
	// 解析成map
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal([]byte(str), &m); err != nil {
		panic(fmt.Sprintf("json unmarshal error: %v", err))
	}

	// 排序字段名
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 构建签名字符串
	var sb strings.Builder
	for _, k := range keys {
		// 跳过指定字段
		skip := false
		for _, p := range pass {
			if k == p {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		val := string(m[k])

		// 去掉json字符串的双引号
		if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
			val = val[1 : len(val)-1]
		}
		if len(val) == 0 || val == "null" {
			continue
		}
		// 空值处理：若参数值为空字符串，仍参与签名
		sb.WriteString(fmt.Sprintf("%s|", val))
	}
	signStr := sb.String()
	if len(signStr) > 0 {
		signStr = signStr[:len(signStr)-1] // 去掉最后的 |
	}
	return signStr
}
