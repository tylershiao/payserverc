package aipay2

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
	b.SignKey = key
	b.HttpType = base.HttpTypeJson
	b.ShouldSignUpper = true
	b.PassKeyName = true
	b.CallbackResp.Msg = "SUCCESS"
	b.SignPassStr = []string{"sign"}
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
	if s.Base.Opt == 1 {
		resp := s.Base.Resp.(*PayResp)
		if resp.Code != 0 || resp.Data.H5Url == "" {
			return errors.New("pay fail")
		}
		s.Base.PayResp = &common.PayResp{APIOrderID: resp.Data.OrderID, URL: resp.Data.H5Url, Channel: int(s.Base.Channel)}
	} else if s.Base.Opt == base.OPTWithdraw {
		resp := s.Base.Resp.(*WithdrawResp)
		if s.Base.Status == 0 && resp.Code != 0 {
			return errors.New("withdraw fail")
		}
		s.Base.WithdrawResp = &common.WithdrawResp{APIOrderID: resp.Data.PayoutID, Message: resp.Data.Msg, Channel: int(s.Base.Channel)}
	} else if s.Base.Opt == base.OPTQueryWithdraw {
		resp := s.Base.Resp.(*QueryWithdrawResp)
		log.Debug("QueryWithdrawResp:%+v", resp)
		s.Base.QueryWithdrawResp.OrderID = resp.Data.MerchantPayoutID
		s.Base.QueryWithdrawResp.APIOrderID = resp.Data.PayoutID
		if s.Base.Status == 0 {
			if resp.Code == 200 {
				if resp.Data.Status == 1 {
					s.Base.QueryWithdrawResp.Status = base.QueryStatusOrderSuccess
				} else if resp.Data.Status == 0 {
					s.Base.QueryWithdrawResp.Status = base.QueryStatusSuccess
				}
			} else {
				if resp.Data.Status == 2 {
					s.Base.QueryWithdrawResp.Status = base.QueryStatusFail
				}
			}
		}
		if s.Base.QueryWithdrawResp.Status == 0 {
			s.Base.QueryWithdrawResp.Status = base.QueryStatusUnknown
		}
	} else if s.Base.Opt == base.OPTQueryPay {
		resp := s.Base.Resp.(*QueryPayResp)
		log.Debug("QueryPayResp:%+v", resp)
		orderID := ""
		apiOrderID := ""
		status := 0
		if len(resp.Data) > 0 {
			orderID = resp.Data[0].MerchantOrderID
			apiOrderID = resp.Data[0].OrderID
			status = resp.Data[0].Status
		}
		s.Base.QueryPayResp.OrderID = orderID
		s.Base.QueryPayResp.APIOrderID = apiOrderID
		if s.Base.Status == 0 {
			if resp.Code == 0 && status == 1 {
				s.Base.QueryPayResp.Status = base.QueryStatusOrderSuccess
			}
		}
		if s.Base.QueryPayResp.Status == 0 {
			s.Base.QueryPayResp.Status = base.QueryStatusUnknown
		}
	} else {
		return errors.New("unknown opt")
	}
	return nil
}

func (s *Sub) PackPayReq() interface{} {
	r := s.Base.PayReq
	send := &PayReq{
		MerchantID:      mid,
		MerchantOrderID: r.OrderID,
		Amount:          fmt.Sprintf("%.2f", r.Amount),
		Timestamp:       time.Now().UnixMilli(),
		// Sign: ,
		PayType:   1,
		NotifyURL: s.Base.GetPayCallbackURL(),
		FirstName: r.Name,
		LastName:  r.Name,
		Mobile:    r.Mobile,
		Email:     r.Email,
	}
	send.Sign = util.CalculateMD5(fmt.Sprintf("merchantId=%d&merchantOrderId=%s&amount=%s&%s", mid, send.MerchantOrderID, send.Amount, key))
	return send
}

