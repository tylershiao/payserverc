package hongxinpay

import (
	"errors"
	"net/http"
	"server/common"
	"server/modules/pay/base"
	"strings"

	"github.com/liangdas/mqant/log"
)

func NewSub(b *base.Base) {
	sub := &Sub{Base: b}
	b.SignKey = key
	b.HttpType = base.HttpTypeJson
	b.ShouldSignUpper = true

	switch b.Opt {
	case base.OPTPay:
		b.ReqURL = baseURL + payURL
		b.Resp = new(PayResp)
	case base.OPTWithdraw:
		b.ReqURL = baseURL + withdrawURL
		b.Resp = new(PayoutResp)
		// 以下字段不参与签名
		b.SignPassStr = []string{"sign"}
	case base.OPTQueryPay:
		b.ReqURL = baseURL + queryPayURL
		b.Resp = new(QueryPayResp)
		b.SignPassStr = []string{"sign"}
	case base.OPTQueryWithdraw:
		b.ReqURL = baseURL + queryWithdrawURL
		b.Resp = new(QueryPayoutResp)
		b.SignPassStr = []string{"sign"}
	case base.OPTPayCB:
		b.SignPassStr = []string{"sign"}
		b.CallbackReq = new(PayCallbackReq)
		b.CallbackResp.Msg = "success"
		b.HttpType = base.HttpTypeGet
	case base.OPTWithdrawCB:
		b.SignPassStr = []string{"sign"}
		b.CallbackReq = new(PayoutCallbackReq)
		b.CallbackResp.Msg = "success"
		b.HttpType = base.HttpTypeGet
	}
	b.Sub = sub
}

type Sub struct {
	Base *base.Base
}

func (s *Sub) PackHeader(header http.Header) {
	// JSON POST，默认 Content-Type 已在 base 里设置
}

func (s *Sub) PackReq() interface{} {
	switch s.Base.Opt {
	case base.OPTPay:
		return s.PackPayReq()
	case base.OPTWithdraw:
		return s.PackWithdrawReq()
	case base.OPTQueryPay:
		return s.PackQueryPayReq()
	case base.OPTQueryWithdraw:
		return s.PackQueryWithdrawReq()
	}
	return nil
}

func (s *Sub) GetResp() error {
	switch s.Base.Opt {
	case base.OPTPay:
		resp := s.Base.Resp.(*PayResp)
		if resp.Code != 0 {
			return errors.New(resp.Msg)
		}
		s.Base.PayResp = &common.PayResp{
			APIOrderID: resp.Data.OrderNo,
			URL:        resp.Data.PayUrl,
			Channel:    int(s.Base.Channel),
		}
	case base.OPTWithdraw:
		resp := s.Base.Resp.(*PayoutResp)
		if resp.Code != 0 {
			return errors.New(resp.Msg)
		}
		s.Base.WithdrawResp = &common.WithdrawResp{
			APIOrderID: resp.Data.OrderNo,
			Channel:    int(s.Base.Channel),
		}
	case base.OPTQueryPay:
		resp := s.Base.Resp.(*QueryPayResp)
		log.Debug("hongxinpay QueryPayResp: %+v", resp)
		s.Base.QueryPayResp.OrderID = resp.Data.MerchantOrderNo
		s.Base.QueryPayResp.APIOrderID = resp.Data.OrderNo
		switch resp.Data.OrderStatus {
		case "ORDER_SUCCESS":
			s.Base.QueryPayResp.Status = base.QueryStatusOrderSuccess
		case "ORDER_FAILED", "ORDER_CLOSED":
			s.Base.QueryPayResp.Status = base.QueryStatusFail
		default:
			s.Base.QueryPayResp.Status = base.QueryStatusUnknown
		}
	case base.OPTQueryWithdraw:
		resp := s.Base.Resp.(*QueryPayoutResp)
		log.Debug("hongxinpay QueryPayoutResp: %+v", resp)
		s.Base.QueryWithdrawResp.OrderID = resp.Data.MerchantOrderNo
		s.Base.QueryWithdrawResp.APIOrderID = resp.Data.OrderNo
		switch resp.Data.OrderStatus {
		case "PAYOUT_DEAL_SUCCESS":
			s.Base.QueryWithdrawResp.Status = base.QueryStatusOrderSuccess
		case "PAYOUT_FAILED", "PAYOUT_DEAL_FAILED":
			s.Base.QueryWithdrawResp.Status = base.QueryStatusFail
		case "PAYOUT_PROCESS":
			s.Base.QueryWithdrawResp.Status = base.QueryStatusSuccess
		default:
			s.Base.QueryWithdrawResp.Status = base.QueryStatusUnknown
		}
	}
	return nil
}

