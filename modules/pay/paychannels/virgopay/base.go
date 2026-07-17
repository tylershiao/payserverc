package virgopay

import (
	"errors"
	"net/http"
	"server/common"
	"server/config"
	"server/modules/pay/base"
	"strings"

	"github.com/liangdas/mqant/log"
)

func NewSub(b *base.Base) {
	sub := &Sub{
		Base: b,
	}
	baseURL := config.GetConfig().Pay.VirgoPay.URL
	b.SignKey = config.GetConfig().Pay.VirgoPay.Key
	b.HttpType = base.HttpTypeJson
	b.ShouldSignUpper = true
	if b.Opt == base.OPTPay {
		b.Resp = new(PayResp)
		b.ReqURL = baseURL + payURL
	} else if b.Opt == base.OPTWithdraw {
		b.Resp = new(WithdrawResp)
		b.ReqURL = baseURL + withdrawURL
		b.SignPassStr = []string{"Sign"}
	} else if b.Opt == base.OPTPayCB {
		b.SignPassStr = []string{"Sign"}
		b.CallbackResp.Msg = "SUCCESS"
		b.CallbackReq = new(PayCallbackReq)
	} else if b.Opt == base.OPTWithdrawCB {
		b.SignPassStr = []string{"Sign"}
		b.CallbackResp.Msg = "SUCCESS"
		b.CallbackReq = new(WithdrawCallbackReq)
	} else if b.Opt == base.OPTQueryWithdraw {
		b.Resp = new(QueryWithdrawResp)
		b.ReqURL = baseURL + queryWithdrawURL
	} else if b.Opt == base.OPTQueryPay {
		b.Resp = new(QueryPayResp)
		b.ReqURL = baseURL + queryPayURL
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
	} else if s.Base.Opt == base.OPTQueryWithdraw {
		return s.PackQueryWithdrawReq()
	} else if s.Base.Opt == base.OPTQueryPay {
		return s.PackQueryPayReq()
	}
	return nil
}

func (s *Sub) GetResp() error {
	if s.Base.Opt == 1 {
		resp := s.Base.Resp.(*PayResp)
		if resp.Code != 0 || resp.Data.PayURL == "" {
			return errors.New("pay fail")
		}
		s.Base.PayResp = &common.PayResp{APIOrderID: resp.Data.OrderId, URL: resp.Data.PayURL, Channel: int(s.Base.Channel)}
		return nil
	} else if s.Base.Opt == base.OPTWithdraw {
		resp := s.Base.Resp.(*WithdrawResp)
		if s.Base.Status == 0 && resp.Code != 0 {
			return errors.New("withdraw fail")
		}
		s.Base.WithdrawResp = &common.WithdrawResp{APIOrderID: resp.Data.OrderId, Channel: int(s.Base.Channel), Message: resp.Message}
		return nil
	} else if s.Base.Opt == base.OPTQueryWithdraw {
		resp := s.Base.Resp.(*QueryWithdrawResp)
		log.Debug("QueryWithdrawResp:%+v", resp)
		s.Base.QueryWithdrawResp.OrderID = resp.Data.MOrderId
		s.Base.QueryWithdrawResp.APIOrderID = resp.Data.OrderId
		if s.Base.Status == 0 {
			if resp.Code == 200 {
				if resp.Data.Status == 2 {
					s.Base.QueryWithdrawResp.Status = base.QueryStatusOrderSuccess
				} else if resp.Data.Status == 1 {
					s.Base.QueryWithdrawResp.Status = base.QueryStatusSuccess
				}
			} else {
				if resp.Data.Status == -2 || strings.Contains(resp.Message, "order is not exist") {
					s.Base.QueryWithdrawResp.Status = base.QueryStatusFail
				}
			}
		}
		if s.Base.QueryWithdrawResp.Status == 0 {
			s.Base.QueryWithdrawResp.Status = base.QueryStatusUnknown
		}
		return nil
	} else {
		resp := s.Base.Resp.(*QueryPayResp)
		log.Debug("QueryPayResp:%+v", resp)
		s.Base.QueryPayResp.OrderID = resp.Data.MOrderId
		s.Base.QueryPayResp.APIOrderID = resp.Data.OrderId
		if s.Base.Status == 0 {
			if resp.Code == 200 && resp.Data.Status == 1 {
				s.Base.QueryPayResp.Status = base.QueryStatusOrderSuccess
			}
		}
		if s.Base.QueryPayResp.Status == 0 {
			s.Base.QueryPayResp.Status = base.QueryStatusUnknown
		}
		return nil
	}
}

