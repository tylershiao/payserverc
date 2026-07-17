package propay

import (
	"fmt"
	"net/http"
	"server/common"
	"server/modules/pay/base"
	"server/modules/pay/values"
	"server/sign"
	"server/util"
	"strings"

	"github.com/liangdas/mqant/log"
)

// NewSub 初始化ProPay渠道
// 根据不同的操作类型设置对应的响应结构体和请求URL
func NewSub(b *base.Base) {
	sub := &Sub{
		Base: b,
	}
	b.SignKey = key
	b.HttpType = base.HttpTypeJson
	b.ShouldSignUpper = true

	switch b.Opt {
	case base.OPTPay:
		b.Resp = new(PayResp)
		b.ReqURL = baseURL + payURL
	case base.OPTWithdraw:
		b.Resp = new(WithdrawResp)
		b.ReqURL = baseURL + withdrawURL
	case base.OPTPayCB:
		b.CallbackResp.Msg = "success"
		b.CallbackReq = new(PayCallbackReq)
		b.HttpType = base.HttpTypeForm
		b.PostformType = base.PostformTypeBasic
	case base.OPTWithdrawCB:
		b.CallbackResp.Msg = "success"
		b.CallbackReq = new(WithdrawCallbackReq)
		b.HttpType = base.HttpTypeForm
		b.PostformType = base.PostformTypeBasic
	case base.OPTQueryWithdraw:
		b.Resp = new(QueryWithdrawResp)
		b.ReqURL = baseURL + queryWithdrawURL
	case base.OPTQueryPay:
		b.Resp = new(QueryPayResp)
		b.ReqURL = baseURL + queryPayURL
	}
	b.Sub = sub
}

// Sub 是ProPay渠道的实现结构体
type Sub struct {
	Base *base.Base
}

// PackHeader 为请求头添加必要的信息
// ProPay使用application/json作为内容类型
func (s *Sub) PackHeader(header http.Header) {
	// ProPay使用application/json
}

// PackReq 根据操作类型打包请求数据
// 根据不同的操作类型返回对应的请求结构体
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

// GetResp 处理响应数据并设置返回结果
// 根据不同的操作类型解析响应并设置相应的返回值
func (s *Sub) GetResp() error {
	switch s.Base.Opt {
	case base.OPTPay:
		resp := s.Base.Resp.(*PayResp)
		if resp.Code != 0 {
			return fmt.Errorf("pay failed: %s", resp.Msg)
		}
		s.Base.PayResp = &common.PayResp{
			APIOrderID: resp.Data.PayOrderId,
			URL:        resp.Data.PayData,
			Channel:    int(s.Base.Channel),
		}
		return nil

	case base.OPTWithdraw:
		resp := s.Base.Resp.(*WithdrawResp)
		if resp.Code != 0 {
			return fmt.Errorf("withdraw failed: %s", resp.Msg)
		}
		s.Base.WithdrawResp = &common.WithdrawResp{
			APIOrderID: resp.Data.TransferID,
			Channel:    int(s.Base.Channel),
		}
		return nil

	case base.OPTQueryWithdraw:
		resp := s.Base.Resp.(*QueryWithdrawResp)
		s.Base.QueryWithdrawResp.OrderID = resp.Data.MchOrderNo
		s.Base.QueryWithdrawResp.APIOrderID = resp.Data.TransferID

		if resp.Code == 0 {
			switch resp.Data.State {
			case 2: // 转账成功
				s.Base.QueryWithdrawResp.Status = base.QueryStatusOrderSuccess
			case 3: // 转账失败
				s.Base.QueryWithdrawResp.Status = base.QueryStatusFail
			case 1: // 转账中
				s.Base.QueryWithdrawResp.Status = base.QueryStatusSuccess
			default:
				s.Base.QueryWithdrawResp.Status = base.QueryStatusUnknown
			}
		} else {
			s.Base.QueryWithdrawResp.Status = base.QueryStatusUnknown
		}
		return nil

	case base.OPTQueryPay:
		resp := s.Base.Resp.(*QueryPayResp)
		s.Base.QueryPayResp.OrderID = resp.Data.MchOrderNo
		s.Base.QueryPayResp.APIOrderID = resp.Data.PayOrderId

		if resp.Code == 0 && resp.Data.State == 2 { // 支付成功
			s.Base.QueryPayResp.Status = base.QueryStatusOrderSuccess
		} else if resp.Code == 0 && resp.Data.State == 1 { // 支付中
			s.Base.QueryPayResp.Status = base.QueryStatusSuccess
		} else if resp.Data.State == 3 { // 支付失败
			s.Base.QueryPayResp.Status = base.QueryStatusFail
		} else {
			s.Base.QueryPayResp.Status = base.QueryStatusUnknown
		}
		return nil
	}

	return nil
}

// PackPayReq 打包支付请求数据
// 构造支付请求参数并生成签名
func (s *Sub) PackPayReq() interface{} {
	r := s.Base.PayReq
	send := &PayReq{
		MchNo:         mchNo,
		AppId:         appId,
		MchOrderNo:    r.OrderID,
		Amount:        int64(r.Amount * 100), // 转换为分
		CustomerName:  r.Name,
		CustomerEmail: r.Email,
		CustomerPhone: r.Mobile,
		NotifyUrl:     s.Base.GetPayCallbackURL(),
		ReturnUrl:     values.GetFrontCallback(),
		//NotifyUrl: "https://your-test-domain.com/callback/pay/3",
		//ReturnUrl: "https://your-test-domain.com/common/front/payCallback",
	}

	// 生成签名
	signStr := s.getSignStr(send)
	send.Sign = s.calculateSign(signStr)

	return send
}

