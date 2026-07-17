package worth_pay

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
	} else if b.Opt == base.OPTPayCB {
		b.CallbackReq = new(PayCallbackReq)
	} else if b.Opt == base.OPTWithdrawCB {
		b.CallbackReq = new(WithdrawCallbackReq)
	} else if b.Opt == base.OPTQueryWithdraw {
		b.Resp = new(QueryWithdrawResp)
		b.ReqURL = baseURL + withdrawURL
		b.HttpType = base.HttpTypeGet
	} else if b.Opt == base.OPTQueryPay {
		b.Resp = new(QueryPayResp)
		b.ReqURL = baseURL + payURL
		b.HttpType = base.HttpTypeGet
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
	header.Set("mid", mid)
	header.Set("signature", s.Base.SignMD5(s.Base.PostBody))
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
	log.Debug("worthPay opt code: %d", s.Base.Opt)
	if s.Base.Opt == base.OPTPay {
		resp := s.Base.Resp.(*PayResp)
		log.Debug("worthPay payin resp: %+v", resp)
		if resp.Code != 0 || resp.PaymentLink == "" {
			return errors.New("payin fail")
		}
		s.Base.PayResp = &common.PayResp{APIOrderID: fmt.Sprintf("%v", resp.MerchantOrderId), URL: resp.PaymentLink, Channel: int(s.Base.Channel)}
		return nil
	} else if s.Base.Opt == base.OPTWithdraw {
		resp := s.Base.Resp.(*WithdrawResp)
		log.Debug("worthPay payout resp: %+v", resp)
		if resp.Code != 0 {
			return errors.New("withdraw fail")
		}
		s.Base.WithdrawResp = &common.WithdrawResp{APIOrderID: fmt.Sprintf("%v", resp.MerchantOrderId), Channel: int(s.Base.Channel)}
		return nil
	} else if s.Base.Opt == base.OPTQueryWithdraw {
		resp := s.Base.Resp.(*QueryWithdrawResp)
		log.Debug("QueryWithdrawResp:%+v", resp)
		s.Base.QueryWithdrawResp.OrderID = resp.TxId
		s.Base.QueryWithdrawResp.APIOrderID = resp.MerchantOrderId
		if s.Base.Status == 0 {
			if resp.Code == 0 {
				if resp.OrderStatus == "0" {
					s.Base.QueryWithdrawResp.Status = base.QueryStatusOrderSuccess
				} else if resp.OrderStatus == "3" {
					s.Base.QueryWithdrawResp.Status = base.QueryStatusSuccess
				}
			} else {
				log.Error("payout query resp return non-zero status code: %d, order_id: %s", resp.Code, resp.MerchantOrderId)
				s.Base.QueryWithdrawResp.Status = base.QueryStatusUnknown
			}
		}
		if s.Base.QueryWithdrawResp.Status == 0 {
			s.Base.QueryWithdrawResp.Status = base.QueryStatusUnknown
		}
		return nil
	} else {
		resp := s.Base.Resp.(*QueryPayResp)
		log.Debug("QueryPayResp:%+v", resp)
		s.Base.QueryWithdrawResp.OrderID = resp.TxId
		s.Base.QueryWithdrawResp.APIOrderID = resp.MerchantOrderId
		if s.Base.Status == 0 {
			if resp.Code == 0 && resp.OrderStatus == "0" {
				s.Base.QueryPayResp.Status = base.QueryStatusOrderSuccess
			} else if resp.Code == 0 && resp.OrderStatus == "2" {
				s.Base.QueryPayResp.Status = base.QueryStatusSuccess
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
		CustomerName:    r.Name,
		CustomerPhone:   r.Mobile,
		CustomerEmail:   CheckOrGenerateGmail(r.Email),
		MerchantOrderId: r.OrderID,
		Amount:          int64(r.Amount * 100),
		Purpose:         "games",
		ReturnUrl:       values.GetFrontCallback(),
		NotifyUrl:       s.Base.GetPayCallbackURL(),
		OrderNote:       r.OrderID,
	}
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
		MerchantOrderId: r.OrderID,
		Amount:          int64(r.Amount * 100),
		NotifyUrl:       s.Base.GetWithdrawCallbackURL(),
		BeneName:        r.Name,
		BenePhone:       r.Mobile,
		BeneEmail:       r.Email,
		Remarks:         r.OrderID,
	}

	if r.PayType == common.WithdrawTypeBank {
		send.TransferMode = "BANKTRANSFER"
		send.BeneIfsc = r.IFSC
		send.AccountNumber = r.Account
	} else if r.PayType == common.WithdrawTypeUPI {
		send.TransferMode = "UPI"
		// 生成随机的11位数字字符串，这里是因为支付方需求，所以这里生成了11位随机数字
		randomDigits := fmt.Sprintf("%011d", rand.Intn(10000000000))
		send.BeneIfsc = randomDigits
		send.AccountNumber = r.Account
	}

	return send
}

// PackQueryWithdrawReq 打包查询提现请求数据
// 构造查询提现请求结构体，包括商户ID和订单号，并生成签名
func (s *Sub) PackQueryWithdrawReq() interface{} {
	r := s.Base.QueryWithdrawReq
	send := &QueryWithdrawReq{
		Nonce:           fmt.Sprintf("%011d", rand.Intn(10000000000)),
		MerchantOrderId: r.OrderID,
	}
	return send
}

// PackQueryPayReq 打包查询支付请求数据
// 构造查询支付请求结构体，包括商户ID和订单号，并生成签名
func (s *Sub) PackQueryPayReq() interface{} {
	r := s.Base.QueryPayReq
	send := &QueryPayReq{
		Nonce:           fmt.Sprintf("%011d", rand.Intn(10000000000)),
		MerchantOrderId: r.OrderID,
	}
	return send
}

// CheckSign 验证签名
// 根据回调类型验证签名是否正确，设置回调响应信息
func (s *Sub) CheckSign(str string) bool {
	mySign := ""
	if s.Base.Opt == base.OPTPayCB {
		req := s.Base.CallbackReq.(*PayCallbackReq)
		log.Debug("checkSign pay:%+v", *req)
		s.Base.CallbackResp.OrderID = req.MerchantOrderId
		s.Base.CallbackResp.Success = req.OrderStatus == "0"
		mySign = s.Base.SignMD5(req)
	} else if s.Base.Opt == base.OPTWithdrawCB {
		req := s.Base.CallbackReq.(*WithdrawCallbackReq)
		log.Debug("checkSign withdraw:%+v", *req)
		s.Base.CallbackResp.OrderID = req.MerchantOrderId
		s.Base.CallbackResp.Success = req.OrderStatus == "0"
		mySign = s.Base.SignMD5(req)
	}
	return mySign == s.Base.C.GetHeader("signature")
}