func (s *Sub) PackPayReq() interface{} {
	r := s.Base.PayReq
	send := &PayReq{
		MchId:     config.GetConfig().Pay.VirgoPay.MID,
		MOrderId:  r.OrderID,
		UserId:    r.UID,
		Amount:    int64(r.Amount * 100),
		IP:        r.IP,
		Name:      r.Name,
		Mobile:    r.Mobile,
		Email:     r.Email,
		NotifyURL: s.Base.GetPayCallbackURL(),
		DeviceID:  r.DeviceID,
	}
	send.Sign = s.Base.SignMD5(send)
	return send
}

func (s *Sub) PackWithdrawReq() interface{} {
	r := s.Base.WithdrawReq
	if r.PayType != common.WithdrawTypeBank {
		return nil
	}
	send := &WithdrawReq{
		MchId:    config.GetConfig().Pay.VirgoPay.MID,
		MOrderId: r.OrderID,
		UserId:   r.UID,
		Amount:   int64(r.Amount * 100),
		IP:       r.IP,
		Name:     r.Name,
		Mobile:   r.Mobile,
		Email:    r.Email,
		DeviceID: r.DeviceID,
		// PayType: ,
		// Account: ,
		// IFSC: ,
		// BankName: ,
		NotifyURL: s.Base.GetWithdrawCallbackURL(),
	}

	if r.PayType == common.WithdrawTypeBank {
		send.PayType = 1
		send.IFSC = r.IFSC
		send.Account = r.Account
		send.BankName = r.BankName
	} else if r.PayType == common.WithdrawTypeUPI {
		send.PayType = 2
		send.Account = r.Account
		return nil
	}
	send.Sign = s.Base.SignMD5(send)
	return send
}

func (s *Sub) PackQueryWithdrawReq() interface{} {
	r := s.Base.QueryWithdrawReq
	send := &QueryWithdrawReq{
		MchId:    config.GetConfig().Pay.VirgoPay.MID,
		MOrderId: r.OrderID,
	}
	send.Sign = s.Base.SignMD5(send)
	return send
}

func (s *Sub) PackQueryPayReq() interface{} {
	r := s.Base.QueryPayReq
	send := &QueryPayReq{
		MchId:    config.GetConfig().Pay.VirgoPay.MID,
		MOrderId: r.OrderID,
	}
	send.Sign = s.Base.SignMD5(send)
	return send
}

func (s *Sub) CheckSign(str string) bool {
	// str += "&key=" + key
	checkSign := ""
	s.Base.CallbackResp.Msg = "success"
	mySign := ""
	if s.Base.Opt == base.OPTPayCB {
		req := s.Base.CallbackReq.(*PayCallbackReq)
		log.Debug("checkSign pay:%+v", *req)
		checkSign = req.Sign
		s.Base.CallbackResp.OrderID = req.MOrderId
		switch req.Status {
		case 3:
			s.Base.CallbackResp.Success = true
		case 4:
			s.Base.CallbackResp.Success = false
		default:
			s.Base.CallbackStatus = base.CallbackStatusPass
			return false
		}
		mySign = s.Base.SignMD5(req)
	} else if s.Base.Opt == base.OPTWithdrawCB {
		req := s.Base.CallbackReq.(*WithdrawCallbackReq)
		log.Debug("checkSign withdraw:%+v", *req)
		checkSign = req.Sign
		s.Base.CallbackResp.OrderID = req.MOrderId
		switch req.Status {
		case 3:
			s.Base.CallbackResp.Success = true
		case 4, 5:
			s.Base.CallbackResp.Success = false
		default:
			s.Base.CallbackStatus = base.CallbackStatusPass
			return false
		}
		s.Base.CallbackResp.APIOrderID = req.OrderId
		s.Base.CallbackResp.FailMessage = req.Message
		s.Base.CallbackResp.Utr = req.Utr
		mySign = s.Base.SignMD5(req)
	}
	return mySign == checkSign
}