// PackWithdrawReq 打包提现请求数据
// 构造提现请求参数并生成签名
func (s *Sub) PackWithdrawReq() interface{} {
	r := s.Base.WithdrawReq
	if r.PayType != common.WithdrawTypeBank {
		return nil
	}

	send := &WithdrawReq{
		MchNo:        mchNo,
		AppId:        appId,
		MchOrderNo:   r.OrderID,
		Amount:       int64(r.Amount * 100), // 转换为分
		EntryType:    "IMPS",                // 目前只支持IMPS
		AccountNo:    r.Account,
		AccountCode:  r.IFSC,
		AccountName:  r.Name,
		AccountEmail: r.Email,
		AccountPhone: r.Mobile,
		NotifyUrl:    s.Base.GetWithdrawCallbackURL(),
		BankName:     r.BankName,
	}

	// 生成签名
	signStr := s.getSignStr(send)
	send.Sign = s.calculateSign(signStr)

	return send
}

// PackQueryWithdrawReq 打包查询提现请求数据
// 构造查询提现请求参数并生成签名
func (s *Sub) PackQueryWithdrawReq() interface{} {
	r := s.Base.QueryWithdrawReq
	send := &QueryWithdrawReq{
		MchNo:      mchNo,
		AppId:      appId,
		MchOrderNo: r.OrderID,
	}

	// 如果有APIOrderID，使用transferId
	if r.APIOrderID != "" {
		send.TransferID = r.APIOrderID
		send.MchOrderNo = ""
	}

	// 生成签名
	signStr := s.getSignStr(send)
	send.Sign = s.calculateSign(signStr)

	return send
}

// PackQueryPayReq 打包查询支付请求数据
// 构造查询支付请求参数并生成签名
func (s *Sub) PackQueryPayReq() interface{} {
	r := s.Base.QueryPayReq
	send := &QueryPayReq{
		MchNo:      mchNo,
		AppId:      appId,
		MchOrderNo: r.OrderID,
	}

	// 如果有APIOrderID，使用payOrderId
	if r.APIOrderID != "" {
		send.PayOrderId = r.APIOrderID
		send.MchOrderNo = ""
	}

	// 生成签名
	signStr := s.getSignStr(send)
	send.Sign = s.calculateSign(signStr)

	return send
}

// CheckSign 验证回调签名
// 根据不同的回调类型验证签名并设置回调响应
func (s *Sub) CheckSign(str string) bool {
	s.Base.CallbackResp.Msg = "success"

	switch s.Base.Opt {
	case base.OPTPayCB:
		req := s.Base.CallbackReq.(*PayCallbackReq)
		log.Debug("pay callback: %+v", *req)

		// 验证签名
		mySign := s.calculateSign(s.getCallbackSignStr(req))
		if mySign != req.Sign {
			log.Error("sign check failed, expected: %s, got: %s", mySign, req.Sign)
			return false
		}

		// 设置回调响应
		s.Base.CallbackResp.OrderID = req.MchOrderNo
		s.Base.CallbackResp.APIOrderID = req.PayOrderId
		s.Base.CallbackResp.Success = req.State == 2 // 2表示支付成功
		s.Base.CallbackResp.FailMessage = req.ErrMsg

		return true

	case base.OPTWithdrawCB:
		req := s.Base.CallbackReq.(*WithdrawCallbackReq)
		log.Debug("withdraw callback: %+v", *req)

		// 验证签名
		mySign := s.calculateSign(s.getCallbackSignStr(req))
		if mySign != req.Sign {
			log.Error("sign check failed, expected: %s, got: %s", mySign, req.Sign)
			return false
		}

		// 设置回调响应
		s.Base.CallbackResp.OrderID = req.MchOrderNo
		s.Base.CallbackResp.APIOrderID = req.TransferID
		s.Base.CallbackResp.Success = req.State == 2 // 2表示转账成功
		s.Base.CallbackResp.FailMessage = req.ErrMsg
		s.Base.CallbackResp.Utr = req.ChannelExtra

		return true
	}

	return false
}

// getSignStr 获取签名字符串
// 按ASCII码排序并排除空值
func (s *Sub) getSignStr(data interface{}) string {
	// ProPay要求按ASCII码排序并排除空值
	signStr := sign.GetSignStr(data, "sign")
	return signStr
}

// getCallbackSignStr 获取回调签名字符串
// 回调签名不包含sign字段
func (s *Sub) getCallbackSignStr(data interface{}) string {
	// 回调签名不包含sign字段
	signStr := sign.GetSignStr(data, "sign")
	return signStr
}

// calculateSign 计算MD5签名
// ProPay签名规则：stringA + "&key=" + 私钥，然后MD5并转大写
func (s *Sub) calculateSign(signStr string) string {
	// ProPay签名规则：stringA + "&key=" + 私钥，然后MD5并转大写
	finalStr := signStr + "&key=" + key
	log.Debug("final sign string: %s", finalStr)

	md5Result := util.CalculateMD5(finalStr)
	return strings.ToUpper(md5Result)
}
