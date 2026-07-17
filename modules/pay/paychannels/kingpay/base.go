package kingpay

import (
	"errors"
	"fmt"
	"net/http"
	"server/common"
	"server/modules/pay/base"
	"server/sign"
	"server/util"
	"sort"
	"strings"

	"github.com/liangdas/mqant/log"
)

func NewSub(b *base.Base) {
	sub := &Sub{
		Base: b,
	}
	b.SignKey = key
	b.HttpType = base.HttpTypeJson
	// b.ShouldSignUpper = true
	// b.KeyName = "null"
	// b.PassKeyName = true
	if b.Opt == base.OPTPay {
		b.Resp = new(PayResp)
		b.ReqURL = baseURL + payURL
	} else if b.Opt == base.OPTWithdraw {
		b.Resp = new(WithdrawResp)
		b.ReqURL = baseURL + withdrawURL
	} else if b.Opt == base.OPTPayCB {
		b.SignPassStr = []string{"sign"}
		b.CallbackResp.Msg = "success"
		b.CallbackReq = new(PayCallbackReq)
	} else if b.Opt == base.OPTWithdrawCB {
		b.SignPassStr = []string{"sign"}
		b.CallbackResp.Msg = "success"
		b.CallbackReq = new(WithdrawCallback)
	} else if b.Opt == base.OPTQueryWithdraw {
		b.Resp = new(QueryWithdrawResp)
		b.ReqURL = baseURL + queryWithdrawURL
	} else if b.Opt == base.OPTQueryPay {
		b.ReqURL = baseURL + queryPayURL
		b.Resp = new(QueryPayResp)
	}
	b.Sub = sub
}

type Sub struct {
	Base *base.Base
}

func (s *Sub) PackHeader(header http.Header) {
}

func (s *Sub) PackReq() interface{} {
	if s.Base.Opt == base.OPTPay {
		return s.PackPayReq()
	} else if s.Base.Opt == base.OPTWithdraw {
		return s.PackWithdrawReq()
	} else if s.Base.Opt == base.OPTQueryPay {
		return s.PackQueryPayReq()
	} else if s.Base.Opt == base.OPTQueryWithdraw {
		return s.PackQueryWithdrawReq()
	}
	return nil
}

func (s *Sub) GetResp() error {
	if s.Base.Opt == 1 {
		resp := s.Base.Resp.(*PayResp)
		if resp.Code != 200 || resp.Data.PayUrl == "" {
			return errors.New("pay fail")
		}
		s.Base.PayResp = &common.PayResp{APIOrderID: resp.Data.OrderNo, URL: resp.Data.PayUrl, Channel: int(s.Base.Channel)}
		return nil
	} else if s.Base.Opt == base.OPTWithdraw {
		resp := s.Base.Resp.(*WithdrawResp)
		if s.Base.Status == 0 && (resp.Code != 200 || resp.Data.OrderNo == "") {
			return errors.New("withdraw fail")
		}
		s.Base.WithdrawResp = &common.WithdrawResp{APIOrderID: resp.Data.OrderNo, Channel: int(s.Base.Channel)}
		return nil
	} else if s.Base.Opt == base.OPTQueryWithdraw {
		resp := s.Base.Resp.(*QueryWithdrawResp)
		log.Debug("QueryWithdrawResp:%+v", resp)
		s.Base.QueryWithdrawResp.OrderID = resp.Data.OrderNo
		s.Base.QueryWithdrawResp.APIOrderID = resp.Data.PaymentOrderNo
		if s.Base.Status == 0 {
			if resp.Data.OrderStatus == 2 {
				s.Base.QueryWithdrawResp.Status = base.QueryStatusOrderSuccess
			} else if resp.Data.OrderStatus == 1 {
				s.Base.QueryWithdrawResp.Status = base.QueryStatusSuccess
			} else if resp.Data.OrderStatus == 3 {
				s.Base.QueryWithdrawResp.Status = base.QueryStatusFail
			}
		}
		if s.Base.QueryWithdrawResp.Status == 0 {
			s.Base.QueryWithdrawResp.Status = base.QueryStatusUnknown
		}
		return nil
	} else if s.Base.Opt == base.OPTQueryPay {
		resp := s.Base.Resp.(*QueryPayResp)
		log.Debug("QueryPayResp:%+v", resp)
		s.Base.QueryPayResp.OrderID = resp.Data.PaymentOrderNo
		if s.Base.Status == 0 {
			if resp.Data.OrderStatus == 2 {
				s.Base.QueryPayResp.Status = base.QueryStatusOrderSuccess
			}
		}
		if s.Base.QueryPayResp.Status == 0 {
			s.Base.QueryPayResp.Status = base.QueryStatusUnknown
		}
		return nil
	}
	return nil
}

