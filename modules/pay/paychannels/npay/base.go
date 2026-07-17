package npay

import (
	"errors"
	"net/http"
	"server/common"
	"server/modules/pay/base"
	"server/util"
	"strings"

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
	b.CallbackResp.Msg = "SUCCESS"
	// b.SignPassStr = []string{"pay_md5sign"}
	switch b.Opt {
	case base.OPTWithdraw:
		// b.SignPassStr = []string{}
		b.Resp = new(WithdrawResp)
		b.ReqURL = baseURL + withdrawURL
	case base.OPTWithdrawCB:
		b.CallbackReq = new(WithdrawCallbackReq)
	case base.OPTQueryWithdraw:
		b.Resp = new(QueryWithdrawResp)
		b.ReqURL = baseURL + queryWithdrawURL
	default:
		return
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
	// case base.OPTPay:
	// 	return s.PackPayReq()
	case base.OPTWithdraw:
		return s.PackWithdrawReq()
	case base.OPTQueryWithdraw:
		return s.PackQueryWithdrawReq()
		// case base.OPTQueryPay:
		// return s.PackQueryPayReq()
	}
	return nil
}

// GetResp 处理响应数据
// 根据不同的操作类型处理响应结果，设置相应的返回值
func (s *Sub) GetResp() error {
	switch s.Base.Opt {
	case base.OPTWithdraw:
		resp := s.Base.Resp.(*WithdrawResp)
		if s.Base.Status == 0 && resp.Status != "2" {
			return errors.New("withdraw fail")
		}
		s.Base.WithdrawResp = &common.WithdrawResp{APIOrderID: resp.SerialID, Channel: int(s.Base.Channel), Message: resp.Message}
		return nil
	case base.OPTQueryWithdraw:
		resp := s.Base.Resp.(*QueryWithdrawResp)
		log.Debug("QueryWithdrawResp:%+v", resp)
		s.Base.QueryWithdrawResp.OrderID = resp.WithdrawID
		s.Base.QueryWithdrawResp.APIOrderID = resp.SerialID
		if s.Base.Status == 0 {
			switch resp.Status {
			case "2":
				s.Base.QueryWithdrawResp.Status = base.QueryStatusOrderSuccess
			case "3":
				s.Base.QueryWithdrawResp.Status = base.QueryStatusSuccess
			case "0":
				s.Base.QueryWithdrawResp.Status = base.QueryStatusFail
			}
		}
		if s.Base.QueryWithdrawResp.Status == 0 {
			s.Base.QueryWithdrawResp.Status = base.QueryStatusUnknown
		}
		return nil
	default:
		return errors.New("unknown opt")
	}
}

// PackPayReq 打包支付请求数据
// 构造支付请求结构体，包括商户ID、订单号、支付类型、通知地址等信息，并生成签名
func (s *Sub) PackPayReq() interface{} {
	// r := s.Base.PayReq
	// send := &PayReq{
	// 	PayMemberID:    mid,
	// 	PayOrderID:     r.OrderID,
	// 	PayApplyDate:   time.Now().Format("2006-01-02 15:04:05"),
	// 	PayBankCode:    "953",
	// 	PayNotifyURL:   s.Base.GetPayCallbackURL(),
	// 	PayCallbackURL: values.GetFrontCallback(),
	// 	PayAmount:      fmt.Sprintf("%.2f", r.Amount),
	// 	Email:          r.Email,
	// 	IP:             r.IP,
	// 	Phone:          r.Mobile,
	// 	PayProductName: "shop",
	// 	FirstName:      r.Name,
	// 	LastName:       r.Name,
	// 	ReturnJSON:     "1",
	// }
	// send.PayMd5Sign = s.Base.SignMD5(send)
	// return send
	return nil
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
		Merchant:   mid,
		WithdrawID: r.OrderID,
		DrawAmount: util.RoundFloat(r.Amount, 2),
		// AccountNo: ,
		DrawName: r.Name,
		// DrawBankName: r.BankName,
		// DrawBankID: ,
		// IfscNo: ,
		NoticeURL: s.Base.GetWithdrawCallbackURL(),
		// Md5Info: ,
		// Sign: ,
	}

	if r.PayType == common.WithdrawTypeBank {
		send.AccountNo = r.Account
		send.DrawBankName = "IFSC"
		send.DrawBankID = "ifsc"
		send.IfscNo = r.IFSC
	} else {
		return nil
	}
	md5Str := send.WithdrawID + send.AccountNo + send.DrawName + send.DrawBankID + payoutkey + send.NoticeURL + send.DrawAmount
	send.Md5Info = util.CalculateMD5(md5Str)

	signStr := send.Merchant + send.WithdrawID + send.DrawAmount + payoutkey + send.AccountNo + send.NoticeURL + send.DrawName
	send.Sign = strings.ToUpper(util.CalculateMD5(signStr))
	return send
}

// PackQueryWithdrawReq 打包查询提现请求数据
// 构造查询提现请求结构体，包括商户ID和订单号，并生成签名
func (s *Sub) PackQueryWithdrawReq() interface{} {
	r := s.Base.QueryWithdrawReq
	send := &QueryWithdrawReq{
		Merchant:   mid,
		WithdrawID: r.OrderID,
	}
	signStr := send.Merchant + send.WithdrawID + payoutkey
	send.Sign = strings.ToUpper(util.CalculateMD5(signStr))
	return send
}

// PackQueryPayReq 打包查询支付请求数据
// 构造查询支付请求结构体，包括商户ID和订单号，并生成签名
func (s *Sub) PackQueryPayReq() interface{} {
	// r := s.Base.QueryPayReq
	// send := &QueryPayReq{
	// 	PayMemberID: mid,
	// 	PayOrderID:  r.OrderID,
	// }
	// send.PayMd5Sign = s.Base.SignMD5(send)
	return nil
}

// CheckSign 验证签名
// 根据回调类型验证签名是否正确，设置回调响应信息
func (s *Sub) CheckSign(str string) bool {
	// str += "&key=" + key
	checkSign := ""
	s.Base.CallbackResp.Msg = "success"
	mySign := ""
	switch s.Base.Opt {
	case base.OPTPayCB:
		// req := s.Base.CallbackReq.(*PayCallbackReq)
		// log.Debug("checkSign pay:%+v", *req)
		// checkSign = req.Sign
		// s.Base.CallbackResp.OrderID = req.OrderID
		// s.Base.CallbackResp.APIOrderID = req.TransactionID
		// s.Base.CallbackResp.Success = req.ReturnCode == "00"
		// mySign = s.Base.SignMD5(req)
	case base.OPTWithdrawCB:
		req := s.Base.CallbackReq.(*WithdrawCallbackReq)
		log.Debug("checkSign withdraw:%+v", *req)
		checkSign = req.Sign
		s.Base.CallbackResp.OrderID = req.WithdrawID
		s.Base.CallbackResp.Success = req.Status == "2"
		s.Base.CallbackResp.APIOrderID = req.SerialID
		// s.Base.CallbackResp.FailMessage = req.RefMsg
		s.Base.CallbackResp.Utr = req.UTR

		signStr := req.Merchant + req.WithdrawID + req.SerialID + req.Status + req.Amount + req.RecordDate + req.AccountBankName + req.AccountName + req.AccountNo + payoutkey
		mySign = strings.ToUpper(util.CalculateMD5(signStr))
	}
	return mySign == checkSign
}
