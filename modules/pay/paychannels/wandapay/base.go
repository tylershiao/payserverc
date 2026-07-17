package wandapay

import (
	"errors"
	"github.com/liangdas/mqant/log"
	"github.com/shopspring/decimal"
	"net/http"
	"server/common"
	"server/modules/pay/base"
	"strconv"
)

// NewSub 初始化子支付模块
// 根据不同的操作类型（支付、提现、回调等）设置相应的配置和响应结构体
func NewSub(b *base.Base) {
	sub := &Sub{
		Base: b,
	}

	_, key := getWandaConfig()
	b.SignKey = key
	b.HttpType = base.HttpTypeJson
	b.ShouldSignUpper = true

	if b.Opt == base.OPTPay {
		b.Resp = new(PayResp)
		b.ReqURL = baseURL + payURL
	} else if b.Opt == base.OPTWithdraw {
		b.Resp = new(WithdrawResp)
		b.ReqURL = baseURL + withdrawURL
		b.SignPassStr = []string{"sign"}
	} else if b.Opt == base.OPTPayCB {
		b.SignPassStr = []string{"sign"}
		b.CallbackResp.Msg = "SUCCESS"
		b.CallbackReq = new(PayCallbackReq)
	} else if b.Opt == base.OPTWithdrawCB {
		b.SignPassStr = []string{"sign"}
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

// Sub 是 Wanda 支付模块的子结构体
type Sub struct {
	Base *base.Base
}

// PackHeader 设置 HTTP 请求头信息
func (s *Sub) PackHeader(header http.Header) {
	header.Set("User-Agent", "WandaPay-Go-Client/1.0")
	header.Set("Accept", "application/json")
	header.Set("Accept-Encoding", "deflate")
	header.Set("Accept-Language", "en-US,en;q=0.9")
	header.Set("Connection", "keep-alive")
	header.Set("Cache-Control", "no-cache")
}

// PackReq 根据操作类型打包请求数据
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

// GetResp 处理响应数据并填充到基础结构体中
func (s *Sub) GetResp() error {
	if s.Base.Opt == base.OPTPay {
		resp := s.Base.Resp.(*PayResp)
		if !resp.IsSuccess() {
			return errors.New("wanda pay create order failed: " + resp.Message)
		}
		s.Base.PayResp = &common.PayResp{
			APIOrderID: resp.TxnID,
			URL:        resp.PayLink,
			Channel:    int(s.Base.Channel),
		}
		return nil
	} else if s.Base.Opt == base.OPTWithdraw {
		resp := s.Base.Resp.(*WithdrawResp)
		if s.Base.Status == 0 && !resp.IsSuccess() {
			return errors.New("withdraw fail: " + resp.Message)
		}
		s.Base.WithdrawResp = &common.WithdrawResp{
			APIOrderID: resp.TxnID,
			Channel:    int(s.Base.Channel),
		}
		return nil
	} else if s.Base.Opt == base.OPTQueryWithdraw {
		resp := s.Base.Resp.(*QueryWithdrawResp)
		log.Debug("QueryWithdrawResp:%+v", resp)
		s.Base.QueryWithdrawResp.OrderID = resp.TxnID
		s.Base.QueryWithdrawResp.APIOrderID = resp.TxnID
		if s.Base.Status == 0 {
			switch resp.Status {
			case "SUCCESS":
				s.Base.QueryWithdrawResp.Status = base.QueryStatusOrderSuccess
			case "PENDING":
				s.Base.QueryWithdrawResp.Status = base.QueryStatusSuccess
			case "FAILED":
				s.Base.QueryWithdrawResp.Status = base.QueryStatusFail
			}
		}
		if s.Base.QueryWithdrawResp.Status == 0 {
			s.Base.QueryWithdrawResp.Status = base.QueryStatusUnknown
		}
		return nil
	} else {
		resp := s.Base.Resp.(*QueryPayResp)
		log.Debug("QueryPayResp:%+v", resp)
		s.Base.QueryPayResp.OrderID = resp.TxnID
		s.Base.QueryPayResp.APIOrderID = resp.TxnID
		if s.Base.Status == 0 {
			if resp.Status == "SUCCESS" {
				s.Base.QueryPayResp.Status = base.QueryStatusOrderSuccess
			} else if resp.Status == "FAILED" {
				s.Base.QueryPayResp.Status = base.QueryStatusFail
			}
		}
		if s.Base.QueryPayResp.Status == 0 {
			s.Base.QueryPayResp.Status = base.QueryStatusUnknown
		}
		return nil
	}
}

// PackPayReq 打包支付请求数据
func (s *Sub) PackPayReq() interface{} {
	r := s.Base.PayReq
	merchantNo, _ := getWandaConfig()

	// 金额转换为分
	amountInFen := decimal.NewFromFloat(r.Amount).Mul(decimal.NewFromInt(100)).IntPart()
	amountStr := strconv.FormatInt(amountInFen, 10)

	send := &PayReq{
		MerchantNo: merchantNo,
		TxnID:      r.OrderID,
		Amount:     amountStr,
		Type:       1,
		Callback:   s.Base.GetPayCallbackURL(),
		Currency:   "INR",
	}
	send.Sign = s.Base.SignMD5(send)
	return send
}

// PackWithdrawReq 打包提现请求数据
func (s *Sub) PackWithdrawReq() interface{} {
	r := s.Base.WithdrawReq
	merchantNo, _ := getWandaConfig()

	if r.PayType != common.WithdrawTypeBank && r.PayType != common.WithdrawTypeUPI {
		return nil
	}

	// 金额转换为分
	amountInFen := decimal.NewFromFloat(r.Amount).Mul(decimal.NewFromInt(100)).IntPart()
	amountStr := strconv.FormatInt(amountInFen, 10)

	send := &WithdrawReq{
		MerchantNo: merchantNo,
		TxnID:      r.OrderID,
		Amount:     amountStr,
		Type:       "1",
		IFSC:       r.IFSC,
		CardNum:    r.Account,
		Name:       r.Name,
		Email:      r.Email,
		Mobile:     r.Mobile,
		Callback:   s.Base.GetWithdrawCallbackURL(),
		Currency:   "INR",
	}
	send.Sign = s.Base.SignMD5(send)
	return send
}

// PackQueryWithdrawReq 打包查询提现请求数据
func (s *Sub) PackQueryWithdrawReq() interface{} {
	r := s.Base.QueryWithdrawReq
	merchantNo, _ := getWandaConfig()

	send := &QueryWithdrawReq{
		MerchantNo: merchantNo,
		TxnID:      r.OrderID,
	}
	send.Sign = s.Base.SignMD5(send)
	return send
}

// PackQueryPayReq 打包查询支付请求数据
func (s *Sub) PackQueryPayReq() interface{} {
	r := s.Base.QueryPayReq
	merchantNo, _ := getWandaConfig()

	send := &QueryPayReq{
		MerchantNo: merchantNo,
		TxnID:      r.OrderID,
	}
	send.Sign = s.Base.SignMD5(send)
	return send
}

// CheckSign 验证回调签名
func (s *Sub) CheckSign(str string) bool {
	s.Base.CallbackResp.Msg = "SUCCESS"

	log.Info("开始验证")

	if s.Base.Opt == base.OPTPayCB {
		req := s.Base.CallbackReq.(*PayCallbackReq)
		log.Debug("checkSign pay:%+v", *req)

		s.Base.CallbackResp.OrderID = req.TxnID
		s.Base.CallbackResp.Success = req.Status == "SUCCESS"
		s.Base.CallbackResp.APIOrderID = req.ID

		// 验证签名
		mySign := s.Base.SignMD5(req)
		return mySign == req.Sign
	} else if s.Base.Opt == base.OPTWithdrawCB {
		req := s.Base.CallbackReq.(*WithdrawCallbackReq)
		log.Debug("checkSign withdraw:%+v", *req)

		s.Base.CallbackResp.OrderID = req.TxnID
		s.Base.CallbackResp.Success = req.Status == "SUCCESS" || req.Status == "PENDING"
		s.Base.CallbackResp.APIOrderID = req.ID
		s.Base.CallbackResp.FailMessage = req.Message
		s.Base.CallbackResp.Utr = req.UTR

		// 验证签名
		mySign := s.Base.SignMD5(req)
		return mySign == req.Sign
	}
	return false
}
