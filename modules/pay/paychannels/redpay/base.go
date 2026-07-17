package redpay

import (
	"errors"
	"fmt"
	"net/http"
	"server/common"
	"server/modules/pay/base"
	"server/modules/pay/values"
	"server/sign"
	"strings"
	"time"

	"github.com/liangdas/mqant/log"
)

// NewSub 初始化子支付模块
// 根据不同的操作类型（支付、提现、回调等）设置相应的请求URL、响应结构体和签名参数
func NewSub(b *base.Base) {
	sub := &Sub{
		Base: b,
	}
	b.SignKey = key
	b.HttpType = base.HttpTypeJson
	b.ShouldSignUpper = true
	b.CallbackResp.Msg = "success"
	b.SignPassStr = []string{"sign"}
	if b.Opt == base.OPTPay {
		b.Resp = new(PayResp)
		b.ReqURL = baseURL + payURL
	} else if b.Opt == base.OPTWithdraw {
		b.Resp = new(WithdrawResp)
		b.ReqURL = baseURL + withdrawURL
	} else if b.Opt == base.OPTPayCB {
		b.HttpType = base.HttpTypeGet
		b.CallbackReq = new(PayCallbackReq)
	} else if b.Opt == base.OPTWithdrawCB {
		b.HttpType = base.HttpTypeGet
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
	if s.Base.Opt == base.OPTPay {
		resp := s.Base.Resp.(*PayResp)
		if resp.Code != 200 || resp.Data.PayURL == "" {
			return errors.New("pay fail")
		}
		s.Base.PayResp = &common.PayResp{APIOrderID: resp.Data.PlatformOrderID, URL: resp.Data.PayURL, Channel: int(s.Base.Channel)}
		return nil
	} else if s.Base.Opt == base.OPTWithdraw {
		resp := s.Base.Resp.(*WithdrawResp)
		if s.Base.Status == 0 && resp.Code != 200 {
			return errors.New("withdraw fail")
		}
		s.Base.WithdrawResp = &common.WithdrawResp{APIOrderID: resp.Data.PlatformOrderID, Channel: int(s.Base.Channel)}
		return nil
	} else if s.Base.Opt == base.OPTQueryWithdraw {
		resp := s.Base.Resp.(*QueryWithdrawResp)
		log.Debug("QueryWithdrawResp:%+v", resp)
		s.Base.QueryWithdrawResp.OrderID = resp.Data.MerchantOrderNo
		s.Base.QueryWithdrawResp.APIOrderID = resp.Data.PlatformOrderID
		if s.Base.Status == 0 {
			if resp.Code == 200 {
				if resp.Data.Status == "1" {
					s.Base.QueryWithdrawResp.Status = base.QueryStatusOrderSuccess
				} else if resp.Data.Status == "0" {
					s.Base.QueryWithdrawResp.Status = base.QueryStatusSuccess
				}
			} else {
				if resp.Data.Status == "2" || resp.Data.Status == "5" {
					s.Base.QueryWithdrawResp.Status = base.QueryStatusFail
				}
			}
		}
		if s.Base.QueryWithdrawResp.Status == 0 {
			s.Base.QueryWithdrawResp.Status = base.QueryStatusUnknown
		}
		return nil
	} else if s.Base.Opt == base.OPTQueryPay {
		resp := s.Base.Resp.(*QueryPayResp)
		log.Debug("QueryPayResp:%+v", resp)
		s.Base.QueryPayResp.OrderID = resp.Data.MerchantOrderNo
		s.Base.QueryPayResp.APIOrderID = resp.Data.PlatformOrderID
		if s.Base.Status == 0 {
			if resp.Code == 200 && resp.Data.Status == "1" {
				s.Base.QueryPayResp.Status = base.QueryStatusOrderSuccess
			}
		}
		if s.Base.QueryPayResp.Status == 0 {
			s.Base.QueryPayResp.Status = base.QueryStatusUnknown
		}
		return nil
	} else {
		return errors.New("unknown opt")
	}
}

// PackPayReq 打包支付请求数据
// 构造支付请求结构体，包括商户ID、订单号、支付类型、通知地址等信息，并生成签名
func (s *Sub) PackPayReq() interface{} {
	r := s.Base.PayReq
	send := &PayReq{
		MerchantCid:     mid,
		MerchantOrderNo: r.OrderID,
		Country:         "IN",
		Currency:        "INR",
		Amount:          fmt.Sprintf("%d", int(r.Amount*100)),
		Channel:         "1001",
		Name:            r.Name,
		Mobile:          r.Mobile,
		Email:           r.Email,
		NotifyUrl:       s.Base.GetPayCallbackURL(),
		ReturnUrl:       values.GetFrontCallback(),
		Subject:         "shop",
		Body:            "shop buy",
		Timestamp:       fmt.Sprintf("%d", time.Now().UnixMilli()),
	}
	send.Sign, _ = GenSign(send)
	return send
}

// CheckOrGenerateGmail 输入任意邮箱字符串：
//   - 如果是 Gmail（不区分大小写），原样返回；
//   - 否则返回一个随机生成的 Gmail 地址。
func CheckOrGenerateGmail(email string) string {
	trimmed := strings.TrimSpace(strings.ToLower(email))
	if strings.HasSuffix(trimmed, "@gmail.com") {
		local := strings.Split(trimmed, "@")[0]
		if len(local) > 0 && len(local) <= 64 && len(trimmed) <= 254 {
			return trimmed
		}
	}
	return sign.SafeRandomGmail()
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
		MerchantCid:     mid,
		MerchantOrderNo: r.OrderID,
		Country:         "IN",
		Currency:        "INR",
		Amount:          int(r.Amount * 100),
		Channel:         "1001",
		Name:            r.Name,
		Mobile:          r.Mobile,
		Email:           r.Email,
		// BankCode: ,
		// BankName: ,
		// BankCardNum: ,
		NotifyUrl: s.Base.GetWithdrawCallbackURL(),
		Timestamp: fmt.Sprintf("%d", time.Now().UnixMilli()),
	}

	if r.PayType == common.WithdrawTypeBank {
		send.BankCode = r.IFSC
		send.BankCardNum = r.Account
		send.BankName = r.BankName
	} else {
		return nil
	}
	send.Sign, _ = GenSign(send)
	return send
}

