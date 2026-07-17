package decentpay

import (
	"errors"
	"fmt"
	"net/http"
	"server/common"
	"server/modules/pay/base"
	"server/sign"
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
	// b.SignKey = key
	b.HttpType = base.HttpTypeJson
	b.ShouldSignUpper = true
	b.CallbackResp.Msg = "SUCCESS"
	b.SignPassStr = []string{"sign"}
	b.SignKey = key
	// b.KeyName = "secretKey"
	b.PassKeyName = true
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
		b.HttpType = base.HttpTypeForm
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
	now := fmt.Sprintf("%d", time.Now().Unix())
	header.Set("appId", appID)
	header.Set("timestamp", now)
	header.Set("merchantNo", mid)
	// merchantNo，appId，timestamp依次排在前三个位，与所有入参按照参数按照字段名的 ASCII 码从小到大排序后使用 QueryString 的格式（即key1=value1&key2=value2…）拼接的字符串组合，再加上分配的Access Key Secret。
	// merchantNo=value1&appId=value2&timestamp=value3&key1=value4&key2=value5&key3=value6{Access Key Secret}
	str := sign.GetSignStr(s.Base.PostBody)
	sign := GetSign(mid, appID, now, str)
	header.Set("sign", sign)
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
		if !resp.Success || resp.Data.Target == "" {
			return errors.New("pay fail")
		}
		s.Base.PayResp = &common.PayResp{APIOrderID: resp.Data.OrderNo, URL: resp.Data.Target, Channel: int(s.Base.Channel)}
	case base.OPTQueryPay:
		resp := s.Base.Resp.(*QueryPayResp)
		log.Debug("QueryPayResp:%+v", resp)
		s.Base.QueryPayResp.OrderID = resp.Data.MerchantOrderNo
		s.Base.QueryPayResp.APIOrderID = resp.Data.OrderNo
		if s.Base.Status == 0 {
			if resp.ErrCode == "200" && resp.Data.Status == "received" {
				s.Base.QueryPayResp.Status = base.QueryStatusOrderSuccess
			}
		}
		if s.Base.QueryPayResp.Status == 0 {
			s.Base.QueryPayResp.Status = base.QueryStatusUnknown
		}
	case base.OPTWithdraw:
		resp := s.Base.Resp.(*WithdrawResp)
		if s.Base.Status == 0 && resp.Success != nil && !*resp.Success {
			return errors.New("withdraw fail")
		}
		s.Base.WithdrawResp = &common.WithdrawResp{APIOrderID: resp.Data.OrderNo, Message: resp.ErrMsg, Channel: int(s.Base.Channel)}
	case base.OPTQueryWithdraw:
		resp := s.Base.Resp.(*QueryWithdrawResp)
		log.Debug("QueryWithdrawResp:%+v", resp)
		s.Base.QueryWithdrawResp.OrderID = resp.Data.MerchantOrderNo
		s.Base.QueryWithdrawResp.APIOrderID = resp.Data.OrderNo
		if s.Base.Status == 0 {
			if resp.ErrCode == "200" {
				switch resp.Data.Status {
				case "success":
					s.Base.QueryWithdrawResp.Status = base.QueryStatusOrderSuccess
				case "creating", "processing":
					s.Base.QueryWithdrawResp.Status = base.QueryStatusSuccess
				}
			} else {
				if resp.Data.Status == "fail" {
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
		Amount:          r.Amount,
		NotifyURL:       s.Base.GetPayCallbackURL(),
		MerchantOrderNo: r.OrderID,
		PayProduct:      "india-upi-h5",
		Description:     "shop",
		UserID:          r.UID,
		FullName:        r.Name,
		DateOfReg:       time.Now().Format("2006-01-02 15:04:05"), // 2025-07-08 00:00:00
	}

	return send
}

// PackWithdrawReq 打包提现请求数据
// 构造提现请求结构体，根据支付类型（银行、UPI）设置不同的字段，并生成签名
func (s *Sub) PackWithdrawReq() interface{} {
	r := s.Base.WithdrawReq
	send := &WithdrawReq{
		// AccountNumber: ,
		NotifyURL: s.Base.GetWithdrawCallbackURL(),
		Amount:    r.Amount,
		// BankID: ,
		// BankCode: ,
		CustomerName:    r.Name,
		Description:     "shop out",
		MerchantOrderNo: r.OrderID,
	}

	switch r.PayType {
	case common.WithdrawTypeBank:
		send.AccountNumber = r.Account
		send.BankID = 1
		send.BankCode = r.IFSC
	// case common.WithdrawTypeUPI:
	// 	send.Type = "1"
	// 	send.UPI = r.Account
	default:
		return nil
	}
	return send
}

// PackQueryWithdrawReq 打包查询提现请求数据
// 构造查询提现请求结构体，包括商户ID和订单号，并生成签名
func (s *Sub) PackQueryWithdrawReq() interface{} {
	r := s.Base.QueryWithdrawReq
	send := &QueryWithdrawReq{
		// MerchantID:       mid,
		// MerchantPayoutID: r.OrderID,
	}
	// send.Sign = s.Base.SignMD5(send)
	s.Base.ReqURL += r.OrderID
	return send
}

// PackQueryPayReq 打包查询支付请求数据
// 构造查询支付请求结构体，包括商户ID和订单号，并生成签名
func (s *Sub) PackQueryPayReq() interface{} {
	r := s.Base.QueryPayReq
	send := &QueryPayReq{
		// MerchantID:      mid,
		// MerchantOrderID: r.OrderID,
	}
	// send.Sign = s.Base.SignMD5(send)
	s.Base.ReqURL += r.OrderID
	return send
}

// CheckSign 验证签名
// 根据回调类型验证签名是否正确，设置回调响应信息
func (s *Sub) CheckSign(str string) bool {
	checkSign := s.Base.C.GetHeader("sign")
	switch s.Base.Opt {
	case base.OPTPayCB:
		req := s.Base.CallbackReq.(*PayCallbackReq)
		log.Debug("checkSign pay:%+v", *req)
		// checkSign = req.Sign
		s.Base.CallbackResp.OrderID = req.MerchantOrderNo
		s.Base.CallbackResp.Success = req.Status == "received"
	case base.OPTWithdrawCB:
		req := s.Base.CallbackReq.(*WithdrawCallbackReq)
		log.Debug("checkSign withdraw:%+v", *req)
		// checkSign = req.Sign
		s.Base.CallbackResp.OrderID = req.MerchantOrderNo
		switch req.Status {
		case "success":
			s.Base.CallbackResp.Success = true
		case "fail":
			s.Base.CallbackResp.Success = false
		default:
			s.Base.CallbackStatus = base.CallbackStatusPass
			return false
		}
		// s.Base.CallbackResp.Utr = req
	default:
		return false
	}
	signStr := sign.GetSignStrNull(str)
	mySign := GetSign(s.Base.C.GetHeader("merchantNo"), s.Base.C.GetHeader("appId"), s.Base.C.GetHeader("timestamp"), signStr)
	pass := mySign == checkSign
	if !pass {
		log.Error("check sign fail:mySign:%s,checkSign:%s", mySign, checkSign)
	}
	return pass
}

func GetSign(mid, appId, timeStamp, signStr string) string {
	finalStr := fmt.Sprintf("merchantNo=%s&appId=%s&timestamp=%s&%s%s", mid, appID, timeStamp, signStr, key)
	log.Debug("finalStr:%v", finalStr)
	sign := strings.ToUpper(util.CalculateMD5(finalStr))
	return sign
}
