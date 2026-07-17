package wkpay

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"server/common"
	"server/modules/pay/base"
	"server/modules/pay/values"
	"server/util"
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
	b.HttpType = base.HttpTypeForm
	b.ShouldSignUpper = true
	b.CallbackResp.Msg = "OK"
	b.SignPassStr = []string{"pay_md5sign"}
	if b.Opt == base.OPTPay {
		b.Resp = new(PayResp)
		b.SignPassStr = []string{"pay_md5sign", "email", "pay_attach", "IP", "phone", "pay_productname",
			"pay_productnum", "pay_productdesc", "pay_producturl", "return_json", "firstName", "lastName"}
		b.ReqURL = baseURL + payURL
	} else if b.Opt == base.OPTWithdraw {
		b.SignPassStr = []string{"pay_md5sign"}
		b.Resp = new(WithdrawResp)
		b.ReqURL = baseURL + withdrawURL
	} else if b.Opt == base.OPTPayCB {
		b.SignPassStr = []string{"sign", "attach"}
		b.CallbackReq = new(PayCallbackReq)
	} else if b.Opt == base.OPTWithdrawCB {
		b.SignPassStr = []string{"pay_md5sign"}
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
		if resp.Code != 0 || resp.Data.PayURL == "" {
			return errors.New("pay fail")
		}
		s.Base.PayResp = &common.PayResp{URL: resp.Data.PayURL, Channel: int(s.Base.Channel)}
		return nil
	} else if s.Base.Opt == base.OPTWithdraw {
		resp := s.Base.Resp.(*WithdrawResp)
		if s.Base.Status == 0 && resp.Status != "success" {
			return errors.New("withdraw fail")
		}
		s.Base.WithdrawResp = &common.WithdrawResp{APIOrderID: resp.TransactionID, Channel: int(s.Base.Channel), Message: resp.Msg}
		return nil
	} else if s.Base.Opt == base.OPTQueryWithdraw {
		resp := s.Base.Resp.(*QueryWithdrawResp)
		log.Debug("QueryWithdrawResp:%+v", resp)
		s.Base.QueryWithdrawResp.OrderID = resp.OutTradeNo
		s.Base.QueryWithdrawResp.APIOrderID = resp.TransactionID
		if s.Base.Status == 0 && resp.Status == "success" {
			// if resp.Status == "success" {
			if resp.RefCode == "1" {
				s.Base.QueryWithdrawResp.Status = base.QueryStatusOrderSuccess
			} else if resp.RefCode == "3" || resp.RefCode == "4" || resp.RefCode == "6" {
				s.Base.QueryWithdrawResp.Status = base.QueryStatusSuccess
			} else if resp.Status == "2" {
				// } else {
				s.Base.QueryWithdrawResp.Status = base.QueryStatusFail
			}
			// }
		}
		if s.Base.QueryWithdrawResp.Status == 0 {
			s.Base.QueryWithdrawResp.Status = base.QueryStatusUnknown
		}
		return nil
	} else if s.Base.Opt == base.OPTQueryPay {
		resp := s.Base.Resp.(*QueryPayResp)
		log.Debug("QueryPayResp:%+v", resp)
		s.Base.QueryPayResp.OrderID = resp.OrderID
		s.Base.QueryPayResp.APIOrderID = resp.TradeState
		if s.Base.Status == 0 {
			if resp.ReturnCode == "00" && resp.TradeState == "SUCCESS" {
				s.Base.QueryPayResp.Status = base.QueryStatusOrderSuccess
			}
		}
		if s.Base.QueryPayResp.Status == 0 {
			s.Base.QueryPayResp.Status = base.QueryStatusUnknown
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
		PayMemberID:    mid,
		PayOrderID:     r.OrderID,
		PayApplyDate:   time.Now().Format("2006-01-02 15:04:05"),
		PayBankCode:    "953",
		PayNotifyURL:   s.Base.GetPayCallbackURL(),
		PayCallbackURL: values.GetFrontCallback(),
		PayAmount:      fmt.Sprintf("%.2f", r.Amount),
		Email:          r.Email,
		IP:             r.IP,
		Phone:          r.Mobile,
		PayProductName: "shop",
		FirstName:      r.Name,
		LastName:       r.Name,
		ReturnJSON:     "1",
	}
	send.PayMd5Sign = s.Base.SignMD5(send)
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
		OutTradeNo: r.OrderID,
		Money:      fmt.Sprintf("%d", int(r.Amount)),
		// BankName: ,
		// SubBranch: ,
		AccountName: r.Name,
		// CardNumber: ,
		Province: "Maharashtra",
		City:     "Mumbai",
		// Extends: ,
		NotifyURL: s.Base.GetWithdrawCallbackURL(),
	}

	if r.PayType == common.WithdrawTypeBank {
		send.CardNumber = r.Account
		send.BankName = r.BankName
		send.SubBranch = r.BankName
		extData := ExtendsData{
			Phone: r.Mobile,
			IFSC:  r.IFSC,
		}
		byt, _ := json.Marshal(extData)
		fmt.Println(string(byt))
		send.Extends = base64.StdEncoding.EncodeToString(byt)
	} else {
		return nil
	}

	send.PayMd5Sign = s.Base.SignMD5(send)
	return send
}

// PackQueryWithdrawReq 打包查询提现请求数据
// 构造查询提现请求结构体，包括商户ID和订单号，并生成签名
func (s *Sub) PackQueryWithdrawReq() interface{} {
	r := s.Base.QueryWithdrawReq
	send := &QueryWithdrawReq{
		MchID:      mid,
		OutTradeNo: r.OrderID,
	}
	send.PayMd5Sign = s.Base.SignMD5(send)
	return send
}

// PackQueryPayReq 打包查询支付请求数据
// 构造查询支付请求结构体，包括商户ID和订单号，并生成签名
func (s *Sub) PackQueryPayReq() interface{} {
	r := s.Base.QueryPayReq
	send := &QueryPayReq{
		PayMemberID: mid,
		PayOrderID:  r.OrderID,
	}
	send.PayMd5Sign = s.Base.SignMD5(send)
	return send
}

// CheckSign 验证签名
// 根据回调类型验证签名是否正确，设置回调响应信息
func (s *Sub) CheckSign(str string) bool {
	str += "&key=" + key
	checkSign := ""
	mySign := strings.ToUpper(util.CalculateMD5(str))
	if s.Base.Opt == base.OPTPayCB {
		req := s.Base.CallbackReq.(*PayCallbackReq)
		log.Debug("checkSign pay:%+v", *req)
		checkSign = req.Sign
		s.Base.CallbackResp.OrderID = req.OrderID
		s.Base.CallbackResp.APIOrderID = req.TransactionID
		s.Base.CallbackResp.Success = req.ReturnCode == "00"
		// mySign = s.Base.SignMD5(req)
	} else if s.Base.Opt == base.OPTWithdrawCB {
		req := s.Base.CallbackReq.(*WithdrawCallbackReq)
		log.Debug("checkSign withdraw:%+v", *req)
		checkSign = req.PayMd5Sign
		s.Base.CallbackResp.OrderID = req.OutTradeNo
		s.Base.CallbackResp.Success = req.RefCode == "1"
		s.Base.CallbackResp.APIOrderID = req.TransactionID
		s.Base.CallbackResp.FailMessage = req.RefMsg
		s.Base.CallbackResp.Utr = req.UTR
		// mySign = s.Base.SignMD5(req)
	}
	return mySign == checkSign
}
