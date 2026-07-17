package mqr

import (
	"fmt"
	"net/http"
	"server/common"
	"server/modules/pay/base"
	"server/modules/pay/values"
	"strconv"
	"time"

	"github.com/liangdas/mqant/log"
)

// NewSub 初始化MQR渠道
// 根据不同的操作类型设置对应的响应结构体和请求URL
func NewSub(b *base.Base) {
	sub := &Sub{
		Base: b,
	}
	b.SignKey = APISecret
	b.HttpType = base.HttpTypeJson
	b.ShouldSignUpper = false

	switch b.Opt {
	case base.OPTPay:
		b.Resp = new(PayResp)
		b.ReqURL = BaseURL + PayURL
	case base.OPTPayCB:
		b.CallbackResp.Msg = "success"
		b.CallbackReq = new(PayCallbackReq)
		b.HttpType = base.HttpTypeJson
	case base.OPTQueryPay:
		b.Resp = new(QueryPayResp)
		b.ReqURL = BaseURL + QueryPayURL
	// MQR支付方式不支持代付功能，所以不处理 OPTWithdraw 和 OPTWithdrawCB
	case base.OPTWithdraw, base.OPTWithdrawCB, base.OPTQueryWithdraw:
		// 不支持代付功能
		b.Resp = new(PayResp)
		b.Resp.(*PayResp).Code = 1
		b.Resp.(*PayResp).Msg = "MQR payment does not support withdrawal function"
	}
	b.Sub = sub
}

// Sub 是MQR渠道的实现结构体
type Sub struct {
	Base *base.Base
}

// PackHeader 为请求头添加必要的信息
// MQR使用application/json作为内容类型，并添加特定的请求头
func (s *Sub) PackHeader(header http.Header) {
	// 设置基本Content-Type
	header.Set("Content-Type", "application/json")

	// MQR支付需要特定的请求头
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := GenerateRandomString(16)

	header.Set("x-merchant-id", MerchantID)
	header.Set("x-api-key", APIKey)
	header.Set("x-timestamp", timestamp)
	header.Set("x-nonce", nonce)
	header.Set("x-sign-type", SignType)

	// 获取请求体
	body := s.Base.PostBody

	// 保存请求头参数，用于签名计算
	s.Base.ReqData = map[string]interface{}{
		"x-merchant-id": MerchantID,
		"x-api-key":     APIKey,
		"x-timestamp":   timestamp,
		"x-nonce":       nonce,
		"x-sign-type":   SignType,
	}

	// 生成签名
	sign := s.calculateSign(body)

	// 添加到请求头
	header.Set("x-signature", sign)

}

// PackReq 根据操作类型打包请求数据
// 根据不同的操作类型返回对应的请求结构体
func (s *Sub) PackReq() interface{} {
	switch s.Base.Opt {
	case base.OPTPay:
		return s.PackPayReq()
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
			APIOrderID: resp.Data.OrderNo,
			URL:        resp.Data.PayURL,
			Channel:    int(s.Base.Channel),
		}
		return nil

	case base.OPTQueryPay:
		resp := s.Base.Resp.(*QueryPayResp)
		if resp.Code != 200 {
			return fmt.Errorf("query pay failed: %s", resp.Message)
		}

		s.Base.QueryPayResp.OrderID = resp.Data.MerchantOrderNo
		s.Base.QueryPayResp.APIOrderID = resp.Data.OrderNo

		// 根据状态设置查询结果
		switch resp.Data.Status {
		case 2: // 支付成功
			s.Base.QueryPayResp.Status = base.QueryStatusOrderSuccess
		case 3: // 支付失败
			s.Base.QueryPayResp.Status = base.QueryStatusFail
		case 1: // 支付中
			s.Base.QueryPayResp.Status = base.QueryStatusSuccess
		default:
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
		ProjectCode:     MerchantID,
		MerchantOrderNo: r.OrderID,
		UserUID:         r.UID,
		Amount:          strconv.FormatFloat(r.Amount, 'f', 2, 64),
		Currency:        Currency,
		PaymentMode:     PaymentModeMQR,
		NotifyURL:       s.Base.GetPayCallbackURL(),
		ReturnURL:       values.GetFrontCallback(),
		ExtraData:       "",
	}

	return send
}