func (s *Sub) PackWithdrawReq() interface{} {
	r := s.Base.WithdrawReq
	if r.PayType != common.WithdrawTypeBank {
		return nil
	}
	send := &WithdrawReq{
		MerchantID:      mid,
		MerchantOrderID: r.OrderID,
		Amount:          fmt.Sprintf("%.2f", r.Amount),
		Timestamp:       time.Now().UnixMilli(),
		// Sign: ,
		NotifyURL: s.Base.GetWithdrawCallbackURL(),
	}
	send.FundAccount.Contact.Name = r.Name
	send.FundAccount.Contact.Email = r.Email
	send.FundAccount.Contact.Mobile = r.Mobile
	send.FundAccount.Contact.Address = r.Address
	switch r.PayType {
	case common.WithdrawTypeBank:
		send.FundAccount.AccountType = "bank_account"
		send.FundAccount.BankAccount.IFSC = r.IFSC
		send.FundAccount.BankAccount.AccountNumber = r.Account
		send.FundAccount.BankAccount.Name = r.BankName
	case common.WithdrawTypeUPI:
		send.FundAccount.AccountType = "vpa"
		send.FundAccount.VPA.Address = r.Account
	default:
		return nil
	}
	send.Sign = util.CalculateMD5(fmt.Sprintf("merchantId=%d&merchantOrderId=%s&amount=%s&%s", mid, send.MerchantOrderID, send.Amount, key))
	return send
}

func (s *Sub) PackQueryWithdrawReq() interface{} {
	r := s.Base.QueryWithdrawReq
	send := &QueryWithdrawReq{
		MerchantID:      mid,
		MerchantOrderID: r.OrderID,
		Timestamp:       time.Now().UnixMilli(),
	}
	send.Sign = util.CalculateMD5(fmt.Sprintf("merchantId=%d&merchantOrderId=%s&%s", mid, send.MerchantOrderID, key))
	return send
}

func (s *Sub) PackQueryPayReq() interface{} {
	r := s.Base.QueryPayReq
	send := &QueryPayReq{
		MerchantID:       mid,
		MerchantOrderIDs: []string{r.OrderID},
	}
	send.Sign = util.CalculateMD5(fmt.Sprintf("merchantId=%d&%s", mid, key))
	return send
}

func (s *Sub) CheckSign(str string) bool {
	// str += "&key=" + key
	checkSign := ""
	s.Base.CallbackResp.Msg = "success"
	signStr := ""
	mySign := ""
	switch s.Base.Opt {
	case base.OPTPayCB:
		req := s.Base.CallbackReq.(*PayCallbackReq)
		log.Debug("checkSign pay:%+v", *req)
		checkSign = req.Sign
		s.Base.CallbackResp.OrderID = req.MerchantOrderID
		switch req.Status {
		case 1:
			s.Base.CallbackResp.Success = true
		case 2:
			s.Base.CallbackResp.Success = false
		default:
			s.Base.CallbackStatus = base.CallbackStatusPass
			return false
		}
		signStr = fmt.Sprintf("merchantId=%d&merchantOrderId=%s&amount=%v&%s", mid, req.MerchantOrderID, req.Amount, key)
		mySign = util.CalculateMD5(signStr)
	case base.OPTWithdrawCB:
		req := s.Base.CallbackReq.(*WithdrawCallbackReq)
		log.Debug("checkSign withdraw:%+v", *req)
		checkSign = req.Sign
		s.Base.CallbackResp.OrderID = req.MerchantOrderID
		switch req.Status {
		case 1:
			s.Base.CallbackResp.Success = true
		case 2:
			s.Base.CallbackResp.Success = false
		default:
			s.Base.CallbackStatus = base.CallbackStatusPass
			return false
		}
		s.Base.CallbackResp.APIOrderID = fmt.Sprintf("%d", req.OrderID)
		s.Base.CallbackResp.FailMessage = req.Msg
		signStr = fmt.Sprintf("merchantId=%d&merchantOrderId=%s&amount=%v&%s", mid, req.MerchantOrderID, req.Amount, key)
		mySign = util.CalculateMD5(signStr)
	}
	pass := mySign == checkSign
	if !pass {
		log.Error("sign err,mySign:%s,sign:%s,signStr", mySign, checkSign, signStr)
	}
	return pass
}