func (s *Sub) PackPayReq() interface{} {
	r := s.Base.PayReq
	req := &PayReq{
		MerchantCode:    merchantCode,
		MerchantOrderNo: r.OrderID,
		OrderAmount:     r.Amount,
		BuyerPhone:      r.Mobile,
		BuyerName:       r.Name,
		BuyerEmail:      r.Email,
		NotifyUrl:       s.Base.GetPayCallbackURL(),
	}
	req.Sign = s.Base.SignMD5(req)
	return req
}

func (s *Sub) PackWithdrawReq() interface{} {
	r := s.Base.WithdrawReq
	req := &PayoutReq{
		MerchantCode:    merchantCode,
		MerchantOrderNo: r.OrderID,
		OrderAmount:     r.Amount,
		BuyerPhone:      r.Mobile,
		BuyerName:       r.Name,
		BuyerEmail:      r.Email,
		BuyerBankNumber: r.Account,
		BuyerIfscCode:   r.IFSC,
		NotifyUrl:       s.Base.GetWithdrawCallbackURL(),
		PayoutType:      "IMPS",
	}
	req.Sign = s.Base.SignMD5(req)
	return req
}

func (s *Sub) PackQueryPayReq() interface{} {
	r := s.Base.QueryPayReq
	req := &QueryPayReq{
		MerchantCode:    merchantCode,
		MerchantOrderNo: r.OrderID,
	}
	req.Sign = s.Base.SignMD5(req)
	return req
}

func (s *Sub) PackQueryWithdrawReq() interface{} {
	r := s.Base.QueryWithdrawReq
	req := &QueryPayoutReq{
		MerchantCode:    merchantCode,
		MerchantOrderNo: r.OrderID,
	}
	req.Sign = s.Base.SignMD5(req)
	return req
}

// callback 验签
func (s *Sub) CheckSign(str string) bool {
	signRemote := s.Base.GetResp.Get("sign")
	mySign := s.Base.SignMD5WithStr(str)
	log.Debug("hongxinpay callback check signRaw=%s mySign=%s", str, mySign)

	if strings.EqualFold(signRemote, mySign) {
		if s.Base.Opt == base.OPTPayCB {
			req := s.Base.CallbackReq.(*PayCallbackReq)
			s.Base.CallbackResp.OrderID = req.MerchantOrderNo
			s.Base.CallbackResp.APIOrderID = req.OrderNo
			s.Base.CallbackResp.Success = req.OrderStatus == "ORDER_SUCCESS"
		} else if s.Base.Opt == base.OPTWithdrawCB {
			req := s.Base.CallbackReq.(*PayoutCallbackReq)
			s.Base.CallbackResp.OrderID = req.MerchantOrderNo
			s.Base.CallbackResp.APIOrderID = req.OrderNo
			if req.OrderStatus == "PAYOUT_DEAL_SUCCESS" {
				s.Base.CallbackResp.Success = true
			} else if req.OrderStatus == "PAYOUT_DEAL_FAILED" {
				s.Base.CallbackResp.Success = false
			}
			s.Base.CallbackResp.Utr = req.Utr
		}
		return true
	}
	return false
}