// PackQueryPayReq 打包查询支付请求数据
// 构造查询支付请求参数并生成签名
func (s *Sub) PackQueryPayReq() interface{} {
	r := s.Base.QueryPayReq
	send := &QueryPayReq{
		OrderNo: r.APIOrderID,
	}

	// 如果没有APIOrderID，使用OrderID
	if r.APIOrderID == "" {
		send.OrderNo = r.OrderID
	}

	return send
}

// CheckSign 验证回调签名
// 根据不同的回调类型验证签名并设置回调响应
func (s *Sub) CheckSign(str string) bool {
	s.Base.CallbackResp.Msg = "success"

	switch s.Base.Opt {
	case base.OPTPayCB:
		req := s.Base.CallbackReq.(*PayCallbackReq)
		log.Debug("MQR pay callback: %+v", *req)

		// 验证签名
		mySign := s.calculateCallbackSign(req)
		if mySign != req.Signature {
			log.Error("sign check failed, expected: %s, got: %s", mySign, req.Signature)
			return false
		}

		// 设置回调响应
		s.Base.CallbackResp.OrderID = req.OrderNo
		s.Base.CallbackResp.APIOrderID = req.MerchantOrderNo
		s.Base.CallbackResp.Success = req.Status == 2 // 2表示支付成功
		s.Base.CallbackResp.FailMessage = ""

		return true
	}

	return false
}

// calculateSign 计算请求签名
// MQR签名规则：将参数按ASCII码排序，然后使用HMAC-SHA256或MD5签名
func (s *Sub) calculateSign(data interface{}) string {

	// 构建签名参数
	params := make(map[string]string)

	// 添加请求头参数
	if reqData, ok := s.Base.ReqData.(map[string]interface{}); ok {
		for k, v := range reqData {
			params[k] = convertToString(v)
		}
	}

	// 添加请求体参数
	switch v := data.(type) {
	case *PayReq:
		params["project_code"] = v.ProjectCode
		params["merchant_order_no"] = v.MerchantOrderNo
		params["user_uid"] = v.UserUID
		params["amount"] = v.Amount
		params["currency"] = v.Currency
		params["payment_mode"] = strconv.Itoa(v.PaymentMode)
		params["notify_url"] = v.NotifyURL
		params["return_url"] = v.ReturnURL
		params["extra_data"] = v.ExtraData
	}

	// 计算签名
	if SignType == SignTypeHMACSHA256 {
		return calculateHMACSHA256Signature(params, APISecret)
	} else {
		return calculateMD5Signature(params, APISecret)
	}
}

// calculateQuerySign 计算查询请求签名
func (s *Sub) calculateQuerySign(data interface{}) string {
	// 构建签名参数
	params := make(map[string]string)

	// 添加请求头参数
	if reqData, ok := s.Base.ReqData.(map[string]interface{}); ok {
		for k, v := range reqData {
			params[k] = convertToString(v)
		}
	}

	// 添加请求体参数
	switch v := data.(type) {
	case *QueryPayReq:
		params["order_no"] = v.OrderNo
	}

	// 计算签名
	if SignType == SignTypeHMACSHA256 {
		return calculateHMACSHA256Signature(params, APISecret)
	} else {
		return calculateMD5Signature(params, APISecret)
	}
}

// calculateCallbackSign 计算回调签名
func (s *Sub) calculateCallbackSign(data interface{}) string {
	// 构建签名参数
	params := make(map[string]string)

	// 添加回调参数
	switch v := data.(type) {
	case *PayCallbackReq:
		if v.OrderNo != "" {
			params["order_no"] = v.OrderNo
		}
		if v.MerchantOrderNo != "" {
			params["merchant_order_no"] = v.MerchantOrderNo
		}
		if v.ProjectCode != "" {
			params["project_code"] = v.ProjectCode
		}
		if v.UserUID != "" {
			params["user_uid"] = v.UserUID
		}
		if v.Amount != "" {
			params["amount"] = v.Amount
		}
		if v.Currency != "" {
			params["currency"] = v.Currency
		}
		if v.Status != 0 {
			params["status"] = strconv.Itoa(v.Status)
		}
		if v.PaymentMode != 0 {
			params["payment_mode"] = strconv.Itoa(v.PaymentMode)
		}
		if v.CreatedAt != 0 {
			params["created_at"] = strconv.FormatInt(v.CreatedAt, 10)
		}
		if v.ExtraData != "" {
			params["extra_data"] = v.ExtraData
		}
	}

	// 计算签名
	if SignType == SignTypeHMACSHA256 {
		return calculateHMACSHA256Signature(params, APISecret)
	} else {
		return calculateMD5Signature(params, APISecret)
	}
}
