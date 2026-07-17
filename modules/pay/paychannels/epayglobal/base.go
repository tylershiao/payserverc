package epayglobal

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"server/common"
	"server/config"
	"server/modules/pay/base"
	"server/util"
	"sort"

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
	b.SignPassStr = []string{"sign"}
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
		if resp.Code != 200 || resp.PayURL == "" {
			return errors.New("pay fail")
		}
		s.Base.PayResp = &common.PayResp{APIOrderID: resp.DisOrderNo, URL: resp.PayURL, Channel: int(s.Base.Channel)}
	case base.OPTQueryPay:
		resp := s.Base.Resp.(*QueryPayResp)
		log.Debug("QueryPayResp:%+v", resp)
		s.Base.QueryPayResp.OrderID = resp.OrderNo
		s.Base.QueryPayResp.APIOrderID = resp.DisOrderNo
		if s.Base.Status == 0 {
			if resp.Code == 200 && resp.Status == 2 {
				s.Base.QueryPayResp.Status = base.QueryStatusOrderSuccess
			}
		}
		if s.Base.QueryPayResp.Status == 0 {
			s.Base.QueryPayResp.Status = base.QueryStatusUnknown
		}
	case base.OPTWithdraw:
		resp := s.Base.Resp.(*WithdrawResp)
		if s.Base.Status == 0 && resp.Code != 200 {
			return errors.New("withdraw fail")
		}
		s.Base.WithdrawResp = &common.WithdrawResp{APIOrderID: resp.DisOrderNo, Channel: int(s.Base.Channel)}
	case base.OPTQueryWithdraw:
		resp := s.Base.Resp.(*QueryWithdrawResp)
		log.Debug("QueryWithdrawResp:%+v", resp)
		s.Base.QueryWithdrawResp.OrderID = resp.OrderNo
		s.Base.QueryWithdrawResp.APIOrderID = resp.DisOrderNo
		if s.Base.Status == 0 && resp.Code == 200 {
			switch resp.Status {
			case 2:
				s.Base.QueryWithdrawResp.Status = base.QueryStatusOrderSuccess
			case 1:
				s.Base.QueryWithdrawResp.Status = base.QueryStatusSuccess
			case 3:
				s.Base.QueryWithdrawResp.Status = base.QueryStatusFail
			}
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
		TradeNo: mid,
		AppID:   appID,
		PayCode: config.GetConfig().Pay.EpayGlobal.PayCode,
		// PayCode: 2063,
		Price:   int(r.Amount * 100),
		OrderNo: r.OrderID,
		// SuccessURL: values.GetFrontCallback(),
		// FailURL: ,
		PayNoticeURL: s.Base.GetPayCallbackURL(),
		UserID:       r.UID,
		UserIP:       r.IP,
		// Attach: ,
		// Sign: ,
	}
	send.Sign = s.Base.SignMD5(send)
	return send
}

// PackWithdrawReq 打包提现请求数据
// 构造提现请求结构体，根据支付类型（银行、UPI）设置不同的字段，并生成签名
func (s *Sub) PackWithdrawReq() interface{} {
	r := s.Base.WithdrawReq
	send := &WithdrawReq{
		TradeNo: mid,
		OrderNo: r.OrderID,
		AppID:   appID,
		PayCode: config.GetConfig().Pay.EpayGlobal.WithdrawPayCode,
		// PayCode: 2066,
		Price: int(r.Amount * 100),
		// AccountNo: r.Account,
		// AccountType: ,
		AccountName: r.Name,
		// BankCode: ,
		// IdentifyNum: ,
		PayNoticeURL: s.Base.GetWithdrawCallbackURL(),
		// Attach: ,
	}
	switch r.PayType {
	case common.WithdrawTypeBank:
		send.AccountType = "BANK"
		send.BankCode = "INDIAN"
		send.AccountNo = r.Account
		send.IdentifyNum = r.IFSC
	case common.WithdrawTypeUPI:
		send.AccountType = "UPI"
		send.BankCode = "UPI"
		send.AccountNo = r.Account
	default:
		return nil
	}
	send.Sign = s.Base.SignMD5(send)
	return send
}

