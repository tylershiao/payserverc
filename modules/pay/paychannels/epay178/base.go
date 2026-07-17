package epay178

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"server/common"
	"server/modules/pay/base"
	"server/sign"
	"server/util"
	"strings"

	"github.com/liangdas/mqant/log"
)

func NewSub(b *base.Base) {
	sub := &Sub{
		Base: b,
	}
	b.HttpType = base.HttpTypeJson
	b.CallbackResp.Msg = "success"
	b.SignPassStr = []string{"sign", "data"}
	b.SignKey = key
	b.KeyName = "secret"
	b.ShouldSignUpper = true
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

type Sub struct {
	Base *base.Base
}

func (s *Sub) PackHeader(header http.Header) {
}

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

func (s *Sub) GetResp() error {
	switch s.Base.Opt {
	case base.OPTPay:
		resp := s.Base.Resp.(*PayResp)
		if resp.ReasonCode != "000000" || resp.Data.PayURL == "" {
			return errors.New("pay fail")
		}
		s.Base.PayResp = &common.PayResp{APIOrderID: resp.Data.TransactionID, URL: resp.Data.PayURL, Channel: int(s.Base.Channel)}
	case base.OPTQueryPay:
		resp := s.Base.Resp.(*QueryPayResp)
		log.Debug("QueryPayResp:%+v", resp)
		s.Base.QueryPayResp.OrderID = resp.Data.OutTradeNo
		s.Base.QueryPayResp.APIOrderID = resp.Data.TransactionID
		if s.Base.Status == 0 {
			if resp.ReasonCode == "000000" && resp.Data.Paid == "1" {
				s.Base.QueryPayResp.Status = base.QueryStatusOrderSuccess
			}
		}
		if s.Base.QueryPayResp.Status == 0 {
			s.Base.QueryPayResp.Status = base.QueryStatusUnknown
		}
	case base.OPTWithdraw:
		resp := s.Base.Resp.(*WithdrawResp)
		if s.Base.Status == 0 && resp.ReasonCode != "000000" {
			return errors.New("withdraw fail")
		}
		s.Base.WithdrawResp = &common.WithdrawResp{APIOrderID: resp.Data.TransactionID, Channel: int(s.Base.Channel)}
	case base.OPTQueryWithdraw:
		resp := s.Base.Resp.(*QueryWithdrawResp)
		log.Debug("QueryWithdrawResp:%+v", resp)
		s.Base.QueryWithdrawResp.OrderID = resp.Data.OutTradeNo
		s.Base.QueryWithdrawResp.APIOrderID = resp.Data.TransactionID
		if s.Base.Status == 0 {
			if resp.ReasonCode == "000000" {
				switch resp.Data.Result {
				case "1":
					s.Base.QueryWithdrawResp.Status = base.QueryStatusOrderSuccess
				case "0":
					s.Base.QueryWithdrawResp.Status = base.QueryStatusSuccess
				}
			} else if resp.Data.Result == "2" || resp.Data.Result == "4" {
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

func (s *Sub) PackPayReq() interface{} {
	r := s.Base.PayReq
	send := &PayReq{
		MchID:      mid,
		OutTradeNo: r.OrderID,
		Type:       "1004",
		Amount:     int64(r.Amount * 100),
		Currency:   "INR",
		Subject:    "shop",
		NotifyURL:  s.Base.GetPayCallbackURL(),
		Nonce:      fmt.Sprintf("%d", rand.Intn(9000000)+1000000),
	}
	payExtra := PayExtra{}
	payExtra.Email = r.Email
	payExtra.Name = r.Name
	payExtra.Phone = r.Mobile
	str, _ := json.Marshal(payExtra)
	send.Extra = string(str)
	send.Sign = s.Base.SignMD5(send)
	return send
}

func (s *Sub) PackWithdrawReq() interface{} {
	r := s.Base.WithdrawReq
	send := &WithdrawReq{
		MchID:      mid,
		OutTradeNo: r.OrderID,
		Type:       "2001",
		Amount:     int64(r.Amount * 100),
		Currency:   "INR",
		Subject:    "shop",
		NotifyURL:  s.Base.GetWithdrawCallbackURL(),
		Nonce:      fmt.Sprintf("%d", rand.Intn(9000000)+1000000),
	}
	withdrawExtra := WithdrawExtra{}
	withdrawExtra.Email = r.Email
	withdrawExtra.Name = r.Name
	withdrawExtra.Phone = r.Mobile

	switch r.PayType {
	case common.WithdrawTypeBank:
		withdrawExtra.TransType = "BANK_TRANSFER"
		withdrawExtra.TransMode = "IMPS"
		withdrawExtra.Account = r.Account
		withdrawExtra.IFSC = r.IFSC
		withdrawExtra.BankName = r.BankName
	case common.WithdrawTypeUPI:
		withdrawExtra.TransType = "UPI"
		withdrawExtra.VPA = r.Account
	default:
		return nil
	}
	str, _ := json.Marshal(withdrawExtra)
	send.Extra = string(str)
	send.Sign = s.Base.SignMD5(send)
	return send
}

func (s *Sub) PackQueryWithdrawReq() interface{} {
	r := s.Base.QueryWithdrawReq
	send := &QueryWithdrawReq{
		MchID:         mid,
		OutTradeNo:    r.OrderID,
		TransactionID: r.APIOrderID,
		Nonce:         fmt.Sprintf("%d", rand.Intn(9000000)+1000000),
	}
	send.Sign = s.Base.SignMD5(send)
	return send
}

func (s *Sub) PackQueryPayReq() interface{} {
	r := s.Base.QueryPayReq
	send := &QueryPayReq{
		MchID:         mid,
		OutTradeNo:    r.OrderID,
		TransactionID: r.APIOrderID,
		Nonce:         fmt.Sprintf("%d", rand.Intn(9000000)+1000000),
	}
	send.Sign = s.Base.SignMD5(send)
	return send
}

func (s *Sub) CheckSign(str string) bool {
	checkSign := ""
	switch s.Base.Opt {
	case base.OPTPayCB:
		req := s.Base.CallbackReq.(*PayCallbackReq)
		log.Debug("checkSign pay:%+v", *req)
		switch req.Paid {
		case "1":
			s.Base.CallbackResp.Success = true
		case "2":
			s.Base.CallbackResp.Success = false
		default:
			s.Base.CallbackStatus = base.CallbackStatusPass
			return false
		}
		checkSign = req.Sign
		s.Base.CallbackResp.OrderID = req.OutTradeNo
		s.Base.CallbackResp.APIOrderID = req.TransactionID
		s.Base.CallbackResp.Utr = req.UTR
	case base.OPTWithdrawCB:
		req := s.Base.CallbackReq.(*WithdrawCallbackReq)
		log.Debug("checkSign withdraw:%+v", *req)
		switch req.Result {
		case "1":
			s.Base.CallbackResp.Success = true
		case "2", "4":
			s.Base.CallbackResp.Success = false
		default:
			s.Base.CallbackStatus = base.CallbackStatusPass
			return false
		}
		checkSign = req.Sign
		s.Base.CallbackResp.OrderID = req.OutTradeNo
		s.Base.CallbackResp.APIOrderID = req.TransactionID
		s.Base.CallbackResp.Utr = req.UTR
		s.Base.CallbackResp.FailMessage = req.Remark
	default:
		return false
	}
	signStr := sign.GetSignStrNull(str, s.Base.SignPassStr...)
	mySign := strings.ToUpper(util.CalculateMD5(signStr + "&" + s.Base.KeyName + "=" + key))
	pass := mySign == checkSign
	if !pass {
		log.Error("check sign fail:mySign:%s,checkSign:%s", mySign, checkSign)
	}
	return pass
}
