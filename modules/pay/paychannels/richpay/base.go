package richpay

import (
	"errors"
	"fmt"
	"net/http"
	"server/common"
	"server/config"
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
	} else if b.Opt == base.OPTQueryWithdraw { // 查询
		b.Resp = new(QueryWithdrawResp)
		b.ReqURL = baseURL + queryWithdrawURL
	} else if b.Opt == base.OPTQueryPay { // 查询
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
	if accessToken == "" {
		GetAccessToken()
	}
	header.Add("Access-Token", accessToken)
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
		if resp.BusCode == "AO-000103" { // token错误
			GetAccessToken()
			return errors.New("pay fail")
		}
		if resp.Code != 200 || resp.Data.PayURL == "" {
			return errors.New("pay fail")
		}
		s.Base.PayResp = &common.PayResp{APIOrderID: fmt.Sprintf("%v", resp.Data.OrderNo), URL: resp.Data.PayURL, Channel: int(s.Base.Channel)}
		return nil
	} else if s.Base.Opt == base.OPTWithdraw {
		resp := s.Base.Resp.(*WithdrawResp)
		if resp.BusCode == "AO-000103" { // token错误
			GetAccessToken()
			return errors.New("withdraw fail")
		}
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
		if resp.BusCode == "AO-000103" { // token错误
			GetAccessToken()
			s.Base.QueryWithdrawResp.Status = base.QueryStatusUnknown
			return nil
		}
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
		if resp.BusCode == "AO-000103" { // token错误
			GetAccessToken()
			s.Base.QueryWithdrawResp.Status = base.QueryStatusUnknown
			return nil
		}
		if resp.Code != 200 {
			s.Base.QueryPayResp.Status = base.QueryStatusUnknown
			return nil
		}
		s.Base.QueryWithdrawResp.OrderID = resp.Data.MerchantOrderNo
		s.Base.QueryPayResp.Msg = resp.Msg
		if s.Base.Status != 0 {
			s.Base.QueryPayResp.Status = base.QueryStatusUnknown
		} else {
			if resp.Data.Status == "2" {
				s.Base.QueryPayResp.Status = base.QueryStatusOrderSuccess
			} else if resp.Data.Status == "3" {
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
		MerchantID:      mid,
		MerchantOrderNo: r.OrderID,
		OrderAmount:     util.RoundFloat(r.Amount, 2),
		Email:           r.Email,
		Name:            r.Name,
		Phone:           r.Mobile,
		DeeplinkSwitch:  "1",
		NotifyURL:       s.Base.GetPayCallbackURL(),
		NonceStr:        util.GenerateRandomString(6),
		Action:          "payin",
		Timestamp:       fmt.Sprintf("%d", time.Now().Unix()),
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
		MerchantID:      mid,
		MerchantOrderNo: r.OrderID,
		OrderAmount:     int64(r.Amount),
		Name:            r.Name,
		Email:           r.Email,
		Phone:           r.Mobile,
		NotifyURL:       s.Base.GetWithdrawCallbackURL(),
		NonceStr:        util.GenerateRandomString(6),
		Action:          "payout",
		Timestamp:       fmt.Sprintf("%d", time.Now().Unix()),
		// Type: ,
		// VPA: ,
		// Account: ,
		// IFSC: ,
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
	send := &QueryWithdrawReq{
		MerchantOrderNo: r.OrderID,
	}
	return send
}

func (s *Sub) PackQueryPayReq() interface{} {
	r := s.Base.QueryPayReq
	send := &QueryPayReq{
		MerchantOrderNo: r.OrderID,
	}
	return send
}

func GetAccessToken() {
	url := baseURL + tokenURL
	resp := &TokenResp{}
	header := map[string]string{"clientid": mid, "secret": secret}
	util.HttpPost(url, nil, resp, header)
	if resp.Code == 200 {
		accessToken = resp.AccessToken
		refreshToken = resp.RefreshToken
		next := resp.ExpiresIn - 200
		if next <= 0 {
			next = 60
		}
		if config.GetConfig().WorkID == 1 { // 只有a服去刷新token
			log.Debug("next refresh token:%d", next)
			tokenTimerLock.Lock()
			if tokenTimer != nil {
				tokenTimer.Stop()
			}
			tokenTimer = time.AfterFunc(time.Duration(next)*time.Second, RefreshAccessToken)
			tokenTimerLock.Unlock()
		} else {
			log.Debug("next GetAccessToken token:%d", next+1)
			tokenTimerLock.Lock()
			if tokenTimer != nil {
				tokenTimer.Stop()
			}
			tokenTimer = time.AfterFunc(time.Duration(next+1)*time.Second, GetAccessToken)
			tokenTimerLock.Unlock()
		}
	} else {
		tokenTimerLock.Lock()
		if tokenTimer != nil {
			tokenTimer.Stop()
		}
		log.Debug("GetAccessToken fail resp:%+v", resp)
		tokenTimer = time.AfterFunc(10*time.Second, GetAccessToken)
		tokenTimerLock.Unlock()
	}
}

func RefreshAccessToken() {
	url := baseURL + refreshTokenURL
	resp := &TokenResp{}
	header := map[string]string{"Refresh-Token": refreshToken}
	util.HttpPost(url, nil, resp, header)
	if resp.Code == 200 {
		accessToken = resp.AccessToken
		refreshToken = resp.RefreshToken
		tokenTimerLock.Lock()
		if tokenTimer != nil {
			tokenTimer.Stop()
		}
		tokenTimer = time.AfterFunc(time.Duration(resp.ExpiresIn-200)*time.Second, RefreshAccessToken)
		tokenTimerLock.Unlock()
	} else {
		tokenTimerLock.Lock()
		if tokenTimer != nil {
			tokenTimer.Stop()
		}
		tokenTimer = time.AfterFunc(10*time.Second, GetAccessToken)
		tokenTimerLock.Unlock()
	}
}