func (s *Sub) PackPayReq() interface{} {
	r := s.Base.PayReq
	send := &PayReq{
		Amount:      fmt.Sprintf("%.2f", r.Amount),
		DownOrderNo: r.OrderID,
		KeyId:       KeyID,
		NotifyUrl:   s.Base.GetPayCallbackURL(),
	}
	send.Sign = s.SignMD5(send)
	return send
}

func (s *Sub) PackWithdrawReq() interface{} {
	r := s.Base.WithdrawReq
	send := &WithdrawReq{
		AccountName: r.Name,
		AccountNo:   r.Account,
		Ifsc:        r.IFSC,
		DownOrderNo: r.OrderID,
		Amount:      fmt.Sprintf("%.2f", r.Amount),
		KeyId:       KeyID,
		NotifyUrl:   s.Base.GetWithdrawCallbackURL(),
	}
	if r.PayType == common.WithdrawTypeBank {
	} else {
		return nil
	}
	send.Sign = s.SignMD5(send)
	return send
}

func (s *Sub) PackQueryWithdrawReq() interface{} {
	r := s.Base.QueryWithdrawReq
	send := &QueryPayReq{
		DownOrderNo: r.OrderID,
		KeyId:       KeyID,
	}
	send.Sign = s.SignMD5(send)
	return send
}

func (s *Sub) PackQueryPayReq() interface{} {
	r := s.Base.QueryPayReq
	send := &QueryPayReq{
		DownOrderNo: r.OrderID,
		KeyId:       KeyID,
	}
	send.Sign = s.SignMD5(send)
	return send
}

func (s *Sub) CheckSign(str string) bool {
	// str += key
	// log.Debug("str:%s", str)
	checkSign := ""
	// mySign := strings.ToUpper(util.CalculateMD5(str))
	mySign := ""
	if s.Base.Opt == base.OPTPayCB {
		req := s.Base.CallbackReq.(*PayCallbackReq)
		log.Debug("checkSign pay:%+v", *req)
		checkSign = req.Sign
		s.Base.CallbackResp.OrderID = req.DownOrderNo
		s.Base.CallbackResp.Success = req.OrderStatus == 1
		mySign = s.SignMD5(req)
		// mySign = util.CalculateMD5(appID + req.OrderNumber + req.ActualAmount + req.Status + key)
	} else if s.Base.Opt == base.OPTWithdrawCB {
		req := s.Base.CallbackReq.(*WithdrawCallback)
		log.Debug("checkSign withdraw:%+v", *req)
		// if req.OrderStatus ==1{
		// 	return false
		// }
		checkSign = req.Sign
		s.Base.CallbackResp.OrderID = req.DownOrderNo
		s.Base.CallbackResp.Success = req.OrderStatus == 1
		mySign = s.SignMD5(req)
		// s.Base.CallbackResp.FailMessage = req.Msg
		// mySign = s.Base.SignMD5WithStr(str)
	}
	pass := mySign == checkSign
	if !pass {
		log.Debug("final sign str:%s,sign:%s", str, mySign)
	}
	return pass
}

func (s *Sub) SignMD5(send interface{}) string {
	b := s.Base
	signStr := GetSignStr(send, b.SignPassStr...)
	signStr += key
	log.Debug("final signStr:%s", signStr)
	ret := util.CalculateMD5(signStr)
	if b.ShouldSignUpper {
		ret = strings.ToUpper(ret)
	}
	return ret
}

func GetSignStr(send interface{}, pass ...string) string {
	m := sign.StructToMapJson(send)
	str := []string{}
	for i := range m {
		if i == "sign" {
			continue
		}
		shouldPass := false
		for _, v := range pass {
			if v == i {
				shouldPass = true
				break
			}
		}
		if shouldPass {
			continue
		}
		str = append(str, i)
	}
	sort.Strings(str)
	signStr := ""
	for _, v := range str {
		signStr += fmt.Sprintf("%v", m[v])
	}
	log.Debug("signStr:%s", signStr)
	return signStr
}
