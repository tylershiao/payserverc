package dkpay

import (
	"errors"
	"fmt"
	"net/http"
	"server/common"
	"server/modules/pay/base"
	"server/util"
	"time"

	"github.com/liangdas/mqant/log"
)

func NewSub(b *base.Base) {
	sub := &Sub{
		Base: b,
	}
	b.HttpType = base.HttpTypeJson
	// b.ShouldQueryWithdraw = true
	if b.Opt == base.OPTPay {
		b.Resp = new(PayResp)
		b.ReqURL = baseURL + payURL
	} else if b.Opt == base.OPTWithdraw {
		b.Resp = new(WithdrawResp)
		b.ReqURL = baseURL + withdrawURL
	} else if b.Opt == base.OPTPayCB {
		b.CallbackReq = new(PayCallbackReq)
		b.CallbackResp.Msg = "success"
	} else if b.Opt == base.OPTWithdrawCB {
		b.CallbackReq = new(WithdrawCallbackReq)
		b.CallbackResp.Msg = "success"
	} else if b.Opt == base.OPTQueryWithdraw {
		b.Resp = new(QueryWithdrawResp)
		b.ReqURL = baseURL + queryWithdrawURL
	} else if b.Opt == base.OPTQueryPay {
		b.Resp = new(QueryPayResp)
		b.ReqURL = baseURL + queryPayURL
	} else {
		return
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
	log.Debug("resp:%v", s.Base.Resp)
	if s.Base.Opt == base.OPTPay {
		resp := s.Base.Resp.(*PayResp)
		if resp.Code != 200 || resp.Data.PayURL == "" {
			return errors.New("pay fail")
		}
		s.Base.PayResp = &common.PayResp{APIOrderID: fmt.Sprintf("%v", resp.Data.OrderNo), URL: resp.Data.PayURLH5, Channel: int(s.Base.Channel)}
		return nil
	} else if s.Base.Opt == base.OPTWithdraw {
		resp := s.Base.Resp.(*WithdrawResp)
		if resp.Msg == "SV-000000" {
			return errors.New("withdraw fail")
		}
		if s.Base.Status == 0 && resp.Code != 200 {
			return errors.New("withdraw fail")
		}
		s.Base.WithdrawResp = &common.WithdrawResp{APIOrderID: fmt.Sprintf("%v", resp.Data.OrderNo), Channel: int(s.Base.Channel)}
		return nil
	} else if s.Base.Opt == base.OPTQueryWithdraw {
		resp := s.Base.Resp.(*QueryWithdrawResp)
		if resp.Code != 200 {
			s.Base.QueryWithdrawResp.Status = base.QueryStatusUnknown
			return nil
		}
		s.Base.QueryWithdrawResp.OrderID = resp.Data.MerchantOrderNo
		s.Base.QueryWithdrawResp.Msg = resp.Msg
		if s.Base.Status != 0 {
			s.Base.QueryWithdrawResp.Status = base.QueryStatusUnknown
		} else {
			if resp.Data.Status == "2" {
				s.Base.QueryWithdrawResp.Status = base.QueryStatusOrderSuccess
			} else if resp.Data.Status == "3" || resp.Data.Status == "4" || resp.Data.Status == "5" {
				s.Base.QueryWithdrawResp.Status = base.QueryStatusFail
			} else {
				s.Base.QueryWithdrawResp.Status = base.QueryStatusSuccess
			}
		}
		return nil
	} else if s.Base.Opt == base.OPTQueryPay {
		resp := s.Base.Resp.(*QueryPayResp)
		if resp.Code != 200 {
			s.Base.QueryPayResp.Status = base.QueryStatusUnknown
			return nil
		}
		s.Base.QueryPayResp.OrderID = resp.Data.MerchantOrderNo
		s.Base.QueryPayResp.Msg = resp.Msg
		if s.Base.Status != 0 {
			s.Base.QueryPayResp.Status = base.QueryStatusUnknown
		} else {
			if resp.Data.Status == "2" {
				s.Base.QueryPayResp.Status = base.QueryStatusOrderSuccess
			} else if resp.Data.Status == "3" || resp.Data.Status == "4" {
				s.Base.QueryPayResp.Status = base.QueryStatusFail
			} else {
				s.Base.QueryPayResp.Status = base.QueryStatusSuccess
			}
		}
		return nil
	}
	return errors.New("unknown opt")
}

func (s *Sub) PackPayReq() interface{} {
	r := s.Base.PayReq
	send := &PayReq{
		MerchantID:      util.GetInt64(mid),
		MerchantOrderNo: r.OrderID,
		OrderAmount:     util.RoundFloat(r.Amount, 2),
		Email:           r.Email,
		Name:            r.Name,
		Phone:           r.Mobile,
		DeeplinkSwitch:  "1",
		NotifyURL:       s.Base.GetPayCallbackURL(),
		NonceStr:        util.GenerateRandomString(6),
		Action:          "payin",
		Timestamp:       fmt.Sprintf("%d", time.Now().UnixMilli()),
	}
	sign, err := base.SignDsa(mid+send.MerchantOrderNo+send.OrderAmount+send.NonceStr+send.Timestamp+send.Action, privateKey)
	if err != nil {
		log.Error("err:%v", err)
		return nil
	}
	send.Sign = sign
	return send
}

func (s *Sub) PackWithdrawReq() interface{} {
	r := s.Base.WithdrawReq
	send := &WithdrawReq{
		MerchantID:      util.GetInt64(mid),
		MerchantOrderNo: r.OrderID,
		OrderAmount:     int64(r.Amount),
		Name:            r.Name,
		Email:           r.Email,
		Phone:           r.Mobile,
		NotifyURL:       s.Base.GetWithdrawCallbackURL(),
		NonceStr:        util.GenerateRandomString(6),
		Action:          "payout",
		Timestamp:       fmt.Sprintf("%d", time.Now().UnixMilli()),
	}
	if r.PayType == common.WithdrawTypeBank {
		send.Type = "BANK"
		send.Account = r.Account
		send.IFSC = r.IFSC
	} else if r.PayType == common.WithdrawTypeUPI {
		send.Type = "UPI"
		send.VPA = r.Account
	} else {
		return nil
	}
	str := fmt.Sprintf("%s%s%d%s%s%s", mid, send.MerchantOrderNo, send.OrderAmount, send.NonceStr, send.Timestamp, send.Action)
	sign, err := base.SignDsa(str, privateKey)
	if err != nil {
		log.Error("err:%v", err)
		return nil
	}
	send.Sign = sign
	return send
}

func (s *Sub) CheckSign(str string) bool {
	log.Debug("callback:%v", s.Base.CallbackReq)
	if s.Base.Opt == base.OPTPayCB {
		req := s.Base.CallbackReq.(*PayCallbackReq)
		s.Base.CallbackResp.OrderID = req.MerchantOrderNo
		s.Base.CallbackResp.APIOrderID = req.OrderNo
		s.Base.CallbackResp.Utr = req.UTR
		if req.Status == "0" || req.Status == "1" {
			return false
		}
		s.Base.CallbackResp.Success = req.Status == "2"
		message := req.MerchantOrderNo + req.OrderNo + req.OrderAmount + req.NonceStr + req.Timestamp
		pass, err := base.VerifyDsa(message, req.Sign, publicKey)
		if err != nil {
			log.Error("err:%v", err)
			return false
		}
		if !pass {
			log.Error("check sign fail,message:%s", message)
		}
		return pass
	} else if s.Base.Opt == base.OPTWithdrawCB {
		req := s.Base.CallbackReq.(*WithdrawCallbackReq)
		s.Base.CallbackResp.OrderID = req.MerchantOrderNo
		s.Base.CallbackResp.APIOrderID = req.OrderNo
		s.Base.CallbackResp.Utr = req.UTR
		if req.Status == "0" || req.Status == "1" {
			return false
		}
		s.Base.CallbackResp.Success = req.Status == "2"
		message := req.MerchantOrderNo + req.OrderNo + req.OrderAmount + req.NonceStr + req.Timestamp
		pass, err := base.VerifyDsa(message, req.Sign, publicKey)
		if err != nil {
			log.Error("err:%v", err)
			return false
		}
		if !pass {
			log.Error("check sign fail,message:%s", message)
		}
		return pass
	}
	return false
}

func (s *Sub) PackQueryWithdrawReq() interface{} {
	r := s.Base.QueryWithdrawReq
	sign, err := base.SignDsa(mid+r.OrderID, privateKey)
	if err != nil {
		log.Error("err:%v", err)
		return nil
	}
	send := &QueryWithdrawReq{
		MerchantID:      util.GetInt64(mid),
		MerchantOrderNo: r.OrderID,
		Sign:            sign,
	}
	return send
}

func (s *Sub) PackQueryPayReq() interface{} {
	r := s.Base.QueryPayReq
	sign, err := base.SignDsa(mid+r.OrderID, privateKey)
	if err != nil {
		log.Error("err:%v", err)
		return nil
	}
	send := &QueryPayReq{
		MerchantID:      util.GetInt64(mid),
		MerchantOrderNo: r.OrderID,
		Sign:            sign,
	}
	return send
}
