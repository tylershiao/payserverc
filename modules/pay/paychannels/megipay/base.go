package megipay

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"reflect"
	"server/common"
	"server/modules/pay/base"
	"server/modules/pay/values"
	"server/util"
	"sort"
	"strings"

	"github.com/liangdas/mqant/log"
)

func NewSub(b *base.Base) {
	sub := &Sub{Base: b}
	b.HttpType = base.HttpTypeJson
	b.CallbackResp.Msg = "success"
	b.SignPassStr = []string{"sign"}
	b.SignKey = merchantSignKey
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
		b.Resp = new(QueryResp)
		b.ReqURL = baseURL + queryURL
	case base.OPTQueryPay:
		b.Resp = new(QueryResp)
		b.ReqURL = baseURL + queryURL
	}
	b.Sub = sub
}

type Sub struct {
	Base *base.Base
}

func (s *Sub) PackHeader(header http.Header) {
}

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

func (s *Sub) GetResp() error {
	switch s.Base.Opt {
	case base.OPTPay:
		resp := s.Base.Resp.(*PayResp)
		if resp.Code != http.StatusOK || resp.PayURL == "" {
			return responseError("megipay pay fail", resp.Msg)
		}
		s.Base.PayResp = &common.PayResp{
			APIOrderID: resp.DisOrderNo,
			URL:        resp.PayURL,
			Channel:    int(s.Base.Channel),
			Message:    resp.Msg,
		}
	case base.OPTWithdraw:
		resp := s.Base.Resp.(*WithdrawResp)
		if resp.Code != http.StatusOK || isFinalFailStatus(resp.Status) {
			return responseError("megipay withdraw fail", resp.Msg)
		}
		s.Base.WithdrawResp = &common.WithdrawResp{
			APIOrderID: resp.DisOrderNo,
			Channel:    int(s.Base.Channel),
			Message:    resp.Msg,
		}
	case base.OPTQueryPay:
		resp := s.Base.Resp.(*QueryResp)
		log.Debug("megipay QueryPayResp:%+v", resp)
		s.Base.QueryPayResp.Msg = resp.Msg
		s.Base.QueryPayResp.OrderID = resp.OrderNo
		s.Base.QueryPayResp.APIOrderID = resp.DisOrderNo
		if s.Base.Status == 0 && resp.Code == http.StatusOK {
			switch resp.Status {
			case 2:
				s.Base.QueryPayResp.Status = base.QueryStatusOrderSuccess
			case 3:
				s.Base.QueryPayResp.Status = base.QueryStatusFail
			}
		}
		if s.Base.QueryPayResp.Status == 0 {
			s.Base.QueryPayResp.Status = base.QueryStatusUnknown
		}
	case base.OPTQueryWithdraw:
		resp := s.Base.Resp.(*QueryResp)
		log.Debug("megipay QueryWithdrawResp:%+v", resp)
		s.Base.QueryWithdrawResp.Msg = resp.Msg
		s.Base.QueryWithdrawResp.OrderID = resp.OrderNo
		s.Base.QueryWithdrawResp.APIOrderID = resp.DisOrderNo
		if s.Base.Status == 0 && resp.Code == http.StatusOK {
			switch resp.Status {
			case 2:
				s.Base.QueryWithdrawResp.Status = base.QueryStatusOrderSuccess
			case 10:
				s.Base.QueryWithdrawResp.Status = base.QueryStatusSuccess
			case 3, 7, 9:
				s.Base.QueryWithdrawResp.Status = base.QueryStatusFail
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

func (s *Sub) PackPayReq() interface{} {
	r := s.Base.PayReq
	attach := marshalAttach(&PayAttach{Name: r.Name})
	send := &PayReq{
		TradeNo:      mid,
		AppID:        appID,
		PayCode:      payCodePayIn,
		PayMethod:    payMethod,
		Price:        amountToPaise(r.Amount),
		OrderNo:      r.OrderID,
		SuccessURL:   values.GetFrontCallback(),
		FailURL:      values.GetFrontCallbackFail(),
		PayNoticeURL: s.Base.GetPayCallbackURL(),
		UserID:       r.UID,
		UserIP:       r.IP,
		Attach:       attach,
	}
	send.Sign = signMegipay(send)
	return send
}

func (s *Sub) PackWithdrawReq() interface{} {
	r := s.Base.WithdrawReq
	send := &WithdrawReq{
		TradeNo:      mid,
		OrderNo:      r.OrderID,
		AppID:        appID,
		PayCode:      payCodePayOut,
		Price:        amountToPaise(r.Amount),
		AccountNo:    r.Account,
		AccountName:  r.Name,
		PayNoticeURL: s.Base.GetWithdrawCallbackURL(),
		UserIP:       r.IP,
	}
	attach := &WithdrawAttach{
		Email:    r.Email,
		Phone:    r.Mobile,
		BankName: r.BankName,
	}
	switch r.PayType {
	case common.WithdrawTypeBank:
		send.AccountType = bankAccountType
		send.BankCode = bankCode
		send.IdentifyType = "IFSC"
		send.IdentifyNum = r.IFSC
		attach.PayType = bankPayType
	case common.WithdrawTypeUPI:
		send.AccountType = upiAccountType
		send.BankCode = upiBankCode
		attach.PayType = upiPayType
	default:
		return nil
	}
	send.Attach = marshalAttach(attach)
	send.Sign = signMegipay(send)
	return send
}

func (s *Sub) PackQueryWithdrawReq() interface{} {
	r := s.Base.QueryWithdrawReq
	send := &QueryReq{
		OrderType: payOutOrderType,
		TradeNo:   mid,
		AppID:     appID,
		OrderNo:   r.OrderID,
	}
	send.Sign = signMegipay(send)
	return send
}

func (s *Sub) PackQueryPayReq() interface{} {
	r := s.Base.QueryPayReq
	send := &QueryReq{
		OrderType: payInOrderType,
		TradeNo:   mid,
		AppID:     appID,
		OrderNo:   r.OrderID,
	}
	send.Sign = signMegipay(send)
	return send
}

func (s *Sub) CheckSign(str string) bool {
	checkSign := ""
	switch s.Base.Opt {
	case base.OPTPayCB:
		req := s.Base.CallbackReq.(*PayCallbackReq)
		log.Debug("megipay pay callback:%+v", *req)
		switch req.Status {
		case 2:
			s.Base.CallbackResp.Success = true
		case 3:
			s.Base.CallbackResp.Success = false
		default:
			s.Base.CallbackStatus = base.CallbackStatusPass
			return false
		}
		checkSign = req.Sign
		s.Base.CallbackResp.OrderID = req.OrderNo
		s.Base.CallbackResp.APIOrderID = req.DisOrderNo
	case base.OPTWithdrawCB:
		req := s.Base.CallbackReq.(*WithdrawCallbackReq)
		log.Debug("megipay withdraw callback:%+v", *req)
		switch req.Status {
		case 2:
			s.Base.CallbackResp.Success = true
		case 3, 7, 9:
			s.Base.CallbackResp.Success = false
		default:
			s.Base.CallbackStatus = base.CallbackStatusPass
			return false
		}
		checkSign = req.Sign
		s.Base.CallbackResp.OrderID = req.OrderNo
		s.Base.CallbackResp.APIOrderID = req.DisOrderNo
		s.Base.CallbackResp.FailMessage = req.Remark
		s.Base.CallbackResp.Utr = extractUTR(req.PayInfo)
	default:
		return false
	}
	mySign, err := signMegipayJSON([]byte(str))
	if err != nil {
		log.Error("megipay check sign parse fail:%v", err)
		return false
	}
	pass := mySign == checkSign
	if !pass {
		log.Error("megipay check sign fail:mySign:%s,checkSign:%s", mySign, checkSign)
	}
	return pass
}

const (
	payInOrderType  = "pay_in"
	payOutOrderType = "pay_out"
)

func signMegipay(send interface{}) string {
	body, err := json.Marshal(send)
	if err != nil {
		log.Error("megipay marshal sign body fail:%v", err)
		return ""
	}
	sign, err := signMegipayJSON(body)
	if err != nil {
		log.Error("megipay sign body fail:%v", err)
		return ""
	}
	return sign
}

func signMegipayJSON(body []byte) (string, error) {
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.UseNumber()
	params := map[string]interface{}{}
	if err := dec.Decode(&params); err != nil {
		return "", err
	}
	keys := make([]string, 0, len(params))
	for k, v := range params {
		if k == "sign" || isEmptySignValue(v) {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var builder strings.Builder
	for i, k := range keys {
		if i > 0 {
			builder.WriteString("&")
		}
		builder.WriteString(k)
		builder.WriteString("=")
		builder.WriteString(formatSignValue(params[k]))
	}
	builder.WriteString("&key=")
	builder.WriteString(merchantSignKey)
	signStr := builder.String()
	log.Debug("megipay signStr:%s", signStr)
	return util.CalculateMD5(signStr), nil
}

func isEmptySignValue(v interface{}) bool {
	if v == nil {
		return true
	}
	if s, ok := v.(string); ok {
		return s == ""
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Slice, reflect.Map:
		return rv.Len() == 0
	}
	return false
}

func formatSignValue(v interface{}) string {
	switch val := v.(type) {
	case json.Number:
		return val.String()
	case string:
		return val
	case bool:
		if val {
			return "true"
		}
		return "false"
	case map[string]interface{}, []interface{}:
		data, _ := json.Marshal(val)
		return string(data)
	default:
		return fmt.Sprintf("%v", val)
	}
}

func amountToPaise(amount float64) int64 {
	return int64(math.Round(amount * 100))
}

func marshalAttach(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		log.Error("megipay marshal attach fail:%v", err)
		return "{}"
	}
	return string(data)
}

func isFinalFailStatus(status int) bool {
	return status == 3 || status == 7 || status == 9
}

func responseError(defaultMsg, msg string) error {
	if msg != "" {
		return errors.New(msg)
	}
	return errors.New(defaultMsg)
}

func extractUTR(payInfo string) string {
	if payInfo == "" {
		return ""
	}
	data := map[string]interface{}{}
	if err := json.Unmarshal([]byte(payInfo), &data); err != nil {
		return ""
	}
	for _, key := range []string{"utr2", "utr", "UTR", "bank_ref_no", "bankRefNo"} {
		if v, ok := data[key]; ok {
			return fmt.Sprintf("%v", v)
		}
	}
	return ""
}
