package gogopay

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"server/common"
	"server/modules/pay/base"
	"server/util"
	"time"

	"github.com/liangdas/mqant/log"
)

// NewSub 初始化子支付模块
// 根据不同的操作类型（支付、提现、回调等）设置相应的请求URL、响应结构体和签名参数
func NewSub(b *base.Base) {
	sub := &Sub{
		Base: b,
	}
	// b.SignKey = key
	b.HttpType = base.HttpTypeForm
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
		if resp.Code != 0 || resp.Data.URL == "" {
			return errors.New("pay fail")
		}
		s.Base.PayResp = &common.PayResp{APIOrderID: resp.Data.OrderNumber, URL: resp.Data.URL, Channel: int(s.Base.Channel)}
	case base.OPTQueryPay:
		resp := s.Base.Resp.(*QueryPayResp)
		log.Debug("QueryPayResp:%+v", resp)
		s.Base.QueryPayResp.OrderID = resp.Data.OutTradeNo
		s.Base.QueryPayResp.APIOrderID = resp.Data.OrderNumber
		if s.Base.Status == 0 {
			if resp.Code == 0 && resp.Data.Status == "payin_success" {
				s.Base.QueryPayResp.Status = base.QueryStatusOrderSuccess
			}
		}
		if s.Base.QueryPayResp.Status == 0 {
			s.Base.QueryPayResp.Status = base.QueryStatusUnknown
		}
	case base.OPTWithdraw:
		resp := s.Base.Resp.(*WithdrawResp)
		if s.Base.Status == 0 && resp.Code != 0 {
			return errors.New("withdraw fail")
		}
		data := &WidhrawRespData{}
		err := json.Unmarshal(resp.Data, data)
		if err != nil {
			log.Error("err:%v", err)
			return err
		}
		if data.OrderNumber == "" {
			return errors.New("withdraw fail")
		}
		s.Base.WithdrawResp = &common.WithdrawResp{APIOrderID: data.OrderNumber, Channel: int(s.Base.Channel)}
	case base.OPTQueryWithdraw:
		resp := s.Base.Resp.(*QueryWithdrawResp)
		log.Debug("QueryWithdrawResp:%+v", resp)
		s.Base.QueryWithdrawResp.OrderID = resp.Data.OutTradeNo
		s.Base.QueryWithdrawResp.APIOrderID = resp.Data.OrderNumber
		if s.Base.Status == 0 {
			if resp.Code == 0 {
				switch resp.Data.Status {
				case "payout_success":
					s.Base.QueryWithdrawResp.Status = base.QueryStatusOrderSuccess
				case "payout_ing":
					s.Base.QueryWithdrawResp.Status = base.QueryStatusSuccess
				}
			} else {
				if resp.Data.Status == "payout_fail" {
					s.Base.QueryWithdrawResp.Status = base.QueryStatusFail
				}
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
		MerchantID:     mid,
		OutTradeNo:     r.OrderID,
		TotalFee:       fmt.Sprintf("%.2f", r.Amount),
		NotifyURL:      s.Base.GetPayCallbackURL(),
		ReplyType:      "URL",
		Timestamp:      fmt.Sprintf("%d", time.Now().Unix()),
		CustomerName:   r.Name,
		CustomerMobile: r.Mobile,
		CustomerEmail:  r.Email,
	}
	send.Sign = s.Base.SignMD5(send)
	return send
}

// PackWithdrawReq 打包提现请求数据
// 构造提现请求结构体，根据支付类型（银行、UPI）设置不同的字段，并生成签名
func (s *Sub) PackWithdrawReq() interface{} {
	r := s.Base.WithdrawReq
	send := &WithdrawReq{
		MerchantID:  mid,
		OutTradeNo:  r.OrderID,
		TotalFee:    fmt.Sprintf("%.2f", r.Amount),
		NotifyURL:   s.Base.GetWithdrawCallbackURL(),
		Timestamp:   fmt.Sprintf("%d", time.Now().Unix()),
		PaymentMode: "IMPS",
		// AccountNumber: ,
		// IFSC: ,
		CustomerName:   r.Name,
		CustomerMobile: r.Mobile,
		CustomerEmail:  r.Email,
	}

	switch r.PayType {
	case common.WithdrawTypeBank:
		send.AccountNumber = r.Account
		send.IFSC = r.IFSC
	// case common.WithdrawTypeUPI:
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
		MerchantID: mid,
		OutTradeNo: r.OrderID,
		Timestamp:  fmt.Sprintf("%d", time.Now().Unix()),
	}
	send.Sign = s.Base.SignMD5(send)
	return send
}

// PackQueryPayReq 打包查询支付请求数据
// 构造查询支付请求结构体，包括商户ID和订单号，并生成签名
func (s *Sub) PackQueryPayReq() interface{} {
	r := s.Base.QueryPayReq
	send := &QueryPayReq{
		OutTradeNo: r.OrderID,
		Timestamp:  fmt.Sprintf("%d", time.Now().Unix()),
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
		case "payin_success":
			s.Base.CallbackResp.Success = true
		case "payin_fail":
			s.Base.CallbackResp.Success = false
		default:
			s.Base.CallbackStatus = base.CallbackStatusPass
			return false
		}
		checkSign = req.Sign
		s.Base.CallbackResp.OrderID = req.OutTradeNo
	case base.OPTWithdrawCB:
		req := s.Base.CallbackReq.(*WithdrawCallbackReq)
		log.Debug("checkSign withdraw:%+v", *req)
		switch req.Status {
		case "payout_success":
			s.Base.CallbackResp.Success = true
		case "payout_fail":
			s.Base.CallbackResp.Success = false
		default:
			s.Base.CallbackStatus = base.CallbackStatusPass
			return false
		}
		checkSign = req.Sign
		s.Base.CallbackResp.OrderID = req.OutTradeNo
		s.Base.CallbackResp.Utr = req.UTR
	default:
		return false
	}
	mySign := util.CalculateMD5(str + "&" + s.Base.KeyName + "=" + key)
	pass := mySign == checkSign
	if !pass {
		log.Error("check sign fail:mySign:%s,checkSign:%s", mySign, checkSign)
	}
	return pass
}