// PackQueryWithdrawReq 打包查询提现请求数据
// 构造查询提现请求结构体，包括商户ID和订单号，并生成签名
func (s *Sub) PackQueryWithdrawReq() interface{} {
	r := s.Base.QueryWithdrawReq
	send := &QueryWithdrawReq{
		OrderType: "pay_out",
		TradeNo:   mid,
		OrderNo:   r.OrderID,
		// BankOrderNo: ,
	}
	send.Sign = s.Base.SignMD5(send)
	return send
}

// PackQueryPayReq 打包查询支付请求数据
// 构造查询支付请求结构体，包括商户ID和订单号，并生成签名
func (s *Sub) PackQueryPayReq() interface{} {
	r := s.Base.QueryPayReq
	send := &QueryPayReq{
		OrderType: "pay_in",
		TradeNo:   mid,
		OrderNo:   r.OrderID,
		// BankOrderNo: ,
	}
	send.Sign = s.Base.SignMD5(send)
	return send
}

// CheckSign 验证签名
// 根据回调类型验证签名是否正确，设置回调响应信息
func (s *Sub) CheckSign(str string) bool {
	checkSign := ""
	switch s.Base.Opt {
	case base.OPTPayCB:
		req := s.Base.CallbackReq.(*PayCallbackReq)
		log.Debug("checkSign pay:%+v", *req)
		switch req.Status {
		case 2:
			s.Base.CallbackResp.Success = true
		case 3:
			s.Base.CallbackResp.Success = false
		default:
			s.Base.CallbackStatus = base.CallbackStatusPass
			return false
		}
		checkSign = req.Sign
		s.Base.CallbackResp.OrderID = req.OrderNo
	case base.OPTWithdrawCB:
		req := s.Base.CallbackReq.(*WithdrawCallbackReq)
		log.Debug("checkSign withdraw:%+v", *req)
		switch req.Status {
		case 2:
			s.Base.CallbackResp.Success = true
		case 3, 7, 8:
			s.Base.CallbackResp.Success = false
		default:
			s.Base.CallbackStatus = base.CallbackStatusPass
			return false
		}
		checkSign = req.Sign
		s.Base.CallbackResp.OrderID = req.OrderNo
		// s.Base.CallbackResp.Utr = req.ReferenceNo
	default:
		return false
	}
	signStr := GetSignStrNull(str, s.Base.SignPassStr...)
	mySign := util.CalculateMD5(signStr + "&" + s.Base.KeyName + "=" + key)
	pass := mySign == checkSign
	if !pass {
		log.Error("check sign fail:mySign:%s,checkSign:%s", mySign, checkSign)
	}
	return pass
}

// 根据body里的字段直接拼接出签名字符串(去除null)
func GetSignStrNull(str string, pass ...string) string {
	m := map[string]json.RawMessage{}
	sortStrs := []string{}
	json.Unmarshal([]byte(str), &m)
	for i := range m {
		sortStrs = append(sortStrs, i)
	}
	signStr := ""
	sort.Strings(sortStrs)
	for _, v := range sortStrs {
		shouldPass := false
		for _, s := range pass {
			if v == s {
				shouldPass = true
				break
			}
		}
		if shouldPass {
			continue
		}

		raw := m[v]
		var val string
		if len(raw) > 1 && raw[0] == 34 {
			var s string
			if err := json.Unmarshal(raw, &s); err == nil {
				// 成功反解，说明是被转义的字符串
				val = s
			} else {
				// 反解失败，兜底逻辑
				val = string(raw[1 : len(raw)-1])
			}
		} else {
			val = string(raw)
		}
		if len(val) == 0 || string(val) == "null" {
			continue
		}
		signStr += fmt.Sprintf("%v=%v", v, val)
		signStr += "&"
	}
	signStr = signStr[:len(signStr)-1]
	log.Debug("signStr:%v", signStr)
	return signStr
}
