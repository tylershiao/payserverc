package tkpay

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"server/common"
	"server/modules/pay/base"
	"server/modules/pay/values"
	"server/sign"
	"strings"

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
	if b.Opt == base.OPTPay {
		b.Resp = new(PayResp)
		b.ReqURL = baseURL + payURL
	} else if b.Opt == base.OPTWithdraw {
		b.Resp = new(WithdrawResp)
		b.ReqURL = baseURL + withdrawURL
		b.SignPassStr = []string{"withdraw_type", "notify_url", "phone", "email", "currency", "sign", "pix_type", "identity_card_num"}
	} else if b.Opt == base.OPTPayCB {
		b.SignPassStr = []string{"sign"}
		b.CallbackResp.Msg = "SUCCESS"
		b.CallbackReq = new(PayCallbackReq)
	} else if b.Opt == base.OPTWithdrawCB {
		b.SignPassStr = []string{"sign", "msg"}
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
	if s.Base.Opt == 1 {
		resp := s.Base.Resp.(*PayResp)
		if resp.Code != 200 || resp.Data.PayUrl == "" {
			return errors.New("pay fail")
		}
		s.Base.PayResp = &common.PayResp{APIOrderID: fmt.Sprintf("%v", resp.Data.ID), URL: resp.Data.PayUrl, Channel: int(s.Base.Channel)}
		return nil
	} else if s.Base.Opt == base.OPTWithdraw {
		resp := s.Base.Resp.(*WithdrawResp)
		if s.Base.Status == 0 && resp.Code != 200 {
			return errors.New("withdraw fail")
		}
		s.Base.WithdrawResp = &common.WithdrawResp{APIOrderID: fmt.Sprintf("%v", resp.Data.ID), Channel: int(s.Base.Channel)}
		return nil
	} else if s.Base.Opt == base.OPTQueryWithdraw {
		resp := s.Base.Resp.(*QueryWithdrawResp)
		log.Debug("QueryWithdrawResp:%+v", resp)
		s.Base.QueryWithdrawResp.OrderID = resp.Data.OrderID
		s.Base.QueryWithdrawResp.APIOrderID = resp.Data.ID
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
		s.Base.QueryPayResp.OrderID = resp.Data.OrderID
		s.Base.QueryPayResp.APIOrderID = resp.Data.ID
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

// PackPayReq 打包支付请求数据
// 构造支付请求结构体，包括商户ID、订单号、支付类型、通知地址等信息，并生成签名
func (s *Sub) PackPayReq() interface{} {
	r := s.Base.PayReq
	send := &PayReq{
		MerchantID: mid,
		OrderID:    r.OrderID,
		PayType:    "1",
		NotifyUrl:  s.Base.GetPayCallbackURL(),
		ReturnUrl:  values.GetFrontCallback(),
		Amount:     int64(r.Amount * 100),
		Currency:   "INR",
		Name:       r.Name,
		Phone:      r.Mobile,
		Email:      CheckOrGenerateGmail(r.Email),
	}
	send.Sign = s.Base.SignMD5(send)
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
		MerchantID: mid,
		OrderID:    r.OrderID,
		Amount:     int64(r.Amount * 100),
		// WithdrawType: ,
		NotifyUrl: s.Base.GetWithdrawCallbackURL(),
		Name:      r.Name,
		Phone:     r.Mobile,
		Email:     r.Email,
		// BankCode: ,
		Currency: "INR",
		// Account: ,
	}

	if r.PayType == common.WithdrawTypeBank {
		send.WithdrawType = "1"
		send.BankCode = r.IFSC
		send.Account = r.Account
	} else if r.PayType == common.WithdrawTypeUPI {
		send.WithdrawType = "2"
		// 生成随机的11位数字字符串，这里是因为支付方需求，所以这里生成了11位随机数字
		randomDigits := fmt.Sprintf("%011d", rand.Intn(10000000000))
		send.BankCode = randomDigits
		send.Account = r.Account
	}

	str := sign.GetSignStr(send, s.Base.SignPassStr...) + "&key=" + key
	var err error
	send.Sign, err = base.RsaEncrypt(str, myPublicKey)
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
		MerchantID: mid,
		OrderID:    r.OrderID,
	}
	send.Sign = s.Base.SignMD5(send)
	return send
}

// PackQueryPayReq 打包查询支付请求数据
// 构造查询支付请求结构体，包括商户ID和订单号，并生成签名
func (s *Sub) PackQueryPayReq() interface{} {
	r := s.Base.QueryPayReq
	send := &QueryPayReq{
		MerchantID: mid,
		OrderID:    r.OrderID,
	}
	send.Sign = s.Base.SignMD5(send)
	return send
}

// CheckSign 验证签名
// 根据回调类型验证签名是否正确，设置回调响应信息
func (s *Sub) CheckSign(str string) bool {
	// str += "&key=" + key
	checkSign := ""
	s.Base.CallbackResp.Msg = "success"
	mySign := ""
	if s.Base.Opt == base.OPTPayCB {
		req := s.Base.CallbackReq.(*PayCallbackReq)
		log.Debug("checkSign pay:%+v", *req)
		checkSign = req.Data.Sign
		s.Base.CallbackResp.OrderID = req.Data.OrderID
		s.Base.CallbackResp.Success = req.Code == 200
		mySign = s.Base.SignMD5(req.Data)
	} else if s.Base.Opt == base.OPTWithdrawCB {
		req := s.Base.CallbackReq.(*WithdrawCallbackReq)
		log.Debug("checkSign withdraw:%+v", *req)
		checkSign = req.Data.Sign
		s.Base.CallbackResp.OrderID = req.Data.OrderID
		s.Base.CallbackResp.Success = req.Code == 200
		s.Base.CallbackResp.APIOrderID = req.Data.ID
		s.Base.CallbackResp.FailMessage = req.Message
		s.Base.CallbackResp.Utr = req.Data.UTR
		mySign = s.Base.SignMD5(req.Data)
	}
	return mySign == checkSign
}