// PackQueryWithdrawReq 打包查询提现请求数据
// 构造查询提现请求结构体，包括商户ID和订单号，并生成签名
func (s *Sub) PackQueryWithdrawReq() interface{} {
	r := s.Base.QueryWithdrawReq
	send := &QueryWithdrawReq{
		MerchantCid:     mid,
		Timestamp:       fmt.Sprintf("%d", time.Now().UnixMilli()),
		MerchantOrderNo: r.OrderID,
	}
	send.Sign, _ = GenSign(send)
	return send
}

// PackQueryPayReq 打包查询支付请求数据
// 构造查询支付请求结构体，包括商户ID和订单号，并生成签名
func (s *Sub) PackQueryPayReq() interface{} {
	r := s.Base.QueryPayReq
	send := &QueryPayReq{
		MerchantCid:     mid,
		Timestamp:       fmt.Sprintf("%d", time.Now().UnixMilli()),
		MerchantOrderNo: r.OrderID,
	}
	send.Sign, _ = GenSign(send)
	return send
}

// CheckSign 验证签名
// 根据回调类型验证签名是否正确，设置回调响应信息
func (s *Sub) CheckSign(str string) bool {
	// str += "&key=" + key
	checkSign := ""
	mySign := ""
	if s.Base.Opt == base.OPTPayCB {
		req := s.Base.CallbackReq.(*PayCallbackReq)
		log.Debug("checkSign pay:%+v", *req)
		checkSign = req.Sign
		s.Base.CallbackResp.OrderID = req.MerchantOrderNo
		switch req.Status {
		case "1":
			s.Base.CallbackResp.Success = true
		case "2":
			s.Base.CallbackResp.Success = false
		default:
			s.Base.CallbackStatus = base.CallbackStatusPass
			return false
		}
		mySign, _ = GenSign(req)
	} else if s.Base.Opt == base.OPTWithdrawCB {
		req := s.Base.CallbackReq.(*WithdrawCallbackReq)
		log.Debug("checkSign withdraw:%+v", *req)
		switch req.Status {
		case "1":
			s.Base.CallbackResp.Success = true
		case "2", "5", "6":
			s.Base.CallbackResp.Success = false
		default:
			s.Base.CallbackStatus = base.CallbackStatusPass
			return false
		}
		checkSign = req.Sign
		s.Base.CallbackResp.OrderID = req.MerchantOrderNo
		s.Base.CallbackResp.Success = req.Status == "1"
		s.Base.CallbackResp.APIOrderID = req.PlatformOrderID
		// s.Base.CallbackResp.FailMessage = req.
		s.Base.CallbackResp.Utr = req.UTR
		mySign, _ = GenSign(req)
	}
	return mySign == checkSign
}
