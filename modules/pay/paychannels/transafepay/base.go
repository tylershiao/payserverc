package transafepay

import (
	"errors"
	"net/http"
	"server/common"
	"server/modules/pay/base"
	"server/modules/pay/values"
	"server/sign"
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
	b.HttpType = base.HttpTypeJson
	// b.ShouldSignUpper = true
	b.CallbackResp.Msg = "SUCCESS"
	b.SignPassStr = []string{"sign"}
	// b.KeyName = "secretKey"
	// b.PassKeyName = true
	if b.Opt == base.OPTPay {
		b.Resp = new(PayResp)
		b.ReqURL = baseURL + payURL
	} else if b.Opt == base.OPTWithdraw {
		b.Resp = new(WithdrawResp)
		b.ReqURL = baseURL + withdrawURL
	} else if b.Opt == base.OPTPayCB {
		b.CallbackReq = new(PayCallbackReq)
	} else if b.Opt == base.OPTWithdrawCB {
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

// Sub 支付模块的子结构体
type Sub struct {
	Base *base.Base
}

// PackHeader 用于打包请求头信息
// 当前实现为空，可根据需要扩展
func (s *Sub) PackHeader(header http.Header) {
	header.Set("lang", "en")
}

// PackReq 根据操作类型打包请求数据
// 根据不同的操作类型（支付、提现、查询等）返回对应的请求结构体
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

// GetResp 处理响应数据
// 根据不同的操作类型处理响应结果，设置相应的返回值
func (s *Sub) GetResp() error {
	if s.Base.Opt == 1 {
		resp := s.Base.Resp.(*PayResp)
		if resp.Status != 200 || resp.Data.Link == "" {
			return errors.New("pay fail")
		}
		s.Base.PayResp = &common.PayResp{APIOrderID: resp.Data.PlatOrderNo, URL: resp.Data.Link, Channel: int(s.Base.Channel)}
		return nil
	} else if s.Base.Opt == base.OPTQueryPay {
		resp := s.Base.Resp.(*QueryPayResp)
		log.Debug("QueryPayResp:%+v", resp)
		s.Base.QueryPayResp.OrderID = resp.Data.MchOrderNo
		s.Base.QueryPayResp.APIOrderID = resp.Data.PlatOrderNo
		if s.Base.Status == 0 {
			if resp.Status == 200 && resp.Data.OrderStatus == "SUCCESS" {
				s.Base.QueryPayResp.Status = base.QueryStatusOrderSuccess
			}
		}
		if s.Base.QueryPayResp.Status == 0 {
			s.Base.QueryPayResp.Status = base.QueryStatusUnknown
		}
		return nil
	} else if s.Base.Opt == base.OPTWithdraw {
		resp := s.Base.Resp.(*WithdrawResp)
		if s.Base.Status == 0 && resp.Status != 200 {
			return errors.New("withdraw fail")
		}
		s.Base.WithdrawResp = &common.WithdrawResp{APIOrderID: resp.Data.PlatOrderNo, Channel: int(s.Base.Channel)}
		return nil
	} else if s.Base.Opt == base.OPTQueryWithdraw {
		resp := s.Base.Resp.(*QueryWithdrawResp)
		log.Debug("QueryWithdrawResp:%+v", resp)
		s.Base.QueryWithdrawResp.OrderID = resp.Data.MchOrderNo
		s.Base.QueryWithdrawResp.APIOrderID = resp.Data.PlatOrderNo
		if s.Base.Status == 0 {
			if resp.Status == 200 {
				if resp.Data.OrderStatus == "SUCCESS" {
					s.Base.QueryWithdrawResp.Status = base.QueryStatusOrderSuccess
				} else if resp.Data.OrderStatus == "PENDING" {
					s.Base.QueryWithdrawResp.Status = base.QueryStatusSuccess
				}
			} else {
				if resp.Data.OrderStatus == "FAILED" || resp.Data.OrderStatus == "REFUND" {
					s.Base.QueryWithdrawResp.Status = base.QueryStatusFail
				}
			}
		}
		if s.Base.QueryWithdrawResp.Status == 0 {
			s.Base.QueryWithdrawResp.Status = base.QueryStatusUnknown
		}
		return nil
	} else {
		return nil
	}
}

// PackPayReq 打包支付请求数据
// 构造支付请求结构体，包括商户ID、订单号、支付类型、通知地址等信息，并生成签名
func (s *Sub) PackPayReq() interface{} {
	r := s.Base.PayReq
	send := &PayReq{
		MchID:       mid,
		TxChannel:   "TX_INDIA_001",
		AppID:       appID,
		Timestamp:   time.Now().Unix(),
		MchOrderNo:  r.OrderID,
		BankCode:    "UPI",
		Amount:      r.Amount,
		Name:        r.Name,
		Phone:       r.Mobile,
		Email:       r.Email,
		ProductInfo: "shop buy",
		NotifyUrl:   s.Base.GetPayCallbackURL(),
		ReturnUrl:   values.GetFrontCallback(),
	}
	str := sign.GetSignStr(send, s.Base.SignPassStr...)
	var err error
	send.Sign, err = base.SignWithSHA256([]byte(str), privateKey)
	if err != nil {
		log.Error("err:%v", err)
	}
	return send
}

// PackWithdrawReq 打包提现请求数据
// 构造提现请求结构体，根据支付类型（银行、UPI）设置不同的字段，并生成签名
func (s *Sub) PackWithdrawReq() interface{} {
	r := s.Base.WithdrawReq
	// 这里支付可能存在问题，但是需要使用tkpay所以先关掉测试
	//if r.PayType != common.WithdrawTypeBank {
	//	return nil
	//}
	send := &WithdrawReq{
		MchID:      mid,
		TxChannel:  "TX_INDIA_001",
		AppID:      appID,
		Timestamp:  time.Now().Unix(),
		MchOrderNo: r.OrderID,
		Name:       r.Name,
		Phone:      r.Mobile,
		Email:      r.Email,
		BankCode:   "BANK_IN",
		// Account: ,
		Amount:    r.Amount,
		NotifyUrl: s.Base.GetWithdrawCallbackURL(),
		// IFSC: ,
	}

	if r.PayType == common.WithdrawTypeBank {
		send.Account = r.Account
		send.IFSC = r.IFSC
	} else {
		return nil
	}
	str := sign.GetSignStr(send, s.Base.SignPassStr...)
	var err error
	send.Sign, err = base.SignWithSHA256([]byte(str), privateKey)
	if err != nil {
		log.Error("err:%v", err)
	}
	return send
}

// PackQueryWithdrawReq 打包查询提现请求数据
// 构造查询提现请求结构体，包括商户ID和订单号，并生成签名
func (s *Sub) PackQueryWithdrawReq() interface{} {
	r := s.Base.QueryWithdrawReq
	send := &QueryWithdrawReq{
		MchID:      mid,
		Timestamp:  time.Now().Unix(),
		MchOrderNo: r.OrderID,
		TrxID:      r.APIOrderID,
	}
	str := sign.GetSignStr(send, s.Base.SignPassStr...)
	var err error
	send.Sign, err = base.SignWithSHA256([]byte(str), privateKey)
	if err != nil {
		log.Error("err:%v", err)
	}
	return send
}

// PackQueryPayReq 打包查询支付请求数据
// 构造查询支付请求结构体，包括商户ID和订单号，并生成签名
func (s *Sub) PackQueryPayReq() interface{} {
	r := s.Base.QueryPayReq
	send := &QueryPayReq{
		MchID:      mid,
		Timestamp:  time.Now().Unix(),
		MchOrderNo: r.OrderID,
		TrxID:      r.APIOrderID,
	}
	str := sign.GetSignStr(send, s.Base.SignPassStr...)
	var err error
	send.Sign, err = base.SignWithSHA256([]byte(str), privateKey)
	if err != nil {
		log.Error("err:%v", err)
	}
	return send
}

// CheckSign 验证签名
// 根据回调类型验证签名是否正确，设置回调响应信息
func (s *Sub) CheckSign(str string) bool {
	checkSign := ""
	if s.Base.Opt == base.OPTPayCB {
		req := s.Base.CallbackReq.(*PayCallbackReq)
		log.Debug("checkSign pay:%+v", *req)
		checkSign = req.Sign
		s.Base.CallbackResp.OrderID = req.MchOrderNo
		s.Base.CallbackResp.Success = req.OrderStatus == "SUCCESS"
	} else if s.Base.Opt == base.OPTWithdrawCB {
		req := s.Base.CallbackReq.(*WithdrawCallbackReq)
		log.Debug("checkSign withdraw:%+v", *req)
		checkSign = req.Sign
		s.Base.CallbackResp.OrderID = req.MchOrderNo
		s.Base.CallbackResp.Success = req.OrderStatus == "SUCCESS"
		s.Base.CallbackResp.Utr = req.UTR
	}
	checkStr, _ := base.GetSignStringFromJSON(str, s.Base.SignPassStr...)
	pass, err := base.VerifyWithSHA256([]byte(checkStr), publicKey, checkSign)
	if err != nil {
		log.Error("err:%v", err)
	}
	return pass
}
